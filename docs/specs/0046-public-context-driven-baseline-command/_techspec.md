---
spec: 0046-public-context-driven-baseline-command
prd: _prd.md
created: 2026-07-23
---

# Public Context-Driven Baseline Command — Technical Spec

## Executive Summary

Roundfix will replace the Python engine distributed in the
`setup-context-driven` skill with a deterministic Go engine behind the public
`roundfix baseline` command family. Humans use one linear,
confirmation-gated state machine; automations exchange a portable Baseline
Plan between explicit `plan` and `apply` subcommands. The primary trade-off is
to carry exact postimages and bounded preimage identities in a larger JSON
artifact rather than persist hidden pending state, because reviewability,
cross-process reproducibility, and safe CI handoff matter more than a smaller
payload. Semantic ACP analysis remains an optional, sealed proposal boundary;
all authority stays in deterministic validation, human decisions, the approved
plan digest, and the recoverable Go apply transaction.

## System Architecture

The command entry point remains `cmd/roundfix/main.go`, with stdlib
`flag.FlagSet` dispatch in `internal/cli`. A new cohesive
`internal/baseline` package owns the embedded catalog, Baseline Profiles,
repository inventory, Decision Plan, Change Plan, rendering, digests,
verification, and transaction protocol. It is deterministic, network-free,
and independent of `internal/cli`, `internal/agent`, `internal/config`,
`internal/store`, and the installed skill.

Canonical assets move to `internal/baseline/assets/**` and compile through
`go:embed`. Profiles load only from
`.roundfix/baseline/profiles/<id>.json`, compose embedded entries, and cannot
add code or assets. The skill keeps only tested CLI guidance after cutover.

`internal/baselineacp` adapts the Baseline proposal interface to the existing
Exact Agent Selection Proof and ACPX session plumbing in `internal/agent`.
It never uses the Batch-oriented `agent.Runner.Run`, Run Database, Run Event
Journal, Agent log, or configured Agent Work Categories. The adapter runs each
attempt in a private empty directory and exposes no checkout path.

```text
roundfix baseline ─┬─ linear prompt driver ─────────────┐
                   ├─ plan/apply JSON drivers ──────────┤
                   ├─ profile/skills/assets operations  │
                   ▼                                    ▼
              internal/baseline ◄──────────── internal/baselineacp
                │       │       │                    │
        embedded catalog│ repository adapter         └─ internal/agent → ACPX
                        └─ transaction + verifier
```

A Git worktree with at least one commit is required so the plan can derive a
clone-stable lineage identity. Dirty worktrees, detached HEAD, and missing
upstream are allowed. Repository safety comes from bounded preimages rather
than global cleanliness.

## Implementation Design

### Interfaces

The workflow keeps interaction outside the engine and normalizes human and
automation answers into the same decision input:

```go
type Workflow interface {
	Inspect(context.Context, Repository) (State, error)
	Decide(State, DecisionInput) (State, error)
	Plan(context.Context, State) (PlanDocument, error)
	Apply(context.Context, Repository, PlanDocument, Digest) (Result, error)
}
```

Repository access is root-anchored, no-follow, and records every consulted
entry plus each possible mutation target:

```go
type Repository interface {
	Identity(context.Context) (RepositoryIdentity, error)
	Snapshot(context.Context, InventoryRequest) (RepositorySnapshot, error)
	BeginTransaction(context.Context, PlanDocument) (Transaction, error)
}

type Transaction interface {
	Apply(context.Context) (VerificationEvidence, error)
	Rollback(context.Context) error
	Close() error
}
```

Semantic work is proposal-only and receives canonical bytes rather than a
filesystem capability:

```go
type ProposalAnalyzer interface {
	Classify(context.Context, AnalysisSnapshot) (ClassificationProposal, error)
	Revise(context.Context, RevisionSnapshot) (RevisionProposal, error)
}
```

The CLI injects the prompt reader, output writers, Git runner, clock, catalog,
repository adapter, and proposal analyzer. Production code receives no
test-only hooks.

### Data Models

`Catalog` contains immutable modules, decisions, templates, Repository
Capabilities, retention transitions, and built-in Baseline Profiles.
`CatalogDigest` is computed from canonical, path-ordered embedded bytes.
`CustomProfile` contains an ID, catalog schema, selected entry IDs, values, and
its own digest. Unknown IDs, duplicate keys, profile composition, remote
references, and custom assets are invalid.

`RepositoryIdentity` hashes the Git object format and sorted root commit IDs
reachable from `HEAD`. It is clone-stable but changes with unrelated or
rewritten root history. Plans contain only normalized relative paths.

