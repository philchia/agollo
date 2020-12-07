package agollo

import (
	"testing"
)

func TestChangeType(t *testing.T) {
	var tps = []ChangeType{ADD, MODIFY, DELETE, ChangeType(-1)}
	var strs = []string{"ADD", "MODIFY", "DELETE", "UNKNOW"}
	for i, tp := range tps {
		if tp.String() != strs[i] {
			t.FailNow()
		}
	}
}

func TestMakeDeleteChange(t *testing.T) {
	change := makeDeleteChange("key", "val")
	t.Log(change.String())
	if change.ChangeType != DELETE || change.OldValue != "val" {
		t.FailNow()
	}
}

func TestMakeModifyChange(t *testing.T) {
	change := makeModifyChange("key", "old", "new")
	t.Log(change.String())
	if change.ChangeType != MODIFY || change.OldValue != "old" || change.NewValue != "new" {
		t.FailNow()
	}
}

func TestMakeAddChange(t *testing.T) {
	change := makeAddChange("key", "value")
	t.Log(change.String())
	if change.ChangeType != ADD || change.NewValue != "value" {
		t.FailNow()
	}
}

func TestChangeEvent_String(t *testing.T) {
	type fields struct {
		Namespace string
		Changes   map[string]*Change
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "case1",
			fields: fields{
				Namespace: "application",
				Changes: map[string]*Change{
					"key": {
						Key:        "key",
						ChangeType: ADD,
						NewValue:   "value",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ChangeEvent{
				Namespace: tt.fields.Namespace,
				Changes:   tt.fields.Changes,
			}
			got := e.String()
			t.Log(got)
		})
	}
}
