---
task: task_05
spec: 0073-skill-versions-decoupled-from-the-binary
status: completed
type: chore
complexity: low
---

# Task 05: Synchronise the Roundfix Skill

## Overview

This Spec changes what `roundfix doctor` reports and what blocks a command, so
the Skill must teach the new contract before the Spec can close. This is one of
the two authorized tooling Tasks.

## Requirements

1. MUST document that owned-skill readiness is a comparison against a declared
   minimum, not a content match, and that a version at or above the minimum
   satisfies.
2. MUST document the three distinct states — satisfies, below the minimum, and
   unversioned-or-unresolvable — and state that an unreachable source is never
   reported as a missing skill.
3. MUST document that a below-minimum skill blocks and names the skill, the
   minimum, the version found, and the upgrade path.
4. MUST state plainly that third-party skills are never held to a version
   Roundfix invented for them.
5. MUST state that editing an owned skill no longer requires a regeneration
   step.
6. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`, run
   `make baseline-digests`, and re-record the two characterization corpora that
   command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the comparison, the three states, and the third-party boundary.
- [ ] Run `make skills-sync`, then the regeneration chain.

## Acceptance Criteria

- [ ] The Skill states readiness is a comparison against a declared minimum.
- [ ] The Skill names all three states and their distinctions.
- [ ] The Skill states the four facts a below-minimum failure names.
- [ ] The Skill states third-party skills are never held to a version.
- [ ] The Skill states an owned-skill edit needs no regeneration step.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0073-skill-versions-decoupled-from-the-binary/task_05\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 4, 5 and 6.
- `_techspec.md` → Integration Points; Build Order 5.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.

## Result

### Implementation

- Replaced the Roundfix Skill's owned-skill content-match description with the
  shipped declared-version-versus-declared-minimum comparison. The reference
  names `satisfies`, `below minimum`, and `unversioned or unresolvable` as
  distinct states and says an unreachable source is never reported as a
  missing skill.
- Documented that `below minimum` blocks and names the skill, minimum, version
  found, and `roundfix skills install --target project` upgrade path.
- Documented that the owned-version comparison never applies to third-party
  skills and that Roundfix never holds them to a version it invented.
- Replaced the obsolete source-repository digest guidance: owned-skill content
  edits need no Baseline digest or characterization-corpus regeneration for
  compatibility readiness, while `make skills-sync` still maintains the
  canonical and embedded source copies.
- Regenerated `skills/roundfix/` from `.agents/skills/roundfix/`. The sanctioned
  digest and characterization commands produced no additional changed files.

### Focused checks

- `rtk make skills-sync` exited 0.
- `rtk make baseline-digests` exited 0 and reported that derived artifacts
  already matched their canonical sources (`changed: false`).
- The first sandboxed
  `rtk go test ./internal/baseline -count=1 -run
  TestBaselinePlanCharacterization -update-baseline-plan-characterization`
  attempt could not read the host Go build cache. The authorized rerun exited
  0 with 7 tests passing.
- `rtk go test ./internal/baseline -count=1 -run
  TestCatalogDiagnosticCharacterization -update-catalog-diagnostics` exited 0
  with 2 tests passing.
- `rtk proxy diff -rq .agents/skills/roundfix skills/roundfix` exited 0 with no
  output, proving the two directories are byte-identical.
- `rtk git diff --quiet HEAD -- '*.go'` exited 0. The complete status listed
  only the canonical Skill, its mirror, and this Task file.
- `rtk git -c core.fsmonitor=false diff --check` exited 0 before this Result
  record was appended; it is rerun after the final Task-file edit.
- The Task's declared `## Verification` commands, including `make verify`, were
  not run. They remain Daemon-owned.

### Acceptance evidence

- **Declared minimum comparison:** the Skill states readiness compares the
  installed declared version with the running binary's declared minimum, not
  content, and that a version at or above the minimum `satisfies`.
- **Three distinct states:** the Skill names `satisfies`, `below minimum`, and
  `unversioned or unresolvable`, including the unreachable-source distinction
  from a missing skill.
- **Four-fact failure:** the `below minimum` entry names the skill, minimum,
  found version, and project-install upgrade path.
- **Third-party boundary:** the Skill states the owned-version comparison never
  applies to third-party skills and Roundfix never invents their version.
- **No regeneration for readiness:** the Baseline section states an owned-skill
  content edit no longer requires Baseline digest or characterization-corpus
  regeneration because readiness compares versions rather than bytes.
- **Byte-identical mirror:** the recursive directory comparison exited 0 with
  no output after synchronization.
- **Repository Verification:** not exercised by the Agent; the Daemon must run
  the declared `make verify` command before settlement.
- **No Go source change:** the Go-only diff check exited 0, and the full status
  audit contains no `.go` path.
