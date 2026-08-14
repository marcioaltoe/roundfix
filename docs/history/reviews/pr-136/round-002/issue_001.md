---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/adr/0096-the-qa-gate-proves-machine-facts-before-it-spends-an-agent-turn.md
line: 20
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5X7,comment:PRRC_kwDOS0qyts7eY0ie
review_hash: 461fff9ade6f2a91e83d7918ea6f31b17b84dd878b4371d480b567d110dde2d9
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 001: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Correct the measurement denominator.**

Lines [18]-[20] say that five rounds each found one defect, then say that two of the three were known facts. State whether the denominator is five defects, three distinct defects, or three measured facts.

As per coding guidelines, technical documents must use concrete, source-backed numbers.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/adr/0096-the-qa-gate-proves-machine-facts-before-it-spends-an-agent-turn.md`
around lines 18 - 20, Clarify the measurement denominator in the sentence
beginning “Agent turn. Measured on Spec 0079” by explicitly identifying whether
the comparison uses five total findings, three distinct defects, or three
measured facts. Make the reported counts internally consistent and
source-backed, without implying a different denominator.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:594a9b54aa964023b4dd486e -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Replaced the unsupported five-round denominator with the measured
  first three Runs: 92 minutes for the initial seven Tasks plus gate, then 29
  and 30 minutes for the corrective Runs. The sentence now states that those
  Runs exposed three distinct defects and that two were machine-known facts,
  matching the promoted backlog entry's source measurements.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  run -buildvcs=false ./cmd/roundfix spec check` passed with no findings for
  Specs 0080 and 0081; `rtk git diff --check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
