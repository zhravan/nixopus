import { getBasePath, getServerConfigBase } from '@/lib/base-path';

let config: any = null;
let configPromise: Promise<any> | null = null;

async function fetchConfig() {
  if (config) return config;
  if (!configPromise) {
    const base = getBasePath();
    const url =
      typeof window !== 'undefined'
        ? base
          ? `${base}/api/config`
          : '/api/config'
        : `${getServerConfigBase()}/api/config`;
    configPromise = fetch(url)
      .then((r) => r.json())
      .then((c) => {
        config = c;
        return c;
      });
  }
  return configPromise;
}

export async function getBaseUrl() {
  const { baseUrl } = await fetchConfig();
  return baseUrl;
}

export async function getWebsocketUrl() {
  const { websocketUrl } = await fetchConfig();
  return websocketUrl;
}

export async function getWebhookUrl() {
  const { webhookUrl } = await fetchConfig();
  return webhookUrl;
}

export async function getWebsiteDomain() {
  const { websiteDomain } = await fetchConfig();
  return websiteDomain;
}

export async function getPasswordLoginEnabled() {
  const { passwordLoginEnabled } = await fetchConfig();
  return passwordLoginEnabled;
}
