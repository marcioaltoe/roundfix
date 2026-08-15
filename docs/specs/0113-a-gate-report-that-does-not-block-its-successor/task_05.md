---
task: task_05
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 05: Let a Spec's own coined term reach its gate

## Overview

The authoring contract gives the QA gate the glossary update for a term the Spec
coined. The gate's static precondition refuses on that term being undocumented,
before the gate can document it. The repair is assigned to the one actor forbidden
from reaching it, and two Specs stalled there on consecutive days — 0103 on
2026-08-14 and 0114 on 2026-08-15.

## Requirements

1. MUST NOT refuse the gate's static precondition for a term the Spec under gate
   declared in its own Vocabulary Contract and has not yet documented.
2. MUST still refuse a term emitted by code that no Spec declared.
3. MUST still refuse the same Spec's undocumented term outside the gate, so
   authoring stages keep reporting it.
4. MUST report the pending term rather than hiding it, so the gate knows what it
   is expected to document.

## Subtasks

- [ ] Exempt a Spec's own declared term from its gate's precondition.
- [ ] Keep every other undocumented-term refusal.
- [ ] Cover the exempt, the undeclared, and the other-stage cases.

## Acceptance Criteria

- [ ] A Spec whose Vocabulary Contract declares an undocumented term passes its
      gate's static precondition.
- [ ] A term emitted by code and declared by no Spec is still refused.
- [ ] The same undocumented term is still reported outside the gate.
- [ ] The pending term is named in the gate's own input rather than silently
      dropped.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestGateAcceptsItsOwnDeclaredTerm' -v > /tmp/0113-t05.log 2>&1; s=$?; grep -q '^--- PASS: TestGateAcceptsItsOwnDeclaredTerm' /tmp/0113-t05.log || { cat /tmp/0113-t05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0113-t05.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestGateAcceptsItsOwnDeclaredTerm' /tmp/0113-t05.log || { echo 'the exemption test did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t05.log > /tmp/0113-t05-n.txt; test "$(cat /tmp/0113-t05-n.txt)" -ge 3 || { echo "expected the exempt, undeclared, and other-stage cases, got $(cat /tmp/0113-t05-n.txt)"; cat /tmp/0113-t05.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving the exemption is narrow rather than blanket.
- `go build -buildvcs=false -o /tmp/0113-t05-roundfix ./cmd/roundfix && /tmp/0113-t05-roundfix spec check > /tmp/0113-t05-sweep.log 2>&1; grep -q 'SC-VOCABULARY-UNDOCUMENTED' /tmp/0113-t05-sweep.log && { echo 'the corpus gained an undocumented-term finding:'; grep 'SC-VOCABULARY-UNDOCUMENTED' /tmp/0113-t05-sweep.log; exit 1; }; grep -rq 'TestGateAcceptsItsOwnDeclaredTerm' internal/speccheck || { echo 'the sweep is clean, but the exemption does not exist'; exit 1; }` — expected: exits 0, proving the exemption did not loosen the corpus, anchored to the case this Task adds.

## Context

- interface: `internal/speccheck/mechanical.go`
- instruction: `docs/backlog/2026-08-14-a-spec-that-coins-a-term-cannot-pass-its-own-gate.md`

## References

`_techspec.md` → Build Order 5; Risks & Considerations, relaxing the vocabulary
precondition. `_prd.md` → Core Feature 6; Goal 4; User Story 4. ADR-0134.

## Result

Implemented the gate-precondition classifier at the mechanical-stage seam. The
vocabulary detector now retains private provenance only for an undocumented
token selected by the checked Spec's complete Vocabulary Contract.
`GatePrecondition` moves only that finding into its named `Inputs`; every other
finding remains blocking, and the ordinary authoring result is not mutated.

Attempt 1 feedback showed that the corpus sweep had matched a presence-aware
vocabulary skip, not an undocumented-token finding. Text output now names that
skip as the `vocabulary documentation detector`; JSON keeps the stable
`SC-VOCABULARY-UNDOCUMENTED` skip code, and actual text findings still print the
code. This keeps missing-input observability without presenting a skip as a
finding.

Acceptance evidence:

- A current Spec's declared but undocumented `publish:` token produced a
  non-blocking precondition with one input naming `publish:` in
  `TestGateAcceptsItsOwnDeclaredTerm/own_declared_term_becomes_named_gate_input`.
- A synthetic `orphan:` vocabulary diagnostic with no Spec-declaration
  provenance remained a blocking finding and did not become gate input in
  `TestGateAcceptsItsOwnDeclaredTerm/undeclared_emitted_term_still_blocks`.
- `CheckStage(..., StageTechSpec)` continued to return the original
  `SC-VOCABULARY-UNDOCUMENTED` finding for `publish:` in
  `TestGateAcceptsItsOwnDeclaredTerm/authoring_stage_still_reports_declared_term`.
- The positive case asserted the pending token survives in
  `GatePreconditionResult.Inputs`, rather than disappearing when removed from
  the blocking findings.
- `TestRenderVocabularySkipDoesNotLookLikeFinding` proved the vocabulary skip
  stays visible by name and missing artifact, remains coded in JSON, and does
  not collide with the text signature retained by an actual vocabulary
  finding.

Focused checks:

- Pre-change reproduction:
  `rtk env GOCACHE=/private/tmp/roundfix-0113-task05-gocache go test ./internal/speccheck -run '^TestGateAcceptsItsOwnDeclaredTerm$/own_declared_term_becomes_named_gate_input$'`
  failed to compile because `speccheck.GatePrecondition` did not exist.
- After implementation:
  `rtk env GOCACHE=/private/tmp/roundfix-0113-task05-gocache go test ./internal/speccheck -run '^TestGateAcceptsItsOwnDeclaredTerm$'`
  exited 0.
- Attempt 1 regression reproduction:
  `rtk env GOCACHE=/private/tmp/roundfix-0113-task05-gocache go test ./internal/speccheck -run '^TestRenderVocabularySkipDoesNotLookLikeFinding$'`
  failed because the text skip contained the finding code; after the renderer
  repair, the same focused check exited 0.
- After formatting the final implementation and test edits:
  `rtk env GOCACHE=/private/tmp/roundfix-0113-task05-gocache go test -count=1 ./internal/speccheck`
  exited 0 (`ok roundfix/internal/speccheck 2.341s`).
- The focused CLI renderer checks
  `rtk env GOCACHE=/private/tmp/roundfix-0113-task05-gocache go test -count=1 ./internal/cli -run '^(TestRunSpecCheckCleanText|TestSpecCheckWithoutStageIsUnchanged)$'`
  exited 0 (`ok roundfix/internal/cli 0.648s`).
- `rtk git diff --check` exited 0.

The commands under `## Verification` were not run; the Daemon owns them. Task
06 owns consuming this classified input while performing and recording assigned
repairs.
