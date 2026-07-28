---
status: accepted
created_at: 2026-07-28T21:03:33Z
updated_at: 2026-07-28T21:03:33Z
deprecated_at: null
superseded_by: null
---

# QA verdicts distinguish environment-blocked rows

The QA gate runs in a Run Worktree with no commit or push authority, so
acceptance rows whose journeys need a live Pull Request or Final Push are
structurally unreachable, and recording them `blocked` capped the verdict at
`partial` forever — a Spec's most valuable journeys made it unarchivable.
Blocked rows therefore carry a typed cause: a row blocked only by the gate's
environment does not cap the verdict when the report records equivalent
observed or supervised evidence for it, while a finding-blocked, failed, or
unevidenced row caps the verdict exactly as before. Silently loosening the
gate or granting Agents push authority was rejected: the gate still never
credits a journey without evidence, and the distinction makes the evidence
requirement explicit instead of impossible.
