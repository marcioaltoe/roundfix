---
spec: 0308-recorded-paths-are-read-before-carry-forward-provenance
prd: _prd.md
created: 2026-10-22
---

# Recorded paths are read before carry-forward provenance — Technical Spec

## Executive Summary

`recordedPathsSectionOffset` in `internal/spec/recorded_paths.go` accepts `##
Recorded paths` only when it is the last unfenced level-2 heading of a Task
file. `RecordCarryForward` (`internal/spec/spec.go`) appends `## Carry-forward
provenance` after it, so `RecordedTaskPaths` returns nothing for every carried
Task, and `RecordTaskPaths` appends a second section after the provenance. One
helper replaces that offset: it lists the unfenced level-2 sections and walks
back from the end over the Daemon-written settlement sections (`## Result`,
`## Carry-forward provenance`), collecting every `## Recorded paths` section it
passes and stopping at the first other heading. The reader returns the union of
those sections' paths; the writer replaces the first one in place, removes the
others, and keeps every other byte. Every consumer goes through these two
functions, so no consumer changes.

The trade-off accepted: an author who writes a `## Recorded paths` heading
directly before `## Result` has it read as the Daemon's. ADR-0166 already
reserves that heading for the Daemon, and today the same authored heading is
read whenever it ends the file.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; tests use temporary repositories and a
  local fake reviewer. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0166, and
  ADR-0176 does not apply; both are named with their titles in `_prd.md` →
  Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling file and no
  Governed Path changes; the Spec-contained authorization record
  `docs/specs/0308-recorded-paths-are-read-before-carry-forward-provenance/_authorization.md`
  records `paths: []`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Current behavior

Read and characterized on 2026-10-22 at `42f523c2`. Line numbers move when
other Specs merge; the function names do not.

1. `recordedPathsSectionOffset` (`internal/spec/recorded_paths.go:108`)
   tracks fences with `markdownFenceMarker`, keeps the offset and text of the
   last unfenced level-2 heading, and succeeds only when that heading is
   exactly `## Recorded paths` (`:141`). `RecordedTaskPaths` (`:88`) reads
   `- ` + backtick lines from that offset to the end. `RecordTaskPaths` (`:45`)
   truncates the file at that offset, or keeps it whole without one, and
   appends the new section.
2. `RecordCarryForward` (`internal/spec/spec.go:1036`) appends `\n## Carry-forward
   provenance\n\n- Source Run: ...\n- Source commit: ...\n` after every
   existing byte. Oraculum's Spec 0092 `task_02.md` and `task_03.md` at
   `01cf384` and `b12d75b` end with `## Result`, `## Recorded paths`, `##
   Carry-forward provenance` in that order.
3. Characterization, `go test -overlay` probe in `internal/spec` that wrote
   nothing to the repository: for a synthetic carried Task,
   `RecordedTaskPaths` returned `[]`, `DeclaredTaskPaths` returned only the
   `interface:` path, `RecordTaskPaths` with a new path left two `## Recorded
   paths` headings, and `RecordTaskPaths` with no path left the carried
   section in place.
4. Characterization through the corrective retry, probe in `internal/cli`:
   `newDeclaredCorrectionFixture` with a Task whose recorded `review.txt` is
   followed by `## Carry-forward provenance`, reviewed, disposed `--fixed-by`
   and proved with `ProveReviewCorrection`, returned `accepted=false` and
   `changed path review.txt is outside the named archived Specs and the paths
   their Tasks declared`, for `archive-record` and `legacy-folder`.
5. Characterization through reconcile, probe in `internal/worktree`:
   `newArchivedLeftoverFixture` with the archived Task rewritten to end in the
   carried shape and `recorded.txt` dirty in the Run Worktree returned state
   `dirty` with `recorded.txt is outside the archived Spec's declared and
   recorded paths`; the same probe over a prototype of the helper below
   returned `superseded`.

Sweep of every consumer (`rg -n 'RecordedTaskPaths\(|RecordTaskPaths\(|DeclaredTaskPaths\(|recordedPathsSectionOffset\(' internal -g '!*_test.go'`):

- `spec.DeclaredTaskPaths` (`internal/spec/declared_paths.go`) unions a Task's
  `creates:` and `interface:` paths with `RecordedTaskPaths`. Its only caller
  is `archivedSpecDeclaredPaths` in `internal/cli/deliver_review_correction.go`,
  reached by `ProveReviewCorrection` from the Delivery Queue's retry of a
  `corrective-spec-required` item. Carried shape: the recorded paths now count;
  this is the reported defect.
- `dirtyPathsInArchivedSpec` (`internal/worktree/merged_head.go`) adds
  `RecordedTaskPaths` to an archived Spec's scoped paths for reconcile's
  supersession proof. Carried shape: the recorded paths now count, as the
  `reconcile` guide already states.
