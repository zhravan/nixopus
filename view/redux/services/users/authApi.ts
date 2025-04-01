import { AUTHURLS } from '@/redux/api-conf';
import { baseQueryWithReauth } from '@/redux/base-query';
import { AuthResponse, LoginPayload, RefreshTokenPayload, User } from '@/redux/types/user';
import { createApi } from '@reduxjs/toolkit/query/react';
import { fetchBaseQuery } from '@reduxjs/toolkit/query';

export const authApi = createApi({
  reducerPath: 'authApi',
  baseQuery: baseQueryWithReauth,
  tagTypes: ['Authentication'],
  endpoints: (builder) => ({
    loginUser: builder.mutation<AuthResponse, LoginPayload>({
      query(credentials) {
        return {
          url: AUTHURLS.USER_LOGIN,
          method: 'POST',
          body: credentials
        };
      },
      transformResponse: (response: { data: AuthResponse }) => {
        return { ...response.data };
      },
      invalidatesTags: [{ type: 'Authentication', id: 'LIST' }]
    }),
    getUserDetails: builder.query<User, void>({
      query: () => ({
        url: AUTHURLS.USER_DETAILS,
        method: 'GET'
      }),
      providesTags: [{ type: 'Authentication', id: 'LIST' }],
      transformResponse: (response: { data: User }) => {
        return { ...response.data };
      }
    }),
    refreshToken: builder.mutation<AuthResponse, RefreshTokenPayload>({
      query: (payload) => ({
        url: AUTHURLS.REFRESH_TOKEN,
        method: 'POST',
        body: payload
      }),
      transformResponse: (response: { data: AuthResponse }) => {
        return { ...response.data };
      }
    }),
    resetPassword: builder.mutation<void, { token: string; password: string }>({
      query({ token, password }) {
        return {
          url: `${AUTHURLS.RESET_PASSWORD}?token=${token}`,
          method: 'POST',
          body: { password }
        };
      }
    })
  })
});

export const {
  useLoginUserMutation,
  useGetUserDetailsQuery,
  useRefreshTokenMutation,
  useResetPasswordMutation
} = authApi;
