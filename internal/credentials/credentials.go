package credentials

import (
	"encoding/json"
	"fmt"
	"strings"
)

type apiKeyCredentials struct {
	APIKey string `json:"api_key"`
}

func ExtractToken(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("api key credentials are required")
	}

	var credentials apiKeyCredentials
	if err := json.Unmarshal(raw, &credentials); err != nil {
		return "", fmt.Errorf("decode api key credentials: %w", err)
	}

	token := strings.TrimSpace(credentials.APIKey)
	if token == "" {
		return "", fmt.Errorf("api_key is required")
	}

	return token, nil
}

func GetBearerAuthorizationHeaderValue(value string) (string, error) {
	token := strings.TrimSpace(value)
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	return "Bearer " + token, nil
}
