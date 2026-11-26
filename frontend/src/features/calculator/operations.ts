import type { Operation } from './types'

type OperationDefinition = {
  label: string
  shortLabel: string
  symbol: string
  input: string
}

// Display/input tokens only: parsing and evaluation belong to the backend.
export const operations: Record<Operation, OperationDefinition> = {
  add: { label: 'Addition', shortLabel: 'Add', symbol: '+', input: '+' },
  subtract: {
    label: 'Subtraction',
    shortLabel: 'Sub',
    symbol: '−',
    input: '−',
  },
  multiply: {
    label: 'Multiplication',
    shortLabel: 'Mul',
    symbol: '×',
    input: '×',
  },
  divide: { label: 'Division', shortLabel: 'Div', symbol: '÷', input: '÷' },
  power: {
    label: 'Exponentiation',
    shortLabel: 'Pow',
    symbol: 'xʸ',
    input: '^',
  },
  sqrt: { label: 'Square root', shortLabel: 'Sqrt', symbol: '√', input: '√' },
  percentage: {
    label: 'Percentage',
    shortLabel: 'Perc',
    symbol: '%',
    input: '%',
  },
}

export const operationOrder: Operation[] = [
  'add',
  'subtract',
  'multiply',
  'divide',
  'power',
  'sqrt',
  'percentage',
]
