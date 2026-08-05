const apiURL = import.meta.env.VITE_API_URL;

if (!apiURL) {
  throw new Error('VITE_API_URL must be configured');
}

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${apiURL}${path}`);
  if (!response.ok) {
    const payload = await response.json().catch(() => ({}));
    throw new Error((payload as { error?: string }).error || `Request failed with ${response.status}`);
  }
  return response.json() as Promise<T>;
}

export async function fetchServers() {
  return request<import('../types/api').Server[]>('/servers');
}

export async function fetchServerMetrics(serverId: number) {
  return request<import('../types/api').Metric[]>(`/servers/${serverId}/metrics`);
}
