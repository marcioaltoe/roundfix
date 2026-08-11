---
spec: 0040-mandatory-context-driven-adrs
prd: _prd.md
created: 2026-07-17
---

# Mandatory Context-Driven ADRs — Technical Spec

## Executive Summary

Extend `setup-context-driven` with a first-class `decisionRecords` artifact
type and two unconditional templates owned by the existing `context-workflow`
module. Audit treats those Baseline ADRs like other managed artifacts. Apply
adds an ADR migration planner before normal rendering: it inventories the root
ADR namespace, classifies exact or managed baseline files, builds a `+2`
historical rename map, rewrites proven repository-owned references, and merges
those operations into the existing atomic `ChangePlan`.

The main trade-off is migration breadth. Reserving identities `0001` and `0002`
in mature repositories necessarily changes paths and references beyond the two
new files. The implementation therefore favors proof over convenience: Git
identifies the versioned files eligible for reference rewriting, protected or
ambiguous references block the plan, the full future tree is validated before
writes, and rollback preserves original bytes and file modes.

## System Architecture

- `.agents/skills/setup-context-driven/assets/modules/context-workflow.json`
  owns `decisionRecords` entries for the two Baseline ADRs and rule identifiers
  for colocated ubiquitous language, English ADRs, and UUIDv7 Internal
  Identifiers.
- `.agents/skills/setup-context-driven/assets/templates/adrs/` stores the
  project-agnostic English template bodies supplied in this Spec's
  `references/` directory. The generated target blocks receive the normal
  setup ownership markers.
- `context_assets.py` validates ADR artifact IDs, paths, versions, template
  kinds, rule references, and uniqueness across modules.
- `context_setup.py` gains root ADR inventory, baseline classification,
  historical renumbering, reference discovery/rewrite, migration findings,
  preview operations, and mode-preserving transactional application.
- Root and domain-guide templates point agents to the two Baseline ADRs without
  duplicating their full content in `AGENTS.md`.
- `docs/agents/setup-context.json` continues to inventory managed artifact IDs,
  paths, templates, versions, and digests. The existing managed artifacts are
  sufficient to prove later migration completion; no new setup question is
  stored.
- `make skills-sync` copies the canonical skill into
  `skills/setup-context-driven/`; the embedded copy is not edited directly.
- Spec 0036 remains responsible only for Doctor's Repository Skill Set check and
  owned-skill synchronization.

## Implementation Design

### First-class ADR assets

Advance the module asset contract so every module declares
`decisionRecords`, using an empty list when it owns none. An ADR declaration has
the same stable ownership fields as a supporting guide:

```json
{
  "id": "adr.ubiquitous-language",
  "version": 1,
  "path": "docs/adr/0001-ubiquitous-language-in-colocated-context-md.md",
  "template": "template.adr.ubiquitous-language",
  "rules": [
    "rule.context.domain-docs",
    "rule.context.adr-language"
  ]
}
```

The second artifact is `adr.internal-identifiers` at
`docs/adr/0002-uuidv7-required-for-internal-identifiers.md`, backed by
`template.adr.internal-identifiers` and
`rule.context.uuidv7-internal-identifiers`. Template catalog entries use kind
`adr`. Asset validation enforces:

- safe root paths exactly under `docs/adr/`;
- four-digit canonical prefixes;
- English canonical slugs supplied by the catalog;
- unique ADR artifact IDs and target paths;
- template kind `adr` for every decision record;
- rule ownership and valid template versions.

`ExpectedArtifact.kind` accepts `adr`. Rendering uses shared-file semantics:
the generated managed block is appended to an empty file, replaces its matching
block on refresh, and preserves content outside the markers. Preview renders
`create decision record` or `refresh managed decision` rather than calling an
ADR a guide.

Because every supported profile includes `context-workflow`, no new module,
profile variant, entry decision, or prompt is needed.

### Baseline classification

Before planning ordinary managed content, classify each canonical path as one
of:

- `managed-current` — expected marker, version, and digest are current;
- `managed-stale` — expected marker exists but needs normal refresh;
- `exact-unmarked` — the complete unmarked file equals the canonical template
  body after the repository's normal final-newline normalization;
- `missing` — no file occupies the path;
- `historical` — a regular Markdown file occupies the path but is neither a
  recognized managed Baseline ADR nor an exact template;
- `invalid` — unreadable, non-UTF-8, symlinked, malformed, duplicate, or
  otherwise unsafe state.

`exact-unmarked` is automatically wrapped in the managed marker and requires no
`adoption.*` decision. `managed-current` and `managed-stale` prove that the
baseline migration has already occurred, so ordinary refresh applies and the
remaining ADR sequence is never shifted again.

When neither canonical artifact is managed, exact unmarked artifacts are
reserved in place and all other root ADRs are historical migration inputs. A
divergent occupant of a reserved path is historical and moves with the `+2`
map. Mixed manual or partially managed state that cannot prove one coherent
prior migration yields `adr.baseline.state-ambiguous`; apply performs no
changes.

