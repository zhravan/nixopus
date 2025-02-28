export interface User {
    id: string;
    username: string;
    email: string;
    avatar: string;
    createdAt: Date;
    updatedAt: Date;
    deletedAt: Date | null;
    isVerified: boolean;
    resetToken: string;
}

export interface AuthResponse {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    user: User;
}

export interface RefreshTokenPayload {
    refresh_token: string;
}

export interface LoginPayload {
    email: string;
    password: string;
}