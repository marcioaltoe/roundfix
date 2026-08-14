---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/adr/0097-a-qa-row-carries-forward-only-on-declared-unmoved-evidence.md
line: 16
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5X_,comment:PRRC_kwDOS0qyts7eY0ii
review_hash: 649dc3ea6d7014bc32386775c1071080e0ece9f99742356857f969ade8d9af7a
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Make non-repository evidence explicit and non-carriable.**

Lines [3]-[8] make carry-forward depend only on repository paths and ancestry. Lines [15]-[16] assume that an externally dependent row cannot also declare repository inputs, but the contract does not enforce that assumption. A row could cite unchanged files and still depend on elapsed time or a live service, causing a stale pass to carry.

Add a typed evidence-input declaration or an explicit `reobserve_required` flag. Reject carry-forward when any non-repository input exists.

As per coding guidelines, acceptance conditions must be observable and checkable.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/adr/0097-a-qa-row-carries-forward-only-on-declared-unmoved-evidence.md`
around lines 3 - 16, Add an explicit, typed evidence-input declaration or
reobserve_required marker to the QA row contract, and update the carry-forward
eligibility rule so any non-repository dependency—such as external repositories,
live services, or elapsed time—forces re-observation even when repository paths
are unchanged. Make the acceptance conditions observable and checkable by
documenting or validating that mixed repository/non-repository evidence is never
carried forward.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:33fde5917de954388f313c9c -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The original rule could not represent mixed evidence. ADR-0097 and
  Spec 0080 now type every input as `repository_path`, `external_repository`,
  `live_service`, or `elapsed_time`; carry-forward requires a non-empty list
  containing only repository paths. Any mixed or non-repository input forces
  re-observation, and the adversarial test contract names that refusal.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test -count=1 -parallel=1 ./internal/speccheck -run
  '^TestCheckCorpusBudget$'` passed; `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go run -buildvcs=false
  ./cmd/roundfix spec check` passed with no findings for Specs 0080 and 0081.
- Daemon Verification: `make verify` not run; Daemon-owned.
