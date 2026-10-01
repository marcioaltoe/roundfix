---
spec: 0189-profiles-that-follow-the-current-models
status: archived
created: 2026-09-30
surfaces: [backend, cli, docs]
archived: "2026-10-01"
source_slug: 0189-profiles-that-follow-the-current-models
---


# Profiles that follow the current models

Roundfix names its models in four places, and on 2026-09-30 all four were
behind the ACP Runtimes they describe:

- The Model Catalog offers `gpt-5.4` and `gpt-5.4-mini`, which left Codex on
  2026-08-31, and `gpt-5.3-codex-spark`, which left on 2026-09-14. It offers no
  GPT-6 model and describes `opus` as Opus 5. The picker's Claude reasoning
  efforts are `default`, `high` and `maximum`; the adapter advertises
  `default`, `low`, `medium`, `high`, `xhigh` and `max`.
- The built-in fallback, the legacy Codex default and the Baseline semantic
  analysis fallback are `gpt-5.5`. It retires from Codex on 2026-10-14.
  Preflight proves every configured tuple, so from that day a repository on
  built-in profiles cannot start a Run.
- The recommendation ranking is dated 2026-08-07, covers five categories, and
  carries a DeepSWE result and an average cost for every row. No model that
  leads today has both figures published.
- The adapter floors are `codex-acp` 1.1.5 and `claude-agent-acp` 0.63.0. The
  current tuples were proved only on 2.0.1 and 0.84.0.

The maintainer adopted the current models in this repository's own profiles on
2026-09-30 (#292). This Spec makes the shipped binary agree with that
decision, and gives it one place where the next model change is made.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Agent Models keep
  the identifiers each adapter advertises, and the snapshot is named by its
  date. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change is static data, local
  configuration and documentation. No credential is read and no network call
  is added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0180 (this Spec) derives built-in
  selections from one dated Recommended Profile. ADR-0037 keeps Roundfix the
  owner of model and reasoning selection, and ADR-0049 keeps a profile atomic
  with a non-empty Fallback Chain; the Recommended Profile has both for every
  category. ADR-0050 and ADR-0140 prove every configured tuple and substitute
  none, which is why a retired built-in fallback stops every Run. ADR-0107
  proves only the categories a configuration defines, so the built-ins still
  define the five required categories only. ADR-0151 and ADR-0153 bind
  `profiles.review` to the pre-PR review provider's runtime, so the recommended
  `review` profile stays on Codex. ADR-0079 and ADR-0147 keep advertised
  identifiers opaque and let the adapter's refusal stand; the catalog stays
  picker data, never an allowlist. ADR-0069 keeps the Baseline semantic
  analysis read-only and supervised, and ADR-0180 replaces only its two model
  names. This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096,
  ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and ADR-0093 and ADR-0094 check
  its consistency by citation and artifact presence. ADR-0166 records a Task's
  undeclared paths, ADR-0167 keeps the pre-PR Pull Request row from deciding a
  qualifying partial, and ADR-0176 reads citations only from authored text.
  All hold.
  ADR-0097 cites ADR-0080 but carries a QA row forward on unmoved evidence.
  ADR-0114 cites ADR-0050 but decides what opening an Agent Session counts as.
  ADR-0165 and ADR-0169 cite ADR-0153 but decide what a blocking review after
  archive does and which diff the review reads.
  ADR-0168 cites ADR-0093 but sets the horizon of the related-ADR check.
  ADR-0182 (Spec 0190) cites ADR-0096 but moves mechanical facts to Task
  settlement. ADR-0183 and ADR-0184 (Spec 0191) cite ADR-0093 and ADR-0156 but
  add receipts for attributed claims and transcripts for command surfaces.
  ADR-0181 (Spec 0196) cites ADR-0180 but compares a configuration with the
  recommendation, which this Spec leaves out.
  This Spec changes none of the nine, so they do not apply.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se necessário"
  for the skill files, and "Autorizar os dois" for `.roundfixrc.yml`), with
  the standing grant of 2026-09-21 for the two governed test files, recorded
  in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/agents/openai.yaml`,
  `skills/roundfix/agents/openai.yaml`, `internal/cli/cli_test.go`,
  `internal/docscontract/publicdocs_test.go`, `.roundfixrc.yml`. Sanctioned
  regeneration: `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A repository on built-in profiles, or on a generated config, still passes
  Preflight after `gpt-5.5` leaves Codex on 2026-10-14.
