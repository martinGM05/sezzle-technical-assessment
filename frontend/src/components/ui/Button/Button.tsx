import type { ComponentPropsWithoutRef } from 'react'
import './Button.css'

export type ButtonProps = ComponentPropsWithoutRef<'button'> & {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
}

export function Button({
  variant = 'secondary',
  type = 'button',
  className = '',
  ...props
}: ButtonProps) {
  return (
    <button
      {...props}
      type={type}
      className={`button button--${variant} ${className}`.trim()}
    />
  )
}
