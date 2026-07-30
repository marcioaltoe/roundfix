---
source: coderabbit
pr: "52"
round: 2
round_created_at: "2026-07-30T19:37:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/loop-qa-ordering
head_sha: fc884f849005781a3e88472fbbabc6a249f7083d
file: internal/baseline/assets/modules/autonomous-work.json
line: 93
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6VOND-,comment:PRRC_kwDOS0qyts7br3BB
review_hash: 2d6aaae7c0da76f7f0aaa7da799bd2751cbbf62eb2586a9a0a387bf7750f0743
duplicate_of: ""
source_review_id: "4822588027"
source_review_submitted_at: "2026-07-30T19:37:10Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Require QA to pass before merging.**

The prescribed sequence ends with “request QA once, merge” without explicitly waiting for a terminal passing result. State that merging is allowed only after the QA gate returns `pass`; stop on failure.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/assets/modules/autonomous-work.json` at line 93, The
guidance in the autonomous workflow must explicitly require a terminal QA result
before merging: update the sequence around “request QA once, merge” to wait for
the QA gate, permit merge only when it returns pass, and stop the workflow on
failure.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:99ee036ae8057123209b2bfc -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - The finding applies. The autonomous-work clause ordered the workflow to
    request QA and then merge without requiring the terminal QA verdict, so a
    failing or incomplete gate did not explicitly stop delivery.
  - Updated the canonical Baseline module and this repository's autonomous-work
    guide to wait for the terminal QA verdict, permit merge only on `pass`, and
    stop the workflow on every non-`pass` verdict.
  - Bumped the `autonomous-work` module and supporting-guide versions from 6
    to 7, then regenerated the deterministic formatter fixture and catalog
    pins with `rtk make baseline-digests`.
  - Focused evidence:
    - The pre-fix `rtk jq -e` assertion for both `merge only` and
      `non-\`pass\`` returned `false` with exit 1; the same assertion after the
      edit returned `true` with exit 0.
    - `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/baseline -run 'Test(FormatterComposition|BaselineCompatibilityCorpus|CatalogCompatibility)$' -count=1`
      passed.
    - A second `rtk make baseline-digests` passed and reported
      `changed:false`, proving the generated artifacts match their canonical
      sources.
  - The first focused Go test attempt used the inaccessible host Go cache and
    failed with `open /Users/marcio/Library/Caches/go-build/...: operation not
    permitted`; the identical focused suite passed with the repository-local
    writable cache. The Daemon owns the unchanged authoritative `make verify`
    run after this Agent turn.
