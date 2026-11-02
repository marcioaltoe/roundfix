---
spec: 0338-the-release-decision-travels-with-the-spec
status: active
created: 2026-11-01
surfaces: [backend, cli, docs]
type: feat
---

# The release decision travels with the Spec

A maintainer decides two release facts while a Spec is planned, and Roundfix
loses both before the release. The Delivery Queue titles every pull request
`feat: deliver <slug>`, so a fix Spec squashes as a feature and
`roundfix release plan` refuses the patch the maintainer approved; Argus had
to ship 0.4.0 instead of 0.3.1. And when the maintainer authorizes a release
with the Spec, the release plan still stops at `approval_required` and a new
session asks again, because no rule lets the Spec's authorization record
answer it. After this Spec the PRD declares the Spec's Delivery Type, the
Archive Record keeps it and the queue's pull request title uses it, and an
approved authorization record that grants `release` can carry a Release
Ceiling that answers the release plan's approval question while the fresh
plan stays at or below it and nothing is breaking or classified by hand.

## Sources

Two Backlog Entries, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-01-the-delivery-pull-request-type-follows-the-spec.md`](references/2026-11-01-the-delivery-pull-request-type-follows-the-spec.md),
   `type: fix`, Argus, ai-memory message `01a193d4`.
2. [`references/2026-11-01-a-spec-authorization-carries-the-release-decision.md`](references/2026-11-01-a-spec-authorization-carries-the-release-decision.md),
   `type: feat`, Argus, ai-memory message `01a19920`.

Both change how the release plan reads what a delivered Spec decided, so one
Spec owns them, as the maintainer's Scope answer states. The two other open
Backlog Entries of 2026-11-01 belong to Specs 0339 and 0340.

## Prerequisites

None in `requires`. Authored against `dd3b2b62` (v0.128.0), where no other
Spec is active. Specs 0339 and 0340 are authored in parallel and deliver after
this one; `_techspec.md` → Build Order names the files they may share.

## Project Constraints

- Identifier strategy: not applicable — no identifier or slug is minted; the
  PRD and the Archive Record gain one optional `type` key, the authorization
  record two optional keys, and the release plan JSON one optional object
  inside `roundfix.release-plan/0.0.1`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the release plan and the queue title read local Git objects.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0143 keeps release planning
  read-only, ADR-0143: "A generic release request must never silently
  authorize a version decision"; the Release Ceiling is an explicit, recorded
  decision and the plan stays read-only. ADR-0309 fixes the queue title's
  shape, ADR-0309: "`feat: deliver <slug>` is used when it is at most 91
  long"; this Spec amends it so the type comes from the Spec. ADR-0302, which
  ADR-0309 cites, keeps naming the Spec by its number. ADR-0256 keeps a
  Lagging Surface blocking any proposed release, answered or not, ADR-0256:
  "A lagging surface or an unreadable map makes the new". ADR-0340 keeps the
  authorization record's `executables` rules, whose path shape the new
  `release_paths` key reuses, ADR-0340: "A Spec's `_authorization.md` may
  carry `executables:`, a list of exact repository-relative file paths."
  The release plan reads no `source_revision`, which after a squash merge
  survives only on the pull request head. ADR-0343 cites ADR-0340 for
  settle's executable check and governs QA acceptance; this Spec does not
  touch it. ADR-0349, written by this Spec, records the
  Release Ceiling. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AU answers of 2026-11-01 quoted in the authorization record, which keep the
  cycle AT rules and the cycle AJ grant for the Governed Paths each Spec
  declares; the governed files are the Spec Check's `constraints.go`, the two
  Baseline modules and their derived guides, golden, profile and Setup
  Manifest, and the Roundfix, `write-prd`, `write-techspec` and `write-tasks`
  skill files the Tasks change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0338-the-release-decision-travels-with-the-spec/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/release.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/speccheck/constraints.go`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A Spec's pull request and squash subject carry the Spec's Delivery Type, so
  a fix Spec can ship as a patch.
- A Spec without a declared type delivers exactly as today, `feat: deliver`.
- An approved authorization record that grants `release` can record a
  Release Ceiling and the release paths, and the release plan treats the
  records of the Specs it delivers as the answer to `approval_required` while
  the plan stays within every ceiling.
- Above a ceiling, breaking, classified by hand, or with any change outside a
  ceiling, the plan still asks, with the same bytes as today when no record
  in the range carries a ceiling.
- The authoring skills fill the type and ask the release question at
  planning time, and the Baseline wording that reaches adopters says so.

## Core Features

1. Delivery Type. The PRD frontmatter `type` names one of `feat`, `fix`,
   `perf` or `refactor`. `write-prd` writes `feat` (or `perf`), and
   `write-techspec` writes `fix` or `refactor` in the minimal PRD of its
   correction route. `roundfix spec check` reports any other value.
2. The type reaches the title. The archive copies the PRD's `type` into the
   Archive Record, and the Delivery Queue titles the pull request
   `<type>: deliver …` from it, `feat` when the record has none.
   `roundfix release plan` keeps classifying by the commit type.
