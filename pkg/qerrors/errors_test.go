package qerrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	var (
		base *Error
	)
	err := Newf(2, "state")
	err2 := Newf(2, "state")
	err3 := err.WithMetadata(map[string]string{
		"foo": "bar",
	})
	werr := fmt.Errorf("wrap %w", err)

	if errors.Is(err, new(Error)) {
		t.Errorf("should not be equal: %v", err)
	}
	if !errors.Is(werr, err) {
		t.Errorf("should be equal: %v", err)
	}
	if !errors.Is(werr, err2) {
		t.Errorf("should be equal: %v", err)
	}

	if !errors.As(err, &base) {
		t.Errorf("should be matchs: %v", err)
	}
	t.Logf("%s", err)

	if reason := State(err); reason != err3.State {
		t.Errorf("got %d want: %s", reason, err)
	}

	if err3.Metadata["foo"] != "bar" {
		t.Error("not expected metadata")
	}
}
