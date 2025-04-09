import {
  useCopyFileOrDirectoryMutation,
  useMoveOrRenameDirectoryMutation
} from '@/redux/services/file-manager/fileManagersApi';
import { FileData } from '@/redux/types/files';
import { toast } from 'sonner';
import { useTranslation } from '@/hooks/use-translation';
export const useFileManagerActionsHook = () => {
  const { t } = useTranslation();
  const [
    moveOrRenameDirectory,
    { isLoading: isMoveOrRenameDirectoryLoading, error: moveOrRenameDirectoryError }
  ] = useMoveOrRenameDirectoryMutation();
  const [
    copyFileOrDirectory,
    { isLoading: isCopyFileOrDirectoryLoading, error: copyFileOrDirectoryError }
  ] = useCopyFileOrDirectoryMutation();

  const handleCopyFile = (
    data: FileData,
    setFileToCopy: React.Dispatch<React.SetStateAction<FileData | undefined>>
  ) => {
    try {
      setFileToCopy(data);
    } catch (error) {
      toast.error(t('toasts.errors.setFileToCopy'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  const handleFileMove = async (from_path: string, to_path: string) => {
    try {
      await moveOrRenameDirectory({ from_path, to_path });
    } catch (error) {
      toast.error(t('toasts.errors.moveFile'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  const handleFilePaste = async (from_path: string, to_path: string) => {
    try {
      await copyFileOrDirectory({ from_path, to_path });
    } catch (error) {
      toast.error(t('toasts.errors.copyFile'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  return {
    handleCopyFile,
    handleFileMove,
    handleFilePaste,
    copyFileOrDirectoryError,
    moveOrRenameDirectoryError,
    isCopyFileOrDirectoryLoading,
    isMoveOrRenameDirectoryLoading
  };
};
