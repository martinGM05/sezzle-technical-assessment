package calculator

import (
	"math"
	"testing"
)

// Exact comparison is the default, including zeros and subnormal boundaries.
// Irrational/repeating/rounded results use a relative tolerance with no absolute
// floor, so zero cannot accidentally match a tiny nonzero expected result.
// Finite 0^0 is defined as 1; underflow follows IEEE 754 rounding to zero.
type calculationCase struct {
	name              string
	operation         string
	operands          []float64
	want              float64
	wantCode          string
	relativeTolerance float64
}

func runCalculationCases(t *testing.T, cases []calculationCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.operation, tt.operands)
			if tt.wantCode != "" {
				if err == nil {
					t.Fatalf("Calculate(%q, %v) = %g, nil; want error code %s", tt.operation, tt.operands, got, tt.wantCode)
				}
				if err.Code != tt.wantCode {
					t.Fatalf("Calculate(%q, %v): error code = %q; want %q (detail: %q)", tt.operation, tt.operands, err.Code, tt.wantCode, err.Detail)
				}
				if err.Detail == "" || err.Error() != err.Detail {
					t.Fatalf("Calculate(%q, %v): error must expose its nonempty detail; got %#v", tt.operation, tt.operands, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Calculate(%q, %v): unexpected error %s (%s)", tt.operation, tt.operands, err.Code, err.Detail)
			}
			if math.IsNaN(got) || math.IsInf(got, 0) {
				t.Fatalf("Calculate(%q, %v) = %g; want finite result %.17g", tt.operation, tt.operands, got, tt.want)
			}
			if got == tt.want {
				return
			}
			if tt.relativeTolerance > 0 && math.Abs(got-tt.want) <= tt.relativeTolerance*math.Abs(tt.want) {
				return
			}
			t.Fatalf("Calculate(%q, %v) = %.17g; want %.17g (relative tolerance %g)", tt.operation, tt.operands, got, tt.want, tt.relativeTolerance)
		})
	}
}

