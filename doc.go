/*
Package ember is a simple back-end for http://emberjs.com applications which use the RESTAdapter.

Using Ember

Getting Ember

	go get github.com/kolbasov/ember

Running Tests

	go test

TodoMVC Demo Application
	The example folder contains TodoMVC demo application.

	cd example
	go run main.go

Defining routes

Let's assume that our applications shows an editable list of people. Here is the project's structure.

Project structure:
	app
	  public
	    assets
	      app.cs
	      app.js
	      ...
	  index.html
	main.go

We need to define routes for the index.html, assets and REST API. Also, let's prefix the API with "api/v1".

Static resources:
	index.html: GET /
	assets:     GET /assets

API:
	Find:    GET    api/v1/people/1
	FindAll: GET    api/v1/people
	Update:  PUT    api/v1/people/1
	Create:  POST   api/v1/people
	Delete:  DELETE api/v1/people/1

First, we need to create an Ember instance.

	// Create an ember instance.
	e := ember.New()

Then we can define routes for the index.html and assets.

	// Register a route for public/index.html
	e.Index("public/index.html")

	// Register a route for public/assets
	e.Assets("/assets", "public/assets")

We use Ember.Namespace method to create a namespace.

	// Register a namespace.
	api := e.Namespace("/api/v1")

Now we can add a model to the namespace and define its routes.

	// Create routes for a model.
	people := api.Model("people")

	// GET /api/v1/people
	people.FindAll(getPeople)
	...

Call the Ember.Run method to start the application.

	// Serve requests.
	e.Run(":8080")

Here is a part of the main.go file:

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


The application serves the following routes:

index.html and assets:
	GET: http://localhost:8080

API:
	GET/POST   http://localhost:8080/api/v1/people
	PUT/DELETE http://localhost:8080/api/v1/people/1

It's possible to create a route without a namespace.
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

The application will serve the following route:
	GET http://localhost:8080/stats
*/
package ember
