#!/usr/bin/env bash
set -euo pipefail

# Test that select_issue in bin/hgf does NOT redirect /dev/tty into fzf stdin (which breaks stdin pipe)
if grep -q '</dev/tty' bin/hgf || grep -q '< /dev/tty' bin/hgf; then
  echo "FAIL: bin/hgf redirects /dev/tty into fzf stdin, breaking piped issue list" >&2
  exit 1
fi

echo "All select_issue fzf stdin tests passed."
exit 0
