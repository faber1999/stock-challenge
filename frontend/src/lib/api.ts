import type {
  ListStocksParams,
  ListStocksResponse,
  RecommendationsResponse,
  Stock,
  SyncStocksResponse,
} from "../types/stocks";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL?.toString().trim();

async function parseJSON<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Request failed with status ${response.status}`);
  }
  return (await response.json()) as T;
}

export async function syncStocks(limit = 10): Promise<SyncStocksResponse> {
  const response = await fetch(`${API_BASE_URL}/stocks/sync`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ limit }),
  });
  return parseJSON<SyncStocksResponse>(response);
}

export async function listStocks(
  params: ListStocksParams,
): Promise<ListStocksResponse> {
  const query = new URLSearchParams({
    q: params.q,
    action: params.action,
    sort_by: params.sortBy,
    order: params.order,
    limit: String(params.limit),
    offset: String(params.offset),
  });
  const response = await fetch(`${API_BASE_URL}/stocks?${query.toString()}`);
  return parseJSON<ListStocksResponse>(response);
}

export async function getStock(ticker: string): Promise<Stock> {
  const response = await fetch(
    `${API_BASE_URL}/stocks/${encodeURIComponent(ticker)}`,
  );
  return parseJSON<Stock>(response);
}

export async function listRecommendations(
  limit = 5,
): Promise<RecommendationsResponse> {
  const response = await fetch(
    `${API_BASE_URL}/stocks/recommendations?limit=${limit}`,
  );
  return parseJSON<RecommendationsResponse>(response);
}
