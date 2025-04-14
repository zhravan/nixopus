export enum AUTHURLS {
  USER_LOGIN = 'v1/auth/login',
  USER_DETAILS = 'v1/user',
  REFRESH_TOKEN = 'v1/auth/refresh-token',
  REQUEST_PASSWORD_RESET = 'v1/auth/request-password-reset',
  RESET_PASSWORD = 'v1/auth/reset-password',
  SEND_VERIFICATION = 'v1/auth/send-verification-email',
  VERIFY_EMAIL = 'v1/auth/verify-email',
  LOGOUT = 'v1/auth/logout',
  SETUP_TWO_FACTOR = 'v1/auth/setup-2fa',
  VERIFY_TWO_FACTOR = 'v1/auth/verify-2fa',
  DISABLE_TWO_FACTOR = 'v1/auth/disable-2fa',
  TWO_FACTOR_LOGIN = 'v1/auth/2fa-login'
}

export enum USERURLS {
  USER_ORGANIZATIONS = 'v1/user/organizations',
  CREATE_ORGANIZATION = 'v1/organizations',
  ADD_USER_TO_ORGANIZATION = 'v1/organizations/add-user',
  REMOVE_USER_FROM_ORGANIZATION = 'v1/organizations/remove-user',
  UPDATE_USER_ROLE = 'v1/organizations/update-user-role',
  UPDATE_USER_NAME = 'v1/user/name',
  REQUEST_PASSWORD_RESET_LINK = 'v1/auth/request-password-reset',
  ORGANIZATION_USERS = 'v1/organizations/users',
  CREATE_USER = 'v1/auth/create-user',
  GET_RESOURCES = 'v1/organizations/resources'
}

export enum USER_NOTIFICATION_SETTINGS {
  ADD_SMTP = 'v1/notification/smtp',
  GET_SMTP = 'v1/notification/smtp',
  UPDATE_SMTP = 'v1/notification/smtp',
  DELETE_SMTP = 'v1/notification/smtp',

  GET_PREFERENCES = 'v1/notification/preferences',
  UPDATE_PREFERENCES = 'v1/notification/preferences',

  GET_WEBHOOK = 'v1/notification/webhook',
  CREATE_WEBHOOK = 'v1/notification/webhook',
  UPDATE_WEBHOOK = 'v1/notification/webhook',
  DELETE_WEBHOOK = 'v1/notification/webhook'
}

export enum DOMAIN_SETTINGS {
  ADD_DOMAIN = 'v1/domain',
  GET_DOMAINS = 'v1/domains',
  UPDATE_DOMAIN = 'v1/domain',
  DELETE_DOMAIN = 'v1/domain',
  GENERATE_RANDOM_SUBDOMAIN = 'v1/domain/generate'
}

export enum GITHUB_CONNECTOR {
  ADD_GITHUB_CONNECTOR = 'v1/github-connector',
  GET_GITHUB_CONNECTORS = 'v1/github-connector/all',
  UPDATE_GITHUB_CONNECTOR = 'v1/github-connector',
  DELETE_GITHUB_CONNECTOR = 'v1/github-connector',
  ALL_REPOSITORIES = 'v1/github-connector/repositories'
}

export enum DEPLOY {
  GET_APPLICATIONS = 'v1/deploy/applications',
  CREATE_APPLICATION = 'v1/deploy/application',
  GET_APPLICATION = 'v1/deploy/application',
  UPDATE_APPLICATION = 'v1/deploy/application',
  REDEPLOY_APPLICATION = 'v1/deploy/application/redeploy',
  DEPLOYMENT = 'v1/deploy/application/deployments',
  DELETE_APPLICATION = 'v1/deploy/application',
  ROLLBACK_APPLICATION = 'v1/deploy/application/rollback',
  RESTART_APPLICATION = 'v1/deploy/application/restart',
  GET_APPLICATION_LOGS = 'v1/deploy/application/logs/{application_id}',
  GET_DEPLOYMENT_LOGS = 'v1/deploy/application/deployments/{deployment_id}/logs',
  GET_APPLICATION_DEPLOYMENTS = 'v1/deploy/application/deployments'
}

export enum SOCKET_EVENTS {
  MONITOR_APPLICATION_DEPLOYMENT = 'monitor_application_deployment'
}

export enum SOCKET_ACTIONS {
  SUBSCRIBE = 'subscribe'
}

export enum FILEMANAGERURLS {
  LIST_FILES_AT_PATH = 'v1/file-manager',
  CREATE_DIRECTORY = 'v1/file-manager/create-directory',
  DELETE_DIRECTORY = 'v1/file-manager/delete-directory',
  MOVE_FOLDER_FILES_RECURSIVELY_OR_RENAME = 'v1/file-manager/move-directory',
  COPY_FOLDER_FILES_RECURSIVELY = 'v1/file-manager/copy-directory',
  CALCULATE_DIRECTORY_SIZE = 'v1/file-manager/directories/size',
  UPLOAD_FILE = 'v1/file-manager/upload'
}

export enum AUDITURLS {
  GET_RECENT_AUDIT_LOGS = 'v1/audit/logs'
}
