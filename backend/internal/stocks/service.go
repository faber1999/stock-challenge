package stocks

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo        *repository
	apiClient   *apiClient
	syncTimeout time.Duration
	syncMaxPage int
}

func NewService(pool *pgxpool.Pool, apiURL, apiToken string, syncTimeout time.Duration, syncMaxPages int) *Service {
	if syncTimeout <= 0 {
		syncTimeout = 20 * time.Second
	}
	if syncMaxPages <= 0 {
		syncMaxPages = 50
	}

	return &Service{
		repo:        newRepository(pool),
		apiClient:   newAPIClient(apiURL, apiToken, syncTimeout),
		syncTimeout: syncTimeout,
		syncMaxPage: syncMaxPages,
	}
}

func (s *Service) Sync(ctx context.Context, limit int) (SyncResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.syncTimeout)
	defer cancel()

	limit = s.normalizeSyncLimit(limit)
	seenTokens := map[string]bool{}
	byTicker := map[string]Stock{}

	nextPage := ""
	pages := 0
	for i := 0; i < limit; i++ {
		resp, err := s.apiClient.FetchPage(ctx, nextPage)
		if err != nil {
			return SyncResult{}, fmt.Errorf("fetch page %d: %w", i+1, err)
		}
		pages++

		for _, raw := range resp.Items {
			stock, err := parseStock(raw)
			if err != nil {
				log.Printf("skip invalid stock on page %d: %v", i+1, err)
				continue
			}
			stock.RecommendScore = recommendationScore(stock)
			byTicker[stock.Ticker] = stock
		}

		if strings.TrimSpace(resp.NextPage) == "" {
			break
		}
		if seenTokens[resp.NextPage] {
			log.Printf("sync cycle detected for next_page=%s", resp.NextPage)
			break
		}

		seenTokens[resp.NextPage] = true
		nextPage = resp.NextPage
	}

	stocks := make([]Stock, 0, len(byTicker))
	for _, stock := range byTicker {
		stocks = append(stocks, stock)
	}
	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].Ticker < stocks[j].Ticker
	})

	if err := s.repo.ReplaceAll(ctx, stocks); err != nil {
		return SyncResult{}, err
	}

	return SyncResult{
		PagesProcessed: pages,
		StocksSaved:    len(stocks),
	}, nil
}

func (s *Service) List(ctx context.Context, params ListParams) (ListResult, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) GetByTicker(ctx context.Context, ticker string) (Stock, error) {
	stock, err := s.repo.GetByTicker(ctx, ticker)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Stock{}, ErrNotFound
		}
		return Stock{}, err
	}
	return stock, nil
}

func (s *Service) Recommend(ctx context.Context, limit int) ([]Recommendation, error) {
	if limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	items, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	recs := make([]Recommendation, 0, len(items))
	for _, st := range items {
		upside := upsidePct(st)
		score := st.RecommendScore
		recs = append(recs, Recommendation{
			Stock:     st,
			Score:     score,
			UpsidePct: upside,
		})
	}

	sort.Slice(recs, func(i, j int) bool {
		if recs[i].Score == recs[j].Score {
			return recs[i].UpsidePct > recs[j].UpsidePct
		}
		return recs[i].Score > recs[j].Score
	})

	if len(recs) > limit {
		recs = recs[:limit]
	}

	return recs, nil
}

func (s *Service) normalizeSyncLimit(limit int) int {
	if limit <= 0 {
		return 1
	}
	if limit > s.syncMaxPage {
		return s.syncMaxPage
	}
	return limit
}

func actionScore(action string) float64 {
	a := strings.ToLower(strings.TrimSpace(action))
	switch {
	case strings.Contains(a, "strong buy"):
		return 3.5
	case strings.Contains(a, "buy"):
		return 3.0
	case strings.Contains(a, "outperform"), strings.Contains(a, "overweight"):
		return 2.5
	case strings.Contains(a, "hold"), strings.Contains(a, "neutral"):
		return 1.0
	case strings.Contains(a, "sell"), strings.Contains(a, "underperform"), strings.Contains(a, "underweight"):
		return -1.5
	default:
		return 0.0
	}
}

func recommendationScore(s Stock) float64 {
	upside := upsidePct(s)
	return actionScore(s.Action) + ratingScore(s.RatingTo) + upside/20.0
}

func upsidePct(s Stock) float64 {
	if s.TargetFrom <= 0 {
		return 0
	}
	return ((s.TargetTo - s.TargetFrom) / s.TargetFrom) * 100
}

func ratingScore(rating string) float64 {
	switch strings.ToLower(strings.TrimSpace(rating)) {
	case "strong buy", "strong-buy":
		return 3.0
	case "buy", "outperform", "overweight":
		return 2.0
	case "hold", "neutral":
		return 0.5
	case "sell", "underperform", "underweight":
		return -1.5
	default:
		return 0
	}
}
