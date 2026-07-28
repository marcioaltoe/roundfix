---
task: task_08
spec: 0042-verification-capacity-and-daemon-task-settlement
status: completed
type: docs
complexity: low
---

# Task 08: Align the protected authorial Skill pairs

## Overview

Publish the completed Daemon-owned status, Verification Capacity, temporary
failure, and Task Type-routing contracts in the two protected authorial Skills.
This is a tooling-only slice bounded to the four exact canonical/generated
files authorized by the maintainer and this Task file.

## Requirements

1. MUST align `.agents/skills/implement-task/SKILL.md` with the
   implementation-ready handoff: the Agent does not edit Task status, run the
   declared `## Verification`, or claim the terminal Task verdict.
2. MUST align `.agents/skills/roundfix/SKILL.md` with independent capacities,
   Daemon-owned settlement, observable Verification phases, exit `75`, one
   exclusive retry, and ADR-0051 Task Type-selected Agent Sessions.
3. MUST apply identical content to `skills/implement-task/SKILL.md` and
   `skills/roundfix/SKILL.md`.
4. MUST limit repository mutations to those four authorized `SKILL.md` files,
   this `task_08.md` file, and the maintainer-authorized derived Skill-digest
   pins that the authorized Skill edit deterministically invalidates.
5. MUST NOT run `make skills-sync`, because it rewrites every owned Skill
   directory; use `make skills-sync-check` as read-only verification.
6. MUST leave code, tests, configuration, manifests, public docs, other
   Roundfix-owned Skills, upstream-managed Skills, and lock files unchanged
   beyond those maintainer-authorized derived digest pins.

## Subtasks

- [x] Update the canonical `implement-task` Skill.
- [x] Apply the identical generated `implement-task` copy.
- [x] Update the canonical Roundfix Skill.
- [x] Apply the identical generated Roundfix copy.
- [x] Verify changed-file scope, byte identity, shipped Skill contracts, and
      full-gate compatibility.

## Acceptance Criteria

- [x] No supported Skill tells an Implement Agent to run declared Task
      Verification, edit Task status, or settle its terminal verdict.
- [x] The Roundfix Skill describes both capacities, the exit-75 retry contract,
      and per-Task Task Type routing consistently with shipped behavior.
- [x] Each canonical/generated Skill pair is byte-identical.
- [x] Git changed-file evidence for this Task contains only the four authorized
      Skill files, `task_08.md`, and the maintainer-authorized derived
      Skill-digest pins named in the active Spec artifacts' Tooling authority
      entries.
- [x] No other protected or upstream-managed tooling change.
- [x] Shipped Skill validation and the complete repository gate pass.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/skill-dispatch.md`
- instruction: `docs/agents/autonomous-work.md`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `skills/implement-task/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `rtk cmp .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md`
  — expected: no output and exit zero.
- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — expected: no output and exit zero.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/implement-task/SKILL.md" && path != "skills/implement-task/SKILL.md" && path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "docs/specs/0042-verification-capacity-and-daemon-task-settlement/task_08.md" && path != "internal/baseline/assets/setups/go-cli.json" && path != "internal/baseline/assets/setups/rust-cli.json" && path != "internal/baseline/assets/setups/typescript-bun.json" && path != "internal/baseline/testdata/catalog.normalized.json" && path != "internal/baseline/testdata/catalog.digest" && path != "internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json" && path != "internal/baseline/testdata/parity-corpus/v1/manifest.json") {print; bad=1}} END {exit bad}'`
  — expected: no changed path outside the four authorized Skill files, this
  Task file, and the maintainer-authorized derived digest pins (the three
  Baseline setup assets, the catalog identity fixtures, and the parity-corpus
  fixture/manifest rows that pin the `implement-task` and `roundfix` Skill
  `contentDigest`).
- `rtk make skills-sync-check` — expected: every canonical/generated owned
  Skill pair has no drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: every
  shipped Roundfix Skill contract passes.
- `rtk git diff --check` — expected: no whitespace errors.
- `rtk make verify` — expected: formatting, tests, Skill synchronization,
  shipped Skill validation, and build all pass.

## References

- `_prd.md` → Core Feature 10; Decisions; Project Constraints.
- `_techspec.md` → Integration Points; Build Order 8; Decisions.
- `task_02.md` → implemented Daemon status and handoff contract.
- `task_07.md` → completed operator wording and ADR-0051 alignment.
- `docs/agents/spec-routing.md` → tooling authorization and changed-file
  postflight.

## Result

