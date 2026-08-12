---
status: pending
created_at: 2026-08-12
updated_at: 2026-08-12
kind: finding
---

# Run lifecycle — a queue of eight Specs shows where the loop breaks (2026-08-12)

A session in `oraculum` on 2026-08-07/08 ran 25 Runs across 8 Specs: 3 delivered
and merged (0026, 0028, 0029), 1 ready and blocked, 2 frozen, 2 deferred on a
temporal dependency, and release `v0.4.2` published. Minted here from the Inbox
Entry `inbox/roundfix/2026-08-08-fila-0026-0035-defeitos-e-lacunas-medidos.md`
in the Secondbrain. Every item below has evidence in that repository's Runs and
handoffs; the ordering is by real cost, not by theoretical severity.

## 1. The declared-only `partial` path does not close, and blocks delivery today

- Symptom / evidence: `roundfix archive --help` promises archival after
  "every Task is completed" **or a partial verdict whose blocked rows are covered
  only by declared Unreachable Acceptance". The two conditions are cumulative and
  the second never occurs alone: the Daemon settles the QA Task from the verdict,
  `partial` is not `pass`, so the Task stays `failed` and the first condition
  fails. Spec 0027's fourth report was exactly the described case — 15 pass, 9
  declared-blocked, 0 environment-blocked, 0 fail, 0 finding-blocked, 0 skipped,
  0 pending — and archive refused:
  `Task "task_07" is "failed"; archive requires every Task to be "completed"`.
  There is no bypass flag; `archive` accepts only `<slug>`, and `qa_override` is
  explicitly forbidden for the declared case.
- Root cause: the settlement rule and the archive rule read the same verdict and
  disagree about what a declared-only `partial` means.
- Action / suggestion: have the Daemon settle the QA Task as `completed` when the
  `partial` is declared-only — the same check `archive` already knows how to do.

## 2. Parallel bootstrap collides on the shared `.git`

- Symptom / evidence: with `worktree.concurrency > 1`, sibling Task Worktrees run
  bootstrap simultaneously against the repository's shared `.git`, and a
  `prepare` writing `git config` collides on the lock:
  `error: could not lock config file …/.git/config: File exists` followed by
  `error: prepare script from "oraculum" exited with 255`. The symptom misleads:
  the bootstrap **reports failure after having done all the work** — the torn
  down Task Worktrees had both `.env` files at mode 600 and 578 packages
  installed. Cost: two Tasks of Spec 0027.
- Root cause: bootstrap is about the worktree, and concurrent installation over a
  shared cache and a shared `.git` is a known hazard.
- Action / suggestion: serialize Worktree Bootstrap even when Tasks run in
  parallel.

## 3. An integration conflict is discovered too late

- Symptom / evidence: two Specs (0032, 0035) lost whole Runs — all the Agent's
  work done, failing only at cherry-pick:
  `reason: integration conflict: docs/references/planejamento/scripts/build_pdf.py, build_xlsx.py, …`
  and `reason: integration conflict: packages/backend/src/infra/controllers/mcp/schemas-comandas.ts`.
  In both cases the graph declared independent Tasks editing the same file. At
  `concurrency: 1` it never appeared; raising concurrency made visible an error
  the graphs already had.
- Root cause: wave independence is declared by the author and never checked
  against the files each Task touches.
- Action / suggestion: Roundfix knows which files each Task Worktree changed.
  Detect set intersection before dispatching the wave and serialize, or at least
  report "these Tasks touched the same files" so the Supervisor fixes the graph
  before the next Run.

## 4. Preflight proves fallbacks, making an unused runtime a hard dependency

- Symptom / evidence: with cross-runtime fallbacks, every Run must prove every
  configured tuple and substitutes none. In a configuration with 5 distinct
  tuples, 3 of them claude, an intermittent claude adapter refused several Runs
  whose preferred (codex) was perfect:
  `profile proof failed for runtime "claude", model "opus", reasoning_effort "high"; affected categories: backend fallback[1], qa fallback[1]`.