- The Spec Consistency Check's Verification path provision
  (`internal/speccheck/verification_path.go`) reads `RecordedTaskPaths` only
  for pending non-QA Tasks. A carried Task is `completed`, so the result
  changes only for a pending Task that still holds provenance, whose recorded
  paths now count as the rule documents.
- `RecordTaskPaths` writers: the Daemon's `commitTask`
  (`internal/daemon/task_engine.go`) and `settleTaskAndCommit`
  (`internal/cli/settle.go`). Carried shape: the existing section is replaced
  where it stands and the provenance kept, instead of a second section after
  it. Every other shape writes the same bytes as today.
- Not consumers: `nonAuthorialTaskSection` (`internal/speccheck/citations.go`),
  `reviewTaskSettlementRegion` (`internal/cli/review_conventions.go`) and the
  judge's planning pairs (`internal/judge/pairs.go`) classify each section by
  its heading, whatever the order; unchanged.
- Text: `docs/user-guide/commands/review.md`, `deliver.md`, `reconcile.md`,
  `implement.md`, `settle.md` and the Roundfix references say recorded paths
  count and that the Daemon writes the section at the end of the Task file;
  both stay true, so no guide changes.

## System Architecture

- `internal/spec/recorded_paths.go`: one unexported helper,
  `recordedPathsSections`, replaces `recordedPathsSectionOffset`;
  `RecordedTaskPaths` and `RecordTaskPaths` keep their signatures.
- Tests in `internal/spec`, `internal/cli` and `internal/worktree`.

No package, import edge, exported symbol or schema changes.

## Implementation Design

### Interfaces

```go
// recordedPathsSections returns, in file order, every unfenced
// "## Recorded paths" section that only Daemon-written settlement sections
// ("## Result", "## Carry-forward provenance", "## Recorded paths") follow.
// start includes the newline before the heading; end is the next section's start.
func recordedPathsSections(content []byte) []taskFileSection

type taskFileSection struct {
	heading    string
	start, end int
}
```

### Data Models

None changes. The Task file keeps its sections and their bytes.

### API Contracts

1. API Contract: the trailing settlement run. `recordedPathsSections` lists
   the unfenced level-2 headings with the fence rule
   `recordedPathsSectionOffset` uses today (the line, trimmed, opens or closes
   a fence by `markdownFenceMarker`; headings by `markdownHeading` at depth
   2; a heading matches by its exact untrimmed line). Each section starts at
   its heading, minus the newline before it when there is one, and ends where
   the next section starts or at the end of the file. Walking back from the
   last section, it collects each `## Recorded paths` section, passes over
   `## Result` and `## Carry-forward provenance`, and stops at any other
   heading.
2. API Contract: the reader. `RecordedTaskPaths` returns, in file order and
   without duplicates, the `- ` + backtick lines after the heading line of
   every collected section, and no path when none is collected. A file that
   ends with the recorded section reads exactly as today.
3. API Contract: the writer. `RecordTaskPaths` keeps its input checks and
   section bytes. With no collected section it appends the new section to the
   whole file, as today. Otherwise it writes the bytes before the first
   collected section, the new section (nothing when there are no paths), and
   every byte outside the collected sections, so a following `##
   Carry-forward provenance` or `## Result` keeps its bytes. It does not
   write the file when the result equals the content.
4. API Contract: unchanged consumers. `DeclaredTaskPaths`,
   `ProveReviewCorrection` and its refusal texts, `dirtyPathsInArchivedSpec`
   and its reason text, the Verification path provision, `commitTask`,
   `settleTaskAndCommit`, `RecordCarryForward` and the settlement heading
   classifiers keep their code and bytes.

### Surface Transcripts

None. No command, flag, exit code, message or stream placement changes; the
corrective retry, reconcile and the strict check accept more inputs of the
carried shape with today's texts.

### Invariants

```text
1. A recorded section followed by an authored level-2 section, or inside a fence, is never read or replaced.
2. A Task file that ends with its recorded section reads and rewrites exactly as before.
3. Recording paths never changes a byte outside the recorded sections and never leaves two recorded sections.
4. A changed path that no Task of the named archived Specs declares or records is refused with today's text.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → the reader through the corrective retry (API Contracts 1, 2 and 4)
- Goal 2, Success Metric 2, Invariant 4 → unchanged corrective-scope refusal (API Contract 4)
- Goal 3, Core Feature 3, Success Metric 5 → the reader through reconcile (API Contracts 2 and 4)
- Goal 4, Core Feature 2, Success Metric 4, Invariants 2-3 → the writer (API Contract 3)
- Goal 5, Success Metric 3, Invariant 1 → the trailing settlement run (API Contract 1)
- Success Metric 6 → the PRD's Unreachable Acceptance

## Absence Evidence

1. `recordedPathsSections` and `taskFileSection` — command: `grep -rn -e recordedPathsSections -e taskFileSection internal`; result: no match at `42f523c2`; the only section locator is `recordedPathsSectionOffset`, which reads the last level-2 heading only.

## ADR Alternatives

- ADR-0166: no rejected alternative is taken; the Daemon still never edits the authored Context, never records a Governed Path, and still appends a new section at the end of the Task file; only the reader and the in-place replacement learn that carry-forward provenance may follow it.
- ADR-0176: unrelated to its alternatives; no citation check changes, and `## Result`, `## Recorded paths` and `## Carry-forward provenance` stay the sections the citation projection treats as not written by the Spec's authors.

