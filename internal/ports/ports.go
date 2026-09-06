// Package ports isolates external CLI dependencies behind interfaces.
// See docs/rfc-22-language-migration.md §5 for the design rationale.
package ports

import (
	"context"
	"errors"
	"fmt"
)

// IssueSummary is a single row of `gh issue list`.
type IssueSummary struct {
	Number int
	Title  string
}

// Issue is the detail of `gh issue view`.
type Issue struct {
	Number int
	Title  string
	Body   string
	State  string
}

// GhClient abstracts `gh issue list/view` (transparent auth via gh CLI).
type GhClient interface {
	ListOpen(ctx context.Context) ([]IssueSummary, error)
	View(ctx context.Context, number int) (Issue, error)
}

// Direction is the herdr pane split direction.
type Direction string

const (
	DirectionRight Direction = "right"
	DirectionDown  Direction = "down"
)

// HerdrRunner abstracts `herdr pane split/send-text`.
// The agent command is prepared in the new pane for human review;
// it is never auto-sent (no `pane run`).
type HerdrRunner interface {
	Split(ctx context.Context, dir Direction, ratio float64) (paneID string, err error)
	SendText(ctx context.Context, paneID string, text string) error
}

// AgentLauncher runs a coding agent directly in the current terminal
// (inline fallback when Herdr is unavailable).
type AgentLauncher interface {
	Launch(ctx context.Context, agent string, prompt string) error
}

// BinaryNotFoundError reports a missing external binary on PATH.
// Callers use it to fall back gracefully (e.g. InlineRunner).
type BinaryNotFoundError struct {
	Binary string
}

func (e *BinaryNotFoundError) Error() string {
	return fmt.Sprintf("%s: executable not found on PATH", e.Binary)
}

// IsBinaryNotFound reports whether err wraps a *BinaryNotFoundError.
func IsBinaryNotFound(err error) bool {
	var target *BinaryNotFoundError
	return errors.As(err, &target)
}
