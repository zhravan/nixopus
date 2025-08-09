import React, { useState } from 'react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import { TrashIcon, ChevronDownIcon, ChevronUpIcon, PencilIcon } from 'lucide-react';
import { DotsVerticalIcon } from '@radix-ui/react-icons';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useAppSelector } from '@/redux/hooks';
import EditUserDialog from './EditUserDialog';
import { OrganizationUsers, UserTypes } from '@/redux/types/orgs';
import { DeleteDialog } from '@/components/ui/delete-dialog';
import { useTranslation } from '@/hooks/use-translation';
import { User } from '@/redux/types/user';
import { ResourceGuard } from '@/components/rbac/PermissionGuard';
import { TypographySmall, TypographyMuted } from '@/components/ui/typography';

type EditUser = {
  id: string;
  name: string;
  email: string;
  avatar: string;
  role: 'Owner' | 'Admin' | 'Member' | 'Viewer';
  permissions: string[];
};

type RoleType = 'owner' | 'admin' | 'member' | 'viewer';

const roleHierarchy: Record<RoleType, number> = {
  owner: 4,
  admin: 3,
  member: 2,
  viewer: 1
};

interface TeamMembersProps {
  users: EditUser[];
  handleRemoveUser: (userId: string) => void;
  getRoleBadgeVariant: (role: string) => 'default' | 'secondary' | 'destructive' | 'outline';
  onUpdateUser: (userId: string, role: UserTypes) => Promise<void>;
}

const MAX_VISIBLE_PERMISSIONS = 3;

