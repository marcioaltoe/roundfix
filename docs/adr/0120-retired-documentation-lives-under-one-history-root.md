---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-12T12:41:53Z
updated_at: 2026-08-12T14:16:46Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Retired documentation lives under one history root

Spec 0085 built the single archive root this repository asked for and anchored it
at the repository root, where an archive of documentation sits beside the Go source
trees; it also deferred relocating retired ADRs, leaving inactive decisions inside
the directory an Agent loads as decision context, and it never covered Review
Artifacts, whose orphan case still accumulates beside the active Spec Root. The
root is `docs/history/`, holding `specs/`, `findings/`, `adr/`, `backlog/`, and
`reviews/`, because retired documentation is documentation and one root for every
retired family is what makes a single exclusion sufficient. `history` is chosen
over `archived` and `retired` because the tree holds three things a reader asks
for by different names — what was built, what was decided and undone, and what a
review said — and over `references`, which this repository already uses for
durable, active, upstream reference. The name carries no underscore, so the review
tool's exclusion becomes path-anchored rather than suffix-matched, which is the
reasoning that file already applies to its other unprefixed tree. This supersedes
Spec 0085's goal that retired material leave the directories an Agent loads by
default — an Agent reaches `docs/agents/` through `CLAUDE.md` rather than loading
`docs/` wholesale — while preserving that Spec's other two goals, one archive root
and one path filter, unchanged.
