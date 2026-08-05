import { useEffect, useState } from 'react';
import { fetchServerMetrics } from '../services/api';
import type { Metric } from '../types/api';

export function useServerMetrics(serverId: number) {
  const [metrics, setMetrics] = useState<Metric[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    fetchServerMetrics(serverId)
      .then((data) => {
        if (!cancelled) {
          setMetrics(data);
        }
      })
      .catch((err: Error) => {
        if (!cancelled) {
          setError(err.message);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [serverId]);

  return { metrics, loading, error };
}
