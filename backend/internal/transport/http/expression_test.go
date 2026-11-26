package http_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	transport "github.com/martinGM05/sezzle-technical-assessment/backend/internal/transport/http"
)

const expressionPath = "/api/v1/expressions"

func expressionBody(t *testing.T, input string) string {
	t.Helper()
	body, err := json.Marshal(map[string]string{"expression": input})
	if err != nil {
		t.Fatalf("encode expression request: %v", err)
	}
	return string(body)
}

func assertExpressionResponse(t *testing.T, rec *httptest.ResponseRecorder, method string, wantStatus int, wantCode string, wantResult float64) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("expression HTTP status=%d; want=%d; body=%.200s", rec.Code, wantStatus, rec.Body.String())
	}
	media := "application/json"
	if wantStatus >= 400 {
		media = "application/problem+json"
	}
	if rec.Header().Get("Content-Type") != media {
		t.Fatalf("Content-Type=%q; want %q", rec.Header().Get("Content-Type"), media)
	}
	assertDocumentedResponse(t, method, expressionPath, rec.Code, media, rec.Body.String())
	if method == http.MethodHead {
		return
	}
	var body map[string]any
	decodeOne(t, rec.Body.String(), &body)
	if wantCode != "" {
		title := http.StatusText(wantStatus)
		if wantStatus == 422 {
			title = "Unprocessable Content"
		}
		detail, ok := body["detail"].(string)
		if len(body) != 5 || body["type"] != "about:blank" || body["title"] != title || body["status"] != float64(wantStatus) || body["code"] != wantCode || !ok || strings.TrimSpace(detail) == "" {
			t.Fatalf("invalid Problem Details: %v; want status %d code %s", body, wantStatus, wantCode)
		}
		return
	}
	result, ok := body["result"].(float64)
	if len(body) != 1 || !ok || math.IsNaN(result) || math.IsInf(result, 0) || math.Abs(result-wantResult) > 1e-14*math.Max(1, math.Abs(wantResult)) {
		t.Fatalf("response=%v; want single finite result %.17g", body, wantResult)
	}
}

func TestHTTPExpressionsSuccess(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "addition",
			input: "2+3",
			want:  5,
		},
		{
			name:  "subtraction_with_sign",
			input: "-5+3",
			want:  -2,
		},
		{
			name:  "multiplication_precedence",
			input: "2+3*4",
			want:  14,
		},
		{
			name:  "parentheses",
			input: "(2+3)*4",
			want:  20,
		},
		{
			name:  "division",
			input: "10/4",
			want:  2.5,
		},
		{
			name:  "right_associative_exponentiation",
			input: "2^3^2",
			want:  512,
		},
		{
			name:  "unary_minus_precedence",
			input: "-2^2",
			want:  -4,
		},
		{
			name:  "negative_base",
			input: "(-2)^2",
			want:  4,
		},
		{
			name:  "negative_exponent",
			input: "2^-3",
			want:  .125,
		},
		{
			name:  "square_root",
			input: "sqrt(9+7)",
			want:  4,
		},
		{
			name:  "percentage",
			input: "200*15%",
			want:  30,
		},
		{
			name:  "maximum_decoded_expression_length",
			input: "1" + strings.Repeat(" ", 4095),
			want:  1,
		},
	}
	handler := transport.NewHandler()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, expressionPath, strings.NewReader(expressionBody(t, tt.input)))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assertExpressionResponse(t, rec, http.MethodPost, 200, "", tt.want)
		})
	}
}

