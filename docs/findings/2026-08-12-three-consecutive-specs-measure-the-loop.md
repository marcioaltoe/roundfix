---
status: pending
created_at: 2026-08-12
updated_at: 2026-08-12
kind: finding
---

# Autonomous loop — three consecutive Specs measure the friction (2026-08-12)

An autonomous session in `vortex` on 2026-08-07/08 delivered Specs 0023, 0024
and 0025 through the Roundfix loop and recorded the friction with numbers,
separating tool defects from the Supervisor's own authoring errors. Minted here
from the Inbox Entry
`inbox/roundfix/2026-08-08-atrito-medido-do-loop-autonomo-em-tres-specs.md` in
the Secondbrain, whose provenance is that session.

The framing number: fourteen QA gate executions for three Specs, eleven of them
without `pass`. **None of the eleven failed on wrong business logic** — every
one failed on a contract: boundary, vocabulary, Verification or row
classification. The same gate found four real defects no suite would have
caught, so any proposal to loosen it must be read against that.

| | 0023 | 0024 | 0025 |
| --- | --- | --- | --- |
| QA gate executions | 2 | 4 | 8 |
| Verdicts without `pass` before the `pass` | 1 | 3 | 7 |

## 1. The Pull Request row costs reruns in every Spec

- Symptom / evidence: the authored gate is the terminal Task of the graph and
  runs before a Pull Request exists (ADR-0091), so the row covering PR checks,
  approval and threads is unreachable by design in every Spec. Behaviour was
  inconsistent across Specs identical in that respect: 0023 and 0024 passed on
  equivalent evidence; 0025 returned `partial`/`fail` until the task file
  invoked the ADR-0080 path explicitly. This one row alone accounted for six of
  the eight gate executions of Spec 0025.
- Root cause: the equivalent-evidence path exists, and the route to it is not
  discoverable from inside a task file. Two attempts to settle it from the task
  file — reclassifying the row, and changing its effect on the verdict — were
  correctly refused, because both rules belong to the Run prompt and the
  `qa-gate` skill.
- Action / suggestion: apply the equivalent-evidence path to the Pull Request
  row by default in `qa-gate`, requiring the evidence to be recorded (head SHA,
  ancestry, changed files, local checks). Minimal alternative: the QA task
  template is born carrying the clause.

## 2. The resolve Agent rewrites `## Verification` and introduces environment dependence

- Symptom / evidence: review rounds rewrote `## Verification` blocks of a Spec
  that was not even being implemented — its documents were in the Pull Request
  through a wide commit scope — and introduced
  `test -n "${VORTEX_POSTGRES_INTEGRATION_URL:-}" && <tests>` and
  `if (!process.env.TASK_04_BASELINE) process.exit(1)`. Neither variable exists
  in the repository, so both fail on the first line. Two Tasks failed with the
  work complete and settled through `settle`.
- Root cause: `## Verification` is Supervisor authorship and nothing protects it
  from the review-resolution Agent. The review intent was right in both cases —
  the original assertion was too broad — and the form was wrong.
- Action / suggestion: make `## Verification` immutable for the review
  resolution Agent, or at minimum surface an edit to it as a contract change
  highlighted in the Round report.

## 3. `spec check` does not catch a non-hermetic Verification

- Symptom / evidence: `SC-VERIFY-WORK-INDEPENDENT` exists; the sibling rule does
  not. Three forms appeared in this session and all three produced a false red:
  a Verification referencing an environment variable the repository does not
  declare, one chaining `test -n "$VAR" &&` in front, and one depending on a
  temporary directory or a tree snapshot outside the repository.
- Root cause: the detector set checks work-independence, not hermeticity.
  Separately, `SC-REF-UNRESOLVED` does not distinguish a file the Task
  **creates**: `## Context` accepts only `interface:` and `instruction:`, both
  requiring an existing path, so a Task cannot declare its own output. It hit
  four task files.
- Action / suggestion: add `SC-VERIFY-NON-HERMETIC`, and accept a `creates:`
  entry in `## Context` that is declarative and not checked for existence.

## 4. `include_nitpicks` and the branch ruleset have no coherence check

- Symptom / evidence: a Pull Request stuck at `BLOCKED` with green CI,
  `mergeable: MERGEABLE`, zero actionable findings and an empty
  `reviewDecision`. No isolated signal pointed at the problem. The `main`
  ruleset required `required_review_thread_resolution: true` while
  `.roundfixrc.yml` carried `include_nitpicks: false`, so nitpick threads were
  never fetched, never resolved, and the merge never unblocked.
- Root cause: two configurations contradicting each other in files nobody reads
  together. An auxiliary diagnostic was missing too: the classic
  `branches/<b>/protection` API returns **404 Branch not protected** when the
  protection is a ruleset, which reads as "there is no protection".
- Action / suggestion: Roundfix already has a coherence Preflight for
  `pushTriggersReview` × `review_source.request_review`. Add the analogous one
  and refuse naming both settings and the correction.

## 5. A transient `gh` failure ends the Run as `Failed`

- Symptom / evidence: `watch --until-clean` ended `Failed` with
  `fetch CodeRabbit commit statuses: gh command failed`, and `gh` worked in the
  same shell seconds later. Same family: one issue was marked `failed` with
  `Review Source propagation failed during resolve: … gh command failed` — the
  fix had been applied; only the call marking the thread resolved failed.
- Root cause: the retry taxonomy covers deadline, DNS, reset, `429` and `5xx`,
  and not a non-zero exit from the `gh` CLI.
