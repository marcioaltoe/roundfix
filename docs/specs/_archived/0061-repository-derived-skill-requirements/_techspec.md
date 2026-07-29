---
spec: 0061-repository-derived-skill-requirements
prd: _prd.md
created: 2026-07-29
---

# Repository-derived skill requirements — Technical Spec

## Executive Summary

`skills.CheckRepository` resolves its external requirement with
`external := Recommended()`, a fixed list embedded in the binary. The fix
moves that resolution to the caller, because the derivation needs the Baseline
catalog and the `skills` package must not depend on it — `skills` is the
embedded bundle and stays a leaf. The trade-off accepted is one new exported
entry point on `skills` plus a small resolver in `internal/cli`, rather than
inverting the dependency or duplicating the catalog inside `skills`.

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier
  is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 keep profile
  proof independent from skill readiness; ADR-0066 keeps Baseline execution in
  the Go CLI, which is why the resolver lives there. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-29 the maintainer expressly
  authorizes exactly `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, and the deterministic Skill-digest fallout
  `make baseline-digests` rewrites. No other protected tooling mutation is
  authorized. Source: `docs/agents/agent-instructions.md`.

## System Architecture

- `skills` keeps owning the embedded bundle and the readiness comparison. It
  gains an entry point that accepts the external requirement instead of
  computing it, and keeps the existing one for callers that have no manifest
  context.
- `internal/cli` gains the resolver: read the repository's Setup Manifest,
  load the embedded Baseline catalog, union the selected modules'
  `requiredSkills`, subtract the owned names, and pass the result to the
  readiness check. It already imports both packages.
- No new package. `skills` does not import `internal/baseline`, preserving the
  leaf position that makes the embedded bundle safe to reuse.

## Implementation Design

### Interfaces

```go
// skills — accepts the requirement instead of deciding it.
func CheckRepositoryWithExternal(ctx context.Context, root string, external []string) (RepositoryReadiness, error)

// internal/cli — derives the requirement from repository state.
// Returns an empty set and ok=false when the manifest is absent or unreadable.
func resolveExternalSkillRequirement(root string) (external []string, ok bool, err error)
```

`CheckRepository` keeps its signature and behavior for any caller that has no
repository context; the Doctor path stops using it.

### Data Models

No schema changes. Two documents are read:

- the Setup Manifest at `docs/agents/setup-context.json`, whose `modules`
  array names the selected modules;
- each module entry in the embedded catalog, whose `requiredSkills` array
  names the skills that module requires.

The external requirement is the union of those arrays minus `skills.Names()`,
deduplicated and sorted so the readiness line is deterministic.

### API Contracts

- `roundfix doctor` reports the same `skills:` line shape with counts derived
  per repository. With no manifest it reports zero external requirements and a
  next action naming Baseline adoption.
- A missing-external failure changes from first-miss-wins to naming every
  missing skill, each with `bunx skills add marcioaltoe/skills@<skill>`.
  The per-skill form is what the upstream CLI supports; the package-wide form
  the current message prints installs the whole catalog.
- The Doctor Command remains diagnosis-only.

## Coverage Map

- Story 1 and Feature 1 → `resolveExternalSkillRequirement` plus the Doctor
  wiring.
- Story 2 and Feature 4 → the missing-skill error carrying the full set and
  per-skill commands.
- Story 3 and Feature 3 → the absent-manifest branch.
- Story 4 → covered by the same derivation: this repository's manifest selects
  `go`, `cli-surface`, and `tui-surface`.
- Feature 2 and 5 → unchanged owned resolution and the untouched
  diagnosis-only contract, asserted by the existing no-mutation test.

## Integration Points

Filesystem only: the repository's Setup Manifest and the embedded catalog.

## Testing Approach

Existing seams. `skills/repository_test.go` covers the new entry point with an
explicit external set, including the empty set. `internal/cli/doctor_test.go`
covers the derivation: a manifest selecting TypeScript modules requires no Go
skill, a manifest selecting `go`/`tui-surface` requires them, an absent
manifest requires none and names adoption, and a missing skill lists every
gap with its per-skill command. The existing
`TestRunDoctorRealRepositoryCheckDoesNotMutateState` keeps proving the
diagnosis-only contract.

## Build Order

1. Add the external-accepting entry point in `skills`, keeping the existing
   one intact.
2. Add the Setup Manifest resolver and wire the Doctor path to it (depends
   on: 1).
3. Report every missing external skill with its per-skill install command
   (depends on: 1).
4. Align the Roundfix Skill pair and regenerate derived digests with the
   sanctioned command (depends on: 1, 2, 3).

## Risks & Considerations

- **Repositories that pass today may change requirement counts.** That is the
  point, and it is a readiness change, not a behavior change: no operational
  command consults this check.
- **A manifest listing a module the catalog does not know** must fail legibly
  rather than silently requiring nothing; the resolver reports the unknown
  module.
- **`autoresearch` and `exa-web-search` stop being required anywhere**, since
  no module declares them. They remain installed where present and remain
  advisory recommendations.

## Decisions

- The resolver lives in `internal/cli`, not in `skills`, so the embedded
  bundle stays a leaf package.
- No manifest means no external requirement — stated, not inferred.
- `recommended.txt` remains recommendation data and stops being the readiness
  authority.
