package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/hydn-co/mesh-hashicorp/internal/normalization"
)

const (
	vaultTokenHeader     = "X-Vault-Token"
	vaultNamespaceHeader = "X-Vault-Namespace"
	vaultRequestHeader   = "X-Vault-Request"
	vaultKVType          = "kv"
	vaultKVVersion1      = "1"
	vaultKVVersion2      = "2"
)

type VaultClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	Namespace  string
}

func NewVaultClient(httpClient *http.Client, address string, namespace string, token string) (*VaultClient, error) {
	baseClient, err := NewClient(httpClient, "vault", address, token)
	if err != nil {
		return nil, err
	}

	return &VaultClient{
		BaseURL:    baseClient.BaseURL,
		HTTPClient: baseClient.HTTPClient,
		Token:      baseClient.Token,
		Namespace:  strings.TrimSpace(namespace),
	}, nil
}

func (c *VaultClient) ListIdentityEntityIDs(ctx context.Context) ([]string, error) {
	return c.listKeysWithGetQuery(ctx, "/v1/identity/entity/id")
}

func (c *VaultClient) GetIdentityEntity(ctx context.Context, entityID string) (VaultIdentityEntity, error) {
	response := vaultIdentityEntityResponse{}
	if err := c.get(
		ctx,
		"/v1/identity/entity/id/"+url.PathEscape(strings.TrimSpace(entityID)),
		nil,
		&response,
	); err != nil {
		return VaultIdentityEntity{}, err
	}

	return response.Data, nil
}

func (c *VaultClient) ListIdentityGroupIDs(ctx context.Context) ([]string, error) {
	return c.listKeysWithGetQuery(ctx, "/v1/identity/group/id")
}

func (c *VaultClient) GetIdentityGroup(ctx context.Context, groupID string) (VaultIdentityGroup, error) {
	response := vaultIdentityGroupResponse{}
	if err := c.get(
		ctx,
		"/v1/identity/group/id/"+url.PathEscape(strings.TrimSpace(groupID)),
		nil,
		&response,
	); err != nil {
		return VaultIdentityGroup{}, err
	}

	return response.Data, nil
}

func (c *VaultClient) ListAuthMethods(ctx context.Context) (map[string]VaultAuthMethod, error) {
	response := vaultAuthMethodsResponse{}
	if err := c.get(ctx, "/v1/sys/auth", nil, &response); err != nil {
		return nil, err
	}
	if response.Data == nil {
		return map[string]VaultAuthMethod{}, nil
	}

	return response.Data, nil
}

func (c *VaultClient) ListPolicyNames(ctx context.Context) ([]string, error) {
	response := vaultPolicyListResponse{}
	if err := c.get(ctx, "/v1/sys/policy", nil, &response); err != nil {
		return nil, err
	}

	return response.Policies, nil
}

func (c *VaultClient) GetMount(ctx context.Context, mountPath string) (VaultMount, error) {
	normalizedMountPath, err := normalization.NormalizeVaultMountPath(mountPath)
	if err != nil {
		return VaultMount{}, err
	}

	response := vaultMountResponse{}
	if err := c.get(ctx, "/v1/sys/mounts/"+escapeVaultPath(normalizedMountPath), nil, &response); err != nil {
		return VaultMount{}, err
	}

	if response.Data.Type != "" || response.Data.Options != nil {
		return response.Data, nil
	}

	return VaultMount{
		Type:    response.Type,
		Options: response.Options,
	}, nil
}

func (c *VaultClient) ListMounts(ctx context.Context) (map[string]VaultMount, error) {
	response := vaultMountsResponse{}
	if err := c.get(ctx, "/v1/sys/mounts", nil, &response); err != nil {
		return nil, err
	}
	if response.Data == nil {
		return map[string]VaultMount{}, nil
	}

	mounts := make(map[string]VaultMount, len(response.Data))
	for mountPath, mount := range response.Data {
		normalizedMountPath, err := normalization.NormalizeVaultMountPath(mountPath)
		if err != nil {
			return nil, fmt.Errorf("normalize vault mount %q: %w", mountPath, err)
		}
		mounts[normalizedMountPath] = mount
	}

	return mounts, nil
}

