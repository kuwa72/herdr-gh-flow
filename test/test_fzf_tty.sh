#!/usr/bin/env bash
set -euo pipefail

# Behavior test for issue selection: verify that fzf receives gh issue list output via stdin

tmp_dir=$(mktemp -d "/tmp/test-fzf-mock-XXXXXX")
trap 'rm -rf "$tmp_dir"' EXIT

mock_log="$tmp_dir/mock.log"
fzf_stdin="$tmp_dir/fzf_stdin.txt"
export MOCK_LOG="$mock_log"
export FZF_STDIN="$fzf_stdin"

# Create mock gh binary that outputs distinct issues
cat << 'EOF' > "$tmp_dir/gh"
#!/usr/bin/env bash
set -euo pipefail
cmd="${1:-}"
subcmd="${2:-}"

if [ "$cmd" = "issue" ] && [ "$subcmd" = "list" ]; then
  # Simulate template output
  echo "#42 Test Issue Forty Two [bug]"
  echo "#99 Test Issue Ninety Nine [feature]"
  exit 0
fi

if [ "$cmd" = "issue" ] && [ "$subcmd" = "view" ]; then
  echo '{"title":"mock title","body":"mock body","number":42}'
  exit 0
fi

exit 0
EOF
chmod +x "$tmp_dir/gh"

# Create mock fzf binary that captures stdin and selects the first issue
cat << 'EOF' > "$tmp_dir/fzf"
#!/usr/bin/env bash
set -euo pipefail

cat > "$FZF_STDIN"
echo "fzf called with args: $@" >> "$MOCK_LOG"

# Simulate selecting the first item with Enter key (first line empty/default key, second line selected item)
echo ""
head -n 1 "$FZF_STDIN"
exit 0
EOF
chmod +x "$tmp_dir/fzf"

# Create mock git binary
cat << 'EOF' > "$tmp_dir/git"
#!/usr/bin/env bash
set -euo pipefail
exit 0
EOF
chmod +x "$tmp_dir/git"

# Run hgf in direct mode with mocks in PATH
export HERDR_ENV="0"
PATH="$tmp_dir:$PATH" ./bin/hgf work "" --agent testagent > "$tmp_dir/output.txt" 2>&1

# Assert that fzf actually received the issues from gh via stdin
if ! grep -q "#42 Test Issue Forty Two" "$fzf_stdin"; then
  echo "FAIL: fzf did not receive issue #42 on stdin. Stdin was:" >&2
  cat "$fzf_stdin" >&2
  exit 1
fi

if ! grep -q "#99 Test Issue Ninety Nine" "$fzf_stdin"; then
  echo "FAIL: fzf did not receive issue #99 on stdin. Stdin was:" >&2
  cat "$fzf_stdin" >&2
  exit 1
fi

echo "All issue selection fzf stdin behavior tests passed."
exit 0
