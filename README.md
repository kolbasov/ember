[![Build Status](https://travis-ci.org/kolbasov/ember.svg?branch=master)](https://travis-ci.org/kolbasov/ember)
[![GoDoc](https://godoc.org/github.com/kolbasov/ember?status.svg)](https://godoc.org/github.com/kolbasov/ember)

## Ember

Go package **ember** is a simple back-end for http://emberjs.com applications which use the RESTAdapter.

## Using Ember

### Getting Ember
```
go get github.com/kolbasov/ember
```

### Running Tests
```
go test
```

### TodoMVC Demo Application
The ```example``` folder contains TodoMVC demo application.
```
cd example
go run main.go
```


### Defining routes

Let's assume that our applications shows an editable list of people. Here is the project's structure.
```	
app
  public
    assets
      app.cs
      app.js
      ...
  index.html
main.go
```

We need to define routes for the index.html, assets and REST API. Also, let's prefix the API with ```api/v1```.

<table>
  <thead>
    <tr><th>Action</th><th>HTTP Verb</th><th>URL</th></tr>
  </thead>
  <tbody>
  	<tr><th>index.html</th><td>GET</td><td>/</td></tr>
  	<tr><th>assets</th><td>GET</td><td>/assets</td></tr>
  	<tr><th colspan="3">REST API</th></tr>
    <tr><th>Find</th><td>GET</td><td>api/v1/people/1</td></tr>
    <tr><th>Find All</th><td>GET</td><td>api/v1/people</td></tr>
    <tr><th>Update</th><td>PUT</td><td>api/v1/people/1</td></tr>
    <tr><th>Create</th><td>POST</td><td>api/v1/people</td></tr>
    <tr><th>Delete</th><td>DELETE</td><td>api/v1/people/1</td></tr>
  </tbody>
</table>

First, we need to create an Ember instance.
```go
// Create an ember instance.
e := ember.New()
```

Then we can define routes for the index.html and assets.
```go
// Register a route for public/index.html
e.Index("public/index.html")

// Register a route for public/assets
e.Assets("/assets", "public/assets")
```

We use ```Ember.Namespace``` method to create a namespace.
```go
// Register a namespace.
api := e.Namespace("/api/v1")
```

Now we can add a model to the namespace and define its routes.
```go
// Create routes for a model.
people := api.Model("people")

// GET /api/v1/people
people.FindAll(getPeople)
...
```

Call the ```Ember.Run``` method to start the application.
```go
// Serve requests.
e.Run(":8080")
```

Here is a part of the main.go file:
```go
	import (
		"net/http"
		"github.com/kolbasov/ember"
	)

	func main() {
		// Create an ember instance.
		e := ember.New()

		// Register a route for public/index.html
		e.Index("public/index.html")

		// Register a route for public/assets
		e.Assets("/assets", "public/assets")

		// Register a namespace.
		api := e.Namespace("/api/v1")

		// Create routes for a model.
		people := api.Model("people")

		// GET /api/v1/people
		people.FindAll(getPeople)

		// GET /api/v1/people/{id}
		people.Find(getPerson)

		// POST /api/v1/people
		people.Create(addPerson)

		// PUT /api/v1/people/{id}
		people.Update(updatePerson)

		// DELETE /api/v1/people/{id}
		people.Delete(deletePerson)

		// Serve requests.
		e.Run(":8080")
	}

// Functions getPeople, getPerson, updatePerson and deletePerson.
...
```

The application serves the following routes:

<table>
  <thead>
    <tr><th>HTTP Verb</th><th>URL</th></tr>
  </thead>
  <tbody>
  	<tr>
  		<th>GET</th>
  		<td>http://localhost:8080</td>  		
  	</tr>
  	<tr>
  		<th>GET/POST</th>
  		<td>http://localhost:8080/api/v1/people</td>  		
  	</tr>
  	<tr>
  		<th>PUT/DELETE</th>
  		<td>http://localhost:8080/api/v1/people/1</td>  		
  	</tr>
  </tbody>
</table>

Also, it's possible to define a route for a model without a namespace.

```go
	import (
		"net/http"
		"github.com/kolbasov/ember"
	)

	func main() {
		// Create an ember instance.
		e := ember.New()

		// Create routes for a model without a namespace and register a route.
		// GET /stats
		e.Model("stats").FindAll(getStats)

		// Serve requests.
		e.Run(":8080")
	}
```

The application will serve the following route:

<table>
  <thead>
    <tr><th>HTTP Verb</th><th>URL</th></tr>
  </thead>
  <tbody>
  	<tr>
  		<th>GET</th>
  		<td>http://localhost:8080/stats</td>  		
  	</tr>
  </tbody>
</table>