`RepositorySnapshot` is the complete bounded preimage. Every entry records
path, existence, regular-file/symlink kind, mode, safe link target when
applicable, byte length, and `sha256:<hex>` content identity. It contains all
bytes consulted by audit or decisions and every target that apply may create,
replace, or remove. Whole-worktree state is deliberately excluded. Each path
component is revalidated with `Lstat` before apply to reject escaping,
external, cyclic, unreadable, or special-file carriers.

Root regular instruction carriers receive immutable adjacent backups named
`AGENTS.<64-lowercase-sha256>.md` or
`CLAUDE.<64-lowercase-sha256>.md`. Creation is exclusive; an existing path is
accepted only when its bytes match its name. A safe root symlink records its
in-repository target and backs up the target bytes once. Nested carriers are
inventory evidence and warning sources only.

`PlanDocument` is strict `roundfix/baseline-plan/v1` JSON containing repository
and catalog identity, selected profile, normalized decisions, retention
accounting, warnings, preimages, ordered `ManagedEntry` records, exact
postimages encoded without platform-dependent paths, and `planDigest`.
`fileChanges` is a serialized projection grouped by path; strict parsing
recomputes it from the canonical ledger and postimages, so it cannot become a
second authority. The digest covers the canonical payload except its own
field and the reproducible projection.

`Result` uses `roundfix/baseline-result/v1` and reports state, verified
postimages, recommendations, warnings, refusal category, and `nextAction`.
Volatile ACP attempt metadata may appear here but never affects a Plan Digest.
Raw prompts, model messages, thought, and invalid payloads are never persisted.
Go writes the current `setup-context-driven/manifest/0.0.1` Setup Manifest,
while its readers retain every manifest and transition currently accepted by
the Python runtime until a separately approved deprecation.

The apply transaction uses a Git-private lock and recovery directory. It
journals exact preimages, stages and fsyncs postimages, rechecks the bounded
preimage, then replaces paths deterministically. It verifies postimages and
carrier relationships; failure rolls back in reverse order and fsyncs affected
directories. Another apply must recover or report an interrupted transaction.

### API Contracts

Human workflow:

```text
roundfix baseline [--repo <path>] [--format text|json]
```

It requires interactive input, uses numbered linear prompts, shows a
consolidated review, and asks once for the final digest-bound apply. Rejecting
the plan returns to a selected decision area. Without an interactive terminal,
it exits with a structured next action directing the caller to `baseline
plan`; it never guesses or prompts through redirected input.

Decision and preservation prompts expose one catalog or Setup Manifest default
and accept Enter as confirmation. Still-valid manifest values win when a
changed Profile Digest forces adoption. Catalog defaults are English,
`rtk make verify`, Post-only, spec artifacts enabled, single-context, external
triage disabled, autonomous work enabled, `codex gpt-5.6-sol`,
`claude fable high`, Secondbrain enabled, and the Repository-Specific Normative
Rules carrier permitted. The carrier and its root pointer are emitted only when
non-empty rules exist. Existing root instructions default to Preservation; an empty inventory
defaults to Greenfield. A recoverable existing profile is the profile default.
Classification, plan approval, and apply require a non-empty explicit choice.

Automation:

```text
roundfix baseline plan --repo <path> [--profile <id>] \
  [--decision <id=value> ...] [--decision-file <path> ...] \
  --format text|json
roundfix baseline apply --repo <path> --plan <file> \
  --confirm-plan <digest> --format text|json
```

JSON planning writes one complete `roundfix/baseline-plan/v1` document to
stdout. Apply parses with duplicate-key and unknown-field rejection, verifies
the document digest, resolves the same Git lineage, validates every preimage,
and applies exactly the supplied postimages. It never silently recalculates or
substitutes a plan. Text output is a human projection; only JSON is a portable
apply input. Diagnostics and progress go to stderr.

Profile and maintenance operations are:

```text
roundfix baseline profile init --id <id> [--from <built-in-id>]
roundfix baseline profile show <id> [--format text|json]
roundfix baseline profile validate [<id>|<path>] [--format text|json]
roundfix baseline skills restore --repo <path> --profile <id> \
  [--skill <name> ...] [--source-dir <path>] \
  [--confirm-plan <digest>] [--format text|json]
roundfix baseline assets sync --source-dir <path> \
  [--check] [--format text|json]
```

`profile init` writes a repository-owned declarative profile, while `show` and
`validate` resolve it against the embedded catalog. `skills restore` and
`assets sync` are Go ports of the maintained Python responsibilities but are
not states in adoption or update.

