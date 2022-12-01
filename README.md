# dittofunc

HTTP server implementation of Ditto <https://github.com/juanvillacortac/ditto> code generation.

## Usage

```go
// main.go

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ethanthatonekid/dittofunc/ditto/httpserver/ditto"
)

var origin = "http://localhost:8080"
var owner = "juanvillacortac"
var repo = "ditto"
var path = "/examples/angular/config.yml"
var ref = "main"
var token string

func main() {
	flag.StringVar(&origin, "origin", origin, "Ditto origin")
	flag.StringVar(&owner, "owner", owner, "GitHub owner")
	flag.StringVar(&repo, "repo", repo, "GitHub repo")
	flag.StringVar(&path, "path", path, "GitHub path")
	flag.StringVar(&ref, "ref", ref, "GitHub ref")
    flag.StringVar(&token, "token", token, "GitHub token")
	flag.Parse()

	// Create a new client.
	client := ditto.New(*http.DefaultClient, origin)

	// Create a new query.
	q := ditto.NewGenQuery(origin, owner, repo, path, ref)

	// Execute the query.
	generated, err := client.Gen(q)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the generated code.
	log.Println(generated)
}
```

```bash
# Spin up the server locally.
go run .

# Make a request to the server.
go run examples/hello_world/main.go -token=...
```

## Server

The Ditto `httpserver.Handler` is a simple HTTP handler that receives a Ditto configuration file from GitHub and generates the code.

### Self-hosting

Requirements:

- [Go](https://go.dev/dl/)

```bash
go run .
```

## Deployment

The Ditto `httpserver.Handler` is a simple HTTP handler that receives a Ditto configuration file from GitHub and generates the code.

Netlify is a great option for hosting the server.

### Netlify

Visit the [Netlify](https://www.netlify.com/) website and create an account.

[Get started with a new site from GitHub](https://app.netlify.com/start).

## License

[MIT](LICENSE)
