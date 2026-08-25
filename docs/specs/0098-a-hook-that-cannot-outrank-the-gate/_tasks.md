---
spec: 0098-a-hook-that-cannot-outrank-the-gate
qa: task_08
---

# Task Graph: Hook Strictness Invariant and Recovery

## Dependency Graph

```
task_01 (hook detection) →
  task_02 (settle acceptance) →
    task_03 (settle verification) →
      task_04 (settle integration) →
        task_05 (deletion handling) →
  task_06 (invariant documentation) →
task_07 (acceptance verification) →
task_08 (qa gate) [terminal, depends on all]
```

## Tasks

### task_01 — Detect and Record Hook Refusal
**Type**: backend  
**Status**: pending

Implement hook refusal detection in the Daemon's commit path. After Verification
passes, detect when `git commit` exits non-zero due to a hook refusal (not a
legitimate git error).

**Work**:
- Add hook detection logic to commit boundary (internal/daemon)
- Classify refusal: parse stderr for hook markers
- Log refusal with Run ID, Task ID, hook output, exit code
- Publish Run Event: `category: hook_refused, reason: <hook output summary>`
- Leave staged changes in place (do not revert)
- Record Task as `completed` (not `failed`)
- End Run with outcome `Unresolved` and guidance to `settle`

**Verification**:
```bash
# Verify hook refusal is detected and event published
grep -r "hook_refused" internal/daemon/*_test.go | wc -l | grep -q "^[1-9]"
```

**Context**:
- Current commit path: internal/daemon/implement.go (or similar)
- Hook marker patterns: "hook refused", "pre-commit hook", "commit-msg hook"
- Must not suppress legitimate git errors

### task_02 — Extend Settle to Accept Completed-But-Uncommitted
**Type**: backend  
**Status**: pending  
**Depends**: task_01

Update settle command contract to recognize Tasks that are `completed` but have
uncommitted work.

**Work**:
- Modify settle preflight: accept `completed` status (currently rejects)
- Resolve settle surface in priority order:
  1. Task Worktree (highest priority: has work)
  2. Run Worktree (secondary)
  3. User checkout (fallback)
- Load Task file from selected surface
- Verify Task status matches selected surface

**Verification**:
```bash
# Verify settle accepts completed status
roundfix settle --spec 0098 --task task_01 2>&1 | grep -q "Settle surface:"
```

**Context**:
- File: internal/cli/settle.go
- Current error: "Task has no failed settle surface"
- New path: handle both `failed` and `completed`

### task_03 — Re-run Verification in Settle
**Type**: backend  
**Status**: pending  
**Depends**: task_02

Execute Task Verification in the selected settle surface without editing code.

**Work**:
- Reuse Verification logic from Implement Command
- Run commands in settle surface verbatim (no edits)
- Collect stdout/stderr diagnostics
- On pass: proceed to staging
- On fail: stop, print diagnostics, leave surface unchanged
- No Agent session, no repair prompt

**Verification**:
```bash
# Verify settle runs verification correctly
make test -k "TestSettleVerification" | grep -q "ok"
```

**Context**:
- Reuse: Verification Capacity and mechanics from implement.go
- Hermetic: commands run in worktree, not checkout
- Diagnostics: point to failed command and log

### task_04 — Integrate Settled Commits
**Type**: backend  
**Status**: pending  
**Depends**: task_03

After settle verification passes, commit and integrate onto Run Branch.

**Work**:
- Stage all changes using `git add --all` (handles deletions)
- Create standard Task commit (same shape as Implement commits)
- Integrate onto Run Branch if present
- If Run Worktree: just update checkout to task commit
- If Task Worktree: integrate onto Run Branch queue
- Remove Task Worktree and Task Branch on success

**Verification**:
```bash
# Verify settle commits and integrates cleanly
git log --oneline -n 1 | grep -q "^[a-f0-9]*"
```

