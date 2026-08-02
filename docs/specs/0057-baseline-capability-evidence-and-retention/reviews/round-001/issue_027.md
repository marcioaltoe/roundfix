---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/cli/cli.go
line: 4807
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ym-,comment:PRRC_kwDOS0qyts7cjgFg
review_hash: c49fc618fedbd6fc3a22009584866189c83bcb8fd5f53fc98386ca524ce78f74
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:31Z"
---

# Issue 027: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Document the output-failure exit code.**

`runBaselineCapabilitiesCheckCommand` returns `exitRunFailed` when `writeBaselineCapabilitiesCheckResult` fails, at `internal/cli/baseline_profile.go` lines 128-131. This help block lists only 0, 2, and 3.

The sibling help blocks document that code. `baseline apply` at line 4859 lists "1  apply, verification, output, rollback, or recovery failure", and `baseline skills restore` at line 4930 does the same. An agent that follows this documented contract cannot classify the missing code.





<details>
<summary>📝 Proposed fix</summary>

```diff
 Exit codes:
   0  capability evidence evaluated with no blocking divergence
+  1  output failure
   2  invalid arguments, repository failure, or no resolvable Baseline Profile
   3  capability evidence evaluated with a blocking divergence
```
</details>

As per coding guidelines: "AP-E1: Use differentiated exit codes so agents can distinguish usage, authentication, retryable, and other errors."

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
Exit codes:
  0  capability evidence evaluated with no blocking divergence
  1  output failure
  2  invalid arguments, repository failure, or no resolvable Baseline Profile
  3  capability evidence evaluated with a blocking divergence
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 4804 - 4807, Update the exit-code
documentation for runBaselineCapabilitiesCheckCommand to include exitRunFailed
for output failures from writeBaselineCapabilitiesCheckResult, using the same
documented code and wording convention as the sibling baseline apply and
baseline skills restore help blocks.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:21ebe83a4cab65fc760a52fa -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Baseline capability-check help now documents exit code 1 for output failures, and the public help-contract test requires it. `go test ./internal/baseline ./internal/cli` passed with a writable GOCACHE.
