---
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# A transient profile-proof timeout does not park the Delivery Queue

A profile proof whose disposable Agent Session setup times out now gets a
third attempt after a backoff of 5 and then 10 seconds, its setup deadline
becomes the config key `defaults.profile_proof_setup_timeout` (default `30s`),
and the cleanup of a Session that never opened stops logging a misleading
warning; every configured tuple is still proven.

Minimal PRD for a bug fix, minted by `write-techspec`. It carries the contract
the downstream skills parse; the TechSpec holds the design.

Between 2026-10-19 19:30 and 2026-10-20 00:27 BRT, the Delivery Queue on the
oraculum machine parked five times on `profile proof setup timed out twice:
context deadline exceeded; classification: temporary`. The timeouts fell on
different runtimes and models, in the implement Preflight and the pre-PR
review, at a load average of about 1.6, and the same tuple passed each time on
the next `roundfix deliver retry`. Profile readiness proves every configured
tuple, and each attempt has a fixed 30-second setup deadline, so one slow
Session setup among five or six tuples parks the item and costs a Supervisor
cycle of about 30 minutes. One of those failures also logged `disposable Agent
Session cleanup found no Session to close`, with a recovery step that did not
apply.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md`](references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md),
   origin oraculum, type hint `fix`. It records the five parks, the message,
   the cleanup warning, the Delivery Owners, the versions and the three
   remedies it proposes.

No open Backlog Entry or unresolved Finding shares this context: `docs/backlog/`
and `docs/findings/` hold no entry about profile proof, Agent Session setup or
Preflight timeouts.

## Prerequisites

None. This Spec is authored against `9448aa64` (v0.95.0). Specs 0286 and 0287
are authored beside it and deliver after it; neither changes the profile
proof, `internal/config` or the `profiles` and `runtime` references. Another
Spec merged first may raise the Roundfix skill version, the Behavior Surface
Record or the coverage record, so this plan is rebased and checked again
before its Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new config key and the messages are fixed strings. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Every test uses an injected Agent runner
  or the existing compiled fake acpx, and no test starts a real ACP agent.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the proof, its tests and its
  documentation follow recorded decisions. ADR-0298 (this Spec): "A timed-out
  profile proof gets a third attempt, and its setup deadline is
  configurable". ADR-0050: "Configured fallbacks activate after
  notification". ADR-0114: "Opening an Agent Session is not Agent work".
  ADR-0140: "An Agent Selection is proven as the exact advertised tuple".
  ADR-0107: "Profile readiness covers every configured Agent Work Category".
  ADR-0184: "A TechSpec states a command surface as a transcript". ADR-0259:
  "A test runs in parallel unless it names why it cannot". ADR-0271: "A test
  bounds by real time only what it asserts". ADR-0285: "Every test file
  compiles for Windows, and a repository contract proves it". ADR-0250: "A
  module version is chosen when it is recorded, and the coverage record lists
  every platform". ADR-0189: "An owned skill's version names its content, and
  its minimum is the bundle". ADR-0256: "A release waits for the skills that
  describe a changed surface". ADR-0280: "A guide Task records its own surface
  fingerprint, a PRD records who decided, and the glossary is held to its own
  rule". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0240:
  "One QA partial policy, and rows a Run sandbox cannot reach". ADR-0276: "A
  row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "A Spec proves what it claims,
  and every active Spec is held to it". ADR-0183: "An attributed claim
  carries a receipt the check proves". Set no obligation here, because this Spec changes none of their rules: ADR-0069 (baseline semantic analysis is read-only and supervised); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0180 (built-in selections derive from one dated Recommended Profile); ADR-0181 (a configuration is compared with the recommendation only where a person asked); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0217 (a Cursor selection names the advertised model value, and the login stays the maintainer's); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0238 (a light tier runs low-complexity Tasks on an open model); ADR-0242 (a blocked review names the protocol step and retries once before the prompt); ADR-0245 (a lost rollout is runtime infrastructure); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0257 (inside a Run the Daemon is the only full gate, and the Baseline states each rule once); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0299 (the history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0274 (a turn that did not finish is not accepted); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof); ADR-0295 (reachability reads the fetched default branch, and origin names it when origin/HEAD is unset). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AK approval of 2026-10-20 through AskUserQuestion, option "Sim, mesmas
  regras do AJ (Recomendado)", which keeps the cycle AJ grant "Mexer em
  Governed Paths declarados" ("'Concedo' para os caminhos que cada Spec
  declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/profiles.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `docs/references/coverage-record.json`, `skills/roundfix/SKILL.md`. No
  Makefile, workflow, Go module, lint, formatter, Baseline module or Project
  Config file change.

## Goals

- A proof whose setup times out once or twice under transient load passes
  without parking the item, in every profile-readiness consumer.
- An operator on a slow machine can lengthen the setup deadline in config
  instead of retrying by hand.
- A tuple that really cannot open a Session still fails, classified
  `temporary`, in bounded time.
- Every configured tuple, preferred and fallback, is still proven before a
  Run.
- The cleanup of a Session that never opened writes nothing misleading.

## Core Features

1. The configurable setup deadline. `defaults.profile_proof_setup_timeout`,
   default `30s`, a duration greater than zero and at most `10m`, read from
   User Config and Project Config and passed by every profile-readiness
   consumer to each proof attempt.
2. The third attempt. A setup timeout is retried after 5 seconds and again
   after 10 seconds; the third timeout keeps the `temporary` classification
   and names the attempts, the deadline and the config key.
3. Silent cleanup of an absent Session. A disposable Session that acpx reports
   absent at cleanup is closed without a warning.
4. Documentation. The configuration guide, the `profiles` command guide, the
   Roundfix `profiles` and `runtime` references and **Proof Setup Deadline**
   in `CONTEXT.md` describe the key, the attempts and the backoff.

## Non-Goals / Out of Scope

- Proving only the Preferred Selection, or deferring fallback proofs to Run
  time: the maintainer kept every tuple proven.
- Retrying any failure other than a setup timeout, or retrying after the
  command's own context ends.
- Changing the Fallback Chain probe made after Run creation, the cleanup
  deadline, the classification of other proof failures, the Delivery Queue's
  Park Classes or `roundfix deliver retry`.
- Adding the key to the `roundfix init` template or to this repository's
  `.roundfixrc.yml`.

## Success Metrics

1. Success Metric: with an injected runner, a tuple whose setup times out on
   its first two attempts passes on the third, after waits of 5 and 10
   seconds recorded by an injected sleeper, through `profiles validate` and
   the operational Preflight, with nothing on standard error.
2. Success Metric: a tuple that times out three times makes `roundfix
   profiles validate --category backend` exit `2`, classified `temporary`,
   with the exact standard error the TechSpec's API Contract 2 fixes, naming
   the attempts, the deadline and the config key.
3. Success Metric: `defaults.profile_proof_setup_timeout: 90s` in Project
   Config reaches every proof attempt as a 90-second setup deadline, and the
   agent's setup context carries that deadline; a value outside the range
   makes the built `roundfix profiles show` reproduce the TechSpec's Surface
   Transcript 1 whole.
4. Success Metric: a proof whose setup never opened its Session returns its
   setup error alone and logs no cleanup warning.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- The adopted Inbox Entry, a measurement this Spec did not design: five parks
  on the oraculum machine between 2026-10-19 19:30 and 2026-10-20 00:27 BRT,
  each cleared by the next `roundfix deliver retry` with the same tuple.
- Google, "Site Reliability Engineering" (2016), chapter "Addressing Cascading
  Failures" (https://sre.google/sre-book/addressing-cascading-failures/):
  "Limit retries per request. Don't retry a given request indefinitely." and
  "Always use randomized exponential backoff when scheduling retries."
- Amazon Builders' Library, "Timeouts, retries, and backoff with jitter"
  (https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/):
  "The most common pattern is an exponential backoff, where the wait time is
  increased exponentially after every attempt." and "In almost all cases, our
  solution is to limit the number of times that the client retries".

## Unreachable Acceptance

None. Every criterion is proven with injected runners, sleepers and the
compiled fake acpx; no criterion needs the network, a provider, Docker or an
open Pull Request.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Secondbrain, read only. Read `wiki/index.md`, then ran `qmd query "profile
  proof timeout retry backoff agent session setup deadline preflight" --all
  --files --min-score 0.3` and `qmd query "retry with backoff for transient
  timeouts exponential jitter" --all --files --min-score 0.3`. The first
  returned the mirrors of this repository's `runtime` and `profiles`
  references, the `profiles` guide, ADR-0114, ADR-0140, ADR-0107 and
  `CONTEXT.md`; none records a decision on the retry count or the deadline.
  The second returned
  `raw/files/books/site-reliability-engineering-2016/27-addressing-cascading-failures.md`,
  whose retry passages are quoted above. No Inbox Entry was written, because
  the briefing forbids writing there.
