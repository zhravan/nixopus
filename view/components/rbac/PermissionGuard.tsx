import React from 'react';
import { RBACGuard } from './RBACGuard';
import { AccessDenied } from './AccessDenied';
import { Action, Permission, Resource } from '@/lib/rbac';

interface PermissionGuardProps {
    children: React.ReactNode;
    permission: Permission;
    permissions?: Permission[];
    resource?: Resource;
    action?: Action;
    fallback?: React.ReactNode;
    loadingFallback?: React.ReactNode;
}

export const PermissionGuard: React.FC<PermissionGuardProps> = ({
    children,
    permission,
    fallback,
    loadingFallback,
}) => (
    <RBACGuard
        permission={permission}
        fallback={fallback || <AccessDenied />}
        loadingFallback={loadingFallback}
    >
        {children}
    </RBACGuard>
);

export const AnyPermissionGuard: React.FC<PermissionGuardProps> = ({
    children,
    permissions,
    fallback,
    loadingFallback,
}) => (
    <RBACGuard
        permissions={permissions}
        requireAll={false}
        fallback={fallback || <AccessDenied />}
        loadingFallback={loadingFallback}
    >
        {children}
    </RBACGuard>
);

export const AllPermissionsGuard: React.FC<PermissionGuardProps> = ({
    children,
    permissions,
    fallback,
    loadingFallback,
}) => (
    <RBACGuard
        permissions={permissions}
        requireAll={true}
        fallback={fallback || <AccessDenied />}
        loadingFallback={loadingFallback}
    >
        {children}
    </RBACGuard>
);

export const ResourceGuard: React.FC<PermissionGuardProps> = ({
    children,
    resource,
    action,
    fallback,
    loadingFallback,
}) => (
    <RBACGuard
        resource={resource}
        action={action}
        fallback={fallback || <AccessDenied />}
        loadingFallback={loadingFallback}
    >
        {children}
    </RBACGuard>
); 