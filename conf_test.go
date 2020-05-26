package agollo

import "testing"

func TestNewConf(t *testing.T) {
	var tcs = []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "fakename",
			wantErr: true,
		},
		{
			name:    "./LICENSE",
			wantErr: true,
		},
		{
			name:    "./testdata/app.properties",
			wantErr: false,
		},
	}

	for _, tc := range tcs {
		if _, err := NewConf(tc.name); (err == nil) == tc.wantErr {
			t.FailNow()
		}
	}
}

func TestConfNormalize(t *testing.T) {
	cases := []struct {
		namespaces []string
	}{
		{
			namespaces: []string{"a", "b"},
		},
	}

	for _, c := range cases {
		conf := &Conf{NameSpaceNames: c.namespaces}
		conf.normalize()
		if !strIn(conf.NameSpaceNames, defaultNamespace) {
			t.Fatal("application should be in conf's namespaces")
		}
	}
}
