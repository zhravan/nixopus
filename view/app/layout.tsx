'use client';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import { ThemeProvider } from '@/components/theme-provider';
import { Provider } from 'react-redux';
import { PersistGate } from 'redux-persist/integration/react';
import { store, persistor } from '@/redux/store';
import { Toaster } from '@/components/ui/sonner';
import { useAppSelector } from '@/redux/hooks';
import DashboardLayout from '@/components/dashboard-layout';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin']
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin']
});

export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <Layout>{children}</Layout>;
}

const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        suppressHydrationWarning
      >
        <Provider store={store}>
          <PersistGate loading={null} persistor={persistor}>
            <ChildrenWrapper>{children}</ChildrenWrapper>
          </PersistGate>
        </Provider>
      </body>
    </html>
  );
};

const ChildrenWrapper = ({ children }: { children: React.ReactNode }) => {
  const authenticated = useAppSelector((state) => state.auth.isAuthenticated);

  return (
    <>
      <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
        {authenticated ? <DashboardLayout>{children}</DashboardLayout> : children}
      </ThemeProvider>
      <Toaster />
    </>
  );
};
