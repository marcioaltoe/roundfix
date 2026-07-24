---
task: task_07
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 07: Emit portable Baseline Plans

## Overview

Complete the non-interactive planning slice by turning equivalent repository
state and decisions into one portable, digest-bound JSON artifact. Humans and
automations can inspect the same canonical ledger through a concise file-level
projection.

## Requirements

1. MUST implement strict `roundfix/baseline-plan/v1` and
   `roundfix/baseline-result/v1` serializers with deterministic ordering,
   duplicate-key rejection, unknown-field rejection, and stable exit
   categories.
2. MUST include repository, catalog, profile, decision, retention, preimage,
   postimage, warning, Setup Manifest, and ordered managed-entry evidence.
3. MUST derive one `fileChanges` row per affected path from the canonical
   ledger and reject any projection mismatch.
4. MUST compute the Plan Digest from the approved canonical payload and exact
   postimages without volatile ACP metadata or absolute paths.
5. MUST allow `roundfix baseline plan` to emit text or a complete portable JSON
   artifact and never prompt.
6. MUST reproduce exact maintained planned bytes and digests and remain
   idempotent for equivalent snapshots and decisions.

## Subtasks

- [x] Port Decision Plan, rendering, retention, and Setup Manifest generation.
- [x] Define strict plan and result document codecs.
- [x] Produce the canonical managed-entry ledger and derived file projection.
- [x] Implement repository-relative postimages and Plan Digest calculation.
- [x] Expose complete text/JSON planning results and stable failures.
- [x] Add cross-clone, determinism, and compatibility tests.

## Acceptance Criteria

- [x] Equivalent human-normalized and file-based inputs produce identical Plan Digests.
- [x] JSON plans contain no absolute checkout path or hidden pending-state reference.
- [x] Another clone with matching identity and bounded preimages accepts the same plan document.
- [x] `fileChanges` has one row per path and always matches the canonical ledger.
- [x] Missing decisions produce result-schema next actions and exit 3 without a partial plan.
- [x] Exact maintained fixtures reproduce planned bytes, manifests, ledgers, and digests.
- [x] Planning performs zero repository mutations.

## Context

- instruction: `docs/adr/0068-baseline-command-uses-one-confirmation-gated-workflow.md`
- instruction: `docs/adr/0071-baseline-plans-are-portable-and-preimage-bound.md`
- interface: `internal/releaseplan/plan.go`
- interface: `internal/cli/releaseplan_command.go`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestPlanDocument|TestBaselinePlanCommand|TestFileChangesProjection|TestPlanDigest|TestCrossClonePlan|TestPlanDeterminism'` — expected: schemas, output, digest, portability, and no-write cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline plan --help` — expected: help exposes the approved profile, decision, decision-file, repository, and format flags without interactive options.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 2, 4, 7–8; Core Features 10, 12, 15–16, 18–20.
- `_techspec.md` → Data Models: PlanDocument and Result; API Contracts: Automation; Build Order 4 and 7.
- ADR-0068 → explicit automation planning stage.
- ADR-0071 → portable plan and bounded-preimage contract.

## Result

Roundfix now emits a complete non-interactive
`roundfix/baseline-plan/v1` document from one clone-stable repository identity,
the embedded catalog identity, one resolved Baseline Profile, normalized
decisions, maintained retention accounting, complete bounded preimages, exact
repository-relative postimages, warnings, the current
`setup-context-driven/manifest/0.0.1` Setup Manifest, and an ordered canonical
managed-entry ledger. The concise `fileChanges` surface is derived from that
ledger and the preimage/postimage sets; strict parsing rejects any projection
mismatch.

The plan and `roundfix/baseline-result/v1` codecs reject duplicate keys at
every JSON nesting level, unknown fields, trailing values, invalid collection
shapes, identity drift, non-portable paths, unordered or duplicate evidence,
manifest/ledger disagreement, postimage byte mismatch, and Plan Digest
mismatch. Results carry stable operation state, refusal category,
`nextAction`, explicit `verifiedPostimages`, warnings, and recommendations
arrays. Missing decisions and unsupported manifest identities return no
partial Plan Document.

Planning renders the maintained managed blocks and supporting guides in
Decision Plan order, preserves the exact ownership-marker byte contract,
produces the profile-bound 0.0.1 Setup Manifest without volatile confirmation
timestamps, accounts for maintained Upgrade Retention and Readoption
dispositions, and derives the digest from the canonical payload plus exact
postimage bytes while excluding only its own digest and the reproducible
file-level projection.

Acceptance evidence:

- `TestBaselinePlanCommandEmitsPortableJSONAndNormalizesDecisionFiles`
  produced the same Plan Digest from equivalent inline and strict
  Decision-Document inputs, emitted no absolute checkout path, and observed no
  repository mutation.
- `TestPlanDeterminismAndNoMutation` produced byte-identical plan JSON and the
  same digest twice from one snapshot, rejected checkout-path leakage, and
  compared the complete repository tree before and after planning.
- `TestCrossClonePlanAcceptsMatchingIdentityAndPreimages` accepted the same
  parsed plan in another clone with matching lineage and bounded preimages,
  then rejected one changed consulted preimage as stale.
- `TestFileChangesProjectionRejectsMismatch` proved the strict parser
  recomputes one affected-path projection from the canonical ledger and
  rejects invented managed-entry evidence.
- `TestPlanDocumentMissingDecisionsReturnsResultWithoutPartialPlan` and the
  CLI preflight tests returned `roundfix/baseline-result/v1`,
  `action_required`, a decision next action, explicit empty verified
  postimages, exit 3, and no partial plan.
- `TestPlanDeterminismMatchesMaintainedManagedEntryFixture` reproduced the
  maintained Go CLI/TUI managed-entry order, artifact metadata and digests,
  and every exact non-manifest planned byte identity. Strict manifest
  validation independently recomputes its resolved modules, decisions,
  artifact order, templates, versions, and content digests before accepting
  the plan.
- `TestPlanDocumentIncludesMaintainedUpgradeRetention` reproduced the complete
  ordered maintained transition ledger, including each source clause,
  enforcement, disposition, targets, and reason;
  `TestPlanDocumentRejectsUnknownManifestRetentionWithoutPartialPlan` kept an
  unknown identity at classification-required without emitting a partial
  plan.
- `TestPlanDigestBindsExactPostimagesAndIgnoresProjection` proved exact
  postimage changes invalidate the digest while changes to the reproducible
  `fileChanges` projection do not create a second authority.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache go test -count=1
  ./internal/baseline ./internal/cli -run
  'TestPlanDocument|TestBaselinePlanCommand|TestFileChangesProjection|TestPlanDigest|TestCrossClonePlan|TestPlanDeterminism'`:
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache go run
  -buildvcs=false ./cmd/roundfix baseline plan --help`: passed; help exposes
  profile, repeatable decision and decision-file, repository, and text/JSON
  flags, the stable exit categories, and no interactive option.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache go test -count=1
  ./internal/baseline ./internal/cli`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache go vet
  ./internal/baseline ./internal/cli`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-cache make verify`: passed.
- `rtk git diff --check`: passed.

The isolated `GOCACHE` keeps build artifacts inside the Task Worktree sandbox.
The Daemon remains responsible for the task file's verbatim authoritative
Verification commands. No other Task file or Task Graph manifest was edited,
and no commit, push, or pull request was created.
