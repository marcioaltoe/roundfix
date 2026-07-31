---
source: coderabbit
pr: "56"
round: 3
round_created_at: "2026-07-31T14:59:09Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/write-prd/SKILL.md
line: 138
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vc6Tt,comment:PRRC_kwDOS0qyts7cBOhn
review_hash: d1a058b00d35f77653aa9d70fae12b7662dbb3e7fc9b21f6f65dc737198245bb
duplicate_of: ""
source_review_id: "4829633138"
source_review_submitted_at: "2026-07-31T14:58:04Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Exclude archived Specs from automatic link rewrites.**

“Search the whole repository” includes `docs/specs/_archived/`. Rewriting those files changes historical artifacts during another Spec’s adoption. Restrict automatic rewrites to mutable paths and report archived links for an explicit policy.

Based on learnings, archived Specs must not be edited; unarchive only through an explicit directory move and metadata change.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-prd/SKILL.md around lines 131 - 138, Update step 8's
link-rewrite procedure to exclude the `docs/specs/_archived/` directory from the
automatic search and rewrite scope. Instead of rewriting links within archived
Specs during adoption, collect and report those links separately for explicit
policy review, preserving the byte-for-byte integrity of historical artifacts.
Keep the existing rewrite flow unchanged for all non-archived repository paths
and for link destinations within the adopted source itself.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8f1ffc122eeebb1ce7963b10 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Step 8 required a whole-repository automatic rewrite without excluding
  `docs/specs/_archived/`, which conflicts with the repository rule that
  completed and archived legacy Specs remain byte-identical.

## Resolution

- Excluded `docs/specs/_archived/` from automatic link rewrites and identified
  archived Specs as immutable historical artifacts.
- Required links found in archived Specs to be reported separately for explicit
  policy review while preserving the existing rewrite and resolution gate for
  every non-archived linking file, including adopted sources.
- Added focused Skill contract assertions for both the exclusion and separate
  reporting requirements.

## Focused evidence

- Before the fix,
  `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -run TestSpecReferenceLifecycleSkillContracts -count=1`
  exited 1 and reported the missing archived-Spec exclusion and reporting
  contracts. The same command exited 0 after the fix.
- `rtk make skills-sync-check` — exit 0.
- `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -count=1`
  — exit 0.
- `rtk make baseline-digests` — exit 0; the required repeat run reported
  `ok:true, changed:false`.
- `make verify` was not run; authoritative Verification is Daemon-owned.
