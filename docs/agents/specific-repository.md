# Agent instructions

This repository is Roundfix: a local-first Go CLI and daemon that picks up work
items — PR review issues and Spec Task Graphs — resolves them through the user's
selected ACP runtime, and pushes only when nothing unresolved remains. Stdlib
`flag` dispatch and a Bubble Tea v2 TUI.

Project map: `cmd/roundfix/` is the thin CLI entry point; behavior lives in
`internal/...` — `internal/cli/` owns parsing, output, and exit behavior;
`internal/app/` holds app metadata.

## Repository rules

- A top-level test in a Parallel Test Package calls `t.Parallel()` as its first statement or carries a `// Sequential: <reason>` comment naming the process-wide state it changes; the Parallel Test Package rule in `internal/testfixture` enforces it.
- Keep the project KISS: prefer the smallest behavior that satisfies the
  documented product contract.
- **NEVER** copy names, branding, package names, comments, examples, or
  generated artifacts from reference projects into this repository.
- **HARD RULE — repository contracts validate at the pull request boundary**:
  `make verify` excludes the checks whose inputs are the repository itself —
  markdown contracts and the derived-artifact regeneration gates — so an
  ordinary commit does not re-run them. `make verify-changed`, the Verification
  of Runs and of the QA gate, runs the Repository Contract Tests that their
  Contract Relevance selects (ADR-0252). `make verify-docs` runs the
  `internal/docscontract` tests, the contracts `repo-test` lists and `roundfix spec
  check`, and it **MUST** pass before any pull request opens. The Full Contract Run,
  `make verify-contracts`, runs every Repository Contract Test whatever its Contract
  Relevance, and CI runs it on every push to main and before every release (ADR-0253).
  Platform-only failures stay a known limit because CI is the
  Linux gate. Nothing under `docs/history/` is ever validated as live work.
- **HARD RULE — skill ownership**: repo-owned authorial workflow skills may be
  adapted locally; every other skill is upstream-managed and **MUST NOT** be
  modified here.
- An owned skill's content changes only together with its version. Run
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
  the record command raises both version fields one patch above the highest
  recorded version when needed, records the result, and requires
  `skills/testdata/owned-skill-versions.json` to be declared in the Task.
- A Baseline module's content changes only together with its version.
  After editing a module under `internal/baseline/assets/modules/`, run
  `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
  before `make baseline-digests`; the record command keeps a version above
  every recorded one or writes the next free version into the module's own
  version line, records it in `internal/baseline/module-versions.json`, and
  requires that record to be declared in the Task. A Spec names the next
  version, never a number. `docs/references/coverage-record.json` is
  re-recorded only with
  `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`,
  which writes the same bytes on any host.
- **HARD RULE — sanctioned digest regeneration**: after an expressly authorized
  Roundfix-owned Skill or Baseline module edit, run `make baseline-digests`.
  Every derived pin rewritten by that command is deterministic fallout of the
  authorized source edit and needs no separate express authorization. A
  hand-edited pin value remains an unauthorized mutation.
- **HARD RULE — spawning test packages install the suite guard**: a Go test package whose tests start a process installs `suiteguard.Main` in its `TestMain`, carries `TestEverySpawningPackageInstallsTheSuiteGuard` in a `suiteguard_repocontract_test.go`, and is listed in `guardedSpawningPackages` in `internal/suiteguardcontract/contract.go`. A Task that adds such a package or its first spawning test declares all three files.
- Never add or change text inside the `### QA settlement` section of a skill: `TestSettlementGuidanceIsOneTable` requires it byte-identical in `qa-gate`, `archive-spec` and `roundfix`. New skill text goes under its own heading.
- A reworded Baseline clause keeps its identity. A removed clause is named in the `replaces` list of the clause that absorbs it, with the same enforcement, so the Source Baseline transition (`classifySourceClauseTransition` in `internal/baseline/plan.go`) records it `replaced`. A new clause in a rule the Standard TypeScript profile requires needs a Source Baseline row, so prefer extending an existing clause. A Task that removes a clause declares `internal/cli/baseline_update_test.go` and `internal/baseline/plan_test.go` when their fleet fixtures pin it.
- The canonical copy of an owned skill is `.agents/skills/<name>/`, and `make skills-sync` copies it to `skills/<name>/`. A Task that changes an owned skill declares the canonical file, its mirror, the `SKILL.md` pair whose version fields rise, and `skills/testdata/owned-skill-versions.json`.
- Derived Baseline files (the profile digest pin, the catalog snapshots and plan goldens under `internal/baseline/testdata/`, the formatter goldens and `docs/agents/setup-context.json`) change only through `make baseline-digests` and the Managed Refresh. A Task that edits module, template or profile source declares every derived file those commands rewrite, measured on a scratch copy.
- A Spec that edits a Baseline module declares two commands in the Sanctioned regeneration section of its `_authorization.md`: `make baseline-digests`, and the module record command with an explicit `outputs:` list holding `internal/baseline/module-versions.json` and each changed `internal/baseline/assets/modules/<name>.json`, without globs. Suiteguard refuses the module and record writes of an undeclared command.
- Measure the `paths:` of `_authorization.md`, never guess them: run `speccheck.GovernedPath` over every file the Tasks declare in a `go test -overlay` probe that writes nothing to the repository, and record `paths: []` when none is governed.
- The release plan the baseline requires is `roundfix release plan`, and the
  maintainer decisions it defers to live in `docs/user-guide/release-runbook.md`.

## Anti-patterns (immediate rejection)

1. Introducing Cobra, testify, or any dependency the stdlib covers.
