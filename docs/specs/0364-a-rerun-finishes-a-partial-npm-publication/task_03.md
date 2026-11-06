---
task: task_03
spec: 0364-a-rerun-finishes-a-partial-npm-publication
status: pending
type: infra
complexity: medium
---

# Task 03: A tag run compares each published coordinate with its build and publishes only what is missing

## Overview

After task_02 the preflight lets a coordinate already published at the
target version through and lists it in
`$RUNNER_TEMP/roundfix-release-published`. This Task adds the tag-run step
"Compare published packages", which accepts such a coordinate only when the
decompressed tar stream of its published tarball equals that of this run's
packed build, and makes "Publish to npm" skip each coordinate it accepted.
On 2026-11-06 the v0.132.1 platform tarballs rebuilt on another machine with
the CI's Go 1.26.10 had byte-identical tar streams but different gzip bytes,
so the comparison never uses the tarball integrity (ADR-0374).

## Requirements

1. MUST add to `.github/workflows/release.yml` a step named exactly
   `Compare published packages`, with
   `if: ${{ github.event_name == 'push' }}`, right after "Cross-compile and
   stage" and right before "Publish to npm", preceded by a comment that says
   it accepts an already published coordinate only on byte-identical
   decompressed tar streams and why not the integrity. Under
   `set -euo pipefail`, with `TAG="${GITHUB_REF_NAME#v}"`, it:
   - empties `$RUNNER_TEMP/roundfix-release-complete`, works in a
     `mktemp -d` directory removed by an `EXIT` trap, and appends
     `printf '\n### Already published\n\n'` to `$GITHUB_STEP_SUMMARY`;
   - when `$RUNNER_TEMP/roundfix-release-published` is empty, appends
     `No coordinate was already published.` and exits 0;
   - for each listed coordinate, derives the package directory as the
     publish step does (`dist/npm/roundfix` for `roundfix`, else
     `dist/npm/packages/cli-${coordinate#@roundfix/cli-}`), runs
     `npm pack --json --pack-destination <dir>` in it and
     `npm pack "$coordinate@$TAG" --json` in a separate directory, reads
     `.[0].filename` and `.[0].integrity` of each with `jq -er`, decompresses
     both with `gzip -dc` and compares the two tar streams with `cmp -s`;
   - prints exactly, for an equal pair,
     `printf '%-36s %-12s %s\n' "$coordinate" "$TAG" "same content"`,
     appends the coordinate to the complete file and appends
     ``- `<coordinate>@<version>`: same content as this build`` to the
     summary;
   - prints exactly, for a different pair,
     `echo "::error::registry: $coordinate@$TAG is already published with different content (registry $registry_integrity, this build $local_integrity)"`;
     for a failed local pack or field read
     `echo "::error::undetermined: $coordinate@$TAG: this build could not be packed"`;
     for a failed registry pack or field read
     `echo "::error::undetermined: $coordinate@$TAG: the published tarball could not be fetched"`;
     for a failed decompression
     `echo "::error::undetermined: $coordinate@$TAG: a tarball could not be decompressed"`;
     each of those keeps the coordinate out of the complete file, the loop
     goes on to the next coordinate, and the step exits 1 after the loop.
2. MUST change "Publish to npm": set `TAG="${GITHUB_REF_NAME#v}"` and
   `COMPLETE="$RUNNER_TEMP/roundfix-release-complete"` before
   `publish_coordinate`, whose first statement returns 0 after printing
   exactly `echo "skipping $coordinate@$TAG: already published with the same content"`
   when the complete file exists and holds the coordinate as a whole line
   (`grep -qxF --`). Everything else in the step keeps its bytes: the order
   (platforms, then `roundfix`), the identity grep, the fallback window and
   the summary.
