---
source: coderabbit
pr: "137"
round: 1
round_created_at: "2026-08-07T03:22:12Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/promoted-backlog-entries-have-a-home
head_sha: ea93c68b70d066c1ee7f322e40ac1d547420e8be
file: docs/findings/2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md
line: 100
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XKdmq,comment:PRRC_kwDOS0qyts7eg8KI
review_hash: f0d1c170a72e538ba21785a0fc28d76d8654aa8a3fa07b430ef25d0b8782d0f8
duplicate_of: ""
source_review_id: "4879615443"
source_review_submitted_at: "2026-08-07T03:12:07Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Align the detector description with its fix messages.**

In `internal/speccheck/backlog.go` (Lines [40-98]), the missing-Spec branch does not report a concrete destination. The unresolvable-Spec branch also asks for a corrected slug before a destination is known. Limit the claim to resolvable Specs, or say that the fix provides destination guidance.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/findings/2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md`
around lines 94 - 100, Update the detector description around SC-BACKLOG-UNMOVED
to match the actual fix messages: either limit the claim about reporting the
destination to resolvable Specs, or revise it to state that the fix provides
destination guidance. Ensure the missing-Spec and unresolvable-Spec cases are
not described as reporting a concrete destination before one is known.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:f340d4938488668d05b6f931 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The addendum said every fix line reports a destination, but the
  missing-Spec and unresolvable-Spec branches can only explain how to establish
  one. The description now limits the concrete-destination claim to resolvable
  Specs and describes the other branches as destination guidance.

## Resolution

- Changed `docs/findings/2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md`
  to match the three detector branches without overstating their output.
- Focused evidence: `rtk rg -n "For a resolvable Spec|Missing and unresolvable Spec cases" docs/findings/2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md`
  exited 0 and located the corrected description at lines 98-99.
- Scope evidence: `rtk git diff --check` exited 0.
- Authoritative Verification `make verify` was not run; the Daemon owns it
  after this Agent turn.
