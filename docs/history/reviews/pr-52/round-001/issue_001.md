---
source: coderabbit
pr: "52"
round: 1
round_created_at: "2026-07-30T19:23:52Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/loop-qa-ordering
head_sha: fc884f849005781a3e88472fbbabc6a249f7083d
file: internal/baseline/assets/modules/autonomous-work.json
line: 93
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6VN9n2,comment:PRRC_kwDOS0qyts7brgTQ
review_hash: 07b40d74146337a5675c1df43d520e829475e1f4eb194b817629b28380a7ccbd
duplicate_of: ""
source_review_id: "4822478570"
source_review_submitted_at: "2026-07-30T19:23:02Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Scope PR-review requirements to specs that observe review surfaces.**

This guidance unconditionally requires opening a PR and reaching Clean review before QA, but the PR objectives state that specs whose acceptance does not observe review surfaces remain unaffected. Scope this requirement to those acceptance flows; otherwise, unrelated specs may be blocked unnecessarily.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/assets/modules/autonomous-work.json` at line 93, The
guidance value should require opening and reaching Clean Pull Request review
before QA only for Specs whose acceptance observes review surfaces. Update the
ordering and QA instructions in the guidance string to preserve the existing
flow for those Specs while leaving unrelated Specs unaffected.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e0b7567b84e1256b6c989838 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - The finding applies. The clause conditioned the general QA precondition on
    observed surfaces, but its Pull Request and Clean-review directives were
    unconditional, so a Spec with no review-surface acceptance could still be
    ordered to wait for review before QA.
  - Scoped the Pull Request creation, Clean review, and
    `implement -> Pull Request -> review -> QA -> merge` ordering to Specs
    whose acceptance observes Pull Request review surfaces. The clause now
    explicitly permits other Specs to request QA after the Task Graph closes
    without waiting for a Pull Request or review.
  - Bumped the `autonomous-work` module and supporting-guide versions from 5
    to 6, then regenerated the deterministic Baseline snapshots and catalog
    pins with `rtk make baseline-digests`.
  - Focused evidence:
    - `rtk proxy rg -n "When a Spec's acceptance observes Pull Request review surfaces|For those Specs, the order is|Specs whose acceptance does not observe Pull Request review surfaces" internal/baseline/assets/modules/autonomous-work.json internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
      exited 0 and found all three scope statements in the canonical module and
      generated guide.
    - `rtk proxy env GOCACHE=/private/tmp/roundfix-pr52-batch001-gocache go test ./internal/baseline -run 'Test(FormatterComposition|BaselineCompatibilityCorpus|CatalogCompatibility)$' -count=1`
      passed.
    - A second `rtk make baseline-digests` passed and reported
      `changed:false`, proving the generated artifacts are idempotent.
  - The first focused Go test attempt used the inaccessible host Go cache and
    failed with `operation not permitted`; the identical suite above passed
    with a writable focused-check cache. The Daemon owns the unchanged
    authoritative `make verify` run after this Agent turn.
