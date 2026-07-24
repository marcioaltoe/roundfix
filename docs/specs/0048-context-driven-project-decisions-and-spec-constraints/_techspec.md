---
spec: 0048-context-driven-project-decisions-and-spec-constraints
prd: _prd.md
created: 2026-07-24
---

# Context-Driven project decisions and Spec constraints — Technical Spec

## Executive Summary

This feature adds two typed catalog decisions and renders them through the
existing deterministic Baseline Plan: `identifier.strategy` owns new
project-controlled identifiers, while `auth.provider` owns the selected Better
Auth contract and derives its HTTP exception. The Setup Manifest remains the
machine authority, and generated domain and backend guides carry complete
operative rules without requiring an ADR. Repository-owned authorial skills
also gain a mandatory human-readable `Project Constraints` snapshot; they stop
before finalizing a new artifact when that section is missing or incomplete.
The design accepts duplicated authentication projection across two decisions
and resolves the risk with deterministic conflict validation.

## Project Constraints

- Identifier strategy: applicable. The Baseline decision suggests
  `{"kind":"uuid-v7"}` but this implementation creates no application
  identifier and migrates no existing value.
- Authentication and HTTP: applicable. The Baseline models Better Auth and a
  derived route exception but does not install the provider or mutate
  application routes.
- Active ADR obligations: ADRs 0063, 0066, 0068, 0071, 0073, 0076, and 0077
  govern HTTP ownership, CLI authority, explicit confirmation, portable
  planning, apply, decision derivation, and Spec snapshots.
- Tooling authority: no tooling configuration, script, ignore file, plugin
  declaration, or version pin is authorized for modification by this Spec.
- Sources: `docs/agents/agent-instructions.md`,
  `docs/agents/domain.md`, `docs/agents/backend.md`, and
  `docs/agents/spec-routing.md`.

## System Architecture

Spec 0048 starts only after Spec 0047 has a passing QA verdict because its
render targets and semantic ownership are prerequisites. It extends four
existing components:

1. `internal/baseline/assets/decisions.json` declares the identifier and
   authentication decision types, suggestions, conditional requirements, and
   render bindings.
2. `internal/baseline` validates the discriminated objects, derives one
   Better Auth HTTP exception, resolves conflicts, persists normalized values,
   and renders self-contained domain and backend rules.
3. The core and Spec workflow modules render the universal tooling-authority
   clause and the Project Constraint authoring and execution contract.
4. The repo-owned `write-prd`, `write-techspec`, `write-tasks`,
   `implement-task`, and `qa-gate` skills enforce the new artifact contract.
   Upstream-managed skills, including `domain-modeling`, remain unchanged.

The root human command and `baseline plan` continue to call the same
`ResolveDecisionInput` and `BuildPlan` path. Suggested values exist in catalog
metadata for interactive presentation only; normalization never treats a
default as an automation answer.

## Implementation Design

### Interfaces

The two decision values use explicit domain types at validation and rendering
boundaries:

```go
type IdentifierStrategy struct {
	Kind     string `json:"kind"`
	Guidance string `json:"guidance,omitempty"`
}

type AuthProviderDecision struct {
	Kind           string        `json:"kind"`
	RouteException HTTPException `json:"routeException"`
}
```

Derivation is pure and runs before Plan digest calculation:

```go
func NormalizeProjectDecisions(
	profile ResolvedProfile,
	input []DecisionValue,
	catalog *Catalog,
) ([]DecisionValue, error)

func DeriveHTTPContract(
	contract HTTPContract,
	auth *AuthProviderDecision,
) (HTTPContract, error)
```

The implementation can keep these functions package-private if no second
caller needs them. Their observable contract is the normalized Decision Plan,
Setup Manifest, rendered guides, and portable Plan digest.

### Data Models

`identifier.strategy` has decision type `identifier-strategy` and accepts
exactly:

```json
{"kind":"uuid-v7"}
```

or:

```json
{"kind":"repository-defined","guidance":"<non-empty operative rule>"}
```

Unknown fields, an empty `guidance`, or `guidance` on `uuid-v7` fail strict
validation. The rule applies only to new project-owned Internal Identifiers.
The generated domain guide states that provider identifiers, protocol
identifiers, natural keys, and business codes preserve their source contracts.

`auth.provider` has decision type `auth-provider`. It is selected only when the
resolved Profile retains `capability.stack.better-auth` and accepts
`kind == "better-auth"` plus one strict `routeException`:

```json
{
  "kind": "better-auth",
  "routeException": {
    "scope": "/api/auth/*",
    "methods": ["GET", "POST"],
    "owner": "Better Auth",
    "reason": "<non-empty provider-protocol reason>"
  }
}
```

