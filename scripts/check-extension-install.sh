#!/usr/bin/env bash
# Install check for the spec-status Spec Kit extension (feature 005, FR-018).
# Development tooling only: creates a throwaway Spec Kit project in a temporary
# directory, installs extensions/spec-status into it with `specify extension
# add --dev`, and checks that exactly the five skills are rendered for Claude.
# It never touches this repository's own .specify/ or .claude/.
#
# Usage: scripts/check-extension-install.sh
# Needs: specify (Spec Kit CLI, 1.2.0 or newer) on PATH.
set -euo pipefail

repo="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ext="$repo/extensions/spec-status"
command -v specify >/dev/null || { echo "FAIL: specify (Spec Kit CLI) not found on PATH" >&2; exit 1; }

tmp="$(mktemp -d)"
trap 'rm -rf "${tmp:?}"' EXIT
start=$(date +%s)

cd "${tmp:?}"
specify init --here --force --non-interactive --integration claude >"$tmp/init.log" 2>&1 \
  || { echo "FAIL: specify init" >&2; tail -20 "$tmp/init.log" >&2; exit 1; }
specify extension add --dev "$ext" >"$tmp/add.log" 2>&1 \
  || { echo "FAIL: specify extension add --dev" >&2; tail -20 "$tmp/add.log" >&2; exit 1; }

skills_dir="$tmp/.claude/skills"
# Portable to bash 3.2 (macOS): no mapfile.
found="$(find "$skills_dir" -maxdepth 2 -path '*speckit-spec-status-*/SKILL.md' 2>/dev/null | sort)"
count="$(printf '%s' "$found" | grep -c . || true)"
want=(review explain reconcile unblock verify-readiness)
if [ "$count" -ne "${#want[@]}" ]; then
  echo "FAIL: expected ${#want[@]} spec-status skills, found $count:" >&2
  printf '%s\n' "$found" >&2
  ls "$skills_dir" >&2 || true
  exit 1
fi
for name in "${want[@]}"; do
  f="$skills_dir/speckit-spec-status-$name/SKILL.md"
  [ -f "$f" ] || { echo "FAIL: missing skill speckit-spec-status-$name" >&2; exit 1; }
done
# Roles must be discoverable from the rendered skill. Spec Kit 1.2.0 does not
# carry argument-hint into extension skills (observed 2026-11-12), so the
# description names the roles instead.
grep -q '^description:.*pm, developer, qa or architect' "$skills_dir/speckit-spec-status-review/SKILL.md" \
  || { echo "FAIL: review skill description does not name the roles" >&2; exit 1; }
grep -q 'reference/rules.md' "$skills_dir/speckit-spec-status-review/SKILL.md" \
  || { echo "FAIL: review skill lost its pointer to reference/rules.md" >&2; exit 1; }
[ -f "$tmp/.specify/extensions/spec-status/reference/rules.md" ] \
  || { echo "FAIL: reference files were not installed" >&2; exit 1; }

elapsed=$(( $(date +%s) - start ))
echo "OK: 5 spec-status skills rendered (review, explain, reconcile, unblock, verify-readiness); roles named; reference files installed; ${elapsed}s"
specify --version 2>/dev/null | head -1
