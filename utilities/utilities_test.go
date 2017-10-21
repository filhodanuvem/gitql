package utilities

import "testing"

func TestIsFieldPresentInArray(t *testing.T) {
	testArray := []string{"first", "second", "third", "fourth"}

	if IsFieldPresentInArray(testArray, "fifth") {
		t.Error("Fifth should not be in the test Array")
	}

	if !IsFieldPresentInArray(testArray, "first") {
		t.Error("First should bin in the testArray")
	}

	if IsFieldPresentInArray(testArray, "second") {
		return
	}

}