- Exa: the Amazon Builders' Library article quoted above. Supporting: it
  names a bounded retry count and exponential backoff. Challenging: both
  sources recommend jitter. They address many clients retrying against one
  server; a Queue Owner proves its tuples serially on one machine, and the
  same article notes that Amazon picks scheduled-work jitter by "a consistent
  method that produces the same number every time on the same host", so a
  deterministic wait is kept (ADR-0298).
- Characterization on 2026-10-20 at `9448aa64`, recorded in the TechSpec.

## Glossary

- adds: **Proof Setup Deadline**

## Decisions

- 2026-10-20 (about 06:30 BRT): the cycle AK decisions through
  AskUserQuestion — next cycle "Duas Specs + um lote (Recomendado)"; this
  Spec's design "Retry interno + prazo configurável (Recomendado)" (option
  text: "Um timeout classificado 'temporary' ganha uma nova tentativa
  automática dentro do preflight antes de estacionar, e o prazo de criação da
  sessão passa a ser configurável (documentado em profiles.md). Continua
  provando todas as tuplas."); authorization "Sim, mesmas regras do AJ
  (Recomendado)".
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the operator's cycle AK briefing for Spec 0285 of 2026-10-20 and the characterization in the TechSpec.

Taken from those sources, with no further interview:

- Three attempts in all, with waits of 5 and then 10 seconds, and the
  `temporary` classification kept. See ADR-0298.
- The key is `defaults.profile_proof_setup_timeout`, default `30s`, greater
  than zero and at most `10m`. See ADR-0298.
- The cleanup warning is a symptom of the same timeout path, so an absent
  Session at cleanup is silent. See ADR-0298.

## Open Questions

None.
