---
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# The queue shows a skipped Linux Lane and waits for a Run slot

Minimal PRD for two Delivery Queue fixes, minted by `write-techspec`. On
2026-10-20 the queue skipped an enabled Linux Lane for three Specs because its
image had been pruned, and the only trace was a console notice; the same day
it parked Spec 0289 as `delivery-error` only because other repositories' Runs
held every slot of the Active Implement Run ceiling. After this Spec, a lane
skipped for an environment reason is shown on the item in `roundfix deliver
status`, in the queue's Pull Request body and in `roundfix doctor`, and an item
whose only refusal is the ceiling waits for a free slot up to a configurable
bound before it parks as it does today. The TechSpec holds the design.

## Sources

Two Backlog Entries, written by this Spec's planning change directly into
`references/` and adopted at birth (see `references/_index.md`):

1. [`references/2026-10-20-the-queue-skips-the-linux-lane-without-a-trace.md`](references/2026-10-20-the-queue-skips-the-linux-lane-without-a-trace.md),
   type `fix`: with `delivery.linux_lane: true`, the gate of Specs 0285, 0286
   and 0287 skipped the lane on a missing image and merged them with one
   console notice.
2. [`references/2026-10-20-the-queue-parks-an-item-at-the-run-ceiling.md`](references/2026-10-20-the-queue-parks-an-item-at-the-run-ceiling.md),
   type `fix`: the queue parked Spec 0289 as `delivery-error` when the child
   `roundfix implement` met the machine-wide ceiling, and it waited for a
   manual `roundfix deliver retry`.

Both change how the Delivery Queue reports and advances an item between its
Run and its Pull Request, in the same workflow and engine, so one Spec owns
them, as the cycle AM decision assigned. No other open Backlog Entry or
unresolved Finding shares this context: on 2026-10-20 `docs/backlog/` held
only the `roundfix runs causes` entry another cycle AM Spec owns, and
`docs/findings/` held none.

## Prerequisites

None that the Task Graph enforces. The other cycle AM Specs are authored
beside this one and deliver first: the Pull Request Record Commit Spec changes
`internal/cli/deliver_workflow.go` (the record commit subject and the park
reason of a refused Git hook), the user guide `deliver.md` and the Roundfix
`deliver` reference, and raises the Roundfix skill version. This Spec changes
other functions of that file (`Gate`, `RunSpec`, `deliveryRunResult`,
`newCommandDeliveryEngine`) and other sections of those documents. Before this
plan enters a Run, the operator rebases it on the default branch and runs
`roundfix spec check 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot --strict --run-verification`
again.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema, table or record
  field is minted. The lane skip and the wait reuse the existing item warning
  and blocker texts; the new configuration key `runs.slot_wait` is a Go
  duration. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. The Pull Request body is written through
  the existing GitHub boundary, `doctor` runs only `docker info` and `docker
  image inspect`, and every test uses fake Docker, fake runners and temporary
  homes and repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the lane visibility, the doctor check, the wait and their documentation follow recorded decisions. ADR-0303: "A Linux Lane skipped for its environment is shown where the delivery is read". ADR-0304: "A queue item waits for a Run slot up to runs.slot_wait". ADR-0288: "The Delivery Queue gate runs the changed Go packages in a Linux Lane". ADR-0294: "The Linux Lane runs like CI". ADR-0297: "The Linux Lane reads the repository history read-only". ADR-0199: "A queue's token ceiling stops new starts and never a running Task". ADR-0279: "Status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor". ADR-0184: "A TechSpec states a command surface as a transcript". ADR-0259: "A test runs in parallel unless it names why it cannot". ADR-0271: "A test bounds by real time only what it asserts". ADR-0285: "Every test file compiles for Windows, and a repository contract proves it". ADR-0189: "An owned skill's version names its content, and its minimum is the bundle". ADR-0256: "A release waits for the skills that describe a changed surface". ADR-0280: "A guide Task records its own surface fingerprint, a PRD records who decided, and the glossary is held to its own rule". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0240: "One QA partial policy, and rows a Run sandbox cannot reach". ADR-0276: "A row the authorization keeps offline qualifies, and a byte invariant compares with the Delivery Base". ADR-0265: "A Spec proves what it claims, and every active Spec is held to it". ADR-0183: "An attributed claim carries a receipt the check proves". Set no obligation here, because this Spec changes none of their rules: ADR-0193 (a Spec names its prerequisites, and the owner waits for them); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0295 (reachability reads the fetched default branch, and origin names it when origin/HEAD is unset); ADR-0298 (a timed-out profile proof gets a third attempt, and its setup deadline is configurable); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0306 (archive and cancel refuse a path whose resolved target leaves the repository); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof); ADR-0299 (the history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a pre-PR review finding parks only after validation). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20 through AskUserQuestion, authorization "the same
  rules as cycle AL", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
  go.mod, lint, que pedem autorização específica)."), and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- An enabled Linux Lane that the queue gate skipped for an environment reason
  is visible on the item in `roundfix deliver status`, in the Pull Request body
  the queue writes, and in `roundfix doctor`; the skip still passes the gate.
