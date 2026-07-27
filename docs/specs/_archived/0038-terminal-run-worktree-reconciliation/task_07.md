---
task: task_07
spec: 0038-terminal-run-worktree-reconciliation
status: completed
type: docs
complexity: low
---

# Task 07: Align the protected Roundfix Skill pair

## Overview

Publish proof-based Run Worktree reconciliation in the canonical and generated
Roundfix Skill. This tooling-only slice is restricted to the two exact files
authorized by the maintainer and this Task file.

## Requirements

1. MUST align `.agents/skills/roundfix/SKILL.md` with dry-run-first
   reconciliation, five states, and explicit safe-only apply.
2. MUST apply byte-identical content to `skills/roundfix/SKILL.md`.
3. MUST limit repository changes to those two authorized `SKILL.md` files and
   this `task_07.md` file.
4. MUST NOT run `make skills-sync`, because it rewrites every owned Skill
   directory; use the read-only sync check.
5. MUST leave code, tests, manifests, public docs, other owned Skills,
   upstream-managed Skills, locks, and recommendation files unchanged.

## Subtasks

- [x] Update the canonical Roundfix Skill reconciliation contract.
- [x] Apply the identical generated Roundfix Skill copy.
- [x] Verify exact changed-file scope and byte identity.
- [x] Confirm shipped Skill and full-gate compatibility.

## Acceptance Criteria

- [x] The Skill tells an Agent to inspect before apply and never use manual Git
      deletion as the supported workflow.
- [x] Dirty, unintegrated, and unknown results are preserved in Skill guidance.
- [x] Canonical and generated files are byte-identical.
- [x] Git evidence contains only the two authorized Skill paths, this Task
      file, and the maintainer-authorized derived Skill-digest pins named in
      the active Spec artifacts' Tooling authority entries.
