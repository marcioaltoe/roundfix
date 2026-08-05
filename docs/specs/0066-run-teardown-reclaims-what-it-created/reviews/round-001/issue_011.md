---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/store/process.go
line: 277
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9YJ,comment:PRRC_kwDOS0qyts7dnSbq
review_hash: 7aa37e40091fa68ad29652089d5118b4024aee8d032ee6fb43f150e692f36645
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 011: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Honor context cancellation in the rediscovery loop.**

The loop repeats while each enumeration discovers a new PID. It never checks `ctx.Err()`. A tree that keeps spawning children keeps the loop running after the caller cancels the context. `TerminateAndWait` returns quickly once the context is done, so each new PID only appends an unproven outcome, and the loop continues. `reconcile --apply` then cannot be interrupted, and `outcomes` grows without bound.

Add a cancellation check at the top of each iteration.





<details>
<summary>🛡️ Proposed fix</summary>

```diff
 	for {
+		if err := ctx.Err(); err != nil {
+			return outcomes, ownerProcessControlError(pid, "terminate owned process tree", err)
+		}
 		owned, err := controller.ownedProcesses(pid)
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	for {
		if err := ctx.Err(); err != nil {
			return outcomes, ownerProcessControlError(pid, "terminate owned process tree", err)
		}
		owned, err := controller.ownedProcesses(pid)
		if err != nil {
			treeErr := ownerProcessControlError(pid, "enumerate owned process tree", err)
			if _, ownerRecorded := seen[pid]; !ownerRecorded {
				seen[pid] = struct{}{}
				outcomes = append(outcomes, unprovenTerminationOutcome(pid, treeErr))
			}
			return outcomes, treeErr
		}
		owned = normalizeProcessTree(pid, owned)

		discovered := 0
		for _, ownedPID := range owned {
			if _, ok := seen[ownedPID]; ok {
				continue
			}
			seen[ownedPID] = struct{}{}
			discovered++

			identity, absent, identityErr := controller.ownedProcessIdentity(ctx, pid, ownedPID, recordedIdentity)
			if identityErr != nil {
				outcomes = append(outcomes, unprovenTerminationOutcome(ownedPID, identityErr))
				continue
			}
			if absent {
				outcomes = append(outcomes, TerminationOutcome{PID: ownedPID, Proven: true})
				continue
			}
			if err := controller.TerminateAndWait(ctx, ownedPID, identity); err != nil {
				outcomes = append(outcomes, unprovenTerminationOutcome(ownedPID, err))
				continue
			}
			outcomes = append(outcomes, TerminationOutcome{PID: ownedPID, Proven: true})
		}
		if discovered == 0 {
			return outcomes, nil
		}
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/process.go` around lines 239 - 277, Add a ctx.Err()
cancellation check at the start of the rediscovery loop before calling
controller.ownedProcesses, returning the accumulated outcomes and cancellation
error when the context is done. Keep the existing enumeration, termination, and
discovered == 0 behavior unchanged for active contexts.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c8a07f4410c8c53d6c3e65d8 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Confirmed the process-tree rediscovery loop could enumerate again after cancellation.
  - Added a context check before every rediscovery and a regression proving canceled termination returns `context.Canceled`, preserves the already-proven owner outcome, and performs zero enumerations.
  - Focused evidence: the owner-process focused suite passed (27 tests); complete affected package suites passed (1,247 tests).
  - The Daemon owns authoritative `make verify` after this Agent turn.
