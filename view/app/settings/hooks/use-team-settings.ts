import React, { useState } from 'react';

function useTeamSettings() {
  const [users, setUsers] = useState<any>([
    {
      id: 1,
      name: 'John Doe',
      email: 'john@example.com',
      role: 'Admin',
      permissions: ['Read', 'Write', 'Delete'],
      avatar: '/api/placeholder/30/30'
    },
    {
      id: 2,
      name: 'Jane Smith',
      email: 'jane@example.com',
      role: 'Member',
      permissions: ['Read', 'Write'],
      avatar: '/api/placeholder/30/30'
    },
    {
      id: 3,
      name: 'Bob Johnson',
      email: 'bob@example.com',
      role: 'Viewer',
      permissions: ['Read'],
      avatar: '/api/placeholder/30/30'
    }
  ]);

  const [isAddUserDialogOpen, setIsAddUserDialogOpen] = useState(false);
  const [newUser, setNewUser] = useState({ name: '', email: '', role: 'Member' });

  const handleAddUser = () => {
    const id = users.length + 1;
    let permissions: string[] = [];

    switch (newUser.role) {
      case 'Admin':
        permissions = ['Read', 'Write', 'Delete', 'Manage'];
        break;
      case 'Member':
        permissions = ['Read', 'Write'];
        break;
      case 'Viewer':
        permissions = ['Read'];
        break;
    }

    setUsers([
      ...users,
      {
        id,
        name: newUser.name,
        email: newUser.email,
        role: newUser.role,
        permissions,
        avatar: '/api/placeholder/30/30'
      }
    ]);

    setNewUser({ name: '', email: '', role: 'Member' });
    setIsAddUserDialogOpen(false);
  };

  const handleRemoveUser = (userId: string) => {
    setUsers(users.filter((user: any) => user.id !== userId));
  };

  const getRoleBadgeVariant = (role: string) => {
    switch (role) {
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

  return {
    users,
    isAddUserDialogOpen,
    setIsAddUserDialogOpen,
    newUser,
    setNewUser,
    handleAddUser,
    handleRemoveUser,
    getRoleBadgeVariant
  };
}

export default useTeamSettings;
