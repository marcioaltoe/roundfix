---
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
prd: _prd.md
created: 2026-07-21
---

# Context-Driven Baseline coverage and Repository Skill Set restoration — Technical Spec

## Executive Summary

Extend the repo-owned `setup-context-driven` asset catalog so a selected profile proves semantic rule coverage, exact managed references, complete mutation planning, and restorable external skills. The design keeps the compact-root architecture: detailed portable rules remain in setup-owned guides, while project architecture and policy remain repository-authored. The primary trade-off is a stricter two-pass mutation flow: planning gains one extra invocation, but a deterministic plan digest binds user authorization to the exact file-tree preimage and prevents preview/apply drift. Audit stays local and read-only; network access occurs only in the explicit skill-restoration command.

## System Architecture

- `.agents/skills/setup-context-driven/assets` remains the declarative authority. Profiles gain required rule sets; modules gain renderable rule guidance, exact skill-dispatch mappings, and typed references; setup snapshots gain immutable external source records and complete-directory digests.
- `context_assets.py` loads the versioned schemas and rejects incomplete semantic coverage, unreachable rule IDs, dispatch gaps, unresolved reference targets, mutable provenance, and unsafe paths before repository inspection begins.
- `context_setup.py` keeps `DecisionPlan` as the composition model, but makes one concrete `ChangePlan` the authority for both public `plannedChanges` and writes. A plan digest binds selection, resolved decision values, catalog content, logical operations, and path preimage/postimage digests.
- The existing audit surface validates the selected future artifact graph as well as current files. The apply surface becomes preview-first and requires the exact plan digest before any write. A new `restore-skills` surface restores only missing or drifted external members of the selected Repository Skill Set.
- Canonical `.agents/skills/setup-context-driven` remains the only authorial source. `make skills-sync` refreshes `skills/setup-context-driven`; no upstream-managed skill content is edited or bundled into Roundfix.
- No Roundfix Go CLI, Doctor Command, Baseline ADR set, nested package instructions, or general skill-package management changes are included. Spec 0036 retains ownership of Doctor Skill Readiness and its `skills-lock.json` compatibility proof.

The flow is:

```text
assets + decisions + repository preimage
              -> DecisionPlan
              -> concrete ChangePlan + planDigest
              -> confirm same digest
              -> staged write + postwrite delta proof
```

## Implementation Design

### Interfaces

Catalog rule and reference contracts become explicit inputs to rendering:

```python
@dataclass(frozen=True)
class RuleContract:
    rule_id: str
    coverage: tuple[str, ...]
    guidance: str

@dataclass(frozen=True)
class ArtifactReference:
    token: str
    target_managed_id: str | None
    repository_path: Path | None
    ownership: str
```

One exact plan owns the observable preview and the executable bytes:

```python
@dataclass(frozen=True)
class FileMutation:
    path: Path
    before_digest: str | None
    after_digest: str | None
    content: bytes | None
    operations: tuple[PlannedChange, ...]

@dataclass(frozen=True)
class ChangePlan:
    kind: str
    mutations: tuple[FileMutation, ...]
    digest: str
```

Acquisition and lock compatibility stay replaceable test seams:

```python
class SkillSource(Protocol):
    def acquire(self, source: SkillSourceRef, target: Path) -> None: ...

class SkillLockAdapter(Protocol):
    def entry_for(self, source: SkillSourceRef, tree: Path) -> dict: ...
```

The production source uses argv-only Git subprocesses with prompting disabled. The lock adapter operates only on an already acquired and verified temporary tree; it cannot select a branch, fetch content, or write the target repository.

### Data Models

The asset contract adds a stable `coverage` catalog. Every profile declares `requiredRules`; every rule declares one or more coverage category IDs plus its portable Markdown guidance. Every required rule must belong to a selected module, be carried by a selected supporting guide, and render through that guide's required `{{artifact.rules}}` binding. This makes generated text derive from the rule authority instead of allowing metadata and prose to drift. Universal categories cover safety, selected Verification, verification-configuration integrity, skill dispatch, research sources, dependency changes, Git and delivery, and security/configuration. Language and enabled application-surface categories are profile-dependent.

