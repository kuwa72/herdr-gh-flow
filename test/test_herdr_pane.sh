#!/usr/bin/env bash
set -euo pipefail

# Behavior test for Herdr pane split & agent command dispatch using mock herdr and git

tmp_dir=$(mktemp -d "/tmp/test-herdr-mock-XXXXXX")
trap 'rm -rf "$tmp_dir"' EXIT

mock_log="$tmp_dir/mock.log"
export MOCK_LOG="$mock_log"

# Create mock herdr binary
cat << 'EOF' > "$tmp_dir/herdr"
#!/usr/bin/env bash
set -euo pipefail
echo "herdr called with: $@" >> "$MOCK_LOG"
cmd="${1:-}"
subcmd="${2:-}"

if [ "$cmd" = "pane" ] && [ "$subcmd" = "split" ]; then
  # Return JSON with mock pane_id
  echo '{"id":"cli:pane:split","result":{"pane":{"pane_id":"mock:p99"}},"type":"pane_info"}'
  exit 0
fi

if [ "$cmd" = "pane" ] && [ "$subcmd" = "send-text" ]; then
  echo "send-text target: $3" >> "$MOCK_LOG"
  echo "send-text body: $4" >> "$MOCK_LOG"
  exit 0
fi

echo "mock herdr unknown: $@" >> "$MOCK_LOG"
exit 0
EOF
chmod +x "$tmp_dir/herdr"

# Create mock gh binary
cat << 'EOF' > "$tmp_dir/gh"
#!/usr/bin/env bash
set -euo pipefail
if [ "${1:-}" = "issue" ] && [ "${2:-}" = "view" ]; then
  echo '{"title":"mock issue title","body":"mock issue body","number":1}'
  exit 0
fi
exit 0
EOF
chmod +x "$tmp_dir/gh"

# Create mock git binary to avoid altering repository during test
cat << 'EOF' > "$tmp_dir/git"
#!/usr/bin/env bash
set -euo pipefail
exit 0
EOF
chmod +x "$tmp_dir/git"

# Run hgf work with HERDR_ENV=1 and mocks in PATH
export HERDR_ENV="1"
PATH="$tmp_dir:$PATH" ./bin/hgf work 1 --agent agy > "$tmp_dir/output.txt" 2>&1

# Assertions
if ! grep -q "send-text target: mock:p99" "$mock_log"; then
  echo "FAIL: send-text was not targeted to new pane_id (mock:p99)" >&2
  cat "$mock_log" >&2
  exit 1
fi

if ! grep -q 'send-text body: agy -i' "$mock_log"; then
  echo "FAIL: send-text body does not use 'agy -i'" >&2
  cat "$mock_log" >&2
  exit 1
fi

echo "All herdr pane behavior tests passed."
exit 0
