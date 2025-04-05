'use client';
import React from 'react';
import ListRepositories from '../components/github-repositories/list-repositories';
import { useAppSelector } from '@/redux/hooks';
import { hasPermission } from '@/lib/permission';

function page() {
  const user = useAppSelector((state) => state.auth.user);
  const activeOrg = useAppSelector((state) => state.user.activeOrganization);
  const canCreate = hasPermission(user, 'deploy', 'create', activeOrg?.id);

  if (!canCreate) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold">Access Denied</h2>
          <p className="text-muted-foreground">
            You don't have permission to create self-host applications
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl">
      <ListRepositories />
    </div>
  );
}

export default page;
