package stocks

import (
	"fmt"
	"strconv"
	"strings"
)

func parseStock(item map[string]any) (Stock, error) {
	ticker := strings.ToUpper(firstString(item, "ticker", "symbol"))
	if ticker == "" {
		return Stock{}, fmt.Errorf("missing ticker")
	}
	targetFrom, ok := firstNumber(item, "target_from", "targetFrom")
	if !ok {
		return Stock{}, fmt.Errorf("missing target_from")
	}
	targetTo, ok := firstNumber(item, "target_to", "targetTo")
	if !ok {
		return Stock{}, fmt.Errorf("missing target_to")
	}
	currency := strings.ToUpper(firstString(item, "currency"))
	if currency == "" {
		currency = "USD"
	}

	return Stock{
		Ticker:     ticker,
		Company:    firstString(item, "company", "name"),
		Brokerage:  firstString(item, "brokerage"),
		Action:     firstString(item, "action", "rating_action"),
		RatingFrom: firstString(item, "rating_from", "ratingFrom"),
		RatingTo:   firstString(item, "rating_to", "ratingTo"),
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Currency:   currency,
	}, nil
}

func firstString(item map[string]any, keys ...string) string {
	for _, key := range keys {
		raw, ok := item[key]
		if !ok || raw == nil {
			continue
		}
		switch v := raw.(type) {
		case string:
			s := strings.TrimSpace(v)
			if s != "" {
				return s
			}
		default:
			s := strings.TrimSpace(fmt.Sprintf("%v", v))
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

func firstNumber(item map[string]any, keys ...string) (float64, bool) {
	for _, key := range keys {
		raw, ok := item[key]
		if !ok || raw == nil {
			continue
		}

		if n, ok := toFloat64(raw); ok {
			return n, true
		}
	}
	return 0, false
}

func toFloat64(raw any) (float64, bool) {
	switch v := raw.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case string:
		clean := strings.TrimSpace(v)
		clean = strings.ReplaceAll(clean, ",", "")
		clean = strings.ReplaceAll(clean, "$", "")
		clean = strings.ReplaceAll(clean, "%", "")
		if clean == "" {
			return 0, false
		}
		n, err := strconv.ParseFloat(clean, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}
