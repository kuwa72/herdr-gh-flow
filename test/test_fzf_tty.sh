#!/usr/bin/env bash
set -euo pipefail

# Test that select_issue in bin/hgf passes /dev/tty for interactive fzf
if ! grep -q '</dev/tty' bin/hgf && ! grep -q '< /dev/tty' bin/hgf; then
  echo "FAIL: bin/hgf select_issue does not redirect input from /dev/tty for interactive fzf" >&2
  exit 1
fi

echo "All select_issue TTY tests passed."
exit 0
