---
source: coderabbit
pr: "54"
round: 1
round_created_at: "2026-07-31T12:25:20Z"
status: invalid
terminal_reason: "SysctlKinfoProc returned EIO for a reaped PID, violating Task 01's required ESRCH process-gone signal."
head_repository: marcioaltoe/roundfix
head_branch: ma/0055-owner-identity-without-fork
head_sha: 995943945afc356f26b6d42995cb09ed8391f661
file: internal/store/process_darwin.go
line: 27
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vaa8r,comment:PRRC_kwDOS0qyts7b9pvE
review_hash: 16a4033506d8f3d0fc01a6d5a2d3fd7c8f5ff6f173855c13102e64e3524a9438
duplicate_of: ""
source_review_id: "4828408546"
source_review_submitted_at: "2026-07-31T12:20:25Z"
---

# Issue 002: _ Security & Privacy_ _ Trivial_ _ Quick win_

## Review Comment

_🔒 Security & Privacy_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🌐 Web query:

`Does golang.org/x/sys/unix SysctlKinfoProc return EIO or ESRCH when the sysctl kern.proc.pid returns zero bytes for a nonexistent PID?`

💡 Result:

No, `golang.org/x/sys/unix.SysctlKinfoProc` does not explicitly return `EIO` or `ESRCH` when the underlying `sysctl` call for `kern.proc.pid` returns zero bytes for a nonexistent PID [1]. Instead, the behavior depends on the result of the `sysctl` system call itself, which is wrapped by `SysctlKinfoProc` [1]. When querying `kern.proc.pid` for a nonexistent PID on Darwin, the `sysctl` call typically returns an error (often `ESRCH`, "No such process"), which is then returned directly by `SysctlKinfoProc` [1]. If the `sysctl` call were to succeed but return zero bytes, the `SysctlKinfoProc` implementation would likely return `EIO` (I/O error) [1]. This occurs because `SysctlKinfoProc` performs a check after the system call: if the amount of data written to the buffer (`n`) does not equal the `SizeofKinfoProc`, it returns `EIO` [1]. In summary: 1. If the underlying `sysctl` fails (e.g., with `ESRCH`), `SysctlKinfoProc` propagates that error [1]. 2. If the `sysctl` returns zero bytes (or any size other than `SizeofKinfoProc`), `SysctlKinfoProc` returns `EIO` due to its internal size validation check [1].

Citations:

- 1: https://go.googlesource.com/sys.git/+/refs/tags/v0.29.0/unix/syscall_darwin.go

---

🏁 Script executed:

```shell
printf '%s\n' '--- module version ---'
rg -n 'golang.org/x/sys' go.mod go.sum
printf '%s\n' '--- related source and tests ---'
rg -n -C 8 'processStartIdentity|SysctlKinfoProc|SizeofKinfoProc|unix\.ESRCH' internal/store .
```

Repository: marcioaltoe/roundfix

Length of output: 46794

---

🏁 Script executed:

```shell
set -eu
url='https://raw.githubusercontent.com/golang/sys/v0.45.0/unix/syscall_darwin.go'
tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
curl -fsSL "$url" -o "$tmp"
printf '%s\n' '--- v0.45.0 SysctlKinfoProc implementation ---'
rg -n -A 18 -B 4 'func SysctlKinfoProc|func SysctlRaw' "$tmp"
printf '%s\n' '--- referenced constants and types ---'
curl -fsSL 'https://raw.githubusercontent.com/golang/sys/v0.45.0/unix/zerrors_darwin_amd64.go' | rg -n 'EIO|ESRCH' | head
curl -fsSL 'https://raw.githubusercontent.com/golang/sys/v0.45.0/unix/types_darwin.go' | rg -n -A 5 -B 3 'KinfoProc|SizeofKinfoProc' | head -80
```

Repository: marcioaltoe/roundfix

Length of output: 1856

---

🌐 Web query:

`Darwin sysctl kern.proc.pid nonexistent process returns ESRCH or zero bytes KERN_PROC_PID`

💡 Result:

On Darwin (macOS) and related BSD-derived systems, the sysctl interface using KERN_PROC_PID to query information about a nonexistent process is designed to return the error ESRCH ("No such process") [1][2][3]. When invoking sysctl(2) with the MIB (Management Information Base) identifying KERN_PROC_PID and a specific PID, the kernel performs a lookup for that process [4][3][5]. If no process matches the specified PID, the system call fails and sets errno to ESRCH [1][2][3]. Historically, there was some variation in this behavior across different BSD implementations [1][2][6]. In older versions of some systems, the call might return zero (indicating success) but write zero bytes to the provided buffer if no matching process was found, which required callers to check the buffer or size to detect the failure [1][2][6]. Modern implementations have largely standardized the behavior to explicitly return ESRCH when the target process cannot be found, ensuring that the interface is predictable and that callers can easily distinguish between a successful query (which populates the provided buffer with process information) and a failure [3][7][8]. In summary, the expected behavior for KERN_PROC_PID on a nonexistent process is for the sysctl call to return -1 and set errno to ESRCH [3][7].

Citations:

- 1: https://groups.google.com/g/muc.lists.freebsd.hackers/c/tXrid4AyznE
- 2: https://lists.freebsd.org/pipermail/freebsd-current/2004-March/023256.html
- 3: https://mail-index.netbsd.org/source-changes-hg/2018/08/01/msg028973.html
- 4: https://github.com/openbsd/src/blob/d1ecb0c527101f98d2aabc716a8cdf6be4f8bea8/sys/kern/kern_sysctl.c
- 5: https://man.openbsd.org/OpenBSD-7.4/sysctl.2
- 6: https://lists.freebsd.org/pipermail/freebsd-current/2004-March/023224.html
- 7: https://mail-index.netbsd.org/source-changes-hg/2018/07/23/msg003449.html
- 8: https://lists.freebsd.org/pipermail/svn-src-head/2017-December/107458.html

---



</details>

**Use `unix.SysctlKinfoProc` for the Darwin process lookup.**

In `golang.org/x/sys v0.45.0`, the function propagates `ESRCH` from `sysctl` and returns `EIO` only when the successful result size differs from `SizeofKinfoProc`. Do not map every `EIO` to `ESRCH`, because that would hide other size mismatches.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/process_darwin.go` around lines 13 - 27, Update
processStartIdentity to use unix.SysctlKinfoProc for the Darwin process lookup
instead of SysctlRaw and manual unsafe decoding. Preserve the existing formatted
process identity and contextual error behavior, while relying on
SysctlKinfoProc’s ESRCH and successful-size mismatch errors without remapping
EIO.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:dce4f96bd3ab346c8beb90d8 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The proposed helper does not preserve the Spec contract on the current
  Darwin host. `golang.org/x/sys/unix.SysctlKinfoProc` returns `EIO` when the
  kernel successfully returns a zero-length result, while Task 01 requirement
  6 requires a nonexistent process to remain distinguishable as `ESRCH`. The
  current implementation checks the raw length before decoding, maps only the
  zero-length absence result to `ESRCH`, and leaves other size mismatches as
  contextual errors. No production change is retained.

## Focused evidence

- With the proposed `SysctlKinfoProc` substitution applied temporarily,
  `rtk env GOCACHE=/private/tmp/roundfix-review-bb7640a06ea4290d-gocache go test -count=1 ./internal/store`
  failed at `TestOwnerProcessIdentityFailsForAbsentProcess`: the reaped PID
  returned `input/output error` (`EIO`), but the contract expected `no such
  process` (`ESRCH`).
- After restoring the contract-preserving implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-review-bb7640a06ea4290d-gocache go test -count=1 ./internal/store`
  exited `0`.
- `golang.org/x/sys@v0.45.0/unix/syscall_darwin.go` confirms that
  `SysctlKinfoProc` returns `EIO` whenever the successful byte count differs
  from `SizeofKinfoProc`; it does not classify a zero-length result as
  `ESRCH`.
- Authoritative `make verify` is Daemon-owned and was not run in this Agent
  turn.