func TestCalculateArithmetic(t *testing.T) {
	runCalculationCases(t, []calculationCase{
		{
			name:      "add_positive",
			operation: "add",
			operands:  []float64{2, 3},
			want:      5,
		},
		{
			name:      "add_negative",
			operation: "add",
			operands:  []float64{-2, -3},
			want:      -5,
		},
		{
			name:      "add_mixed_sign_cancellation",
			operation: "add",
			operands:  []float64{-2, 2},
			want:      0,
		},
		{
			name:      "add_mixed_sign_fraction",
			operation: "add",
			operands:  []float64{-2, 0.5},
			want:      -1.5,
		},
		{
			name:      "add_decimal",
			operation: "add",
			operands:  []float64{1.25, 2.5},
			want:      3.75,
		},
		{
			name:      "add_zero_left",
			operation: "add",
			operands:  []float64{0, 7},
			want:      7,
		},
		{
			name:      "add_zero_right",
			operation: "add",
			operands:  []float64{7, 0},
			want:      7,
		},
		{
			name:      "add_both_zero",
			operation: "add",
			operands:  []float64{0, 0},
			want:      0,
		},
		{
			name:      "subtract_positive",
			operation: "subtract",
			operands:  []float64{5, 3},
			want:      2,
		},
		{
			name:      "subtract_negative",
			operation: "subtract",
			operands:  []float64{-5, -3},
			want:      -2,
		},
		{
			name:      "subtract_mixed_sign",
			operation: "subtract",
			operands:  []float64{2, -3},
			want:      5,
		},
		{
			name:      "subtract_mixed_sign_negative",
			operation: "subtract",
			operands:  []float64{-2, 3},
			want:      -5,
		},
		{
			name:      "subtract_decimal_negative_result",
			operation: "subtract",
			operands:  []float64{1.5, 2},
			want:      -0.5,
		},
		{
			name:      "subtract_zero_left",
			operation: "subtract",
			operands:  []float64{0, 3},
			want:      -3,
		},
		{
			name:      "subtract_zero_right",
			operation: "subtract",
			operands:  []float64{3, 0},
			want:      3,
		},
		{
			name:      "subtract_both_zero",
			operation: "subtract",
			operands:  []float64{0, 0},
			want:      0,
		},
		{
			name:      "multiply_positive",
			operation: "multiply",
			operands:  []float64{2, 3},
			want:      6,
		},
		{
			name:      "multiply_negative",
			operation: "multiply",
			operands:  []float64{-2, -3},
			want:      6,
		},
		{
			name:      "multiply_mixed_sign",
			operation: "multiply",
			operands:  []float64{-2, 3},
			want:      -6,
		},
		{
			name:      "multiply_decimal",
			operation: "multiply",
			operands:  []float64{1.5, 0.25},
			want:      0.375,
		},
		{
			name:      "multiply_zero_left",
			operation: "multiply",
			operands:  []float64{0, -3},
			want:      0,
		},
		{
			name:      "multiply_zero_right",
			operation: "multiply",
			operands:  []float64{-3, 0},
			want:      0,
		},
		{
			name:      "multiply_both_zero",
			operation: "multiply",
			operands:  []float64{0, 0},
			want:      0,
		},
		{
			name:      "divide_positive",
			operation: "divide",
			operands:  []float64{8, 2},
			want:      4,
		},
		{
			name:      "divide_negative",
			operation: "divide",
			operands:  []float64{-8, -2},
			want:      4,
		},
		{
			name:      "divide_mixed_sign",
			operation: "divide",
			operands:  []float64{-8, 2},
			want:      -4,
		},
		{
			name:      "divide_negative_denominator",
			operation: "divide",
			operands:  []float64{8, -2},
			want:      -4,
		},
		{
			name:      "divide_decimal",
			operation: "divide",
			operands:  []float64{1.5, 0.5},
			want:      3,
		},
		{
			name:      "divide_fractional_result",
			operation: "divide",
			operands:  []float64{1, 4},
			want:      0.25,
		},
		{
			name:      "divide_zero_numerator",
			operation: "divide",
			operands:  []float64{0, 2},
			want:      0,
		},
		{
			name:              "divide_repeating",
			operation:         "divide",
			operands:          []float64{1, 3},
			want:              0.3333333333333333,
			relativeTolerance: 1e-14,
		},
		{
			name:              "divide_negative_repeating",
			operation:         "divide",
			operands:          []float64{-2, 3},
			want:              -0.6666666666666666,
			relativeTolerance: 1e-14,
		},
		{
			name:      "divide_positive_zero",
			operation: "divide",
			operands:  []float64{1, 0},
			wantCode:  "DIVISION_BY_ZERO",
		},
		{
			name:      "divide_negative_zero",
			operation: "divide",
			operands:  []float64{1, math.Copysign(0, -1)},
			wantCode:  "DIVISION_BY_ZERO",
		},
		{
			name:      "divide_zero_over_zero",
			operation: "divide",
			operands:  []float64{0, 0},
			wantCode:  "DIVISION_BY_ZERO",
		},
	})
}

func TestCalculatePower(t *testing.T) {
	runCalculationCases(t, []calculationCase{
		{
			name:      "positive_exponent",
			operation: "power",
			operands:  []float64{2, 3},
			want:      8,
		},
		{
			name:      "negative_exponent",
			operation: "power",
			operands:  []float64{2, -2},
			want:      0.25,
		},
		{
			name:      "zero_exponent",
			operation: "power",
			operands:  []float64{3, 0},
			want:      1,
		},
		{
			name:      "fractional_exponent",
			operation: "power",
			operands:  []float64{9, 0.5},
			want:      3,
		},
		{
			name:      "negative_base_even",
			operation: "power",
			operands:  []float64{-2, 4},
			want:      16,
		},
		{
			name:      "negative_base_odd",
			operation: "power",
			operands:  []float64{-2, 3},
			want:      -8,
		},
		{
			name:      "negative_base_negative_exponent",
			operation: "power",
			operands:  []float64{-2, -3},
			want:      -0.125,
		},
		{
			name:      "zero_base_positive_exponent",
			operation: "power",
			operands:  []float64{0, 3},
			want:      0,
		},
		{
			name:      "zero_base_fractional_exponent",
			operation: "power",
			operands:  []float64{0, 0.5},
			want:      0,
		},
		{
			name:      "zero_to_zero_is_one",
			operation: "power",
			operands:  []float64{0, 0},
			want:      1,
		},
		{
			name:      "negative_zero_to_zero_is_one",
			operation: "power",
			operands:  []float64{math.Copysign(0, -1), 0},
			want:      1,
		},
		{
			name:      "decimal_base",
			operation: "power",
			operands:  []float64{0.5, 3},
			want:      0.125,
		},
		{
			name:              "irrational_result",
			operation:         "power",
			operands:          []float64{2, 0.5},
			want:              1.4142135623730951,
			relativeTolerance: 1e-14,
		},
		{
			name:      "negative_base_fractional_exponent",
			operation: "power",
			operands:  []float64{-1, 0.5},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "negative_base_fractional_reciprocal",
			operation: "power",
			operands:  []float64{-8, 1.0 / 3},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "zero_negative_exponent",
			operation: "power",
			operands:  []float64{0, -1},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "overflow",
			operation: "power",
			operands:  []float64{math.MaxFloat64, 2},
			wantCode:  "INVALID_RESULT",
		},
	})
}

