---
task: task_05
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: backend
complexity: high
---

# Task 05: Carry a row forward only on unmoved evidence

## Overview

The one mechanism in this Spec that can make something green wrongly, so it is
designed to fail closed in every direction and tested that way: the refusal
cases outnumber the happy path six to one, because the risk is entirely on the
permissive side.

A row is carried because nothing it depends on moved — never because it passed
recently. Recency would inherit a verdict; this inherits an observation, and
says so in the report.

## Requirements

1. MUST implement `Carriable` with every condition of ADR-0097 holding at
   once: the earlier report established the row as `pass`; the row declared a
   non-empty typed input list and every entry is `repository_path`; the
   earlier report's head is an ancestor of the current head; no declared input
   appears in the changed-path set between those heads; and every cited
   evidence path still resolves with unchanged content.
2. MUST refuse to carry a row that failed, was blocked by any cause, was
   skipped, or declared no inputs.
3. MUST refuse a mixed input list: any `external_repository`, `live_service`,
   or `elapsed_time` entry makes the row re-observe, even alongside
   repository paths.
4. MUST record, on every carried row, the report and the head that established
   it, so a reader sees which evidence is fresh and which is inherited.
5. MUST NOT make any verdict more permissive than a fresh observation would;
   ADR-0080 keeps verdict semantics and a carried row cannot change a count's
   meaning.
6. MUST reuse the existing changed-path primitive rather than adding a second
   way to compute what moved.

## Subtasks

- [ ] Implement the resolver and its ancestry and changed-path checks.
- [ ] Author the refusal suite before the happy path.
- [ ] Emit the establishing citation into the carried row.

## Acceptance Criteria

- [ ] Each refusal condition has its own named test and each refuses.
- [ ] A row meeting every condition carries, citing its establishing report
      and head.
- [ ] A mixed input list never carries.
- [ ] No verdict computation changed.

## Rehearsal Cases

- Case: prior status is not `pass`; Observation: `Carriable` returns false and
  the row is re-observed.
- Case: row declares no inputs; Observation: refused, so declaration is
  opt-in.
- Case: a declared `repository_path` appears in the changed-path set;
  Observation: refused.
- Case: the establishing head is not an ancestor of the current head;
  Observation: refused.
- Case: a cited evidence file resolves but its content changed; Observation:
  refused.
- Case: inputs mix `repository_path` with `elapsed_time`; Observation:
  refused.
- Case: every condition holds; Observation: carried, with the establishing
  report path and head recorded on the row.

## Context

- interface: internal/speccheck/citations.go
- interface: internal/worktree/worktree.go
- interface: internal/daemon/task_engine.go

## Verification

- `output="$(go test -count=1 ./internal/... -run 'Carriable' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the carry-forward tests are selected and pass.
- `output="$(go test -count=1 ./internal/... -run 'Carriable' -v 2>&1)"; printf '%s\n' "$output" | grep -cE -- '--- PASS: [^ ]+/' | { read count; [ "$count" -ge 7 ]; }`
  — expected: exit 0; at least the seven rehearsal cases are present as
  passing subtests, so the refusal suite cannot be reduced to one happy path.
  — expected: exit 0; the existing changed-path primitive is the one in use and
  both packages stay green.

These commands are deliberately absent: `go build -buildvcs=false ./...` and a
whole-package `go test` sweep both pass against a tree where no work has
happened, so each approves the Task before it starts. Compilation and
regression are the Run-level gate's job; the commands above name cases that
do not exist yet.

## References

- `_prd.md` → Core Feature 5; User Stories 2, 6; Success Metrics.
- `_techspec.md` → Implementation Design (Carry-forward); Testing Approach;
  Build Order 5.
- ADR-0097, ADR-0080.

## Result

Implemented evidence-scoped carry-forward inside the existing pre-QA
mechanical stage. `Carriable` now refuses every non-pass status, absent or
mixed inputs, unproven ancestry, intersecting changed paths, incomplete or
malformed snapshots, changed path sets or digests, and cited evidence outside
the declared snapshots. The repository resolver parses typed row declarations
and canonical snapshot frontmatter, proves the establishing head's ancestry,
reuses `worktree.PriorChangedFiles` for the exact establishing-head-to-`HEAD`
range, expands literal and `*`/`?`/`**` repository inputs from tracked Git
blobs, and hashes their bytes at the current head. A carried row retains the
original report and head even when the immediately previous report already
carried it.

Focused checks run after the final implementation edit:

- Pre-change signal: `rtk rg -n '^func Carriable|TestCarriable' internal/speccheck`
  exited 1 because neither the API nor its named suite existed. After the
  refusal suite was authored first,
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/speccheck -run '^TestCarriable$' -count=1`
  failed to compile on the deliberately missing `ReportRow`, snapshot types,
  and `Carriable` implementation.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/speccheck -run '^(TestCarriable|TestMechanical.*|TestMaterializeMechanicalResult|TestCheckCorpusBudget)$' -count=1 -v`
  passed. It reported fourteen named refusal subtests, two happy-path subtests,
  four real-Git resolver cases, the existing detector cases, materialization,
  and corpus non-regression.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/daemon -run '^(TestMechanicalStage|TestQAMechanicalRequest|TestWriteMechanicalQAReport|TestTaskCycleQA|TestPerWorkAgentSessionMixedTaskTypesAndQA|TestQAPullRequest|TestPriorChangedFilesUseCurrentWorktreeHeadAndIgnoreSiblingBranch)' -count=1`
  passed, including the shared changed-path primitive's current-worktree-head
  behavior and the existing Daemon QA-stage boundaries.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/spec -run '^TestQAVerdict' -count=1`
  passed, preserving the existing pass, partial, fail, missing-report,
  unreadable-report, blocked-count, and report-recency semantics.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go vet ./internal/speccheck ./internal/worktree ./internal/daemon`
  passed with no diagnostics.
- `rtk git -c core.fsmonitor=false diff --check` passed.

Acceptance evidence:

1. `TestCarriable` names and refuses failed, blocked, skipped, inputless,
   changed-path, non-ancestor, changed-content, mixed-input, missing-snapshot,
   missing-citation, uncovered-citation, glob-intersection, and changed-path-set
   cases independently.
2. `TestMechanicalStageCarriableCarriesUnchangedEvidenceWithCitation` proves
   an unchanged row carries across an unrelated commit and materializes the
   establishing report path and head. The recursive-glob happy path also
   carries, and `TestMechanicalStageCarriablePreservesOriginalEstablishingCitation`
   proves a later carry does not replace the original citation.
3. `TestCarriable/mixed_repository_and_elapsed_inputs_refuse_carry` proves a
   repository path combined with `elapsed_time` never carries; the production
   predicate rejects every non-`repository_path` kind.
4. No verdict package or verdict computation changed. The focused
   `TestQAVerdict*` selection passed, while carry-forward only populates the
   verdict-free `MechanicalResult.Carried` rows.

The commands under this Task's `## Verification` were not run; the Daemon owns
that complete selection and settlement evidence.
