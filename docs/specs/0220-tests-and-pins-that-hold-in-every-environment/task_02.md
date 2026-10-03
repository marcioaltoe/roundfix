---
task: task_02
spec: 0220-tests-and-pins-that-hold-in-every-environment
status: pending
type: chore
complexity: medium
---

# Task 02: Restore the trailing skills and hold them to the lock and the snapshot

## Overview

Eleven required external skills trail the Setup Snapshot, and the supported
restore that would bring them to it fails three skill contract tests that pin
one hand-written digest of every upstream-managed tree. This Task runs the
restore and, in the same change, makes those tests read the records the restore
writes: each `skills-lock.json` entry against its installed tree, and every
required external skill against the Setup Snapshot (ADR-0224). Afterwards Doctor
reports `skills: ok` and `make verify` passes.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-02, "Restoring trailing skills
   breaks the pinned skill digests" (adopted at
   [references/2026-10-02-restoring-trailing-skills-breaks-the-pinned-skill-digests.md](references/2026-10-02-restoring-trailing-skills-breaks-the-pinned-skill-digests.md)).
2. MUST restore the eleven trailing skills with the supported restore,
   `roundfix baseline update --repo . --yes`, using a binary built from this
   tree. When the Agent turn has no network, the restore MUST take a local
   checkout of the upstream skills repository at `b3c45a4` through
   `--source-dir`. A second run MUST report `Skills restored: 0`. The restored
   trees and `skills-lock.json` MUST hold only the bytes the restore wrote; no
   upstream skill text is edited by hand.
3. MUST change only the Governed Paths listed in `_authorization.md` `paths`.
   If the restore writes any other path under `.agents/skills/`, the Task MUST
   stop and report that path instead of committing it.
4. MUST remove `upstreamManagedSkillTreeDigest` and every comparison with it
   from `skills/baseline_skill_contract_test.go`. `TestUpstreamADRFormatUnchanged`,
   `TestAuthoringConstraintOwnership` and `TestAuthorialSkillSync` MUST instead
   compare each `skills-lock.json` entry's `computedHash` with
   `SkillFolderHash` of its installed tree, and fail naming every skill that
   differs. The `ADR-FORMAT.md` digest check MUST stay unchanged.
5. MUST add `TestUpstreamManagedSkillLockNamesAnEditedTree` in the skills
   package. It copies one upstream-managed tree and its lock entry into a
   temporary repository, changes one byte of the copy, and asserts that the
   lock comparison names that skill and no other.
6. MUST add `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` in
   `internal/baseline`. It reads the Baseline Profile from
   `docs/agents/setup-context.json` and the skill names from `skills-lock.json`,
   calls `TrailingSetupSkills`, and fails when any name is returned, naming
   each skill and the command `roundfix baseline update`.
7. MUST NOT change production code, the Setup Snapshot or any derived Baseline
   file.

## Subtasks

- [ ] Run the restore and confirm a second run restores nothing.
- [ ] Replace the pinned digest with the lock comparison in the three tests.
- [ ] Add the lock-mismatch test and the snapshot test.
- [ ] Confirm Doctor's `skills` line and the repository Verification.

## Acceptance Criteria

- [ ] No test names the digest of an upstream-managed skill tree.
- [ ] A one-byte edit of an upstream-managed tree fails the lock comparison,
      naming that skill.
- [ ] `TrailingSetupSkills` over this repository returns no name, and Doctor's
      `skills` line is `ok` with no `DR-SKILL-TRAILS-SNAPSHOT`.
- [ ] `go test ./skills` and `go test ./internal/baseline` pass on the restored
      tree.

## Context

