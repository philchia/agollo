package openapi

import (
	"flag"
	"os"
	"testing"
)

var _api OpenAPI

func TestMain(m *testing.M) {
	var token = flag.String("token", "", "token")
	var portal = flag.String("portal", "", "portal")
	var appid = flag.String("appid", "", "app id")
	var env = flag.String("env", "", "env")
	var cluster = flag.String("cluster", "", "cluster")
	if *token == "" ||
		*portal == "" ||
		*appid == "" ||
		*env == "" ||
		*cluster == "" {
		os.Exit(0)
	}
	flag.Parse()
	_api = New(*portal, *appid, *env, *cluster, *token)

	os.Exit(m.Run())
}

func TestEnvs(t *testing.T) {

	envs, err := _api.Envs()
	if err != nil {
		t.Error(err)
	}
	for _, env := range envs {
		t.Logf("%s: %s", env.Name, env.Clusters)
	}
}

func TestNamespaces(t *testing.T) {
	namespaces, err := _api.Namespaces()
	if err != nil {
		t.Error(err)
	}
	for _, namespace := range namespaces {
		t.Logf("%s, %s, %s", namespace.ClusterName, namespace.NamespaceName, namespace.CreateBy)
	}
}

func TestNamespace(t *testing.T) {
	namespace, err := _api.NamespaceInfo("application")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s, %s, %s", namespace.ClusterName, namespace.NamespaceName, namespace.CreateBy)
	for _, item := range namespace.Items {
		t.Log(item.Key, item.Value)
	}
}

func TestCreateNamespace(t *testing.T) {
	if err := _api.CreateNamespace("testtest", "json", false, "", "zhaifei@hellobike.com"); err != nil {
		t.Error(err)
	}
}

func TestGetLock(t *testing.T) {

	lock, err := _api.GetLock("application")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("lock: %s,%v,%s", lock.NamespaceName, lock.Locked, lock.LockedBy)
}

func TestAddConfig(t *testing.T) {
	if err := _api.AddConfig("application", "testkey", "testvalue", "text", "zhaifei@hellobike.com"); err != nil {
		t.Error(err)
		return
	}
}

func TestUpdateConfig(t *testing.T) {
	if err := _api.UpdateConfig("application", "testkey", "testvalue1", "update", "zhaifei@hellobike.com"); err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteConfig(t *testing.T) {
	if err := _api.DeleteConfig("application", "testkey", "zhaifei@hellobike.com"); err != nil {
		t.Error(err)
		return
	}
}

func TestRelease(t *testing.T) {
	if err := _api.Release("application", "test release", "", "zhaifei@hellobike.com"); err != nil {
		t.Error(err)
		return
	}
}

func TestGetRelease(t *testing.T) {
	release, err := _api.GetRelease("application")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%#v", release)
}