**Context**:
- Commit message: standard format from Implement Command
- Integration: same conflict handling as Task integration
- Cleanup: remove worktree only on success

### task_05 — Handle Deleted File Staging
**Type**: backend  
**Status**: pending  
**Depends**: task_04

Ensure `git add --all` correctly stages file deletions (no pathspec errors).

**Work**:
- Verify `git add --all` invocation handles deleted files
- Test with a Task that correctly removes a file
- Confirm deletion commits without error
- No fallback to manual pathspec filtering

**Verification**:
```bash
# Create a file, delete it in a task, verify settle commits the deletion
rm -f test-deletion.txt && git add --all && git diff --cached | grep -q "deleted"
```

**Context**:
- This is the same loss as hook refusal (verified work discarded at commit)
- Reached through different call: Task deletes, settle stages
- Must not mask with temporary git shims

### task_06 — Document Hook Strictness Invariant
**Type**: docs  
**Status**: pending  
**Depends**: task_01

Write the invariant into the repository's Baseline module. This work is a copy
of text into a managed region, not new authoring.

**Work**:
- Add section to docs/agents/autonomous-work.md (managed region)
- Text: "A commit hook must not be stricter than the authoritative Verification"
- Cite: ADR-0098
- Explain: consequence when hook refusal occurs
- Provide: recovery path (roundfix settle)
- Use exact managed marker boundaries (setup-context-driven)

**Verification**:
```bash
# Verify invariant is in rendered guidance
grep -q "commit hook must not be stricter" docs/agents/autonomous-work.md
```

**Context**:
- Baseline module: docs/agents/autonomous-work.md
- Managed region: setup-context-driven markers
- Maintenance copy: persists through Baseline refresh

### task_07 — Acceptance Verification
**Type**: test  
**Status**: pending  
**Depends**: task_05, task_06

Verify the three measured hook refusal cases resolve without losing work.

**Work**:
- Test case 1: 82-line function over 80-char limit
  - Implement the function correctly (82 chars per line)
  - Hook refuses
  - Settle accepts, re-runs verification
  - Verification passes (function is correct)
  - Settle commits deletion of violation

- Test case 2: 2462-line generated file over 500-line limit
  - Generate correctly (2462 lines)
  - Hook refuses
  - Settle accepts, re-runs verification
  - Verification passes (file is correct)
  - Settle commits the file

- Test case 3: `Array#sort()` instead of `toSorted()`
  - Use `sort()` correctly (work is correct)
  - Hook refuses
  - Settle accepts, re-runs verification
  - Verification passes (behavior is correct)
  - Settle commits the change

**Verification**:
```bash
# All three cases resolve via settle
make test -k "TestHookRefusalRecovery" | grep -q "ok.*3 subtests"
```

**Context**:
- These are real cases from prior runs (Spec 0098 PRD evidence)
- Verification: same gate that runs on implement
- No manual edits between hook refusal and settle

### task_08 — QA Gate
**Type**: qa  
**Status**: pending  
**Depends**: task_01, task_02, task_03, task_04, task_05, task_06, task_07

Verify all deliverables and document any external acceptance evidence.

**Acceptance Rows**:

| # | Requirement | Evidence | Status |
| - | --- | --- | --- |
| 1 | Hook refusal detection with Run Event | Test: hook_refused event published | TBD |
| 2 | Settle accepts completed-but-uncommitted | Test: task status completed → settle accepted | TBD |
| 3 | Settle re-runs Verification | Test: verification re-runs, passes, commits | TBD |
| 4 | Three measured cases resolve | Test: all three cases: function limit, file limit, array method | TBD |
| 5 | Invariant in rendered docs | Grep: hook strictness clause in autonomous-work.md | TBD |
| 6 | No external acceptance required | - | N/A |

**Verify**:
```bash
roundfix spec check 0098 --strict && \
make test -k "TestHookStrictness" && \
grep -q "commit hook must not be stricter" docs/agents/autonomous-work.md
```
