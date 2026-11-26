package http_test

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// The contract is JSON-form YAML 1.2. This deliberately checks only the JSON
// Schema keywords used here, not arbitrary OpenAPI documents; no dependency is
// introduced. Unknown validation keywords fail rather than silently passing.
func loadContract(t *testing.T) map[string]any {
	t.Helper()
	data, err := os.ReadFile("../../../../docs/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("OpenAPI must be valid JSON-form YAML: %v", err)
	}
	return doc
}

// contractObject reports a useful path instead of panicking on contract drift.
func contractObject(t *testing.T, object map[string]any, keys ...string) map[string]any {
	t.Helper()
	for index, key := range keys {
		next, ok := object[key].(map[string]any)
		if !ok {
			t.Fatalf("OpenAPI %s must be an object; got %T", strings.Join(keys[:index+1], "/"), object[key])
		}
		object = next
	}
	return object
}

func checkSchema(doc map[string]any, schema map[string]any, value any) error {
	// Validate every branch before matching data: an unsupported alternative
	// must not be mistaken for an ordinary oneOf mismatch.
	if err := validateSchemaSubset(doc, schema, 0); err != nil {
		return err
	}
	return checkSchemaValue(doc, schema, value)
}

func validateSchemaSubset(doc map[string]any, schema map[string]any, depth int) error {
	if depth > 32 {
		return fmt.Errorf("schema reference cycle or nesting exceeds supported depth")
	}
	nested := func(value any) error {
		child, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("expected schema object, got %T", value)
		}
		return validateSchemaSubset(doc, child, depth+1)
	}
	for keyword, value := range schema {
		valid := false
		switch keyword {
		case "$ref":
			reference, ok := value.(string)
			if !ok || len(schema) != 1 {
				return fmt.Errorf("only standalone string $ref is supported")
			}
			target, err := resolveReference(doc, reference)
			if err != nil {
				return err
			}
			return nested(target)
		case "description":
			_, valid = value.(string)
		case "type":
			switch value {
			case "object", "array", "string", "number", "integer":
				valid = true
			}
		case "const":
			valid = true
		case "enum", "required", "oneOf":
			values, ok := value.([]any)
			valid = ok && len(values) > 0
			for _, item := range values {
				if keyword == "oneOf" {
					if err := nested(item); err != nil {
						return err
					}
				}
				if keyword == "required" {
					if _, ok := item.(string); !ok {
						return fmt.Errorf("required entries must be strings")
					}
				}
			}
		case "properties":
			properties, ok := value.(map[string]any)
			valid = ok
			for _, child := range properties {
				if err := nested(child); err != nil {
					return err
				}
			}
		case "items":
			if err := nested(value); err != nil {
				return err
			}
			valid = true
		case "additionalProperties":
			_, valid = value.(bool)
		case "minimum", "maximum", "minItems", "maxItems", "minLength":
			number, ok := value.(float64)
			valid = ok && !math.IsNaN(number) && !math.IsInf(number, 0)
			if keyword != "minimum" && keyword != "maximum" {
				valid = valid && number >= 0 && math.Trunc(number) == number
			}
		}
		if !valid {
			return fmt.Errorf("unsupported or malformed schema keyword %s", keyword)
		}
	}
	return nil
}

func checkSchemaValue(doc map[string]any, schema map[string]any, value any) error {
	if reference, ok := schema["$ref"].(string); ok {
		target, err := resolveReference(doc, reference)
		if err != nil {
			return err
		}
		return checkSchemaValue(doc, target, value)
	}
	for keyword := range schema {
		switch keyword {
		case "description", "type", "const", "enum", "oneOf", "additionalProperties", "required", "properties", "minItems", "maxItems", "items", "minimum", "maximum", "minLength":
		default:
			return fmt.Errorf("unsupported schema keyword %s", keyword)
		}
	}
	if choices, ok := schema["oneOf"].([]any); ok {
		matches := 0
		for _, choice := range choices {
			if checkSchemaValue(doc, choice.(map[string]any), value) == nil {
				matches++
			}
		}
		if matches != 1 {
			return fmt.Errorf("oneOf matched %d alternatives", matches)
		}
	}
	if expected, ok := schema["const"]; ok && !reflect.DeepEqual(value, expected) {
		return fmt.Errorf("value %v differs from const %v", value, expected)
	}
	if values, ok := schema["enum"].([]any); ok {
		found := false
		for _, candidate := range values {
			found = found || reflect.DeepEqual(value, candidate)
		}
		if !found {
			return fmt.Errorf("value %v not in enum", value)
		}
	}
	switch schema["type"] {
	case "object":
		object, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
		properties, _ := schema["properties"].(map[string]any)
		required, _ := schema["required"].([]any)
		for _, key := range required {
			if _, ok := object[key.(string)]; !ok {
				return fmt.Errorf("missing %s", key)
			}
		}
		for key, item := range object {
			property, ok := properties[key]
			if !ok {
				if schema["additionalProperties"] == false {
					return fmt.Errorf("unexpected property %s", key)
				}
				continue
			}
			if err := checkSchemaValue(doc, property.(map[string]any), item); err != nil {
				return fmt.Errorf("%s: %w", key, err)
			}
		}
	case "array":
		array, ok := value.([]any)
		if !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
		if min, ok := schema["minItems"].(float64); ok && len(array) < int(min) {
			return fmt.Errorf("too few items")
		}
		if max, ok := schema["maxItems"].(float64); ok && len(array) > int(max) {
			return fmt.Errorf("too many items")
		}
		for _, item := range array {
			if err := checkSchemaValue(doc, schema["items"].(map[string]any), item); err != nil {
				return err
			}
		}
	case "string":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
		if min, ok := schema["minLength"].(float64); ok && len(str) < int(min) {
			return fmt.Errorf("string too short")
		}
	case "number", "integer":
		number, ok := value.(float64)
		if !ok || math.IsNaN(number) || math.IsInf(number, 0) {
			return fmt.Errorf("expected finite number, got %v", value)
		}
		if schema["type"] == "integer" && math.Trunc(number) != number {
			return fmt.Errorf("expected integer")
		}
		if min, ok := schema["minimum"].(float64); ok && number < min {
			return fmt.Errorf("below minimum")
		}
		if max, ok := schema["maximum"].(float64); ok && number > max {
			return fmt.Errorf("above maximum")
		}
	}
	return nil
}

