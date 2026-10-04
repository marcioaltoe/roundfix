---
task: task_04
spec: 0222-a-review-through-claude-that-keeps-its-verdict
status: completed
type: backend
complexity: medium
---

# Task 04: The review bounds a Claude prompt in estimated tokens against half of its window

## Overview

Spec 0212 bounds every review by the bytes of its diff, a limit measured on
Codex. A Claude review sends the diff, its instructions and the Spec context
to a provider whose window is 1,000,000 tokens, and on 2026-10-02 a Claude
review overflowed that window while the command printed only a protocol error
(the adopted Backlog Entry
[2026-10-03](references/2026-10-03-the-review-through-claude-blocks-on-a-refused-permission.md)).
This Task bounds the whole Claude prompt in estimated tokens against half of
that window, before any provider call, and leaves Codex's byte bound as it is.
It is verifiable on its own through the review tests and the built binary.

## Requirements

1. MUST add `estimateReviewPromptTokens(prompt string) int`, returning the
   prompt's byte length divided by two, rounded up, and the constants
   `claudeReviewContextTokens = 1000000` and a budget of half of it.
2. MUST, in `runConfiguredReviewSession`, read the Spec context right after
   the diff and before the readiness probe, assemble the final prompt of round
   one or round two with its Spec context, and, for the `claude` provider,
   record the estimate in the new record field `EstimatedPromptTokens` (JSON
   `estimatedPromptTokens`, omitted when zero) and refuse with a
   `reviewPrePromptError` when it exceeds the budget, with the reason of API
   Contract 4:
   `review prompt too large: <n> estimated tokens after omitting <m> path(s) exceeds the claude review budget of 500000 tokens, half of its 1000000-token context window`.
   The refusal MUST happen before the readiness probe, any Agent Session
   preparation, any prompt and any fallback.
3. MUST stop applying the 917,504-byte diff bound to the `claude` provider,
   and MUST keep it, its place and its reason unchanged for `codex`, with no
   estimate recorded.
4. MUST add the tests named in Verification to the new file
   `internal/cli/review_prompt_bound_test.go`, sizing candidates from the
   constants and `reviewDiffBound`, not from literals:
   - a `claude` candidate whose prompt estimate exceeds the budget: exit `2`,
     Surface Transcript 1's stderr line, `estimatedPromptTokens` above the
     budget, and zero probe, prepare and prompt calls on the fake runner;
   - a `claude` candidate whose diff exceeds `reviewDiffBound` and whose
     estimate fits: exactly one prepared prompt, and the recorded estimate;
   - a `codex` candidate with that same diff: the existing
     `review diff too large` reason and exit `2`, with no estimate.
5. MUST keep the existing bound tests in `internal/cli/review_scope_test.go`
   unchanged and passing.
6. MUST NOT change the Codex review path, the review prompt text or the
   lineage rules.

## Subtasks

- [ ] Add the estimate, the window constant and the budget.
- [ ] Read the Spec context before readiness and bound the Claude prompt.
- [ ] Limit the byte bound to Codex.
- [ ] Add the bound tests.

## Acceptance Criteria

- [ ] A Claude prompt over half the window is refused before any runner call,
      naming its estimate and the budget.
