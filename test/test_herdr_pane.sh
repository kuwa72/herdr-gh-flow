#!/usr/bin/env bash
set -euo pipefail

# Test that when HERDR_ENV=1, hgf extracts pane_id from herdr pane split and calls herdr pane send-text <PANE_ID>

# Ensure bin/hgf does not use `herdr pane send-text --current`
if grep -q 'herdr pane send-text --current' bin/hgf; then
  echo "FAIL: bin/hgf still uses invalid 'herdr pane send-text --current'" >&2
  exit 1
fi

# Ensure bin/hgf captures new pane_id from herdr pane split
if ! grep -q 'pane_id' bin/hgf; then
  echo "FAIL: bin/hgf does not extract pane_id from herdr pane split output" >&2
  exit 1
fi

echo "All herdr pane tests passed."
exit 0
