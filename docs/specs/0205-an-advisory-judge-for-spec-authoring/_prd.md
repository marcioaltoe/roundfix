---
spec: 0205-an-advisory-judge-for-spec-authoring
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# An advisory judge for Spec authoring

The Spec Consistency Check proves that a cited ADR is listed, that a claim
shares words with the record it names, and that every PRD goal has a Coverage
Map line. It cannot tell whether the record actually says what the Spec
attributes to it, or whether the TechSpec section a Coverage Map line names
describes a way to reach the goal. Both defects reach review, where a reviewer
settles them only by opening the cited text.

On 2026-09-30 five TypeSafe Jev judgments were measured against this
repository's Spec archive. The finding that records it,
[2026-09-30-jev-judgments-measured-against-the-spec-archive.md](../../history/findings/2026-09-30-jev-judgments-measured-against-the-spec-archive.md),
passed two of them as advisory checks:

- **Citation support.** Does the cited ADR say what the sentence attributes to
  it? AUROC 0.89 against the word-overlap baseline's 0.80. Raised at its
  operating point on 3 of 200 real citations, and all three were a known false
  citation.
- **Goal to mechanism.** Does the TechSpec section a Coverage Map line names
  describe a mechanism for the PRD goal? AUROC 0.86 against 0.77. Raised on 3
  of 143 real mapped sections, with modest recall.

Task Type, the overlap shortlist and Task lint were not adopted.

This Spec adds `roundfix spec judge`, a command that asks those two questions
about one Spec and prints what it raises. The `write-prd` and `write-techspec`
skills run it after the Spec Consistency Check, and the authoring model
answers each raised judgment. It never gates anything, and it fails open.

On 2026-10-01, before delivery, the maintainer moved the judge's primary
transport to OpenRouter: "Podemos consumir através do openrouter. Estou muito
disposto a seguir por esse caminho", with a key of Roundfix's own: "Gosto de
ter chaves separadas para determinar o real custo de cada trabalho/projeto no
openrouter". The command calls Jev through OpenRouter's System One API when
`ROUNDFIX_OPENROUTER_API_KEY` is set, and calls TypeSafe directly only
when that key is absent and `ROUNDFIX_TYPESAFE_API_KEY` is set.

## Prerequisites

This Spec is delivered after Spec 0194, which splits the Roundfix Skill and
the command reference into one file per command. Its third Task edits the
`spec` Skill reference that Spec 0194 creates, `.agents/skills/roundfix/references/spec.md`,
and creates the `spec` user-guide command guide, `docs/user-guide/commands/spec.md`,
which Spec 0194 did not create, with its row in the command index. Its
Verification reads them. The Delivery Queue does not enforce this order, so
the operator queues Spec 0194 first.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. A judgment is
  named by its fixed kind (`citation-support`, `goal-mechanism`), a request by
  the SHA-256 of its state, and a Judge Log file by its UTC calendar month.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the command sends HTTPS requests with
  a bearer key to OpenRouter's System One endpoint when
  `ROUNDFIX_OPENROUTER_API_KEY` is set, or else to the TypeSafe endpoint
  when `ROUNDFIX_TYPESAFE_API_KEY` is set. Each key is read only from the command's
  environment, sent only in the authorization header of its own endpoint, and
  never printed, logged or stored; the generic `OPENROUTER_API_KEY` is never
  read; without either key the command sends nothing. The request and
  response fields are the ones TypeSafe's published API reference defines,
  which OpenRouter's System One API implements and extends with `id`,
  `provider` and `usage.cost`. A request carries only this repository's Spec
  artifacts. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0200 (this Spec) decides that the
  judge is advisory, fails open and compares thresholds only for the pinned
  model version. ADR-0201 (this Spec) decides what a request may carry, that it
  goes to OpenRouter first and TypeSafe second, where each key comes from,
  that every call is logged, and the monthly ceiling.
  ADR-0116 reads the cited record for a claim, and ADR-0183 adds the Claim
  Receipt; both stay unchanged and the judge complements them. ADR-0176 reads
  only authored text for citations, and the judge reads only the authored PRD,
  TechSpec and ADRs. ADR-0035 lets the Spec Root live outside the repository,
  and the judge reads the Spec where the Spec Root resolves. ADR-0089 passes
  the environment to code under test explicitly, which is how every test
  supplies the key and the transport. ADR-0184 has the TechSpec state the new
  command as Surface Transcripts. ADR-0187 splits the Roundfix Skill and the
  command reference by command, and the new command is described in the
  `spec` files. ADR-0189 ties an owned skill's version to its content, so each
  edited skill raises its version. ADR-0193 names a Spec's prerequisites; this
  Spec states Spec 0194 as one in its PRD and Build Order. ADR-0081 and
  ADR-0149 sanction the regeneration that follows the skill edits. ADR-0182
  runs Settlement Checks before each Task commit, ADR-0178 authorizes a Task
  commit by the grant it ran under, and ADR-0166 records undeclared paths;
  every Task here declares its paths. ADR-0167 keeps the pre-PR Pull Request
  row from deciding a qualifying partial; this gate aims at `pass`. This
  Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104,
  ADR-0117, ADR-0155 and ADR-0156, and ADR-0093 and ADR-0094 check its
  consistency by citation and artifact presence. ADR-0036 cites ADR-0035 but
  decides how review artifacts are committed, ADR-0192 cites ADR-0149 but
  decides how a conflict on declared derived paths is resolved, ADR-0097
  cites ADR-0080 but carries a QA row forward, and ADR-0168 cites ADR-0093 but
  narrows the related-ADR check. ADR-0029 and ADR-0142 cite ADR-0036 but
  decide where review artifacts live and what decides a watch outcome, and
  ADR-0195 cites ADR-0097 but decides when a QA row is observed again. This
  Spec changes none of these seven, so none of them applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se necessário"),
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- An author learns, before review, which ADR attributions in a PRD or TechSpec
  the cited record may not support, and which Coverage Map lines point at a
  section that may not deliver the goal.
