export interface Domain {
  id: string;
  domain: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
  is_wildcard: boolean;
}

export interface Server {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
  is_primary: boolean;
}

export interface ServerAndDomain {
  id: string;
  server_id: string;
  domain_id: string;
  server: Server;
  domains: Domain[];
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}
