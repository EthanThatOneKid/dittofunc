package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ethanthatonekid/dittofunc/dittofunc"
)

var addr = "localhost:8080"

func main() {
	flag.StringVar(&addr, "addr", addr, "HTTP address to listen on")
	flag.Parse()

	log.Println("listening on address", addr)

	if err := http.ListenAndServe(addr, dittofunc.NewHandler()); err != nil {
		log.Fatalln(err)
	}
}
