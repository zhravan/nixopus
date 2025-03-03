import React from 'react';
import { Button } from '@/components/ui/button';

interface DomainActionsProps {
  domainId: string;
}

export function DomainActions({ domainId }: DomainActionsProps) {
  const handleEdit = () => {
    console.log(`Editing domain ${domainId}`);
  };

  const handleDelete = () => {
    console.log(`Deleting domain ${domainId}`);
  };

  const handleManage = () => {
    console.log(`Managing domain ${domainId}`);
  };

  return (
    <div className="flex justify-end gap-2">
      <Button variant="ghost" className="text-primary p-0 m-0" onClick={handleEdit}>
        Edit
      </Button>
      <Button variant="ghost" className="text-destructive p-0 m-0" onClick={handleDelete}>
        Delete
      </Button>
      <Button variant="ghost" className="text-muted-foreground p-0 m-0" onClick={handleManage}>
        Manage
      </Button>
    </div>
  );
}
