# Recommended Spec implementation sequence

The recommended sequence for making Roundfix usable and reliable as quickly as possible is:

1. `0041-agent-selection-runtime-readiness`
2. `0036-doctor-skill-readiness`
3. `0037-terminal-outcome-integrity`
4. `0038-terminal-run-worktree-reconciliation`
5. `0039-review-source-evidence-and-detached-outcomes`
6. `0040-mandatory-context-driven-adrs`

## Primary dependencies

```text
0041 ──→ 0036

0037 ──→ 0038
  └────→ 0039

0040 is independent
```

## Rationale

- `0041` resolves the current model, adapter, Setup, and Preflight Validation blocker. It is the first step toward allowing Roundfix itself to execute the remaining Specs reliably.
- `0036` explicitly depends on `0041`. After both land, Doctor can validate Agent Selection Profiles and the Repository Skill Set.
- `0037` establishes outcome, Force Stop, Stop Request, and Agent Session integrity. It is the required foundation for `0038` and `0039`.
- `0038` precedes `0039` because it improves Spec Run worktree recovery and sanitation, making later implementations safer.
- `0039` depends on `0037` and is the largest expansion of review, retry, evidence, notification, and Detached Run behavior.
- `0040` is independent and important, but it does not directly improve Roundfix's ability to execute Runs. It belongs after the runtime and recovery infrastructure.

## Recommended branch strategy

- Use one branch and pull request per Spec.
- Squash merge and synchronize `main` between Specs.
- `0038` and `0039` can start in parallel after `0037`, but sequential work is recommended to avoid conflicts in the CLI, Roundfix Skill, documentation, and tests.
- Do not renumber the Specs again. Their dependencies are documented, and numbering does not need to match implementation order.

## Current state

- `0041` and `0036` already have Task Graphs.
- `0037`, `0038`, `0039`, and `0040` must pass through `write-tasks` before `roundfix implement`.
