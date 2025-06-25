import React from 'react';
import { useRBAC, Resource, Action, Permission } from '@/lib/rbac';

interface RBACGuardProps {
  children: React.ReactNode;
  resource?: Resource;
  action?: Action;
  permission?: Permission;
  permissions?: Permission[];
  requireAll?: boolean;
  fallback?: React.ReactNode;
  loadingFallback?: React.ReactNode;
}

export const RBACGuard: React.FC<RBACGuardProps> = ({
  children,
  resource,
  action,
  permission,
  permissions,
  requireAll = false,
  fallback = null,
  loadingFallback = null,
}) => {
  const { 
    canAccessResource, 
    hasPermission, 
    hasAnyPermission, 
    hasAllPermissions, 
    isLoading 
  } = useRBAC();

  if (isLoading) {
    return <>{loadingFallback}</>;
  }

  let shouldRender = false;

  if (resource && action) {
    shouldRender = canAccessResource(resource, action);
  } else if (permission) {
    shouldRender = hasPermission(permission);
  } else if (permissions && permissions.length > 0) {
    shouldRender = requireAll 
      ? hasAllPermissions(permissions)
      : hasAnyPermission(permissions);
  } else {
    shouldRender = true;
  }

  return shouldRender ? <>{children}</> : <>{fallback}</>;
}; 