func TestHTTPExpressionsValidation(t *testing.T) {
	cases := []struct {
		name   string
		body   string
		status int
		code   string
	}{
		{
			name:   "empty_body",
			body:   "",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "malformed_json",
			body:   "{",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "unknown_property",
			body:   "{\"expression\":\"1\",\"extra\":2}",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "case_sensitive_field",
			body:   "{\"Expression\":\"1\"}",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "trailing_json",
			body:   "{\"expression\":\"1\"} {}",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "root_null",
			body:   "null",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "root_array",
			body:   "[]",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "root_string",
			body:   "\"1\"",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "number_field",
			body:   "{\"expression\":2}",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "boolean_field",
			body:   "{\"expression\":true}",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "array_field",
			body:   "{\"expression\":[1]}",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "object_field",
			body:   "{\"expression\":{}}",
			status: 400,
			code:   "MALFORMED_JSON",
		},
		{
			name:   "missing_field",
			body:   "{}",
			status: 422,
			code:   "INVALID_EXPRESSION",
		},
		{
			name:   "null_field",
			body:   "{\"expression\":null}",
			status: 422,
			code:   "INVALID_EXPRESSION",
		},
		{
			name:   "empty_expression",
			body:   "{\"expression\": \"\"}",
			status: 422,
			code:   "INVALID_EXPRESSION",
		},
		{
			name:   "blank_expression",
			body:   "{\"expression\": \" \\t\\n\"}",
			status: 422,
			code:   "INVALID_EXPRESSION",
		},
		{
			name:   "invalid_token",
			body:   "{\"expression\": \"2 @ 3\"}",
			status: 422,
			code:   "SYNTAX_ERROR",
		},
		{
			name:   "missing_parenthesis",
			body:   "{\"expression\": \"(2+3\"}",
			status: 422,
			code:   "SYNTAX_ERROR",
		},
		{
			name:   "invalid_operator_sequence",
			body:   "{\"expression\": \"2+*3\"}",
			status: 422,
			code:   "SYNTAX_ERROR",
		},
		{
			name:   "division_by_zero",
			body:   "{\"expression\": \"10/(5-5)\"}",
			status: 422,
			code:   "DIVISION_BY_ZERO",
		},
		{
			name:   "negative_zero_divisor",
			body:   "{\"expression\": \"1/-0\"}",
			status: 422,
			code:   "DIVISION_BY_ZERO",
		},
		{
			name:   "sqrt_domain",
			body:   "{\"expression\": \"sqrt(-1)\"}",
			status: 422,
			code:   "INVALID_OPERANDS",
		},
		{
			name:   "power_domain",
			body:   "{\"expression\": \"(-1)^.5\"}",
			status: 422,
			code:   "INVALID_RESULT",
		},
		{
			name:   "overflow",
			body:   "{\"expression\": \"10^309\"}",
			status: 422,
			code:   "INVALID_RESULT",
		},
		{
			name:   "literal_overflow",
			body:   "{\"expression\": \"1e309\"}",
			status: 422,
			code:   "INVALID_RESULT",
		},
		{
			name:   "nan_spelling",
			body:   "{\"expression\": \"NaN\"}",
			status: 422,
			code:   "SYNTAX_ERROR",
		},
		{
			name:   "excess_length",
			body:   expressionBody(t, strings.Repeat("1", 4097)),
			status: 422,
			code:   "EXPRESSION_LIMIT",
		},
		{
			name:   "excess_nesting",
			body:   expressionBody(t, strings.Repeat("(", 65)+"1"+strings.Repeat(")", 65)),
			status: 422,
			code:   "EXPRESSION_LIMIT",
		},
	}
	handler := transport.NewHandler()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, expressionPath, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assertExpressionResponse(t, rec, http.MethodPost, tt.status, tt.code, 0)
		})
	}
}

func TestHTTPExpressionsConstraints(t *testing.T) {
	validBody := `{"expression":"2+3"}`
	cases := []struct {
		name   string
		method string
		media  string
		body   string
		status int
		code   string
		allow  string
	}{
		{
			name:   "missing_media_type",
			method: http.MethodPost,
			body:   validBody,
			status: 415,
			code:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:   "wrong_media_type",
			method: http.MethodPost,
			media:  "text/plain",
			body:   validBody,
			status: 415,
			code:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:   "malformed_media_type",
			method: http.MethodPost,
			media:  "application/json;=bad",
			body:   validBody,
			status: 415,
			code:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:   "get",
			method: http.MethodGet,
			status: 405,
			code:   "METHOD_NOT_ALLOWED",
			allow:  "POST",
		},
		{
			name:   "head_has_no_body",
			method: http.MethodHead,
			status: 405,
			code:   "METHOD_NOT_ALLOWED",
			allow:  "POST",
		},
		{
			name:   "put",
			method: http.MethodPut,
			status: 405,
			code:   "METHOD_NOT_ALLOWED",
			allow:  "POST",
		},
		{
			name:   "patch",
			method: http.MethodPatch,
			status: 405,
			code:   "METHOD_NOT_ALLOWED",
			allow:  "POST",
		},
		{
			name:   "delete",
			method: http.MethodDelete,
			status: 405,
			code:   "METHOD_NOT_ALLOWED",
			allow:  "POST",
		},
		{
			name:   "options_without_origin",
			method: http.MethodOptions,
			status: 403,
			code:   "CORS_DENIED",
		},
		{
			name:   "exact_body_limit",
			method: http.MethodPost,
			media:  "application/json",
			body:   validBody + strings.Repeat(" ", (1<<20)-len(validBody)),
			status: 200,
		},
		{
			name:   "excess_body_limit",
			method: http.MethodPost,
			media:  "application/json",
			body:   validBody + strings.Repeat(" ", (1<<20)-len(validBody)+1),
			status: 413,
			code:   "BODY_TOO_LARGE",
		},
		{
			name:   "size_precedes_syntax",
			method: http.MethodPost,
			media:  "application/json",
			body:   strings.Repeat("x", (1<<20)+1),
			status: 413,
			code:   "BODY_TOO_LARGE",
		},
		{
			name:   "media_precedes_size",
			method: http.MethodPost,
			media:  "text/plain",
			body:   strings.Repeat("x", (1<<20)+1),
			status: 415,
			code:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:   "duplicate_fields_keep_last_value",
			method: http.MethodPost,
			media:  "application/json",
			body:   `{"expression":"bad","expression":"2+3"}`,
			status: 200,
		},
		{
			name:   "decoded_escaped_string",
			method: http.MethodPost,
			media:  "application/json",
			body:   `{"expression":"\u0032+\u0033"}`,
			status: 200,
		},
	}
	handler := transport.NewHandler()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, expressionPath, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.media)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assertExpressionResponse(t, rec, tt.method, tt.status, tt.code, 5)
			if rec.Header().Get("Allow") != tt.allow {
				t.Fatalf("Allow=%q; want %q", rec.Header().Get("Allow"), tt.allow)
			}
		})
	}
}