- [ ] A Claude prompt within the budget is sent even when its diff exceeds the
      Codex byte bound, and a Codex review keeps its byte bound.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_prompt_bound_test.go`
- instruction: `internal/cli/review_scope_test.go`
- instruction: `internal/cli/review_test.go`
- instruction: `.agents/skills/roundfix/references/review.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewCommandRefusesAClaudePromptOverHalfItsWindow|TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound|TestReviewCommandKeepsTheCodexDiffByteBound)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewCommandRefusesAClaudePromptOverHalfItsWindow TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound TestReviewCommandKeepsTheCodexDiffByteBound; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `tmp="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix || exit 1; repo="$tmp/repo"; mkdir -p "$repo" "$tmp/home" || exit 1; git -C "$repo" init -q -b main || exit 1; printf '%s\n' 'defaults:' "  artifact_dir: $tmp/artifacts" 'pre_pr_review:' '  provider: claude' 'profiles:' '  review:' '    preferred:' '      runtime: claude' '      model: opus' '      reasoning_effort: high' '    fallbacks:' '      - runtime: claude' '        model: sonnet' '        reasoning_effort: high' > "$repo/.roundfixrc.yml" || exit 1; git -C "$repo" add . && git -C "$repo" -c user.name=qa -c user.email=qa@example.invalid commit -qm base || exit 1; git -C "$repo" checkout -qb feature || exit 1; awk 'BEGIN { for (i = 0; i < 20000; i++) printf "line %05d of a candidate larger than half of the claude context window\n", i }' > "$repo/large.txt" || exit 1; git -C "$repo" add . && git -C "$repo" -c user.name=qa -c user.email=qa@example.invalid commit -qm large || exit 1; out="$(cd "$repo" && env -u NODE_OPTIONS HOME="$tmp/home" PATH="/usr/bin:/bin" "$tmp/roundfix" review --base main 2>&1 >/dev/null)"; code=$?; test "$code" -eq 2 || { printf 'exit %s: %s\n' "$code" "$out" >&2; exit 1; }; printf '%s\n' "$out" | grep -qF -- "review prompt too large:" || { printf 'unexpected stderr: %s\n' "$out" >&2; exit 1; }; printf '%s\n' "$out" | grep -qF -- "exceeds the claude review budget of 500000 tokens, half of its 1000000-token context window" || { printf 'unexpected stderr: %s\n' "$out" >&2; exit 1; }` — expected: exit 0; before this Task the built binary refuses the candidate with the byte bound's `review diff too large` reason, so the command fails. The candidate is refused before the readiness probe, so no runtime is reached.

## References

- `_prd.md` → Core Feature 4; User Story 4; Success Metric 4; Goal 3; Goal 4
- `_techspec.md` → The Claude prompt bound; API Contract 4; Surface Transcript 1; Testing Approach 3; Build Order 4
- ADR-0227
- ADR-0169

## Result

Implemented the provider-specific admission bound in `internal/cli/review.go`.
Claude estimates the assembled round-one or round-two prompt, including Spec
context, as its byte length divided by two rounded up. Its window constant is
1,000,000 tokens and its budget is half that window. The optional
`estimatedPromptTokens` record field carries the estimate. An over-budget
prompt returns the specified `reviewPrePromptError` before readiness,
preparation, prompting or fallback. Codex keeps the existing diff bound at
its existing admission point and records no token estimate. Spec context is
now read before readiness; prompt text and lineage rules are unchanged.

Focused evidence for the acceptance criteria:

- Over-budget refusal: `TestReviewCommandRefusesAClaudePromptOverHalfItsWindow`
  exercises both rounds, checks exit `2`, the exact Surface Transcript 1
  stderr line and reason, an estimate above the budget in the persisted
  record, no answer path, and no additional probe, prepare, prompt or session
  end calls, with a fallback configured.
- Provider-specific admission: `TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound`
  checks that a diff above `reviewDiffBound` reaches exactly one prepared
  Claude prompt, includes Spec context, and records that whole prompt's
  estimate within budget. `TestReviewCommandKeepsTheCodexDiffByteBound` uses
  the same candidate content and checks exit `2`, the unchanged byte-bound
  reason, no runner calls, and absence of the estimate JSON field.
- `TestEstimateReviewPromptTokens` checks empty, even-byte, odd-byte and
  non-ASCII inputs. Existing scope and review tests remain unchanged.

Focused checks run:

- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test ./internal/cli -count=1 -run 'TestReviewCommand.*(HalfItsWindow|CodexByteBound|DiffByteBound)'`
  initially exited `1`: both Claude cases reproduced the old byte-bound
  refusal instead of the new behavior; the Codex case passed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test ./internal/cli -count=1 -run 'TestReview|TestEstimateReviewPromptTokens'`
  exited `0` after the final implementation and test edits (`ok`, 32.558s).
  This includes the unchanged bound tests in `review_scope_test.go`.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited `0`.

The incoming `status: in_progress` is Daemon-owned and was preserved. No
declared Verification command was run; verification and settlement remain
with the Daemon. No commit, push, or pull request was made. No follow-up scope
was identified.

## Carry-forward provenance

- Source Run: `run_20261004T145620Z_ee61ac9af5ab0ce2`
- Source commit: `586e52662a3fc928c41e0df2b9faf6ec860985d8`
