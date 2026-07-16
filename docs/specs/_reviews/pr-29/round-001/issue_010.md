---
source: coderabbit
pr: "29"
round: 1
round_created_at: "2026-07-16T20:45:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/setup-context-driven-validator
head_sha: 49cdc07dcdf5b8fcb40eb459f27383b00995c0e3
file: skills/setup-context-driven/tests/test_apply.py
line: 127
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYfw,comment:PRRC_kwDOS0qyts7WgS5W
review_hash: eaed22e122d34cf1f06d42aac609a89a214d46963266a4765cfb7e4d8ce8a053
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:20Z"
---

# Issue 010: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Replace permission-based failure simulation with deterministic fault injection.**

Lines 118-123 rely on `chmod(0o500)`, but privileged Unix users can still write and Windows does not enforce these permissions equivalently. The test can therefore return success without exercising `managed.apply.failed`, making verification environment-dependent. Inject an `OSError` into the in-process write/replace path instead and assert rollback there.

As per coding guidelines, “Never use workarounds in production code or tests; fix the root cause.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/tests/test_apply.py` around lines 111 - 127,
Replace the chmod-based setup in
test_failure_before_commit_preserves_target_files with deterministic fault
injection by mocking the in-process write or replace operation to raise OSError
during run_apply. Keep the temporary repository and snapshot assertions, then
verify managed.apply.failed and unchanged target files after the injected
failure, restoring the mock afterward.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5fcb920e9f13f11155ea08b7 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The test depended on platform-specific directory permissions. Replaced it with deterministic in-process `Path.replace` fault injection and kept rollback assertions.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
