import {
  CalculatorKey,
  type CalculatorKeyProps,
} from '../CalculatorKey/CalculatorKey'
import type { CalculatorAction } from '../../types'
import './CalculatorKeypad.css'

export type CalculatorKeypadProps = {
  onAction: (action: CalculatorAction) => void
  disabled?: boolean
}

type KeyDefinition = Pick<
  CalculatorKeyProps,
  'label' | 'children' | 'action' | 'variant'
> & { id: string }

const keys: KeyDefinition[] = [
  {
    id: 'clear',
    label: 'Clear',
    children: 'C',
    action: { type: 'clear' },
    variant: 'danger',
  },
  {
    id: 'backspace',
    label: 'Backspace',
    children: '⌫',
    action: { type: 'backspace' },
  },
  {
    id: 'open-parenthesis',
    label: 'Open parenthesis',
    children: '(',
    action: { type: 'input', value: '(' },
  },
  {
    id: 'close-parenthesis',
    label: 'Close parenthesis',
    children: ')',
    action: { type: 'input', value: ')' },
  },
  {
    id: '7',
    label: 'Digit 7',
    children: '7',
    action: { type: 'input', value: '7' },
  },
  {
    id: '8',
    label: 'Digit 8',
    children: '8',
    action: { type: 'input', value: '8' },
  },
  {
    id: '9',
    label: 'Digit 9',
    children: '9',
    action: { type: 'input', value: '9' },
  },
  {
    id: 'divide',
    label: 'Divide',
    children: '÷',
    action: { type: 'operation', operation: 'divide' },
  },
  {
    id: '4',
    label: 'Digit 4',
    children: '4',
    action: { type: 'input', value: '4' },
  },
  {
    id: '5',
    label: 'Digit 5',
    children: '5',
    action: { type: 'input', value: '5' },
  },
  {
    id: '6',
    label: 'Digit 6',
    children: '6',
    action: { type: 'input', value: '6' },
  },
  {
    id: 'multiply',
    label: 'Multiply',
    children: '×',
    action: { type: 'operation', operation: 'multiply' },
  },
  {
    id: '1',
    label: 'Digit 1',
    children: '1',
    action: { type: 'input', value: '1' },
  },
  {
    id: '2',
    label: 'Digit 2',
    children: '2',
    action: { type: 'input', value: '2' },
  },
  {
    id: '3',
    label: 'Digit 3',
    children: '3',
    action: { type: 'input', value: '3' },
  },
  {
    id: 'subtract',
    label: 'Subtract',
    children: '−',
    action: { type: 'operation', operation: 'subtract' },
  },
  {
    id: '0',
    label: 'Digit 0',
    children: '0',
    action: { type: 'input', value: '0' },
  },
  {
    id: 'decimal',
    label: 'Decimal point',
    children: '.',
    action: { type: 'input', value: '.' },
  },
  {
    id: 'equals',
    label: 'Equals',
    children: '=',
    action: { type: 'equals' },
    variant: 'primary',
  },
  {
    id: 'add',
    label: 'Add',
    children: '+',
    action: { type: 'operation', operation: 'add' },
  },
]

export function CalculatorKeypad({
  onAction,
  disabled = false,
}: CalculatorKeypadProps) {
  return (
    <div
      className="calculator-keypad"
      role="group"
      aria-label="Calculator keypad"
      data-testid="calculator-keypad"
    >
      {keys.map(({ id, ...key }) => (
        <CalculatorKey
          key={id}
          {...key}
          testId={`calculator-key-${id}`}
          onAction={onAction}
          disabled={disabled}
        />
      ))}
    </div>
  )
}
