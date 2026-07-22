<!-- setup-context-driven:begin id=guide.autonomous-work version=3 -->

# Autonomous work

Default backend work uses `codex gpt-5.5 xhigh`. Design, UI, UX, and
frontend-dominant work uses `claude opus xhigh` when the Task Graph routes that
surface.

- **mandatory**: The Supervisor authors Specs, starts and monitors Runs, and orchestrates outcomes. Delegate implementation to the selected ACP Runtime through a Roundfix Run.

- **prohibited**: The Supervisor must not write feature code or tests.

- **mandatory**: The Daemon runs each Task's declared Verification verbatim. A Task can settle `completed` only after that Verification passes; failed diagnostics return to the same Agent Session for the bounded retry policy.

<!-- setup-context-driven:end id=guide.autonomous-work -->
