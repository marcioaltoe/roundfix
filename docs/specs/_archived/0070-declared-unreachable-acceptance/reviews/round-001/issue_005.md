---
source: coderabbit
pr: "110"
round: 1
round_created_at: "2026-08-04T22:55:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0070-implementation
head_sha: a588c6ca3ab9d977284ba1f9e80a89b0e6336786
file: internal/spec/qa.go
line: 227
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WeYYQ,comment:PRRC_kwDOS0qyts7dggq7
review_hash: ce506823e42fc20a0faa9b3a92574cd3a380f5a2f14d99e9dbf8c1d34b84b8f3
duplicate_of: ""
source_review_id: "4859094834"
source_review_submitted_at: "2026-08-04T21:23:49Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Reject `pass` reports with declared-blocked rows.**

`readQAReport` accepts `verdict: pass` when `rows_blocked_declared` is positive. The QA policy prohibits `pass` for declared-blocked rows. Callers can therefore accept an invalid passing report.

- `internal/spec/qa.go#L219-L227`: reject `VerdictPass` when `RowsBlockedDeclared > 0`.
- `internal/spec/qa_test.go#L73-L78`: change this case to expect the `rows_blocked_declared` validation error.

<details>
<summary>Proposed fix</summary>

```diff
 case VerdictPass:
 	if report.RowsBlockedFinding > 0 {
 		return QAReport{}, QAReportError{
 			Path: path,
 			Err:  fmt.Errorf("rows_blocked_finding must be zero when verdict is %q", VerdictPass),
 		}
 	}
+	if report.RowsBlockedDeclared > 0 {
+		return QAReport{}, QAReportError{
+			Path: path,
+			Err:  fmt.Errorf("rows_blocked_declared must be zero when verdict is %q", VerdictPass),
+		}
+	}
 	return report, nil
```

```diff
- name:             "declared-blocked pass does not change finding validation",
+ name:             "declared-blocked pass is unreadable",
  verdict:          VerdictPass,
  extraFrontmatter: []string{"rows_blocked_declared: 3"},
- wantVerdict:      VerdictPass,
+ wantError:        "rows_blocked_declared must be zero when verdict is \"pass\"",
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	switch report.Verdict {
	case VerdictPass:
		if report.RowsBlockedFinding > 0 {
			return QAReport{}, QAReportError{
				Path: path,
				Err:  fmt.Errorf("rows_blocked_finding must be zero when verdict is %q", VerdictPass),
			}
		}
		if report.RowsBlockedDeclared > 0 {
			return QAReport{}, QAReportError{
				Path: path,
				Err:  fmt.Errorf("rows_blocked_declared must be zero when verdict is %q", VerdictPass),
			}
		}
		return report, nil
```

```suggestion
		{
			name:             "declared-blocked pass is unreadable",
			verdict:          VerdictPass,
			extraFrontmatter: []string{"rows_blocked_declared: 3"},
			wantError:        "rows_blocked_declared must be zero when verdict is \"pass\"",
		},
```

</details>

<!-- suggestion_end -->

<details>
<summary>📍 Affects 2 files</summary>

- `internal/spec/qa.go#L219-L227` (this comment)
- `internal/spec/qa_test.go#L73-L78`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/qa.go` around lines 219 - 227, Update the VerdictPass
validation in internal/spec/qa.go lines 219-227 to reject reports when
RowsBlockedDeclared is greater than zero, returning the corresponding
rows_blocked_declared validation error. Update the affected test case in
internal/spec/qa_test.go lines 73-78 to expect this validation error instead of
accepting the report.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/spec/qa.go</file>
<line_range>219-227</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/spec/qa_test.go</file>
<line_range>73-78</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a0b3beb9419df96dadb7227e -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `readQAReport` rejected finding-blocked `pass` reports but omitted the parallel declared-blocked rule required by the canonical QA policy. `VerdictPass` now returns `QAReportError` when `RowsBlockedDeclared` is positive, while the intentional environment-blocked `pass` behavior remains unchanged.
- Evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/spec -run 'Test(QAVerdictValidatesBlockedCounts|UnreachableRejectsMalformedDeclaration|ArchivedQAReportCorpusRemainsReadable|ArchivedPassCorpusRemainsArchiveEligible)$'` passed; `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/cli -run '^TestRunArchive'` passed.
