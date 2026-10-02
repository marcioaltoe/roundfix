---
type: fix
status: open
created: 2026-10-02
spec: null
reason: null
---

# Restoring trailing skills breaks the pinned skill digests

## Symptom

Spec 0215's Doctor reports `DR-SKILL-TRAILS-SNAPSHOT` in this repository: 11 installed external skills differ from the embedded Setup Snapshot (`b3c45a4`): `bubbletea`, `domain-modeling`, `golang-concurrency` and others. `roundfix baseline update --repo . --yes` restores them (`Skills restored: 11`, the second run reports `current`, and the Doctor reports `skills: ok`). Then `make verify` fails in three tests that pin the upstream-managed skill tree digest:

```text
baseline_skill_contract_test.go:1264: upstream-managed skill tree digest = "27bc59d6…", want "8832b7ac…"
```

`TestUpstreamADRFormatUnchanged`, `TestAuthoringConstraintOwnership` and `TestAuthorialSkillSync` fail. The restore was not merged (PR #350, closed 2026-10-02), so the Doctor keeps warning.

## Where

`skills/baseline_skill_contract_test.go` (a Governed Path), which pins the digest; the installed `.agents/skills/` trees and `skills-lock.json`.

## Expected

Restoring the installed skills to the Setup Snapshot and re-pinning the digest land together, under a grant for the governed test, and `make verify` passes afterwards. Better still, the pin is derived from the Setup Snapshot itself, so a supported restore never needs a hand-edited constant.
