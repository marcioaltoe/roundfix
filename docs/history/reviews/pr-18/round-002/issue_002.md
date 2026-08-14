---
source: coderabbit
pr: "18"
round: 2
round_created_at: "2026-07-07T14:05:56Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 4237143afdd7097e755e14b962156aaf6c6e6654
file: internal/cli/cli_test.go
line: 4548
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6O6PJ2,comment:PRRC_kwDOS0qyts7SyYRm
review_hash: f56b11b2a1c3c598a79da1ce17b1d6db0e36dbdb0a140293b59b452e8d7e49a4
duplicate_of: ""
source_review_id: "4645087962"
source_review_submitted_at: "2026-07-07T12:31:07Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Use `ExecContext` instead of `Exec` for the raw SQLite update.**

golangci-lint's `noctx` flags this call; a `context.Context` is already used elsewhere in this same test file (`context.Background()`).





<details>
<summary>🐛 Proposed fix</summary>

```diff
 func setListedRunCreatedAt(t *testing.T, homeDir string, runID string, createdAt time.Time) {
 	t.Helper()
+	ctx := context.Background()
 	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
 	if err != nil {
 		t.Fatalf("open raw Run Database: %v", err)
 	}
 	defer func() {
 		if err := db.Close(); err != nil {
 			t.Fatalf("close raw Run Database: %v", err)
 		}
 	}()
-	result, err := db.Exec(
+	result, err := db.ExecContext(
+		ctx,
 		`UPDATE runs SET created_at = ?, updated_at = ? WHERE id = ?`,
 		createdAt.UTC().Format(time.RFC3339Nano),
 		createdAt.UTC().Format(time.RFC3339Nano),
 		runID,
 	)
```
</details>

As per static analysis, "(*database/sql.DB).Exec must not be called. use (*database/sql.DB).ExecContext".

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func setListedRunCreatedAt(t *testing.T, homeDir string, runID string, createdAt time.Time) {
	t.Helper()
	ctx := context.Background()
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open raw Run Database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close raw Run Database: %v", err)
		}
	}()
	result, err := db.ExecContext(
		ctx,
		`UPDATE runs SET created_at = ?, updated_at = ? WHERE id = ?`,
		createdAt.UTC().Format(time.RFC3339Nano),
		createdAt.UTC().Format(time.RFC3339Nano),
		runID,
	)
	if err != nil {
		t.Fatalf("set listed Run creation time: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("read listed Run timestamp update result: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected to update one listed Run timestamp, updated %d", affected)
	}
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 golangci-lint (2.12.2)</summary>

[error] 4294-4294: (*database/sql.DB).Exec must not be called. use (*database/sql.DB).ExecContext

(noctx)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 4283 - 4310, The raw SQLite update in
setListedRunCreatedAt uses DB.Exec without a context, which triggers the noctx
lint warning. Update this helper to use ExecContext with an appropriate context
value already available in this test file, keeping the existing SQL and
timestamp parameters unchanged. Make the change in setListedRunCreatedAt so the
raw database update follows the same context-aware pattern used elsewhere in
cli_test.go.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3dd845fdc200a41ccde24203 -->

_Sources: Coding guidelines, Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - `setListedRunCreatedAt` now uses `db.ExecContext` with the test context for the raw SQLite timestamp update.
  - Verification: configured command `make verify` was run as `rtk make verify` and passed: Go tests, skills check, and build completed.
