import { MastraClient } from '@mastra/client-js';
import { getBasePath } from '@/lib/base-path';

const AGENT_URL = process.env.NEXT_PUBLIC_AGENT_URL || '';

export const AGENT_ID = 'deploy-agent';

export function isAgentConfigured(): boolean {
  return Boolean(AGENT_URL);
}

export function createAgentClient(authHeaders: Record<string, string> = {}): MastraClient {
  const base = getBasePath();
  return new MastraClient({
    baseUrl: base ? `${base}/api/agent/` : '/api/agent/',
    headers: authHeaders
  });
}