func (c *VaultClient) ListKVSecrets(ctx context.Context, mountPath string) ([]VaultSecret, error) {
	normalizedMountPath, err := normalization.NormalizeVaultMountPath(mountPath)
	if err != nil {
		return nil, err
	}

	mount, err := c.GetMount(ctx, normalizedMountPath)
	if err != nil {
		return nil, fmt.Errorf("get vault mount %s: %w", normalizedMountPath, err)
	}
	if mount.Type != vaultKVType {
		return nil, fmt.Errorf("vault mount %s is not a kv secrets engine", normalizedMountPath)
	}

	version, err := resolveVaultKVVersion(mount)
	if err != nil {
		return nil, fmt.Errorf("resolve kv version for mount %s: %w", normalizedMountPath, err)
	}

	switch version {
	case vaultKVVersion1:
		return c.listKVV1Secrets(ctx, normalizedMountPath, "")
	case vaultKVVersion2:
		return c.listKVV2Secrets(ctx, normalizedMountPath, "")
	default:
		return nil, fmt.Errorf("unsupported vault kv version %s", version)
	}
}

func (c *VaultClient) SetKVV1Secret(
	ctx context.Context,
	mountPath string,
	secretPath string,
	data json.RawMessage,
) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("vault secret data is required")
	}

	normalizedMountPath, err := c.requireKVVersion(ctx, mountPath, vaultKVVersion1)
	if err != nil {
		return err
	}
	normalizedSecretPath, err := normalization.NormalizeVaultSecretPath(secretPath)
	if err != nil {
		return err
	}

	if err := c.postJSON(
		ctx,
		"/v1/"+escapeVaultPath(normalizedMountPath)+"/"+escapeVaultPath(normalizedSecretPath),
		data,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (c *VaultClient) SetKVV2Secret(
	ctx context.Context,
	mountPath string,
	secretPath string,
	data json.RawMessage,
	cas *int,
) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("vault secret data is required")
	}

	normalizedMountPath, err := c.requireKVVersion(ctx, mountPath, vaultKVVersion2)
	if err != nil {
		return err
	}
	normalizedSecretPath, err := normalization.NormalizeVaultSecretPath(secretPath)
	if err != nil {
		return err
	}

	request := vaultKVV2WriteRequest{Data: data}
	if cas != nil {
		request.Options = &vaultKVV2WriteOptions{CAS: *cas}
	}
	if err := c.postJSON(
		ctx,
		"/v1/"+escapeVaultPath(normalizedMountPath)+"/data/"+escapeVaultPath(normalizedSecretPath),
		request,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (c *VaultClient) listKeysWithGetQuery(ctx context.Context, path string) ([]string, error) {
	query := url.Values{}
	query.Set("list", "true")

	return c.listKeys(ctx, http.MethodGet, path, query)
}

func (c *VaultClient) listKeysWithListMethod(ctx context.Context, path string) ([]string, error) {
	return c.listKeys(ctx, "LIST", path, nil)
}

func (c *VaultClient) listKeys(ctx context.Context, method string, path string, query url.Values) ([]string, error) {
	body, statusCode, status, err := c.doRequest(ctx, method, path, query, nil, "")
	if err != nil {
		return nil, err
	}
	if statusCode == http.StatusNotFound && vaultListWasEmpty(body) {
		return nil, nil
	}
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("vault request failed with status %s%s", status, formatVaultErrorBody(body))
	}
	if len(body) == 0 {
		return nil, nil
	}

	response := vaultListResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("decode vault response: %w", err)
	}

	return response.Data.Keys, nil
}

