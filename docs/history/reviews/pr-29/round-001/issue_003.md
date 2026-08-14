---
source: coderabbit
pr: "29"
round: 1
round_created_at: "2026-07-16T20:45:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/setup-context-driven-validator
head_sha: 49cdc07dcdf5b8fcb40eb459f27383b00995c0e3
file: skills/setup-context-driven/scripts/context_assets.py
line: 238
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYfb,comment:PRRC_kwDOS0qyts7WgS45
review_hash: 3f1c8b2713b9c555abdf61b71ee530c70b2d2073ff5d1ccbe2e2f52b45b368da
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:20Z"
---

# Issue 003: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Return diagnostics for malformed collection shapes instead of crashing.**

`_read_collection()` calls `.get()` on any parsed JSON value, while `_index_assets()` calls `.get()` on unvalidated collection items. Arrays, `null`, or scalar entries therefore raise `AttributeError` instead of `AssetValidationError` and controlled exit code 2.

Validate object, list, and item types before dereferencing them.






Also applies to: 249-262

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/scripts/context_assets.py` around lines 219 -
238, Update _read_collection() and _index_assets() to validate parsed collection
data and each collection item before calling .get(). Treat arrays, null,
scalars, and other non-object shapes as validation failures by recording
diagnostics and raising AssetValidationError, preserving the controlled exit
code 2 instead of allowing AttributeError to escape.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:79a09ccb9381396aebac96f1 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Malformed collection documents and collection items could reach `.get()` calls. Added object/list/item shape checks that report diagnostics instead of raising `AttributeError`.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
