import { useAppSelector } from '@/redux/hooks';
import {
  useCreateUserMutation,
  useGetOrganizationUsersQuery,
  useRemoveUserFromOrganizationMutation,
  useUpdateOrganizationDetailsMutation,
  useUpdateUserRoleMutation
} from '@/redux/services/users/userApi';
import { UserTypes } from '@/redux/types/orgs';
import { useState, useEffect } from 'react';
import { toast } from 'sonner';

function useTeamSettings() {
  const [users, setUsers] = useState<any>([]);
  const [isAddUserDialogOpen, setIsAddUserDialogOpen] = useState(false);
  const [newUser, setNewUser] = useState({ name: '', email: '', role: 'Member', password: '' });
  const [isEditTeamDialogOpen, setEditTeamDialogOpen] = useState(false);
  const [teamName, setTeamName] = useState('');
  const [teamDescription, setTeamDescription] = useState('');
  const [createUser, { isLoading: isCreating }] = useCreateUserMutation();
  const [removeUserFromOrganization] = useRemoveUserFromOrganizationMutation();
  const [updateUserRole] = useUpdateUserRoleMutation();
  const activeOrganization = useAppSelector((state) => state.user.activeOrganization);
  const {
    data: apiUsers,
    isLoading,
    error
  } = useGetOrganizationUsersQuery(activeOrganization?.id, {
    skip: !activeOrganization
  });
  const [updateOrganizationDetails, { isLoading: isUpdating, error: updateError }] =
    useUpdateOrganizationDetailsMutation();

  useEffect(() => {
    if (apiUsers) {
      const transformedUsers = apiUsers.map((user) => {
        const roleName = user.role?.name || 'Unknown';
        const permissions =
          user.role?.permissions?.map(
            (permission) => `${permission.resource.toUpperCase()}:${permission.name}`
          ) || [];
        return {
          id: user.user.id,
          name: user.user?.username || 'Unknown User',
          email: user.user?.email || '',
          role: roleName,
          permissions
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
    const tempUser = {
      username: newUser.name || '',
      email: newUser.email || '',
      password: newUser.password || '',
      organization: activeOrganization?.id || '',
      type: newUser.role.toLowerCase() as UserTypes
    };

    setUsers([...users, { id: newId, ...tempUser, name: newUser.name }]);
    try {
      const user = await createUser(tempUser as any);
      toast.success('User added successfully');
    } catch (error) {
      toast.error('Failed to add user');
    }
    setNewUser({ name: '', email: '', role: 'Member', password: '' });
    setIsAddUserDialogOpen(false);
  };

  const handleRemoveUser = async (userId: string) => {
    try {
      await removeUserFromOrganization({
        user_id: userId,
        organization_id: activeOrganization?.id || ''
      });
      setUsers(users.filter((user: any) => user.id !== userId));
      toast.success('User removed successfully');
    } catch (error) {
      toast.error('Failed to remove user');
    }
  };

  const handleUpdateUser = async (userId: string, role: UserTypes) => {
    try {
      await updateUserRole({
        user_id: userId,
        organization_id: activeOrganization?.id || '',
        role_name: role
      });

      const updatedUsers = users.map((user: any) => {
        if (user.id === userId) {
          return {
            ...user,
            role
          };
        }
        return user;
      });

      setUsers(updatedUsers);
      toast.success('User updated successfully');
    } catch (error) {
      toast.error('Failed to update user');
    }
  };

  const getRoleBadgeVariant = (role: string) => {
    switch (role) {
      case 'Owner':
        return 'default';
      case 'Admin':
        return 'destructive';
      case 'Member':
        return 'default';
      case 'Viewer':
        return 'secondary';
      default:
        return 'outline';
    }
  };

  const handleUpdateTeam = async () => {
    setEditTeamDialogOpen(false);
    if (teamName.length <= 0 || teamDescription.length <= 0) {
      toast.error('Team name and description are required');
      setTeamName(activeOrganization?.name || '');
      setTeamDescription(activeOrganization?.description || '');
      return;
    }

    if (
      teamName !== activeOrganization?.name ||
      teamDescription !== activeOrganization?.description
    ) {
      await updateOrganizationDetails({
        id: activeOrganization?.id || '',
        name: teamName,
        description: teamDescription
      });
      toast.success('Team details updated successfully');
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
