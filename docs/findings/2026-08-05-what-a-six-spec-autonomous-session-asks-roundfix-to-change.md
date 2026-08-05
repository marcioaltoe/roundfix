# 2026-08-05 — What a six-Spec autonomous session asks Roundfix to change

status: pending

## What was observed

One Supervisor drove six Specs from findings to merged Pull Requests in the
fluxus repository between 2026-08-03 and 2026-08-05: 30 Tasks, five merged Pull
Requests (#78–#82), 57 review issues resolved. Roundfix implemented every Task
correctly. The session's cost sat almost entirely elsewhere.

Ten Implement Runs were started. Measured by settled work:

| Run | Outcome | Cause |
| --- | --- | --- |
| 0012 #3 gate | fail | one missing ADR citation named |
| 0012 #4 gate | fail | a *different* missing ADR citation named |
| 0015 #1 | 4 of 7 Tasks failed | Verification used the wrong test runner |
| 0015 #2 | 0 completed, 3 failed | `agent_full_access` rejected by the adapter |

Four Runs produced no settled Task. None of them found a defect in production
code. Every one was a property Roundfix could have checked before the Run
existed.

Five specific findings from this session are already filed:
[[2026-08-04-fail-fast-verification-spends-the-single-repair-turn-on-the-first-of-n-defects]],
[[2026-08-04-a-static-gate-row-reported-one-instance-per-cycle]],
[[2026-08-04-review-runs-halt-autonomous-delivery-on-unrelated-dirty-files]],
[[2026-08-04-a-spec-archives-with-pass-while-a-user-story-was-never-exercised]],
[[2026-08-04-watch-derives-a-review-head-it-never-checks-is-reachable]],
[[2026-08-05-agent-full-access-passes-config-validation-and-fails-every-task]],
[[2026-08-05-archive-refuses-a-spec-whose-graph-declined-the-qa-gate]], and
[[2026-08-05-preflight-prescribes-integrating-a-superseded-run-branch]]. This
one is the behavioural argument they share.

## The three behaviours worth changing

### 1. Report the class, not the instance

The dominant cost was checkers that stop at the first problem while feeding a
loop with a bounded number of attempts.

Three consecutive QA gates failed the same row, `PC-03`, each naming exactly one
missing ADR. Patching what the report named bought exactly one more cycle. A
manual sweep found **eleven** missing, two of them structural. Gate #4 ran ten
minutes to report a single citation.

The same shape appears one altitude down. `make verify` is fail-fast, so Task 07
of Spec 0012 got its one repair turn against an architecture failure, fixed it
honestly, and then attempt 2 revealed a second, independent defect the Agent had
never been shown. There is no third attempt. The Task settled failed on a defect
nobody told it about.

**Behaviour asked for:** where a check can enumerate its violations, enumerate
them. Where Verification is a pipeline, run every step before the verdict so the
single repair turn sees the whole picture. Exhaustiveness is not a reporting
nicety here — it is the difference between one Run and five.

### 2. Prove a setting against the runtime that must honour it

`defaults.agent_full_access: true` passed config validation and then failed
every Task in the Run, before any Agent work, with
`session/set_mode "full-access"` answering ACP -32602. `roundfix doctor`
reported every check green immediately beforehand.

Doctor's readiness contract is unusually thorough — Node, acpx, adapter lineage
and version, every Agent Selection Profile tuple proved through disposable
sessions, the Repository Skill Set, codex hygiene. It opens real sessions
already. It just never asks whether the session mode the config demands is
supported.

The asymmetry is the point: a malformed config value is caught before a Run
exists; a well-formed value the adapter cannot honour costs the entire Run, with
no fallback, because this setting sits outside the Agent Selection lifecycle
that fallback protects.

**Behaviour asked for:** every configuration value that a runtime must accept
gets proved the way profiles are proved, and fails at preflight with
`No side effects` rather than per Task after Run creation.

### 3. Never prescribe an action that can destroy work

Branch Integrity Preflight refused a review and printed:

```
Next action: inspect each pending Run Worktree, then run the listed integration
command from the repository root when it is safe.
integration_command="git merge --ff-only roundfix/run-..."
```

Running it would have been wrong. The pending commit was a Task from a failed
Run; the surviving version on the target branch was 272 lines further along,
having passed through three later Tasks and a QA gate. Both commits carry the
same subject. Only the skill's prose warning — verify a stranded branch before
discarding it — prevented an operator from following the printed instruction and
overwriting newer work.

`reconcile` could not help: a failed Run almost always leaves a dirty Worktree,
which classifies `dirty` and is preserved, and `--apply` acts only on `safe` or
`superseded`. The supported exit was `--skip-branch-integrity`.

**Behaviour asked for:** a prescribed command must be safe to run, or must not
be prescribed. Where supersession can be proven — a pending branch whose Tasks a
later Run settled `completed` on the target — classify it and let reconcile
release it. Where it cannot, name the risk in the refusal instead of printing a
bare fast-forward.

## What made the difference, and what Roundfix could own

The session ran in fits and starts while every policy decision stopped the loop:
whether a `partial` verdict is acceptable, whether migrations are authorized,
whether a Pull Request may merge, how many corrective Tasks are allowed, what to
do when a review does not converge. Nine such stops.

Then the maintainer was asked for all of them at once. After three answered
questions, **two complete Specs were implemented, gated, reviewed, and merged
with no further stop**.

Seven of the nine were policies, not case-by-case judgements. They are stable
per project, they are exactly the escalation triggers `autonomous-work.md`
enumerates, and today they live only in conversation — so every session
rediscovers them, and every session pauses to ask.

**Behaviour asked for:** let a repository declare its autonomy policy as
configuration Roundfix reads, in the same place Verification and profiles live.
A minimum set, drawn from what actually stopped this session:

```yaml
autonomy:
  accept_qa_verdict: [pass, partial_environment_only]
  corrective_task_budget: 2
  on_non_convergence: merge_if_checks_green   # or: leave_open | stop
  migrations: authorized                       # or: ask | forbidden
  merge: squash_to_default_branch              # or: open_pull_request_only
```

Roundfix already owns the moments where each of these is consulted: the gate
verdict, the corrective loop, the review outcome, the Task that carries a
migration, the merge. Reading them from config instead of stopping is what turns
a loop that works into a loop that chains.

The stop conditions do not disappear — they narrow to what genuinely deserves a
person: a Spec that is wrong rather than thin, an unauthorized tooling mutation,
an acceptance that needs a credential the repository cannot hold. Those are the
triggers worth waking someone for. The other seven are not.

## Evidence

- fluxus `main`: `ad7162fc`, `8676d8fc`, `7b9d23c3`, `6ad191f3`, `619cac5c`.
- Runs `run_20260803T235114Z_635c89f240d6c6e4`,
  `run_20260804T111316Z_fbfe4853530bf2c5`,
  `run_20260804T181023Z_f9fef8003e217a60`,
  `run_20260804T225948Z_affa3833d7d175d6`.
- QA reports under fluxus `docs/specs/_archived/0012-…/qa/`,
  `_archived/0014-…/qa/`, `_archived/0015-…/qa/`.
- Measured on codex-acp 1.1.9, acpx 0.13.0, with `roundfix doctor` green.

## Spec pointer

`0064-spec-artifact-consistency-gate` covers the static half of behaviour 1.
Behaviours 2, 3, and the autonomy policy have none.
