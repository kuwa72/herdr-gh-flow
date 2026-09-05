#!/usr/bin/env bash
set -euo pipefail

# TDD Test for bin/hgf

# 1. Check help option
output=$(./bin/hgf --help 2>&1 || true)
if ! echo "$output" | grep -q "Usage: hgf"; then
  echo "FAIL: bin/hgf --help does not contain Usage" >&2
  exit 1
fi

# 2. Check prompt generation function
prompt_output=$(./bin/hgf prompt 1 2>&1 || true)
if ! echo "$prompt_output" | grep -q "TDD"; then
  echo "FAIL: bin/hgf prompt does not contain TDD requirement" >&2
  exit 1
fi

echo "All hgf basic tests passed."
exit 0
