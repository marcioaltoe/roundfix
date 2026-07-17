# Agent instructions

Generated root instructions must stay short: mandatory invariants, links, and
decision pointers only. Put conditional rule bodies in `docs/agents/` guides so
the root file remains scannable and reusable across profiles.

Use root-cause fixes. Do not silence failing tests, suppress lint, swallow
errors, or add timing hacks to close a task. If verification fails, report the
command and diagnostic instead of weakening the contract.

Completion needs fresh evidence from the current session. Record the commands
that prove each acceptance criterion and keep follow-up work out of the current
slice.
