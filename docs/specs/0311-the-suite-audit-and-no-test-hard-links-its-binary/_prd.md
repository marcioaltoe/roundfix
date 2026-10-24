---
spec: 0311-the-suite-audit-and-no-test-hard-links-its-binary
status: active
created: 2026-10-24
surfaces: [backend]
---

# The suite audit, and no test hard-links its binary

Minimal PRD for a test-suite fix, minted by `write-techspec`. Two Runs lost
their QA precondition on 2026-10-23 to `internal/agent` tests whose fake acpx
child, the test binary re-executed through `os.Args[0]`, died before writing
`args.json`. The audit this planning change commits as
`docs/references/2026-10-24-test-suite-audit.md` found the sender: macOS
`syspolicyd` killed each child because the path it resolved for the child's
executable had just disappeared, and the `internal/agent` `TestMain` creates
and removes hard links to the test binary while other children of that binary
start. After this Spec the package adapters are copies, the re-executed
fixture children make none, and a repository rule refuses any hard link a
test creates, as ADR-0324 decides. The same audit measures every package and
names the slow, redundant and hidden-input tests that Specs 0312 to 0314
and later work cut; this Spec sets the program's targets.

## Sources

1. The maintainer's request of 2026-10-23, about 23:55 BRT, quoted in
   `_authorization.md`, and the audit committed by this planning change,
   [`docs/references/2026-10-24-test-suite-audit.md`](../../references/2026-10-24-test-suite-audit.md):
   the two Run failures of 2026-10-23
   (`run_20261023T205240Z_cb3296cbba44e77f`, Spec 0304, and
   `run_20261024T020951Z_493078906e450701`, Spec 0306), the `syspolicyd`
   records that match them, the load harness that reproduced them, the
   standalone reproduction, and the per-package measurements.

No open Backlog Entry or Finding shares this root cause. The open Backlog
Entry `docs/backlog/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md`
is a cost of the CI job, not of a test, and Spec 0312 adopts it.

## Prerequisites

