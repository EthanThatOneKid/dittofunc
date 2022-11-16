// Run: go run example/hello_world/main.go -origin=http://localhost:8080 -token=... -owner=... -repo=... -path=... -ref=...
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ethanthatonekid/dittofunc/dittofunc"
)

var origin = "http://localhost:8080"
var owner string
var repo string
var path string
var ref string
var token string

func main() {
	flag.StringVar(&origin, "origin", origin, "Ditto origin")
	flag.StringVar(&token, "token", token, "GitHub token")
	flag.StringVar(&owner, "owner", owner, "GitHub owner")
	flag.StringVar(&repo, "repo", repo, "GitHub repo")
	flag.StringVar(&path, "path", path, "GitHub path")
	flag.StringVar(&ref, "ref", ref, "GitHub ref")
	flag.Parse()

	// Create a new client.
	e := MustEnv()
	client := dittofunc.NewHTTPClient(*http.DefaultClient, e.Origin)

	// Create a new query.
	q := dittofunc.NewGenQuery(e.GithubToken, e.Owner, e.Repo, e.Path, e.Ref)

	// Execute the query.
	generated, err := client.Gen(q)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the generated code.
	log.Println(generated)
}

func MustEnv() env {
	var e env

	e.Origin = origin
	if e.Origin == "" {
		e.Origin = "http://localhost:8080"
	}

	e.GithubToken = token
	if e.GithubToken == "" {
		log.Fatalln("token is not set")
	}

	e.Owner = owner
	if e.Owner == "" {
		log.Fatalln("owner is not set")
	}

	e.Repo = repo
	if e.Repo == "" {
		log.Fatalln("repo is not set")
	}

	e.Path = path
	if e.Path == "" {
		log.Fatalln("path is not set")
	}

	e.Ref = ref
	if e.Ref == "" {
		log.Fatalln("ref is not set")
	}

	return e
}

type env struct {
	// Origin is the origin of the DittoFunc server.
	// Defaults to "http://localhost:8080".
	Origin string

	// GithubToken is the GitHub token to use.
	// Required.
	GithubToken string

	// Owner is the owner of the repository.
	// Required.
	Owner string

	// Repo is the repository name.
	// Required.
	Repo string

	// Path is the path to the file.
	// Required.
	Path string

	// Ref is the ref to use.
	// Defaults to zero value.
	Ref string
}
