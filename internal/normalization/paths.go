package normalization

import (
	"fmt"
	"strings"
)

func NormalizeVaultMountPath(mountPath string) (string, error) {
	return normalizeSlashTrimmedRequiredValue("vault mount path", mountPath)
}

func NormalizeVaultSecretPath(secretPath string) (string, error) {
	return normalizeSlashTrimmedRequiredValue("vault secret path", secretPath)
}

func normalizeSlashTrimmedRequiredValue(fieldName string, value string) (string, error) {
	normalizedValue := strings.Trim(strings.TrimSpace(value), "/")
	if normalizedValue == "" {
		return "", fmt.Errorf("%s is required", fieldName)
	}

	return normalizedValue, nil
}
