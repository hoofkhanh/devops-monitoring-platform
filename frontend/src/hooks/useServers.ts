import { useEffect, useState } from 'react';
import { fetchServers } from '../services/api';
import type { Server } from '../types/api';

export function useServers() {
  const [servers, setServers] = useState<Server[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    fetchServers()
      .then((data) => {
        if (!cancelled) {
          setServers(data);
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
  }, []);

  return { servers, loading, error };
}
