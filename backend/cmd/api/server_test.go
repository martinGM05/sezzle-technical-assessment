package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServeShutdown(t *testing.T) {
	cases := []struct {
		name          string
		timeout       time.Duration
		finishRequest bool
	}{
		{
			name:          "drains_active_request",
			timeout:       2 * time.Second,
			finishRequest: true,
		},
		{
			name:          "deadline_closes_active_connection",
			timeout:       50 * time.Millisecond,
			finishRequest: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			started := make(chan struct{})
			release := make(chan struct{})
			shutdownStarted := make(chan struct{})
			server := &http.Server{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					close(started)
					select {
					case <-release:
						if _, err := io.WriteString(w, "completed"); err != nil {
							t.Errorf("write drained response: %v", err)
						}
					case <-r.Context().Done():
					}
				}),
			}
			server.RegisterOnShutdown(func() { close(shutdownStarted) })
			ctx, cancel := context.WithCancel(context.Background())
			result := make(chan error, 1)
			stopped := make(chan struct{})
			go func() {
				defer close(stopped)
				result <- serve(ctx, server, listener, tt.timeout)
			}()
			t.Cleanup(func() {
				cancel()
				if err := server.Close(); err != nil {
					t.Errorf("close server: %v", err)
				}
				select {
				case <-stopped:
				case <-time.After(3 * time.Second):
					t.Error("server goroutine did not exit")
				}
			})
			client := &http.Client{Timeout: 3 * time.Second}
			t.Cleanup(client.CloseIdleConnections)
			type response struct {
				body string
				err  error
			}
			responseDone := make(chan response, 1)
			go func() {
				res, err := client.Get("http://" + listener.Addr().String())
				if err != nil {
					responseDone <- response{err: err}
					return
				}
				body, readErr := io.ReadAll(res.Body)
				closeErr := res.Body.Close()
				responseDone <- response{body: string(body), err: errors.Join(readErr, closeErr)}
			}()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("request did not reach the handler")
			}
			cancel()
			select {
			case <-shutdownStarted:
			case <-time.After(time.Second):
				t.Fatal("cancellation did not start graceful shutdown")
			}
			if tt.finishRequest {
				select {
				case err := <-result:
					t.Fatalf("server exited before active request completed: %v", err)
				default:
				}
				close(release)
			}
			select {
			case err := <-result:
				if tt.finishRequest && err != nil {
					t.Fatalf("graceful shutdown: %v", err)
				}
				if !tt.finishRequest && !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("shutdown error=%v; want deadline exceeded", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("shutdown exceeded its bound")
			}
			select {
			case response := <-responseDone:
				if tt.finishRequest && (response.err != nil || response.body != "completed") {
					t.Fatalf("drained response=%q error=%v", response.body, response.err)
				}
				if !tt.finishRequest && response.err == nil {
					t.Fatalf("expected forced connection close, got %q", response.body)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("client connection remained open after shutdown")
			}
		})
	}
}

func TestServeTerminalErrors(t *testing.T) {
	cases := []struct {
		name        string
		closeServer bool
		wantError   error
	}{
		{
			name:        "expected_server_closed",
			closeServer: true,
		},
		{
			name:      "unexpected_listener_failure",
			wantError: net.ErrClosed,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			server := &http.Server{}
			if tt.closeServer {
				if err := server.Close(); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := listener.Close(); err != nil {
					t.Fatal(err)
				}
			}
			err = serve(context.Background(), server, listener, time.Second)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("serve error=%v; want %v", err, tt.wantError)
			}
		})
	}
}

func TestServeCanceledBeforeStart(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := serve(ctx, &http.Server{}, listener, time.Second); err != nil {
		t.Fatalf("already-canceled lifecycle: %v", err)
	}
	connection, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
	if err == nil {
		if closeErr := connection.Close(); closeErr != nil {
			t.Errorf("close unexpected connection: %v", closeErr)
		}
		t.Fatal("listener still accepts connections after shutdown")
	}
}
