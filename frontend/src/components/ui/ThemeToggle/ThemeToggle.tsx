import { useId } from 'react'
import './ThemeToggle.css'

export type Theme = 'light' | 'dark' | 'system'
export type ThemeToggleProps = {
  value: Theme
  onChange: (theme: Theme) => void
}

const options: { value: Theme; label: string; symbol: string }[] = [
  { value: 'light', label: 'Light', symbol: '☀' },
  { value: 'dark', label: 'Dark', symbol: '☾' },
  { value: 'system', label: 'System', symbol: '◐' },
]

export function ThemeToggle({ value, onChange }: ThemeToggleProps) {
  const name = useId()
  return (
    <fieldset className="theme-toggle" data-testid="theme-toggle">
      <legend className="theme-toggle__legend">Color theme</legend>
      <div className="theme-toggle__options">
        {options.map((option) => (
          <label className="theme-toggle__option" key={option.value}>
            <input
              className="theme-toggle__input"
              type="radio"
              name={name}
              value={option.value}
              checked={value === option.value}
              onChange={() => onChange(option.value)}
              data-testid={`theme-option-${option.value}`}
            />
            <span className="theme-toggle__label">
              <span aria-hidden="true">{option.symbol}</span>
              {option.label}
            </span>
          </label>
        ))}
      </div>
    </fieldset>
  )
}
