---
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: active
created: 2026-10-01
surfaces: [backend, cli, docs]
---

# Gates and Run storage that let a correct delivery finish

Between 2026-09-25 and 2026-10-01, four deliveries whose work was correct
stopped on a gate or on the disk, and each needed an operator. The pre-PR
review of Spec 0200 failed three times with `agent/protocol error`, twice in
the Delivery Queue and once by hand, on a candidate diff of 1,295,055 bytes,
and someone reviewed it by hand. In Spec 0203's delivery a QA Agent wrote
its evidence as a Go test file; the next QA pass imported it byte for byte,
the repository gate refused to format it before any row ran, and every retry
would have imported it again. Spec 0200's task_04 deletes a skill file, but
a Task's `## Context` has no kind for a deleted path, so the Spec Consistency
Check refused the correct Task once it ran, and the operator relabelled the
entry as a path the Task creates. And `roundfix reconcile` keeps every Run of
a Spec merged by squash: on 2026-10-01 this repository held 32 Run Worktrees
of merged Specs, 19.1 GB, which contributed to a disk-full stop of a delivery.
This Spec lets each of those correct deliveries finish: the review reads the
diff it can review and says what it left out, the QA import leaves compiled
source behind, a Task can declare a deletion, and a merged Spec releases the
Runs it superseded.

## Prerequisites

