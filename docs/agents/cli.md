<!-- setup-context-driven:begin id=guide.cli-surface version=0.0.1 -->

# CLI surface

- Treat command names, flags, stdout and stderr placement, machine-readable fields, and exit codes as public API. Keep stdout for requested output and stderr for diagnostics, progress, and warnings.

- Make automation deterministic and non-interactive. Write operations must be explicit, replayable, safe by default, and observable; use dry-run, confirmation, or idempotency contracts where the repository requires them.

<!-- setup-context-driven:end id=guide.cli-surface -->
