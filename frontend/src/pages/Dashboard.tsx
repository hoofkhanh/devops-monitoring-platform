import { ServerCard } from '../components/ServerCard';
import { useServers } from '../hooks/useServers';

export function DashboardPage() {
  const { servers, loading, error } = useServers();

  return (
    <main className="container">
      <header className="page-header">
        <div>
          <h1>DevOps Monitoring Dashboard</h1>
          <p>Live server overview with CPU, memory, and disk usage.</p>
        </div>
      </header>

      {loading ? (
        <p className="state">Loading servers…</p>
      ) : error ? (
        <p className="state">Unable to load server data.</p>
      ) : (
        <section className="card-grid" aria-label="Server list">
          {servers.map((server) => (
            <ServerCard key={server.id} server={server} />
          ))}
        </section>
      )}
    </main>
  );
}
