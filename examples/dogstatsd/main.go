package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/catkins/science"
)

const statsdAddress = "127.0.0.1:8125"

var scienceReporter science.Reporter

func main() {
	statsdClient, err := statsd.New(statsdAddress)
	scienceReporter = &DatadogReporter{client: statsdClient}

	if err != nil {
		log.Panicf("unable to make statsd client: %s", err.Error())
	}

	fetchValueFromServer("user-123")
	fetchValueFromServer("user-456")
}

type user struct {
	id   string
	name string
}

// the function being refactored
func fetchValueFromServer(id string) (user, error) {
	experiment := science.NewExperiment("fetch_users",
		science.WithReporter(scienceReporter),
		science.WithComparer(science.DeepEqualityCompare))

	experiment.Control(func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return user{id, "Chris"}, nil
	})

	experiment.Candidate(func() (interface{}, error) {
		if id == "user-456" {
			time.Sleep(5 * time.Millisecond)
			return user{id, "Chris"}, nil
		}
		time.Sleep(60 * time.Millisecond)
		return user{}, errors.New("unable to fetch user")
	})

	value, err := experiment.Run()

	if user, ok := value.(user); ok {
		return user, err
	}

	return user{}, err
}

// DatadogReporter is a science.Reporter implementation that sends metrics back to datadog
type DatadogReporter struct {
	client *statsd.Client
}

// Report sends results of experiment to datadog
//
// metrics sent:
// - "science.control_time": time taken to run the control function
// - "science.candidate_time": time taken to run the candidate function
// - "science.match": counter of matching retults
// - "science.mismatch": counter of mis-matching retults
func (reporter *DatadogReporter) Report(experimentName string, result bool, control science.Run, candidate science.Run) {
	resultString := "match"
	if !result {
		resultString = "mismatch"
	}
	resultTag := fmt.Sprintf("result:%s", resultString)
	experimentTag := fmt.Sprintf("exeriment:%s", experimentName)

	reporter.client.Timing("science.control_time",
		control.Duration,
		[]string{experimentTag, resultTag},
		1)

	reporter.client.Timing("science.candidate_time",
		candidate.Duration,
		[]string{experimentTag, resultTag},
		1)

	reporter.client.Incr(fmt.Sprintf("science.%s", resultString), []string{experimentTag}, 1)

}
