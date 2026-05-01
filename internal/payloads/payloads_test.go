package payloads

import (
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRegisterPolymorphicPayloads(t *testing.T) {
	testkit.TestPolymorphicRegistrations(t, map[string]any{
		"mesh://hashicorp/actions/terraform_team_provision_payload":         &TerraformTeamProvisionPayload{},
		"mesh://hashicorp/actions/terraform_workspace_provision_payload":    &TerraformWorkspaceProvisionPayload{},
		"mesh://hashicorp/actions/terraform_team_membership_assign_payload": &TerraformTeamMembershipAssignPayload{},
		"mesh://hashicorp/actions/terraform_team_access_assign_payload":     &TerraformTeamAccessAssignPayload{},
	})
}

func TestShouldRejectTerraformTeamProvisionPayloadWhenNameMissing(t *testing.T) {
	err := (&TerraformTeamProvisionPayload{}).Validate()

	require.Error(t, err)
	assert.EqualError(t, err, "name is required")
}

func TestShouldRejectTerraformTeamMembershipAssignPayloadWhenFieldsMissing(t *testing.T) {
	err := (&TerraformTeamMembershipAssignPayload{}).Validate()

	require.Error(t, err)
	assert.EqualError(t, err, "team_id is required")
}
