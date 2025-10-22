import { useEffect, useMemo, useState } from 'react';
import Session from 'supertokens-web-js/recipe/session';
import { UserRoleClaim, PermissionClaim } from 'supertokens-web-js/recipe/userroles';

export type Resource =
  | 'organization'
  | 'user'
  | 'role'
  | 'permission'
  | 'deploy'
  | 'file-manager'
  | 'dashboard'
  | 'settings'
  | 'audit'
  | 'notification'
  | 'domain'
  | 'feature-flags'
  | 'github-connector'
  | 'container'
  | 'terminal'
  | 'update'
  | 'extensions';

export type Action = 'create' | 'read' | 'update' | 'delete';

export type Permission = `${Resource}:${Action}`;

export const useRBAC = () => {
  const [roles, setRoles] = useState<string[] | undefined>(undefined);
  const [permissions, setPermissions] = useState<string[] | undefined>(undefined);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    let isMounted = true;
    (async () => {
      try {
        const hasSession = await Session.doesSessionExist();
        if (!hasSession) {
          if (isMounted) {
            setRoles(undefined);
            setPermissions(undefined);
            setIsLoading(false);
          }
          return;
        }

        const [sessionRoles, sessionPerms] = await Promise.all([
          Session.getClaimValue({ claim: UserRoleClaim }),
          Session.getClaimValue({ claim: PermissionClaim })
        ]);

        if (isMounted) {
          setRoles(sessionRoles ?? undefined);
          setPermissions(sessionPerms ?? undefined);
          setIsLoading(false);
        }
      } catch (_err) {
        if (isMounted) {
          setRoles(undefined);
          setPermissions(undefined);
          setIsLoading(false);
        }
      }
    })();
    return () => {
      isMounted = false;
    };
  }, []);

  const isAdmin = useMemo(() => {
    if (!Array.isArray(roles)) return false;
    return roles.some((role) => {
      // Strip orgid_ prefix to get base role name
      if (role.startsWith('orgid_')) {
        const lastUnderscore = role.lastIndexOf('_');
        if (lastUnderscore !== -1 && lastUnderscore < role.length - 1) {
          const baseRole = role.substring(lastUnderscore + 1);
          return baseRole === 'admin';
        }
      }
      return role === 'admin';
    });
  }, [roles]);

  const canAccessResource = (resource: Resource, action: Action): boolean => {
    if (isAdmin) return true;
    if (!permissions) return false;
    const permissionString: Permission = `${resource}:${action}`;
    return permissions.includes(permissionString);
  };

  const hasPermissionCheck = (permission: Permission): boolean => {
    if (isAdmin) return true;
    if (!permissions) return false;
    return permissions.includes(permission);
  };

  const hasAnyPermission = (permissions: Permission[]): boolean => {
    return permissions.some(hasPermissionCheck);
  };

  const hasAllPermissions = (permissions: Permission[]): boolean => {
    return permissions.every(hasPermissionCheck);
  };

  return {
    canAccessResource,
    hasPermission: hasPermissionCheck,
    hasAnyPermission,
    hasAllPermissions,
    isLoading,
    roles,
    permissions,
    isAdmin
  };
};
