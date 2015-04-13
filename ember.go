package ember

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	root = "/"
)

// New returns a new ember instance.
func New() *Ember {
	e := Ember{Router: mux.NewRouter()}
	http.Handle(root, e.Router)
	return &e
}

func Vars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

// Ember provides simple routing for http://emberjs.com applications.
type Ember struct {
	// Router for extenting default functionality.
	Router *mux.Router
}

// Run listens on the TCP network address addr and then handles requests on incoming connections.
func (e *Ember) Run(addr string) {
	http.ListenAndServe(addr, nil)
}

// Index registers a handler for the index.html file.
func (e *Ember) Index(path string) *Ember {
	index := func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, path)
	}

	e.Router.HandleFunc(root, index).Name("index")
	e.Router.NotFoundHandler = http.HandlerFunc(index)
	return e
}

// Assets registers a handler for assets directory.
func (e *Ember) Assets(route string, path string) *Ember {
	fs := http.StripPrefix(route, http.FileServer(http.Dir(path)))
	e.Router.PathPrefix(route).Handler(fs)
	return e
}

// Namespace creates a new namespace instance.
func (e *Ember) Namespace(path string) *Namespace {
	return &Namespace{
		path:   path,
		Router: e.Router.PathPrefix(path).Subrouter(),
	}
}

// Model creates a new model instance.
func (e *Ember) Model(name string) *Model {
	return &Model{
		name:   name,
		Router: e.Router,
	}
}