- A skip because the lane has nothing to test stays as quiet as today.
- An item whose only preflight refusal is the Active Implement Run ceiling
  waits for a free slot, shows the wait with its holders and deadline in
  `roundfix deliver status`, starts its Run when a slot frees, and parks
  exactly as today when the wait ends.
- The store's ceiling check stays the only authority: two waiting queues never
  hold Runs beyond the ceiling.

## Core Features

1. Lane Environment Skip on the item. The gate returns the reason of an
   environment skip with its passing result; the Delivery Engine records it as
   the item warning `linux-lane-skipped: <reason>`, which status prints, and
   adds a paragraph naming it to the Pull Request body.
2. Doctor's `linux-lane` check. With `delivery.linux_lane: true` in the
   checkout's effective configuration, `roundfix doctor` reports whether Docker
   is on PATH and reachable and whether the lane's image is present locally,
   as warnings with a next action; the line is absent when the lane is off.
3. Run Slot Wait in the engine. On a ceiling-only refusal the engine records
   the wait on the running item, polls the Run slots, starts the Run when one
   frees, and parks at the earlier of its own deadline and the queue's
   deadline; `runs.slot_wait` in User Config bounds it.
4. Refusal recognition and wiring. The queue's Implement step recognizes the
   ceiling-only refusal and reads the Run slots the way `implement` counts
   them, and the Queue Owner's engine receives the configured wait.
5. Documentation. The user guide pages for `deliver`, `doctor` and
   configuration, the Roundfix `deliver`, `setup` and `runtime` references and the
   glossary describe both behaviors.

## Non-Goals / Out of Scope

- Parking an item because the lane was skipped, or pulling an image.
- Changing the console and lane log notices, the skip reasons, their order, or
  `roundfix deliver lane`'s exit codes.
- Waiting on any other preflight refusal, raising or changing
  `runs.max_active`, or changing `roundfix implement`'s output or exit codes.
- A new Delivery Stage, Run Database column or schema version.
- Showing the wait or the lane skip in `roundfix deliver plan`.

## Success Metrics

1. Success Metric: through the gate with the compiled Docker stand-in and an
   absent image, the gate passes and returns the reason `image golang:1.26.1
   is not present locally`; the engine records `linux-lane-skipped: image
   golang:1.26.1 is not present locally` on the item; `roundfix deliver
   status` prints it as Surface Transcript 1; and the Pull Request body ends
   with API Contract 3's paragraph naming that reason.
2. Success Metric: `no go.mod at the repository root` and `no changed Go
   package` record no warning and leave the body `Delivers Spec <slug>.`, a
   disabled lane records nothing, and a later gate that runs the lane removes
   an earlier `linux-lane-skipped` warning.
3. Success Metric: with `delivery.linux_lane: true`, `roundfix doctor` prints
   `linux-lane: warn (DR-LANE-IMAGE-MISSING: …; next: docker pull <image>)`
   for an absent image, `DR-LANE-DOCKER-MISSING` and
   `DR-LANE-DOCKER-UNREACHABLE` for Docker, and `linux-lane: ok (…)` when both
   are fine; with the lane off the line is absent and every existing doctor
   transcript is unchanged.
4. Success Metric: on a ceiling-only refusal the item stays `running` with the
   blocker text of Surface Transcript 2, starts its Run after a slot frees,
   parks `delivery-error` with today's exact reason at `runs.slot_wait`, parks
   `queue-deadline` when the queue's deadline comes first, and parks at once
   with `runs.slot_wait: 0` or any other preflight refusal.
5. Success Metric: two queues waiting on one free slot over one Run Database
   create exactly one Run, and the other keeps waiting with its blocker text.
