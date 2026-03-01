package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"stockchallenge/backend/internal/stocks"
)

type Handler struct {
	pool         *pgxpool.Pool
	stockService *stocks.Service
}

func NewRouter(pool *pgxpool.Pool, stockService *stocks.Service, corsAllowedOrigins string) http.Handler {
	h := &Handler{
		pool:         pool,
		stockService: stockService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /db-time", h.dbTime)
	mux.HandleFunc("POST /stocks/sync", h.syncStocks)
	mux.HandleFunc("GET /stocks/recommendations", h.recommendStocks)
	mux.HandleFunc("GET /stocks/{ticker}", h.getStockByTicker)
	mux.HandleFunc("GET /stocks", h.listStocks)
	mux.HandleFunc("GET /swagger", h.swaggerUI)
	mux.HandleFunc("GET /swagger/", h.swaggerUI)
	mux.HandleFunc("GET /swagger/doc.json", h.swaggerDoc)

	return withCORS(mux, corsAllowedOrigins)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) dbTime(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	err := h.pool.QueryRow(r.Context(), "SELECT NOW()").Scan(&now)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to query database"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"now": now.UTC().Format(time.RFC3339)})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func (h *Handler) syncStocks(w http.ResponseWriter, r *http.Request) {
	limit, err := syncLimitFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.stockService.Sync(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) listStocks(w http.ResponseWriter, r *http.Request) {
	params := stocks.ListParams{
		Query:  strings.TrimSpace(r.URL.Query().Get("q")),
		Action: strings.TrimSpace(r.URL.Query().Get("action")),
		SortBy: strings.TrimSpace(r.URL.Query().Get("sort_by")),
		Order:  strings.TrimSpace(r.URL.Query().Get("order")),
		Limit:  queryIntOrDefault(r, "limit", 20),
		Offset: queryIntOrDefault(r, "offset", 0),
	}

	result, err := h.stockService.List(r.Context(), params)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list stocks"})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getStockByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := strings.TrimSpace(r.PathValue("ticker"))
	if ticker == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ticker is required"})
		return
	}

	stock, err := h.stockService.GetByTicker(r.Context(), ticker)
	if err != nil {
		if errors.Is(err, stocks.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "stock not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch stock"})
		return
	}

	writeJSON(w, http.StatusOK, stock)
}

func (h *Handler) recommendStocks(w http.ResponseWriter, r *http.Request) {
	limit := queryIntOrDefault(r, "limit", 5)
	result, err := h.stockService.Recommend(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to build recommendations"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": result,
		"total": len(result),
	})
}

func queryIntOrDefault(r *http.Request, key string, fallback int) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}

type syncRequest struct {
	Limit int `json:"limit"`
}

func syncLimitFromRequest(r *http.Request) (int, error) {
	limit := queryIntOrDefault(r, "limit", 10)

	if r.ContentLength == 0 {
		return limit, nil
	}

	var req syncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return 0, errors.New("invalid json body")
	}

	if req.Limit == 0 {
		return limit, nil
	}
	if req.Limit < 0 {
		return 0, errors.New("limit must be a positive integer")
	}

	return req.Limit, nil
}

func withCORS(next http.Handler, allowedOriginsRaw string) http.Handler {
	allowed := parseAllowedOrigins(allowedOriginsRaw)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if allowOrigin(origin, allowed) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			} else if containsWildcard(allowed) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseAllowedOrigins(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{"*"}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		out = append(out, origin)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func allowOrigin(origin string, allowed []string) bool {
	if containsWildcard(allowed) {
		return true
	}
	for _, item := range allowed {
		if strings.EqualFold(item, origin) {
			return true
		}
	}
	return false
}

func containsWildcard(allowed []string) bool {
	for _, item := range allowed {
		if item == "*" {
			return true
		}
	}
	return false
}
