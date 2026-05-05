package payloads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func validateVaultSecretSetPayloadFields(mountPath string, secretPath string, data json.RawMessage) error {
	if strings.TrimSpace(mountPath) == "" {
		return fmt.Errorf("mount_path is required")
	}
	if strings.TrimSpace(secretPath) == "" {
		return fmt.Errorf("secret_path is required")
	}
	trimmedData := bytes.TrimSpace(data)
	if len(trimmedData) == 0 {
		return fmt.Errorf("data is required")
	}
	if trimmedData[0] != '{' {
		return fmt.Errorf("data must be a JSON object")
	}

	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(trimmedData, &decoded); err != nil {
		return fmt.Errorf("data must be a valid JSON object: %w", err)
	}

	return nil
}
