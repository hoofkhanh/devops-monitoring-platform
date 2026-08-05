export interface Server {
  id: number;
  hostname: string;
  ip: string;
  os: string;
  status: string;
  last_seen?: string | null;
  created_at: string;
}

export interface Metric {
  id: number;
  server_id: number;
  cpu: number;
  memory: number;
  disk: number;
  timestamp: string;
}