Each module declares `skillDispatch` entries whose keys exactly equal `requiredSkills`. Values describe when the skill is mandatory; duplicate skill names across active modules retain all distinct triggers. A new `guide.skill-dispatch` renders the ordered, deduplicated map through the computed `active-modules.skill-dispatch` binding. `verification.gate` becomes an entry decision for every profile and renders into the universal agent-instructions guide, independent of autonomous-work selection.

Root blocks and guides gain typed `references`. A setup-owned reference names a target managed ID and a template token; a repository-owned reference names a safe repository-relative path. Template loading proves the token exists and binds it to the declared target. Decision resolution rejects a definite source whose setup-owned target is absent from the exact definite artifact set. Audit checks repository-owned targets against the repository and rejects paths outside it. Existing Markdown-link scanning remains defense in depth, but an unrelated stale file cannot satisfy a managed reference. `guide.monorepo` becomes unconditional whenever the monorepo module is active, and frontend guidance declares the repository-owned `DESIGN.md` contract without generating it.

`PlannedChange` retains `action`, `path`, `managedId`, `state`, and optional `condition`, and adds a required `reason`, optional `fromPath`, and typed reference-edit detail. `FileMutation` groups all logical operations affecting one path and carries SHA-256 digests of the exact bytes before and after the change; absence uses `null`. A managed ID moving between paths is reported as a rename rather than an unrelated deletion and creation. The plan digest is SHA-256 over canonical JSON containing plan kind, profile/setup selection, resolved decision values, catalog digest, ordered public operations, and preimage/postimage digests. Volatile timestamps and the digest field itself are excluded.

Only the old Setup Manifest inventory plus valid ownership markers can authorize automatic removal. An unmarked path outside that inventory is never a mutation candidate. An inventoried path with missing or invalid ownership markers blocks until a dynamic `removal.<managed-id>` decision records either `preserve` or `remove`; either outcome appears in the recomputed plan before confirmation. Adoption decisions continue to govern unmarked content that a selected artifact would begin managing.

External setup entries move to `setup-context-driven/setup-snapshot-v2` and contain a safe GitHub `owner/repo`, a full immutable commit object ID, a source-relative directory, and a lowercase complete-tree digest. The snapshot digest covers canonical serialization of the complete normalized records, not only paths. The portable tree digest hashes every regular file in bytewise POSIX-path order with length framing for paths and contents, excludes `.git` and `node_modules`, and rejects symlinks and special files. Roundfix-owned `source.type: repo` entries remain distinct and keep their existing ownership-safe synchronization rule.

`Finding` gains optional `remediation`. Missing or drifted external-skill findings include provider, skill name, exact source/ref/path, expected tree digest, and preview argv. Audit computes the same complete-tree digest offline and never calls Git or the lock adapter. The external `skills-lock.json` `computedHash` remains compatibility metadata rather than the setup content authority; Spec 0036's Doctor fixture and the isolated lock adapter must agree before restoration can ship.

### API Contracts

`audit` accepts repeated `--decision ID=VALUE` so maintainers and agents can inspect a fully resolved Decision Plan without writes. Existing `setup-context-driven/audit-v1` fields remain; `reason`, `remediation`, and `planDigest` are additive. Unresolved decisions still return conditional possibilities and exit `3`; only a fully resolved plan carries an authorizable digest.

`apply` uses the existing profile and decision inputs plus `--confirm-plan <sha256>`. With non-empty changes and no confirmation, it writes the complete plan, emits `plan.confirmation.required`, and exits `3`. A malformed digest exits `2`. A repository, decision, catalog, or source-state change produces a new plan, emits `plan.confirmation.stale`, performs no write, and exits `3`. A matching digest stages all bytes, applies them, verifies the observed affected-path delta equals the authorized mutations, and rolls back on mismatch or I/O failure. An empty idempotent plan exits `0` without requiring confirmation.