function TeamMembers({
  users,
  handleRemoveUser,
  getRoleBadgeVariant,
  onUpdateUser
}: TeamMembersProps) {
  const { t } = useTranslation();
  const loggedInUser = useAppSelector((state) => state.auth.user) as User;
  const activeOrganization = useAppSelector((state) => state.user.activeOrganization);
  const [expandedUsers, setExpandedUsers] = useState<Set<string>>(new Set());
  const [editingUser, setEditingUser] = useState<EditUser | null>(null);
  const [userToRemove, setUserToRemove] = useState<EditUser | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  const getCurrentUserRole = (): RoleType | null => {
    if (!loggedInUser || !activeOrganization) return null;
    const orgUser = loggedInUser.organization_users.find(
      (ou: OrganizationUsers) => ou.organization_id === activeOrganization.id
    );
    return (orgUser?.role?.name?.toLowerCase() as RoleType) || null;
  };

  const canModifyUser = (targetUser: EditUser) => {
    if (!loggedInUser || !targetUser || !activeOrganization) {
      return false;
    }

    const currentUserRole = getCurrentUserRole();
    const targetUserRole = targetUser.role?.toLowerCase() as RoleType;

    if (!currentUserRole || !targetUserRole) {
      return false;
    }

    const canModify = roleHierarchy[currentUserRole] >= roleHierarchy[targetUserRole];
    return canModify;
  };

  const toggleUserPermissions = (userId: string) => {
    setExpandedUsers((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(userId)) {
        newSet.delete(userId);
      } else {
        newSet.add(userId);
      }
      return newSet;
    });
  };

  const handleEditUser = (user: EditUser) => {
    if (!canModifyUser(user)) {
      return;
    }
    setEditingUser({
      ...user,
      permissions: user.permissions
    });
  };

  const handleSaveUser = (userId: string, role: UserTypes) => {
    const user = users.find((u) => u.id === userId);
    if (!user || !canModifyUser(user)) {
      return;
    }
    onUpdateUser(userId, role);
    setEditingUser(null);
  };

  const renderPermissions = (permissions: string[], userId: string) => {
    const isExpanded = expandedUsers.has(userId);
    const visiblePermissions = isExpanded
      ? permissions
      : permissions.slice(0, MAX_VISIBLE_PERMISSIONS);
    const hasMore = permissions.length > MAX_VISIBLE_PERMISSIONS;

    return (
      <div className="flex items-center gap-2">
        <div className="flex flex-wrap gap-1.5">
          {visiblePermissions.map((permission, index) => (
            <Badge key={index} variant="outline" className="bg-primary/10 text-primary rounded-full">
              {permission}
            </Badge>
          ))}
        </div>
        {hasMore && (
          <Button
            variant="ghost"
            size="sm"
            className="h-6 px-2 text-xs font-medium text-primary hover:text-primary/80"
            onClick={() => toggleUserPermissions(userId)}
          >
            {isExpanded ? (
              <>
                {t('settings.teams.members.actions.showLess')}{' '}
                <ChevronUpIcon className="ml-1 h-3 w-3" />
              </>
            ) : (
              <>
                {t('settings.teams.members.actions.showMore')}{' '}
                <ChevronDownIcon className="ml-1 h-3 w-3" />
              </>
            )}
          </Button>
        )}
      </div>
    );
  };

  const hasEditableActions = (user: EditUser) => {
    if (loggedInUser.id === user.id) {
      return false;
    }
    return canModifyUser(user);
  };

  return (
    <>
      <Card>
        <CardHeader>
          <TypographySmall>{t('settings.teams.members.title')}</TypographySmall>
          <TypographyMuted>{t('settings.teams.members.description')}</TypographyMuted>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('settings.teams.members.table.headers.user')}</TableHead>
                <TableHead>{t('settings.teams.members.table.headers.role')}</TableHead>
                <TableHead>{t('settings.teams.members.table.headers.permissions')}</TableHead>
                {users.some((user) => hasEditableActions(user)) && (
                  <TableHead className="text-right">
                    {t('settings.teams.members.table.headers.actions')}
                  </TableHead>
                )}
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center space-x-3">
                      <Avatar className="w-12 h-12 shadow-md">
                        {user.avatar ? (
                          <AvatarImage src={user.avatar} alt="Profile avatar" />
                        ) : (
                          <AvatarFallback className="bg-secondary text-foreground text-xl font-medium">
                            {user.name.slice(0, 2).toUpperCase()}
                          </AvatarFallback>
                        )}
                      </Avatar>
                      <div>
                        <div className="font-medium">{user.name}</div>
                        <TypographyMuted>{user.email}</TypographyMuted>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={getRoleBadgeVariant(user.role)}>{user.role}</Badge>
                  </TableCell>
                  <TableCell>{renderPermissions(user.permissions, user.id)}</TableCell>
                  {hasEditableActions(user) && (
                    <TableCell className="text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <DotsVerticalIcon className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <ResourceGuard resource="user" action="update">
                            {canModifyUser(user) && (
                              <DropdownMenuItem onClick={() => handleEditUser(user)}>
                                <PencilIcon className="h-4 w-4 mr-2" />
                                {t('settings.teams.members.actions.edit')}
                              </DropdownMenuItem>
                            )}
                          </ResourceGuard>
                          <ResourceGuard resource="user" action="update">
                            <ResourceGuard resource="user" action="delete">
                              {canModifyUser(user) && <DropdownMenuSeparator />}
                            </ResourceGuard>
                          </ResourceGuard>
                          <ResourceGuard resource="user" action="delete">
                            {canModifyUser(user) && (
                              <DropdownMenuItem
                                className="text-destructive focus:text-destructive"
                                onClick={() => {
                                  setUserToRemove(user);
                                  setIsDeleteDialogOpen(true);
                                }}
                              >
                                <TrashIcon className="h-4 w-4 mr-2" />
                                {t('settings.teams.members.actions.remove')}
                              </DropdownMenuItem>
                            )}
                          </ResourceGuard>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {editingUser && (
        <EditUserDialog
          isOpen={!!editingUser}
          onClose={() => setEditingUser(null)}
          user={editingUser}
          onSave={handleSaveUser}
        />
      )}

      {userToRemove && (
        <DeleteDialog
          title={t('settings.teams.members.deleteDialog.title').replace(
            '{name}',
            userToRemove.name
          )}
          description={t('settings.teams.members.deleteDialog.description').replace(
            '{name}',
            userToRemove.name
          )}
          onConfirm={() => {
            handleRemoveUser(userToRemove.id);
            setUserToRemove(null);
            setIsDeleteDialogOpen(false);
          }}
          confirmText={t('settings.teams.members.deleteDialog.confirm')}
          variant="destructive"
          icon={TrashIcon}
          open={isDeleteDialogOpen}
          onOpenChange={setIsDeleteDialogOpen}
        />
      )}
    </>
  );
}

export default TeamMembers;
