---
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# A new Delivery Queue resumes an item whose Spec is already archived on its branch

Minimal PRD for a bug fix, minted by `write-techspec`. When a Delivery Queue
item stops after its archive commit, for example by using up its retries
while its Pull Request waits for checks, a new queue continues the same item
branch and parks the item `delivery-error` before doing anything, because
Delivery Revalidation reads a Spec folder the archive commit already removed.
On 2026-10-21 that left Spec 0292 for the operator to merge by hand. After
this Spec, a new queue whose item branch already carries an exact archive of
the Spec resumes the item at review and carries it through the gate,
publication, checks and merge without running implementation or the archive
again, and a branch that cannot be proven parks with a reason that names the
cause and the manual path. The TechSpec holds the design.

## Sources

One Backlog Entry, written by this Spec's planning change directly into
`references/` and adopted at birth (see `references/_index.md`):

1. [`references/2026-10-22-a-new-delivery-queue-parks-an-item-whose-spec-is-archived-on-its-branch.md`](references/2026-10-22-a-new-delivery-queue-parks-an-item-whose-spec-is-archived-on-its-branch.md),
   type `fix`: the new queue for Spec 0292 printed `Continuing item branch`
   and parked `delivery-error` at revalidation with a stat ENOENT on
   `docs/specs/<slug>`; the operator merged Pull Request #619 by hand and ran
   `roundfix deliver stop`.

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` held only the reconcile process-record entry
that Spec 0300 owns, and `docs/findings/` held none.

## Prerequisites

None that the Task Graph enforces. This Spec is authored against `34b4e246`
(v0.107.0) beside the cycle AO Specs 0297, 0298 and 0300 and delivers third,
after 0297 and 0298. Spec 0298 changes the check wait in
`internal/delivery/engine.go` (`checkCandidate`, `createPullRequest`,
`NewEngine`); this Spec changes other functions of that file
(`advanceItem`, `Retry`) and adds one engine dependency. Spec 0297 may touch
the archived-candidate reading of the Task Graph; this Spec reads no Task
Graph. Each documentation Task raises the Roundfix skill version over whatever
the tree holds when it starts. Before this plan enters a Run, the operator
rebases it on the default branch and runs `roundfix spec check
0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
--strict --run-verification` again.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema, table or record
  field is minted; the new blocker `archived-item-unproven` is a text in the
  existing item blocker field, like every other blocker. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; a resumed item reaches the forge through
  the existing push, pull request and merge calls, and every test uses
  temporary repositories, homes and fake forges. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0223: "A delivery keeps its work
  across archive, requeue and review"; ADR-0229: "An operator archive resumes
  any park from the Run start"; ADR-0237: "A Delivery Retry records a merge
  made outside the queue"; ADR-0240: "One QA partial policy, and rows a Run
  sandbox cannot reach"; ADR-0312: "A new Delivery Queue resumes an archived
  item branch at review"; and, because they cite ADR-0240, the QA gate's
  rules ADR-0104: "A Spec accepts on evidence it did not author", ADR-0167:
  "The pre-PR Pull Request row never decides a qualifying partial" and
  ADR-0276: "A row the authorization keeps offline qualifies, and a byte
  invariant compares with the Delivery Base", which this Spec's QA gate
  follows unchanged; and ADR-0295: "Reachability reads the fetched default
  branch, and origin names it when origin/HEAD is unset", which cites
  ADR-0276 and which the proof leaves unchanged, since it reads only the item
  branch. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 through AskUserQuestion, authorization "the same
  rules as cycle AN", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module or Project Config file change.

## Goals

- A new queue whose continued item branch carries an exact archive of the
  Spec delivers it: no `delivery-error` park, no Implement Run, no second
  archive commit and no second Pull Request.
- The resumed head is reviewed under the configured pre-PR review policy
  before publication, and passes the repository gate and the required checks
  like any candidate.
- An item branch whose archive cannot be proven parks before any Run with a
  reason that names the cause and the manual path, and a Delivery Retry of it
  proves the branch again.
- An item branch without an archive keeps today's revalidation and Run.

## Core Features

1. Archived Item Resume. A queued item whose worktree head holds the Spec's
   archive and not its folder, and whose commit that removed the folder is an
   exact archive of the Spec, records its head as the candidate and enters
   `reviewing`; the archive step then finds the Spec archived and commits
   nothing.
2. Unproven archive park. A head that holds neither the folder nor an
   archive, no removing commit, or a removing commit that is not exact parks
   the item `archived-item-unproven: <cause>` in the `finding` Park Class,
   whose next action names the manual path; a Delivery Retry returns it to
   `queued`.
3. A missing Specs Root. When the archived Spec was the only active one, the
   item worktree has no Specs Root directory, and every archive path the queue
   resolves for the item still resolves.
4. Documentation. The `deliver` guide and the Roundfix `deliver` reference
   describe the resume, the new blocker and its retry; the glossary adds
   **Archived Item Resume** and revises **Delivery Revalidation**.

## Non-Goals / Out of Scope

- Skipping the pre-PR review for a resumed head, or reusing a review record
  from an earlier queue.
- Changing ADR-0223's continuation rule, the archive proof itself, the
  Delivery Retry of an archived item that has a recorded candidate (ADR-0223,
  ADR-0229) or a retry's merge observation.
- Detecting a Spec already merged on the default branch at `deliver start`;
  `deliver start` still loads the Spec from the checkout.
- The check wait (Spec 0298) and the Wave collision reading of a corrective
  Task (Spec 0297).

## Success Metrics

1. Success Metric: through `roundfix deliver start` in a temporary repository
   whose continued item branch carries an exact archive commit and a Pull
   Request Record Commit for an open Pull Request of a fake forge, the item
   ends `merged` with that Pull Request, no Implement Run, no new commit on
   the item branch, no Pull Request created and one review of the item head,
   where the Delivery Base parks it `delivery-error` at revalidation.
2. Success Metric: the same holds when the archived Spec was the only active
   one, so the item worktree has no Specs Root directory.
3. Success Metric: a removing commit that also changes an unrelated file, a
   folder removed without an archive, and an archive whose removing commit
   is absent from the item branch each park `archived-item-unproven: <cause>`
   with no Run and no review, and `roundfix deliver status` prints the
   `finding` Park Class with the manual path.
4. Success Metric: a Delivery Retry of an `archived-item-unproven` park
   returns the item to `queued` and counts one retry.
5. Success Metric: an item branch whose head still holds the Spec folder is
   revalidated and runs as on the Delivery Base.
6. Success Metric: after release, the operator's next new queue over an item
   branch that already archived its Spec reaches `merged` without a hand
   merge.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- This repository's history: the squash commit of Spec 0292 on the default
  branch, `feat: deliver 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it (#619)`,
  whose tree holds the Archive Record
  `docs/history/specs/0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it.md`
  with `pull_request: "619"`: the Pull Request the operator merged by hand
  carried both the archive and the Pull Request Record Commit, which is the
  shape Success Metric 1 reproduces.
- GitHub's documentation of protected branches
  (https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches):
  "GitHub records the state of the diff at the point when a pull request is
  approved. This state represents the set of changes that the reviewer
  approved. If the diff changes from this state […] the approving review is
  dismissed as stale". A review covers the head it saw, which is why a
  resumed head is reviewed again.

## Unreachable Acceptance

- criterion: Success Metric 6, after release, the operator's next new queue
  over an item branch that already archived its Spec reaches `merged` without
  a hand merge.
  reason: only a real Delivery Queue against a real forge produces it, and no
  Verification or QA row may reach the network or run `roundfix deliver`
  against a forge.
  satisfied-by: the operator's `roundfix deliver status` showing that item
  `merged` after the first such queue following this Spec's release.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `34b4e246`, through `roundfix deliver start` in a
  temporary repository with a bare origin, an item branch
  `roundfix/deliver-0001-widget-flow-1111111111111111` holding a completed
  Task, a passing QA report, an archive made by the Archive Command, its
  commit and a Pull Request Record Commit for #619, and an owner built from
  the real workflow with a fake forge. The item parked `delivery-error:
  revalidate item: run strict Spec Consistency Check: read Spec
  "0001-widget-flow": stat …/docs/specs/0001-widget-flow: no such file or
  directory`; without a second active Spec it parked `delivery-error:
  revalidate item: resolve item Specs Root: specs.root resolved to
  "…/docs/specs", which does not exist`. A prototype of the TechSpec's design
  on a scratch copy took both to `merged` with no Implement Run and one review
  of the Pull Request Record Commit, and parked an archive commit with an
  extra file; three sabotages (no resume, no exactness proof, a strict Specs
  Root) each made one case fail. The TechSpec records the code path.
- Secondbrain, read only. Read `wiki/index.md` and
  `wiki/concepts/agendas-e-reexecucao-em-motores-de-workflow-durable.md`
  (durable workflow engines resume without redoing a completed step; a
  catch-up is opt-in), and ran `qmd query "delivery queue resume archived
  spec item branch new queue revalidation" --all --files --min-score 0.3`,
  which returned no file at that score. The page confirms
  that a resume skips completed steps; no prior decision covers re-review on
  resume. No Inbox Entry was written, because the briefing forbids writing
  there.
- Exa: GitHub's protected-branch documentation above (a review covers the
  diff it saw), and durable-execution documentation such as River's resumable
  jobs (https://riverqueue.com/docs/resumable-jobs: "the next attempt can skip
  already-completed steps and pick up where it left off") and a merge queue's
  crash-safety notes (https://lumaprism.github.io/mergequeue/docs/how-it-works:
  a recorded step "re-checks base and completes … (idempotent — if it already
  happened, it's a no-op)"). Challenging: those engines resume from a step
  they recorded, while a new Delivery Queue has no record of the old one, so
  the resume must prove its premise from the branch itself (the exact
  archive), and the one step without such proof, the review, runs again.

## Glossary

- adds: **Archived Item Resume**
- changes: **Delivery Revalidation**

## Decisions

- 2026-10-22 (about 00:55 BRT): the cycle AO scope (multi-select, all four,
  in this order) "Z: portão e corretiva", "Espera de checks", "Retomar item
  arquivado", "Y: reconcile PID" — one Spec and one minor each — and the
  authorization "the same rules as cycle AN".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AO briefing for Spec 0299 of 2026-10-22 and the characterization recorded in the TechSpec.

The maintainer chose the item without a design option. The operator's reading
in the briefing: when the item branch already carries the Spec's archive
commit, adopted by change, a new queue resumes from the stage after archive —
publish or refresh the Pull Request, wait for checks, review per policy,
merge — never re-runs implement or the archive, and a branch whose archive
commit is not exact, or whose Spec folder is missing for another reason, still
parks with a reason that names the cause and the manual path. Taken from that
reading, with no further interview:

- "Review per policy" is the configured pre-PR review of the resumed head, so
  the item enters `reviewing` and the archive step that follows commits
  nothing; the new queue has no evidence of an earlier review of that head.
- "Adopted by change" is the queue's existing archive proof applied to the
  commit that removed the Spec's `_prd.md`, against its first parent.
- The park is a new blocker in the `finding` Park Class, and its retry
  returns the item to `queued`, as a dependency park's retry does.
- One new ADR, ADR-0312, records the rule; the glossary gains **Archived Item
  Resume** and **Delivery Revalidation** says when it does not run.

## Open Questions

None.
