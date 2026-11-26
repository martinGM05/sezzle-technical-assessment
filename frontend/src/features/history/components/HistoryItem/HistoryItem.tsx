import { Button } from '../../../../components/ui/Button/Button'
import { IconButton } from '../../../../components/ui/IconButton/IconButton'
import type { HistoryEntry } from '../../types'
import './HistoryItem.css'

export type HistoryItemProps = {
  entry: HistoryEntry
  onReuse: (entry: HistoryEntry) => void
  onCopy?: (entry: HistoryEntry) => void
  onDelete: (id: string) => void
}

export function HistoryItem({
  entry,
  onReuse,
  onCopy,
  onDelete,
}: HistoryItemProps) {
  const date = new Date(entry.timestamp)
  const validTimestamp = Number.isFinite(date.getTime())

  return (
    <li className="history-item" data-testid={`history-item-${entry.id}`}>
      <div className="history-item__calculation">
        <p className="history-item__expression">{entry.expression}</p>
        <p className="history-item__result">{entry.result}</p>
        {validTimestamp ? (
          <time
            className="history-item__time"
            dateTime={entry.timestamp}
            data-testid="history-timestamp"
          >
            {date.toLocaleString(undefined, {
              dateStyle: 'medium',
              timeStyle: 'short',
            })}
          </time>
        ) : (
          <span className="history-item__time">Time unavailable</span>
        )}
      </div>
      <div
        className="history-item__actions"
        role="group"
        aria-label="Entry actions"
      >
        <Button
          variant="ghost"
          onClick={() => onReuse(entry)}
          aria-label="Reuse expression"
          data-testid="history-reuse"
        >
          Reuse
        </Button>
        <IconButton
          variant="ghost"
          aria-label="Copy result"
          disabled={!onCopy}
          onClick={() => onCopy?.(entry)}
          data-testid="history-copy"
        >
          ⧉
        </IconButton>
        <IconButton
          variant="ghost"
          aria-label="Delete entry"
          onClick={() => onDelete(entry.id)}
          data-testid="history-delete"
        >
          ×
        </IconButton>
      </div>
    </li>
  )
}