The catalog suggestion contains the confirmed scope, methods, owner, and the
session/OAuth/callback/provider-protocol reason. A human can change the typed
scope, methods, or reason during the explicit decision step, but the Better
Auth kind requires owner `Better Auth`. Changing provider requires a Profile
adaptation that removes the Better Auth capability rather than an internally
contradictory value.

The derivation merges the route exception into `http.contract` by normalized
`(owner, scope)`. An identical exception is idempotent. A different exception
with the same key, duplicate methods, an unsupported method, or a missing
reason stops planning and identifies both decision IDs. Final exceptions sort
by scope, owner, then methods before serialization and rendering.

`SetupManifest.Decisions` already stores arbitrary strict JSON values, so both
objects fit the existing manifest and Plan schemas. The derived
`http.contract` is persisted beside `auth.provider`; every audit re-derives and
compares them before reuse.

The tooling rule is a catalog `Normative Clause`, not a Decision Value. It
covers creating, editing, renaming, moving, or deleting:

- linter, formatter, typechecker, test-runner, architecture-checker, build,
  package-manager, or code-generator configuration;
- tooling scripts and ignore files;
- plugin declarations and version pins; and
- other repository-tooling policy files named by the repository.

No Profile effect can exclude this core clause.

### API Contracts

For a Profile with unresolved values, the human flow presents:

1. `identifier.strategy` with UUID version 7 as the visible keep-or-change
   suggestion;
2. `auth.provider` only when Better Auth is selected, showing the complete
   route exception rather than a provider name alone; and
3. `http.contract`, with the derived exception visible in its complete value.

Enter on a marked choice counts as an explicit human answer. Existing
compatible Manifest values appear as keep-or-change choices and outrank catalog
suggestions. An incomplete or conflicting stored value cannot be reused.

Automation supplies the same objects through a strict Decision Document:

```json
{
  "schemaVersion": "setup-context-driven/decisions/0.0.1",
  "version": "0.0.1",
  "decisions": [
    {"id": "identifier.strategy", "value": {"kind": "uuid-v7"}},
    {"id": "http.contract", "value": {"mode": "REST"}},
    {
      "id": "auth.provider",
      "value": {
        "kind": "better-auth",
        "routeException": {
          "scope": "/api/auth/*",
          "methods": ["GET", "POST"],
          "owner": "Better Auth",
          "reason": "Preserve provider-owned session, OAuth, callback, and protocol semantics."
        }
      }
    }
  ]
}
```

If either selected decision is unresolved, `baseline plan` returns the existing
`roundfix/baseline-result/v1` action-required document, exit `3`, all missing
IDs, and no partial Plan. Inline `--decision` remains available for JSON
objects because its current parser already admits object values.

The domain render binding converts `identifier.strategy` into an operative
Markdown clause. The backend binding renders HTTP mode, every ordered
exception, and the Better Auth rule. Structured renderers escape Markdown
content and reject marker-shaped or non-canonical text before postimage
assembly.

Project Constraint snapshots are Markdown sections, not frontmatter. Every new
PRD and TechSpec records:

- identifier strategy as applicable or not applicable with a reason;
- authentication and HTTP policy as applicable or not applicable with a
  reason;
- active ADR obligations relevant to the change;
- tooling authority, including express authorization and exact bounded files
  only when granted; and
- the source paths under `docs/agents/`.

`write-prd` and `write-techspec` must not report completion until the section is
present and complete. `write-tasks` refuses decomposition when a non-archived
Spec lacks the section or proposes tooling work without recorded authority.
`implement-task` refuses tooling mutation outside the listed files.
`qa-gate` verifies the section and authorization evidence. Existing archived
or otherwise completed Specs are not rewritten.

## Coverage Map

- Goal 1 → typed decision declarations, human prompts, strict automation input.
- Goal 2 → Setup Manifest persistence and domain/backend structured renderers.
- Goal 3 → Project Constraint templates and authorial skill completion gates.
- Goal 4 → universal core tooling clause and Task execution refusal.
- Goal 5 → shared normalization, derivation, and Plan digest path.
- Story 1 → `identifier.strategy` validation and UUID version 7 suggestion.
- Story 2 → `auth.provider` and deterministic `http.contract` derivation.
- Story 3 → missing-decision result and strict Decision Document validation.
- Story 4 → PRD/TechSpec templates, `write-tasks`, and QA enforcement.
- Story 5 → universal tooling clause and bounded authorization checks.
- Story 6 → self-contained rendered clauses and concise Spec snapshots.

## Integration Points

The embedded catalog and setup snapshots remain the only Baseline content
source. Catalog validation must prove that each selected Profile includes
`identifier.strategy`, that Better Auth capabilities require `auth.provider`,
and that each structured decision has exactly one supported renderer.

