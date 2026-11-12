# Quickstart: validate dependency update automation

Run these checks after implementation. The expected results come from the [contracts](contracts/) and the [spec](spec.md).

## 1. Local, before pushing

```sh
grep '^go ' go.mod                                   # expect: go 1.26.0
test -z "$(gofmt -l .)" && go mod tidy -diff && go mod verify
CGO_ENABLED=0 go build -o spec ./cmd/spec && ./spec status --color=never --fail-on error   # exit 0
for t in windows/amd64 linux/arm64; do CGO_ENABLED=0 GOOS=${t%/*} GOARCH=${t#*/} go build ./... || echo FAIL $t; done
GOTOOLCHAIN=go1.26.11 CGO_ENABLED=0 go test -count=1 ./...   # all packages ok
GOTOOLCHAIN=go1.26.11 go vet ./...                   # no output
GOTOOLCHAIN=go1.26.11 go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...   # "No vulnerabilities found."
git diff main -- go.mod go.sum                       # only the go line, plus tidy changes; no new module paths (R7)
```

Check the workflow file against [contracts/ci-workflow.md](contracts/ci-workflow.md):

```sh
grep -E 'uses: ' .github/workflows/ci.yml            # every line: @<40 hex> # vX.Y.Z
grep -E 'secrets\.|pull_request_target|continue-on-error' .github/workflows/ci.yml   # no output
grep -A1 '^permissions' .github/workflows/ci.yml     # contents: read
```

## 2. Byte-identical product output (SC-005)

Build before (`main`) and after (branch), both under go1.26.11, then compare:

```sh
for r in testdata/example-roadmap testdata/speckit-roadmap-upstream; do
  ./spec-before status --root $r --json > a.json; ./spec-after status --root $r --json > b.json; cmp a.json b.json
done
```

## 3. On GitHub, after pushing the branch

- The pull request shows six checks (four test cells, `hygiene`, `dependency-review`), all passing, plus CodeQL's code-scanning check once default setup is configured.
- Seeded hygiene faults (SC-007): on a scratch branch, misformat one file, then separately break a roadmap field. Only `hygiene` fails, each time naming the cause.
- A deliberate failure test: push a temporary commit that breaks one test. All four checks fail. Revert it, and all four pass. Record the run URLs in the PR.
- The repository's Insights → Dependency graph → Dependabot page lists both ecosystems after merge. The first weekly run happens the following Monday.
- Security: `gh api repos/tosin2013/spec-status/dependabot/alerts --jq '.[].security_advisory.ghsa_id'` lists any open alerts. An x/sys alert (GO-2026-5024) is expected. Its fix PR is classified by the [merge rule](data-model.md): it changes only versions, the `go` line is already ≥1.25, so it is merge after review once CI passes.

## 4. Optional owner step

Branch protection on `main` requiring the CI checks (four test cells and `hygiene`, plus `dependency-review` on pull requests) makes the merge rule's "CI passes" mechanical. It is a repository setting, outside this feature's files. Offer it, and don't apply it without approval.