`restore-skills --repo PATH --profile ID [--skill NAME ...] --format text|json [--source-dir PATH] [--confirm-plan SHA256]` selects only missing or drifted external skills required by that profile. Without `--skill`, it selects all such required skills. Preview acquires each unique `(source, ref)` once into temporary storage, rejects a commit mismatch, verifies every source subtree digest, and reports created, refreshed, and removed files plus the lock edit under schema `setup-context-driven/restore-v1`. `--source-dir` may supply an offline Git object store but is never persisted. Confirmation recomputes the plan, atomically swaps staged sibling directories and the lock temp file, preserves unrelated lock entries, verifies final tree digests, and rolls everything back on failure. Unsupported provider, unavailable exact commit, incompatible lock adapter, unsafe archive entry, or digest mismatch exits `1` before repository mutation; invalid arguments exit `2`. There is no fallback to a branch or default revision.

## Coverage Map

- Goal 1 → profile `requiredRules`, coverage catalog, rule reachability validation, and computed rule rendering.
- Goal 2 → compact root blocks, setup-owned supporting guides, ownership language, and repository-owned design reference.
- Goal 3 → typed references, exact Decision Plan graph validation, and future-tree audit.
- Goal 4 → authoritative `ChangePlan`, plan digest confirmation, removal decisions, and postwrite delta proof.
- Goal 5 → snapshot v2 provenance, complete-tree audit, structured remediation, and `restore-skills`.
- Story 1 → mandatory universal and applicable coverage categories across every profile.
- Story 2 → exact `requiredSkills`/`skillDispatch` equality and generated dispatch guide.
- Story 3 → managed-ID references resolved against definite selected artifacts, not incidental files.
- Story 4 → complete logical operations and path digests exposed before `--confirm-plan` writes.
- Story 5 → exact-commit acquisition, tree verification, portable lock update, and idempotent restoration.
- Story 6 → explicit setup/repository ownership in generated guidance, references, and manifest behavior.

## Integration Points

- Git is the immutable-source adapter for external restoration. It runs non-interactively against declared GitHub identities or an explicit local object store and never executes downloaded content.
- `skills-lock.json` remains an external compatibility boundary. A pinned, isolated adapter derives the versioned lock entry from verified temporary bytes; normalized GitHub identity, full ref, and source path come from the bundled snapshot, never from machine-local paths.
- The repository filesystem remains the apply boundary. Existing staging and rollback behavior expands to directory swaps and postwrite digest verification.
- Spec 0036 remains the Doctor integration boundary. Its Repository Skill Set lock-hash compatibility fixture must agree with the restoration adapter; this Spec does not add Doctor checks.
- Canonical-to-embedded skill synchronization remains the distribution boundary through `make skills-sync` and its existing equality checks.

## Testing Approach

- Asset mutation tests remove one required rule, rule carrier, dispatch mapping, reference target, template token, immutable ref, safe source path, or tree digest and assert stable loader diagnostics. Each profile must satisfy all universal and applicable coverage categories.
- Decision tests enumerate every finite boolean/enum branch that includes or excludes an artifact for every profile. Every definite setup-owned reference must resolve inside the selected set; every declared repository-owned reference must remain inside the repository.
- Rendering tests prove the selected Verification appears in universal guidance, dispatch equals the active modules' required skills, frontend points to repository-owned `DESIGN.md`, and repository-authored bytes outside markers remain unchanged.
- Change-plan tests compare public operations and unique `(path, beforeDigest, afterDigest)` triples with a complete before/after tree snapshot. Fixtures include Go-to-Rust omitted removals, the unmarked false-positive removal, marked conditional cleanup, missing legacy markers, adoption, explicit preserve/remove decisions, shared root blocks, rename, reference edits, stale confirmation, rollback, and idempotent reapply.
- Tree-digest tests cover nested edits, additions, removals, ordering, empty files, excluded directories, symlinks, and special files. A cross-Spec fixture distinguishes the portable tree digest from external lock compatibility metadata.
- Restoration unit tests fake only acquisition and lock-adapter boundaries. Integration tests use a disposable local Git repository and the real command path to prove exact-commit acquisition, removal preview, multi-skill fetch grouping, portable lock provenance, stale-plan rejection, rollback, no mutation on digest/fetch/adapter failure, final audit success, and an empty second restore.
- Macro tests apply every profile, traverse each artifact-changing decision transition, audit clean, and reapply with no changes. Run `make skills-sync`, then the mandatory `make verify`; canonical and embedded setup skills must match.

