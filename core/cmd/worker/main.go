package main

import (
	"beelder/internal/worker"
	"beelder/internal/worker/builder"
	"log"
)

func main() {

	builder := builder.NewBuilder()
	worker := worker.NewWorker(builder)

	if err := worker.Start(); err != nil {
		log.Fatal(err)
	}
}
