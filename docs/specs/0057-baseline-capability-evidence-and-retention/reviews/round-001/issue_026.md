---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/cli/baseline_profile.go
line: 231
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ym7,comment:PRRC_kwDOS0qyts7cjgFc
review_hash: 90acaf190f627bf534c1f4cf5bf7d04b0014cf96f29abd47daf60d9a25605a41
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 026: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Match the single-dash `-format` form so failures stay machine-readable.**

Go's `flag` package accepts `-format=json` and `-format json` exactly like the double-dash forms. `baselineCapabilitiesCheckJSONRequested` matches only `--format`. When parsing fails, `request.format` is the zero value, so this helper is the only source of JSON intent.

The result is inconsistent error output for equivalent invocations:

- `baseline capabilities check --format=json --bogus` writes the structured failure body to stdout.
- `baseline capabilities check -format=json --bogus` writes nothing to stdout.

An agent that requested JSON receives no machine-readable error in the second case.





<details>
<summary>🐛 Proposed fix</summary>

```diff
 func baselineCapabilitiesCheckJSONRequested(args []string) bool {
 	for index, arg := range args {
-		if arg == "--format" && index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
+		if (arg == "--format" || arg == "-format") &&
+			index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
 			return true
 		}
-		if strings.HasPrefix(arg, "--format=") && strings.TrimSpace(strings.TrimPrefix(arg, "--format=")) == "json" {
-			return true
+		for _, prefix := range []string{"--format=", "-format="} {
+			if value, found := strings.CutPrefix(arg, prefix); found &&
+				strings.TrimSpace(value) == "json" {
+				return true
+			}
 		}
 	}
 	return false
 }
```
</details>

As per coding guidelines: "AP-O2: Use a consistent JSON response schema for success, empty-result, and error cases."

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func baselineCapabilitiesCheckJSONRequested(args []string) bool {
	for index, arg := range args {
		if (arg == "--format" || arg == "-format") &&
			index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
			return true
		}
		for _, prefix := range []string{"--format=", "-format="} {
			if value, found := strings.CutPrefix(arg, prefix); found &&
				strings.TrimSpace(value) == "json" {
				return true
			}
		}
	}
	return false
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_profile.go` around lines 221 - 231, Update
baselineCapabilitiesCheckJSONRequested to recognize both single-dash and
double-dash format arguments, including separate-value and equals forms such as
-format json and -format=json. Preserve the existing whitespace trimming and
JSON-value matching so parsing failures continue producing the structured JSON
error response.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:257f6144b52e6752e436ab47 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: JSON intent detection now recognizes -format json and -format=json as well as both double-dash forms. A table-driven public-command regression proves every form returns the structured preflight failure body. `go test ./internal/baseline ./internal/cli` passed with a writable GOCACHE.
