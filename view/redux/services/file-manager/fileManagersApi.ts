import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from '@/redux/base-query';
import { FileData } from '@/redux/types/files';
import { FILEMANAGERURLS } from '@/redux/api-conf';

export const fileManagersApi = createApi({
  reducerPath: 'fileManagersApi',
  baseQuery: baseQueryWithReauth,
  tagTypes: ['fileManager'],
  endpoints: (builder) => ({
    getFilesInPath: builder.query<FileData[], { path: string }>({
      query: ({ path }) => ({
        url: `${FILEMANAGERURLS.LIST_FILES_AT_PATH}?path=${encodeURIComponent(path)}`,
        method: 'GET'
      }),
      transformResponse: (response: { data: FileData[] }) => response.data
    }),
    createDirectory: builder.mutation<null, { path: string; name: string }>({
      query: ({ path, name }) => ({
        url: FILEMANAGERURLS.CREATE_DIRECTORY,
        method: 'POST',
        body: { path: path + '/' + name }
      }),
      transformResponse: (response: any) => response
    }),
    deleteDirectory: builder.mutation<any, { path: string }>({
      query: ({ path }) => ({
        url: FILEMANAGERURLS.DELETE_DIRECTORY,
        method: 'DELETE',
        body: { path }
      }),
      transformResponse: (response: any) => response
    }),
    deleteFile: builder.mutation<any, { path: string }>({
      query: ({ path }) => ({
        url: FILEMANAGERURLS.DELETE_FILE,
        method: 'DELETE',
        body: { path }
      }),
      transformResponse: (response: any) => response
    }),
    moveOrRenameDirectory: builder.mutation<any, { from_path: string; to_path: string }>({
      query: ({ from_path, to_path }) => ({
        url: FILEMANAGERURLS.MOVE_FOLDER_FILES_RECURSIVELY_OR_RENAME,
        method: 'POST',
        body: { from_path, to_path }
      }),
      transformResponse: (response: any) => response
    }),
    copyFileOrDirectory: builder.mutation<any, { from_path: string; to_path: string }>({
      query: ({ from_path, to_path }) => ({
        url: FILEMANAGERURLS.COPY_FOLDER_FILES_RECURSIVELY,
        method: 'POST',
        body: { from_path, to_path }
      }),
      transformResponse: (response: any) => response
    }),
    createFile: builder.mutation<any, { path: string; name: string }>({
      query: ({ path, name }) => ({
        url: FILEMANAGERURLS.CREATE_FILE,
        method: 'POST',
        body: { path, name }
      }),
      transformResponse: (response: any) => response
    }),
    calculateDirectorySize: builder.mutation<any['data'], { path: string }>({
      query: ({ path }) => ({
        url: FILEMANAGERURLS.CALCULATE_DIRECTORY_SIZE,
        method: 'POST',
        body: { path }
      }),
      transformResponse: (response: any) => response.data
    }),
    uploadFile: builder.mutation<any, { file: File; path: string }>({
      query: ({ file, path }) => {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('path', path);

        return {
          url: FILEMANAGERURLS.UPLOAD_FILE,
          method: 'POST',
          body: formData
        };
      },
      transformResponse: (response: any) => response
    })
  })
});

export const {
  useGetFilesInPathQuery,
  useCreateDirectoryMutation,
  useDeleteDirectoryMutation,
  useDeleteFileMutation,
  useMoveOrRenameDirectoryMutation,
  useCopyFileOrDirectoryMutation,
  useCreateFileMutation,
  useCalculateDirectorySizeMutation,
  useUploadFileMutation
} = fileManagersApi;
