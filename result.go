package science

// Result represents the result of a single experiment run
type Result struct {
	// The output and timing from the control path
	Control Run

	// The output and timing from the candidate path
	Candidate Run
}