- The picker offers the models and reasoning efforts the adapters advertise
  today, and nothing they have retired.
- One dated Recommended Profile per Agent Work Category is what Roundfix shows
  and what its built-ins are.
- The reference document and the binary state the same snapshot, and a test
  fails when they stop agreeing.

## Core Features

1. **The Model Catalog follows today's adapters.** The Codex catalog lists
   `gpt-6.1-sol`, `gpt-6-astra`, `gpt-6-sol`, `gpt-6-luna`, `gpt-5.6-sol`,
   `gpt-5.6-terra` and `gpt-5.6-luna`, in that order. The Claude catalog lists
   `opus`, `sonnet`, `claude-fable-5-1`, `haiku` and `default`, and describes
   `opus` as Opus 5.5 and `sonnet` as Sonnet 5.5. The picker's reasoning
   efforts are the ones every catalog model of the runtime advertises: `low`
   to `max` for Codex, and `default` plus `low` to `max` for Claude.
2. **The adapter floors are the proven versions.** Adapter Readiness requires
   `codex-acp` 2.0.1 or newer and `claude-agent-acp` 0.84.0 or newer. Setup
   and Doctor name those versions in their install actions.
3. **The recommendation is a dated Recommended Profile.** Roundfix ships one
   Preferred Selection and one Fallback Chain for each of the ten Agent Work
   Categories, under the snapshot date 2026-09-30, with a rationale for every
   selection. `roundfix profiles show` prints that profile for each category
   in place of the five-row ranking. It stays advisory.
4. **Built-in selections are the Recommended Profile.** The built-in profiles
   of the five required categories, the config `roundfix setup` and
   `roundfix init` generate, and the legacy Codex runtime default all come
   from the Recommended Profile. No production code names `gpt-5.5`.
5. **The Baseline semantic analysis uses current models.** It tries
   `gpt-6.1-sol` and then `gpt-5.6-sol`, both at `xhigh`. Its read-only,
   supervised contract does not change.
6. **The reference document states the shipped snapshot.**
   `docs/references/model-selection.md` is refreshed to 2026-09-30 with the
   sources of every figure, and a test requires its date and its recommended
   selections to match the binary.

## Declared breaks

- `roundfix profiles show --json` moves from schema `roundfix/profiles/v1` to
  `roundfix/profiles/v2`. Recommendation rows lose `benchmark`,
  `result_percent`, `average_cost_usd` and `category_specific`, and gain
  `role`. A profile loses `recommendation_source`.
- The built-in `review` profile leads with `codex / gpt-5.6-luna / max`
  instead of copying `general`. Spec 0041 kept Luna out of the generated
  defaults; ADR-0180 reverses that.
- The built-in `frontend` Preferred Selection uses effort `high`, not `xhigh`.
- Every built-in implementation profile now falls back to Claude. The built-in
  `frontend` profile already required the Claude runtime.
- Adapter Readiness refuses `codex-acp` below 2.0.1 and `claude-agent-acp`
  below 0.84.0, and names the install command.
- The picker no longer offers `gpt-5.5`, `gpt-5.4`, `gpt-5.4-mini`,
  `gpt-5.3-codex-spark`, `claude-fable-5` or the Claude effort `maximum`.
- A legacy config that relies on the built-in Codex runtime default resolves
  to `gpt-6.1-sol` with `high`, not `gpt-5.5` with `xhigh`.

## Non-Goals / Out of Scope

- Comparing a configured profile with the recommendation, warning on
  `roundfix upgrade`, applying the recommendation, or pinning a deliberate
  deviation. Spec 0196 adds those on top of the Recommended Profile.
- Per-model reasoning efforts in the picker. `ultra` stays reachable through a
  configured profile or a complete override, where Exact Agent Selection Proof
  decides.
- A catalog for the `opencode` runtime, a fourth ACP Runtime, or fetching a
  model list from the network.
- The Baseline decision `runtime.backend`, whose default names `gpt-5.6-sol`.
  That model is not retiring, and the decision belongs to the Baseline.
- A measured A/B before adoption. The maintainer chose direct adoption with
  a revert if Runs get worse.

## Success Metrics

1. Success Metric: with no User Config and no Project Config, the effective
   profile of each required category equals the Recommended Profile of that
   category, and no tuple names `gpt-5.5`.
