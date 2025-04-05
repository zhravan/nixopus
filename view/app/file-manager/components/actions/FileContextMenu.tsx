import React from 'react';
import { Button } from '@/components/ui/button';
import { Upload, FolderPlus, Eye, EyeOff, MoreVertical } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import { Dialog, DialogTrigger } from '@/components/ui/dialog';
import FileUpload from './Upload';
import { FileData } from '@/redux/types/files';
import { useCreateDirectoryMutation } from '@/redux/services/file-manager/fileManagersApi';

interface FileContextMenuProps {
  refetch: () => void;
  showHidden: boolean;
  setShowHidden: React.Dispatch<React.SetStateAction<boolean>>;
  currentPath: string;
  setSelectedPath: React.Dispatch<React.SetStateAction<string>>;
  selectedPath: string;
  files: FileData[];
}

const FileContextMenu: React.FC<FileContextMenuProps> = ({
  refetch,
  showHidden,
  setShowHidden,
  currentPath,
  setSelectedPath,
  selectedPath,
  files
}) => {
  const [isUploadOpen, setIsUploadOpen] = React.useState(false);
  const [createDirectory] = useCreateDirectoryMutation();

  const handleCreateDirectory = async () => {
    try {
      const newFolderName = `New Folder ${files.filter((f) => f.name.startsWith('New Folder')).length + 1}`;
      await createDirectory({ path: currentPath, name: newFolderName });
      refetch();
    } catch (error) {
      console.error('Failed to create directory:', error);
    }
  };

  return (
    <div className="flex items-center space-x-2">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="icon">
            <MoreVertical className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <Dialog open={isUploadOpen} onOpenChange={setIsUploadOpen}>
            <DialogTrigger asChild>
              <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                <Upload className="mr-2 h-4 w-4" />
                Upload File
              </DropdownMenuItem>
            </DialogTrigger>
            <FileUpload setIsDialogOpen={setIsUploadOpen} />
          </Dialog>

          <DropdownMenuItem onSelect={handleCreateDirectory}>
            <FolderPlus className="mr-2 h-4 w-4" />
            Create Directory
          </DropdownMenuItem>

          <DropdownMenuItem onSelect={() => setShowHidden(!showHidden)}>
            {showHidden ? (
              <>
                <EyeOff className="mr-2 h-4 w-4" />
                Hide Hidden Files
              </>
            ) : (
              <>
                <Eye className="mr-2 h-4 w-4" />
                Show Hidden Files
              </>
            )}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

export default FileContextMenu;
