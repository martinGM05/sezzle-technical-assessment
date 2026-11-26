import { Button } from '../../../../components/ui/Button/Button'
import { operationOrder, operations } from '../../operations'
import type { Operation } from '../../types'
import './OperationSelector.css'

export type OperationSelectorProps = {
  selectedOperation: Operation | null
  onSelect: (operation: Operation) => void
  disabled?: boolean
}

export function OperationSelector({
  selectedOperation,
  onSelect,
  disabled = false,
}: OperationSelectorProps) {
  return (
    <section
      className="operation-selector"
      aria-label="Operations"
      data-testid="operation-selector"
    >
      <h2 className="operation-selector__heading">Operations</h2>
      <div className="operation-selector__options">
        {operationOrder.map((operation) => {
          const definition = operations[operation]
          return (
            <Button
              key={operation}
              className="operation-selector__option"
              aria-label={definition.label}
              aria-pressed={selectedOperation === operation}
              disabled={disabled}
              onClick={() => onSelect(operation)}
              data-testid={`operation-${operation}`}
            >
              <span className="operation-selector__symbol" aria-hidden="true">
                {definition.symbol}
              </span>
              <span className="operation-selector__short" aria-hidden="true">
                {definition.shortLabel}
              </span>
              <span className="operation-selector__name" aria-hidden="true">
                {definition.label}
              </span>
            </Button>
          )
        })}
      </div>
    </section>
  )
}
