---
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: active
created: 2026-10-01
surfaces: [backend, data, docs]
---

# The Jev Router on docs and chore Tasks

Roundfix runs every Task of a work category on the same model and reasoning
effort, even where most of the work is routine, as in the `docs` and `chore`
categories. OpenRouter publishes `typesafe/jev-router`, which uses Jev to
pick a model and an effort for each request. If it keeps quality on routine
work while it picks cheaper models or lower effort, those categories cost less
per Task.

The adopted Backlog Entry
([references/2026-10-01-try-the-jev-router-for-cheap-agent-categories.md](references/2026-10-01-try-the-jev-router-for-cheap-agent-categories.md))
proposed the experiment: OpenCode with an OpenRouter provider whose model is
the router, selected only for `docs` and `chore`, on a key of Roundfix's own.
The authoring session measured, on 2026-10-01, without a key and without a
call that spends:

- OpenRouter's model listing names `typesafe/jev-router` ("TypeSafe: Jev
  Router"), with a 1,000,000-token context, `tools`, `tool_choice`,
  `parallel_tool_calls` and `reasoning_effort` among its supported
  parameters, and its prompt and completion prices given as `-1`, variable.
  Its endpoint listing is empty. Its model page says it "picks the best model
  and reasoning effort for each request" and "adapts as your conversation
  evolves", and its FAQ calls the router's own pricing zero.
- TypeSafe's documentation index (<https://docs.typesafe.ai/llms.txt>) has
  no page for the router, so its contract is OpenRouter's chat completions
  API.
- The installed OpenCode `1.18.34` accepts an inline configuration in
  `OPENCODE_CONFIG_CONTENT`. A custom provider `roundfix-openrouter` there,
  using `@ai-sdk/openai-compatible` with base URL
  `https://openrouter.ai/api/v1`, key `{env:ROUNDFIX_OPENROUTER_API_KEY}` and
  model `typesafe/jev-router`, is listed by `opencode models
  roundfix-openrouter` as `roundfix-openrouter/typesafe/jev-router`, and
  `opencode debug config` masks the key. This ran with the key unset and
  model fetching disabled.
- OpenRouter's key endpoint, `GET https://openrouter.ai/api/v1/key`, reports
  the key's `usage_monthly` in US dollars for the current UTC month, and a
  per-key `limit` the account owner can set.

On 2026-10-01 the maintainer allowed agent prompts that carry this
repository's code and diffs, and only those, to go to OpenRouter on
`ROUNDFIX_OPENROUTER_API_KEY`, under the US$5 monthly ceiling shared with
all Jev use, of which about US$0.40 was spent. This Spec makes the router an
opt-in OpenCode selection that only a repository's Project Config can name,
gates every routed prompt on that ceiling, records what each one cost, and
measures the router on four replayed `docs` and `chore` Tasks against the
current default. Nothing becomes a default.

## Prerequisites

This Spec is delivered after Spec 0205, which adds the Judge Log, its
schema `roundfix/judge-log/v1`, its path under Roundfix Home, the judge
package's loaded US$5 ceiling and the transport that reads
`ROUNDFIX_OPENROUTER_API_KEY`. Task 02 reads and appends that log through
the judge package, and Task 03 reads its ceiling. It is also delivered after
Spec 0217, because both change the Roundfix Skill's runtime reference and
raise the skill's version. The Delivery Queue does not enforce this order, so
the operator queues Spec 0205 and Spec 0217 first.

## Project Constraints

