<!-- setup-context-driven:begin id=guide.cli-surface version=0.0.1 -->

# CLI surface

- **mandatory**: Treat command names, flags, stdout and stderr placement, machine-readable fields, and exit codes as public API.

- **mandatory**: Keep stdout for requested output and stderr for diagnostics, progress, and warnings.

- **mandatory**: When the repository ships a skill or a command reference that describes its commands, a change to command behavior updates that skill or reference in the same pull request.

- **mandatory**: Make automation deterministic and non-interactive.

- **mandatory**: Make write operations explicit, replayable, safe by default, and observable; use dry-run, confirmation, or idempotency contracts where the repository requires them.

<!-- setup-context-driven:end id=guide.cli-surface -->
