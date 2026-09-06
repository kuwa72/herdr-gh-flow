// Package ghcli implements ports.GhClient via the `gh` CLI
// (transparent auth; go-gh library adoption is deferred).
package ghcli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kuwa72/lead-cli/internal/ports"
)

// ListLimit mirrors legacy bin/hgf `gh issue list --limit 50`.
const ListLimit = 50

// Client calls the `gh` CLI. Zero value is usable.
type Client struct {
	// Bin is the gh binary name or path. Empty means "gh".
	Bin string
	// LookPath resolves the binary; defaults to exec.LookPath.
	// A missing binary yields *ports.BinaryNotFoundError so callers
	// can fall back gracefully.
	LookPath func(name string) (string, error)
}

// New returns a Client with defaults.
func New() *Client { return &Client{} }

var _ ports.GhClient = (*Client)(nil)

func (c *Client) bin() string {
	if c.Bin != "" {
		return c.Bin
	}
	return "gh"
}

func (c *Client) lookPath() func(string) (string, error) {
	if c.LookPath != nil {
		return c.LookPath
	}
	return exec.LookPath
}

// run executes bin with args, returning stdout. Exit failures surface
// gh's stderr so callers see the real cause.
func (c *Client) run(ctx context.Context, args ...string) ([]byte, error) {
	bin := c.bin()
	if _, err := c.lookPath()(bin); err != nil {
		return nil, &ports.BinaryNotFoundError{Binary: bin}
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return nil, fmt.Errorf("gh %s: %w: %s", strings.Join(args, " "), err, msg)
		}
		return nil, fmt.Errorf("gh %s: %w", strings.Join(args, " "), err)
	}
	return stdout.Bytes(), nil
}

// ListOpen runs `gh issue list --state open --limit 50 --json number,title`.
func (c *Client) ListOpen(ctx context.Context) ([]ports.IssueSummary, error) {
	out, err := c.run(ctx, "issue", "list",
		"--state", "open",
		"--limit", strconv.Itoa(ListLimit),
		"--json", "number,title")
	if err != nil {
		return nil, err
	}
	var raw []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &raw); err != nil {
		return nil, fmt.Errorf("gh issue list: decode JSON: %w", err)
	}
	summaries := make([]ports.IssueSummary, len(raw))
	for i, r := range raw {
		summaries[i] = ports.IssueSummary{Number: r.Number, Title: r.Title}
	}
	return summaries, nil
}

// View runs `gh issue view <n> --json number,title,body,state`.
func (c *Client) View(ctx context.Context, number int) (ports.Issue, error) {
	out, err := c.run(ctx, "issue", "view", strconv.Itoa(number),
		"--json", "number,title,body,state")
	if err != nil {
		return ports.Issue{}, err
	}
	var raw struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		State  string `json:"state"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &raw); err != nil {
		return ports.Issue{}, fmt.Errorf("gh issue view %d: decode JSON: %w", number, err)
	}
	return ports.Issue{
		Number: raw.Number,
		Title:  raw.Title,
		Body:   raw.Body,
		State:  raw.State,
	}, nil
}
