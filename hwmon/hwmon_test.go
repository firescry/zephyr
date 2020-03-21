package hwmon

import (
	"testing"
)

func TestJoinPaths(t *testing.T) {
	pathA := "/usr/"
	pathB := "/bin/"
	exp := "/usr/bin/"
	got := joinPaths(pathA, pathB)
	if got != exp {
		t.Errorf("joinPaths(%s, %s) = %s; expected %s", pathA, pathB, got, exp)
	}
}
