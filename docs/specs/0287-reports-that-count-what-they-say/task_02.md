---
task: task_02
spec: 0287-reports-that-count-what-they-say
status: completed
type: backend
complexity: medium
---

# Task 02: The storage report measures each Run under the root it wrote to

## Overview

`roundfix storage report` listed Artifact Root `.` as `missing` with 508 Run
rows and 0 bytes, because it measures each Run's recorded `artifact_dir`
literally and every Implement Run records an empty one (`_prd.md` → Sources
2). This Task resolves each recorded root the way the Run resolved it when it
wrote its artifacts, inside `Store.StorageReport`, and measures that root. It
is verifiable on its own on temporary Roundfix Homes.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-the-storage-report-misses-the-default-artifact-roots.md`](references/2026-10-20-the-storage-report-misses-the-default-artifact-roots.md)
   (`_prd.md` → Sources 2), by implementing API Contract 3 of `_techspec.md`
   in `internal/store/journal.go`, with `StorageReport(ctx, homeDir)` and
   `storageArtifactRoot` as `_techspec.md` → Interfaces sketches them, and by
   passing `homeDir` from `measureStorageReport` in `internal/cli/gc.go`.
2. MUST resolve each Run's root as API Contract 3 states:
   - a recorded root that is not blank through
     `config.ResolveArtifactDirectory(recorded, gitRoot, homeDir)`, keeping
     the cleaned recorded value when it fails;
   - a blank recorded root to the repository key's default,
     `config.ResolveArtifactDirectory("", repositoryRoot, homeDir)` with
     `repositoryRoot` read as `COALESCE(NULLIF(repository_root, ''),
     git_root)`; to the recorded checkout's default,
     `config.DefaultArtifactDirectoryForPath(gitRoot, homeDir)`, only when it
     differs and holds `runs/<run-id>` while the key's default does not, or
     when the key's default cannot be resolved; and to `.`, reported
     `missing`, when neither resolves.
   The Artifact Root, repository and state groups all measure the resolved
   root and `<resolved root>/runs/<run-id>`.
3. MUST state the hazard in the code's behavior: the data comes from the Run
   Database and the recorded repositories' Git metadata, and goes to standard
   output only. The report MUST NOT call `ValidateArtifactDirectory` or
   create, write or remove any directory or file, in Roundfix Home or under
   any root. The Run ID is checked before any path is built, as
   `storageRunArtifactBytes` checks it, and a stat error other than not-exist
   fails the report with `inspect Run artifact directory <path>: <error>`.
4. MUST keep everything else: the report's sections, columns, totals,
   reconciliation and Storage Snapshot reader; the Repositories grouping by
   recorded Git root; `gc compact`'s live-writer fallback; and
   `DiscoverArtifactRoots`, `roundfix gc sanitize` and Run Retention, which
   this Task does not change.
5. MUST declare this break: the five `StorageReport(ctx)` calls in
   `internal/store/journal_test.go` pass the test's home; every existing test
   keeps its name and passes.
6. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows, using temporary homes and repositories, and never
   reading `~/.roundfix`:
   - in `internal/store/journal_test.go`,
     `TestStorageReportMeasuresTheArtifactRootARunWroteTo`, with subtests:
     - an Implement Run with no recorded root whose repository default root
       holds `runs/<run-id>` bytes: that root is listed present with them, and
       the Run's repository and state groups carry the same bytes;
     - a Run recorded from a linked worktree that was then removed, whose
       repository root is the main checkout: the main checkout's default root;
     - a `fetch` Run whose recorded root is relative: measured under the Run's
       Git root;
     - a Run whose directory exists only under its checkout's default root:
       that root;
     - a blank root whose default does not exist: listed `missing` with 0
       bytes, and the directory still absent after the report;
   - in `internal/cli/gc_test.go`,
     `TestRunStorageReportListsTheDefaultRootOfARunWithoutOne`, through
     `storage report` on a temporary home: the Artifact Roots section holds the
     line made of two spaces, `<home>/.roundfix/artifacts/<repository id>`, a
     tab, `present`, a tab, `1`, a tab and the bytes; no `.` root is listed;
     `Measured bytes` equals those bytes; and the Run Database bytes are
     unchanged.
7. MUST document the change in `docs/user-guide/commands/storage-report.md`
   and `.agents/skills/roundfix/references/storage.md`, each with these two
   sentences: "A Run that recorded no Artifact Root is measured under its
   repository's default root, `~/.roundfix/artifacts/<repository id>`." and "A
   relative Artifact Root is measured under the Run's Git root." MUST search
   the guides, the agent guides and the skill copies for wording that says the
   report measures the recorded root as written, and change each. No text
   inside any `### QA settlement` section may change.
