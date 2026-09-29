---
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
status: active
created: 2026-09-29
surfaces: [backend, cli, docs]
---

# A pre-PR review that keeps its verdict and its own record

## Executive Summary

Keep the Agent's message boundaries where the runner reads the ACP stream, and
read the review verdict from the last non-blank message. When both chunks carry
a `messageId`, a different one starts a new message. Without one, a message
chunk that follows a thought, a tool call or a plan starts a new message. `ExecuteResult` reports the
messages in order. The review classifies `ExecuteResult.Answer()`, and a
sealed prompt parses its final message. The review record and its answer move
into a per-checkout directory of the Artifact Directory, and every reader looks
only there. The disposition ledger stays where it is.

The trade-off accepted: a review that ends in `No findings.` after commentary
now passes, because commentary is no longer part of the verdict. Two chunks
with no boundary still read as one message. No Run Database or record schema
changes.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; the record
  directory name is a path digest no other surface reads. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the local
  Agent runtime; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0174 (this Spec) is implemented
  here. ADR-0017, ADR-0020, ADR-0080, ADR-0091, ADR-0093, ADR-0094, ADR-0096,
  ADR-0104, ADR-0117, ADR-0151, ADR-0153, ADR-0155, ADR-0156, ADR-0165,
  ADR-0166, ADR-0168 and ADR-0169 hold as the PRD states. ADR-0097, ADR-0167,
  ADR-0170, ADR-0171, ADR-0172 and ADR-0173 do not apply, for the reasons the
  PRD records. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

Two existing seams change, and one file is new.

- `internal/agent`:
  - `stream.go` gains `StreamUpdate.MessageID`, and `acp_stream.go` reads it
    from `agent_message_chunk`.
  - `agent.go` gains `ExecuteResult.Messages` and `ExecuteResult.Answer`.
  - The new `agent_messages.go` holds the message log.
  - `acpx_runner.go` (`readPromptStream`, `handleStdoutLine`) and `sealed.go`
    (`sealedStreamParser`) feed every session update to the log.
- `internal/cli/review.go` — the classifier reads `result.Answer()`. The
  record, the answer, reuse, stale-answer removal and `review dispose` resolve
  the current checkout's record directory.

The Delivery Queue keeps reading the verdict from the command's standard output
(`internal/cli/deliver_workflow.go` `runReview`) and does not change.

## Implementation Design

### Interfaces

```go
// internal/agent
type ExecuteResult struct {
	// ...existing fields...
	Message  string   // every Agent message, separated by one blank line
	Messages []string // each Agent message, in order
}
// Answer returns the last message that is not blank, or Message when the
// runner reported no Messages.
func (result ExecuteResult) Answer() string

type agentMessageLog struct{ /* messages, last update kind, last messageId, bytes */ }
func (log *agentMessageLog) observe(update StreamUpdate)
func (log *agentMessageLog) messages() []string
func (log *agentMessageLog) final() string

// internal/cli
func reviewCheckoutDir(artifactDir, checkoutRoot string) string
```

### Data Models

No schema change. `StreamUpdate` gains `MessageID string`, set only for a
message chunk whose update carries a non-empty `messageId`. The review record's
fields are unchanged; only where the file lives and `answerPath`'s value change.

### Agent messages keep their boundaries

`agentMessageLog.observe` receives every parsed session update:

- When a message chunk and the previous message chunk both carry a non-empty
  `MessageID`, the identifiers alone decide. A different identifier starts a
  new message, and the same identifier appends, even across other updates.
  That is the Agent Client Protocol's own rule.
- Otherwise a message chunk starts a new message when the last observed update
  was a thought, a tool call, a tool update or a plan. Any other update, such
  as a status or an unrecognized raw update, never splits a message.
- Otherwise it appends to the current message.

`messages()` returns the texts in order. `final()` returns the last message
whose trimmed text is not empty. The log also counts the bytes appended, for
the sealed cap.

- **ACP runner.** `handleStdoutLine` takes the log instead of a
  `strings.Builder` and observes every update before publishing it.
  `readPromptStream` returns the log. The runner sets `Messages` to its
  messages and `Message` to them joined with `"\n\n"`. One message therefore
  leaves `Message` exactly as before.
- **Sealed prompts.** `sealedStreamParser.consumeSessionUpdate` observes
  message, thought, tool and plan updates, the same boundary set the prompt
  parser uses. The cap `SealedPromptMaxOutputBytes` applies to
  the log's total message bytes. A tool update still sets `ErrSealedToolUse`.
  `SealedPromptResult.Output` is the bytes of `final()`.
- **The review.** `classifyReviewCommandResult` classifies `result.Answer()`
  with today's rules, and the `empty agent output` check applies to it.
  `finishReviewCommandWithAnswer` still persists `result.Message`, so the answer
  file keeps every message.

### A review record per checkout

`reviewCheckoutDir(artifactDir, checkoutRoot)` returns
`<artifactDir>/pre-pr-review/<key>`. The key is the first 16 hex characters of
the SHA-256 of `filepath.Clean(checkoutRoot)`, and `checkoutRoot` is the root
`preflight.InspectGit` reports, the same value the record's `repository` field
holds. The following all use that directory:

- `persistReviewRecord` and `persistReviewAnswer`, each of which creates the
  directory with mode `0755` before its temporary file
- `removeReviewAnswer`
- `reusableReviewRecord`
- `runReviewDisposeCommand`

