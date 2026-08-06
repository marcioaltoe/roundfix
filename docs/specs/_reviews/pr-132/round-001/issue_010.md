---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: skills/baseline_skill_contract_test.go
line: 575
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kay,comment:PRRC_kwDOS0qyts7eJyl7
review_hash: f221925c093f0a75ebc15f743ae8ae11a7f467e729c3bc37a580aa86942f90b5
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:10Z"
---

# Issue 010: _ Functional Correctness_ _ Major_ _ Heavy lift_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _🏗️ Heavy lift_

**Validate owned digest fields, not current digest values.**

An old or arbitrary owned-skill digest passes this test after the skill content changes. The test only searches for the current `SkillFolderHash` result.

Decode each corpus record and reject `treeDigest` and `contentDigest` when the skill is Roundfix-owned. This verifies the invariant for both current and stale digest values.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/baseline_skill_contract_test.go` around lines 557 - 575, Update the
corpus validation around ownedDigests to decode each baseline skill record and
inspect its treeDigest and contentDigest fields instead of searching for current
hash strings. For every Roundfix-owned skill, reject any non-empty digest
values, including stale or arbitrary ones, while preserving the existing corpus
iteration and failure reporting.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7797f6b9238fb9c339f0fe36 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Searching for only the current folder hashes allowed stale or arbitrary owned-skill digests to survive. Corpus JSON and embedded base64 JSON content are now decoded recursively; any non-empty `treeDigest` or `contentDigest` on a Roundfix-owned Skill is reported by field and JSON path. A stale-digest fixture proves both fields are detected. Focused evidence: `rtk go test ./skills -count=1 -run '^TestCharacterizationCorporaDoNotRecordOwnedSkillDigests$'` with the repository-local Go cache exited 0. Authoritative `make verify` remains Daemon-owned.
