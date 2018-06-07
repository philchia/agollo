package agollo

import "testing"

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
	if change.ChangeType != DELETE || change.OldValue != "val" {
		t.FailNow()
	}
}

func TestMakeModifyChange(t *testing.T) {
	change := makeModifyChange("key", "old", "new")
	if change.ChangeType != MODIFY || change.OldValue != "old" || change.NewValue != "new" {
		t.FailNow()
	}
}

func TestMakeAddChange(t *testing.T) {
	change := makeAddChange("key", "value")
	if change.ChangeType != ADD || change.NewValue != "value" {
		t.FailNow()
	}
}
