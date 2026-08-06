---
status: done
created_at: 2026-08-04
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-run-lifecycle-and-branch-integrity.md
---

# 2026-08-04 — Watch derives a review head it never checks is reachable

## What was observed

`roundfix watch --source coderabbit --pr 78 --until-clean` started a Run,
reached `WaitingForReview`, and died with an untyped API error:

```
Review Source status: pending; phase=WaitingForReview;
expected_head=3738038ce1e4e10200ac7c1f3b5d334af1ee87ef; evidence_kind=none
Watch Run run_20260804T142044Z_03cddc1461cfbf6e reached Failed after 0 Round(s).
Review Issues: unknown — fetch did not complete.
roundfix: watch failed after Run start: fetch CodeRabbit check runs: gh command failed with HTTP 422
```

The cause was one line up in its own startup report, already printed and
already correct:

```
Git: clean, 1 unpushed commit(s), upstream origin/ma/0012-reconciliation-assurance
```

The expected head was a local commit. GitHub had never seen that SHA, so
asking for its check runs is a request for a resource that cannot exist, and
the API answered `422`. Pushing the commit and rerunning the identical command
worked immediately.

The Run had already been created, so the failure consumed a Run id, a Run
Event Journal, and an artifact directory to discover something knowable before
any of them existed.

## Root cause

Branch Integrity Preflight validates plenty about the checkout — clean tracked
tree, pending Run Branch work, competing Active Runs — but not the one
property the review flow depends on first: **the head it will ask the Review
Source about must be reachable by the Review Source.**

Roundfix even computes the evidence. `1 unpushed commit(s)` is in the startup
report, printed before the failure, derived from the same upstream comparison a
preflight check would use. The datum is present and unused.

The failure then presents as an infrastructure error rather than a
precondition. `HTTP 422` names the transport, not the mistake, so the operator
has to reason backwards from an opaque status code to "the commit is local".
For an autonomous Supervisor this is worse than for a person: a typed
precondition failure is a next action, while a `422` after Run start is an
incident to investigate.

## What would settle it

Add a review-head reachability row to Branch Integrity Preflight, before Run
creation: if the resolved expected head is not present on the configured
remote, refuse with exit `2` and the actionable message, in the same shape as
the existing dirty-tree refusal.

```
review Run Preflight Validation refused because the expected head
<sha> is not on <remote>; 1 unpushed commit(s) on <branch>.
Next action: push the branch, then rerun.
```

Two smaller improvements stand on their own:

- **Classify the `422`.** A Review Source error naming a head the remote does
  not have should surface as an unreachable-head reason rather than a generic
  `gh command failed`, so the diagnosis survives even if the preflight row is
  never added.
- **Do not spend a Run on a precondition.** The refusal belongs before Run
  creation, exactly like the dirty-tree case, which already gets this right and
  reports `No side effects`.

## Related

[[2026-08-04-review-runs-halt-autonomous-delivery-on-unrelated-dirty-files]] is
the neighbouring row in the same preflight: one requirement is enforced too
broadly, and this one is not enforced at all.

## Spec pointer

None yet.
