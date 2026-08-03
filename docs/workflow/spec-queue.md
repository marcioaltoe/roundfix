# Spec queue

The one dependency-and-risk-ordered queue of approved Specs that
`docs/agents/autonomous-work.md` requires the Supervisor to maintain. The
Supervisor works this list top to bottom. A Spec leaves the list when it
archives.

## Order

| # | Spec | Why here |
| --- | --- | --- |
| 1 | `0071-verification-cost` | Measured first. A third of Spec 0057's five hours was one whole-package command repeated in every Task, and the suite runs almost sequentially on twelve cores. Every Spec after this one is cheaper. |
| 2 | `0072-qa-is-a-task-not-a-flag` | Follows 0071. The gate is a flag on the invocation, so a graph that grows after it reported leaves no trace — three closings on 0057 read as three normal cycles. Authoring it into the graph changes every Spec after it. |
| 3 | `0074-git-spawn-economy` | Measured, like 0071: the suite is spawn-bound (36% utilization, kernel time 3× user, 13,926 git spawns per run), and ~6k of those are production-issued — every real Run pays them against the user's repository. Test-side waste is already removed; deleting the twelve heaviest tests was measured to buy only 7–14s. The maintainer's <60s full-fresh target lands here. |
| 4 | `0073-skill-versions-decoupled-from-the-binary` | Requested during the 0.0.3 release. Skill content is pinned by digest into the binary, so a skill edit changes what Roundfix claims — three failures in one session came from it. Compatibility becomes a declared minimum version instead, and the 0.3.1 release cut waits for it. |
| 5 | `0059-run-storage-compaction-and-global-sanitation` | Operational: the Run Database grows with every cycle and its pages are never reclaimed. |
| 6 | `0064-spec-artifact-consistency-gate` | Highest leverage per unit of work. Half of this repository's QA findings were artifact contradictions catchable in seconds; every Spec after this one costs less. |
| 7 | `0063-qa-cycle-economics` | Largest throughput win. A cold gate costs 148s against 5s warm, and one stale assertion hides every behavioral finding for a whole cycle. |
| 8 | `0065-loop-order-and-verification-honesty` | The loop states two contradictory orders, and a Task can settle `completed` having done nothing. |
| 9 | `0070-declared-unreachable-acceptance` | Follows 0065. A Spec whose acceptance no hermetic Verification can reach cannot archive without `qa_override`, which spends the mechanism reserved for failed evidence on a Spec that has none. |
| 10 | `0066-run-teardown-reclaims-what-it-created` | Failed cycles leave Run Branches that block review Runs, and adapter children that outlive their Run by days. |
| 11 | `0068-spec-close-audit` | Pairs with 0066. A Spec cycle leaves branches and worktrees nobody audits, and an unmerged Pull Request lets delivered work stay invisible on the default branch. |
| 12 | `0069-review-run-targets-its-pull-request` | Pairs with 0068. A Review Run resolves its branch from the checkout instead of from the Pull Request it names, and the mismatch check runs after the Review Source query rather than at Preflight. |
| 13 | `0067-derived-artifact-regeneration-boundary` | Smallest of the group. Recurs on every owned-skill edit, but has a known manual workaround. |

## Why this order

0071 is first because it was measured, not estimated: `go test ./internal/baseline
-count=1` costs 109s warm, every one of Spec 0057's fourteen Tasks carried one
such command, and adding `t.Parallel()` to one test's subtests took it from 29s
to 17s. Until that changes, every later Spec pays the same tax multiplied by its
Task count.

0056 and 0057 archived on 2026-08-02 and left this list. 0059 was approved and
ordered before this group existed and carries live defects in shipped
behavior, which outranks improving the loop that ships them.

0074 sits directly after 0072 by maintainer direction (2026-08-03): the
complete fresh suite must run under 60 seconds, and after the test-side
campaign the remaining floor is production-issued git spawns — work that also
speeds every real Run. The release cut (0.3.1) still waits for 0073.

Within the new group, 0064 comes before 0063 despite being smaller. It is the
cheaper of the two and pays back immediately: the contradictions it catches are
exactly what cost four gate cycles across Specs 0056 and 0058, and every Spec
implemented after it — including 0063 itself — is authored under its check.
0063 then removes the cost of the cycles that remain.

0065 follows because it changes how Task Graphs are authored and how the loop
sequences itself; landing it after 0064 means its own graph is authored under
the consistency check. 0070 sits immediately after it because both concern what
the loop is allowed to call finished: 0065 fixes a Task settling `completed`
without doing the work, 0070 fixes a Spec that did all reachable work and still
cannot close. 0066, 0068, and 0069 are adjacent and run together: both close
the gap between what a Run creates and what survives it, one for process trees
and Run Branches, 0068 for the Supervisor's own branches and worktrees, and
0069 for a Review Run acting on the wrong branch entirely. 0067 is last: it recurs often but its workaround is one
documented command.

## Prerequisites the maintainer owns

- **0058 is merged and archived under override.** Its remaining QA row needs a
  real tagged release against six live npm trusted-publisher bindings, plus the
  repository variable `NPM_TRUSTED_PUBLISHING_FALLBACK=1`. Recorded in the
  release runbook.
- **0058 archived under override on 2026-08-02** (PR #65); its release prerequisites below remain open.
- **0063, 0064, 0065, and 0067 each need tooling authorization with bounded
  files before decomposition.** None is authorized today; each Spec's Tooling
  authority row says so.