Exit codes preserve the CLI contract: `0` success; `1` execution, apply,
verification, output, or incomplete-rollback failure; `2` usage, invalid
schema, unsafe repository, or Preflight Validation failure; `3` missing
decision, manual classification, renewed approval, stale preimage, or another
actionable unresolved state; and `130` cancellation. JSON mode returns a
`baseline-result/v1` failure on stdout when arguments were parsed far enough
to select that contract, with the diagnostic on stderr.

ACP classification and revision use separate strict schemas. Preflight eagerly
proves Codex `gpt-5.6-sol`/`xhigh` and Codex `gpt-5.5`/`xhigh`. Each actual
attempt receives byte-identical canonical input, uses a fresh session,
`--deny-all`, an empty allowed-tools set, no terminal, one turn, and a
two-minute deadline. A snapshot is limited to 2 MiB and 256 entries; output is
limited to 512 KiB. Oversized or unprovable work moves to manual review.
Tools, extra prose, unknown or missing IDs, duplicate dispositions, unsupported
destinations, digest mismatch, timeout, or invalid JSON discard the attempt.
Parent cancellation does not trigger fallback. Failure to prove a session
closed prevents the next attempt.

## Coverage Map

- Goal 1 → CLI drivers, Workflow, embedded Catalog, Go maintenance operations.
- Goal 2 → normalized DecisionInput, PlanDocument, deterministic digests.
- Goal 3 → inventory, capability audit, proposal validation, final QA corpus.
- Goal 4 → embedded assets, thin skill, documentation contract tests.
- Goal 5 → one Workflow, recoverable Transaction, parity matrix.
- Goal 6 → dedicated documentation delivery step and public examples.
- Story 1 → root `roundfix baseline` linear state-machine driver.
- Story 2 → `baseline plan`, `baseline apply`, strict JSON schemas and exits.
- Story 3 → root-carrier inventory, backups, classification and repository rules.
- Story 4 → embedded Catalog and repository-owned CustomProfile loader.
- Story 5 → sealed ProposalAnalyzer and consolidated editable review.
- Story 6 → revision state, scoped revision proposal and plan recalculation.
- Story 7 → canonical ManagedEntry ledger and derived `fileChanges`.
- Story 8 → Repository Capability evaluator and recommendation-only verifier.
- Story 9 → thin setup skill calling the public command family.
- Story 10 → user guide, CLI reference, migration, recovery and troubleshooting.

## Integration Points

- **Git:** a narrow command adapter discovers the root, Git-private transaction
  path, object format, and root commits. It does not require branch, upstream,
  or clean-state policy.
- **Filesystem:** root-anchored repository and transaction adapters perform
  bounded reads, safe-link inspection, staging, fsync, replacement, rollback,
  and postimage verification.
- **ACPX/Codex ACP:** `internal/baselineacp` reuses adapter lineage and Exact
  Agent Selection Proof, then owns a non-Run sealed session lifecycle.
- **Embedded assets:** `go:embed` makes the CLI binary the only runtime catalog
  authority; maintainer sync validates provenance and canonical digests.
- **Skill distribution:** the existing skill sync gate distributes only
  guidance after Python removal and checks its recipes against CLI help.

## Testing Approach

Before porting behavior, the 240 existing Python tests are classified into
`exact`, `semantic`, `designed-delta`, `ancillary`, or explicitly justified
`retired` rows. Characterization fixtures capture repository preimage,
decision input, normalized output, refusal, planned bytes, ledger,
`fileChanges`, manifest, digest, post-state, and rollback. Go and Python run
side by side until every maintained row passes; the checked-in corpus remains
after Python removal so CI no longer depends on Python.

The Fluxus finding receives seven named macro assertions: safe aliases retain
source evidence before remediation; decisions are collected without repeated
full-plan noise; emitted Decision Documents parse with their required schema;
HTTP candidates include bounded route facts and source digests without policy
inference; the 256-source and 2 MiB budgets count only files that produce HTTP
candidates; PostgreSQL diagnostics separate implementation evidence from a
missing contract; file review leads with `fileChanges`; and only locally
validated commands are labeled repository-executable. This mapping is part of
the final QA matrix rather than narrative-only coverage.

Stdlib table tests cover schemas, catalog references, ordering, digests,
aliases, backups, identity, preimages, transitions, projection, and reapply.
Filesystem tests use `t.TempDir` Git repositories to inject failure around
each transaction phase, verify modes and bytes, recover crashes, and assert
rollback. Mutation tests corrupt relationships to prove validators fail for
the intended reason.

