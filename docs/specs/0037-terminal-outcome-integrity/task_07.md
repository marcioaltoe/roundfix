---
task: task_07
spec: 0037-terminal-outcome-integrity
status: completed
type: docs
complexity: low
---

# Task 07: Align the protected Roundfix Skill pair

## Overview

Publish Terminal Outcome Integrity in the canonical and generated Roundfix
Skill after product behavior and public guidance are complete. This is a
tooling-only slice bounded to the two exact files authorized by the maintainer
and this Task file.

## Requirements

1. MUST align `.agents/skills/roundfix/SKILL.md` with proof-before-Force-Stop,
   graceful wait interruption, registered-session cleanup, and winner-only
   terminal publication.
2. MUST apply byte-identical content to `skills/roundfix/SKILL.md`.
3. MUST limit repository changes to those two authorized `SKILL.md` files and
   this `task_07.md` file.
4. MUST NOT run `make skills-sync`, because it rewrites every owned Skill
   directory; use the read-only sync check.
5. MUST leave code, tests, manifests, public documentation, other owned Skills,
   upstream-managed Skills, locks, and recommendation files unchanged.

## Subtasks

- [x] Update the canonical Roundfix Skill terminal-outcome contract.
- [x] Apply the identical generated Roundfix Skill copy.
- [x] Verify exact changed-file scope and byte identity.
- [ ] Confirm shipped Skill and full-gate compatibility.

## Acceptance Criteria

- [x] The Roundfix Skill never tells an Agent to accept terminal overwrite or
      report Force Stop before owner exit proof.
- [x] The Skill surfaces graceful interruption and registered-session cleanup
      consistently with supported docs.
- [x] Canonical and generated files are byte-identical.
- [x] Git evidence contains only the two authorized Skill paths, this Task
      file, and the maintainer-authorized derived Skill-digest pins named in
      the active Spec artifacts' Tooling authority entries.
- [x] No other protected or upstream-managed Skill changes.
- [ ] Shipped Skill validation and the complete repository gate pass.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/skill-dispatch.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — expected: no output and exit zero.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "docs/specs/0037-terminal-outcome-integrity/task_07.md" && path != "internal/baseline/assets/setups/typescript-bun.json" && path != "internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json" && path != "internal/baseline/testdata/catalog.normalized.json" && path != "internal/baseline/testdata/catalog.digest" && path != "internal/baseline/testdata/parity-corpus/v1/manifest.json") {print; bad=1}} END {exit bad}'`
  — expected: no changed path outside the authorized pair, this Task file, and
  the supervisor-authorized derived digest pins (the baseline setup asset, its
  catalog identity, and the parity-corpus fixture/manifest rows that pin the
  roundfix Skill `contentDigest`).
- `rtk make skills-sync-check`
  — expected: every canonical/generated owned Skill pair has no drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check`
  — expected: every shipped Roundfix Skill contract passes.
- `rtk git diff --check`
  — expected: no whitespace errors.
- `rtk make verify`
  — expected: formatting, tests, Skill checks, and build pass.

## References

- `_prd.md` → Goals; User Experience; Decisions; Project Constraints.
- `_techspec.md` → Build Order 7; Decisions.
- `docs/agents/spec-routing.md` → tooling authorization and changed-file
  postflight.

## Result

Aligned the canonical and generated Roundfix Skill with failed-closed Force
Stop, graceful Review Source wait interruption, registered active Agent Session
cleanup, immutable terminal outcomes, and winner-only terminal event and
notification publication. The two Skill files remain byte-identical, and no
path outside the Task's exact mutation allowlist changed.

The task first settled failed because the complete repository gate required an
out-of-scope tooling mutation: `TestAuthorialSkillSync/typescript-bun.json`
expected the previous Roundfix Skill digest pinned in
`internal/baseline/assets/setups/typescript-bun.json`, and changing that
snapshot was outside this Task's exact mutation allowlist.

The supervisor subsequently granted explicit authorization to update that
derived digest pin as mechanical fallout of the authorized Skill edit. Under
that authorization the roundfix skill entry's `contentDigest` in
`internal/baseline/assets/setups/typescript-bun.json` was updated from
`1e4ea59e0d6e553e42c6c63e16ad95165a78be8bbf31b8e0cd8b56ce13cc9146` to the
canonical
`a1f01156d2ef6ecb020b54d2559965f626f791301c32f14d269f6a751c779cf9`. Because
that value feeds further derived identity pins, the same mechanical fallout
propagated to: the setup's own canonical `digest` in the same file
(`48592a56…` → `76d74ab2…`, the value the catalog validator reports as
canonical), the roundfix `contentDigest` row in
`internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, the
catalog identity fixtures
`internal/baseline/testdata/catalog.normalized.json` and
`internal/baseline/testdata/catalog.digest` (regenerated from
`Catalog.Normalized()`/`Catalog.Digest()` of the embedded catalog), and the
`fixtures/asset-sync.json` sha256 row in
`internal/baseline/testdata/parity-corpus/v1/manifest.json`. No SKILL.md,
manifest task file, or configuration outside those derived pins changed in
this settlement pass.

