import type { ReactNode } from 'react'
import {
  Button,
  type ButtonProps,
} from '../../../../components/ui/Button/Button'
import type { CalculatorAction } from '../../types'
import './CalculatorKey.css'

export type CalculatorKeyProps = {
  label: string
  children: ReactNode
  action: CalculatorAction
  onAction: (action: CalculatorAction) => void
  disabled?: boolean
  variant?: ButtonProps['variant']
  testId?: string
}

export function CalculatorKey({
  label,
  children,
  action,
  onAction,
  disabled = false,
  variant = 'secondary',
  testId,
}: CalculatorKeyProps) {
  return (
    <Button
      className={`calculator-key${action.type === 'operation' ? ' calculator-key--operation' : ''}`}
      aria-label={label}
      data-testid={testId}
      disabled={disabled}
      variant={variant}
      onClick={() => onAction(action)}
    >
      {children}
    </Button>
  )
}