This Spec is delivered after Spec 0211 (a delivery queue that finishes without
intervention). Both raise the Roundfix Skill's version, so this Spec's Task
Graph names Spec 0211 in `requires` and the queue owner waits for it to merge.
Spec 0210 changes the Evidence Snapshot that the QA import carries; this Spec
changes which files the import copies and leaves the snapshot alone, so it
needs neither order nor rebase against Spec 0210.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Runs keep their
  Run identifiers, Specs their slugs and review records their checkout key.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the review keeps its configured
  ACP Runtime and reads the candidate from local Git; the QA import, the Spec
  Consistency Check and reconcile read only local files and Git. No credential
  and no network call is added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0212 (this Spec) refines ADR-0053
  and ADR-0161 for a merged Spec's Runs, ADR-0212: "its archived directory now
  represents every path a Run changed under the Spec's directory".
  ADR-0115's Branch Disposition, ADR-0170 and ADR-0052 hold unchanged. The
  review keeps ADR-0153's provider policy, ADR-0169's merge-base diff,
  ADR-0174's record per checkout, ADR-0196's validation and ADR-0197's two
  rounds. The QA import keeps ADR-0194's rules, and ADR-0210's per-input
  digest and ADR-0097's carry conditions decide whether a row citing a file
  left behind still carries. The `deletes` kind joins the declarations
  ADR-0166 counts and ADR-0178 authorizes. ADR-0184 has the TechSpec state
  the changed command surface as a transcript, ADR-0184: "A TechSpec now
  declares numbered Surface Transcripts", and ADR-0187 and ADR-0189 govern
  the skill edits. The gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0104, ADR-0155, ADR-0156 and ADR-0167. ADR-0195's always-observed rows
  hold unchanged for this Spec's own gate, and ADR-0096's mechanical stage
  keeps its role. ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  this Spec's consistency by citation and receipt. ADR-0182 does not apply,
  because no Task settlement fact changes, and ADR-0211 does not apply,
  because it governs the agent environment Spec 0211 changes. ADR-0165 and
  ADR-0192 hold unchanged: a blocking review after archive still parks for a
  corrective Spec, and a derived-path conflict is still resolved by
  regeneration. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix, qa-gate and write-tasks
  skills under `.agents/skills/` and their `SKILL.md` mirrors are Governed
  Paths, and the maintainer authorized skill edits ("considere autorizado a
  ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `skills/roundfix/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Goals

- A delivery candidate whose reviewable diff fits the bound measured on this
  repository is reviewed by the configured provider, and the review record
  lists every path it did not read and why.
- A QA Agent's evidence can never fail a later pass's repository gate through
  the failed-pass import.
- A Task that deletes a file declares it as a deletion, and the Spec
  Consistency Check, the scope audit and the Wave collision rule treat it as
  declared.
- A terminal Run of a merged, archived Spec is released by reconcile when its
  only unrepresented work is under the Spec's directory or inside the Spec's
  declared scope.

## User Stories

1. As the operator of the Delivery Queue, I want the pre-PR review to omit
   QA evidence and upstream-managed skill copies from the diff it sends, so
   that a Spec that refreshes vendored skills or records large evidence is
   still reviewed by the configured provider.
2. As the operator, I want a blocked review to say how large the diff was and
   what the runtime last wrote, so that I can tell a size failure from a
   runtime outage.
3. As the operator, I want the next QA pass to leave behind a source file the
   previous pass wrote as evidence, so that the delivery does not loop on a
   gate failure no Task caused.
4. As a Spec author, I want to declare a path my Task deletes, so that the
   declaration says what the Task does and the gate accepts it.
5. As the maintainer, I want reconcile to release the Runs of a Spec the
   default branch merged and archived, so that storage stays bounded without
   deleting Roundfix worktrees by hand.

## Core Features

1. **The review omits what it cannot usefully read.** The pre-PR review
   leaves out of the diff it sends the files under any Spec's `qa/evidence/`
   directory and the files under each upstream-managed skill directory the
   repository's skills lock names. The prompt lists the omitted paths by
   reason, so the reviewer does not raise their absence.
2. **A bounded diff, and a named cause when it does not fit.** When the
   remaining diff exceeds the review bound, the review records `blocked`
   before any provider call, with a reason naming the diff size and the
   bound. The bound sits below the largest diff this repository reviewed
   successfully.
3. **The record says what it did not read.** The review record lists each
   omitted path with its reason, the reviewed diff's size in bytes, and, for
   a runtime failure, the last lines the runtime wrote to standard error.
4. **The import leaves compiled source behind.** When the next QA pass imports
   a failed pass, it skips each file whose extension a supported toolchain
   compiles or formats, records each skipped path with the reason
   `compiled source`, and imports the rest. A row whose evidence was skipped
   is re-run, never carried.
5. **The qa-gate skill keeps runnable evidence out of compiled paths.** The
   skill tells the QA Agent to store runnable evidence under a non-compiled
   extension, such as `.go.txt`.
6. **A Task can declare a deletion.** A Task's `## Context` accepts
   `deletes: <path>`. The unresolved-path check skips it; the Daemon does not
   settle a Task while its `deletes` path still exists; the scope audit, the
   Governed Path declaration audit and the Wave collision rule count it as
   declared. The write-tasks skill and its template teach it.
7. **A merged Spec releases its Runs.** When the merged head holds the Spec
   archived, reconcile treats every path a Run changed under the Spec's
   directory as represented, and a Run Worktree's uncommitted changes as
   superseded when each changed path is declared or recorded by one of the
   Spec's Tasks or lies under its directory (ADR-0212). A Run whose Run Branch
   was left without a target branch gets the same merged-head proof. Every
   other unrepresented path still preserves the Run.
8. **The skill and the guides say so.** The Roundfix Skill's `review` and
   `reconcile` references and the `review` and `reconcile` command guides
   describe the omission, the bound, the record fields and the release rule.

## User Experience

`roundfix review` prints its outcome as today. A blocked review over the
bound prints a reason naming the size and the bound, and its record carries
the omitted paths and the diff size. `roundfix reconcile` reports a merged
Spec's released Runs as `superseded` with a reason naming the archived Spec
at the merged head; `--apply` removes them as it removes any superseded Run.
A Task whose `deletes` path still exists fails its settlement attempt with
`deletes: <path> still exists`, and the Agent gets the same bounded retry as
for a failed Verification command. Nothing else on the command line changes.

## Non-Goals / Out of Scope

- Splitting a review into parts, or reviewing a diff above the bound by any
  other means; a later Spec may add that once a case needs it.
- Changing the review provider policy, its model, its two-round lineage or its
  finding validation.
- Rewriting an imported QA Report, or the Evidence Snapshot Spec 0210 owns.
- Releasing any Run whose Task commit is not represented, or any uncommitted
  path outside the Spec's declared or recorded scope.
- Reading Pull Request state from GitHub in reconcile, or deleting a Delivery
  Queue item's directory.
- The Delivery Retry, checking and merging stages and the agent environment:
  Spec 0211 owns them.
- Editing `CONTEXT.md`; the glossary check is the QA gate's.

## Success Metrics

1. Success Metric: Spec 0200's reviewed candidate of 1,295,055 bytes, with
   786,524 bytes of QA evidence and 48,124 bytes of skill copies under
   `.agents/skills/`, sends no more than 460,407 bytes of diff once evidence
   and upstream-managed skill copies are omitted, and the record lists every
   omitted path.
2. Success Metric: a candidate whose remaining diff exceeds the review bound
   is recorded `blocked` with the size and the bound, and no provider call is
   made.
3. Success Metric: a failed QA pass that holds a Go source file under
   `qa/evidence/` is imported without that file, the skip is recorded with
   `compiled source`, and the next pass's repository gate is not affected by
   it.
4. Success Metric: a Spec whose Task declares `deletes: <path>` for a removed
   path reports no `SC-REF-UNRESOLVED`, and a Task that leaves the path in
   place does not settle `completed`.
5. Success Metric: a terminal Run of a merged Spec whose only differences are
   Spec-directory commits and uncommitted changes to the Spec's declared paths
   is classified `superseded`; one with an uncommitted change to an
   undeclared file outside the Spec's directory stays preserved.

## Acceptance evidence

The outside-evidence rows rest on records this Spec did not produce:

- The review records under this repository's Artifact Directory, read on
  2026-10-01: Spec 0194's candidate (`c49a0128`..`e8fad6f7`) was reviewed with
  outcome `reviewed` on a diff of 919,745 bytes, and Spec 0200's candidate
  (`64aff3f7`..`954ad599`) was recorded `blocked` with `review runtime
  failure: Agent Selection failed for runtime "codex": agent/protocol error`
  on a diff of 1,295,055 bytes. Both commit pairs are in this repository's
  history, so `git diff` reproduces the sizes and the breakdown in Success
  Metric 1.
- Spec 0203's archived QA evidence: `qa/evidence/ledger_replay_test.go` is
  the compiling placeholder the operator committed, and
  `ledger_replay_test.go.txt` keeps the original the QA Agent wrote.
- Spec 0200's archived task_04, which declares the deleted
  `.agents/skills/context7/SKILL.md` as `creates`.
- The Secondbrain note of 2026-09-25 from the Vortex repository,
  `inbox/roundfix/_triaged/2026-09-25-reconcile-nunca-libera-run-de-spec-mergeada-por-squash.md`:
  13 preserved Runs and 11 GB, classified `unintegrated` and `dirty`, after
  squash merges and archive moves.
- This repository's Run Worktrees, measured read-only on 2026-10-01: 32 Run
  Worktrees of 24 merged and archived Specs, 19.1 GB.

## Decisions

- The review omits two classes of path instead of splitting the diff: the
  measurement shows QA evidence, not vendored skills, made Spec 0200's diff
  too large, and the reviewable remainder fits the measured bound.
- The bound is a constant set below the largest diff reviewed successfully on
  this repository, and exceeding it blocks with a named cause rather than
  reaching the runtime.
- The import skips compiled source rather than refusing the whole pass, so the
  rest of the failed pass's evidence still carries.
- A deleted path is a fourth Context kind rather than a flag on `creates`,
  so a declaration says what the Task does.
- A merged Spec supersedes its Runs' Spec-directory work and declared
  leftovers. See ADR-0212.
- The four sources share one Spec under the rule that sources sharing a
  context share a Spec: each is a gate or a resource a correct delivery
  stopped on.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
