---
task: task_02
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
status: pending
type: backend
complexity: high
---

# Task 02: GC sanitation classifies the Artifact Root each Run wrote to

## Overview

`roundfix gc sanitize` groups Runs by the Artifact Directory they recorded, so
every Run that recorded none lands in one group `""` that its first proof
refuses, and the default roots under `~/.roundfix/artifacts/` are never reached
(`_prd.md` → Sources 1). This Task groups Runs by the root the shared resolver
returns, classifies each resolved root with the existing proofs, checks a
directory proven only by absence again just before removing it, and makes the
dry run name what it resolved. It is verifiable on its own on temporary
Roundfix Homes.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-gc-sanitize-preserves-every-artifact-root-a-run-recorded-empty.md`](references/2026-10-20-gc-sanitize-preserves-every-artifact-root-a-run-recorded-empty.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 4 to 7 of
   `_techspec.md` in `internal/store/journal.go` and `internal/cli/gc.go`,
   with `DiscoverArtifactRoots(ctx, runStore, homeDir)` and the `Unresolved`
   field as `_techspec.md` → Interfaces sketches them.
2. MUST group Runs by the root `store.ResolveRunArtifactRoot` returns and add
   no other resolution: an unsafe Run ID joins the group of its recorded value
   so today's proofs report it; any other resolver error sets the group's
   `Unresolved`, and that group is classified `unsafe` with the evidence
   `preserved because the Artifact Root of <Unresolved> cannot be resolved`
   before every other proof.
3. MUST keep every existing proof, in its order and with its evidence text, on
   the resolved root: clean absolute path, safe Run IDs, inside Roundfix Home
   and not equal to it, existing, physical root, physical `runs/`, and the
   default-root proof against every Run of the group. Only `active` and
   `orphaned` roots yield candidates, only `runs/<run-id>` entries are ever
   candidates, and `allRunIDs` maps each Run ID to its group's resolved root.
4. MUST state the hazard in the code's behavior, as API Contract 6 does. The
   data source is the Run Database and the directories under each resolved
   root. With `--apply`, the destination is the removal of `runs/<run-id>`
   directories only. Immediately before removing a directory proven only by
   its Run's absence, the command reads that Run ID through a reader it opens
   with the GC dependencies' reader opener; a Run found there keeps its
   directory, reported `preserved: <path> (Run <id> appeared in durable Run
   metadata before removal)` and not counted. A read error fails the command
   before that removal.
5. MUST print the lines of API Contract 7: one `Recorded as: "<value>" by <n>
   Run(s)` line per distinct recorded value that differs from the root,
   sorted, after the `Artifact Root:` line; and, after the candidate and
   preserved lines, `Directories reclaimable here:` and `Artifact bytes
   reclaimable here:` in a dry run, or `Directories removed here:` and
   `Artifact bytes reclaimed here:` with `--apply`. The header, cutoff, totals
   and closing lines keep their text.
6. MUST declare this break: `TestDiscoverArtifactRootsReturnsEveryRecordedRepository`
   in `internal/store/journal_test.go` passes the test's home and keeps its
   assertions. Every other existing test keeps its name and passes unchanged,
   including `TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories`,
   `TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval` and the five tests in
   `internal/cli/gc_sanitize_pre_key_root_test.go`.
7. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows, using temporary homes and repositories, and never
   reading `~/.roundfix`:
   - in `internal/store/journal_test.go`,
     `TestDiscoverArtifactRootsGroupsRunsByTheRootTheyWroteTo`, with the cases
     `_techspec.md` → Testing Approach lists;
   - in `internal/cli/gc_sanitize_resolved_root_test.go`:
     - `TestGCSanitizeDryRunReachesTheDefaultRootTranscript`, asserting Surface
       Transcript 1 of `_techspec.md` whole: every standard output line in
       order with the home, repository ID, repository, Run ID, cutoff and byte
       count substituted, empty standard error, and exit `0`;
     - `TestGCSanitizeReachesTheDefaultRootAndKeepsEveryProof`, with the cases
       `_techspec.md` → Testing Approach lists, asserting each kept path still
       exists after `--apply`, each removed path is gone, and a second
       `--apply` reports `Directories removed: 0`;
     - `TestGCSanitizeKeepsADirectoryWhoseRunAppearsBeforeRemoval`, which
       replaces the GC dependencies' reader opener so that its second call
       first gives the orphan's Run ID a row in the Run Database, and asserts
       that the directory stays and is reported with the `appeared in durable
       Run metadata before removal` reason;
     - `TestGCSanitizePreservesARunWhoseArtifactRootCannotBeResolved`, a Run
       whose Git and repository roots were blanked in the Run Database,
       asserting `Classification: unsafe`, the new evidence, and no removal.
8. MUST document the change in `docs/user-guide/commands/gc.md` in a new `##
   Sanitation` section that holds the sentences "`roundfix gc sanitize`
   resolves each Run's recorded Artifact Directory to the Artifact Root the
   Run wrote to and classifies that root with the proofs below." and "A
   directory proven only because its Run is absent from the Run Database is
   checked again immediately before removal.", lists the six classifications
   and what each keeps, and shows the `Recorded as` and per-root lines. MUST
   search the guides, the agent guides and the skill copies for wording that
   says sanitation classifies the recorded root as written, and change each
   one this Task owns; the Roundfix `storage` reference is task_03's.
9. MUST re-record the Behavior Surface Record, because this Task edits a user
   guide page: `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back; no record is hand-edited.

