import { DOMAIN_SETTINGS } from '@/redux/api-conf';
import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from '@/redux/base-query';
import {
    CreateSMTPConfigRequest,
    SMTPConfig,
    UpdateSMTPConfigRequest
} from '@/redux/types/notification';

export const domainsApi = createApi({
    reducerPath: 'domainsApi',
    baseQuery: baseQueryWithReauth,
    tagTypes: ['Domains'],
    endpoints: (builder) => ({
        getAllDomains: builder.query<SMTPConfig, void>({
            query: () => ({
                url: DOMAIN_SETTINGS.GET_DOMAINS,
                method: 'GET'
            }),
            providesTags: [{ type: 'Domains', id: 'LIST' }],
            transformResponse: (response: { data: SMTPConfig }) => {
                return response.data;
            }
        }),
        createDomain: builder.mutation<null, CreateSMTPConfigRequest>({
            query: (data) => ({
                url: DOMAIN_SETTINGS.ADD_DOMAIN,
                method: 'POST',
                body: data
            }),
            invalidatesTags: [{ type: 'Domains', id: 'LIST' }],
            transformResponse: (response: { data: null }) => {
                return response.data;
            }
        }),
        updateDomain: builder.mutation<null, UpdateSMTPConfigRequest>({
            query: (data) => ({
                url: DOMAIN_SETTINGS.UPDATE_DOMAIN,
                method: 'PUT',
                body: data
            }),
            invalidatesTags: [{ type: 'Domains', id: 'LIST' }],
            transformResponse: (response: { data: null }) => {
                return response.data;
            }
        }),
        deleteDomain: builder.mutation<null, string>({
            query: (id) => ({
                url: DOMAIN_SETTINGS.DELETE_DOMAIN,
                method: 'DELETE',
                params: { id }
            }),
            invalidatesTags: [{ type: 'Domains', id: 'LIST' }],
            transformResponse: (response: { data: null }) => {
                return response.data;
            }
        }),
    })
});

export const {
    useGetAllDomainsQuery,
    useCreateDomainMutation,
    useUpdateDomainMutation,
    useDeleteDomainMutation
} = domainsApi;
