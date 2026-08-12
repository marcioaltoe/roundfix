---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-12T14:16:46Z
updated_at: 2026-08-12T14:16:46Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A Review retires on local reachability, and never writes into history

A Review Artifact for a Pull Request with no owning Spec has no terminal
disposition today, and the accurate liveness signal — the provider's own Pull
Request state — cannot be a requirement, because reading it needs authentication
and a network the Baseline migration must run without. An orphan Review Artifact
therefore retires on local Git reachability of the head its round recorded: an
ancestor of the default branch retires as merged, an unreachable head whose branch
is gone retires as abandoned, and a reachable non-ancestor head stays live.
Reachability can lag the provider — a merged Pull Request whose commits were
squashed leaves no ancestor — so an undecidable head stays live rather than
retiring, since keeping a live artifact costs a directory and retiring a live one
breaks the Round the loop is still running. Independently, the Review Artifact root
resolver never resolves into the history root: a new Round for an already-retired
Spec's Pull Request writes to the live root and the migration may retire it later,
because a tool that writes into its own history makes history mutable.
