---
status: accepted
created_at: 2026-10-07T00:00:00Z
updated_at: 2026-10-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Legacy Archive Folder is read leniently, and a failed QA keeps its verdict

ADR-0248 converts each Legacy Archive Folder through `BuildArchiveRecord`, the
builder the Archive Command uses, and the History Sanitize Command refuses the
whole plan at the first folder that builder rejects. On 2026-10-07 an adopter
running 0.55.0 got no plan at all. Four of its folders failed checks written
after they were archived. Two Task Graph projection tables named Tasks the
author had moved out of the graph, one table used a Task type the closed set no
longer holds, and one Spec was archived by maintainer decision with a failing
QA Report and no override stamp. That last folder's record took the
disposition `fail`, which `ParseArchiveRecord` refuses. The adopter cannot
repair any of them: archived metadata is never hand-edited, and
`archive --qa-override` acts only on an active Spec. On 2026-10-07 the
maintainer chose all three repairs ("Sim, as três frentes") and, for a batch
that meets a folder it cannot convert, "Pular e seguir".

**Lenient Legacy Reading.** A Legacy Archive Folder's Task Graph is read for
what an Archive Record needs, its QA Task, as a tolerant reader does. The
manifest front matter keeps every check it has today. In the projection table,
a row that names a Task outside the graph and a type outside today's closed set
are tolerated and named, never rewritten. Malformed and duplicate rows still
refuse. The rule applies wherever a Legacy Archive Folder becomes a record: the
History Sanitize Command and the folder branch of `ReadArchivedSpec`, so the
plan and every later reader of that folder agree. It never applies to an active
Spec. `spec.Load`, the Implement Command and the Archive Command keep the
current schema.

**The `failed-qa` disposition.** A folder whose newest QA Report says `fail`
and that carries no override receives the disposition `failed-qa`. The record
keeps `qa_report` and `qa_verdict: fail`, so the verdict reads as it was
recorded. It is never `pass`, it is never `qa-override` (ADR-0154 reserves that
for recorded user authority), and its QA Task never counts as completed. The
name pairs with `no-qa` from ADR-0248, and the record's `status: archived`
already says "archived". The schema stays `roundfix/archive-record/v1`. The
Archive Command still refuses to archive an active Spec on a failing QA without
an override, so in practice only the History Sanitize Command writes this
disposition.

**Refused Units.** A unit the plan cannot convert is a Refused Unit. The plan
names it with its reason and goes on to the next unit. `--apply --batch <n>`
converts the next `n` units that can be converted and leaves each Refused Unit
untouched and listed. Refused Units do not count toward `n`. They sit at the
head of the fixed order until someone repairs them, so if they counted, every
later batch would spend its places on the same refusals and convert less. A
batch whose every examined unit is refused writes nothing and exits 2, which
gives a loop over batches a stopping point. Two refusals still stop the whole
command. A unit named by `--promote` refuses the command when it is refused,
because the operator asked for that unit. A Git failure while reading delivery
history also refuses the command, because it belongs to the environment and
not to the unit.

## Consequences

Two alternatives were rejected. Converting a failed QA as an implicit
`qa-override` would claim a user authority no record holds. Rewriting the
legacy manifests would change archived bytes to satisfy a schema written
later. The record readers that branch on disposition were checked, and only
`ArchivedTaskCompleted` changes. Delivery treats a `failed-qa` record as an
archive without an override. Review, reconcile, the archive-license check and
cause reports read it as any other record.
