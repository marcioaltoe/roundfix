---
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
status: active
created: 2026-10-17
surfaces: [backend, cli, docs]
---

# A Run never reaches the judge, and an amended plan keeps its settled Tasks

Five gaps of cycle AH, each observed by the operator on 2026-10-16 and
2026-10-17, cost a QA gate, a release gate, a Run or a hand repair. The QA
Agents of Specs 0273 and 0275 ran the archive plan, and its Archive Advice sent
Spec files to the judge with the judge's own key, which the Agent's shell held
although Roundfix had removed it. The v0.86.0 release gate failed once on a
worktree bootstrap test that passed on rerun. The command tests do not compile
for Windows, a release platform. An amendment of a Spec's planning files threw
away every settled Task's carry-forward three times in one day. And the
history cleanup refused an adopter's superseded ADR whose successor already
named it, because the two used a form of ADR reference the plan does not read.
After this Spec no Agent or Verification of a Run can make the judge send
anything, the bootstrap test has no race, every test compiles for Windows and a
contract keeps it so, a settled Task survives an amendment that leaves its own
contract alone, and the cleanup reads every form of ADR reference adopters use.

## Sources

There are no open Backlog Entries or Findings; the maintainer named these five
sources on 2026-10-17. Each is an operator observation, cited here with its
evidence.

1. A judge key reaches a Run. The QA gates of Spec 0273 (2026-10-16, its first
   QA Report, finding QA-001) and Spec 0275 (2026-10-17, its first QA Report,
   finding F2) ran the archive plan, whose Archive Advice made two
   authenticated archive-value requests each, using
   `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` from the Agent's environment. Both
   reports are in Git with their Spec folders at the `source_revision` of the
   Archive Records `docs/history/specs/0273-*.md` and
   `docs/history/specs/0275-*.md`. ADR-0273: "The judge's keys never enter a
   Run."
2. A flaky worktree test. `TestBootstrapSerializesAcrossSiblings` failed the
   v0.86.0 release gate in GitHub Actions run 39044083565 on 2026-10-17 with
   "stopped waiting for round 1 sibling task_02 bootstrap release; siblings
   started so far: [task_01 task_02]: watched work ended", and passed on
   rerun. ADR-0271: "A test bounds by real time only the budget it asserts."