- The judge never blocks authoring and never changes the outcome of another
  command: every failure to judge ends as a named skip.
- A request carries only this repository's Spec artifacts in English, and a
  key never leaves the environment except in the authorization header of its
  own endpoint.
- Every call is on record with its answer and cost, and spending stops at
  US$5 per calendar month.
- The questions, thresholds and model version the measurement validated are
  the ones the command uses, kept in one reviewable place.

## User Stories

1. As an authoring model finishing a PRD, I want the ADR attributions the
   cited records may not support listed with the answer and its confidence, so
   that I correct or defend each one before the TechSpec builds on it.
2. As an authoring model finishing a TechSpec, I want each Coverage Map line
   whose named section may not deliver its goal listed with the probability,
   so that I fix the mapping or the section before decomposition.
3. As a maintainer, I want the command to skip, not fail, when no Jev key is
   set, the service is down, the ceiling is reached or the Spec is not in
   English, so that authoring never waits on a paid external service.
6. As a maintainer, I want Roundfix's Jev calls billed to an OpenRouter key of
   their own, so that I can see what this work costs apart from every other
   project on the same account.
4. As a maintainer, I want a monthly log of every call with its answer, tokens
   and cost, so that I can audit what was sent, re-measure thresholds and see
   the month's spend.
5. As a maintainer reviewing the data boundary, I want the command to be
   unable to send anything but the Spec's PRD, TechSpec and ADR text, so that
   the authorization holds without anyone remembering it.

## Core Features

1. **`roundfix spec judge <slug>` judges one active Spec.** With `--stage prd`
   it judges the PRD's citations; with `--stage techspec` the TechSpec's
   citations and its Coverage Map lines; with no stage, both. It prints each
   raised judgment and one summary line, as text or `--format json`.
2. **Citation support.** Every sentence of the PRD or TechSpec that cites
   exactly one accepted ADR with an attribution verb is asked whether the
   ADR's text supports what the sentence attributes to it. The judgment is
   raised when the answer is not `supports` at confidence 0.8 or more.
3. **Goal to mechanism.** Every Coverage Map line that maps a PRD goal to a
   named TechSpec section is asked whether that section describes a mechanism
   that would deliver the goal. The judgment is raised when the probability is
   below 0.3.
4. **The measured questions, one reviewable source.** The question texts,
   answer criteria, thresholds, extraction rules, model version, price and
   ceiling are the ones the 2026-09-30 measurement used, kept together in one
   file that the TechSpec reproduces byte for byte.
5. **One pinned model.** Requests name Jev 1.13 by the ID each endpoint
   accepts: `jev-1.13` on OpenRouter and `jev-1.13.0` on TypeSafe. An
   answer is compared with a threshold only when the model it reports
   normalizes to Jev 1.13: `jev-1.13.<patch>` from TypeSafe, or
   `typesafe/jev-1.13` with an optional `-YYYYMMDD` snapshot suffix from
   OpenRouter. An answer from any other model is shown as skipped with the
   model that answered, logged, and never compared with a threshold.
