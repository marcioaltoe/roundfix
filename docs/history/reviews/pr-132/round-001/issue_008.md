---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: skills/baseline_skill_contract_integration_test.go
line: 48
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kad,comment:PRRC_kwDOS0qyts7eJylf
review_hash: cd877089cd96e99fc56626f62fa51ed88b84b9217f9161124b9562b709ac54dc
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:10Z"
---

# Issue 008: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Run artifact regeneration before the byte comparison.**

This test only edits files in a temporary copy. It does not run the snapshot regeneration path or another derived-artifact producer.

Run the regeneration workflow after the owned-skill edit, then compare the generated artifacts. The current test passes even if regeneration reintroduces an owned content digest.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/baseline_skill_contract_integration_test.go` around lines 16 - 48, The
test TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical must invoke the
repository’s snapshot or derived-artifact regeneration workflow after editing
the owned skill files and before the assertArtifactBytesEqual comparisons. Reuse
the existing regeneration helper or command used by the project, targeting
verificationRoot, so comparisons validate regenerated outputs and detect
reintroduced owned-content digests.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:b3c892b782934b8d2dcfeb04 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Comparing untouched derived files immediately after editing only the Skill files could not prove that a producer would keep those artifacts stable. The test now runs the real `make baseline-digests` workflow in the isolated tracked repository before comparing derived, characterization, and archived bytes. Focused evidence: `rtk go test ./skills -count=1 -run '^TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical$' -v` with the repository-local Go cache exited 0. Authoritative `make verify` remains Daemon-owned.