6. Success Metric: on the maintainer's machine, with `delivery.linux_lane:
   true` and the lane image removed, `roundfix doctor` reports
   `DR-LANE-IMAGE-MISSING` against the real Docker daemon.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- GitHub's documentation of status checks
  (https://docs.github.com/en/pull-requests/reference/status-checks): "A job
  that is skipped will report its status as "Success". It will not prevent a
  pull request from merging, even if it is a required check." A skipped check
  passing silently is the published behavior of the forge this queue merges
  through; the lane's skip keeps passing, so its visibility has to come from
  the item, the Pull Request body and doctor.
- GitHub's documentation of Actions concurrency
  (https://docs.github.com/en/actions/how-tos/write-workflows/choose-when-workflows-run/control-workflow-concurrency):
  "When a concurrent job or workflow is queued, if another job or workflow
  using the same concurrency group in the repository is in progress, the
  queued job or workflow will be `pending`." A run that meets a capacity limit
  waits in a visible pending state rather than failing, which is the shape of
  the Run Slot Wait.
- This repository's ceiling race test, which no Task of this Spec authors:
  `TestRunCeilingIsCountedInsideTheCreateTransaction` in
  `internal/store/store_test.go` proves that two concurrent `CreateRun` calls
  at a ceiling of one create exactly one Run. It must pass unchanged at the
  audited head, because the wait relies on it.
- The maintainer's observations of 2026-10-20 recorded in Sources 1 and 2.

## Unreachable Acceptance

- criterion: Success Metric 6, on the maintainer's machine with
  `delivery.linux_lane: true` and the lane image removed, `roundfix doctor`
  reports `DR-LANE-IMAGE-MISSING` against the real Docker daemon.
  reason: no Verification or QA row may run real Docker, and removing an
  image changes the maintainer's machine; every row uses the compiled Docker
  stand-in.
  satisfied-by: the maintainer runs `roundfix doctor` in this repository after
  the release that carries this Spec, with the image absent or after `docker
  image rm golang:1.26` followed by `docker pull golang:1.26`.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Secondbrain, read only. Read `wiki/index.md`, then ran `qmd query "roundfix
  linux lane skipped image delivery queue run ceiling max_active wait" --all
  --files --min-score 0.3`; it returned no hit above the threshold, so no wiki
  page holds a decision on lane skips or ceiling waits. No Inbox Entry was
  written, because the briefing forbids writing to the Secondbrain.
- Exa: the two GitHub pages quoted above. Challenging: GitHub's pending state
  is unbounded by default and can cancel a pending run; the maintainer chose a
  bound ("Espera com limite"), and a waiting item is never replaced.
- Characterization on 2026-10-20 at `70c78d15` (v0.100.0), recorded in the
  TechSpec, including an overlay probe that wrote nothing to the repository.

## Glossary

- adds: **Lane Environment Skip**
- adds: **Run Slot Wait**

## Decisions

- 2026-10-20 (about 18:15–18:35 BRT): the cycle AM scope "R → S → T
  (Recomendado)" (option text: "Três Specs, um minor cada (v0.101–v0.103),
  nesta ordem. R primeiro porque já afeta o oraculum. Mesmas regras do AL.");
  the minor items "Distribuir + U (Recomendado)" (option text: "'runs causes'
  vai em T (mesmo resolvedor e storage). 'Esperar vaga' vai em S (gate e
  preflight da fila). O endurecimento de symlink vira a Spec U, no fim
  (v0.104.0)."); for the wait, "Espera com limite" (option text: "Espera até
  um prazo configurável (padrão 2h) e depois estaciona como hoje."); the
  authorization: the same rules as cycle AL.
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the operator's cycle AM briefing for Spec 0291 of 2026-10-20 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- A skip stays a skip. Only an enabled lane's environment skip is shown;
  `no go.mod at the repository root` and `no changed Go package` stay quiet
  because the lane had nothing to test. A fetched default whose configuration
  cannot be read is shown, because nobody can prove the lane was off.
- The skip rides on the item's existing warning and the wait on its existing
  blocker text, so no schema changes and older item binaries keep reading the
  Run Database. Recorded in ADR-0303 and ADR-0304.
- `doctor` reads the checkout's effective `delivery.linux_lane`, not the
  fetched default, and prints its line only when the lane is on.
- The wait's key is `runs.slot_wait` in User Config only, beside
  `runs.max_active`, because the ceiling it waits on is machine-wide and only
  User Config can set it; a Project Config value is ignored with the existing
  warning.

## Open Questions

None.