6. **Advisory and fail-open.** The command exits `0` whenever it ran, whether
   it raised judgments or not. A missing key, a service or network failure, a
   reached ceiling, an unreadable Judge Log and a refused request each end as
   a named skip, and the summary says how many judgments were not asked. Only
   a usage error, such as an unknown Spec, exits `2`.
7. **Only Spec artifacts, only English.** A request's state is built only
   from the Spec's own PRD and TechSpec and this repository's accepted ADRs (status `accepted`, or a legacy ADR without lifecycle front matter that the repository treats as active; never `proposed`, `rejected`, `deprecated` or `superseded`),
   read as regular files inside their directories. An artifact that is not
   English is reported as skipped and sent nowhere.
8. **OpenRouter first, on a key of its own.** When
   `ROUNDFIX_OPENROUTER_API_KEY` is set, every request of the run goes to
   OpenRouter's System One API with that key; otherwise, when
   `ROUNDFIX_TYPESAFE_API_KEY` is set, to TypeSafe directly with that key; otherwise
   the run is skipped with a reason that names
   `ROUNDFIX_OPENROUTER_API_KEY`. The generic `OPENROUTER_API_KEY` is
   never read. No key appears in any output, log or file.
9. **Every call is logged.** Each request appends one line to the Judge Log
   in Roundfix Home, one file per UTC month: the time, repository, Spec,
   judgment, artifact and line, state hash, question, the transport, the
   answer with its probabilities and confidence, latency, tokens, the model
   that answered, the cost (OpenRouter's reported `usage.cost`, or else input
   tokens at the pinned price) and the outcome.
10. **A monthly ceiling.** Before each request the command sums the month's
    logged cost across both transports. At or above US$5 it asks nothing more
    and reports `monthly ceiling reached`.
11. **The authoring skills answer the judge.** `write-prd` and
    `write-techspec` run the command under a new heading after their checker
    step, and tell the authoring model to answer each raised judgment by
    correcting the artifact or by stating in its report why the text stands.
    The Roundfix Skill and the command reference describe the command.

## User Experience

The author runs the command, or the skill runs it, and reads a few lines: one
per raised judgment, naming the file and line, the ADR or the goal and
section, the answer and its confidence or probability, and the claim text.
Clear judgments are counted in the summary, not listed. When nothing can be
asked, one summary line says why. The JSON form lists every judgment.

## Declared breaks

- `roundfix spec --help` and the top-level help list one more command.
- Roundfix makes its first outbound network request, and only from this
  command.

## Non-Goals / Out of Scope

- Gating anything: no exit code of `spec check`, `implement`, `deliver` or
  the QA gate depends on a judgment.
- Task Type, the overlap shortlist and Task lint, which the measurement did
  not adopt.
- Judging Task files, findings or Backlog Entries, although the maintainer's
  authorization covers them.
- Portuguese or any other non-English Spec.
- Calls through a metering gateway, or through any service other than
  OpenRouter's System One API and the direct TypeSafe endpoint.
- Falling back to the generic `OPENROUTER_API_KEY`, or switching transport in
  the middle of a run.
- An agent or a skill for the judge.
- Re-measuring thresholds or moving the pinned model.
- Any live call to OpenRouter or TypeSafe from a test, a Verification or the
  QA gate.

## Success Metrics

1. Success Metric: on a Spec whose TechSpec attributes to an ADR something
   that record does not say, and maps a goal to an unrelated section, the
   command prints both as `advisory`, with the answer and its confidence or
   probability, and exits `0`.
2. Success Metric: with no Jev key, the ceiling reached, the Judge Log
   unreadable, or the service failing, the command exits `0`, sends no further
   request, and names the reason.
3. Success Metric: every string a request carries comes from the Spec's PRD
   or TechSpec or an accepted ADR, and a Spec whose PRD is a symbolic link, or
   is not English, sends nothing.
4. Success Metric: an answer whose reported model does not normalize to Jev
   1.13 (`jev-1.13.<patch>`, or `typesafe/jev-1.13` with an optional
   `-YYYYMMDD` suffix) is never raised or cleared, and one that does, in
   either shape, is.
