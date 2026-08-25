---
spec: 0098-a-hook-that-cannot-outrank-the-gate
status: active
created: 2026-08-25
---

# TechSpec: A Hook That Cannot Outrank the Gate

## Project Constraints

- Identifier strategy: applicable — Verification, Task Worktree, Run, and Task status vocabulary are glossary terms. The closing node checks whether the work introduced or changed one. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or request. The work is commit orchestration and recovery. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0014 (Daemon verification ownership), ADR-0038 (one Verification repair), ADR-0020, ADR-0057, ADR-0056, ADR-0096, ADR-0117 apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — autonomous-work module. Express maintainer authorization: `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`. Bounded files: `internal/daemon`, `internal/cli/settle.go`, `docs/agents/autonomous-work.md`. Source: `docs/agents/agent-instructions.md`.

## Coverage Map

| PRD Item | TechSpec Section |
| --- | --- |
| User Story 1 | Hook Refusal Detection |
| User Story 2 | Settle Extension |
| User Story 3 | Recovery Path |
| Core Feature 1 | Hook Refusal Detection |
| Core Feature 2 | Settle Recovery |
| Core Feature 3 | Deleted File Handling |
| Core Feature 4 | Invariant Documentation |

## Context

The Daemon runs the authoritative Verification and then commits. When the
repository's commit hook refuses, the work is reverted and the Run ends failed.
Three Runs died this way in one Spec, every time with work that was correct and
already verified, left staged in a Task Worktree.

The repair loop covers a Verification failure, not a hook failure. The command
built for recovery (`settle`) refuses work that is completed but uncommitted,
because its contract assumes lost work is always failed.

## Root Cause Analysis

1. **Two verification gates where one cannot be satisfied**: The Daemon is the
   verification authority (ADR-0014). After verification passes, a commit hook
   becomes a second authority. If the hook is stricter than the Verification,
   verified work is lost and cannot be recovered.

2. **No invariant stating the conflict**: This relationship — that a commit hook
   may never be stricter than the authoritative Verification — is written
   nowhere, so repositories configure the conflict by accident.

3. **Settle command incompleteness**: Settle assumes lost work is always `failed`.
   When a Task is `completed` but uncommitted (due to hook refusal), settle
   refuses because the status is not `failed`.

## Solution Design

### 1. Define the Hook Strictness Invariant

Write into `docs/agents/autonomous-work.md` (Baseline module, managed region):

```markdown
## Hook Strictness Invariant

The repository's commit hook must never be stricter than the authoritative
Verification. When the Daemon commits after a passing Verification, the commit
succeeds or the repository is misconfigured. A hook refusal is a configuration
conflict, not a legitimate gate.
```

This invariant:
- Makes the conflict explicit during repository adoption
- Prevents accidental misconfiguration
- Establishes that hook refusal is exceptional, not ordinary

### 2. Handle Hook Refusal at Commit Time

In `internal/daemon/implement.go` (or the commit boundary):

When a commit hook refuses after Verification passes:

1. **Log the refusal** with Run ID, Task ID, hook output, and exit code
2. **Classify** the refusal as a configuration error, not a Task failure
3. **Options for outcome** (pick one design decision):
   - **A) Daemon commits as the authority**: The hook's refusal is overridden
     because the Daemon already passed the authoritative gate. Uses `git commit
     --no-verify` or equivalent.
   - **B) Record as repairable**: Mark Task as `hook_refused`, keep staged
     changes, and let settle or a repair Task fix it.
   - **C) Record and stop**: End the Run with clear diagnostics naming the hook
     refusal and the recovery path.

**Design choice for this Spec**: **Option C** — record the refusal and stop
with a clear recovery path. This is the safest default:
- Does not override repository configuration
- Requires explicit Supervisor action (either reconfigure hook or run settle)
- Leaves Verification as the primary authority
- Recoverable without losing work

### 3. Extend Settle Command Recovery

Update `internal/cli/settle.go` to accept a Task with status `completed` but
uncommitted:

**Current contract:**
```go
if taskStatus != "failed" {
  return fmt.Errorf("no failed settle surface; candidates: ...")
}
```

**New contract:**
```go
if taskStatus != "failed" && taskStatus != "completed" {
  return fmt.Errorf("cannot settle Task %s: status is %s (expected failed or completed-but-uncommitted)", taskID, taskStatus)
}

if taskStatus == "completed" {
  // Assume work is already verified; just commit and integrate
  return settleCompletedButUncommitted(taskID, workSurface)
}
```

This handles the state: "verified, settled in worktree, not committed" without
requiring an Agent turn or additional verification.

### 4. Staging Safety for Deleted Files

When a Task correctly removes a file and settle stages changes:

**Current behavior**: `git add -A` can fail if the pathspec matches a deleted
file in some implementations.

**New behavior**: Use `git add --all` (or equivalent) that handles deletions
correctly. Stage the deletion as a removal, not as an error.

This is the same loss as hook refusal — verified work discarded at the commit
step — reached through a different call path.

## Specification

### Hook Refusal Detection

1. After `Verification` passes, the Daemon calls commit
2. If commit exits non-zero and stderr names a hook refusal:
   - Log the refusal with full context
   - Publish a Run Event with classification `hook_refused`
   - Leave staged changes in the worktree
   - Record Task status as `completed` (not `failed`)
   - End the Run with outcome `Unresolved` and exit code `1`
   - Print recovery action: `roundfix settle --spec <slug> --task <task_id>`

### Settle Recovery Path

1. Settle accepts Task with status `completed` and no `failed` settle surface
2. Selects the surface with completed work (Task Worktree, Run Worktree, or checkout)
3. Loads the Task file from that surface
4. Re-runs Verification in place (no edit)
5. On pass:
   - Stages all changes in the surface (including deletions)
   - Creates standard Task commit
   - Integrates onto Run Branch (or checks out if Run Worktree)
   - Completes the Task
6. On failure:
   - Stops, leaves surface unchanged
   - Prints diagnostics pointing to failed command

### Baseline Module Update

Add to `docs/agents/autonomous-work.md` under `setup-context-driven:begin`:

## Acceptance Criteria

### Functional

1. **Hook refusal is detected and recorded** (not silent failure)
   - Run Event published with classification `hook_refused`
   - Diagnostics include hook command, stderr, and exit code
   - Task remains `completed` and uncommitted

2. **Settle accepts completed-but-uncommitted Tasks**
   - No Agent turn required
   - Verification re-runs in the same worktree
   - On pass, commits and integrates
   - On fail, stops with clear diagnostics

3. **Deleted files are staged correctly**
   - `git add --all` or equivalent handles deletions
   - No pathspec errors
   - Deletion commits cleanly

4. **Three measured cases resolve without losing work**
   - 82-line function over 80-char limit → hook refused, settle fixes
   - 2462-line generated file over 500-char limit → hook refused, settle fixes
   - `Array#sort()` instead of `toSorted()` → hook refused, settle fixes

### Authoring

5. **Invariant is written in rendered guidance**
   - Appears in repository's `docs/agents/autonomous-work.md`
   - Maintenance copy confirms the rule during Baseline refresh
   - No local edit overwrites it

### Evidence

6. **No external acceptance row required**
   - Hook strictness is internal to the repository
   - Verification and settle are within Roundfix's own contract
   - Evidence is the three measured cases