func (c *VaultClient) get(ctx context.Context, path string, query url.Values, out any) error {
	body, statusCode, status, err := c.doRequest(ctx, http.MethodGet, path, query, nil, "")
	if err != nil {
		return err
	}
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("vault request failed with status %s%s", status, formatVaultErrorBody(body))
	}
	if out == nil || statusCode == http.StatusNoContent || len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode vault response: %w", err)
	}

	return nil
}

func (c *VaultClient) postJSON(ctx context.Context, path string, payload any, out any) error {
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode vault request: %w", err)
	}

	body, statusCode, status, err := c.doRequest(
		ctx,
		http.MethodPost,
		path,
		nil,
		bytes.NewReader(encodedPayload),
		"application/json",
	)
	if err != nil {
		return err
	}
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("vault request failed with status %s%s", status, formatVaultErrorBody(body))
	}
	if out == nil || statusCode == http.StatusNoContent || len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode vault response: %w", err)
	}

	return nil
}

func (c *VaultClient) doRequest(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	requestBody io.Reader,
	contentType string,
) ([]byte, int, string, error) {
	req, err := c.newRequest(ctx, method, path, query, requestBody, contentType)
	if err != nil {
		return nil, 0, "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, "", fmt.Errorf("execute vault request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	if err != nil {
		return nil, 0, "", fmt.Errorf("read vault response: %w", err)
	}

	return responseBody, resp.StatusCode, resp.Status, nil
}

func (c *VaultClient) newRequest(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	body io.Reader,
	contentType string,
) (*http.Request, error) {
	requestURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse vault base url: %w", err)
	}
	requestURL.Path = joinAPIPath(requestURL.Path, path)
	if query != nil {
		requestURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create vault request: %w", err)
	}
	req.Header.Set(vaultTokenHeader, c.Token)
	req.Header.Set(vaultRequestHeader, "true")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if c.Namespace != "" {
		req.Header.Set(vaultNamespaceHeader, c.Namespace)
	}

	return req, nil
}

func escapeVaultPath(path string) string {
	segments := strings.Split(path, "/")
	for index, segment := range segments {
		segments[index] = url.PathEscape(segment)
	}

	return strings.Join(segments, "/")
}

func (c *VaultClient) listKVV1Secrets(ctx context.Context, mountPath string, prefix string) ([]VaultSecret, error) {
	keys, err := c.listKeysWithListMethod(ctx, buildVaultKVV1ListPath(mountPath, prefix))
	if err != nil {
		return nil, err
	}
	return c.expandVaultKeys(ctx, mountPath, prefix, vaultKVVersion1, keys)
}

func (c *VaultClient) listKVV2Secrets(ctx context.Context, mountPath string, prefix string) ([]VaultSecret, error) {
	keys, err := c.listKeysWithListMethod(ctx, buildVaultKVV2ListPath(mountPath, prefix))
	if err != nil {
		return nil, err
	}
	return c.expandVaultKeys(ctx, mountPath, prefix, vaultKVVersion2, keys)
}

func (c *VaultClient) expandVaultKeys(
	ctx context.Context,
	mountPath string,
	prefix string,
	version string,
	keys []string,
) ([]VaultSecret, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	sortedKeys := append([]string(nil), keys...)
	sort.Strings(sortedKeys)

	secrets := make([]VaultSecret, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		joinedPath, err := joinVaultSecretPath(prefix, key)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(strings.TrimSpace(key), "/") {
			var nested []VaultSecret
			switch version {
			case vaultKVVersion1:
				nested, err = c.listKVV1Secrets(ctx, mountPath, joinedPath)
			case vaultKVVersion2:
				nested, err = c.listKVV2Secrets(ctx, mountPath, joinedPath)
			default:
				return nil, fmt.Errorf("unsupported vault kv version %s", version)
			}
			if err != nil {
				return nil, err
			}
			secrets = append(secrets, nested...)
			continue
		}

		secrets = append(secrets, newVaultSecret(mountPath, joinedPath, version))
	}

	return secrets, nil
}

