---
spec: 0001-implement-command
status: active
created: 2026-07-04
surfaces: [cli, data]
---

# Implement Command

Roundfix today resolves review feedback on Open Pull Requests, but the Specs produced by the planning workflow — PRD, Task Graph, task files — are still executed by hand, Task by Task. The Implement Command makes Roundfix execute a Spec's Task Graph: it runs Agents over the Tasks in dependency order, verifies and commits each one through the Daemon, and optionally ends with a QA gate — so a completed Spec becomes a verified local branch with one command.

## Goals

- A developer executes a whole Spec locally with one command: each completed Task ends as one commit whose Verification passed.
- A spec Run is a first-class Run: journaled, attachable, and visible in the Live Run View exactly like a review Run.
- Re-running a Spec resumes where work stopped, with no manual bookkeeping.
- The review path (fetch, resolve, watch) behaves exactly as before; its public API is untouched.

## User Stories

1. As a developer with an approved Spec, I want to run the Implement Command with the spec slug, so that the Task Graph executes in dependency order without me driving each Task by hand.
2. As a developer, I want each successful Task to become one commit carrying the code and the updated task file, so that the branch history documents what was done and the evidence for it.
3. As a developer whose Run ended Unresolved, I want a new Run of the same Spec to pick up every non-completed Task, so that resuming costs one command and no manual status edits.
4. As a developer, I want the Run to refuse to start when the Spec is malformed, the working tree is dirty, or I am on the default branch, so that contract problems surface before any Agent runs.
5. As a developer shipping a full feature, I want an opt-in QA step at the end of the Run, so that the implementation is validated against the PRD and the QA Report is committed as evidence.
6. As a developer, I want to Attach to a running or finished spec Run, so that I can follow Tasks progressing or replay what happened.
7. As a developer who omits the spec flag, I want Interactive Input to list the active Specs, so that I can pick the target without remembering slugs.

## Core Features

1. **Implement Command.** Running the Implement Command with a spec slug creates a Run over that Spec in the current repository. Each invocation creates a new Run; spec Runs have no Rounds.
2. **Preflight Validation.** Before a Run is created, Roundfix validates: the Spec exists and is active; the Task Graph manifest parses and is acyclic; every referenced task file exists with parseable frontmatter; every Task has a Verification section; the working tree is clean; the current branch is not the repository default; and no Active Run exists for the same work target or the same working tree. Failures name the offending Task or check and the next useful action.
3. **Sequential execution.** Tasks run one at a time in a topological order of the Task Graph. Each Task is assigned to one Agent invocation as a Batch of one, using the ACP Runtime selection and configuration that review Runs already use.
4. **Task prompt.** The Agent receives the task file's content plus the execution invariants: edit code, run Verification while working, update task status and append the Result section, never commit. The invariants mirror the implement-task contract so agent-side and daemon-side expectations cannot drift.
5. **Daemon-owned verification and commit.** After the Agent finishes a Task, the Daemon runs that Task's Verification commands verbatim. Passing verification gates one commit containing the code changes and the updated task file, with a Conventional Commits message derived from the Task's frontmatter type and title and a trailer naming the Spec and Task. See ADR-0013, ADR-0014.
6. **Status settling.** The Agent writes task status and Result while working; the Daemon settles the final status. Completed requires passing verification; anything else is settled failed with the reason journaled. See ADR-0014.
7. **Failure policy.** A failed Task produces no commit; its dependents stay pending and are never executed with unsatisfied needs; independent Tasks continue. A Run that ends with any non-completed Task ends in the Unresolved Outcome; Stop Requests and infrastructure errors end it as Failed, generalizing the ADR 0010 semantics.
8. **Resume.** A new Run over the same Spec picks up every non-completed Task: pending, failed, and in_progress left behind by a dead Run. There is no retry inside a single Run.
9. **QA gate, opt-in.** With the QA flag, and only when every Task is completed, the Run ends with a qa-gate step. The verdict is read from the QA Report frontmatter; a failing or unreadable verdict ends the Run in the Unresolved Outcome; the QA Report is committed in its own commit either way. See ADR-0015.
10. **No push.** Spec Runs never push and never open pull requests; handing the branch to a pull request — and from there to the review path — is the developer's explicit decision after the Run. See ADR-0013.
11. **Run visibility.** Spec Runs write Run Events to the Run Event Journal; Attach and the Live Run View render Tasks as the Work Items in place of Review Issues.
12. **Interactive Input.** Without the spec flag, Interactive Input lists the repository's active Specs for selection, following the same flow as the existing commands.
13. **Review path unchanged.** fetch, resolve, and watch keep their exact behavior, flags, exit codes, and outputs; the work-target generalization is invisible to review users.

## User Experience

- In non-interactive mode the command reports the Run outcome and per-Task results deterministically on stdout; diagnostics and progress go to stderr; outcomes map to stable exit codes consistent with the existing commands.
- In interactive mode the Live Run View shows the Run: Tasks with their statuses, the Run Event timeline with Follow Mode, and the terminal outcome.
- Preflight failures read as one actionable message naming the check that failed and, when known, the fix — for example, which task file lacks a Verification section.

## Non-Goals / Out of Scope

- Parallel Task execution, ready-set scheduling, and git worktree isolation (work-plan item 4).
- Retry budgets and Agent escalation inside a Run (item 7).
- Versioned or per-repo templated prompts (item 5).
- Structured result artifacts beyond the task file's Result section (item 6).
- A long-lived daemon process, work queue, filesystem or webhook triggers, or a watch mode for Specs (item 3).
- Permission policy changes (item 8).
- ACP layer migration to acpx (item 2).
- Pushing, opening pull requests, or driving the review path from a spec Run.
- Running QA over partially implemented Specs.
- A unified Task Source interface shared by the review and spec paths — the paths stay siblings behind their own seams.

## Success Metrics

- Dogfood: a real Spec in this repository executes end-to-end via the Implement Command, producing one passing-verification commit per completed Task.
- A resumed Run after an induced Task failure picks up only non-completed Tasks and finishes the graph.
- The full existing test suite passes unchanged for the review path after the core generalization.
- QA opt-in produces a committed QA Report whose verdict settles the Run outcome.

## Decisions

- "Run" generalizes to work targets (Open Pull Request or Spec); one Active Run per target; a single-working-tree lock applies until worktree-per-task lands. See ADR-0012.
- Spec Runs commit per Task and never push; they work on the current branch with a default-branch veto; the commit message derives from the task frontmatter. See ADR-0013.
- The Daemon runs Task Verification verbatim and settles task status; a missing Verification section fails Preflight Validation. See ADR-0014.
- QA runs inside the Run behind an opt-in flag, only with all Tasks completed; the verdict comes from the QA Report frontmatter; the report is always committed. See ADR-0015.
- Round remains a review-path concept; a spec Run is one cycle over the graph per invocation.
- Each Task is one Batch of one; grouping several Work Items into one Batch does not apply to Specs.
- The task prompt is a minimal dedicated builder mirroring the implement-task contract; templating waits for work-plan item 5.
- Interactive Input gains a Spec picker for parity with the existing commands.

## Open Questions

- Task frontmatter `type` values (frontend, backend, data, infra, docs, test, chore) are surfaces, not commit types. Default mapping until revisited: docs → docs, test → test, chore → chore, everything else → feat. Owner: Marcio.
- The QA Report verdict frontmatter (field name and values) must be added to the upstream qa-gate skill contract. Default until agreed: the Daemon treats a missing or unreadable verdict as fail. Owner: Marcio.
