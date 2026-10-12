---
spec: 0259-proven-signals-and-a-shared-build-cache
status: active
created: 2026-10-12
surfaces: [backend, cli, infra, docs]
---

# Proven signals and a shared build cache

Two defects make the repository's own test runs unreliable and expensive. In
the full suite a test's child process is sometimes killed by a signal nobody
meant for it: Spec 0254 shipped with that failure, F6, as a known defect, and
helper processes have outlived their tests by 45 hours
([the F6 Backlog Entry](references/2026-10-09-a-test-child-killed-by-an-unknown-signal.md)).
Meanwhile every Run Worktree, Task Worktree and Delivery Queue item worktree
builds its own Go cache, one of them 8 GB, and on 2026-10-09 the maintainer's
disk fell to 2 GB free three times
([the build cache Backlog Entry](references/2026-10-09-each-worktree-keeps-its-own-go-build-cache.md)).
This Spec makes every signal go only to a process proven to be the one
intended, gives every waiting helper a way to notice that its test binary is
gone, has `make` share one Go build cache, and lets `roundfix doctor` and
`roundfix gc` report and reclaim the build caches left in ended worktrees.

The maintainer decided on 2026-10-11, finally: fix the class, not the case; F6
counts as fixed when it does not reappear in a bounded number of full suites,
and naming the sender is not required; share or reclaim the cache, report its
bytes in the `storage` line of the Doctor Command and reclaim it in the GC
Command. The cycle grant covers the Makefile and the CI workflows for this Spec.

## Prerequisites

None. The other Specs of cycle AE (0255 to 0258 and 0260) are authored at the
same time. Spec 0256 sweeps the readers of the archive root and Spec 0260
changes help exit codes; neither touches a file this Spec's Verification
reads. No Verification here depends on another Spec's artifacts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record format
  changes; the Run Database is read, never migrated. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; tests and Verification use temporary homes, repositories and
  fake runners. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0266 (this Spec) decides Proven
  Signals, Helper Lifelines, Helper Residue, the shared cache and Worktree
  Build Caches. ADR-0213: "must end when the test binary that started it
  ends", which the Helper Lifeline makes true for helpers that wait for a
  signal. ADR-0172: "It never deletes, compacts or writes anything", and the
  storage check stays read-only; ADR-0266 amends its byte rule for Worktree Build Caches
  only, with a two-second bound. ADR-0127: "Process residue is a readiness
  fact, not a Run record"; Helper Residue is a test-suite fact and changes
  nothing about Process Residue. ADR-0044: "a provably dead owner makes the
  lock orphaned"; the legacy PID-only path for a Run without a recorded
  identity stays. ADR-0126: "The guard therefore runs per package", so the
  residue check runs inside the suite guard each spawning package installs.
  ADR-0125: "A spawned fixture is a compiled binary, never a written script";
  no Task writes an executable. ADR-0259: a converted test is proved by the
  race detector and a shuffled run, and "a ceiling only falls". ADR-0171: "A
  retention prune reports only what it reclaimed", and ADR-0255's Run Retention
  report is left as it is. ADR-0256: "A
  release waits for the skills that describe a changed surface", so the
  `doctor` and `gc` changes update the skill that describes them. ADR-0244: "A
  Spec declares the domain terms it introduces". ADR-0257: "Inside a Run the
  Daemon is the only full gate", so Task agents run focused tests. ADR-0258:
  "A Spec is authored against the QA rerun classes", which this Spec follows.
  Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AE grant of 2026-10-11 ("Concedo" for the Governed Paths each Spec declares,
  with the Makefile covered for this Spec only); bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/storage.md`, `.roundfixrc.yml`,
  `Makefile`, `docs/agents/specific-repository.md`,
  `internal/gittest/gittest.go`, `skills/roundfix/SKILL.md`. No CI workflow,
  lint, formatter or module file change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0259-proven-signals-and-a-shared-build-cache/_authorization.md`.

## Goals

- No Roundfix code or test sends a signal to a process it has not proven to be
  the one it means, and a repository rule keeps it that way.
- No test helper outlives the test binary that started it, and the suite
  reports any that does.
- F6 does not reappear in ten consecutive full suites at the audited head.
- `make` builds into one Go cache shared by every worktree, writable by the
  Daemon and by a sandboxed Agent alike.
