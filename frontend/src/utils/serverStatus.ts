import type { Server } from '../types/api';

export function getServerStatus(server: Server): 'healthy' | 'stale' | 'offline' {
  if (server.status === 'online') {
    return 'healthy';
  }

  if (server.status === 'offline') {
    return 'offline';
  }

  return 'stale';
}
