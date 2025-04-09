'use client';
import React from 'react';
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger
} from '@/components/ui/context-menu';
import {
  Copy,
  EyeIcon,
  EyeOffIcon,
  FolderPlusIcon,
  Move,
  TrashIcon,
  UploadCloudIcon
} from 'lucide-react';
import Skeleton from './components/skeleton/Skeleton';
import Header from './components/layout/Header';
import { SearchBar } from '@/components/ui/search-bar';
import { Breadcrumbs } from './components/breadcrumbs/Breadcrumbs';
import SortMethods from './components/actions/SortMethods';
import LayoutSwitcher from './components/layout/LayoutSwitcher';
import Actions from './components/actions/Actions';
import { FileItem } from './components/file-list/FileItem';
import useFileManager from './hooks/ui/useFileManager';
import { useAppSelector } from '@/redux/hooks';
import { hasPermission } from '@/lib/permission';
import { useTranslation } from '@/hooks/use-translation';
import { DeleteDialog } from '@/components/ui/delete-dialog';
import { FileData } from '@/redux/types/files';

function FileManager() {
  const { t } = useTranslation();
  const user = useAppSelector((state) => state.auth.user);
  const activeOrg = useAppSelector((state) => state.user.activeOrganization);
  const [fileToDelete, setFileToDelete] = React.useState<FileData | null>(null);

  const {
    currentPath,
    layout,
    showHidden,
    selectedPath,
    fileToCopy,
    fileToMove,
    isCopyFileOrDirectoryLoading,
    isMoveOrRenameDirectoryLoading,
    handleFileSelect,
    handleFilePaste,
    handleFileMove,
    handleSearchChange,
    handleSortChange,
    visibleFiles,
    createNewFolder,
    fileClicked,
    isLoading,
    refetch,
    searchTerm,
    setLayout,
    setShowHidden,
    setFileToCopy,
    setFileToMove,
    setSelectedPath,
    files,
    handleFileUpload,
    handleDelete
  } = useFileManager();

  const canRead = hasPermission(user, 'file-manager', 'read', activeOrg?.id);
  const canCreate = hasPermission(user, 'file-manager', 'create', activeOrg?.id);
  const canUpdate = hasPermission(user, 'file-manager', 'update', activeOrg?.id);
  const canDelete = hasPermission(user, 'file-manager', 'delete', activeOrg?.id);

  const handleFileDrop = (e: React.DragEvent) => {
    if (!canCreate) return;
    e.preventDefault();
    e.stopPropagation();

    const files = Array.from(e.dataTransfer.files);
    files.forEach(handleFileUpload);
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!canCreate) return;
    const files = Array.from(e.target.files || []);
    files.forEach(handleFileUpload);
  };

  if (!canRead) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold">{t('common.accessDenied')}</h2>
          <p className="text-muted-foreground">{t('common.noPermissionView')}</p>
        </div>
      </div>
    );
  }

  if (isLoading || isCopyFileOrDirectoryLoading || isMoveOrRenameDirectoryLoading) {
    return <Skeleton />;
  }

  return (
    <div onDrop={handleFileDrop} onDragOver={(e) => e.preventDefault()}>
      <div className="mx-auto max-w-7xl">
        <div className="mb-6 flex items-center justify-between px-0 lg:px-6">
          <Header />
          <div className="relative">
            <SearchBar
              searchTerm={searchTerm}
              handleSearchChange={handleSearchChange}
              label={t('common.searchFiles')}
            />
          </div>
        </div>
        <div className="mb-6 flex items-center justify-between px-0 lg:px-6">
          <Breadcrumbs
            breadcrumbs={currentPath.split('/').filter(Boolean)}
            fileClicked={fileClicked}
          />
          <div className="flex items-center gap-4">
            <SortMethods files={visibleFiles} onSortChange={handleSortChange} />
            <LayoutSwitcher layout={layout} setLayout={setLayout} />
            <Actions
              refetch={refetch}
              showHidden={showHidden}
              setShowHidden={setShowHidden}
              currentPath={currentPath}
              setSelectedPath={setSelectedPath}
              selectedPath={selectedPath}
              files={files || []}
              canCreate={canCreate}
              canUpdate={canUpdate}
              canDelete={canDelete}
            />
          </div>
        </div>

        <ContextMenu>
          <ContextMenuTrigger>
            <div
              className={`grid w-full ${
                layout === 'grid'
                  ? 'grid-cols-2 gap-2 sm:grid-cols-3 sm:gap-3 md:grid-cols-4 md:gap-4 lg:grid-cols-6 xl:grid-cols-8'
                  : 'grid-cols-1 gap-2'
              }`}
            >
              {visibleFiles.map((file, index) => (
                <FileItem
                  key={file.path}
                  file={file}
                  onFolderClick={fileClicked}
                  type={file.file_type === 'Directory' ? 'folder' : 'file'}
                  layout={layout}
                  activePath={selectedPath}
                  onFolderClickActive={handleFileSelect}
                  refetch={refetch}
                  setFileToCopy={setFileToCopy}
                  setFileToMove={setFileToMove}
                  index={index}
                  canUpdate={canUpdate}
                  canDelete={canDelete}
                />
              ))}
              {visibleFiles.length === 0 && (
                <div className="col-span-full text-center text-5xl text-muted-foreground">
                  {t('fileManager.noFiles')}
                </div>
              )}
            </div>
          </ContextMenuTrigger>
          <ContextMenuContent>
            {canCreate && (
              <>
                <ContextMenuItem>
                  <label className="flex cursor-pointer items-center">
                    <UploadCloudIcon className="mr-2 h-5 w-5" />
                    <span>{t('fileManager.actions.upload')}</span>
                    <input type="file" className="hidden" onChange={handleFileInput} multiple />
                  </label>
                </ContextMenuItem>
                <ContextMenuItem onSelect={createNewFolder}>
                  <FolderPlusIcon className="mr-2 h-5 w-5" />
                  <span>{t('fileManager.actions.createDirectory')}</span>
                </ContextMenuItem>
              </>
            )}
            {canUpdate && fileToMove && (
              <ContextMenuItem
                onSelect={() => {
                  handleFileMove(fileToMove.path, currentPath + '/' + fileToMove.name);
                  refetch();
                  setFileToMove(undefined);
                }}
              >
                <Move className="mr-2 h-5 w-5" />
                <span>{t('fileManager.item.actions.moveHere')}</span>
              </ContextMenuItem>
            )}
            {canUpdate && fileToCopy && (
              <ContextMenuItem
                onSelect={() => {
                  handleFilePaste(fileToCopy.path, currentPath + '/' + fileToCopy.name);
                  refetch();
                  setFileToCopy(undefined);
                }}
              >
                <Copy className="mr-2 h-5 w-5" />
                <span>{t('fileManager.item.actions.pasteHere')}</span>
              </ContextMenuItem>
            )}
            <ContextMenuItem onSelect={() => setShowHidden(!showHidden)}>
              {showHidden ? (
                <EyeOffIcon className="mr-2 h-5 w-5" />
              ) : (
                <EyeIcon className="mr-2 h-5 w-5" />
              )}
              <span>
                {showHidden
                  ? t('fileManager.actions.hideHidden')
                  : t('fileManager.actions.showHidden')}
              </span>
            </ContextMenuItem>
            {canDelete && (
              <ContextMenuItem
                onSelect={() =>
                  setFileToDelete(
                    selectedPath ? files?.find((f) => f.path === selectedPath) || null : null
                  )
                }
              >
                <TrashIcon className="mr-2 h-5 w-5" />
                <span>{t('fileManager.item.actions.delete')}</span>
              </ContextMenuItem>
            )}
          </ContextMenuContent>
        </ContextMenu>
        <DeleteDialog
          title={t('fileManager.deleteDialog.title')}
          description={
            fileToDelete?.file_type === 'Directory'
              ? t('fileManager.deleteDialog.descriptionDirectory', {
                  name: fileToDelete.name || ''
                })
              : t('fileManager.deleteDialog.descriptionFile', { name: fileToDelete?.name || '' })
          }
          onConfirm={() => {
            if (fileToDelete) {
              handleDelete(fileToDelete.path);
              setFileToDelete(null);
            }
          }}
          open={!!fileToDelete}
          onOpenChange={(open) => !open && setFileToDelete(null)}
          variant="destructive"
          confirmText={t('fileManager.deleteDialog.confirm')}
          cancelText={t('fileManager.deleteDialog.cancel')}
          icon={TrashIcon}
        />
      </div>
    </div>
  );
}

export default FileManager;
