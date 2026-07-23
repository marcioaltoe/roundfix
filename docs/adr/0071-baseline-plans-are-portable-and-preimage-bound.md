# ADR-0071: Baseline Plans are portable and preimage-bound

Status: Accepted

A `roundfix/baseline-plan/v1` document is a self-contained, path-relative
artifact that carries the exact planned postimages, canonical managed-entry
ledger, repository identity, catalog identity, and every bounded preimage that
influenced or can be changed by the plan. `roundfix baseline apply` can consume
the artifact in another clone of the same Git repository only when the stable
repository identity, plan digest, and all bounded preimages match. Hidden
pending-plan state, absolute checkout paths, whole-worktree invalidation, and
applying a plan solely because target files match were rejected because they
either break automation portability or weaken approval scope.
