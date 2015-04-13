package ember

import "github.com/gorilla/mux"

type Namespace struct {
	path   string
	Router *mux.Router
}

func (ns *Namespace) Model(name string) *Model {
	return &Model{
		name:   name,
		Router: ns.Router,
	}
}
