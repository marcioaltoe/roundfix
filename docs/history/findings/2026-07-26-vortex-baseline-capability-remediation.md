---
status: done
created_at: 2026-07-26
updated_at: 2026-09-08
absorbed_by: 0057-baseline-capability-evidence-and-retention
---

# Context-Driven Baseline — capability divergences do not carry enough evidence to be remediated (2026-07-26)

A Baseline alignment against `standard-typescript-monorepo` in `/Users/marcio/dev/vortex` returned `action_required` with two blocking and five advisory divergences. Every divergence was individually correct, but the report did not carry enough information to act on any of them. Resolving the two blocking items required reading `internal/baseline/assets/profiles/standard-typescript-monorepo.json` and `internal/baseline/profile_alignment.go` to learn what each probe actually inspects. That source reading, not the remediation itself, was the dominant cost of the session.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-26-vortex-baseline-capability-remediation.md`.