func resolveReference(doc map[string]any, reference string) (map[string]any, error) {
	if !strings.HasPrefix(reference, "#/") {
		return nil, fmt.Errorf("unsupported reference %s", reference)
	}
	var current any = doc
	for _, segment := range strings.Split(strings.TrimPrefix(reference, "#/"), "/") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid reference %s", reference)
		}
		current = object[segment]
	}
	target, ok := current.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing reference %s", reference)
	}
	return target, nil
}

func assertResponseSchema(t *testing.T, status int, body string, successSchema string) {
	t.Helper()
	doc := loadContract(t)
	name := successSchema
	if status >= 400 {
		name = "Problem" + strconv.Itoa(status)
	}
	schema, err := resolveReference(doc, "#/components/schemas/"+name)
	if err != nil {
		t.Fatal(err)
	}
	var value any
	decodeOne(t, body, &value)
	if err := checkSchema(doc, schema, value); err != nil {
		t.Fatalf("response violates OpenAPI %s: %v; body=%s", name, err, body)
	}
}

func assertRouteContract(t *testing.T, method, path string, status int) {
	t.Helper()
	doc := loadContract(t)
	paths := contractObject(t, doc, "paths")
	pathItem, found := paths[path].(map[string]any)
	if !found {
		if status != 404 {
			t.Fatalf("undocumented route %s status %d", path, status)
		}
		return
	}
	operation, found := pathItem[strings.ToLower(method)].(map[string]any)
	if !found {
		if status != 405 {
			t.Fatalf("undocumented method %s status %d", method, status)
		}
		return
	}
	responses := contractObject(t, operation, "responses")
	if _, ok := responses[strconv.Itoa(status)]; !ok {
		t.Fatalf("%s %s status %d absent from OpenAPI", method, path, status)
	}
}

func TestHTTPOpenAPIRequestContract(t *testing.T) {
	doc := loadContract(t)
	if doc["openapi"] != "3.1.0" {
		t.Fatal("expected OpenAPI 3.1.0")
	}
	paths := contractObject(t, doc, "paths")
	if len(paths) != 2 {
		t.Fatalf("unexpected endpoint count %d", len(paths))
	}
	post := contractObject(t, paths, "/api/v1/calculations", "post")
	content := contractObject(t, post, "requestBody", "content")
	schema := contractObject(t, content, "application/json", "schema")
	arities := []struct {
		operation string
		arity     int
	}{
		{
			operation: "add",
			arity:     2,
		},
		{
			operation: "subtract",
			arity:     2,
		},
		{
			operation: "multiply",
			arity:     2,
		},
		{
			operation: "divide",
			arity:     2,
		},
		{
			operation: "power",
			arity:     2,
		},
		{
			operation: "sqrt",
			arity:     1,
		},
		{
			operation: "percentage",
			arity:     1,
		},
	}
	reference, ok := schema["$ref"].(string)
	if !ok {
		t.Fatal("request schema must reference its component")
	}
	resolved, err := resolveReference(doc, reference)
	if err != nil {
		t.Fatal(err)
	}
	choices, ok := resolved["oneOf"].([]any)
	if !ok || len(choices) != len(arities) {
		t.Fatal("operation count drift")
	}
	for _, tt := range arities {
		for count := 0; count <= 3; count++ {
			t.Run(fmt.Sprintf("%s_arity_%d", tt.operation, count), func(t *testing.T) {
				operands := make([]any, count)
				for i := range operands {
					operands[i] = float64(1)
				}
				value := map[string]any{"operation": tt.operation, "operands": operands}
				err := checkSchema(doc, schema, value)
				if (err == nil) != (count == tt.arity) {
					t.Fatalf("schema arity=%d want=%d: %v", count, tt.arity, err)
				}
			})
		}
	}
}

