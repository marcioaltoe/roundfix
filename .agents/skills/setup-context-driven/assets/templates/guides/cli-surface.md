# CLI surface

Design commands as automation protocols. Keep stdout for requested output,
stderr for diagnostics, and JSON schemas stable when machine output exists.

Write operations must be explicit and safe by default. Prefer dry-run, confirm,
or idempotency controls over prompts that block non-interactive agents.
