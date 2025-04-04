import { combineReducers, configureStore, ConfigureStoreOptions } from '@reduxjs/toolkit';
import {
  persistStore,
  persistReducer,
  FLUSH,
  REHYDRATE,
  PAUSE,
  PERSIST,
  PURGE,
  REGISTER
} from 'redux-persist';
import createWebStorage from 'redux-persist/lib/storage/createWebStorage';
import { authApi } from '@/redux/services/users/authApi';
import authReducer from '@/redux/features/users/authSlice';
import { userApi } from '@/redux/services/users/userApi';
import userSlice from '@/redux/features/users/userSlice';
import { notificationApi } from '@/redux/services/settings/notificationApi';
import { domainsApi } from '@/redux/services/settings/domainsApi';
import { GithubConnectorApi } from '@/redux/services/connector/githubConnectorApi';
import githubConnector from './features/github-connector/githubConnectorSlice';
import { deployApi } from './services/deploy/applicationsApi';
import { fileManagersApi } from './services/file-manager/fileManagersApi';
import { auditApi } from './services/audit';

const createNoopStorage = () => {
  return {
    getItem(_key: string) {
      return Promise.resolve(null);
    },
    setItem(_key: string, value: any) {
      return Promise.resolve(value);
    },
    removeItem(_key: string) {
      return Promise.resolve();
    }
  };
};

const storage = typeof window !== 'undefined' ? createWebStorage('local') : createNoopStorage();

const persistConfig = {
  key: 'root',
  storage,
  whitelist: ['auth']
};

const rootReducer = combineReducers({
  [authApi.reducerPath]: authApi.reducer,
  auth: authReducer,
  [userApi.reducerPath]: userApi.reducer,
  notificationApi: notificationApi.reducer,
  [domainsApi.reducerPath]: domainsApi.reducer,
  [GithubConnectorApi.reducerPath]: GithubConnectorApi.reducer,
  githubConnector: githubConnector,
  [deployApi.reducerPath]: deployApi.reducer,
  user: userSlice,
  fileManagersApi: fileManagersApi.reducer,
  [auditApi.reducerPath]: auditApi.reducer
});

const persistedReducer = persistReducer(persistConfig, rootReducer);

const storeOptions: ConfigureStoreOptions = {
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER]
      }
    }).concat(
      authApi.middleware,
      userApi.middleware,
      notificationApi.middleware,
      domainsApi.middleware,
      GithubConnectorApi.middleware,
      deployApi.middleware,
      fileManagersApi.middleware,
      auditApi.middleware
    )
};

export const store = configureStore(storeOptions);

export const persistor = persistStore(store);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
