#!/usr/bin/env bash
# Upstream drift check (feature 013, roadmap 010).
# Development and CI tooling only: for each pinned upstream project in
# testdata/upstream/pins.json, resolve a release (the latest published,
# non-prerelease one, or a given tag), download its template files into a
# temporary directory, and run the existing drift check (TestDrift) against
# the pin. Downloaded files are untrusted text: they are only read by the Go
# test and are never executed, sourced or committed. The spec binary itself
# never fetches anything (constitution Art. II).
#
# Usage:
#   scripts/check-upstream-drift.sh [--project NAME] [--tag TAG]
#                                   [--candidate NAME=FILE] [--path NAME=REPO_PATH]
# Environment (used when the matching flag is absent; the workflow sets them):
#   DRIFT_PROJECT  project name, or "all" / empty for every project
#   DRIFT_TAG      release tag (requires a single project)
#   GH_TOKEN       token for gh api (optional locally if gh is logged in)
#   GITHUB_STEP_SUMMARY  if set, the Markdown summary is appended to it
# Exit: 0 when every project is "same tag as the pin" or "no unresolved
# change"; 1 when any project has drift, a missing template, a fetch error
# or a check error (the drift test itself failed to run),
# or on a usage error.
# Needs: go and jq (pins.json is read with jq); plus gh and curl unless every
# selected project uses --candidate.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
pins="$repo_root/testdata/upstream/pins.json"

usage() {
  sed -n '11,13p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//' >&2
  exit 1
}

project="${DRIFT_PROJECT:-}"
tag="${DRIFT_TAG:-}"
candidates=() # NAME=FILE
paths=()      # NAME=REPO_PATH
while [ $# -gt 0 ]; do
  case "$1" in
    --project) [ $# -ge 2 ] || usage; project="$2"; shift 2 ;;
    --tag) [ $# -ge 2 ] || usage; tag="$2"; shift 2 ;;
    --candidate) [ $# -ge 2 ] || usage; candidates+=("$2"); shift 2 ;;
    --path) [ $# -ge 2 ] || usage; paths+=("$2"); shift 2 ;;
    -h|--help) usage ;;
    *) echo "unknown argument: $1" >&2; usage ;;
  esac
done
[ "$project" = "all" ] && project=""

command -v jq >/dev/null || { echo "FAIL: jq not found" >&2; exit 1; }
command -v go >/dev/null || { echo "FAIL: go not found" >&2; exit 1; }

all_projects="$(jq -r '.projects[].name' "$pins")"
if [ -n "$project" ]; then
  printf '%s\n' "$all_projects" | grep -qxF -- "$project" \
    || { echo "unknown project: $project (known: $(printf '%s' "$all_projects" | tr '\n' ' '))" >&2; exit 1; }
  selected="$project"
else
  selected="$all_projects"
fi
if [ -n "$tag" ] && [ -z "$project" ]; then
  echo "--tag (or DRIFT_TAG) needs a single project" >&2
  exit 1
fi
# tag_problem TAG prints why a tag is unsafe to use in an API path or URL,
# or nothing when it is fine. Used for --tag and for the tag that
# releases/latest returns (upstream-controlled text).
tag_problem() {
  if ! printf '%s' "$1" | grep -Eq '^[A-Za-z0-9._/-]+$'; then
    echo "tag has unsupported characters"; return 0
  fi
  case "$1" in
    /*|*/|*..*|*//*) echo "tag must not contain '..', '//' or a leading or trailing '/'" ;;
  esac
}
if [ -n "$tag" ]; then
  problem="$(tag_problem "$tag")"
  [ -z "$problem" ] || { echo "$problem: $tag" >&2; exit 1; }
fi

# #6, #4: every --candidate and --path must name a selected, single-template
# project, and --candidate cannot be combined with --tag for that project.
check_overrides() {
  local flag="$1"; shift
  local kv n
  for kv in "$@"; do
    case "$kv" in *=*) ;; *) echo "$flag needs NAME=VALUE: $kv" >&2; exit 1 ;; esac
    n="${kv%%=*}"
    printf '%s\n' "$selected" | grep -qxF -- "$n" \
      || { echo "$flag names a project that is not selected: $n" >&2; exit 1; }
    if [ "$(jq -r --arg n "$n" '.projects[] | select(.name==$n) | .templates | length' "$pins")" != "1" ]; then
      echo "$flag supports only projects with one template: $n" >&2; exit 1
    fi
    if [ "$flag" = "--candidate" ] && [ -n "$tag" ]; then
      echo "--candidate and --tag cannot be combined: the local file has no tag" >&2; exit 1
    fi
  done
}
check_overrides --candidate ${candidates[@]+"${candidates[@]}"}
check_overrides --path ${paths[@]+"${paths[@]}"}

# lookup NAME KEY in a NAME=VALUE list; prints VALUE or nothing.
lookup() {
  local name="$1"; shift
  local kv
  for kv in "$@"; do
    case "$kv" in "$name="*) printf '%s' "${kv#*=}"; return 0 ;; esac
  done
  return 0
}

tmp="$(mktemp -d)"
trap 'rm -rf "${tmp:?}"' EXIT

rows=""
details=""
failed=0

add_row() { rows="${rows}| $1 | $2 | $3 | $4 | $5 |"$'\n'; }

