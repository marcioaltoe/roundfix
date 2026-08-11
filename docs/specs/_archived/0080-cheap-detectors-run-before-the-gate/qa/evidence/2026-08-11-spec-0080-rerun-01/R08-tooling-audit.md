# R08 — protected-tooling authorization and chronology

Build: `c2372a9f709c9197aa5c5e89fd71da1ab46f07e6`.

## Authorization chronology

The exact grant is
`docs/workflow/authorizations/2026-08-06-proof-cost.md`. Its corrected record
landed in commit `58b6881d` on 2026-08-07. `git merge-base --is-ancestor`
exited 0 for authorization → Task 04 (`a7df02ba`), authorization → Task 06
(`9d36349a`), Task 06 → consequent Task 07 (`1d863255`), and Task 07 →
corrective Task 09 (`c2372a9f`). Authorization, protected edits, consequent
fix, and later harness fix are separate chronological commits.

## Actual Task commit paths

Every Task commit was read with
`git diff-tree --no-commit-id --name-only -r <commit>`.

- Task 01 `98e1467e`: `internal/speccheck/**`, mechanical fixtures, and its
  Task file only.
- Task 02 `88453ed6`: Agent prompt/Daemon context Go files and its Task file.
- Task 03 `1432dd1e`: Daemon/mechanical Go files and its Task file.
- Task 04 `a7df02ba`: exactly `.agents/skills/qa-gate/SKILL.md`,
  `skills/qa-gate/SKILL.md`, and its Task file.
- Task 05 `78d53524`: mechanical/report Go files and its Task file.
- Task 06 `9d36349a`: the granted Makefile, two modules, two adopted guides;
  `docs/agents/setup-context.json`; sanctioned Baseline assets/testdata; and
  its Task file. No `.github/**` path appears.
- Task 07 `1d863255`: only
  `internal/baseline/preservation_test.go` and its Task file. The one-line
  declared entry expectation moves 132 → 134; no authorized cause path is
  folded into this consequent commit.
- Task 09 `c2372a9f`: only two test files and its Task file.

`git diff --quiet` confirmed Task 06 did not change
`internal/baseline/preservation_test.go`, and Task 07 did not change any Task
06 protected cause path. The Task 06 Makefile diff adds only
`verify-incremental` and its `.PHONY` token; the existing `verify` recipe is
unchanged.

## Regeneration and postimage evidence

- `rtk make skills-sync-check` exited 0.
- Both qa-gate skill copies have SHA-256
  `7a3ee44d32fe9644ffd5b8acc1963abbe96e07fdbe702d69525c7c4a4892ca08`.
- Focused `TestCatalogCompatibility`, `TestFormatterComposition`,
  `TestCatalogRegenerationMode`, and `TestDerivedOwnershipIsExhaustive` cases
  exited 0.
- The copy-based repository-contract test
  `TestDeclaredStepRegenerationAndFrozenBoundaries` ran with the
  `repocontract` build tag and exited 0 after 25.14 seconds. Its sanctioned
  command rewrote sanctioned artifacts, preserved frozen boundaries, and its
  negative carriers rejected absent, wrong, unchanged, and frozen-path
  regeneration shapes.

The complete command audit found no missing, late, folded, untraceable, or
out-of-bound protected-tooling change.
