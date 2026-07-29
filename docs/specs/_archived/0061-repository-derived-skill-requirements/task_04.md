---
task: task_04
spec: 0061-repository-derived-skill-requirements
status: completed
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

## Result

Implementation:

- Rewrote Repository Skill Set readiness so the running binary owns the
  Roundfix-owned set while the repository's Setup Manifest and selected
  modules derive the external set.
- Documented the absent-or-unreadable Setup Manifest outcome: zero external
  requirements, a failed `skills:` line, and `roundfix baseline` in the next
  action.
- Replaced package-wide external installation guidance with
  `bunx skills add marcioaltoe/skills@<skill>` for each named missing or
  outdated external skill.
- Applied the same edit to the canonical and embedded Skill copies.
- Ran `rtk make baseline-digests`; it passed and regenerated exactly:
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`.

Focused-check evidence by acceptance criterion:

1. External requirement follows selected modules: the pre-change marker probe
   reported `manifest=0 per_skill=0 package_wide=1`; after the edit, the
   two-file `rtk awk` contract assertion passed with
   `manifest=1 per_skill=1 package_wide=0` for each Skill copy.
2. Per-skill remediation and absent-manifest outcome: the same post-edit
   assertion proved both copies name the Setup Manifest and
   `bunx skills add marcioaltoe/skills@<skill>`, with no
   `bunx skills experimental_install` guidance.
3. Byte identity and embedded catalog validity:
   `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
   passed;
   `rtk go test -count=1 ./skills -run '^TestAuthorialSkillSync$'` passed
   (18 tests);
   `rtk go test -count=1 ./internal/baseline -run
   '^(TestCatalogCompatibility|TestBaselineCompatibilityCorpus)$'` passed
   (2 tests).
4. Sanctioned pin production: the first `rtk make baseline-digests` run
   reported `ok:true`, `changed:true`, and the five generated paths above; an
   immediate second run passed with `ok:true`, `changed:false`, proving the
   generated artifacts already matched their canonical sources.

Additional focused evidence:

- `rtk git diff --check` passed.
- Changed-file postflight contains only the two expressly authorized Skill
  paths, the five generator-reported derived-pin paths, and this Task file.

Daemon Verification was not run by the Agent. The commands under
`## Verification` remain for the Daemon.
