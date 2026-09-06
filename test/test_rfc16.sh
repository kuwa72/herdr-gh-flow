#!/usr/bin/env bash
set -euo pipefail

rfc="docs/rfc-16-product-ux.md"
design="docs/design.md"

if [ ! -f "$rfc" ]; then
  echo "FAIL: $rfc does not exist" >&2
  exit 1
fi

if [ ! -s "$rfc" ]; then
  echo "FAIL: $rfc is empty" >&2
  exit 1
fi

first_line=$(head -n 1 "$rfc")
if [[ "$first_line" != "#"* ]]; then
  echo "FAIL: $rfc does not start with a Markdown title" >&2
  exit 1
fi

for keyword in "提供形態" "基本UXモデル" "後続タスク"; do
  if ! grep -q "$keyword" "$rfc"; then
    echo "FAIL: $rfc missing keyword '$keyword'" >&2
    exit 1
  fi
done

if ! grep -q "rfc-16-product-ux.md" "$design"; then
  echo "FAIL: $design does not reference the RFC 16 document" >&2
  exit 1
fi

echo "RFC 16 artifacts validated."
exit 0
