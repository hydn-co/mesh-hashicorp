package main

import (
	"log"

	"github.com/hydn-co/mesh-hashicorp/internal/actions"
	activitycollectors "github.com/hydn-co/mesh-hashicorp/internal/collectors/activity"
	entitycollectors "github.com/hydn-co/mesh-hashicorp/internal/collectors/entity"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func main() {
	runner.Run(WithManifest())
}

func WithManifest() *runner.Manifest {
	manifest := runner.CreateManifest(
		"mesh-hashicorp",
		"",
		"HashiCorp",
		"Mesh integration with HashiCorp",
	)

	// Register collectors
	if err := manifest.RegisterFeature(
		"hashicorp_terraform_account_entity_collector",
		"Collect HCP Terraform Accounts",
		"Collect HCP Terraform organization users.",
		true,
		runner.FeatureTypeCollector,
		new(options.TerraformAccountEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewTerraformAccountEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_team_entity_collector",
		"Collect HCP Terraform Teams",
		"Collect HCP Terraform teams and memberships.",
		true,
		runner.FeatureTypeCollector,
		new(options.TerraformTeamEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewTerraformTeamEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_workspace_entity_collector",
		"Collect HCP Terraform Workspaces",
		"Collect HCP Terraform workspaces as applications.",
		true,
		runner.FeatureTypeCollector,
		new(options.TerraformWorkspaceEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewTerraformWorkspaceEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_policy_entity_collector",
		"Collect HCP Terraform Policies",
		"Collect HCP Terraform policy sets and policies.",
		true,
		runner.FeatureTypeCollector,
		new(options.TerraformPolicyEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewTerraformPolicyEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_team_access_entity_collector",
		"Collect HCP Terraform Team Access",
		"Collect HCP Terraform workspace access as permissions and group-permission links.",
		true,
		runner.FeatureTypeCollector,
		new(options.TerraformTeamAccessEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewTerraformTeamAccessEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_audit_trail_activity_collector",
		"Collect HCP Terraform Audit Activity",
		"Collect HCP Terraform audit trail activity.",
		true,
		runner.FeatureTypeCollector,
		new(options.TerraformAuditTrailActivityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorLastActivity,
		runner.APIKeyCredential,
		runner.Factory(activitycollectors.NewTerraformAuditTrailActivityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_vault_identity_entity_collector",
		"Collect Vault Identity Entities",
		"Collect Vault identity entities, aliases, groups, and memberships.",
		true,
		runner.FeatureTypeCollector,
		new(options.VaultIdentityEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultIdentityEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_vault_policy_entity_collector",
		"Collect Vault Policies",
		"Collect Vault policies.",
		true,
		runner.FeatureTypeCollector,
		new(options.VaultPolicyEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultPolicyEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_vault_auth_method_entity_collector",
		"Collect Vault Auth Methods",
		"Collect Vault auth methods as applications.",
		true,
		runner.FeatureTypeCollector,
		new(options.VaultAuthMethodEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultAuthMethodEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_vault_secret_entity_collector",
		"Collect Vault Secrets",
		"Collect Vault secret references from KV mounts without reading secret values.",
		true,
		runner.FeatureTypeCollector,
		new(options.VaultSecretEntityCollectorOptions),
		nil,
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(entitycollectors.NewVaultSecretEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	// Register actions
	if err := manifest.RegisterFeature(
		"hashicorp_terraform_team_provision_action",
		"Provision HCP Terraform Team",
		"Create an HCP Terraform team.",
		false,
		runner.FeatureTypeAction,
		new(options.TerraformTeamProvisionActionOptions),
		new(payloads.TerraformTeamProvisionPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewTerraformTeamProvisionAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_workspace_provision_action",
		"Provision HCP Terraform Workspace",
		"Create an HCP Terraform workspace.",
		false,
		runner.FeatureTypeAction,
		new(options.TerraformWorkspaceProvisionActionOptions),
		new(payloads.TerraformWorkspaceProvisionPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewTerraformWorkspaceProvisionAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_team_membership_assign_action",
		"Assign User To HCP Terraform Team",
		"Assign an HCP Terraform user to a team.",
		false,
		runner.FeatureTypeAction,
		new(options.TerraformTeamMembershipAssignActionOptions),
		new(payloads.TerraformTeamMembershipAssignPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewTerraformTeamMembershipAssignAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_terraform_team_access_assign_action",
		"Assign HCP Terraform Team Access",
		"Assign an HCP Terraform team permission to a workspace.",
		false,
		runner.FeatureTypeAction,
		new(options.TerraformTeamAccessAssignActionOptions),
		new(payloads.TerraformTeamAccessAssignPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewTerraformTeamAccessAssignAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_vault_kv_v1_secret_set_action",
		"Set Vault KV v1 Secret",
		"Create or update a Vault KV v1 secret.",
		false,
		runner.FeatureTypeAction,
		new(options.VaultKVV1SecretSetActionOptions),
		new(payloads.VaultKVV1SecretSetPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewVaultKVV1SecretSetAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"hashicorp_vault_kv_v2_secret_set_action",
		"Set Vault KV v2 Secret",
		"Create or update a Vault KV v2 secret.",
		false,
		runner.FeatureTypeAction,
		new(options.VaultKVV2SecretSetActionOptions),
		new(payloads.VaultKVV2SecretSetPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewVaultKVV2SecretSetAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.Validate(); err != nil {
		log.Fatal(err)
	}

	return manifest
}
