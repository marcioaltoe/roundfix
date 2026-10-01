---
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: archived
created: 2026-09-30
surfaces: [backend, cli, data, docs]
archived: "2026-10-01"
source_slug: 0204-a-run-that-reports-the-tokens-and-spend-it-used
qa_override: true
qa_override_approval: 'maintainer standing approval of 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'QA partial with environment rows only: rows 04, 06, 07, 09 and 10 need the built binary under a disposable Home with a seeded Run Database, which the QA Agent''s instructions forbid; their in-process equivalents pass exact text, schema, totals and read-only checks (binary-equivalents.log), and the operator ran the unseeded journeys with the built binary (qa/evidence/2026-10-01-operator/binary-journeys.txt). Row 15 is the pre-PR Pull Request row.'
qa_override_qa_outcome: missing
qa_override_qa_task_status: pending
qa_override_revision: 9000a443d784f279493e5946052213f6b53f0f56
---


# A Run that reports the tokens and spend it used

Roundfix cannot say what a Run consumed. `roundfix deliver status` ends its
limits line with `spend not measured`, and a Delivery Queue can be bounded only
by a deadline and a retry allowance. The questions a maintainer asks after a
wave cannot be answered from the Run Database: how many tokens a Spec took,
which Task took most of them, and whether a queue should stop before it drains
a quota window.

- For a week the only measurement was a local metering gateway on the Codex
  traffic. It logged tokens per request, outside Roundfix, and Runs were
  attributed to it by time window. Between 2026-09-25 and 2026-09-29 it
  failed 500 of 5,426 requests with a server error, killed Agent Sessions in the
  middle of Tasks, and was taken out of the Run path.
- The ACP adapters Roundfix drives already report usage on the stream it
  reads, and Roundfix drops it.
- A backlog entry of 2026-08-08, deferred on 2026-08-26, asked for the same
  record per Agent Session and left open whether the adapters report usage at
  all. This Spec answers that question with measurement.

This Spec records, for every prompt of every Agent Session a Run opens, the
tokens and cost its adapter reported. It prints them per Task and per Run, per
queue in `deliver status`, and adds an optional token ceiling to a Delivery
Queue.

## Prerequisites

This Spec is delivered after Spec 0194, which splits the Roundfix Skill and
the command reference into one file per command. Its Tasks edit the per-command
files for `deliver`, `events`, `implement` and `runs` that Spec 0194 creates,
and its Verification reads them. The Delivery Queue does not enforce this
order, so the operator queues Spec 0194 first.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; usage rows are
  keyed by the existing Run ID, Work Item scope and Agent Session name.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Roundfix reads the JSON-RPC lines
  its ACP adapters already print and the Run Database; no credential is read,
  no model traffic is proxied and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0198 (this Spec) decides what a
  token count means per adapter lineage, keeps an absent report absent,
  records only adapter-reported cost and leaves metering gateways out of
  Roundfix. ADR-0199 (this Spec) makes the token ceiling a Delivery Queue
  Limit that parks new starts and never stops a running Task. ADR-0017 keeps
  acpx as the only agent layer, so usage is read from its JSON-RPC output.
  ADR-0020 keeps the parsed prompt result authoritative over the acpx exit
  code, and usage is read from that same result. ADR-0051 gives each Task and
  the QA gate its own Agent Session, which is the scope a count is recorded
  under. ADR-0008 stores raw producer JSON in Agent Run Events; the usage
  Run Event is Daemon-owned and defines its own small payload. ADR-0033 prunes
  the Run Event Journal by retention, which is why usage is kept in its own
  table that retention never deletes. ADR-0098 appends Run Events in batches;
  the usage event takes the ordinary batched path. ADR-0158 and ADR-0164 keep
  the Run Budget time-based and are unchanged. ADR-0166 records a Task's
  undeclared paths, and every Task here declares its paths. ADR-0167 keeps
  the pre-PR Pull Request row from deciding a qualifying partial; this Spec's
  gate aims at `pass`. ADR-0176 reads citations only from authored text.
  ADR-0182 runs Settlement Checks before each Task commit, so each Task keeps
  the repository Verification green on its own. This Spec's gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and
  ADR-0156, and ADR-0093 and ADR-0094 check its consistency by citation and
  artifact presence. ADR-0184 has the TechSpec state each changed command as
  a Surface Transcript. ADR-0097 cites ADR-0080 but carries a QA row forward,
  ADR-0168 cites ADR-0093 but narrows the related-ADR check, and ADR-0183
  cites ADR-0093 but adds receipts for attributed claims; this Spec changes
  none of them, so they do not apply. ADR-0194 and ADR-0195 cite ADR-0097 but
  decide how a QA row's evidence is snapshotted and when it is re-observed,
  which this Spec does not touch. ADR-0030 cites ADR-0008 but makes per-Batch
  agent log files opt-in; usage is read from the stream, not from those files.
  ADR-0170 cites ADR-0158 but decides what Task Carry-Forward skips, and
  ADR-0171 cites ADR-0033 but decides what a retention prune reports, and
  ADR-0172 cites ADR-0171 but decides how Doctor reports reclaimable storage;
  none of the three changes, and the new tables are outside retention. ADR-0197 cites ADR-0051
  but concerns the pre-PR reviewer's Agent Session, whose tokens this Spec
  does not count. None of these seven applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se necessário"),
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/runs.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- Every prompt a Run sends through an Agent Session leaves a usage record in
  the Run Database: the tokens its adapter reported, counted by the report's
  own scope, and the cost the adapter reported, or an explicit unreported
  mark.
