---
status: done
created_at: 2026-07-30
updated_at: 2026-09-08
absorbed_by: 0066-run-teardown-reclaims-what-it-created
---

# Failed QA runs accumulate Run Branches nothing can release (2026-07-30)

Spec 0053 gives `roundfix reconcile` a `superseded` classification so a Run Branch holding nothing but an obsolete QA report stops being preserved forever. Implementing that Spec produced the case its design does not cover: four QA runs, none passing, and four Run Branches that neither `reconcile --apply` nor Branch Integrity Preflight can clear.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-30-failed-qa-runs-accumulate-unreleasable-run-branches.md`.
