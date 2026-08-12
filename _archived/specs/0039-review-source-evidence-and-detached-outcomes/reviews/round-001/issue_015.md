---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/watch/watch_test.go
line: 1462
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIp9,comment:PRRC_kwDOS0qyts7aUVER
review_hash: 01f802a023523d4d093178f26f7ccc70420ce25cf3e482051b293da40282d50e
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 015: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Document the "last result repeats" queue semantics.**

`Evidence` pops entries only while more than one remains, so the final `reviewEvidenceResult` is returned indefinitely. Several tests depend on this (steady-state pending or verified evidence) while others depend on exact call counts, so the behavior is load-bearing but implicit.



<details>
<summary>♻️ Proposed comment</summary>

```diff
+// Evidence returns queued results in order and repeats the final result for
+// every later call, so steady-state phases need only one trailing entry.
 func (source *fakeReviewEvidenceSource) Evidence(_ context.Context, req ReviewEvidenceRequest) (reviewsource.Evidence, error) {
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
type fakeReviewEvidenceSource struct {
	evidence reviewsource.Evidence
	err      error
	results  []reviewEvidenceResult
	requests []ReviewEvidenceRequest
}

// Evidence returns queued results in order and repeats the final result for
// every later call, so steady-state phases need only one trailing entry.
func (source *fakeReviewEvidenceSource) Evidence(_ context.Context, req ReviewEvidenceRequest) (reviewsource.Evidence, error) {
	source.requests = append(source.requests, req)
	if len(source.results) > 0 {
		result := source.results[0]
		if len(source.results) > 1 {
			source.results = source.results[1:]
		}
		return result.evidence, result.err
	}
	return source.evidence, source.err
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/watch/watch_test.go` around lines 1445 - 1462, Document the queue
semantics in fakeReviewEvidenceSource.Evidence: results are consumed while
multiple entries remain, while the final reviewEvidenceResult is retained and
returned for all subsequent requests. Add a concise comment near the results
handling, preserving the existing behavior and request tracking.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:6d55e9179f794580973b61b9 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Documented that the fake evidence source consumes queued results and repeats the terminal result, clarifying the fixture contract. `go test ./internal/watch` passed.