3. MUST create `internal/verifyselect/release_workflow_compare_test.go`,
   headed with a `// Suite:` comment like its neighbours and using the
   harness task_02 wrote in `release_workflow_test.go` (same package; it
   does not change that file), with fixture tarballs built by `archive/tar` and
   `compress/gzip`:
   - `TestReleaseCompareAcceptsOnlyTheSameContent`, a table over the step
     "Compare published packages" for version `0.132.1`: (a) the same tar
     stream, the local copy at `gzip.BestSpeed` and the registry copy at
     `gzip.BestCompression`, with different integrity strings: exit 0, a
     `same content` row, the coordinate in the complete file; (b) a
     registry tar stream whose file body differs: exit 1 with the exact
     `different content` line naming both integrities, the coordinate absent
     from the complete file; (c) no registry tarball: exit 1 with the
     `could not be fetched` line; (d) an empty published file: exit 0 and the
     summary line `No coordinate was already published.`;
   - `TestReleaseRerunPublishesOnlyTheMissingCoordinates`, which runs
     "Publication preflight", "Compare published packages" and "Publish to
     npm" in sequence in one temporary repository and `RUNNER_TEMP`: the
     v0.132.1 shape (five platform coordinates published with the same
     content, `roundfix` eligible) exits 0 with five `skipping` lines and a
     publish log of exactly one entry, `roundfix`; the v0.130.0 shape
     (`@roundfix/cli-darwin-arm64` published with the same content, the rest
     eligible) publishes `cli-darwin-x64`, `cli-linux-arm64`,
     `cli-linux-x64`, `cli-win32-x64` and `roundfix` in that order; a
     published coordinate with different content stops at the compare step
     and the publish log stays empty;
   - `TestReleaseComparesPublishedPackagesAfterTheBuildAndBeforePublishing`,
     which asserts from the parsed workflow that the step exists once, right
     after "Cross-compile and stage" and right before "Publish to npm", with
     the push-only `if`.
4. MUST follow the rules task_02's Requirement 2 lists for the new file
   (`t.Parallel()` first, the RequireTool path, no literal executable mode)
   and keep `go test -timeout 20m -count=1 ./internal/testfixture` passing.
5. MUST keep task_02's tests and `TestMainAndReleaseRunEveryContract`
   passing; MUST NOT change `internal/verifyselect/release_workflow_test.go`, "Publication preflight", "Cross-compile and
   stage", "GitHub Release", any other workflow, the Makefile, `go.mod`,
   `go.sum`, or any file under `internal/suiteguardcontract`.

## Subtasks

- [ ] Add the step "Compare published packages" with its lines and summary.
- [ ] Make the publish step skip a coordinate that compared equal.
- [ ] Add the compare table, the rerun sequence and the step-order tests.

## Acceptance Criteria

- [ ] A published coordinate with the same decompressed content is skipped, and the rerun publishes only the missing coordinates in Release Set order.
- [ ] A published coordinate with different content, or a tarball that cannot be read, stops the job before any publish.
- [ ] Two tarballs that differ only in their gzip bytes compare equal.

## Context

- instruction: `docs/adr/0374-a-rerun-completes-a-release-set-from-the-same-content.md`
- instruction: `docs/user-guide/release-runbook.md`
- instruction: `internal/verifyselect/full_contract_run_test.go`
- interface: `.github/workflows/release.yml`
- creates: `internal/verifyselect/release_workflow_compare_test.go`

## Consumers

- The step "Compare published packages", `$RUNNER_TEMP/roundfix-release-complete` and the publish skip — consumers: `./internal/verifyselect`

## Verification

- `for t in TestReleaseCompareAcceptsOnlyTheSameContent TestReleaseRerunPublishesOnlyTheMissingCoordinates TestReleaseComparesPublishedPackagesAfterTheBuildAndBeforePublishing; do grep -q "func $t(" internal/verifyselect/release_workflow_compare_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; grep -qF -- 'name: Compare published packages' .github/workflows/release.yml || exit 1; go test -timeout 20m -count=1 ./internal/testfixture || exit 1; out="$(mktemp)"; go test -timeout 20m -count=1 -v -run '^(TestReleaseCompareAcceptsOnlyTheSameContent|TestReleaseRerunPublishesOnlyTheMissingCoordinates|TestReleaseComparesPublishedPackagesAfterTheBuildAndBeforePublishing|TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne|TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry|TestMainAndReleaseRunEveryContract)$' ./internal/verifyselect > "$out" 2>&1; status=$?; cat "$out"; passed=0; grep -q -- '--- PASS: TestReleaseCompareAcceptsOnlyTheSameContent' "$out" && grep -q -- '--- PASS: TestReleaseRerunPublishesOnlyTheMissingCoordinates' "$out" && passed=1; rm -f "$out"; test "$status" -eq 0 && test "$passed" -eq 1` — expected: exit 0. Before this Task the three tests and the step do not exist, so the first check fails; afterwards the repository test rules accept the new file, the compare table, the rerun sequence and the step order pass and are not skipped, and task_02's tests still pass.

## References

- `_prd.md` → Goals; User Stories 1-2; Core Features 2-3; Success Metrics 1-3
- `_techspec.md` → Interfaces; API Contracts 2-3; Testing Approach; Build Order 3
- ADR-0374
