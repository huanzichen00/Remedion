package jev

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.typesafe.ai",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) systemOne(ctx context.Context, input SystemOneRequest) (SystemOneResponse, error) {
	var result SystemOneResponse

	body, err := json.Marshal(input)
	if err != nil {
		return result, fmt.Errorf("marshal Jev request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v1/systemone",
		bytes.NewReader(body),
	)
	if err != nil {
		return result, fmt.Errorf("create Jev request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("send Jev request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if readErr != nil {
			return result, fmt.Errorf("read Jev error response (HTTP %d): %w", resp.StatusCode, readErr)
		}

		message := strings.TrimSpace(string(responseBody))

		return result, &APIError{StatusCode: resp.StatusCode, Message: message}
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode Jev response: %w", err)
	}

	return result, nil
}
