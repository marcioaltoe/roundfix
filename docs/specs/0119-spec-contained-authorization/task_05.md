---
task: task_05
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: high
---

# Task 05: Audit the consuming commit against the grant's ancestor

## Overview

An executor must not be able to widen its own grant to make its change pass.
Read the grant from the ancestor that authorized the consuming commit and audit
the commit's actual governed changes against that ancestor's bounded paths. The
slice is verifiable on its own against a real temporary Git history: a commit
that creates or widens its own grant refuses, and one whose grant already
existed passes.

## Requirements

1. MUST resolve the operative grant from the ancestry that precedes the
   consuming commit, and MUST record which record path and revision the audit
   read so a later reader can retrace the decision.
2. MUST refuse a consuming commit that creates or widens the grant it depends
   on, and MUST NOT let a later amendment authorize an earlier change
   retroactively.
3. MUST accept an earlier bounded amendment as authority for work that comes
   after it.
4. MUST discover operative records in active Specs, in archived Specs, and in
   the preserved legacy location, and MUST NOT narrow the governed set the audit
   judges against.
5. MUST keep the existing changed-path refusals and their reported tokens
   working for the conditions they already own, including a governed change
   outside the bounded set and a grant edited in the commit that consumes it.
6. MUST report an unavailable revision or unreadable record as an unresolved
   audit input rather than as a pass.

## Subtasks

- [ ] Resolve the authorizing ancestor for a consuming commit.
- [ ] Record the record path and revision the audit read.
- [ ] Refuse self-approval and retroactive approval; accept earlier amendments.
- [ ] Discover records in active, archived, and preserved legacy locations.
- [ ] Prove each case against a real temporary Git history.

## Acceptance Criteria

- [ ] A commit that adds or widens its own grant in the same commit refuses,
      naming the grant path.
- [ ] A commit whose grant landed in an earlier ancestor passes, and the audit
      result carries the record path and revision it read.
- [ ] A grant amended after the consuming commit does not authorize it; the same
      amendment authorizes a later commit.
- [ ] A grant recorded in an archived Spec and one in the preserved legacy
      location both resolve for a commit that consumes them.
- [ ] A governed change outside the bounded set still refuses with the token it
      reports today.
- [ ] An unavailable revision reports an unresolved audit input, distinguishable
      from a pass.

## Context

- interface: `internal/speccheck/mechanical.go`
- interface: `internal/gittest/gittest.go`

## Verification

- `grep -q 'func TestAuditReadsTheAuthorizingAncestor' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestAuditReadsTheAuthorizingAncestor$'` — the audit resolves the grant from the authorizing ancestor and records the path and revision it read.
- `grep -q 'func TestAuditRefusesSelfApprovalAndRetroactiveGrants' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestAuditRefusesSelfApprovalAndRetroactiveGrants$'` — same-commit self-approval and later retroactive approval refuse against a real temporary Git history, while an earlier amendment authorizes later work.
- `grep -q 'func TestAuditDiscoversRecordsInEveryLocation' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestAuditDiscoversRecordsInEveryLocation$'` — active, archived and preserved legacy records all resolve, and an unavailable revision reports unresolved.
- `grep -q 'func TestAuditReadsTheAuthorizingAncestor' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck` — the package suite runs with the ancestor audit present.

## References

- `_prd.md` → User Stories 3; Core Features 2; Goals 2, 4; Success Metrics.
- `_techspec.md` → Implementation Design: Audit and compatibility; Build Order 3; Risks & Considerations.
- ADR-0117, ADR-0130.
