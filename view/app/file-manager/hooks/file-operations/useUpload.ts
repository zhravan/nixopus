import { useState, useRef } from 'react';
import { useUploadFileMutation } from '@/redux/services/file-manager/fileManagersApi';
import useFileManager from '../ui/useFileManager';
import { toast } from 'sonner';
import { useTranslation } from '@/hooks/use-translation';

interface FileObject {
  id: string;
  file: File;
  progress: number;
}

export default function use_file_upload() {
  const { t } = useTranslation();
  const [files, setFiles] = useState<FileObject[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { currentPath, refetch } = useFileManager();
  const [uploadFile] = useUploadFileMutation();

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const droppedFiles = Array.from(e.dataTransfer.files);
    addFiles(droppedFiles);
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = Array.from(e.target.files || []);
    addFiles(selectedFiles);
  };

  const addFiles = (newFiles: File[]) => {
    const fileObjects = newFiles.map((file) => ({
      id: Math.random().toString(36).substr(2, 9),
      file,
      progress: 0
    }));

    setFiles((prev) => [...prev, ...fileObjects]);
    handleUpload();
  };

  const handleUpload = async () => {
    for (const fileObj of files) {
      try {
        await uploadFile({ file: fileObj.file, path: currentPath });
        setFiles((prev) => prev.filter((f) => f.id !== fileObj.id));
      } catch (error) {
        toast.error(t('toasts.errors.uploadFile'), {
          description: error instanceof Error ? error.message : 'Unknown error'
        });
      }
    }
    refetch();
  };

  const deleteFile = (id: string) => {
    try {
      setFiles((prev) => prev.filter((file) => file.id !== id));
    } catch (error) {
      toast.error(t('toasts.errors.deleteFile'), {
        description: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  };

  return {
    files,
    isDragging,
    fileInputRef,
    handleDragEnter,
    handleDragLeave,
    handleDragOver,
    handleDrop,
    handleFileInput,
    deleteFile
  };
}
