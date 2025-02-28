import { authApi } from '@/redux/services/users/authApi';
import { User } from '@/redux/types/user';
import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { createAsyncThunk } from '@reduxjs/toolkit';

interface AuthState {
    user: User | null;
    token: string | null;
    refreshToken: string | null;
    expiresIn: number | null;
    isAuthenticated: boolean;
    isInitialized: boolean;
}

interface AuthPayload {
    user: User | null;
    token: string | null;
    refreshToken: string | null;
    expiresIn: number | null;
}

const isClient = typeof window !== 'undefined';

const getTokenFromStorage = () => {
    if (!isClient) return null;
    try {
        return localStorage.getItem('token');
    } catch (error) {
        console.error('Failed to get token from localStorage:', error);
        return null;
    }
};

const getRefreshTokenFromStorage = () => {
    if (!isClient) return null;
    try {
        return localStorage.getItem('refreshToken');
    } catch (error) {
        console.error('Failed to get refresh token from localStorage:', error);
        return null;
    }
};

const getTokenExpiry = () => {
    if (!isClient) return null;
    try {
        const expiryStr = localStorage.getItem('tokenExpiry');
        return expiryStr ? parseInt(expiryStr, 10) : null;
    } catch (error) {
        console.error('Failed to get token expiry from localStorage:', error);
        return null;
    }
};

const isTokenExpired = () => {
    const expiry = getTokenExpiry();
    if (!expiry) return true;

    return (expiry - 10000) < Date.now();
};

export const initializeAuth = createAsyncThunk<AuthPayload | null>(
    'auth/initialize',
    async (_, { dispatch, rejectWithValue }) => {
        if (!isClient) return null;

        try {
            const token = getTokenFromStorage();
            const refreshToken = getRefreshTokenFromStorage();

            if (!token || !refreshToken) {
                console.log('No tokens found in storage');
                return null;
            }

            if (isTokenExpired()) {
                console.log('Token expired, attempting refresh');
                try {
                    const refreshResult = await dispatch(
                        authApi.endpoints.refreshToken.initiate({ refresh_token: refreshToken })
                    ).unwrap();

                    if (refreshResult?.accessToken) {
                        console.log('Token refreshed successfully');

                        const userResult = await dispatch(
                            authApi.endpoints.getUserDetails.initiate(undefined)
                        ).unwrap();

                        return {
                            user: userResult,
                            token: refreshResult.accessToken,
                            refreshToken: refreshResult.refreshToken || refreshToken,
                            expiresIn: refreshResult.expiresIn
                        };
                    } else {
                        console.error('Invalid refresh response');
                        return null;
                    }
                } catch (refreshError) {
                    console.error('Failed to refresh token:', refreshError);
                    removeTokensFromStorage();
                    return rejectWithValue('Token refresh failed');
                }
            } else {
                console.log('Token valid, fetching user details');
                try {
                    const userResult = await dispatch(
                        authApi.endpoints.getUserDetails.initiate(undefined)
                    ).unwrap();

                    const expiresIn = getTokenExpiry();
                    return {
                        user: userResult,
                        token,
                        refreshToken,
                        expiresIn: expiresIn ? Math.floor((expiresIn - Date.now()) / 1000) : null
                    };
                } catch (error) {
                    console.error('Failed to fetch user details:', error);
                    return rejectWithValue('Failed to fetch user details');
                }
            }
        } catch (error) {
            console.error('Failed to initialize auth:', error);
            return rejectWithValue('Auth initialization failed');
        }
    },
);

const initialState: AuthState = {
    user: null,
    token: null,
    refreshToken: null,
    expiresIn: null,
    isAuthenticated: false,
    isInitialized: false,
};

const setTokensToStorage = (token: string, refreshToken?: string, expiresIn?: number) => {
    if (isClient && token) {
        try {
            localStorage.setItem('token', token);
            localStorage.setItem('lastLogin', new Date().toISOString());

            if (refreshToken) {
                localStorage.setItem('refreshToken', refreshToken);
            }

            if (expiresIn) {
                const expiryTime = Date.now() + expiresIn * 1000;
                localStorage.setItem('tokenExpiry', expiryTime.toString());
            }
        } catch (error) {
            console.error('Failed to save tokens to localStorage:', error);
        }
    }
};

const removeTokensFromStorage = () => {
    if (isClient) {
        try {
            localStorage.removeItem('token');
            localStorage.removeItem('refreshToken');
            localStorage.removeItem('lastLogin');
            localStorage.removeItem('tokenExpiry');
        } catch (error) {
            console.error('Failed to remove tokens from localStorage:', error);
        }
    }
};

export const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setCredentials: (state, action: PayloadAction<{ user: any; token: string; refreshToken?: string; expiresIn?: number }>) => {
            const { user, token, refreshToken, expiresIn } = action.payload;
            if (token) {
                state.token = token;
                state.isAuthenticated = true;
                setTokensToStorage(token, refreshToken, expiresIn);
            }
            if (refreshToken) {
                state.refreshToken = refreshToken;
            }
            if (expiresIn !== undefined) {
                state.expiresIn = expiresIn;
            }
            state.user = user;
        },
        logout: (state) => {
            state.user = null;
            state.token = null;
            state.refreshToken = null;
            state.expiresIn = null;
            state.isAuthenticated = false;
            removeTokensFromStorage();
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(initializeAuth.fulfilled, (state, action) => {
                if (action.payload) {
                    state.user = action.payload.user;
                    state.token = action.payload.token;
                    state.refreshToken = action.payload.refreshToken;
                    state.expiresIn = action.payload.expiresIn;
                    state.isAuthenticated = true;
                }
                state.isInitialized = true;
            })
            .addCase(initializeAuth.rejected, (state) => {
                state.isInitialized = true;
            })
            .addMatcher(authApi.endpoints.loginUser.matchFulfilled, (state, { payload }) => {
                if (payload?.accessToken) {
                    state.user = payload.user;
                    state.token = payload.accessToken;
                    state.refreshToken = payload.refreshToken || null;
                    state.expiresIn = payload.expiresIn || null;
                    state.isAuthenticated = true;
                    setTokensToStorage(payload.accessToken, payload.refreshToken, payload.expiresIn);
                }
            })
            .addMatcher(authApi.endpoints.refreshToken.matchFulfilled, (state, { payload }) => {
                if (payload?.accessToken) {
                    state.token = payload.accessToken;
                    if (payload.refreshToken) {
                        state.refreshToken = payload.refreshToken;
                    }
                    if (payload.expiresIn) {
                        state.expiresIn = payload.expiresIn;
                    }
                    setTokensToStorage(payload.accessToken, payload.refreshToken, payload.expiresIn);
                }
            })
            .addMatcher(authApi.endpoints.getUserDetails.matchFulfilled, (state, { payload }) => {
                if (payload) {
                    state.user = payload;
                }
            });
    },
});

export const { setCredentials, logout } = authSlice.actions;
export default authSlice.reducer;