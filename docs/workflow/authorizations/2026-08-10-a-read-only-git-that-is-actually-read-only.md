---
granted: 2026-08-10
action: stop-test-git-reads-writing-the-index
paths:
  - internal/gittest/gittest.go
  - skills/baseline_skill_contract_test.go
  - docs/backlog/2026-08-10-one-reader-in-cli-still-couples-verify-to-the-docs-tree.md
  - internal/baseline/derived_ownership_test.go
  - internal/baseline/derived_regeneration_repocontract_test.go
  - skills/owned_skill_edit_repocontract_test.go
  - Makefile
  - .github/workflows/ci-verify.yml
  - docs/agents/specific-repository.md
  - docs/references/coverage-record.json
consuming: direct
---

# Tooling authorization — a read-only git that is actually read-only (2026-08-10)

On 2026-08-10 the maintainer showed two consecutive `make verify` runs where
`internal/cli` (57s) and `skills` (14s) re-executed with nothing changed:

> Ainda estamos com problemas de performance no verify

## What the investigation established

Go refuses to save a test result whose input changed while the test ran.
`internal/cli` tests record the repository's `.git` directory as an input —
`findGitRoot` in `internal/config` stats it while resolving configuration from
the test process's working directory. And during a `go test ./...`, the skills
contract test's `copyTrackedRepository` runs `git ls-files --cached` against
the real repository with a plain environment; git's optional index refresh
renames `.git/index`, which updates the `.git` directory's mtime mid-run. The
cli result is then discarded as unsaveable, every verify, forever.

Measured: an input-mtime diff across a neighbours-only test run named exactly
one changed cli input — `/Users/marcio/dev/roundfix/.git`. Individually,
`go test ./internal/cli` twice gives `(cached)`; under `./...` it never does.

This is also the mechanism behind the open backlog entry about a phantom docs
reader in cli: the invalidating step there was `git checkout --` of the
touched file — an index write — not the markdown itself.

## Authorized change

`gittest.IsolatedEnv` gains `GIT_OPTIONAL_LOCKS=0`, and `gittest.ConfigArgs`
gains `-c core.fsmonitor=false`, so every test git invocation reads without
writing the index — the same discipline the production `ExecGitRunner` already
applies. `copyTrackedRepository` adopts the gittest helpers instead of a plain
environment. The backlog entry is closed with the root cause written in.

## Addendum — 2026-08-10 — the whole-repo regeneration tests move to the boundary

Measured after the `.git` fix: changing `internal/tui`, a leaf whose own tests
cost 1.2s, still took `make verify` 1m59 — `internal/baseline` (101s) and
`skills` (67s) re-ran because their regeneration tests copy the repository tree
to run `make baseline-digests` inside it, which makes every file under
`internal/` an input of theirs.

The snapshot idea recorded in the campaign does not work and is withdrawn: any
snapshot of the whole tree changes whenever the tree changes, so the tests
would keep invalidating. These three tests are repository-consistency gates —
their verdict changes only when Baseline assets or owned skills change — so
they move to the pull-request boundary behind a `repocontract` build tag,
joining the documentation contracts already there. `go test ./...` no longer
compiles them; `make verify-docs` runs them explicitly.

The maintainer chose this route over per-test work inside `internal/cli`:

> Rota A e voltamos para o trabalho nas spec.

## Bounded by purpose

Reads stop writing; nothing else changes. No test is deleted or weakened, and
temp-repository invocations keep working — the flags are read-safety, not
behavior.

## Commit choreography

This record lands as its own commit, before the fix.
