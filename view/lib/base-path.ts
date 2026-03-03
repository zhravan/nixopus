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

export function getServerConfigBase(): string {
  const origin = process.env.VIEW_APP_ORIGIN || `http://localhost:${process.env.PORT || '7443'}`;
  const base = getBasePath();
  return base ? `${origin}${base}` : origin;
}