- A maintainer reads the tokens and reported cost of each Task, each Run and
  each Delivery Queue from the commands they already use.
- An operator can bound a Delivery Queue by tokens, and the bound stops new
  work without discarding work in progress.
- No number Roundfix prints claims more precision than the adapter gave it.

## User Stories

1. As a maintainer reviewing a finished Run, I want its tokens per Task and in
   total, with the reported cost where the adapter gave one, so that I know
   which Task cost what.
2. As a Supervisor following a Run through its event stream, I want one
   record per prompt with its tokens, so that I can watch consumption while
   the Run is active.
3. As an operator of a Delivery Queue, I want `deliver status` to show the
   tokens the queue's Runs used instead of `spend not measured`, so that I can
   judge a queue while it runs.
4. As an operator starting an unattended queue, I want a token ceiling that
   parks the remaining items once the queue has used that many tokens, so that
   one queue cannot drain a quota window, and I want the Run in progress left
   to finish its Tasks.
5. As a maintainer who wants per-request figures, I want the user guide to
   describe the metering gateway as an operator option with its known
   failure, so that I can choose it knowingly.

## Core Features

1. **Adapter usage is read, not dropped.** For each prompt, Roundfix reads the
   `usage` object of the prompt response and every `usage_update` notification
   the adapter sends during the prompt, including a prompt that fails or is
   stopped after its stream was read.
2. **A turn is counted by its report's scope.** For an adapter lineage whose
   prompt response covers only its last model request (`codex-acp`), the
   turn's tokens are the sum of its context readings, basis `request-sum`,
   with no split. For every other lineage, the turn's tokens are the prompt
   response's total, basis `turn`, with the input, output, cache and thought
   split it reported. A `request-sum` report whose total exceeds its last
   reading is counted as `turn`.
3. **An absent report stays absent.** A prompt with neither report is recorded
   as unreported. Every total names how many prompts it leaves out, and no
   surface prints an unreported prompt as zero.
4. **Reported cost is kept as the adapter said it.** A `usage_update` cost is
   recorded with its currency. An Agent Session's cost is the sum of the
   increases of its cumulative readings, and a lower reading starts a new
   count. Roundfix computes no price.
5. **Usage is recorded in the Run Database per prompt.** Each record names the
   Run, the Work Item scope (Task, QA gate or review Batch), the Agent Session,
   the Agent Selection that ran it, the basis, the tokens and the cost. The
   records live in their own table, which Journal Retention never prunes. A
   failed write never changes the prompt's outcome and is reported as a
   warning.
