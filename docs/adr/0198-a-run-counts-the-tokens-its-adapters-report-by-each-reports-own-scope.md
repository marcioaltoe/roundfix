---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Run counts the tokens its adapters report, by each report's own scope

Roundfix recorded nothing about what an Agent Session consumed, and
`roundfix deliver status` said `spend not measured`. A local metering gateway
counted tokens per request for a week, then failed 500 of 5,426 requests with
a server error and was taken out of the Run path. The ACP adapters already report
usage on the stream Roundfix reads, but not with one meaning:

- The Agent Client Protocol stabilizes `usage_update` as session state: `used`
  is the tokens currently in context, and the optional `cost` is the
  cumulative session cost. The token usage of a completed turn, the `usage`
  object on the prompt response, is still a draft that leaves open whether it
  covers the turn or the session.
- `codex-acp` fills the prompt response's `usage` with its last model
  request only, and its `usage_update.used` with that request's total. On
  2026-09-25 to 2026-09-28 the sum of a turn's `used` readings was 92.5% to
  100% of what the gateway counted on the wire for the same session, while the
  reported `usage` was 2% to 3% of it.
- `claude-agent-acp` fills the prompt response's `usage` with the whole turn,
  sends `usage_update` readings while a response streams, and puts the
  cumulative session cost in `usage_update.cost`.

Roundfix records what the adapter reported for every prompt, and counts it
under these rules:

- **The adapter lineage names the scope.** An adapter lineage whose prompt
  response covers only its last model request counts a turn as the sum of the
  turn's `used` readings, basis `request-sum`, and records no input, output or
  cache split, because none was reported for the turn. Every other lineage
  counts the prompt response's `usage.totalTokens`, basis `turn`, with the
  split it reported. If a `request-sum` adapter's reported total ever exceeds
  its last reading, the report covers more than one request and is counted as
  `turn`.
- **Absent stays absent.** A prompt whose adapter reported neither is recorded
  as unreported. No surface prints it as zero, and every total says how many
  prompts it leaves out.
- **Spend is what the adapter said.** A reported `cost` is kept as the
  adapter's cumulative reading. An Agent Session's spend is the sum of its
  increases, and a lower reading starts a new count. Roundfix keeps no price
  table: the only adapter that reports model identifiers without a cost does
  not report a turn's input, output and cache split, so a price computed from
  its totals would misapply the cache discount.
- **No gateway logic.** Roundfix neither proxies nor meters model traffic. A
  metering gateway stays an operator option described in the user guide.

## Consequences

- Codex turns are counted within a few percent of the wire, without their
  split. Claude turns carry their split and a spend.
- A new adapter lineage is counted as `turn` until measurement shows its
  prompt response covers less.
- A change in an adapter's reporting reaches Roundfix as a changed number,
  not as an error. The runtime guard above catches the one change that would
  undercount a `request-sum` adapter.
- Tokens spent outside a Run, such as the pre-PR review, are not counted.
