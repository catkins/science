package science

import "reflect"

// DeepEqualityCompare implements a deep equality check between results
func DeepEqualityCompare(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
