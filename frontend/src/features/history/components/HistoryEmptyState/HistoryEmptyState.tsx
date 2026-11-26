import './HistoryEmptyState.css'

export function HistoryEmptyState() {
  return (
    <div className="history-empty-state" data-testid="history-empty-state">
      <span className="history-empty-state__icon" aria-hidden="true">
        ↺
      </span>
      <h3>No calculations yet</h3>
      <p>Your completed calculations will appear here.</p>
    </div>
  )
}
