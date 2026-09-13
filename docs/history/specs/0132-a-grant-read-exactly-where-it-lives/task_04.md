---
task: task_04
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
type: backend
complexity: high
---

# Task 04: Resolve a citation against the artifact that carries it

## Overview

For an external Spec Root the constraint reader produces a display path that
reaches outside the code repository, so joining a Spec-relative link to it
yields a path the containment check rejects and a valid record is refused. This
slice resolves the reference against the artifact's own location and anchors
containment to the Spec Root. It is verifiable alone: the external citation
resolves and the default one is unchanged.

This is an authorized tooling Task. It may change only
`internal/speccheck/constraints.go`,
`internal/speccheck/constraints_characterization_test.go`, and this Task file.
Stop before any other mutation. The bounded set comes from
[_authorization.md](_authorization.md).

## Requirements

1. MUST resolve a Spec-relative citation against the artifact that carries it,
   not against the code repository root.
2. MUST anchor the containment check to the resolved Spec Root, so a record
   beside its PRD in an external root is inside the boundary rather than
   outside it.
3. MUST keep rejecting a reference that escapes the Spec Root, including one
   that traverses upward past it or resolves through a symlink out of it.
4. MUST keep the default root's resolution identical, including which record a
   row citing several references resolves to.
5. MUST expose the resolved reference so the mechanical consumer can read it in
   Task 08, without changing that consumer here.
6. MUST update only the characterization rows this Task intentionally moves.

## Subtasks

- [ ] Resolve the reference relative to the carrying artifact.
- [ ] Anchor containment to the Spec Root.
- [ ] Prove escape is still rejected.
- [ ] Prove the default root's answers are unchanged.

## Acceptance Criteria

- [ ] A record cited as `[_authorization.md](_authorization.md)` beside a PRD in
      an external Spec Root resolves; it was refused before this Task.
- [ ] A reference traversing upward past the Spec Root, and one resolving
      through a symlink out of it, are still rejected.
- [ ] The default root resolves the same record it resolves today, including for
      a row citing an approved grant beside a proposed record.
Task 08 owns the mechanical consumer. This Task resolves the reference; making
the audit consume it requires `internal/speccheck/mechanical.go`, which this
Task's bounded scope forbids, so demanding it here would ask for a state the
same Task file refuses.

## Context

- interface: `internal/speccheck/constraints.go`

## Verification

- `grep -q 'func TestCitationResolvesInExternalSpecRoot' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestCitationResolvesInExternalSpecRoot$'` — an external Spec-relative citation resolves; this fails today.
- `grep -q 'func TestCitationRejectsEscapeFromSpecRoot' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestCitationRejectsEscapeFromSpecRoot$'` — upward traversal and symlink escape are still rejected.
- `grep -q 'func TestCitationResolvesInExternalSpecRoot' internal/speccheck/constraints_characterization_test.go || exit 1; go test -count=1 ./internal/speccheck` — the package suite passes with the default-root answers unmoved.

## References

- `_prd.md` → Goals 2; User Stories 2; Core Features 3.
- `_techspec.md` → Implementation Design: Paths derived from the resolved root; Build Order 4.

## Result

Implementation evidence:

- A Markdown authorization citation now resolves from the carrying artifact's
  filesystem path. The reader retains the repository-relative display path,
  while the selected reference also carries its resolved filesystem path and
  the root-relative path used to read it.
- Containment is anchored to the Spec Root derived from the carrying PRD or
  TechSpec. The check canonicalizes the longest existing path prefix, so an
  upward traversal and a symlink whose target leaves that root cannot become
  authorization candidates.
- Default-root citations continue to read from the project repository with
  their existing repository-relative paths. Reference order and operative
  versus proposed selection are unchanged.
- The resolved path, read root, and read path remain on the selected
  `authorizationReference`, making the exact result available to the
  mechanical consumer in Task 08 without changing that consumer here.
- The Task 01 external-root characterization moved from the recorded refusal
  to the resolved grant. The only new characterization is its negative
  companion for upward and symlink escape.

Focused-check evidence:

- The first two unchanged `rtk go test` attempts were blocked before test
  execution by `operation not permitted` from the sandboxed Go build cache.
  Re-running with an isolated cache reached the intended pre-change signal:
  `TestCitationResolvesInExternalSpecRoot` failed because the valid external
  citation produced `SC-TOOLING-UNAPPROVED` and did not identify one record.
- `rtk proxy env GOCACHE=/var/folders/p7/wnjmrqxd6kzcb0zs8c958j7c0000gn/T/tmp.TvTwLT2Siv go test -count=1 ./internal/speccheck -run '^Test(CitationResolvesInExternalSpecRoot|CitationRejectsEscapeFromSpecRoot|ConstraintReaderCharacterizesGrantCitation|ConstraintsResolveSpecContainedRecord|ConstraintsAcceptHonestProposalDeclaration|ConstraintsRefuseNonOperativeGrant|CheckReplay.*|AuditReadsTheResolvedReference|AuditRefusesOutOfGrantUnderBothCitationForms|MechanicalAuthorizationReadsThePRDBoundedDeclaration)$'`
  exited 0 after the last production edit. This covers both external escape
  cases, the default-root citation forms, approved-versus-proposed selection,
  relative-root replay paths, and the unchanged mechanical consumer contract.
- `rtk make fmt-check build GO=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt`
  exited 0.
- The required `rtk make verify-incremental` first reached a sandbox block on
  its read-only `api.github.com` request. The approved unchanged rerun reached
  the full test stage and exited 2: three `internal/daemon` fixtures supply
  non-Git external Spec Roots to the Task 03 resolver, which refuses them at
  `git rev-parse --show-toplevel`. `internal/speccheck` passed in that run.
- The Task's declared `## Verification` commands were not run; Daemon
  Verification owns them after this handoff.

Acceptance-criterion evidence:

- External Spec Root: `TestCitationResolvesInExternalSpecRoot` passed against
  separate committed project and Spec repositories, with no tooling finding
  or authorization skip.
- Escape rejection: `TestCitationRejectsEscapeFromSpecRoot/upward_traversal`
  and `/symlink_outside_root` passed and required an exact-record refusal.
- Default-root compatibility: the focused run passed the existing citation
  form characterization and
  `TestConstraintsAcceptHonestProposalDeclaration/approved_narrow_record_wins_over_proposed_record`,
  preserving the selected approved grant when the row also cites a proposal.

Follow-up:

- The failing `internal/daemon` external-root fixtures must use Git-backed Spec
  Roots or be reconciled with the Task 03 resolver contract in their owning
  slice. Those files are outside this Task's authorized mutation set.
