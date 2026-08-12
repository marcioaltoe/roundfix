# Non-Goals and regressions

Current build `1b1bfc3` retains the declared boundaries:

- Verification Capacity is config-only and scoped to one Implement Run; docs
  explicitly disclaim coordination with other Runs, CI, or external
  processes.
- The full declared Verification commands and timeouts remain unchanged.
- Exit 1 with timeout, listener, database, and port text stays deterministic;
  only child exit 75 receives temporary classification.
- One Task-scoped temporary retry and one same-Session Agent repair remain
  independently bounded. No retry count, backoff, third numbered attempt, or
  new terminal Run outcome is exposed.
- Review Batch Verification, Task Graph readiness, failed Task Worktree
  preservation, dependency blocking, Settle recovery, serialized integration,
  and stdout/stderr ownership remain covered by the current full suite.
- Capacity does not claim protection for Agent work or Worktree Bootstrap
  resources.

The full repository gate, full race suite, focused public macro matrix, 20
repeated capacity/retry/cancellation runs, built help, documentation
contracts, and Skill checks all passed on the current build.
