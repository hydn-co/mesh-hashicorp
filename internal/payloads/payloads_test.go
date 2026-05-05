package payloads

import (
	"encoding/json"
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
		"mesh://hashicorp/actions/vault_kv_v1_secret_set_payload":           &VaultKVV1SecretSetPayload{},
		"mesh://hashicorp/actions/vault_kv_v2_secret_set_payload":           &VaultKVV2SecretSetPayload{},
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

func TestShouldRejectVaultKVV1SecretSetPayloadWhenMountPathMissing(t *testing.T) {
	err := (&VaultKVV1SecretSetPayload{SecretPath: "app/config", Data: json.RawMessage(`{"foo":"bar"}`)}).Validate()

	require.Error(t, err)
	assert.EqualError(t, err, "mount_path is required")
}

func TestShouldRejectVaultKVV1SecretSetPayloadWhenDataIsJSONString(t *testing.T) {
	err := (&VaultKVV1SecretSetPayload{
		MountPath:  "secret",
		SecretPath: "app/config",
		Data:       json.RawMessage(`"{\"foo\":\"bar\"}"`),
	}).Validate()

	require.Error(t, err)
	assert.EqualError(t, err, "data must be a JSON object")
}

func TestShouldRejectVaultKVV2SecretSetPayloadWhenCASNegative(t *testing.T) {
	cas := -1
	err := (&VaultKVV2SecretSetPayload{
		MountPath:  "secret",
		SecretPath: "app/config",
		Data:       json.RawMessage(`{"foo":"bar"}`),
		CAS:        &cas,
	}).Validate()

	require.Error(t, err)
	assert.EqualError(t, err, "cas cannot be negative")
}
