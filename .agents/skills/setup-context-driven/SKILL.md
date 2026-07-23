---
name: setup-context-driven
description: Configure a repo for CONTEXT-driven development (the method explained in docs/user-guide/context-driven-development.md) — scaffold the full docs/ layout (_inbox, adr, agents, design, findings, handoffs, references, specs, user-guide) and the CONTEXT.md glossary, and seed the docs/agents/ usage guides (docs layout with the findings template, issue tracker, spec routing, domain docs, triage labels, autonomous work model, and optional Secondbrain guidance). Run when preparing a repo for the write-prd/write-tasks/implement pipeline or for Roundfix-driven autonomous work; re-run to audit and refresh only setup-owned managed content.
disable-model-invocation: true
metadata:
  category: setup
  tags: [workflow, prd, issues, planning, triage, repository-context, agents]
  version: 0.9.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Setup Context-Driven

Configure a repository for CONTEXT-driven development through the portable asset catalog and `scripts/context_setup.py`. The script is the source of truth for audit, apply, setup snapshots, managed markers, decisions, and finding codes.

## Asset map

- `assets/profiles/` selects a supported profile and canonical skill setup snapshot: `typescript-bun-monorepo`, `go-cli-tui`, or `rust-cli`.
- `assets/source-baselines/` publishes strict `0.0.1` project-agnostic governed corpora. Each independent index pins every Normative Clause, recommendation, and Operational Contract; `accounting.json` records individual retained or rejected prior clauses.
- `assets/coverage.json` defines stable semantic coverage categories. A profile's `requiredRules` must prove every universal and applicable category; coverage is not a line-count target.
- `assets/modules/` owns compact root pointers, supporting guides, portable rule guidance, required decisions, exact `requiredSkills`/`skillDispatch` mappings, and typed references.
- `assets/templates/` stores generated repository content. Root blocks must stay short and point to `docs/agents/` guides.
- `assets/setups/` stores `setup-context-driven/setup-snapshot-v2` records. External skills carry a GitHub repository, immutable commit, safe source path, and complete-tree digest; Roundfix-owned skills remain separate.
- `references/` is workflow guidance for agents, not generated output. Maintainers must read [`references/asset-maintenance.md`](references/asset-maintenance.md) before changing catalog or snapshot data.

The baseline owns only declared managed blocks, setup-owned guides, the Setup Manifest, and portable workflow rules. Repository-authored architecture and policy remain outside setup ownership. In particular, frontend guidance can require a repository-owned `DESIGN.md`, but setup never generates project-specific architecture.

Source Baseline carriers keep root instructions compact and render full guide fragments. Corpus validation rejects denied project tokens, machine-specific paths, and copied generated managed markers before any carrier can be planned.

## Audit-first workflow

1. Inspect the repository just enough to pick the likely profile. Use local files only: root instructions, `CONTEXT.md`, `docs/`, package files, and language/runtime markers.
2. Run audit before asking setup questions:

   ```bash
   rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --format json
   ```

   If a profile is known but the manifest is missing or stale, rerun with `--profile <profile-id>` to preview the selected composition.
3. Require JSON schema `setup-context-driven/audit-v1`. Audit is local, read-only, and network-free. Read the result and present only:
   - selected profile and ordered modules;
   - selected canonical skill setup name;
   - active or conditional modules and their trigger decisions;
   - blocking findings by `code`, `managedId`, `path`, `message`, `action`, and structured `remediation` when present;
   - every `plannedChanges` entry, including `action`, `path`, `managedId`, `state`, `reason`, `beforeDigest`, `afterDigest`, and optional `condition`, `fromPath`, and `referenceEdits`;
   - `retentionAccounting`, when the repository has a recognized source-baseline transition, including every entry's `fromClause`, `enforcement`, `disposition`, `targets`, and `reason`;
   - optional cleanup information only when audit was run with `--show-extra-skills`.
4. Do not dump generated Markdown by default. Mention that full templates live under `assets/templates/` if the user asks to inspect them.

## Decisions

Use stored compatible decisions from `docs/agents/setup-context.json` first. Ask only for `decision.required` findings, one decision code at a time, in the order the CLI reports them. `verification.gate` is an entry decision for every profile. `runtime.backend` and `runtime.design` become required after `autonomous.enabled=true`. Do not ask again for a stored compatible value.

Question routing:

- `spec.scaffold` — confirm local `docs/specs/<feature-slug>/` is the planning source.
- `domain.layout` — ask whether the repository is `single-context` or `multi-context`.
- `triage.external` — ask only when external forge issues are relevant.
- `autonomous.enabled` — ask whether Supervisor-to-ACP Runtime delegation applies.
- `runtime.backend` — ask for the backend/default implementation runtime and model.
- `runtime.design` — ask for the design, UI, UX, or frontend runtime and model.
- `verification.gate` — ask for the command agents must run before completion claims.
- `language.generated` — generated repository content must be `English`.
- `secondbrain.enabled` — ask whether read-only local Secondbrain guidance must be generated.
- `adoption.*` — ask only after showing the existing unmarked file that would become setup-owned.

Record each answer as a repeated `--decision ID=VALUE` argument. For newly introduced decisions, ask the new code once, then let the manifest carry it on later runs.

After every required answer is known, rerun audit with the selected profile and the complete decision set:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --format json --profile <profile-id> --decision <id=value> ...
```

Only a fully resolved Decision Plan has an authorizable `planDigest`. If audit still returns `decision.required`, keep the repository unchanged and ask only the next reported decision.

## Upgrade Retention Contract

Before a Context-Driven Baseline version transition can mutate the repository, the Upgrade Retention Contract accounts for every previously managed mandatory clause. Each clause must map to a reachable current clause with the same enforcement strength, a Repository-Owned Extension, or an explicit rejection with a recorded reason. The ordered `retentionAccounting` output is the normative part of the Change Plan: it reports each `fromClause`, its `enforcement`, the `retained`, `moved`, `replaced`, or `rejected` disposition, any `targets`, and the `reason`. The same accounting appears in text output and contributes to `planDigest`, so a normative mapping change makes an earlier confirmation stale even when file bytes are unchanged.

Audit and apply fail closed with exit `1` and no writes when the source baseline is unknown, a prior clause is unaccounted, a target is absent from the selected future artifact graph, or a target weakens enforcement. Report the finding code and affected clause or target. Do not infer a baseline, weaken a clause, edit generated output, or bypass the transition ledger; correct the source-baseline identity or the canonical transition asset, rerun audit, review the new complete Change Plan, and obtain confirmation for its new `planDigest`.

## Preview and apply

Apply is explicit, local, and network-free. It never runs before the user sees the complete managed change summary from the resolved audit.

Before applying, tell the user:

- the profile, ordered modules, and canonical skill setup;
- every managed file or block that will be created, refreshed, or removed;
- every Upgrade Retention Contract entry, when `retentionAccounting` is present;
- that only `docs/agents/setup-context.json` and declared setup-owned Markdown boundaries can change;
- that repository-authored bytes outside managed markers remain untouched;
- that extra installed skills are informational only and this workflow never removes skills.

Ask for confirmation of the exact `planDigest`. After confirmation, rerun with the same profile and decisions and bind the write to that digest:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py apply --repo <repo> --format json --profile <profile-id> --decision <id=value> ... --confirm-plan <planDigest>
```

If apply returns `plan.confirmation.stale`, the preimage or inputs changed: no write occurred. Present the recomputed plan and ask for confirmation of its new digest. If apply returns `decision.required`, present the unchanged selection and preview and ask only the next unresolved decision. For any other blocking finding, report its code, managed identity, path, reason, and action; do not patch around it manually.

After apply, rerun the same resolved audit in JSON mode. A documentation-clean result exits `0`. Do not treat required-skill drift as permission to fetch or replace anything; restoration is the separate, explicit workflow below.

## Repository ownership and delegation floor

`repository.extension.enabled=true` authorizes setup to create the declared Repository-Owned Extension once when it is absent and include that creation in the confirmed Change Plan. The extension file is unmarked, never enters `managedArtifacts`, and remains outside setup management: audit, reapply, and profile transitions do not compare, rewrite, format, or remove its bytes. Setup owns only the generated typed reference to it. If a previously selected extension is later missing, report the blocking `reference.repository.missing` finding; do not recreate repository content automatically.

Audit also scans bounded root and nested agent-instruction documents for delegation to setup-managed categories. `delegation.baseline-floor` is an informational baseline-floor finding: it names repository-authored documents that delegate to a category absent from the active catalog, but it does not enter the Change Plan, block apply, or grant setup authority over those documents. Explain that the generated baseline is a floor, not a replacement for project-specific policy.

## Skill dispatch and Formatter-Stable Output

