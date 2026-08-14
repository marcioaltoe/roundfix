---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: invalid
terminal_reason: "The test-only PRAGMA value is constant-derived, and repository policy forbids adding lint or security suppressions merely to silence a false positive."
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/store/store_test.go
line: 2554
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIp2,comment:PRRC_kwDOS0qyts7aUVEK
review_hash: 9cac6a403709d5b752768f16b3e43a60ab4afde67f8423aee701af1acb0ddee3
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 013: _ Security & Privacy_ _ Trivial_ _ Low value_

## Review Comment

_🔒 Security & Privacy_ | _🔵 Trivial_ | _💤 Low value_

**Annotate the flagged `PRAGMA` interpolation so the scanners stay quiet.**

Both OpenGrep (`coderabbit.sql-injection.go-query-format`) and ast-grep (`sql-injection-exec-sprintf-go`) flag this line. It is a false positive — `PRAGMA user_version` cannot take a bind parameter and `newerVersion` is derived from the `schemaVersion` constant — but an unannotated finding is indistinguishable from a real one in CI output.

<details>
<summary>🔧 Suggested annotation</summary>

```diff
-	newerVersion := schemaVersion + 1
-	if _, err := db.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, newerVersion)); err != nil {
+	newerVersion := schemaVersion + 1
+	// PRAGMA does not accept bind parameters and newerVersion derives from the
+	// schemaVersion constant, so no untrusted input reaches the statement.
+	//nolint:gosec // constant-derived PRAGMA value; parameters are not supported here
+	if _, err := db.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, newerVersion)); err != nil {
```
</details>

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 ast-grep (0.44.1)</summary>

[warning] 2550-2550: Detected a SQL statement built with 'fmt.Sprintf' and passed directly to 'db.Exec'/'db.ExecContext'. Interpolating values into a query string lets an attacker inject arbitrary SQL. Use parameterized queries instead: pass the SQL with placeholders ('?' or '') as the query argument and supply the values as separate arguments, e.g. 'db.Exec("UPDATE t SET x = ? WHERE id = ?", x, id)'.
Context: db.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, newerVersion))
Note: [CWE-89] Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection').

(sql-injection-exec-sprintf-go)

</details>
<details>
<summary>🪛 OpenGrep (1.25.0)</summary>

[ERROR] 2551-2551: SQL query built via fmt.Sprintf or string concatenation passed to a database method. Use parameterized queries with placeholder arguments.


(coderabbit.sql-injection.go-query-format)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/store_test.go` around lines 2550 - 2554, Annotate the
fmt.Sprintf PRAGMA statement in the schema-version fixture test around
newerVersion with the repository’s recognized suppression comment for both
OpenGrep coderabbit.sql-injection.go-query-format and ast-grep
sql-injection-exec-sprintf-go. Keep the existing constant-derived value and
error handling unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:bc06f7e364d00ae8bc3b8e34 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The interpolated value is `schemaVersion + 1`, so no untrusted data reaches SQL. Adding the requested suppression would conflict with the repository rule against lint/security suppressions used only to silence warnings; no production or test behavior requires a change.
