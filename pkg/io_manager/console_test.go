package io_manager

import (
	"errors"
	"github.com/pkg/term"
	"testing"
)

func TestAskPlainReturnsValidChoiceAfterRetry(t *testing.T) {
	reads := []string{"wrong", "2"}

	choice := askPlain("Choose action", map[string]string{
		"1": "Save",
		"2": "Restore",
	}, func(string) string {
		current := reads[0]
		reads = reads[1:]
		return current
	})

	if choice != "2" {
		t.Fatalf("expected fallback prompt to return the valid choice, got %q", choice)
	}
}

func TestGetInputReturnsZeroWhenTTYUnavailable(t *testing.T) {
	previousOpenTTY := openTTY
	openTTY = func() (*term.Term, error) {
		return nil, errors.New("tty unavailable")
	}
	defer func() {
		openTTY = previousOpenTTY
	}()

	if got := getInput(); got != 0 {
		t.Fatalf("expected getInput to return 0 when tty is unavailable, got %d", got)
	}
}
