import { describe, expect, it } from 'vitest'
import { normalizeExpression } from './normalizeExpression'

describe('normalizeExpression', () => {
  it('maps every visual operator, including repeated symbols', () => {
    expect(normalizeExpression('√(9)×(−2−3)÷2×4')).toBe('sqrt(9)*(-2-3)/2*4')
  })

  it('leaves ASCII and unrecognized input for backend validation', () => {
    const expression = ' sqrt(2.25)^2-1%+4/2*3 + $é '
    expect(normalizeExpression(expression)).toBe(expression)
  })

  it('does not repair incomplete expressions', () => {
    expect(normalizeExpression('√9−')).toBe('sqrt9-')
  })
})
