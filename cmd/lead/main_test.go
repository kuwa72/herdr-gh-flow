package main

import (
	"strings"
	"testing"
)

func TestVersionStringContainsAllFields(t *testing.T) {
	Version, Commit, Date = "v0.0.0-test", "deadbee", "2026-09-06"
	t.Cleanup(func() {
		Version, Commit, Date = "dev", "none", "unknown"
	})

	got := versionString()
	for _, want := range []string{"lead", "v0.0.0-test", "deadbee", "2026-09-06"} {
		if !strings.Contains(got, want) {
			t.Fatalf("versionString() = %q, missing %q", got, want)
		}
	}
}

func TestRunVersionCommands(t *testing.T) {
	for _, args := range [][]string{{"version"}, {"--version"}, {"-v"}} {
		if code := run(args); code != 0 {
			t.Fatalf("run(%v) = %d, want 0", args, code)
		}
	}
}

func TestRunUnknownCommandFails(t *testing.T) {
	if code := run([]string{"no-such-command"}); code == 0 {
		t.Fatal("run(unknown) = 0, want non-zero")
	}
}
