package stocks

import "testing"

func TestParseStock(t *testing.T) {
	item := map[string]any{
		"ticker":      "aapl",
		"company":     "Apple Inc.",
		"brokerage":   "Broker",
		"action":      "upgraded by",
		"rating_from": "Hold",
		"rating_to":   "Buy",
		"target_from": "120",
		"target_to":   "130.25",
	}

	stock, err := parseStock(item)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if stock.Ticker != "AAPL" {
		t.Fatalf("expected ticker AAPL, got %s", stock.Ticker)
	}
	if stock.RatingFrom != "Hold" {
		t.Fatalf("unexpected rating_from: %v", stock.RatingFrom)
	}
	if stock.TargetTo != 130.25 {
		t.Fatalf("unexpected target_to: %v", stock.TargetTo)
	}
}

func TestUpsidePct(t *testing.T) {
	upside := upsidePct(Stock{
		TargetFrom: 140.0,
		TargetTo:   160.0,
	})
	if upside < 14.2 || upside > 14.3 {
		t.Fatalf("unexpected upside %f", upside)
	}
}
