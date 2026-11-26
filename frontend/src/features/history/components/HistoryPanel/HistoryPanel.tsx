import { useEffect, useId, useRef, useState } from 'react'
import { Alert } from '../../../../components/ui/Alert/Alert'
import { Button } from '../../../../components/ui/Button/Button'
import type { HistoryEntry } from '../../types'
import { HistoryEmptyState } from '../HistoryEmptyState/HistoryEmptyState'
import { HistoryItem, type HistoryItemProps } from '../HistoryItem/HistoryItem'
import { HistoryToggle } from '../HistoryToggle/HistoryToggle'
import './HistoryPanel.css'

export type HistoryPanelProps = {
  entries: readonly HistoryEntry[]
  open: boolean
  onOpenChange: (open: boolean) => void
  onClear: () => void
} & Pick<HistoryItemProps, 'onReuse' | 'onCopy' | 'onDelete'>

const PAGE_SIZE = 20

type HistoryContentProps = Omit<HistoryPanelProps, 'open' | 'onOpenChange'>

function HistoryContent({
  entries,
  onClear,
  onReuse,
  onCopy,
  onDelete,
}: HistoryContentProps) {
  const [confirming, setConfirming] = useState(false)
  const [visibleCount, setVisibleCount] = useState(PAGE_SIZE)
  const actionsRef = useRef<HTMLDivElement>(null)
  const confirmationRef = useRef<HTMLDivElement>(null)
  const confirmationId = useId()
  const shownCount = Math.min(visibleCount, entries.length)

  useEffect(() => {
    if (confirming) confirmationRef.current?.querySelector('button')?.focus()
  }, [confirming])

  function closeConfirmation() {
    setConfirming(false)
    actionsRef.current?.querySelector('button')?.focus()
  }

  return (
    <>
      <div className="history-panel__toolbar" ref={actionsRef}>
        <p className="history-panel__count">
          {entries.length} {entries.length === 1 ? 'entry' : 'entries'}
        </p>
        <Button
          variant="ghost"
          disabled={entries.length === 0}
          onClick={() => setConfirming(true)}
          data-testid="history-clear"
        >
          Clear history
        </Button>
      </div>
      {confirming ? (
        <div
          className="history-panel__confirmation"
          ref={confirmationRef}
          role="group"
          aria-labelledby={confirmationId}
          onKeyDown={(event) => {
            if (event.key === 'Escape') {
              event.preventDefault()
              closeConfirmation()
            }
          }}
        >
          <Alert variant="warning">
            <span id={confirmationId}>Clear all history?</span> This removes
            every entry.
          </Alert>
          <div className="history-panel__confirmation-actions">
            <Button onClick={closeConfirmation}>Cancel</Button>
            <Button
              variant="danger"
              onClick={() => {
                closeConfirmation()
                onClear()
              }}
            >
              Clear all entries
            </Button>
          </div>
        </div>
      ) : null}
      {entries.length === 0 ? (
        <HistoryEmptyState />
      ) : (
        <>
          <ol className="history-panel__list" aria-label="Calculation entries">
            {entries.slice(0, shownCount).map((entry) => (
              <HistoryItem
                key={entry.id}
                entry={entry}
                onReuse={onReuse}
                onCopy={onCopy}
                onDelete={onDelete}
              />
            ))}
          </ol>
          {entries.length > PAGE_SIZE ? (
            <p className="history-panel__count" role="status">
              Showing {shownCount} of {entries.length} entries
            </p>
          ) : null}
          {shownCount < entries.length ? (
            <Button
              className="history-panel__more"
              onClick={() => setVisibleCount((count) => count + PAGE_SIZE)}
            >
              Show more history
            </Button>
          ) : null}
        </>
      )}
    </>
  )
}

export function HistoryPanel({
  open,
  onOpenChange,
  ...contentProps
}: HistoryPanelProps) {
  const id = useId()
  const toggleRef = useRef<HTMLHeadingElement>(null)
  const contentRef = useRef<HTMLDivElement>(null)
  const focusInContent = useRef(false)
  const empty = contentProps.entries.length === 0

  useEffect(() => {
    if ((!open || empty) && focusInContent.current) {
      toggleRef.current?.querySelector('button')?.focus()
      focusInContent.current = false
    }
  }, [open, empty])

  return (
    <section
      className="history-panel"
      aria-label="Calculation history"
      data-testid="history-panel"
      onFocusCapture={(event) => {
        focusInContent.current =
          contentRef.current?.contains(event.target) ?? false
      }}
      onBlurCapture={(event) => {
        if (
          event.relatedTarget &&
          !event.currentTarget.contains(event.relatedTarget)
        )
          focusInContent.current = false
      }}
    >
      <h2 ref={toggleRef}>
        <HistoryToggle open={open} onOpenChange={onOpenChange} panelId={id} />
      </h2>
      <div
        id={id}
        ref={contentRef}
        className="history-panel__content"
        hidden={!open}
        data-testid="history-content"
      >
        {/* Reset transient confirmation and disclosure state on collapse or empty history. */}
        {open ? (
          <HistoryContent
            key={empty ? 'empty' : 'populated'}
            {...contentProps}
          />
        ) : null}
      </div>
    </section>
  )
}