### Root ADR inventory and rename map

Inventory only regular files matching `docs/adr/NNNN-slug.md` directly below
the root directory. Do not recurse. The parser returns the four-digit integer,
slug, relative path, and baseline classification. Before constructing a map,
reject:

- two files with the same numeric identity;
- files in root `docs/adr/` that look like numbered ADRs but do not match the
  canonical filename grammar;
- historical numbers above `9997` because adding two would exceed four digits;
- a destination occupied by a non-migrating path;
- symlinks or non-regular paths;
- case-folded path collisions on case-insensitive filesystems.

For each historical root ADR, the destination number is `old + 2`, formatted as
four digits, with the slug and extension unchanged. This preserves relative
order and intentional gaps. Context-local paths such as
`src/<context>/docs/adr/` never enter the inventory.

Use temporary logical names while validating the future tree so cycles and
overlapping source/destination paths cannot cause an intermediate collision.
The final `ChangePlan` represents each rename as creation of the new path plus
deletion of the old path; no source is deleted before all destination content
has been staged.

### Reference discovery and rewriting

The migration planner runs `git -C <repo> ls-files -z` to obtain the versioned
repository paths eligible for rewriting. Git failure is a blocking input error
because the script can no longer prove which files are repository-owned. Skip
`.git`, installed or embedded Agent Skill trees, dependency/vendor trees,
build outputs, and Roundfix runtime state. If a protected tracked file contains
a reference that would become stale, emit `adr.reference.protected` and abort
instead of editing upstream-managed content.

For ordinary regular UTF-8 text files, rewrite only proven forms:

1. exact old ADR filenames and repository-relative paths using the complete
   old-to-new path map;
2. explicit labels `ADR-NNNN` and `ADR NNNN` when the old number identifies
   exactly one root ADR;
3. slash-separated explicit ADR lists such as `ADRs 0029/0030/0031`, resolving
   every member through the same unique number map.

Do not replace unqualified four-digit numbers, substrings, Spec numbers, schema
versions, or context-local ADR references. A syntactically explicit reference
that cannot be mapped uniquely yields `adr.reference.ambiguous`. After
rewriting the in-memory future tree, rescan for every old exact path and
explicit old ADR label. Any remaining repository-owned match yields
`adr.reference.stale` and blocks apply.

Reference edits preserve original encoding, newline style, executable bit, and
other permission bits. Binary files are not rewritten. A binary or symlink
that contains or represents an old ADR reference cannot be proven safe and
blocks the migration with its path.

### Change planning, preview, and rollback

Extend the current `ChangePlan` so each file mutation records original bytes,
target bytes, original mode, and operation type. The combined plan contains:

1. historical ADR target creations;
2. historical ADR source deletions;
3. versioned reference edits;
4. Baseline ADR creation, exact adoption, or managed refresh;
5. root/domain instruction refreshes;
6. Setup Manifest refresh.

Preview exposes stable ordered actions: `rename decision record`,
`update decision reference`, `create decision record`, and
`refresh managed decision`. It also emits one informational
`adr.language.review-required` entry containing the migrated historical ADR
paths. This is a manual review inventory, not language detection, and never
blocks.

Validate the complete plan before disk mutation: safe paths, no duplicate
targets, expected managed blocks and digests, no stale references, preserved
unmanaged bytes around each block, and a valid future Setup Manifest. Stage all
target contents as sibling temporary files with their intended modes. Replace
targets and remove sources only after every stage succeeds. On any `OSError`,
restore every original path, byte sequence, and mode and remove every newly
created path. An injected failure at each write/delete position must prove
rollback, including when old and new ADR paths overlap.

### Audit findings

Add stable findings to the existing JSON schema:

- `adr.baseline.missing` — canonical Baseline ADR file absent;
- `adr.baseline.stale` — marker version or digest differs from the asset;
- `adr.baseline.state-ambiguous` — partial or contradictory prior adoption;
- `adr.filename.invalid` — root numbered ADR cannot be parsed safely;
- `adr.number.duplicate` — more than one root ADR claims a number;
- `adr.number.overflow` — `+2` cannot remain four digits;
- `adr.rename.collision` — a planned destination is occupied;
- `adr.reference.ambiguous` — an explicit reference lacks one unique target;
- `adr.reference.protected` — a required edit falls under protected ownership;
- `adr.reference.stale` — future validation still finds an old reference;
- `adr.language.review-required` — non-blocking historical review inventory.

Missing, stale, ambiguous, invalid, overflow, collision, protected, and stale
reference findings are blocking. The language-review inventory is
informational. Audit remains read-only and produces the same ordered plan as
apply for the same repository state.

### Canonical ADR content

The two files under this Spec's `references/` directory are the approved
project-agnostic English bodies. Copy them into
`assets/templates/adrs/` without source-project product names, libraries,
frameworks, services, repository-specific paths, or migration assumptions.

