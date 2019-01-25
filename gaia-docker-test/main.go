package main

import (
	"log"

	"github.com/michelvocks/gaia-docker-test/golang"
)

// Job1Handler is a job.
func Job1Handler(args golang.Arguments) error {
	log.Println("Job1 triggered!")
	return nil
}

func main() {
	// Create example job
	jobs := golang.Jobs{
		golang.Job{
			Title:   "Job1",
			Handler: Job1Handler,
		},
	}

	if err := golang.Serve(jobs); err != nil {
		panic(err)
	}
}
