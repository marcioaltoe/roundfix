---
spec: 0237-review-selection-failures-that-say-why
prd: _prd.md
created: 2026-10-06
---

# Review selection failures that say why — Technical Spec

## Executive Summary

An `acpx prompt` process repeats the Agent Client Protocol from the start:
`initialize`, `session/load` or `session/new`, `session/set_model` when a model
is passed, then `session/prompt`. In JSON mode acpx echoes each request on
stdout, and a failure at any step arrives as a JSON-RPC error line on stdout
with the adapter's message, an empty stderr and exit `1`. The ACPX Runner
parses that line and then maps only the exit code, so every step yields the
same `agent/protocol error`, which is what Pantheon saw. The fix keeps the
first error line, names the step by the request it answers, and records
whether `session/prompt` had been sent. `roundfix review` puts the step and the
message in its reason and record, and retries a failure once on the same
selection when it happened before the review prompt was sent. The trade-off
accepted is one more Agent Session preparation before a configured fallback,
and a retry decision that rests on acpx echoing its requests, over a reason
that says nothing and a maintainer who skips the review (ADR-0242).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  review record gains three optional camel-case fields. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the review still drives the local
  acpx; no request, credential or forge read is added, and tests use fake
  runners and a fake acpx. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0242 (this Spec) governs the named
  step, the adapter's message and the single retry before the prompt.
  ADR-0020: "Without a parsed result, a nonzero exit stays a Batch failure
  exactly as before", so the exit mapping and error types stay and only gain
  the protocol detail. ADR-0227: "Every other exit after a parsed result stays
  a transport anomaly", untouched. ADR-0050: "once Agent work begins, Roundfix
  fails the Work Item instead of switching models", and ADR-0114: "Roundfix
  therefore treats opening an Agent Session as selection, not work"; the
  review's retry stays before the prompt and Runs do not change. ADR-0153:
  "Pre-PR review is an explicit provider policy", so a retry or a block never
  selects `none`. ADR-0151: "Configured review profiles select the independent
  reviewer", whose fallback chain still follows the preferred selection.
  ADR-0174: "A pre-PR review reads the final message", and ADR-0197: "A
  pre-PR reviewer lineage spans at most two rounds"; round two's resume keeps
  its rule and gains the same retry. ADR-0187 and ADR-0189 govern the Roundfix
  Skill edit and its version. ADR-0184: "A TechSpec now declares numbered
  Surface Transcripts", answered below with the reason none applies. The gate
  is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and
  ADR-0167; ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  consistency; ADR-0166, ADR-0178 and ADR-0182 bind each Task commit. ADR-0169's merge-base diff and ADR-0196's finding validation apply unchanged to a review that reaches the reviewer, and ADR-0233 regenerates the raised skill version at merge. ADR-0069 and ADR-0238 cite ADR-0050 but decide Baseline analysis and the light tier; ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry; ADR-0165 cites ADR-0153 but decides how a blocking review after archive parks publication; ADR-0192 cites ADR-0178 but decides how a conflict confined to declared derived paths is resolved; ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park; ADR-0180 cites ADR-0069 but decides the built-in selections, and ADR-0181 and ADR-0217 cite ADR-0180 but decide where a configuration is compared with the recommendation and how a Cursor selection names its model; ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot; ADR-0237 cites ADR-0229 but decides how a Delivery Retry records a merge made outside the queue; this Spec changes none of them, so none applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths; express maintainer
  authorization: "considere autorizado a ajustar todas as skills se
  necessário", and the maintainer's answer "Concedo" of 2026-10-06 for the
  Governed Paths this Spec declares; bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0237-review-selection-failures-that-say-why/_authorization.md`.

## System Architecture

No new package, command or flag. Two existing components change:

- **ACPX Runner** (`internal/agent`): `readPromptStream` already parses every
  stdout line; it also tracks the client requests acpx echoes, their
  responses and the first JSON-RPC error. `RunPrompt` attaches that trace to
  the `BatchFailureError` it maps from a non-zero exit without a parsed
  result, and `classifyNoOutputFailure` carries it into the
  `SelectionFailureError`. A new file holds the `ProtocolFailure` type and
  `DescribeProtocolFailure`, which also places a session-preparation failure.
- **Pre-PR Review Command** (`internal/cli`): `runConfiguredReviewSession`
  retries a failure before the prompt once per selection, and
  `classifyReviewCommandResult` fills the new record fields.

The daemon calls the same runner, so a Run's selection-failure messages gain
the step text; its fallback logic and reason codes do not change, because the
protocol detail travels in a new field and not in the wrapped error chain.

## Implementation Design

### Interfaces

```go
// internal/agent/protocol_failure.go
const ProtocolFailureMessageLimit = 512