Aligned both protected authorial Skill pairs with the implementation-ready
handoff and shipped Daemon contracts, propagated the maintainer-authorized
derived Skill-digest fallout through its full chain, and the complete
repository gate is green.

Daemon-assigned Implement Agents no longer edit Task status, run declared Task
Verification, claim a terminal verdict, or commit. The Roundfix Skill now
documents independent Task and Verification capacities, observable Verification
phases, deterministic Verification Feedback, the exit-75 exclusive retry, and
Task Type-selected Agent Sessions.

The Task first settled `failed` because the complete repository gate required
tooling mutations outside the original five-file allowlist: three Baseline
setup snapshots still pinned the pre-edit `contentDigest` of the two changed
authorial Skills, so `TestAuthorialSkillSync` reported four stale-digest cases.

### First-pass acceptance evidence (before maintainer authorization)

- `rtk go test ./internal/cli -run 'TestRunSkills(Check|ListSeparatesOwnedFromExternal)$' -count=1`
  passed: both operational embedded-Skill tests passed.
- `rtk go test ./skills -count=1` failed: 109 tests passed, 4
  `TestAuthorialSkillSync` cases failed, and 1 test skipped. The only remaining
  diagnostics were stale Baseline setup snapshot digests.
- `rtk proxy go test ./skills -run '^TestAuthorialSkillSync$' -count=1 -v`
  reproduced the blocker: all canonical/generated/embedded pair subtests
  passed, while the three setup snapshots expected the old `implement-task`
  digest `9c71ad...f6119` instead of `009377...d99ef9`; the TypeScript/Bun
  snapshot also expected the old Roundfix digest `0e2324...edfec` instead of
  `136a65...547d4`.
- `rtk shasum -a 256 .agents/skills/implement-task/SKILL.md
  skills/implement-task/SKILL.md .agents/skills/roundfix/SKILL.md
  skills/roundfix/SKILL.md` confirmed matching pair digests:
  `f18bc2...f9628` for `implement-task` and `0cf295...06e0` for `roundfix`.
- `rtk git -c core.fsmonitor=false status --short` and
  `rtk git -c core.fsmonitor=false diff --name-only` listed only the four
  authorized Skill files and this Task file.
- A changed-line whitespace scan over `rtk git -c core.fsmonitor=false diff
  --no-color` exited zero.
- `make skills-sync` was not run.

### Maintainer-authorized derived digest propagation

On 2026-07-28 the maintainer expressly authorized the deterministic
Skill-digest fallout of this Spec's authorized Skill edit in exactly the seven
derived paths now named in this Task's changed-file allowlist. Every value
below was taken from what the validator or the tests reported as canonical, or
regenerated from the embedded catalog; none were invented.

- The `implement-task` entry's `contentDigest` moved from
  `9c71ad3969ef4394b399421419437a9f00843b23b7fbdc524072be8e242f6119` to the
  canonical
  `009377be9160cbd9447c8a165e1c8d97df6981c2306fbbfdc6708f5294d99ef9` reported
  by `TestAuthorialSkillSync` in all three Baseline setup snapshots:
  `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`, and
  `internal/baseline/assets/setups/typescript-bun.json`.
- The `roundfix` entry's `contentDigest` in
  `internal/baseline/assets/setups/typescript-bun.json` moved from
  `0e2324f9458c861a0c1d932cd6b2bf7246ced43b2eefbcab6c570408410edfec` to the
  canonical
  `136a65201dfc3e15e2589506ed1400c123458c1f5c3ff9037fa20fc1ddc547d4`. The other
  two snapshots do not carry the operational Roundfix Skill, so only their
  `implement-task` row moved.
- Each setup's own canonical top-level `digest` moved to the value the catalog
  validator reports through `catalog.setup.digest.mismatch`:
  `go-cli` from `a63d4f17fb66a345d48ae64e8f199d28f52bf7727e92e47c0a3bbf80fae967cc`
  to `a27991b2a6dedede72258a2b4836c1c656bebcd6e49703bab7799a265b7d79f9`;
  `rust-cli` from `cdd07e328d08b7e24708f499cb7e40884b994a72df046ea64e73e1d078a85c05`
  to `b94ded3566f90f27f5a304cfce797993c6bb4175a4cdbadebb7807e909558ca0`;
  `typescript-bun` from `87d4652856c35c9cd37e3b0523167d02d20e3507fcfe445b8debcd2ae19038b3`
  to `ffcaa02c9934deed7c7061e152aad840456a4c892bae7c4684b56e921f441622`.
