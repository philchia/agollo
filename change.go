package agollo

import (
	"bytes"
	"fmt"
)

// ChangeType for a key
type ChangeType int

const (
	// ADD a new value
	ADD ChangeType = iota
	// MODIFY a old value
	MODIFY
	// DELETE ...
	DELETE
)

func (c ChangeType) String() string {
	switch c {
	case ADD:
		return "ADD"
	case MODIFY:
		return "MODIFY"
	case DELETE:
		return "DELETE"
	}

	return "UNKNOW"
}

// ChangeEvent change event
type ChangeEvent struct {
	Namespace string
	Changes   map[string]*Change
}

func (e *ChangeEvent) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n[CHANGE]====\n")
	_, _ = fmt.Fprintf(&buf, "	[namespace]%s\n", e.Namespace)
	for _, change := range e.Changes {
		_, _ = fmt.Fprintf(&buf, "%s\n", change)
		_, _ = fmt.Fprintf(&buf, "--------\n")
	}

	buf.WriteString("============\n")
	return buf.String()
}

// Change represent a single key change
type Change struct {
	Key        string
	OldValue   string
	NewValue   string
	ChangeType ChangeType
}

func (c *Change) String() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "	[KEY]%s\n", c.Key)
	_, _ = fmt.Fprintf(&buf, "	[%s]\n", c.ChangeType)
	_, _ = fmt.Fprintf(&buf, "	[OLD]%s\n", c.OldValue)
	_, _ = fmt.Fprintf(&buf, "	[NEW]%s\n", c.NewValue)

	return buf.String()
}

func makeDeleteChange(key, value string) *Change {
	return &Change{
		Key:        key,
		ChangeType: DELETE,
		OldValue:   value,
	}
}

func makeModifyChange(key, oldValue, newValue string) *Change {
	return &Change{
		Key:        key,
		ChangeType: MODIFY,
		OldValue:   oldValue,
		NewValue:   newValue,
	}
}

func makeAddChange(key, value string) *Change {
	return &Change{
		Key:        key,
		ChangeType: ADD,
		NewValue:   value,
	}
}