5. Success Metric: the question file the command embeds is byte-identical to
   the block in the TechSpec, and its texts and thresholds equal the
   measurement's.
6. Success Metric: every language gate decision on this repository's archived
   PRDs, TechSpecs and accepted ADRs is English, and on a Portuguese Spec it is
   not.
7. Success Metric: with both keys set every request goes to OpenRouter and
   carries only the OpenRouter key; with only `ROUNDFIX_TYPESAFE_API_KEY` every request
   goes to TypeSafe; with only the generic `OPENROUTER_API_KEY` nothing is
   sent.

## Recorded limits

- The thresholds belong to `jev-1.13.0`, measured through the direct TypeSafe
  API, and were chosen on the data they were reported on, so their real
  precision may be lower. OpenRouter's `typesafe/jev-1.13-20260917` was
  measured equivalent on the same states on 2026-10-01; a later dated 1.13
  snapshot, which the pin also accepts, was not.
- About 2% of advisory flags flip between two runs of the same states, on
  either transport, so a raised judgment can come and go between runs.
- Recall is modest: 69.5% of random citation swaps and 10.5% of nearest-ADR
  swaps were caught, and recall on real coverage defects was not measured.
- Two commands started at the same moment can together pass the ceiling by
  their own requests.
- A failed request is logged with no tokens and no cost.
- The ceiling sums two bills, the OpenRouter account's and the TypeSafe
  account's, so it bounds the machine's Jev spend, not either account's.

## Decisions

- **Advisory, pinned and fail-open.** See ADR-0200.
- **The data boundary, the transports, the keys, the log and the ceiling.**
  See ADR-0201.
- **OpenRouter first, on a dedicated key.** The maintainer chose OpenRouter
  on 2026-10-01, the comparison of that day found it equivalent to TypeSafe
  direct in quality and cost, and a key that only Roundfix uses keeps its Jev
  cost separable on the OpenRouter account; reading the generic key would
  blur that.
- **A command, not a skill or an agent.** The skills call a command, so the
  judge works the same from any runtime and needs no skill name.
- **The log lives in Roundfix Home, not in one repository.** The keys and the
  bills belong to the maintainer, not to one repository, so the ceiling sums
  every repository's calls on both transports.
- **Only the two measured judgments, asked as measured.** The extraction
  rules and truncation limits are the measurement's, because the thresholds
  hold only for the inputs they were measured on.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- **A measurement this Spec did not design.** The shadow benchmark of
  2026-09-30 sent 5,240 requests to `jev-1.13.0`, for US$0.218, with zero
  failures, and scored the judgments against labels drawn from 171 archived
  Specs and 146 accepted ADRs. Its question texts and thresholds are the ones
  this Spec adopts.
