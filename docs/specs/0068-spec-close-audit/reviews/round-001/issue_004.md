---
source: coderabbit
pr: "113"
round: 1
round_created_at: "2026-08-05T02:12:07Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0068-implementation
head_sha: c9af2617f988bd63e1bd8f22c6178758a8e5fd40
file: internal/specaudit/audit_test.go
line: 263
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WhZqN,comment:PRRC_kwDOS0qyts7dk962
review_hash: a9f26862f5db831e00f7e168153907d54f53497cbcaf2c31a304b37fb4b5315e
duplicate_of: ""
source_review_id: "4860420451"
source_review_submitted_at: "2026-08-05T02:11:27Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Run the table cases as named subtests.**

The table already carries a `name` field, but the loop asserts inline with `t.Fatalf`. The first failing target stops the test, so the second target is never checked. The coding guidelines require every table case to pass its `name` to `t.Run`.





<details>
<summary>♻️ Proposed fix</summary>

```diff
 	} {
-		survivor := requireSurvivor(t, result, target.name, target.isWorktree)
-		if survivor.Kind == KindResidue {
-			t.Fatalf("Active Run survivor %q classified as residue", target.name)
-		}
-		if !strings.Contains(survivor.Evidence, "Active Run "+run.ID) {
-			t.Fatalf("Active Run survivor evidence = %q, want Run %s", survivor.Evidence, run.ID)
-		}
-		if survivor.Reclaim != "" {
-			t.Fatalf("Active Run survivor reclaim = %q, want empty", survivor.Reclaim)
-		}
+		t.Run(target.name, func(t *testing.T) {
+			survivor := requireSurvivor(t, result, target.name, target.isWorktree)
+			if survivor.Kind == KindResidue {
+				t.Fatalf("Active Run survivor %q classified as residue", target.name)
+			}
+			if !strings.Contains(survivor.Evidence, "Active Run "+run.ID) {
+				t.Fatalf("Active Run survivor evidence = %q, want Run %s", survivor.Evidence, run.ID)
+			}
+			if survivor.Reclaim != "" {
+				t.Fatalf("Active Run survivor reclaim = %q, want empty", survivor.Reclaim)
+			}
+		})
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	for _, target := range []struct {
		name       string
		isWorktree bool
	}{
		{name: branch},
		{name: worktreePath, isWorktree: true},
	} {
		t.Run(target.name, func(t *testing.T) {
			survivor := requireSurvivor(t, result, target.name, target.isWorktree)
			if survivor.Kind == KindResidue {
				t.Fatalf("Active Run survivor %q classified as residue", target.name)
			}
			if !strings.Contains(survivor.Evidence, "Active Run "+run.ID) {
				t.Fatalf("Active Run survivor evidence = %q, want Run %s", survivor.Evidence, run.ID)
			}
			if survivor.Reclaim != "" {
				t.Fatalf("Active Run survivor reclaim = %q, want empty", survivor.Reclaim)
			}
		})
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/specaudit/audit_test.go` around lines 246 - 263, Wrap each iteration
of the table-driven loop in a named t.Run using target.name, moving the
requireSurvivor call and all assertions inside the subtest so every target is
executed and reported independently.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:1c2c5f8a59f57d4780736755 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The table declared case names but ran assertions in the parent test, so the first fatal assertion could mask the second target.

## Resolution

- Wrapped each Active Run survivor case in `t.Run(target.name, ...)` with all case assertions inside the subtest.
- Focused evidence: `rtk env GOCACHE=/private/tmp/roundfix-review0068-audit-cache go test ./internal/specaudit -run '^(TestAuditClassifiesResidueBranch|TestAuditScratchWorktreeReclaimCommandRuns|TestAuditPreservesActiveRunSurvivors)$' -count=1` exited 0; the full affected-package command also exited 0.
- Daemon Verification: `make verify` was not run; the Daemon owns that command.
