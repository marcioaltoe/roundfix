---
task: task_05
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 05: Plan root-instruction preservation

## Overview

Extend the read-only planner through greenfield, preservation, Source
Baseline, and Readoption decisions. The result names every root backup and
repository-rules disposition before any instruction carrier can change.

## Requirements

1. MUST offer explicit greenfield and preservation modes for unconfigured or
   incompatible repositories.
2. MUST plan immutable content-addressed root backups and preserve a safe alias
   target exactly once.
3. MUST import no existing rule in greenfield mode and require one validated
   disposition for every root rule in preservation mode.
4. MUST retain every maintained Source Baseline, Readoption, Upgrade Retention,
   Decision Document, and Repository-Specific Normative Rules contract.
5. MUST leave nested carriers unchanged and expose their conflicts only as
   warnings.
6. MUST emit complete Decision Document skeletons with the required schema
   fields and stable next actions.

## Subtasks

- [x] Port Source Baseline, Readoption, and retention decision models.
- [x] Resolve greenfield and preservation planning states.
- [x] Derive content-addressed root backup identities.
- [x] Render consolidated manual classification inputs and warnings.
- [x] Add parity and refusal tests for every preservation state.

## Acceptance Criteria

- [x] Greenfield plans backups but moves zero existing rules into repository rules.
- [x] Preservation cannot become ready with an unclassified root rule.
- [x] Backup paths use the full raw-content SHA-256 identity and reject collisions.
- [x] Safe aliases back up one target; unsafe aliases remain blocking.
- [x] Decision skeletons pass the strict runtime parser without manual schema repair.
- [x] Exact maintained Readoption and retention fixtures match the Python corpus.

## Context

- instruction: `docs/adr/0058-upgrade-retention-requires-clause-level-accounting.md`
- instruction: `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_baseline.py`
- interface: `docs/agents/specific-repository.md`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestGreenfieldPlan|TestPreservationPlan|TestRootBackupIdentity|TestDecisionDocumentSkeleton|TestReadoptionCompatibility'` — expected: preservation modes, backups, strict decision inputs, and parity cases pass.
- `rtk go test -count=1 ./internal/baseline -run 'TestNestedCarrierWarning|TestPreservationRequiresEveryDisposition'` — expected: nested files stay out of mutation plans and incomplete classifications remain unresolved.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 3 and 5; Core Features 5–9 and 19; User Experience.
- `_techspec.md` → Data Models: RepositorySnapshot and backups; Testing Approach; Build Order 3–4.
- ADR-0058 → clause-level retention.
- ADR-0070 → root-only automatic preservation.

## Result

Roundfix now resolves root-instruction preservation as a read-only planner
slice over the immutable repository inspection. Greenfield and preservation
are explicit modes. Both derive one exclusive, content-addressed backup per
trusted root source using the complete raw-content SHA-256 identity; an
existing backup is accepted only when its regular-file bytes match, and a
collision blocks planning. Safe aliases share one source backup, preferring a
root regular target when one exists. Root-unsafe carriers remain blocking,
while nested regular, opaque, and unsafe carriers remain outside backup and
classification plans and surface only the existing conflict warning.

Preservation inventories byte-evidenced root Source Baseline Entries, including
managed-block boundaries, and stays `action_required` until every entry has one
current disposition. The emitted consolidated Decision Document skeleton is
strict-parser-valid: unmarked rule bytes are proposed for the canonical
Repository-Specific Normative Rules destination, while structural or old
managed evidence is proposed as an individually reasoned rejection. Accepted
Repository-Specific Normative Rules bytes remain exact, canonical-base64,
digest-bound, unmarked decision input.

The Go models now load and verify the maintained Source Baseline identity,
60-entry manifest, 51-entry source-accounting ledger, Readoption Decision
Document union, and both clause-complete Upgrade Retention Contracts from the
embedded catalog. Duplicate JSON keys, stale source identities and entry
digests, incomplete or duplicate dispositions, unsafe typed destinations,
non-canonical bytes, managed markers, and backup collisions fail closed.

Acceptance evidence:

- `TestGreenfieldPlanBacksUpWithoutImport` produced the full-SHA
  `AGENTS.<64-hex>.md` backup and no dispositions or Repository-Specific
  Normative Rules bytes.
- `TestPreservationRequiresEveryDisposition` kept an omitted root entry
  unresolved; `TestPreservationPlanAcceptsCompleteDecisionDocument` became
  ready only with one current disposition per entry and retained the exact
  proposed bytes.
- `TestRootBackupIdentityRejectsCollisions` blocked mismatched bytes at the
  digest-derived path; `TestRootBackupIdentitySafeAliasesBackUpTargetOnce`
  reduced two safe aliases to one trusted source backup.
- Existing root unsafe-alias tests still block, while
  `TestNestedCarrierWarningLeavesNestedSourcesOutOfPreservation` and
  `TestNestedCarrierWarningUnsafeAliasRemainsNonBlocking` prove nested paths
  never enter the mutation slice.
- `TestDecisionDocumentSkeletonPassesStrictParser` round-tripped an emitted
  skeleton containing unmarked and old managed content through the strict
  parser; the duplicate-key refusal test rejected malformed input.
- `TestReadoptionCompatibilityMaintainedFixture` parsed the exact 19-entry
  Python parity fixture and verified the embedded Source Baseline, source
  accounting, and every maintained Upgrade Retention transition.

Verification:

- Pre-change focused tests failed to compile because the preservation planner,
  backup ledger, Decision Document parser, Readoption models, and typed Source
  Baseline and retention accessors did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache go test -count=1
  ./internal/baseline ./internal/cli -run
  'TestGreenfieldPlan|TestPreservationPlan|TestRootBackupIdentity|TestDecisionDocumentSkeleton|TestReadoptionCompatibility'`:
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache go test -count=1
  ./internal/baseline -run
  'TestNestedCarrierWarning|TestPreservationRequiresEveryDisposition'`:
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache go test -count=1
  ./internal/baseline ./internal/cli`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache go vet
  ./internal/baseline ./internal/cli`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache make verify`: passed,
  including the repository tests, format check, shipped-skill check, and
  build.

The isolated `GOCACHE` keeps build artifacts inside the Task Worktree sandbox.
The Daemon remains responsible for the task file's verbatim authoritative
Verification commands. No commit, push, or pull request was created.
