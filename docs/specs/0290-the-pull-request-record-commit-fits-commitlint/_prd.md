---
spec: 0290-the-pull-request-record-commit-fits-commitlint
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# The Pull Request Record Commit fits commitlint

Minimal PRD for a bug fix, minted by `write-techspec`. The Delivery Queue
writes the Pull Request Record Commit with the subject `docs: record pull
request #<number> in the <slug> Archive Record`, so a long Spec slug makes the
subject longer than the 100 characters a conventional commitlint configuration
allows. A repository whose `commit-msg` hook runs commitlint refuses the
commit, the queue parks the item with a generic `delivery-error`, and a retry
repeats the same subject. After this Spec the subject names the Spec by its
number, so it fits for every slug, and a commit the queue makes that a hook
still refuses parks with a reason naming the hook and the hand commit that
recovers it. The TechSpec holds the design.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-20-o-commit-do-registro-do-pull-request-estoura-o-commitlint-com-slug-longo.md`](references/2026-10-20-o-commit-do-registro-do-pull-request-estoura-o-commitlint-com-slug-longo.md),
   origin oraculum, type hint `fix`. On 2026-10-20 at 16:02 BRT the oraculum
   Delivery Queue parked a Spec whose slug has 54 characters with
   `delivery-error: record pull request: commit pull request record: git commit
   -m docs: record pull request #280 in the <slug> Archive Record: …
   commitlint … exit status 1`. The subject had 107 characters; the
   repository's `commit-msg` hook runs commitlint with
   `@commitlint/config-conventional`, whose `header-max-length` is 100. The
   operator committed the staged change by hand with the subject `docs: record
   pull request #280 in the 0088 Archive Record` and retried, and the queue
   adopted that commit.

No open Backlog Entry or unresolved Finding shares this context: the only open
Backlog Entry at the authoring base,
`docs/backlog/2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md`,
is about `roundfix runs causes`. `docs/findings/` holds no entry. This
planning change writes one new open Backlog Entry,
`docs/backlog/2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md`,
for the two other queue messages that grow with the slug, which the
maintainer's decision leaves out of this Spec.

## Prerequisites

None. This Spec is authored against `70c78d15` (v0.100.0). The other cycle AM
plans merge with it and deliver after it; one of them may also change
`internal/cli/deliver_workflow.go`, for the Linux Lane and the Run slot wait,
and this Spec keeps its change inside the three queue commits and the record
subject.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the record keeps `roundfix/archive-record/v1` and its
  `pull_request` field. The Spec number in the new subject is read from the
  existing slug. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the queue's Git commands run in the item
  worktree as before, and every test uses temporary repositories and homes
  with no remote beyond a local bare repository. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0302: "A queue commit names the
  Spec by number, and a hook refusal parks with the hand commit". Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20 through AskUserQuestion, authorization "the same
  rules as cycle AL", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0290-the-pull-request-record-commit-fits-commitlint/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- The Pull Request Record Commit's subject fits 100 characters for every
  slug and every legal pull request number, and names the Spec by its number.
- A record commit already made by hand, under the old long subject or any
  other, keeps being adopted, so items worked around before this Spec still
  pass.
- A commit the Delivery Queue makes that a commit hook refuses parks with a
  reason that names the hook, the refused subject, and the recovery, instead of
  a generic Git error; a repository without a commit hook keeps today's error.
- The `deliver` user guide and the Roundfix skill's `deliver` reference state
  the subject and the recovery.

## Core Features

1. The record subject. `docs: record pull request #<number> in the <spec
   number> Archive Record`, where the Spec number is the run of one to ten
   ASCII digits that starts the slug when a `-` follows it; a slug without
   one gives `docs: record pull request #<number> in the Archive Record`.
2. Adoption by proof. The queue keeps accepting an existing record commit by
   comparing its change with the candidate, never by reading its message.
3. A named hook refusal. The Pull Request Record Commit, the archive commit
   and the derived merge commit classify a failed `git commit` the way the
   Daemon classifies a Task commit refusal; the park keeps its blocker, class
   and stage, and its reason names the hook, the subject, the exit status, the
   hook's output and the recovery.
4. Documentation. The `deliver` user guide and the Roundfix `deliver`
   reference give the subject and the recovery, and the glossary's Hook
   Refusal entry covers a Delivery Queue commit.

## Non-Goals / Out of Scope

- Changing the archive commit subject `docs: archive <slug>`, the pull request
  title `feat: deliver <slug>`, or the derived merge message; the open Backlog
  Entry this planning change writes tracks the first two, which exceed 100
  characters for a slug longer than 86 characters.
- Passing `--no-verify`, retrying with another subject, or any automatic
  rewrite of a refused commit.
- A new blocker or Park Class, or a change to the Daemon's Task commit
  refusal.
- Rewriting any Archive Record or record commit already made.

## Success Metrics

1. Success Metric: in a temporary repository whose `commit-msg` hook refuses
   a header longer than 100 characters, the Pull Request Record Commit of a
   Spec whose slug has 54 characters, for pull request 280, is created with
   the 57-character subject `docs: record pull request #280 in the 0088
   Archive Record` and proven exact, where the Delivery Base refuses it.
2. Success Metric: a slug without a numeric prefix gives `docs: record pull
   request #7 in the Archive Record`, and the subject is at most 70
   characters for a ten-digit pull request number and a ten-digit Spec
   number.
