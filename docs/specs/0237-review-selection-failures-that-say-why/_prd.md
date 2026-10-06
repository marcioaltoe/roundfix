---
spec: 0237-review-selection-failures-that-say-why
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# Review selection failures that say why

On 2026-10-02 in Pantheon (Roundfix 0.26), `roundfix review --base main` with
the default `codex / gpt-5.6-luna / max` reviewer blocked twice with
`review runtime failure: Agent Selection failed for runtime "codex":
agent/protocol error` and an empty answer file, while the Doctor Command and
`codex exec` answered. Three days later the same review ran. The reason named
neither the protocol step that failed nor the adapter's message, and the
maintainer skipped the review for a Pull Request
([the adopted Backlog Entry](references/2026-10-06-a-review-agent-selection-failure-says-nothing.md)).

A measurement during authoring reproduced the reason byte for byte: the
installed acpx 0.19.4 reports a failure at `initialize`, `session/set_model` or
`session/prompt` as a JSON-RPC error line on stdout with the adapter's message,
an empty stderr and exit `1`, and the ACPX Runner discards that line.

This is a bug fix: the review keeps its contract and its provider policy. It
mints this minimal PRD for the downstream artifact contract; the design lives
in the [_techspec.md](_techspec.md) and ADR-0242.

## Prerequisites

None. Specs 0235 and 0236 are authored in the same cycle; if either also raises
the Roundfix Skill's version, the operator orders the queue so that the later
Spec raises it from the earlier one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  review record gains three optional camel-case fields and finding ids keep
  the `F<n>` sequence. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the review still drives the
  locally installed acpx; no request, credential or forge read is added, and
  tests use fake runners and a fake acpx. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0242 (this Spec) decides that a
  blocked review names the protocol step and the adapter's message and that a
  failure before the prompt is retried once on the same selection. ADR-0020:
  "Without a parsed result, a nonzero exit stays a Batch failure exactly as
  before", which this Spec keeps while it keeps the adapter's message.
  ADR-0227: "Every other exit after a parsed result stays a transport
  anomaly", unchanged. ADR-0050: "once Agent work begins, Roundfix fails the
  Work Item instead of switching models", and ADR-0114: "Roundfix therefore
  treats opening an Agent Session as selection, not work", the boundary the
  retry respects; Runs keep their fallback behavior. ADR-0153: "Pre-PR review
  is an explicit provider policy", so no failure selects `none`. ADR-0151:
  "Configured review profiles select the independent reviewer", and the
  profile and its fallback chain stay as they are. ADR-0174: "A pre-PR review
  reads the final message", unchanged for a review that reaches the reviewer,
  and ADR-0197: "A pre-PR reviewer lineage spans at most two rounds", whose
  round-two resume keeps its rule. ADR-0187 splits the Roundfix Skill by
  command and ADR-0189 ties an owned skill's version to its content, so the
  skill edit raises the version. ADR-0184: "A TechSpec now declares numbered
  Surface Transcripts", answered in the TechSpec with the reason none applies.
  ADR-0240: "One QA partial policy, and rows a Run sandbox cannot reach"
  decides when its QA partial qualifies. The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155,
  ADR-0156 and ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and
  ADR-0183 check this Spec's consistency by citation and receipt. ADR-0178
  authorizes each Task commit by its grant, ADR-0182 runs Settlement Checks
  before it, and ADR-0166 records undeclared paths; every Task declares its
  paths. ADR-0169's merge-base diff and ADR-0196's finding validation apply unchanged to a review that reaches the reviewer, and ADR-0233 regenerates the raised skill version at merge. ADR-0069 and ADR-0238 cite ADR-0050 but decide Baseline analysis and the light tier; ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry; ADR-0165 cites ADR-0153 but decides how a blocking review after archive parks publication; ADR-0192 cites ADR-0178 but decides how a conflict confined to declared derived paths is resolved; ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park; ADR-0180 cites ADR-0069 but decides the built-in selections, and ADR-0181 and ADR-0217 cite ADR-0180 but decide where a configuration is compared with the recommendation and how a Cursor selection names its model; ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot; ADR-0237 cites ADR-0229 but decides how a Delivery Retry records a merge made outside the queue; this Spec changes none of them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths, and the maintainer authorized skill
  edits ("considere autorizado a ajustar todas as skills se necessário") and,
  on 2026-10-06, the Governed Paths this Spec declares ("Concedo"). No other
  Governed Path changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0237-review-selection-failures-that-say-why/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A blocked review names the protocol step that failed and the adapter's own
  message, bounded and without request data.
- A failure during Agent Selection, before the review prompt was sent, is
  retried once on the same selection, and the record says so.
