package ember

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var addr string

type routeTest struct {
	method   string
	path     string
	expected string
}

func init() {
	flag.StringVar(&addr, "addr", ":8080", "Address, default :8080")
	flag.Parse()

	e := New()

	// Register a route for index.html file.
	e.Index("example/dist/index.html")

	// Register a route for assets.
	e.Assets("/assets", "example/dist/assets")

	// Create a model and register a route for /stats
	e.Model("stats").FindAll(responce)

	// Register a namespace.
	api := e.Namespace("/api/v1")

	// Create routes for a model.
	people := api.Model("people")

	// GET /api/v1/people
	people.FindAll(responce)

	// GET /api/v1/people/{id}
	people.Find(responce)

	// POST /api/v1/people
	people.Create(responce)

	// PUT /api/v1/people/{id}
	people.Update(responce)

	// DELETE /api/v1/people/{id}
	people.Delete(responce)

	// Not the best idea...
	go e.Run(addr)
}

func responce(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%v %v", req.Method, req.URL)
}

func TestIndex(t *testing.T) {
	testFile(t, "/")
}

func TestAssets(t *testing.T) {
	tests := []string{
		"/assets/example.css",
		"/assets/example.js",
	}

	for _, test := range tests {
		testFile(t, test)
	}
}

func TestNotFound(t *testing.T) {
	testFile(t, "/devnull")
}

func TestGetModel(t *testing.T) {
	tests := []routeTest{
		{
			method:   "GET",
			path:     "/stats",
			expected: "GET /stats",
		},
	}

	for _, test := range tests {
		testRoute(t, test)
	}
}

func TestNamespace(t *testing.T) {
	tests := []routeTest{
		{
			method:   "GET",
			path:     "/api/v1/people",
			expected: "GET /api/v1/people",
		},
		{
			method:   "GET",
			path:     "/api/v1/people/1",
			expected: "GET /api/v1/people/1",
		},
		{
			method:   "PUT",
			path:     "/api/v1/people/1",
			expected: "PUT /api/v1/people/1",
		},
		{
			method:   "POST",
			path:     "/api/v1/people",
			expected: "POST /api/v1/people",
		},
		{
			method:   "DELETE",
			path:     "/api/v1/people/1",
			expected: "DELETE /api/v1/people/1",
		},
	}

	for _, test := range tests {
		testRoute(t, test)
	}
}

func testRoute(t *testing.T, test routeTest) {
	res, err := sendRequest(test.method, "http://localhost"+addr+test.path)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	actual := strings.Trim(string(data), "\n")

	if test.expected == actual {
		return
	}

	t.Errorf(
		"[%v %v] expected: [%v], got: [%v]\n",
		test.method,
		test.path,
		test.expected,
		actual,
	)
}

func testFile(t *testing.T, path string) {
	res, err := sendRequest("GET", "http://localhost"+addr+path)
	if err != nil {
		panic(err)
	}

	if res.StatusCode == http.StatusOK {
		return
	}

	t.Errorf("[%v] expected: [%d], got: [%d]\n", path, http.StatusOK, res.StatusCode)
}

func sendRequest(method string, url string) (res *http.Response, err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}
