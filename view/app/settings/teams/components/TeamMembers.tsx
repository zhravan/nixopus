import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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
import { useResourcePermissions } from '@/lib/permission';
import EditUserDialog from './EditUserDialog';
import { UserTypes } from '@/redux/types/orgs';
import { DeleteDialog } from '@/components/ui/delete-dialog';
import { useTranslation } from '@/hooks/use-translation';

type EditUser = {
  id: string;
  name: string;
  email: string;
  avatar: string;
  role: 'Owner' | 'Admin' | 'Member' | 'Viewer';
  permissions: string[];
};

interface TeamMembersProps {
  users: EditUser[];
  handleRemoveUser: (userId: string) => void;
  getRoleBadgeVariant: (role: string) => 'destructive' | 'default' | 'secondary' | 'outline';
  onUpdateUser: (userId: string, role: UserTypes) => void;
}

const MAX_VISIBLE_PERMISSIONS = 3;

function TeamMembers({
  users,
  handleRemoveUser,
  getRoleBadgeVariant,
  onUpdateUser
}: TeamMembersProps) {
  const { t } = useTranslation();
  const loggedInUser = useAppSelector((state) => state.auth.user);
  const activeOrganization = useAppSelector((state) => state.user.activeOrganization);
  const { canUpdate: canUpdateUser, canDelete: canDeleteUser } = useResourcePermissions(
    loggedInUser,
    'user',
    activeOrganization?.id
  );
  const [expandedUsers, setExpandedUsers] = useState<Set<string>>(new Set());
  const [editingUser, setEditingUser] = useState<EditUser | null>(null);
  const [userToRemove, setUserToRemove] = useState<EditUser | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

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

  const handleEditUser = (user: any) => {
    setEditingUser({
      ...user,
      permissions: user.permissions
    });
  };

  const handleSaveUser = (userId: string, role: UserTypes) => {
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
            <Badge key={index} variant="outline" className="bg-background">
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

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>{t('settings.teams.members.title')}</CardTitle>
          <CardDescription>{t('settings.teams.members.description')}</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('settings.teams.members.table.headers.user')}</TableHead>
                <TableHead>{t('settings.teams.members.table.headers.role')}</TableHead>
                <TableHead>{t('settings.teams.members.table.headers.permissions')}</TableHead>
                {(canUpdateUser || canDeleteUser) && (
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
                        <div className="text-sm text-muted-foreground">{user.email}</div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={getRoleBadgeVariant(user.role)}>{user.role}</Badge>
                  </TableCell>
                  <TableCell>{renderPermissions(user.permissions, user.id)}</TableCell>
                  {(canUpdateUser || canDeleteUser) && loggedInUser.id !== user.id && (
                    <TableCell className="text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <DotsVerticalIcon className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          {canUpdateUser && (
                            <>
                              <DropdownMenuItem onClick={() => handleEditUser(user)}>
                                <PencilIcon className="h-4 w-4 mr-2" />
                                {t('settings.teams.members.actions.edit')}
                              </DropdownMenuItem>
                            </>
                          )}
                          {canUpdateUser && canDeleteUser && <DropdownMenuSeparator />}
                          {canDeleteUser && (
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
