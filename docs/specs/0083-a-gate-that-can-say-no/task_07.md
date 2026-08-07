---
task: task_07
spec: 0083-a-gate-that-can-say-no
status: completed
type: qa
complexity: high
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate for this Spec. Its distinguishing obligation is
adversarial: this Spec's whole subject is a gate that reported success while the
suite failed, so the QA evidence must include making the gate fail on purpose
and observing it say no. A green reading is exactly the evidence that was
untrustworthy before this Spec.

## Requirements

1. MUST derive the QA matrix from the PRD's goals and core features and the
   completed task evidence, and execute it through the real gate.
2. MUST prove the gate says no: introduce a failing test, run the repository's
   verification gate by its documented name, and observe a non-zero exit.
3. MUST prove the gate says yes on a clean tree, so the repair did not simply
   make the gate always fail.
4. MUST prove the retired gates no longer fire on unrelated work: add an ADR and
   a Spec folder in a throwaway copy and observe no gate failure from the
   archived corpus counter.
5. MUST prove the retained gates still fire: `spec check`, the published-example
   contract, and the coverage invariant must each still fail on an induced
   defect.
6. MUST prove both stabilized tests under induced load, independently of the
   evidence their own Tasks recorded.
7. MUST classify every finding by user impact and record auditable evidence.
8. MUST operate only on copies for any destructive rehearsal, create them under
   the system temporary directory with the shared prefix `roundfix-qa-0083-`,
   delete every copy before finishing including after a failed case, and record
   each path and its removal in the report.
9. MUST read every gate result from the gate's documented command, not from a
   hand-built approximation — the Spec exists because a wrapper's summary lied.

## Rehearsal Cases

- Case: introduce a failing test and run the verification gate; Observation: non-zero exit and the failing package named.
- Case: introduce a failing test that emits high-volume log output alongside its error; Observation: non-zero exit, proving the original masking condition is covered.
- Case: run the verification gate on a clean tree; Observation: exit 0.
- Case: add an ADR and a Spec folder in a copy, then run the gate; Observation: no failure attributable to the archived corpus counter.
- Case: run the corpus budget check under induced CPU load; Observation: passes.
- Case: induce an inefficiency in the corpus sweep; Observation: the performance signal still fails.
- Case: induce a Spec consistency defect; Observation: `spec check` still fails.
- Case: document a CLI example the parser rejects; Observation: the published-example contract still fails.
- Case: remove a recorded test; Observation: the coverage invariant still fails, and its record resolves from outside `docs/specs/`.
- Case: run each stabilized test twenty times under induced load; Observation: no failure.

## Subtasks

- [ ] Derive the resumable QA matrix from the PRD and task evidence.
- [ ] Execute every rehearsal case and capture its observation.
- [ ] Prove each retired gate is quiet and each retained gate still bites.
- [ ] Classify findings by user impact and write the dated QA report.
- [ ] Delete every copy created, and record each path and its removal.

## Acceptance Criteria

- [ ] Every rehearsal case above is executed and its observation recorded.
- [ ] Every PRD goal and core feature has a matrix row with a verdict and evidence.
- [ ] The gate is shown failing on an induced defect and passing on a clean tree.
- [ ] Every retained gate is shown still failing on an induced defect.
- [ ] Any environment-blocked row carries equivalent evidence and a stated reason.
- [ ] The dated QA report is written under the Spec's `qa/` directory.
- [ ] No repository copy remains on disk after the gate finishes, including for
      cases that failed.
- [ ] The QA report lists every copy path used and states that each was removed.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `.agents/skills/evidence-gate/SKILL.md`
- instruction: `docs/findings/2026-08-07-the-only-gate-reports-green-on-a-red-suite.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `ls docs/specs/0083-a-gate-that-can-say-no/qa/ | grep -q '\.md$'` — expected: exits 0, proving the dated QA report exists.
- `grep -q -i 'non-zero' docs/specs/0083-a-gate-that-can-say-no/qa/*.md` — expected: exits 0, proving the gate was observed failing rather than assumed able to.
- `ls -d "${TMPDIR:-/tmp}"/roundfix-qa-0083-* 2>/dev/null | grep . ; test $? -eq 1` — expected: exits 0, proving every copy the gate created was removed.
- `grep -q 'roundfix-qa-0083' docs/specs/0083-a-gate-that-can-say-no/qa/*.md` — expected: exits 0, proving the report records the copy paths it used.
- `go test ./... -count=1` — expected: exits 0.

## References

- `_prd.md` → all Goals and Core Features.
- `_techspec.md` → Testing Approach; Risks & Considerations.
- ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0097.
