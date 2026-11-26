package http_test

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	transport "github.com/martinGM05/sezzle-technical-assessment/backend/internal/transport/http"
)

func TestHTTPErrorPrecedence(t *testing.T) {
	cases := []httpCase{
		{
			name:       "route_before_method_and_media",
			method:     "DELETE",
			path:       "/missing",
			body:       "broken",
			wantStatus: 404,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "method_before_media_and_json",
			method:     "GET",
			path:       "/api/v1/calculations",
			body:       "broken",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "POST",
		},
		{
			name:       "media_before_json",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       "broken",
			wantStatus: 415,
			wantCode:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:       "unknown_property_before_required_fields",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"extra":true}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operand_type_before_missing_operation",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operands":"bad"}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "null_element_before_unsupported_operation",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"unknown","operands":[null,2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			runHTTPCase(t, transport.NewHandler(), tt)
		})
	}
}

type brokenWriter struct {
	*httptest.ResponseRecorder
	writes       int
	panicOnWrite bool
}

func (w *brokenWriter) Write(body []byte) (int, error) {
	w.writes++
	if w.panicOnWrite {
		_, _ = w.ResponseRecorder.Write(body[:1]) // bytes.Buffer cannot fail.
		panic("panic-secret")
	}
	return 0, errors.New("client disconnected")
}

func TestHTTPWriteFailures(t *testing.T) {
	cases := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{
			name:   "health",
			method: "GET",
			path:   "/health",
			status: 200,
		},
		{
			name:   "calculation",
			method: "POST",
			path:   "/api/v1/calculations",
			body:   `{"operation":"add","operands":[1,2]}`, status: 200,
		},
		{
			name:   "problem",
			method: "GET",
			path:   "/missing",
			status: 404,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			previous := slog.Default()
			slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
			t.Cleanup(func() { slog.SetDefault(previous) })
			writer := &brokenWriter{ResponseRecorder: httptest.NewRecorder()}
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			transport.NewHandler().ServeHTTP(writer, req)
			if writer.writes != 1 || writer.Code != tt.status {
				t.Fatalf("write attempts=%d status=%d; want one attempt and status=%d", writer.writes, writer.Code, tt.status)
			}
			if !strings.Contains(logs.String(), "response write failed") {
				t.Fatalf("response write error was not logged: %q", logs.String())
			}
		})
	}
}

func capturePanic(action func()) (value any) {
	defer func() { value = recover() }()
	action()
	return nil
}

func TestHTTPPanicAfterResponseStarted(t *testing.T) {
	writer := &brokenWriter{
		ResponseRecorder: httptest.NewRecorder(),
		panicOnWrite:     true,
	}
	got := capturePanic(func() {
		transport.NewHandler().ServeHTTP(writer, httptest.NewRequest("GET", "/health", nil))
	})
	if got != http.ErrAbortHandler {
		t.Fatalf("panic=%v; want http.ErrAbortHandler to terminate the partial response", got)
	}
	if writer.writes != 1 || writer.Body.Len() != 1 || writer.Code != 200 {
		t.Fatalf("partial response was rewritten: writes=%d status=%d body=%q", writer.writes, writer.Code, writer.Body.String())
	}
}

type abortedBody struct{}

func (abortedBody) Read([]byte) (int, error) { panic(http.ErrAbortHandler) }
func (abortedBody) Close() error             { return nil }

func TestHTTPPreservesAbortHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/calculations", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = abortedBody{}
	rec := httptest.NewRecorder()
	got := capturePanic(func() { transport.NewHandler().ServeHTTP(rec, req) })
	if got != http.ErrAbortHandler || rec.Body.Len() != 0 {
		t.Fatalf("abort panic=%v body=%q; want unchanged sentinel and no response", got, rec.Body.String())
	}
}
