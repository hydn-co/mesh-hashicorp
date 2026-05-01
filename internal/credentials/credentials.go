package credentials

import (
	"encoding/json"
	"fmt"
	"strings"
)

type APITokenCredentials struct {
	APIKey string `json:"api_key,omitempty"`
	Token  string `json:"token,omitempty"`
}

func ParseAPIToken(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("api key credentials are required")
	}

	var credentials APITokenCredentials
	if err := json.Unmarshal(raw, &credentials); err != nil {
		return "", fmt.Errorf("decode api key credentials: %w", err)
	}

	if token := strings.TrimSpace(credentials.Token); token != "" {
		return token, nil
	}
	if apiKey := strings.TrimSpace(credentials.APIKey); apiKey != "" {
		return apiKey, nil
	}

	return "", fmt.Errorf("token or api_key is required")
}
