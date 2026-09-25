---
task: task_03
spec: 0163-baseline-decisions-and-regeneration
status: completed
type: backend
complexity: medium
---

# Task 03: Skill regeneration declares its outputs

## Overview

`OutputsFor` in `internal/baseline/derived_ownership.go` reads ownership records only below `internal/baseline`, so a grant naming `make skills-sync` covers none of the mirrors that command writes and every grant enumerates them.

## Requirements

1. MUST add `skills/_ownership.yml` declaring `owner: dedicated` and `command: make skills-sync`.
2. MUST make `OutputsFor(root, "make skills-sync")` return every regular file under `skills/<owned>/` for each member of the Makefile's `OWNED_SKILLS`, reading the Makefile without writing it.
3. MUST exclude `.agents/skills`, Go files, `testdata`, `recommended.txt` and any skill directory not in `OWNED_SKILLS` from that set.
4. MUST keep `OutputsFor` for `make baseline-digests` returning exactly what it returns today.
5. MUST make the mechanical audit accept a changed owned mirror under a grant naming `make skills-sync`, and still refuse it under a grant that names neither the path nor the command.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The `make skills-sync` output set equals the owned mirrors' files in a fixture tree.
- [ ] Authorial, Go, `testdata`, `recommended.txt` and unowned files are excluded.
- [ ] `make baseline-digests` outputs are unchanged.
- [ ] The audit accepts a mirror change under `make skills-sync` and refuses it without.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/derived_ownership.go`
- interface: `internal/speccheck/mechanical.go`
- creates: `skills/_ownership.yml`

## Verification

- `out="$(go test -count=1 -v -run "^(TestOutputsForSkillsSyncListsOwnedSkillMirrors|TestOutputsForSkillsSyncExcludesAuthorialAndUnownedFiles|TestOutputsForBaselineDigestsUnchangedBySkillOwnership|TestMechanicalAuditAcceptsSkillMirrorUnderSkillsSync|TestMechanicalAuditRefusesSkillMirrorWithoutSkillsSync)$" ./internal/baseline ./internal/speccheck 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestOutputsForSkillsSyncListsOwnedSkillMirrors TestOutputsForSkillsSyncExcludesAuthorialAndUnownedFiles TestOutputsForBaselineDigestsUnchangedBySkillOwnership TestMechanicalAuditAcceptsSkillMirrorUnderSkillsSync TestMechanicalAuditRefusesSkillMirrorWithoutSkillsSync; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Skill regeneration ownership
- [2026-09-08-skill-regeneration-declares-its-owned-outputs.md](../0121-baseline-decisions-and-complete-regeneration/references/2026-09-08-skill-regeneration-declares-its-owned-outputs.md)

## Result

### Implementation

- `skills/_ownership.yml` now declares the shipped skill mirrors as dedicated
  outputs of `make skills-sync`.
- `OutputsFor` handles that command by validating the ownership record, reading
  `OWNED_SKILLS` from the Makefile, and returning the sorted regular files under
  only those `skills/<owned>/` directories. Other commands retain the existing
  `internal/baseline` ownership resolution path.
- The mechanical audit needs no production change: its command-only grant path
  already consumes `OutputsFor`, and the new integration regressions prove both
  acceptance of an owned mirror and refusal without path or command authority.

### Focused checks

- Before the implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^TestOutputsForSkillsSyncListsOwnedSkillMirrors$' ./internal/baseline`
  failed because `OutputsFor` tried to resolve `make skills-sync` only below the
  fixture's `internal/baseline/derived` tree. The matching mechanical audit
  regression failed through the same unresolved-output path.
- After the implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^(TestOutputsForSkillsSync|TestOutputsForBaselineDigests)' ./internal/baseline`
  and
  `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^TestMechanicalAudit' ./internal/speccheck`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/baseline ./internal/speccheck`
  passed (`internal/baseline` 55.255s; `internal/speccheck` 18.921s).
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache make verify-incremental`
  passed after the existing GitHub-backed test boundary received network access.
  `rtk git diff --check` also passed.

### Acceptance evidence

1. `TestOutputsForSkillsSyncListsOwnedSkillMirrors` compares the complete sorted
   set of files below two Makefile-owned mirror directories, including a nested
   regular file.
2. `TestOutputsForSkillsSyncExcludesAuthorialAndUnownedFiles` proves that the
   canonical `.agents/skills` file, root Go file, `skills/testdata`,
   `recommended.txt`, and an unowned skill directory do not enter the set.
3. `TestOutputsForBaselineDigestsUnchangedBySkillOwnership` proves the existing
   command still returns exactly its fixture's Baseline-owned regular file and
   no skill mirror.
4. `TestMechanicalAuditAcceptsSkillMirrorUnderSkillsSync` proves a command-only
   grant admits an owned mirror change, while
   `TestMechanicalAuditRefusesSkillMirrorWithoutSkillsSync` proves a grant naming
   neither the mirror nor the command still raises the path-escape finding.

The Task's declared `## Verification` command was not run; Verification remains
Daemon-owned.

## Carry-forward provenance

- Source Run: `run_20260925T135958Z_cb4c08f055bb883c`
- Source commit: `65ba70b79df26d43966b1d93bb0ddcf5e8277780`
