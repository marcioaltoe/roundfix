---
task: task_05
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
type: backend
complexity: high
---

# Task 05: Report Review Issue knowledge and terminal context

## Overview

Carry whether Review Issues were actually fetched and one bounded terminal
context through text reports and the durable outcome stream. Pre-fetch failure
becomes explicitly unknown instead of appearing as five misleading zero
counts, and every non-Clean outcome gains an actionable next step.

## Requirements

1. MUST mark Review Issue knowledge true only after fetch completes.
2. MUST report `Review Issues: unknown — fetch did not complete.` when
   knowledge is false and omit all zero-valued status counts.
3. MUST retain valid zero-issue summaries after a successful fetch.
4. MUST carry bounded terminal reason, next action, Evidence, verified head,
   Console Log, and Attach command through one terminal context.
5. MUST project additive optional fields through `roundfix-events/v1`.
6. MUST preserve existing stdout, stderr, filter, and unknown-field contracts.
7. MUST publish terminal context only through the completion winner.

## Subtasks

- [ ] Extend watch result with issue knowledge and terminal context.
- [ ] Render unknown versus known-zero reports.
- [ ] Build one terminal-context boundary.
- [ ] Extend additive outcome stream projection.
- [ ] Preserve existing JSONL and filter behavior.
- [ ] Add pre-fetch, fetched-zero, and terminal-outcome cases.

## Acceptance Criteria

- [ ] Failed status discovery or fetch emits one unknown line and no numeric
      Review Issue counts.
- [ ] A completed zero-issue fetch remains known and reports valid zeros.
- [ ] Every non-Clean outcome event contains reason and next action.
- [ ] Console Log, Attach command, issue knowledge, and Evidence are present
      when available and omitted safely when unavailable.
- [ ] Existing `roundfix-events/v1` consumers can ignore additive fields.
- [ ] Requested output remains on stdout and diagnostics remain on stderr.
- [ ] A losing completion path emits no duplicate terminal context.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/event_test.go`
- interface: `internal/runevent/stream.go`

## Verification

- `rtk go test ./internal/watch ./internal/cli -run 'Test.*(ReviewIssuesKnown|ReviewIssuesUnknown|TerminalContext|FetchedZero)' -count=1`
  — expected: pre-fetch failures are unknown and successful zero fetches remain
  known.
- `rtk go test ./internal/runevent -run 'Test.*(Outcome.*Context|ReviewIssues|Evidence)' -count=1`
  — expected: additive outcome records carry bounded actionable context.
- `rtk go test -race ./internal/watch ./internal/cli ./internal/runevent -run 'Test.*(TerminalContext|ReviewIssues)' -count=1`
  — expected: terminal context publication is race-free.

## References

- `_prd.md` → Goals 3–4; User Stories 3–4 and 6; Core Features 6–9; User
  Experience; Success Metrics.
- `_techspec.md` → Data Models: watch.Result; API Contracts: Terminal report,
  stream, and notification; Build Order 5.
- `../../adr/0052-run-completion-is-compare-and-set.md` → winner-only outcome
  publication.
