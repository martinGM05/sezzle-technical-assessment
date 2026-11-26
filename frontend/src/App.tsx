import { useState } from 'react'
import { Alert } from './components/ui/Alert/Alert'
import { ThemeToggle } from './components/ui/ThemeToggle/ThemeToggle'
import { Calculator } from './features/calculator/components/Calculator/Calculator'
import { useCalculator } from './features/calculator/hooks/useCalculator'
import { useCalculatorKeyboard } from './features/calculator/hooks/useCalculatorKeyboard'
import { HistoryPanel } from './features/history/components/HistoryPanel/HistoryPanel'
import { useClipboard } from './features/history/hooks/useClipboard'
import { useHistory } from './features/history/hooks/useHistory'
import { useTheme } from './features/theme/hooks/useTheme'
import './App.css'

function App() {
  const theme = useTheme()
  const history = useHistory()
  const calculator = useCalculator(history.add)
  const clipboard = useClipboard()
  const [historyOpen, setHistoryOpen] = useState(false)
  useCalculatorKeyboard(calculator.handleAction)
  const storageWarning = [history.storageError, theme.storageError]
    .filter(Boolean)
    .join(' ')

  return (
    <main className="app-shell">
      <header className="app-shell__header">
        <p className="app-shell__brand">
          <strong>Sezzle</strong>
          <span>Technical Assessment</span>
        </p>
        <ThemeToggle value={theme.theme} onChange={theme.setTheme} />
      </header>
      <Calculator
        expression={calculator.expression}
        result={calculator.result}
        selectedOperation={calculator.selectedOperation}
        loading={calculator.loading}
        error={calculator.error}
        onAction={calculator.handleAction}
      />
      {storageWarning ? (
        <div className="app-shell__notice">
          <Alert variant="warning" aria-label="Storage status">
            {storageWarning}
          </Alert>
        </div>
      ) : null}
      {clipboard.message || clipboard.error ? (
        <div className="app-shell__notice">
          <Alert
            variant={clipboard.error ? 'error' : 'info'}
            aria-label="Clipboard status"
          >
            {clipboard.error ?? clipboard.message}
          </Alert>
        </div>
      ) : null}
      <div className="app-shell__history">
        <HistoryPanel
          entries={history.entries}
          open={historyOpen}
          onOpenChange={setHistoryOpen}
          onReuse={(entry) => calculator.reuse(entry.expression)}
          onCopy={(entry) => {
            void clipboard.copy(entry.result)
          }}
          onDelete={history.remove}
          onClear={history.clear}
        />
      </div>
      <footer className="app-shell__footer">
        Evaluated by Go · IEEE 754 arithmetic
      </footer>
    </main>
  )
}

export default App
