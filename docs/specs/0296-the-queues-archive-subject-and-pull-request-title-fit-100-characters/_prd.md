---
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
status: active
created: 2026-10-21
surfaces: [backend, cli, docs]
---

# The queue's archive subject and pull request title fit 100 characters

Minimal PRD for a bug fix, minted by `write-techspec`. The Delivery Queue
writes the archive commit `docs: archive <slug>` and opens the pull request
titled `feat: deliver <slug>`, which a squash merge turns into the subject
`feat: deliver <slug> (#<number>)`. Both grow with the slug, so a long slug
passes the 100 characters a conventional commitlint configuration allows: a
`commit-msg` hook refuses the archive commit, and the squash commit lands on
the default branch with an over-long subject. After this Spec both messages
keep the slug while it fits and otherwise name the Spec by its Spec Number,
the title followed by the Spec title cut at a word, so neither passes 100
characters. The TechSpec holds the design.

## Sources

One open Backlog Entry, adopted by this Spec's planning change into
`references/` (see `references/_index.md`):

1. [`references/2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md`](references/2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md),
   type `fix`, written by Spec 0290's planning change for the two queue
   messages ADR-0302 left out. Each message is 14 characters plus the slug,
   so a slug longer than 86 characters passes 100; this repository's squash
   commit `9e92ac39` has a 117-character subject. Source Revision Recovery,
   which reads the trailing `(#N)` of the squash subject, must keep working.

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base this entry was the only file in `docs/backlog/`, and
`docs/findings/` holds no entry.

## Prerequisites

None in `requires`. This Spec is authored against `8fb2da1f` (v0.102.0 plus
the CI timeout change), where Spec 0291 has merged, and delivers after the
active Specs 0292 and 0293 and the cycle AN Specs 0294 and 0295. Spec 0291
appends its Linux Lane paragraph to the Pull Request body in the Delivery
Engine (`internal/delivery`); this Spec changes only the `Archive` subject and
the `Publication` title in `internal/cli/deliver_workflow.go` and keeps the
body `Delivers Spec <slug>.` as it is. Its documentation Task raises skill
versions over whatever the tree holds when it starts, so it rebases over the
other Specs' version bumps.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the Spec Number is read from the existing slug by ADR-0302's
  rule, and the Spec title from the existing Archive Record. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the title goes to the same `gh pr create`
  call as before, and every test uses temporary repositories, homes and fake
  forges. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0302: "A queue commit names the
  Spec by number, and a hook refusal parks with the hand commit"; ADR-0309:
  "A queue archive subject and pull request title keep the slug while it fits
  100 characters". Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AN approval of 2026-10-21 through AskUserQuestion, authorization "the same
  rules as cycle AM", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`. No Makefile,
  workflow, Go module, lint, formatter, Baseline module or Project Config file
  change.

## Goals

- The archive commit subject is at most 100 characters for every slug, and
  unchanged byte for byte whenever `docs: archive <slug>` fits.
- The pull request title plus the squash suffix ` (#<number>)` is at most 100
  characters for every slug and every pull request number below 100,000, and
  the title is unchanged byte for byte whenever `feat: deliver <slug>` leaves
  room for that suffix.
- Source Revision Recovery, the release plan, the queue's adoption of an
  existing archive commit and its lookup of an existing pull request keep
  working, and an archive commit made under the old long subject is still
  adopted.
- The `deliver` user guide, the Roundfix `deliver` reference and the
  `archive-spec` skill state both messages.

## Core Features

1. The archive subject. `docs: archive <slug>` when it is at most 100 UTF-16
   code units; otherwise `docs: archive the <spec number> Spec`, or `docs:
   archive the Spec` for a slug without a Spec Number.
2. The pull request title. `feat: deliver <slug>` when it is at most 91
   units, leaving 9 for ` (#99999)`; otherwise `feat: deliver <spec number>:
   <Spec title>`, or `feat: deliver <Spec title>` without a Spec Number, the
   Spec title read from the Archive Record and cut at a word so the title is
   at most 91 units. The Pull Request body still names the full slug.
