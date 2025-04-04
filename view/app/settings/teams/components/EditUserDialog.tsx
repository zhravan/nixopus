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
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { UserTypes } from '@/redux/types/orgs';
import { toast } from 'sonner';
import { Input } from '@/components/ui/input';
import { Search } from 'lucide-react';

interface Permission {
  resource: string;
  action: 'read' | 'create' | 'update' | 'delete';
}

interface EditUserDialogProps {
  isOpen: boolean;
  onClose: () => void;
  user: {
    id: string;
    name: string;
    email: string;
    avatar: string;
    role: 'Owner' | 'Admin' | 'Member' | 'Viewer';
    permissions: string[];
  };
  resources: string[];
  onSave: (userId: string, role: UserTypes, permissions: Permission[]) => void;
}

const AVAILABLE_ROLES: { value: UserTypes; label: string }[] = [
  { value: 'admin', label: 'Admin' },
  { value: 'member', label: 'Member' },
  { value: 'viewer', label: 'Viewer' }
];

const ACCESS_LEVELS = ['read', 'create', 'update', 'delete'] as const;

function EditUserDialog({ isOpen, onClose, user, resources, onSave }: EditUserDialogProps) {
  const [selectedRole, setSelectedRole] = React.useState<UserTypes>('admin');
  const [selectedPermissions, setSelectedPermissions] = React.useState<Permission[]>([]);
  const [searchQuery, setSearchQuery] = React.useState('');
  const [selectedResource, setSelectedResource] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (isOpen) {
      const role = user.role.toLowerCase() as UserTypes;
      setSelectedRole(role);

      const existingPermissions = user.permissions.map((p) => {
        const [resource, action] = p.split(':');
        return {
          resource: resource.toLowerCase(),
          action: action.toLowerCase() as Permission['action']
        };
      });
      setSelectedPermissions(existingPermissions);

      if (existingPermissions.length > 0) {
        setSelectedResource(existingPermissions[0].resource);
      }
    }
  }, [isOpen, user]);

  const handleRoleChange = (value: string) => {
    const newRole = value as UserTypes;
    setSelectedRole(newRole);
    setSelectedPermissions([]);
  };

  const filteredResources = resources.filter((resource) =>
    resource.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handlePermissionToggle = (resource: string, action: Permission['action']) => {
    setSelectedPermissions((prev) => {
      const exists = prev.some((p) => p.resource === resource && p.action === action);
      if (exists) {
        return prev.filter((p) => !(p.resource === resource && p.action === action));
      }
      return [...prev, { resource, action }];
    });
  };

  const handleSave = () => {
    if (!selectedRole) {
      toast.error('Please select a role');
      return;
    }
    onSave(user.id, selectedRole, selectedPermissions);
    onClose();
  };

  const isPermissionSelected = (resource: string, action: Permission['action']) => {
    return selectedPermissions.some((p) => p.resource === resource && p.action === action);
  };

  const getResourcePermissions = (resource: string) => {
    return selectedPermissions.filter((p) => p.resource === resource);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[800px]">
        <DialogHeader>
          <DialogTitle>Edit User Permissions</DialogTitle>
          <DialogDescription>Update role and permissions for {user.name}</DialogDescription>
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

          <div className="space-y-4">
            <div className="flex items-center space-x-2">
              <div className="relative flex-1">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search resources..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-8"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label className="text-base font-medium">Resources</Label>
                <div className="space-y-1">
                  {filteredResources.map((resource) => {
                    const resourcePermissions = getResourcePermissions(resource);
                    return (
                      <Button
                        key={resource}
                        variant={selectedResource === resource ? 'secondary' : 'ghost'}
                        className="w-full justify-start"
                        onClick={() => setSelectedResource(resource)}
                      >
                        <Badge variant="outline" className="mr-2">
                          {resource.toUpperCase()}
                        </Badge>
                        {resourcePermissions.length > 0 && (
                          <Badge variant="secondary" className="ml-auto">
                            {resourcePermissions.length} selected
                          </Badge>
                        )}
                      </Button>
                    );
                  })}
                </div>
              </div>

              <div className="space-y-2">
                {selectedResource ? (
                  <div className="space-y-4">
                    <div className="space-y-2 pl-4">
                      {ACCESS_LEVELS.map((action) => {
                        const isSelected = isPermissionSelected(selectedResource, action);
                        return (
                          <div
                            key={`${selectedResource}-${action}`}
                            className="flex items-center space-x-2"
                          >
                            <Checkbox
                              id={`${selectedResource}-${action}`}
                              checked={isSelected}
                              onCheckedChange={() =>
                                handlePermissionToggle(selectedResource, action)
                              }
                            />
                            <Label
                              htmlFor={`${selectedResource}-${action}`}
                              className="text-sm font-normal capitalize"
                            >
                              {action}
                            </Label>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                ) : (
                  <div className="text-sm text-muted-foreground">
                    Select a resource to view its access levels
                  </div>
                )}
              </div>
            </div>
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
