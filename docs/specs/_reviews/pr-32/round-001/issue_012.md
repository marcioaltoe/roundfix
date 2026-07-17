---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/store/agent_selection_test.go
line: 35
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tF,comment:PRRC_kwDOS0qyts7Wt95g
review_hash: 577e331497aaad33c2c4bb79f0a0df7c210a10143f35acc85848b2b962c1a25e
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-32/round-002/issue_012.md
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---


# Issue 012: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Use `ExecContext` in the migration fixture.**

`db.Exec` triggers the configured `noctx` check. Pass `ctx` into `buildV8Fixture` and execute each statement with `db.ExecContext(ctx, statement)`.

<details>
<summary>Proposed fix</summary>

```diff
-	buildV8Fixture(t, homeDir)
+	buildV8Fixture(t, ctx, homeDir)

-func buildV8Fixture(t *testing.T, homeDir string) {
+func buildV8Fixture(t *testing.T, ctx context.Context, homeDir string) {
...
-		if _, err := db.Exec(statement); err != nil {
+		if _, err := db.ExecContext(ctx, statement); err != nil {
```
</details>

As per coding guidelines, database operations must use context-first APIs.








Also applies to: 293-299, 365-368

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/agent_selection_test.go` at line 35, Update buildV8Fixture and
its call sites to accept a context parameter, then replace each db.Exec
invocation in the fixture with db.ExecContext(ctx, statement). Ensure the
migration test calls buildV8Fixture with the existing ctx and apply the same
change to the additional referenced statement blocks.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:9e696bb8ecd8136a282d2103 -->

_Sources: Coding guidelines, Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
