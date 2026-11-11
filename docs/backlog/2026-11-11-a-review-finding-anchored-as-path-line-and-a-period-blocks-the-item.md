---
type: fix
status: open
created: 2026-11-11
spec: null
reason: null
---

# A review finding anchored as path:line and a period blocks the item

## Symptom

On 2026-11-11 the pre-PR review of Spec 0376 (Roundfix 0.139.1, provider `codex`) answered two findings, each starting with an anchor followed by a period, for example `docs/user-guide/commands/skills.md:34. Failure: …`. Roundfix read neither anchor, so the review ended `blocked` with the reason `findings name no file and line` and the Delivery Queue parked the item as `review-blocked`. Both findings were real and anchored in the candidate diff.

## Where

`parseReviewFindingAnchor` in `internal/cli/review_validation.go`: after the `path:line[-line]` match it accepts only a space, `:`, `,` or `)` as the next character, so a sentence-ending `.` (and likewise `;`) makes the anchor unreadable.

## Expected

An anchor followed by a period or a semicolon is read like one followed by a space, so such findings are validated against the candidate diff instead of blocking the item. A path that itself ends in a period-like suffix must still parse as the whole path.

## Evidence

Review record `~/.roundfix/artifacts/339f8dac2b687a04/pre-pr-review/69ca9a33cab51858/` on the mini: `pre-pr-review-answer.txt` holds the two findings; `pre-pr-review.json` has `outcome` `blocked` and `reason` `findings name no file and line`. The operator fixed both findings on the item branch and retried.
