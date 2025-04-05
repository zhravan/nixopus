import { BASE_URL } from '@/redux/conf';
import { RootState } from '@/redux/store';
import {
  BaseQueryFn,
  FetchArgs,
  FetchBaseQueryError,
  fetchBaseQuery
} from '@reduxjs/toolkit/query/react';
import { Mutex } from 'async-mutex';

const mutex = new Mutex();

export const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  const baseQuery = fetchBaseQuery({
    baseUrl: BASE_URL,
    prepareHeaders: (headers, { getState }) => {
      const token = (getState() as RootState).auth.token;
      const organizationId =
        (getState() as RootState).user.activeOrganization?.id ||
        (getState() as RootState).auth.user?.organization_users?.[0]?.organization_id;
      if (token) {
        headers.set('authorization', `Bearer ${token}`);
      }
      if (organizationId) {
        headers.set('X-Organization-Id', organizationId);
      }
      return headers;
    },
    credentials: 'include'
  });

  await mutex.waitForUnlock();
  let result = await baseQuery(args, api, extraOptions);

  if (result.error && result.error.status === 401) {
    console.warn('Unauthorized request, logging out');
    api.dispatch({ type: 'auth/logout' });
  }

  return result;
};
