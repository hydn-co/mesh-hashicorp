package main

import (
	"log"
	"os"
	"strconv"

	"github.com/hydn-co/mesh-hashicorp/internal/actions"
	activitycollectors "github.com/hydn-co/mesh-hashicorp/internal/collectors/activity"
	entitycollectors "github.com/hydn-co/mesh-hashicorp/internal/collectors/entity"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

var experimentalFlagString string = os.Getenv("MESH_CONNECTOR_ALPHA_FUNCTIONS_ENABLED")

func main() {
	runner.Run(WithManifest())
}

func WithManifest() *runner.Manifest {
	experimentalFlag, err := strconv.ParseBool(experimentalFlagString)
	if err != nil {
		experimentalFlag = false
	}

	manifest := runner.CreateManifest(
		"mesh-hashicorp",
		"",
		"HashiCorp",
		"Mesh integration with HashiCorp",
	)

	if experimentalFlag {
		// Register terraform collectors
		manifest.MustRegisterFeature(
			"hashicorp_terraform_account_entity_collector",
			"Collect Terraform Accounts",
			"Collect Terraform organization users.",
			runner.FeatureSchedulable,
			runner.FeatureTypeCollector,
			new(options.TerraformAccountEntityCollectorOptions),
			(*connector.NoPayload)(nil),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(entitycollectors.NewTerraformAccountEntityCollector),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_team_entity_collector",
			"Collect Terraform Teams",
			"Collect Terraform teams and memberships.",
			runner.FeatureSchedulable,
			runner.FeatureTypeCollector,
			new(options.TerraformTeamEntityCollectorOptions),
			(*connector.NoPayload)(nil),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(entitycollectors.NewTerraformTeamEntityCollector),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_workspace_entity_collector",
			"Collect Terraform Workspaces",
			"Collect Terraform workspaces as applications.",
			runner.FeatureSchedulable,
			runner.FeatureTypeCollector,
			new(options.TerraformWorkspaceEntityCollectorOptions),
			(*connector.NoPayload)(nil),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(entitycollectors.NewTerraformWorkspaceEntityCollector),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_policy_entity_collector",
			"Collect Terraform Policies",
			"Collect Terraform policy sets and policies.",
			runner.FeatureSchedulable,
			runner.FeatureTypeCollector,
			new(options.TerraformPolicyEntityCollectorOptions),
			(*connector.NoPayload)(nil),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(entitycollectors.NewTerraformPolicyEntityCollector),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_team_access_entity_collector",
			"Collect Terraform Team Access",
			"Collect Terraform workspace access as permissions and group-permission links.",
			runner.FeatureSchedulable,
			runner.FeatureTypeCollector,
			new(options.TerraformTeamAccessEntityCollectorOptions),
			(*connector.NoPayload)(nil),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(entitycollectors.NewTerraformTeamAccessEntityCollector),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_audit_trail_activity_collector",
			"Collect Terraform Audit Activity",
			"Collect Terraform audit trail activity.",
			runner.FeatureSchedulable,
			runner.FeatureTypeCollector,
			new(options.TerraformAuditTrailActivityCollectorOptions),
			(*connector.NoPayload)(nil),
			runner.FeatureResumeBehaviorLastActivity,
			runner.APIKeyCredential,
			runner.Factory(activitycollectors.NewTerraformAuditTrailActivityCollector),
		)
	}

	// Register vault collectors
	manifest.MustRegisterFeature(
		"hashicorp_vault_identity_account_entity_collector",
		"Collect Vault Accounts",
		"Collect Vault identity entities as accounts.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.VaultIdentityAccountEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultIdentityAccountEntityCollector),
	)

	manifest.MustRegisterFeature(
		"hashicorp_vault_identity_group_entity_collector",
		"Collect Vault Groups",
		"Collect Vault identity groups and memberships.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.VaultIdentityGroupEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultIdentityGroupEntityCollector),
	)

	manifest.MustRegisterFeature(
		"hashicorp_vault_policy_entity_collector",
		"Collect Vault Policies",
		"Collect Vault policies.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.VaultPolicyEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultPolicyEntityCollector),
	)

	manifest.MustRegisterFeature(
		"hashicorp_vault_auth_method_entity_collector",
		"Collect Vault Auth Methods",
		"Collect Vault auth methods as applications.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.VaultAuthMethodEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultAuthMethodEntityCollector),
	)

	manifest.MustRegisterFeature(
		"hashicorp_vault_secret_entity_collector",
		"Collect Vault Secrets",
		"Collect Vault secret references from KV mounts without reading secret values.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.VaultSecretEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultSecretEntityCollector),
	)

	if experimentalFlag {
		// Register terraform actions
		manifest.MustRegisterFeature(
			"hashicorp_terraform_team_provision_action",
			"Provision Terraform Team",
			"Create a Terraform team.",
			runner.FeatureUnschedulable,
			runner.FeatureTypeAction,
			new(options.TerraformTeamProvisionActionOptions),
			new(payloads.TerraformTeamProvisionPayload),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(actions.NewTerraformTeamProvisionAction),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_workspace_provision_action",
			"Provision Terraform Workspace",
			"Create a Terraform workspace.",
			runner.FeatureUnschedulable,
			runner.FeatureTypeAction,
			new(options.TerraformWorkspaceProvisionActionOptions),
			new(payloads.TerraformWorkspaceProvisionPayload),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(actions.NewTerraformWorkspaceProvisionAction),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_team_membership_assign_action",
			"Add User To Terraform Team",
			"Assign a Terraform user to a team.",
			runner.FeatureUnschedulable,
			runner.FeatureTypeAction,
			new(options.TerraformTeamMembershipAssignActionOptions),
			new(payloads.TerraformTeamMembershipAssignPayload),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(actions.NewTerraformTeamMembershipAssignAction),
		)

		manifest.MustRegisterFeature(
			"hashicorp_terraform_team_access_assign_action",
			"Assign Terraform Team Access",
			"Assign a Terraform team permission to a workspace.",
			runner.FeatureUnschedulable,
			runner.FeatureTypeAction,
			new(options.TerraformTeamAccessAssignActionOptions),
			new(payloads.TerraformTeamAccessAssignPayload),
			runner.FeatureResumeBehaviorNone,
			runner.APIKeyCredential,
			runner.Factory(actions.NewTerraformTeamAccessAssignAction),
		)
	}

	// Register vault actions
	manifest.MustRegisterFeature(
		"hashicorp_vault_kv_v1_secret_set_action",
		"Upsert v1 vault secret",
		"Create or update a Vault KV v1 secret.",
		runner.FeatureUnschedulable,
		runner.FeatureTypeAction,
		new(options.VaultKVV1SecretSetActionOptions),
		new(payloads.VaultKVV1SecretSetPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewVaultKVV1SecretSetAction),
	)

	manifest.MustRegisterFeature(
		"hashicorp_vault_kv_v2_secret_set_action",
		"Upsert v2 vault secret",
		"Create or update a Vault KV v2 secret.",
		runner.FeatureUnschedulable,
		runner.FeatureTypeAction,
		new(options.VaultKVV2SecretSetActionOptions),
		new(payloads.VaultKVV2SecretSetPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewVaultKVV2SecretSetAction),
	)

	if err := manifest.Validate(); err != nil {
		log.Fatal(err)
	}

	return manifest
}
