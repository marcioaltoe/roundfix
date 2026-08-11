---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/cli/cli.go
line: 2652
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpT,comment:PRRC_kwDOS0qyts7aUVDW
review_hash: 3f2ac224eae5b9f024f0c48457aae48b9a647c29d182caff32fec17a46d218f3
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:18Z"
---

# Issue 003: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**A git-proof invocation failure fails the Run after the artifact commit was already pushed.**

Every *structural* mismatch inside `inheritReviewArtifactEvidence` (wrong parent, wrong subject, non-artifact paths, unverified parent) returns `false, nil` and correctly falls back to polling the new head. But a failure of any proof *command* — `rev-parse`, `rev-list`, `show`, `diff` — returns an error here, which aborts the publication and drives the Run to `Failed`, even though `maybeRunFinalPush` already succeeded on line 2633. A momentary git failure therefore converts a healthy, pushed round into a failed one instead of degrading to the existing fallback.

Treat proof-command failures the same as proof mismatches: fall back to normal head polling and surface the git error as a warning.



<details>
<summary>🛡️ Proposed fix</summary>

```diff
 			evidence, inherited, err := inheritReviewArtifactEvidence(ctx, reviewArtifactEvidenceRequest{
 				...
 			})
 			if err != nil {
-				return watch.ArtifactPublication{}, err
+				fmt.Fprintf(ui.progress, "Review artifact Evidence inheritance unavailable; falling back to head polling: %v\n", err)
+				inherited = false
 			}
 			if !inherited {
 				evidence = reviewsource.Evidence{}
 			}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
			evidence, inherited, err := inheritReviewArtifactEvidence(ctx, reviewArtifactEvidenceRequest{
				Source:         req.source,
				PRNumber:       preflightResult.PullRequest.Number,
				BaseRepository: preflightResult.PullRequest.BaseRepository,
				HeadRepository: preflightResult.PullRequest.HeadRepository,
				HeadBranch:     preflightResult.PullRequest.HeadBranch,
				GitRoot:        preflightResult.Git.Root,
				Commit:         commit,
				ParentEvidence: artifactReq.ParentEvidence,
				ParentHeadSHA:  artifactReq.ParentHeadSHA,
			})
			if err != nil {
				fmt.Fprintf(ui.progress, "Review artifact Evidence inheritance unavailable; falling back to head polling: %v\n", err)
				inherited = false
			}
			if !inherited {
				evidence = reviewsource.Evidence{}
			}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 2636 - 2652, Update the error handling
around inheritReviewArtifactEvidence so proof-command failures do not return an
error after maybeRunFinalPush has succeeded. Log the git error as a warning,
reset evidence to the empty value, and continue through the existing normal
head-polling fallback; preserve the current inherited-evidence path and
structural-mismatch behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:35b80f0596558355cfc1fbfe -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: A post-push artifact-proof transport failure now warns, clears inherited evidence, and falls back to current-head polling instead of aborting the Run. `TestRunWatchArtifactEvidenceProofFailureFallsBackAfterPush` and the successful-inheritance regression passed.
