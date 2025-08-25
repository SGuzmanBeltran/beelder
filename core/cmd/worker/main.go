package main

import (
	"beelder/internal/worker"
	"log"
)

func main() {

	builder := worker.NewBuilder()
	worker := worker.NewWorker(builder)

	if err := worker.Start(); err != nil {
		log.Fatal(err)
	}
}