- interface: `.agents/skills/bubbletea/references/components.md`
- interface: `.agents/skills/bubbletea/references/emoji-width-fix.md`
- interface: `.agents/skills/bubbletea/references/golden-rules.md`
- interface: `.agents/skills/bubbletea/references/troubleshooting.md`
- interface: `.agents/skills/bubbletea/SKILL.md`
- deletes: `.agents/skills/domain-modeling/CONTEXT-FORMAT.md`
- creates: `.agents/skills/domain-modeling/GLOSSARY-FORMAT.md`
- interface: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/golang-concurrency/references/channels-and-select.md`
- interface: `.agents/skills/golang-concurrency/references/pipelines.md`
- interface: `.agents/skills/golang-concurrency/references/sync-primitives.md`
- interface: `.agents/skills/golang-concurrency/SKILL.md`
- interface: `.agents/skills/golang-context/references/cancellation.md`
- interface: `.agents/skills/golang-context/references/http-services.md`
- interface: `.agents/skills/golang-context/SKILL.md`
- interface: `.agents/skills/golang-error-handling/references/error-creation.md`
- interface: `.agents/skills/golang-error-handling/references/error-handling.md`
- interface: `.agents/skills/golang-error-handling/references/error-wrapping.md`
- interface: `.agents/skills/golang-error-handling/SKILL.md`
- interface: `.agents/skills/golang-lint/references/linter-reference.md`
- interface: `.agents/skills/golang-lint/SKILL.md`
- creates: `.agents/skills/golang-testing/references/benchmarks.md`
- creates: `.agents/skills/golang-testing/references/coverage.md`
- creates: `.agents/skills/golang-testing/references/examples.md`
- interface: `.agents/skills/golang-testing/references/integration-testing.md`
- interface: `.agents/skills/golang-testing/references/mocking.md`
- interface: `.agents/skills/golang-testing/SKILL.md`
- interface: `.agents/skills/no-workarounds/SKILL.md`
- interface: `.agents/skills/systematic-debugging/condition-based-waiting-example.ts`
- interface: `.agents/skills/systematic-debugging/condition-based-waiting.md`
- interface: `.agents/skills/systematic-debugging/CREATION-LOG.md`
- interface: `.agents/skills/systematic-debugging/defense-in-depth.md`
- interface: `.agents/skills/systematic-debugging/find-polluter.sh`
- interface: `.agents/skills/systematic-debugging/root-cause-tracing.md`
- interface: `.agents/skills/systematic-debugging/SKILL.md`
- interface: `.agents/skills/systematic-debugging/test-pressure-1.md`
- interface: `.agents/skills/systematic-debugging/test-pressure-2.md`
- interface: `.agents/skills/systematic-debugging/test-pressure-3.md`
- interface: `.agents/skills/testing-boss/references/ai-writes-tests.md`
- interface: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/tui-design/references/app-patterns.md`
- interface: `.agents/skills/tui-design/references/visual-catalog.md`
- interface: `.agents/skills/tui-design/SKILL.md`
- interface: `skills-lock.json`
- interface: `skills/baseline_skill_contract_test.go`
- creates: `internal/baseline/repository_skills_snapshot_test.go`
- instruction: `internal/baseline/skills_trailing.go`
- instruction: `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestUpstreamADRFormatUnchanged|TestAuthoringConstraintOwnership|TestAuthorialSkillSync|TestUpstreamManagedSkillLockNamesAnEditedTree)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestUpstreamADRFormatUnchanged TestAuthoringConstraintOwnership TestAuthorialSkillSync TestUpstreamManagedSkillLockNamesAnEditedTree; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the lock-mismatch test does not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -qF -- '--- PASS: TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot (' || { printf '%s\n' "$out"; exit 1; }` — expected: exit 0; before this Task the test does not exist, so the command fails.
- `! grep -Eq 'upstreamManagedSkillTreeDigest|8832b7acd7fb65ec' skills/baseline_skill_contract_test.go && test ! -e .agents/skills/domain-modeling/CONTEXT-FORMAT.md && test -f .agents/skills/domain-modeling/GLOSSARY-FORMAT.md` — expected: exit 0; before this Task the pinned constant and the trailing tree are present, so the command fails.

## References

- `_prd.md` → Goals 3-4; Core Feature 2; Success Metrics 2-3; Acceptance
  evidence
- `_techspec.md` → Testing Approach 3-4; Build Order 2; Risks & Considerations
- [references/2026-10-02-restoring-trailing-skills-breaks-the-pinned-skill-digests.md](references/2026-10-02-restoring-trailing-skills-breaks-the-pinned-skill-digests.md)
- ADR-0224; ADR-0221; ADR-0191; ADR-0179
