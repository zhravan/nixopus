import { useAppSelector } from '@/redux/hooks';
import {
  useCreateInviteMutation,
  useGetInvitedOrganizationUsersQuery,
  useRemoveUserFromOrganizationMutation,
  useUpdateOrganizationDetailsMutation,
  useUpdateUserRoleMutation
} from '@/redux/services/users/userApi';
import { UserTypes } from '@/redux/types/orgs';
import { useState, useEffect } from 'react';
import { toast } from 'sonner';
import { useTranslation } from '@/hooks/use-translation';

function useTeamSettings() {
  const { t } = useTranslation();
  const [users, setUsers] = useState<any>([]);
  const [isAddUserDialogOpen, setIsAddUserDialogOpen] = useState(false);
  const [newUser, setNewUser] = useState({ name: '', email: '', role: 'Member' });
  const [isEditTeamDialogOpen, setEditTeamDialogOpen] = useState(false);
  const [teamName, setTeamName] = useState('');
  const [teamDescription, setTeamDescription] = useState('');
  const [createInvite, { isLoading: isCreating }] = useCreateInviteMutation();
  const [removeUserFromOrganization] = useRemoveUserFromOrganizationMutation();
  const [updateUserRole] = useUpdateUserRoleMutation();
  const activeOrganization = useAppSelector((state) => state.user.activeOrganization);
  const {
    data: apiUsers,
    isLoading,
    error,
    refetch: refetchUsers
  } = useGetInvitedOrganizationUsersQuery(activeOrganization?.id as string, {
    skip: !activeOrganization
  });
  const [updateOrganizationDetails, { isLoading: isUpdating, error: updateError }] =
    useUpdateOrganizationDetailsMutation();

  useEffect(() => {
    if (apiUsers) {
      const transformedUsers = apiUsers.map((row: any) => {
        const roleName: string = row.role?.name || row.invite_role;
        const permissions =
          row?.role?.permissions?.map(
            (permission: { resource: string; name: string }) =>
              `${permission.resource.toUpperCase()}:${permission.name}`
          ) || [];
        const isVerified: boolean = Boolean(row?.user?.is_verified || false);
        const acceptedAt: string | null = row.accepted_at || null;
        const expiresAt: string | null = row.expires_at || null;
        let status: string = '-';
        if (!acceptedAt && !isVerified) {
          if (expiresAt && new Date(expiresAt).getTime() < Date.now()) {
            status = 'Expired';
          } else {
            status = 'Pending';
          }
        }

        const id = row?.user?.id || row?.user_id;
        const name = row?.user?.username || row?.invite_name || '';
        const email = row?.user?.email || row?.invite_email || '';
        return {
          id,
          name,
          email,
          role: normalizeRole(roleName),
          permissions,
          status,
          invite: {
            email: row.invite_email || null,
            name: row.invite_name || null,
            role: row.invite_role || null,
            expires_at: row.expires_at || null,
            accepted_at: row.accepted_at || null,
            invited_by: row.invited_by || null
          }
        };
      });

      setUsers(transformedUsers);
    }
  }, [apiUsers]);

  useEffect(() => {
    if (activeOrganization) {
      setTeamName(activeOrganization.name);
      setTeamDescription(activeOrganization.description);
    }
  }, [activeOrganization]);

  const handleAddUser = async () => {
    const newId = crypto.randomUUID();
    const permissions = newUser.role === 'Member' ? ['READ', 'UPDATE'] : ['READ'];
    setUsers([
      ...users,
      {
        id: newId,
        name: newUser.name,
        email: newUser.email,
        role: newUser.role,
        permissions,
        status: 'Pending'
      }
    ]);
    try {
      await createInvite({
        email: newUser.email || '',
        name: newUser.name || '',
        role: newUser.role
      }).unwrap();
      await refetchUsers();
      toast.success(t('settings.teams.messages.userInvited'));
    } catch (error) {
      toast.error(t('settings.teams.messages.userInviteFailed'));
    }
    setNewUser({ name: '', email: '', role: 'Member' });
    setIsAddUserDialogOpen(false);
  };

  const handleRemoveUser = async (userId: string) => {
    try {
      await removeUserFromOrganization({
        user_id: userId,
        organization_id: activeOrganization?.id || ''
      });
      await refetchUsers();
      toast.success(t('settings.teams.messages.userRemoved'));
    } catch (error) {
      toast.error(t('settings.teams.messages.userRemoveFailed'));
    }
  };

  const handleUpdateUser = async (userId: string, role: UserTypes) => {
    try {
      await updateUserRole({
        user_id: userId,
        organization_id: activeOrganization?.id || '',
        role_name: role
      });
      await refetchUsers();
      toast.success(t('settings.teams.messages.userUpdated'));
    } catch (error) {
      toast.error(t('settings.teams.messages.userUpdateFailed'));
    }
  };

  const getRoleBadgeVariant = (role: string) => {
    const r = normalizeRole(role);
    switch (r) {
      case 'Owner':
        return 'default';
      case 'admin':
        return 'destructive';
      case 'member':
        return 'default';
      case 'viewer':
        return 'secondary';
      default:
        return 'outline';
    }
  };

  function normalizeRole(r: string): string {
    return (r || '').toLowerCase();
  }

  const handleUpdateTeam = async () => {
    setEditTeamDialogOpen(false);
    if (teamName.length <= 0 || teamDescription.length <= 0) {
      toast.error(t('settings.teams.messages.requiredFields'));
      setTeamName(activeOrganization?.name || '');
      setTeamDescription(activeOrganization?.description || '');
      return;
    }

    if (
      teamName !== activeOrganization?.name ||
      teamDescription !== activeOrganization?.description
    ) {
      try {
        await updateOrganizationDetails({
          id: activeOrganization?.id || '',
          name: teamName,
          description: teamDescription
        });
        await refetchUsers();
        toast.success(t('settings.teams.messages.teamUpdated'));
      } catch (error) {
        toast.error(t('settings.teams.messages.teamUpdateFailed'));
        setTeamName(activeOrganization?.name || '');
        setTeamDescription(activeOrganization?.description || '');
      }
    }
  };

  return {
    users,
    isLoading,
    error,
    isAddUserDialogOpen,
    setIsAddUserDialogOpen,
    newUser,
    setNewUser,
    handleAddUser,
    handleRemoveUser,
    handleUpdateUser,
    getRoleBadgeVariant,
    handleUpdateTeam,
    setEditTeamDialogOpen,
    setTeamName,
    setTeamDescription,
    isEditTeamDialogOpen,
    teamName,
    teamDescription,
    isUpdating
  };
}

export default useTeamSettings;
