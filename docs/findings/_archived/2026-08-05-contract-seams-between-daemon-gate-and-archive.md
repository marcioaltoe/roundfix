---
status: done
created_at: 2026-08-05
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
---

# 2026-08-05 — Contract seams between Daemon, gate, and Archive

status: pending

Source: three consecutive Specs delivered end to end in the `fiscus` repository
(`0001-fundacao-de-contextos`, `0002-auth-staff-e-directory`,
`0003-design-system-graphite`) — 35 Tasks, 17 Implement Runs, 14 gate executions, 7 merged pull
requests. Companion to
`2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md` and
`2026-08-05-authored-verification-gates-are-untested-code.md`, which cover the gate-authoring
and worktree classes. This one records what the **components disagree about**: the Daemon, the
QA gate, and Archive each hold a valid contract, and the seams between them cost this session
five extra Runs and four supervisor edits to Daemon-owned artifacts.

Cross-checked against `fluxus`'s same-day finding
(`docs/findings/2026-08-05-o-que-travou-a-entrega-autonoma-de-seis-specs.md`), which measured
the same classes independently — the observations below are not one repository's accident.

## 1. Three ways the settle/verdict/archive contracts disagree

Each was hit in sequence; each forced the supervisor to hand-edit an artifact the Daemon owns.

