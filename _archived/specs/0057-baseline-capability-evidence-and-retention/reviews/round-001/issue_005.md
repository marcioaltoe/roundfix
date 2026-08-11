---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/classification.go
line: 843
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0YmZ,comment:PRRC_kwDOS0qyts7cjgE0
review_hash: a3d0d54d7577437811bed71be19a533a02aedb439d1687d7420bbf03ce92c81b
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:29Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Do not discard the `os.OpenRoot` failure without a signal.**

If `os.OpenRoot` fails, `classifyCarriers` returns `nil` and every carrier stays unclassified. The behavior is fail-safe, because all nested-carrier warnings remain. The failure cause is lost, however: the function neither returns nor logs the error. A maintainer who sees unexpected warnings has no evidence of the root cause.

Return an `error` from `classifyCarriers` and propagate it, or record the failure as a `Finding` so the caller can report it.






As per coding guidelines: "Errors must be either logged or returned, never both." and "Always check returned errors; never discard them with `_`."

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/classification.go` around lines 840 - 843, Update
classifyCarriers to preserve the os.OpenRoot error instead of returning nil
silently: either change its signature and callers to return and propagate the
error, or convert the failure into a reportable Finding. Ensure the error is
handled exactly once—returned or logged, not both—and update every caller to
check the revised result.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:1c54b6618484c16d52c57ca8 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `classifyCarriers` now returns the `os.OpenRoot` failure with operation context, and all callers propagate it. A focused missing-root regression and the full Baseline package test passed.
