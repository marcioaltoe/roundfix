---
task: task_07
spec: 0073-skill-versions-decoupled-from-the-binary
status: completed
type: backend
complexity: medium
---

# Task 07: Assert the property directly, and record the external provenance

## Overview

Corrective Task from the QA gate's F-001 (`Blocks-Completion`) and F-002.

**F-001 is a harness defect, not a product defect.**
`TestOwnedSkillEditLeavesMakeVerifyGreen` proves the Spec's headline property
by running the complete `make verify` **nested inside a test**. That makes the
Spec's own acceptance hostage to every other test in the repository: the gate
observed it fail on
`TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`,
which passes in isolation and failed under concurrent load. The property being
proven is narrow — an owned-skill edit needs no regeneration step — and it
should be asserted narrowly.

**F-002 is a Success Metric with no Task.** The PRD requires that
`golang-dependency-management`, `golang-safety`, and `golang-structs-interfaces`
carry recorded provenance without becoming setup-snapshot entries. No Task
requirement ever asked for it, so nothing delivered it — the same shape as Spec
0059's F-001 in the same session.

## Requirements

1. MUST replace the nested-`make verify` proof with a direct assertion of the
   property: after editing an owned skill's content, the derived Baseline
   artifacts and both characterization corpora are unchanged, so no
   regeneration step is required.
2. MUST assert that property from the artifacts themselves — compare their
   bytes before and after the edit — rather than by invoking a repository-wide
   gate from inside a test.
3. MUST NOT weaken what is proven. The claim is unchanged; only the instrument
   changes, from "the whole suite passes" to "these artifacts did not move".
4. MUST NOT record provenance for the three external Go Skills. That work
   needs `skills-lock.json`, which this Spec's tooling authority does not
   reach, and the Task assigned it on 2026-08-06 correctly refused to widen its
   own boundary. The PRD now declares it under `## Unreachable Acceptance` with
   the grant that would satisfy it.
5. MUST leave every owned-skill readiness behaviour from task_02 and task_04
   unchanged, asserted over the existing tests.

## Subtasks

- [ ] Replace the nested gate with a direct artifact-stability assertion.
- [ ] Record provenance for the three external Go Skills.

## Acceptance Criteria

- [ ] Editing an owned skill leaves every derived Baseline artifact and both
      characterization corpora byte-identical, asserted directly.
- [ ] No test in this Spec invokes `make verify` from inside a test.
- [ ] None of the three external Go Skills is compared against an owned
      minimum, which remains true without recording their provenance.
- [ ] Owned-skill readiness behaviour is unchanged, asserted over the existing
      tests.

## Context

- interface: `skills/baseline_skill_contract_integration_test.go`
- interface: `skills/repository.go`
- interface: `skills-lock.json`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `grep -rn "make verify" skills/*_test.go | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no test in this package invokes the repository gate.
- `output="$(go test ./skills -count=1 -run 'OwnedSkillEdit|Provenance' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the property and provenance tests ran and passed.
- `git diff --name-only HEAD | grep -q "^skills-lock.json$" && exit 1 || exit 0`
  — expected: exit 0; this Task does not touch the protected lock file it has
  no authority over.
- `go test ./skills ./internal/baseline -count=1` — expected: exit 0.
- `make verify` — expected: exit 0.

## References

- `qa/qa-report-2026-08-06.md` → F-001, F-002.
- `_prd.md` → Success Metrics 1 and 7; Which skills the contract covers.
- `_techspec.md` → Testing Approach.

## Result

### Implementation

- The current Task base already contains commit `e4f07b1b`, which replaced
  `TestOwnedSkillEditLeavesMakeVerifyGreen` with
  `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`. The focused history
  inspection confirmed that its parent still imported `os/exec` and nested the
  repository-wide gate, while the current test invokes no command.
- The direct test copies the tracked repository, snapshots every file under the
  `Makefile`'s `DERIVED_DIGEST_PATHS`, snapshots both characterization corpora
  explicitly, edits the canonical and mirrored `roundfix` Skill, and compares
  every artifact's bytes before and after. Its comparison rejects changed,
  removed, and newly created artifacts.
- The archived-Spec byte assertion from the prior proof remains in the direct
  test, so changing the instrument did not drop that preservation check.
- This Run made no implementation-code change because the corrective test was
  already present at `HEAD`. It did not change owned-skill readiness code or
  `skills-lock.json`.

### Focused checks

- Red signal: `rtk git show
  e4f07b1b^:skills/baseline_skill_contract_integration_test.go` showed the old
  build-tagged test importing `os/exec` and invoking `make verify` through
  `exec.Command`.
- `rtk go test ./skills -count=1 -run
  '^(TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical|TestRepositoryReadinessNeverComparesThirdPartySkillVersions|TestOwnedSkillBundleReadinessKeepsStatesDistinct|TestCheckRepositoryClassifiesMissingAndOutdatedSkills)$'
  -v`: authorized rerun exited 0; RTK reported 15 passing tests. The first
  sandboxed attempt reached no tests because the host Go build cache was not
  writable from the sandbox.
- `rtk grep -n "make verify|exec\\.Command|os/exec"
  skills/baseline_skill_contract_integration_test.go`: exit 1 with no matches.
- `rtk grep -n "make verify" skills/*_test.go`: exit 1 with no matches.
- `rtk grep -n
  "golang-dependency-management|golang-safety|golang-structs-interfaces"
  skills-lock.json`: exit 1 with no matches, confirming the protected
  provenance file remains unchanged for those Skills.
- `rtk git diff HEAD -- skills/baseline_skill_contract_integration_test.go
  skills/repository.go skills-lock.json`: exit 0 with an empty diff.

### Acceptance evidence

- The focused direct test passed after reading the derived roots from the
  repository-owned `DERIVED_DIGEST_PATHS` declaration and separately reading
  `catalog.diagnostics.golden.json` and `plan-characterization`; it compares
  the captured bytes after the owned-Skill edit without a regeneration step.
- Both source searches found no nested `make verify`, `exec.Command`, or
  `os/exec` use in the target test, and no `make verify` text in any
  `skills/*_test.go` file.
- `TestRepositoryReadinessNeverComparesThirdPartySkillVersions` passed with a
  call ledger containing exactly the owned Skills and no third-party Skill.
  The three named external Go Skills still have no protected-lock entry, as
  required by this Task's authority boundary.
- `TestOwnedSkillBundleReadinessKeepsStatesDistinct` and
  `TestCheckRepositoryClassifiesMissingAndOutdatedSkills` passed in the same
  focused run, preserving the existing owned-skill readiness states and
  repository classifications from task_02 and task_04.

### Follow-up

- The Subtask line that says to record provenance contradicts Requirement 4,
  the third acceptance criterion, and the PRD's `## Unreachable Acceptance`.
  It remains unchecked. Recording that provenance requires a maintainer grant
  naming `skills-lock.json` and a separately bounded Task.
- This Task's declared `## Verification` commands were not run; terminal
  Verification and status settlement remain Daemon-owned.
