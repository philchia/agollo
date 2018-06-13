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
			name:    "./testdata/" + defaultConfName,
			wantErr: false,
		},
	}

	for _, tc := range tcs {
		if _, err := NewConf(tc.name); (err == nil) == tc.wantErr {
			t.FailNow()
		}
	}
}
