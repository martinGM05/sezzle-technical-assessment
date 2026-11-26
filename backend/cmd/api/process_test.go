package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

// This subprocess exercises the same exit-code entry point as main, while
// keeping the fixed public port out of tests and isolating process signals.
func TestMain(m *testing.M) {
	if os.Getenv("CALCULATOR_TEST_PROCESS") != "1" {
		os.Exit(m.Run())
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	server := newServer()
	server.Addr = os.Getenv("CALCULATOR_TEST_ADDRESS")
	os.Exit(run(server))
}

func TestProcessStartupFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("starts a subprocess")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := listener.Close(); err != nil {
			t.Errorf("close occupied listener: %v", err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^$")
	command.Env = append(os.Environ(), "CALCULATOR_TEST_PROCESS=1", "CALCULATOR_TEST_ADDRESS="+listener.Addr().String())
	output, err := command.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("startup failure did not terminate: %v; output=%s", ctx.Err(), output)
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
		t.Fatalf("startup failure: error=%v; want process exit 1; output=%s", err, output)
	}
}

func TestProcessSignals(t *testing.T) {
	if testing.Short() {
		t.Skip("starts subprocesses and sends signals")
	}
	cases := []struct {
		name   string
		signal os.Signal
	}{
		{
			name:   "SIGINT",
			signal: os.Interrupt,
		},
		{
			name:   "SIGTERM",
			signal: syscall.SIGTERM,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			command := exec.CommandContext(ctx, os.Args[0], "-test.run=^$")
			command.Env = append(os.Environ(), "CALCULATOR_TEST_PROCESS=1", "CALCULATOR_TEST_ADDRESS=127.0.0.1:0")
			stdout, err := command.StdoutPipe()
			if err != nil {
				cancel()
				t.Fatal(err)
			}
			var stderr bytes.Buffer
			command.Stderr = &stderr
			if err := command.Start(); err != nil {
				cancel()
				t.Fatal(err)
			}
			ready := make(chan string, 1)
			scanned := make(chan error, 1)
			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					var entry struct {
						Message string `json:"msg"`
						Address string `json:"address"`
					}
					if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
						scanned <- err
						return
					}
					if entry.Message == "starting server" {
						ready <- entry.Address
					}
				}
				scanned <- scanner.Err()
			}()
			finished := make(chan struct{})
			result := make(chan error, 1)
			go func() {
				// Finish reading the pipe before Wait closes it.
				scanErr := <-scanned
				result <- errors.Join(scanErr, command.Wait())
				close(finished)
			}()
			t.Cleanup(func() {
				cancel()
				<-finished
			})
			var address string
			select {
			case address = <-ready:
			case <-finished:
				t.Fatalf("server exited before ready: %v; stderr=%s", <-result, stderr.String())
			case <-ctx.Done():
				t.Fatal("server never became ready")
			}
			client := &http.Client{Timeout: 2 * time.Second}
			t.Cleanup(client.CloseIdleConnections)
			response, err := client.Get("http://" + address + "/health")
			if err != nil {
				t.Fatal(err)
			}
			body, readErr := io.ReadAll(response.Body)
			closeErr := response.Body.Close()
			if err := errors.Join(readErr, closeErr); err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != http.StatusOK || string(body) != "{\"status\":\"ok\"}\n" {
				t.Fatalf("health: status=%d body=%s", response.StatusCode, body)
			}
			if err := command.Process.Signal(tt.signal); err != nil {
				t.Fatal(err)
			}
			select {
			case <-finished:
				if err := <-result; err != nil {
					t.Fatalf("%s did not exit cleanly: %v; stderr=%s", tt.name, err, stderr.String())
				}
			case <-ctx.Done():
				t.Fatalf("%s did not stop the server within its bound", tt.name)
			}
		})
	}
}
