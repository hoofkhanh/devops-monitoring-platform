import { useServerMetrics } from '../hooks/useServerMetrics';
import type { Server } from '../types/api';
import { ServerStatusBadge } from './ServerStatusBadge';

interface ServerCardProps {
  server: Server;
}

export function ServerCard({ server }: ServerCardProps) {
  const { metrics, loading, error } = useServerMetrics(server.id);
  const latestMetric = metrics[0];

  return (
    <article className="server-card" aria-label={`${server.hostname} card`}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>{server.hostname}</h2>
        <ServerStatusBadge server={server} />
      </div>
      <div className="server-meta">
        <div>{server.ip}</div>
        <div>{server.os}</div>
        <div>Last seen: {server.last_seen ? new Date(server.last_seen).toLocaleString() : 'n/a'}</div>
      </div>

      {loading ? (
        <div className="state">Loading metrics…</div>
      ) : error ? (
        <div className="state">Unable to load metrics.</div>
      ) : latestMetric ? (
        <div className="metric-list">
          <div className="metric-pill">
            <strong>{latestMetric.cpu.toFixed(1)}%</strong>
            <span>CPU</span>
          </div>
          <div className="metric-pill">
            <strong>{latestMetric.memory.toFixed(1)}%</strong>
            <span>Memory</span>
          </div>
          <div className="metric-pill">
            <strong>{latestMetric.disk.toFixed(1)}%</strong>
            <span>Disk</span>
          </div>
        </div>
      ) : (
        <div className="state">No metrics yet.</div>
      )}
    </article>
  );
}
