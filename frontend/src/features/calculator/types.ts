export type Operation =
  'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'sqrt' | 'percentage'
export type InputCharacter =
  '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | '.' | '(' | ')'
export type CalculatorAction =
  | { type: 'input'; value: InputCharacter }
  | { type: 'operation'; operation: Operation }
  | { type: 'clear' }
  | { type: 'backspace' }
  | { type: 'equals' }
