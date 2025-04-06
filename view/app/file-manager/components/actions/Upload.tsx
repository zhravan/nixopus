import React from 'react';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import { Upload, Trash } from 'lucide-react';
import { DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import use_file_upload from '../../hooks/file-operations/useUpload';
import { useTranslation } from '@/hooks/use-translation';

interface FileUploadProps {
  setIsDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const FileUpload: React.FC<FileUploadProps> = ({ setIsDialogOpen }) => {
  const { t } = useTranslation();
  const {
    files,
    isDragging,
    fileInputRef,
    handleDragEnter,
    handleDragLeave,
    handleDragOver,
    handleDrop,
    handleFileInput,
    deleteFile
  } = use_file_upload();

  return (
    <DialogContent>
      <DialogHeader className="space-y-8">
        <DialogTitle>{t('fileManager.upload.title')}</DialogTitle>
        <div className="mx-auto w-full">
          <div
            className={`cursor-pointer rounded-lg border-2 border-dashed p-10 text-center ${
              isDragging ? 'border-blue-500 bg-muted' : 'border-gray-300'
            }`}
            onDragEnter={handleDragEnter}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            onClick={() => fileInputRef.current?.click()}
          >
            <Upload className="mx-auto h-12 w-12 text-muted-foreground" />
            <p className="text-md mt-2 text-muted-foreground">{t('fileManager.upload.dragDrop')}</p>
          </div>
          <input
            type="file"
            ref={fileInputRef}
            className="hidden"
            onChange={handleFileInput}
            multiple
          />
          <div className="mt-4 space-y-2">
            {files.map((fileObj) => (
              <div key={fileObj.id} className="rounded bg-muted p-4 shadow">
                <div className="flex items-center justify-between">
                  <span className="truncate">{fileObj.file.name}</span>
                  <Button variant="ghost" size="icon" onClick={() => deleteFile(fileObj.id)}>
                    <Trash className="h-4 w-4" />
                  </Button>
                </div>
                <Progress value={fileObj.progress} className="mt-2 h-2 bg-muted-foreground" />
              </div>
            ))}
          </div>
        </div>
      </DialogHeader>
    </DialogContent>
  );
};

export default FileUpload;
