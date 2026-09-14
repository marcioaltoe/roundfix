# Scope, authority, and Pull Request evidence

## Task commit scope

`git diff-tree --no-commit-id --name-status -r` established each Task commit:

- `c03a9f30`: `task_01.md` and `internal/daemon/task_engine_test.go`;
- `b8869c7a`: `task_02.md`, `internal/cli/settle.go`,
  `internal/cli/settle_test.go`, `internal/daemon/task_engine.go`, and
  `internal/daemon/task_engine_test.go`;
- `4f6dabf9`: `task_03.md`, `internal/speccheck/mechanical.go`, and
  `internal/speccheck/mechanical_test.go`.

The full non-QA range `e9746a45..4f6dabf9` contains only those nine paths.
`git diff --exit-code e9746a45..4f6dabf9 --
internal/spec/authorization.go internal/speccheck/governed.go` exited 0, so the
operation vocabulary and Governed Path set did not move. The same range contains
no `Makefile`, linter, formatter, test-runner, build-tool, package-manager,
plugin, version-pin, `.github`, `.agents`, `skills`, `go.mod`, or `go.sum`
change.

The approved authorization commit `f2d3f0ce` is an ancestor of consuming commit
`4f6dabf9`, and the authorization record is byte-unchanged between that grant
commit and audited HEAD. The only governed test path changed by Tasks 01-03 was
the exactly bounded `internal/speccheck/mechanical_test.go`.

No external or symlinked Spec Root path, authorization/operation-vocabulary
owner, Governed Path owner, or residual Daemon wall-clock behavior entered the
non-QA Task range. The feature reuses `QA-AUTH-PATHS` and the existing
authorization-operation refusal. The TechSpec's Vocabulary Contract coins no
term, so the domain close check requires no glossary edit.

## Pull Request equivalent controls

The QA assignment states that no Pull Request is open for
`feat/spec-contained-authorization`; the per-Run branch has no Pull Request of
its own. Equivalent controls for audited head
`b2c1d03ecb7cf63e905f94bc850b615218347471`:

- approval: the maintainer's exact-path authorization is commit `f2d3f0ce`, an
  ancestor of every non-QA Task, and the maintainer-authored Task 04 amendment
  is current-head commit `b2c1d03e`; no Pull Request acceptance was obtained;
- checks and status: pinned `make verify`, scoped `go vet`, and strict Spec
  checking all pass on the audited head;
- unresolved review threads: no Pull Request exists, so no Pull Request thread
  stands or is claimed resolved;
- Merge-Ready acceptance: none exists yet;
- review-artifact ancestry: no accepted Pull Request review artifact exists for
  this candidate; this QA evidence is bound directly to `b2c1d03e`.

Unblocking requires the later authorized delivery stage to integrate this QA
report, obtain the configured pre-Pull-Request review for the resulting
candidate, open the Pull Request on the target branch, and observe its checks
and review state.
