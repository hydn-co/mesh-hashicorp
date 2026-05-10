package main

import (
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var terraformFeatureNames = []string{
	"hashicorp_terraform_account_entity_collector",
	"hashicorp_terraform_team_entity_collector",
	"hashicorp_terraform_workspace_entity_collector",
	"hashicorp_terraform_policy_entity_collector",
	"hashicorp_terraform_team_access_entity_collector",
	"hashicorp_terraform_audit_trail_activity_collector",
	"hashicorp_terraform_team_provision_action",
	"hashicorp_terraform_workspace_provision_action",
	"hashicorp_terraform_team_membership_assign_action",
	"hashicorp_terraform_team_access_assign_action",
}

var vaultFeatureNames = []string{
	"hashicorp_vault_identity_account_entity_collector",
	"hashicorp_vault_identity_group_entity_collector",
	"hashicorp_vault_policy_entity_collector",
	"hashicorp_vault_auth_method_entity_collector",
	"hashicorp_vault_secret_entity_collector",
	"hashicorp_vault_kv_v1_secret_set_action",
	"hashicorp_vault_kv_v2_secret_set_action",
}

func setExperimentalFlagString(t *testing.T, value string) {
	t.Helper()

	original := experimentalFlagString
	experimentalFlagString = value
	t.Cleanup(func() {
		experimentalFlagString = original
	})
}

func TestDescribe(t *testing.T) {
	testkit.InvokeDescribe(t, WithManifest())
}

func TestList(t *testing.T) {
	testkit.InvokeList(t, WithManifest())
}

func TestShouldOmitTerraformFeaturesWhenAlphaFunctionsDisabled(t *testing.T) {
	// Arrange
	setExperimentalFlagString(t, "false")

	// Act
	manifest := WithManifest()

	// Assert
	require.NotNil(t, manifest)
	for _, featureName := range terraformFeatureNames {
		assert.NotContains(t, manifest.Features, featureName)
	}
	for _, featureName := range vaultFeatureNames {
		assert.Contains(t, manifest.Features, featureName)
	}
}

func TestShouldIncludeTerraformFeaturesWhenAlphaFunctionsEnabled(t *testing.T) {
	// Arrange
	setExperimentalFlagString(t, "true")

	// Act
	manifest := WithManifest()

	// Assert
	require.NotNil(t, manifest)
	for _, featureName := range terraformFeatureNames {
		assert.Contains(t, manifest.Features, featureName)
	}
	for _, featureName := range vaultFeatureNames {
		assert.Contains(t, manifest.Features, featureName)
	}
}
