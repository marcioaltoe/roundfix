---
source: coderabbit
pr: "56"
round: 1
round_created_at: "2026-07-31T14:02:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: CONTEXT.md
line: 56
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8ye,comment:PRRC_kwDOS0qyts7b_18a
review_hash: 3a890fd2315e2cf269549c2d488d8dc457bfb3737b1a3a1e34e15233736c8af7
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 008: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Use `Spec` instead of `feature` in the artifact-set definition.**

The workflow creates Specs for features, refactors, and bug fixes. Calling this “one feature's planning artifact set” excludes valid Specs. Use “one Spec's planning artifact set” to match the workflow contract.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@CONTEXT.md` at line 56, Update the artifact-set definition in CONTEXT.md to
refer to “one Spec's planning artifact set” instead of “one feature's planning
artifact set,” preserving the existing list of artifacts and provenance
requirements.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:12c789dbaa880f49a9b68a38 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `Spec` is the glossary term for the artifact set and includes feature,
  refactor, and bug-fix routes; `feature` was narrower than the workflow.
- Fix: the definition now reads "One Spec's planning artifact set" while
  preserving its artifact and provenance list.
- Focused evidence: direct inspection of `CONTEXT.md`; `rtk git diff --check`
  passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
