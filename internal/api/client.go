package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
)

const maxErrorBodyBytes = 8 * 1024

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(httpClient *http.Client, baseURL string, token string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	normalizedBaseURL, err := normalizeBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	if _, err := credentials.GetBearerAuthorizationHeaderValue(token); err != nil {
		return nil, fmt.Errorf("validate terraform token: %w", err)
	}

	return &Client{
		BaseURL:    normalizedBaseURL,
		HTTPClient: httpClient,
		Token:      token,
	}, nil
}

func normalizeBaseURL(baseURL string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		return "", fmt.Errorf("terraform base url is required")
	}
	if strings.Contains(baseURL, "://") {
		parsed, err := url.Parse(baseURL)
		if err != nil {
			return "", fmt.Errorf("parse terraform base url: %w", err)
		}
		if parsed.Scheme == "" || parsed.Host == "" {
			return "", fmt.Errorf("terraform base url must include a scheme and host")
		}
		return parsed.String(), nil
	}
	parsed, err := url.Parse("https://" + baseURL)
	if err != nil {
		return "", fmt.Errorf("parse terraform base url: %w", err)
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("terraform base url must include a host")
	}
	return parsed.String(), nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	requestURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("parse terraform base url: %w", err)
	}
	requestURL.Path = strings.TrimRight(requestURL.Path, "/") + path
	requestURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return fmt.Errorf("create terraform request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.api+json")
	authorizationHeaderValue, err := credentials.GetBearerAuthorizationHeaderValue(c.Token)
	if err != nil {
		return fmt.Errorf("build authorization header: %w", err)
	}
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute terraform request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		if readErr != nil {
			return fmt.Errorf("terraform request failed with status %s: read error body: %w", resp.Status, readErr)
		}
		if len(body) == 0 {
			return fmt.Errorf("terraform request failed with status %s", resp.Status)
		}
		return fmt.Errorf("terraform request failed with status %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode terraform response: %w", err)
	}

	return nil
}