3. Release Ceiling. `_authorization.md` may carry `release_ceiling`
   (`patch`, `minor` or `major`) and `release_paths` beside `release` in
   `operations`; a malformed or contradictory pair refuses the record.
4. The release plan reads the ceiling. For a plan that would ask, the plan
   finds each Spec the range delivers, reads its record at the parent of the
   delivering commit, and answers the approval when every delivered Spec's
   record grants a ceiling at or above the plan's impact, nothing is
   breaking or classified by hand, and every change that sets the impact is a
   Spec delivery. The plan names the records it read, in text and JSON.
   Backlog Entries become planning evidence beside `docs/specs/`,
   `docs/adr/`, `docs/findings/` and `docs/handoffs/`, so a plan commit that
   adopts or files them no longer needs manual classification.
5. The rules in writing. The Baseline's release and authorization clauses,
   the Roundfix skill's release, deliver, spec and archive references, the
   `write-prd`, `write-techspec` and `write-tasks` skills, the user guides,
   `CONTEXT.md` and ADR-0349 state them.

## Non-Goals / Out of Scope

- Reading the Delivery Type in the release plan; it reads commit types, as
  the maintainer's PR type answer says.
- Answering a manual classification, a breaking plan or a reset plan from a
  record; they still ask.
- Performing a release from a Run: `release` stays outside the Delivery
  Queue's operations, and the record only answers the plan's question.
- Backfilling `type` into existing Archive Records or `release_ceiling` into
  existing authorization records.
- Fetching a pull request head from a forge to read a record; the plan stays
  local.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, lint or formatter
  configuration, and any `### QA settlement` section.

## Success Metrics

1. Success Metric: a PRD with `type: fix` archives to a record carrying
   `type: fix`, and the queue's title for a fitting slug is
   `fix: deliver <slug>`; with no `type` the record has no `type` key and the
   title is `feat: deliver <slug>`, byte for byte as before.
2. Success Metric: a long slug keeps the numbered form within 91 UTF-16
   units for each of the four types, and the title of each type is classified
   by the release plan as that type.
3. Success Metric: `roundfix spec check --strict` reports `SC-PRD-TYPE` for a
   PRD whose `type` is `docs`, and no finding for each of the four types or an
   absent key.
4. Success Metric: an authorization record with `release_ceiling` and no
   `release` operation, an unknown ceiling, a ceiling without
   `release_paths`, or `release_paths` without a ceiling is refused with a
   reason naming the field; a well-formed pair grants and exposes both values.
5. Success Metric: in a temporary repository whose range delivers one Spec
   through a `feat:` squash and whose record grants `release` with ceiling
   `minor`, `roundfix release plan` exits 0 with decision `ready`, the
   proposed minor version and the record named; with ceiling `patch`, with a
   breaking commit, with `--impact minor --reason`, with a direct `feat:`
   commit outside the Spec, or with a second delivered Spec without a ceiling
   it exits 3 with `approval_required` and names why.
6. Success Metric: a range whose records carry no ceiling prints the same text
   and JSON as before, and a commit that changes only `docs/specs/` and
   `docs/backlog/` paths is classified maintenance-only instead of blocking.
7. Success Metric: the Baseline's release clause and authorization clause,
   their rendered guides and the Standard TypeScript goldens carry the new
   sentences, and a second Managed Refresh changes nothing.
8. Success Metric: the authoring skills, the Roundfix skill references and the
   user guides describe the Delivery Type, the Release Ceiling and the release
   question, and the mirrors equal their canonical copies.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The two Argus reports quoted in the adopted references, written before this
  Spec existed: the `feat: deliver` squash of a pure fix refused as a patch,
  and a release authorized with the Spec asked again.
- The Conventional Commits 1.0.0 specification
  (https://www.conventionalcommits.org/en/v1.0.0/, read through Exa on
  2026-11-01): "`fix` type commits should be translated to `PATCH` releases.
  `feat` type commits should be translated to `MINOR` releases." and "types
  other than `fix:` and `feat:` are allowed, for example
  @commitlint/config-conventional (based on the Angular convention)
  recommends `build:`, `chore:`, `ci:`, `docs:`, `style:`, `refactor:`,
  `perf:`, `test:`, and others."
