// Package calculator implements stateless arithmetic over finite float64 values.
package calculator

import "math"

const (
	OperationAdd        = "add"
	OperationSubtract   = "subtract"
	OperationMultiply   = "multiply"
	OperationDivide     = "divide"
	OperationPower      = "power"
	OperationSqrt       = "sqrt"
	OperationPercentage = "percentage"
)

const (
	CodeInvalidOperands  = "INVALID_OPERANDS"
	CodeInvalidOperation = "INVALID_OPERATION"
	CodeDivisionByZero   = "DIVISION_BY_ZERO"
	CodeInvalidResult    = "INVALID_RESULT"
)

// Error carries a stable client-facing code and a safe explanation.
type Error struct {
	Code   string
	Detail string
}

func (e *Error) Error() string {
	return e.Detail
}

// Calculate validates arity and finite operands before dispatching the operation.
// Power follows math.Pow, including finite 0^0 = 1; percentage means value / 100.
func Calculate(operation string, operands []float64) (float64, *Error) {
	if err := validateOperands(operation, operands); err != nil {
		return 0, err
	}

	var result float64
	switch operation {
	case OperationAdd:
		result = operands[0] + operands[1]
	case OperationSubtract:
		result = operands[0] - operands[1]
	case OperationMultiply:
		result = operands[0] * operands[1]
	case OperationDivide:
		if operands[1] == 0 {
			return 0, &Error{
				Code:   CodeDivisionByZero,
				Detail: "Cannot divide by zero.",
			}
		}
		result = operands[0] / operands[1]
	case OperationPower:
		result = math.Pow(operands[0], operands[1])
	case OperationSqrt:
		if operands[0] < 0 {
			return 0, &Error{
				Code:   CodeInvalidOperands,
				Detail: "Square root requires a non-negative operand.",
			}
		}
		result = math.Sqrt(operands[0])
	case OperationPercentage:
		result = operands[0] / 100
	default:
		return 0, &Error{
			Code:   CodeInvalidOperation,
			Detail: "Operation is not supported.",
		}
	}

	if !isFinite(result) {
		return 0, &Error{
			Code:   CodeInvalidResult,
			Detail: "Calculation result must be finite.",
		}
	}
	return result, nil
}

func validateOperands(operation string, operands []float64) *Error {
	// Unknown operations retain binary arity validation before dispatch. This
	// preserves the established error precedence for multiply-invalid inputs.
	arity := 2
	detail := "Exactly two operands are required."
	if operation == OperationSqrt || operation == OperationPercentage {
		arity = 1
		detail = "Exactly one operand is required."
	}
	if len(operands) != arity {
		return &Error{Code: CodeInvalidOperands, Detail: detail}
	}
	for _, operand := range operands {
		if !isFinite(operand) {
			return &Error{
				Code:   CodeInvalidOperands,
				Detail: "Operands must be finite.",
			}
		}
	}
	return nil
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
