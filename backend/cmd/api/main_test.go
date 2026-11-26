package main

import (
	"net/http/httptest"
	"testing"
	"time"
)

// Characterize the existing composition without opening the fixed public port.
func TestServerConfiguration(t *testing.T) {
	server := newServer()
	if server.Addr != ":8080" {
		t.Fatalf("address=%q; want :8080", server.Addr)
	}
	cases := []struct {
		name string
		got  time.Duration
		want time.Duration
	}{
		{
			name: "read_headers",
			got:  server.ReadHeaderTimeout,
			want: 5 * time.Second,
		},
		{
			name: "read_body",
			got:  server.ReadTimeout,
			want: 10 * time.Second,
		},
		{
			name: "write_response",
			got:  server.WriteTimeout,
			want: 10 * time.Second,
		},
		{
			name: "idle_connection",
			got:  server.IdleTimeout,
			want: 60 * time.Second,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("timeout=%s; want %s", tt.got, tt.want)
			}
		})
	}
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)
	if rec.Code != 200 || rec.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("default handler: status=%d body=%q", rec.Code, rec.Body.String())
	}
}
