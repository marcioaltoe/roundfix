# Spec 0058 replay provenance

The replay copies its base artifacts at test runtime from
`docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight/`.
Its accepted QA Report reproduces
`docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight/qa/qa-report-2026-08-01-04.md`
after the normal integration and review journey, leaving the real tagged OIDC
release as the only unmet row.

The original Spec 0058 PRD has no `## Unreachable Acceptance` section.
The declaration overlay is added by Spec 0070; it was not present when Spec
0058 ran or when its archived report was written. The replay removes the
original `qa_override` stamp before invoking the Archive Command.

The `wrongly-declared` variant records the release journey as reachable and a
wrongly-declared-row finding. The `unmatched` variant reports one declared
blocked row while omitting the declaration overlay, preserving the archive
refusal for an unmatched row.
