package http_test

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	transport "github.com/martinGM05/sezzle-technical-assessment/backend/internal/transport/http"
)

type httpCase struct {
	name       string
	method     string
	path       string
	body       string
	mediaType  string
	wantStatus int
	wantCode   string
	wantResult float64
	wantAllow  string
}

func runHTTPCase(t *testing.T, handler http.Handler, tt httpCase) {
	t.Helper()
	req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
	if tt.mediaType != "" {
		req.Header.Set("Content-Type", tt.mediaType)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != tt.wantStatus {
		t.Fatalf("%s %s: status=%d want=%d body=%s", tt.method, tt.path, rec.Code, tt.wantStatus, rec.Body.String())
	}
	assertDocumentedResponse(t, tt.method, tt.path, rec.Code, rec.Header().Get("Content-Type"), rec.Body.String())
	if got := rec.Header().Get("Allow"); got != tt.wantAllow {
		t.Errorf("Allow=%q want=%q", got, tt.wantAllow)
	}
	if tt.method == http.MethodHead {
		if rec.Body.Len() != 0 {
			t.Fatalf("HEAD body=%q; want empty", rec.Body.String())
		}
		wantMedia := "application/json"
		if tt.wantCode != "" {
			wantMedia = "application/problem+json"
		}
		if got := rec.Header().Get("Content-Type"); got != wantMedia {
			t.Fatalf("Content-Type=%q want=%q", got, wantMedia)
		}
		return
	}
	if tt.wantCode != "" {
		assertProblem(t, rec, tt.wantStatus, tt.wantCode)
		return
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type=%q; want application/json", got)
	}
	schemaName := "CalculationResult"
	if tt.path == "/health" {
		schemaName = "Health"
	}
	assertResponseSchema(t, rec.Code, rec.Body.String(), schemaName)
	var body map[string]json.RawMessage
	decodeOne(t, rec.Body.String(), &body)
	if len(body) != 1 {
		t.Fatalf("unexpected success keys: %s", rec.Body.String())
	}
	if tt.path == "/health" {
		if string(body["status"]) != `"ok"` {
			t.Fatalf("health=%s; want status ok", rec.Body.String())
		}
		return
	}
	var result float64
	if err := json.Unmarshal(body["result"], &result); err != nil {
		t.Fatal(err)
	}
	if math.IsNaN(result) || math.IsInf(result, 0) || math.Abs(result-tt.wantResult) > 1e-14*math.Abs(tt.wantResult) {
		t.Fatalf("result=%.17g want=%.17g", result, tt.wantResult)
	}
}

func decodeOne(t *testing.T, body string, target any) {
	t.Helper()
	d := json.NewDecoder(strings.NewReader(body))
	if err := d.Decode(target); err != nil {
		t.Fatalf("invalid JSON %q: %v", body, err)
	}
	if err := d.Decode(new(any)); err != io.EOF {
		t.Fatalf("trailing response JSON: %v", err)
	}
}

func assertProblem(t *testing.T, rec *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status=%d want=%d", rec.Code, status)
	}
	if rec.Header().Get("Content-Type") != "application/problem+json" {
		t.Fatalf("problem media type=%q", rec.Header().Get("Content-Type"))
	}
	assertResponseSchema(t, rec.Code, rec.Body.String(), "")
	var body map[string]any
	decodeOne(t, rec.Body.String(), &body)
	title := http.StatusText(status)
	if status == 422 {
		title = "Unprocessable Content"
	}
	if len(body) != 5 || body["type"] != "about:blank" || body["title"] != title || body["status"] != float64(status) || body["code"] != code {
		t.Fatalf("invalid Problem Details: %s; want status %d code %s title %s", rec.Body.String(), status, code, title)
	}
	detail, ok := body["detail"].(string)
	if !ok || strings.TrimSpace(detail) == "" {
		t.Fatalf("missing detail: %v", body)
	}
	for _, secret := range []string{"do-not-leak", "panic-secret", "ReadCloser", "json.Unmarshal", "goroutine"} {
		if strings.Contains(rec.Body.String(), secret) {
			t.Fatalf("response leaks internal detail: %s", rec.Body.String())
		}
	}
}

