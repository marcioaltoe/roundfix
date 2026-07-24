---
task: task_06
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 06: Resolve profile alignment decisions

## Overview

Turn repository evidence and one selected Baseline Profile into explicit
blocking or advisory decisions. The planner provides useful HTTP, PostgreSQL,
and Verification evidence without inventing repository policy.

## Requirements

1. MUST resolve exactly one valid built-in or repository-owned Baseline
   Profile and compare repository evidence with its requirements.
2. MUST block on unresolved required divergence and keep advisory divergence
   visible without treating it as policy.
3. MUST emit bounded HTTP route candidates, observed methods and scopes, and
   source digests without assigning contract mode, owner, or rationale.
4. MUST distinguish PostgreSQL implementation evidence from a missing accepted
   repository contract and name accepted contract paths when required.
5. MUST label commands repository-executable only after local declaration
   validation; portable roles and profile expectations remain distinct.
6. MUST remain local, read-only, network-free, and command-execution-free.

## Subtasks

- [x] Port Repository Capability evaluation and evidence ranking.
- [x] Resolve profile alignment and divergence decision states.
- [x] Add HTTP route-candidate and source-digest projection.
- [x] Separate implementation, contract, portable-role, and executable-command evidence.
- [x] Add finding regression and profile parity tests.

## Acceptance Criteria

- [x] Required divergence prevents a ready Plan until explicitly resolved.
- [x] Advisory divergence never blocks and never becomes inferred policy.
- [x] HTTP candidates contain facts but no inferred Normative Clause.
- [x] PostgreSQL diagnostics report found implementation evidence separately from contract absence.
- [x] A nonexistent formatter or Verification script is never labeled executable.
- [x] Equivalent evidence and answers produce equivalent normalized decisions across interaction modes.

## Context

- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_capabilities.py`
- interface: `.agents/skills/setup-context-driven/assets/profiles/standard-typescript-monorepo.json`
- interface: `docs/findings/2026-07-23-setup-context-driven-adoption-process-improvements.md`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestProfileAlignment|TestRequiredDivergence|TestHTTPRouteCandidates|TestPostgreSQLEvidence|TestExecutableVerificationCommand'` — expected: alignment, evidence classification, and finding regressions pass.
- `rtk go test -count=1 ./internal/baseline -run TestCapabilityAuditNoExecution` — expected: audit uses declared local evidence and invokes no repository or network command.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 4 and 8; Core Features 3–4, 10, 14, 17–18.
- `_techspec.md` → Data Models: Catalog and RepositorySnapshot; Testing Approach: Fluxus assertions; Build Order 3–4.
- ADR-0063 → repository-owned HTTP contract policy.

## Result

Roundfix now resolves one exact built-in or repository-owned Baseline Profile
into a deterministic, read-only alignment result. Explicit decision answers
are validated against the selected profile and normalized by ID. Missing
required answers or Repository Capabilities produce blocking divergences and
`action_required`; recommended and optional gaps remain visible advisory
divergences and never enter the normalized decision set.

Repository Capability evaluation ports the maintained evidence-strength order
(`none < declared < discovered < verified`) over bounded declared files,
installed Repository Skills, and PATH discovery that inspects executable file
metadata without launching a process. The audit is context-aware, root
confined, size/count bounded, network-free, and performs no repository writes.

For the Standard TypeScript Monorepo Profile, the result now includes:

- bounded HTTP route candidates with only observed scope, methods, source path,
  and `sha256:` source identity; the model has no mode, owner, reason,
  rationale, or Normative Clause field;
- PostgreSQL implementation evidence from local driver, adapter,
  configuration, or compose declarations, kept separate from the accepted
  repository contract paths and contract evidence;
- portable Verification role expectations kept distinct from the selected
  repository Verification command, with `repositoryExecutable: true` only
  after an exact root `package.json` script or Make target declaration is
  validated locally.

Acceptance evidence:

- `TestRequiredDivergencePreventsReadyPlan` held the result at
  `action_required` while Better Auth evidence was absent and became ready only
  after the declared evidence was added.
- `TestProfileAlignmentAdvisoryDivergenceNeverBlocksOrInfersPolicy` kept a
  missing recommended Firecrawl capability advisory and absent from normalized
  decisions.
- `TestHTTPRouteCandidatesContainFactsWithoutNormativeClause` projected GET,
  POST, and unclassified route-scope facts with one source digest and rejected
  every policy-bearing field from the serialized shape.
- `TestPostgreSQLEvidenceSeparatesImplementationAndContract` reported package
  and compose implementation evidence while blocking on the absent accepted
  contract, named all accepted paths, then became ready after `DATABASE.md`
  supplied the repository contract.
- `TestExecutableVerificationCommandRequiresLocalDeclaration` kept nonexistent
  formatter, workspace, and selected Verification commands non-executable;
  only an exact `Makefile` target changed the selected gate to executable.
- `TestProfileAlignmentEquivalentNormalizedDecisions` produced byte-identical
  JSON for equivalent answers supplied in opposite orders.
- `TestProfileAlignmentResolvesExactlyOneProfile` covered built-in and
  repository-owned profiles and rejected missing or unknown selection;
  `TestProfileAlignmentCapabilityEvidenceRanking` covered satisfied,
  insufficient, blocking, and advisory evidence ranks.
- `TestCapabilityAuditNoExecution` placed executable trap files on PATH,
  observed zero execution, and compared the complete repository tree before
  and after the audit with no change.

Verification:

- Pre-change focused tests failed to compile because the profile-alignment
  engine and result models did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test -count=1
  ./internal/baseline ./internal/cli -run
  'TestProfileAlignment|TestRequiredDivergence|TestHTTPRouteCandidates|TestPostgreSQLEvidence|TestExecutableVerificationCommand'`:
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test -count=1
  ./internal/baseline -run TestCapabilityAuditNoExecution`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test -count=1
  ./internal/baseline ./internal/cli`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go vet
  ./internal/baseline ./internal/cli`: passed.
- One earlier full package run observed the pre-existing
  `TestBaselinePlanPreflightJSONActionRequired` tree snapshot racing with a
  disappearing Git `maintenance.lock`. The exact test immediately passed in
  isolation, no Task 06 code touches Git or that helper, and the fresh complete
  package run passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache make verify`: passed,
  including 1,812 Go tests in 21 packages, both 256-test setup-skill suites,
  catalog asset validation, shipped-skill checks, and the binary build.

The isolated `GOCACHE` keeps build artifacts inside the Task Worktree sandbox.
The Daemon remains responsible for the task file's verbatim authoritative
Verification commands. No other Task file or Task Graph manifest was edited,
and no commit, push, or pull request was created.