- `roundfix doctor` reports and `roundfix gc` reclaims the build caches left
  in worktrees whose work has ended.

## User Stories

1. As the maintainer, I want a full `make test` to pass or fail on its merits,
   so that a red suite means a defect and not a stray signal.
2. As a Task agent, I want the test I write to be unable to signal a process
   it did not prove, so that my test cannot kill another package's child.
3. As the maintainer, I want a test binary killed by a timeout to leave no
   helper behind, so that I stop finding orphans days later.
4. As an Agent in the Run sandbox, I want `make verify-changed` to build
   without my setting `GOCACHE`, so that the shared cache costs me nothing.
5. As the maintainer, I want `roundfix doctor` to tell me how many bytes of
   build cache ended worktrees hold, and `roundfix gc` to remove them, so that
   the disk does not fill while a queue runs.

## Core Features

1. **Proven Signal.** A signal other than the zero probe goes either to an
   unreaped child the sender started, or through operations that compare the
   target's start identity with the one recorded when its PID was learned and
   send nothing on a mismatch. A process-group signal proves its leader the
   same way. Force Stop proves the owner again before escalating to a forced
   kill, and an identity that changed during the grace counts as proven exit.
   The production detach and QA-format group kills use the same operations.
2. **Signal rule.** A repository rule test fails, naming file and line, on any
   direct kill with a non-zero signal, any lookup of a process by bare PID, and
   any shell `kill`, `pkill` or `killall` other than the zero probe, outside
   the listed Proven Signal operations. It covers production and test code.
3. **Helper Lifeline.** Every test helper that waits ends when its test binary
   is gone: through a pipe the binary holds, or, where the code under test owns
   the helper's streams, by comparing the owner's PID and start identity.
4. **Helper Residue.** The suite guard marks its test binary and, after the
   package's tests and cleanups, reports each live process carrying the mark as
   Helper Residue, fails the package and kills it.
5. **F6 test parallel again.** The test that saw F6 returns to running in
   parallel and the package's ceiling of Sequential Tests falls by one.
6. **Shared build cache.** Without `GOCACHE` in the environment, `make` uses
   one cache under the temporary directory, named for the user, and falls back
   to the worktree's `.gocache` only when that directory cannot be written.
7. **Worktree Build Caches.** A Project Config key lists the
   repository-relative directories that hold rebuildable build output. The
   Doctor Command's `storage` line reports how many ended worktrees hold one
   and their bytes, within a two-second bound, and the GC Command's dry run
   lists them and its live run removes them.
8. **Glossary and guides.** `CONTEXT.md` defines **Proven Signal**, **Helper
   Lifeline**, **Helper Residue** and **Worktree Build Cache**, and revises the
   **GC Command**; the user guides and the Roundfix skill describe the new
   `doctor` and `gc` output.

## User Experience

A repository that declares no Worktree Build Cache sees no change in `doctor`
or `gc`. One that declares them sees one more clause on the `storage` line and
one more block in the `gc` report:

- `storage: found (Runs reclaimable: 0; Run Database free bytes: 0; Worktree Build Caches reclaimable: 2 (8589934592 bytes); next: roundfix gc)`
- In `roundfix gc --dry-run`, after the Run Retention lines, the declared
  directories, how many are reclaimable, their bytes, how many were kept
  because their work is still in progress, and the worktree paths.

A suite that leaves a helper alive fails its package with one suite guard line
per process, naming its PID and command.

## Non-Goals / Out of Scope

- Naming the process that sent F6's signal; `dtrace` needs root.
- Changing the PID-only path for legacy Runs without a recorded identity, or
  the `--owner-identity-unreadable` flag.
- Process residue of real Runs (ADR-0127) and the `residue:` Doctor line.
- The shell fake runners' `kill -0` owner probes: a shell cannot read a start
  identity without `ps`, which the sandbox denies; the Helper Residue check
  catches one that outlives its package.
- Reclaiming the shared cache itself; the Go toolchain trims it, and
  `go clean -cache` with that `GOCACHE` empties it.
- A `.gocache` in the user's own checkout; it is not a Roundfix worktree.
- Changing the CI workflows, the race detector's place, or the item binary's
  `migrate --check` notice.

## Success Metrics

1. Success Metric: ten consecutive full suites
   (`go test -count=1 -parallel 16 ./...`) at the audited head show no child
   that exited by a signal, and no suite guard Helper Residue line.
