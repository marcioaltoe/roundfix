# Agent instructions

This repository is Roundfix: a local-first Go CLI and daemon that picks up work
items — PR review issues and Spec Task Graphs — resolves them through the user's
selected ACP runtime, and pushes only when nothing unresolved remains. Stdlib
`flag` dispatch and a Bubble Tea v2 TUI.

Project map: `cmd/roundfix/` is the thin CLI entry point; behavior lives in
`internal/...` — `internal/cli/` owns parsing, output, and exit behavior;
`internal/app/` holds app metadata.

## Repository rules

- Keep the project KISS: prefer the smallest behavior that satisfies the
  documented product contract.
- **NEVER** copy names, branding, package names, comments, examples, or
  generated artifacts from reference projects into this repository.
- **HARD RULE — repository contracts validate at the pull request boundary**:
  `make verify` excludes the checks whose inputs are the repository itself —
  markdown contracts and the derived-artifact regeneration gates — so an
  ordinary commit does not re-run them. `make verify-changed`, the Verification
  of Runs and of the QA gate, runs the Repository Contract Tests that their
  Contract Relevance selects (ADR-0252). `make verify-docs` still runs all Repository
  Contract Tests and `roundfix spec check`, and it **MUST** pass before any pull
  request opens. Platform-only failures stay a known limit because CI is the
  Linux gate. Nothing under `docs/history/` is ever validated as live work.
- **HARD RULE — roundfix skill sync**: before opening any PR, confirm the
  roundfix skill still matches the shipped CLI behavior; a PR that changes CLI
  behavior ships the skill update too.
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
- The release plan the baseline requires is `roundfix release plan`, and the
  maintainer decisions it defers to live in `docs/user-guide/release-runbook.md`.

## Anti-patterns (immediate rejection)

1. Introducing Cobra, testify, or any dependency the stdlib covers.
