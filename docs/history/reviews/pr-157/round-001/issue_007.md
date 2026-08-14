---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: internal/cli/archive.go
line: 22
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdQ,comment:PRRC_kwDOS0qyts7gSdxx
review_hash: 08261abda7185ef029e2d4ec7917a4cc3550be0074c49c0e36421efd2fc6ce7a
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 007: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Make archive help describe external Spec roots.**

When `specs.root` is external, the command archives from the configured root, but this help text still says `docs/specs/<slug>` moves to `_archived/specs/<slug>`. State that these are default paths or describe the configured Spec root instead.

The PR objective explicitly includes support for external Spec roots.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/archive.go` around lines 15 - 22, Update the archiveUsage text
to avoid presenting docs/specs/<slug> and the archive destination as universal
paths; identify them as default paths or describe that the configured Spec root
is used, while preserving the existing archive behavior description.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:4d3fe4f0e1a4eabfc49e2320 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Updated `internal/cli/archive.go`'s `archiveUsage` help text so it no longer presents `docs/specs/<slug>` and `_archived/specs/<slug>` as universal paths. It now states that these are the destinations for the built-in `docs/specs` root and that a non-default or external Spec Root archives beside the active root at `<spec-root>/_archived/<slug>/`, matching the archive command's behavior. Existing behavior description (verifies Tasks and QA, stamps metadata, creates no Run, never pushes) is unchanged. Focused evidence: `TestRunArchiveHelp` still checks all required substrings and the full `make verify` gate passes.
