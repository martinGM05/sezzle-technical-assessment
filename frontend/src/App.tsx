import { useEffect, useState } from 'react'
import { Alert } from './components/ui/Alert/Alert'
import {
  ThemeToggle,
  type Theme,
} from './components/ui/ThemeToggle/ThemeToggle'
import { Calculator } from './features/calculator/components/Calculator/Calculator'
import { operations } from './features/calculator/operations'
import type { CalculatorAction, Operation } from './features/calculator/types'
import { HistoryPanel } from './features/history/components/HistoryPanel/HistoryPanel'
import type { HistoryEntry } from './features/history/types'
import './App.css'

function App() {
  const [theme, setTheme] = useState<Theme>('system')
  const [expression, setExpression] = useState('')
  const [selectedOperation, setSelectedOperation] = useState<Operation | null>(
    null,
  )
  const [integrationNotice, setIntegrationNotice] = useState(false)
  const [historyOpen, setHistoryOpen] = useState(false)
  const [history, setHistory] = useState<readonly HistoryEntry[]>([])

  useEffect(() => {
    document.documentElement.dataset.theme = theme
  }, [theme])

  function handleAction(action: CalculatorAction) {
    setIntegrationNotice(false)
    switch (action.type) {
      case 'input':
        setExpression((current) => current + action.value)
        break
      case 'operation':
        setSelectedOperation(action.operation)
        setExpression((current) => current + operations[action.operation].input)
        break
      case 'clear':
        setExpression('')
        setSelectedOperation(null)
        break
      case 'backspace':
        setExpression((current) => current.slice(0, -1))
        break
      case 'equals':
        setIntegrationNotice(true)
        break
    }
  }

  return (
    <main className="app-shell">
      <header className="app-shell__header">
        <p className="app-shell__brand">
          <strong>Sezzle</strong>
          <span>Technical Assessment</span>
        </p>
        <ThemeToggle value={theme} onChange={setTheme} />
      </header>
      <Calculator
        expression={expression}
        result={null}
        selectedOperation={selectedOperation}
        onAction={handleAction}
      />
      {integrationNotice ? (
        <div className="app-shell__notice">
          <Alert aria-label="Integration notice">
            Calculation will be available when backend integration is added.
          </Alert>
        </div>
      ) : null}
      <div className="app-shell__history">
        <HistoryPanel
          entries={history}
          open={historyOpen}
          onOpenChange={setHistoryOpen}
          onReuse={(entry) => {
            setExpression(entry.expression)
            setSelectedOperation(null)
            setIntegrationNotice(false)
          }}
          onDelete={(id) =>
            setHistory((entries) => entries.filter((entry) => entry.id !== id))
          }
          onClear={() => setHistory([])}
        />
      </div>
      <footer className="app-shell__footer">
        Interface preview · Calculations are not connected yet.
      </footer>
    </main>
  )
}

export default App
