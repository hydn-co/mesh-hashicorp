package credentials

import (
	"encoding/json"
	"fmt"
	"strings"
)

type apiKeyCredentials struct {
	APIKey string `json:"api_key"`
}

func NormalizeToken(value string) (string, error) {
	return normalizeRequiredValue("token", value)
}

func ExtractToken(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("api key credentials are required")
	}

	var credentials apiKeyCredentials
	if err := json.Unmarshal(raw, &credentials); err != nil {
		return "", fmt.Errorf("decode api key credentials: %w", err)
	}

	return normalizeRequiredValue("api_key", credentials.APIKey)
}

func GetBearerAuthorizationHeaderValue(value string) (string, error) {
	token, err := NormalizeToken(value)
	if err != nil {
		return "", err
	}

	return "Bearer " + token, nil
}

func normalizeRequiredValue(fieldName string, value string) (string, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return "", fmt.Errorf("%s is required", fieldName)
	}

	return trimmedValue, nil
}