Final verification, run in the Run worktree with `GOFLAGS=-buildvcs=false`:

- `go test ./skills/... -count=1` — passed: `ok  roundfix/skills  0.600s`;
  `TestAuthorialSkillSync` and all its subtests, including
  `typescript-bun.json`, report PASS.
- `go test ./internal/baseline -count=1` — passed:
  `ok  roundfix/internal/baseline  88.931s`.
- `go test ./...` — full suite green: every package reports `ok` (skills
  `9.622s`, internal/baseline `126.952s`, internal/cli `157.899s`, and all
  remaining packages), with no failures.

The skill pair is aligned, the derived digest pins are canonical, and the full
suite is green.

### Acceptance criterion evidence

- The stale promises that Force Stop completes Stopped and releases its lock
  immediately are absent. The Skill now states that owner exit proof precedes
  the Stopped report and lock release, while failed proof leaves the Run Active
  with its lock retained.
- The Skill states that graceful stop is observed before Review Source access
  and after each interruptible sleep, no later than the next configured poll
  boundary. Cleanup targets only registered Agent Sessions whose latest
  lifecycle is active and treats an already-absent registered session as
  idempotent.
- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` passed
  with no output.
- Git status and diff evidence contained only
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, and this Task
  file.
- No other protected, owned, or upstream-managed Skill path changed.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed, but
  `rtk make verify` failed in `TestAuthorialSkillSync/typescript-bun.json`;
  therefore the combined shipped-Skill and complete-gate criterion is not
  satisfied.

### First-pass verification evidence (before supervisor authorization)

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — passed.
- The declared `rtk git status --porcelain | rtk awk ...` scope command exited
  `0` but emitted the worktree's `fsmonitor_ipc__send_query` diagnostic, so it
  was not treated as standalone proof. `rtk git -c core.fsmonitor=false status
  --short` and `rtk git -c core.fsmonitor=false diff --name-only` both passed
  and listed exactly the three authorized paths.
- `rtk make skills-sync-check` — passed; 4 focused Skill tests passed.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — passed for every
  shipped Roundfix Skill contract.
- `rtk git diff --check` — passed after the final Result update.
- `rtk make verify` — failed after 2,475 tests passed, 2 failed, and 2 skipped.
  Both failures are the parent and subtest for `TestAuthorialSkillSync`; the
  stored digest is
  `1e4ea59e0d6e553e42c6c63e16ad95165a78be8bbf31b8e0cd8b56ce13cc9146`,
  while the changed canonical Skill digest is
  `a1f01156d2ef6ecb020b54d2559965f626f791301c32f14d269f6a751c779cf9`.
- The Daemon owns the authoritative execution of the declared Verification
  after this Agent turn.

### Follow-up

Resolved. The supervisor explicitly authorized updating the derived
`contentDigest` pin in `internal/baseline/assets/setups/typescript-bun.json`
as mechanical fallout of the authorized Skill edit. The pin and its
downstream derived identity fixtures now carry the canonical values, and
`go test ./skills/... -count=1`, `go test ./internal/baseline -count=1`, and
`go test ./...` all pass in the Run worktree.