- Identifier strategy: not applicable — the routed selection is named by the
  fixed OpenCode model value `roundfix-openrouter/typesafe/jev-router`; a
  router line in the Judge Log reuses the judgment kind vocabulary with the
  fixed kind `router-prompt`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — Roundfix sends one HTTPS
  `GET https://openrouter.ai/api/v1/key` with
  `Authorization: Bearer <ROUNDFIX_OPENROUTER_API_KEY>` before and after each
  routed prompt, and nothing else. OpenCode sends the prompts to
  `https://openrouter.ai/api/v1` with the same key, which it reads from its
  own environment through the `{env:ROUNDFIX_OPENROUTER_API_KEY}` placeholder.
  Roundfix reads that variable only to check that it is set and to send it in
  the key endpoint's authorization header; it never writes the value into an
  argument, a file, a log, a Run Event or the inline configuration, and it
  never reads `OPENROUTER_API_KEY`. The prompts carry the repository's own
  code and diffs, so only a repository's committed Project Config can select
  the router. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0218 (this Spec) decides that the
  router is reached only as one OpenCode selection with an inline provider,
  that only Project Config can select it, and that each routed prompt runs
  under the shared Jev ceiling and is recorded in the Judge Log. ADR-0201:
  "a spending ceiling of US$5 per calendar month", which this Spec applies to
  the router as well, read from the judge's own question file. ADR-0200 adds
  the advisory judge, and this Spec does not change what it judges. ADR-0050:
  "once Agent work begins, Roundfix fails the Work Item instead of switching
  models over potentially modified state", which decides what a ceiling
  refusal does after work began. ADR-0114: "A Fallback Selection may switch
  ACP Runtime automatically only while Agent work has not begun". ADR-0017:
  "Roundfix replaces its hand-rolled ACP client layer", making acpx the only
  agent layer, so the router is reached through acpx's `opencode` agent. ADR-0037: "passes both explicitly for every Agent
  Session, and fails Preflight Validation when the runtime does not support
  them", which the key check extends to a routed selection. ADR-0049: "every
  configured tuple must be proven through the installed ACP adapter rather
  than accepted from a static compatibility assumption", so the routed tuple
  is proved like any other. ADR-0105 lets the projection read OpenCode's large
  catalog, which now includes the routed model. ADR-0108 warms an OpenCode
  session to apply an effort; the routed selection carries none, so it is not
  warmed. ADR-0180 gives one dated Recommended Profile, and ADR-0181: "a
  repository may keep an older model on purpose"; this Spec changes neither
  the profile nor the comparison. ADR-0198 counts each adapter's reported
  tokens, which the measurement reads beside the router lines. ADR-0199
  cites ADR-0198 but decides a queue's token ceiling, which this Spec does
  not change. ADR-0208 groups work items that share a context; this Spec and
  Spec 0217 were split because a routed model on an existing runtime and a
  fourth runtime do not share one. ADR-0209 cites ADR-0208 but decides how
  the judge suggests grouping, which is unchanged. ADR-0211 drops a missing
  Node preload from the agent environment and holds for OpenCode. ADR-0125:
  "Fixtures are therefore compiled once", which binds every fake OpenCode and
  acpx in the tests. ADR-0187 and ADR-0189 bind the Roundfix Skill's runtime
  reference and its version. ADR-0193: a Spec names its prerequisites, as
  above. ADR-0184: "A TechSpec described a command's new behavior in prose, and
  each reader rebuilt the exact output from it", so the changed surface is a
  Surface Transcript. This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and ADR-0093, ADR-0117,
  ADR-0168, ADR-0176 and ADR-0183 check its consistency by citation and
  receipt. ADR-0182 runs Settlement Checks before each Task commit, ADR-0178
  authorizes a Task commit by its grant, and ADR-0166 records undeclared
  paths. ADR-0069 cites ADR-0050 but decides the Baseline semantic analysis,
  ADR-0096, ADR-0097, ADR-0194, ADR-0195 and ADR-0210 decide the gate's
  machine stage, row carry and evidence, and ADR-0192 decides conflicts on
  derived paths; this Spec changes none of them, so none applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and the cycle authorization of 2026-10-01 ("Tudo, de A a F"),
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`. Sanctioned regeneration:
  `make skills-sync`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- A repository can opt in, in its committed Project Config, to running a
  category's Tasks on the Jev Router through OpenCode, and the selection is
  proved before any Run depends on it.
- The router never receives a prompt from a repository that did not commit
  that choice, and Roundfix's OpenRouter key never leaves the environment
  except in an authorization header.
- No routed prompt is sent once the month's Jev spend reaches US$5, counting
  the judge and the router together, and each routed prompt's cost is on
  record.
- The router's cost, outcome and repairs on real `docs` and `chore` Tasks are
  on record next to the current default, so a later decision about defaults
  rests on measurement.

## User Stories

1. As a maintainer, I want to name the Jev Router as the preferred selection
   of `docs` and `chore` in one repository's Project Config, with the current
   default as fallback, so that I can try it without changing any default.
2. As a maintainer, I want User Config and a one-Run override to refuse the
   router, so that my key never follows a setting into another repository.
3. As a maintainer, I want routed prompts to stop at the same US$5 monthly
   ceiling the judge uses, counted together, so that the experiment cannot
   overspend.
4. As a maintainer deciding about defaults, I want the cost and outcome of
   routed Tasks next to the current default's, so that the decision rests on
   measurement.

## Core Features

1. **One routed selection.** The OpenCode selection
   `{runtime: opencode, model: "roundfix-openrouter/typesafe/jev-router", reasoning_effort: ""}`
   is the Jev Router. A non-empty effort on it is refused when the
   configuration loads. Any other OpenCode model keeps today's behavior.
2. **Project Config only.** A routed selection in User Config is a load error
   naming the rule, and a one-Run override naming it is refused. Project
   Config may name it as a preferred or fallback selection of any category.
3. **An inline provider, the key by placeholder.** For a routed selection's
   sessions only, Roundfix sets `OPENCODE_CONFIG_CONTENT` to a fixed provider
   definition whose key is `{env:ROUNDFIX_OPENROUTER_API_KEY}`. Proof and a
   Run's prompts refuse the selection with `jev_router_key_missing` when the
   variable is empty. No Roundfix output, file, log or Run Event holds the
   key, and `OPENROUTER_API_KEY` is never read.
4. **The shared ceiling gates each routed prompt.** Before each routed
   prompt, Roundfix computes the month's Jev spend as the Judge Log's
   TypeSafe cost plus the larger of its OpenRouter cost and the key's
   `usage_monthly`. At or above the judge's ceiling it refuses with
   `jev_ceiling_reached`; when the log or the key endpoint cannot be read,
   with `jev_spend_unreadable`. Before Agent work begins the refusal moves
   the Fallback Chain on after notification; after work began the Work Item
   fails with that classification.
5. **Each routed prompt is recorded.** After a routed prompt, Roundfix reads
   the key's usage again and appends a `router-prompt` line to the Judge Log
   with the Run, Spec, Task, category, attempt, the change in usage as its
   cost, its tokens and its outcome, so the judge's ceiling counts it too.
6. **The router, measured.** With the key set and spend below the ceiling,
   the measurement replays four archived Tasks, two `docs` and two `chore`,
   on the router and on the current default, under the protocol in the
   TechSpec, and records cost, outcome, prompts, Verification repairs and
   wall time. It stops when the gate refuses.
7. **Documented.** The configuration guide and the Roundfix Skill's runtime
   reference describe the routed selection, its Project Config rule, its key
   and its ceiling.

## User Experience

A maintainer adds to a repository's `.roundfixrc.yml` a `docs` profile whose
preferred selection is the router and whose fallback is the current default,
and runs `roundfix profiles validate --category docs`, which proves it. A
routed Task runs like any other; when the month's Jev spend reaches US$5, the
next routed Task starts on its fallback after a notification that names
`jev_ceiling_reached`. Putting the router in `~/.roundfix/config.yml` makes
every command refuse to load that configuration and say that only Project
Config may select the router.

## Declared breaks

- A selection whose model is `roundfix-openrouter/typesafe/jev-router` is
  refused from User Config and from a one-Run override, where today any
  OpenCode model OpenCode advertises is accepted.
- The Judge Log, until now written only by `roundfix spec judge`, also
  receives `router-prompt` lines.

## Non-Goals / Out of Scope

- Any change to the Recommended Profile, the built-in profiles or this
  repository's Project Config. A default change is a follow-up Backlog Entry.
- The router on the `codex` or `claude` runtimes, which reach subscription
  accounts, not an OpenAI-compatible endpoint.
- Any other OpenRouter model, provider or key, and any TypeSafe endpoint for
  the router.
- Reading `OPENROUTER_API_KEY`, setting a per-key limit on OpenRouter, or any
  account operation.
- Routing a repository other than this one in the measurement, or any
  Secondbrain content anywhere.
- Changing what the judge judges or how it asks.

## Success Metrics

1. Success Metric: a Project Config naming the router with an empty effort
   loads and proves through a fake OpenCode that advertises it; the same
   selection in User Config, with a non-empty effort, or as a one-Run
   override is refused.
2. Success Metric: with a sentinel value in `ROUNDFIX_OPENROUTER_API_KEY`, a
   routed session's acpx environment carries the placeholder configuration
   and the sentinel appears in no argument, configuration value, Run Event or
   Judge Log line; with the variable empty the selection is refused with
   `jev_router_key_missing`; `OPENROUTER_API_KEY` is read by no production
   file.
3. Success Metric: with a fake key endpoint reporting a `usage_monthly` that
   puts the spend at or above US$5, or failing, the first routed prompt is
   not sent and the fallback starts; a later refusal fails the Work Item; and
   below the ceiling the prompt runs and one `router-prompt` line with the
   usage change is appended.
4. Success Metric: the judge package's Judge Log reader reads a
   `router-prompt` line and counts its cost in the month's spend.
5. Success Metric: the measurement record names, for each replayed Task and
   selection, cost, outcome, prompts, repairs and wall time, or the gate's
   refusal that stopped it.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- **OpenRouter's router listing and page.** `GET
  https://openrouter.ai/api/v1/models` lists `typesafe/jev-router` with the
  context, parameters and `-1` prices quoted above, and
  <https://openrouter.ai/typesafe/jev-router> describes it, both read
  2026-10-01.
