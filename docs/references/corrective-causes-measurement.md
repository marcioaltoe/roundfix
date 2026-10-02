# Corrective cause measurement

Window: 2026-09-17 through 2026-10-02 UTC, inclusive at the recorded date
boundaries. The read-only `roundfix runs causes` report contains 177 terminal
Runs and 132 failed Verification or corrective items. Its table digest is
`53e5e5a764109b92806c6efe4d88603c13e5de2c6d21816755406deb11b6dd5f`.

## Counts

| class | count |
| --- | ---: |
| scope-or-authorization | 3 |
| shared-section-contract | 5 |
| repository-convention | 15 |
| implementation-defect | 9 |
| environment | 0 |
| unclassified | 100 |
| total | 132 |

Repository knowledge is the first three classes: 23 items, or 17.42% of all
items and 71.88% of the 32 classified items. The unclassified share is 75.76%.

The trigger counts are derived from each item’s associated Task text using the
fixed trigger table. They are: `pre-pr-review` 40, `qa-gate` 37,
`verification` 26, and `unknown` 24. Five items could not be assigned a Task
text because their Task files were unavailable in the active or archived Spec
roots; these are included in the total and are not silently assigned a trigger.

The five most frequent classified signatures were `record-or-golden` (13),
`review-defect` (9), `shared-section` (5), `authorization-or-scope` (3), and
`help-or-guide` (2). There are only five observed signatures in this record.

Because unclassified items exceed one third of the record, the required
unclassified-check list is shown instead of pretending that unmatched items
have a signature:

| unclassified check | count |
| --- | ---: |
| `make verify-changed` | 24 |
| `pre-pr-review` | 18 |
| `qa-gate` | 10 |
| `grep -rq "verify-incremental" docs/agents/` | 4 |
| `settlement check: spec consistency` | 2 |
| Archive eligibility test command | 2 |
| Project Config authorization test command | 2 |
| Spec judge and skills contract test command | 2 |
| Profiles recommendation contract test command | 2 |
| all other unmatched checks | 34 |

## Samples

These are five representative rows per classified class (or every row when a
class has fewer than five), showing the record’s signature.

| class | Spec / Task | kind | signature |
| --- | --- | --- | --- |
| scope-or-authorization | 0170-authoring-that-fails-before-dispatch / task_06 | corrective | authorization-or-scope |
| scope-or-authorization | 0181-gates-that-refuse-only-what-someone-can-act-on / task_07 | corrective | authorization-or-scope |
| scope-or-authorization | 0187-a-queue-that-recovers-without-a-supervisor / task_06 | corrective | authorization-or-scope |
| shared-section-contract | 0177-runs-that-fit-their-budget-and-park-honestly / task_07 | corrective | shared-section |
| shared-section-contract | 0181-gates-that-refuse-only-what-someone-can-act-on / task_07 | verification | shared-section |
| shared-section-contract | 0181-gates-that-refuse-only-what-someone-can-act-on / task_07 | verification | shared-section |
| shared-section-contract | 0188-a-grant-that-authorizes-delivery-without-governed-paths / task_01 | verification | shared-section |
| shared-section-contract | 0192-owned-skills-that-describe-the-product-as-it-is / task_03 | verification | shared-section |
| repository-convention | 0169-override-and-ordinal-follow-ups / task_04 | corrective | help-or-guide |
| repository-convention | 0170-authoring-that-fails-before-dispatch / task_02 | verification | record-or-golden |
| repository-convention | 0170-authoring-that-fails-before-dispatch / task_02 | verification | record-or-golden |
| repository-convention | 0172-a-qa-gate-that-tells-the-truth / task_07 | corrective | record-or-golden |
| repository-convention | 0174-operator-surfaces-that-tell-the-truth / task_01 | verification | record-or-golden |
| implementation-defect | 0149-a-supported-way-to-reopen-a-settled-gate / task_05 | corrective | review-defect |
| implementation-defect | 0149-a-supported-way-to-reopen-a-settled-gate / task_06 | corrective | review-defect |
| implementation-defect | 0151-a-supersession-the-archive-can-see / task_05 | corrective | review-defect |
| implementation-defect | 0151-a-supersession-the-archive-can-see / task_06 | corrective | review-defect |
| implementation-defect | 0153-a-reviewer-the-workflow-runs / task_06 | corrective | review-defect |
| environment | none observed | — | — |

## Decision

The fixed `_techspec.md` → Decision rules, rule 1, says `inconclusive` when
there are fewer than 30 items or when unclassified items exceed one third.
Here there are 132 items, but 100 are unclassified, so the rule yields
`inconclusive`. The published `_prd.md` → Acceptance evidence says that
repository-context evidence points away from a memory layer: one study found
context files increased cost and reduced success, and a 288-run ablation
attributed failures to implementation skill rather than missing repository
knowledge. That published evidence is comparison context only; it does not
override the predeclared rule or the measured counts.

Verdict: inconclusive