- A failure after the prompt was sent is never retried, so no candidate is
  reviewed twice, and no failure ever selects `none`.

## User Stories

1. As an operator whose review blocked, I want the reason to say which
   protocol step failed and what the adapter said, so that I can tell a
   refused model from a failed turn and decide whether to rerun.
2. As an operator, I want a transient failure before the prompt to be retried
   once automatically on the reviewer my repository chose, so that I do not
   skip a required review over a hiccup.
3. As an operator, I want the record to list every automatic retry with its
   step and message, so that a review that passed after a retry is not
   mistaken for one that never failed.

## Core Features

1. **The runner keeps the adapter's message.** When an `acpx prompt` process
   exits non-zero without a parsed result, the ACPX Runner reads the first
   JSON-RPC error line on stdout, names the step by the request it answers,
   keeps its `message` to one line of at most 512 bytes, and records whether a
   `session/prompt` request had been sent (ADR-0242).
2. **The review names the step.** A blocked review's reason carries the step
   and the message, and its record carries `failedStep` and `adapterMessage`.
3. **One retry before the prompt.** A selection failure while preparing the
   Agent Session, or a prompt process that failed before it sent
   `session/prompt`, is retried once on the same selection; the record lists
   it under `selectionRetries`. A configured fallback activates only after the
   retry also failed.
4. **The skill and the guide say so.** The Roundfix Skill's `review` reference
   and the `review` command guide describe the step, the message, the retry
   and the record fields.

## Non-Goals / Out of Scope

- Changing how Runs handle a selection failure: they keep activating the next
  fallback without a same-selection retry (ADR-0050, ADR-0114).
- Retrying a failure after the prompt was sent, or deciding by the adapter's
  wording which failures are transient.
- Changing the exit-zero paths, the transport anomaly after a parsed result,
  or the permission-refusal classification of ADR-0227.
- Reading the JSON-RPC error's `data`, the request parameters or the prompt
  into any reason or record.
- Explaining why the Pantheon review failed on 2026-10-02; the evidence is
  gone, and this Spec makes the next failure say why.

## Success Metrics

1. Success Metric: a fake acpx that prints the stdout recorded from acpx 0.19.4
   for a failure at `initialize`, at `session/set_model` and at
   `session/prompt`, each with exit `1` and an empty stderr, yields three
   different errors naming their step and the adapter's message; before this
   Spec all three read `Agent Selection failed for runtime "codex":
   agent/protocol error`.
2. Success Metric: the first two report that no prompt was sent, and the third
   and an adapter disconnect after `session/prompt` report that it was.
3. Success Metric: a review whose session preparation fails once with a
   selection failure, or whose prompt process fails once before
   `session/prompt`, retries the same selection, exits `0` on a clean answer,
   and records one retry; a failure after `session/prompt` blocks after one
   prompt call with `failedStep` and `adapterMessage` set; two failures before
   the prompt activate a configured fallback, or block with the retry noted.
4. Success Metric: every existing review and runner test passes, and the only
   changed expectations are the three fallback tests that now need two
   failures before the fallback.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Pantheon report of 2026-10-05, now in the Secondbrain at
  `inbox/roundfix/_triaged/2026-10-05-revisao-pre-pr-com-codex-falha-por-protocolo-de-forma-intermitente.md`,
  with the incident's reason and the Doctor Command's adapter line
  (`codex-acp` 2.1.1).
- A measurement on 2026-10-06 with the installed acpx 0.19.4 and a fake ACP
  adapter that fails one chosen method, under a temporary `HOME`: every
  failure inside `acpx prompt --format json` (at `initialize`,
  `session/set_model` and `session/prompt`) printed a JSON-RPC error line on
  stdout naming the failing request's id with the adapter's message, an empty
  stderr and exit `1`; an adapter crash during `session/prompt` printed an
  error with a null id and `AGENT_DISCONNECTED`; `acpx sessions ensure`,
  `set-mode` and `set` printed the adapter's message on stderr and exited `1`.
  The same stdout fed to the current ACPX Runner through `go test -overlay`
  produced the incident's reason for all three steps.
- The Agent Client Protocol overview
  (<https://agentclientprotocol.com/protocol/v1/overview>): the phases are
  initialization (`initialize`, `authenticate`), session setup (`session/new`
  or `session/load`) and the prompt turn (`session/prompt`), and errors carry
  an `error` object with `code` and `message`.

## Decisions

- Name the step by the request the JSON-RPC error answers, and keep only its
  `message`; see ADR-0242.
- Retry once on the same selection only before the prompt was sent, and only
  then try a fallback; see ADR-0242.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
