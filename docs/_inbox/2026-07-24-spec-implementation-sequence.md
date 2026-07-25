# Recommended Spec implementation sequence

Implement Spec 0047 next, obtain its passing QA verdict, migrate the remaining
active Specs to the `Project Constraints` contract, and only then implement
Spec 0048. After the Context-Driven Baseline lane, complete the runtime Specs
in dependency order.

This note supersedes the
[2026-07-17 implementation sequence](2026-07-17-spec-implementation-sequence.md).
That earlier sequence remains as historical context.

## Current state

- [Spec 0040](../specs/_archived/0040-mandatory-context-driven-adrs/_prd.md)
  is archived as superseded by Specs 0047 and 0048. It will not be
  implemented independently.
- [Spec 0045](../specs/_archived/0045-context-driven-baseline-0-0-1-reset/_prd.md)
  and
  [Spec 0046](../specs/_archived/0046-public-context-driven-baseline-command/_prd.md)
  are archived after their accepted QA.
- Specs 0047 and 0048 have complete Task Graphs.
- Specs 0036 and 0042 have complete Task Graphs.
- Specs 0037, 0038, and 0039 require `write-tasks` before implementation.
- The archived Spec 0041 prerequisite for Spec 0036 is already satisfied.

## Recommended sequence

1. **[0047 — Context-Driven guidance composition](../specs/0047-context-driven-guidance-composition/_prd.md).**
   Implement its six Task waves, run the final QA gate, and require the newest
   QA Report to contain `verdict: pass`.
2. **Project Constraints migration gate.** Add a complete
   `Project Constraints` section to the PRD and TechSpec of active Specs 0036,
   0037, 0038, 0039, and 0042. Record the applicable repository decisions,
   their source paths, unresolved decisions, and the tooling-change
   authorization rule. This migration prepares existing work for the
   downstream enforcement introduced by Spec 0048; it is not a separate Spec.
3. **[0048 — Context-Driven project decisions and Spec constraints](../specs/0048-context-driven-project-decisions-and-spec-constraints/_prd.md).**
   Start only after the newest Spec 0047 QA Report passes. Implement its seven
   Task waves to make confirmed project decisions and `Project Constraints`
   part of authoring, execution, and QA.
4. **[0036 — Doctor skill readiness](../specs/0036-doctor-skill-readiness/_prd.md).**
   Its Spec 0041 prerequisite is already archived as completed, so its existing
   Task Graph is ready after the baseline lane.
5. **[0037 — Terminal outcome integrity](../specs/0037-terminal-outcome-integrity/_prd.md).**
   Generate and approve its Task Graph, then implement it before Specs 0038 and
   0039. It owns the guarded terminal transitions and stop-aware behavior those
   Specs consume.
6. **[0042 — Verification capacity and Daemon Task settlement](../specs/0042-verification-capacity-and-daemon-task-settlement/_prd.md).**
   Implement its existing Task Graph after terminal outcome integrity. This is
   a recommended safety sequence rather than a declared hard dependency.
7. **[0038 — Terminal Run Worktree reconciliation](../specs/0038-terminal-run-worktree-reconciliation/_prd.md).**
   Generate and approve its Task Graph, then implement it after Spec 0037,
   whose guarded reconciliation transition it requires.
8. **[0039 — Review Source evidence and Detached outcomes](../specs/0039-review-source-evidence-and-detached-outcomes/_prd.md).**
   Generate and approve its Task Graph, then implement it after Spec 0037.
   Keeping it after Spec 0038 is recommended because both change adjacent CLI,
   lifecycle, skill, documentation, and test surfaces.

## Dependency map

```text
Hard dependencies

0047 -- passing QA --> 0048
0041 -- completed --> 0036
0037 -------------> 0038
0037 -------------> 0039

Recommended delivery order

0047 --> Project Constraints migration --> 0048
                                      then
0036 --> 0037 --> 0042 --> 0038 --> 0039
```

The Project Constraints migration must happen after Spec 0047 establishes the
final guidance composition and before Spec 0048 begins enforcing the new
contract. Completed archived Specs do not need retroactive migration.

## Branch strategy

- The current `ma/context-driven-baseline-specs` branch owns the 0040, 0045,
  and 0046 archive changes and the planned 0047–0048 Context-Driven Baseline
  lane.
- After that lane, use one `ma/` branch and one pull request per remaining
  Spec.
- Squash merge and synchronize `main` between Specs.
- Keep Specs 0038 and 0039 sequential even though they become independently
  eligible after Spec 0037; this reduces conflicts across their shared
  surfaces.

## Next action

Execute the Task Graph for Spec 0047. Do not start Spec 0048 until the newest
0047 QA Report records `verdict: pass`.