Temporary files are created inside it. The disposition ledger and its lock stay
at `<artifactDir>/pre-pr-review-dispositions.jsonl`. The existing check that
refuses a record whose `repository` differs from the checkout stays, as a guard
against a digest collision. Nothing reads, writes or removes
`<artifactDir>/pre-pr-review.json` or `<artifactDir>/pre-pr-review-answer.txt`.

## API Contracts

1. API Contract: `roundfix review` classifies the reviewer's final non-blank
   message with today's rules, and keeps every message in
   `pre-pr-review-answer.txt`, separated by a blank line.
2. API Contract: `roundfix review` writes `pre-pr-review.json` and
   `pre-pr-review-answer.txt` under
   `<Artifact Directory>/pre-pr-review/<checkout key>/`, and `answerPath`
   names that file. Reuse and `roundfix review dispose` read only that
   directory. The ledger stays at
   `<Artifact Directory>/pre-pr-review-dispositions.jsonl`.

## Coverage Map

- Goals 1–2 → Agent messages keep their boundaries (ACP runner, the review);
  API Contract 1.
- Goal 3 → Agent messages keep their boundaries (sealed prompts).
- Goals 4–5 → A review record per checkout; API Contract 2.
- Core Features 1–3 → Agent messages keep their boundaries.
- Core Feature 4 → A review record per checkout.
- Success Metrics 1–4 → Testing Approach 1–2.
- Success Metric 5 → Testing Approach 3.

## Integration Points

- **acpx and the ACP adapters.** The runner reads their `session/update`
  lines unchanged. `messageId` is optional in ACP v1, and codex-acp 2.0.0 sends
  one per message.
- **The Baseline's ACP analyzer** (`internal/baselineacp`) parses sealed
  output as a JSON proposal and gains the final-message rule without a change
  of its own.
- **Spec 0182's merge base.** The record's `baseCommit` and reuse key are what
  Spec 0182 defines; this Spec moves only the file.
- **The Delivery Queue.** It reads the record from standard output, so its
  verdict path is unchanged. Its item worktree now keeps its own record, which
  `roundfix review dispose` finds there.

## Testing Approach

Every test uses fake ACP streams, fake Agent runners, temporary repositories
and temporary Artifact Directories. None reaches a reviewer or `~/.roundfix`.

1. **Message log.** Unit tests of `agentMessageLog`:
   - a tool call or a thought between chunks without `messageId` splits them;
   - a changed `messageId` splits them, and the same `messageId` across a tool
     call does not;
   - chunks with neither boundary, or with only a status update between them,
     stay one message;
   - `Answer()` skips a trailing blank message and falls back to `Message`.

   Through `runFakeACPXPrompt`, a message, a tool call and a message report two
   `Messages` and a joined `Message`. Through `parseSealedPromptStream`,
   commentary, a thought and a JSON message yield exactly the JSON.
   `TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason`,
   `TestSealedACPXPromptReturnsMessageAndClosesSession` and
   `TestSealedPromptStreamIgnoresNonTerminalJSONRPCResults` stay green.
2. **Final-message verdict.** Through `newReviewCommandFixture` with a fake
   runner whose result carries `Messages`:
   - a progress message then `Findings:` classifies `findings` with one item;
   - a progress message then `No findings.` classifies `reviewed`;
   - a final message with both verdicts still blocks;
   - the answer file holds both messages.

   `TestReviewClassifiesVerdictVariants`,
   `TestReviewPassesOnlyAWholeAnswerVerdict` and `TestReviewKeepsTheRawAnswer`
   stay green.
3. **Record per checkout.** A repository and a `git worktree add` checkout
   share one configured Artifact Directory:
   - each review writes only under its own checkout directory;
   - a review in the second checkout leaves the first checkout's record and
     answer bytes unchanged, and the first still reuses and disposes its
     record;
   - `review dispose` in a checkout without a record is refused even though
     the other has one;
   - a matching record at the old shared location is neither reused nor
     removed.

   Existing tests that read the record or answer at the old path are updated
   to `reviewCheckoutDir`, and their names stay.
4. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Agent messages keep their boundaries, the sealed final message and the
   final-message verdict (depends on: none).
2. A review record per checkout (depends on: 1 — both edit
   `internal/cli/review.go`, `docs/user-guide/commands.md` and the Roundfix
   skill).
3. Terminal QA (depends on: 1, 2).

## Risks & Considerations

- **A verdict split across messages.** An Agent that sends `Findings:` and its
  list as two messages would be read from the list alone and block, which is no
  worse than today. The prompt already asks for one answer. The QA gate records
  the classification of real-shaped streams.
- **Commentary that passes.** Reading the final message means commentary
  before `No findings.` no longer blocks. The classifier still refuses any other
  content beside that verdict inside the final message.
- **Unattended readers.** The Delivery Queue reads only standard output, and
  `review dispose` refuses a missing or foreign record. No reader falls back to
  another checkout's directory.
- **Directory growth.** One directory per checkout that ever ran a review,
  holding two small files. It is recorded as a limit and not pruned.

## Decisions

- Keep message boundaries in the runner; never search inside a line. See
  ADR-0174.
- The final non-blank message is the answer, and the answer file keeps all
  messages. See ADR-0174.
- One record directory per checkout, located by path digest and guarded by the
  record's `repository`. See ADR-0174.
- Ignore the old shared location; do not migrate or delete it.
