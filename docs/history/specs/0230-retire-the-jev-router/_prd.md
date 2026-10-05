---
spec: 0230-retire-the-jev-router
status: archived
created: 2026-10-05
surfaces: [backend, cli, docs]
archived: "2026-10-05"
source_slug: 0230-retire-the-jev-router
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: Q7a needs the operator''s OpenRouter activity export (docs/_inbox/openrouter_activity_2026-10-05.csv in the main checkout; the measurement addendum records its figures), Q7c and Q7d would need OpenRouter requests that this Spec forbids, and Q10 has no Pull Request yet; the queue opens it. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 73587c4d7ea9ccf5531168868093219ba53b3f82
---


# Retire the Jev Router

The Jev Router let a repository hand a Task to OpenRouter's
`typesafe/jev-router`, which picks a model for each request. OpenRouter's
activity export for Roundfix's key from 2026-10-01 to 2026-10-05 shows that
the routed sessions were billed, pay-per-use, for `openai/gpt-6-astra`
(US$15.35), `anthropic/claude-opus-5.5` (US$10.12) and `openai/gpt-6.1-sol`
(US$4.81), models this repository already reaches through the codex and
claude subscriptions, while 456 direct judge calls cost US$0.014. On
2026-10-05 the maintainer set the rule "os modelos da openai e anthropic
devem ser utilizado exclusivamente pela assinatura que o codex e claude
fornecem", decided "Aposentar o router", and, after weighing a restricted
router, confirmed "Vamos seguir com a remoção". This Spec removes the router
and everything that exists only for it, and makes Roundfix refuse any Agent
Selection that would reach an OpenAI or Anthropic model through OpenRouter,
with a message that names the rule. The advisory judge keeps working exactly
as it does today.

## Prerequisites