func TestHTTPOpenAPIRejectsInvalidRequests(t *testing.T) {
	doc := loadContract(t)
	schema, err := resolveReference(doc, "#/components/schemas/CalculationRequest")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name string
		body string
	}{
		{
			name: "missing_both",
			body: `{}`,
		},
		{
			name: "missing_operation",
			body: `{"operands":[1,2]}`,
		},
		{
			name: "missing_operands",
			body: `{"operation":"add"}`,
		},
		{
			name: "null_operation",
			body: `{"operation":null,"operands":[1,2]}`,
		},
		{
			name: "null_operands",
			body: `{"operation":"add","operands":null}`,
		},
		{
			name: "null_element",
			body: `{"operation":"sqrt","operands":[null]}`,
		},
		{
			name: "unknown_field",
			body: `{"operation":"add","operands":[1,2],"extra":1}`,
		},
		{
			name: "unknown_operation",
			body: `{"operation":"unknown","operands":[1,2]}`,
		},
		{
			name: "uppercase_operation",
			body: `{"operation":"ADD","operands":[1,2]}`,
		},
		{
			name: "operand_string",
			body: `{"operation":"sqrt","operands":["2"]}`,
		},
		{
			name: "root_null",
			body: `null`,
		},
		{
			name: "root_array",
			body: `[]`,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var value any
			decodeOne(t, tt.body, &value)
			if err := checkSchema(doc, schema, value); err == nil {
				t.Fatalf("contract accepted invalid request %s", tt.body)
			}
		})
	}
}

// Follow the actual path/method response reference, rather than merely checking
// a similarly named component. This catches response-schema wiring drift too.
func assertDocumentedResponse(t *testing.T, method, path string, status int, media, body string) {
	t.Helper()
	assertRouteContract(t, method, path, status)
	doc := loadContract(t)
	pathItem, ok := contractObject(t, doc, "paths")[path].(map[string]any)
	if !ok {
		// Unknown routes are verified through Problem404.
		return
	}
	operation, ok := pathItem[strings.ToLower(method)].(map[string]any)
	if !ok {
		key := "post"
		if path == "/health" {
			key = "get"
		}
		operation = contractObject(t, pathItem, key)
	}
	response := contractObject(t, operation, "responses", strconv.Itoa(status))
	content, hasContent := response["content"].(map[string]any)
	if method == "HEAD" && path == "/health" && hasContent {
		t.Fatalf("HEAD %s status %d must not document response content", path, status)
	}
	if method == "HEAD" || status == 204 {
		if body != "" {
			t.Fatalf("body forbidden for %s/%d: %q", method, status, body)
		}
		return
	}
	if !hasContent {
		t.Fatalf("missing content schema for %s %s %d", method, path, status)
	}
	entry, ok := content[media].(map[string]any)
	if !ok {
		t.Fatalf("undocumented media %q for %s %s %d", media, method, path, status)
	}
	var value any
	decodeOne(t, body, &value)
	if err := checkSchema(doc, contractObject(t, entry, "schema"), value); err != nil {
		t.Fatalf("%s %s %d violates documented response: %v", method, path, status, err)
	}
}

func TestContractCheckerRejectsUnsupportedSchemas(t *testing.T) {
	cases := []struct {
		name   string
		schema string
	}{
		{
			name:   "reference_sibling_constraint",
			schema: `{"$ref":"#/number","minimum":10}`,
		},
		{
			name:   "unknown_type",
			schema: `{"type":"unexpected"}`,
		},
		{
			name:   "unsupported_alternative_cannot_be_hidden_by_a_match",
			schema: `{"oneOf":[{"type":"number"},{"type":"string","pattern":"x"}]}`,
		},
		{
			name:   "malformed_minimum",
			schema: `{"type":"number","minimum":"10"}`,
		},
		{
			name:   "unsupported_absent_property",
			schema: `{"type":"object","properties":{"optional":{"type":"string","pattern":"x"}}}`,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var schema map[string]any
			decodeOne(t, tt.schema, &schema)
			doc := map[string]any{"number": map[string]any{"type": "number"}}
			value := any(float64(1))
			if tt.name == "unsupported_absent_property" {
				value = map[string]any{}
			}
			if err := checkSchema(doc, schema, value); err == nil {
				t.Fatalf("unsupported schema silently accepted: %s", tt.schema)
			}
		})
	}
}

func TestOpenAPIHeadResponsesHaveNoContent(t *testing.T) {
	doc := loadContract(t)
	responses := contractObject(t, doc, "paths", "/health", "head", "responses")
	for _, status := range []string{"200", "403", "500"} {
		t.Run(status, func(t *testing.T) {
			response := contractObject(t, responses, status)
			if _, exists := response["content"]; exists {
				t.Fatalf("HEAD %s advertises a response body; HEAD never carries content", status)
			}
			description, ok := response["description"].(string)
			want := "application/problem+json"
			if status == "200" {
				want = "application/json"
			}
			if !ok || !strings.Contains(description, "Content-Type: "+want) {
				t.Fatalf("HEAD %s must describe its Content-Type header without a body schema", status)
			}

		})
	}
}
