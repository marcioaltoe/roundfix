<!-- source-baseline-entry: contract.autonomous.supervisor-runtime -->
# Autonomous work protocol

1. The Supervisor authors Specs, starts and monitors Runs, and orchestrates outcomes.
2. The Supervisor MUST delegate implementation to the selected ACP Runtime through the repository's autonomous execution boundary.
3. The Supervisor MUST NOT write feature code or tests during an autonomous Run.
4. The Daemon runs each Task's declared Verification verbatim.
5. A Task can settle `completed` only after Verification passes. Failed diagnostics return to the same Agent Session under the bounded retry policy.
<!-- /source-baseline-entry: contract.autonomous.supervisor-runtime -->

<!-- source-baseline-entry: clause.autonomous.delegate-through-roundfix -->
The Supervisor authors Specs, starts and monitors Runs, and orchestrates outcomes. Delegate implementation to the selected ACP Runtime through a Roundfix Run.
<!-- /source-baseline-entry: clause.autonomous.delegate-through-roundfix -->

<!-- source-baseline-entry: clause.autonomous.supervisor-prohibition -->
The Supervisor must not write feature code or tests.
<!-- /source-baseline-entry: clause.autonomous.supervisor-prohibition -->

<!-- source-baseline-entry: clause.autonomous.daemon-verifies-task -->
The Daemon runs each Task's declared Verification verbatim. A Task can settle `completed` only after that Verification passes; failed diagnostics return to the same Agent Session for the bounded retry policy.
<!-- /source-baseline-entry: clause.autonomous.daemon-verifies-task -->
