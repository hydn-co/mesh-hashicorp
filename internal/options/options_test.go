package options

import (
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRegisterPolymorphicOptions(t *testing.T) {
	testkit.TestPolymorphicRegistrations(t, map[string]any{
		"mesh://hashicorp/collectors/terraform_account_entity_collector_options":       &TerraformAccountEntityCollectorOptions{},
		"mesh://hashicorp/collectors/terraform_team_entity_collector_options":          &TerraformTeamEntityCollectorOptions{},
		"mesh://hashicorp/collectors/terraform_workspace_entity_collector_options":     &TerraformWorkspaceEntityCollectorOptions{},
		"mesh://hashicorp/collectors/terraform_policy_entity_collector_options":        &TerraformPolicyEntityCollectorOptions{},
		"mesh://hashicorp/collectors/terraform_team_access_entity_collector_options":   &TerraformTeamAccessEntityCollectorOptions{},
		"mesh://hashicorp/collectors/terraform_audit_trail_activity_collector_options": &TerraformAuditTrailActivityCollectorOptions{},
		"mesh://hashicorp/collectors/vault_identity_entity_collector_options":          &VaultIdentityEntityCollectorOptions{},
		"mesh://hashicorp/collectors/vault_policy_entity_collector_options":            &VaultPolicyEntityCollectorOptions{},
		"mesh://hashicorp/collectors/vault_auth_method_entity_collector_options":       &VaultAuthMethodEntityCollectorOptions{},
		"mesh://hashicorp/actions/terraform_team_provision_action_options":             &TerraformTeamProvisionActionOptions{},
		"mesh://hashicorp/actions/terraform_workspace_provision_action_options":        &TerraformWorkspaceProvisionActionOptions{},
		"mesh://hashicorp/actions/terraform_team_membership_assign_action_options":     &TerraformTeamMembershipAssignActionOptions{},
		"mesh://hashicorp/actions/terraform_team_access_assign_action_options":         &TerraformTeamAccessAssignActionOptions{},
	})
}

func TestShouldRejectTerraformOptionsWhenRequiredFieldsMissing(t *testing.T) {
	err := ValidateTerraformOptions(&TerraformOptionsCore{})

	require.Error(t, err)
	assert.EqualError(t, err, "hostname is required in feature options")
}

func TestShouldRejectVaultOptionsWhenAddressMissing(t *testing.T) {
	err := ValidateVaultOptions(&VaultOptionsCore{})

	require.Error(t, err)
	assert.EqualError(t, err, "address is required in feature options")
}
