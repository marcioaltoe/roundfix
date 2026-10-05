---
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
status: active
created: 2026-10-04
surfaces: [backend, docs]
---

# A Jev Router that checks its credit and names its model

On 2026-10-04 a routed Task ended at OpenRouter instead of in Roundfix. The
Jev Router gate read the month's Jev spend (US$9.69 of the US$50 ceiling) and
the key's remaining limit (US$40.31) and sent a Verification repair prompt.
The OpenRouter account behind the key held about US$5.11, so OpenRouter
refused the request with HTTP 402 ("This request would exceed your available
credits given your current in-flight requests"). Roundfix reported
`agent/protocol error`, no fallback ran, and both of the run's Judge Log
lines recorded an empty model, so the ten-fold price spread between the
2026-10-02 and 2026-10-04 routed prompts cannot be traced to a model. The
maintainer approved cost controls for the router that day, after asking for
the work to be "implementado e testado sem limitações. Depois vemos o custo e
como reduzir o custo." This Spec makes the gate read the account's credit and
refuse below a floor before the prompt, names an OpenRouter credit refusal so
the Fallback Chain takes the Task when it still can, and records on each
`router-prompt` Judge Log line the model the router reported.

## Prerequisites

None. Specs 0227 and 0228 are authored in the same cycle; if either also
raises the Roundfix Skill's version, the operator orders the queue so that the
later Spec raises it from the earlier one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the new
  key `jev.router_min_credit_usd` follows the dotted configuration names, and
  the reasons `openrouter_credit_low` and `openrouter_credit_refused` follow
  the existing refusal codes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the gate sends one more HTTPS read
  before each routed prompt, `GET https://openrouter.ai/api/v1/credits`, with
  the same bearer key as the existing key read, and the key is read only from
  `ROUNDFIX_OPENROUTER_API_KEY`. Routed model requests that OpenCode sends now
  pass through a relay that Roundfix owns on the loopback interface and that
  forwards them unchanged to `https://openrouter.ai/api/v1`; the key crosses
  Roundfix's memory and never its arguments, files, logs, Run Events or the
  inline configuration. No test or Verification reaches the network: every
  OpenRouter answer comes from a local stand-in. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0234 (this Spec) decides the
  credit floor, the named credit refusal and the relay. ADR-0218: "Roundfix
  gives OpenCode an inline provider definition", which ADR-0234 points at the
  relay while the key placeholder, the project-only selection and the
  refusal rule of that decision stand, ADR-0218: "At or above the ceiling, or
  when either source cannot be read, the prompt is not sent". ADR-0231: "The
  ceiling is `jev.monthly_ceiling_usd` in User Config", the scope the floor
  follows for the same reason. ADR-0201: "the model that answered", the
  Judge Log field the router lines now fill. ADR-0050: "once Agent work
  begins, Roundfix fails the Work Item instead of switching models over
  potentially modified state", and ADR-0114: "A Fallback Selection may switch
  ACP Runtime automatically only while Agent work has not begun"; both decide
  what a credit refusal does. ADR-0002: "Roundfix uses YAML for User Config
  at", the file the operator edits. ADR-0027: "truly unknown keys keep
  failing strict validation", so an older binary refuses the new key.
  ADR-0187 splits the Roundfix Skill by command and ADR-0189 ties an owned
  skill's version to its content, so the skill edit raises the version.
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts"; this Spec
  changes no command's arguments or output, so it declares none. The gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and
  ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  this Spec's consistency by citation and receipt. ADR-0178 authorizes each
  Task commit by its grant, ADR-0182 runs Settlement Checks before it, and
  ADR-0166 records undeclared paths; every Task declares its paths. ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park, ADR-0069 cites ADR-0050 but decides the Baseline semantic analysis, ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0192 cites ADR-0178 but decides derived-path conflict regeneration, ADR-0200 and ADR-0209 cite ADR-0201 but decide that the judge never gates and that it only suggests which sources share a Spec, ADR-0208 cites ADR-0209 but decides how sources are grouped, ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot, and ADR-0180 cites ADR-0069 but decides the dated Recommended Profile, which this Spec leaves unchanged with the router outside it, ADR-0181 cites ADR-0180 but decides where a configuration is compared with that profile, and ADR-0217 cites ADR-0180 but decides how a Cursor selection is named; this Spec changes none of them, so none applies.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths, and the maintainer authorized skill
  edits ("considere autorizado a ajustar todas as skills se necessário") and
  this cost-control cycle on 2026-10-04 ("Aprovado"). No other Governed Path
  changes, and the repository's own Project Config is not edited. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0229-a-jev-router-that-checks-credit-and-names-its-model/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A routed prompt is not sent when the OpenRouter credit left, the lower of
  the account's balance and the key's remaining limit, is below a floor the
  maintainer sets; before Agent work that refusal hands the Task to the
  Fallback Chain.
- An OpenRouter credit refusal of a routed request is reported under its own
  name, never as `agent/protocol error`.
- Every `router-prompt` Judge Log line names the models and providers the
  router reported for that prompt and its last response id.
- A machine whose User Config sets nothing gets a US$15 floor, and a
  repository cannot change it.

## User Stories

1. As the maintainer, I want a routed prompt refused before it is sent when my
   OpenRouter account cannot cover it, so that a Task falls back instead of
   paying for a prompt OpenRouter will cut off.
2. As the maintainer, I want to set the credit floor in my User Config, so
   that it matches what one routed prompt costs on my account.
3. As the maintainer, I want an OpenRouter credit refusal named as one, so
   that the Run tells me to add credit instead of reporting a protocol error.
4. As the maintainer, I want each routed prompt's Judge Log line to name the
   model the router chose, so that I can trace what a routed Task cost to the
   models that ran.
5. As a user of a cloned repository, I want a Project Config value for the
   floor ignored with a warning, so that a repository cannot lower what my
   account must hold before my key is spent.

## Core Features

1. **The account credit in the gate.** Before each routed prompt, after the
   key-limit and ceiling checks, the gate reads the account's total credits
   and total usage from OpenRouter. The credit left is the account balance
   (credits less usage), or the key's remaining limit when the key reports
   one and it is lower.
2. **A floor.** When the credit left is below the floor, the prompt is not
   sent and the refusal is `openrouter_credit_low`, naming the credit left
   and the floor. Before Agent work the Fallback Chain takes the Task; after
   it, the Work Item fails with that reason. An unreadable credit answer
   refuses the prompt as `jev_spend_unreadable`.
3. **A configured floor.** User Config accepts `jev.router_min_credit_usd`, a
   finite number of US dollars greater than zero; unset, the floor is US$15.
   A Project Config value is removed before it is read, with the warning
   Roundfix gives for a User Config-only setting; an invalid value is a
   configuration error naming the key.
4. **A named credit refusal.** When OpenRouter answers a routed request with
   HTTP 402, the prompt fails with `openrouter_credit_refused`, carrying
   OpenRouter's `limit_source` when it gives one. Before Agent work it is a
   failed selection and the Fallback Chain takes the Task; after it, the Work
   Item fails with that reason. The prompt's Judge Log line records it.
5. **The routed model on the Judge Log.** A routed session's model requests
   pass through a loopback relay that forwards them unchanged and notes what
   each response reports. The prompt's `router-prompt` line records the
   distinct models and providers in first-seen order and the last response
   id; a prompt whose responses report none leaves them empty, as today.
6. **The guides and the skill say so.** The configuration guide and the
   Roundfix Skill's `runtime` reference describe the floor, its key, default
   and scope, both new reasons, the relay and the new Judge Log fields.
   ADR-0218 carries a note that ADR-0234 supersedes it in part.

## User Experience

A machine with enough credit sees nothing new except the model on each
`router-prompt` line. With less credit than the floor, a routed Task falls
back before its first prompt with the fallback notice naming
`openrouter_credit_low`, or fails after work began with that reason. A
request OpenRouter refuses for credit ends with `openrouter_credit_refused`
instead of `agent/protocol error`. A Project Config that sets the floor
prints one warning on standard error and keeps the User Config value.

## Non-Goals / Out of Scope

- Changing the `docs` and `chore` default profiles or any Recommended
  Profile: the router stays project-selected (maintainer decision of
  2026-10-04).
- Bounding one prompt's cost: a per-request output cap, a price ceiling or a
  candidate restriction sent to the router, or a dedicated measurement key.
- Retrying a request OpenRouter refused for its in-flight budget, or switching
  models after Agent work began.
- Taking each prompt's cost from the responses instead of the key's usage
  change.
- Resuming a routed session in a Roundfix process other than the one that
  opened it.
- A floor or a relay for `roundfix spec judge`, whose requests cost cents.
- Writing the operator's User Config, adding the key to `roundfix config init`
  or relaxing strict validation for older binaries.

## Success Metrics

1. Success Metric: with a local OpenRouter stand-in reporting an account
   balance of US$5.11 and a key remaining limit of US$40.31, and no floor
   configured, the gate refuses the routed prompt with
   `openrouter_credit_low` before it is sent, and a Task whose Agent work had
   not begun is taken by its fallback; with a balance of US$20 the prompt is
   sent; with a key remaining limit of US$3 and a balance of US$200 it is
   refused.
2. Success Metric: a routed prompt whose request the stand-in answers with
   HTTP 402 and `limit_source` `openrouter_credits` ends with
   `openrouter_credit_refused`, never `agent/protocol error`; before Agent work
   the fallback takes the Task.
3. Success Metric: a routed prompt whose stand-in responses report the models
   `a/one` and `b/two` writes one `router-prompt` line whose model is
   `a/one, b/two`, with the providers and the last response id, and the key
   appears in no log, file or Run Event.
4. Success Metric: a Project Config floor leaves the effective floor at the
   User Config value or US$15 and prints the User Config-only warning; zero, a
   negative value and an infinite value fail with the key named.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- The 2026-10-04 Jev Router measurement addendum, a measurement this Spec did
  not design, and the run's own records that it cites: OpenCode's message
  record of the refused request (HTTP 402 and OpenRouter's message) and the
  acpx session stream, whose prompt answer was a JSON-RPC error carrying the
  same message and no model, and whose usage reports carried no model either.
- OpenRouter's published limits page
  (<https://openrouter.ai/docs/api_reference/limits>), read 2026-10-04: credit
  limits come from the account balance, the key's limit and an in-flight
  budget; each answers HTTP 402, and `error.metadata.limit_source` is
  `openrouter_in_flight_budget`, `openrouter_key_limit` or
  `openrouter_credits`.
- OpenRouter's credits reference
  (<https://openrouter.ai/docs/api/api-reference/credits/get-credits>), read
  2026-10-04: `GET /api/v1/credits` returns `data.total_credits` and
  `data.total_usage` and is documented as needing a management key.
- OpenRouter's router guide
  (<https://openrouter.ai/docs/guides/routing/routers/auto-router>), read
  2026-10-04: a router's response carries "the `model` field showing which
  model was actually used", with the generation `id`.
- Read-only queries on 2026-10-04 with `ROUNDFIX_OPENROUTER_API_KEY`, an
  ordinary key whose `is_management_key` is false: `GET /api/v1/key` answered
  200 with `data.limit`, `data.limit_reset`, `data.limit_remaining` and
  `data.usage_monthly` among its fields, and `GET /api/v1/credits` answered
  200 with `data.total_credits` and `data.total_usage` only. No value was
  recorded.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "OpenRouter credits balance in-flight 402 router model" --all
--files --min-score 0.3`; it holds this repository's own mirrored Specs and
a model-price monitoring page (`wiki/concepts/modelos-custos-e-selecao.md`),
and nothing on OpenRouter's account credit, its 402 sources or a router's
reported model. Exa found OpenRouter's limits page, its errors page, its
credits reference and its router guide cited above.

## Decisions

- The gate takes the lower of the account balance and the key's remaining
  limit and refuses below `jev.router_min_credit_usd`, a User Config value
  with the US$15 default. See ADR-0234.
- An OpenRouter HTTP 402 on a routed request is `openrouter_credit_refused`;
  the Fallback Chain boundary of ADR-0114 is unchanged. See ADR-0234.
- A loopback relay owned by Roundfix carries routed requests, because only
  OpenRouter's responses name the routed model. See ADR-0234.
- An unreadable credit answer refuses the prompt, as an unreadable spend
  does.

## Open Questions

- When does the operator set a floor other than US$15? Default: never until
  a measurement shows US$15 refuses routed prompts the account could have
  paid for; set earlier, an older binary refuses the User Config (ADR-0027).
- What if OpenRouter enforces the management-key requirement on its credits
  read? Default: every routed prompt is refused as `jev_spend_unreadable`,
  and the decision is revisited with a Backlog Entry.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
