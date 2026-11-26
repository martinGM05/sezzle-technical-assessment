package expression_test

import (
	"math"
	"strings"
	"testing"

	"github.com/martinGM05/sezzle-technical-assessment/backend/internal/expression"
)

func assertValue(t *testing.T, input string, want float64) {
	t.Helper()
	got, err := expression.Evaluate(input)
	if err != nil {
		t.Fatalf("Evaluate(%q): unexpected error: %v", input, err)
	}
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-want) > 1e-14*math.Max(1, math.Abs(want)) {
		t.Fatalf("Evaluate(%q) = %.17g; want %.17g (finite, relative/absolute tolerance 1e-14)", input, got, want)
	}
}

func TestLiterals(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "zero",
			input: "0",
			want:  0,
		},
		{
			name:  "integer",
			input: "42",
			want:  42,
		},
		{
			name:  "decimal",
			input: "12.75",
			want:  12.75,
		},
		{
			name:  "leading_decimal_point",
			input: ".5",
			want:  .5,
		},
		{
			name:  "trailing_decimal_point",
			input: "5.",
			want:  5,
		},
		{
			name:  "scientific_small",
			input: "1.25e-2",
			want:  .0125,
		},
		{
			name:  "scientific_large",
			input: "2E+3",
			want:  2000,
		},
		{
			name:  "ascii_whitespace",
			input: " \t\r\n42 \n",
			want:  42,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertValue(t, tt.input, tt.want) })
	}
}

func assertError(t *testing.T, input, wantCode string) {
	t.Helper()
	got, err := expression.Evaluate(input)
	if err == nil {
		t.Fatalf("Evaluate(%q) = %v without error; want %s", input, got, wantCode)
	}
	if err.Code != wantCode || err.Detail == "" || got != 0 {
		t.Fatalf("Evaluate(%q) = %v, error=%+v; want zero and safe %s error", input, got, err, wantCode)
	}
}

func TestBasicArithmetic(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "addition",
			input: "2 + 3",
			want:  5,
		},
		{
			name:  "negative_subtraction",
			input: "3 - 5",
			want:  -2,
		},
		{
			name:  "multiplication",
			input: "3 * 4",
			want:  12,
		},
		{
			name:  "division",
			input: "5 / 2",
			want:  2.5,
		},
		{
			name:  "multiplication_before_addition",
			input: "2 + 3 * 4",
			want:  14,
		},
		{
			name:  "left_associative_subtraction",
			input: "10 - 3 - 2",
			want:  5,
		},
		{
			name:  "left_associative_division",
			input: "20 / 2 / 5",
			want:  2,
		},
		{
			name:  "same_priority_multiply_divide",
			input: "8 / 2 * 3",
			want:  12,
		},
		{
			name:  "decimal_precision",
			input: "0.1 + 0.2",
			want:  .3,
		},
		{
			name:  "nonterminating_division",
			input: "1 / 3",
			want:  1.0 / 3,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertValue(t, tt.input, tt.want) })
	}
}

func TestLexicalAndBasicSyntaxErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
		code  string
	}{
		{
			name:  "empty",
			input: "",
			code:  "INVALID_EXPRESSION",
		},
		{
			name:  "blank",
			input: " \t\n\r",
			code:  "INVALID_EXPRESSION",
		},
		{
			name:  "unknown_character",
			input: "2 @ 3",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "not_a_number",
			input: "NaN",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "infinity",
			input: "Inf",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "hexadecimal_not_supported",
			input: "0x1p2",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "underscore_not_supported",
			input: "1_000",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "unicode_operator",
			input: "2 × 3",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "unicode_space",
			input: "\u00a02",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "decimal_point_only",
			input: ".",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "missing_exponent",
			input: "1e+",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "adjacent_numbers",
			input: "1 2",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "repeated_decimal_point",
			input: "1.2.3",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "missing_right_operand",
			input: "2 +",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "invalid_operator_sequence",
			input: "2 + * 3",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "division_by_zero",
			input: "10 / 0",
			code:  "DIVISION_BY_ZERO",
		},
		{
			name:  "literal_overflow",
			input: "1e309",
			code:  "INVALID_RESULT",
		},
		{
			name:  "arithmetic_overflow",
			input: "1e308 * 10",
			code:  "INVALID_RESULT",
		},
		{
			name:  "complete_syntax_precedes_arithmetic",
			input: "1/0 +",
			code:  "SYNTAX_ERROR",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertError(t, tt.input, tt.code) })
	}
}

