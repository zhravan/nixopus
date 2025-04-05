export enum AUTHURLS {
  USER_LOGIN = 'v1/auth/login',
  USER_DETAILS = 'v1/user',
  REFRESH_TOKEN = 'v1/auth/refresh-token',
  REQUEST_PASSWORD_RESET = 'v1/auth/request-password-reset',
  RESET_PASSWORD = 'v1/auth/reset-password',
  SEND_VERIFICATION = 'v1/auth/send-verification-email',
  VERIFY_EMAIL = 'v1/auth/verify-email'
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
  UPDATE_PREFERENCES = 'v1/notification/preferences'
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
  RESTART_APPLICATION = 'v1/deploy/application/restart'
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
  DELETE_FILE = 'v1/files',
  CREATE_FILE = 'v1/files',
  DELETE_DIRECTORY = 'v1/files/directories',
  MOVE_FOLDER_FILES_RECURSIVELY_OR_RENAME = 'v1/file-manager/move-directory',
  COPY_FOLDER_FILES_RECURSIVELY = 'v1/files/copy',
  CALCULATE_DIRECTORY_SIZE = 'v1/files/directories/size',
  UPLOAD_FILE = 'v1/file-manager/upload'
}

export enum AUDITURLS {
  GET_RECENT_AUDIT_LOGS = 'v1/audit/logs'
}
