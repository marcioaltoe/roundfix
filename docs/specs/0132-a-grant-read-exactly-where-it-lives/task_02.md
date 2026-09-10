---
task: task_02
spec: 0132-a-grant-read-exactly-where-it-lives
status: pending
type: backend
complexity: medium
---

# Task 02: Close the frontmatter on a complete marker line

## Overview

A purported closing line such as `---evil` currently closes the frontmatter by
prefix, leaving the rest of the line discarded as body while every other field
parses cleanly, so a malformed record resolves to a grant. This slice makes the
delimiter exact: the marker and nothing else. It is verifiable alone — the
malformed record refuses and the well-formed one still grants.

## Requirements

1. MUST treat a frontmatter delimiter as a complete line: the marker alone,
   allowing only a trailing carriage return.
2. MUST classify a record whose opening or closing delimiter line carries extra
   characters as malformed, resolving to a refusal that names the offending
   line, never to a grant.
3. MUST keep a well-formed record resolving exactly as it does today, including
   its fields, paths and operations.
4. MUST keep an unreadable record and an unavailable revision reporting
   unresolved rather than refused, so the new refusal fires on a positive
   observation of malformation and never on missing evidence.
5. MUST NOT change any other parse rule, field name or reason code.
6. MUST update only the characterization row this Task intentionally moves.

## Subtasks

- [ ] Require a complete marker line at both delimiters.
- [ ] Classify a malformed delimiter as a refusal naming the line.
- [ ] Prove a well-formed record is unchanged.
- [ ] Move the one recorded characterization row.

## Acceptance Criteria

- [ ] A record whose closing line is `---evil` refuses, and the refusal names
      that line; it granted before this Task.
- [ ] A record whose opening line carries extra characters refuses the same way.
- [ ] A well-formed record grants the identical fields, paths and operations it
      granted before.
- [ ] A record with a trailing carriage return on its delimiter still parses.
- [ ] An unreadable record still reports unresolved, distinguishable from
      refused.

## Context

- interface: `internal/authorization/authorization.go`

## Verification

- `grep -q 'func TestMalformedDelimiterRefuses' internal/authorization/authorization_test.go && go test -count=1 ./internal/authorization -run '^TestMalformedDelimiterRefuses$'` — a malformed delimiter refuses and names the line; this fails today because the record grants.
- `grep -q 'func TestMalformedDelimiterRefuses' internal/authorization/authorization_test.go || exit 1; go test -count=1 ./internal/authorization -run '^TestAuthorizationParseCharacterization$'` — the recorded answers still hold with exactly the one intended row moved.
- `grep -q 'func TestMalformedDelimiterRefuses' internal/authorization/authorization_test.go || exit 1; go test -count=1 ./internal/authorization ./internal/spec ./internal/speccheck ./internal/suiteguardcontract` — every consumer of the reader still passes.

## References

- `_prd.md` → Goals 1; User Stories 1; Core Features 1; Decisions: Declared intentional breaks 1.
- `_techspec.md` → Implementation Design: Exactness at the parse boundary; Build Order 2.
