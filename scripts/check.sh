#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

step() {
	printf '\n==> %s\n' "$1"
}

fail() {
	printf 'gate: %s\n' "$1" >&2
	exit 1
}

for tool in go golangci-lint jq; do
	command -v "$tool" >/dev/null 2>&1 || fail "missing $tool; install the tools listed in CONTRIBUTING.md"
done

step "format"
unformatted=""
while IFS= read -r -d '' file; do
	out="$(gofmt -l "$file")"
	if [[ -n "$out" ]]; then
		unformatted="${unformatted}${out}"$'\n'
	fi
done < <(find . -name '*.go' -not -path './.git/*' -not -path './bin/*' -not -path './dist/*' -print0)
if [[ -n "$unformatted" ]]; then
	printf '%s' "$unformatted" >&2
	fail "gofmt required; run 'make fmt'"
fi

step "vet"
go vet ./...

step "test"
go test ./...

step "build"
go build ./...

step "lint"
golangci-lint run ./...

step "config"
shopt -s nullglob
configs=(config/*.json)
[[ ${#configs[@]} -gt 0 ]] || fail "no config/*.json files found"
for file in "${configs[@]}"; do
	jq -e . "$file" >/dev/null
done
test -f config/AGENTS.md || fail "config/AGENTS.md is missing"
while IFS= read -r -d '' file; do
	grep -q '^description:' "$file" || fail "missing description: $file"
done < <(find config/skills -name SKILL.md -print0)

step "agent contract"
tmp="$(mktemp -d "${TMPDIR:-/tmp}/rutter-check.XXXXXX")"
trap 'rm -rf "$tmp"' EXIT
go build -trimpath -o "$tmp/rutter" ./cmd/rutter
spec="$("$tmp/rutter" agent --json)"
printf '%s' "$spec" | jq -e '.ok == true and (.data.commands | length) > 0' >/dev/null
printf '%s' "$spec" | jq -e '[.data.commands[] | (.flags // [])[] | .name] | index("json") != null' >/dev/null
set +e
RUTTER_DATA_DIR="$tmp/data" "$tmp/rutter" status --check --json >/dev/null 2>&1
code=$?
set -e
[[ "$code" -eq 3 ]] || fail "expected 'status --check' to exit 3, got $code"

printf '\nlocal gate passed\n'
