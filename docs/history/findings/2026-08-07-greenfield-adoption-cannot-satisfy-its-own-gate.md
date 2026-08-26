---
status: done
absorbed_by: 2026-08-06-rollup-baseline-and-derived-tooling.md
created_at: 2026-08-07
updated_at: 2026-08-26
kind: finding
---

# Greenfield adoption cannot satisfy its own gate (2026-08-07)

A repository with stale managed carriers cannot complete adoption in Greenfield
mode through the interactive command. The run reaches profile alignment, then
refuses with a next action the Greenfield path can never supply.

## What was observed

Reported from a repository whose adoption reached `alignment: ready` with zero
divergences and then failed on "review and edit every proposed root-rule
classification". Three code facts produce it:

- `internal/cli/baseline_human.go` — `promptBaselineClassification` returns
  immediately for any mode other than Preservation, handing back a
  `RootPreservationRequest` whose `Decisions` and `SourceBaseline` are nil. The
  interactive workflow therefore never offers classification in Greenfield.
- `internal/baseline/preservation.go` — the Greenfield early return is guarded
  by `request.Mode == PreservationModeGreenfield && len(staleManagedSources) == 0`.
  With stale managed carriers present, planning continues past it.
- The Source Baseline is then non-empty (the stale managed sources became
  sources), so the empty-entries return does not fire either. With
  `request.Decisions == nil`, planning sets `PreservationStateActionRequired`
  and a decision skeleton, which `BuildPlan` turns into an action outcome.

The gate asks for classification input that the Greenfield interactive path is
structured never to collect.

## Why it matters

Any repository adopted under an older catalog has stale managed carriers by
construction, which is precisely the population most likely to re-run adoption.
For those repositories Greenfield is not a slower route — it is unreachable, and
the refusal names a review step the maintainer was never offered. The only way
through is Preservation, which costs the full classification pass.

Confirmed on `roundfix` 0.3.1.

## Route

Not fixed here. Spec 0082 removes classification from the *update* path but
explicitly leaves first adoption unchanged, so it does not reach this. Fixing it
needs a decision first: whether Greenfield should account for stale managed
carriers without classifying them, or whether the mode should refuse up front
with a next action that names Preservation instead of failing at a gate it
cannot pass.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
