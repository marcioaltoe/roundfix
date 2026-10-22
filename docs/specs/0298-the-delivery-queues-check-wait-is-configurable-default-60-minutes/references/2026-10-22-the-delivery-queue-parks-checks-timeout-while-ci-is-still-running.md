---
type: fix
status: promoted
created: 2026-10-22
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
reason: null
---

# The Delivery Queue parks checks-timeout while CI is still running

## Symptom

On 2026-10-21 the Delivery Queue item of Spec 0292 parked `checks-timeout`
twice while its Pull Request's CI was still running, and each time it needed a
manual `roundfix deliver retry` after the checks turned green. That day the
`ci-verify` Verification gate of this repository took 30 to 33 minutes with
five active Specs: run `39639595662` was cancelled at the job's 30-minute
timeout after 30.4 minutes and run `39645653242` passed after 32.4 minutes.
The queue waits at most 30 minutes, and the operator cannot change that
without a release.

## Where

- `internal/delivery/engine.go`: `defaultCheckTimeout = 30 * time.Minute`;
  `EngineDependencies.CheckTimeout` overrides it, and only tests set it
  (`newCommandDeliveryEngine` in `internal/cli/deliver_workflow.go` passes
  nothing).
- The deadline is set per wait from `engine.clock.Now().Add(engine.checkTimeout)`
  at three call sites: `Engine.createPullRequest`, while an existing Pull
  Request still reports an earlier candidate head; `Engine.checkCandidate`,
  when the item enters `checking`; and `Engine.checkCandidate` again, once,
  when it re-runs a failed check.
- `Engine.checkCandidate` parks the item with the bare blocker
  `checks-timeout`, which names neither the bound nor how to change it.
- `.github/workflows/ci-verify.yml`: the `verify` job's `timeout-minutes` was
  raised from 30 to 60 on 2026-10-21 (`c7bc381d`, #620).

## Expected

A User Config key, such as `delivery.check_timeout`, sets the wait as a
validated Go duration with a default of 60 minutes, the CI job's timeout, so
the operator adjusts it without a release. Each wait still starts its own
deadline, and the `checks-timeout` park names the key and the value it waited.

## Evidence

The Spec 0292 parks and the CI run durations above, recorded by the operator
on 2026-10-21 and read with `gh run list --workflow ci-verify.yml` on
2026-10-22.
