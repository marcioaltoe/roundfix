---
spec: 0234-an-openrouter-key-per-stage
status: active
created: 2026-10-05
surfaces: [backend, cli, docs]
---

# An OpenRouter key per stage

Roundfix spends on OpenRouter in two stages: the Jev judge that advises Spec
authoring, and, once Spec 0233 ships the light tier, implementation of a Task
on an open model through OpenCode. Both read the one Roundfix key
`ROUNDFIX_OPENROUTER_API_KEY`, so OpenRouter's activity export, which groups
spend by API key, shows one total and cannot say what each stage cost. On
2026-10-05 the maintainer decided to create one key per stage, each with its
own monthly limit, so the export measures each stage exactly. This Spec lets
each stage read its own key first, keeps the shared key as the documented
fallback until the new keys exist, and records which variable each stage used
without ever recording a key.

## Prerequisites

Spec 0233 (the light tier on open models) must be merged first. It adds the
one place that reads the implementation stage's OpenRouter key, the
implementation spend record and the `openrouter.implement_monthly_ceiling_usd`
ceiling; the implementation-stage Task of this Spec changes that place and
that record, and its Verification runs tests that live beside them. The
Delivery Queue does not order Specs by itself, so the operator queues this
Spec after 0233.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the new
  names are environment variable names and one JSON field, `key_variable`.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — each stage still authenticates to
  OpenRouter with a bearer key read only from the process environment, and
  only the choice of variable changes. A key value is never printed, logged,
  stored or written to a file; records and reports carry the variable's name.
  The judge's destinations are unchanged, and no test or Verification command
  reaches OpenRouter or TypeSafe. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0239 (this Spec) decides the rule:
  "Each OpenRouter stage now reads its own Roundfix variable first and the
  shared key second." It supersedes ADR-0201 in part, whose key scope says a
  later feature "reads these two names, never the generic"
  `OPENROUTER_API_KEY`; the generic key stays unread. The judge's boundary
  is unchanged, ADR-0201: "No other service receives a request, and a run
  never switches recipient". The judge's ceiling stays a User Config value,
  ADR-0231: "Project Config cannot set it". The subscription rule stands and
  the implementation key changes no selectable model, ADR-0235: "OpenAI and
  Anthropic models run only through the codex and claude subscriptions".
  ADR-0238 applies: "A light tier runs low-complexity Tasks on an open model",
  and this Spec routes the implementation key helper it introduced.
  ADR-0200 and ADR-0209 stay as they are: the judge stays advisory and asks
  the same questions. ADR-0208 is honored by adopting nothing, since Spec
  0233 owns the one Backlog Entry that shares this context. ADR-0229 and ADR-0237 do not apply, because no Delivery Queue
  park or archive behavior changes. ADR-0184 binds the changed `spec judge` output as Surface Transcripts.
  ADR-0179 bounds the Governed Paths in `_authorization.md`, and ADR-0189 and
  ADR-0233 bind each owned skill's version raise. The authored QA gate
  follows ADR-0080: "QA verdicts distinguish environment-blocked rows", and
  ADR-0091: "required to be terminal and to depend on every leaf"; ADR-0096,
  ADR-0097 and ADR-0167 bind its machine stage, its row carry and its pre-PR
  Pull Request row, ADR-0104: "Every Spec therefore rests at least one named"
  acceptance row on outside evidence, ADR-0194, ADR-0195 and ADR-0210 bind
  what a QA row records, when it is observed again and its evidence snapshot,
  and ADR-0182 settles each Task on the facts its gate checks. ADR-0093,
  ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this Spec's
  consistency by citation and receipt. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's spec and runtime
  references, both copies of its `SKILL.md`, and both copies of the
  write-prd and write-techspec `SKILL.md` are Governed Paths. The
  maintainer's standing grant of 2026-09-30, "considere autorizado a ajustar
  todas as skills se necessário", and the grant of 2026-10-05, "Concedo",
  cover them; this Spec changes no other Governed Path. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0234-an-openrouter-key-per-stage/_authorization.md`; bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`,
  `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`,
  `skills/write-techspec/SKILL.md`.

## Goals

- The judge uses `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` when it is set, and
  otherwise behaves exactly as today: `ROUNDFIX_OPENROUTER_API_KEY`, then
  `ROUNDFIX_TYPESAFE_API_KEY`.
- Implementation on an open model through OpenCode uses
  `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` when it is set, and otherwise
  `ROUNDFIX_OPENROUTER_API_KEY`.
- Every Judge Log line, every `spec judge` report and every implementation
  spend record names the variable it used, and none contains a key value.
- `roundfix doctor` shows, by name, which variables each stage would read
  and whether each is set.
- Each stage's monthly ceiling keeps applying to that stage alone, whichever
  key it reads.

## User Stories

1. As the maintainer, I want the judge and implementation to bill to two
   different OpenRouter keys, so that the activity export shows what each
   stage cost.
