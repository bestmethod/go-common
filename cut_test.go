package gocommon

import (
	"testing"
)

func TestCut(t *testing.T) {
	if cut("Testing this  function", 1, " ") != "Testing" {
		t.Errorf("String: 'Testing this  function': Split pos 1 should be 'Testing'")
	}

	if cut("Testing this  function", 2, " ") != "this" {
		t.Errorf("String: 'Testing this  function': Split pos 2 should be 'this'")
	}

	if cut("Testing this  function", 3, " ") != "function" {
		t.Errorf("String: 'Testing this  function': Split pos 3 should be 'function'")
	}

	if cut("Testing this  function", 4, " ") != "" {
		t.Errorf("String: 'Testing this  function': Split pos 4 should be empty string")
	}
}