func TestCalculateSquareRoot(t *testing.T) {
	runCalculationCases(t, []calculationCase{
		{
			name:      "positive",
			operation: "sqrt",
			operands:  []float64{9},
			want:      3,
		},
		{
			name:      "zero",
			operation: "sqrt",
			operands:  []float64{0},
			want:      0,
		},
		{
			name:      "negative_zero",
			operation: "sqrt",
			operands:  []float64{math.Copysign(0, -1)},
			want:      0,
		},
		{
			name:      "decimal",
			operation: "sqrt",
			operands:  []float64{0.25},
			want:      0.5,
		},
		{
			name:              "irrational_result",
			operation:         "sqrt",
			operands:          []float64{2},
			want:              1.4142135623730951,
			relativeTolerance: 1e-14,
		},
		{
			name:      "negative_integer",
			operation: "sqrt",
			operands:  []float64{-1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "negative_decimal",
			operation: "sqrt",
			operands:  []float64{-0.25},
			wantCode:  "INVALID_OPERANDS",
		},
	})
}

func TestCalculatePercentage(t *testing.T) {
	runCalculationCases(t, []calculationCase{
		{
			name:      "positive",
			operation: "percentage",
			operands:  []float64{50},
			want:      0.5,
		},
		{
			name:      "negative",
			operation: "percentage",
			operands:  []float64{-50},
			want:      -0.5,
		},
		{
			name:      "zero",
			operation: "percentage",
			operands:  []float64{0},
			want:      0,
		},
		{
			name:      "fractional",
			operation: "percentage",
			operands:  []float64{0.5},
			want:      0.005,
		},
		{
			name:      "greater_than_100",
			operation: "percentage",
			operands:  []float64{150.5},
			want:      1.505,
		},
	})
}

func TestCalculateValidation(t *testing.T) {
	runCalculationCases(t, []calculationCase{
		{
			name:      "unsupported",
			operation: "unknown",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERATION",
		},
		{
			name:      "empty",
			operation: "",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERATION",
		},
		{
			name:      "uppercase",
			operation: "ADD",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERATION",
		},
		{
			name:      "leading_space",
			operation: " add",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERATION",
		},
		{
			name:      "trailing_space",
			operation: "add ",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERATION",
		},
		{
			name:      "whitespace_only",
			operation: " ",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERATION",
		},
		{
			name:      "add_nil",
			operation: "add",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_empty",
			operation: "add",
			operands:  []float64{},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_excessive",
			operation: "add",
			operands:  []float64{1, 2, 3},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_unary",
			operation: "add",
			operands:  []float64{1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_nil",
			operation: "subtract",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_empty",
			operation: "subtract",
			operands:  []float64{},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_excessive",
			operation: "subtract",
			operands:  []float64{1, 2, 3},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_unary",
			operation: "subtract",
			operands:  []float64{1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_nil",
			operation: "multiply",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_empty",
			operation: "multiply",
			operands:  []float64{},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_excessive",
			operation: "multiply",
			operands:  []float64{1, 2, 3},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_unary",
			operation: "multiply",
			operands:  []float64{1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_nil",
			operation: "divide",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_empty",
			operation: "divide",
			operands:  []float64{},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_excessive",
			operation: "divide",
			operands:  []float64{1, 2, 3},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_unary",
			operation: "divide",
			operands:  []float64{1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_nil",
			operation: "power",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_empty",
			operation: "power",
			operands:  []float64{},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_excessive",
			operation: "power",
			operands:  []float64{1, 2, 3},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_unary",
			operation: "power",
			operands:  []float64{1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "sqrt_nil",
			operation: "sqrt",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "sqrt_empty",
			operation: "sqrt",
			operands:  []float64{},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "sqrt_excessive",
			operation: "sqrt",
			operands:  []float64{1, 2, 3},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "sqrt_binary",
			operation: "sqrt",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "percentage_nil",
			operation: "percentage",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "percentage_empty",
			operation: "percentage",
			operands:  []float64{},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "percentage_excessive",
			operation: "percentage",
			operands:  []float64{1, 2, 3},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "percentage_binary",
			operation: "percentage",
			operands:  []float64{1, 2},
			wantCode:  "INVALID_OPERANDS",
		},
	})
}

func TestCalculateNumericBoundaries(t *testing.T) {
	runCalculationCases(t, []calculationCase{
		{
			name:      "add_nan_operand_1",
			operation: "add",
			operands:  []float64{math.NaN(), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_positive_infinity_operand_1",
			operation: "add",
			operands:  []float64{math.Inf(1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_negative_infinity_operand_1",
			operation: "add",
			operands:  []float64{math.Inf(-1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_nan_operand_2",
			operation: "add",
			operands:  []float64{1, math.NaN()},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_positive_infinity_operand_2",
			operation: "add",
			operands:  []float64{1, math.Inf(1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "add_negative_infinity_operand_2",
			operation: "add",
			operands:  []float64{1, math.Inf(-1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_nan_operand_1",
			operation: "subtract",
			operands:  []float64{math.NaN(), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_positive_infinity_operand_1",
			operation: "subtract",
			operands:  []float64{math.Inf(1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_negative_infinity_operand_1",
			operation: "subtract",
			operands:  []float64{math.Inf(-1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_nan_operand_2",
			operation: "subtract",
			operands:  []float64{1, math.NaN()},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_positive_infinity_operand_2",
			operation: "subtract",
			operands:  []float64{1, math.Inf(1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "subtract_negative_infinity_operand_2",
			operation: "subtract",
			operands:  []float64{1, math.Inf(-1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_nan_operand_1",
			operation: "multiply",
			operands:  []float64{math.NaN(), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_positive_infinity_operand_1",
			operation: "multiply",
			operands:  []float64{math.Inf(1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_negative_infinity_operand_1",
			operation: "multiply",
			operands:  []float64{math.Inf(-1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_nan_operand_2",
			operation: "multiply",
			operands:  []float64{1, math.NaN()},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_positive_infinity_operand_2",
			operation: "multiply",
			operands:  []float64{1, math.Inf(1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "multiply_negative_infinity_operand_2",
			operation: "multiply",
			operands:  []float64{1, math.Inf(-1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_nan_operand_1",
			operation: "divide",
			operands:  []float64{math.NaN(), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_positive_infinity_operand_1",
			operation: "divide",
			operands:  []float64{math.Inf(1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_negative_infinity_operand_1",
			operation: "divide",
			operands:  []float64{math.Inf(-1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_nan_operand_2",
			operation: "divide",
			operands:  []float64{1, math.NaN()},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_positive_infinity_operand_2",
			operation: "divide",
			operands:  []float64{1, math.Inf(1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "divide_negative_infinity_operand_2",
			operation: "divide",
			operands:  []float64{1, math.Inf(-1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_nan_operand_1",
			operation: "power",
			operands:  []float64{math.NaN(), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_positive_infinity_operand_1",
			operation: "power",
			operands:  []float64{math.Inf(1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_negative_infinity_operand_1",
			operation: "power",
			operands:  []float64{math.Inf(-1), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_nan_operand_2",
			operation: "power",
			operands:  []float64{1, math.NaN()},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_positive_infinity_operand_2",
			operation: "power",
			operands:  []float64{1, math.Inf(1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "power_negative_infinity_operand_2",
			operation: "power",
			operands:  []float64{1, math.Inf(-1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "sqrt_nan_operand_1",
			operation: "sqrt",
			operands:  []float64{math.NaN()},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "sqrt_positive_infinity_operand_1",
			operation: "sqrt",
			operands:  []float64{math.Inf(1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "sqrt_negative_infinity_operand_1",
			operation: "sqrt",
			operands:  []float64{math.Inf(-1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "percentage_nan_operand_1",
			operation: "percentage",
			operands:  []float64{math.NaN()},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "percentage_positive_infinity_operand_1",
			operation: "percentage",
			operands:  []float64{math.Inf(1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "percentage_negative_infinity_operand_1",
			operation: "percentage",
			operands:  []float64{math.Inf(-1)},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "addition_overflow",
			operation: "add",
			operands:  []float64{math.MaxFloat64, math.MaxFloat64},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "negative_addition_overflow",
			operation: "add",
			operands:  []float64{-math.MaxFloat64, -math.MaxFloat64},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "subtraction_overflow",
			operation: "subtract",
			operands:  []float64{math.MaxFloat64, -math.MaxFloat64},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "multiplication_overflow",
			operation: "multiply",
			operands:  []float64{math.MaxFloat64, 2},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "division_overflow",
			operation: "divide",
			operands:  []float64{math.MaxFloat64, 0.5},
			wantCode:  "INVALID_RESULT",
		},
		{
			name:      "largest_finite_sum",
			operation: "add",
			operands:  []float64{math.MaxFloat64, 0},
			want:      math.MaxFloat64,
		},
		{
			name:      "smallest_finite_sum",
			operation: "add",
			operands:  []float64{-math.MaxFloat64, 0},
			want:      -math.MaxFloat64,
		},
		{
			name:      "largest_finite_product",
			operation: "multiply",
			operands:  []float64{math.MaxFloat64, 1},
			want:      math.MaxFloat64,
		},
		{
			name:      "largest_finite_division",
			operation: "divide",
			operands:  []float64{math.MaxFloat64, 1},
			want:      math.MaxFloat64,
		},
		{
			name:      "smallest_subnormal_preserved",
			operation: "add",
			operands:  []float64{math.SmallestNonzeroFloat64, 0},
			want:      math.SmallestNonzeroFloat64,
		},
		{
			name:      "division_subnormal_preserved",
			operation: "divide",
			operands:  []float64{math.SmallestNonzeroFloat64, 1},
			want:      math.SmallestNonzeroFloat64,
		},
		{
			name:      "subnormal_addition",
			operation: "add",
			operands:  []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
			want:      2 * math.SmallestNonzeroFloat64,
		},
		{
			name:      "division_underflow_rounds_to_zero",
			operation: "divide",
			operands:  []float64{math.SmallestNonzeroFloat64, 2},
			want:      0,
		},
		{
			name:      "power_underflow_rounds_to_zero",
			operation: "power",
			operands:  []float64{2, -1075},
			want:      0,
		},
		{
			name:      "percentage_underflow_rounds_to_zero",
			operation: "percentage",
			operands:  []float64{math.SmallestNonzeroFloat64},
			want:      0,
		},
		{
			name:      "precision_integer_boundary",
			operation: "add",
			operands:  []float64{9007199254740992, 1},
			want:      9007199254740992,
		},
		{
			name:      "adjacent_floats_difference",
			operation: "subtract",
			operands:  []float64{math.Nextafter(1, 2), 1},
			want:      0x1p-52,
		},
		{
			name:              "decimal_sum_rounding",
			operation:         "add",
			operands:          []float64{0.1, 0.2},
			want:              0.3,
			relativeTolerance: 1e-14,
		},
		{
			name:              "decimal_product_rounding",
			operation:         "multiply",
			operands:          []float64{0.1, 0.2},
			want:              0.02,
			relativeTolerance: 1e-14,
		},
		{
			name:              "decimal_subtraction_rounding",
			operation:         "subtract",
			operands:          []float64{0.3, 0.2},
			want:              0.1,
			relativeTolerance: 1e-14,
		},
		{
			name:              "sqrt_largest_finite",
			operation:         "sqrt",
			operands:          []float64{math.MaxFloat64},
			want:              1.3407807929942596e154,
			relativeTolerance: 1e-14,
		},
		{
			name:              "sqrt_smallest_subnormal",
			operation:         "sqrt",
			operands:          []float64{math.SmallestNonzeroFloat64},
			want:              2.2227587494850775e-162,
			relativeTolerance: 1e-14,
		},
		{
			name:              "percentage_largest_finite",
			operation:         "percentage",
			operands:          []float64{math.MaxFloat64},
			want:              1.7976931348623157e306,
			relativeTolerance: 1e-14,
		},
	})
}

// Preserve validation precedence when more than one input is invalid.
func TestCalculateErrorPrecedence(t *testing.T) {
	runCalculationCases(t, []calculationCase{
		{
			name:      "arity_before_unknown_operation",
			operation: "unknown",
			operands:  nil,
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "finite_operands_before_unknown_operation",
			operation: "unknown",
			operands:  []float64{math.NaN(), 1},
			wantCode:  "INVALID_OPERANDS",
		},
		{
			name:      "finite_operands_before_division_by_zero",
			operation: "divide",
			operands:  []float64{math.Inf(1), 0},
			wantCode:  "INVALID_OPERANDS",
		},
	})
}
