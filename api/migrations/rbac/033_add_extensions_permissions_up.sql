INSERT INTO permissions (id, name, description, resource) VALUES
(uuid_generate_v4(), 'create', 'Create extensions', 'extensions'),
(uuid_generate_v4(), 'read', 'Read extension information', 'extensions'),
(uuid_generate_v4(), 'update', 'Update extension information', 'extensions'),
(uuid_generate_v4(), 'delete', 'Delete extensions', 'extensions'),
(uuid_generate_v4(), 'install', 'Install extensions', 'extensions'),
(uuid_generate_v4(), 'uninstall', 'Uninstall extensions', 'extensions');

WITH admin_role AS (
    SELECT id FROM roles WHERE name = 'admin'
)
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT uuid_generate_v4(), admin_role.id, permissions.id
FROM admin_role, permissions
WHERE permissions.resource = 'extensions';

WITH viewer_role AS (
    SELECT id FROM roles WHERE name = 'viewer'
),
read_permissions AS (
    SELECT id FROM permissions 
    WHERE name = 'read' 
    AND resource = 'extensions'
)
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT uuid_generate_v4(), viewer_role.id, read_permissions.id
FROM viewer_role, read_permissions;

WITH member_role AS (
    SELECT id FROM roles WHERE name = 'member'
),
member_permissions AS (
    SELECT id FROM permissions 
    WHERE name IN ('read', 'install', 'uninstall') 
    AND resource = 'extensions'
)
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT uuid_generate_v4(), member_role.id, member_permissions.id
FROM member_role, member_permissions;