None. No other Spec is active.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the new
  refusal reason `subscription_only` follows the existing selection refusal
  codes, and the retired key keeps its dotted name in the deprecated-key
  list. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the change removes outbound HTTP:
  the router gate's key and credit reads and the loopback relay that
  forwarded OpenCode's requests to OpenRouter. The judge's requests, its two
  recipients and its keys are unchanged. No test or Verification reaches the
  network. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0235 (this Spec) decides the
  retirement and the refusal: "Roundfix therefore retires the Jev Router
  entirely". It supersedes ADR-0218: "Only Project Config may select the
  router", whose selection and gate stop applying, and ADR-0234: "A relay sees
  what the router reports", whose relay goes with it. ADR-0231: "It must be a
  finite number greater than zero", the judge's ceiling, stays; only its
  router parts are superseded. ADR-0201: "Each request appends one line to a
  Judge Log in Roundfix Home", unchanged for the judge, and ADR-0200: "Spec
  authoring gets an advisory judge that never gates", unchanged. ADR-0027:
  "recognized deprecated keys are ignored with one stderr warning naming the
  replacement", so the retired floor key degrades to a warning. ADR-0050:
  "once Agent work begins, Roundfix fails the Work Item instead of switching
  models over potentially modified state", and ADR-0114: "A Fallback Selection
  may switch ACP Runtime automatically only while Agent work has not begun"; a
  refused selection is refused before any prompt, so its Fallback Chain takes
  the Task. ADR-0049: "each present profile replaces the complete
  lower-precedence profile", so a refused selection fails its whole profile at
  validation. ADR-0187 splits the Roundfix Skill by command and ADR-0189 ties
  an owned skill's version to its content, so the skill edit raises the
  version. ADR-0184: "A TechSpec now declares numbered Surface Transcripts";
  this Spec declares the transcripts of the new refusal. The authored QA gate
  follows ADR-0080: "QA verdicts distinguish environment-blocked rows", and
  ADR-0091: "required to be terminal and to depend on every leaf"; ADR-0096,
  ADR-0097 and ADR-0167 bind its machine stage, its row carry and its pre-PR
  Pull Request row, ADR-0104: "Every Spec therefore rests at least one named"
  acceptance row on outside evidence, ADR-0194, ADR-0195 and ADR-0210 bind
  what a QA row records, when it is observed again and its evidence snapshot,
  and ADR-0182 settles each Task on the facts its gate checks. ADR-0093,
  ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this Spec's
  consistency by citation and receipt, and ADR-0233 regenerates the raised
  skill version at merge. ADR-0069 cites ADR-0050 but decides the Baseline
  semantic analysis, ADR-0180 cites ADR-0069 but decides the dated Recommended
  Profile, which this Spec leaves unchanged, ADR-0181 cites ADR-0180 but
  decides where a configuration is compared with that profile, ADR-0209 cites
  ADR-0200 but decides how the judge suggests sources and ADR-0208 cites
  ADR-0209 but decides how sources are grouped, ADR-0217 cites ADR-0049 but
  decides how a Cursor selection is named, and ADR-0229 cites ADR-0167 but
  decides how an operator archive resumes a park; this Spec changes none of
  them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files, its
  `SKILL.md` mirror and the repository's Project Config comment are Governed
  Paths. The maintainer authorized skill edits on 2026-09-30 ("considere
  autorizado a ajustar todas as skills se necessário") and edits to the
  guides and `.roundfixrc.yml` ("Autorizar os dois"), and decided the
  retirement on 2026-10-05 ("Aposentar o router", "Vamos seguir com a
  remoção"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0230-retire-the-jev-router/_authorization.md`; bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `skills/roundfix/SKILL.md`, `.roundfixrc.yml`.

## Goals

- No Agent Selection reaches an OpenAI or Anthropic model through OpenRouter
  from Roundfix, and a configuration that tries is refused with a message
  that names the subscription rule.
- The Jev Router, its gate, its credit floor, its relay and its Judge Log
  lines no longer exist in the product, the guides or the Roundfix Skill.
- A machine whose configuration still carries the retired credit floor keeps
  working, with one warning.
- The advisory judge, its ceiling and its Judge Log behave exactly as before.

## User Stories

1. As the maintainer, I want Roundfix to refuse a profile that sends Tasks to
   an OpenAI or Anthropic model through OpenRouter, so that those models are
   only ever billed through my codex and claude subscriptions.
2. As the maintainer, I want the refusal to tell me which field and model
   broke the rule and what the rule is, so that I can fix the profile without
   reading the source.
3. As an operator whose User Config still sets the router's credit floor, I
   want Roundfix to load it with a warning, so that retiring the router does
   not break my machine.
4. As the maintainer, I want `roundfix spec judge` to keep working with the
   same keys, ceiling and Judge Log, so that the cheap direct judge stays.
5. As a reader of the guides and the Roundfix Skill, I want them to describe
   the rule and no longer describe the router, so that I do not configure a
   feature that is gone.

## Core Features

1. **The router is gone.** Roundfix no longer gives OpenCode an OpenRouter
   provider of its own, no longer runs a relay, no longer gates a prompt on
   OpenRouter spend or credit, and no longer writes `router-prompt` lines to
   the Judge Log. Its refusal reasons (`jev_router_key_missing`,
   `jev_router_key_unbounded`, `jev_spend_unreadable`, `jev_ceiling_reached`,
   `openrouter_credit_low`, `openrouter_credit_refused`) are no longer
   produced by Agent work.
2. **The subscription rule in configuration.** An Agent Selection on the
   `opencode` runtime is refused when its model is under the retired
   `roundfix-openrouter` provider, or under OpenRouter with an OpenAI or
   Anthropic model, or with an OpenRouter router or saved preset, because a
   router can pick an OpenAI or Anthropic model per request. The refusal
   names the field, the model and the rule. It holds for Project Config,
   User Config, a one-Run override and `roundfix profiles configure`, which
   writes nothing when it refuses. Other OpenRouter models stay selectable.
3. **The subscription rule before a session.** A selection that reaches the
   runner without passing configuration validation is refused before any
   prompt is sent, with the reason `subscription_only`; before Agent work the
   Fallback Chain takes the Task.
4. **The retired key degrades.** `jev.router_min_credit_usd` in User Config
   or Project Config is ignored with one warning that says the router was
   retired and the key can be removed.
5. **The judge is unchanged.** `roundfix spec judge`, its keys, its
   recipients, `jev.monthly_ceiling_usd` and its Judge Log lines behave as
   before; lines a routed prompt wrote earlier stay in their month's file and
   still count toward that month's ceiling.
6. **The record says so.** The configuration guide, the Roundfix Skill's
   `runtime` reference, the model reference and the repository's Project
   Config comment state the rule. ADR-0218 and ADR-0234 were retired as
   superseded by ADR-0235, and ADR-0231 notes the part ADR-0235 supersedes,
   when this Spec was authored.

## User Experience

A profile that names a refused model stops every command that loads
configuration with one error on standard error, naming the field, the model
and the rule, and a non-zero exit; `roundfix profiles configure` refuses the
fragment and leaves the file untouched. A machine whose configuration still
sets the floor sees one warning per command and nothing else. A repository
that selects any other model sees no change.

## Non-Goals / Out of Scope

- OpenAI or Anthropic models served by OpenCode's own `openai`, `anthropic`,
  `opencode` and `opencode-go` providers, or by the Cursor runtime; ADR-0235
  leaves them to a later decision.
- An OpenCode provider the operator defines outside Roundfix that points at
  OpenRouter, and third-party routers under other authors.
- Choosing a model tier per Task, which returns as a separate Backlog Entry
  the operator owns.
- Changing the judge, its Judge Log schema, its ceiling or its keys, or
  removing past `router-prompt` lines from the Judge Log.
- Changing any built-in or Recommended Profile, the `docs` and `chore`
  defaults or the changelog.

## Success Metrics

1. Success Metric: Project Config, User Config, a one-Run override and
   `roundfix profiles configure` each refuse an `opencode` selection naming
   `roundfix-openrouter/typesafe/jev-router`, `openrouter/openai/gpt-6.1-sol`,
   `openrouter/anthropic/claude-opus-5.5` or `openrouter/openrouter/auto`,
   with the field, the model and the rule in the message, and the configure
   command writes nothing; `openrouter/deepseek/deepseek-v4-pro` on
   `opencode` and `gpt-6.1-sol` on `codex` still load.
2. Success Metric: the runner refuses a session for
   `openrouter/anthropic/claude-opus-5.5` on `opencode` before the ACP
   adapter is started, and a Task whose preferred selection is refused
   before Agent work is taken by its fallback with the reason
   `subscription_only`.
3. Success Metric: outside the judge, no Roundfix source names OpenRouter's
   API host or sets `OPENCODE_CONFIG_CONTENT`, nothing writes a
   `router-prompt` line, and the judge's tests pass with the judge package
   byte-identical to the starting tree.
4. Success Metric: a User Config and a Project Config that set
   `jev.router_min_credit_usd: 20` each load with exactly one warning naming
   the key.
5. Success Metric: no live guide, ADR or Roundfix Skill reference describes
   the router as current, and the configuration guide and the `runtime`
   reference quote the rule.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- OpenRouter's activity export for the `roundfix_jev` key, 2026-10-01 to
  2026-10-05, 727 generations, kept by the maintainer outside the repository
  (SHA-256 `51dd529854af6844f41d88dabbc6c0dac8bb4acf18683e26d9a55cb738f95f09`).
  Summed by model: `openai/gpt-6-astra` US$15.35, `anthropic/claude-opus-5.5`
  US$10.12, `openai/gpt-6.1-sol` US$4.81, `deepseek/deepseek-v4.1-flash`
  US$0.10 and `typesafe/jev-1.13` US$0.0138 over 456 calls.
- The "Relay confirmation, 2026-10-05" and "Maintainer decision, 2026-10-05"
  sections of Spec 0218's archived measurement record of 2026-10-04, a
  measurement this Spec did not design: one routed `chore` replay recorded
  `deepseek/deepseek-v4.1-flash, openai/gpt-6.1-sol,
  anthropic/claude-opus-5.5` as the routed models, at US$9.86.
- OpenRouter's Auto Router guide
  (<https://openrouter.ai/docs/guides/routing/routers/auto-router>), read
  2026-10-05: the slugs `openrouter/auto` and `openrouter/auto-beta` route
  each request to a model chosen at request time, with `anthropic/*` among
  its documented candidates, so a router slug can reach the models the rule
  reserves.
- OpenRouter's model variants guide
  (<https://openrouter.ai/docs/guides/routing/model-variants/overview>), read
  2026-10-05: a variant is a `:` suffix on the same author and slug, so the
  author segment identifies the vendor whatever the variant.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "OpenRouter router OpenAI Anthropic subscription codex claude"
--all --files --min-score 0.3`. Its synthesis on subscription authentication
for local agents records the same boundary for this machine: an OpenAI or
Anthropic key bills the API pay-per-use and is a conscious fallback only,
while the codex and claude CLIs consume the subscriptions. Its model-cost
page and this repository's mirrored model reference list OpenRouter's
OpenAI and Anthropic identifiers as selectable through OpenCode, which is the
path this Spec closes. Exa found OpenRouter's Auto Router and model variants
guides cited above. The repository has no other active Spec and no open
Backlog Entry or unresolved Finding that shares this context; the open
Backlog Entry on the Judge Log's routed-spend lag is declined with this
Spec, because nothing will write routed spend any more.

## Decisions

- Retire the router whole, rather than keep it with an exclusion list. See
  ADR-0235.
- Refuse by OpenRouter author (`openai`, `anthropic`, and the router authors
  `openrouter` and `typesafe` and `@` presets) rather than refuse OpenRouter
  as a whole. See ADR-0235.
- Check in configuration validation and again before a session. See
  ADR-0235.
- Keep the retired floor key as a deprecated key with a warning (ADR-0027).
- Leave the judge package byte-identical, including the two Judge Log
  helpers only the router used, so the judge's unchanged behavior is
  provable from the diff.

## Open Questions

- Should OpenCode's own OpenAI and Anthropic providers and the Cursor
  runtime fall under the rule? Default: not in this Spec; ADR-0235 leaves
  them to a later decision.
- Should the operator restrict OpenAI and Anthropic models on the OpenRouter
  account as well? Default: recommended as an operator step after merge; it
  is outside the repository.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
