import { Button } from '../../../../components/ui/Button/Button'
import './HistoryToggle.css'

export type HistoryToggleProps = {
  open: boolean
  panelId: string
  onOpenChange: (open: boolean) => void
}

export function HistoryToggle({
  open,
  panelId,
  onOpenChange,
}: HistoryToggleProps) {
  return (
    <Button
      variant="ghost"
      className="history-toggle"
      aria-expanded={open}
      aria-controls={panelId}
      onClick={() => onOpenChange(!open)}
      data-testid="history-toggle"
    >
      <span>Calculation history</span>
      <span aria-hidden="true">{open ? '−' : '+'}</span>
    </Button>
  )
}
