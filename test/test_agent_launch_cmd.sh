#!/usr/bin/env bash
set -euo pipefail

# Test that agy invocation uses -i / --prompt-interactive instead of bare positional argument

# 1. Check Herdr pane send-text
if grep -E 'send-text .* "agy [^"-]' bin/hgf; then
  echo "FAIL: bin/hgf sends 'agy \"\$prompt_text\"' without -i or --prompt-interactive" >&2
  exit 1
fi

# 2. Check direct terminal launch
if grep -E 'agy "\$prompt_text"' bin/hgf; then
  echo "FAIL: bin/hgf launches 'agy \"\$prompt_text\"' without -i or --prompt-interactive" >&2
  exit 1
fi

echo "All agent launch command tests passed."
exit 0