Repository-owned authorial skills are updated locally and in the canonical
Repository Skill Set snapshots governed by
`docs/agents/skill-governance.md`. The implementation must compare
upstream-managed skill digests before and after and fail if
`domain-modeling/ADR-FORMAT.md` or another external skill changes.

Public documentation and the thin `setup-context-driven` skill show the exact
human suggestions, Decision Document shapes, derived exception behavior,
Project Constraint section, and refusal cases. They call only the public
Baseline commands.

## Testing Approach

Unit tests table-drive every accepted and rejected discriminator shape,
unknown field, empty guidance, Better Auth owner, method normalization,
duplicate exception, conflict, stable sort, and idempotent derivation. Tests
assert external identifiers remain outside the rendered UUID rule and that the
tooling clause cannot be disabled by any Profile decision.

CLI tests capture stdin, stdout, and stderr for visible suggestions, explicit
Enter confirmation, change flows, compatible-value reuse, missing automation
inputs, and conflict diagnostics. Equivalent human and Decision Document
answers must produce byte-identical normalized decisions, postimages, and Plan
Digests.

Catalog and golden tests cover domain, backend, agent-instructions,
spec-routing, source corpora, formatter fixtures, and canonical setup
snapshots. Skill contract tests remove one Project Constraint row, tooling
authorization condition, source path, or hard-stop instruction and must fail.
An external-skill mutation test proves upstream bytes are unchanged.

Macro tests run greenfield and update for every affected Profile, apply the
plan, run formatter and repository Verification recommendations outside
Baseline, audit, and require an empty reapply. Spec QA repeats separately
authorized Fluxus greenfield and update journeys and authors one new fixture
PRD and TechSpec through the updated skills. The full repository gate is
`rtk make verify`.

## Build Order

1. Add the strict identifier and authentication decision types, validators,
   catalog declarations, Profile requirements, and conflict tests after Spec
   0047 reaches a passing QA verdict.
2. Add deterministic Better Auth-to-HTTP derivation, Setup Manifest reuse, and
   human/automation decision parity (depends on: 1).
3. Add structured domain and backend renderers plus complete self-contained
   guide clauses and formatter fixtures (depends on: 1, 2).
4. Add the universal tooling-authority clause to the core module, templates,
   source accounting, and retention contracts (depends on: 1).
5. Update PRD and TechSpec templates and the repo-owned `write-prd`,
   `write-techspec`, `write-tasks`, `implement-task`, and `qa-gate` contracts
   for mandatory readable Project Constraint snapshots (depends on: 3, 4).
6. Update affected Profiles, canonical setup snapshots, user documentation,
   and the thin setup skill; add skill ownership and documentation guards
   (depends on: 2, 3, 4, 5).
7. Run all-profile macro coverage and separately authorized Fluxus greenfield,
   update, Spec-authoring, tooling-refusal, Verification, audit, and empty
   reapply journeys (depends on: 6).

## Risks & Considerations

Persisting the derived Better Auth exception in `http.contract` duplicates
information from `auth.provider`. The strict derivation and audit comparison
turn divergence into an action-required result rather than choosing one value
silently.

Repository-defined identifier guidance is intentionally less structured than
UUID version 7. It preserves maintainer choice but can be harder to validate;
the renderer treats it as opaque operative text, rejects unsafe marker content,
and requires a non-empty explicit human or automation answer.

Readable Project Constraint sections cannot provide a compile-time guarantee.
The authorial, task-decomposition, execution, and QA skills therefore enforce
the same rows at successive boundaries. The trade-off is repeated validation,
but no additional runtime or frontmatter schema is introduced.

Rollout changes the catalog digest, so existing repositories enter the normal
adoption/update path and review their new decisions and generated clauses.
Declining the Plan writes nothing. Apply rollback uses the existing recoverable
transaction; reverting an accepted project decision later requires a fresh
Baseline Plan and updated active Specs.

## Decisions

- Model identifier strategy as a discriminated object with UUID version 7 as a
  visible suggestion, not an inferred answer. See
  [ADR-0076](../../adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md).
- Model Better Auth separately and derive its typed route exception into the
  repository-owned HTTP Contract Decision. See ADR-0076.
- Store tooling authorization only in the readable `Project Constraints`
  section; do not add authorization frontmatter.
- Require a concise snapshot with values and `docs/agents/` sources, and block
  authorial completion when it is absent. See
  [ADR-0077](../../adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md).
- Update only repo-owned authorial skills; preserve all upstream-managed skill
  bytes.
- Start implementation only after Spec 0047 has a passing QA verdict.
