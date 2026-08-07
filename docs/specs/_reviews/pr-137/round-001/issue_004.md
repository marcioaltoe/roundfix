---
source: coderabbit
pr: "137"
round: 1
round_created_at: "2026-08-07T03:22:12Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/promoted-backlog-entries-have-a-home
head_sha: ea93c68b70d066c1ee7f322e40ac1d547420e8be
file: internal/speccheck/backlog.go
line: 141
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XKdm3,comment:PRRC_kwDOS0qyts7eg8KW
review_hash: 81ed32066ce46f04148bfb60717ddab90bc8c0d3cbb5317066e26f830505e568
duplicate_of: ""
source_review_id: "4879615443"
source_review_submitted_at: "2026-08-07T03:12:07Z"
---

# Issue 004: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Accept CRLF frontmatter delimiters.**

The exact `---\n` and `\n---` checks do not match CRLF files. A promoted entry with `---\r\n` returns empty frontmatter, so `SC-BACKLOG-UNMOVED` does not report it.

Normalize line endings before delimiter parsing. Add a `Check` test with CRLF frontmatter that requires the finding.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/backlog.go` around lines 132 - 141, Update
parseBacklogFrontmatter to normalize CRLF line endings before checking the
opening and closing frontmatter delimiters, while preserving existing parsing
behavior. Add a Check test covering CRLF frontmatter for a promoted entry and
assert that SC-BACKLOG-UNMOVED is reported.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7cefaf382b1aa0589f9ae65b -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: A CRLF opening delimiter failed the exact LF prefix check before YAML
  parsing, so `Check` could not report the promoted entry. This is the same
  production defect independently described by issue 003.

## Resolution

- Normalized CRLF to LF before frontmatter delimiter recognition.
- Added a public `Check` regression case that writes CRLF frontmatter and
  requires the `SC-BACKLOG-UNMOVED` finding with destination guidance.
- Reproduction evidence before the production fix: the regression failed with
  no `SC-BACKLOG-UNMOVED` finding for the CRLF entry.
- Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/speccheck -run TestCheckBacklogUnmoved -count=1` passed after the
  fix; `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/speccheck -count=1` also passed.
- Authoritative Verification `make verify` was not run; the Daemon owns it
  after this Agent turn.
