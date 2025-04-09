import { useState, useEffect } from 'react';
import { FileData } from '@/redux/types/files';
import {
  useCalculateDirectorySizeMutation,
  useDeleteDirectoryMutation,
  useMoveOrRenameDirectoryMutation
} from '@/redux/services/file-manager/fileManagersApi';
import { toast } from 'sonner';
import { useTranslation } from '@/hooks/use-translation';

export const useFileOperations = (file: FileData, refetch: () => void) => {
  const { t } = useTranslation();
  const [moveOrRenameDirectory] = useMoveOrRenameDirectoryMutation();
  const [deleteDirectory] = useDeleteDirectoryMutation();
  const [calculateDirectorySize, { isLoading: isSizeLoading, data: fileSize }] =
    useCalculateDirectorySizeMutation();

  const [isEditing, setIsEditing] = useState(false);
  const [editedFileName, setEditedFileName] = useState(file.name);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  useEffect(() => {
    if (isDialogOpen) {
      calculateDirectorySize({ path: file.path });
    }
  }, [isDialogOpen, file.path, calculateDirectorySize]);

  const handleRename = async () => {
    if (editedFileName !== file.name) {
      try {
        const from_path = file.path;
        const to_path = file.path.replace(file.name, editedFileName);
        await moveOrRenameDirectory({ from_path, to_path });
        refetch();
      } catch (error) {
        toast.error(t('toasts.errors.renameFile'), {
          description: error instanceof Error ? error.message : 'Unknown error'
        });
        setEditedFileName(file.name);
      }
    } else {
      setIsEditing(false);
    }
  };

  const handleDelete = async () => {
    try {
      await deleteDirectory({ path: file.path });
      refetch();
    } catch (error) {
      toast.error(t('toasts.errors.deleteFile'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  const startRenaming = () => {
    try {
      setIsEditing(true);
      setEditedFileName(file.name);
    } catch (error) {
      toast.error(t('toasts.errors.startRenaming'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleRename();
    } else if (e.key === 'Escape') {
      setIsEditing(false);
      setEditedFileName(file.name);
    }
  };

  const handleTextDoubleClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      setIsEditing(true);
    } catch (error) {
      toast.error(t('toasts.errors.startRenaming'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  const onDeleteFolder = async () => {
    try {
      await deleteDirectory({ path: file.path });
      refetch();
    } catch (error) {
      toast.error(t('toasts.errors.deleteFile'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  return {
    isEditing,
    setIsEditing,
    editedFileName,
    setEditedFileName,
    isDialogOpen,
    setIsDialogOpen,
    isSizeLoading,
    fileSize,
    handleRename,
    handleDelete,
    startRenaming,
    handleKeyDown,
    handleTextDoubleClick,
    onDeleteFolder
  };
};