## Integration Points

None outside the repository: Git on temporary repositories and the local fake
reviewer the existing correction tests use. No test, Verification command or
QA row reads the real `~/.roundfix`, another repository or the network.

## Testing Approach

Every new top-level test calls `t.Parallel()` first and uses temporary
directories and repositories; fixtures are synthetic.

- `internal/spec/recorded_paths_section_test.go` (existing; its `Invariant:`
  header becomes "Only an exact unfenced Recorded paths heading followed by
  nothing but Daemon-written settlement sections is Daemon-owned."):
  - `TestRecordedPathsAreReadBeforeTrailingSettlementSections`, rows:
    `carry-forward-provenance` (recorded then provenance: the path is read
    and `DeclaredTaskPaths` holds it); `result-then-provenance`;
    `duplicate-sections` (recorded, provenance, recorded: the union in file
    order); `authored-section-between` (recorded, `## Notes`, provenance: no
    path); `fenced-heading-before-provenance` (a fenced recorded heading in
    `## Result`, then provenance: no path).
  - `TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance`, rows:
    `same-paths` (file bytes and inode unchanged), `new-paths` (exact bytes,
    provenance kept), `no-paths` (section removed, provenance kept),
    `duplicate-sections` (one section at the first one's place, provenance
    kept).
- `internal/cli/review_declared_correction_test.go`:
  `TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared`, with
  `archive-record` and `legacy-folder` subtests over a Task built by
  `declaredCorrectionTask("", true)` plus a provenance section: a correction
  of `review.txt` is accepted; a second correction that also changes
  `other.txt` is refused with the exact text of Success Metric 2.
- `internal/worktree/merged_spec_leftovers_test.go`:
  `TestArchivedSpecLeftoverRecordedBeforeCarryForwardProvenance`: the archived
  Task of `newArchivedLeftoverFixture` rewritten to end with `## Recorded
  paths` then `## Carry-forward provenance`, `recorded.txt` dirty, state
  `superseded`.
- Unchanged and run by name: `TestRecordedPathsKeepAFencedHeadingInAuthoredText`,
  `TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections`,
  `TestRecordedPathsReplaceTheTrailingDaemonSection`,
  `TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte`,
  `TestRecordTaskPathsWithNoPathsRemovesTheSection`,
  `TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths`,
  `TestARecordedPathCountsAsDeclared`,
  `TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath`,
  `TestDeclaredPathsAreReadAtTheRecordedRevision`,
  `TestArchivedSpecLeftoverScopeKinds`,
  `TestTaskCommitRecordsAnUndeclaredTestFile`,
  `TestTaskCommitRecordsNothingWhenEveryPathIsDeclared`,
  `TestSettleRecordsTheUndeclaredPathsOfItsTaskCommit` and
  `TestVerificationPathUndeclaredIgnoresWhatThePlanProvides`.
- A prototype of API Contracts 1-3 passed these tests at authoring time, and
  five sabotages of it each failed a named row: dropping the pass over
  settlement headings (`carry-forward-provenance`), never stopping at an
  authored heading (`authored-section-between` and
  `TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections`), ignoring fences
  (`fenced-heading-before-provenance` and
  `TestRecordedPathsKeepAFencedHeadingInAuthoredText`), dropping the bytes
  after the replaced section (`new-paths`), and always appending (`same-paths`
  and `TestRecordedPathsReplaceTheTrailingDaemonSection`).

## Build Order

1. The section helper, the reader, the writer and their tests in
   `internal/spec`, `internal/cli` and `internal/worktree` (no dependency).
2. QA gate (depends on: 1).

No Spec is listed in `requires`.

## Risks & Considerations

- An authored `## Recorded paths` directly before `## Result` is read as the
  Daemon's; the heading is Daemon-owned (ADR-0166), and today an authored one
  at the end of the file is read the same way.
- A level-2 heading inside an agent-written `## Result` ends the trailing run
  and hides an earlier recorded section, as today; the oraculum shape uses
  level-3 headings there.
- Rewriting duplicates into one section drops nothing the reader would have
  returned before the fix, which read only the last section.

## Glossary

None.

## Decisions

- The section is the Daemon's when only Daemon-written settlement sections
  follow it, the ones convention C2 names; no new ADR, ADR-0166 stands.
- The reader unions duplicate sections and the writer collapses them at the
  first one's place.
- No consumer, guide or skill changes; the fix lives in the one reader.
