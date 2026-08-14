---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/releaseplan/classifier.go
line: 55
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5s5,comment:PRRC_kwDOS0qyts7Wt95U
review_hash: fa30e5a1b0486b96f500c991a24c0edfa9f253268b0f368720ff03909d012dc4
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-32/round-002/issue_009.md
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---


# Issue 009: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Apply the maintenance-only boundary before aggregating commit impact.**

A `fix:` commit touching only `_test.go` files currently contributes patch impact because Line 49 wins before Line 51. That contradicts the documented no-release maintenance boundary and can produce an unnecessary release plan.

<details>
<summary>Proposed fix</summary>

```diff
 		evidence := ClassifyCommit(commit)
 		changes = append(changes, evidence)
-		automaticMinimum = MaxImpact(automaticMinimum, evidence.AutomaticImpact)
-		breaking = breaking || evidence.Breaking
 
 		switch {
-		case evidence.AutomaticImpact != ImpactNone:
+		case !evidence.CrossesMaintenanceOnlyBoundary:
+			sources[SourceMaintenanceOnly] = true
+		case evidence.AutomaticImpact != ImpactNone:
+			automaticMinimum = MaxImpact(automaticMinimum, evidence.AutomaticImpact)
+			breaking = breaking || evidence.Breaking
 			sources[SourceConventionalCommit] = true
-		case !evidence.CrossesMaintenanceOnlyBoundary:
-			sources[SourceMaintenanceOnly] = true
 		default:
 			blockingCommits = append(blockingCommits, evidence.CommitSHA)
 		}
```

Add regression cases for `fix:` and `feat!:` commits containing only maintenance paths.
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/releaseplan/classifier.go` around lines 42 - 55, Update the commit
classification loop around ClassifyCommit so CrossesMaintenanceOnlyBoundary is
evaluated before aggregating AutomaticImpact, preventing maintenance-only
commits from contributing release impact. Preserve blocking/source behavior and
add regression coverage for fix: and feat!: commits containing only maintenance
paths.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:624ecab0872fcccc1028d9d5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
