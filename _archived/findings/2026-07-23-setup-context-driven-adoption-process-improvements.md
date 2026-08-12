---
status: done
created_at: 2026-07-23
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-baseline-and-derived-tooling.md
---

# Context-Driven setup — live adoption exposed avoidable operator friction (2026-07-23)

A live adoption of the `standard-typescript-monorepo` profile in
`/Users/marcio/dev/fluxus` completed with the current `0.0.1` setup contract,
but the path to a confirmed plan required manual repository cleanup, repeated
large audit responses, hand-built structured evidence, and a capability
document that duplicated facts already present in the repository. This report
keeps the safety properties delivered by
[Spec 0045](../specs/0045-context-driven-baseline-0-0-1-reset/_prd.md) and
records process improvements exposed by a real repository adoption.

Session evidence:

- Initial audit blocked with `source-inventory.carrier.symlink` for root and
  backend `CLAUDE.md` carriers. At the pre-adoption commit, both paths had Git
  mode `120000` and pointed to `AGENTS.md`.
- The resolved Decision Plan contained 11 repository decisions and produced 27
  `plannedChanges`. Its confirmed digest was
  `511e4d7cfbdc155c2c187055c5baf0cc1baf5212372b3b542c5911d86498c46b`.
- Apply completed only after `decision-file.schema.invalid` and
  `capability.required.missing` had been resolved.
- The final sequence passed: `rtk bun run fmt`, `rtk make verify`, resolved
  audit, unconfirmed empty apply, and `rtk git diff --check`.

## 1. Symlink rejection can remove source evidence before Readoption inventories it

- **Symptom / evidence**: Audit rejected `CLAUDE.md -> AGENTS.md` and
  `packages/backend/CLAUDE.md -> AGENTS.md` with
  `source-inventory.carrier.symlink`. Setup could not continue until the
  carriers were removed. After the backend carrier and its 478-line target
  were removed by explicit maintainer choice, the next audit could no longer
  include those bytes in the Source Baseline inventory.
- **Root cause**:
  [`_discover_inventory_carriers`](../../.agents/skills/setup-context-driven/scripts/context_baseline.py)
  correctly refuses to follow instruction symlinks, but its diagnostic records
  only the carrier path and the message `instruction carrier must not be a
  symlink`. The preflight offers no bounded target fingerprint or disposition
  path before requiring repository remediation.
- **Action / suggestion**: Preserve fail-closed behavior for apply. For a
  symlink whose resolved target is a regular file inside the repository, make
  audit report the alias, safe target path, target digest, and carrier
  relationship without treating the alias as an independent source. Present
  explicit convert, remove, or retain-and-block choices before any carrier
  disappears. Keep external, cyclic, escaping, special-file, and unreadable
  targets opaque and blocking.

## 2. One-at-a-time decisions cause repeated plan noise and manual state handling

- **Symptom / evidence**: The Fluxus adoption required 11 decisions:
  `autonomous.enabled`, `domain.layout`, `http.contract`,
  `language.generated`, `repository.extension.enabled`, `runtime.backend`,
  `runtime.design`, `secondbrain.enabled`, `spec.scaffold`,
  `triage.external`, and `verification.gate`. Audit was rerun after each
  answer, and each response repeated the provisional `plannedChanges` payload.
  The operator also had to maintain a temporary Decision File outside the
  setup-owned manifest.
- **Root cause**: The
  [setup skill](../../.agents/skills/setup-context-driven/SKILL.md) requires
  the agent to ask one decision at a time while also presenting every
  `plannedChanges` entry from each audit. The machine contract has unresolved
  decision identifiers, but no compact `nextDecision`, resumable Decision File
  skeleton, or aggregate plan projection for this progressive phase.
- **Action / suggestion**: Replace repeated audit invocations with one stateful
  Baseline workflow. Collect each required repository decision without
  inferring an answer, then present one consolidated, editable review and final
  Change Plan. Non-interactive use must emit the complete plan in one
  machine-readable response and apply only its exact approved digest.

## 3. The skill's Decision File description omits a required schema field

- **Symptom / evidence**: A Decision File containing `version` and
  `decisions` was rejected with `decision-file.schema.invalid`. Adding
  `"schemaVersion": "setup-context-driven/decisions/0.0.1"` made the same
  input parse.
- **Root cause**: The setup skill describes the file as strict
  `setup-context-driven/decisions/0.0.1` JSON with `version`, `decisions`, and
  optional `readoption`, but does not state that `schemaVersion` is a required
  property. The
  [user guide](../user-guide/context-driven-development.md) contains the
  complete JSON example, while the executable parser requires all three
  top-level fields.
- **Action / suggestion**: Make the skill, user guide, CLI help, emitted
  Decision File skeleton, and parser describe the same exact fields. Add a
  documentation contract test that checks the skill example against the real
  parser rather than only checking that the schema identifier appears in
  prose.

## 4. HTTP contract capture requires avoidable manual evidence assembly