- **OpenRouter's System One API, published.** OpenRouter's TypeSafe SDK
  guide (<https://openrouter.ai/docs/guides/community/typesafe-sdk>) states
  that `POST https://openrouter.ai/api/v1/systemone` takes
  `Authorization: Bearer <OpenRouter API key>`, implements TypeSafe's request
  and response shapes (requests require `model`, `state` and `questions`;
  responses carry `model`, `answers` and `usage`) and adds `id`, `provider`
  and `usage.cost`; that the bare `jev-1.13` is routed as `typesafe/jev-1.13`
  and `jev-latest` as `~typesafe/jev-latest`; and that the response `model` is
  the OpenRouter ID that served the request, as in its example
  `typesafe/jev-1.13-20260917`. Its Jev hub
  (<https://openrouter.ai/docs/guides/community/jev>) states that Jev is
  billed to the OpenRouter account per input token with free output, that
  `usage.cost` is the response's cost in US dollars, that the context is
  32,000 tokens and that no TypeSafe account is needed. Its Jev tutorial
  (<https://openrouter.ai/docs/guides/community/jev-tutorial>) says the
  response names the dated snapshot that served the request, so a dated
  suffix is expected. The model page (<https://openrouter.ai/typesafe/jev-1.13>)
  prices it at US$0.042 per million input tokens with TypeSafe as the only
  provider, and the errors reference
  (<https://openrouter.ai/docs/api/reference/errors-and-debugging>) assigns
  `402` to a key or account without credits. All read 2026-10-01.
- **Both endpoints, probed.** On 2026-10-01 the maintainer's session probed
  both endpoints live, outside this Spec's Tasks: TypeSafe refused
  `jev-1.13` (HTTP 400, `Unknown model: jev-1.13`) and answered `jev-1.13.0`
  as `jev-1.13.0`; OpenRouter refused `jev-1.13.0` (HTTP 400,
  `Model typesafe/jev-1.13.0 does not exist`) and answered `jev-1.13` as
  `typesafe/jev-1.13-20260917`, provider `TypeSafe`, with a `usage.cost`
  equal to its input tokens at US$0.042 per million. The TechSpec records
  them under Measured outside evidence.
- **The transports, compared.** The same session asked both transports the
  885 unique benchmark states with this Spec's questions: 885 of 885 answered
  on each, citation support at AUROC 0.8900 (TypeSafe) and 0.8899
  (OpenRouter), goal to mechanism at 0.8583 and 0.8644, US$0.0269 each, p50
  latency 335 and 366 ms. The two transports' answers differ no more than the
  direct API's answers differ across days (advisory flag agreement 97.5 to
  98.2%, mean absolute probability change 0.020). This is the measured basis
  for OpenRouter as the primary transport.
- **The API, published.** The TypeSafe API reference defines the request's
  `state`, `model` and `questions`, the Choice answer's `choice`,
  `probabilities` and `confidence`, the Noul answer's `noul`, the response's
  `model` and `usage`, and the `401`, `422`, `429` and `529` errors
  (<https://docs.typesafe.ai/api.md>). The models page prices `jev-1.13.0` at
  US$0.042 per million input tokens with free output, advises pinning the
  versioned ID once thresholds are tuned, and says English is where accuracy
  is best (<https://docs.typesafe.ai/models.md>). The confidence page defines
  confidence as a statistic of a Choice's distribution that a Noul does not
  carry (<https://docs.typesafe.ai/confidence.md>). The citation-check cookbook
  gates a citation verdict at confidence 0.8
  (<https://docs.typesafe.ai/cookbooks/citation_check.md>). All read
  2026-09-30.
- **The skill name rule, published.** The Agent Skills specification allows
  only lowercase letters, digits and hyphens in a skill's `name`
  (<https://agentskills.io/specification>), so `the_judge` is not a valid
  skill name, while a Claude Code subagent name only excludes `:`
  (<https://code.claude.com/docs/en/sub-agents>). Read 2026-09-30.
- **A language measurement.** The share of 29 English function words is at
  least 0.103 in every archived PRD and TechSpec and 0.137 in every ADR of
  this repository, and 0.008 at the median of the Secondbrain's Portuguese
  wiki pages, measured locally on 2026-09-30 with nothing sent anywhere.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the query `qmd query
"juiz semântico consultivo para autoria de especificações citação de ADR
TypeSafe Jev" --all --files --min-score 0.3`. The digest
`wiki/sources/digest-roundfix-orquestradores-e-ecosistema-jev-o-que-se-transfere-ao-roundfix-2026-09-30.md`
concludes that Jev serves Roundfix as an advisory check of artifacts already
written, never as a mechanical gate, and describes a citation verifier in
which code proves the quote exists and Jev only judges support; that shaped
the advisory decision and the complement to the Claim Receipt. The entity page
`wiki/entities/jev-typesafe.md` records the first-hand latency of 498 ms at
p50 through a gateway, against 330 ms direct in the measurement, which first
kept the command off any gateway. Exa and Context7 found the TypeSafe pages,
the Agent Skills specification and the Claude Code subagent reference cited
above; the skill name rule is why the design creates no skill.

For the 2026-10-01 amendment the query `qmd query "Jev via OpenRouter System
One API custo latência gateway" --all --files --min-score 0.3` returned the
same entity page first. It lists OpenRouter (`typesafe/jev-1.13`) among Jev's
providers with the same questions and answers, the account and the billing
being the only difference; that is why the pin accepts OpenRouter's model ID
for 1.13 and nothing more. An advisory command run once per authoring stage
tolerates the extra hop. Exa found OpenRouter's model page, its System One
API reference and its Jev tutorial, and the OpenRouter guides cited above were
read with a plain HTTP fetch.

## Open Questions

- The maintainer chose `the_judge` as the name of an agent or skill for the
  judge. This design creates neither, so the name is unused; a future skill
  would need a hyphenated name such as `the-judge`. Default until answered: no
  skill and no agent.

## Technical candidate

The [_techspec.md](_techspec.md) records the question file, the extraction
rules, the command transcripts, coverage and build order.
