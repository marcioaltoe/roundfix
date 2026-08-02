---
status: pending
created_at: 2026-07-29
updated_at: 2026-07-29
---

# QA cycles — the gate serializes findings and sits far from the defects it must catch (2026-07-29)

Spec 0052 needed four QA cycles at roughly twenty minutes each, and a parallel
Vortex session needed seven on one spec. Neither number came from bad
implementation: in both repositories the Agents built what the Spec asked for.
The cost came from two structural properties of the gate — it stops at the
first static failure, so one stale assertion hides every flow finding and buys
a whole cycle; and it is the only detector that exercises real external
surfaces, so a defect a thirty-second probe would expose waits for a
twenty-minute gate. This report separates what Roundfix must change from what
the supervisor must change, because only the first belongs in a Spec.

## 1. A static-gate failure costs a full cycle and hides every flow finding

- **Symptom / evidence**: Spec 0052's first QA cycle reported
  `rtk make verify: failed — 2,827 passed, 6 failed, 2 skipped` and closed with
  "Public-flow QA stopped at the code-caused static failure, as required by the
  gate", recording 2 passing rows, 1 failing static gate, and 8 blocked flow
  rows. The single reproducible failure was one assertion pinning the literal
  `"1.1.4"` while the Spec legitimately raised that pin — a one-line stale
  expectation. The next cycle, with that line fixed, found four real product
  and documentation defects. Those four were discoverable during the first
  cycle; the gate simply never looked.
- **Root cause**: the gate treats the repository-wide Verification as an
  all-or-nothing precondition for the flow matrix. The repository's gate runs
  the production build *after* the test suite, so a failing test also prevents
  the binary from being built, and the flow rows — which drive the built
  binary — become unreachable. The coupling is real but incidental: the
  product compiled fine, and a separately built binary could have executed
  every flow row.
- **Action / suggestion**: when the static gate fails but the product still
  builds, build the binary independently and execute the flow matrix anyway,
  reporting the static failure as its own finding alongside whatever the flows
  find. Keep the current stop only for a failure that genuinely prevents
  execution, such as a compile error. One cycle should surface everything a
  build can surface.

## 2. The only detector for external contracts is the slowest one

- **Symptom / evidence**: Spec 0052's `Blocks-Completion` defect was that the
  official Claude adapter echoes the requested model back into its advertised
  list, so requesting `opus` returns both `opus[1m]` and `opus`, which collided
  under the Spec's new opaque parsing and failed the built-in frontend tuple on
  every real Run. Every Task's Verification was `go build` plus `go test`, and
  every capability fixture in the suite modelled an advertised list that no
  adapter returns. The defect could only appear against the installed adapter.
  A `roundfix profiles validate --category frontend --json` probe answers the
  same question in about thirty seconds; the QA gate took about twenty minutes
  to reach it.
- **Root cause**: Task Verification commands in this repository prove that the
  code compiles and that unit tests pass. Nothing in the authoring contract
  asks whether a Task changes a contract with an external surface — an ACP
  adapter, a Review Source, a registry — or requires a probe against that
  surface when it does. The detector therefore lands at the end of the Spec
  instead of inside the Task that introduced the risk.
- **Action / suggestion**: this is primarily an authoring discipline, recorded
  in [the loop instructions](../workflow/spec-implementation-loop.md). The
  product-side half worth considering: have `write-tasks` require an
  external-surface probe in the Verification whenever a Task's declared scope
  touches an adapter, Review Source, or registry contract, so the requirement
  is a contract rather than a habit.

## 3. Corroboration from a second repository

A Vortex session on 2026-07-29 reported seven QA cycles on one Spec, and its
own analysis matches this one on both points: the detector's latency did not
match the defect class (operational defects only visible to a container-backed
gate), and the corrective cycles were serial because each cycle surfaced one
defect at a time. Its third cause — a test seam that never reached the
production driver branch — is the same shape as the fixture gap in finding 2:
the suite exercised a path the product does not take in production.

Its remaining causes are supervisor-side and belong to the loop instructions,
not to Roundfix: specifications that named a destination without asking who
reads it, corrective work that introduced the next defect, and narration
overhead.

## 4. Concurrency is limited by a stale note, not by the design

- **Symptom / evidence**: this repository pins `worktree.concurrency: 1` with
  the comment "Raise to 2 once integration performs a 3-way merge instead of a
  plain cherry-pick". Spec 0052 ran eight Tasks strictly serially because of
  it.
- **Root cause**: the comment contradicts
  [ADR-0026](../adr/0026-task-integration-is-a-serialized-cherry-pick-queue.md),
  which states that the serialized cherry-pick is deliberate: "Tasks that
  write-tasks declared independent are file-disjoint by contract, making a
  conflict a real graph defect to surface, not noise to auto-resolve." Under
  that design a cherry-pick conflict is the intended signal that the graph
  declared two Tasks independent when they were not, and the Task settles
  `failed` with its Worktree kept for inspection — visible and recoverable,
  never silent corruption. No 3-way merge is required or wanted.
- **Action / suggestion**: raising concurrency is safe exactly when the Task
  Graph is genuinely file-disjoint, which is the graph author's
  responsibility. Correct the stale comment when a bounded authorization
  covers that file, and record the real rule in the loop instructions: raise
  concurrency, and treat any cherry-pick conflict as a defect in the
  decomposition rather than a reason to serialize everything forever.

## What worked — keep

- The QA gate found a defect that no fixture in the repository would ever have
  caught, because it insisted on driving the real adapter. That insistence is
  the reason the Spec did not ship with the frontend broken on every fresh
  installation; nothing here should weaken it.
- The gate audited the supervisor with the same rigor as the Agents, refusing a
  recovery commit that touched Project Config outside the Spec's bounded
  authorization. One of the four cycles was that refusal, and it was correct.

## Routing — 2026-08-01

Routed to [Spec 0063](../specs/0063-qa-cycle-economics/_prd.md) on 2026-08-01.
