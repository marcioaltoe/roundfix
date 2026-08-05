---
task: task_06
spec: 0067-derived-artifact-regeneration-boundary
status: pending
type: backend
complexity: medium
---

# Task 06: Make every record state what is measurably true

## Overview

The gate found the Spec failing its own first Success Metric, and its records
asserting things that are not so.

F-001: one sanctioned invocation still leaves `make verify` red, because
`BASELINE_DIGEST_STEPS` gained the catalog-diagnostic corpus and not the
plan-characterization one. task_03 named a single corpus where Core Feature 2
requires the command to cover **everything it claims**, and there are two.

F-002: the catalog-diagnostic record still says it is not in
`BASELINE_DIGEST_STEPS` after task_03 put it there, and the parity record says
nothing regenerates the corpus while the sanctioned run demonstrably rewrites
`fixtures/asset-sync.json` and `manifest.json` beneath it.

The second half matters beyond this Spec: the PRD asserts the parity corpus is
frozen, and the evidence says it is not. This Task records what is measurably
true and does not decide whether that state is desirable — see Requirement 6.

## Requirements

1. MUST add the plan-characterization corpus to `BASELINE_DIGEST_STEPS`, using
   the exact invocation its ownership record declares, so one sanctioned run
   leaves `make verify` green after an owned-Skill edit.
2. MUST update the catalog-diagnostic record to state that it **is** covered by
   the sanctioned command, since task_03 made that true.
3. MUST derive each record's `owner` from measurement rather than from
   assumption: run the sanctioned command in a fixture and record which
   directories it rewrites.
4. MUST correct the parity record to state what the measurement shows. If the
   sanctioned command rewrites artifacts beneath it, the record MUST NOT claim
   nothing regenerates them.
5. MUST make the declared-step test actually fail when a `frozen` directory is
   rewritten. It passed while the parity corpus was being rewritten, so it is
   not currently asserting what task_02 claimed.
6. MUST NOT change which artifacts the sanctioned command rewrites beyond
   Requirement 1, and MUST NOT make the parity corpus stop being rewritten.
   Whether it *should* be frozen is a product decision the PRD asserts and the
   evidence contradicts; this Task makes the record honest and leaves the
   decision recorded for the maintainer.
7. MUST NOT change any digest value or artifact content as a result of this
   Task.

## Subtasks

- [ ] Add the plan-characterization step to the sanctioned list.
- [ ] Measure which directories the sanctioned run rewrites.
- [ ] Correct every record whose `owner` or `reason` the measurement
      contradicts.
- [ ] Make the frozen assertion able to fail, and prove it fails.

## Acceptance Criteria

- [ ] Editing an owned Skill, running `make skills-sync` and then
      `make baseline-digests` once, leaves `make verify` green.
- [ ] A second sanctioned run reports no changes and `make verify` stays green.
- [ ] Every record's `owner` matches the measured behaviour of the sanctioned
      command, asserted by a test that compares record to measurement.
- [ ] The declared-step test fails when a `frozen` directory is rewritten,
      proven by a fixture that rewrites one and observes the failure.
- [ ] No digest value or artifact content changed by this Task.

## Context

- interface: `Makefile`
- instruction: `docs/agents/agent-instructions.md`

## Verification

- `make skills-sync && make baseline-digests && make verify` — expected: exit 0
  for all three; one sanctioned run leaves the gate green.
- `make baseline-digests && make verify` — expected: exit 0; the second run is
  a no-op and the gate stays green.
- `go test ./internal/baseline -count=1 -run 'Ownership|DeclaredStep|Frozen|Measured' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the record-versus-measurement tests ran and passed.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD | grep -vE "^(Makefile$|internal/baseline/testdata/.*_ownership\.yml$|internal/baseline/.*_test\.go$|internal/baseline/.*\.go$|docs/specs/0067-derived-artifact-regeneration-boundary/task_06\.md$)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded tooling file, the records, the tests they
  drive, and this Task file changed.

## References

- `_prd.md` → Core Features 1, 2, 3 and 5; Success Metric 1; Decisions.
- `_techspec.md` → Build Order 1 and 3; Risks & Considerations ("a manifest can
  lie").
- `qa/qa-report-2026-08-05.md` → F-001 and F-002.
- `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.
- ADR-0081, ADR-0085.
