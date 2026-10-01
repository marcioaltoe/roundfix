---
status: accepted
created_at: 2026-10-01T21:00:00Z
updated_at: 2026-10-01T21:00:00Z
deprecated_at: null
superseded_by: null
---

# A measurement lands as a record, and archived evidence leaves only with approval

Three open questions each proposed a change on a belief: reopening a memory
layer for Agent Sessions, adopting a Task acceptance judgment, and dropping
archived QA evidence from the history root. On 2026-10-01 the maintainer
decided that this work measures and proposes, and that any removal of
archived evidence waits for the maintainer's explicit approval in a later
change. Roundfix therefore keeps measuring separate from changing:

- **A record and a rule.** A measurement lands in the repository as a record
  that a test can recompute, with a written recommendation whose decision
  rule is stated before the data are read.
- **No change rides on it.** A measurement changes no gate, Verification,
  judgment or skill. A change the measurement supports is written as a
  Backlog Entry for a later Spec.
- **Nothing archived is removed.** No file under the history root is deleted,
  moved or rewritten by a measurement or by the Spec that delivers it. A
  removal is its own later change, and the maintainer's explicit approval of
  that removal is recorded in it. A measurement that needs a tree without
  archived evidence builds that tree in a disposable copy outside the
  repository.
- **A re-measured judgment keeps the judge's boundary.** It sends only this
  repository's Task files, which the maintainer's authorization of the judge
  already covers, through the judge's transports, keys, Judge Log and monthly
  ceiling. It runs from a test harness only when its flag is passed, so it is
  not compiled into the binary and never runs in a Verification or a QA gate.

## Consequences

- A measurement can end without a change, and that is a complete result.
- A recommendation to remove archived evidence can stand open until the
  maintainer answers it.
- Removing a file from the tree does not shrink the repository's history: a
  clone still downloads every version of it.
