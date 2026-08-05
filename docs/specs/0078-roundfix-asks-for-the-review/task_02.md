---
task: task_02
spec: 0078-roundfix-asks-for-the-review
status: pending
type: backend
complexity: medium
---

# Task 02: Ask at the seam where the Round pushes

## Overview

`deps.Resolver.Resolve` performs a Round's fixes, commit, and Final Push and
returns the new `resolved.HeadSHA`. Every path after it waits for Evidence
bound to that head. This slice publishes the request in between.

It is the only Task that can produce a duplicate request, so it carries the
one-per-Round assertions — including the case that looks like two pushes: the
artifact-only docs commit, whose descendant inherits its parent's Evidence
under ADR-0036 and is not worth a review of its own.

## Requirements

1. MUST publish exactly one review request per Round, for the head
   `Resolve` reports, when `review_source.request_review` is enabled.
2. MUST publish it after the Final Push and before the Run waits for Evidence
   on that head, on both paths out of `Resolve`: the merge-ready confirmation
   and the next Round's wait.
3. MUST NOT publish for the artifact-only docs commit created after the Final
   Push.
4. MUST NOT publish when the Round pushed no new head.
5. MUST publish from `resolve` after its push, under the same configuration.
6. MUST NOT publish from `fetch` under any configuration.
7. MUST leave the Spec 0077 Evidence classification untouched: a refused
   request still resolves `skipped`, still ends the Run naming the refusal, and
   is followed by no second request.
8. MUST preserve today's control flow exactly when the configuration is
   disabled, including a nil requester.

## Subtasks

- [ ] Add the optional requester dependency and the call at the `Resolve` seam.
- [ ] Add the `resolve` command call site.
- [ ] Assert one request per Round across both post-`Resolve` paths.
- [ ] Assert the disabled default changes nothing.

## Acceptance Criteria

- [ ] A Round that pushes publishes exactly one request, for the pushed head.
- [ ] A Round whose Final Push is followed by the artifact-only docs commit
      still publishes exactly one request, for the fix head.
- [ ] A Round that pushes nothing publishes no request.
- [ ] `resolve` publishes one request after its push.
- [ ] `fetch` publishes none, asserted rather than assumed.
- [ ] A refused request ends the Run Review Skipped naming the refusal, with no
      second request in the same Run.
- [ ] With `request_review` disabled, every existing watch, resolve, and fetch
      test passes unchanged.

## Context

- interface: `internal/watch/watch.go`
- interface: `internal/cli/cli.go`
- instruction: `docs/adr/0036-review-artifacts-are-committed-in-a-separate-docs-commit.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/watch -count=1 -run 'Request|Round|Seam|Artifact' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the seam tests ran and passed.
- `go test ./internal/watch ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD | grep -qE "^(\.coderabbit\.yaml|\.roundfixrc\.yml)$" && exit 1 || exit 0`
  — expected: exit 0; turning the flow on is task_04's scope.

## References

- `_prd.md` → Core Features 1, 3, 5 and 6; Success Metrics 1, 2 and 4.
- `_techspec.md` → System Architecture; Build Order 2.
- ADR-0036, ADR-0054.
