---
spec: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
status: active
created: 2026-10-16
surfaces: [backend]
---

# The F6 sender is named, and a rebuilt item binary is replaced

This is an engineering-framed minimal PRD for a bug fix that enters through the
TechSpec (`docs/agents/spec-routing.md`). F6 is the stdout subtest of
`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` in `internal/cli`.
Its `migrate --check` child has ended by a signal in four recorded full
suites, and nobody could say who sent it
([the Backlog Entry of 2026-10-13](references/2026-10-13-trace-the-sender-of-the-f6-signal.md)).
Every occurrence was the subtest whose build copies the fixture over the
`bin/roundfix` the previous subtest ran. On 2026-10-16 the author measured that
macOS kills a signed executable rewritten in place while another process still
maps it, with exactly F6's notice, and that the kill leaves a diagnostic report
naming the kernel's code-signing enforcement, readable without `sudo`. This
Spec records how the probe ended in the test itself, runs full suites to catch
a recurrence and read its report, and makes every delivery test replace the
item binary instead of writing into it, so ten full suites pass.

## Prerequisites

None. Specs 0272, 0273 and 0274 of cycle AH are authored beside this one and
touch other files; no Verification here depends on them, and this Spec is last
in the queue (release v0.86.0). Rebase on the default branch and run the strict
check again before the Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema, record or
  persisted format changes; only test fixtures and test failure messages
  change. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — nothing reaches a network, a
  provider or a credential; the tests run in temporary repositories and homes.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0281 (this Spec) decides the rule
  and the diagnosis: "A test that rebuilds an executable replaces it." ADR-0266
  keeps its rule, "A Proven Signal is the only signal Roundfix code or tests
  send", and its acceptance: "F6 counts as fixed when the class is fixed and
  ten consecutive full suites at the audited head show no child killed by a
  signal". Its rejected alternative stays rejected for the reason it gives,
  that tracing tools "need root, and the maintainer decided that F6 is fixed by
  fixing the class"; this Spec names the sender without root. ADR-0125: "so
  no test executes a file that something may still hold open"; the rename
  extends that concern to a file a process still maps. ADR-0213: "must end
  when the test binary that started it ends", so the held copy in the new test
  reads its standard input to end of file. ADR-0259: "A test runs in parallel
  unless it names why it cannot". ADR-0271: "A test bounds by real time only
  what it asserts". Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized. The Tasks change only Go test files in `internal/cli` and a
  diagnosis record in this Spec folder; no Makefile, workflow, `go.mod`, lint,
  formatter, test-runner, skill or Baseline module changes. The maintainer
  approved cycle AH on 2026-10-16 with delivery operations. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced/_authorization.md`;
  bounded files: none (`paths: []`).

## Goals

- When the item binary probe of a delivery test ends by a signal, the test's
  failure says how: the Signal Exit, whether the probe recorded its start with
  its process, parent and group IDs, and whether its build rewrote the item
  binary in place.
- A recurrence of F6 is attributed to a sender from evidence: the kernel's
  code-signing enforcement when a diagnostic report says so, another process
  when none exists.
- No delivery test rebuilds its item binary by writing into the existing file,
  and ten consecutive full suites pass.

## User Stories

1. As the maintainer, I want F6's sender named from evidence, so that I stop
   publishing releases over a known failure.
2. As an Agent reading a failed suite, I want the item binary test's failure
   to say how its probe ended, so that the next occurrence carries its own
   diagnosis.

## Core Features

1. Probe report. The delivery item fixture records its process, parent and
   group IDs with its start, and the item binary tests report, on failure, the
   inode of the item binary before and after its build, whether the probe
   recorded a start, and the step context's state.
2. Diagnosis. Up to ten full suites run at the instrumented tree before the
   fix; a recurrence is attributed through its diagnostic report, and the
   outcome is recorded in this Spec folder.
3. Replacement. Every delivery test's item binary build copies to a new file
   and renames it over the item binary path, and a test that rebuilds while a
   copy of the previous build still runs proves it.

## Non-Goals / Out of Scope

- Changing the Delivery Queue, `stepExecutable`, the Signal Exit wording or any
  production code. A user's item binary build that writes in place can meet
  the same kill; that is a follow-up, not this Spec.
- Running `sudo`, `eslogger`, `dtrace` or `dtruss`.
- Proving which process still mapped F6's item binary when its build copied
  over it.
- Retrying, waiting for or serializing the probe.

## Success Metrics

1. Success Metric: ten consecutive full suites at the audited head exit 0
   with no child killed by a signal and no Helper Residue line.
2. Success Metric: with the fixture build returned to copying in place, the
   new rebuild test fails on the audited head with the item binary probe killed
   by `SIGKILL`, and passes again with the replacement.
3. Success Metric: the diagnosis record names a sender for every F6
   recurrence it saw, or states that none occurred in its suites.

## Unreachable Acceptance

- criterion: naming a process that sent F6's signal, when a recurrence leaves
  no code-signing diagnostic report.
  reason: tracing a signal one process sends another needs `eslogger`, which
  needs `sudo` and Full Disk Access that only the maintainer holds.
  satisfied-by: the Supervisor sends the maintainer, by push, the command
  `sudo eslogger signal exit cs_invalidated > /tmp/f6-eslogger.jsonl`, run in
  a terminal with Full Disk Access while ten full suites run, and then
  `jq -c 'select(.event.signal != null and (.event.signal.target.executable.path | endswith("/bin/roundfix")))' /tmp/f6-eslogger.jsonl`;
  Spec D stops until the maintainer answers.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
read through Exa on 2026-10-16:

- Apple Developer Documentation, "Updating Mac Software"
  (https://developer.apple.com/documentation/security/updating-mac-software):
  "macOS caches information about the code’s signature in the kernel. It
  doesn’t flush that cache when you modify the file’s contents. Modifying the
  file in place yields a mismatch between the file’s contents and the
  in-kernel cache, which can cause a hard-to-reproduce code-signing crash the
  next time you run the tool." and "write the updated code to a temporary file
  and replace the existing file with that temporary one". The page renders
  only with JavaScript; its text was read through two copies that cite it, the
  `rules_pkg` pull request 941 (https://github.com/bazelbuild/rules_pkg/pull/941)
  and Stack Overflow question 76561373
  (https://stackoverflow.com/questions/76561373).
- Open Radar FB8914243 (https://openradar.appspot.com/FB8914243): "If a vnode
  contains a Mach-O executable, replacing the contents with a different Mach-O
  executable produces a file that incorrectly fails code signature
  verification in xnu and won’t run." and "It’s killed by the kernel with a
  SIGKILL."
- The `eslogger` manual page of macOS 27.0.1: "Like all Endpoint Security
  clients, eslogger must be run as super-user and requires the responsible
  process to have TCC Full Disk Access authorization."

The local measurements of 2026-10-16 are recorded in ADR-0281 and in the
TechSpec's Current behavior; they are this Spec's own and do not count as the
outside row.

## Validation

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "processo filho morto por sinal SIGKILL teste roundfix F6 binario sobrescrito code signing macOS" --all --files --min-score 0.3`.
It returned `wiki/sources/digest-roundfix-toolchain-formatter-e-subprocessos-go-no-macos-2026-09-09.md`,
which had already named Go issue 42684 (https://github.com/golang/go/issues/42684),
macOS caching a signature by inode and kills by invalid signature as "uma
hipótese diagnóstica pertinente" without proof; this Spec measured it.
`wiki/concepts/sinalizacao-segura-de-processos-por-pid.md` and
`wiki/sources/digest-roundfix-pid-reutilizado-e-cache-go-no-sandbox-2026-10-12.md`
hold the stale-PID research behind ADR-0266, the hypothesis this Spec found
does not explain the recurrence after that class was fixed. The mirrors of
ADR-0266, Spec 0269 and the adopted Backlog Entry agreed with the repository.
Exa supplied the Apple, Open Radar, `rules_pkg` and Stack Overflow pages above;
the Apple page itself returned no text to the fetch. No Inbox Entry was
written: the brief for this Spec forbids Secondbrain writes, so the Supervisor
captures the research digest.

## Glossary

None.

## Decisions

- The sender is named from the macOS diagnostic report, without `sudo`; see
  ADR-0281.
- Every delivery test rebuild replaces the item binary by rename; see
  ADR-0281.
- The diagnosis runs at most ten full suites, before the fix, and its outcome
  is a record in this Spec folder.
- The acceptance stays ADR-0266's ten clean full suites.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the characterization, the contracts
and the build order.
