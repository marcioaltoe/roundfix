---
task: task_01
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
status: pending
type: backend
complexity: high
---

# Task 01: Agent messages keep their boundaries and the review reads the final message

## Overview

The Agent runner appends every `agent_message_chunk` of a session to one string with no boundary, so a progress message and the final answer reach the Pre-PR Review classifier as one line. On 2026-09-29 a real five-finding answer was recorded `blocked` because it began `…security issues.Findings:`. The same concatenation feeds the sealed prompts the Baseline parses as JSON.

This Task keeps message boundaries where the runner reads the ACP stream, reports each message, and makes the review classify the last non-blank message and a sealed prompt return it. The stream comes from the local acpx process. Its only readers are the review classifier, the review answer file and the Baseline's sealed-prompt parser. The classifier's rules inside the final message do not change.

## Requirements

1. MUST add `MessageID string` to `StreamUpdate` in `internal/agent/stream.go`, and MUST set it in `internal/agent/acp_stream.go` from an `agent_message_chunk` update's optional `messageId` field. An absent or null value leaves it empty, and no other update kind sets it.
2. MUST add an unexported message log in the new `internal/agent/agent_messages.go`, as `_techspec.md` → Agent messages keep their boundaries states:
   - When a message chunk and the previous message chunk both carry a non-empty `MessageID`, the identifiers alone decide: a different one starts a new message, and the same one appends, even across other updates.
   - Otherwise a message chunk starts a new message only when the last observed update was a thought, a tool call, a tool update or a plan. A status or raw update never splits a message.
   - It reports the messages in order, the last message whose trimmed text is not empty, and the total message bytes.
3. MUST add `Messages []string` to `ExecuteResult` and the method `func (result ExecuteResult) Answer() string`, which returns the last non-blank element of `Messages`, or `Message` when `Messages` is empty.
4. MUST make `handleStdoutLine` and `readPromptStream` in `internal/agent/acpx_runner.go` observe every parsed session update through the log. The runner MUST set `Messages` to the log's messages and `Message` to them joined with one blank line, so a session with one message keeps its `Message` byte-identical.
5. MUST make `sealedStreamParser` in `internal/agent/sealed.go` observe message, thought, tool and plan updates through the log, so every update kind that splits a message elsewhere splits it in a sealed prompt too. It MUST return `SealedPromptResult.Output` as the bytes of the final non-blank message. It MUST apply `SealedPromptMaxOutputBytes` to the total message bytes and still refuse a tool update with `ErrSealedToolUse`.
6. MUST make `classifyReviewCommandResult` in `internal/cli/review.go` classify `result.Answer()` with today's rules, including the `empty agent output` check. The answer file MUST still receive `result.Message`, so it keeps every message.
7. MUST state in the `roundfix review` section of `docs/user-guide/commands.md` and in `.agents/skills/roundfix/SKILL.md` that Roundfix `reads the verdict from the reviewer's final message` and that the answer file `keeps every message, separated by a blank line`, then MUST run `make skills-sync` so `skills/roundfix/SKILL.md` matches.
8. MUST change no exported function signature, rename or remove no top-level test, and keep `TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason`, `TestSealedACPXPromptReturnsMessageAndClosesSession`, `TestSealedPromptStreamIgnoresNonTerminalJSONRPCResults`, `TestReviewClassifiesVerdictVariants`, `TestReviewPassesOnlyAWholeAnswerVerdict` and `TestReviewKeepsTheRawAnswer` green without editing them.
9. MUST put the new tests in `internal/agent/agent_messages_test.go` and `internal/cli/review_final_message_test.go`, using fake ACP streams (`runFakeACPXPrompt`, `parseSealedPromptStream`) and the review command's fake runner (`newReviewCommandFixture`), never a real reviewer.

## Subtasks

