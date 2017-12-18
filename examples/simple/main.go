package main

import (
	"errors"
	"time"

	"github.com/catkins/science"
)

func main() {
	experiment := science.NewExperiment("segse")

	experiment.Control(func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return "some value", nil
	})

	experiment.Candidate(func() (interface{}, error) {
		time.Sleep(60 * time.Millisecond)
		return nil, errors.New("unable to fetch value")
	})

	experiment.Run()
}
