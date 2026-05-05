package payloads

import (
	"fmt"

	"github.com/hydn-co/mesh-hashicorp/internal/validation"
)

type PayloadValidator interface {
	Validate() error
}

func ValidatePayload(payload PayloadValidator, payloadName string) error {
	if validation.IsNil(payload) {
		return fmt.Errorf("%s is required but not provided", payloadName)
	}

	if err := payload.Validate(); err != nil {
		return fmt.Errorf("invalid %s: %w", payloadName, err)
	}

	return nil
}
