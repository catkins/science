package science

import (
	"log"
	"os"
)

// Experiment captures a single codepath refactoring
type Experiment struct {
	name      string
	comparer  Comparer
	reporter  Reporter
	control   ExperimentalFunc
	candidate ExperimentalFunc
}

// Comparer is used to check for equality between two results
type Comparer func(a, b interface{}) bool

// DefaultComparer is the default Comparer provided when creating a new Experiment via NewExperiment
var DefaultComparer = SimpleEqualityCompare

// Reporter is used to report the results of experiment runs
// eg. printing out to logs, or sending metrics to statsd or prometheus
type Reporter interface {
	Report(experimentName string, result bool, control Run, candidate Run)
}

// DefaultReporter is the default Reporter provided when creating a new Experiment via NewExperiment
var DefaultReporter Reporter = SimpleReporter{os.Stderr}

// ExperimentalFunc is a function which returns a value or an error
type ExperimentalFunc func() (value interface{}, err error)

// NewExperiment builds and initializes a new Experiment
func NewExperiment(name string, options ...ExperimentOption) *Experiment {
	experiment := Experiment{
		name:     name,
		comparer: DefaultComparer,
		reporter: DefaultReporter,
	}

	for _, option := range options {
		option(&experiment)
	}

	return &experiment
}

// Control sets the experiments control function
// This would be the "old" code that you are refactoring away
func (experiment *Experiment) Control(controlFunc ExperimentalFunc) {
	experiment.control = controlFunc
}

// Candidate sets the experiments candidate function
// This would be the "new" code that you are testing out
func (experiment *Experiment) Candidate(candidateFunc ExperimentalFunc) {
	experiment.candidate = candidateFunc
}

// Run executes both the control and candidate functions, and returns the results from the control function
//
// If Control hasn't been set, then Run() will panic
// If control function panics, it will propagate out
func (experiment *Experiment) Run() (interface{}, error) {
	if experiment.control == nil {
		log.Panicf("experiment %s has no control function provided", experiment.name)
	}

	controlRun := runExperimentalFunc(experiment.control)
	candidateRun := runExperimentalFunc(experiment.candidate)

	result := experiment.comparer(controlRun.Value, candidateRun.Value)
	result = result && experiment.comparer(controlRun.Value, candidateRun.Value)

	experiment.reporter.Report(experiment.name, result, controlRun, candidateRun)

	if controlRun.Panicked {
		log.Panicf("experiment %s panicked: %s", experiment.name, controlRun.Err)
	}

	return controlRun.Value, controlRun.Err
}

// ExperimentOption represents a functional option used to configure an Experiment instance
type ExperimentOption func(*Experiment)

// WithReporter is a constructor option to allow you to provide your own Reporter to a new Experiment
func WithReporter(reporter Reporter) ExperimentOption {
	return func(experiment *Experiment) {
		experiment.reporter = reporter
	}
}

// WithComparer is a constructor option to allow you to provide your own Comparer to a new Experiment
func WithComparer(comparer Comparer) ExperimentOption {
	return func(experiment *Experiment) {
		experiment.comparer = comparer
	}
}
