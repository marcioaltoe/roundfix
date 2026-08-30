---
status: pending
type: qa
---

# Task: QA gate

Verify every deliverable of this Spec against the running commands.

## Work

- An author following any of the three pre-work authoring skills — write-prd,
  write-techspec, write-tasks — reaches the probing form of the check
- The terminal QA gate's own precondition stays non-probing, and the skill says
  why the difference exists. A gate that probed would report every completed
  Task vacuous and refuse every finished graph
- A clean verdict from the non-probing form states its own coverage on the
  verdict line, and the replaced trailing note is gone
- A clean verdict from the probing form states that the commands ran
- A verdict with findings is unchanged
- A QA Report names its auditing binary, distinct from the audited tree
- A precondition-refusal report names its auditor too
- A released-build identity produces a complete record rather than empty fields
- Staleness reports `current`, `stale`, and `unknown`, and `unknown` is never
  silently rendered as current
- A report written without the new keys still validates
- No verdict rule, row contract, or blocked-cause count changed
- The skill changes stay inside the recorded authorization's bounded files,
  checked from Git evidence, with the generated copies matching `make
  skills-sync` output
- The glossary check: whether this Spec introduced, changed, or retired a term
  the domain context should carry

## Outside evidence

One acceptance row rests on evidence this Spec did not author. Three freshly
decomposed Specs in `fluxus` passed `roundfix spec check <slug>` — the full
unscoped sweep — while carrying **eight** authored Verification commands that
could not fail: three naming a suite or file that already existed, five invoking
the repository verification gate. A repository this Spec did not build, measured
before this Spec existed, recorded at
`references/2026-08-14-a-clean-checker-left-eight-vacuous-commands-standing.md`
and provenanced in `references/_index.md`. The row records that this measurement
is what establishes the requirement, rather than a rehearsal of the Spec's own
premise.

The same document records eleven vacuous commands measured in **this**
repository while implementing Spec 0094. That number corroborates the class but
is not outside evidence, and the row must not cite it as such — an earlier draft
of this Task did, which is the attribution the gate caught as F-002.

## References

- All user stories and core features

## Verification
- `newest="$(ls -1 docs/specs/0116-a-verdict-that-states-its-own-scope/qa/qa-report-*.md 2>/dev/null | sort | tail -1)"; test -n "$newest" && grep -q "^verdict: pass" "$newest" && roundfix spec check 0116-a-verdict-that-states-its-own-scope --strict && go test -count=1 ./internal/app ./internal/spec ./internal/speccheck ./internal/cli`