ADR `0001` distinguishes documentation language from domain language: ADRs are
English, while ubiquitous terms follow domain experts and may record English
code aliases. ADR `0002` uses the glossary's Internal Identifier boundary and
requires UUIDv7 format, version, and RFC variant validation; existing data
migration remains project-owned.

### Documentation and skill synchronization

- Update root `CONTEXT.md` with **Context-Driven Baseline**, **Baseline ADR**,
  and **Internal Identifier** without adding implementation mechanics.
- Update generated root/domain guidance so agents read the two Baseline ADRs
  and write every new or modified ADR in English.
- Update `.agents/skills/setup-context-driven/SKILL.md` with the mandatory ADR
  audit, preview, exact adoption, and transactional migration behavior.
- Update the Context-Driven user guide and docs-layout guidance with the
  reserved identities and root-only namespace rule.
- Synchronize the canonical skill to `skills/setup-context-driven/` through
  `make skills-sync`; never author the embedded copy directly.
- Keep Spec 0036 and the canonical Roundfix Skill aligned with the fact that
  Doctor verifies the installed skill version but does not inspect project
  ADRs.

## Coverage Map

- Goals 1–3 and Stories 1–4 → first-class ADR assets and canonical templates.
- Goals 4–6 and Stories 5–7 → managed blocks, exact adoption, root inventory,
  `+2` mapping, reference rewrite, and transactional apply.
- Goals 7–9 → audit findings, root-only namespace, instruction/skill sync, and
  the explicit Doctor boundary.
- Historical language requirement → English baseline templates plus
  informational manual-review inventory without heuristic translation.

## Testing Approach

- Asset tests cover the new module schema, ADR template kind, safe canonical
  paths, unique artifact IDs/paths, rule references, and all three profiles.
- Planner unit tests cover empty roots, gaps, exact template adoption, managed
  current/stale blocks, divergent reserved occupants, duplicate numbers,
  malformed names, overflow, collisions, case-fold collisions, and mixed
  partial state.
- Reference table tests cover exact filenames, relative Markdown paths,
  `ADR-NNNN`, `ADR NNNN`, plural slash lists, false-positive numbers, Specs,
  context-local ADRs, protected paths, binary files, symlinks, and stale
  post-plan scans.
- Apply tests use real temporary Git repositories and assert preview/apply
  parity, bytes outside managed blocks, newline and mode preservation,
  root-only scope, exact automatic adoption, and no second diff.
- Failure-injection tests exercise every staging, replacement, deletion, and
  mode-restoration boundary and compare a full pre/post repository snapshot.
- Macro fixtures cover TypeScript/Bun monorepo, Go CLI/TUI, and Rust CLI
  profiles with both new and mature ADR trees.
- Run the setup Python suite, asset validation for canonical and embedded
  copies, `make skills-sync-check`, `git diff --check`, and `make verify`.

## Build Order

1. Add the approved English ADR templates, first-class asset schema, module
   declarations, and asset validation tests.
2. Implement baseline classification, root ADR inventory, rename-map
   validation, and focused planner tests.
3. Implement Git-tracked reference discovery, exact rewrite grammar,
   protected-path handling, and future-tree stale-reference validation.
4. Extend preview and atomic apply with rename/delete operations, mode
   preservation, rollback, and idempotency coverage.
5. Integrate blocking audit findings, automatic exact adoption, and the manual
   language-review inventory across all profiles.
6. Refresh generated instruction templates, `CONTEXT.md`, user docs, and the
   canonical setup skill; synchronize the embedded copy and run full QA.

## Risks & Considerations

- Renumbering touches historical links and prose broadly. Exact grammars and a
  future-tree stale-reference scan prevent silent partial migration.
- A repository may contain generated or upstream-managed files with ADR
  references. Protected references block instead of violating ownership.
- Filesystem case behavior differs by platform. Validate both exact and
  case-folded targets before staging.
- Existing atomic writes did not need to preserve executable modes because
  they targeted Markdown and JSON. Reference edits can touch scripts, so mode
  preservation becomes part of the transaction contract.
- Language detection is unreliable without external dependencies. The setup
  enforces English in its owned content and instructions and reports historical
  review work without claiming automated language proof.
- Installing ADR `0002` can expose existing UUIDv4 or ad hoc identifiers. The
  setup records the architectural invariant but does not pretend documentation
  audit proves code or data conformance.
- Fixed identities are intentionally exceptional to append-only ADR numbering.
  They belong to the portable baseline; all later project ADRs continue from
  the migrated sequence.

## Decisions

- Add a first-class ADR artifact type rather than disguise decisions as agent
  guides.
- Keep both ADRs unconditional through the existing `context-workflow` module.
- Use managed blocks and exact unmarked auto-adoption.
- Migrate only the root ADR namespace and preserve context-local namespaces.
- Use Git's tracked-file inventory as the repository-ownership boundary for
  reference rewriting.
- Treat ambiguous or protected references as blockers.
- Validate and apply the complete rename/reference/content plan atomically,
  preserving bytes and modes.
- Do not add a language detector or translate historical content.
- Do not extend Doctor with project-document inspection.

## Open Questions

None.
