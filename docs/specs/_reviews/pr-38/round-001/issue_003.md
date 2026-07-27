---
source: coderabbit
pr: "38"
round: 1
round_created_at: "2026-07-27T15:34:32Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-outcome-integrity
head_sha: 9ed57622bb92f138aa3e23d4d59e260ebbff0116
file: internal/cli/orphan_unix_test.go
line: 282
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UG-O4,comment:PRRC_kwDOS0qyts7aENCO
review_hash: e2bd6984b73215ef552a5ae37fee3530a302f3e2d9b1eaab82902b26046e401e
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260727T152947Z_936cd84aa803ba5d/verification/batch-001-attempt-2.log'
source_review_id: "4788632386"
source_review_submitted_at: "2026-07-27T15:23:14Z"
---


# Issue 003: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🌐 Web query:

`Go exec.Cmd StdoutPipe Wait before reads complete documentation`

💡 Result:

In Go, it is incorrect to call Cmd.Wait before all reads from a pipe created by StdoutPipe (or StderrPipe) have completed [1][2][3]. This is because Cmd.Wait is responsible for waiting for the command to exit and then closing the associated pipes [1][4]. If you call Wait while there is still data to be read, the pipe may be closed prematurely, causing your read operations to fail with an error (e.g., "file already closed") or resulting in lost output [1][5][4]. To use StdoutPipe correctly: 1. Call Start to begin the command [6][3]. 2. Read from the pipe until you reach EOF or have consumed all necessary data [1][6][4]. 3. Finally, call Wait to clean up the command's resources [1][4]. Because of this requirement, it is also incorrect to use Cmd.Run when using StdoutPipe or StderrPipe, as Run internally calls Start followed immediately by Wait [1][2][3]. For streaming output that requires concurrent processing, you can read from the pipe in a separate goroutine while the main goroutine calls Wait [6]. Top results: [1][5][4][3]

Citations:

- 1: https://pkg.go.dev/os/exec
- 2: https://go.dev/src/os/exec/exec.go?m=text
- 3: https://tip.golang.org/pkg/exec/
- 4: https://github.com/golang/go/issues/19685
- 5: https://github.com/golang/go/issues/60908
- 6: https://stackoverflow.com/questions/46723308/streaming-exec-command-stdoutpipe

---

🏁 Script executed:

```shell
sed -n '220,320p' internal/cli/orphan_unix_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 3479

---



</details>

**Delay `cmd.Wait()` until after the readiness line is read.** Starting the waiter before `scanner.Scan()` can close the pipe before the helper’s `ready` line is drained, turning an early helper failure into a misleading readiness error.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/orphan_unix_test.go` around lines 248 - 282, The
startCLIForceStopOwnerProcess helper currently starts cmd.Wait in a goroutine
before reading the helper’s readiness line. Move waiter creation until after
scanner.Scan and the "ready" validation complete, while preserving cleanup’s
ability to wait for the process and report early helper failures accurately.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:0a6651e223c269e9c65c2e17 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The helper started `cmd.Wait` before draining the readiness line from `StdoutPipe`. It now starts the waiter after validating `ready`; cleanup retains a once-guarded fallback that kills and waits when readiness fails early. Focused evidence: `rtk proxy env GOCACHE=/tmp/roundfix-run-936cd84aa803ba5d-gocache go test ./internal/cli -run '^TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion$' -count=1` passed.
