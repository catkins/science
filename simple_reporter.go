package science

import (
	"fmt"
	"io"
)

// SimpleReporter logs result out to provided io.Writer
type SimpleReporter struct {
	Output io.Writer
}

// Report writes out a the result to the SimpleReporter's output
func (reporter SimpleReporter) Report(experimentName string, result bool, control Run, candidate Run) {
	fmt.Fprintf(
		reporter.Output,
		"experiment: %s, equal: %t, control: %s, candidate: %s\n",
		experimentName,
		result,
		formatRun(control),
		formatRun(candidate))
}

func formatRun(run Run) string {
	if run.Err != nil {
		return fmt.Sprintf("returned error %q in %s", run.Err, run.Duration)
	}

	return fmt.Sprintf("returned %+v in %s", run.Value, run.Duration)
}
