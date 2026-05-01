package collectors

import (
	"fmt"
	"net/http"

	"github.com/hydn-co/mesh-hashicorp/internal/api"
)

func newTerraformClient(baseURL string, token string) (*api.Client, error) {
	client, err := api.NewClient(http.DefaultClient, baseURL, token)
	if err != nil {
		return nil, fmt.Errorf("build terraform client: %w", err)
	}

	return client, nil
}
