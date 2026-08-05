import { getServerStatus } from '../utils/serverStatus';
import type { Server } from '../types/api';

interface ServerStatusBadgeProps {
  server: Server;
}

export function ServerStatusBadge({ server }: ServerStatusBadgeProps) {
  const status = getServerStatus(server);

  return <span className={`badge ${status}`}>{status}</span>;
}