CLI tests call `RunContext(args, stdout, stderr)` and assert help, no-TTY
behavior, stdout/stderr separation, schemas, stable exits, plan/apply handoff,
profile commands, and actionable failures. Macro tests build the real binary
and execute greenfield, preservation, update, profile change, stale plan,
cross-clone apply, unsafe alias, and empty-reapply journeys.

ACP unit tests use a hand-rolled fake at the proposal seam. They assert exact
selection proof, sealed flags, byte-identical fallback input, strict output
validation, cancellation, timeout, cleanup, no tool events, zero checkout
mutation, and no Run Database or Agent log. An opt-in local ACPX integration
test proves the real adapter contract; model quality is never judged by another
LLM. Documentation contract tests parse every published command and fixture.
Final QA runs maintained formatter/Verification commands only as external
composition evidence and performs the live Fluxus journey only after separate
authorization.

## Build Order

1. Freeze the Python compatibility matrix and characterization corpus,
   including explicit designed deltas.
2. Establish `internal/baseline`, move the canonical embedded catalog, and
   implement strict schemas, profiles, serializers, and digest domains
   (depends on: 1).
3. Implement Git identity, bounded inventory, safe carriers, Source Baseline,
   Readoption, Repository Capabilities, and complete preimages
   (depends on: 2).
4. Implement Decision Plan, rendering, retention, Setup Manifest,
   ManagedEntry ledger, file projection, and portable PlanDocument
   (depends on: 2, 3).
5. Implement the recoverable transaction, immutable backups, exact apply,
   rollback, recovery, and Baseline verifier (depends on: 3, 4).
6. Implement the sealed ACP proposal adapter, preferred/fallback supervisor,
   classification review, and scoped plan-revision proposals
   (depends on: 3, 4).
7. Deliver the human state machine plus `plan`, `apply`, and `profile`
   CLI contracts with text/JSON renderers and documentation-test seams
   (depends on: 4, 5, 6).
8. Port `baseline skills restore` and maintainer `baseline assets sync` to Go
   with their parity suites (depends on: 2, 5, 7).
9. Complete the dedicated Roundfix documentation Task: user guide, CLI
   reference, automation schema, migration, recovery, troubleshooting,
   examples, and thin setup skill (depends on: 7, 8).
10. Run the full parity and composition matrix, remove the Python runtime and
    Python fallback, and tighten skill-sync checks around the Go authority
    (depends on: 1–9).
11. Execute final QA, including the separately authorized Fluxus journey, and
    collect release-ready evidence (depends on: 10).

## Risks & Considerations

- Portable plans can contain repository policy and exact generated bytes.
  Documentation must treat them as review artifacts that may be sensitive; ACP
  raw material and messages are excluded.
- Root-commit lineage distinguishes unrelated repositories but not every fork.
  Full bounded preimage validation remains mandatory and is the final
  protection against applying an unrelated plan.
- Multi-file replacement is not a single filesystem primitive. A durable
  phase journal, saved preimages, directory sync, reverse rollback, and
  mandatory recovery make interruption observable and recoverable.
- Sandbox flags can drift with ACPX. Adapter lineage/version proof and an exact
  argument contract must fail closed; unavailable analysis falls back to
  manual classification without blocking deterministic Baseline use.
- Semantic input limits can force manual review in unusually large carriers.
  Limits are explicit in results and never justify truncation or partial
  approval.
- Full legacy parity is expensive but prevents silent loss. The matrix must
  name every intentionally changed behavior before Python removal.
- Built-in-only custom profile composition is deliberately less extensible
  than arbitrary repository assets. It keeps v1 deterministic and non-
  executable; expanding that trust boundary requires a later ADR.

## Decisions

- The CLI Go engine is the sole Baseline runtime authority. See ADR-0066.
- Custom profiles live under `.roundfix/baseline/profiles` and compose only
  the embedded catalog. See ADR-0067.
- Humans use one linear workflow; automation uses portable `plan` and `apply`
  stages with separate strict schemas. See ADR-0068.
- ACP receives a sealed snapshot without checkout or tool access and remains
  proposal-only. See ADR-0069.
- Audit covers bounded carriers, but automatic preservation changes only root
  carriers. See ADR-0070.
- Plans are path-relative, cross-clone portable, and bounded-preimage guarded.
  See ADR-0071.
- Go cutover preserves every maintained Python contract and ports auxiliary
  operations under `roundfix baseline`. See ADR-0072.
- Apply uses a recoverable multi-file transaction. See ADR-0073.
- Git is required, but dirty state, detached HEAD, and missing upstream are
  allowed.
- Human interaction uses linear prompts, not Bubble Tea.
- Baseline verification never executes repository formatter or Verification
  commands.
