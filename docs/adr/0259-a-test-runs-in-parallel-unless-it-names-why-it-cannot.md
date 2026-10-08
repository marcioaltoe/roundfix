---
status: accepted
created_at: 2026-10-08T00:00:00Z
updated_at: 2026-10-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A test runs in parallel unless it names why it cannot

On 2026-10-08 a full `go test -count=1 -parallel 16 ./...` on the maintainer's
10-core machine took 375 s to 384 s, and `internal/cli` alone accounted for
all of it. Of its 1,409 top-level tests, 592 never called `t.Parallel()`. Go
runs those one at a time before any parallel test of the package starts, so
their durations add up: 279 s to 324 s of the package's wall time. The same
held one level down: 450 more sequential tests sat in `internal/daemon`,
`internal/spec`, `internal/speccheck`, `internal/worktree`, `internal/store`
and `internal/baseline`, whose sequential phases measured 110 s to 145 s,
84 s to 119 s and 103 s to 138 s for the first three. On GitHub's runners the suite measured 185 s to 223 s
against a 240 s budget. ADR-0089 already made code under test take its
environment and working directory as arguments, and the command environment
carries per-test overrides, so almost none of these tests had a reason to
stay sequential. They were simply never marked.

**Every top-level test in a Parallel Test Package calls `t.Parallel()` as its
first statement, or it is a Sequential Test.** A Sequential Test carries a
`// Sequential: <reason>` comment, in its doc comment or before its first
statement, naming the process-wide state it changes: an environment variable,
the working directory, a package-level variable or a signal handler. That is
the form ADR-0089 asked for, and 31 comments in the repository already use it. A repository test parses
the test files of each Parallel Test Package and fails, naming the file and
line, on a top-level test that does neither. It also holds each package to a
ceiling on its Sequential Tests, so the comment cannot become a way around the
rule; a ceiling only falls. The rule starts with the seven packages above.
Smaller packages join it when a later change converts them.

**Converting a test proves it shares nothing.** A test that becomes parallel
can race another on a variable Go does not guard, as two delivery tests did on
`app.BuildCommit` in the authoring prototype. Go panics when a parallel test
calls `t.Setenv` or `t.Chdir`, but it does not notice a package variable. Each
conversion therefore runs the package once under the race detector and once
with `-shuffle=on`. The race detector stays out of `make test`, as decided
earlier; it runs only in the Task that converts the package. Published
evidence points the same way: Candido et al. (ASE 2017) traced about 97.5 % of
the failures that parallel runs caused in 468 Java projects to race
conditions, and only 0.8 % to broken order dependence.

**Fixtures copy what a test reads.** A test that materializes pinned history
asks Git for the paths it reads, not the whole archived corpus. One of them
wrote 6,643 files to read 759.

**One `go test` runs every selected Repository Contract Test.** ADR-0253 had
`cmd/verify-select` print one invocation per build tag, and the Makefile ran
them one after the other. One invocation that names both tags and the union
of test names and packages lets Go run the two sets' packages side by side.
Two names never collide across the tags, and the selector checks that before
it merges them.

**Rejected.** Raising `-parallel` or `-p` changes nothing while sequential
tests hold a package's queue. Splitting `internal/cli` into several packages
would let `go test` overlap them, but it moves a thousand tests and leaves
every sequential chain in place. A third-party linter such as `paralleltest`
would add a dependency for one `go/ast` walk the repository already does for
written executables. Keeping the race detector in `make test` was rejected
earlier because it multiplies the suite's time on every run. Fixture templates
copied instead of recreated were measured and left out: the suite creates
about 2,000 repositories with `git init`, a small part of the at least 83,660
Git processes it starts.

## Consequences

- A new test in a Parallel Test Package is parallel by default, or it states
  why it is not, and the repository test enforces both.
- After the conversion the suite is bound by the work it does, not by
  waiting: on the same machine and load the authoring prototype measured
  227 s against 384 s for the unchanged tree. The at least 83,660 Git processes
  per run and their files are the next cost, and they belong to a later Spec
  that batches Git in production code.
- The CI suite budget stays at 240 s until CI measures the converted suite.
