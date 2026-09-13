---
task: task_08
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
type: backend
complexity: medium
---

# Task 08: Let the mechanical audit consume the resolved reference

## Overview

Task 04 resolves a Spec-relative citation against the artifact that carries it.
The changed-path audit is the consumer that must read that resolved reference
and judge bounded paths against the project repository, and it lives in a file
Task 04's bounded scope forbids. This slice is that consumer. It is verifiable
alone: an external Spec Root produces a real audit where it produces a skip
today.

`internal/speccheck/mechanical.go` is not a governed path, so this Task needs no
tooling grant for it. Its regression lives in
`internal/speccheck/constraints_characterization_test.go`, which the approved
grant covers, rather than in `internal/speccheck/mechanical_test.go`, which
belongs to Spec 0119's grant and not to this one.

## Requirements

1. MUST carry the resolved reference from the constraint reader into the
   mechanical authorization request, so the audit never re-derives it with a
   narrower rule.
2. MUST judge bounded paths and ancestor authority against the project
   repository root and delivery target, never against the Spec repository, per
   the two-root boundary the TechSpec names.
3. MUST produce a real changed-path audit for a Spec whose record lives in an
   external Spec Root, where it records a skip today.
4. MUST keep every audit outcome that works today: the out-of-grant refusal, the
   grant-edited-in-the-consuming-commit refusal, and the presence-aware skip for
   a Spec that genuinely declares no authorization.
5. MUST NOT let an unresolved reference read as a skip, because a skip says
   there was nothing to audit and that is the failure this pairing removes.

## Subtasks

- [ ] Carry the resolved reference into the mechanical request.
- [ ] Judge paths against the project root while reading the record from the Spec root.
- [ ] Prove the external case produces a real audit.
- [ ] Prove the preserved outcomes and the unresolved-is-not-a-skip rule.

## Acceptance Criteria

- [ ] A Spec whose record lives in an external Spec Root produces a changed-path
      audit with Task commits; it records a skip today.
- [ ] The audit judges bounded paths against the project repository, proven by a
      case where the two roots hold different trees and judging the Spec root
      would reach the wrong answer.
- [ ] A governed change outside the grant still refuses, and a Spec declaring no
      authorization still records the presence-aware skip.
- [ ] An unresolved reference refuses rather than skipping.

## Context

- interface: `internal/speccheck/mechanical.go`

## Verification

- `grep -q 'func TestMechanicalAuditConsumesResolvedReference' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestMechanicalAuditConsumesResolvedReference$'` — an external Spec Root produces a real audit; this fails today.
- `grep -q 'func TestMechanicalAuditJudgesTheProjectRoot' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestMechanicalAuditJudgesTheProjectRoot$'` — bounded paths are judged against the project repository, proven where judging the Spec root gives the wrong answer.
- `grep -q 'func TestMechanicalAuditConsumesResolvedReference' internal/speccheck/constraints_characterization_test.go || exit 1; go test -count=1 ./internal/speccheck` — the package suite passes with the preserved outcomes unmoved.

## References

- `_prd.md` → Goals 2; Core Features 3; User Stories 2.
- `_techspec.md` → Implementation Design: Two roots, named separately; The boundary both roots travel in; Testing Approach observation 2.

## Result

Implementation evidence:

- The mechanical authorization resolver now preserves the constraint reader's
  selected citation as one request value containing the display path, Spec
  repository root and revision, repository-relative record path, project
  repository root, and delivery target. The path-only reader remains available
  to callers that do not have a committed delivery target.
- The changed-path audit reads committed authorization bytes from the Spec
  repository. It resolves the Task commit, delivery target, authorizing
  ancestor, changed paths, bounded paths, and sanctioned regeneration outputs
  against the project repository.
- Authorization-path self-approval checks run only when the Spec and project
  roots share Git history. A record in a separate repository cannot collide
  with an identically named project path.
- The Daemon's QA request carries the resolved reference whenever a Run has a
  delivery target. Target-less fixture requests retain their existing
  path-only, presence-aware behavior.

Focused-check evidence:

- Before the production edit,
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task08-gocache /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go test -count=1 ./internal/speccheck -run 'TestMechanicalAudit(ConsumesResolvedReference|JudgesTheProjectRoot)'`
  exited 1 because the resolved-reference API and mechanical request field did
  not exist.
- After the final production edit,
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task08-gocache /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go test -count=1 ./internal/speccheck -run 'Test(MechanicalAuditConsumesResolvedReference|MechanicalAuditJudgesTheProjectRoot|AuditRefusesOutOfGrantUnderBothCitationForms|UnresolvedReferenceIsNotASkip|AuditRefusesSelfApprovalAndRetroactiveGrants|AuditDiscoversRecordsInEveryLocation|MechanicalAuthPathsAgainstGitFixtures)'`
  exited 0. This exercised the two new external-root cases and the preserved
  out-of-grant, same-commit grant edit, unresolved-input, and record-discovery
  outcomes.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task08-gocache /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go test -count=1 ./internal/daemon -run 'TestQAMechanicalRequest(SelectsTheAuthorizedTaskCommit|CarriesTheGatePrecondition|CarriesAssignedRepairs)'`
  exited 0 after the final Daemon edit, covering request assembly and the
  target-less compatibility path.
- The focused default-root compatibility run for
  `TestAuditReadsTheResolvedReference`,
  `TestMechanicalAuthorizationReadsThePRDBoundedDeclaration`, and
  `TestCitationResolvesInExternalSpecRoot` exited 0.
- The compile-only check for `./internal/speccheck` and `./internal/daemon`
  exited 0, and `rtk git diff --check` exited 0.
- `rtk make fmt-check build GO=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt`
  exited 0 after the final code edit.
- The required `rtk make verify-incremental` first reached a sandbox refusal on
  its read-only `api.github.com` request. The unchanged authorized retry
  reached the repository test stage and exited 2: three pre-existing
  `internal/daemon` fixtures provide non-Git external Spec Roots to the Task 03
  operation resolver, and five concurrent TaskCycle cases timed out around the
  failing package run. `internal/speccheck` passed in that run. Those fixture
  files belong to another Task slice.
- The Task's declared `## Verification` commands were not run; Daemon
  Verification owns them after this handoff.
- The changed-file postflight lists only the mechanical consumer, its Daemon
  request assembler, the authorized characterization suite, and this Task
  file. The Task status change remains the Daemon's pre-existing edit.

Acceptance-criterion evidence:

- External changed-path audit:
  `TestMechanicalAuditConsumesResolvedReference/external_record_audits_the_Task_commit`
  passed with one granted authorization read and no authorization skip.
- Project-root judgment: `TestMechanicalAuditJudgesTheProjectRoot` passed where
  the project Task commit does not exist in the Spec repository. The audit read
  the external grant and reported the project-only `.golangci.yml` change as
  outside its `Makefile` bound.
- Preserved outcomes: the project-root case proves the out-of-grant refusal;
  `TestMechanicalAuditConsumesResolvedReference/genuine_no-authorization_declaration_keeps_its_skip`
  proves the presence-aware skip; the focused legacy cases preserve refusal
  when the consuming commit edits its own grant.
- Unresolved reference:
  `TestMechanicalAuditConsumesResolvedReference/unresolved_external_record_refuses`
  passed with an `unresolved` authorization read and `QA-AUTH-PATHS` finding,
  and with no authorization skip.

Follow-up:

- The non-Git external Spec Root fixtures reported by the incremental profile
  remain for their owning Task 03 contract; this Task does not change them.
