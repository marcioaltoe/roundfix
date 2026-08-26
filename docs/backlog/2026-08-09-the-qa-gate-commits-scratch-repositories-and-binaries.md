---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-09
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The QA gate commits scratch repositories and binaries into the Spec

## Opportunity

A gate that proves a CLI behaviour builds isolated repositories to run it in, and
builds the binary it exercises. Spec 0090's gate committed both into the Spec's
evidence directory: nine directories recorded as Gitlinks (mode `160000`) and six
executables — `roundfix` twice, `codex-acp` twice, `acpx` twice.

A Gitlink is a submodule reference with no submodule behind it. Every later
`git status` in a fresh clone reports them, and nothing can resolve them.

## Value

Measured on 2026-08-09 while closing Spec 0090. The gate's own finding (F-006)
caught the Gitlinks and classified them Blocks-Completion; the binaries it did
not catch, and they reached a commit on the branch.

This is not specific to that Spec. Any authored gate whose acceptance needs a
real Run in an isolated repository produces the same residue, so it recurs once
per Spec that takes evidence seriously — which is the kind the contract asks for.

## Shape

Non-binding. The question worth settling is where a gate's scratch state should
live so that evidence stays committed and the machinery producing it does not.
An evidence directory that names paths outside the repository, or an ignore rule
the gate owns, both answer it; the second is cheaper and the first is more
auditable.

Worth settling in the same work: whether the gate should refuse to commit a file
it did not author as evidence — a built binary is never evidence of anything a
reader can check, and it is large.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
