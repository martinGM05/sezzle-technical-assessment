import { Button, type ButtonProps } from '../Button/Button'
import './IconButton.css'

export type IconButtonProps = Omit<ButtonProps, 'aria-label'> & {
  'aria-label': string
}

export function IconButton({
  children,
  className = '',
  ...props
}: IconButtonProps) {
  return (
    <Button {...props} className={`icon-button ${className}`.trim()}>
      <span className="icon-button__icon" aria-hidden="true">
        {children}
      </span>
    </Button>
  )
}