- Root cause: the root cause of the intermittence was an asymmetry in
  `~/.acpx/config.json` — codex by local binary, claude by `npx -y …`, paying
  package resolution on every disposable session — but the design amplifies it: a
  fallback that will never be used blocks the start.
- Action / suggestion: always prove the preferred Selection, and prove fallbacks
  lazily or tolerantly, or allow an opt-out by configuration.

## 5. The largest lever: `spec check --run-verification`

- Symptom / evidence: a mode that executes each `## Verification` line in a
  disposable checkout and reports the exit code would alone have caught **six
  defects in one night**, all of assumed shell semantics:

  | Defect | Effect |
  | --- | --- |
  | `grep -rc PATTERN file` | prints `file:0`, not `0` |
  | `grep -q … -l \| head -1` | exit code of `head`, always 0 — empty verification |
  | `wc -l` on macOS | pads with spaces: `"       0" != "0"` |
  | wrong file path | assertion on the schema, code in the use case |
  | command with `pytest` absent | fails on tooling, not on a defect |
  | glob `qa-report-*.md` inside `grep` | passes if **one** file matches; an old report satisfies a rerun |

  Four of them **failed work that was correct** and cost Runs.
- Root cause: `SC-VERIFY-WORK-INDEPENDENT` checks the form; nothing checks the
  execution.
- Action / suggestion: add the execution mode. This is corroborated
  independently by the `fiscus` report of the same date.

## 6. Roundfix should own the QA Task's Verification

- Symptom / evidence: the Daemon already derives the verdict. Letting the author
  write a predicate over the report produced, in Spec 0035, a gate that
  **passed itself having failed**: `^verdict: (pass|blocked)` — rejecting `fail`
  and `partial`, which are valid verdicts, and accepting `blocked`, which does
  not belong to the domain — plus a glob matching any report in the folder. The
  gate itself caught both (F-004 and F-005 of the 0035 report).
- Root cause: a derived verdict is re-asserted by hand at the one place where the
  author has the least information.
- Action / suggestion: Roundfix owns the QA Task's Verification.

## 7. Force stop leaves an orphan grandchild

- Symptom / evidence: a `bun test` survived `roundfix stop --force` and ran for 47
  minutes at 99% CPU (`STAT R`, zero network connections).
- Root cause: Force Stop proves the exit of the registered owner, not of the
  processes it spawned.
- Action / suggestion: prove the process tree, not the owner alone. Sibling of
  `2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md`.

## 8. `reconcile` is scoped by repo-id, and a worktree of the same repository gets its own

- Symptom / evidence: a `git worktree` of the same repository received a distinct
  `repo-id` (`oraculum-wt-0029-3e814276` against `oraculum-d1255e75`).
  Consequences: `runs list` without `--all` does not see the worktree's Runs, the
  Branch Integrity Preflight on one side sees Run Branches from the other
  (forcing `--skip-branch-integrity`), and removing the checkout left an
  irreconcilable Run that had to be recreated just to run `reconcile`.
- Root cause: repository identity is derived per checkout rather than per
  repository.
- Action / suggestion: resolve one identity for a repository and its worktrees.

## 9. Instruction gaps in the context-driven baseline

- Symptom / evidence: `## Unreachable Acceptance` is not documented —
  `autonomous-work.md` cites ADR-0080 and Spec 0078, but where the section lives
  (the PRD), the entry format and the meaning of `satisfied-by` all had to be
  inferred, and it is the mechanism that decides whether a Spec delivers or
  freezes. A clause is missing that a Verification line is executed before it is
  written (see finding 5). `write-tasks` declares wave independence without
  looking at files (see finding 3). A Spec whose acceptance is post-release does
  not declare itself as such — three Specs (0030, 0031 and part of 0027) only
  close after a release and an observation window, and nothing in routing signals
  it at authoring time.
- Root cause: rules learned in execution have no upstream owner in the Baseline
  modules.
- Action / suggestion: carry each as a Baseline module clause, and add a temporal
  prerequisite field so the queue orders itself.
