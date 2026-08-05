---
date: 2026-08-05
surface: docs/specs, .agents/skills
status: open
---

# An absence grep rejects the work it was written to protect

Three Task Verification commands failed in one night. None of them failed on
the work. All three were hand-written `grep` invocations with exclusion
filters, authored by the Supervisor, that rejected correct implementations.

## What was observed

**Spec 0070, task_04.** A changed-file postflight allow-listed the two
authorized skill paths:

```bash
git diff --name-only HEAD | grep -vE "^(\.agents/skills/qa-gate/|skills/qa-gate/|…task_04\.md$)" | grep -q . && exit 1 || exit 0
```

Editing the skill moves its tree digest, which moves the setup snapshots and
the catalog digest. That fallout is sanctioned by ADR-0081 and named in the
Spec's own authorization record — and the allow-list omitted it. The Task
failed on exactly the files its authorization blessed.

**Spec 0068, task_02.** The coverage detector reported that Core Feature 4 had
no Task referencing it. The Supervisor answered by adding a sentence to a
References section saying the feature was satisfied — while the Requirement
three lines above still said the opposite. The implementation followed the
`MUST`, correctly, and the QA gate found the contradiction.

**Spec 0068, task_08.** A negative grep meant to prove the reclaim stays a
string, never a call:

```bash
if grep -rn "os.RemoveAll\|worktree remove" internal/specaudit --include="*.go" \
  | grep -v "_test.go" | grep -v "fmt.Sprintf\|const \|reclaim" | grep -q .; then exit 1; fi
```

The implementer named a constant `worktreeRemoveReclaimPrefix` holding exactly
that string — the right shape. The filter excluded lowercase `reclaim`; the
constant carries `Reclaim`. The check rejected the correct answer.

## What this says

An absence assertion built from substring exclusions is **guessing at a
property instead of testing it**, and it fails in both directions: it passes
work that names things differently, and it fails work that names them well. The
exclusion list encodes the author's imagination of what correct code looks
like, which is precisely what a Task must not do — the Task states the
contract; the implementer chooses the shape.

The repository already knows the general form of this. `docs/agents/autonomous-work.md`
warns that a bare negative grep exits non-zero when the pattern is absent and
prescribes the portable guard form. That fixes the *exit status*; it does not
fix the deeper problem, which is that the grep was standing in for a behaviour.

Each of the three had a behavioural assertion available:

| Grep intent | Behavioural proof |
| --- | --- |
| only bounded paths changed | run the sanctioned regeneration chain, then assert the gate is green |
| this feature has a Task | make a Requirement carry it |
| the audit executes nothing | a fixture whose Git state is byte-identical before and after |

The third was fixed that way and passed on the next run. Naming cannot defeat
a state comparison.

## Why it matters for autonomous work

The cost is asymmetric and invisible. A grep that wrongly *passes* costs a QA
cycle later; a grep that wrongly *fails* costs an Agent turn, a repair turn,
and a relaunch — and it teaches nothing, because the diagnostic points at
correct code. In one night this class consumed three relaunches while
twenty-eight implementation Tasks settled on attempt 1.

The rule this suggests: **a Verification command asserts a behaviour the code
must exhibit, not a shape the source must have.** Where only a source-shape
check is possible, it belongs in review, not in a gate that stops a Task.

## Evidence

- Spec 0070 `task_04`, run `run_20260804T174057Z_50ce6689051f1eab`, settled
  `failed` after one repair turn.
- Spec 0068 `task_02` Requirement 3 versus its References note; QA finding
  F-001 in `docs/specs/0068-spec-close-audit/qa/qa-report-2026-08-04.md`.
- Spec 0068 `task_08`, run `run_20260805T011411Z_b40689103917dbad`, settled
  `failed`; the offending constant is `worktreeRemoveReclaimPrefix` in
  `internal/specaudit/audit.go`.
- The replacement behavioural assertion passed on
  `run_20260805T013130Z_31398370e8ba8670`, outcome `Clean`.
