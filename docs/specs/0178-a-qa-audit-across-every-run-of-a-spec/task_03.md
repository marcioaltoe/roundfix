---
task: task_03
spec: 0178-a-qa-audit-across-every-run-of-a-spec
status: pending
type: backend
complexity: high
---

# Task 03: Staleness is measured against the Delivery Base, and stale warns

## Overview

`ResolveAuditorEvidence` in `internal/spec/auditor_evidence.go` compares the
auditing binary's build commit with the audited repository's `HEAD`. The Daemon
that runs the QA mechanical stage is built from the default branch, and the
audited head carries the candidate's commits, so in a Roundfix self-audit the
`auditor_staleness` line the Daemon seeds reads `stale` by construction, and no
code reads it. The build identity comes from the running binary's link-time
stamp; the evidence comes from the Run Worktree's Git history; the line is read
by the QA Agent and by the maintainer, and after this Task a real staleness also
reaches the Run journal as a warning an operator acts on. The warning never
stops the gate: a multi-item `roundfix deliver` queue keeps one binary while
each later item starts from a newer main, and a refusal would stop every item
after the first.

## Requirements

1. MUST change `ResolveAuditorEvidence` to
   `ResolveAuditorEvidence(ctx context.Context, repoRoot, deliveryBase string, binary app.AuditingBinary)`.
   Ancestry MUST compare the build commit, with `-dirty` removed, with
   `deliveryBase` instead of `HEAD`: equal is `AncestryNotOlder`, a strict
   ancestor is `AncestryOlder`, a non-ancestor is `AncestryNotOlder`. An empty
   `deliveryBase` MUST leave ancestry `AncestryUnknown`, never fall back to
   `HEAD`; the declared-version fallback stays as it is. `AuditorEvidence` MUST
   gain `Binary app.AuditingBinary` and `DeliveryBase string`, set to the binary
   and the base it was resolved for.
2. MUST update the call sites of `TestResolveAuditorEvidence` in
   `internal/spec/auditor_evidence_test.go`, which this signature change
   invalidates, to pass the fixture's `HEAD` commit as the base; its expected
   results stay unchanged. No `internal/**/testdata/*.txt` harness calls
   `ResolveAuditorEvidence`; the only production caller is
   `writeMechanicalQAReport`.
3. MUST word the two commit-ancestry reasons of `CompareToTree` in
   `internal/app/version.go` as
   `commit ancestry: build commit predates the delivery base` and
   `commit ancestry: build commit does not predate the delivery base`, and
   update those two expectations in `TestCompareToTree` in
   `internal/app/version_test.go`. `TestStalenessLineNeverRepeatsItsState`
   stays green and unchanged.
4. MUST add `Auditor func() app.AuditingBinary` to `daemon.Dependencies` in
   `internal/daemon/engine.go`; a nil value means `app.Auditor`. The QA step and
   `writeMechanicalQAReport` read the auditing binary only through it.
   `writeMechanicalQAReport` keeps its signature and resolves evidence with
   `qaDeliveryBase`'s base when resolved, `""` otherwise.
5. MUST make `mechanicalQAReportContent` in `internal/daemon/task_engine.go` and
   `spec.WritePreconditionRefusalReport` in `internal/spec/qa.go` write
   `evidence.Binary` when its `Version` is set and `app.Auditor()` otherwise, so
   every existing caller keeps its bytes.
6. MUST treat a `stale` state as a warning, never a refusal. When the state
   `CompareToTree` returns is `stale`, `mechanicalQAReportContent` MUST append,
   in both its seed form and its refusal form, a `## Auditor staleness warning`
   section written as a list, never a table, with the items
   `- auditor_staleness: <staleness line>` and
   `- action: rebuild roundfix from delivery base <base> and restart it before the next gate`.
   After `writeMechanicalQAReport`, the QA step MUST publish exactly one
   `runevent.KindDaemonQA` event whose payload carries
   `phase: auditor_staleness`, `auditor_staleness`, `delivery_base`, `action`
   (the same instruction) and `report`. It MUST NOT mark the mechanical result
   blocking or refused, so repository Verification and the Agent turn run as
   for `current` and `unknown`. MUST NOT add an `auditing binary age`
   Precondition Refusal or an `AuditorStalenessCheck` constant. `current` and
   `unknown` MUST publish no such event and write no such section.
