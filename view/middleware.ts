import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import {
  getRefreshToken,
  getToken,
  isTokenExpired,
  refreshAccessToken,
  setAuthTokens,
  clearAuthTokens
} from './lib/auth';

export async function middleware(request: NextRequest) {
  console.log(`Processing request for path: ${request.nextUrl.pathname}`);
  const publicPaths = ['/login', '/register', '/api/auth', '/_next', '/static', '/favicon.ico'];

  if (publicPaths.some((path) => request.nextUrl.pathname.startsWith(path))) {
    console.log('Public path accessed, skipping auth check');
    return NextResponse.next();
  }

  const isGitHubFlow =
    (request.nextUrl.pathname.startsWith('/self-host') ||
      request.nextUrl.pathname.startsWith('/github-callback')) &&
    (request.nextUrl.searchParams.has('code') ||
      request.nextUrl.searchParams.has('installation_id') ||
      request.nextUrl.searchParams.has('setup_action'));

  const isPasswordResetFlow =
    request.nextUrl.pathname.startsWith('/reset-password') &&
    request.nextUrl.searchParams.has('token');

  const isVerificationFlow =
    request.nextUrl.pathname.startsWith('/verify-email') &&
    request.nextUrl.searchParams.has('token');

  if (isGitHubFlow || isPasswordResetFlow || isVerificationFlow) {
    console.log('Special flow detected, skipping auth check');
    return NextResponse.next();
  }

  const token = getToken(request);
  const refreshToken = getRefreshToken(request);

  if (!token) {
    console.log('No token found, redirecting to login');
    return NextResponse.redirect(new URL('/login', request.url));
  }

  if (isTokenExpired(token) && refreshToken) {
    console.log('Token expired, attempting refresh');
    try {
      const newTokens = await refreshAccessToken(refreshToken);
      console.log('Token refresh successful');
      const response = NextResponse.next();
      setAuthTokens(newTokens, response);
      return response;
    } catch (error) {
      console.error('Token refresh failed:', error);

      const response = NextResponse.redirect(new URL('/login?error=session_expired', request.url));
      clearAuthTokens(response);
      console.log('Cleared auth cookies due to refresh failure');

      return response;
    }
  }

  if (isTokenExpired(token)) {
    console.log('Token expired without refresh token, redirecting to login');
    const response = NextResponse.redirect(new URL('/login?error=session_expired', request.url));
    response.cookies.delete('token');
    response.cookies.delete('refreshToken');
    return response;
  }

  console.log('Token valid, proceeding with request');
  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api/v1/auth|_next/static|_next/image|favicon.ico).*)',
    '/self-host/:path*',
    '/github-callback/:path*'
  ]
};
