package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"
)

// The public HTTPS origin and backend HTTP listener deliberately differ. A
// reverse proxy may terminate TLS but must not widen the backend allowlist.
func TestSameOriginProxyDeployment(t *testing.T) {
	cases := []struct {
		name           string
		allowedOrigins string
		origin         string
		wantStatus     int
	}{
		{
			name:       "unconfigured_proxy_reproduces_review_finding",
			origin:     "https://calculator.example",
			wantStatus: http.StatusForbidden,
		},
		{
			name:           "configured_public_origin",
			allowedOrigins: "https://calculator.example",
			origin:         "https://calculator.example",
			wantStatus:     http.StatusOK,
		},
		{
			name:           "comma_separated_origins_trim_whitespace",
			allowedOrigins: " https://other.example, https://calculator.example ",
			origin:         "https://calculator.example",
			wantStatus:     http.StatusOK,
		},
		{
			name:           "forwarded_headers_do_not_authorize_attacker_origin",
			allowedOrigins: "https://calculator.example",
			origin:         "https://attacker.example",
			wantStatus:     http.StatusForbidden,
		},
		{
			name:           "wildcard_configuration_remains_denied",
			allowedOrigins: "*",
			origin:         "https://calculator.example",
			wantStatus:     http.StatusForbidden,
		},
		{
			name:       "originless_clients_still_work",
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ALLOWED_ORIGINS", tt.allowedOrigins)
			backend := httptest.NewServer(newServer().Handler)
			t.Cleanup(backend.Close)
			target, err := url.Parse(backend.URL)
			if err != nil {
				t.Fatal(err)
			}
			proxyHandler := httputil.NewSingleHostReverseProxy(target)
			proxy := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Header.Set("X-Forwarded-Proto", "https")
				r.Header.Set("X-Forwarded-Host", "calculator.example")
				proxyHandler.ServeHTTP(w, r)
			}))
			t.Cleanup(proxy.Close)
			request, err := http.NewRequest(http.MethodPost, proxy.URL+"/api/v1/calculations", strings.NewReader(`{"operation":"add","operands":[1,2]}`))
			if err != nil {
				t.Fatal(err)
			}
			request.Host = "calculator.example"
			request.Header.Set("Content-Type", "application/json")
			if tt.origin != "" {
				request.Header.Set("Origin", tt.origin)
			}
			response, err := proxy.Client().Do(request)
			if err != nil {
				t.Fatal(err)
			}
			body, readErr := io.ReadAll(response.Body)
			closeErr := response.Body.Close()
			if readErr != nil || closeErr != nil {
				t.Fatalf("read=%v close=%v", readErr, closeErr)
			}
			if response.StatusCode != tt.wantStatus {
				t.Fatalf("proxy status=%d want=%d body=%s", response.StatusCode, tt.wantStatus, body)
			}
			if tt.wantStatus == http.StatusOK && string(body) != "{\"result\":3}\n" {
				t.Fatalf("result=%s; want 3", body)
			}
			if tt.wantStatus == http.StatusForbidden && !strings.Contains(string(body), `"code":"CORS_DENIED"`) {
				t.Fatalf("missing stable denial code: %s", body)
			}
			wantOrigin := tt.origin
			if tt.wantStatus == http.StatusForbidden {
				wantOrigin = ""
			}
			if got := response.Header.Get("Access-Control-Allow-Origin"); got != wantOrigin {
				t.Fatalf("allow origin=%q want=%q", got, wantOrigin)
			}
		})
	}
}
