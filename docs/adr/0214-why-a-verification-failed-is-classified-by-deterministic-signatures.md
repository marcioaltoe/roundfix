---
status: accepted
created_at: 2026-10-01T21:00:00Z
updated_at: 2026-10-01T21:00:00Z
deprecated_at: null
superseded_by: null
---

# Why a Verification failed is classified by deterministic signatures

Whether a memory layer for Agent Sessions is worth reopening depends on how
many failed Verification attempts and corrective Tasks come from missing
repository knowledge rather than plain implementation error, and nothing
recorded that cause. A model asked to judge the cause from logs gives a
different answer from one run to the next, and published work that extracted
failure causes with a model needed a second filtering stage to remove claims
the logs did not support. Roundfix therefore classifies the cause by
signatures it can show, never by a model judgment:

- **A closed list.** Every item gets one class: `scope-or-authorization`,
  `shared-section-contract`, `repository-convention`,
  `implementation-defect` or `environment`. The first three are the
  repository-knowledge classes. An item that no signature matches is
  `unclassified`, which is counted and reported, never folded into a class.
- **Signatures, in one reviewable table.** A signature is an ordered entry
  with an identifier, a class, the one evidence source it reads (the tail of
  the Verification diagnostic, the Verification command, or the corrective
  Task's title and overview) and a regular expression. The first entry that
  matches decides. The table is embedded in the binary, and every report names
  the digest of the table that produced it.
- **A corrective Task is found by its place in the graph.** It is a Task
  that is not the QA Task and whose number follows the QA Task's in its Task
  Graph, which is how corrective Tasks have always been appended.
- **Read-only.** The report opens the Run Database with the reader that
  never migrates and writes nothing, here or anywhere else.

## Consequences

- A keyword can match for the wrong reason. Each item names the signature
  that classified it, so a reviewer can see and dispute the reason.
- Evidence that Journal Retention or the GC Command already removed cannot be
  classified, so the item is `unclassified` with no evidence.
- A signature added after the data were read is disclosed in the record that
  used it, because it was fitted to that data.
