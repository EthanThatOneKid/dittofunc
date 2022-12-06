package main

import (
	"log"

	"github.com/apex/gateway"

	"github.com/ethanthatonekid/dittofunc/ditto/httpserver"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	return gateway.ListenAndServe("", httpserver.New())
}
