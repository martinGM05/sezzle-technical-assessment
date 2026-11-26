package expression

import (
	"math"
	"strconv"
	"strings"

	"github.com/martinGM05/sezzle-technical-assessment/backend/internal/calculator"
)

const (
	tokenEnd    byte = 0
	tokenSqrt   byte = 's'
	tokenNumber byte = 'n'
)

type token struct {
	kind  byte
	value float64
}

func tokenize(input string) ([]token, *calculator.Error) {
	var tokens []token
	for position := 0; position < len(input); {
		char := input[position]
		switch {
		case char == ' ' || char == '\t' || char == '\r' || char == '\n':
			position++
		case char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' || char == '^' || char == '%':
			tokens = append(tokens, token{kind: char})
			position++
		case strings.HasPrefix(input[position:], "sqrt"):
			tokens = append(tokens, token{kind: tokenSqrt})
			position += len("sqrt")
		case isDigit(char) || char == '.':
			start := position
			digits := 0
			for position < len(input) && isDigit(input[position]) {
				digits++
				position++
			}
			if position < len(input) && input[position] == '.' {
				position++
				for position < len(input) && isDigit(input[position]) {
					digits++
					position++
				}
			}
			if digits == 0 {
				return nil, syntaxError()
			}
			if position < len(input) && (input[position] == 'e' || input[position] == 'E') {
				position++
				if position < len(input) && (input[position] == '+' || input[position] == '-') {
					position++
				}
				exponentStart := position
				for position < len(input) && isDigit(input[position]) {
					position++
				}
				if position == exponentStart {
					return nil, syntaxError()
				}
			}
			value, err := strconv.ParseFloat(input[start:position], 64)
			if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
				return nil, &calculator.Error{Code: calculator.CodeInvalidResult, Detail: "Number must be finite and within float64 range."}
			}
			tokens = append(tokens, token{kind: tokenNumber, value: value})
		default:
			return nil, syntaxError()
		}
	}
	return append(tokens, token{kind: tokenEnd}), nil
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}