type ProtocolFailure struct {
	Step       string // Invariants 2-4
	Code       int    // JSON-RPC error code; 0 without an error line
	Message    string // JSON-RPC error message, Invariant 5; "" without an error line
	PromptSent bool   // a session/prompt request was echoed before the failure
}

type ProtocolDescription struct {
	Step, Message string
	PromptSent    bool
	FromPrompt    bool // placed from a Protocol field, Invariant 4
}

func DescribeProtocolFailure(err error) (ProtocolDescription, bool)

// internal/agent: one field on each existing error type.
type BatchFailureError struct{ /* ...existing... */ Protocol *ProtocolFailure }
type SelectionFailureError struct{ /* ...existing... */ Protocol *ProtocolFailure }
```

```go
// internal/cli: record fields.
type reviewSelectionRetry struct {
	Selection int    `json:"selection"`
	Step      string `json:"step"`
	Message   string `json:"message,omitempty"`
}
// reviewRecord gains:
//   FailedStep       string                 `json:"failedStep,omitempty"`
//   AdapterMessage   string                 `json:"adapterMessage,omitempty"`
//   SelectionRetries []reviewSelectionRetry `json:"selectionRetries,omitempty"`
```

### Invariants

1. The runner tracks a stdout line as a client request when it has a non-null
   `id` and a `method` in `initialize`, `authenticate`, `session/new`,
   `session/load`, `session/set_model`, `session/set_mode`,
   `session/set_config_option` or `session/prompt`; a line with an `id`, no
   `method` and a `result` or `error` answers the request with that id.
   `PromptSent` is true once a `session/prompt` request was tracked.
2. The first line with an `error` object is the failure. Its step is the
   method of the tracked request whose id it carries; when its id is null or
   unmatched, the method of the latest unanswered tracked request.
3. When no tracked request is unanswered, the step is `session/prompt` if one
   was sent and `session setup` otherwise; with no tracked request at all it is
   `adapter startup`. Without an error line the step is derived the same way
   and `Code` and `Message` stay empty.
4. `DescribeProtocolFailure` first looks for a `Protocol` on a
   `SelectionFailureError` or `BatchFailureError` in the chain, which comes
   from an acpx prompt process (the review prompt or the deferred-effort
   warm-up), and reports it with `FromPrompt` true. Otherwise it places a
   session-preparation failure with `FromPrompt` false: a
   `SelectionRejectedError` whose operation starts with `ensure` and a
   `ModelNotAdvertisedError` are `sessions ensure`; one whose operation starts
   with `set ` is that operation (`set model`, `set reasoning_effort`); an
   `AccessPolicyError` is `set-mode`; an `AdapterProbeError` is
   `adapter startup`. Any other error is not placed.
5. A message keeps only the JSON-RPC `message` string, never `data`, the
   request parameters or the prompt. Without an error line, the description's
   message is acpx's stderr tail. Either is collapsed to one line (every run
   of whitespace becomes one space) and cut to `ProtocolFailureMessageLimit`
   bytes on a rune boundary.
6. `Protocol` is set only on a `BatchFailureError` mapped from a non-zero exit
   without a parsed prompt result, and copied by `classifyNoOutputFailure`.
   `Reason`, `Err`, `Stderr`, the exit mapping, the error types and the
   daemon's reason codes are unchanged.
7. The review retries at most once per selection. It retries a
   `PrepareSession` error that `reviewSelectionCanFallback` accepts, after
   ending that session, by preparing the same selection again; and a
   `RunPrepared` error that `DescribeProtocolFailure` places with
   `FromPrompt` true and `PromptSent` false, by calling `RunPrepared` again on
   the prepared session. Context cancellation and deadline are never retried.
8. A failure after `session/prompt` was sent, or one not placed, is never
   retried. A configured fallback activates only after the selection's retry
   failed before the prompt, as today. No path selects `none`.

### Error text

`BatchFailureError` and `SelectionFailureError` print, after their reason, ` at
<step>`, then `: <message> (JSON-RPC <code>)` when the error line had a
message, then the existing `: <err>` part. Pantheon's case becomes, for a
failure answered to `session/prompt`:
`Agent Selection failed for runtime "codex": agent/protocol error at session/prompt: <message> (JSON-RPC -32603)`.

### Data Models

The review record gains `failedStep`, `adapterMessage` and `selectionRetries`,
all omitted when empty; older records stay readable, and reuse and dispose
read only the fields they read today.

### API Contracts

1. API Contract: protocol failure from a prompt process. `ACPXRunner.RunPrompt`
   with stdout that ends in a JSON-RPC error line, no parsed result and exit
   `1` returns a `SelectionFailureError` with `Reason == "agent/protocol error"`
   whose `Protocol` names the step, code and message per Invariants 1 to 5,
   and whose text follows Error text. Without an error line the step is still
   named and the message is empty.
2. API Contract: description. `DescribeProtocolFailure` returns the step,
   the bounded adapter message, `PromptSent` and `FromPrompt` for API
   Contract 1's errors and for the preparation failures of Invariant 4, and
   `false` for any other error.
3. API Contract: review record. A review blocked by a runtime failure that
   `DescribeProtocolFailure` places carries `failedStep` and, when non-empty,
   `adapterMessage`; its reason stays `review runtime failure: <err>`, and when
   the failing selection was retried it ends with
   ` (after one automatic retry before the prompt)`.
4. API Contract: retry. Per Invariants 7 and 8, before each retry the command
   writes to stderr
   `roundfix: review Agent Selection failed before the prompt at <step> (<message>); retrying selection <index> once.`
   and appends `{selection, step, message}` to `selectionRetries`; a review
   that succeeds after a retry keeps the entry and its normal exit code.

### Surface Transcripts

None. The changed surfaces are the stderr reason and the review record of a
review whose ACP runtime fails mid-protocol, which the built product reaches
only through a live adapter. They are proved at the command entry point with
the fake review runner, at the runner seam with stdout recorded from acpx
0.19.4, and in QA against the installed acpx with a fake ACP adapter.

## Coverage Map

- Goal 1 → Invariants 1 to 5; API Contracts 1, 2 and 3.
- Goal 2 → Invariant 7; API Contract 4.
- Goal 3 → Invariant 8; API Contracts 3 and 4.
- User Story 1 → API Contracts 1 and 3.
- User Story 2 → API Contract 4.
- User Story 3 → API Contract 4 (`selectionRetries`).
- Core Feature 1 → API Contracts 1 and 2.
- Core Feature 2 → API Contract 3.
- Core Feature 3 → API Contract 4.
- Core Feature 4 → Build Order 1.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 2 and 3.

## Integration Points

- acpx 0.19.4 in JSON mode echoes each client request and prints each
  response, notification and error as one JSON-RPC line on stdout; it
  reconnects to a named session with `initialize` and `session/load` for every
  `prompt`, and applies `--model` with `session/set_model` before
  `session/prompt`. Its own failures carry a null id and a `data.acpxCode`,
  which the runner does not read.
- The Agent Client Protocol's error object carries `code` and `message`;
  only those two reach Roundfix's text.

Research record. The Secondbrain was read first (`wiki/index.md`, then a
`qmd` query on the defect): it returned the Roundfix mirror of ADR-0114 and
the Spec 0041 readiness work, which set the boundary that opening a session is
selection, and the triaged Pantheon report, which gave the reason and the
adapter version. Exa found the Agent Client Protocol overview, initialization
and prompt-turn pages, which fixed the step names and the error object. The
measurement against the installed acpx 0.19.4 with a fake adapter, recorded in
`_prd.md` → Acceptance evidence, decided the design: the message was never
missing, only discarded.

## Testing Approach

1. Runner seam: `runFakeACPXPrompt` in `internal/agent/acpx_runner_test.go`
   drives the real `RunPrompt` against a fake acpx. A new
   `internal/agent/protocol_failure_test.go` feeds stdout recorded from acpx
   0.19.4 for failures at `initialize`, `session/set_model` and
   `session/prompt`, an adapter disconnect after `session/prompt`, an exit
   without any line, and a message longer than the limit, asserting API
   Contract 1, Invariants 1 to 6 and that the three step errors differ; and it
   asserts `DescribeProtocolFailure` for each preparation type of Invariant 4.
   `TestACPXPromptExitClassificationMatrix` stays unchanged and green.
2. Review seam: `newReviewCommandFixture` and the fake `reviewCommandRunner`,
   in a new `internal/cli/review_selection_retry_test.go`: a preparation
   selection failure then success; a prompt-process failure before
   `session/prompt` then success; a failure after `session/prompt`; two
   failures before the prompt without and with a configured fallback; a
   generic preparation error; and a round-two resume whose prompt process
   fails before `session/prompt`.
3. The three tests that pin a fallback after one failure,
   `TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt`,
   `TestReviewRoundTwoContinuesTheRecordedFallbackSelection` and
   `TestReviewValidatorUsesTheSuccessfulFallbackSelection`, now give the
   preferred selection two failures and expect one more preparation; nothing
   else in them changes.

## Build Order

1. The Roundfix Skill's `review` reference, its mirror and the `review`
   command guide describe the step, the message, the retry and the record
   fields; the skill version is raised and recorded.
2. The ACPX Runner keeps the failing step, the adapter's message and whether
   the prompt was sent (API Contracts 1 and 2).
3. The review names the step, records it and retries once before the prompt
   (API Contracts 3 and 4) (depends on: 1, 2).
4. The final QA gate (depends on: 3).

## Risks & Considerations

- A later acpx that stops echoing requests makes every prompt-process failure
  `adapter startup` or `session setup` with `PromptSent` false, so it would be
  retried once even after a prompt; the review is read-only, and the record
  lists the retry. The QA gate re-measures the installed acpx.
- An adapter message could quote something sensitive the adapter chose to
  print; Roundfix reads only the `message` field, never request data or the
  prompt, and bounds it to one line of 512 bytes.
- A selection that fails twice before the prompt costs one more preparation
  before its fallback; the incident's failures were transient.

## Vocabulary Contract

- emits: `internal/cli/review.go`
  pattern: `retrying selection`
  documented-in: `docs/user-guide/commands/review.md`
- emits: `internal/cli/review.go`
  pattern: `selectionRetries`
  documented-in: `docs/user-guide/commands/review.md`
- emits: `internal/cli/review.go`
  pattern: `failedStep`
  documented-in: `docs/user-guide/commands/review.md`

No glossary term is adopted; "protocol step" is used in its plain sense beside
**Agent Selection** and **Fallback Selection**.

## Decisions

- Name the step by the request the error answers and keep only its message;
  see ADR-0242.
- Carry the protocol detail in a new field, not the wrapped error, so Runs'
  reason codes and fallback stay as they are.
- Retry once per selection, only before the prompt, before any fallback; see
  ADR-0242.
