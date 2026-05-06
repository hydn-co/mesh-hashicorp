package credentials

import (
	"encoding/json"

	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

func NormalizeToken(value string) (string, error) {
	return connectorutil.RequireTrimmedString("token", value)
}

func ExtractToken(raw json.RawMessage) (string, error) {
	return connectorutil.ExtractAPIKey(raw)
}

func GetBearerAuthorizationHeaderValue(value string) (string, error) {
	token, err := NormalizeToken(value)
	if err != nil {
		return "", err
	}

	return "Bearer " + token, nil
}
