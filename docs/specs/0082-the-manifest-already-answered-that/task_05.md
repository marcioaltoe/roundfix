---
task: task_05
spec: 0082-the-manifest-already-answered-that
status: pending
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

- [ ] Add the skills stage and its seam for test injection.
- [ ] Reinstall the owned bundle into the project skill directory.
- [ ] Drive external restoration through preview and confirm.
- [ ] Degrade unreachable sources to per-skill warnings.
- [ ] Project the outcome into the result document on its own status axis.
- [ ] Add the suppression and offline-source flags with tests.

## Acceptance Criteria

- [ ] An applied update reinstalls the owned skill bundle and the result names
      the count installed.
- [ ] An external skill that is missing or drifted is restored, and the result
      names it.
- [ ] With the upstream source unreachable, the command still applies the
      guidance refresh, exits successfully, and names the drifted skill and the
      reason it was not restored.
- [ ] The skills axis in the result is distinguishable from the apply axis; a
      warning on one does not change the other.
- [ ] With the skills stage suppressed, no skill directory changes and the result
      says the stage was skipped.
- [ ] The command's own tests pass with no network access.

## Context

- interface: `skills/skills.go`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/skills_restore_git.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/cli/ -run 'BaselineUpdate.*Skill|Skill.*BaselineUpdate' -v 2>&1 | grep -q '^--- PASS: '` — expected: exits 0, proving the skills-stage cases exist and pass.
- `go test ./internal/cli/ -run 'BaselineUpdate' -v 2>&1 | grep -q -i 'unreachable'` — expected: exits 0, proving the degradation case ran.
- `go test ./internal/baseline/ ./internal/cli/ ./skills/ -count=1` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix baseline update --help` — expected: exits 0 and the usage names the skills suppression and offline-source flags.

## References

- `_techspec.md` → Build Order 6; Integration Points; Data Models.
- `_prd.md` → Core Feature 5; User Story 4; Goal 3.
- ADR-0068, ADR-0087.