- **OpenRouter's key reference.** <https://openrouter.ai/docs/api/api-reference/api-keys/get-current-api-key>
  and <https://openrouter.ai/docs/api_reference/limits>, read 2026-10-01:
  `GET /api/v1/key` with a bearer key returns `usage_monthly`, the key's
  credit usage in US dollars for the current UTC month, with `limit` and
  `limit_remaining`; `402` means a credit limit was reached.
- **OpenCode's configuration reference.** <https://opencode.ai/docs/config/>
  and <https://opencode.ai/docs/providers/>, read 2026-10-01: inline
  configuration in `OPENCODE_CONFIG_CONTENT` is a runtime override,
  `{env:VAR}` substitutes an environment variable, and a custom provider uses
  `@ai-sdk/openai-compatible` with `options.baseURL` and `options.apiKey`.
  The installed OpenCode `1.18.34`, which this Spec did not build, behaved as
  documented in the probe above.
- **The live measurement** of Core Feature 6, against OpenRouter, is a
  measurement of a service this Spec did not build.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the query
`qmd query "Jev Router OpenRouter roteamento de modelo custo por tarefa"
--all --files --min-score 0.3`. The source
`wiki/sources/jev-ecosistema-verificadores-routers-2026-09-25.md` describes a
different, CLI-side `jev-router` project, whose policy refuses a downgrade in
a large conversation because switching models wastes more prompt cache than
it saves; that is why this Spec records each routed prompt's cost rather
than assuming a cheaper model is cheaper per Task, and why the measurement
compares whole Tasks. The digest
`wiki/sources/digest-roundfix-orquestradores-e-ecosistema-jev-o-que-se-transfere-ao-roundfix-2026-09-30.md`
notes no published evidence for that project. Exa found OpenRouter's key
reference and limits page, OpenCode's configuration and providers pages, and
the router page cited above; TypeSafe's index was fetched directly.

## Decisions

- **One OpenCode selection, an inline provider.** See ADR-0218. A provider
  of Roundfix's own name keeps OpenCode's built-in `openrouter` provider, and
  any key it would read, out of the path.
- **Project Config only.** See ADR-0218. The maintainer allowed only this
  repository's content to reach the router; a committed per-repository
  choice is the boundary Roundfix can hold.
- **One ceiling for all Jev use, the key's own usage as its floor.** See
  ADR-0218. The key's `usage_monthly` counts every OpenRouter call on that
  key even when a line is missing, and the Judge Log counts the direct
  TypeSafe spend the key cannot see.
- **The Judge Log records router prompts.** A second log would leave the
  judge blind to router spend.
- **Opt-in only.** The maintainer decided on 2026-10-01 that both
  experiments deliver measurement and an opt-in entry, never a default
  change.

## Open Questions

- Whether OpenCode's OpenAI-compatible provider keeps tool calls working
  when the router changes the model between turns is unknown until the
  measurement. Default until measured: a routed Task that fails on a tool
  call counts as not settled, and the reading records it.

## Technical candidate

The [_techspec.md](_techspec.md) records the inline provider, the gate, the
router line, the measurement protocol, coverage and build order.
