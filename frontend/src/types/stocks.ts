export type SortBy =
  | "ticker"
  | "company"
  | "brokerage"
  | "action"
  | "target_from"
  | "target_to"
  | "recommend_score"
  | "synced_at";

export type SortOrder = "asc" | "desc";

export interface Stock {
  id: number;
  ticker: string;
  company: string;
  brokerage: string;
  action: string;
  rating_from: string;
  rating_to: string;
  target_from: number;
  target_to: number;
  currency: string;
  recommend_score: number;
  synced_at: string;
}

export interface ListStocksResponse {
  items: Stock[];
  total: number;
  limit: number;
  offset: number;
}

export interface RecommendationsResponse {
  items: Array<Stock & { score: number; upside_pct: number }>;
  total: number;
}

export interface SyncStocksResponse {
  pages_processed: number;
  stocks_saved: number;
}

export interface ListStocksParams {
  q: string;
  action: string;
  sortBy: SortBy;
  order: SortOrder;
  limit: number;
  offset: number;
}