6. **The Run Event Stream carries one `usage` record per prompt.** A new
   stream category, `usage`, is on by default and selectable with `--filter`.
7. **`roundfix runs show <run-id>` prints a Run's usage.** One line per Work
   Item scope and a total, as text or `--json`. It is read-only.
8. **The Implement Run summary ends with a `Tokens:` line.** It prints the
   Run's total in the same words as `runs show`.
9. **`deliver status` prints a `Usage:` line.** It reports the tokens and cost
   of every Run the queue recorded for its items, and the limits line names
   the token ceiling instead of `spend not measured`.
10. **`deliver start --max-tokens <n>` sets a token ceiling.** At or above it,
    a queued item parks as `queue-token-ceiling` without a worktree, a retry
    is refused, and the Pending Question tells the operator to record a new
    queue. An item past `queued`, and a Run in progress, continue.
11. **The user guide describes what is measured and the gateway option.** It
    states each adapter's basis, the ceiling's boundary, and how an operator
    can route Codex through a metering gateway, with the failure measured on
    2026-09-29 and the transport caveat of the current Node.js.

## Declared breaks

- `deliver start` and `deliver status` print `tokens <n>` or `tokens none`
  where the limits line printed `spend not measured`, and `deliver status`
  prints one more line.
- The Implement Run summary prints one more line after the outcome line.
- `roundfix events` prints `usage` records by default. A consumer that
  rejects an unknown category must pass `--filter`.
- A Run Database written by this Spec has a newer schema, and an older
  Roundfix refuses it and names `upgrade`, as for every schema change.
- A Delivery Queue item can park with the new blocker `queue-token-ceiling`.

## Non-Goals / Out of Scope

- A price table, or any cost Roundfix computes from token counts.
- Proxying, metering or rewriting model traffic, or any code for a metering
  gateway.
- Counting tokens spent outside a Run: the pre-PR review, the Doctor's and
  `profiles configure`'s disposable Agent Sessions, and any prompt a person
  sends by hand.
- Stopping or cancelling a Run, or a Task, because of tokens. The Run Budget
  stays time-based.
- A per-Run or per-Task token ceiling.
- Showing usage in the Live Run View.
- A per-model split of a turn, even where the adapter reports one.
- Backfilling Runs that finished before this Spec.

## Success Metrics

1. Success Metric: fed the recorded Codex prompt fixture, whose final
   reported total equals its last context reading, Roundfix records the sum of
   the readings with basis `request-sum` and no split; fed the recorded Claude
   fixture, it records the reported turn total with basis `turn` and its
   split.
2. Success Metric: a prompt whose adapter reports nothing is recorded as
   unreported, and `runs show`, the Implement Run summary and `deliver status`
   each name it as unreported and never print it as `0`.
3. Success Metric: for a Run of two Tasks and a QA gate, `runs show` prints
   each scope's tokens and a total equal to their sum, and the Implement Run
   summary prints the same total.
4. Success Metric: `deliver status` no longer prints `spend not measured`, and
   its `Usage:` line sums every Run the queue recorded, including a Run of an
   earlier retry.
5. Success Metric: with `--max-tokens` below the queue's recorded tokens, the
   next queued item parks as `queue-token-ceiling` without a worktree, a retry
   is refused, and the item whose Run was in progress is not stopped; below
   the ceiling, items start as before.
6. Success Metric: a failed usage write leaves the prompt's result, the Task's
   status and the Run's outcome unchanged.

## Recorded limits

- A Codex turn is counted without its input, output and cache split, because
  the adapter reports the split only for its last request.
- The Codex count measured 92.5% to 100% of the gateway's wire count. The
  missing requests are ones the adapter does not report.
- A Run in progress counts toward `deliver status` and the ceiling only once
  the queue records it, which happens when the Run ends.
- A queue can end above its ceiling by up to one item's Run.
- An adapter that reports nothing never reaches the ceiling.

