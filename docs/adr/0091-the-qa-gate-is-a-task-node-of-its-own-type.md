---
status: accepted
created_at: 2026-08-03T00:00:00Z
updated_at: 2026-08-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The QA gate is a Task node of its own type

ADR-0088 moved the QA gate from the Implement Command's `--qa` parameter into
the Spec's artifacts; this records how it is represented: as an ordinary
`task_NN.md` with a new Task Type `qa`, named by a `qa:` declaration in the
manifest frontmatter, required to be terminal and to depend on every leaf —
while a Spec that needs no gate declares `qa: declined` with a reason, and a
graph with neither declaration is a legacy graph that keeps working
unchanged. A manifest-only property was rejected because the graph already
owns ordering, status, and history, and a node inherits all three for free —
including the invalidation semantics: appending work after the gate reported
either leaves the gate non-terminal or puts a pending dependency beneath a
settled terminal, and both are load-time validation errors that name the
insertion. The node routes to the Daemon's existing gate step, so report,
verdict, and typed blocked-row counts stay exactly as ADR-0080 defines them.
