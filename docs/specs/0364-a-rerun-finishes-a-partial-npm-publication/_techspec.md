---
spec: 0364-a-rerun-finishes-a-partial-npm-publication
prd: _prd.md
created: 2026-11-06
---

# A rerun finishes a partial npm publication — Technical Spec

## Executive Summary

The Publication Preflight in `.github/workflows/release.yml` gains a
`published` state for a coordinate whose packument lists the target version
under `versions`, writes those coordinates to
`$RUNNER_TEMP/roundfix-release-published` and lets them through; a new
tag-run step, "Compare published packages", runs after "Cross-compile and
stage" and before "Publish to npm", and accepts a published coordinate only
when the decompressed tar stream of its registry tarball equals that of this
run's `npm pack`, writing it to `$RUNNER_TEMP/roundfix-release-complete`;
`publish_coordinate` skips each coordinate in that file. A Go test in
`internal/verifyselect` runs the workflow's own step scripts, read from the
YAML, under `bash` with stub `curl` and `npm` commands, so the rerun path is
proven without the network and without publishing. The trade-off this design
accepts is that a rerun after a new Go patch release is refused rather than
completed, in exchange for a Release Set that never mixes two builds.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the preflight keeps its `curl` read
  of `https://registry.npmjs.org/<coordinate>`; the new step reads published
  tarballs with `npm pack <coordinate>@<version>`, which verifies them
  against the registry's integrity; Trusted Publishing, the token fallback
  and its window are unchanged, and the new step reads no credential.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0374 (new) records the rule.
  ADR-0146: "no package publishes unless every coordinate is eligible for the
  exact target version", amended by ADR-0374 to "eligible or already
  published with the same content"; the fallback is unchanged. The
  contract step stays directly after `make verify` and before the preflight
  (`TestMainAndReleaseRunEveryContract`). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 5 answers of 2026-11-06 quoted in the authorization record;
  no Makefile, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0364-a-rerun-finishes-a-partial-npm-publication/_authorization.md`;
  bounded files: `.github/workflows/release.yml`.

## System Architecture

- `.github/workflows/release.yml`, step "Publication preflight": the
  classifier, a `published` case, the published list file and the
  all-published refusal. It stays after "Verify every contract" and before
  the build, on both triggers.
- `.github/workflows/release.yml`, new step "Compare published packages",
  `if: ${{ github.event_name == 'push' }}`, between "Cross-compile and stage"
  and "Publish to npm".
- `.github/workflows/release.yml`, step "Publish to npm":
  `publish_coordinate` returns early for a coordinate in the complete list.
  The OIDC publish, the identity grep, the fallback window and the summary
  are unchanged.
- `internal/verifyselect/release_workflow_test.go` (new, task_02): the
  harness, which reads the steps from the workflow with `gopkg.in/yaml.v3`
  (already a dependency, as in `full_contract_run_test.go`) and runs their
  `run` text with `bash -e`, and the preflight and 409 tests.
- `internal/verifyselect/release_workflow_compare_test.go` (new, task_03):
  the compare, rerun and step-order tests, on the same harness.
- Docs: `docs/user-guide/release-runbook.md`, `CONTEXT.md` (**Release Set**,
  **Publication Preflight**), ADR-0146's amendment notice, the two Coverage
  Reviews in `docs/references/skill-coverage.json` and the Behavior Surface
  Record.

## Implementation Design

### Interfaces

The classifier the preflight uses (`CLASSIFY_JQ`), with `published` ahead of
`used`:

```jq
if type != "object" or ((.versions? | type) != "object") or ((.time? | type) != "object") then error("invalid packument")
elif (.versions | has($tag)) then "published"
elif (.time | has($tag)) then "used"
elif (.time | has("unpublished")) then "cooldown"
else "eligible" end
```

Files the steps share inside one job, under `$RUNNER_TEMP`:

- `roundfix-release-set` (exists): the six coordinates, platforms first.
- `roundfix-release-published` (new): written by the preflight, one
  coordinate per line, emptied first.
- `roundfix-release-complete` (new): written by the compare step, one
  coordinate per line, emptied first.

The compare step, per published coordinate (package directory as the publish
step derives it: `dist/npm/roundfix` for `roundfix`,
`dist/npm/packages/cli-<suffix>` otherwise):

1. `( cd "$pkgdir" && npm pack --json --pack-destination "$work/local" )`,
   reading `.[0].filename` and `.[0].integrity` with `jq`.
2. `( cd "$work/registry" && npm pack "$coordinate@$TAG" --json )`, reading
   the same two fields.
3. `gzip -dc` of each tarball, then `cmp -s` of the two tar streams.

### Data Models

No changes. The step summary gains a section; no file outside
`$RUNNER_TEMP` is written.

### API Contracts

1. API Contract: preflight lines. A `published` coordinate prints its table
   row with state `published`, adds the same row to the step summary and
   prints `::notice::registry: <coordinate>@<version> is already published;
   a tag run compares it with its own build`; it does not set `blocked`.
   `used` keeps `::error::registry: <coordinate>@<version> is already used`,
   now only for a version under `time` and not under `versions`. After the
   loop, when nothing blocked and every coordinate is `published`, the step
   prints `::error::registry: every coordinate of the Release Set is already
   published at <version>` and exits 1. Every other line keeps its bytes.
2. API Contract: compare step. With an empty published list it appends
   `No coordinate was already published.` under a `### Already published`
   summary heading and exits 0. For each equal coordinate it prints
   `printf '%-36s %-12s %s\n' "$coordinate" "$TAG" "same content"`, appends
   the coordinate to the complete list and a summary bullet
   ``- `<coordinate>@<version>`: same content as this build``. A different
   tar stream prints `::error::registry: <coordinate>@<version> is already
   published with different content (registry <registry integrity>, this
   build <local integrity>)`; a failed local pack prints
   `::error::undetermined: <coordinate>@<version>: this build could not be
   packed`; a failed registry pack prints `::error::undetermined:
   <coordinate>@<version>: the published tarball could not be fetched`; a
   failed decompression prints `::error::undetermined:
   <coordinate>@<version>: a tarball could not be decompressed`. Each of
   those keeps the coordinate out of the complete list; the step checks every
   coordinate and then exits 1 if any failed.
