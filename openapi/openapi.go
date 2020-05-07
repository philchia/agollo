package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// OpenAPI contains api to manage configs
type OpenAPI interface {
	Envs() ([]*Env, error)
	Namespaces() ([]*NamespaceInfo, error)
	NamespaceInfo(namespaceName string) (*NamespaceInfo, error)
	CreateNamespace(namespaceName string, format string, public bool, comment string, dataChangeCreatedBy string) error
	GetLock(namespaceName string) (*Lock, error)
	AddConfig(namespaceName string, key, value string, comment string, dataChangeCreatedBy string) error
	UpdateConfig(namespaceName string, key, value string, comment string, dataChangeLastModifiedBy string) error
	DeleteConfig(namespaceName string, key, operator string) error
	Release(namnespaceName string, releaseTitle string, releaseComment string, releaseBy string) error
	GetRelease(namespaceName string) (*Release, error)
}

// New create an OpenAPI instance
func New(portal_address string, appid string, env string, cluster string, token string) OpenAPI {
	ret := &api{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		portalAddr: portal_address,
		appid:      appid,
		env:        env,
		cluster:    cluster,
		token:      token,
	}
	return ret
}

type api struct {
	client     *http.Client
	portalAddr string
	appid      string
	env        string
	cluster    string
	token      string
}

func (a *api) request(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", a.token)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var e Error
		if err := json.Unmarshal(bts, &e); err != nil {
			return nil, err
		}

		return nil, e
	}

	return bts, nil
}

// Envs get all env info
func (a *api) Envs() ([]*Env, error) {
	// http://{portal_address}/openapi/v1/apps/{appId}/envclusters
	url := fmt.Sprintf("%s/openapi/v1/apps/%s/envclusters", a.portalAddr, a.appid)
	bts, err := a.request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var envs []*Env
	if err := json.Unmarshal(bts, &envs); err != nil {
		return nil, err
	}

	return envs, nil
}

func (a *api) Namespaces() ([]*NamespaceInfo, error) {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces", a.portalAddr, a.env, a.appid, a.cluster)

	bts, err := a.request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var namespaces []*NamespaceInfo

	if err := json.Unmarshal(bts, &namespaces); err != nil {
		return nil, err
	}

	return namespaces, nil
}

func (a *api) NamespaceInfo(namespaceName string) (*NamespaceInfo, error) {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", a.portalAddr, a.env, a.appid, a.cluster, namespaceName)

	bts, err := a.request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var namespace NamespaceInfo

	if err := json.Unmarshal(bts, &namespace); err != nil {
		return nil, err
	}

	return &namespace, nil
}

func (a *api) CreateNamespace(namespaceName string, format string, public bool, comment string, dataChangeCreatedBy string) error {
	// http://{portal_address} /openapi/v1/apps/{appId}/appnamespaces
	url := fmt.Sprintf("%s/openapi/v1/apps/%s/appnamespaces", a.portalAddr, a.appid)
	params := map[string]interface{}{
		"name":                namespaceName,
		"appId":               a.appid,
		"format":              format,
		"isPublic":            public,
		"comment":             comment,
		"dataChangeCreatedBy": dataChangeCreatedBy,
	}

	bts, err := json.Marshal(params)
	if err != nil {
		return err
	}

	if _, err := a.request("POST", url, bytes.NewReader(bts)); err != nil {
		return err
	}

	return nil
}

func (a *api) GetLock(namespaceName string) (*Lock, error) {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/lock
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/lock",
		a.portalAddr, a.env, a.appid, a.cluster, namespaceName)

	bts, err := a.request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var lock Lock

	if err := json.Unmarshal(bts, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

func (a *api) AddConfig(namespaceName string, key, value string, comment string, dataChangeCreatedBy string) error {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items",
		a.portalAddr, a.env, a.appid, a.cluster, namespaceName)
	var params = map[string]interface{}{
		"key":                 key,
		"value":               value,
		"comment":             comment,
		"dataChangeCreatedBy": dataChangeCreatedBy,
	}

	bts, err := json.Marshal(&params)
	if err != nil {
		return err
	}

	if _, err := a.request("POST", url, bytes.NewReader(bts)); err != nil {
		return err
	}

	return nil
}

func (a *api) UpdateConfig(namespaceName string, key, value string, comment string, dataChangeLastModifiedBy string) error {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s",
		a.portalAddr, a.env, a.appid, a.cluster, namespaceName, key)
	var params = map[string]interface{}{
		"key":                      key,
		"value":                    value,
		"comment":                  comment,
		"dataChangeLastModifiedBy": dataChangeLastModifiedBy,
	}

	bts, err := json.Marshal(&params)
	if err != nil {
		return err
	}

	if _, err := a.request("PUT", url, bytes.NewReader(bts)); err != nil {
		return err
	}

	return nil
}

func (a *api) DeleteConfig(namespaceName string, key, operator string) error {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}?operator={operator}
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s?operator=%s",
		a.portalAddr, a.env, a.appid, a.cluster, namespaceName, key, operator,
	)

	_, err := a.request("DELETE", url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *api) Release(namnespaceName string, releaseTitle string, releaseComment string, releasedBy string) error {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/releases
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/releases",
		a.portalAddr, a.env, a.appid, a.cluster, namnespaceName,
	)
	params := map[string]interface{}{
		"releaseTitle":   releaseTitle,
		"releaseComment": releaseComment,
		"releasedBy":     releasedBy,
	}

	bts, err := json.Marshal(&params)
	if err != nil {
		return err
	}

	if _, err := a.request("POST", url, bytes.NewReader(bts)); err != nil {
		return err
	}
	return nil
}

func (a *api) GetRelease(namespaceName string) (*Release, error) {
	// http://{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/releases/latest
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/releases/latest",
		a.portalAddr, a.env, a.appid, a.cluster, namespaceName,
	)

	bts, err := a.request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var release Release
	if err := json.Unmarshal(bts, &release); err != nil {
		return nil, err
	}

	return &release, nil
}