8. MUST run these commands, because this Task edits a user guide page and the
   Roundfix reference. No record or mirror is hand-edited.
   Before them, raise both version fields of the Roundfix skill's canonical
   `SKILL.md` (`metadata.version` and the top-level `version`) by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   these commands itself before handing back: the Daemon only runs the
   Verification and never regenerates a mirror or a record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Read each Run's Git root and repository root with its recorded root.
- [ ] Resolve the root a Run wrote to and measure the groups under it.
- [ ] Pass Roundfix Home from the command and update the store test calls.
- [ ] Add the store and command tests.
- [ ] Document the resolution in the guide and the reference; sync and record.

## Acceptance Criteria

- [ ] A Run that recorded no Artifact Root is measured under its repository's default root, in the Artifact Root, repository and state groups.
- [ ] A relative recorded root is measured under the Run's Git root, and an absolute one as recorded.
- [ ] The report creates no directory and changes no byte.
- [ ] The guide and the reference describe the resolution.

## Consumers

- `store.StorageReport` and its Artifact Root, repository and state groups — consumers: `./internal/store`, `./internal/cli`

## Context

- instruction: `docs/adr/0296-the-storage-report-reconciles-every-page-of-one-storage-snapshot.md`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `internal/config/config.go`
- instruction: `internal/cli/run_retention.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/cli/gc.go`
- interface: `internal/cli/gc_test.go`
- interface: `docs/user-guide/commands/storage-report.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestStorageReportMeasuresTheArtifactRootARunWroteTo|TestStorageReportReconcilesMeasuredTotalsWithoutMutation|TestStorageReportReadsALiveWALSnapshot|TestStorageReportGroupsPointerMapPages|TestStorageReportReleasesReadTransactionOnError)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestStorageReportMeasuresTheArtifactRootARunWroteTo TestStorageReportReconcilesMeasuredTotalsWithoutMutation; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestRunStorageReportListsTheDefaultRootOfARunWithoutOne|TestRunStorageReportOutsideGitRepository|TestRunStorageReportWithoutDatabaseTranscript|TestRunStorageReportWhileAWriterHoldsTheWAL|TestRunGCCompactPreviewReportsStorageMeasurementWhenWriterAdvances)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunStorageReportListsTheDefaultRootOfARunWithoutOne TestRunStorageReportOutsideGitRepository TestRunStorageReportWithoutDatabaseTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/storage-report.md .agents/skills/roundfix/references/storage.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A Run that recorded no Artifact Root is measured under its repository' || { printf 'missing blank-root sentence in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A relative Artifact Root is measured under the Run' || { printf 'missing relative-root sentence in %s\n' "$f" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/store ./internal/cli || exit 1; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the two new tests do not exist, so the command fails at its first pass check. Afterwards the resolution tests pass with the existing storage report, `gc compact` fallback and store tests, the guide and the reference describe the resolution, both packages build for Windows, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 2; Core Feature 4; Success Metric 2
- `_techspec.md` → API Contract 3; API Contract 5; Invariants 3-4; Interfaces; Data Models; Current behavior; Testing Approach; Build Order 2
- ADR-0296; ADR-0255; ADR-0259; ADR-0271; ADR-0285; ADR-0189; ADR-0256; ADR-0280


## Result

Implemented the assigned storage-report slice; status and declared Verification
remain Daemon-owned.

- Empty recorded roots resolve from the stored repository key, with the recorded
  checkout default selected only when its Run directory exists and the key's
  does not, or when repository metadata cannot resolve the key's default.
  `COALESCE(NULLIF(repository_root, ''), git_root)` preserves legacy records.
- Nonblank roots resolve through `config.ResolveArtifactDirectory`; resolution
  errors retain the cleaned recorded value. The CLI passes its explicit home.
  Artifact Root, repository and state groups measure the resolved location;
  repository groups retain the recorded Git root.
- Run IDs are checked before root/path construction. Non-not-exist inspection
  errors carry `inspect Run artifact directory <path>: <error>`. The source
  comment states that database and repository Git metadata feed standard output;
  resolution never invokes writable-directory validation.

Acceptance evidence:

- Default roots: `TestStorageReportMeasuresTheArtifactRootARunWroteTo` passes nine
  subtests covering Implement defaults, removed linked worktrees, relative fetch
  roots, checkout-only defaults, missing defaults, absolute roots, repository
  precedence when both directories exist, differing defaults with neither
  directory present, and unavailable repository metadata. Each asserts Artifact
  Root, repository and state bytes. The new CLI test asserts the resolved root's
  exact tab-separated line, no `.` root, and the measured artifact bytes.
- Relative and absolute roots: the named store test's fetch-relative and absolute
  cases pass; all five existing report calls now receive their test home without
  renaming an existing test.
- No mutation: the new store cases compare the full temporary home and repository
  trees before/after and check missing defaults stay absent. The new CLI test
  compares database bytes. The existing reconciliation/no-mutation test and
  live-writer compaction fallback also pass in the focused check below.
- Documentation: both required sentences appear in the guide and canonical skill
  reference. A search across `docs/user-guide`, `docs/agents`, `.agents/skills`
  and `skills` found no remaining wording claiming this report measures every
  recorded root literally. Sanitation and retention descriptions stay intact.
  Roundfix's two canonical version fields rose from `0.1.114` to `0.1.115`;
  the mirror and records were regenerated, never hand-edited. QA settlement
  sections were left unchanged.

Focused checks and generation (all exit 0 after implementation):

- `go test ./internal/store ./internal/cli -run 'Test(StorageReportMeasuresTheArtifactRootARunWroteTo|RunStorageReportListsTheDefaultRootOfARunWithoutOne|StorageReportReconcilesMeasuredTotalsWithoutMutation|RunGCCompactPreviewReportsStorageMeasurementWhenWriterAdvances)$' -count=1`
- `make skills-sync`
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
- `make baseline-digests` — repository-mandated regeneration; reported no changes.
- `GOOS=windows go test -c -o <temporary directory>/store.test.exe ./internal/store`
  and `GOOS=windows go test -c -o <temporary directory>/cli.test.exe ./internal/cli`
  — both changed test packages compile for Windows.
- `git diff --check`

Before the production change, the new store test reproduced empty roots shown
as `.` and relative roots measured literally. After correcting a missing
fixture-directory setup, the CLI regression independently reproduced
`.\tmissing\t1\t0` instead of its default-root line with 21 bytes. Both
regressions pass after resolution was implemented.

Environment: the supplied writable worktree path contains a symlink unsupported
by the sandbox executor. Commands used authorized escalation in the resolved
assigned checkout. All test homes and repositories were temporary; no test read
the real Roundfix Home.

Declared Verification and repository Verification were not run. No Task status,
Task Graph, other Task file, commit, push or pull request was authored here.

## Carry-forward provenance

- Source Run: `run_20261020T144117Z_258265979932c5f1`
- Source commit: `151468d6271426a711d0b6cb0ace501242977427`
