---
source: coderabbit
pr: "29"
round: 1
round_created_at: "2026-07-16T20:45:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/setup-context-driven-validator
head_sha: 49cdc07dcdf5b8fcb40eb459f27383b00995c0e3
file: skills/setup-context-driven/scripts/context_setup.py
line: 1549
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYfg,comment:PRRC_kwDOS0qyts7WgS5C
review_hash: 7594d76e2b3beca41986854eb773c14c47260e9470b5a557648c6fa95eac1d01
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:20Z"
---

# Issue 006: _ Security & Privacy_ _ Major_ _ Quick win_

## Review Comment

_🔒 Security & Privacy_ | _🟠 Major_ | _⚡ Quick win_

**Reject manifest artifact paths outside the repository before any I/O.**

Absolute and `..` paths from `managedArtifacts` are accepted. Stale marked artifacts can then flow through `remove_obsolete_artifacts()` into `apply_change_plan()`, rewriting or deleting files outside `--repo`.

Require safe relative paths and enforce resolved repository containment at the write boundary.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/scripts/context_setup.py` around lines 1542 -
1549, Update manifest_artifact_paths() to reject absolute paths and paths
containing traversal components from managedArtifacts, and validate each
candidate against the resolved repository root before any filesystem I/O.
Preserve only safe relative artifact paths, then enforce resolved repository
containment at the write boundary in
remove_obsolete_artifacts()/apply_change_plan() so stale artifacts cannot modify
files outside --repo.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5670759feee031922b22e284 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Manifest artifact paths were trusted before file access. Added safe relative path validation, repository containment checks, and write-boundary enforcement.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