7. MUST put new spec tests in `internal/spec/auditor_delivery_base_test.go`:
   - `TestAuditorEvidenceIsCurrentWhenTheBuildIsTheDeliveryBase`: build commit
     equals the base while `HEAD` carries two more commits;
   - `TestAuditorEvidenceIsStaleWhenTheBuildPredatesTheDeliveryBase`;
   - `TestAuditorEvidenceIsUnknownWithoutADeliveryBase`: empty base, build
     commit present in the repository, ancestry `AncestryUnknown`;
   - `TestAuditorEvidenceCarriesTheBinary`: the evidence carries the binary and
     the base it was resolved for.
8. MUST put new daemon tests in `internal/daemon/qa_auditor_staleness_test.go`,
   each injecting `Dependencies.Auditor` with a binary whose build commit is
   chosen per case, in a fixture repository with `refs/remotes/origin/main` and
   the `refs/remotes/origin/HEAD` symbolic ref:
   - `TestSelfAuditSeedRecordsACurrentAuditor`: a build of the Delivery Base,
     with candidate commits after it, seeds
     `auditor_staleness: "current: commit ancestry: build commit does not predate the delivery base"`;
   - `TestStaleAuditorWarnsAndTheGateProceeds`: a build of an ancestor of the
     Delivery Base produces exactly one `daemon.qa` event with phase
     `auditor_staleness` whose `delivery_base` is the base and whose `action`
     names it, a seeded report carrying the `## Auditor staleness warning`
     section and no `precondition_check`, a repository Verification that runs,
     and a fake runner that records the QA Agent call;
   - `TestStaleAuditorWarningRidesARefusedGate`: with a contradicted Spec and a
     stale auditor, the report's `precondition_check` is still
     `spec check --strict`, the report carries the warning section, and the
     warning event is published once;
   - `TestCurrentAuditorEmitsNoStalenessWarning`: a build of the Delivery Base
     publishes no `auditor_staleness` event, writes no warning section and
     reaches the QA Agent;
   - `TestUnknownAuditorEmitsNoStalenessWarning`: a build commit absent from the
     repository publishes no such event, writes no such section and reaches the
     QA Agent.
9. MUST update the Auditing Binary entry of `CONTEXT.md` to say it is the
   Daemon's binary that runs the mechanical stage, that `auditor_staleness`
   compares its build commit with the Delivery Base, and that a stale Auditing
   Binary is published as a warning (that phrase) while the gate proceeds.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A Daemon built from the Delivery Base reads `current` although the audited
      head carries the candidate's commits, and publishes no warning.
- [ ] A Daemon built from an ancestor of the Delivery Base publishes one
      `daemon.qa` warning, records it in the seeded report, and the gate still
      runs repository Verification and the Agent turn.
- [ ] Every existing report writer keeps its bytes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- creates: `internal/spec/auditor_delivery_base_test.go`
- creates: `internal/daemon/qa_auditor_staleness_test.go`
- interface: `internal/spec/auditor_evidence.go`
- interface: `internal/spec/auditor_evidence_test.go`
- interface: `internal/spec/qa.go`
- interface: `internal/app/version.go`
- interface: `internal/app/version_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAuditorEvidenceIsCurrentWhenTheBuildIsTheDeliveryBase|TestAuditorEvidenceIsStaleWhenTheBuildPredatesTheDeliveryBase|TestAuditorEvidenceIsUnknownWithoutADeliveryBase|TestAuditorEvidenceCarriesTheBinary|TestResolveAuditorEvidence|TestCompareToTree|TestStalenessLineNeverRepeatsItsState|TestSelfAuditSeedRecordsACurrentAuditor|TestStaleAuditorWarnsAndTheGateProceeds|TestStaleAuditorWarningRidesARefusedGate|TestCurrentAuditorEmitsNoStalenessWarning|TestUnknownAuditorEmitsNoStalenessWarning)$" ./internal/spec ./internal/app ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAuditorEvidenceIsCurrentWhenTheBuildIsTheDeliveryBase TestAuditorEvidenceIsStaleWhenTheBuildPredatesTheDeliveryBase TestAuditorEvidenceIsUnknownWithoutADeliveryBase TestAuditorEvidenceCarriesTheBinary TestResolveAuditorEvidence TestCompareToTree TestStalenessLineNeverRepeatsItsState TestSelfAuditSeedRecordsACurrentAuditor TestStaleAuditorWarnsAndTheGateProceeds TestStaleAuditorWarningRidesARefusedGate TestCurrentAuditorEmitsNoStalenessWarning TestUnknownAuditorEmitsNoStalenessWarning; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "a stale Auditing Binary is published as a warning"` — expected: exit 0; before this Task none of the nine new tests exists and the glossary does not say a stale Auditing Binary is published as a warning, so the command fails.

## References

- [_techspec.md](_techspec.md) — Staleness against the Delivery Base
