package payloads

import (
	"fmt"
	"strings"
)

func validateVaultSecretSetPayloadFields(mountPath string, secretPath string, data map[string]any) error {
	if strings.TrimSpace(mountPath) == "" {
		return fmt.Errorf("mount_path is required")
	}
	if strings.TrimSpace(secretPath) == "" {
		return fmt.Errorf("secret_path is required")
	}
	if data == nil {
		return fmt.Errorf("data is required")
	}

	return nil
}
