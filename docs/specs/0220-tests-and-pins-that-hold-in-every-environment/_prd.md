---
spec: 0220-tests-and-pins-that-hold-in-every-environment
status: active
created: 2026-10-03
surfaces: [backend]
---

# Tests and pins that hold in every environment

Two checks in `make verify` hold only in the environment they were written in.
This Spec adopts both Backlog Entries that name them:

- [references/2026-10-02-the-detached-child-test-meets-eperm-on-its-process-group.md](references/2026-10-02-the-detached-child-test-meets-eperm-on-its-process-group.md)
- [references/2026-10-02-restoring-trailing-skills-breaks-the-pinned-skill-digests.md](references/2026-10-02-restoring-trailing-skills-breaks-the-pinned-skill-digests.md)

**A process-group probe that reads permission as life.**
`TestImplementDetachChildEndsWhenItsTestBinaryDies` decides that its fixture
ended by sending signal 0 to the fixture's process group. On macOS that call
returns `EPERM`, not `ESRCH`, while the group still exists but every member is
an exited, unreaped process. XNU's `killpg1` filters zombies out of the group
before it counts signalable members, then returns `EPERM` when it counted none.
The authoring session reproduced it twice. A child that leads its own group and
exits unreaped makes the probe return `EPERM` every time. In a detached,
terminal-less run of the death test, one probe in five returned `EPERM` while
the fixture's PID still answered as alive, and a moment later the group was
gone. The Daemon's Verification met this window on 2026-10-02 (Runs
`run_20261002T094236Z_2384f2babb94acdf` and
`run_20261002T142316Z_75b99bf121221798`). The stopgap of PR #347 now skips the
test whenever it happens, so the property Spec 0213 delivered goes unproved on
every macOS Verification that meets the window.

**A pin that a supported restore breaks.** Doctor reports eleven required
external skills trailing the Setup Snapshot `b3c45a4`, and the supported
restore (`roundfix baseline update --repo . --yes`) brings them to it. The
restore then fails three tests in `skills/baseline_skill_contract_test.go`,
because they compare one hand-pinned digest of all upstream-managed skill trees
with the installed trees. Measured in a disposable clone on 2026-10-03:
`Skills restored: 11`, a second run restored 0, Doctor printed
`skills: ok (43 required: 14 Roundfix-owned, 29 external)`, and `go test
./skills` failed exactly those three tests with digest `27bc59d6…` against the
pinned `8832b7ac…`.

This is a bug fix with no product behavior change. No production code changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier is created or changed;
  the Spec edits tests and restores vendored skill content. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network call or
  HTTP surface is touched by a test; the restore's source fetch happens once in
  the Task's Agent turn, never in Verification. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0213 makes a fixture process end
  with its test binary, and Core Feature 1 proves that ending without reading
  signal permission. ADR-0213: "A process a test starts, directly or through a
  detached Run, must end when the test binary that started it ends." ADR-0221 holds a required external skill to its Setup
  Snapshot, and ADR-0224 (this Spec) holds this repository's skills to that
  snapshot and to their lock instead of a pinned digest. ADR-0191 makes the
  snapshot follow its upstream by name, and the restore takes `b3c45a4` as
  that upstream. ADR-0191: "A setup snapshot mirrors one upstream list, and the
  asset sync is its only writer." ADR-0125 and ADR-0126 bind the test fixtures and the suite
  guard, and both hold unchanged. ADR-0028 owns detach, which keeps its
  behavior. ADR-0179 bounds the Governed Paths this Spec changes in
  `_authorization.md`. ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183
  check this Spec's consistency by citation and receipt. The gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167.
  ADR-0189 does not apply, because no owned skill changes. ADR-0204 and
  ADR-0206 do not apply, because no profile composition and no skill removal
  from a setup changes. ADR-0030, ADR-0096, ADR-0097, ADR-0098, ADR-0194,
  ADR-0195 and ADR-0210 do not apply, because no agent run log, journal batch,
  gate machine stage, row carry rule or evidence snapshot changes. ADR-0219
  does not apply, because no profile draft is created. ADR-0182 and
  ADR-0184 do not apply, because no Task settlement fact and no command
  surface change. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md): the unattended program of 2026-09-30,
  extended 2026-10-01, the maintainer's "Continue" to cycle G/H/I on 2026-10-03,
  "Autorizar os dois" for the Baseline source and the guides, and "considere
  autorizado a ajustar todas as skills se necessário" for skills. No Makefile,
  lint, formatter, test-runner or CI file and no dependency changes. Bounded
  files: `.agents/skills/bubbletea/SKILL.md`,
  `.agents/skills/bubbletea/references/components.md`,
  `.agents/skills/bubbletea/references/emoji-width-fix.md`,
  `.agents/skills/bubbletea/references/golden-rules.md`,
  `.agents/skills/bubbletea/references/troubleshooting.md`,
  `.agents/skills/domain-modeling/CONTEXT-FORMAT.md`,
  `.agents/skills/domain-modeling/GLOSSARY-FORMAT.md`,
  `.agents/skills/domain-modeling/SKILL.md`,
  `.agents/skills/golang-concurrency/SKILL.md`,
  `.agents/skills/golang-concurrency/references/channels-and-select.md`,
  `.agents/skills/golang-concurrency/references/pipelines.md`,
  `.agents/skills/golang-concurrency/references/sync-primitives.md`,
  `.agents/skills/golang-context/SKILL.md`,
  `.agents/skills/golang-context/references/cancellation.md`,
  `.agents/skills/golang-context/references/http-services.md`,
  `.agents/skills/golang-error-handling/SKILL.md`,
  `.agents/skills/golang-error-handling/references/error-creation.md`,
  `.agents/skills/golang-error-handling/references/error-handling.md`,
  `.agents/skills/golang-error-handling/references/error-wrapping.md`,
  `.agents/skills/golang-lint/SKILL.md`,
  `.agents/skills/golang-lint/references/linter-reference.md`,
  `.agents/skills/golang-testing/SKILL.md`,
  `.agents/skills/golang-testing/references/benchmarks.md`,
  `.agents/skills/golang-testing/references/coverage.md`,
  `.agents/skills/golang-testing/references/examples.md`,
  `.agents/skills/golang-testing/references/integration-testing.md`,
  `.agents/skills/golang-testing/references/mocking.md`,
  `.agents/skills/no-workarounds/SKILL.md`,
  `.agents/skills/systematic-debugging/CREATION-LOG.md`,
  `.agents/skills/systematic-debugging/SKILL.md`,
  `.agents/skills/systematic-debugging/condition-based-waiting-example.ts`,
  `.agents/skills/systematic-debugging/condition-based-waiting.md`,
  `.agents/skills/systematic-debugging/defense-in-depth.md`,
  `.agents/skills/systematic-debugging/find-polluter.sh`,
  `.agents/skills/systematic-debugging/root-cause-tracing.md`,
  `.agents/skills/systematic-debugging/test-pressure-1.md`,
  `.agents/skills/systematic-debugging/test-pressure-2.md`,
  `.agents/skills/systematic-debugging/test-pressure-3.md`,
  `.agents/skills/testing-boss/SKILL.md`,
  `.agents/skills/testing-boss/references/ai-writes-tests.md`,
  `.agents/skills/tui-design/SKILL.md`,
  `.agents/skills/tui-design/references/app-patterns.md`,
  `.agents/skills/tui-design/references/visual-catalog.md`,
  `skills/baseline_skill_contract_test.go`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- The detached-child death tests never skip, and they pass in every
  environment where they previously failed or skipped, including a detached
  run with no terminal.