for name in $selected; do
  repo="$(jq -r --arg n "$name" '.projects[] | select(.name==$n) | .repo' "$pins")"
  pinned_tag="$(jq -r --arg n "$name" '.projects[] | select(.name==$n) | .tag' "$pins")"
  ntemplates="$(jq -r --arg n "$name" '.projects[] | select(.name==$n) | .templates | length' "$pins")"
  local_file="$(lookup "$name" ${candidates[@]+"${candidates[@]}"})"
  override_path="$(lookup "$name" ${paths[@]+"${paths[@]}"})"
  cand_tag="local file"
  commit="-"
  dir="$tmp/$name"
  mkdir -p "$dir"

  if [ -z "$local_file" ]; then
    want="$tag"
    if [ -z "$want" ]; then
      if ! want="$(gh api "repos/$repo/releases/latest" -q .tag_name 2>"$dir/err")" || [ -z "$want" ]; then
        add_row "$name" "$pinned_tag" "?" "-" "fetch error: latest release of $repo: $(head -c 200 "$dir/err" | tr '\n|' '  ')"
        failed=1; continue
      fi
    fi
    problem="$(tag_problem "$want")"
    if [ -n "$problem" ]; then
      add_row "$name" "$pinned_tag" "(rejected)" "-" "fetch error: release tag of $repo rejected: $problem"
      failed=1; continue
    fi
    if ! commit="$(gh api "repos/$repo/commits/$want" -q .sha 2>"$dir/err")" || ! printf '%s' "$commit" | grep -Eq '^[0-9a-f]{40}$'; then
      add_row "$name" "$pinned_tag" "$want" "-" "fetch error: tag $want of $repo: $(head -c 200 "$dir/err" | tr '\n|' '  ')"
      failed=1; continue
    fi
    cand_tag="$want"
  fi

  result=""
  unresolved=0
  lines=""
  i=0
  while [ "$i" -lt "$ntemplates" ]; do
    rpath="$(jq -r --arg n "$name" --argjson i "$i" '.projects[] | select(.name==$n) | .templates[$i].path' "$pins")"
    [ -n "$override_path" ] && rpath="$override_path"
    file="$dir/template-$i.md"
    if [ -n "$local_file" ]; then
      cp -- "$local_file" "$file"
    else
      code="$(curl -sS --max-time 30 -o "$file" -w '%{http_code}' "https://raw.githubusercontent.com/$repo/$commit/$rpath" 2>"$dir/err" || true)"
      if [ "$code" = "404" ]; then
        result="template not found: $rpath"; break
      elif [ "$code" != "200" ]; then
        result="fetch error: $rpath at $cand_tag (HTTP ${code:-none}: $(head -c 200 "$dir/err" | tr '\n|' '  '))"; break
      fi
    fi
    out="$dir/drift-$i.log"
    status=0
    (cd "$repo_root" && CGO_ENABLED=0 go test ./internal/conformance -run 'TestDrift$' -count=1 -v \
        -args -project="$name" -candidate="$file" -template="$i") >"$out" 2>&1 || status=$?
    # The test must have run: "no tests to run" or a skip also exits 0.
    if ! grep -Eq '^--- (PASS|FAIL): TestDrift ' "$out"; then
      result="check error: TestDrift did not run for $rpath ($(grep -m1 -E 'no tests to run|SKIP|FAIL|panic|rror' "$out" | head -c 200 | tr '|' ' '))"; break
    fi
    if [ "$status" -ne 0 ] && ! grep -q 'unresolved shape changes' "$out"; then
      result="check error: drift check failed for $rpath ($(grep -m1 -E 'FAIL|panic|rror' "$out" | head -c 200 | tr '|' ' '))"; break
    fi
    found="$(grep -E '(added|removed):.*\(unresolved\)$' "$out" | sed -E 's/^[[:space:]]*drift_test\.go:[0-9]+: //' || true)"
    if [ -n "$found" ]; then
      n="$(printf '%s\n' "$found" | grep -c .)"
      unresolved=$((unresolved + n))
      lines="${lines}${rpath}:"$'\n'"${found}"$'\n'
    fi
    i=$((i + 1))
  done

  if [ -z "$result" ]; then
    if [ "$unresolved" -gt 0 ]; then
      if [ "$unresolved" -eq 1 ]; then result="1 unresolved change"; else result="$unresolved unresolved changes"; fi
      # Upstream text is untrusted: use a fence longer than any backtick run
      # in it, so it cannot close the code block in the summary.
      longest="$(printf '%s' "$lines" | grep -o '`*' | awk '{ if (length($0) > m) m = length($0) } END { print m + 0 }')"
      fence='```'
      while [ "${#fence}" -le "$longest" ]; do fence="${fence}\`"; done
      details="${details}"$'\n'"### $name: unresolved changes ($cand_tag)"$'\n\n'"${fence}"$'\n'"${lines}${fence}"$'\n'
    elif [ "$cand_tag" = "$pinned_tag" ]; then
      result="same tag as the pin"
    else
      result="no unresolved change"
    fi
  fi
  case "$result" in
    "same tag as the pin"|"no unresolved change") ;;
    *) failed=1 ;;
  esac
  short="$commit"; [ "$commit" != "-" ] && short="${commit:0:7}"
  add_row "$name" "$pinned_tag" "$cand_tag" "$short" "$result"
done

summary="## Upstream drift check"$'\n\n'"| Project | Pinned | Candidate | Commit | Result |"$'\n'"|---|---|---|---|---|"$'\n'"${rows}${details}"
if [ "$failed" -ne 0 ]; then
  summary="${summary}"$'\n'"A pin moves only through the manual procedure in testdata/upstream/README.md."$'\n'
fi
printf '%s' "$summary"
if [ -n "${GITHUB_STEP_SUMMARY:-}" ]; then
  printf '%s' "$summary" >>"$GITHUB_STEP_SUMMARY"
fi
exit "$failed"