- **Symptom / evidence**: Resolving `http.contract` required a manual route
  inventory, an SHA-256 digest of
  `packages/backend/src/infra/controllers/http/app.ts`, and four ordered typed
  exceptions for Better Auth, health, OpenAPI, and Scalar routes. Setup
  correctly refused to infer REST versus Post-only, but it supplied no
  repository-derived candidate evidence to help construct the typed value.
- **Root cause**: The profile owns a strict policy decision, exception schema,
  and source binding, but audit exposes only the unresolved decision. Route
  discovery and source-digest preparation are left entirely to the operator or
  agent.
- **Action / suggestion**: Keep policy selection explicit. Add a bounded,
  read-only evidence projection that lists candidate route carriers, observed
  methods and scopes, and their current digests. Emit those facts as
  unclassified candidates in the Decision File template; never assign the
  contract mode, owner, or reason automatically.

## 5. PostgreSQL capability evidence is narrower than the diagnostic implies

- **Symptom / evidence**: Audit reported `capability.required.missing` and
  `PostgreSQL has no compatible local evidence` even though Fluxus already
  declared `postgres:18-alpine` in `docker-compose.yml`, the `postgres` driver
  and `drizzle-orm` in `packages/backend/package.json`, PostgreSQL-backed
  Drizzle imports, and `DATABASE_URL`. Setup became ready only after adding a
  new root `DATABASE.md`.
- **Root cause**: The
  [Standard TypeScript Monorepo profile](../../.agents/skills/setup-context-driven/assets/profiles/standard-typescript-monorepo.json)
  probes only `DATABASE.md`, `docs/architecture/database.json`, and
  `docs/architecture/postgresql.json` for the literal `PostgreSQL`. The
  resulting diagnostic does not distinguish “no implementation evidence” from
  “no accepted repository contract document.”
- **Action / suggestion**: Decide which claim the capability must prove. If it
  proves stack presence, accept ranked local evidence from the driver,
  adapter, compose image, or typed database configuration. If it requires an
  explicit repository contract, rename the diagnostic and next action to say
  that the contract document is missing, list every accepted path, and report
  the implementation evidence that was found but was insufficient.

## 6. Change Plan detail is exact but difficult to review at file level

- **Symptom / evidence**: The resolved Fluxus plan contained 27 entries.
  Multiple root `AGENTS.md` managed-block entries repeated the same final-file
  digest, while composite guide changes were split by managed identity. This is
  useful for ownership accounting but makes the maintainer reconstruct the
  actual file-level patch before confirming one digest.
- **Root cause**: `plannedChanges` is both the canonical managed-entry ledger
  and the reader-facing plan. It has no separate aggregation by path.
- **Action / suggestion**: Preserve the current complete entry ledger and its
  contribution to `planDigest`. Add a derived `fileChanges` projection with
  one row per path, file action, before/after digest, and ordered managed IDs.
  The skill can present this file view first and retain the full ledger for
  retention and machine review.

## 7. Persisted verification commands can disagree with the adopted repository

- **Symptom / evidence**: The generated Fluxus Setup Manifest persisted
  `verification.format.command: "bun run format"` and
  `verification.workspace.command: "bun run verify"`. The root
  `package.json` defines `fmt` and does not define `format` or `verify`.
  Separately, the selected repository Verification was correctly persisted as
  `rtk make verify`, and the successful formatter command was
  `rtk bun run fmt`.
- **Root cause**: The profile snapshot copies canonical verification command
  names without validating them against the selected repository scripts. The
  repository-owned `verification.gate` decision does not reconcile or label
  those profile-level commands.
- **Action / suggestion**: Distinguish portable verification roles from
  repository-executable commands in the schema. Either resolve each role to a
  validated local script during decisions or label unresolved canonical
  commands as profile expectations rather than runnable repository commands.
  Audit must report a mismatch before apply. Completion guidance can emit the
  exact resolved formatter and Verification commands, while audit and apply
  remain network-free and continue to execute no repository scripts.

## What worked — keep

- Audit and apply failed closed and wrote nothing for symlink, invalid
  Decision File, missing capability, incomplete decision, and unconfirmed plan
  states.
- Exact `planDigest` confirmation bound the write to the reviewed repository
  preimage and normalized decisions.
- The Repository-Owned Extension remained outside setup-managed bytes.
- Formatter, Verification, post-apply audit, and empty reapply composed without
  a generated delta.
- The final generated root remained compact while the detailed contracts lived
  under `docs/agents/`.

## Routing

[Spec 0046](../specs/0046-public-context-driven-baseline-command/_prd.md)
owns the inventory-remediation, decision UX, capability-evidence,
plan-presentation, verification-command, public CLI authority, custom-profile,
and Python-removal contracts identified here. Its approved product direction
uses one `roundfix baseline` workflow, immutable root-instruction backups,
supervised read-only ACP proposals, one selected Baseline Profile, consolidated
plan review, and digest-bound non-interactive apply. Audit remains
repository-wide even though automatic preservation changes only root
instruction carriers.

The Spec's future Task Graph must contain a dedicated Roundfix documentation
Task covering the user guide, CLI reference and examples, automation,
migration from the Python-backed skill, recovery, troubleshooting, and the thin
setup skill. This finding is triaged and closed; Spec implementation and QA
remain tracked by that Spec.
