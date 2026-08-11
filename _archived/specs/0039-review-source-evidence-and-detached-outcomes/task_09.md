---
task: task_09
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
type: docs
complexity: low
---

# Task 09: Align the protected Roundfix Skill pair

## Overview

Publish Review Source Evidence and Detached outcome handling in the canonical
and generated Roundfix Skill. This tooling-only slice is restricted to the two
exact files authorized by the maintainer and this Task file.

## Requirements

1. MUST align `.agents/skills/roundfix/SKILL.md` with Review Skipped, bounded
   transient retry, unknown issue knowledge, Detached monitoring, notification
   receipts, and artifact-only Evidence inheritance.
2. MUST apply byte-identical content to `skills/roundfix/SKILL.md`.
3. MUST limit repository changes to those two authorized `SKILL.md` files and
   this `task_09.md` file.
4. MUST NOT run `make skills-sync`, because it rewrites every owned Skill
   directory; use the read-only sync check.
5. MUST leave code, tests, manifests, public docs, other owned Skills,
   upstream-managed Skills, locks, and recommendation files unchanged.

## Subtasks

- [x] Update the canonical Roundfix Skill review-evidence contract.
- [x] Apply the identical generated Roundfix Skill copy.
- [x] Verify exact changed-file scope and byte identity.
- [x] Confirm shipped Skill and full-gate compatibility.

## Acceptance Criteria

- [x] The Skill distinguishes Review Skipped and unknown issue knowledge from
      successful zero-issue review.
- [x] The Skill surfaces the Supervisor monitor command and exact next actions.
- [x] Canonical and generated files are byte-identical.
- [x] Git evidence contains only the two authorized Skill paths, this Task
      file, the maintainer-authorized derived Skill-digest pins, and the
      schema-version test files the gate required.
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
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "docs/specs/0039-review-source-evidence-and-detached-outcomes/task_09.md" && path != "internal/baseline/assets/setups/typescript-bun.json" && path != "internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json" && path != "internal/baseline/testdata/catalog.normalized.json" && path != "internal/baseline/testdata/catalog.digest" && path != "internal/baseline/testdata/parity-corpus/v1/manifest.json" && path != "internal/cli/cli_test.go" && path != "internal/store/store_test.go") {print; bad=1}} END {exit bad}'`
  — expected: no changed path outside the authorized pair, this Task file, the
  maintainer-authorized derived digest pins (the baseline setup asset, its
  catalog identity, and the parity-corpus fixture/manifest rows that pin the
  roundfix Skill `contentDigest`), and the two schema-version test files whose
  migration assertions were made durable against the supported schema version.
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
- `_techspec.md` → Build Order 9; Decisions.
- `docs/agents/spec-routing.md` → tooling authorization and changed-file
  postflight.

## Result

The canonical and generated Roundfix Skill now publish Review Source Evidence
and Detached outcome handling, the derived Skill-digest pins are canonical
again, and the complete repository gate is green.

The Task first settled failed because `make verify` had three failures outside
the original mutation allowlist, from two independent causes. Both are now
resolved: the Run Database migration assertions no longer hardcode a schema
version, and the maintainer-authorized Skill-digest fallout was propagated
through its full derived chain.

### Changes

- The canonical and generated Roundfix Skills now distinguish explicit Review
  Skipped, unknown pre-fetch Review Issue knowledge, and a successful fetched
  zero-issue result.
- The Skills document head-bound Review Source Evidence, bounded typed
  transient retry episodes, the five-line Detached report, the exact
  Supervisor outcome subscription, terminal notification context and receipt
  states, and exact Daemon artifact-only Evidence inheritance.
- Run Database migration assertions now compare against the schema version the
  binary actually supports instead of a literal that every future migration
  invalidates.
- The maintainer-authorized derived Skill-digest pins carry the canonical
  roundfix `contentDigest` and the identity values it feeds.

### Acceptance evidence

- Review Skipped and issue knowledge: the Skill prints
  `Review Issues: unknown — fetch did not complete.` only before fetch
  completion, preserves count lines for known-zero fetches, and documents the
  dedicated Review Skipped reason plus
  `Next action: Reduce or split the pull request, then request another Review Source review.`
- Supervisor recovery: the Detached report contains
  `Supervisor monitor: roundfix events <run-id> --follow --filter outcome`;
  non-Clean outcome records carry bounded reason and next action when present.
