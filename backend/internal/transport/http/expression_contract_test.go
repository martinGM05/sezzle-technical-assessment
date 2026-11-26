package http_test

import (
	"strings"
	"testing"
)

func TestHTTPOpenAPIExpressionRequest(t *testing.T) {
	doc := loadContract(t)
	schema := contractObject(t, doc, "paths", "/api/v1/expressions", "post", "requestBody", "content", "application/json", "schema")
	cases := []struct {
		name  string
		value any
		valid bool
	}{
		{
			name:  "valid_expression",
			value: map[string]any{"expression": "2+3*4"},
			valid: true,
		},
		{
			name:  "maximum_length",
			value: map[string]any{"expression": "1" + strings.Repeat(" ", 4095)},
			valid: true,
		},
		{
			name:  "excess_length",
			value: map[string]any{"expression": strings.Repeat("1", 4097)},
			valid: false,
		},
		{
			name:  "missing",
			value: map[string]any{},
			valid: false,
		},
		{
			name:  "null",
			value: map[string]any{"expression": nil},
			valid: false,
		},
		{
			name:  "empty",
			value: map[string]any{"expression": ""},
			valid: false,
		},
		{
			name:  "wrong_type",
			value: map[string]any{"expression": float64(2)},
			valid: false,
		},
		{
			name:  "unknown_property",
			value: map[string]any{"expression": "1", "extra": true},
			valid: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := checkSchema(doc, schema, tt.value)
			if (err == nil) != tt.valid {
				t.Fatalf("ExpressionRequest validation error=%v; want valid=%t for %+v", err, tt.valid, tt.value)
			}
		})
	}
}
