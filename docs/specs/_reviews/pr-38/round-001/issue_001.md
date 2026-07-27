---
source: coderabbit
pr: "38"
round: 1
round_created_at: "2026-07-27T15:34:32Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-outcome-integrity
head_sha: 9ed57622bb92f138aa3e23d4d59e260ebbff0116
file: internal/cli/cli_test.go
line: 4049
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UG-Oi,comment:PRRC_kwDOS0qyts7aENBm
review_hash: 10e8a5818f2835343f1fc2ac79651e469653386bd5f4e9f3cbef38eeec8ee25a
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260727T152947Z_936cd84aa803ba5d/verification/batch-001-attempt-2.log'
source_review_id: "4788632386"
source_review_submitted_at: "2026-07-27T15:23:13Z"
---


# Issue 001: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**`close(agentStarted)` panics if the fake runner is invoked more than once.**

`release` is guarded by `sync.Once`, but `agentStarted` is not. `--round all` can drive more than one Agent execution (another round or batch), and a second `onRun` call would `close` an already-closed channel, panicking the whole test binary rather than failing this test. Mirror the `releaseOnce` guard.

<details>
<summary>🛡️ Proposed fix</summary>

```diff
 	agentStarted := make(chan struct{})
+	var startedOnce sync.Once
 	releaseAgent := make(chan struct{})
 	var releaseOnce sync.Once
@@
 	runner := &fakeAgentRunner{onRun: func(agent.ExecuteRequest) error {
-		close(agentStarted)
+		startedOnce.Do(func() { close(agentStarted) })
 		<-releaseAgent
 		return nil
 	}}
```

</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	agentStarted := make(chan struct{})
	var startedOnce sync.Once
	releaseAgent := make(chan struct{})
	var releaseOnce sync.Once
	release := func() {
		releaseOnce.Do(func() {
			close(releaseAgent)
		})
	}
	t.Cleanup(release)
	runner := &fakeAgentRunner{onRun: func(agent.ExecuteRequest) error {
		startedOnce.Do(func() { close(agentStarted) })
		<-releaseAgent
		return nil
	}}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 4036 - 4049, Guard the agentStarted
channel close in the fakeAgentRunner onRun callback with its own sync.Once,
mirroring releaseOnce, so repeated Agent executions under --round all do not
panic. Keep the existing release synchronization and callback behavior
unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7b52c1f0c9d2eb12d592b0cb -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The callback could close `agentStarted` more than once when the runner executes another Batch. Added a dedicated `sync.Once` guard without changing the release synchronization. Focused evidence: `rtk proxy env GOCACHE=/tmp/roundfix-run-936cd84aa803ba5d-gocache go test ./internal/cli -run '^(TestCompletionWinnerOwnerVersusForceStopPublishesOneTerminalOutcome|TestRunForceStopOwnerFailurePreservesAgentSessions)$' -count=1` passed.