- Action / suggestion: classify a `gh` failure by stderr rather than by exit code
  alone, and separate "fix applied" from "thread propagated", with a retry queue
  for the propagation.

## 6. Branch Integrity Preflight prescribes the wrong recovery for a superseded branch

- Symptom / evidence: a Run Branch from a stopped Run carried deliberately
  rejected work — a destructive `pgEnum` derivation removing values and
  reordering the rest with no migration. Preflight refused the commands and
  suggested `git merge --ff-only roundfix/run-<id>`, which would have regressed
  the repository. `reconcile` refused to release it:
  `Run Branch does not carry QA Report evidence`. After the Spec merged, the
  target branch stopped existing and the classification stayed `unintegrated`.
  Cost: `--skip-branch-integrity` on every `fetch`/`resolve`/`watch` of the
  session.
- Root cause: classification proves integration by ancestry against a target
  branch that may no longer exist, and offers no way to record supersession.
- Action / suggestion: when a terminal Run's target branch no longer exists
  (merged and deleted), classify the Run Branch as releasable; or offer a way to
  record supersession with a reason for `reconcile` to act on. Today the only
  exit disables the whole guard, including for legitimately pending work.

## 7. `settle` reads the task file from the surface, not from the checkout

- Symptom / evidence: fixing an unsatisfiable Verification required editing the
  file in **two** places — checkout and Task Worktree — because `settle` reads
  from the selected surface. Discovered by trial.
- Root cause: `settle` prints `Settle surface:` and does not say the task file
  comes from there.
- Action / suggestion: name the surface the task file was read from in the same
  output.

## 8. Authoring rules the baseline guides should carry

- Symptom / evidence: five rules, each derived from a measured failure in this
  session, belong in guides that are Roundfix assets (`autonomous-work.md`,
  `backend.md`, `domain.md`, `agent-instructions.md`), so they cannot be added
  in the consuming repository without being overwritten by the next
  `baseline update`.
- Root cause: the rules were learned downstream and have no upstream owner.
- Action / suggestion: carry them in the Baseline modules —
  (1) a Verification is hermetic: no environment variable, temporary directory
  or state outside the repository;
  (2) a requirement describes the property, not the magnitude of the data — one
  requirement asked for cleanup proven with a set above the Postgres parameter
  ceiling, 65,536 rows per execution, unstable under suite load and unable to
  distinguish 2 batches from 200, where the real property proves in three lines;
  (3) a use case test is born against ports and doubles, a persistence proof is
  born in infra — asking for both in one file manufactures an architectural
  boundary violation;
  (4) every Task that changes an already recorded contract updates the record in
  the same slice — the incomplete version, "a Spec that *creates* vocabulary
  updates the canonical document", produced a recurring finding;
  (5) commit scope per Spec — a wide `git add` carried artifacts of two future
  Specs into another's Pull Request; twelve of the sixteen majors of that Pull
  Request were about documents of Specs that were not being implemented, and
  every correction re-triggered review over them. Direct measurement: **6 rounds
  and 56 findings** in the unscoped Pull Request against **1 round and 18
  findings** in the scoped one that followed.

## 9. A "Heavy lift" finding induces the Agent to build test infrastructure that does not exist

- Symptom / evidence: observed three times in the same session, and the third
  was dangerous. Facing a finding that presupposes infrastructure the repository
  does not have, the resolve Agent **creates the dependency** instead of
  refusing: `VORTEX_POSTGRES_INTEGRATION_URL`, `TASK_04_BASELINE`, and a MySQL
  integration test running `DROP TABLE IF EXISTS` against **real ERP table
  names** (`fa7`, `fa9`, `log008fin2512`) in a database named `sglinx`, guarded
  only by "the variable is set". Pointing that variable at a client mirror would
  drop real fiscal tables. The reviewer caught the danger in the next round and
  the Agent hardened it with four guards, so the loop self-corrected — after
  passing through a state where the Pull Request contained destructive DDL
  against client table names.
- Root cause: choosing a test substrate is Spec scope, not a Pull Request
  correction round, and nothing says so at the point of resolution.
- Action / suggestion: a review finding that requires test infrastructure absent
  from the repository should not be resolvable by the review Agent — classify it
  as backlog instead of resolving. Common signal across all three: **the Agent
  introduced a new environment variable**, which is cheap to detect and serves as
  a trigger for human review or automatic refusal.

## What worked — keep

The gate found four real defects that no suite would have caught: a null
monetary value becoming a zero-value purchase; a version guard over a column
that is not a version; a Spec invariant **impossible** as specified, with an
executable repro; and a concatenated coupon key merging distinct tuples.

## Suggested priority by measured cost

1. Finding 1 — the Pull Request row in the gate.
2. Finding 4 — `include_nitpicks` × ruleset coherence.
3. Finding 8 — authoring rules in the baseline guides.
4. Findings 2 and 3 — protect `## Verification` and catch non-hermeticity.
5. Finding 6 — superseded Run Branch classification.
6. Finding 9 — the destructive DDL case justifies priority above its frequency.

## Evidence

In the `vortex` repository, under `docs/specs/_archived/0023-sync-purchase-invoices`,
`docs/specs/_archived/0024-sync-catalogo-de-produto` and
`docs/specs/0025-sync-derivada-horaria-e-cupons`: the reports under `qa/`
(fourteen executions) and the review artifacts under `docs/specs/_reviews/pr-138`,
`pr-139` and `pr-140`.
