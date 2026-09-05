#!/usr/bin/env bash
set -euo pipefail

# Behavior test for agent command dispatch without Herdr (direct execution in terminal)

tmp_dir=$(mktemp -d "/tmp/test-agent-cmd-mock-XXXXXX")
trap 'rm -rf "$tmp_dir"' EXIT

mock_log="$tmp_dir/mock.log"
export MOCK_LOG="$mock_log"

# Create mock agy binary
cat << 'EOF' > "$tmp_dir/agy"
#!/usr/bin/env bash
set -euo pipefail
echo "agy invoked with: $@" >> "$MOCK_LOG"
exit 0
EOF
chmod +x "$tmp_dir/agy"

# Create mock devin binary
cat << 'EOF' > "$tmp_dir/devin"
#!/usr/bin/env bash
set -euo pipefail
echo "devin invoked with: $@" >> "$MOCK_LOG"
exit 0
EOF
chmod +x "$tmp_dir/devin"

# Create mock opencode binary
cat << 'EOF' > "$tmp_dir/opencode"
#!/usr/bin/env bash
set -euo pipefail
echo "opencode invoked with: $@" >> "$MOCK_LOG"
exit 0
EOF
chmod +x "$tmp_dir/opencode"

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

# Create mock git binary
cat << 'EOF' > "$tmp_dir/git"
#!/usr/bin/env bash
set -euo pipefail
exit 0
EOF
chmod +x "$tmp_dir/git"

# Test 1: Run agy in direct mode (HERDR_ENV=0)
export HERDR_ENV="0"
PATH="$tmp_dir:$PATH" ./bin/hgf work 1 --agent agy > "$tmp_dir/output_agy.txt" 2>&1

if ! grep -q "agy invoked with: -i" "$mock_log"; then
  echo "FAIL: direct agy invocation did not pass '-i' flag" >&2
  cat "$mock_log" >&2
  exit 1
fi

# Test 2: Run devin in direct mode
PATH="$tmp_dir:$PATH" ./bin/hgf work 1 --agent devin > "$tmp_dir/output_devin.txt" 2>&1
if ! grep -q "devin invoked with:" "$mock_log"; then
  echo "FAIL: direct devin invocation was not executed" >&2
  cat "$mock_log" >&2
  exit 1
fi

echo "All agent launch command behavior tests passed."
exit 0
