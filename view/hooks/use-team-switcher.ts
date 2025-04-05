import { useAppSelector } from '@/redux/hooks';
import { useCreateOrganizationMutation } from '@/redux/services/users/userApi';
import React from 'react';
import { toast } from 'sonner';

function useTeamSwitcher() {
  const [open, setOpen] = React.useState(false);
  const user = useAppSelector((state) => state.auth.user);
  const isAdmin = React.useMemo(() => user?.type === 'admin', [user]);
  const [teamName, setTeamName] = React.useState('');
  const [teamDescription, setTeamDescription] = React.useState('');
  const [createOrganization, { isLoading }] = useCreateOrganizationMutation();

  const toggleAddTeamModal = () => {
    setOpen(!open);
  };

  const handleTeamNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setTeamName(event.target.value);
  };

  const handleTeamDescriptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setTeamDescription(event.target.value);
  };

  const validateTeamName = (name: string) => {
    if (!name) {
      return false;
    }
    return name.length <= 50;
  };

  const validateTeamDescription = (description: string) => {
    if (!description) {
      return false;
    }
    return description.length <= 100;
  };

  const onCreateTeam = async () => {
    try {
      if (!isAdmin) {
        toast.error('You are not an admin');
        return;
      }

      if (!validateTeamName(teamName)) {
        toast.error('Team name is required and must be less than 50 characters');
        return;
      }

      if (!validateTeamDescription(teamDescription)) {
        toast.error('Team description is required and must be less than 100 characters');
        return;
      }

      const res = await createOrganization({
        name: teamName,
        description: teamDescription
      }).unwrap();

      if (!res.id) {
        toast.error('Failed to create team');
        return;
      }
      toast.success('Team created successfully');
      setOpen(false);
      setTeamName('');
      setTeamDescription('');
    } catch (error) {
      console.error('Failed to create team:', error);
      toast.error('Failed to create team');
    }
  };

  return {
    addTeamModalOpen: open,
    setAddTeamModalOpen: setOpen,
    toggleAddTeamModal,
    createTeam: onCreateTeam,
    teamName,
    teamDescription,
    handleTeamNameChange,
    handleTeamDescriptionChange,
    isLoading
  };
}

export default useTeamSwitcher;
