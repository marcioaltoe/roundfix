---
source: coderabbit
pr: "137"
round: 1
round_created_at: "2026-08-07T03:22:12Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/promoted-backlog-entries-have-a-home
head_sha: ea93c68b70d066c1ee7f322e40ac1d547420e8be
file: docs/workflow/authorizations/2026-08-06-promoted-backlog-destination.md
line: 21
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XKdmt,comment:PRRC_kwDOS0qyts7eg8KM
review_hash: 3a3e0209004e31fab9636edcae861b6b29171b97731b2d3e09ab59d6d56d8639
duplicate_of: ""
source_review_id: "4879615443"
source_review_submitted_at: "2026-08-07T03:12:07Z"
---

# Issue 002: _ Security & Privacy_ _ Major_ _ Quick win_

## Review Comment

_🔒 Security & Privacy_ | _🟠 Major_ | _⚡ Quick win_

**Authorize exact files, not directory globs.**

Lines [12-20] grant access to entire directory trees. This can authorize unrelated future changes under those trees. Enumerate the exact skill files, mirror files, detector files, tests, backlog entries, reference files, and index files changed by this work.

As per coding guidelines, protected-tooling authorization must list exact bounded repository-relative files, not directory globs.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/workflow/authorizations/2026-08-06-promoted-backlog-destination.md`
around lines 10 - 21, The Bounded files section uses directory globs instead of
exact authorization boundaries. Replace each glob with the specific
repository-relative skill, mirror, detector, fixture, test, backlog, reference,
and index files changed by this work, while retaining only the explicitly
required files and the make skills-sync relationship.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:781946c020240f8b8e5c8c56 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The four directory globs authorized unrelated future files and did not
  satisfy the exact bounded-file requirement. The authorization now names the
  18 repository-relative source, destination, skill, mirror, detector, test,
  finding, and authorization paths changed by the original feature commit.

## Resolution

- Replaced every directory glob in
  `docs/workflow/authorizations/2026-08-06-promoted-backlog-destination.md`
  with exact repository-relative files and retained the `make skills-sync`
  mirror relationship.
- Focused evidence: an `awk` extraction of backticked paths from `## Bounded
  files` returned 18 exact paths; `rtk git diff --name-status --no-renames
  main...HEAD` returned the same 18-path set for the feature commit.
- Scope evidence: `rtk git diff --check` exited 0.
- Authoritative Verification `make verify` was not run; the Daemon owns it
  after this Agent turn.
