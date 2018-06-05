package agollo

// Queryer query updated namespaces
type Queryer interface {
	Query(notifications []*notification)
}

type configQueryer struct {
}

func (q configQueryer) Query(notifications []*notification) (map[string]ChangeEvent, error) {
	return nil, nil
}
