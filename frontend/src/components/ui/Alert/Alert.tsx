import type { ComponentPropsWithoutRef } from 'react'
import './Alert.css'

export type AlertProps = Omit<ComponentPropsWithoutRef<'div'>, 'role'> & {
  variant?: 'info' | 'success' | 'warning' | 'error'
  role?: 'status' | 'alert'
}

export function Alert({
  variant = 'info',
  role = variant === 'error' ? 'alert' : 'status',
  className = '',
  ...props
}: AlertProps) {
  return (
    <div
      {...props}
      role={role}
      className={`alert alert--${variant} ${className}`.trim()}
    />
  )
}
