---
task: task_05
spec: 0119-spec-contained-authorization
status: completed
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

1. MUST resolve the operative grant from the ancestry of the delivery target,
   not merely from some earlier commit on the consuming branch, and MUST take
   that target revision as an input rather than assuming it. A grant commit and
   its consuming commit on the same branch collapse into one commit under
   squash delivery, so an earlier-sibling grant is self-approval wearing two
   commits.
2. MUST record which record path and revision the audit read, and MUST persist
   that provenance in the mechanical result and its written report, so the
   grant object a passing audit relied on stays retraceable after the
   in-memory result is gone.
3. MUST refuse a consuming commit that creates or widens the grant it depends
   on, and MUST NOT let a later amendment authorize an earlier change
   retroactively.
4. MUST accept an earlier bounded amendment as authority for work that comes
   after it, when that amendment is already in the delivery target.
5. MUST discover operative records in active Specs, in archived Specs, and in
   the preserved legacy location, and MUST NOT narrow the governed set the audit
   judges against.
6. MUST keep the existing changed-path refusals and their reported tokens
   working for the conditions they already own, including a governed change
   outside the bounded set and a grant edited in the commit that consumes it.
7. MUST report an unavailable revision or unreadable record as an unresolved
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
- [ ] A grant commit followed by its consuming commit on the same branch, with
      neither yet in the delivery target, refuses; the same grant refuses to
      authorize its sibling even though it is an earlier ancestor.
- [ ] A commit whose grant is already in the delivery target passes, and both
      the audit result and its written report carry the record path and
      revision the audit read.
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

## Result

Implemented an ancestor-bound changed-path audit. Each consuming Task commit is
now judged against the authorization record at the merge base of the supplied
delivery target and the consuming commit's parent. The Daemon supplies the
Run-start head as that delivery target and selects every commit that touches a
Governed Path, including commits whose governed paths are all outside the
grant.

The audit reads committed record bytes through the shared typed authorization
reader, carries the asking Spec into that read, and discovers the same
Spec-owned record across active and archived locations. The preserved legacy
location retains its bounded adapter for dated, frontmatter-free records with
no declared consumer; a legacy record that declares another consumer still
refuses. Each audit read persists its Task, outcome, record path, resolved
revision, and refusal or unresolved reason in `MechanicalResult`, and
`WriteMechanicalResult` writes that provenance into the QA Report.

Acceptance evidence:

1. `TestAuditRefusesSelfApprovalAndRetroactiveGrants/same_commit_creates_its_grant`
   and `/same_commit_widens_its_grant` create real commits that add or widen
   their own authorization. Both emit `QA-AUTH-PATHS` naming the record; the
   create case also asserts the existing `changes authorization grant` detail.
2. `TestAuditRefusesSelfApprovalAndRetroactiveGrants/earlier_sibling_grant_is_absent_from_target_ancestry`
   commits a grant and consumer on one branch while the supplied delivery
   target remains at their base. The audit reads the base revision and refuses
   the sibling grant rather than accepting its commit order on that branch.
3. `TestAuditReadsTheAuthorizingAncestor/grant_already_in_delivery_target`
   places the grant in the delivery target before the consuming commit. The
   audit reports `granted`, retains the exact record path and target SHA, and
   the written mechanical report contains both values.
4. `TestAuditRefusesSelfApprovalAndRetroactiveGrants/later_amendment_authorizes_only_later_work`
   branches an early consumer from the narrow grant, then widens the target
   grant. The early audit reads the narrow revision and refuses; a later commit
   based on the amendment reads the amendment revision and produces no
   `QA-AUTH-PATHS` finding.
5. `TestAuditDiscoversRecordsInEveryLocation` resolves grants from
   `docs/specs/`, `docs/history/specs/`, and
   `docs/workflow/authorizations/`. The archived case begins from the active
   record path and records the discovered archived path. Its legacy cases
   preserve a frontmatter-free grant while refusing a record that explicitly
   names another consumer.
6. `TestAuditJudgesTheGrant/governed_path_outside_the_grant_is_refused_by_name`
   preserves the existing `QA-AUTH-PATHS` token and escaped path detail.
   `TestQAMechanicalRequestSelectsTheAuthorizedTaskCommit/ungranted_governed_path_blocks`
   proves the Daemon selects a commit that changes only `.golangci.yml`, even
   though the grant bounds only `Makefile`.
7. `TestAuditDiscoversRecordsInEveryLocation/unavailable_delivery_target_is_unresolved`
   supplies an absent forty-character revision and observes an
   `unresolved` audit read with reason `unavailable_revision` plus a blocking
   `QA-AUTH-PATHS` finding. The existing missing-record case now records
   `unreadable_record` instead of a presence skip.

Focused checks:

- Pre-change, `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/speccheck -run '^TestAuditReadsTheAuthorizingAncestor/grant_already_in_delivery_target$' -count=1`
  failed to compile because `DeliveryTargetRevision` and
  `AuthorizationReads` did not exist. An earlier attempt using the shared Go
  cache was denied by the sandbox and is not implementation evidence.
- After the final implementation edit, `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/speccheck -run '^(TestAuditReadsTheAuthorizingAncestor|TestAuditRefusesSelfApprovalAndRetroactiveGrants|TestAuditDiscoversRecordsInEveryLocation|TestMechanicalAuthPaths|TestMechanicalAuthPathsAcceptsDeclaredRegenerationOutput|TestMechanicalAuthPathsStillRefusesAnUndeclaredPath|TestMechanicalAuthPathsRefusesInvalidRegenerationDeclaration|TestAuditJudgesTheGrant|TestMechanicalReportsAllFindings|TestMaterializeMechanicalResult)(/.*)?$' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/daemon -run '^(TestQAMechanicalRequestSelectsTheAuthorizedTaskCommit|TestMechanicalStageSeedsReportBeforeAgentSession|TestWriteMechanicalQAReportWritesThePreconditionRefusal|TestWriteMechanicalQAReportPreservesSameDayNamingAndPriorReport|TestWriteMechanicalQAReportRecordsTheRefusal)(/.*)?$' -count=1`
  passed after the final Daemon and report changes.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test -tags repocontract ./internal/speccheck -run '^TestCleanupHistoricalGrantEvidence(/.*)?$' -count=1`
  passed, preserving the recovered frontmatter-free legacy grant behavior.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go vet ./internal/speccheck ./internal/daemon`
  passed.

The Task's declared `## Verification` commands were not run; Daemon
Verification remains pending.
