package main

import (
	"log"

	"github.com/apex/gateway"
	"github.com/ethanthatonekid/dittofunc/dittofunc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	return gateway.ListenAndServe("", dittofunc.NewHandler())
}