2. Success Metric: the signal rule passes on the repository and fails, naming
   file and line, on each seeded violation: a direct kill, a bare-PID lookup,
   a group kill and a shell `kill`.
3. Success Metric: a test binary killed with `SIGKILL` while its helpers wait
   leaves none of them alive after the bound, for each helper family.
4. Success Metric: with `GOCACHE` unset, `make` exports the shared cache path
   when it is writable and the worktree's `.gocache` when it is not, and an
   explicit `GOCACHE` wins in both cases.
5. Success Metric: in a temporary home with an ended worktree holding a
   declared cache, `doctor` reports it found with its bytes, `gc --dry-run`
   lists it and changes nothing, and `gc` removes it and leaves the worktree of
   an Active Run untouched.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- LWN, "Rethinking race-free process signaling"
  (https://lwn.net/Articles/784831/, found and read through Exa): "a stale PID
  could be used to send a signal to the wrong process". This is why the
  rule forbids a signal to a bare PID.
- Unix Stack Exchange, "How can I kill a process and be sure the PID hasn't
  been reused" (https://unix.stackexchange.com/questions/181001): "There is
  only one case in which you can safely use the pid to send signals: when the
  target process is a direct child of the process that will be sending the
  signal, and the parent has not yet waited on it." This is the case the rule
  allows without a proof.
- psutil issue 2400 (https://github.com/giampaolo/psutil/issues/2400): "I think
  the `pid + create_time` pair strategy used by psutil is the best we can do",
  and "macOS doesn't seem to have them though" of process descriptors. This is
  why a Proven Signal compares start identities.
- The Codex permissions guide (https://developers.openai.com/codex/permissions):
  the built-in workspace profile "allows writes inside the active workspace
  roots and system temp directories". The installed Go toolchain, read offline
  with `go help cache`: "The cache is safe for concurrent invocations of the go
  command." Together they say one cache under the temporary directory is
  writable from the sandbox and safe for concurrent Runs.
- The Secondbrain, read during authoring through `wiki/index.md` and two
  `qmd query` searches: the mirror of the 2026-07-27 finding
  `projects/roundfix/mirror/docs/history/findings/2026-07-27-sandboxed-agents-cannot-reach-the-default-go-cache.md`
  records that Agents in the sandbox could not write the default Go cache and
  that Spec 0037's QA built with a cache under `/private/tmp`, which is why the
  shared cache lives under the temporary directory and keeps a fallback; the
  mirror of the 2026-08-06 finding
  `projects/roundfix/mirror/docs/history/findings/2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md`
  shows helper leaks recur, which is why the suite guard checks for residue
  instead of trusting each cleanup. No wiki page covered PID reuse; the Exa
  sources above were captured as a pending research Inbox Entry in the
  Secondbrain's own namespace.

## Glossary

- adds: **Proven Signal**
- adds: **Helper Lifeline**
- adds: **Helper Residue**
- adds: **Worktree Build Cache**
- changes: **GC Command**
- not a term: **Signal rule** — the name of one repository rule test, not a domain concept
- not a term: **Shared build cache** — a Makefile default of this repository, not a Roundfix concept
- not a term: **Worktree Build Caches** — plural of the declared term

## Skills

- unchanged: command resolve — task_01 only routes the detached child's process-group kill through a Proven Signal; the command's flags, output and exit codes do not change.
- unchanged: command watch — task_01 only routes the detached child's process-group kill through a Proven Signal; the command's flags, output and exit codes do not change.
- unchanged: guide docs/user-guide/commands/detached-runs.md — task_01 adds one sentence saying that a killed detached child's process group is signalled only after its leader is proven; the skill references describe no teardown signal, so their text stays.

## Decisions

- A signal goes only to an unreaped own child or through a start-identity
  proof, and a repository rule enforces it; see ADR-0266.
- Every waiting helper has a lifeline, and the suite guard reports and kills
  Helper Residue; see ADR-0266.
- `make` shares one Go cache under the temporary directory, which the sandbox
  can write; see ADR-0266.
- `doctor` and `gc` measure and remove declared Worktree Build Caches of ended
  worktrees, with a two-second bound on Doctor's measurement; see ADR-0266.
- F6 is accepted as fixed after ten clean consecutive full suites; see
  ADR-0266.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the sweep, the measurements, the
contracts and the build order.