- The Cocogitto guide "Conventional commits"
  (https://docs.cocogitto.io/guide/commit, read through Exa on 2026-11-01):
  "By default `cog commit` supports standard conventional commits type
  `feat`, `fix` plus the angular commit types: `build`, `ci`, `revert`,
  `docs`, `test`, `style`, `chore`, `perf`." This repository's `cog.toml`
  declares no `commit_types`, so the four Delivery Types are accepted types.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-11-01 against these sources.

- Reading at `dd3b2b62`: `pullRequestTitle` (`internal/cli/deliver_workflow.go`)
  builds `"feat: deliver " + specSlug`; `ClassifyCommit` rates `feat` minor,
  `fix` and `perf` patch, and any other type crossing the maintenance
  boundary as a commit that needs manual classification; the delivery squash
  `dd68aed0` deletes `docs/specs/0337-…/_authorization.md` and adds
  `docs/history/specs/0337-….md` in one first-parent commit, whose parent
  still holds the record; no code consumes the `release` operation.
- Characterization at `dd3b2b62`, in a scratch clone at `dd68aed0` with
  that commit's `bin/roundfix`: `roundfix release plan --from v0.127.0`
  reports `manual_classification_required`, with impact `minor` and one
  blocking commit, the plan commit `619cff66`, whose only paths outside
  `docs/specs/` are new Backlog Entries under `docs/backlog/`. Every plan
  commit of cycles AR, AS and AT (`7c0b0927`, `0c27cdc9`, `7758ea0c`,
  `619cff66`) changes `docs/backlog/` beside its planning files, so every
  queue range since then needed `--impact` and `--reason`, and a ceiling that
  refuses manual classification would never answer one.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix release authorization ceiling approval pull request type" --all --files --min-score 0.3`,
  which returned project mirrors (adopters' `_authorization.md` records, the
  Roundfix release reference) and the wiki page
  `wiki/concepts/cadencia-de-release-gatilhos-vs-lotes.md` on release cadence;
  none decides a release ceiling or the delivery title's type, and the page
  states its Roundfix policy is the project's own inference.
- Exa: the Conventional Commits specification and the Cocogitto guide quoted
  above.

## Glossary

- adds: **Delivery Type** — the Conventional Commits type a Spec delivers as
- adds: **Release Ceiling** — the highest impact a Spec's release grant answers
- not a term: **Delivery Queue** — existing glossary term, unchanged
- not a term: **Archive Record** — existing glossary term, unchanged

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AU brief for Spec 0338 of 2026-11-01 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-01), verbatim:

- Scope: "3 Specs (Recommended)" — option text: "0338 release: tipo do PR
  segue a Spec + teto de release no _authorization.md (os dois mexem no
  release plan e no write-*); 0339 apply aprovado apaga as aposentadas (nova
  ADR, context7-cli/exa-web-search ficam); 0340 links de evidência da QA
  existem no commit. Fila nessa ordem."
- PR type: "Campo no _prd.md (Recommended)" — option text: "Frontmatter
  `type` no _prd.md, preenchido pelo write-prd/write-techspec (rota de
  correção → fix ou refactor, feature → feat). A fila usa no título; sem o
  campo continua feat. O release plan lê o tipo do commit, como hoje."
- Release ceiling (asked earlier the same day, in the Backlog Entry): "Teto de
  impacto (Recommended)" — option text: "Campo com o teto (patch, minor ou
  major). O release plan aceita sem perguntar enquanto o plano novo ficar no
  teto ou abaixo; breaking ou acima do teto pergunta de novo. As skills de
  autoria fazem a pergunta no planejamento quando você pede release." — and
  "Próximo ciclo (Recommended)".
- Release of the cycle: "Uma no fim, v0.129.0 (Recommended)" — option text:
  "Depois das três Specs, com aviso aos adotantes; o apply que apaga
  aposentadas e o teto de release chegam juntos pela Baseline."
- Authorization: "Mesmas do AT (Recommended)" — option text: "Push e PR dos
  planos, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, mudança da Baseline
  (agent-instructions, release, spec-routing, skills de autoria) que chega aos
  adotantes, release v0.129.0 no fim e aviso aos adotantes."

Taken from those answers, ADR-0143, ADR-0309 and the reading above, with no
further interview:

- The Delivery Type set is `feat`, `fix`, `perf` and `refactor`, the Backlog
  Entry type set of `docs/agents/docs-layout.md`; all four are accepted by
  this repository's commit rules (Cocogitto's default types, no custom
  `commit_types`). A `refactor:` delivery that changes released code is a
  commit the release plan asks to classify by hand, as any such commit is
  today; the release reference says so.
- The Archive Record carries the type, because the queue titles the pull
  request after the archive commit removed the PRD and already reads the
  record for the title.
- The record of a delivered Spec is read at the first parent of the commit
  that deletes `<spec-root>/<slug>/_authorization.md` and adds the Archive
  Record `<slug>.md`: that parent is on the default branch and local, and the
  pair is found without reading Roundfix configuration, which the release
  plan never reads.
- Every Spec the range delivers must grant a ceiling at or above the impact,
  and every change that sets the impact must be such a delivery: one Spec's
  consent never answers for another change.
- `docs/backlog/` joins the release plan's planning-evidence paths: a Backlog
  Entry is dated intent, as a Finding already is, and without it every plan
  commit that adopts an entry forces the manual classification a ceiling may
  not answer (Validation). This is the one change of classification this Spec
  makes, and ADR-0349 records it.
- An answered plan reports decision `ready` and exit 0; a plan whose records
  do not answer keeps `approval_required` and names the reason.
- Four implementation Tasks in series (Delivery Type; Release Ceiling and the
  plan; Baseline clauses; authoring skills) and a QA gate.

## Open Questions

None.