- The four `contentDigest` rows in
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` moved
  to the same canonical Skill digests: three `implement-task` rows (one per
  profile) and one `roundfix` row.
- The catalog identity fixtures were regenerated from `Catalog.Normalized()`
  and `Catalog.Digest()` of the embedded catalog through a throwaway in-package
  generator that was deleted afterwards and is absent from `git status`:
  `internal/baseline/testdata/catalog.normalized.json` (file digest
  `setups/go-cli.json` `sha256:28e2ad53…` → `sha256:960dd8f8…` at 14,203 bytes,
  `setups/rust-cli.json` `sha256:82cdba03…` → `sha256:ff07c3e1…` at 12,476
  bytes, and `setups/typescript-bun.json` `sha256:923273a8…` →
  `sha256:6e9a4822…` at 40,649 bytes, all byte counts unchanged) and
  `internal/baseline/testdata/catalog.digest`
  (`sha256:15b0e92799339f4148d51c0da7857075fb55fe301e1810c810824337509a2551` →
  `sha256:18476427463bcaa2e8d4a3d974550e923190dfc906418c5252ab41af2c39f551`).
- The `fixtures/asset-sync.json` sha256 row in
  `internal/baseline/testdata/parity-corpus/v1/manifest.json` moved from
  `9d5e5bad6cb6d8f20b28dfac3fe8a07834399495fd9c74cb0bee08e35a01aaab` to
  `3726b1fc93ac04fd39a953f0246c66af20b3e632dc521c66ae9857a80d7eb85e`, the
  identity `TestBaselineCompatibilityCorpus` reports, with the byte count
  unchanged at 83,662.

No `SKILL.md`, `.roundfixrc.yml`, lock file, recommendation file, manifest Task
file, public document, or configuration outside those derived pins changed in
this settlement pass, and `make skills-sync` was still not run.

### Acceptance evidence

1. Both `implement-task` copies distinguish standalone settlement from a
   Daemon-assigned handoff and prohibit Daemon-assigned status, declared
   Verification, terminal-verdict, and commit authorship.
2. Both Roundfix copies describe `worktree.concurrency` Task Capacity,
   `verification.concurrency` Verification Capacity, `waiting` and `started`
   phases, one Verification Feedback repair, exit `75`, one exclusive retry,
   and Task Type-selected Agent Sessions.
3. `cmp` evidence confirms each canonical/generated pair is byte-identical.
4. Git evidence confirms the changed-file scope contains exactly the four
   authorized Skill paths, this Task file, and the seven maintainer-authorized
   derived digest pins.
5. No other protected Skill, upstream-managed Skill, code, test, config, public
   doc, or lock file changed.
6. Shipped Skill validation and the complete repository gate both pass.

### Final verification evidence

Run in the Run Worktree with `GOFLAGS=-buildvcs=false` and a portable
`GOCACHE`:

- `go test ./skills/... -count=1` — passed: `ok  roundfix/skills  0.674s`.
  `TestAuthorialSkillSync` and all eighteen subtests report PASS, including the
  previously failing `go-cli.json`, `rust-cli.json`, and `typescript-bun.json`
  setup-snapshot cases.
- `go test ./internal/baseline -count=1` — passed:
  `ok  roundfix/internal/baseline  82.373s`.
- `go test ./...` — full suite green: every package reports `ok`
  (`internal/baseline 101.215s`, `internal/cli 138.222s`, `skills 3.513s`, and
  all remaining packages), with no failures.
- `make verify` — passed end to end: the formatting check was clean,
  `go test ./...` reported `2727 passed in 23 packages`, `skills-sync-check`
  reported `4 passed in 1 packages` with no canonical/generated drift,
  `go run -buildvcs=false ./cmd/roundfix skills check` reported
  `Roundfix skill check passed` for all fourteen shipped Skills, and the build
  produced `bin/roundfix`.
- `gofmt -l` over every tracked Go file printed nothing; no Go source changed
  in this settlement pass.
- `cmp .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md`
  and `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — no
  output and exit zero; both pairs are still byte-identical.
- `git -c core.fsmonitor=false diff --check` — exit zero, no whitespace errors.
- `git status --porcelain` filtered through this Task's allowlist exited zero
  with no output.

### Follow-up

Resolved. The maintainer expressly authorized updating the derived
Skill-digest pins as mechanical fallout of the authorized Skill edit, this
Task's changed-file allowlist now names the same seven derived paths, the pins
are canonical again, and the complete repository gate is green.
