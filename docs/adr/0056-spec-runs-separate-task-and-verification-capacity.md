# Spec Runs separate Task Capacity and Verification Capacity

`worktree.concurrency` remains the Task Capacity: the maximum number of Task
Worktree lifecycles an Implement Run may execute concurrently. It no longer
controls how many Task Verification attempts can run at once. A separate
`verification.concurrency` setting controls per-Run Verification Capacity,
defaults to `1`, and must be positive. Each normal attempt holds one capacity
unit for its complete command sequence and releases it before Agent
Verification Feedback. This refines ADR-0025; Task Graph readiness and Task
Worktree isolation continue to define which Tasks may execute concurrently.

Every attempt journals Waiting for Verification before acquisition and starts
no command until cancellation-aware capacity acquisition succeeds. The Run
Event Stream and Live Run View report both effective capacities and each
Task's working, waiting, and verifying phase. The capacity is owned by one
Task cycle, not by a machine-wide service, so it does not coordinate other
Runs, CI, or processes outside Roundfix. ADR-0034 still governs shared
Worktree Bootstrap and database resources; projects must choose a safe Task
Capacity for mutations outside the Daemon Verification boundary.

Exit code `75` is the sole project-authored Temporary Verification Failure
signal. It is classified from the child process exit status, never from logs,
timing, or framework-specific text. A Task may receive one such retry across
its complete Verification lifecycle. The retry releases its normal capacity,
then waits for and consumes the Run's entire Verification Capacity so it runs
alone. It retains separate diagnostics and does not consume the one Agent
repair allowed by ADR-0038. A repeated temporary failure is terminal for the
Task; a deterministic retry failure may use the existing Agent repair if it
has not already been used. This refines ADR-0038 without introducing a third
numbered Verification attempt or a configurable retry budget.
