package ember

import (
	"github.com/gorilla/mux"
	"net/http"
	"path"
)

type Model struct {
	name string
	//namespace string
	// Router for extenting default functionality.
	Router *mux.Router
}

// Find registers a handler for GET /model/{id} route.
func (m *Model) Find(h http.HandlerFunc) *Model {
	return m.get(h, m.name, "/{id}")
}

// Find registers a handler for GET /model route.
func (m *Model) FindAll(h http.HandlerFunc) *Model {
	return m.get(h, m.name)
}

// Update registers a handler for PUT /model/{id} route.
func (m *Model) Update(h http.HandlerFunc) *Model {
	return m.put(h, m.name, "/{id}")
}

// Create registers a handler for POST /model route.
func (m *Model) Create(h http.HandlerFunc) *Model {
	return m.post(h, m.name)
}

// Create registers a handler for DELETE /model/{id} route.
func (m *Model) Delete(h http.HandlerFunc) *Model {
	return m.delete(h, m.name, "/{id}")
}

func (m *Model) register(method string, h http.HandlerFunc, paths ...string) *Model {
	p := path.Join("/", path.Join(paths...))
	m.Router.Methods(method).Path(p).HandlerFunc(h)
	return m
}

func (m *Model) get(h http.HandlerFunc, paths ...string) *Model {
	return m.register("GET", h, paths...)
}

func (m *Model) put(h http.HandlerFunc, paths ...string) *Model {
	return m.register("PUT", h, paths...)
}

func (m *Model) post(h http.HandlerFunc, paths ...string) *Model {
	return m.register("POST", h, paths...)
}

func (m *Model) delete(h http.HandlerFunc, paths ...string) *Model {
	return m.register("DELETE", h, paths...)
}