3. One Spec Number rule. The record subject, the archive subject and the
   title read the Spec Number through one helper, ADR-0302's rule, and the
   record subject keeps its bytes.
4. Readers by proof. Source Revision Recovery and the legacy delivery reader
   still read the number from a squash subject built on the new title, the
   release plan still classifies it as `feat`, and an archive commit made by
   hand under any subject, the old long one included, is still adopted.
5. Documentation. The `deliver` guide, the Roundfix `deliver` reference and
   the `archive-spec` skill give both messages, and the glossary defines
   **Spec Number**.

## Non-Goals / Out of Scope

- Changing the Pull Request Record Commit subject, the derived merge message,
  the Pull Request body or the Linux Lane paragraph the Delivery Engine adds
  to it.
- Renaming a pull request an older binary already opened, or after it opens.
- Rewriting any commit, squash subject or Archive Record already made.
- A slug length limit or a change to slug validation.

## Success Metrics

1. Success Metric: in a temporary repository whose `commit-msg` hook refuses
   a header longer than 100 characters, the queue's archive commit of a Spec
   whose slug has 100 characters and starts with `0296-` succeeds with the
   subject `docs: archive the 0296 Spec` and is proven exact, where the
   Delivery Base's slug subject is refused.
2. Success Metric: the archive subject of an 86-unit slug is `docs: archive
   <slug>` (100 units) and that of an 87-unit slug is `docs: archive the
   <spec number> Spec`; a slug without a Spec Number gives `docs: archive the
   Spec`.
3. Success Metric: the title of a 77-unit slug is `feat: deliver <slug>`
   (91 units) and that of a 78-unit slug is the number form; for Spec 0278's
   slug and title the result is `feat: deliver 0278: A fixed round-2 finding
   gets one Correction Review, and the plan names`, 90 units; every title plus
   ` (#99999)` is at most 100 units, counted in UTF-16 code units.
4. Success Metric: a squash subject made of each title form plus ` (#575)`
   gives `575` to `spec.PullRequestFromSubject` and is classified `feat` with
   a minor impact and no breaking change by `releaseplan.ClassifyCommit`; the
   archive subject forms are classified `docs`.
5. Success Metric: an archive commit made by hand under the old long subject
   is adopted by the archive step with no second commit.
6. Success Metric: `Publication` for a long slug returns the number-form
   title and the body `Delivers Spec <slug>.`, and for a short slug returns
   today's title without reading the Archive Record.
7. Success Metric: after release, the first squash commit the Delivery Queue
   merges for a Spec whose slug is longer than 77 characters has a subject of
   at most 100 characters that names the Spec Number.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-21.

