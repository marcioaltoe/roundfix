---
status: accepted
created_at: 2026-10-13T00:00:00Z
updated_at: 2026-10-13T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A test bounds by real time only what it asserts

On 2026-10-12 the QA gate of Spec 0259 ran ten full suites on a machine shared
with other Runs, at load averages from 16 to 28, and failed three times on tests
whose own subject was sound. The last failure was two tests in
`internal/cli/detach_test.go` that inject a 500 ms liveness budget to assert
something else: a Run-creation timeout, and a Run created after the liveness
deadline. The re-executed test binary did not write its liveness byte within
500 ms, the parent killed it in the wrong phase, and the tests failed with
`signal: killed`. Production allows ten seconds for the same step. The two
earlier failures had the same shape: a real 500 ms Run Budget spent by Git work,
and a 250 ms context that also had to cover a Git fork. Specs 0171 and 0213
fixed this class one test at a time, and one of their fixes raised a Run Budget
from 100 ms to 3 s, which bounds the same Git setup today. A sweep of every test
file on 2026-10-13 found 26 more places of the class, one of them shared
by about forty tests.

Published evidence agrees that load, not code, decides these failures. Silva,
Bell and others measured 52 projects under 27 resource configurations and found
that "46.5% of the flaky tests in our dataset are RAFTs", tests whose failures
depend on the resources available (https://www.jonbell.net/preprint/tse24-raft.pdf).
Pei and others studied 62 front-end tests with async-wait flakiness and found
that tests "relying on an explicit time to wait" were flakier than tests
synchronizing on state, and that "time-based issues are typically addressed by
increasing the time to wait" while the real fix introduces "a missing
synchronization point" (https://ceur-ws.org/Vol-3245/paper7.pdf).

## Decision

**A test bounds by real time only the budget it asserts.** A **Real-Time
Budget** is a wall-clock allowance below ten seconds that a test gives work it
does not control: a child process, a goroutine or Git. Every other wait in a
test ends on the event its assertion reads, or on the test deadline through
`internal/testwait`. A budget the code under test needs but the test does not
assert takes the test deadline or the production constant, never a shorter
number.

**An asserted budget races only work that cannot finish.** When the budget is
the subject, the work it races blocks until it is cancelled or until its test
binary is gone (ADR-0213); it never sleeps for a finite time that a late timer
could outlast. When the assertion needs the budget to expire after an event,
the test waits for that event and then lets the budget expire, through a clock
or timer the code already takes or a dependency added with the production
default, so production behavior does not change.

**A repository rule holds each test file to a ceiling.** A test in
`internal/testfixture` parses every Go test file and counts the Real-Time
Budgets written as constant durations below ten seconds: the duration argument
of `context.WithTimeout`, `time.After`, `time.NewTimer` and `time.AfterFunc`;
a keyed field or an assignment whose name contains timeout, deadline, grace,
budget, liveness, run creation or maximum run duration; and an argument passed
to a function of the same package whose parameter has such a name. A constant
declared in a test file of the same package counts with its value; a
production constant does not. Each file has a ceiling and
the reason its remaining budgets are safe, a file without one has a ceiling of
zero, and a ceiling only falls, as ADR-0259 decided for Sequential Tests.

## Considered options

- Raising each budget until the suite passes. Rejected: Spec 0213's 3 s Run
  Budget is the same fix and the same test class, and the sources above report
  that raised waits recur.
- A comment beside each budget naming why it is safe, in the form of
  `// Sequential:`. Rejected for now: it edits about thirty files, several of
  them rewritten by Specs in delivery, while a ceiling table in the rule's own
  file changes no other file.
- A suite run under synthetic load as a repository test. Rejected: its result
  depends on the machine, so it cannot be a gate. The QA gate of the Spec that
  introduces this rule runs one, as acceptance evidence.
- Slowing every re-executed child at startup to expose the class. Rejected: it
  needs a hook in each package's `TestMain` and covers only child processes.

## Consequences

- A new test that gives work it does not control a short budget fails the rule
  unless its file's ceiling already covers it, so the author has to choose an
  event, the test deadline or a reason.
- The rule cannot see a budget that only an identifier from another package
  carries, such as a production constant, or a sleep; the authoring sweep
  records those and their owners.
- One detach test takes as long as the production liveness budget, about ten
  seconds, because it proves that Run creation may come after that deadline.
