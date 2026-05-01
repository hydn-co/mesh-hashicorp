package options

import (
	"fmt"
	"strings"
)

func ValidateTerraformOptions(opts interface {
	GetHostname() string
	GetOrganization() string
}) error {
	if opts == nil {
		return fmt.Errorf("feature options are required but not provided")
	}
	if strings.TrimSpace(opts.GetHostname()) == "" {
		return fmt.Errorf("hostname is required in feature options")
	}
	if strings.TrimSpace(opts.GetOrganization()) == "" {
		return fmt.Errorf("organization is required in feature options")
	}
	return nil
}

func ValidateVaultOptions(opts interface{ GetAddress() string }) error {
	if opts == nil {
		return fmt.Errorf("feature options are required but not provided")
	}
	if strings.TrimSpace(opts.GetAddress()) == "" {
		return fmt.Errorf("address is required in feature options")
	}
	return nil
}
