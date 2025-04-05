import React from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select';
import { UserTypes } from '@/redux/types/orgs';
import { toast } from 'sonner';

interface EditUserDialogProps {
  isOpen: boolean;
  onClose: () => void;
  user: {
    id: string;
    name: string;
    email: string;
    avatar: string;
    role: 'Owner' | 'Admin' | 'Member' | 'Viewer';
  };
  onSave: (userId: string, role: UserTypes) => void;
}

const AVAILABLE_ROLES: { value: UserTypes; label: string }[] = [
  { value: 'admin', label: 'Admin' },
  { value: 'member', label: 'Member' },
  { value: 'viewer', label: 'Viewer' }
];

function EditUserDialog({ isOpen, onClose, user, onSave }: EditUserDialogProps) {
  const [selectedRole, setSelectedRole] = React.useState<UserTypes>('admin');

  React.useEffect(() => {
    if (isOpen) {
      const role = user.role.toLowerCase() as UserTypes;
      setSelectedRole(role);
    }
  }, [isOpen, user]);

  const handleRoleChange = (value: string) => {
    const newRole = value as UserTypes;
    setSelectedRole(newRole);
  };

  const handleSave = () => {
    if (!selectedRole) {
      toast.error('Please select a role');
      return;
    }
    onSave(user.id, selectedRole);
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit User Role</DialogTitle>
          <DialogDescription>Update role for {user.name}</DialogDescription>
        </DialogHeader>
        <div className="space-y-6 py-4">
          <div className="space-y-2">
            <Label>Role</Label>
            <Select value={selectedRole} onValueChange={handleRoleChange}>
              <SelectTrigger>
                <SelectValue placeholder="Select a role" />
              </SelectTrigger>
              <SelectContent>
                {AVAILABLE_ROLES.map((role) => (
                  <SelectItem key={role.value} value={role.value}>
                    {role.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        <div className="flex justify-end space-x-2">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleSave}>Save Changes</Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default EditUserDialog;
