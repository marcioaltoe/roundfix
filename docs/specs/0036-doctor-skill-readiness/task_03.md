---
task: task_03
spec: 0036-doctor-skill-readiness
status: pending
type: docs
complexity: low
---

# Task 03: Align the protected Roundfix Skill pair

## Overview

Publish Doctor Skill Readiness in the protected Roundfix Skill after the CLI
and public guidance are complete. This is a tooling-only slice: it may change
only the exact canonical/generated files authorized by the maintainer and this
Task file.

## Requirements

1. MUST update the Doctor instructions in
   `.agents/skills/roundfix/SKILL.md` with the blocking Repository Skill Set
   result, ownership-specific remediation, and diagnosis-only behavior.
2. MUST apply the same content to `skills/roundfix/SKILL.md` so the canonical
   and generated copies are byte-identical.
3. MUST limit repository mutations to
   `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
   `internal/baseline/assets/setups/typescript-bun.json`,
   `internal/baseline/testdata/catalog.normalized.json`,
   `internal/baseline/testdata/catalog.digest`,
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
   `internal/baseline/testdata/parity-corpus/v1/manifest.json`, and this
   `task_03.md` file.
4. MUST NOT run `make skills-sync`, because that mutation target rewrites every
   owned skill directory; `make skills-sync-check` is the read-only sync gate.
5. MUST update the five authorized Baseline and parity artifacts only as
   mechanically derived consequences of the canonical Roundfix Skill bytes.
6. MUST leave every external skill, lock file, recommendation file, source code
   file, and user-documentation file unchanged.

## Subtasks

- [ ] Update the canonical Roundfix Skill Doctor contract.
- [ ] Apply the identical change to the generated Roundfix Skill copy.
- [ ] Verify exact changed-file scope and byte identity.

## Acceptance Criteria

- [ ] The Roundfix Skill tells an Agent to surface a failed skill check and its
      printed remediation before work continues, without claiming Doctor
      performs an update.
- [ ] The canonical and generated Roundfix `SKILL.md` files are byte-identical.
- [ ] Git changed-file evidence for this Task contains only the two authorized
      Skill files, five derived Baseline and parity artifacts, and `task_03.md`.
- [ ] No upstream-managed or other Roundfix-owned skill changes.
- [ ] The shipped skill check and complete repository gate pass.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/skill-dispatch.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `internal/baseline/assets/setups/typescript-bun.json`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`

## Verification

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — expected: no output and exit zero.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "internal/baseline/assets/setups/typescript-bun.json" && path != "internal/baseline/testdata/catalog.normalized.json" && path != "internal/baseline/testdata/catalog.digest" && path != "internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json" && path != "internal/baseline/testdata/parity-corpus/v1/manifest.json" && path != "docs/specs/0036-doctor-skill-readiness/task_03.md") {print; bad=1}} END {exit bad}'`
  — expected: no changed path outside the authorized Skill pair, derived
  Baseline and parity artifacts, and this Task file.
- `rtk make skills-sync-check` — expected: every canonical/generated owned
  skill pair has no drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: every
  shipped Roundfix Skill contract passes.
- `rtk git diff --check` — expected: no whitespace errors.
- `rtk make verify` — expected: formatting, tests, skill synchronization,
  shipped skill validation, and build all pass.

## References

- `_prd.md` → Goals; Core Feature 6; User Experience; Decisions.
- `_techspec.md` → Documentation and skill synchronization; Build Order 4.
- `task_02.md` → completed public guidance this Skill must match.
- `docs/agents/spec-routing.md` → tooling authorization and changed-file
  postflight.

## Result

Pending.
