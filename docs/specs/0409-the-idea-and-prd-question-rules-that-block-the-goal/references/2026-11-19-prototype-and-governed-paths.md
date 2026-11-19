# Prototype of the Rule Challenge change, and its Governed Paths

Recorded by the author of the PRD and the TechSpec on 2026-11-19, on the Mac mini. The prototype ran on a scratch worktree of this repository at `348962d2` (the delivery of Spec 0408), outside the planning clone, and the worktree was removed afterwards. It answers three questions: where a new Baseline clause can live, which files the change and its regeneration rewrite, and which of those files are Governed Paths.

## 1. What was applied

The exact texts of the TechSpec's API Contracts 1 to 9, applied by a script: the rule with its two clauses in `internal/baseline/assets/modules/spec-workflow.json`, the guide's rule list and guide version, two lines in the closed list of `internal/baseline/clause_characterization_test.go`, the canonical edits of four skills, two templates and the Roundfix skill's `modes` reference, two glossary entries, the user guide subsection, and ADR-0436 and ADR-0437. The two new test files the TechSpec names were not written.

## 2. A new clause in the route rule is refused

The first attempt added the two clauses to `rule.spec.routing`. `make baseline-digests` stopped at `./internal/baseline:TestReadoptionCompatibilityMaintainedFixture` with exit 2 and this diagnostic, once per clause:

```text
catalog.sourceBaseline.required-clause.missing: baseline.standard-typescript-monorepo-0.0.1: clause.spec.rule-challenge-01-ask; the regenerator maintains manifest rows but never creates them, so add this row first
```

The check is `validateSourceBaselineCoverage` in `internal/baseline/catalog_validate.go`: for every rule a profile lists under `requiredRules`, each clause needs an entry in the profile's Source Baseline manifest. The Standard TypeScript profile lists 33 required rules. Among the rules of the modules it selects, these are required: `rule.spec.routing`, `rule.spec.artifacts`, `rule.spec.docs-layout`, `rule.context.domain-docs`, `rule.context.docs-layout` and `rule.core.decision-authority`. These are not: `rule.spec.project-constraints`, `rule.core.cross-project-continuity`, `rule.autonomous.hook-strictness`, `rule.frontend.recorded-layout` and `rule.secondbrain.research-digests`.

