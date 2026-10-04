---
spec: 0222-a-review-through-claude-that-keeps-its-verdict
prd: _prd.md
created: 2026-10-04
---

# A review through Claude that keeps its verdict — Technical Spec

## Executive Summary

The ACPX Runner turns every non-zero acpx exit after a parsed prompt result,
other than `130`, into a transport anomaly, and `roundfix review` blocks on an
anomaly before it reads the answer. acpx exits `5` whenever every permission
request of a turn was denied, which is what a read-only review session does
when a Claude reviewer asks for a command Claude Code does not consider
read-only. The fix keeps such a turn's result in the runner, marks the
refusal, and lets the review classify the answer by its verdict. The review
also names a runtime's `Prompt is too long` answer in its reason, and bounds a
Claude prompt in estimated tokens against half of the provider's
1,000,000-token window instead of applying Codex's byte bound to the diff
alone. The trade-off accepted is a verdict produced without a tool the
reviewer asked for, recorded as such, over a verdict thrown away; and a token
estimate at two bytes per token that can be off by a few tens of percent,
absorbed by the half-window reserve (ADR-0227).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; finding
  ids keep the `F<n>` sequence and the review record gains two optional
  camel-case fields. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the review still drives the local
  acpx; no request, credential or forge read is added, and tests use fake
  runners and a fake acpx. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0227 (this Spec) governs the
  refused-permission classification and the Claude token bound. ADR-0020:
  "a subsequent nonzero acpx exit is classified as teardown noise", which this
  Spec applies to the review only for the permission exit of a read-only
  `end_turn` turn; every other exit keeps the review's block. ADR-0174: "A
  pre-PR review reads the final message", so the verdict and the
  `Prompt is too long` line come from the final answer. ADR-0196's anchor
  validation and ADR-0197's two-round lineage apply unchanged to the kept
  findings; ADR-0153: "Pre-PR review is an explicit provider policy", kept as
  it is. ADR-0151: "Configured review profiles select the independent
  reviewer", and the profile stays as it is. ADR-0169: "The pre-PR review diffs
  the candidate from its merge base", the diff the bound measures. ADR-0187 and ADR-0189 govern
  the Roundfix Skill edit and its version. ADR-0184: "A TechSpec now declares
  numbered Surface Transcripts", applied to the bound's refusal. The gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and
  ADR-0167; ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  consistency; ADR-0166, ADR-0178 and ADR-0182 bind each Task commit. ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0165 cites ADR-0153 but decides how a blocking review after archive parks publication, and ADR-0192 cites ADR-0178 but decides how a conflict confined to declared derived paths is resolved; ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot; this Spec changes none of them, so none applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths; express maintainer
  authorization: "considere autorizado a ajustar todas as skills se
  necessário", and the maintainer's answer of 2026-10-04 authorizing this
  Spec; bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0222-a-review-through-claude-that-keeps-its-verdict/_authorization.md`.

## System Architecture

No new package, command or flag. Two existing components change:

- **ACPX Runner** (`internal/agent`): `RunPrompt` decides what a non-zero
  exit after a parsed result means. It gains one case, the refused
  permission, before the transport-anomaly case.
- **Pre-PR Review Command** (`internal/cli`): `runConfiguredReviewSession`
  assembles and bounds the prompt, and `classifyReviewCommandResult` turns
  the runner's result into the review record. The first gains the Claude
  token bound and reads the Spec context before the readiness probe; the
  second gains the refused-permission and prompt-too-long cases.

The daemon's Task and review Batches use read-write sessions, so their
handling of an acpx exit does not change.

## Implementation Design

### Interfaces

```go
// internal/agent: one field on the existing result.
type ExecuteResult struct {
	// ...existing fields...
	TransportAnomaly  string
	PermissionRefused bool // read-only turn ended end_turn, acpx exited 5
}

// internal/cli: two optional record fields and the bound.
type reviewRecord struct {
	// ...existing fields...
	PermissionRefused     bool `json:"permissionRefused,omitempty"`
	EstimatedPromptTokens int  `json:"estimatedPromptTokens,omitempty"`
}

