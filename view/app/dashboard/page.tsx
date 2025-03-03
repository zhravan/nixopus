'use client';
import DashboardPageHeader from '@/components/dashboard-page-header';
import { useAppSelector } from '@/redux/hooks';
import { useRouter } from 'next/navigation';
import React from 'react';
import GitHubAppSetup from './components/github-connector/github-app-setup';

function page() {
  const user = useAppSelector((state) => state.auth.user);
  const authenticated = useAppSelector((state) => state.auth.isAuthenticated);
  const router = useRouter();

  if (!user || !authenticated) {
    router.push('/login');
    return null;
  }

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl">
      <DashboardPageHeader label="Dashboard" description={`Welcome, ${user.name}`} />
      <GitHubAppSetup />
    </div>
  )
}

export default page;
