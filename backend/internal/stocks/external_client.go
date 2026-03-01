package stocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type apiClient struct {
	endpoint string
	token    string
	client   *http.Client
}

type apiListResponse struct {
	Items    []map[string]any `json:"items"`
	NextPage string           `json:"next_page"`
}

func newAPIClient(endpoint, token string, timeout time.Duration) *apiClient {
	return &apiClient{
		endpoint: strings.TrimSpace(endpoint),
		token:    strings.TrimSpace(token),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *apiClient) FetchPage(ctx context.Context, nextPage string) (apiListResponse, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return apiListResponse{}, fmt.Errorf("invalid stocks api url: %w", err)
	}

	q := u.Query()
	if strings.TrimSpace(nextPage) != "" {
		q.Set("next_page", nextPage)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return apiListResponse{}, fmt.Errorf("build stocks api request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return apiListResponse{}, fmt.Errorf("call stocks api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiListResponse{}, fmt.Errorf("read stocks api response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return apiListResponse{}, fmt.Errorf("stocks api status %d: %s", resp.StatusCode, string(body))
	}

	var out apiListResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return apiListResponse{}, fmt.Errorf("decode stocks api response: %w", err)
	}

	return out, nil
}
