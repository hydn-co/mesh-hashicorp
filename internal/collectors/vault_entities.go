package collectors

import (
	"fmt"
	"strings"
	"time"

	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/entities"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/types"
)

func NewVaultAccount(entity api.VaultIdentityEntity) (*entities.Account, error) {
	accountRef := strings.TrimSpace(entity.ID)
	if accountRef == "" {
		return nil, fmt.Errorf("vault identity entity id is required")
	}

	accountName := strings.TrimSpace(entity.Name)
	account := entities.NewAccount()
	account.AccountRef = accountRef
	account.AccountType = types.AccountTypeUser
	account.Name = accountName
	account.DisplayName = accountName
	account.Enabled = !entity.Disabled
	account.CreatedAt = parseVaultTimestamp(entity.CreationTime)
	return account, nil
}

func NewVaultGroup(group api.VaultIdentityGroup) (*entities.Group, error) {
	groupRef := strings.TrimSpace(group.ID)
	if groupRef == "" {
		return nil, fmt.Errorf("vault identity group id is required")
	}

	entity := entities.NewGroup()
	entity.GroupRef = groupRef
	entity.Name = strings.TrimSpace(group.Name)
	entity.CreatedAt = parseVaultTimestamp(group.CreationTime)
	return entity, nil
}

func NewVaultGroupMember(groupRef string, accountRef string) (*entities.GroupMember, error) {
	trimmedGroupRef := strings.TrimSpace(groupRef)
	if trimmedGroupRef == "" {
		return nil, fmt.Errorf("vault group member group ref is required")
	}

	trimmedAccountRef := strings.TrimSpace(accountRef)
	if trimmedAccountRef == "" {
		return nil, fmt.Errorf("vault group member account ref is required")
	}

	groupMember := entities.NewGroupMember()
	groupMember.GroupRef = trimmedGroupRef
	groupMember.AccountRef = trimmedAccountRef
	return groupMember, nil
}

func NewVaultApplication(path string, authMethod api.VaultAuthMethod) (*entities.Application, error) {
	applicationRef := strings.TrimSpace(path)
	if applicationRef == "" {
		return nil, fmt.Errorf("vault auth method path is required")
	}

	application := entities.NewApplication()
	application.ApplicationRef = applicationRef
	application.Name = applicationRef
	application.Description = strings.TrimSpace(authMethod.Description)
	return application, nil
}

func NewVaultPolicy(policyName string) (*entities.Policy, error) {
	trimmedPolicyName := strings.TrimSpace(policyName)
	if trimmedPolicyName == "" {
		return nil, fmt.Errorf("vault policy name is required")
	}

	entity := entities.NewPolicy()
	entity.PolicyRef = trimmedPolicyName
	entity.Name = trimmedPolicyName
	return entity, nil
}

func NewVaultSecret(secret api.VaultSecret) (*entities.Secret, error) {
	secretRef := strings.TrimSpace(secret.Ref)
	if secretRef == "" {
		return nil, fmt.Errorf("vault secret ref is required")
	}

	entity := entities.NewSecret()
	entity.SecretRef = secretRef
	entity.Name = strings.TrimSpace(secret.Name)
	entity.Provider = strings.TrimSpace(secret.Provider)
	entity.Path = strings.TrimSpace(secret.Path)
	entity.Type = strings.TrimSpace(secret.Type)
	return entity, nil
}

func parseVaultTimestamp(value string) *time.Time {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	timestamp, err := time.Parse(time.RFC3339Nano, trimmed)
	if err != nil {
		return nil
	}

	return &timestamp
}
