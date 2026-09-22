---
spec: 0153-a-reviewer-the-workflow-runs
status: active
created: 2026-09-21
surfaces: [backend, cli, docs]
---

# A reviewer the workflow runs

## Executive Summary

A command resolves the Pre-PR Review Policy, runs a read-only Codex reviewer
session over the current candidate through the existing Agent Runtime, and
writes a record naming the candidate it examined. `none` records a configured
omission without calling anything. Every other outcome blocks with its reason.

## Project Constraints

- Identifier strategy: applicable — the record names repository, base and head
  exactly. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the Agent Runtime owns transport and
  credentials; this Spec adds neither. Source: `docs/agents/cli.md`,
  `docs/agents/autonomous-work.md`.
- Active ADR obligations: applicable — ADR-0093, ADR-0104, ADR-0130, ADR-0155
  and ADR-0156 all hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — a new command is public CLI surface. Express
  maintainer authorization: recorded in [_authorization.md](_authorization.md);
  bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/cli/cli_test.go`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## What exists already

`Config.PrePRReview` carries `Provider` and `Source`, validated against
`codex, claude, coderabbit, none` at `internal/config/config.go:1676` and
defaulted to `codex`. `doctor` reports its health. Measured by searching the
tree, `doctor.go` is that value's only consumer: nothing runs a review.

`agent.Runner` exposes `Probe`, `Run(ctx, ExecuteRequest, sink)` and
`EndSession`. The Daemon already drives agent sessions through it, so the
reviewer reuses that path and the tests substitute the interface.

## The command

```
roundfix review [--base <ref>]
```

`--base` defaults to the repository's main branch. The candidate is the range
from that base to the current head.

- **Exit 0** — the reviewer ran and returned no findings, or the policy is
  `none` and a configured omission was recorded.
- **Exit 1** — the reviewer ran and returned findings. The record carries them.
- **Exit 2** — Preflight Validation failed, or the selected mode is blocked:
  a runtime failure, a timeout, output the command cannot read, or a provider
  this slice does not execute.

The three exits are distinct on purpose. A reviewer that found problems is not
the same event as a reviewer that could not run, and collapsing them is how a
blocked review becomes indistinguishable from a clean one.

## What the record carries

- the repository, the base commit and the head commit;
- the effective provider and the `Source` that selected it;
- the execution outcome, one of: reviewed, findings, blocked, omitted;
- the findings returned, when there are any;
- the reason, when blocked.

A record names one head. Nothing in this slice consumes the record — gating
publication on it is Spec 0126 Core Feature 6 — so a record that is wrong
blocks no one while the shape is still settling.

## Blocking

`none` is reached only by configuration. A failure never selects it, and a
provider this slice does not execute never becomes it. `claude` and
`coderabbit` are valid policy values, so the command refuses them by naming
them rather than by treating them as invalid input.

## API Contracts

1. `review` accepts `--base`, refuses unknown flags, and exits 0, 1 or 2 as
   above.
2. With the policy at `none`, no reviewer call and no readiness probe occur.
3. A blocked mode produces neither a passing record nor an omitted record.

## Coverage Map

- Goal 1 → The command.
- Goal 2 → What the record carries; API Contract 1.
- Goal 3 → Blocking; API Contract 2.
- Core Feature 1 → The command.
- Core Feature 2 → What the record carries.
- Core Feature 3 → Blocking; API Contract 2.
- Core Feature 4 → Blocking; API Contract 3.
- Core Feature 5 → Blocking.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → The command, Blocking.

## Integration Points

- **`Config.PrePRReview`.** Read, not changed.
- **`agent.Runner`.** The reviewer session runs through it, stubbed in tests.
- **Spec 0126.** Claude and CodeRabbit providers, routing findings into
  corrective work, and gating publication all stay there.

## Testing Approach

1. **Codex runs.** With the policy at `codex` and a stubbed runtime, the command
   starts a session and writes a record naming repository, base and head. Fails
   on the tree as it stands, where the command does not exist.
2. **`none` calls nothing.** The stubbed runtime records no `Run` and no
   `Probe`, the record says omitted, and the exit status is zero.
3. **Blocked stays blocked.** A runtime failure, a timeout and unreadable output
   each exit 2 with the reason named, and the record says blocked in all three.
4. **Unimplemented providers refuse.** `claude` and `coderabbit` exit 2 naming
   the provider, and write no omitted record.
5. **The skill is true.** The shipped skill and its mirror describe the command,
   its exits and its refusals, and the mirror matches the canonical file.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The record type and its writer, with unit tests (depends on: none).
2. The command, its policy resolution and its exits, with the Codex path and
   every refusal, against a stubbed runtime (depends on: 1).
3. The shipped skill and the user guide (depends on: 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **A reviewer that costs money in a gate.** Every test substitutes
  `agent.Runner`; no Verification command starts a real session. The
  authorization states this as a limit, not a preference.
- **An exit status that hides a blocked review.** Collapsing "found problems"
  into "could not run" would make a broken reviewer look like a strict one.
  Three exits and Testing Approach 3 are the control.
- **A record nothing reads.** This slice deliberately produces evidence without
  consuming it. The risk is that the shape turns out wrong for the consumer
  Spec 0126 Core Feature 6 will build; the mitigation is that changing an
  unconsumed record costs nothing.
