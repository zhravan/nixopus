let config: any = null;

async function fetchConfig() {
  if (!config) {
    const response = await fetch('/api/config');
    config = await response.json();
  }
  return config;
}

export async function getBaseUrl() {
  const { baseUrl } = await fetchConfig();
  return baseUrl;
}

export async function getWebsocketUrl() {
  const { websocketUrl } = await fetchConfig();
  return websocketUrl;
}

export const BASE_URL = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:8080/api';
export const WEBSOCKET_URL = process.env.NEXT_PUBLIC_WEBSOCKET_URL || 'ws://localhost:8080/ws';
