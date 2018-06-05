package agollo

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

// Change represent a single key change
type Change struct {
	OldValue   string
	NewValue   string
	ChangeType ChangeType
}
