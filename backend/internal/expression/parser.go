package expression

import "github.com/martinGM05/sezzle-technical-assessment/backend/internal/calculator"

type parser struct {
	tokens   []token
	position int
}

func (p *parser) current() token {
	return p.tokens[p.position]
}

func (p *parser) sum(depth int) (*node, *calculator.Error) {
	left, err := p.product(depth)
	if err != nil {
		return nil, err
	}
	for p.current().kind == '+' || p.current().kind == '-' {
		operator := p.current().kind
		p.position++
		right, err := p.product(depth)
		if err != nil {
			return nil, err
		}
		operation := calculator.OperationAdd
		if operator == '-' {
			operation = calculator.OperationSubtract
		}
		left = &node{operation: operation, children: []*node{left, right}}
	}
	return left, nil
}

func (p *parser) product(depth int) (*node, *calculator.Error) {
	left, err := p.unary(depth)
	if err != nil {
		return nil, err
	}
	for p.current().kind == '*' || p.current().kind == '/' {
		operator := p.current().kind
		p.position++
		right, err := p.unary(depth)
		if err != nil {
			return nil, err
		}
		operation := calculator.OperationMultiply
		if operator == '/' {
			operation = calculator.OperationDivide
		}
		left = &node{operation: operation, children: []*node{left, right}}
	}
	return left, nil
}

func (p *parser) unary(depth int) (*node, *calculator.Error) {
	if depth > MaxNesting {
		return nil, limitError()
	}
	if p.current().kind == '+' || p.current().kind == '-' {
		sign := p.current().kind
		p.position++
		operand, err := p.unary(depth + 1)
		if err != nil {
			return nil, err
		}
		if sign == '+' {
			return operand, nil
		}
		// Multiplication preserves negative zero, unlike subtracting from positive zero.
		return &node{operation: calculator.OperationMultiply, children: []*node{{value: -1}, operand}}, nil
	}
	return p.power(depth)
}

func (p *parser) primary(depth int) (*node, *calculator.Error) {
	current := p.current()
	if current.kind == tokenNumber {
		p.position++
		return &node{value: current.value}, nil
	}
	if current.kind == '(' || current.kind == tokenSqrt {
		if current.kind == tokenSqrt {
			p.position++
			if p.current().kind != '(' {
				return nil, syntaxError()
			}
		}
		p.position++
		inner, err := p.sum(depth + 1)
		if err != nil {
			return nil, err
		}
		if p.current().kind != ')' {
			return nil, syntaxError()
		}
		p.position++
		if current.kind == tokenSqrt {
			return &node{operation: calculator.OperationSqrt, children: []*node{inner}}, nil
		}
		return inner, nil
	}
	return nil, syntaxError()
}

// The right operand accepts unary signs and another power, making ^ right-associative.
func (p *parser) power(depth int) (*node, *calculator.Error) {
	left, err := p.postfix(depth)
	if err != nil {
		return nil, err
	}
	if p.current().kind != '^' {
		return left, nil
	}
	p.position++
	right, err := p.unary(depth + 1)
	if err != nil {
		return nil, err
	}
	return &node{operation: calculator.OperationPower, children: []*node{left, right}}, nil
}

func (p *parser) postfix(depth int) (*node, *calculator.Error) {
	operand, err := p.primary(depth)
	if err != nil {
		return nil, err
	}
	for p.current().kind == '%' {
		p.position++
		operand = &node{operation: calculator.OperationPercentage, children: []*node{operand}}
	}
	return operand, nil
}
