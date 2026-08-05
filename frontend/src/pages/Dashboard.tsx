import { useMetrics } from '../hooks/useMetrics';

export function DashboardPage() {
  const { metrics, loading, error } = useMetrics();
  const latestMetric = metrics[0];

  return (
    <main className="container">
      <header className="page-header">
        <div>
          <h1>DevOps Monitoring Dashboard</h1>
          <p>CPU, memory, and disk usage for this computer.</p>
        </div>
      </header>

      {loading ? (
        <p className="state">Loading metrics…</p>
      ) : error ? (
        <p className="state">Unable to load metrics.</p>
      ) : (
        <section className="metric-panel" aria-label="Latest metrics">
          {latestMetric ? (
            <>
              <div className="metric-list">
                <div className="metric-pill"><strong>{latestMetric.cpu.toFixed(1)}%</strong><span>CPU</span></div>
                <div className="metric-pill"><strong>{latestMetric.memory.toFixed(1)}%</strong><span>Memory</span></div>
                <div className="metric-pill"><strong>{latestMetric.disk.toFixed(1)}%</strong><span>Disk</span></div>
              </div>
              <p className="metric-timestamp">Updated {new Date(latestMetric.timestamp).toLocaleString()}</p>
            </>
          ) : <p className="state">No metrics yet.</p>}
        </section>
      )}
    </main>
  );
}
