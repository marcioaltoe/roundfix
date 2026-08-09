---
task: task_05
spec: 0093-a-spec-that-validates-itself
status: completed
type: infra
complexity: medium
---

# Task 05: Wire the checker into PRD and TechSpec authoring

## Overview

`write-tasks` already ends by running the checker; this Task gives `write-prd`
and `write-techspec` the same closing step, scoped to what each stage can
decide. An author fixes the artifact while it is still open, for the price of a
sub-second command, instead of meeting the same finding at a gate.

## Requirements

1. MUST add a closing validation step to `write-prd` and to `write-techspec`,
   each running the checker scoped to its own stage.
2. MUST make an error-level finding block the stage: neither skill may report
   completion while one stands.
3. MUST state, in each skill, which classes the checker does **not** decide at
   that stage, so a clean result is not read as full coverage.
4. MUST NOT change what either skill decides about its artifact's content; this
   Task adds a verification step and nothing else.
5. MUST record in the commit message which half of the standing grant it serves
   and how reliability was preserved, per that grant's obligations.

## Subtasks

- [ ] Add the closing step to `write-prd`.
- [ ] Add the closing step to `write-techspec`.
- [ ] Name the classes each stage does not decide.

## Acceptance Criteria

- [ ] Both skills instruct the author to run the stage-scoped checker before
      reporting.
- [ ] Both state that an error-level finding blocks the report.
- [ ] Both name what the stage does not cover.
- [ ] The generated copies match their sources.

## Bounded scope

Covered by the standing grant at
`docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`.
This Task may create or modify only:

- `.agents/skills/write-prd/SKILL.md`
- `.agents/skills/write-techspec/SKILL.md`
- `skills/write-prd/SKILL.md` and `skills/write-techspec/SKILL.md`, which `make skills-sync` rewrites, following the authorized edit on ADR-0081's principle
- `docs/specs/0093-a-spec-that-validates-itself/task_05.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q 'spec check' .agents/skills/write-prd/SKILL.md` — expected: exits 0. This string does not exist in the file before this Task.
- `grep -q 'spec check' .agents/skills/write-techspec/SKILL.md` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0, proving the generated copies match their sources.
- `test -z "$(git diff --name-only -- . ':!.agents/skills/write-prd/SKILL.md' ':!.agents/skills/write-techspec/SKILL.md' ':!skills/write-prd/SKILL.md' ':!skills/write-techspec/SKILL.md' ':!docs/specs/0093-a-spec-that-validates-itself/task_05.md')"` — expected: exits 0, proving no path outside the bounded list moved.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Build Order 5.
- ADR-0081, ADR-0117.

## Result

Verification Feedback attempt 1 reported that the first declared check could
not find `spec check` in the canonical PRD skill. The diagnostic artifact was
inspected and contained no log body; the Daemon diagnostic supplied the failed
command and exit status. The root cause was that the prior Agent turn stopped
before making the authorized skill edit.

Both authoring skills now finish their existing Report step by running
`roundfix spec check <slug>` with their own `--stage` value. An error-level
finding or checker execution failure blocks reporting and requires the author
to repair the named artifact and re-run the command. The additions do not
change the questions either skill asks or the content it decides.

Each skill also bounds the meaning of a clean scoped result. The PRD stage
names the TechSpec, Task Graph, commit-dependent, non-mechanical, and product
judgment classes it does not decide. The TechSpec stage names the remaining
Task Graph, commit-dependent, non-mechanical, and user-surface classes. Both
retain omitted work in later authoring stages, the full unscoped sweep, or QA,
and treat named skipped detectors as omissions rather than clean findings.

Focused checks:

- `rtk wc -c
  /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T180804Z_645ff4014195a5b6/verification/batch-003-attempt-1.log`
  — exited `0` and reported `0` bytes; the artifact path was inspected without
  copying a log body into this Result.
- `rtk make skills-sync` — exited `0`; it regenerated the embedded skill
  bundle from the canonical sources.
- `rtk cmp -s .agents/skills/write-prd/SKILL.md
  skills/write-prd/SKILL.md` and the equivalent TechSpec comparison — both
  exited `0` after synchronization.
- `rtk rg -n "roundfix spec check <slug> --stage|error-level finding|not full
  Spec coverage|named skipped detectors" .agents/skills/write-prd/SKILL.md
  .agents/skills/write-techspec/SKILL.md skills/write-prd/SKILL.md
  skills/write-techspec/SKILL.md` — exited `0` and located the command,
  blocking rule, and coverage boundary in both canonical and generated copies.
- `rtk git diff --check -- .agents/skills/write-prd/SKILL.md
  .agents/skills/write-techspec/SKILL.md skills/write-prd/SKILL.md
  skills/write-techspec/SKILL.md
  docs/specs/0093-a-spec-that-validates-itself/task_05.md` — exited `0` after
  the final documentation edits.
- `rtk git -c core.fsmonitor=false status --short --untracked-files=all` —
  reported exactly the two canonical skills, their two generated copies, and
  this Task file; no path outside the bounded list moved.

Acceptance evidence:

- Both skills instruct the author to run the stage-scoped checker before
  reporting: the PRD Report step uses `--stage prd`, and the TechSpec Report
  step uses `--stage techspec`.
- Both state that an error-level finding blocks the report: each Report step
  prohibits reporting while an error finding or checker execution failure
  stands.
- Both name what the stage does not cover: each clean-result paragraph lists
  its later-stage, commit-dependent, non-mechanical, and judgment classes and
  states where those checks remain.
- Generated copies match their sources: both focused byte comparisons exited
  `0` after the sanctioned synchronization.

Standing-grant commit-message obligation for the Daemon-owned commit:

```text
docs: validate PRD and TechSpec authoring stages

Performance: run registered file checks in each artifact's closing authoring
step.
Reliability: block on errors and execution failures; retain omitted classes in
later stages, the full sweep, and QA.

spec: 0093-a-spec-that-validates-itself / task_05
```

Follow-up outside this Task's slice: Task 02 records that the stage registry
does not yet assign `SC-CITATION-UNSUPPORTED` to the PRD stage. This Task did
not widen its protected-tooling allowlist; the unscoped sweep retains that
detector until a separately bounded change registers it for authoring stages.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
