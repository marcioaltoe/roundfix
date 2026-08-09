---
task: task_05
spec: 0082-the-manifest-already-answered-that
status: completed
type: backend
complexity: medium
---

# Task 05: Carry the repository's skills with its guidance

## Overview

Generated guidance dispatches to skills, so a refresh that updates one and not
the other leaves a repository pointing at skills it does not have. This task
adds the skills stage to the update command: the binary-carried owned bundle is
reinstalled, and external Repository Skill Set members that are missing or
drifted are restored through their existing preview-and-confirm contract. An
unreachable upstream degrades to a named warning rather than failing the
guidance refresh.

## Requirements

1. MUST reinstall the Roundfix-owned skill bundle carried in the binary into the
   repository's project skill directory as part of an applied update.
2. MUST restore external Repository Skill Set members that are missing or
   drifted, driving the existing two-step preview-then-confirm restoration
   rather than bypassing it.
3. MUST degrade to a warning, never a nonzero exit, when the immutable upstream
   source for an external skill cannot be reached; the guidance refresh still
   completes and reports success for what it did.
4. MUST report every skill installed, every skill restored, and every skill left
   drifted with the reason for each.
5. MUST keep the skills outcome on its own status axis, so an applied guidance
   refresh with a drifted unreachable skill is not reported as a plain success
   and not as a failure.
6. MUST support suppressing the skills stage entirely, so a maintainer can
   refresh guidance alone.
7. MUST support an explicitly declared offline source for external restoration.
8. MUST keep its own tests hermetic by injecting the skills stage rather than
   reaching the network.

## Subtasks

- [x] Add the skills stage and its seam for test injection.
- [x] Reinstall the owned bundle into the project skill directory.
- [x] Drive external restoration through preview and confirm.
- [x] Degrade unreachable sources to per-skill warnings.
- [x] Project the outcome into the result document on its own status axis.
- [x] Add the suppression and offline-source flags with tests.

## Acceptance Criteria

- [x] An applied update reinstalls the owned skill bundle and the result names
      the count installed.
- [x] An external skill that is missing or drifted is restored, and the result
      names it.
- [x] With the upstream source unreachable, the command still applies the
      guidance refresh, exits successfully, and names the drifted skill and the
      reason it was not restored.
- [x] The skills axis in the result is distinguishable from the apply axis; a
      warning on one does not change the other.
- [x] With the skills stage suppressed, no skill directory changes and the result
      says the stage was skipped.
- [x] The command's own tests pass with no network access.

## Context

- interface: `skills/skills.go`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/skills_restore_git.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/cli/ -run 'BaselineUpdate.*Skill|Skill.*BaselineUpdate' -v > /tmp/task_05-1.log 2>&1 && grep -q '^--- PASS: ' /tmp/task_05-1.log` — expected: exits 0, proving the skills-stage cases exist and pass.
- `go test ./internal/cli/ -run 'BaselineUpdate' -v > /tmp/task_05-2.log 2>&1 && grep -q -i 'unreachable' /tmp/task_05-2.log` — expected: exits 0, proving the degradation case ran.
- `go test ./internal/baseline/ ./internal/cli/ ./skills/ -count=1` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix baseline update --help` — expected: exits 0 and the usage names the skills suppression and offline-source flags.

## References

- `_techspec.md` → Build Order 6; Integration Points; Data Models.
- `_prd.md` → Core Feature 5; User Story 4; Goal 3.
- ADR-0068, ADR-0087.

## Result

Implemented the post-apply skills stage for `roundfix baseline update`. An
approved guidance refresh now reinstalls the binary-carried Roundfix-owned
bundle into the repository's project skill directory, checks the external
Repository Skill Set locally, and sends only missing or drifted members through
the existing restoration preview followed by confirmation of that preview's
exact Plan Digest. Restoration continues per member, so an unreachable
immutable source leaves that named skill drifted with its reason while other
members can still restore.

The `roundfix/baseline-update-result/v1` document now carries a separate
`skills` status axis with installed, restored, and drifted member lists. An
unreachable source produces `skills.status=warning` beside the unchanged
verified apply status and exit 0. `--no-skills` records `skipped` without
calling the stage, while `--skills-source-dir` forwards an absolute offline Git
checkout or bare object-store path. Command tests inject the complete stage;
the stage orchestration tests inject installation, readiness, and restoration
functions, so no focused test reaches a network.

The initial focused compile signal was red with the expected undefined
skills-stage request, result, status, and injection symbols. The sandbox also
denied the default macOS Go build cache, so subsequent focused checks used the
worktree-local ignored `.gocache`:

- `gofmt -w internal/cli/baseline_update.go internal/cli/baseline_update_test.go internal/cli/cli.go` — exit 0.
- The six-test focused skills selection in `./internal/cli` covering install
  reporting, restore reporting, unreachable degradation, independent status
  axes, suppression, offline-source forwarding, and preview-then-confirm — exit
  0 (`ok roundfix/internal/cli`).
- The five-test existing update regression selection covering applied JSON,
  idempotence, reviewed-digest confirmation, help, and exit/cancellation paths
  — exit 0 (`ok roundfix/internal/cli`).
- `git diff --check` — exit 0.

Acceptance evidence:

1. `TestBaselineUpdateSkillStageReportsInstalledAndRestoredSkills` observes
   `installedCount=1` and the installed `roundfix` name after an applied update;
   `TestBaselineUpdateSkillsStageUsesPreviewThenConfirmation` proves the
   production stage requests the `project` install target at the repository
   root.
2. `TestBaselineUpdateSkillsStageUsesPreviewThenConfirmation` records two
   restoration calls for `context7`: an unconfirmed preview, then a second call
   carrying the exact returned Plan Digest. The stage result names `context7`
   in `restored`.
3. `TestBaselineUpdateSkillsStageDegradesUnreachableSourcePerSkill` leaves
   unreachable `context7` in `drifted` with its reason, continues to restore
   `testing-boss`, and returns no error. `TestBaselineUpdateSkillWarningKeepsApplyAxisVerified`
   observes command exit 0 for the assembled warning case.
4. `TestBaselineUpdateSkillWarningKeepsApplyAxisVerified` observes top-level
   apply state `verified` and `approvedPostimages=verified` beside
   `skills.status=warning`, rather than collapsing the two axes into success or
   failure.
5. `TestBaselineUpdateSkipsSkillStageAndPreservesSkillDirectory` injects a
   stage that fails the test if called, runs with `--no-skills`, observes
   `skills.status=skipped`, and compares the complete project skill-directory
   tree digest before and after.
6. All six Task 05-focused command and stage cases use injected functions and
   passed without network access. `TestBaselineUpdatePassesOfflineSourceToSkillStage`
   additionally proves the explicit offline source reaches the stage, and the
   preview/confirm test proves it reaches both restoration calls.

The authored `## Verification` commands were not run; the Daemon owns them and
Task settlement. The spec's later documentation/owned-skill synchronization
slice remains responsible for teaching the shipped skills about these new
flags; this Task did not edit skill sources, mirrors, module assets, sibling
Tasks, or the Task Graph.