func estimateReviewPromptTokens(prompt string) int // ceil(len(prompt) / 2)
func checkClaudeReviewPromptBound(prompt string, omitted []reviewOmittedPath) (int, error)
```

### The refused permission

In `RunPrompt`, after `cmd.Wait` returns exit code `5`, when the stream parsed
a prompt result whose stop reason is `end_turn`, the request's access is
read-only and the request is not inert, the runner publishes the existing
`acpxPermissionDeniedStatus` Run Event, sets `PermissionRefused`, leaves
`TransportAnomaly` empty and returns the result with no error. Any other
combination follows today's code: a parsed result with any exit but `130`
becomes a transport anomaly, and no parsed result maps the exit code as
before.

`classifyReviewCommandResult` reads `result.PermissionRefused` into
`record.PermissionRefused` and classifies the final answer exactly as for a
clean exit. When the record still ends `blocked`, its reason gains the suffix
` (after the read-only session refused a permission request)`.

### The prompt-too-long reason

Before the runtime-failure and transport-anomaly cases, and before verdict
classification, `classifyReviewCommandResult` looks for the first line of
`result.Answer()` that starts with `Prompt is too long` after trimming spaces.
When one exists the record blocks with the reason
`review prompt too long: <line>`, the line cut to at most 512 bytes on a rune
boundary. Admission, timeout and Spec-read failures keep their precedence,
because no answer exists for them.

### The Claude prompt bound

`runConfiguredReviewSession` reads the Spec context right after the diff,
before the readiness probe, and assembles the final prompt, round one's or
round two's, with its Spec context. For the `claude` provider it computes
`estimateReviewPromptTokens` of that prompt and records it; above 500,000
(`claudeReviewContextTokens / 2`, with `claudeReviewContextTokens =
1000000`) it returns a `reviewPrePromptError`, so no readiness probe, Agent
Session, prompt or fallback follows. The `claude` provider no longer applies
the 917,504-byte diff bound. For `codex` nothing changes: the diff bound runs
where it runs today, and no estimate is recorded.

### Data Models

The review record gains `permissionRefused` and `estimatedPromptTokens`, both
omitted when zero; older records stay readable. `review dispose` reads the
record's outcome and findings as today, so a findings record produced after a
refusal is disposable without change to the dispose command.

### API Contracts

1. API Contract: refused-permission result. `ACPXRunner.RunPrompt` with
   read-only access, a parsed `end_turn` result and acpx exit `5` returns a
   nil error, `PermissionRefused == true`, an empty `TransportAnomaly`, and
   publishes `acpxPermissionDeniedStatus`. With exit `1`, with a stop reason
   other than `end_turn`, with read-write access or with an inert request, the
   result is exactly today's.
2. API Contract: review after a refused permission. The record carries
   `"permissionRefused": true` and the outcome its answer's verdict decides:
   `findings` with ids `F1..Fn` and exit `1`, or `reviewed` and exit `0`; a
   `blocked` outcome's reason ends with
   ` (after the read-only session refused a permission request)`.
3. API Contract: prompt-too-long reason. An answer with a line starting
   `Prompt is too long` blocks with exit `2` and the reason
   `review prompt too long: <line>`, whether the runner returned an error or a
   parsed result.
4. API Contract: Claude prompt bound. For `claude`, a prompt whose estimate
   exceeds 500,000 tokens blocks with exit `2`, no runner call of any kind, and
   the reason
   `review prompt too large: <n> estimated tokens after omitting <m> path(s) exceeds the claude review budget of 500000 tokens, half of its 1000000-token context window`;
   the record carries `estimatedPromptTokens`. For `codex`, the 917,504-byte
   diff bound and its reason are unchanged.

### Surface Transcripts

1. Surface Transcript: a Claude review whose prompt exceeds the token budget,
   run in a disposable repository whose Project Config selects `claude` and
   whose candidate adds a file of 1,100,000 bytes, with a disposable home.

   ```transcript
   $ roundfix review --base main
   stdout:
   ...
   stderr:
   roundfix: review blocked: review prompt too large: <n> estimated tokens after omitting <m> path(s) exceeds the claude review budget of 500000 tokens, half of its 1000000-token context window
   exit: 2
   ```

## Coverage Map

- Goal 1 → The refused permission; API Contracts 1 and 2.
- Goal 2 → The prompt-too-long reason; API Contract 3.
- Goal 3 → The Claude prompt bound; API Contract 4; Surface Transcript 1.
- Goal 4 → The refused permission (read-write and inert unchanged); The Claude
  prompt bound (codex unchanged).
- User Story 1 → API Contracts 1 and 2.
- User Story 2 → API Contract 2 (`permissionRefused`).
- User Story 3 → API Contract 3.
- User Story 4 → API Contract 4; Surface Transcript 1.
- Core Feature 1 → API Contract 1.
- Core Feature 2 → API Contract 2.
- Core Feature 3 → API Contract 3.
- Core Feature 4 → API Contract 4.
- Core Feature 5 → Build Order 1.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 3.

## Integration Points

- acpx 0.19.4 exit codes: `5` is `PERMISSION_DENIED`, set when a turn's
  permission requests were all denied or cancelled and none approved. The
  runner reads only the exit code and the parsed stop reason.
- Claude adapter 0.85.0 behind acpx: reports a 1,000,000-token `size` for
  `claude-opus-5-5`; Roundfix declares that window as a constant rather than
  reading it, because the adapter reports it only after a prompt.

Research record. The Secondbrain was read first (`wiki/index.md`, then a
`qmd` query on the defect): the two triaged reports at
`inbox/roundfix/_triaged/2026-10-02-review-claude-bloqueia-quando-o-revisor-pede-terminal.md`
and
`inbox/roundfix/_triaged/2026-10-02-review-claude-perde-achados-e-estoura-o-prompt-no-oraculum.md`
gave the symptoms and the answer text, and the query surfaced ADR-0020 and
Spec 0153's decision that every exit after a parsed result blocks the review,
which this design narrows rather than reverses. Exa found Anthropic's Agent
SDK permissions guide, which changed the plan: the preferred deny-up-front
option needs `disallowedTools` or `dontAsk`, and acpx 0.19.4 forwards neither,
so the design classifies the exit instead. Exa also found Anthropic's
token-counting guide, whose tokenizer note set the two-bytes-per-token
estimate together with a local measurement of 2.66 bytes per token on the
older tokenizer. The digest is captured as a pending research entry in the
Secondbrain inbox.

## Testing Approach

1. Runner seam: `runFakeACPXPrompt` in `internal/agent/acpx_runner_test.go`
   drives the real `RunPrompt` against a fake acpx process. The harness gains
   the request's access and inert flag; new cases print a parsed `end_turn`
   result and exit `5` (read-only, read-write, inert) and a parsed
   `max_tokens` result with exit `5`, and assert API Contract 1 including the
   Run Event.
2. Review seam: `newReviewCommandFixture` with the `claude` provider and the
   fake `reviewCommandRunner`, in a new `internal/cli/review_permission_test.go`.
   Cases: a refused permission with a findings answer, then `review dispose`
   of `F1`; with `No findings.`; with an answer that has no verdict; and an
   answer starting `Prompt is too long`, once with a `BatchFailureError` and
   once with a parsed result. The existing
   `TestReviewCommandBlocksOnTransportAnomaly` stays unchanged and green.
3. Bound seam: the same fixture in a new
   `internal/cli/review_prompt_bound_test.go`, with candidates sized against
   the constants, not literals: a Claude prompt over the budget (no runner
   call, Surface Transcript 1's stderr and exit), a Claude candidate whose
   diff exceeds `reviewDiffBound` but whose estimate fits (one prepared
   prompt), and a Codex candidate with that diff (refused by the byte bound
   as today). The existing bound tests in `review_scope_test.go` use `codex`
   and stay unchanged.

## Build Order

1. The Roundfix Skill's `review` reference, its mirror and the `review`
   command guide describe the four contracts; the skill version is raised and
   recorded.
2. The ACPX Runner keeps a read-only turn that ended after a refused
   permission (API Contract 1).
3. The review classifies that answer by its verdict and names a
   prompt-too-long answer (API Contracts 2 and 3) (depends on: 1, 2).
4. The review bounds a Claude prompt in estimated tokens (API Contract 4)
   (depends on: 3, which shares `internal/cli/review.go`).
5. The final QA gate (depends on: 4).

## Risks & Considerations

- A reviewer that needed the refused command may give a weaker verdict. The
  record names the refusal, and the review prompt already tells the reviewer
  the session is read-only.
- Two bytes per token overestimates Go and Markdown on the newer tokenizer by
  a little and underestimates dense non-ASCII text; the half-window reserve
  absorbs the error, and a residual overflow is now reported by name.
- Reading the Spec context before the readiness probe changes when a Spec
  read failure is reported, not its reason.

## Decisions

- Classify the permission exit only for a read-only `end_turn` turn; see
  ADR-0227.
- Deny-up-front and a prompt instruction were rejected on measured evidence;
  see ADR-0227.
- Bound Claude in estimated tokens at half its window over the whole prompt,
  and keep Codex's byte bound; see ADR-0227.
- Read the prompt-too-long reason from the final answer, not stderr, because
  the runtime writes it as the Agent's message.