## Decisions

- **What a count means, and what spend is.** See ADR-0198.
- **Where the ceiling acts.** See ADR-0199.
- **Recorded per prompt, summed on read.** A total is always the sum of the
  records a reader can list, so the three surfaces cannot disagree.
- **Its own table.** Usage outlives Journal Retention, because the
  maintainer's questions come after a wave, not during it.
- **The gateway stays outside.** The maintainer decided Roundfix must not
  reimplement the gateway's logic, and that the gateway may be used only as an
  operator's choice.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- **The protocol, published.** The Agent Client Protocol's prompt-turn page
  defines `usage_update` with `used`, `size` and an optional cumulative `cost`
  (<https://agentclientprotocol.com/protocol/v1/prompt-turn>), and its
  announcement marks it stable
  (<https://agentclientprotocol.com/announcements/session-usage-stabilized>).
  The End-Turn Token Usage RFD keeps the prompt response's `usage` in draft
  and leaves open whether it is per turn or cumulative
  (<https://agentclientprotocol.com/rfds/end-turn-token-usage>). All three
  read 2026-09-30.
- **The adapters, published.** The npm package `@agentclientprotocol/codex-acp`
  2.0.1 builds the prompt response's `usage` and its `_meta.quota` from its
  last token count, and emits `usage_update` with `used` equal to that count's
  total and no `cost`. The npm package `@agentclientprotocol/claude-agent-acp`
  0.84.0 accumulates the turn's usage into the prompt response, sends
  `usage_update` while a response streams, and attaches `cost` with the
  SDK's cumulative cost in USD. Both identify models in `_meta.quota.model_usage`.
  Read from the installed packages on 2026-09-30.
- **A measurement this Spec did not design.** The acpx session streams on the
  maintainer's machine hold 22,931 Codex `usage_update` readings and 602 Codex
  prompt responses from 2026-09-25 to 2026-09-30, and 43 Claude prompt
  responses with `usage` and cost from August 2026. In 598 Codex prompt
  responses the reported total equals the turn's last reading, and the other
  four carry `usage: null`. For eleven Codex
  sessions that ran alone while the metering gateway logged every request,
  the sum of their readings was between 92.5% and 100% of the gateway's
  input-plus-output count, one of them exactly equal, while the reported
  `usage` was 2% to 3% of it. In Claude sessions the reported cost grew across
  the turns of one session, as a cumulative reading does.
- **The gateway transport, published.** Node.js 26 bundles undici 8
  (<https://nodejs.org/en/blog/release/v26.0.0/>), and undici 8 offers HTTP/2
  by default (<https://github.com/nodejs/undici/pull/4828>). Opting out takes
  a dispatcher with `allowH2: false`
  (<https://undici.nodejs.org/api/Connector>); Node.js offers no flag for its
  built-in `fetch`. The gateway forwards with the built-in `fetch`, and 499 of
  its 500 failures were its own `502` for an upstream it could not reach.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the query `qmd query
"medir tokens e custo por tarefa de agente ACP usage_update PromptResponse
usage teto de tokens" --all --files --min-score 0.3`. It returned
`wiki/concepts/custos-e-limites-de-agentes-de-codigo.md`, which asks any agent
graph to state a token or call limit next to its concurrency and retry policy
and to keep budget for verification; it supports the ceiling and its
item-start boundary. It also returned this repository's deferred backlog entry
of 2026-08-08, whose rule that a missing measurement must stay observably
absent became Core Feature 3. `wiki/sources/digest-roundfix-modelos-custos-e-disponibilidade-2026-09-08.md`
separates API prices from subscription access; with the adapters' missing
split, it is why no price table is computed. Exa found the ACP pages, the
Node.js release note and the undici sources cited above, which changed the
design twice: the draft status of the prompt response's `usage` is why the
basis is recorded per adapter lineage, and the undici default is why the
gateway stays an operator option.

## Technical candidate

The [_techspec.md](_techspec.md) records the counting rule, the data shapes,
the command transcripts, coverage and build order.