3. API Contract: publish step. `publish_coordinate` prints `skipping
   <coordinate>@<version>: already published with the same content` and
   returns 0 for a coordinate in the complete list, before any `npm`
   call; otherwise it is unchanged, so a `TLOG_CREATE_ENTRY_ERROR` still
   prints `::error::publish: <coordinate> failed for a non-identity reason`
   after one `npm publish` and never reaches the token fallback.

### Documentation

The release runbook gains the rehearsal sentence about a `published`
coordinate, the section "Finishing a partial publication" (rerun the failed
job; the compare step; the skip line; why tar streams and not integrity; a
newer Go patch release or an unpublished version means the next patch) and
three vocabulary rows that replace the two Publication Preflight rows.
`CONTEXT.md` rewrites the definitions of **Release Set** and **Publication
Preflight**; ADR-0146 gains the amendment notice that names ADR-0374; the
coverage map records a Coverage Review on the two surfaces whose source is
the runbook, because the covering release reference needs no text change.

### Surface Transcripts

None. The workflow runs only on GitHub Actions; its lines are asserted by
the step-script tests above, and no Roundfix command output changes.

## Coverage Map

- Goal 1, Story 1, Core Features 1-3, Success Metrics 1-2 → Interfaces (the classifier, the shared files and the compare procedure); API Contracts 1-3; Build Orders 2-3.
- Goal 2, Story 2, Core Feature 2, Success Metric 3 → Interfaces (the compare procedure); API Contract 2; Build Order 3.
- Goal 3, Success Metrics 4-5 → Interfaces (the classifier); API Contracts 1 and 3; Build Order 2.
- Goal 4, Story 3, Core Feature 4 → Documentation; Build Order 1.

Deviations from the PRD: None.

## Absence Evidence

1. The `published` state and its list — command: `grep -nw 'published' .github/workflows/release.yml`; result: no match (exit 1); the classifier returns only `used`, `cooldown` and `eligible`, and `unpublished` is the only related word.
2. The step "Compare published packages" and the complete list — command: `grep -nE 'Compare published packages|roundfix-release-complete|npm pack ' .github/workflows/release.yml`; result: no match (exit 1).
3. A test that runs the workflow's step scripts — command: `grep -rl 'Publication preflight' internal`; result: only `internal/verifyselect/full_contract_run_test.go`, which checks the step order and runs no step.

## ADR Alternatives

- ADR-0374: the design takes none of its rejected options: it compares
  decompressed tar streams, not tarball integrity; it reads the registry's
  tarball, not a recorded artifact; it adds no in-job publish retry; and it
  never skips a published coordinate without comparing it.
- ADR-0146: rejected publishing any coordinate before every coordinate is
  proven; the compare step runs before the first publish and stops the job
  on any difference, so that still holds. Its bounded fallback is untouched.

## Integration Points