func buildVaultKVV1ListPath(mountPath string, prefix string) string {
	path := "/v1/" + escapeVaultPath(mountPath)
	if strings.TrimSpace(prefix) != "" {
		path += "/" + escapeVaultPath(prefix)
	}
	return path
}

func buildVaultKVV2ListPath(mountPath string, prefix string) string {
	path := "/v1/" + escapeVaultPath(mountPath) + "/metadata"
	if strings.TrimSpace(prefix) != "" {
		path += "/" + escapeVaultPath(prefix)
	}
	return path
}

func joinVaultSecretPath(prefix string, key string) (string, error) {
	trimmedKey := strings.Trim(strings.TrimSpace(key), "/")
	if trimmedKey == "" {
		return "", fmt.Errorf("vault secret list returned empty key")
	}

	if strings.TrimSpace(prefix) == "" {
		return normalization.NormalizeVaultSecretPath(trimmedKey)
	}

	return normalization.NormalizeVaultSecretPath(prefix + "/" + trimmedKey)
}

func newVaultSecret(mountPath string, secretPath string, version string) VaultSecret {
	fullPath := mountPath + "/" + secretPath
	return VaultSecret{
		Ref:      fullPath,
		Name:     path.Base(secretPath),
		Provider: "HashiCorp Vault",
		Path:     secretPath,
		Type:     "kv-v" + version,
	}
}

func (c *VaultClient) requireKVVersion(ctx context.Context, mountPath string, expectedVersion string) (string, error) {
	normalizedMountPath, err := normalization.NormalizeVaultMountPath(mountPath)
	if err != nil {
		return "", err
	}

	mount, err := c.GetMount(ctx, normalizedMountPath)
	if err != nil {
		return "", fmt.Errorf("get vault mount %s: %w", normalizedMountPath, err)
	}
	if mount.Type != vaultKVType {
		return "", fmt.Errorf("vault mount %s is not a kv secrets engine", normalizedMountPath)
	}

	kvVersion, err := resolveVaultKVVersion(mount)
	if err != nil {
		return "", fmt.Errorf("resolve kv version for mount %s: %w", normalizedMountPath, err)
	}
	if kvVersion != expectedVersion {
		return "", fmt.Errorf("vault mount %s is kv v%s, not kv v%s", normalizedMountPath, kvVersion, expectedVersion)
	}

	return normalizedMountPath, nil
}

func resolveVaultKVVersion(mount VaultMount) (string, error) {
	if mount.Type != vaultKVType {
		return "", fmt.Errorf("mount type %q is not supported", mount.Type)
	}

	version := strings.TrimSpace(mount.Options["version"])
	switch version {
	case "", vaultKVVersion1:
		return vaultKVVersion1, nil
	case vaultKVVersion2:
		return vaultKVVersion2, nil
	default:
		return "", fmt.Errorf("unsupported vault kv version %q", version)
	}
}

func joinAPIPath(basePath string, apiPath string) string {
	trimmedBasePath := strings.TrimRight(basePath, "/")
	if strings.HasSuffix(trimmedBasePath, "/v1") && strings.HasPrefix(apiPath, "/v1/") {
		return trimmedBasePath + strings.TrimPrefix(apiPath, "/v1")
	}
	return trimmedBasePath + apiPath
}

func vaultListWasEmpty(body []byte) bool {
	if len(body) == 0 {
		return true
	}

	response := vaultErrorResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return false
	}

	return len(response.Errors) == 0
}

func formatVaultErrorBody(body []byte) string {
	trimmedBody := strings.TrimSpace(string(body))
	if trimmedBody == "" {
		return ""
	}

	response := vaultErrorResponse{}
	if err := json.Unmarshal(body, &response); err == nil && len(response.Errors) > 0 {
		return ": " + strings.Join(response.Errors, "; ")
	}

	return ": " + trimmedBody
}