- [x] No other protected or upstream-managed Skill changes.
- [x] Shipped Skill validation and the complete repository gate pass.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/skill-dispatch.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — expected: no output and exit zero.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "docs/specs/0038-terminal-run-worktree-reconciliation/task_07.md" && path != "internal/baseline/assets/setups/typescript-bun.json" && path != "internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json" && path != "internal/baseline/testdata/catalog.normalized.json" && path != "internal/baseline/testdata/catalog.digest" && path != "internal/baseline/testdata/parity-corpus/v1/manifest.json") {print; bad=1}} END {exit bad}'`
  — expected: no changed path outside the authorized pair, this Task file, and
  the maintainer-authorized derived digest pins (the baseline setup asset, its
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

The canonical and generated Roundfix Skill now direct Agents to inspect with a
dry-run before apply, define all five Run Worktree Reconciliation states,
preserve `unintegrated`, `dirty`, and `unknown` work, allow cleanup only for
freshly revalidated `safe` entries, and prohibit manual Git or filesystem
deletion as the supported workflow.

### First-pass acceptance evidence (before maintainer authorization)

1. `rtk rg -n '^## Run Worktree reconciliation|Always run a dry-run|safe|unintegrated|dirty|unknown|released|Never substitute manual Git deletion' .agents/skills/roundfix/SKILL.md`
   found the dry-run instruction, all five states, preservation rules, and the
   manual-deletion prohibition.
2. `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited
   zero with no output.
3. `rtk git -c core.fsmonitor=false status --porcelain` and
   `rtk git -c core.fsmonitor=false diff --name-only` reported only the two
   authorized Skill paths and this Task file.
4. `rtk make skills-sync-check` passed: four focused Skill tests passed and
   every canonical/generated owned Skill pair had no drift.
5. `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed for every
   shipped Roundfix-owned Skill.
6. `rtk git diff --check` passed with no whitespace errors.
7. `rtk make verify` failed during `go test ./...`: 2,556 tests passed, seven
   failed, and two skipped. Five store/CLI failures could not read process
   identity because the host could not fork `/usr/bin/ps`; the deterministic
   blocking failure was
   `TestAuthorialSkillSync/typescript-bun.json`.
8. `rtk proxy go test -count=1 ./skills -run 'TestAuthorialSkillSync/typescript-bun.json'`
   reproduced the deterministic failure: the snapshot records Roundfix Skill
   digest `fa774d82f16661c81c235738c78116fba2fdc328bac5797fd663fa31a05f42d6`,
   while the updated canonical Skill requires
   `91b833ef01c723b308604e57dc4075ec8e216880c8d50cf493d7dbced7096f6d`.

### Maintainer-authorized derived digest propagation

The task first settled failed because the complete repository gate required an
out-of-scope tooling mutation: `TestAuthorialSkillSync/typescript-bun.json`
expected the previous Roundfix Skill digest pinned in
`internal/baseline/assets/setups/typescript-bun.json`, and changing that
snapshot was outside this Task's exact mutation allowlist.

On 2026-07-27 the maintainer expressly authorized the deterministic
Skill-digest fallout of the authorized Skill edit in exactly the five derived
paths now named in the `_prd.md` and `_techspec.md` Tooling authority entries.
Under that authorization the roundfix Skill entry's `contentDigest` in
`internal/baseline/assets/setups/typescript-bun.json` was updated from
`fa774d82f16661c81c235738c78116fba2fdc328bac5797fd663fa31a05f42d6` to the
canonical
`91b833ef01c723b308604e57dc4075ec8e216880c8d50cf493d7dbced7096f6d`. Because
that value feeds further derived identity pins, the same mechanical fallout
propagated to:

- the setup's own canonical `digest` in the same file
  (`36f512c6c0370aab357c345a8bf2ceb902738ff837bdd34da7c1c2085567533c` →
  `fe98e52e2a6812b899bd4a048c29afc515e9736ddc4df7120bd8f9b1cf7d9896`, the value
  the catalog validator reports as canonical through
  `catalog.setup.digest.mismatch`);
- the roundfix `contentDigest` row in
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, moved
  to the same canonical Skill digest;
- the catalog identity fixtures
  `internal/baseline/testdata/catalog.normalized.json` (the
  `setups/typescript-bun.json` file digest `sha256:31c1b8c9…` →
  `sha256:411768db…`, byte count unchanged at 40,649) and
  `internal/baseline/testdata/catalog.digest`
  (`sha256:c0435d250c1440584454009211b5f044711382fb4d9a01e9dc4e97b91e3ca014` →
  `sha256:c715657c891c29b108e58a73d20c8bb6b9647cf150852f3600e29b65a950a70f`),
  both regenerated from `Catalog.Normalized()` and `Catalog.Digest()` of the
  embedded catalog;
- the `fixtures/asset-sync.json` sha256 row in
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`
  (`31620de106dfae46414ea11c8ff420e611b1c874c5d4e2936d8f37234142c59b` →
  `cf7dc89818d21d5f7bebb91c2a920f642058db8184285a6ec572c3d1f7827248`, byte count
  unchanged at 83,662).

No `SKILL.md`, manifest task file, or configuration outside those derived pins
changed in this settlement pass, and `make skills-sync` was still not run.

### Final verification evidence

Run in the Run worktree with `GOFLAGS=-buildvcs=false` and a portable
`GOCACHE`:

- `go test ./skills/... -count=1` — passed: `ok  roundfix/skills  0.604s`;
  `TestAuthorialSkillSync` and all its subtests, including
  `typescript-bun.json`, report PASS.
- `go test ./internal/baseline -count=1` — passed:
  `ok  roundfix/internal/baseline  82.706s`.
- `go test ./...` — full suite green: every package reports `ok`
  (`internal/baseline 109.188s`, `internal/cli 169.515s`, `skills 3.538s`, and
  all remaining packages), with no failures. The `/usr/bin/ps` fork failures
  seen in the first pass did not recur.
- `make verify` — passed end to end: formatting check clean, `go test ./...`
  reported `2563 passed in 23 packages`, `skills-sync-check` reported
  `4 passed in 1 packages` with no canonical/generated drift,
  `go run -buildvcs=false ./cmd/roundfix skills check` reported
  `Roundfix skill check passed` for all fourteen shipped Skills, and the build
  produced `bin/roundfix`.
- `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — no output
  and exit zero; the pair is still byte-identical.
- `git status --porcelain` — the only changed code and testdata paths are the
  two authorized `SKILL.md` files and the five maintainer-authorized derived
  digest pins, alongside this Task file and the `_prd.md` and `_techspec.md`
  Tooling authority entries that record the authorization.

### Follow-up

Resolved. The maintainer expressly authorized updating the derived Skill-digest
pins as mechanical fallout of the authorized Skill edit, the `_prd.md` and
`_techspec.md` Tooling authority entries record that authorization, and this
Task's changed-file allowlist now names the same five derived paths. The Skill
pair is aligned, the derived pins are canonical, and the complete repository
gate is green.
