// Package agent implements ports.AgentLauncher: direct (inline) execution
// of coding agents, plus agent resolution and reviewable command strings.
package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/kuwa72/lead-cli/internal/ports"
)

// DefaultAgent is the standard agent (legacy bin/hgf DEFAULT_AGENT).
const DefaultAgent = "agy"

// KnownAgents lists the agents lead can dispatch to
// (docs/design.md §4 plus legacy bin/hgf options).
var KnownAgents = []string{"agy", "claude", "codex", "devin", "opencode", "gemini"}

// Resolve maps a requested agent name to the effective one.
// Empty requests fall back to DefaultAgent.
func Resolve(requested string) string {
	if requested == "" {
		return DefaultAgent
	}
	return requested
}

// IsKnown reports whether name is a recognized agent.
func IsKnown(name string) bool {
	for _, known := range KnownAgents {
		if name == known {
			return true
		}
	}
	return false
}

// Argv builds the exec argv for launching agent with prompt.
// agy takes interactive prompts via -i/--prompt-interactive (closed #19);
// other agents take the prompt positionally.
func Argv(agentName, prompt string) []string {
	name := Resolve(agentName)
	if name == "agy" {
		return []string{"-i", prompt}
	}
	return []string{prompt}
}

// CommandString builds the shell command string prepared in a new Herdr
// pane for human review (legacy bin/hgf behavior; not auto-sent).
func CommandString(agentName, prompt string) string {
	name := Resolve(agentName)
	if name == "agy" {
		return "agy -i " + strconv.Quote(prompt)
	}
	return name + " " + strconv.Quote(prompt)
}

// Launcher runs the agent binary directly in the current terminal
// (inline fallback). Zero value is usable.
type Launcher struct {
	// LookPath resolves the agent binary; defaults to exec.LookPath.
	LookPath func(name string) (string, error)
	// Stdout/Stderr receive the agent's output; nil inherits os.Stdout/os.Stderr.
	Stdout io.Writer
	Stderr io.Writer
}

// New returns a Launcher with defaults.
func New() *Launcher { return &Launcher{} }

var _ ports.AgentLauncher = (*Launcher)(nil)

func (l *Launcher) lookPath() func(string) (string, error) {
	if l.LookPath != nil {
		return l.LookPath
	}
	return exec.LookPath
}

// Launch resolves the agent, verifies its binary is on PATH, and execs it.
// A missing binary yields *ports.BinaryNotFoundError; a failing agent
// propagates its exit status.
func (l *Launcher) Launch(ctx context.Context, agentName string, prompt string) error {
	name := Resolve(agentName)
	if _, err := l.lookPath()(name); err != nil {
		return &ports.BinaryNotFoundError{Binary: name}
	}
	cmd := exec.CommandContext(ctx, name, Argv(name, prompt)...)
	if l.Stdout != nil {
		cmd.Stdout = l.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}
	if l.Stderr != nil {
		cmd.Stderr = l.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("agent %s: %w", name, err)
	}
	return nil
}