2. As the maintainer, I want both stages to keep working on the shared key
   until I create the new keys, so that nothing stops in between.
3. As an operator, I want `roundfix doctor` to name each stage's variables
   and whether they are set, so that I can tell which key a stage will use
   without printing a secret.
4. As an operator reading the Judge Log or a spend record, I want each entry
   to name the variable used, so that Roundfix's own records match the
   export.

## Core Features

1. **A judge key.** The judge reads `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` first,
   then `ROUNDFIX_OPENROUTER_API_KEY`, then `ROUNDFIX_TYPESAFE_API_KEY`. An
   empty variable counts as unset. Without any of them the run skips and
   names the judge key first as the one to set.
2. **An implementation key.** The implementation stage reads
   `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` first, then
   `ROUNDFIX_OPENROUTER_API_KEY`, through the one place Spec 0233 added for it.
3. **One list of names.** Each stage's variables and their order are named
   once in Roundfix and every reader, the doctor included, takes them from
   there. The generic `OPENROUTER_API_KEY` is never read.
4. **The variable is recorded.** The judge's text summary and its JSON
   report name the variable used, as do each Judge Log line and each
   implementation spend record, in `key_variable`.
5. **Doctor names the stage keys.** The `environment:` line lists the judge's
   variables and the implementation stage's variables, each set or not set,
   in preference order.
6. **The guides say it.** The `spec judge` and `doctor` command references,
   the configuration guide, the Roundfix Skill, and the write-prd and
   write-techspec skills name the stage keys, their order and the fallback.

## Non-Goals / Out of Scope

- Creating the OpenRouter keys or setting their limits; the maintainer does
  that in OpenRouter.
- Changing either ceiling's value, scope or sum, the judge's recipients, the
  pinned Jev model, or which open models the light tier may select.
- A stage key for TypeSafe direct; `ROUNDFIX_TYPESAFE_API_KEY` stays the
  judge's last alternative.
- Reading spend from OpenRouter, or importing the activity export.
- The Makefile, `go.mod`, CI workflows, `CONTEXT.md` and `CHANGELOG.md`.

## Success Metrics

1. Success Metric: with `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` and
   `ROUNDFIX_OPENROUTER_API_KEY` both set to different fake values, a judge
   run against a fake transport sends only the judge key, its summary and
   JSON report name `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, and every Judge Log
   line carries `"key_variable":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY"` and
   neither value (before: the judge key is ignored and the shared key is
   sent).
2. Success Metric: with only `ROUNDFIX_OPENROUTER_API_KEY` set, the judge
   sends it and names it, and with no key the run skips with the reason of
   Surface Transcript 1.
3. Success Metric: with `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` set, an
   open-model implementation uses that variable and its spend record names
   it; with only the shared key set it uses and names the shared key; the
   generic `OPENROUTER_API_KEY` alone gives no key.
4. Success Metric: `roundfix doctor`'s `environment:` line lists the judge's
   three variables and the implementation stage's two variables with set or
   not set, and its output contains no key value.
5. Success Metric: every existing judge, doctor and light-tier test passes
   with only the declared output changes.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- OpenRouter's Activity guide
  (<https://openrouter.ai/docs/guides/features/activity>), read 2026-10-05:
  Activity rolls generations up into spend "you can slice by model, provider,
  API key, app, or user", and its Explore results "can be downloaded as CSV".
- OpenRouter's API key reference
  (<https://openrouter.ai/docs/api/api-reference/api-keys/create-a-new-api-key>),
  read 2026-10-05: a key takes an optional `limit` in US dollars and a
  `limit_reset` of `daily`, `weekly` or `monthly`; and the limits reference
  (<https://openrouter.ai/docs/api_reference/limits>) names "Per-key credit
  limits, an optional spending cap configured on an individual API key".
- The Backlog Entry "A judge-assigned model tier per Task" of 2026-10-05,
  whose Measurement section records the shared key billing both DeepSeek
  replays and judge calls in one `usage_monthly` (US$0.4181 for the whole
  measurement).

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "openrouter api key per stage cost attribution"`
(`--all --files --min-score 0.3`); it returned this repository's mirrors of
Spec 0205, ADR-0201 and the Backlog Entry above, none of which separates the
stages. Exa found OpenRouter's Activity guide, its API key reference and its
limits reference, which together show that a key is the unit OpenRouter both
limits and reports by. The only open Backlog Entry is the one above; its
"Separate OpenRouter keys per stage" bullet is this Spec's scope, and Spec
0233, which implements the rest of it, owns the entry, so this Spec adopts
nothing. There is no unresolved Finding.

## Decisions

- One key per stage, the stage key first and the shared key as fallback.
  See ADR-0239.
- Record the variable's name in every record and report, never a value.
- Keep TypeSafe direct as the judge's last alternative with no stage key.

## Open Questions

None. Until the maintainer creates the two keys, both stages read the shared
key and behave as before.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
