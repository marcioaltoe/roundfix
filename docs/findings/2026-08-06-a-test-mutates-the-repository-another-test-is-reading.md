---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
---

# A test mutates the repository another test is reading

**Date:** 2026-08-06
**Found by:** running `make verify` by hand in the main checkout before
relaunching a QA gate, on branch `ma/0079-one-door-for-fleet-knowledge`.

`make verify` failed, and it failed having **deleted tracked files from the
working tree**:

```
 D docs/_inbox/samples/AGENTS.go-cli.md
 D docs/_inbox/samples/AGENTS.rust-cli.md
 D docs/_inbox/samples/AGENTS.typescript-bun.md
?? docs/adr/.gitkeep
?? docs/agents/.gitkeep
?? docs/backlog/.gitkeep
```

The reported failure was in `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`
(`skills/baseline_skill_contract_integration_test.go`), at
`inspect tracked file docs/_inbox/samples/AGENTS.g...`.

## What actually happened

That test is well-behaved on its own: `copyTrackedRepository` runs
`git ls-files -z --cached` in the repository root and copies each tracked path
into `t.TempDir()`, working only on the copy afterwards. It died on
`os.Lstat` of a file that Git says is tracked and the filesystem says is gone.

So the deletion came from **another test running concurrently**. `make verify`
runs `go test -parallel 16 ./...`, and some other case mutates the real
repository root — removing `docs/_inbox/samples/*` and creating `.gitkeep`
files under `docs/adr`, `docs/agents`, and `docs/backlog`, which is the shape
of a Baseline apply executed against the live tree instead of a copy.

Restoring the tree and rerunning `make verify` passed clean, with `git status`
empty. So the failure is a race, not a regression: the outcome depends on
whether the reader reaches a path before the writer removes it.

## Why it matters more than a flake

- **The authoritative gate can report a false failure.** The QA gate's own
  static row runs `make verify`; a lost race there produces a finding about
  work that is correct, costing a full gate round — measured today at roughly
  30 minutes each.
- **The gate can also destroy state it was only supposed to observe.** The
  deletions land in the operator's working tree. A session that trusts a
  clean `git status` after a verify run is trusting a tree another process
  edited.
- **The list of tracked files is a shared resource.** Any test that copies,
  hashes, or digests the tracked corpus races with any test that writes to it.
  This repository has several of both.

## What the fix has to establish

The invariant is one line: **a test may not write inside the repository root.**
Enforcing it needs a way to prove it, not just state it — for instance a
harness that snapshots `git status --porcelain` before and after the suite and
fails on any difference, which would have named the offending case directly
instead of surfacing as a stranger's `Lstat` error.

This finding is a sibling of
`2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md`: both are
test side effects escaping the case that created them — one leaks processes,
one leaks file mutations.
