---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/store/process_windows.go
line: 106
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9YG,comment:PRRC_kwDOS0qyts7dnSbn
review_hash: dc9269608146fbd9fd9a8a756f684e605d65ea1f9c6ff407b3d277e9a6169566
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 010: _ Functional Correctness_ _ Major_ _ Heavy lift_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Check whether the Run launcher already creates a Job object or records process start identities on Windows.
set -euo pipefail

fd -t f -e go . internal/store internal/cli --exec rg -ln 'JobObject|CREATE_NEW_PROCESS_GROUP|SysProcAttr' {} \; || true
rg -n 'GetProcessTimes|JobObject|CreationTime' --type=go || echo "no job-object or creation-time handling found"
rg -n -A 20 'func processStartIdentity' internal/store/process_windows.go
```

Repository: marcioaltoe/roundfix

Length of output: 353

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- process_windows.go outline ---'
ast-grep outline internal/store/process_windows.go
printf '%s\n' '--- process_windows.go relevant source ---'
cat -n internal/store/process_windows.go | sed -n '1,220p'
printf '%s\n' '--- process helpers and call sites ---'
rg -n -A 35 -B 10 'descendantProcessPIDs|TerminateTreeAndWait|ownedProcessIdentity|processStartIdentity|processParent|CreateToolhelp32Snapshot|Process32First|Process32Next' internal --type go
printf '%s\n' '--- launcher process creation ---'
rg -n -A 30 -B 15 'exec\.Cmd|CommandContext|SysProcAttr|CREATE_NEW_PROCESS_GROUP|Start\(' internal --type go
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- process_windows.go outline ---'
ast-grep outline internal/store/process_windows.go
printf '%s\n' '--- process_windows.go relevant source ---'
cat -n internal/store/process_windows.go | sed -n '1,240p'
printf '%s\n' '--- process helpers and call sites ---'
rg -n -A 35 -B 10 'descendantProcessPIDs|TerminateTreeAndWait|ownedProcessIdentity|processStartIdentity|processParent|CreateToolhelp32Snapshot|Process32First|Process32Next' internal --type go
printf '%s\n' '--- launcher process creation ---'
rg -n -A 30 -B 15 'exec\.Cmd|CommandContext|SysProcAttr|CREATE_NEW_PROCESS_GROUP|Start\(' internal --type go
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- Windows implementation ---'
wc -l internal/store/process_windows.go
cat -n internal/store/process_windows.go | sed -n '1,190p'

printf '%s\n' '--- direct callers ---'
rg -n -A 25 -B 15 'TerminateTreeAndWait|OwnerProcessIdentity|ownerProcessIdentity|descendantProcessPIDs' internal/store internal/cli --glob '*.go' --glob '!**/*_test.go'

printf '%s\n' '--- Windows-specific process creation and ownership metadata ---'
rg -n 'SysProcAttr|CREATE_NEW_PROCESS_GROUP|CREATE_SUSPENDED|JobObject|GetProcessTimes|CreationTime|OwnerPID|OwnerIdentity|owner.*identity|identity.*owner' internal --glob '*.go' --glob '!**/*_test.go' || true
```

Repository: marcioaltoe/roundfix

Length of output: 43113

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
from collections import defaultdict

def descendants(owner, parents):
    children = defaultdict(list)
    for pid, parent in parents:
        children[parent].append(pid)
    out = []
    stack = [owner]
    while stack:
        parent = stack.pop()
        for child in children[parent]:
            out.append(child)
            stack.append(child)
    return out

# The table contains an orphan retaining the owner's recycled PID.
parents = [(500, 100), (700, 500), (800, 700), (900, 500)]
got = descendants(500, parents)
print({"owner_pid": 500, "parent_table": parents, "descendants": got,
       "recycled_owner_pid_chain_included": 700 in got and 800 in got})
PY
```

Repository: marcioaltoe/roundfix

Length of output: 318

---



</details>

**Prevent Windows PID reuse from broadening the termination tree.**

`processStartIdentity` is unsupported on Windows, so `TerminateTreeAndWait` uses PID-only ownership. A stale `ParentProcessID` can therefore include an unrelated process and its descendants after the owner PID is reused. Bind Run processes to a Windows Job object, or record and compare process creation times before signalling. Add a Windows regression test for stale parent links.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/process_windows.go` around lines 92 - 106, Update the Windows
process-tree termination flow around descendantProcessPIDs and
TerminateTreeAndWait to prevent stale ParentProcessID values from including
unrelated processes after PID reuse. Bind run processes to a Windows Job object,
or record and validate process creation times before signalling, while
preserving existing descendant termination behavior. Add a Windows-specific
regression test covering a stale parent link.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d308c3e8b8494fa6c71ac7ed -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Confirmed Windows used PID/ParentPID topology without process creation identity, so PID reuse could make a stale parent link appear owned.
  - Windows Runs now record the process creation time as their owner identity. Tree enumeration reads creation times for descendants and accepts only links whose child was created no earlier than its claimed parent.
  - Added platform-neutral creation-bounded topology coverage plus Windows-specific identity parsing and stale-link regressions, while preserving the pre-existing Windows unreaped-child test.
  - Focused evidence: `GOOS=windows GOARCH=amd64 rtk go test -c ./internal/store -o /private/tmp/roundfix-store-windows-run_20260805T050305Z.exe` passed; native store and complete affected package suites passed.
  - The Daemon owns authoritative `make verify` after this Agent turn.