func TestHTTPSuccess(t *testing.T) {
	cases := []httpCase{
		{
			name:       "add",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "add", "operands": [2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: 5,
		},
		{
			name:       "subtract",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "subtract", "operands": [1.5, 2]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: -0.5,
		},
		{
			name:       "multiply",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "multiply", "operands": [-2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: -6,
		},
		{
			name:       "divide",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "divide", "operands": [1, 3]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: 1.0 / 3,
		},
		{
			name:       "power",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": [9, 0.5]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: 3,
		},
		{
			name:       "sqrt",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "sqrt", "operands": [2]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: 1.4142135623730951,
		},
		{
			name:       "percentage",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "percentage", "operands": [150.5]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: 1.505,
		},
		{
			name:       "power_zero_to_zero",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": [0, 0]}`,
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: 1,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPJSONValidation(t *testing.T) {
	cases := []httpCase{
		{
			name:       "empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "whitespace",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       "  \n",
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "malformed",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "unknown_property",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2],"secret":"do-not-leak"}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "trailing_object",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]} {}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "trailing_null",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]} null`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "trailing_garbage",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]} x`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "root_array",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `[]`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "root_string",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `"hi"`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "root_number",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `3`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "root_boolean",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `true`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operation_number",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":1,"operands":[1,2]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operation_array",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":[],"operands":[1,2]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operands_string",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":"1,2"}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operands_object",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":{}}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operand_string",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":["1",2]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operand_boolean",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[true,2]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operand_object",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[{},2]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "operand_array",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[[],2]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "nan_literal",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"sqrt","operands":[NaN]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "infinity_literal",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"sqrt","operands":[Infinity]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "number_outside_float64",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"sqrt","operands":[1e400]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPRequiredFields(t *testing.T) {
	cases := []httpCase{
		{
			name:       "missing_both",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "missing_operation",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operands":[1,2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "missing_operands",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add"}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "null_operation",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":null,"operands":[1,2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "null_operands",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":null}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "null_both",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":null,"operands":null}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "null_first_operand",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[null,2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "null_second_operand",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,null]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "null_unary_operand",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"sqrt","operands":[null]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPArity(t *testing.T) {
	cases := []httpCase{
		{
			name:       "add_empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "add", "operands": []}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "add_wrong_arity",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "add", "operands": [1]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "add_excessive",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "add", "operands": [1, 2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "subtract_empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "subtract", "operands": []}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "subtract_wrong_arity",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "subtract", "operands": [1]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "subtract_excessive",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "subtract", "operands": [1, 2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "multiply_empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "multiply", "operands": []}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "multiply_wrong_arity",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "multiply", "operands": [1]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "multiply_excessive",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "multiply", "operands": [1, 2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "divide_empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "divide", "operands": []}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "divide_wrong_arity",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "divide", "operands": [1]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "divide_excessive",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "divide", "operands": [1, 2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "power_empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": []}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "power_wrong_arity",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": [1]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "power_excessive",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": [1, 2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "sqrt_empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "sqrt", "operands": []}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "sqrt_wrong_arity",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "sqrt", "operands": [1, 2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "sqrt_excessive",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "sqrt", "operands": [1, 2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "percentage_empty",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "percentage", "operands": []}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "percentage_wrong_arity",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "percentage", "operands": [1, 2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "percentage_excessive",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "percentage", "operands": [1, 2, 3]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPDomainErrors(t *testing.T) {
	cases := []httpCase{
		{
			name:       "unsupported",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "unknown", "operands": [1, 2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERATION",
		},
		{
			name:       "empty_operation",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "", "operands": [1, 2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERATION",
		},
		{
			name:       "uppercase",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "ADD", "operands": [1, 2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERATION",
		},
		{
			name:       "padded",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": " add ", "operands": [1, 2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERATION",
		},
		{
			name:       "divide_zero",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "divide", "operands": [1, 0]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "DIVISION_BY_ZERO",
		},
		{
			name:       "sqrt_negative",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "sqrt", "operands": [-1]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "power_nonreal",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": [-1, 0.5]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_RESULT",
		},
		{
			name:       "power_zero_negative",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": [0, -1]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_RESULT",
		},
		{
			name:       "addition_overflow",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "add", "operands": [1.7976931348623157e+308, 1.7976931348623157e+308]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_RESULT",
		},
		{
			name:       "power_overflow",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation": "power", "operands": [1e+308, 2]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "INVALID_RESULT",
		},
		{
			name:       "divide_negative_zero",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"divide","operands":[1,-0.0]}`,
			mediaType:  "application/json",
			wantStatus: 422,
			wantCode:   "DIVISION_BY_ZERO",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPStrictStructure(t *testing.T) {
	cases := []httpCase{
		{
			name:       "root_null",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `null`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:       "case_sensitive_properties",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"Operation":"add","operands":[1,2]}`,
			mediaType:  "application/json",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPMediaType(t *testing.T) {
	cases := []httpCase{
		{
			name:       "missing",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]}`,
			mediaType:  "",
			wantStatus: 415,
			wantCode:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:       "text",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]}`,
			mediaType:  "text/plain",
			wantStatus: 415,
			wantCode:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:       "vendor",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]}`,
			mediaType:  "application/vnd.calculator+json",
			wantStatus: 415,
			wantCode:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:       "malformed",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]}`,
			mediaType:  "application/json; charset",
			wantStatus: 415,
			wantCode:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:       "xml",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]}`,
			mediaType:  "application/xml",
			wantStatus: 415,
			wantCode:   "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:       "utf8_parameter",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]}`,
			mediaType:  "application/json; charset=utf-8",
			wantStatus: 200,
			wantResult: 3,
		},
		{
			name:       "case_insensitive_type",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       `{"operation":"add","operands":[1,2]}`,
			mediaType:  "Application/JSON",
			wantStatus: 200,
			wantResult: 3,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPHealth(t *testing.T) {
	cases := []httpCase{
		{
			name:       "get",
			method:     "GET",
			path:       "/health",
			body:       ``,
			mediaType:  "",
			wantStatus: 200,
		},
		{
			name:       "head",
			method:     "HEAD",
			path:       "/health",
			body:       ``,
			mediaType:  "",
			wantStatus: 200,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

func TestHTTPRouting(t *testing.T) {
	cases := []httpCase{
		{
			name:       "calculation_GET",
			method:     "GET",
			path:       "/api/v1/calculations",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "POST",
		},
		{
			name:       "calculation_HEAD",
			method:     "HEAD",
			path:       "/api/v1/calculations",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "POST",
		},
		{
			name:       "calculation_PUT",
			method:     "PUT",
			path:       "/api/v1/calculations",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "POST",
		},
		{
			name:       "calculation_PATCH",
			method:     "PATCH",
			path:       "/api/v1/calculations",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "POST",
		},
		{
			name:       "calculation_DELETE",
			method:     "DELETE",
			path:       "/api/v1/calculations",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "POST",
		},
		{
			name:       "calculation_OPTIONS",
			method:     "OPTIONS",
			path:       "/api/v1/calculations",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 403,
			wantCode:   "CORS_DENIED",
		},
		{
			name:       "health_POST",
			method:     "POST",
			path:       "/health",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "GET, HEAD",
		},
		{
			name:       "health_PUT",
			method:     "PUT",
			path:       "/health",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "GET, HEAD",
		},
		{
			name:       "health_DELETE",
			method:     "DELETE",
			path:       "/health",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 405,
			wantCode:   "METHOD_NOT_ALLOWED",
			wantAllow:  "GET, HEAD",
		},
		{
			name:       "unknown__missing",
			method:     "GET",
			path:       "/missing",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 404,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "unknown__api_v1_calculations_",
			method:     "GET",
			path:       "/api/v1/calculations/",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 404,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "unknown__health_",
			method:     "GET",
			path:       "/health/",
			body:       ``,
			mediaType:  "application/json",
			wantStatus: 404,
			wantCode:   "NOT_FOUND",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

// The documented one-MiB bound applies to all bytes, even trailing whitespace.
func TestHTTPBodyLimit(t *testing.T) {
	const limit = 1 << 20
	valid := `{"operation":"add","operands":[1,2]}`
	cases := []httpCase{
		{
			name:       "exact_limit",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       valid + strings.Repeat(" ", limit-len(valid)),
			mediaType:  "application/json",
			wantStatus: 200,
			wantResult: 3,
		},
		{
			name:       "oversize_trailing_whitespace",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       valid + strings.Repeat(" ", limit-len(valid)+1),
			mediaType:  "application/json",
			wantStatus: 413,
			wantCode:   "BODY_TOO_LARGE",
		},
		{
			name:       "oversize_before_first_json",
			method:     "POST",
			path:       "/api/v1/calculations",
			body:       strings.Repeat(" ", limit+1) + valid,
			mediaType:  "application/json",
			wantStatus: 413,
			wantCode:   "BODY_TOO_LARGE",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { runHTTPCase(t, transport.NewHandler(), tt) })
	}
}

type failingBody struct{ panicOnRead bool }

func (b failingBody) Read([]byte) (int, error) {
	if b.panicOnRead {
		panic("panic-secret")
	}
	return 0, io.ErrUnexpectedEOF
}
func (b failingBody) Close() error { return nil }

func TestHTTPReadFailures(t *testing.T) {
	cases := []struct {
		name        string
		panicOnRead bool
		wantStatus  int
		wantCode    string
	}{
		{
			name:       "body_read_error",
			wantStatus: 400,
			wantCode:   "MALFORMED_JSON",
		},
		{
			name:        "unexpected_panic_is_sanitized",
			panicOnRead: true,
			wantStatus:  500,
			wantCode:    "INTERNAL_ERROR",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/calculations", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Body = failingBody{panicOnRead: tt.panicOnRead}
			rec := httptest.NewRecorder()
			transport.NewHandler().ServeHTTP(rec, req)
			assertProblem(t, rec, tt.wantStatus, tt.wantCode)
		})
	}
}
