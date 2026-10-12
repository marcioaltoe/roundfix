---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A signal goes only to a proven process, and worktrees share one build cache

In the full suite a test's child is sometimes killed by a signal nobody sent on
purpose (F6 of Spec 0254), helper processes have outlived their tests by 45
hours, and every worktree builds its own multi-gigabyte Go cache. Roundfix now
sends a signal only to a process whose identity it has just proven, gives every
waiting test helper a lifeline to its test binary, and has `make` share one Go
build cache that both the Daemon and a sandboxed Agent can write.

## Context

A sweep of the repository on 2026-10-12 found six places that signal a process
or process group without proving which process now holds the PID: a test
teardown that kills a PID read from a file through `os.FindProcess`, three
group kills keyed on a child's PID after that child may have been reaped, the
detach and QA-format group kills in production, and Force Stop's escalation to
`SIGKILL`, which follows a two-second grace without proving the owner again.
macOS allocates PIDs in sequence: six consecutive spawns on the maintainer's
machine received 91867 to 91872. The suite starts more than 83,000 Git
processes, so the PID counter wraps about once per run, and a PID held across
that wrap names a stranger. Published sources agree on the cause and the
remedy. LWN (https://lwn.net/Articles/784831/) writes that "a stale PID could
be used to send a signal to the wrong process". An answer on Unix Stack
Exchange (https://unix.stackexchange.com/questions/181001) names the one safe
case: "when the target process is a direct child of the process that will be
sending the signal, and the parent has not yet waited on it". psutil's issue
2400 (https://github.com/giampaolo/psutil/issues/2400) keeps the
`pid + create_time` pair as "the best we can do" and notes that macOS has no
process descriptors to signal through.

ADR-0213 already requires that a fixture process "must end when the test binary
that started it ends", but the helpers that wait for a signal watch nothing,
and the detach fixtures watch their owner by PID alone.

The Makefile has set `GOCACHE ?= $(CURDIR)/.gocache` since 2026-07-29, because
the Run sandbox refused writes to `~/Library/Caches/go-build`. A probe on
2026-10-12 under `codex sandbox -P :workspace`, the workspace profile, found
the working directory, `$TMPDIR` and `/tmp` writable and the Go default cache
and every other home path denied. The Codex permissions guide says the same of
the built-in profile, which "allows writes inside the active workspace roots
and system temp directories" (https://developers.openai.com/codex/permissions). `go help cache` states: "The cache is safe
for concurrent invocations of the go command."

## Decision

**A Proven Signal is the only signal Roundfix code or tests send.** A
non-zero signal goes either to the `*os.Process` of a child the sender started
and has not reaped, whose PID the kernel reserves until the parent waits, or
through the store's proven-signal operations. A proven-signal operation reads
the target's start identity immediately before sending and compares it with the
identity recorded when the PID was learned; a mismatch sends nothing and counts
as proof that the recorded process ended. A process-group signal proves the
leader's identity, or, when the leader is gone, sends only while members of the
group remain, because a group ID cannot be reused while it has members. Force
Stop proves the owner again before it escalates to `SIGKILL`. The legacy
PID-only path for a Run without a recorded identity, and the explicit
`--owner-identity-unreadable` flag, are unchanged. A repository rule test parses
every Go file and fails, naming file and line, on `syscall.Kill` or `unix.Kill`
with a signal other than zero, on `os.FindProcess`, and on a string that runs
`kill`, `pkill` or `killall` other than `kill -0`, anywhere outside the
proven-signal operations it lists.

**Every waiting test helper has a Helper Lifeline.** A helper that waits for a
signal, a release or its owner ends when the read end of a pipe whose write end
only its test binary holds reaches end of file. Where the code under test owns
the helper's standard streams, the helper instead compares its owner's PID and
start identity. The suite guard every spawning package installs also marks the
test binary's environment and, after the package's tests and cleanups, looks
for live processes that carry the mark. After a grace period it reports each one
as Helper Residue, fails the package, and kills it; the mark in the process's
environment is the identity proof for that signal.

**`make` shares one Go build cache.** When `GOCACHE` is not set, the Makefile
uses `roundfix-go-build-<uid>` under `$TMPDIR`, or under `/tmp` when `TMPDIR`
is empty, if it can create the directory and write in it, and falls back to
today's `$(CURDIR)/.gocache` otherwise. An explicit `GOCACHE` still wins, so CI
is unchanged. The Go toolchain trims entries unused for five days from that
cache, as it does from any other.

**`doctor` and `gc` see Worktree Build Caches.** Project Config key
`worktree.caches` lists repository-relative directories that hold rebuildable
build output in a worktree. The `storage` line of the Doctor Command reports
how many ended worktrees hold one and their bytes, and the GC Command removes
them. This amends ADR-0172, whose storage check never "or measures directory
bytes": the measurement covers only the declared directories of ended
worktrees, stops after two seconds and then reports a lower bound.

**Acceptance is ten clean full suites.** F6 counts as fixed when the class is
fixed and ten consecutive full suites at the audited head show no child killed
by a signal. F6 failed one of the two head suites of Spec 0254's QA and once
more during its Verification; at a rate of one in four, ten clean suites happen
by chance with a probability of 0.75^10, about 5.6 %.

## Rejected

- **Naming the sender first.** `dtruss` and `dtrace` need root, and the
  maintainer decided that F6 is fixed by fixing the class.
- **Process descriptors.** `pidfd_send_signal` exists only on Linux; macOS,
  where F6 happens, has no equivalent.
- **A shared cache under Roundfix Home or the Go default cache.** The sandbox
  denies both, so a sandboxed Agent would fail before it compiles, which is the
  2026-07-27 finding the repository-local cache fixed.
- **One cache per Run, keyed by Run.** The Go cache is safe for concurrent use,
  so keys would only multiply the bytes.
- **Deleting `.gocache` when a Run settles.** It keeps one copy per live
  worktree, and concurrent Task Worktrees are where the bytes accumulate.

## Consequences

- A test that must signal something it did not start records the identity when
  it learns the PID, and the rule test enforces it for every new test.
- Code that a sandboxed Agent runs through `make` writes its build output to
  the shared cache instead of the worktree; a worktree that still holds a
  `.gocache` is reclaimed by `roundfix gc` once its Run ended.
- Doctor's `storage` line can take up to two seconds longer in a repository
  that declares `worktree.caches`, and a repository that declares none sees no
  change.
