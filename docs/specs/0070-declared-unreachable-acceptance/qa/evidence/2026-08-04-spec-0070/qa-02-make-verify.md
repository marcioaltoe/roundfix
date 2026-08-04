# QA-02 — Full repository gate

Status: pass

Command: `rtk make verify`

The command ran unpiped from the Run Worktree root and exited 0.

- Main suite: 3,287 tests passed across 25 packages.
- Isolated corpus budget: `TestCheckCorpusBudget` passed.
- Skill contract subset: 4 tests passed.
- `roundfix skills check`: passed for all repository-owned skills, including
  `qa-gate` and `evidence-gate`.
- Build: `bin/roundfix` built from `b6ea034-dirty`; the only dirty paths were
  this in-progress QA report and evidence.
- Spec Consistency Check: Spec 0070 reported no findings. Its two skipped
  detectors are artifact-presence-aware skips for a missing Vocabulary
  Contract section and missing `references/_index.md`, not failures.

The gate's known load-sensitive corpus assertion did not reproduce on this
build; the focused budget check also passed in the same gate.
