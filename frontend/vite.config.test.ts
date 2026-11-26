import { afterEach, describe, expect, it, vi } from 'vitest'

afterEach(() => {
  vi.unstubAllEnvs()
  vi.resetModules()
})

describe('local API proxy configuration', () => {
  it('routes relative /api requests to the local Go service by default', async () => {
    vi.stubEnv('API_PROXY_TARGET', undefined)
    vi.resetModules()
    const { default: config } = await import('./vite.config.ts')
    expect(config.server?.proxy?.['/api']).toEqual({
      target: 'http://127.0.0.1:8080',
    })
  })
  it('accepts an explicit local target without changing client URLs', async () => {
    vi.stubEnv('API_PROXY_TARGET', 'http://127.0.0.1:9080')
    vi.resetModules()
    const { default: config } = await import('./vite.config.ts')
    expect(config.server?.proxy?.['/api']).toEqual({
      target: 'http://127.0.0.1:9080',
    })
  })
})
