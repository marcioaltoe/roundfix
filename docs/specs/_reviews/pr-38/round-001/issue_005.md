---
source: coderabbit
pr: "38"
round: 1
round_created_at: "2026-07-27T15:34:32Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-outcome-integrity
head_sha: 9ed57622bb92f138aa3e23d4d59e260ebbff0116
file: internal/store/process_windows.go
line: 28
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UG-PC,comment:PRRC_kwDOS0qyts7aENCb
review_hash: f09ca3d0edf0a06a0ec7cc0d3a4cf86bf036c9213f36c9634f4a81952e89ceea
duplicate_of: ""
source_review_id: "4788632386"
source_review_submitted_at: "2026-07-27T15:23:14Z"
---


# Issue 005: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🌐 Web query:

`Windows OpenProcess succeeds for exited process until all handles closed GetExitCodeProcess STILL_ACTIVE`

💡 Result:

In Windows, a process object is a kernel object that remains in the system as long as there is at least one open handle to it [1][2][3]. When a process terminates, the system updates its internal state (including its exit code) and signals the process object, but the object itself is not destroyed until all handles to that process are closed [1][2]. This behavior explains why OpenProcess can succeed even after a process has exited: as long as any process (including the system or other applications) holds an open handle, the process object persists, and its Process ID (PID) is protected from reuse [1][3]. If you call OpenProcess on such a "zombie" process, the kernel will grant you a handle, allowing you to perform queries like GetExitCodeProcess [1]. Regarding GetExitCodeProcess, it will return the process's actual exit code (e.g., 0 for success, or a specific error code) once the process has terminated [4][2]. It will only return STILL_ACTIVE (259) if the process is still running [4][2]. If you observe GetExitCodeProcess returning STILL_ACTIVE for a process you believe has exited, it indicates that the process has not actually finished execution from the perspective of the kernel [4][2]. To correctly determine if a process has finished: 1. Do not rely solely on OpenProcess success, as it confirms only that the process object still exists, not that it is currently executing [5][6]. 2. Use WaitForSingleObject on the process handle with a timeout of 0 to check if the process has signaled (terminated) [7][6]. 3. Call GetExitCodeProcess and verify that the returned status is not STILL_ACTIVE [4][7]. Because PIDs can be recycled after a process object is finally destroyed (i.e., all handles are closed), you should always hold onto a handle to the process you are monitoring to ensure you are referencing the same instance [1][8].

Citations:

- 1: https://devblogs.microsoft.com/oldnewthing/20110107-00/?p=11803
- 2: https://learn.microsoft.com/en-us/windows/win32/procthread/terminating-a-process
- 3: https://stackoverflow.com/questions/75083746/is-it-safe-to-use-a-process-handle-which-belongs-to-a-terminated-process
- 4: https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getexitcodeprocess
- 5: https://stackoverflow.com/questions/11279035/can-openprocess-with-error-code-error-access-denied-be-used-to-know-if-process-e
- 6: https://stackoverflow.com/questions/592256/fast-way-to-determine-if-a-pid-exists-on-windows
- 7: https://stackoverflow.com/questions/73995991/how-can-one-know-if-a-process-in-windows-has-exited-or-not
- 8: https://stackoverflow.com/questions/68102385/why-does-openprocess-return-a-non-0-value-when-the-process-id-is-no-longer-run

---

🏁 Script executed:

```shell
rg -n "processAbsent|waitForAbsence|force stop|Force Stop|OpenProcess|GetExitCodeProcess|STILL_ACTIVE" internal
```

Repository: marcioaltoe/roundfix

Length of output: 7111

---

🏁 Script executed:

```shell
sed -n '1,220p' internal/store/process_windows.go && printf '\n---\n' && sed -n '1,240p' internal/store/process.go
```

Repository: marcioaltoe/roundfix

Length of output: 7233

---



</details>

**Check the Windows liveness probe before treating a PID as present**   

`OpenProcess` can succeed for a terminated process object, so `processAbsent` can mark a dead owner as present. Because `TerminateAndWait` uses that result as the proof of life, Force Stop can sit in `waitForAbsence` until `stopWindow` expires and then fail with `prove exit after force kill: context deadline exceeded` even though the owner has already exited. Use `GetExitCodeProcess` (or `WaitForSingleObject`) and treat `STILL_ACTIVE` as the only live state.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 ast-grep (0.44.1)</summary>

[warning] 18-18: Narrowing a non-constant integer to a smaller fixed-width type (int8/int16/int32, uint8/uint16/uint32) can silently overflow or wrap, yielding negative or truncated values that are dangerous in size, length, or index logic. Validate the source value is within the target type's range before converting (e.g. bounds-check, or use a checked helper), and avoid narrowing untrusted or len()/parsed values.
Context: uint32(pid)
Note: [CWE-190] Integer Overflow or Wraparound.

(integer-overflow-narrowing-conversion-go)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/process_windows.go` around lines 18 - 28, Update processAbsent
to verify the opened process’s liveness after syscall.OpenProcess succeeds,
using GetExitCodeProcess or WaitForSingleObject; treat only STILL_ACTIVE as
present and return absent for terminated processes, while preserving existing
error handling and closing the handle.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:43e1e4b5431797e8cd4b57ab -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `processAbsent` treated every successfully opened process object as live, including exited objects. It now calls `GetExitCodeProcess` and treats only `STILL_ACTIVE` as present; a Windows regression covers an exited, unreaped child whose process object remains open. Fresh Batch 001 evidence: `rtk proxy env GOOS=windows GOARCH=amd64 GOCACHE=/tmp/roundfix-batch001-windows-cache go test -c -o /tmp/roundfix-store-windows-batch001.test ./internal/store` passed.
