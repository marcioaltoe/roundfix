---
status: accepted
created_at: 2026-08-03T00:00:00Z
updated_at: 2026-08-03T12:00:00Z
deprecated_at: null
superseded_by: null
---

# Triage intent lives in a typed backlog

The documentation layout gave observations a home (`docs/findings/`: dated,
evidence-backed, immutable field reports) and gave intent none, so
suggestions without observed evidence were either forced into findings —
diluting what a finding means — or lost. The layout therefore gains
`docs/backlog/`, holding typed intent entries whose type vocabulary is the
Conventional Commits vocabulary the repository already enforces on commits
and PR titles — `feat`, `fix`, `perf`, `refactor` — so one word carries the
intent from backlog entry to Spec to commit; entries follow a frontmatter
lifecycle of `open`, `promoted` to a named Spec, or `declined` with a
reason, and promotion moves
the entry into the adopting Spec's `references/` exactly as finding adoption
works. The directory is named `backlog` rather than `ideas` because intent
has more than one shape and the name must not change when the next shape
arrives — including the possible future deprecation of `docs/findings/`
into a backlog type, which stays an explicitly deferred decision. A `feat`
entry is upstream raw material the spec pipeline may consume, never the
`write-idea` artifact itself; a finding is never a commitment and a backlog
entry is never evidence. Reusing the commit vocabulary was chosen over a
bespoke type set because the maintainer's own usage showed findings
absorbing problem reports that were really `fix` intent — the shared word
is what keeps the intent legible across every artifact it passes through.
