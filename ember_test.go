package ember

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var e = New()

type routeTest struct {
	method   string
	path     string
	expected string
}

func init() {
	// Register a route for index.html file.
	e.Index("./example/dist/index.html")

	// Register a route for assets.
	e.Assets("/assets", "./example/dist/assets")

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
}

func responce(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%v %v", req.Method, req.URL)
}

func TestIndex(t *testing.T) {
	server := httptest.NewServer(e.Router)
	defer server.Close()

	testFile(t, server, "/")
}

func TestAssets(t *testing.T) {
	server := httptest.NewServer(e.Router)
	defer server.Close()

	tests := []string{
		"/assets/example.css",
		"/assets/example.js",
	}

	for _, test := range tests {
		testFile(t, server, test)
	}
}

func TestNotFound(t *testing.T) {
	server := httptest.NewServer(e.Router)
	defer server.Close()

	testFile(t, server, "/devnull")
}

func TestGetModel(t *testing.T) {
	server := httptest.NewServer(e.Router)
	defer server.Close()

	tests := []routeTest{
		{
			method:   "GET",
			path:     "/stats",
			expected: "GET /stats",
		},
	}

	for _, test := range tests {
		testRoute(t, server, test)
	}
}

func TestNamespace(t *testing.T) {
	server := httptest.NewServer(e.Router)
	defer server.Close()

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
		testRoute(t, server, test)
	}
}

func testRoute(t *testing.T, server *httptest.Server, test routeTest) {
	res, err := sendRequest(test.method, server.URL+test.path)
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

func testFile(t *testing.T, server *httptest.Server, path string) {
	res, err := sendRequest("GET", server.URL+path)
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