func TestHTTPExpressionsReadFailureDoesNotLeak(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, expressionPath, iotest.ErrReader(errors.New("do-not-leak")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	transport.NewHandler().ServeHTTP(rec, req)
	assertExpressionResponse(t, rec, http.MethodPost, 400, "MALFORMED_JSON", 0)
	if strings.Contains(rec.Body.String(), "do-not-leak") {
		t.Fatalf("response leaked internal read failure: %s", rec.Body.String())
	}
}

func TestHTTPExpressionsCORS(t *testing.T) {
	const allowed = "https://calculator.example"
	cases := []struct {
		name             string
		origins          []string
		origin           string
		method           string
		requestedMethod  string
		requestedHeaders string
		status           int
		allowedOrigin    string
	}{
		{
			name:   "restrictive_default",
			origin: allowed,
			method: http.MethodPost,
			status: 403,
		},
		{
			name:          "allowed_origin",
			origins:       []string{allowed},
			origin:        allowed,
			method:        http.MethodPost,
			status:        200,
			allowedOrigin: allowed,
		},
		{
			name:    "disallowed_origin",
			origins: []string{allowed},
			origin:  "https://attacker.example",
			method:  http.MethodPost,
			status:  403,
		},
		{
			name:   "originless_client",
			method: http.MethodPost,
			status: 200,
		},
		{
			name:             "allowed_preflight",
			origins:          []string{allowed},
			origin:           allowed,
			method:           http.MethodOptions,
			requestedMethod:  "POST",
			requestedHeaders: "content-type",
			status:           204,
			allowedOrigin:    allowed,
		},
		{
			name:            "preflight_default_denied",
			origin:          allowed,
			method:          http.MethodOptions,
			requestedMethod: "POST",
			status:          403,
		},
		{
			name:            "preflight_wrong_method",
			origins:         []string{allowed},
			origin:          allowed,
			method:          http.MethodOptions,
			requestedMethod: "DELETE",
			status:          403,
		},
		{
			name:             "preflight_wrong_header",
			origins:          []string{allowed},
			origin:           allowed,
			method:           http.MethodOptions,
			requestedMethod:  "POST",
			requestedHeaders: "Authorization",
			status:           403,
		},
		{
			name:   "cors_precedes_media_validation",
			origin: "https://attacker.example",
			method: http.MethodPost,
			status: 403,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, expressionPath, strings.NewReader(`{"expression":"2+3"}`))
			req.Header.Set("Origin", tt.origin)
			req.Header.Set("Content-Type", "application/json")
			if tt.name == "cors_precedes_media_validation" {
				req.Header.Set("Content-Type", "text/plain")
			}
			req.Header.Set("Access-Control-Request-Method", tt.requestedMethod)
			req.Header.Set("Access-Control-Request-Headers", tt.requestedHeaders)
			rec := httptest.NewRecorder()
			transport.NewHandler(tt.origins...).ServeHTTP(rec, req)
			if tt.status == 204 {
				if rec.Code != 204 || rec.Body.Len() != 0 || rec.Header().Get("Access-Control-Allow-Methods") != "POST" || rec.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
					t.Fatalf("invalid preflight: status=%d headers=%v body=%s", rec.Code, rec.Header(), rec.Body.String())
				}
				assertDocumentedResponse(t, tt.method, expressionPath, rec.Code, "", rec.Body.String())
			} else {
				code := ""
				if tt.status == 403 {
					code = "CORS_DENIED"
				}
				assertExpressionResponse(t, rec, tt.method, tt.status, code, 5)
			}
			if rec.Header().Get("Access-Control-Allow-Origin") != tt.allowedOrigin || rec.Header().Get("Access-Control-Allow-Credentials") != "" {
				t.Fatalf("incorrect CORS permissions: %v", rec.Header())
			}
			if !strings.Contains(strings.Join(rec.Header().Values("Vary"), ","), "Origin") {
				t.Fatalf("missing Origin cache variation: %v", rec.Header())
			}
		})
	}
}

func TestHTTPExpressionsConcurrentRequests(t *testing.T) {
	handler := transport.NewHandler("https://calculator.example")
	for index := 0; index < 32; index++ {
		t.Run(fmt.Sprintf("request_%02d", index), func(t *testing.T) {
			t.Parallel()
			input := fmt.Sprintf("%d*2+sqrt(16)", index)
			want := float64(index*2 + 4)
			req := httptest.NewRequest(http.MethodPost, expressionPath, strings.NewReader(expressionBody(t, input)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Origin", "https://calculator.example")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assertExpressionResponse(t, rec, http.MethodPost, 200, "", want)
		})
	}
}
