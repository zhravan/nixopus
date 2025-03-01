'use client';
import DashboardPageHeader from '@/components/dashboard-page-header';
import React from 'react';
import RecentActivity from './components/RecentActivity';
import TeamStats from './components/TeamStats';
import useTeamSettings from '../hooks/use-team-settings';
import AddMember from './components/AddMember';
import TeamMembers from './components/TeamMembers';
import { useAppSelector } from '@/redux/hooks';

function Page() {
  const {
    users,
    isAddUserDialogOpen,
    setIsAddUserDialogOpen,
    newUser,
    setNewUser,
    handleAddUser,
    handleRemoveUser,
    getRoleBadgeVariant
  } = useTeamSettings();
  const activeOrganization = useAppSelector((state) => state.user.activeOrganization);

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl">
      <DashboardPageHeader label={activeOrganization?.name} description={activeOrganization?.description} />
      <div className="flex justify-end">
        <AddMember
          isAddUserDialogOpen={isAddUserDialogOpen}
          setIsAddUserDialogOpen={setIsAddUserDialogOpen}
          newUser={newUser}
          setNewUser={setNewUser}
          handleAddUser={handleAddUser}
        />
      </div>
      <TeamMembers
        users={users}
        handleRemoveUser={handleRemoveUser}
        getRoleBadgeVariant={getRoleBadgeVariant}
      />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <TeamStats users={users} />
        <RecentActivity />
      </div>
    </div>
  );
}

export default Page;
