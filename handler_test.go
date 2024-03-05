package slogsentry

import (
	"errors"
	"testing"
)

func TestSlogErrorErrorMethod(t *testing.T) {
	tests := []struct {
		input        SlogError
		expectOutput string
	}{
		{SlogError{msg: "the message", err: errors.New("the error")}, "the message: the error"},
		{SlogError{err: errors.New("the error")}, "the error"},
		{SlogError{msg: "the message"}, "the message"},
	}

	for i, test := range tests {
		output := test.input.Error()
		if output != test.expectOutput {
			t.Errorf("test %d: expect: %q, got: %q", i, test.expectOutput, output)
		}
	}
}