3. The command tests do not compile for Windows. The delivery probe report test
   that Spec 0275 added (pull request #565) uses Unix-only system calls without
   a build constraint, and a check on 2026-10-17 of the tree at the v0.85.0
   release commit found the same defect in older test files of three packages.
4. An amended plan drops every settled Task's carry-forward. On 2026-10-17 the
   amendments of Spec 0272 (pull request #551) and Spec 0274 (pull request
   #559) refused Task Carry-Forward for every settled Task, and the operator
   carried them by hand three times.
5. A superseded ADR is refused although its successor names it. On 2026-10-17
   the lead ran the history cleanup of v0.87.0 on a scratch clone of the
   fluxus adopter: its superseded ADR 0002 records its successor by number and
   slug rather than in the `ADR-NNNN` form, and its accepted successor 0064 already
   names 0002 in its `supersedes` front matter, yet the plan refused ADR 0002
   as having no accepted successor.

## Prerequisites

None. Spec 0276 merged as pull request #568 and v0.87.0 is released from the
same tree; this Spec releases as v0.88.0.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; one environment variable name joins the Roundfix
  set, and the carry-forward JSON gains one list that older readers ignore.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the judge's credentials and its
  network request are the subject of the first change, which removes a request
  path and adds none; every test runs on temporary repositories with injected
  or local transports and sends nothing to a provider. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0283 (this Spec): "The judge sends
  nothing while the Run Marker is set." ADR-0284 (this Spec): "A settled Task
  binds only the planning text it names." ADR-0285 (this Spec): "A test that
  needs Unix lives in a file that says so." ADR-0273, which ADR-0283 amends:
  "The judge's keys never enter a Run." ADR-0271: "A test bounds by real time
  only the budget it asserts", so the bootstrap fix adds no sleep and no
  retry. ADR-0266: "A Proven Signal is the only signal Roundfix code or tests
  send", which the moved Unix-only tests keep. ADR-0250: "The coverage record
  lists tests for every release platform." ADR-0253: "Every Repository
  Contract Test runs on main and before a release". ADR-0252: "The selective
  gate runs the repository contracts a change makes relevant". ADR-0170: "A
  Task already completed on its target is nothing to carry", whose whole-set
  refusal stays. ADR-0158: "A budget-expired Run settles BudgetExceeded and
  stays recoverable", whose accepted outcomes stay. ADR-0282: "The cleanup
  writes the lineage a successor lacks." ADR-0278: "A retirement deletes."
  ADR-0254: "Refused Units never block". ADR-0189: "A version names one
  content." ADR-0280: "A guide Task records its own fingerprint." ADR-0256:
  "A release waits for the skills that describe a changed surface". ADR-0259:
  "A test runs in parallel unless it names why it cannot". ADR-0104: "A Spec
  accepts on evidence it did not author". ADR-0276: "A row the authorization
  keeps offline qualifies, and a byte invariant compares with the Delivery
  Base". ADR-0265: "Each claim carries its proof as a section." Set no obligation here, because this Spec changes none of their rules: ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0261 and ADR-0277 (history families, Archive Records, the History Full Tag and legacy folder conversion); ADR-0258 and ADR-0260 (reduced history and live citations); ADR-0164 (Run budget renewal); ADR-0233 (version regeneration at merge); ADR-0240 (the QA partial policy); ADR-0281 (rebuilt executables in tests); ADR-0251, ADR-0262 and ADR-0269 (lenient legacy reading, cancellation and archive refusals); ADR-0165 and ADR-0272 (review parks and records); ADR-0167 and ADR-0229 (QA rows and the operator archive); ADR-0169, ADR-0196 and ADR-0279 (the review's merge base, finding validation and dismissals); ADR-0237 (the Delivery Retry, whose carry-forward rule ADR-0284 refines without changing the retry). Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH approval of 2026-10-16 through AskUserQuestion, option "Governed Paths
  declarados" ("\"Concedo\" como em AE/AF/AG para os Governed Paths que cada
  Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."), and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `docs/references/coverage-record.json`,
  `internal/baseline/repository_test.go`, `skills/roundfix/SKILL.md`, and,
  by the maintainer's answer of 2026-10-18 ("Autorizar 1 linha no Makefile
  (Recomendado)"), `Makefile` for one `REPO_CONTRACT_TESTS` entry. No other
  Makefile, workflow, Go module, lint, formatter or Baseline module change.

## Goals

- No Agent Session and no Verification command of a Run can make the judge
  send a request, whatever the Agent's runtime re-reads from the user's shell
  startup, while the operator's own judge calls outside Runs keep working.
- The bootstrap serialization test passes or fails only on the behavior it
  asserts, never on the order in which a test helper closes a file.
- Every test of the repository compiles for Windows, and a repository contract
  fails when a test file stops compiling there.
- An amendment of a Spec's planning files that leaves a settled Task's own
  contract unchanged no longer refuses that Task's carry-forward.
- The history cleanup reads a superseded ADR's successor, and a successor's
  own record of what it supersedes, in every form adopters write them.

## User Stories

1. As the maintainer, I want the judge to send nothing from inside a Run even
   when an Agent's shell holds a judge key, so that Spec files never leave the
   machine from a QA gate and no QA gate fails on a request it did not mean to
   make.
2. As the maintainer, I want the bootstrap serialization test to have no race,
   so that a release gate fails only on a real defect.
3. As a contributor on any release platform, I want every test to compile for
   Windows and a contract to say so when one stops, so that the coverage
   record's platform lists are true.
4. As the operator of a Spec whose plan I amended after a Run settled some
   Tasks, I want those Tasks carried forward when the amendment did not touch
   the text they implement, so that proved work does not run again and I do
   not carry it by hand.
5. As an adopter's maintainer, I want the history cleanup to accept my ADRs'
   `superseded_by` and `supersedes` values in the forms my log uses, so that a
   correctly linked ADR leaves without a hand edit and gains no duplicate
   lineage line.

## Core Features

1. Run Marker. Every Agent Session a Run prepares or prompts, and every
   Verification command a Run executes, receives the **Run Marker**, the
   variable `ROUNDFIX_RUN_ID` set to the Run's ID.
2. Silent judge inside a Run. While the Run Marker is set, `roundfix spec
   judge`, the Archive Advice of `roundfix archive --plan` and the advice of
   `roundfix history sanitize` send no request, read no key and write no Judge
   Log entry, and their output names the Run as the reason.
3. Withheld keys. The Judge Key Variables also include `TYPESAFE_API_KEY`, and
   no Agent Session or Verification command of a Run receives any of them;
   the Daemon and the Queue Owner keep their own environment.
4. Race-free bootstrap test. Each sibling of the bootstrap serialization test
   waits on its own release, so a release written for one sibling can never
   reach another.
5. Windows-clean tests. Every test or helper that needs a Unix-only system
   call lives in a file constrained to Unix, under its own name, and the
   coverage record lists those tests for Darwin and Linux only.
6. Windows compile contract. A Repository Contract Test compiles every test of
   the repository for Windows and fails with the compiler's diagnostics.
7. Carry-forward over amendments. When a settled Task's PRD, TechSpec or Task
   Graph manifest moved, carry-forward compares only that Task's **Binding
   Text** and carries the Task when it is unchanged, listing the amended files;
   any other moved input, or a changed Binding Text, still refuses the whole
   set with today's reason.
8. ADR reference forms. The history cleanup resolves `superseded_by` written
   as `ADR-NNNN`, `NNNN`, `NNNN-slug`, `NNNN-slug.md` or a path to the ADR, and
   treats a successor whose `supersedes` front matter names the retired ADR in
   any of those forms as already naming it.
9. Documentation. The `spec`, `archive`, `history`, `baseline` and `reconcile`
   user guides, the Roundfix skill references that cover them, and the
   glossary describe each change.

## Non-Goals / Out of Scope

- Configuring any Agent runtime: turning off Codex's shell snapshot, changing
  `ZDOTDIR`, `HOME` or the runtime's shell for Agent sessions.
- Removing judge keys from Git hooks, the QA format command or the Worktree
  Bootstrap, or from the Daemon's and the Queue Owner's own environment.
  ADR-0273: "Git hooks and the QA format command keep the Run's environment".
- Changing the offline archive plan, the Jev Ceiling, the judge's questions or
  any judge answer.
- Changing the bootstrap permit or any production code of the Worktree
  Bootstrap.
- Editing the Makefile or a CI workflow to compile tests for Windows, or
  running any test on Windows.
- Item-level carry-forward, carrying part of a refused set, or tolerating a
  change to a Task's own file, `AGENTS.md`, `CONTEXT.md`, the `implement-task`
  skill or a `## Context` path.
- Editing any adopter repository. The maintainer chose an Inbox note for the
  adopters instead.
- Changing the lineage line's text, the refusal text, or the house form
  `superseded_by: ADR-NNNN` the Baseline template writes.

## Success Metrics

1. Success Metric: a test that runs an Agent Session and a Verification
   command for a Run with every Judge Key Variable and `TYPESAFE_API_KEY` set
   in the Daemon's environment observes none of them, and observes
   `ROUNDFIX_RUN_ID` set to the Run's ID, in each child; through the built
   binary, `roundfix archive <slug> --plan` with `ROUNDFIX_RUN_ID` set prints
   the Run as the reason for every candidate, and in the command and judge
   tests a counting transport sees no request even with a judge key set.
2. Success Metric: the bootstrap serialization test passes 50 consecutive
   runs, and the hazard that failed the v0.86.0 gate, a release writer that
   outlives its sibling, fails the old shared-release design deterministically
   and passes the new one.
3. Success Metric: compiling every test of the repository for Windows, with
   the contract tags, exits 0, the coverage record lists each moved test for
   Darwin and Linux only, and the new contract fails when a Unix-only call is
   put back into an unconstrained test file.
4. Success Metric: a Run whose settled Task names neither Project Constraints
   nor another Task's node is carried forward after an amendment that adds a
   Tooling authority row and a corrective Task node, and a Task whose named
   section changed is still refused with today's reason.
5. Success Metric: on a temporary repository shaped like the fluxus measurement
   (`superseded_by` written as number and slug, successor naming it in
   `supersedes`), the history plan removes the superseded ADR and appends no
   lineage line; each other accepted form resolves its successor.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not write, read on
2026-10-17:

- OpenAI Codex source, `codex-rs/core/src/shell_snapshot.rs`
  (https://github.com/openai/codex/blob/main/codex-rs/core/src/shell_snapshot.rs):
  the zsh snapshot script runs `[[ -r "$rc" ]] && . "$rc"` with `rc` set to
  `$ZDOTDIR/.zshrc` or `$HOME/.zshrc`, then captures `export -p`, and
  `EXCLUDED_EXPORT_VARS` holds only `"PWD", "OLDPWD"`.
- OpenAI Codex issue 30971, "Shell snapshots can trigger credential prompts and
  persist secret env vars in plaintext"
  (https://github.com/openai/codex/issues/30971): "`codex-rs/features/src/lib.rs`
  marks `shell_snapshot` stable and default-enabled" and "runtime execution
  then sources the snapshot before shell commands". A runtime the Agent uses
  re-exports what the user's shell startup exports, which is why the Run
  Marker, not only the removal, is needed.
- POSIX.1-2024, `read()`
  (https://pubs.opengroup.org/onlinepubs/9799919799/functions/read.html):
  "If no process has the pipe open for writing, read() shall return 0 to
  indicate end-of-file." A sibling that opens a shared release FIFO while
  another sibling's release writer is still open sees end-of-file when that
  writer closes, which is the failure of the v0.86.0 gate.
- Go command documentation, `go help buildconstraint`: a build satisfies
  "\"unix\", if GOOS is a Unix or Unix-like system", the constraint the moved
  test files carry.

## Validation

The decisions were checked on 2026-10-17 against two sources.

- Secondbrain. `wiki/index.md`, then two `qmd query` runs from the
  Secondbrain root. "agent runtime shell snapshot re-exports secret
  environment variable judge key run" returned the mirrors of ADR-0273 and of
  the Archive Record of Spec 0267 under `projects/roundfix/mirror/`; they agree
  with the repository and hold no decision on a runtime that re-reads the
  user's shell startup. "carry-forward settled task refused after spec
  amendment planning files" returned the mirrors of `CONTEXT.md`, the
  `reconcile` and `deliver` guides and the `autonomous-work` guides of three
  adopters, which all state today's whole-file rule and none a narrower one.
  The pending inbox `inbox/roundfix/` holds no entry; its triaged entries were
  resolved by earlier Specs.
- Exa. Supporting: the Codex source, the Codex issue, POSIX `read()` and the
  Go build constraint text above. Challenging: the same Codex issue's
  workaround, "Disable the feature in Codex config — `[features]` →
  `shell_snapshot = false`", which would keep the key out of one runtime's
  commands; ADR-0283 records why a runtime setting was rejected. The outside
  research is a standard reference, so no research Inbox Entry was captured.

## Glossary

- adds: **Run Marker**
- adds: **Binding Text**
- changes: **Judge Key Variables**
- changes: **Task Carry-Forward**

## Decisions

- 2026-10-16: which failure rules govern cycle AH — "Regras de AE–AG
  (Recomendado)".
- 2026-10-16: which authority covers the Governed Paths — "Governed Paths
  declarados".
- 2026-10-17: which items Spec F of cycle AH takes — "Chave do juiz entra no
  Run, Teste instável no worktree, Testes no Windows, Carry-forward e emendas".
- 2026-10-17: what to do with Spec F — "Escrever, entregar e lançar
  (Recomendado)".
- 2026-10-17: what the adopters receive — "Só nota no inbox (Recomendado)".

Autonomous route: no person present; decisions taken from the lead's briefing of 2026-10-17, its addendum adding the fluxus ADR reference forms, and the code and evidence characterized in the TechSpec.

Taken from those sources, with no further interview:

- The judge goes silent on the Run Marker, and the removal of the keys stays.
  See ADR-0283.
- Repository-configured commands keep the Run's environment, as ADR-0273
  decided.
- A settled Task binds the sections its References name, its own graph node
  and its table row. See ADR-0284.
- Unix-only tests move into `unix`-constrained files, and a contract of
  default relevance compiles every test for Windows. See ADR-0285.
- The ADR reference forms are read in both directions, and no lineage line is
  written for a successor that already names its predecessor in front matter.

## Open Questions

None.
