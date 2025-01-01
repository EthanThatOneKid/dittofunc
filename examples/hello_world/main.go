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
	q := ditto.NewGenQuery(token, owner, repo, path, ref)

	// Execute the query.
	generated, err := client.Gen(q)
	if err != nil {
		log.Fatalln(err)
	}

	// // Print the generated code.
	log.Println(generated)
}
