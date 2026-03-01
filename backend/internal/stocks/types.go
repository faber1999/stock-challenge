package stocks

import "time"

type Stock struct {
	ID             int64     `json:"id"`
	Ticker         string    `json:"ticker"`
	Company        string    `json:"company"`
	Brokerage      string    `json:"brokerage"`
	Action         string    `json:"action"`
	RatingFrom     string    `json:"rating_from"`
	RatingTo       string    `json:"rating_to"`
	TargetFrom     float64   `json:"target_from"`
	TargetTo       float64   `json:"target_to"`
	Currency       string    `json:"currency"`
	RecommendScore float64   `json:"recommend_score"`
	SyncedAt       time.Time `json:"synced_at"`
}

type ListParams struct {
	Query  string
	Action string
	SortBy string
	Order  string
	Limit  int
	Offset int
}

type ListResult struct {
	Items  []Stock `json:"items"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

type Recommendation struct {
	Stock
	Score     float64 `json:"score"`
	UpsidePct float64 `json:"upside_pct"`
}

type SyncResult struct {
	PagesProcessed int `json:"pages_processed"`
	StocksSaved    int `json:"stocks_saved"`
}
