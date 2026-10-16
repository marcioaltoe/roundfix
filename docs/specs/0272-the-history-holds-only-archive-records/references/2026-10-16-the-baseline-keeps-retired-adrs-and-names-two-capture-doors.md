---
type: fix
status: promoted
created: 2026-10-16
spec: 0272-the-history-holds-only-archive-records
reason: null
---

# The Baseline keeps retired ADRs whole and names two capture doors

## Symptom

The Baseline's docs-layout clause "Give each documentation directory one job"
contradicts what this repository decided and did on 2026-10-14. It says
rejected, deprecated and superseded ADRs "retire whole to `docs/history/adr/`"
and that a history sanitize leaves "retired ADRs stay whole", while the
maintainer had every retired ADR removed (30 files under `docs/history/adr/`,
and three retired ADRs consolidated into the accepted ADRs that restate them)
because each consolidating ADR restates the surviving decision and Git keeps
the text. The same clause lists `docs/_inbox/` "for raw notes" as a
live directory although ADR-0095 makes the Secondbrain inbox the single
capture door and keeps `docs/_inbox/` only for repositories without a brain,
and it calls `docs/handoffs/` "the capture door", the phrase ADR-0095, the
Secondbrain guide and the glossary reserve for the Secondbrain inbox.

## Where

The `context-workflow` Baseline module, rule `rule.context.docs-layout`,
clause `clause.context.docs-one-job-per-directory`, rendered into the
`guide.docs-layout` managed region of `docs/agents/docs-layout.md`; the
glossary term History Root ("holding Archive Records, Reduced History Entries
and retired ADRs").

## Expected

The clause says that a retired ADR leaves the tree once an accepted successor
names it with its date, lists `docs/_inbox/` only for a repository without a
Secondbrain, and uses another word than "capture door" for handoffs; History
Root agrees.

## Evidence

Triaged on 2026-10-16 from the Secondbrain Inbox Entry
`inbox/roundfix/2026-10-14-baseline-text-keeps-a-history-of-retired-adrs-and-names-two-capture-doors.md`,
written during the glossary congruence audit (branch
`docs/glossary-congruence`, commits `a6198793`, `fa4ca587`, `9a13d845`).
`internal/spec/archive_layout_characterization_test.go` already records the
ADR history family as absent.
