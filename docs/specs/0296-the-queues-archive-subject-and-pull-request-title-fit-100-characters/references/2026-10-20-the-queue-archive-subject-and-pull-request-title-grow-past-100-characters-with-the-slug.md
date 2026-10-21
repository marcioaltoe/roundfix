---
type: fix
status: promoted
created: 2026-10-20
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
reason: null
---

# The queue's archive subject and pull request title grow past 100 characters with the slug

## Symptom

The Delivery Queue writes two messages that carry the whole Spec slug: the
archive commit `docs: archive <slug>` (`internal/cli/deliver_workflow.go`,
`Archive`) and the pull request title `feat: deliver <slug>`
(`internal/cli/deliver_workflow.go`, `Publication`), which a squash merge takes
as the commit subject with ` (#<number>)` appended. Each is 14 characters plus
the slug, so a slug longer than 86 characters makes either exceed the
100-character `header-max-length` of `@commitlint/config-conventional`. A
repository whose `commit-msg` hook runs commitlint refuses the archive commit,
and a repository that lints pull request titles or squash subjects refuses the
title.

Roundfix enforces no slug length. A slug is one path component, so the
Archive Record `<slug>.md` bounds it at 252 bytes, and the item branch
`roundfix/deliver-<slug>-<16 hex>` bounds a deliverable slug at 230 bytes,
where its last ref component reaches the common 255-byte file name limit.
At 230 bytes each message is 244 characters long.

## Where

`roundfix deliver`: `Archive` and `Publication` in
`internal/cli/deliver_workflow.go`. Spec
`0290-the-pull-request-record-commit-fits-commitlint` shortened only the Pull
Request Record Commit subject, by the maintainer's decision of 2026-10-20, and
made a refused queue commit park with the refusing hook and the hand-commit
recovery.

## Expected

Both messages stay within 100 characters for every deliverable slug, for
example by naming the Spec by its number as the record subject does, while the
pull request still identifies the Spec in its body. Source Revision Recovery,
which reads the trailing `(#N)` of the squash subject, keeps working.

## Evidence

This repository's own history: the squash commit `9e92ac39` on `main`,
`feat: deliver 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers (#575)`,
has a 117-character subject, and `docs: archive <slug>` for that 96-character
slug is 110 characters. The adopted Inbox Entry of Spec 0290 measured 68 and 67 characters
for the oraculum slug, below the limit only because that slug is 54
characters.