- An exited, unreaped fixture counts as ended, and a live fixture never does.
- This repository's eleven trailing skills are restored to the Setup Snapshot,
  Doctor reports `skills: ok` without `DR-SKILL-TRAILS-SNAPSHOT`, and
  `make verify` passes.
- No test pins the digest of an upstream-managed skill tree by hand.

## Core Features

1. **A process group proved ended by its members.** The death tests decide
   that a fixture ended by reading the process table for live members of its
   process group, excluding exited, unreaped processes, and by the recorded
   start time of the fixture's PID. Cleanup accepts `EPERM` only when that
   reading shows no live member. The `EPERM` skip is removed.
2. **Skills restored and held to their records.** The eleven trailing skills
   are restored through the supported restore. The three skill contract tests
   compare each upstream-managed tree with its own `skills-lock.json` entry,
   and a new test holds every required external skill to its Setup Snapshot
   through the comparison the Doctor Command uses (ADR-0224).

## Non-Goals / Out of Scope

- Any production change, including the Doctor Command, the restore, detach and
  `internal/store` process inspection.
- The release-plan Backlog Entry
  `docs/backlog/2026-09-30-release-plan-does-not-report-skill-and-guide-checks.md`.
  It is a feature that changes the Release Plan Command's frozen JSON contract
  and its documentation contract tests, and it shares no cause or component
  with these two test defects. It stays open for its own Spec.
- Editing any upstream skill's text. The restore writes upstream bytes only.
- Makefile, CI, lint, test-runner or `go.mod` changes, and retries of any kind.

## Success Metrics

1. A child that leads its own process group and exits unreaped makes
   `kill(-pgid, 0)` return `EPERM` on macOS, and the new group reading reports
   that group ended. The death tests pass 10 consecutive iterations each with
   no skip (before: they skip or fail whenever the probe meets the window).
2. `go test ./skills` and the new snapshot test pass on the restored tree, and
   the lock comparison fails naming the skill when one byte of an
   upstream-managed tree changes (before: the restore failed three tests).
3. Doctor's `skills` line is `ok` with no `DR-SKILL-TRAILS-SNAPSHOT` (before:
   `warn` naming eleven skills).

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- XNU's `bsd/kern/kern_sig.c`
  (<https://github.com/apple-oss-distributions/xnu/blob/main/bsd/kern/kern_sig.c>),
  where `killpg1` iterates the group with a filter that excludes
  `p->p_stat == SZOMB` and then returns
  `karg.nfound > 0 ? 0 : (posix ? EPERM : ESRCH)`.
- POSIX `kill()`
  (<https://pubs.opengroup.org/onlinepubs/9699919799/functions/kill.html>),
  which states that the null signal cannot test whether a child has ended,
  because zombies remain processes.
- The Daemon's Verification logs of Runs
  `run_20261002T094236Z_2384f2babb94acdf` and
  `run_20261002T142316Z_75b99bf121221798`, which show the `EPERM` the stopgap
  now skips.
- The upstream skills repository at `b3c45a4`, the commit the embedded Setup
  Snapshot pins, from which the restore takes the eleven trees.

## Decisions

- Prove a fixture's end by its process-group membership and its recorded start
  time, never by the answer to a signal. Treating every `EPERM` as ended was
  rejected, because a live member the sender may not signal also answers
  `EPERM`.
- Restore the skills and replace the pin in one Task, so `make verify` is green
  at each settlement. Hold this repository's skills to the snapshot and the
  lock. See ADR-0224.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
