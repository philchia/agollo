package openapi

import "fmt"

type Env struct {
	Name     string   `json:"env"`
	Clusters []string `json:"clusters"`
}

type NamespaceInfo struct {
	AppID          string   `json:"appId"`
	ClusterName    string   `json:"clusterName"`
	NamespaceName  string   `json:"namespaceName"`
	Comment        string   `json:"comment"`
	Format         string   `json:"format"`
	IsPublic       bool     `json:"isPublic"`
	Items          []KVItem `json:"items"`
	CreateBy       string   `json:"dataChangeCreatedBy"`
	LastModifyBy   string   `json:"dataChangeLastModifiedBy"`
	CreateTime     string   `json:"dataChangeCreatedTime"`
	LastModifyTime string   `json:"dataChangeLastModifiedTime"`
}

type KVItem struct {
	Key            string `json:"key"`
	Value          string `json:"value"`
	CreateBy       string `json:"dataChangeCreatedBy"`
	LastModifyBy   string `json:"dataChangeLastModifiedBy"`
	CreateTime     string `json:"dataChangeCreatedTime"`
	LastModifyTime string `json:"dataChangeLastModifiedTime"`
}

type Lock struct {
	NamespaceName string `json:"namespaceName"`
	Locked        bool   `json:"isLocked"`
	LockedBy      string `json:"lockedBy"`
}

type Release struct {
	AppId                      string            `json:"appId"`
	ClusterName                string            `json:"clusterName"`
	NamespaceName              string            `json:"namespaceName"`
	Name                       string            `json:"name"`
	Configurations             map[string]string `json:"configurations"`
	Comment                    string            `json:"comment"`
	DataChangeCreatedBy        string            `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string            `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string            `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string            `json:"dataChangeLastModifiedTime"`
}

type Error struct {
	Msg  string  `json:"message"`
	Code float64 `json:"status"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code:%f, msg:%s", e.Code, e.Msg)
}
