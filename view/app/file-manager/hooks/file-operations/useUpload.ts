import { useState, useRef } from 'react';
import { useUploadFileMutation } from '@/redux/services/file-manager/fileManagersApi';
import useFileManager from '../ui/useFileManager';

interface FileObject {
  id: string;
  file: File;
  progress: number;
}

export default function use_file_upload() {
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
    uploadFiles(fileObjects);
  };

  const uploadFiles = async (fileObjects: FileObject[]) => {
    for (const fileObj of fileObjects) {
      try {
        const formData = new FormData();
        formData.append('file', fileObj.file);
        formData.append('path', currentPath);

        await uploadFile({ file: fileObj.file, path: currentPath });
        setFiles((prev) => prev.filter((f) => f.id !== fileObj.id));
      } catch (error) {
        console.error('Failed to upload file:', error);
      }
    }
    refetch();
  };

  const deleteFile = (id: string) => {
    setFiles((prev) => prev.filter((file) => file.id !== id));
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
