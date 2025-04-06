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
import { defaultLocale, locales } from './lib/i18n/config';

export async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;

  // Handle auth token clearing
  if (pathname === '/logout') {
    const response = NextResponse.redirect(new URL('/login', request.url));
    clearAuthTokens(response);
    return response;
  }

  // Define public paths that don't require authentication
  const publicPaths = [
    '/login',
    '/register',
    '/api/auth',
    '/_next',
    '/static',
    '/favicon.ico',
    '/reset-password',
    '/verify-email'
  ];

  // Check if the path is public
  const isPublicPath = publicPaths.some((path) => pathname.includes(path));
  if (isPublicPath) {
    return NextResponse.next();
  }

  // Handle special flows that don't require auth
  const isGitHubFlow =
    (pathname.includes('/self-host') || pathname.includes('/github-callback')) &&
    (request.nextUrl.searchParams.has('code') ||
      request.nextUrl.searchParams.has('installation_id') ||
      request.nextUrl.searchParams.has('setup_action'));

  const isPasswordResetFlow =
    pathname.includes('/reset-password') && request.nextUrl.searchParams.has('token');

  const isVerificationFlow =
    pathname.includes('/verify-email') && request.nextUrl.searchParams.has('token');

  if (isGitHubFlow || isPasswordResetFlow || isVerificationFlow) {
    return NextResponse.next();
  }

  // Check for auth token
  const token = request.cookies.get('token');
  if (!token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
    '/self-host/:path*',
    '/github-callback/:path*'
  ]
};

export function i18nMiddleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;
  const pathnameIsMissingLocale = locales.every(
    (locale) => !pathname.startsWith(`/${locale}/`) && pathname !== `/${locale}`
  );

  if (pathnameIsMissingLocale) {
    const locale = defaultLocale;
    return NextResponse.redirect(new URL(`/${locale}${pathname}`, request.url));
  }
}