- The npm registry, read only: the packument (`curl`, unchanged) and each
  published tarball (`npm pack <coordinate>@<version>`, new). npm verifies a
  fetched tarball against the packument's integrity.
- GitHub Actions: `$RUNNER_TEMP` and `$GITHUB_STEP_SUMMARY`, shared by the
  steps of one job; "Re-run failed jobs" reruns the whole job.

## Testing Approach

- `internal/verifyselect/release_workflow_test.go` (new) holds the harness:
  it parses the workflow, writes a step's `run` text to a temporary file and
  runs it with `bash -e` in a temporary repository holding copies of
  `dist/npm/platforms.json`, `dist/npm/roundfix/package.json` and each
  `dist/npm/packages/*/package.json`. The environment is built, not
  inherited: `PATH` starts with an absolute stub directory, and
  `RUNNER_TEMP`, `GITHUB_STEP_SUMMARY`, `TARGET_VERSION`, `GITHUB_REF_NAME`,
  `GITHUB_EVENT_NAME`, `FALLBACK_WINDOW` and a fake `NPM_FALLBACK_TOKEN`
  built by concatenation are set; `NODE_AUTH_TOKEN` is never inherited.
  A relative stub directory would let a step's `cd` reach the real `npm`.
- The stub `curl` answers each packument from a fixture directory (`200`
  with the file, `404` without). The stub `npm` serves `pack` of the current
  package and `pack <coordinate>@<version>` from fixture tarballs, printing
  `[{"filename":…,"integrity":…}]`, records every `publish` with its
  directory and whether `NODE_AUTH_TOKEN` is set, and fails a configured
  directory's publish with the configured output.
- Fixture tarballs are built in Go with `archive/tar` and `compress/gzip`:
  the same tar stream gzip-compressed at `gzip.BestSpeed` and at
  `gzip.BestCompression` gives two different integrities for the same
  content; a different file body gives different content.
- `bash`, `jq`, `gzip` and `cmp` are External Test Tools, resolved with
  `testfixture.RequireTool`, so the Linux Lane, which has no `jq`, skips
  the tests naming the tool, and CI (ubuntu-latest) and macOS run them.
- Tests in `release_workflow_test.go`:
  `TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne` and
  `TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry`
  (characterizes today's publish failure). Tests in
  `release_workflow_compare_test.go`:
  `TestReleaseCompareAcceptsOnlyTheSameContent`,
  `TestReleaseRerunPublishesOnlyTheMissingCoordinates` and
  `TestReleaseComparesPublishedPackagesAfterTheBuildAndBeforePublishing`.
- The package already installs the suite guard and is a guarded spawning
  package, so no guard file changes. No test reaches the network, publishes
  or signals a process.

## Build Order

1. Documentation: ADR-0146's amendment notice, `CONTEXT.md`, the runbook,
   the two Coverage Reviews and the Behavior Surface Record. It comes first
   so each workflow Task names its guide through a dependency.
2. The preflight's `published` state, the all-published refusal, the test
   harness, the preflight test and the 409 characterization (depends on: 1).
3. The compare step, the publish skip and their tests (depends on: 2).
4. Final QA gate (depends on: 3).

## Risks & Considerations

- `actions/setup-go` keeps `check-latest: true`; a rerun after a new Go 1.26
  patch release rebuilds different binaries and is refused. A rerun soon
  after the failure resolves the same Go release, as the measured
  reproducibility shows; the runbook says the next patch is the way out.
- The 409's cause is inferred from published sources and the timing, not
  from the npm debug log, which the runs did not keep; the fix does not
  depend on it, because any failure after some coordinates published leaves
  the same state.
- The Linux Lane skips the new tests (no `jq` in the image); CI's
  ubuntu-latest and the macOS Runs execute them.
- Shared files with Spec 0361, which the queue delivers first: `CONTEXT.md`
  (other entries), `docs/references/skill-coverage.json` (other surfaces)
  and `docs/references/behavior-surfaces.json`, re-recorded by its command.
- task_03 adds its tests in a new file rather than extending task_02's,
  because a Task Context path must exist at the Delivery Base; the two files
  share the package's test helpers.

## Glossary

- changes: **Release Set** — a coordinate already published with the same content counts toward the set.
- changes: **Publication Preflight** — a coordinate already published at the target version no longer stops it.

## Decisions

- Compare decompressed tar streams, not the npm integrity. See ADR-0374.
- Keep the preflight before the build and compare after it, on tag runs
  only.
- Run the workflow's own step scripts in Go tests with stub `curl` and
  `npm`, rather than moving the steps into a separate script.
- Three implementation Tasks and the QA gate, in series, documentation
  first.