- [ ] Read the optional `messageId` of a message chunk.
- [ ] Keep message boundaries in one log shared by the prompt and sealed parsers.
- [ ] Report each message and the final answer from `ExecuteResult`.
- [ ] Classify the review's final message and keep every message in the answer file.
- [ ] Update the guide, the skill and its mirror.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] Chunks without `messageId` split at a thought or a tool call, and stay one message across a status update or with no update between them.
- [ ] Chunks with `messageId` split where the identifier changes, even with no update between them, and stay one message across a tool call when the identifier is the same.
- [ ] `Answer()` returns the last non-blank message, and returns `Message` for a result with no `Messages`.
- [ ] A fake ACP stream with a message, a tool call and a message reports two `Messages` and a `Message` that separates them with a blank line. A one-message stream keeps `Message` unchanged.
- [ ] A sealed stream with a commentary message, a thought or a plan update, and a JSON message returns exactly the JSON, and a stream whose total message bytes exceed the cap is still refused.
- [ ] A review whose runner reports a progress message and then `Findings:` with one list item records `findings` with `F1`. A progress message then `No findings.` records `reviewed`.
- [ ] A final message carrying both verdicts still records `blocked`, and the answer file holds every message.
- [ ] The guide, the skill and its mirror state the final-message rule, and `make skills-sync-check` passes.

## Context

- interface: `internal/agent/stream.go`
- interface: `internal/agent/acp_stream.go`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/sealed.go`
- creates: `internal/agent/agent_messages.go`
- creates: `internal/agent/agent_messages_test.go`
- interface: `internal/cli/review.go`
- creates: `internal/cli/review_final_message_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- instruction: `docs/adr/0174-a-pre-pr-review-reads-the-final-message-and-keeps-a-record-per-checkout.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAgentMessagesSplitAtAThoughtOrToolCallWithoutMessageIDs|TestAgentMessagesFollowTheirMessageIDs|TestAgentMessageChunksWithoutABoundaryStayOneMessage|TestExecuteResultAnswerIsTheLastNonBlankMessage|TestACPXRunPromptReportsEachAgentMessage|TestSealedPromptOutputIsTheFinalMessage|TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason|TestSealedACPXPromptReturnsMessageAndClosesSession|TestSealedPromptStreamIgnoresNonTerminalJSONRPCResults)$" ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAgentMessagesSplitAtAThoughtOrToolCallWithoutMessageIDs TestAgentMessagesFollowTheirMessageIDs TestAgentMessageChunksWithoutABoundaryStayOneMessage TestExecuteResultAnswerIsTheLastNonBlankMessage TestACPXRunPromptReportsEachAgentMessage TestSealedPromptOutputIsTheFinalMessage TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason TestSealedACPXPromptReturnsMessageAndClosesSession TestSealedPromptStreamIgnoresNonTerminalJSONRPCResults; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestReviewClassifiesTheFinalMessageAfterProgressText|TestReviewAcceptsNoFindingsAsTheFinalMessage|TestReviewStillBlocksBothVerdictsInTheFinalMessage|TestReviewAnswerFileKeepsEveryMessage|TestReviewClassifiesVerdictVariants|TestReviewPassesOnlyAWholeAnswerVerdict|TestReviewKeepsTheRawAnswer)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewClassifiesTheFinalMessageAfterProgressText TestReviewAcceptsNoFindingsAsTheFinalMessage TestReviewStillBlocksBothVerdictsInTheFinalMessage TestReviewAnswerFileKeepsEveryMessage TestReviewClassifiesVerdictVariants TestReviewPassesOnlyAWholeAnswerVerdict TestReviewKeepsTheRawAnswer; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four new tests do not exist.
- `for file in docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "reads the verdict from the reviewer's final message" || { printf 'missing phrase in %s: %s\n' "$file" "reads the verdict from the reviewer's final message" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "keeps every message, separated by a blank line" || { printf 'missing phrase in %s: %s\n' "$file" "keeps every message, separated by a blank line" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task the phrases are absent.

## References

- [_prd.md](_prd.md) — Goals 1–3; Core Features 1–3; Success Metrics 1–4
- [_techspec.md](_techspec.md) — Agent messages keep their boundaries; Interfaces; API Contract 1; Testing Approach 1–2; Build Order 1
- [references/2026-09-29-a-review-verdict-after-progress-text-is-unclassifiable.md](references/2026-09-29-a-review-verdict-after-progress-text-is-unclassifiable.md)
- ADR-0174; ADR-0017; ADR-0020; ADR-0153