## Build Order

Cross-Spec integration gate: complete Spec 0036 Task 01's lock compatibility fixture before enabling `restore-skills` writes. Its Doctor behavior remains outside this graph.

1. Versioned coverage, rule-rendering, dispatch, typed-reference, snapshot-v2, and tree-digest asset contracts with mutation tests.
2. Portable baseline rule content, universal Verification decision/rendering, generated skill-dispatch guide, frontend design pointer, and unconditional monorepo guide (depends on: 1).
3. Exact Decision Plan reference validation, repository-owned reference audit, and full profile/decision reference fixtures (depends on: 1, 2).
4. Sole-source concrete `ChangePlan`, audit decision inputs, enriched operations, explicit removal authority, plan digest confirmation, atomic delta verification, and transition parity tests (depends on: 1, 3).
5. Provenance-aware `sync-setups`, complete-directory installed-skill audit, and structured remediation (depends on: 1, 4).
6. Exact-commit acquisition, isolated lock adapter, `restore-skills` preview/confirmation/apply, atomic directory rollback, and restoration integration tests (depends on: 4, 5 and the Spec 0036 integration gate).
7. Macro profile flows, maintainer/user documentation, canonical skill update, embedded synchronization, and repository verification (depends on: 2, 3, 4, 5, 6).

## Risks & Considerations

- A plan can become stale between preview and confirmation. Preimage digests and full recomputation turn this into a safe exit `3`, never implicit reauthorization.
- External repositories are untrusted. Restoration limits provider/path/ref shapes and tree size/file count, rejects traversal, links, devices, and special files, disables Git prompting, and never executes acquired content.
- Current snapshots contain many empty or mutable external sources. Migration must prove bytes at the declared commit; it must block rather than attach a convenient commit to unmatched working-tree content.
- The external skills tool's lock hash ordering is not the portable content proof. Isolation prevents it from choosing source bytes, while the bundled tree digest remains audit authority. If the pinned lock adapter and Spec 0036 disagree, restoration writes stay disabled until that owning contract is corrected.
- Adding `verification.gate` as a universal entry decision creates a one-time question for existing non-autonomous manifests. Stored compatible values migrate automatically; no default command is invented.
- Repository-owned `DESIGN.md` is required only for frontend-enabled profiles. Setup reports its absence but never creates project-specific design content.
- Snapshot v2 and new template versions ship together with the script. Existing repositories keep their manifest and unmarked content; rollback is reinstalling the prior skill version or reverting setup-owned changes through repository version control. Restoration rollback is automatic within the command.
- Audit and documentation apply remain network-free. Only explicit restoration pays fetch and temporary-storage cost, bounded by acquisition limits and grouped by source/ref.

## Decisions

- Preserve compact root pointers and render detailed portable rules into existing setup-owned guides from stable rule contracts.
- Require exact module skill-dispatch coverage and generate one deterministic dispatch guide from the active Decision Plan.
- Make selected Verification universal rather than conditional on autonomous work.
- Use typed managed-ID and repository-path references as the primary contract; keep Markdown scanning only as defense in depth.
- Make one concrete `ChangePlan` authoritative for preview and apply, with digest-bound confirmation and observed-delta verification.
- Never remove unmarked content without a presented adoption or removal decision.
- Restore external skills only from immutable portable provenance, with complete-tree verification and no branch fallback.
- Keep setup tree integrity independent from external lock compatibility metadata.
- Add no ADR because ADR-0046 and ADR-0047 already govern declarative ownership, Decision Plan effects, and shared preview/audit/apply behavior.
