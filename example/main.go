package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kolbasov/ember"
	"net/http"
	"strconv"
	"sync/atomic"
)

// Todo model.
type Todo struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

// Todo DTO.
type TodoDto struct {
	Todo Todo `json:"todo"`
}

// Todos DTO.
type TodosDto struct {
	Todos []Todo `json:"todos"`
}

var (
	count int32 = 3
	Todos       = map[int]Todo{
		1: {ID: 1, Title: "Learn Ember.js", IsCompleted: true},
		2: {ID: 2, Title: "...", IsCompleted: false},
		3: {ID: 3, Title: "Profit!", IsCompleted: false},
	}
	addr string
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "Address, default :8080")
	flag.Parse()

	// Create an ember instance.
	e := ember.New()

	// Register a route for public/index.html
	e.Index("dist/index.html")

	// Register a route for public/assets
	e.Assets("/assets", "dist/assets")

	// Register a namespace.
	api := e.Namespace("/api/v1")

	// Create routes for a model.
	people := api.Model("todos")

	// GET /api/v1/todos
	people.FindAll(getTodos)

	// GET /api/v1/todos/{id}
	people.Find(getTodo)

	// POST /api/v1/todos
	people.Create(addTodo)

	// PUT /api/v1/todos/{id}
	people.Update(updateTodo)

	// DELETE /api/v1/todos/{id}
	people.Delete(deleteTodo)

	fmt.Println("Listening", addr)

	// Serve requests.
	e.Run(addr)
}

// GET /api/v1/todos
func getTodos(w http.ResponseWriter, req *http.Request) {
	var todos []Todo
	for _, v := range Todos {
		todos = append(todos, v)
	}

	e := json.NewEncoder(w)
	e.Encode(TodosDto{todos})
}

// GET /api/v1/todos/{id}
func getTodo(w http.ResponseWriter, req *http.Request) {
	todo, ok := findTodo(req)
	if !ok {
		http.Error(w, "todo not found", http.StatusBadRequest)
		return
	}

	e := json.NewEncoder(w)
	e.Encode(TodoDto{todo})
}

// POST /api/v1/todos
func addTodo(w http.ResponseWriter, req *http.Request) {
	var data TodoDto
	d := json.NewDecoder(req.Body)
	if err := d.Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.Todo.ID = int(atomic.AddInt32(&count, 1))
	Todos[data.Todo.ID] = data.Todo

	e := json.NewEncoder(w)
	e.Encode(data)
}

// PUT /api/v1/todos/{id}
func updateTodo(w http.ResponseWriter, req *http.Request) {
	todo, ok := findTodo(req)
	if !ok {
		http.Error(w, "todo not found", http.StatusBadRequest)
		return
	}

	var data TodoDto
	d := json.NewDecoder(req.Body)
	if err := d.Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.Title = data.Todo.Title
	todo.IsCompleted = data.Todo.IsCompleted

	Todos[todo.ID] = todo

	e := json.NewEncoder(w)
	e.Encode(TodoDto{todo})
}

// DELETE /api/v1/todos/{id}
func deleteTodo(w http.ResponseWriter, req *http.Request) {
	todo, ok := findTodo(req)
	if !ok {
		http.Error(w, "todo not found", http.StatusBadRequest)
		return
	}

	delete(Todos, todo.ID)

	e := json.NewEncoder(w)
	e.Encode(TodoDto{todo})
}

func findTodo(req *http.Request) (todo Todo, ok bool) {
	s, ok := ember.Vars(req)["id"]
	if !ok {
		return todo, false
	}

	id, err := strconv.Atoi(s)
	if err != nil {
		return todo, false
	}

	todo, ok = Todos[id]
	return todo, ok
}
