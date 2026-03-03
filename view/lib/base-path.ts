/**
 * Base path when the app is served under a subpath (e.g. /view in multi-zone).
 * - Server: BASE_PATH (set in Docker/env)
 * - Client: NEXT_PUBLIC_BASE_PATH (inlined at build time)
 */
export function getBasePath(): string {
  if (typeof window !== 'undefined') {
    return process.env.NEXT_PUBLIC_BASE_PATH || '';
  }
  return process.env.BASE_PATH || '';
}