None in `requires`. This Spec is authored against `56dd1a76` (v0.116.0),
with Specs 0306, 0307, 0309 and 0310 active; the operator delivers it right
after 0306 and before 0307, 0309 and 0310. None of those plans changes
`internal/agent/acpx_runner_test.go`, `internal/testfixture`,
`docs/agents/specific-repository.md` or the glossary entries this Spec adds.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted; the change is test code, one guide line and one glossary
  term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary directories
  and the compiled test binary. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0324 (a test never gives a running
  executable a second name), recorded by this planning change, is the
  decision both implementation Tasks carry out. The new tests follow the
  repository's existing parallel-test and real-time rules as the glossary
  states them (Parallel Test Package, Real-Time Budget) and change neither.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized. The Makefile and the CI workflows stay unchanged; the cycle
  AP batch 4 approval of 2026-10-24 grants the delivery operations. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0311-the-suite-audit-and-no-test-hard-links-its-binary/_authorization.md`;
  bounded files: `docs/agents/specific-repository.md`.

## Goals

- No `internal/agent` test child dies of a `syspolicyd` rejection: twelve
  concurrent copies of the package suite, three counts each, pass in every
  one of eight rounds.
- No test in the repository creates a hard link, and a repository rule keeps
  it so.
- A re-executed fixture child leaves nothing under the temporary directory.
- The program's targets for the Run-side Verification and the CI gate are
  recorded with the measurement they start from.

## Core Features

1. Copied package adapters. The three `internal/agent` package adapters
   are copies of the compiled test binary made before any test starts and
   removed after the last one; the `wait` and `fixture` children of
   `TestSpawnedFixtureDeathFailsFast` make none. The fake acpx helper names
   the `RunPrompt` error when the child wrote no arguments, so a later
   failure of this shape says what happened.
2. The Hard Link rule. A test in `internal/testfixture` parses every test
   file under `internal/` and fails on a call that creates a hard link; the
   repository guide states the rule and the glossary defines the
   **Single-Name Test Binary**.

## Non-Goals / Out of Scope

- The CI job's second run of `internal/docscontract`: Spec 0312.
- `internal/cli` tests reading the live documentation trees: Spec 0313.
- The shared Go build cache: Spec 0314, a change to `roundfix gc`.
- The other slow tests the audit lists; this Spec changes none of them.
- Any production code, and any host setting of the maintainer's machines.

## Success Metrics

1. Success Metric: on the audited head, eight rounds of twelve concurrent
   runs of the compiled `internal/agent` test binary, `-test.count=3
   -test.parallel 16`, all exit 0. Measured on 2026-10-24 with the fix
   prototype: 96 of 96 passed at load averages from 143 to 171; with the
   code at `56dd1a76`, four of eight rounds had at least one failing run,
   each matched by a `syspolicyd` rejection.
2. Success Metric: no `internal/agent` package adapter shares the test
   binary's file, and a fixture child run with a fresh temporary directory
   leaves no `roundfix-agent-package-adapters-` entry in it.
3. Success Metric: the Hard Link rule finds no hard link in the repository
   and reports each seeded call with its location.
4. Success Metric: program target for the CI Verification gate of a pull
   request: at most 300 s, from 463 s on PR #654 (Verify changed 271 s,
   Verify docs 173 s); Spec 0312 and later work own the change.
5. Success Metric: program target for `make verify-changed` inside a Run:
   at most 10 minutes of the Task or QA wall clock, from about 15 minutes
   recorded by the operator on 2026-10-23; Specs 0312 to 0314 and later
   work own the change.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The macOS unified log of 2026-10-23, read by the operator's audit and
  copied into `docs/references/2026-10-24-test-suite-audit.md`: the
  `syspolicyd` lines "Unable (errno: 2) to read file at <private> for pid:
  14242" and "Terminating process due to Malware rejection: 14242" at
  23:10:10 BRT, one second before the failure in the Spec 0306 Run log.
- The Apple Developer Forums thread 663456
  (https://developer.apple.com/forums/thread/663456), where a build-system
  maintainer wrote that "in case of 'File not found' macOS kills some other
  running process sharing the inode with the already removed file" and that
  he "was able to avoid inode sharing in my build process by using
  clonefile(2) instead of hardlink".

## Unreachable Acceptance

- criterion: Success Metric 4, the CI Verification gate of a pull request
  takes at most 300 s.
  reason: only GitHub Actions runs the job, and reading its step timings
  needs the network and an open pull request, which no Verification or QA
  row may reach.
  satisfied-by: the operator's reading of the `ci-verify` run of the first
  pull request delivered after Spec 0312.
- criterion: Success Metric 5, `make verify-changed` inside a Run takes at
  most 10 minutes of the Task or QA wall clock.
  reason: only a Daemon-run Verification inside a Run measures it, and a QA
  row runs inside the Run it would measure.
  satisfied-by: the operator's reading of the Verification logs of the Runs
  delivered after Specs 0312 to 0314.

## Validation

The decisions were checked on 2026-10-24 against these sources.

- Reproduction and root cause in scratch clones of `56dd1a76`, recorded in
  the audit: a stdlib overlay that logged every `SIGKILL` sent through
  `os.Process` showed none aimed at the killed child; a marker written as
  the first statement of `TestMain` showed the child died before it; the
  unified log matched each lost child to a `syspolicyd` rejection; a
  standalone program lost 48 of 101,438 children with hard-link churn and 0
  of 115,022 without.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "roundfix test
  suite flaky slow performance" --all --files --min-score 0.3`, which
  returned the Roundfix mirror's records of Spec 0265 (tests that wait on a
  child do not race the machine) and Spec 0103 (a suite that leaks nothing)
  and the wiki digests of 2026-10-12 and 2026-10-13 on load-dependent
  tests. None records `syspolicyd` or hard links. No Inbox Entry was written,
  because the briefing forbids writing there.
- Exa: the Apple Developer Forums thread above, the Flutter reproduction
  `christopherfujino/catalina-crasher-demo` (a process whose binary was
  deleted loses its children to the same rejection), and a 2026 note on
  `syspolicyd` log volume that decodes `errno: 2` as ENOENT. Challenging:
  that note says a rejection usually reaches a process that had already
  exited; here the trace shows the child alive and killed, so it does not
  apply.

## Glossary

- adds: **Single-Name Test Binary**

## Decisions

- 2026-10-24 (about 00:05 BRT): the cycle AP batch 4 answers quoted in
  `_authorization.md` — order "Logo após a 0306 (Recomendado)", scope
  "Auditoria completa", authorization "Sim, só Makefile e workflows
  (Recomendado)".

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 4 briefing of 2026-10-24 and the measurements recorded in the audit and the TechSpec.

Taken from those sources, with no further interview:

- Fix the class at its cause, the hard links, and not by retrying a killed
  child (ADR-0324).
- Copy the binary rather than clone it: a copy needs no platform call, and
  three copies of about 10 MB cost less than the suite's first test.
- Keep the load harness a QA row rather than a test: it needs about ten
  minutes of a loaded machine, which no Task Verification may spend.
- Leave every slow or hidden-input test to the Specs that own its cut; this
  Spec only measures them.

## Open Questions

None.
