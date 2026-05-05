package api

import "encoding/json"

type vaultListResponse struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

type vaultErrorResponse struct {
	Errors []string `json:"errors"`
}

type vaultMountResponse struct {
	Data    VaultMount        `json:"data"`
	Type    string            `json:"type"`
	Options map[string]string `json:"options"`
}

type VaultMount struct {
	Type    string            `json:"type"`
	Options map[string]string `json:"options"`
}

type vaultKVV2WriteOptions struct {
	CAS int `json:"cas"`
}

type vaultKVV2WriteRequest struct {
	Options *vaultKVV2WriteOptions `json:"options,omitempty"`
	Data    json.RawMessage        `json:"data"`
}

type vaultIdentityEntityResponse struct {
	Data VaultIdentityEntity `json:"data"`
}

type VaultIdentityEntity struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Disabled          bool              `json:"disabled"`
	CreationTime      string            `json:"creation_time"`
	LastUpdateTime    string            `json:"last_update_time"`
	Metadata          map[string]string `json:"metadata"`
	Policies          []string          `json:"policies"`
	DirectGroupIDs    []string          `json:"direct_group_ids"`
	GroupIDs          []string          `json:"group_ids"`
	InheritedGroupIDs []string          `json:"inherited_group_ids"`
}

type vaultIdentityGroupResponse struct {
	Data VaultIdentityGroup `json:"data"`
}

type VaultIdentityGroup struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	CreationTime    string            `json:"creation_time"`
	LastUpdateTime  string            `json:"last_update_time"`
	Metadata        map[string]string `json:"metadata"`
	Policies        []string          `json:"policies"`
	MemberEntityIDs []string          `json:"member_entity_ids"`
	MemberGroupIDs  []string          `json:"member_group_ids"`
	ParentGroupIDs  []string          `json:"parent_group_ids"`
}

type vaultAuthMethodsResponse struct {
	Data map[string]VaultAuthMethod `json:"data"`
}

type VaultAuthMethod struct {
	Accessor              string                `json:"accessor"`
	Type                  string                `json:"type"`
	Description           string                `json:"description"`
	UUID                  string                `json:"uuid"`
	DeprecationStatus     string                `json:"deprecation_status"`
	Local                 bool                  `json:"local"`
	SealWrap              bool                  `json:"seal_wrap"`
	PluginVersion         string                `json:"plugin_version"`
	RunningPluginVersion  string                `json:"running_plugin_version"`
	RunningSHA256         string                `json:"running_sha256"`
	ExternalEntropyAccess bool                  `json:"external_entropy_access"`
	Options               map[string]string     `json:"options"`
	Config                VaultAuthMethodConfig `json:"config"`
}

type VaultAuthMethodConfig struct {
	DefaultLeaseTTL int    `json:"default_lease_ttl"`
	ForceNoCache    bool   `json:"force_no_cache"`
	MaxLeaseTTL     int    `json:"max_lease_ttl"`
	TokenType       string `json:"token_type"`
}

type vaultPolicyListResponse struct {
	Policies []string `json:"policies"`
}
