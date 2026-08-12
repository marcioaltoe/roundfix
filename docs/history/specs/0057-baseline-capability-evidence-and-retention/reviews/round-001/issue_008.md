---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/plan_characterization_test.go
line: 163
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymc,comment:PRRC_kwDOS0qyts7cjgE4
review_hash: faab285f1fa12f7b19015f1ba6e7176a85b21dd2cc64572de06dff12a600e4a4
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 008: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Do not run a nested `go test` from a test.**

`TestBaselinePlanCharacterizationPublicCommandJourneys` shells out to `go test ./internal/cli`. This creates several problems. The test requires a Go toolchain and a warm module cache at test time. It ignores `-run`, `-race`, and coverage flags of the outer run. It duplicates work that CI already performs when it runs the whole module. It also couples package `baseline` to package `internal/cli` test selection by name, so a rename of `TestBaselineMacroJourneysPublicCLI` breaks this test silently through a non-obvious failure message.

Let CI run `./internal/cli` directly, or move this cross-package journey assertion behind an integration build tag.

As per coding guidelines: "Use build tags (`//go:build integration`) to separate integration tests from unit tests."

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/plan_characterization_test.go` around lines 145 - 163,
Remove TestBaselinePlanCharacterizationPublicCommandJourneys and its nested
exec.Command go test invocation from the unit test suite. Keep cross-package CLI
journey coverage in CI or relocate it behind an integration build tag using the
existing TestBaselineMacroJourneysPublicCLI coverage, rather than invoking go
test from a test.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:d99759c4cead50dac3b8c57d -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Removed the Baseline-package test that spawned a nested `go test`. The owning CLI journey passed directly with `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/cli -run '^TestBaselineMacroJourneysPublicCLI$' -count=1`.
