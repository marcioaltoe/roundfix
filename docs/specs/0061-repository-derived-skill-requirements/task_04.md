---
task: task_04
spec: 0061-repository-derived-skill-requirements
status: pending
type: docs
complexity: low
---

# Task 04: Align the Roundfix Skill pair and the derived digests

## Overview

The shipped Roundfix Skill describes Repository Skill Set readiness as a fixed
required set. Align it with the derivation this Spec implements and propagate
the deterministic digest fallout with the sanctioned command.

## Requirements

1. MUST describe the required external set as derived from the repository's
   Setup Manifest and its selected modules, and the owned set as coming from
   the running binary.
2. MUST describe the absent-manifest outcome and the per-skill remediation
   command, replacing any package-wide install guidance.
3. MUST keep both Skill copies byte-identical.
4. MUST regenerate every derived digest with `make baseline-digests`; no pin
   may be hand-edited.
5. MUST change only the expressly authorized paths plus this Task file.

## Subtasks

- [ ] Rewrite the readiness section of the canonical Skill.
- [ ] Synchronize the embedded copy.
- [ ] Regenerate the derived digests with the sanctioned command.

## Acceptance Criteria

- [ ] The Skill states that the external requirement follows the repository's
      selected modules and no longer implies a fixed list.
- [ ] The Skill names the per-skill install remediation and the
      absent-manifest outcome.
- [ ] Both copies are byte-identical and the embedded catalog validates.
- [ ] Every derived pin was produced by the sanctioned command.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: the Skill pair is synchronized.
- `go test -count=1 ./internal/baseline/ ./skills/` — expected: pass; the embedded catalog and every derived pin validate.
- `grep -q 'Setup Manifest' skills/roundfix/SKILL.md` — expected: the readiness contract names its source.

## References

`_prd.md` → Core Features 1, 3, 4, Project Constraints: Tooling authority;
`_techspec.md` → Build Order 4.