2. Success Metric: the Model Catalog contains every Codex and Claude model any
   Recommended Profile names, and none of the five retired or replaced
   identifiers. Removing a recommended model from the catalog fails a test.
3. Success Metric: `roundfix profiles show --category docs --json` reports
   schema `roundfix/profiles/v2` and two recommendation rows: the preferred
   `codex / gpt-5.6-luna / max` and the fallback `claude / sonnet / high`.
4. Success Metric: Doctor refuses a Codex adapter that reports version 1.1.5
   and names `@agentclientprotocol/codex-acp@2.0.1` in its install action.
5. Success Metric: changing the snapshot date in the binary without changing
   `docs/references/model-selection.md` fails the documentation contract.

## Recorded limits

- The catalog and the floors are what the installed adapters advertised on
  2026-09-30. They are static data: the next adapter release can outdate them,
  and the next snapshot is a new release.
- The picker's effort list is per runtime. A model that advertises less than
  its runtime's list, such as `haiku`, is refused at proof with the adapter's
  own message.
- Published evidence for `gpt-6.1-sol` is one day old and mixed. The
  Recommended Profile rests on the maintainer's adoption decision and on the
  Exact Agent Selection Proof of every tuple, not on a benchmark.

## Decisions

- **One snapshot, many readers.** See ADR-0180.
- **`gpt-5.5` leaves the picker now.** It is still advertised until
  2026-10-14, but a picker should not offer a model with two weeks left. A
  non-interactive custom value still reaches it until then.
- **Rationale, not figures, in the binary.** Benchmark figures stay in the
  reference document with their sources, where a missing figure can be written
  as missing.
- **Optional categories are recommended but not built in.** All ten categories
  have a Recommended Profile, and only the five required ones are built-in
  profiles, so Profile Readiness proves no tuple a repository did not ask for.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- OpenAI's model page for ChatGPT and Codex
  (<https://learn.chatgpt.com/docs/models>, read 2026-09-30) states that
  GPT-5.5 retires from Codex with ChatGPT sign-in on 2026-10-14 and that
  `gpt-5.4` and `gpt-5.4-mini` retired on 2026-08-31.
- The adapters' own advertisement, read on 2026-09-30 through a refused
  `roundfix profiles configure --dry-run` on `codex-acp` 2.0.1 and
  `claude-agent-acp` 0.84.0:
  - Codex models: `gpt-6.1-sol`, `gpt-6-astra`, `gpt-6-sol`, `gpt-6-luna`,
    `gpt-5.6-sol`, `gpt-5.6-terra`, `gpt-5.6-luna`, `gpt-5.5`.
  - Codex reasoning efforts: `low`, `medium`, `high`, `xhigh`, `max`, plus
    `ultra` on `gpt-6.1-sol` and not on `gpt-5.6-luna`.
  - Claude models: `default`, `opus`, `claude-fable-5-1`, `sonnet`, `haiku`,
    `claude-sonnet-5`, `claude-opus-5`, `claude-fable-5`, `claude-opus-4-8`,
    `claude-opus-4-7`, `claude-opus-4-6`, `claude-sonnet-4-6`.
  - Claude reasoning efforts: `default`, `low`, `medium`, `high`, `xhigh`,
    `max`.
- The proof of the adopted tuples recorded in #292: `roundfix profiles
  configure --scope project --file … --dry-run` exited 0 for every tuple of
  the Recommended Profile on those adapter versions.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and two queries:
`qmd query "modelo recomendado por categoria roundfix perfis de seleção de
agente catálogo defasado"` and `qmd query "aviso de modelo desatualizado
upgrade notice configuração recomendada"`. The strongest hit,
`wiki/concepts/modelos-custos-e-selecao.md`, tracks the OpenRouter model of
another agent and does not govern Roundfix routing. The model research of
2026-09-30 is the entry
`inbox/secondbrain/2026-09-30-modelos-de-codificacao-ago-set-2026-e-runtimes-cursor-grok.md`;
it supplies the release dates and published figures the reference document
will cite, and it marks `gpt-6.1-sol` as unmeasured. Exa located OpenAI's
model page with the retirement dates above, which set this Spec's deadline,
and a release log stating that `gpt-5.3-codex-spark` was deprecated on
2026-09-14.

## Technical candidate

The [_techspec.md](_techspec.md) records the data, the implementation map,
coverage and build order.
