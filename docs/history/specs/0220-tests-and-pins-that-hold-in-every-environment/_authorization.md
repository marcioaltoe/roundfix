---
status: approved
granted: 2026-10-03
action: prove that a detached test fixture ended by its process-group membership and start time instead of a signal's permission answer, and restore this repository's eleven trailing upstream skills to the Setup Snapshot while the skill contract tests read the lock and the snapshot instead of a hand-pinned digest
consuming: 0220-tests-and-pins-that-hold-in-every-environment
paths:
  - .agents/skills/bubbletea/references/components.md
  - .agents/skills/bubbletea/references/emoji-width-fix.md
  - .agents/skills/bubbletea/references/golden-rules.md
  - .agents/skills/bubbletea/references/troubleshooting.md
  - .agents/skills/bubbletea/SKILL.md
  - .agents/skills/domain-modeling/CONTEXT-FORMAT.md
  - .agents/skills/domain-modeling/GLOSSARY-FORMAT.md
  - .agents/skills/domain-modeling/SKILL.md
  - .agents/skills/golang-concurrency/references/channels-and-select.md
  - .agents/skills/golang-concurrency/references/pipelines.md
  - .agents/skills/golang-concurrency/references/sync-primitives.md
  - .agents/skills/golang-concurrency/SKILL.md
  - .agents/skills/golang-context/references/cancellation.md
  - .agents/skills/golang-context/references/http-services.md
  - .agents/skills/golang-context/SKILL.md
  - .agents/skills/golang-error-handling/references/error-creation.md
  - .agents/skills/golang-error-handling/references/error-handling.md
  - .agents/skills/golang-error-handling/references/error-wrapping.md
  - .agents/skills/golang-error-handling/SKILL.md
  - .agents/skills/golang-lint/references/linter-reference.md
  - .agents/skills/golang-lint/SKILL.md
  - .agents/skills/golang-testing/references/benchmarks.md
  - .agents/skills/golang-testing/references/coverage.md
  - .agents/skills/golang-testing/references/examples.md
  - .agents/skills/golang-testing/references/integration-testing.md
  - .agents/skills/golang-testing/references/mocking.md
  - .agents/skills/golang-testing/SKILL.md
  - .agents/skills/no-workarounds/SKILL.md
  - .agents/skills/systematic-debugging/condition-based-waiting-example.ts
  - .agents/skills/systematic-debugging/condition-based-waiting.md
  - .agents/skills/systematic-debugging/CREATION-LOG.md
  - .agents/skills/systematic-debugging/defense-in-depth.md
  - .agents/skills/systematic-debugging/find-polluter.sh
  - .agents/skills/systematic-debugging/root-cause-tracing.md
  - .agents/skills/systematic-debugging/SKILL.md
  - .agents/skills/systematic-debugging/test-pressure-1.md
  - .agents/skills/systematic-debugging/test-pressure-2.md
  - .agents/skills/systematic-debugging/test-pressure-3.md
  - .agents/skills/testing-boss/references/ai-writes-tests.md
  - .agents/skills/testing-boss/SKILL.md
  - .agents/skills/tui-design/references/app-patterns.md
  - .agents/skills/tui-design/references/visual-catalog.md
  - .agents/skills/tui-design/SKILL.md
  - skills/baseline_skill_contract_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0220

The maintainer authorized the unattended program on 2026-09-30 and extended it
on 2026-10-01. On 2026-10-03 the maintainer answered "Continue" to the proposed
next cycle G/H/I; this Spec is item H. "Autorizar os dois" authorized the
Baseline source and the guides, and "considere autorizado a ajustar todas as
skills se necessário" authorized skill changes. The restore of upstream skill
trees under `.agents/skills/` and the replacement of the pinned digest in the
governed `skills/baseline_skill_contract_test.go` are the Governed Path
changes that authority covers here.

## How the paths were measured

The restore ran in a disposable clone of this branch at `305211a3`
(`roundfix baseline update --repo . --yes`, `Skills restored: 11`). Its
`git status --porcelain --untracked-files=all` listed the 43 skill paths above:
38 modified, 4 added and 1 removed
(`.agents/skills/domain-modeling/CONTEXT-FORMAT.md`). `GovernedPath` was then
evaluated for every path this Spec changes through a `go test -overlay` probe
that wrote nothing to the repository. Only the 44 paths above are governed.

These paths are ordinary:

- `skills-lock.json`
- `internal/baseline/repository_skills_snapshot_test.go`
- `internal/cli/implement_detach_teardown_test.go`
- `internal/cli/detach_fixture_group_darwin_test.go`
- `internal/cli/detach_fixture_group_linux_test.go`
- `internal/cli/detach_fixture_group_other_test.go`
- this Spec's directory, ADR-0224 and the two adopted Backlog Entries

## Limits

- The restored bytes are upstream bytes at `b3c45a4`, written only by the
  supported restore. No upstream skill text is edited by hand.
- No production code change, and no change to a command's output or exit code.
- No edit to the Makefile, to lint, formatter or test-runner configuration, to
  a CI workflow or to `go.mod`.
- No retry, sleep or widened deadline in place of a fix.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and no test reaches GitHub or a provider.
- No process is signalled except one the test itself started or recorded.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