3. Success Metric: a record commit made by hand under the old long subject,
   under a short subject and under an arbitrary subject is adopted with no
   second commit.
4. Success Metric: when a commit hook refuses the Pull Request Record
   Commit, the archive commit or the derived merge commit, the step's error,
   and so the `delivery-error` park reason, names the hook (from its own
   output, else the installed `pre-commit` or `commit-msg` hook), the refused
   subject, the exit status, the hook's output and the recovery, and the
   item's Park Class stays `environment`; without an executable hook the error
   is today's raw Git error.
5. Success Metric: after a hook refusal of the record commit, committing the
   staged change by hand with a conforming subject lets the next record step
   adopt it.
6. Success Metric: on the oraculum machine after release, the Delivery
   Queue delivers that repository's next Spec whose slug is longer than 48
   characters without a hand commit.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- The adopted Inbox Entry, written by an oraculum session on 2026-10-20 and
  copied verbatim into `references/`: the park reason, the 107-character
  subject, the `commit-msg` hook running commitlint, and the hand commit the
  queue adopted.
- This repository's history: the squash commit `9e92ac39` on `main` has the
  subject `feat: deliver 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers (#575)`,
  117 characters, written from the pull request title the queue built.
- Git's hook documentation (https://git-scm.com/docs/githooks): for `commit-msg`,
  "Exiting with a non-zero status causes the command to abort.", and for
  `pre-commit`, "Exiting with a non-zero status from this script causes the
  git commit command to abort before creating a commit."
- The commitlint conventional configuration
  (https://github.com/conventional-changelog/commitlint/blob/master/@commitlint/config-conventional/src/index.ts):
  `"header-max-length": [RuleConfigSeverity.Error, "always", 100]`.
- A Stack Overflow answer showing a refused commit's output
  (https://stackoverflow.com/questions/72875815): "header must not be longer
  than 100 characters, current length is 102 [header-max-length]" followed by
  "husky - commit-msg hook exited with code 1 (error)", the line that names
  the hook.

## Unreachable Acceptance

- criterion: Success Metric 6, on the oraculum machine after release, the
  Delivery Queue delivers that repository's next Spec whose slug is longer
  than 48 characters without a hand commit.
  reason: no Verification or QA row may read another repository, the real
  `~/.roundfix` or a remote, and only that repository's queue has its
  commitlint hook and its long slugs.
  satisfied-by: the oraculum operator's next delivery of such a Spec after
  this Spec's release.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Characterization at `70c78d15`, by reading `internal/cli/deliver_workflow.go`
  and running Git 2.56.0 in a temporary repository under a temporary `HOME`
  and `GIT_CONFIG_GLOBAL` whose `commit-msg` hook refuses a header longer than
  100 characters. The 107-character subject of the Inbox Entry was refused
  with exit status 1; Git printed only the hook's line, `header must not be
  longer than 100 characters, current length is 107 [header-max-length]`, and
  the change stayed staged. A commit of the same staged change with the
  57-character subject succeeded. The TechSpec records where the queue's
  three commits, the record proof and the Daemon's hook classification live.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "commitlint
  header-max-length commit-msg hook delivery queue commit subject" --all
  --files --min-score 0.3`, which returned the fluxus mirror's Archive Record
  of its Spec 0036 (a commit hook refusing a generated file), and `qmd query
  "roundfix hook refused commit Task commit hook stricter than verification"
  --all --files --min-score 0.3`, which returned the Roundfix mirror's Archive
  Record of Spec 0098, the origin of the Daemon's Hook Refusal, and the
  `autonomous-work` guides of two adopters. They confirmed that adopters run
  commit hooks the queue meets and that the Daemon already classifies a
  refusal, which this Spec reuses instead of a second classifier. No Inbox
  Entry was written, because the briefing forbids writing there.
- Exa: the commitlint conventional configuration and the Stack Overflow answer
  above, and Git's hook documentation
  (https://git-scm.com/docs/githooks). Challenging: Git prints nothing that
  names a hook, so a hook that does not name itself can only be found by the
  installed hooks, and a repository with both `pre-commit` and `commit-msg`
  cannot be told apart that way; the reason then names both.

## Glossary

- changes: **Hook Refusal**

## Decisions

- 2026-10-20 (about 18:15–18:35 BRT): the cycle AM scope "R → S → T
  (Recomendado)" (option text "Três Specs, um minor cada (v0.101–v0.103),
  nesta ordem. R primeiro porque já afeta o oraculum. Mesmas regras do AL."),
  the authorization "the same rules as cycle AL", and this Spec's design
  "Número + falha nomeada (Recomendado)" (option text "O assunto passa a ser
  'docs: record pull request #N in the 0088 Archive Record'. Se o hook
  recusar mesmo assim, o motivo do park cita o hook de commit e o contorno,
  não um delivery-error genérico.").
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AM briefing for Spec 0290 of 2026-10-20 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- A slug without a numeric prefix drops the identifier rather than keeping the
  slug, so the subject fits for every slug.
- The hook refusal keeps the `delivery-error` blocker, its class and its
  stage, and only its reason text changes, so every reader of the Park Class
  table keeps working.
- The queue reuses the Daemon's hook classification, and exports from the
  Daemon package the list of installed commit hooks it already reads, so a
  hook that does not name itself is still named.
- One new ADR, ADR-0302, records the subject rule, the named refusal and the
  rejected alternatives (`--no-verify`, a cut slug, a new blocker).

## Open Questions

None.
