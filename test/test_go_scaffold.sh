#!/usr/bin/env bash
# issue #35: Go足場 (go.mod・cmd/lead・version・CI配線) の振る舞いテスト。
# ソースのgrepではなく、ビルド・実行・出力をアサートする。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

fail() { echo "FAIL: $*" >&2; exit 1; }

# 1. go build が通ること (CGO無効の静的ビルド含む)
CGO_ENABLED=0 go build -o "$tmp/lead" ./cmd/lead || fail "CGO_ENABLED=0 go build ./cmd/lead failed"
[ -x "$tmp/lead" ] || fail "built binary not executable"

# 2. `lead version` が version/commit/date を表示すること (実行して検証)
out="$("$tmp/lead" version)" || fail "lead version exited non-zero"
[ -n "$out" ] || fail "lead version output empty"
case "$out" in *lead*) ;; *) fail "version output missing 'lead': $out";; esac

# 3. `--version` フラグも同一の版数情報を出すこと
flag_out="$("$tmp/lead" --version)" || fail "lead --version exited non-zero"
[ "$flag_out" = "$out" ] || fail "--version output differs from version subcommand"

# 4. ldflags 版数注入配線の確認 (注入値を実行出力で検証)
CGO_ENABLED=0 go build \
  -ldflags "-X main.Version=test-9.9.9 -X main.Commit=abc1234 -X main.Date=2026-09-06" \
  -o "$tmp/lead-stamped" ./cmd/lead || fail "stamped build failed"
stamped_out="$("$tmp/lead-stamped" version)" || fail "stamped lead version failed"
case "$stamped_out" in *test-9.9.9*) ;; *) fail "stamped version missing: $stamped_out";; esac
case "$stamped_out" in *abc1234*) ;; *) fail "stamped commit missing: $stamped_out";; esac
case "$stamped_out" in *2026-09-06*) ;; *) fail "stamped date missing: $stamped_out";; esac

# 5. go vet / go test が通ること
go vet ./... || fail "go vet failed"
go test ./... || fail "go test failed"

echo "go scaffold behavioral checks passed"
