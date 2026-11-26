export function normalizeExpression(expression: string): string {
  return expression
    .replaceAll('−', '-')
    .replaceAll('×', '*')
    .replaceAll('÷', '/')
    .replaceAll('√', 'sqrt')
}
