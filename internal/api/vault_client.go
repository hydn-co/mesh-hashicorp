package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	vaultTokenHeader     = "X-Vault-Token"
	vaultNamespaceHeader = "X-Vault-Namespace"
	vaultRequestHeader   = "X-Vault-Request"
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
	return c.listKeys(ctx, "/v1/identity/entity/id")
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
	return c.listKeys(ctx, "/v1/identity/group/id")
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

func (c *VaultClient) listKeys(ctx context.Context, path string) ([]string, error) {
	query := url.Values{}
	query.Set("list", "true")

	body, statusCode, status, err := c.doGet(ctx, path, query)
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
	body, statusCode, status, err := c.doGet(ctx, path, query)
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

func (c *VaultClient) doGet(ctx context.Context, path string, query url.Values) ([]byte, int, string, error) {
	req, err := c.newGetRequest(ctx, path, query)
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

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	if err != nil {
		return nil, 0, "", fmt.Errorf("read vault response: %w", err)
	}

	return body, resp.StatusCode, resp.Status, nil
}

func (c *VaultClient) newGetRequest(ctx context.Context, path string, query url.Values) (*http.Request, error) {
	requestURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse vault base url: %w", err)
	}
	requestURL.Path = joinAPIPath(requestURL.Path, path)
	if query != nil {
		requestURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create vault request: %w", err)
	}
	req.Header.Set(vaultTokenHeader, c.Token)
	req.Header.Set(vaultRequestHeader, "true")
	if c.Namespace != "" {
		req.Header.Set(vaultNamespaceHeader, c.Namespace)
	}

	return req, nil
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
