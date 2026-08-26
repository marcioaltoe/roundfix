---
schema: spec-tasks/v1
spec: 0098-a-hook-that-cannot-outrank-the-gate
qa: task_08
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_01]
    - id: task_07
      file: task_07.md
      needs: [task_05, task_06]
    - id: task_09
      file: task_09.md
      needs: [task_07]
    - id: task_10
      file: task_10.md
      needs: [task_07]
    - id: task_08
      file: task_08.md
      needs: [task_07, task_09, task_10]
---

# Tasks — Hook Strictness Invariant and Recovery

| id      | title                                      | type    | complexity | needs       |
| ------- | ------------------------------------------ | ------- | ---------- | ----------- |
| task_01 | Detect and record hook refusal            | backend | high       | —           |
| task_02 | Extend settle to accept completed status  | backend | high       | task_01     |
| task_03 | Re-run verification in settle             | backend | high       | task_02     |
| task_04 | Integrate settled commits                 | backend | high       | task_03     |
| task_05 | Handle deleted file staging               | backend | medium     | task_04     |
| task_06 | Document hook strictness invariant        | docs    | medium     | task_01     |
| task_07 | Acceptance verification (3 cases)         | test    | high       | task_05, task_06 |
| task_09 | Align settle surfaces (QA F-1)            | backend | medium     | task_07     |
| task_10 | Classify unnamed hook refusal (QA F-2)    | backend | high       | task_07     |
| task_08 | QA gate                                   | qa      | medium     | task_07, task_09, task_10 |

Waves: 1 → task_01 · 2 → task_02 · task_03 · task_04 · task_05 · task_06 · 3 → task_07 · 4 → task_09 · task_10 · 5 → task_08

task_09 and task_10 are the two corrective Tasks the 2026-08-25 QA gate's
findings F-1 and F-2 require. They are the whole corrective allowance; a third
finding would recut the TechSpec or promote the work to its own Spec.

## Task: task_01 — Detect and Record Hook Refusal

**Type**: backend **Status**: pending

Implement hook refusal detection in the Daemon's commit path. When `git commit` exits non-zero due to a hook (not a legitimate git error), detect it, log with Run Event, leave staged changes intact, mark Task as `completed`.

**Acceptance**:
- Hook refusal is classified and distinguishable from git errors
- Run Event published with category `hook_refused`
- Task status set to `completed` (not `failed`)
- Staged changes remain in worktree for settle recovery

**Verification**:
```bash
grep -r "hook_refused" internal/daemon/*_test.go | wc -l | grep -qE '^[1-9]'
```

---

## Task: task_02 — Extend Settle to Accept Completed Status

**Type**: backend **Status**: pending **Needs**: task_01

Update settle command contract to recognize Tasks with status `completed` but uncommitted work.

**Acceptance**:
- Settle accepts `completed` status (currently rejects)
- Resolve surface in priority order: Task Worktree → Run Worktree → checkout
- Load Task file from selected surface

**Verification**:
```bash
make test -k TestSettleAcceptsCompleted | grep -q "ok"
```

---

## Task: task_03 — Re-run Verification in Settle

**Type**: backend **Status**: pending **Needs**: task_02

Execute Task Verification in the selected settle surface without edits.

**Acceptance**:
- Verification runs commands verbatim (no code changes)
- Failure stops, leaves surface unchanged
- No Agent session, no repair prompt

**Verification**:
```bash
make test -k TestSettleVerification | grep -q "ok"
```

---

## Task: task_04 — Integrate Settled Commits

**Type**: backend **Status**: pending **Needs**: task_03

After settle verification passes, commit and integrate onto Run Branch.

**Acceptance**:
- `git add --all` stages all changes (handles deletions)
- Standard Task commit created
- Integrated onto Run Branch (or updates checkout)
- Task Worktree removed on success

**Verification**:
```bash
git log --oneline -n 1 | grep -q "^[a-f0-9]*"
```

---

## Task: task_05 — Handle Deleted File Staging

**Type**: backend **Status**: pending **Needs**: task_04

Ensure `git add --all` correctly stages file deletions (no pathspec errors).

**Acceptance**:
- Deleted files stage without error
- Deletion commits cleanly
- No temporary workarounds

**Verification**:
```bash
rm -f /tmp/test-deletion.txt && touch /tmp/test-deletion.txt && \
git add /tmp/test-deletion.txt && rm /tmp/test-deletion.txt && \
git add --all && git diff --cached --name-status | grep -q "^D"
```

---

## Task: task_06 — Document Hook Strictness Invariant

**Type**: docs **Status**: pending **Needs**: task_01

Write invariant into `docs/agents/autonomous-work.md` (Baseline module, managed region).

**Acceptance**:
- Text: "A commit hook must not be stricter than the authoritative Verification"
- Explain consequence and recovery path
- Use managed marker boundaries

**Verification**:
```bash
grep -q "commit hook must not be stricter" docs/agents/autonomous-work.md
```

---

## Task: task_07 — Acceptance Verification

**Type**: test **Status**: pending **Needs**: task_05, task_06

Verify the three measured hook refusal cases resolve without losing work.

**Acceptance**:
- Case 1 (82-line function over 80-char limit) resolves via settle
- Case 2 (2462-line file over 500-line limit) resolves via settle
- Case 3 (`sort()` instead of `toSorted()`) resolves via settle

**Verification**:
```bash
make test -k TestHookRefusalRecovery | grep -q "ok.*3 subtests"
```

---

## Task: task_08 — QA Gate

**Type**: qa **Status**: pending **Needs**: task_07

Verify all deliverables and document acceptance evidence.

**Acceptance Rows**:

| # | Requirement | Evidence |
| - | --- | --- |
| 1 | Hook refusal detection with Run Event | Test: event published |
| 2 | Settle accepts completed-but-uncommitted | Test: status accepted |
| 3 | Settle re-runs verification | Test: verification re-runs, passes |
| 4 | Three measured cases resolve | Test: all three cases pass |
| 5 | Invariant in docs | Grep: hook strictness in autonomous-work.md |

**Verification**:
```bash
roundfix spec check 0098-a-hook-that-cannot-outrank-the-gate --strict && \
make test -k TestHookStrictness | grep -q "ok" && \
grep -q "commit hook must not be stricter" docs/agents/autonomous-work.md
```