The selected catalog is the single source for both the Repository Skill Set and generated dispatch. Each installed skill renders exactly once in `docs/agents/skill-dispatch.md`; that one entry can list multiple distinct triggers. Do not add duplicate dispatch prose to another module or generated guide.

Generated managed Markdown is Formatter-Stable Output for its selected profile. Formatter proof is pinned and profile-specific: the TypeScript/Bun profile binds the declared Oxfmt version, checked-in golden corpus, and provenance digest, while profiles with no selected Markdown formatter declare `none`. Ordinary repository Verification is hermetic and compares generated bytes with that pinned corpus; it downloads nothing and executes no formatter. Final QA separately runs the real pinned formatter probe in a disposable profile fixture, then its selected Verification, a fresh audit, and a second apply. Setup audit and apply never execute a formatter, the repository's Verification, or downloaded content.

## Explicit Repository Skill Set restoration

Audit computes installed complete-tree digests from local files. Missing or drifted external required skills include structured `remediation` with the exact provider, skill, repository, immutable ref, source path, expected tree digest, and preview argv. Audit and documentation apply never execute that argv automatically.

Only after the maintainer explicitly requests restoration, preview the affected external skills. The preview owns Git acquisition and can use the network; `--source-dir` instead points to an offline Git checkout or bare object store that contains the declared exact commit.

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py restore-skills --repo <repo> --profile <profile-id> --skill <name> --format json [--source-dir <git-source>]
```

Require schema `setup-context-driven/restore-v1`. Present its `skills`, `acquisitions`, every created/refreshed/removed file and lock edit in `plannedChanges`, and `planDigest`. The preview exits `3` with `plan.confirmation.required` for a non-empty plan and does not mutate the repository.

Ask for confirmation of that exact digest, then rerun the same preview arguments with the digest:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py restore-skills --repo <repo> --profile <profile-id> --skill <name> --format json [--source-dir <git-source>] --confirm-plan <planDigest>
```

The command restores the complete declared directory, including removal of files absent from the immutable source, updates only the selected portable lock entries, and verifies the final tree digest. A stale plan exits `3` without mutation and requires a new preview and confirmation. Source, proof, safety, lock-adapter, or write failures exit `1` before mutation or roll back all targets. There is no branch, default-revision, or generic skill-refresh fallback.

After restoration, rerun the same resolved audit. A clean result proves the Context-Driven Baseline and selected Repository Skill Set; it does not prove project-specific architecture or policy completeness.

## Exit categories

- exit `0`: audit is clean, apply completed or was already current, or restoration completed or was already current;
- exit `1`: blocking audit or retention findings, or a source, proof, safety, lock-adapter, or write failure;
- exit `2`: invalid arguments, malformed decisions or confirmation digest, invalid skill selection, or malformed lock input;
- exit `3`: decisions are required, or a non-empty plan needs confirmation or became stale.

Text and `setup-context-driven/audit-v1` JSON results go to stdout. Diagnostics go to stderr. Audit and documentation apply are local, use only bundled assets, and remain network-free; only an explicitly requested `restore-skills` operation owns Git acquisition.

## Optional reports

Use `--show-extra-skills` only when the user asks to review installed skills outside the selected setup:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --format json --show-extra-skills
```

Report `skills.extra.installed` and `skills.local.untracked` as review information. Never suggest a removal command.

Use `sync-setups` only as a maintainer operation with an explicit canonical setups directory:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py sync-setups --source-dir <canonical-setups> --check --format json
```

Normal repository setup does not require this checkout.

## Secondbrain

Secondbrain is opt-in through `secondbrain.enabled=true`. When enabled, apply creates one compact root pointer and `docs/agents/secondbrain.md`. When disabled, apply creates neither; if a previous setup owned those artifacts, apply removes only the marked managed Secondbrain block or guide content.

Generated Secondbrain guidance must remain read-only. It must require index-first lookup, `qmd query`, project-mirror caution, file citations, Hermes escalation for durable updates, and secret safety. It must forbid writes to the Secondbrain, `raw/`, and `projects/*/mirror/`.

## Completion

Report every remaining finding by code, managed identity, and path. Do not claim setup is complete without a fresh resolved audit at exit `0`. The full sequence is audit, decision resolution, complete plan review, exact digest confirmation, apply, final audit, and—only when explicitly requested—previewed and digest-confirmed drift restoration followed by another audit.