The second attempt put the same two clauses in a new rule, `rule.spec.rule-challenge`, listed in `guide.spec-routing` between `rule.spec.routing` and `rule.spec.project-constraints`. Everything below is that attempt. No file under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/` changed.

## 3. The commands and the files each one rewrote

Run in this order, each with exit 0.

| command | files it rewrote |
| --- | --- |
| `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1` | `internal/baseline/module-versions.json`, and the version line of `internal/baseline/assets/modules/spec-workflow.json`, from 27 to the next |
| `make baseline-digests` | `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/testdata/catalog.digest`, `internal/baseline/testdata/catalog.normalized.json`, `internal/baseline/testdata/catalog.diagnostics.golden.json`, and under `internal/baseline/testdata/plan-characterization/`: `advisory-only-divergences.golden.json`, `clean-adoption.golden.json`, `idempotent-replan-after-verified-apply.golden.json` and `same-baseline-changed-profile-and-catalog-digests.golden.json` |
| `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` | `docs/agents/spec-routing.md` and `docs/agents/setup-context.json`; a second run printed `File changes: 0` |
| `make skills-sync` | the mirrors under `skills/` of the eight canonical skill files edited |
| `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` | `skills/testdata/owned-skill-versions.json` and the version fields of both copies of the `SKILL.md` of `brainstorming`, `roundfix`, `write-idea`, `write-prd` and `write-techspec` |
| `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` | `docs/references/behavior-surfaces.json`, one fingerprint: `guide docs/user-guide/context-driven-development.md` |

After the skill edits, `make baseline-digests` reported `no changes; derived artifacts already match their canonical sources`, and a third Managed Refresh printed `File changes: 0`: the skill texts move no Baseline file and not the Setup Manifest. The Managed Refresh also wrote the Roundfix version it ran with into the Setup Manifest. In the rendered Spec routing guide the two clauses follow the last route clause and precede the Project Constraints clauses.

The change is 34 files, with 412 lines added and 108 removed, before the two test files and the two ADRs:

- hand-edited: `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/clause_characterization_test.go`, `CONTEXT.md`, `docs/user-guide/context-driven-development.md`, and under `.agents/skills/`: `brainstorming/SKILL.md`, `write-idea/SKILL.md`, `write-idea/references/idea-template.md`, `write-prd/SKILL.md`, `write-prd/references/prd-template.md`, `write-techspec/SKILL.md` and `roundfix/references/modes.md`;
- written by the commands above: every other file, `.agents/skills/roundfix/SKILL.md` included, whose version fields the record command raised.

## 4. Tests on the prototype

Each ran with `GOMAXPROCS=2`, after the six commands, and exited 0:

- `go test ./internal/baseline -count=1`, with the two ids added to the closed list `characterizedBaselineForce`. It includes `TestAnUnlistedBaselineClauseIsReported`, `TestNoTwoBaselineClausesShareText`, `TestShippedGuidanceCitesNoRepositoryRecord`, `TestInstructionHierarchyRendersActivePointersOnce` and the tests that pin the domain clause, the question clause and the route clauses.
- `go test ./skills -count=1`. No sentence an existing test pins changed, and the `--follow-up` section kept eight steps.
- `go test -count=1 -tags docscontract ./internal/docscontract`, the glossary contracts included.
- `go test ./internal/speccheck ./internal/spec -count=1`.
- `go test ./internal/cli -run 'Baseline|GuidanceComposition' -count=1`.
- `make verify-docs`.

`make verify-changed` exited 0 on the same prototype with an earlier wording of the texts; the later wording changed only sentences inside the new clauses and the new sections.

## 5. Governed Paths

Measured on the planning clone with `speccheck.GovernedPath`, in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository.

Governed, and changed by this Spec:

- `.agents/skills/brainstorming/SKILL.md`, `skills/brainstorming/SKILL.md`
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`
- `.agents/skills/roundfix/references/modes.md`
- `.agents/skills/write-idea/SKILL.md`, `skills/write-idea/SKILL.md`
- `.agents/skills/write-idea/references/idea-template.md`
- `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`
- `.agents/skills/write-prd/references/prd-template.md`, `skills/write-prd/references/prd-template.md`
- `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`
- `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- `internal/baseline/assets/modules/spec-workflow.json`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`

Not governed, and changed or created by this Spec: `internal/baseline/module-versions.json`, `internal/baseline/clause_characterization_test.go`, `internal/baseline/rule_challenge_clauses_test.go`, the three catalog files and the four plan goldens under `internal/baseline/testdata/`, `skills/write-idea/references/idea-template.md`, `skills/roundfix/references/modes.md`, `skills/testdata/owned-skill-versions.json`, `skills/rule_challenge_authoring_test.go`, `CONTEXT.md`, `docs/user-guide/context-driven-development.md`, `docs/references/behavior-surfaces.json`, ADR-0436 and ADR-0437.

Controls, governed and not changed: the `Makefile`, `go.mod`, `.roundfixrc.yml`, `.github/workflows/ci-verify.yml`, `docs/agents/docs-layout.md`, `docs/agents/agent-instructions.md`, `docs/agents/specific-repository.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/context-workflow.json`, the Source Baseline manifest of the Standard TypeScript profile, `internal/baseline/plan_test.go`, `skills/baseline_skill_contract_test.go`, both copies of the TechSpec template and `.agents/skills/write-tasks/SKILL.md`. Controls measured not governed and not changed: `docs/agents/domain.md`, `internal/cli/baseline_update_test.go`, `docs/references/skill-coverage.json` and `AGENTS.md`.
