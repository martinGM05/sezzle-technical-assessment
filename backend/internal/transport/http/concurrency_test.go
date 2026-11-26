package http_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	transport "github.com/martinGM05/sezzle-technical-assessment/backend/internal/transport/http"
)

// One shared server and handler exercise isolation across simultaneous success,
// validation-error, and origin-denial responses. This is not a load benchmark.
func TestHTTPConcurrentRequests(t *testing.T) {
	server := httptest.NewServer(transport.NewHandler("https://calculator.example"))
	t.Cleanup(server.Close)
	client := server.Client()
	client.Timeout = 5 * time.Second
	cases := []struct {
		name   string
		body   string
		origin string
		status int
		want   string
		code   string
	}{
		{
			name:   "add",
			body:   `{"operation":"add","operands":[1,2]}`,
			status: 200,
			want:   `{"result":3}`,
		},
		{
			name:   "subtract",
			body:   `{"operation":"subtract","operands":[1,2]}`,
			status: 200,
			want:   `{"result":-1}`,
		},
		{
			name:   "multiply",
			body:   `{"operation":"multiply","operands":[2,3]}`,
			status: 200,
			want:   `{"result":6}`,
		},
		{
			name:   "divide",
			body:   `{"operation":"divide","operands":[1,4]}`,
			status: 200,
			want:   `{"result":0.25}`,
		},
		{
			name:   "power",
			body:   `{"operation":"power","operands":[2,3]}`,
			status: 200,
			want:   `{"result":8}`,
		},
		{
			name:   "sqrt",
			body:   `{"operation":"sqrt","operands":[9]}`,
			status: 200,
			want:   `{"result":3}`,
		},
		{
			name:   "percentage",
			body:   `{"operation":"percentage","operands":[50]}`,
			status: 200,
			want:   `{"result":0.5}`,
		},
		{
			name:   "domain_error",
			body:   `{"operation":"divide","operands":[1,0]}`,
			status: 422,
			code:   "DIVISION_BY_ZERO",
		},
		{
			name:   "cors_denial",
			body:   `{"operation":"add","operands":[1,2]}`,
			origin: "https://evil.example",
			status: 403,
			code:   "CORS_DENIED",
		},
	}
	for round := 0; round < 4; round++ {
		for _, tt := range cases {
			t.Run(fmt.Sprintf("%s_%d", tt.name, round), func(t *testing.T) {
				t.Parallel()
				req, err := http.NewRequest("POST", server.URL+"/api/v1/calculations", strings.NewReader(tt.body))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				origin := tt.origin
				if origin == "" && round%2 == 0 {
					origin = "https://calculator.example"
				}
				if origin != "" {
					req.Header.Set("Origin", origin)
				}
				res, err := client.Do(req)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := res.Body.Close(); err != nil {
						t.Errorf("close response body: %v", err)
					}
				}()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				if res.StatusCode != tt.status {
					t.Fatalf("status=%d want=%d body=%s", res.StatusCode, tt.status, body)
				}
				wantOrigin := origin
				if tt.status == 403 {
					wantOrigin = ""
				}
				if got := res.Header.Get("Access-Control-Allow-Origin"); got != wantOrigin {
					t.Errorf("origin header leaked: got %q want %q", got, wantOrigin)
				}
				rec := httptest.NewRecorder()
				rec.Code = res.StatusCode
				rec.Header().Set("Content-Type", res.Header.Get("Content-Type"))
				if _, err := rec.Body.Write(body); err != nil {
					t.Fatalf("copy response body: %v", err)
				}
				if tt.code != "" {
					assertProblem(t, rec, tt.status, tt.code)
					return
				}
				if strings.TrimSpace(string(body)) != tt.want {
					t.Fatalf("response=%s want=%s", body, tt.want)
				}
				if res.Header.Get("Content-Type") != "application/json" {
					t.Fatal("incorrect success media type")
				}
				assertResponseSchema(t, 200, string(body), "CalculationResult")
			})
		}
	}
}
