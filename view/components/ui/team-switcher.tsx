'use client';

import * as React from 'react';
import { ChevronsUpDown, GroupIcon, Plus, Trash2, Users } from 'lucide-react';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar
} from '@/components/ui/sidebar';
import { UserOrganization } from '@/redux/types/orgs';
import { useAppSelector } from '@/redux/hooks';
import { DeleteDialog } from './delete-dialog';
import { useDeleteOrganizationMutation } from '@/redux/services/users/userApi';
import useTeamSwitcher from '@/hooks/use-team-switcher';

export function TeamSwitcher({
  teams,
  refetch
}: {
  teams?: UserOrganization[] | [];
  refetch: () => void;
}) {
  const { isMobile } = useSidebar();
  const user = useAppSelector((state) => state.auth.user);
  const isAdmin = React.useMemo(() => user?.type === 'admin', [user]);
  const activeTeam = useAppSelector((state) => state.user.activeOrganization);
  const [deleteOrganization] = useDeleteOrganizationMutation();
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = React.useState(false);
  const { toggleAddTeamModal, handleTeamChange } = useTeamSwitcher();

  if (!teams || teams.length === 0) {
    return null;
  }

  const displayTeam = activeTeam || (teams.length > 0 ? teams[0].organization : null);

  if (!displayTeam) {
    return null;
  }

  const handleDeleteOrganization = async () => {
    if (teams.length <= 1) {
      return;
    }

    try {
      await deleteOrganization(displayTeam.id).unwrap();
      const remainingTeams = teams.filter((team) => team.organization.id !== displayTeam.id);
      handleTeamChange(remainingTeams[0]);
      await refetch();
      setIsDeleteDialogOpen(false);
    } catch (error) {
      console.error('Failed to delete organization:', error);
    }
  };

  return (
    <>
      <DeleteDialog
        title="Delete Organization"
        description={`Are you sure you want to delete the organization "${displayTeam.name}"? This action cannot be undone.`}
        onConfirm={handleDeleteOrganization}
        variant="destructive"
        icon={Trash2}
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
      />
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton
                size="lg"
                className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              >
                <div className="bg-primary text-background flex aspect-square size-8 items-center justify-center rounded-lg">
                  <Users className="size-4 text-background" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{displayTeam.name}</span>
                  <span className="truncate text-xs">{displayTeam.description}</span>
                </div>
                <ChevronsUpDown className="ml-auto" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
              align="start"
              side={isMobile ? 'bottom' : 'right'}
              sideOffset={4}
            >
              <DropdownMenuLabel className="text-muted-foreground text-xs">Teams</DropdownMenuLabel>
              {teams.map((team, index) => (
                <DropdownMenuItem
                  key={team.organization.id}
                  onClick={() => handleTeamChange(team)}
                  className="gap-2 p-2"
                >
                  <div className="flex size-6 items-center justify-center rounded-xs border">
                    <GroupIcon className="size-4 shrink-0" />
                  </div>
                  {team.organization.name}
                  {index < 9 && <DropdownMenuShortcut>⌘{index + 1}</DropdownMenuShortcut>}
                </DropdownMenuItem>
              ))}
              <DropdownMenuSeparator />
              {isAdmin && (
                <>
                  <DropdownMenuItem className="gap-2 p-2" onClick={toggleAddTeamModal}>
                    <div className="bg-background flex size-6 items-center justify-center rounded-md border">
                      <Plus className="size-4" />
                    </div>
                    <div className="text-muted-foreground font-medium">Add team</div>
                  </DropdownMenuItem>
                  {teams.length > 1 && (
                    <DropdownMenuItem
                      className="gap-2 p-2 text-destructive"
                      onClick={() => setIsDeleteDialogOpen(true)}
                    >
                      <div className="bg-background flex size-6 items-center justify-center rounded-md border">
                        <Trash2 className="size-4" />
                      </div>
                      <div className="text-muted-foreground font-medium">Delete team</div>
                    </DropdownMenuItem>
                  )}
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
    </>
  );
}