## Subtasks

- [ ] Group discovery by the resolved root and mark unresolved groups.
- [ ] Put the unresolved proof in front of the existing proofs.
- [ ] Check absence again before each removal it proves.
- [ ] Print the recorded values and the per-root reclaim lines.
- [ ] Add the store and command tests, the guide section, and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A Run that recorded no Artifact Root is classified under its repository's default root and its eligible directory is removed by `--apply`.
- [ ] Every directory the existing proofs keep is still kept on the resolved root, and nothing outside `runs/` is touched.
- [ ] A directory whose Run appears before removal is kept.
- [ ] The dry run names each resolved root, its recorded values and what it would reclaim, exactly as Surface Transcript 1.

## Consumers

- `store.DiscoverArtifactRoots` and `store.ArtifactRoot` — consumers: `./internal/store`, `./internal/cli`
- the `gc sanitize` report — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0301-gc-sanitation-classifies-the-artifact-root-each-run-wrote-to.md`
- instruction: `docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md`
- instruction: `internal/cli/gc_test.go`
- instruction: `internal/cli/gc_sanitize_pre_key_root_test.go`
- instruction: `internal/cli/gc_run_retention_test.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/cli/gc.go`
- creates: `internal/cli/gc_sanitize_resolved_root_test.go`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDiscoverArtifactRootsGroupsRunsByTheRootTheyWroteTo|TestDiscoverArtifactRootsReturnsEveryRecordedRepository|TestResolveRunArtifactRootIsTheRootTheRunWroteTo|TestStorageReportMeasuresTheArtifactRootARunWroteTo)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDiscoverArtifactRootsGroupsRunsByTheRootTheyWroteTo TestDiscoverArtifactRootsReturnsEveryRecordedRepository; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestGCSanitizeDryRunReachesTheDefaultRootTranscript|TestGCSanitizeReachesTheDefaultRootAndKeepsEveryProof|TestGCSanitizeKeepsADirectoryWhoseRunAppearsBeforeRemoval|TestGCSanitizePreservesARunWhoseArtifactRootCannotBeResolved|TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories|TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval|TestGCSanitizeReclaimsAPreKeyDefaultRootInABareLayout|TestGCSanitizeStillPreservesARootEqualToNeitherDefault|TestGCSanitizeKeepsAKeyDerivedRootClassification|TestGCSanitizeAcceptsTheCheckoutDefaultWhenTheKeyDefaultCannotBeDerived|TestGCSanitizePreservesUnsafeWhenTheKeyDefaultCannotBeDerived|TestRunRetentionReadsTheRootTheSharedResolverReturns|TestGCKeepsARunUnderAnUnprovenArtifactRoot)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGCSanitizeDryRunReachesTheDefaultRootTranscript TestGCSanitizeReachesTheDefaultRootAndKeepsEveryProof TestGCSanitizeKeepsADirectoryWhoseRunAppearsBeforeRemoval TestGCSanitizePreservesARunWhoseArtifactRootCannotBeResolved TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval TestGCSanitizeStillPreservesARootEqualToNeitherDefault; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/gc.md | grep -qF -- 'classifies that root with the proofs below' || { echo 'missing Sanitation sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/gc.md | grep -qF -- 'is checked again immediately before removal' || { echo 'missing recheck sentence' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/store ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new sanitation tests do not exist, so the command fails at its first pass check. Afterwards discovery groups by the resolved root, the transcript, proof, recheck and unresolved tests pass with every existing sanitation and Run Retention test, the guide documents the section, both packages build for Windows, and the Behavior Surface Record is current.

## References

- `_prd.md` → Sources; Goals; Core Feature 2; Core Feature 3; Core Feature 4; Success Metric 1; Success Metric 2
- `_techspec.md` → API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; Surface Transcript 1; Invariants 3-4, 6; Interfaces; Data Models; Current behavior; Testing Approach; Build Order 2
- ADR-0301; ADR-0033; ADR-0255; ADR-0259; ADR-0271; ADR-0285; ADR-0280
