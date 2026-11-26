# Expression evaluation

`POST /api/v1/expressions` accepts `{"expression":"200 * 15%"}` and returns
`{"result":30}`. The existing calculations endpoint is unchanged.

## Grammar and precedence

```ebnf
expression = sum ;
sum        = product, { ("+" | "-"), product } ;
product    = unary, { ("*" | "/"), unary } ;
unary      = ("+" | "-"), unary | power ;
power      = postfix, [ "^", unary ] ;
postfix    = primary, { "%" } ;
primary    = number | "(", sum, ")" | "sqrt", "(", sum, ")" ;
number     = (digits, [".", [digits]] | ".", digits),
             [ ("e" | "E"), ["+" | "-"], digits ] ;
digits     = digit, {digit} ;
digit      = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;
```

ASCII spaces, tabs, CR and LF are permitted between tokens, not within numbers or
`sqrt`. Names are case-sensitive. Scientific notation is accepted for float64
round trips. No implicit multiplication, commas, variables, arbitrary function
calls or Unicode operator aliases are accepted.

Binding from strongest to weakest: parentheses/sqrt, postfix percentage, power,
unary signs, multiplication/division, addition/subtraction. Power associates
right; other binary operators associate left. Thus `-2^2 = -4`, `(-2)^2 = 4`,
`2^-3 = 0.125`, and `2^3^2 = 512`. Percent always divides its operand by 100:
`200*15% = 30`, `200+15% = 200.15`, `50%% = 0.005`. It is not contextual financial
percentage arithmetic and is not remainder. `2^3%` means `2^(3/100)`.

## Resource and error policy

The decoded expression is limited to **4096 bytes**, including whitespace.
Only ASCII syntax is accepted, so accepted characters each occupy one byte.
At most **64 nested recursive constructs** are allowed along a grammar path:
parentheses, sqrt calls, unary signs and right operands of powers each consume
one level. Flat left-associative chains do not consume this nesting budget;
input length still bounds their tokens, syntax tree, evaluation stack and work.
No mutable global parser state or user-provided executable code is used.

The complete HTTP body has the existing one-MiB limit. Routing/CORS/method/media/
body validation retains the existing precedence. Missing or null expression is
422 `INVALID_EXPRESSION`; wrong JSON types, unknown properties and malformed or
trailing JSON are 400 `MALFORMED_JSON`. Duplicate JSON names retain the last value.

Expression validation checks length first, then tokenization, then complete
syntax before arithmetic. Blank input is 422 `INVALID_EXPRESSION`; invalid
characters or grammar are 422 `SYNTAX_ERROR`; excess length/nesting is 422
`EXPRESSION_LIMIT`. Out-of-range literals are 422 `INVALID_RESULT`. After parsing,
operands are evaluated left-to-right and calculator errors propagate unchanged:
`DIVISION_BY_ZERO`, `INVALID_OPERANDS` (negative sqrt), `INVALID_RESULT`
(non-real power or overflow). Error details never echo the submitted expression.

All arithmetic delegates to the existing calculator. IEEE 754 float64 rounding,
underflow, signed zero, `0^0 = 1`, and rejection of non-finite values remain its
policy. Intermediate overflow is rejected even if later arithmetic would cancel
it. Results are not display-rounded; this is not a financial calculation engine.

## Behavioral test matrix (defined before implementation)

- Literals: zero, signed values, leading/trailing decimal point, scientific
  notation, whitespace, precision, invalid characters/numbers and blank input.
- Basic arithmetic: all four operations, precedence, left association, decimal
  and negative results, missing operands, trailing tokens, division by both zeros.
- Grouping/signs: nested parentheses, sqrt, unary chains, unmatched parentheses,
  missing function delimiters and arguments, exact and excess nesting limits.
- Advanced arithmetic: right-associative power, negative exponents, unary/power
  precedence, percentage binding and repetition, all user-supplied examples.
- Numerical/resource boundaries: finite maxima, overflow, invalid real domains,
  non-finite literal spellings, length/nesting boundaries and long flat chains.
- HTTP: success for every operation; complete Problem Details and OpenAPI schema
  assertions; missing/null/type/JSON failures; media, methods, body limit, CORS,
  HEAD, concurrent requests, and preservation of calculations and health.