func TestGroupingAndUnarySigns(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "parentheses_override_precedence",
			input: "(2 + 3) * 4",
			want:  20,
		},
		{
			name:  "nested_parentheses",
			input: "2 * (3 + (4 - 1))",
			want:  12,
		},
		{
			name:  "leading_negative",
			input: "-5 + 3",
			want:  -2,
		},
		{
			name:  "unary_plus",
			input: "+2",
			want:  2,
		},
		{
			name:  "sign_chain",
			input: "--+-2",
			want:  -2,
		},
		{
			name:  "sign_after_operator",
			input: "2*-3 - -4",
			want:  -2,
		},
		{
			name:  "negative_parentheses",
			input: "-(2+3)",
			want:  -5,
		},
		{
			name:  "exact_parenthesis_depth",
			input: strings.Repeat("(", 64) + "1" + strings.Repeat(")", 64),
			want:  1,
		},
		{
			name:  "exact_unary_depth",
			input: strings.Repeat("-", 64) + "1",
			want:  1,
		},
		{
			name:  "exact_length",
			input: "1" + strings.Repeat(" ", 4095),
			want:  1,
		},
		{
			name:  "flat_chain_does_not_consume_nesting_budget",
			input: strings.Repeat("1+", 2047) + "1 ",
			want:  2048,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertValue(t, tt.input, tt.want) })
	}
}

func TestGroupingAndResourceErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
		code  string
	}{
		{
			name:  "unclosed_parenthesis",
			input: "(2 + 3",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "extra_closing_parenthesis",
			input: "2 + 3)",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "empty_parentheses",
			input: "()",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "implicit_multiplication",
			input: "2(3)",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "sign_only",
			input: "-",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "computed_zero_divisor",
			input: "10 / (5 - 5)",
			code:  "DIVISION_BY_ZERO",
		},
		{
			name:  "negative_zero_divisor",
			input: "1/-0",
			code:  "DIVISION_BY_ZERO",
		},
		{
			name:  "excess_parenthesis_depth",
			input: strings.Repeat("(", 65) + "1" + strings.Repeat(")", 65),
			code:  "EXPRESSION_LIMIT",
		},
		{
			name:  "excess_unary_depth",
			input: strings.Repeat("-", 65) + "1",
			code:  "EXPRESSION_LIMIT",
		},
		{
			name:  "mixed_recursion_budget",
			input: "-(" + strings.Repeat("-", 63) + "1)",
			code:  "EXPRESSION_LIMIT",
		},
		{
			name:  "excess_length",
			input: "1" + strings.Repeat(" ", 4096),
			code:  "EXPRESSION_LIMIT",
		},
		{
			name:  "length_checked_before_blank",
			input: strings.Repeat(" ", 4097),
			code:  "EXPRESSION_LIMIT",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertError(t, tt.input, tt.code) })
	}
}

func TestPower(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "right_associative_power",
			input: "2 ^ 3 ^ 2",
			want:  512,
		},
		{
			name:  "negative_sign_after_power",
			input: "-2 ^ 2",
			want:  -4,
		},
		{
			name:  "parenthesized_negative_base",
			input: "(-2) ^ 2",
			want:  4,
		},
		{
			name:  "negative_exponent",
			input: "2 ^ -3",
			want:  .125,
		},
		{
			name:  "signed_power_rhs",
			input: "2 ^ -3 ^ 2",
			want:  1.0 / 512,
		},
		{
			name:  "unary_plus_power",
			input: "+2 ^ 3",
			want:  8,
		},
		{
			name:  "fractional_exponent",
			input: "9 ^ .5",
			want:  3,
		},
		{
			name:  "power_binds_before_multiply",
			input: "2 * 3 ^ 2 + 1",
			want:  19,
		},
		{
			name:  "zero_to_zero",
			input: "0 ^ 0",
			want:  1,
		},
		{
			name:  "explicit_left_power",
			input: "(2 ^ 3) ^ 2",
			want:  64,
		},
		{
			name:  "exact_exponent_depth",
			input: strings.Repeat("1^", 64) + "1",
			want:  1,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertValue(t, tt.input, tt.want) })
	}
}

func TestPowerErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
		code  string
	}{
		{
			name:  "non_real_power",
			input: "(-2)^.5",
			code:  "INVALID_RESULT",
		},
		{
			name:  "power_overflow",
			input: "10^309",
			code:  "INVALID_RESULT",
		},
		{
			name:  "zero_to_negative",
			input: "0^-1",
			code:  "INVALID_RESULT",
		},
		{
			name:  "missing_exponent",
			input: "2^",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "excess_power_depth",
			input: strings.Repeat("1^", 65) + "1",
			code:  "EXPRESSION_LIMIT",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertError(t, tt.input, tt.code) })
	}
}

func TestSquareRoot(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "expression_argument",
			input: "sqrt(9 + 7)",
			want:  4,
		},
		{
			name:  "nested_calls",
			input: "sqrt(sqrt(16))",
			want:  2,
		},
		{
			name:  "irrational",
			input: "sqrt(2)",
			want:  math.Sqrt2,
		},
		{
			name:  "decimal",
			input: "sqrt(.25)",
			want:  .5,
		},
		{
			name:  "zero",
			input: "sqrt(0)",
			want:  0,
		},
		{
			name:  "spacing_and_sign",
			input: "-sqrt ( 9 ) + 5",
			want:  2,
		},
		{
			name:  "power_after_call",
			input: "sqrt(4)^3",
			want:  8,
		},
		{
			name:  "exact_call_depth",
			input: strings.Repeat("sqrt(", 64) + "1" + strings.Repeat(")", 64),
			want:  1,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertValue(t, tt.input, tt.want) })
	}
}

func TestSquareRootErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
		code  string
	}{
		{
			name:  "negative_domain",
			input: "sqrt(-1)",
			code:  "INVALID_OPERANDS",
		},
		{
			name:  "uppercase_function",
			input: "SQRT(4)",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "missing_parentheses",
			input: "sqrt 4",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "empty_argument",
			input: "sqrt()",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "unclosed_function",
			input: "sqrt(4",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "multiple_arguments",
			input: "sqrt(4,9)",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "identifier_suffix",
			input: "sqrtx(4)",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "missing_multiplication",
			input: "2sqrt(4)",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "excess_call_depth",
			input: strings.Repeat("sqrt(", 65) + "1" + strings.Repeat(")", 65),
			code:  "EXPRESSION_LIMIT",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertError(t, tt.input, tt.code) })
	}
}

func TestPercentage(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "percentage_of_value",
			input: "200 * 15%",
			want:  30,
		},
		{
			name:  "percentage_is_not_contextual",
			input: "200 + 15%",
			want:  200.15,
		},
		{
			name:  "negative_percent",
			input: "-50%",
			want:  -.5,
		},
		{
			name:  "zero_percent",
			input: "0%",
			want:  0,
		},
		{
			name:  "fractional_percent",
			input: ".5%",
			want:  .005,
		},
		{
			name:  "over_one_hundred",
			input: "250%",
			want:  2.5,
		},
		{
			name:  "repeated_postfix",
			input: "50%%",
			want:  .005,
		},
		{
			name:  "postfix_precedes_power",
			input: "50%^2",
			want:  .25,
		},
		{
			name:  "postfix_on_exponent",
			input: "2^3%",
			want:  math.Pow(2, .03),
		},
		{
			name:  "postfix_on_parenthesized_power",
			input: "(2^3)%",
			want:  .08,
		},
		{
			name:  "postfix_after_sqrt",
			input: "sqrt(25)%",
			want:  .05,
		},
		{
			name:  "postfix_flat_chain_at_length_limit",
			input: "1" + strings.Repeat("%", 4095),
			want:  0,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertValue(t, tt.input, tt.want) })
	}
}

func TestPercentageErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "prefix_percent",
			input: "%50",
		},
		{
			name:  "percent_is_not_modulo",
			input: "5%2",
		},
		{
			name:  "missing_operand",
			input: "2 + %",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertError(t, tt.input, "SYNTAX_ERROR") })
	}
}

func TestNumericBoundaries(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "largest_finite_literal",
			input: "1.7976931348623157e308",
			want:  math.MaxFloat64,
		},
		{
			name:  "smallest_subnormal_literal",
			input: "5e-324",
			want:  math.SmallestNonzeroFloat64,
		},
		{
			name:  "literal_underflow",
			input: "1e-999",
			want:  0,
		},
		{
			name:  "arithmetic_underflow",
			input: "5e-324 / 2",
			want:  0,
		},
		{
			name:  "negative_zero",
			input: "-0",
			want:  math.Copysign(0, -1),
		},
		{
			name:  "positive_sign_preserves_negative_zero",
			input: "+-0",
			want:  math.Copysign(0, -1),
		},
		{
			name:  "sqrt_negative_zero",
			input: "sqrt(-0)",
			want:  math.Copysign(0, -1),
		},
		{
			name:  "integer_rounding",
			input: "9007199254740993",
			want:  9007199254740992,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expression.Evaluate(tt.input)
			if err != nil || math.Float64bits(got) != math.Float64bits(tt.want) {
				t.Fatalf("Evaluate(%q) = %.17g (%x), error=%v; want exact %.17g (%x)", tt.input, got, math.Float64bits(got), err, tt.want, math.Float64bits(tt.want))
			}
		})
	}
}

func TestNumericalErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
		code  string
	}{
		{
			name:  "addition_overflow",
			input: "1.7976931348623157e308 + 1.7976931348623157e308",
			code:  "INVALID_RESULT",
		},
		{
			name:  "division_overflow",
			input: "1e308 / 1e-308",
			code:  "INVALID_RESULT",
		},
		{
			name:  "intermediate_overflow_not_cancelled",
			input: "1e308 * 10 / 10",
			code:  "INVALID_RESULT",
		},
		{
			name:  "left_error_first",
			input: "1/0 + sqrt(-1)",
			code:  "DIVISION_BY_ZERO",
		},
		{
			name:  "right_subtree_error",
			input: "2+(1/0)",
			code:  "DIVISION_BY_ZERO",
		},
		{
			name:  "signed_infinity_spelling",
			input: "-Infinity",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "invalid_character_after_domain_failure",
			input: "sqrt(-1) @",
			code:  "SYNTAX_ERROR",
		},
		{
			name:  "null_byte",
			input: "1\x00+2",
			code:  "SYNTAX_ERROR",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertError(t, tt.input, tt.code) })
	}
}

func FuzzEvaluate(f *testing.F) {
	for _, input := range []string{"2+3*4", "-2^2", "sqrt(-1)", "1e309", "", "()", "1/0+", strings.Repeat("-", 65) + "1", "200*15%"} {
		f.Add(input)
	}
	f.Fuzz(func(t *testing.T, input string) {
		got, err := expression.Evaluate(input)
		again, againErr := expression.Evaluate(input)
		if err != nil {
			if againErr == nil || *againErr != *err || got != 0 || again != 0 {
				t.Fatalf("non-deterministic error for %q: %v / %v", input, err, againErr)
			}
			return
		}
		if againErr != nil || math.IsNaN(got) || math.IsInf(got, 0) || math.Float64bits(got) != math.Float64bits(again) {
			t.Fatalf("non-finite or non-deterministic result for %q: %v / %v, %v", input, got, again, againErr)
		}
	})
}
