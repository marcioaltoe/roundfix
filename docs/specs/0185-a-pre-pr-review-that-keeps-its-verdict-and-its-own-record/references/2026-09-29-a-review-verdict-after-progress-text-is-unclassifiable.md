---
type: fix
status: promoted
created: 2026-09-29
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
reason: null
---

# A review verdict that follows the Agent's progress text is unclassifiable

## Symptom

On 2026-09-29, `roundfix review --base main` in `~/dev/roundfix-wave5` (head `3e77ed70`) ended `blocked` with exit 2: `unclassifiable agent output: neither a no-findings nor findings verdict is present`. The answer file held five well-formed findings. The Agent's progress sentence and its verdict arrived joined on one line, `…looking for concrete correctness, regression, or security issues.Findings:`, so no line started with `Findings:`. A full review round was spent, and the findings had to be read from the answer file by hand.

## Where

The verdict classifier in `internal/cli/review.go` requires a line that starts with `Findings:`, using `parseFindingsVerdictLine`. The Agent's separate messages reach it concatenated without a line break, in the review session's message accumulation.

## Expected

Separate Agent messages stay separate lines, or only the final message is classified. A verdict the Agent actually gave, `Findings:` followed by list items, classifies as findings even when progress text came before it in the session.

## Evidence

`~/.roundfix/artifacts/339f8dac2b687a04/pre-pr-review-answer.txt` as written at 14:20 on 2026-09-29, and the `pre-pr-review.json` record for head `3e77ed70`.
