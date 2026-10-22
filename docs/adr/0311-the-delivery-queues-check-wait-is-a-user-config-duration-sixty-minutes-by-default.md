---
status: accepted
created_at: 2026-10-22T12:00:00Z
updated_at: 2026-10-22T12:00:00Z
deprecated_at: null
superseded_by: null
---

# The Delivery Queue's check wait is a User Config duration, sixty minutes by default

The Delivery Queue waited at most 30 minutes, a constant, for a Pull Request's
current-head checks. On 2026-10-21 this repository's own CI Verification gate
needed 30 to 33 minutes, so the item of Spec 0292 parked `checks-timeout`
while CI was still running and needed a manual `deliver retry` after its
checks turned green; that day the gate's job timeout went from 30 to 60
minutes. The maintainer chose "Config com padrão 60m (Recomendado)": a User
Config key, 60 minutes by default, that the operator changes without a
release.

## Decision

**One key bounds every Check Wait.** `delivery.check_timeout` is a Go
duration. It bounds the wait of an item in `checking` for its current-head
checks and a mergeable state, and the wait of publication for an existing Pull
Request to report the candidate head; both used the same constant before.

**Sixty minutes by default.** The built-in value is `60m`, the job timeout of
this repository's CI Verification gate. The change from 30 minutes is a
declared behavior change for every adopter that sets nothing.

**Greater than zero, at most a day.** A value of `0`, a negative value or a
value above `24h` fails configuration validation with
`delivery.check_timeout must be greater than 0 and at most 24h`; an
unparsable value fails as an invalid duration. A zero wait would park every
item at its first poll, because checks are always pending right after a push,
and an unbounded wait has no place in a queue that runs unattended. GitHub
documents a 6-hour job execution limit on its hosted runners and a 24-hour
queue limit for a self-hosted job, so a check still unsettled after a day
needs a person.

**User Config only.** The wait is the operator's tolerance on the machine that
runs the queue, like `runs.slot_wait` (ADR-0304). A `delivery.check_timeout`
in Project Config, the checkout's or the fetched default branch's, is ignored
with the warning `config: delivery.check_timeout in Project Config is
ignored; set delivery.check_timeout in User Config`.

**Each wait starts its own deadline.** A Check Wait begins when the step
begins and, as before, restarts once when the queue re-runs a failed check
(ADR-0236). A Queue Owner or `deliver retry` reads the value when it starts;
an owner already running keeps the value it read.

**The park names what it waited.** At the deadline the item parks with the
blocker `checks-timeout: current-head checks not settled within
delivery.check_timeout <duration>`, the duration in Go's form such as
`1h0m0s`. Its Park Class stays `environment`.

## Considered Options

- Keeping a constant and raising it to 60 minutes was rejected: the next
  slower CI needs another release.
- A Project Config key, like the other `delivery.*` keys, was rejected: those
  describe the repository's build, while this one is the operator's patience,
  and a Project Config change takes effect only after a merge that the slow
  checks themselves hold up.
- Letting `0` mean "no limit" or "park at once" was rejected: the first never
  ends a stuck wait and the second parks every item.
- Deriving the wait from the workflow's `timeout-minutes` was rejected: the
  queue would have to parse every workflow that reports a check, and a job's
  limit does not include its queue time.

## Consequences

- An adopter that sets nothing now waits up to 60 minutes per Check Wait
  before `checks-timeout`.
- An older Roundfix binary refuses a User Config that sets
  `delivery.check_timeout` as an unknown key; the default needs no setting.
- A `checks-timeout` blocker recorded before this change keeps its bare text;
  status and retry read both forms.
