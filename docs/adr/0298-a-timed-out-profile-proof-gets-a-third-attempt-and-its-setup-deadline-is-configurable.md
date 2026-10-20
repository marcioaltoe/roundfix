---
status: accepted
created_at: 2026-10-20T12:00:00Z
updated_at: 2026-10-20T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A timed-out profile proof gets a third attempt, and its setup deadline is configurable

Between 2026-10-19 19:30 and 2026-10-20 00:27 BRT the Delivery Queue on one
machine parked five times because a profile proof's setup timed out twice,
each time on a different tuple, and each time the same tuple passed on the
next `roundfix deliver retry`. Profile readiness proves every configured
tuple, the preferred one and every fallback, and each proof attempt has a
fixed 30-second setup deadline. One slow disposable Agent Session was enough
to park an item, at the cost of one Supervisor cycle. A proof attempt now gets
a third try after a short backoff, and the setup deadline becomes a config key.

## Decision

**Three attempts, with a doubling backoff.** A proof attempt whose setup
deadline expires while the command's own context is still live is retried
after 5 seconds, and again after 10 seconds. The third timeout is reported as
today, classified `temporary`, with the message `profile proof setup timed
out on 3 attempts of <deadline> each`. Any other failure, and a cancellation,
ends the proof at once, and a cancellation during a backoff wait starts no
further attempt. The wait is deterministic, without jitter: one Queue Owner
proves its tuples one after another, so there is no crowd of clients whose
retries could line up.

**The setup deadline is `defaults.profile_proof_setup_timeout`.** It bounds
each attempt's disposable Agent Session setup, defaults to `30s`, the value
before this record, and accepts a Go duration greater than zero and at most
`10m` in User Config or Project Config, Project Config winning. Every
profile-readiness consumer passes it to the proof: `profiles validate`, the
operational Preflight of `implement` and `review`, Doctor, Setup and
`profiles configure`. The Fallback Chain probe made after Run creation keeps
the built-in value.

**Every configured tuple is still proven.** ADR-0050 stands: Preflight proves
the Preferred Selection and every Fallback Selection before a Run, and a
Fallback Chain activates only after Run creation.

**Cleanup of an absent disposable Session is silent.** When a setup deadline
stops `sessions ensure` before acpx records the Session, the cleanup's
`sessions close` reports that no such Session exists. That is the outcome
cleanup wants, so it returns without a warning. Before this record it logged
`disposable Agent Session cleanup found no Session to close` with a recovery
step that did not apply, after every timed-out attempt.

## Considered Options

- Proving only the Preferred Selection and leaving the fallbacks for when
  they are used. Rejected by the maintainer on 2026-10-20 ("Continua provando
  todas as tuplas"): ADR-0050 and ADR-0114 rest on every tuple a Run may
  activate having been proven before the Run exists.
- Raising the built-in deadline instead of retrying. Rejected: every machine
  would wait longer for a failure that a retry clears, and the five parks
  passed on the next attempt with the same deadline.
- Retrying until the command's own context ends. Rejected: published guidance
  limits retries per request, and an unbounded loop hides a tuple that is
  really unavailable.
- Randomized backoff. Rejected for one serial client; it would also make the
  wait untestable without a real clock (ADR-0271).
- Logging the absent Session at a lower severity. Rejected: Roundfix has no
  log levels, and the setup error already reports the timeout.

## Consequences

- A tuple whose setup always times out now fails after about two minutes with
  the default deadline instead of about seventy seconds.
- An operator on a slow machine raises `defaults.profile_proof_setup_timeout`
  in User Config instead of retrying the delivery by hand.
- The console log of a delivery no longer carries a cleanup warning for a
  Session that never opened.
