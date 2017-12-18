package science

import (
	"fmt"
	"time"
)

// Run encapsulates the output and timing from a given experiment run
type Run struct {
	// The value from the run
	Value interface{}

	// The amount to time taken for the function to return
	Duration time.Duration

	// Any error returned by the function
	Err error

	Panicked bool
}

func runExperimentalFunc(experimentalFunc ExperimentalFunc) (run Run) {
	run = Run{}

	startTime := time.Now()

	run.Value, run.Err = experimentalFunc()

	defer func() {
		run.Duration = time.Since(startTime)

		if r := recover(); r != nil {
			run.Panicked = true

			if err := r.(error); err != nil {
				run.Err = err
			} else {
				run.Err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	return run
}
