---
source: coderabbit
pr: "136"
round: 3
round_created_at: "2026-08-06T20:20:19Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: fba018672a8f31a3a4f59e6afd21d2c03c6a220f
file: docs/specs/0080-cheap-detectors-run-before-the-gate/_techspec.md
line: 180
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XGgXD,comment:PRRC_kwDOS0qyts7ebLMB
review_hash: f785b2a486fb64889fa96c05ae3386a399ccb9ee9d48cc9933fe56fb4767da8e
duplicate_of: ""
source_review_id: "4877969817"
source_review_submitted_at: "2026-08-06T20:19:25Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Make carry-forward evidence validation executable.**

`Carriable` promises to verify that each cited evidence path resolves and has unchanged content, but the interface receives only `prior`, `head`, and `changed []string`. It does not define how `repository_path` references, including the glob in Line 179, are expanded or compared with the establishing report. Add verified evidence digests or current-snapshot data to the contract, and cover missing, deleted, changed, and glob inputs.

As per coding guidelines, acceptance conditions must be observable and checkable.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0080-cheap-detectors-run-before-the-gate/_techspec.md` around
lines 160 - 180, Update the Carriable contract to accept verifiable evidence
state, such as establishing-report digests and current-snapshot data, rather
than only prior, head, and changed paths. Define how repository_path
references—including globs—are expanded and compared, and require carry-forward
to reject missing, deleted, or content-changed inputs while accepting unchanged
inputs. Specify observable acceptance coverage for missing, deleted, changed,
and glob references.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:42d7668011aaa03d27f3dd02 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The original `Carriable` interface accepted no content snapshot, so
  it could not execute its promise to expand evidence references and compare
  their bytes. The TechSpec now defines canonical establishing and current
  snapshots, tracked-tree literal and glob expansion, exact path-set and
  SHA-256 comparison, fail-closed malformed or empty inputs, and explicit
  missing, deleted, changed, and glob coverage.
- Focused evidence: the post-change `rtk rg` contract probe found
  `EvidenceSnapshot`, `evidence_snapshots`, missing-literal, glob-change, and
  unchanged-glob coverage. `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go test -count=1 -parallel=1
  ./internal/speccheck -run '^TestCheckCorpusBudget$'` passed. `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go run -buildvcs=false
  ./cmd/roundfix spec check` exited 0 with no findings for Specs 0080 and 0081;
  it separately reported the pre-existing missing Task Graph, reference-index,
  and Vocabulary Contract skips. `rtk git diff --check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
