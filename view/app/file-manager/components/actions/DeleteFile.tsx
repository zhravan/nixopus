import React from 'react';
import { Button } from '@/components/ui/button';
import { DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useDeleteFileMutation } from '@/redux/services/file-manager/fileManagersApi';

interface DeleteFileProps {
  setIsDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
  selectedPath: string;
  refetch: () => void;
}

const DeleteFile: React.FC<DeleteFileProps> = ({ setIsDialogOpen, selectedPath, refetch }) => {
  const [deleteFile] = useDeleteFileMutation();

  const handleDelete = async () => {
    if (selectedPath) {
      await deleteFile({ path: selectedPath });
      refetch();
      setIsDialogOpen(false);
    }
  };

  return (
    <DialogContent>
      <DialogHeader className="space-y-8">
        <DialogTitle>Delete File</DialogTitle>
        <div className="space-y-4">
          <p>Are you sure you want to delete this file?</p>
          <div className="flex justify-end space-x-2">
            <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete} disabled={!selectedPath}>
              Delete
            </Button>
          </div>
        </div>
      </DialogHeader>
    </DialogContent>
  );
};

export default DeleteFile;
