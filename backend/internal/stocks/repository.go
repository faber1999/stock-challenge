package stocks

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func newRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) ReplaceAll(ctx context.Context, stocks []Stock) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "DELETE FROM stocks"); err != nil {
		return fmt.Errorf("clear stocks: %w", err)
	}

	const insertSQL = `
INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, currency, recommend_score, synced_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, now())
`
	for _, s := range stocks {
		if _, err := tx.Exec(ctx, insertSQL,
			s.Ticker, s.Company, s.Brokerage, s.Action, s.RatingFrom, s.RatingTo, s.TargetFrom, s.TargetTo, s.Currency, s.RecommendScore,
		); err != nil {
			return fmt.Errorf("insert stock %s: %w", s.Ticker, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *repository) List(ctx context.Context, params ListParams) (ListResult, error) {
	limit, offset := normalizePageParams(params.Limit, params.Offset)
	sortBy := normalizeSortBy(params.SortBy)
	order := normalizeOrder(params.Order)

	where := make([]string, 0, 2)
	args := make([]any, 0, 4)
	argPos := 1

	if q := strings.TrimSpace(params.Query); q != "" {
		where = append(where, fmt.Sprintf("(ticker ILIKE $%d OR company ILIKE $%d OR brokerage ILIKE $%d)", argPos, argPos, argPos))
		args = append(args, "%"+q+"%")
		argPos++
	}
	if action := strings.TrimSpace(params.Action); action != "" {
		where = append(where, fmt.Sprintf("LOWER(action) = LOWER($%d)", argPos))
		args = append(args, action)
		argPos++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM stocks" + whereClause
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return ListResult{}, fmt.Errorf("count stocks: %w", err)
	}

	query := fmt.Sprintf(`
SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from::FLOAT8, target_to::FLOAT8, currency, recommend_score::FLOAT8, synced_at
FROM stocks
%s
ORDER BY %s %s, ticker ASC
LIMIT $%d OFFSET $%d
`, whereClause, sortBy, order, argPos, argPos+1)

	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return ListResult{}, fmt.Errorf("list stocks: %w", err)
	}
	defer rows.Close()

	items, err := scanStocks(rows)
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (r *repository) GetByTicker(ctx context.Context, ticker string) (Stock, error) {
	const query = `
SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from::FLOAT8, target_to::FLOAT8, currency, recommend_score::FLOAT8, synced_at
FROM stocks
WHERE ticker = UPPER($1)
`
	var out Stock
	err := r.pool.QueryRow(ctx, query, ticker).Scan(
		&out.ID, &out.Ticker, &out.Company, &out.Brokerage, &out.Action, &out.RatingFrom, &out.RatingTo, &out.TargetFrom, &out.TargetTo, &out.Currency, &out.RecommendScore, &out.SyncedAt,
	)
	if err != nil {
		return Stock{}, err
	}
	return out, nil
}

func (r *repository) ListAll(ctx context.Context) ([]Stock, error) {
	const query = `
SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from::FLOAT8, target_to::FLOAT8, currency, recommend_score::FLOAT8, synced_at
FROM stocks
`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list all stocks: %w", err)
	}
	defer rows.Close()

	return scanStocks(rows)
}

func scanStocks(rows pgx.Rows) ([]Stock, error) {
	out := make([]Stock, 0)
	for rows.Next() {
		var s Stock
		if err := rows.Scan(
			&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Currency, &s.RecommendScore, &s.SyncedAt,
		); err != nil {
			return nil, fmt.Errorf("scan stock row: %w", err)
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stocks rows: %w", err)
	}
	return out, nil
}

func normalizePageParams(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func normalizeSortBy(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "ticker":
		return "ticker"
	case "company":
		return "company"
	case "brokerage":
		return "brokerage"
	case "action":
		return "action"
	case "rating_from":
		return "rating_from"
	case "rating_to":
		return "rating_to"
	case "target_from":
		return "target_from"
	case "target_to":
		return "target_to"
	case "recommend_score":
		return "recommend_score"
	case "currency":
		return "currency"
	case "synced_at":
		return "synced_at"
	default:
		return "ticker"
	}
}

func normalizeOrder(order string) string {
	if strings.EqualFold(strings.TrimSpace(order), "desc") {
		return "DESC"
	}
	return "ASC"
}
