// Package expression evaluates bounded arithmetic expressions without shared state.
package expression

import "github.com/martinGM05/sezzle-technical-assessment/backend/internal/calculator"

const (
	MaxLength             = 4096
	MaxNesting            = 64
	CodeExpressionLimit   = "EXPRESSION_LIMIT"
	CodeInvalidExpression = "INVALID_EXPRESSION"
	CodeSyntaxError       = "SYNTAX_ERROR"
)

// Evaluate validates the complete syntax before performing any arithmetic.
// All arithmetic and numerical error semantics belong to calculator.Calculate.
func Evaluate(input string) (float64, *calculator.Error) {
	if len(input) > MaxLength {
		return 0, limitError()
	}
	tokens, err := tokenize(input)
	if err != nil {
		return 0, err
	}
	if tokens[0].kind == tokenEnd {
		return 0, &calculator.Error{Code: CodeInvalidExpression, Detail: "Expression must not be empty."}
	}
	parser := parser{tokens: tokens}
	root, err := parser.sum(0)
	if err != nil {
		return 0, err
	}
	if parser.current().kind != tokenEnd {
		return 0, syntaxError()
	}
	return root.evaluate()
}

func syntaxError() *calculator.Error {
	return &calculator.Error{Code: CodeSyntaxError, Detail: "Expression syntax is invalid."}
}

type node struct {
	value     float64
	operation string
	children  []*node
}

func (n *node) evaluate() (float64, *calculator.Error) {
	if n.operation == "" {
		return n.value, nil
	}
	operands := make([]float64, len(n.children))
	for i, child := range n.children {
		value, err := child.evaluate()
		if err != nil {
			return 0, err
		}
		operands[i] = value
	}
	return calculator.Calculate(n.operation, operands)
}

func limitError() *calculator.Error {
	return &calculator.Error{Code: CodeExpressionLimit, Detail: "Expression exceeds the length or nesting limit."}
}
