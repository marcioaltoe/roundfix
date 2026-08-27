---
status: done
created_at: 2026-08-12
updated_at: 2026-08-27
kind: finding
---

# Autonomous delivery — five Unresolved Runs to deliver one Spec (2026-08-12)

An autonomous session in `fiscus` on 2026-08-07/08 authored Spec
`0004-entrada-de-catalogo` (PRD, TechSpec, four ADRs, a 12-Task graph) and ran it
in **five terminal Runs, all `Unresolved`**, totalling **5h21m**, plus three
corrective Tasks after the QA gate failed with three Blocks-Completion findings.
Minted here from the Inbox Entry
`inbox/roundfix/2026-08-08-atrito-medido-na-entrega-da-0004-do-fiscus.md` in the
Secondbrain. The Supervisor's own authoring error was large and is named in that
repository's finding; what follows is the tool friction.

## 1. Corroboration — `spec check --run-verification`

- Symptom / evidence: the 0004 graph was born with twelve Tasks whose
  Verifications are mostly `bunx vitest run <file-that-does-not-exist-yet>`. Had
  the runner exited `0` without matching a file, **every** one of those commands
  would have passed on the pre-work tree and the whole graph would have been
  vacuous. It was measured deliberately: the runner exits `1` on backend and
  frontend, so the graph was saved by the tool's luck rather than by design. The
  first measurement even reported `exit=0` because `$?` was capturing the `tail`
  of a pipe; it only surfaced on a redo without the pipe.
- Root cause: the same gap the `oraculum` queue names as its largest lever —
  the form of a Verification is checked, its execution is not.
- Action / suggestion: nothing new beyond that proposal. This is a second
  independent case, recorded to add priority rather than to repeat the design.

## 2. The citation parser does not read a Portuguese list, nor an ADR without its prefix

- Symptom / evidence: `SC-COVERAGE-UNTASKED` reported Core Features 6 and 12 as
  referenced by no Task. The Tasks cited them, as `Core Features 4, 5 e 6` and
  `Core Features 8, 9, 10 e 12`. Replacing the `e` with a comma cleared both
  errors **without changing a word of meaning**. In the same report,
  `SC-ADR-UNLISTED` reported six ADRs that the PRD listed as `0026 (…)` and
  `0029/0030/0031`: recognition requires the `ADR-NNNN` token, and that appears
  neither in the failure message nor in the command's documentation.
- Root cause: the parser's accepted forms are undocumented and narrower than the
  authoring language.
- Action / suggestion: accept a conjunction (`e`, `and`) as a separator and an
  ADR number without its prefix when the context is the obligations line; at
  minimum, state in the message which form the parser recognises.

## 3. A missing `worktree.copy` source fails only after the Run exists

- Symptom / evidence: the `fiscus` `.roundfixrc.yml` declares
  `worktree.copy: [packages/backend/.env]` and bootstrap copies that file. **It
  did not exist in the checkout**, and the development Postgres had been down for
  two days. Had the Run been dispatched, bootstrap would have failed in **every**
  Task Worktree and the Run would have died without producing a line of code. It
  was found by manual inspection; nothing warned.
- Root cause: `worktree.copy` entries are declared repository-relative paths —
  verifiable in Preflight without executing anything — and are only read at
  bootstrap.
- Action / suggestion: validate each source's existence in Preflight and refuse
  naming the absent ones, instead of failing worktree by worktree after the Run
  exists.

## 4. An identical repeated failure of the same Work Item is not signalled

- Symptom / evidence: `task_07` failed in **two consecutive Runs with the same
  assertion** — a `Date` compared against the serialization the driver returns
  for `timestamptz` (`"2026-08-08 16:02:00+00"`), with the persisted value
  correct both times. Each Run reported the failure as new; it was only noticed
  by reading the two diagnostic files side by side, after spending a whole Run
  reproducing a known diagnostic.
- Root cause: nothing compares a failure's verdict and diagnostic signature
  against a prior failure of the same Work Item.