- **Non-`pass` gate settles `failed`; Archive requires `completed`.** Spec 0002's final report
  closed `verdict: partial` with zero findings and every unmet row covered by a declaration —
  the exact shape Archive documents as eligible ("an archive-eligible partial Report can still
  leave the Run Unresolved; that outcome does not prevent Archive"). But the Daemon settles any
  non-`pass` gate `failed`, and Archive's first check demands every Task `completed`. The
  documented eligible state is unreachable without editing the gate Task's status by hand.
- **The gate emits `pass` with a declared row; Archive rejects it.** Spec 0003's final report
  closed `verdict: pass`, `rows_blocked_declared: 1`, zero findings — and Archive refused:
  `rows_blocked_declared must be zero when verdict is "pass"`. The gate's own final-verdict
  prose argued the opposite ("a nonzero declared-block count does not contradict this gate's
  pass contract"). Both components are internally consistent and mutually exclusive; the
  supervisor corrected the frontmatter to `partial` with a written note.
- **Declaration matching is gate judgement, but Archive counts the label.** With a declaration
  covering "rows requiring loopback PostgreSQL", the gate still labelled the repository-
  Verification row `blocked (environment: …)` — after itself re-running that command green
  outside the sandbox — and Archive then refused on `rows_blocked_environment is 1; expected 0`.
  The only recourse was rewording the declaration to name the row class and prescribe the label,
  then re-running the whole gate for accounting.

Suggested behavior: make one component the authority. Either the gate computes the verdict from
the counts Archive checks (declared rows ⇒ `partial`; environment rows ⇒ not archivable) and
settles an archive-eligible gate `completed`, or Archive derives eligibility from the counts and
stops requiring a status the Daemon cannot produce. Additionally, when a blocked row's blocker
text matches a Spec declaration, the gate should classify it `declared` rather than leaving the
classification to prose judgement.

## 2. `## Unreachable Acceptance` parses a shape the skill never teaches

- What was observed: Spec 0002's declarations were written as prose bullets and Archive refused
  with `unreachable acceptance declaration … is missing criterion`. The parser
  (`internal/spec/spec.go`) requires `criterion:` / `reason:` / `satisfied-by:` labelled fields,
  and the `criterion` value must match the QA report's `blocked (declared: …)` label for the
  row to be credited. Nothing in the authoring skills states this; it was recovered by reading
  Roundfix's source.
- Suggested behavior: document the declaration grammar where PRDs are authored (a template
  fragment in `write-prd`), and have the Archive refusal print the expected shape rather than
  only the missing field name.

## 3. `SQLITE_BUSY` fails Agent Batches that produced complete work

- What was observed: three times across two Specs, a Batch died with
  `Agent failed: Agent Batch failed: publish Run Events: begin Run Event append: database is
  locked (5) (SQLITE_BUSY)` — 0002/task_08, 0003/task_05, 0003/task_07. Each time the Agent had
  finished: the kept worktree held complete, correct work that `roundfix settle` committed
  without a single change. In every case a supervisor `events --follow` stream was reading the
  Run Database concurrently.
- Suggested behavior: bounded retry with backoff on Run-Event appends (and WAL mode if not
  already on). A journal write that loses a lock race is not an Agent failure and should not
  discard a Batch. This is the single most frequent failure mode of unattended operation in
  this session.

## 4. Squash-merge makes every Run Branch permanently `unintegrated`

- What was observed: after merging each Spec's pull request with squash, `roundfix reconcile`
  classified the Run Branches `unintegrated` forever — their tips can never be ancestors of the
  rewritten default branch. Nine worktrees accumulated; cleanup required manual
  `git worktree remove` plus `git branch -D`, and one leftover directory (64 MB of `.storage`)
  survived even that because no branch pointed at it.
- Suggested behavior: a content-equivalence classification (prove every file of the Run Branch
  exists identically in the target, as the supervisor did by hand with `git ls-tree` before
  deleting) or an explicit `--superseded-by-squash` mode. Squash-only repositories are common;
  today they can never use `reconcile --apply` for their Run Branches.

## 5. Preflight profile proof has no retry where `profiles validate` succeeds

- What was observed: twice, `roundfix implement` refused with
  `profile proof failed for runtime "claude", model "opus", reasoning_effort "xhigh"; adapter
  error: context deadline exceeded` — and running `roundfix profiles validate` immediately
  after passed the same tuple. One occurrence also failed while *closing* the disposable
  session. Both blocked an unattended launch until a human re-ran two commands.
- Suggested behavior: retry the tuple proof on a timeout before refusing, or treat a proof
  timeout as a warning when the same tuple proved recently. A transient adapter deadline is not
  a configuration defect, and preflight explicitly substitutes no fallback.

## 6. The Daemon's staging list can name paths deleted during the turn

- What was observed: a task_04 Agent created `components/.gitkeep`, replaced it with the real
  component, and removed it; the Daemon's commit still listed the file and exited 128
  (`fatal: pathspec … did not match any files`), failing the whole Run **after** the Task
  settled `completed`. Recovery cost a full re-run of a finished Task.
- Suggested behavior: build the staging list from the tree at commit time, or filter
  non-existent paths; and settle the Task `failed` (recoverable through `settle`) instead of
  failing the Run when a commit cannot be created.

## 7. The Baseline should carry an autonomy charter as a first-class artifact

- What was observed: this session stopped sixteen times for maintainer input; ten of those were
  standing policy rather than case judgement — accepting an environment-only `partial`, merging
  when review cannot converge, authority over trivial review threads, bounded tooling authority,
  corrective-cycle budget. The `fluxus` repository measured seven of nine with the same shape
  and answered it by authoring `docs/agents/autonomy-charter.md`, a repository-owned standing
  policy that grants authority without weakening the escalation triggers in
  `autonomous-work.md`. After it existed there, two Specs went through implement → gate →
  archive → PR → review → merge with no stop at all.
- Suggested behavior: promote the charter to a Baseline concept — a generated
  `docs/agents/autonomy-charter.md` scaffold with the decision points enumerated and left
  unanswered, so adopting the Baseline asks the maintainer these questions **once per project**
  instead of once per Spec. The clauses that repeat across fluxus and fiscus (QA verdict
  acceptance, non-convergence merge with disclosure obligation, corrective budget, bounded
  tooling) are strong candidates for the scaffold's defaults.

## What worked — keep

- **`settle` is the reason none of the three `SQLITE_BUSY` failures lost work.** Its contract —
  re-verify in the kept surface, commit, integrate — recovered complete Tasks in seconds each
  time. The recovery path is sound; what needs fixing is how often it is needed.
- **The authored terminal gate earns its cost.** Across three Specs it produced twelve findings
  that unit tests could not, including an architecture-boundary violation, a wrong hosted domain
  proved against a live backend, a measured contrast failure, and two false "unreachable"
  declarations written by the supervisor itself. A gate that refutes its own commissioner is the
  strongest evidence the design works.
- **`--detach` plus `events --follow` is the right unattended shape.** Zero manual polling
  across seventeen Runs; every terminal state arrived as an event, and the JSONL projection was
  stable enough to drive recovery decisions without opening the console log.
