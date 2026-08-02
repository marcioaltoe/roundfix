# Spec queue

The one dependency-and-risk-ordered queue of approved Specs that
`docs/agents/autonomous-work.md` requires the Supervisor to maintain. The
Supervisor works this list top to bottom. A Spec leaves the list when it
archives.

## Order

| # | Spec | Why here |
| --- | --- | --- |
| 1 | `0056-profiles-configure-merge-semantics` | In flight. Destructive defect: configuring one category deleted four others. |
| 2 | `0057-baseline-capability-evidence-and-retention` | Managed clauses disappear with an empty retention ledger in live consumer repositories. Unblocked by 0062. |
| 3 | `0059-run-storage-compaction-and-global-sanitation` | Operational: the Run Database grows with every cycle and its pages are never reclaimed. |
| 4 | `0064-spec-artifact-consistency-gate` | Highest leverage per unit of work. Half of this repository's QA findings were artifact contradictions catchable in seconds; every Spec after this one costs less. |
| 5 | `0063-qa-cycle-economics` | Largest throughput win. A cold gate costs 148s against 5s warm, and one stale assertion hides every behavioral finding for a whole cycle. |
| 6 | `0065-loop-order-and-verification-honesty` | The loop states two contradictory orders, and a Task can settle `completed` having done nothing. |
| 7 | `0066-run-teardown-reclaims-what-it-created` | Failed cycles leave Run Branches that block review Runs, and adapter children that outlive their Run by days. |
| 8 | `0067-derived-artifact-regeneration-boundary` | Smallest of the group. Recurs on every owned-skill edit, but has a known manual workaround. |

## Why this order

The first three were already approved and ordered before this group existed;
0057 and 0059 carry live defects in shipped behavior, which outranks improving
the loop that ships them.

Within the new group, 0064 comes before 0063 despite being smaller. It is the
cheaper of the two and pays back immediately: the contradictions it catches are
exactly what cost four gate cycles across Specs 0056 and 0058, and every Spec
implemented after it — including 0063 itself — is authored under its check.
0063 then removes the cost of the cycles that remain.

0065 follows because it changes how Task Graphs are authored and how the loop
sequences itself; landing it after 0064 means its own graph is authored under
the consistency check. 0066 is operational debris that slows the loop without
blocking correctness. 0067 is last: it recurs often but its workaround is one
documented command.

## Prerequisites the maintainer owns

- **0058 is merged and archived under override.** Its remaining QA row needs a
  real tagged release against six live npm trusted-publisher bindings, plus the
  repository variable `NPM_TRUSTED_PUBLISHING_FALLBACK=1`. Recorded in the
  release runbook.
- **0063, 0064, 0065, and 0067 each need tooling authorization with bounded
  files before decomposition.** None is authorized today; each Spec's Tooling
  authority row says so.
