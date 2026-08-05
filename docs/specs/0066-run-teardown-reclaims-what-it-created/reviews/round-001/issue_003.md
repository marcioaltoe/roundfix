---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: docs/specs/0066-run-teardown-reclaims-what-it-created/qa/evidence/2026-08-05-qa-01/_fixture/setup_reconcile/main.go
line: 127
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9Xt,comment:PRRC_kwDOS0qyts7dnSbL
review_hash: 09148bdf4fb613088c006933df2fb059cddc228fb6fe3bbbf0bbc3416a4dae6b
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 003: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Do not parse `CombinedOutput` as a Git object ID.**

`git` merges stderr into its return value. Line 37 parses that value as a HEAD SHA. Any Git advisory, hint, or warning written to stderr becomes part of the string, and `strings.TrimSpace` does not remove it. The invalid SHA then flows into `store.CreateRun` and `runworktree.Create`, and the fixture fails in a confusing place. Capture stdout separately for value-returning calls.



<details>
<summary>🔧 Proposed fix</summary>

```diff
 func git(dir string, args ...string) string {
 	command := exec.Command("git", args...)
 	command.Dir = dir
-	output, err := command.CombinedOutput()
-	if err != nil {
-		fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
-	}
-	return string(output)
+	var stdout, stderr bytes.Buffer
+	command.Stdout = &stdout
+	command.Stderr = &stderr
+	if err := command.Run(); err != nil {
+		fatalf("git %s: %v\n%s", strings.Join(args, " "), err, stderr.String())
+	}
+	return stdout.String()
 }
```

Add `"bytes"` to the import block.
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func git(dir string, args ...string) string {
	command := exec.Command("git", args...)
	command.Dir = dir
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		fatalf("git %s: %v\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String()
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/specs/0066-run-teardown-reclaims-what-it-created/qa/evidence/2026-08-05-qa-01/_fixture/setup_reconcile/main.go`
around lines 119 - 127, Update the git helper to capture stdout separately from
stderr for value-returning Git commands, rather than returning CombinedOutput
that may include advisory text. Add the bytes buffer import and configure the
command’s standard output and error streams accordingly, while preserving fatal
handling and diagnostic stderr output for command failures.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c6142a28c162f769c490d323 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Confirmed the fixture returned `CombinedOutput`, so harmless Git stderr advisories could contaminate values parsed as revisions.
  - The fixture Git helper now captures stdout and stderr separately, returns stdout only, and includes stderr only in failure diagnostics.
  - Focused evidence: the fixture package command completed successfully and `rtk make fmt-check` passed.
  - The Daemon owns authoritative `make verify` after this Agent turn.
