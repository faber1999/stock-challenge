package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const createStocksTableSQL = `
CREATE TABLE IF NOT EXISTS stocks (
  id INT8 NOT NULL DEFAULT unique_rowid(),
  ticker STRING NOT NULL,
  company STRING NOT NULL DEFAULT '',
  brokerage STRING NOT NULL DEFAULT '',
  action STRING NOT NULL DEFAULT '',
  rating_from STRING NOT NULL DEFAULT '',
  rating_to STRING NOT NULL DEFAULT '',
  target_from DECIMAL(18,4) NOT NULL DEFAULT 0,
  target_to DECIMAL(18,4) NOT NULL DEFAULT 0,
  currency STRING NOT NULL DEFAULT 'USD',
  recommend_score DECIMAL(18,4) NOT NULL DEFAULT 0,
  synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT stocks_pkey PRIMARY KEY (id),
  CONSTRAINT stocks_ticker_unique UNIQUE (ticker)
);
`

// Migrate creates/updates the database schema required by the backend.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	// 1) Ensure base table exists.
	if _, err := pool.Exec(ctx, createStocksTableSQL); err != nil {
		return fmt.Errorf("create stocks table: %w", err)
	}

	// 2) Additive changes for existing schemas.
	additiveStatements := []string{
		"ALTER TABLE IF EXISTS stocks ADD COLUMN IF NOT EXISTS id INT8 NOT NULL DEFAULT unique_rowid()",
		"ALTER TABLE IF EXISTS stocks ADD COLUMN IF NOT EXISTS currency STRING NOT NULL DEFAULT 'USD'",
		"ALTER TABLE IF EXISTS stocks ADD COLUMN IF NOT EXISTS recommend_score DECIMAL(18,4) NOT NULL DEFAULT 0",
		"ALTER TABLE IF EXISTS stocks ADD COLUMN IF NOT EXISTS synced_at TIMESTAMPTZ NOT NULL DEFAULT now()",
	}
	for _, stmt := range additiveStatements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("apply migration statement %q: %w", stmt, err)
		}
	}

	// 3) If legacy schema has incompatible rating types, recreate table.
	ok, err := hasStringRatingColumns(ctx, pool)
	if err != nil {
		return err
	}
	if !ok {
		recreateStatements := []string{
			"DROP TABLE IF EXISTS stocks",
			createStocksTableSQL,
		}
		for _, stmt := range recreateStatements {
			if _, err := pool.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("recreate stocks table with compatible schema: %w", err)
			}
		}
	}

	// 4) Indexes.
	indexStatements := []string{
		"CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks (company)",
		"CREATE INDEX IF NOT EXISTS idx_stocks_action ON stocks (action)",
		"CREATE INDEX IF NOT EXISTS idx_stocks_target_to ON stocks (target_to)",
		"CREATE INDEX IF NOT EXISTS idx_stocks_currency ON stocks (currency)",
		"CREATE INDEX IF NOT EXISTS idx_stocks_recommend_score ON stocks (recommend_score DESC)",
	}
	for _, stmt := range indexStatements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("create index with statement %q: %w", stmt, err)
		}
	}

	return nil
}

func hasStringRatingColumns(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	const query = `
SELECT column_name, data_type
FROM information_schema.columns
WHERE table_schema = current_schema()
  AND table_name = 'stocks'
  AND column_name IN ('rating_from', 'rating_to')
`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return false, fmt.Errorf("check rating column types: %w", err)
	}
	defer rows.Close()

	types := map[string]string{}
	for rows.Next() {
		var columnName string
		var dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return false, fmt.Errorf("scan rating column type: %w", err)
		}
		types[columnName] = strings.ToUpper(strings.TrimSpace(dataType))
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate rating column types: %w", err)
	}

	// If columns are missing or not STRING, consider incompatible.
	if types["rating_from"] != "STRING" || types["rating_to"] != "STRING" {
		return false, nil
	}
	return true, nil
}
