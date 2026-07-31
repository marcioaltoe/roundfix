---
task: task_04
spec: 0055-owner-identity-without-fork
status: completed
type: backend
complexity: medium
---

# Task 04: Fix Stop argument order and add the supervised exit

## Overview

`roundfix stop <run-id> --force` rejects its trailing flag — the same
argument-ordering defect Spec 0042 fixed for Attach. Fix it the same way, and
add the one supervised path out of an unreadable identity, so failing closed on
ignorance has an explicit, operator-driven exit.

## Requirements

1. MUST accept Stop Command flags in any position relative to the Run ID,
   reusing the parsing Spec 0042 introduced for Attach rather than a second
   mechanism.
2. MUST add an explicit `--owner-identity-unreadable` flag that permits the stop
   only when the ownership proof returned `ErrOwnerIdentityUnreadable`.
3. MUST exit `2` and signal nothing when that flag is passed while the identity
   is readable, or while the proof returned a proven mismatch.
4. MUST NOT reach the supervised path through any configuration key, environment
   variable, default, or timeout — the flag is the only entry.
5. MUST keep the proven-mismatch refusal absolute: no flag weakens it.
6. MUST leave Stop Request semantics and force-stop signaling order unchanged.

## Subtasks

- [ ] Fix the argument ordering with the Attach parsing.
- [ ] Add the flag, gated on the unreadable classification.
- [ ] Cover both argument orders, the permitted case, and both refusals.

## Acceptance Criteria

- [ ] `roundfix stop run_x --force` and `roundfix stop --force run_x` behave
      identically.
- [ ] With an unreadable identity, the flag permits the stop.
- [ ] With a readable, matching identity, the flag exits `2` and signals nothing.
- [ ] With a proven mismatch, the flag exits `2` and signals nothing.
- [ ] Without the flag, an unreadable identity still fails closed with its own
      diagnostic.

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/cli/` — expected: pass, including both refusals.
- `go run -buildvcs=false ./cmd/roundfix stop --help` — expected: the flag is
  documented.
- `make verify` — expected: exit 0.

## References

`_prd.md` → Goal 2 Story 4, Goal 4 Story 5, Features 4–5; `_techspec.md` →
Build Order 4, Supervised path, Risks (the supervised flag is a real hazard).

## Result

### Implementation

- Generalized the Attach argument hoister and reused it for Stop, with Stop's
  value-taking flags declared explicitly. Positional Run IDs now parse before
  or after flags without changing Attach parsing.
- Added `--owner-identity-unreadable` as a CLI-only Force Stop option. A
  readable match and a proven mismatch return exit `2` before Agent Session
  cleanup or owner termination. Only `ErrOwnerIdentityUnreadable` authorizes
  the existing PID-only termination proof; cleanup, termination, owner-exit
  proof, and terminal completion retain their prior order.
- Documented the new flag in Stop help and rejected it without `--force`.

### Focused checks

- Before implementation,
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/cli -run 'TestRunForceStop(AcceptsFlagsInAnyPosition|OwnerIdentityUnreadableFlagRequiresUnreadableProof|HelpExplainsProofBeforeCompletion)$'`
  reproduced the defects: 2 cases passed and 6 failed because trailing
  `--force` was positional and `--owner-identity-unreadable` was undefined.
- After implementation, the same focused selection passed 9 cases.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/cli -run '^(TestRunForceStop|TestRunStopByRunIDRecordsStopRequest|TestRunStopGracefulThenForceCompletesImmediately|TestRunStopSpecSelectorRejectsOtherSelectors|TestAttachAcceptsDocumentedFlagOrders|TestParseAttachCommandAcceptsFlagsInAnyPosition)'`
  passed 40 cases, covering the surrounding Force Stop, Stop Request, and
  shared Attach parsing behavior.
- The first focused attempt without the writable cache was blocked by denied
  access to `/Users/marcio/Library/Caches/go-build`; the focused commands above
  used the writable task-local cache. The declared `## Verification` commands
  remain unrun for the Daemon.

### Acceptance evidence

- Both argument orders: `TestRunForceStopAcceptsFlagsInAnyPosition` passed for
  `stop <run-id> --force` and `stop --force <run-id>` and observed the Stopped
  Run state in both cases.
- Unreadable identity with the flag: the `unreadable identity permits
  supervised stop` case passed, observed one termination call with the
  PID-only proof, and observed the Stopped Run state.
- Readable matching identity with the flag: the `readable matching identity
  refuses supervised flag` case passed with exit `2`, zero termination calls,
  and the Run still Active.
- Proven mismatch with the flag: the `proven mismatch refuses supervised
  flag` case passed with exit `2`, zero termination calls, and the Run still
  Active.
- Unreadable identity without the flag: the `unreadable identity without flag
  fails closed` case passed with exit `1`, zero termination calls, the
  `owner process identity is unreadable` diagnostic, and the Run still Active.

### Follow-ups

None discovered within this Task slice.