- Byte identity: `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  exited `0` with no output after the last Skill edit.
- Changed-file scope:
  `rtk git -c core.fsmonitor=false status --porcelain` filtered through this
  Task's allowlist exited `0` with no output. The only changed paths are the
  two authorized Skill files, this Task file, the five maintainer-authorized
  derived digest pins, and the two schema-version test files.
- Protected ownership: the post-gate Git status contains no other owned Skill,
  upstream-managed Skill, lock, recommendation, or public-document path.
- Shipped Skill validation:
  `rtk go run -buildvcs=false ./cmd/roundfix skills check` reported
  `Roundfix skill check passed` for all fourteen shipped Skill contracts.
- Whitespace: `rtk git diff --check` exited `0`.
- Read-only pair synchronization: `rtk make skills-sync-check` reported
  `4 passed in 1 packages` with no canonical/generated drift, and
  `rtk make skills-sync` was never run.

### First-pass gate failures

`rtk make verify` first left three failures from two independent causes:

1. `TestBranchIntegrityPreflightMigratesOutdatedRunDatabase` in
   `internal/cli/cli_test.go` asserted the literal schema version `10` while
   this Spec's Review Skipped migration had already advanced the supported
   version to `11`.
2. `TestAuthorialSkillSync` and its `typescript-bun.json` subtest reported the
   stale Roundfix `contentDigest`
   `91b833ef01c723b308604e57dc4075ec8e216880c8d50cf493d7dbced7096f6d`
   instead of canonical
   `0e2324f9458c861a0c1d932cd6b2bf7246ced43b2eefbcab6c570408410edfec`.

### Durable schema-version assertions

The migration assertion was repaired at its root, not pinned to `11`.
`seedOutdatedV9RunDatabase` now reads `MigrationVersion` from the freshly
created Run Database before reversing the v9→v10 migration and returns that
value, so `TestBranchIntegrityPreflightMigratesOutdatedRunDatabase` asserts the
preflight migrated to the version this binary supports. The helper still seeds
the concrete historical `PRAGMA user_version = 9` fixture, which is correct;
only the result assertion became version-independent. The `store` package's
unexported `schemaVersion` constant is unreachable from `internal/cli` and the
package exports no equivalent accessor, so the supported version is read from
the migrated database itself rather than duplicated as a literal.

The same criterion was applied to the sibling assertions in
`internal/store/store_test.go`, which is in-package and can name the constant
directly: the six `version != 11` result assertions in
`TestOpenCreatesRunDatabaseAndAppliesMigrations` and the v3, v4, v5, v6, and v7
migration tests now compare against `schemaVersion`, and
`TestSchemaReviewSkippedReaderRejectsNewerDatabase` seeds and asserts
`schemaVersion + 1` instead of the literal `12` that the next migration would
turn into a supported version. Historical fixture seeds (`buildV9Fixture`,
`PRAGMA user_version = 3` through `9`) stayed concrete.

### Maintainer-authorized derived digest propagation

On 2026-07-28 the maintainer expressly authorized the deterministic
Skill-digest fallout of the authorized Skill edit. Every value below was taken
from what the validator or test reported as canonical, or regenerated from the
embedded catalog; none were invented.

- The roundfix entry's `contentDigest` in
  `internal/baseline/assets/setups/typescript-bun.json` moved from
  `91b833ef01c723b308604e57dc4075ec8e216880c8d50cf493d7dbced7096f6d` to the
  canonical
  `0e2324f9458c861a0c1d932cd6b2bf7246ced43b2eefbcab6c570408410edfec` reported
  by `TestAuthorialSkillSync`.
- The setup's own canonical `digest` in the same file moved from
  `fe98e52e2a6812b899bd4a048c29afc515e9736ddc4df7120bd8f9b1cf7d9896` to
  `87d4652856c35c9cd37e3b0523167d02d20e3507fcfe445b8debcd2ae19038b3`, the value
  the catalog validator reports through `catalog.setup.digest.mismatch`.
- The roundfix `contentDigest` row in
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` moved
  to the same canonical Skill digest.
- The catalog identity fixtures were regenerated from `Catalog.Normalized()`
  and `Catalog.Digest()` of the embedded catalog:
  `internal/baseline/testdata/catalog.normalized.json` (the
  `setups/typescript-bun.json` file digest `sha256:411768db…` →
  `sha256:923273a8…`, byte count unchanged at 40,649) and
  `internal/baseline/testdata/catalog.digest`
  (`sha256:c715657c891c29b108e58a73d20c8bb6b9647cf150852f3600e29b65a950a70f` →
  `sha256:15b0e92799339f4148d51c0da7857075fb55fe301e1810c810824337509a2551`).
  The throwaway in-package generator used for that regeneration was deleted and
  is absent from `git status`.
- The `fixtures/asset-sync.json` sha256 row in
  `internal/baseline/testdata/parity-corpus/v1/manifest.json` moved from
  `cf7dc89818d21d5f7bebb91c2a920f642058db8184285a6ec572c3d1f7827248` to
  `9d5e5bad6cb6d8f20b28dfac3fe8a07834399495fd9c74cb0bee08e35a01aaab`, the
  identity `TestBaselineCompatibilityCorpus` reports, with the byte count
  unchanged at 83,662.

No `SKILL.md`, `.roundfixrc.yml`, `_prd.md`, `_techspec.md`, manifest task
file, or configuration outside those derived pins changed in this settlement
pass.

### Final verification evidence

Run in the Run Worktree with `GOFLAGS=-buildvcs=false` and a portable
`GOCACHE`:

- `go test ./internal/cli -run 'TestBranchIntegrityPreflight' -count=1` —
  passed: `ok  roundfix/internal/cli  0.608s`.
- `go test ./skills/... ./internal/baseline -count=1` — passed:
  `ok  roundfix/skills  0.705s` and
  `ok  roundfix/internal/baseline  81.723s`.
- `go test ./...` — full suite green: every package reports `ok`, including
  `internal/baseline 106.728s`, `internal/cli 140.892s`,
  `internal/store 3.970s`, and `skills 4.138s`, with no failures.
- `make verify` — passed end to end: the silent `fmt-check` produced no output,
  `go test ./...` reported `2659 passed in 23 packages`, `skills-sync-check`
  reported `4 passed in 1 packages`,
  `go run -buildvcs=false ./cmd/roundfix skills check` reported
  `Roundfix skill check passed` for all fourteen shipped Skills, and the build
  produced `bin/roundfix`.
- `gofmt -l internal/cli/cli_test.go internal/store/store_test.go` — no output
  and exit zero.
- `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — no output
  and exit zero; the pair is still byte-identical.

### Follow-ups

Resolved. Both first-pass causes are closed inside this Task: the migration
assertions are durable against future schema versions, and the
maintainer-authorized derived Skill-digest pins are canonical. This Task's
changed-file allowlist names every path that changed, and the complete
repository gate is green.
