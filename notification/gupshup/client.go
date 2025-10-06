package gupshup

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultBaseURL = "https://api.gupshup.io/sm/api/v1/msg"

// Config represents the configuration needed to interact with the Gupshup SMS API.
type Config struct {
	BaseURL    string
	APIKey     string
	Source     string
	HTTPClient *http.Client
}

// Client can send SMS messages using the Gupshup API.
type Client struct {
	baseURL string
	apiKey  string
	source  string
	client  *http.Client
}

// NewClient creates a new Gupshup client using the provided configuration.
func NewClient(cfg Config) *Client {
	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  cfg.APIKey,
		source:  cfg.Source,
		client:  httpClient,
	}
}

// SendSMS sends a single SMS message using the Gupshup API and returns the provider message ID when available.
func (c *Client) SendSMS(ctx context.Context, destination, message string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("gupshup client is not configured")
	}

	values := url.Values{}
	values.Set("channel", "SMS")
	values.Set("source", c.source)
	values.Set("destination", destination)
	values.Set("message", message)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, strings.NewReader(values.Encode()))
	if err != nil {
		return "", fmt.Errorf("build gupshup request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.apiKey != "" {
		req.Header.Set("apikey", c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send gupshup request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read gupshup response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("gupshup request failed: %s", strings.TrimSpace(string(body)))
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err == nil {
		if msgID, ok := lookupString(data, "messageId"); ok {
			return msgID, nil
		}
		if respObj, ok := lookupMap(data, "response"); ok {
			if msgID, ok := lookupString(respObj, "msgId"); ok {
				return msgID, nil
			}
			if msgID, ok := lookupString(respObj, "messageId"); ok {
				return msgID, nil
			}
		}
	}

	return "", nil
}

func lookupMap(m map[string]any, key string) (map[string]any, bool) {
	v, ok := m[key]
	if !ok {
		return nil, false
	}
	mv, ok := v.(map[string]any)
	return mv, ok
}

func lookupString(m map[string]any, key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	return s, true
}
