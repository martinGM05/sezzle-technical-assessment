package http_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	transport "github.com/martinGM05/sezzle-technical-assessment/backend/internal/transport/http"
)

func TestHTTPCORS(t *testing.T) {
	const allowed = "https://calculator.example"
	cases := []struct {
		name           string
		origins        []string
		origin         string
		method         string
		path           string
		requestMethod  string
		requestHeaders string
		wantStatus     int
		wantCode       string
		wantOrigin     string
		wantMethods    string
		wantHeaders    string
	}{
		{
			name:       "default_without_origin",
			method:     "POST",
			wantStatus: 200,
		},
		{
			name:       "default_denies_origin",
			origin:     allowed,
			method:     "POST",
			wantStatus: 403,
			wantCode:   "CORS_DENIED",
		},
		{
			name:          "default_denies_preflight",
			origin:        allowed,
			method:        "OPTIONS",
			requestMethod: "POST",
			wantStatus:    403,
			wantCode:      "CORS_DENIED",
		},
		{
			name:       "allowed_origin",
			origins:    []string{allowed},
			origin:     allowed,
			method:     "POST",
			wantStatus: 200,
			wantOrigin: allowed,
		},
		{
			name:       "second_allowed_origin",
			origins:    []string{"http://localhost:5173", allowed},
			origin:     allowed,
			method:     "POST",
			wantStatus: 200,
			wantOrigin: allowed,
		},
		{
			name:       "disallowed_origin",
			origins:    []string{allowed},
			origin:     "https://evil.example",
			method:     "POST",
			wantStatus: 403,
			wantCode:   "CORS_DENIED",
		},
		{
			name:       "suffix_attack",
			origins:    []string{allowed},
			origin:     allowed + ".evil.example",
			method:     "POST",
			wantStatus: 403,
			wantCode:   "CORS_DENIED",
		},
		{
			name:       "opaque_origin_rejected",
			origins:    []string{"null"},
			origin:     "null",
			method:     "POST",
			wantStatus: 403,
			wantCode:   "CORS_DENIED",
		},
		{
			name:       "wildcard_not_supported",
			origins:    []string{"*"},
			origin:     allowed,
			method:     "POST",
			wantStatus: 403,
			wantCode:   "CORS_DENIED",
		},
		{
			name:           "allowed_preflight",
			origins:        []string{allowed},
			origin:         allowed,
			method:         "OPTIONS",
			requestMethod:  "POST",
			requestHeaders: "content-type",
			wantStatus:     204,
			wantOrigin:     allowed,
			wantMethods:    "POST",
			wantHeaders:    "Content-Type",
		},
		{
			name:           "header_case_insensitive",
			origins:        []string{allowed},
			origin:         allowed,
			method:         "OPTIONS",
			requestMethod:  "POST",
			requestHeaders: "CONTENT-TYPE",
			wantStatus:     204,
			wantOrigin:     allowed,
			wantMethods:    "POST",
			wantHeaders:    "Content-Type",
		},
		{
			name:          "preflight_without_headers",
			origins:       []string{allowed},
			origin:        allowed,
			method:        "OPTIONS",
			requestMethod: "POST",
			wantStatus:    204,
			wantOrigin:    allowed,
			wantMethods:   "POST",
		},
		{
			name:          "preflight_denies_method",
			origins:       []string{allowed},
			origin:        allowed,
			method:        "OPTIONS",
			requestMethod: "DELETE",
			wantStatus:    403,
			wantCode:      "CORS_DENIED",
		},
		{
			name:           "preflight_denies_headers",
			origins:        []string{allowed},
			origin:         allowed,
			method:         "OPTIONS",
			requestMethod:  "POST",
			requestHeaders: "Content-Type, Authorization",
			wantStatus:     403,
			wantCode:       "CORS_DENIED",
		},
		{
			name:       "preflight_missing_method",
			origins:    []string{allowed},
			origin:     allowed,
			method:     "OPTIONS",
			wantStatus: 403,
			wantCode:   "CORS_DENIED",
		},
		{
			name:          "preflight_missing_origin",
			origins:       []string{allowed},
			method:        "OPTIONS",
			requestMethod: "POST",
			wantStatus:    403,
			wantCode:      "CORS_DENIED",
		},
		{
			name:          "health_preflight",
			origins:       []string{allowed},
			origin:        allowed,
			path:          "/health",
			method:        "OPTIONS",
			requestMethod: "GET",
			wantStatus:    204,
			wantOrigin:    allowed,
			wantMethods:   "GET, HEAD",
		},
		{
			name:          "health_denies_post_preflight",
			origins:       []string{allowed},
			origin:        allowed,
			path:          "/health",
			method:        "OPTIONS",
			requestMethod: "POST",
			wantStatus:    403,
			wantCode:      "CORS_DENIED",
		},
		{
			name:          "unknown_route_preflight",
			origins:       []string{allowed},
			origin:        allowed,
			path:          "/unknown",
			method:        "OPTIONS",
			requestMethod: "POST",
			wantStatus:    404,
			wantCode:      "NOT_FOUND",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if path == "" {
				path = "/api/v1/calculations"
			}
			req := httptest.NewRequest(tt.method, path, strings.NewReader(`{"operation":"add","operands":[1,2]}`))
			req.Header.Set("Content-Type", "application/json")
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.requestMethod != "" {
				req.Header.Set("Access-Control-Request-Method", tt.requestMethod)
			}
			if tt.requestHeaders != "" {
				req.Header.Set("Access-Control-Request-Headers", tt.requestHeaders)
			}
			rec := httptest.NewRecorder()
			transport.NewHandler(tt.origins...).ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status=%d want=%d body=%s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			assertDocumentedResponse(t, tt.method, path, rec.Code, rec.Header().Get("Content-Type"), rec.Body.String())
			for key, want := range map[string]string{
				"Access-Control-Allow-Origin":      tt.wantOrigin,
				"Access-Control-Allow-Methods":     tt.wantMethods,
				"Access-Control-Allow-Headers":     tt.wantHeaders,
				"Access-Control-Allow-Credentials": "",
			} {
				if got := rec.Header().Get(key); got != want {
					t.Errorf("%s=%q want=%q", key, got, want)
				}
			}
			if tt.wantCode != "" {
				assertProblem(t, rec, tt.wantStatus, tt.wantCode)
			}
			if tt.wantStatus == 204 && rec.Body.Len() != 0 {
				t.Fatalf("preflight body=%q; want empty", rec.Body.String())
			}
			if tt.wantStatus == 200 && rec.Body.String() != "{\"result\":3}\n" {
				t.Fatalf("unexpected result: %s", rec.Body.String())
			}
			if tt.wantStatus != 404 {
				vary := strings.Join(rec.Header().Values("Vary"), ",")
				for _, key := range []string{"Origin"} {
					if !strings.Contains(vary, key) {
						t.Errorf("Vary=%q missing %s", vary, key)
					}
				}
				if tt.method == "OPTIONS" {
					for _, key := range []string{"Access-Control-Request-Method", "Access-Control-Request-Headers"} {
						if !strings.Contains(vary, key) {
							t.Errorf("Vary=%q missing %s", vary, key)
						}
					}
				}
			}
		})
	}
}