- Action / suggestion: when verdict and diagnostic signature match a prior
  failure of the same Work Item, say so in the report and in the `task-status`
  Run Event. The repetition is the signal that the task file must change, not the
  Agent. What unblocked this one was adding a requirement to the task file, not a
  third attempt.

## 5. The ceiling of two corrective Tasks does not name the exits when it is reached

- Symptom / evidence: the gate returned `fail` with three Blocks-Completion
  findings, and **two of them were gaps in the contract the Supervisor
  authored** — dispatch recovery and per-row refusal-reason projection — rather
  than implementation slips. The contract says more than two corrective Tasks
  means the decomposition was wrong, and does not say what to do then. The
  session stopped and asked the maintainer: a policy stop.
- Root cause: a ceiling with no sanctioned exit converts an authoring judgement
  into a human interruption.
- Action / suggestion: name the sanctioned exits — amend the TechSpec and redo
  the cut from it, or promote the excess to its own Spec with the gate explicitly
  failing the discovered story. Written down, the excess becomes a decision
  inside the loop's authority.

## 6. Authoring asks for tooling authority by path, and authorization arrives in rounds

- Symptom / evidence: the Spec needed **three express authorizations** from the
  maintainer: the backend set at authoring, the frontend set when the compliance
  Task revealed a missing accessibility library, and
  `packages/backend/.env.example` after the gate failed.
- Root cause: the impact analysis was by path ("does this Spec touch CI?") rather
  than by **class**. The predictable classes are: dependency per package, local
  composition, CI workflow per package, test-runner configuration per package,
  and environment example per package. Asking by path finds what is already
  imagined; asking by class finds what is forgotten.
- Action / suggestion: have `write-techspec` walk the classes and ask for the set
  once, naming the classes the Spec will not touch.

## 7. Gaps in the autonomous-work guide: environment preparation and repeated failure

- Symptom / evidence: `docs/agents/autonomous-work.md` covers authority and the
  end of the cycle, and leaves both operational ends out. Nobody owns environment
  preparation — the missing `.env` and the stopped Postgres were found by habit,
  not by instruction, and the runbook that resolves it lives in the README, which
  the guide does not reference. Nothing says what to do when the same Work Item
  fails identically twice; that amending the task file is Supervisor work was
  inferred from the rule "do not write feature code".
- Root cause: two operational responsibilities with no written owner.
- Action / suggestion: two mandatory lines in the Baseline module that generates
  the guide.

## 8. Observation without an established cause — a Run above the configured budget

- Symptom / evidence: `budget.max_run_duration` is `2h` in the `fiscus`
  `.roundfixrc.yml`, and Run `run_20260808T021810Z_aafc8a06587cbc78` lasted
  **2h34m** according to `roundfix runs list --state all`. The other four stayed
  below (37m, 50m, 20m, 1h0m), so there is no second case to compare.
- Root cause: `unknown`. It was not investigated whether the budget is evaluated
  at Work Item boundaries — which would explain exceeding it by one long Task —
  or whether it was not applied at all.
- Action / suggestion: if it can be exceeded by a Work Item in flight, say so in
  the budget's documentation: whoever configures `2h` is planning the night's
  window.

## What worked — keep

- **A tooling Task first, with bounded files.** `task_01` ran before everything
  else and its commit touched **exactly** the five authorized files plus its own
  task file, verified with `git show --name-only`.
- **The Agent investigated infrastructure instead of loosening a test.** A
  timeout on a PUT against MinIO was traced to an unhealthy container with a
  stale volume from a previous Run; it was recreated from the declared compose
  service and the proof redone (`putStatus: 200`, 50 bytes back), with an explicit
  record that no timeout, retry or assertion was changed.
- **Verification Feedback resolved on the first attempt.** Several failed
  attempt-1 runs were repaired inside the same Agent Session without consuming a
  Run.
- **The gate found what 196 green tests did not** — including a frontend test
  that accepted a generic sentence as proof of "why", encoding the regression in
  its own expectation.

Adopted by Spec `0118-a-task-proved-once-does-not-run-twice`.
