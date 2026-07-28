# Project Constraint audit

Build: `859300203565dc17bfbf01ae4e7a2512e573c17c`

- All eight Task files were read in full and had `status: completed` before QA
  wrote its report.
- Identifier strategy is not applicable because the change reuses Run, Task,
  Batch, and Verification-attempt identities.
- Authentication and HTTP are not applicable because the change remains in
  local Config, Agent Session, Verification, Run Event, CLI, and TUI
  boundaries.
- ADR-0014, ADR-0025, ADR-0038, ADR-0051, ADR-0056, and ADR-0057 were traced
  through the PRD, TechSpec, Tasks, and live checks.
- Task 08 is the only protected-tooling Task. Its Daemon commit is
  `859300203565dc17bfbf01ae4e7a2512e573c17c`.

`git diff-tree --no-commit-id --name-only -r 8593002` reported exactly:

```text
.agents/skills/implement-task/SKILL.md
.agents/skills/roundfix/SKILL.md
docs/specs/0042-verification-capacity-and-daemon-task-settlement/task_08.md
internal/baseline/assets/setups/go-cli.json
internal/baseline/assets/setups/rust-cli.json
internal/baseline/assets/setups/typescript-bun.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
internal/baseline/testdata/parity-corpus/v1/manifest.json
skills/implement-task/SKILL.md
skills/roundfix/SKILL.md
```

Those are the four expressly authorized Skill files, Task 08 itself, and the
seven expressly authorized derived digest pins. The worktree had no delta
before QA created `qa/`.

Task 03 evidence caveat: its historical
`TestVerificationGate.*(Exclusive|Fair|Cancel|Release)` selector matches no
test on the assembled build after later consolidation. QA did not credit that
stale selector; current TaskCycle integration, repeated, public-CLI, and race
checks cover the same fairness and cancellation invariants.