- This repository's history: the squash commits `9e92ac39`, `feat: deliver
  0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers (#575)`,
  117 characters, and `433a7c5a`, `feat: deliver
  0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage (#578)`,
  102 characters, both written from the title the queue built.
- The commitlint conventional configuration
  (https://github.com/conventional-changelog/commitlint/blob/master/@commitlint/config-conventional/src/index.ts):
  `"header-max-length": [RuleConfigSeverity.Error, "always", 100]` and
  `"subject-full-stop": [RuleConfigSeverity.Error, "never", "."]`; and the
  rule's check
  (https://github.com/conventional-changelog/commitlint/blob/master/@commitlint/ensure/src/max-length.ts):
  `typeof value === "string" && value.length <= max`, a JavaScript length in
  UTF-16 code units.
- GitHub's documentation of squash merges
  (https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges):
  for more than one commit the default summary is "The pull request title,
  followed by the pull request number".

## Unreachable Acceptance

- criterion: Success Metric 7, after release, the first squash commit the
  Delivery Queue merges for a Spec whose slug is longer than 77 characters
  has a subject of at most 100 characters that names the Spec Number.
  reason: only a real forge squash-merges a pull request, and no
  Verification or QA row may reach the network or run `roundfix deliver`
  against a forge.
  satisfied-by: the operator's `git log -1 --format=%s` of that squash commit
  on the default branch after the first such delivery following this Spec's
  release.

## Validation

The decisions were checked on 2026-10-21 against these sources.

- Characterization at `956066a7`, rechecked at `8fb2da1f`, by reading
  `internal/cli/deliver_workflow.go` and running Git 2.56.0 in a temporary
  repository under a temporary `HOME` whose `commit-msg` hook refuses a
  header longer than 100 characters. `git commit -m "docs: archive <Spec
  0278's slug>"`, 110 characters, exited 1 with the hook's line `header must
  not be longer than 100 characters, current length is 110
  [header-max-length]` and left the change staged; the same change committed
  with `docs: archive the 0278 Spec`. `git log` at `8fb2da1f` shows 104 `feat:
  deliver` squash subjects, two of them over 100 characters, and 6 archived
  slugs longer than 77 characters. The TechSpec records the readers of both
  messages.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "pull request
  title squash subject commitlint header-max-length slug delivery queue"
  --all --files --min-score 0.3`, which returned the Roundfix mirror's
  `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`
  (squash subjects carry `(#N)`, which Source Revision Recovery reads), the
  mirror of this Backlog Entry and the skills mirror's `conventional-commits`
  skill (`projects/skills/mirror/skills/07-evidence-delivery/conventional-commits/SKILL.md`:
  "Pull request titles must also use Conventional Commits because squash
  merge often uses the PR title as the final commit message"); and `qmd query
  "conventional commit subject truncate long title word boundary" --all
  --files --min-score 0.3`, which returned the same skill and nothing on
  cutting a title. They confirmed that the title is held to the commit
  contract and that the trailing number must survive; no prior decision on
  cutting titles exists. No Inbox Entry was written, because the briefing
  forbids writing there.
- Exa: the commitlint configuration and `max-length` check above, and
  GitHub's squash-merge documentation and its changelog of 2022-08-23
  (https://github.blog/changelog/2022-08-23-new-options-for-controlling-the-default-commit-message-when-merging-a-pull-request/:
  "Pull request title: pull request title and number on the first line").
  Challenging: commitlint counts UTF-16 code units, not bytes or runes, so the
  bound counts the same way; `subject-full-stop` refuses a trailing `.`, so
  the cut title drops one; and with a single-commit pull request GitHub's
  default summary is the commit's title, not the pull request's, which this
  Spec cannot change and the queue's multi-commit candidates do not meet.

## Glossary

- adds: **Spec Number**

## Decisions

- 2026-10-21 (about 01:00 BRT): the cycle AN scope "V → W → X (Recomendado)"
  (option text "Três Specs: v0.105 a v0.107. Y fica no backlog como entrada
  aberta; são só 3 linhas de ruído no reconcile e não custam nada."), the
  authorization "the same rules as cycle AM", and this Spec's design "Slug se
  couber (Recomendado)" (option text "Continua `docs: archive <slug>` e
  `feat: deliver <slug>` quando cabem em 100 (quase sempre). Acima disso, usa
  o número: `docs: archive the 0290 Spec` e `feat: deliver 0290` + título da
  Spec cortado na palavra. Nada muda nos adotantes com slugs curtos.").
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AN briefing for Spec 0296 of 2026-10-21 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The title reserves ` (#99999)`, 9 units, because the pull request number is
  unknown when the title is written; a repository past 99,999 pull requests
  gets one more unit per digit.
- "Fits" for the title counts that reserve, because the squash subject is the
  header commitlint reads on the default branch; a slug of 78 to 86 units now
  gets the number form.
- Length counts UTF-16 code units, as commitlint does.
- The Spec title is the Archive Record's `title`, read only when the number
  form is needed, so a short slug's publication reads nothing new.
- The number form separates the number from the title with `: `; without a
  Spec Number the title follows `feat: deliver ` directly.
- One new ADR, ADR-0309, records the rule; the glossary gains **Spec Number**,
  the rule three messages now share.

## Open Questions

None.
