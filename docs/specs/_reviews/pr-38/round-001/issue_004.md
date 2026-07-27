---
source: coderabbit
pr: "38"
round: 1
round_created_at: "2026-07-27T15:34:32Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-outcome-integrity
head_sha: 9ed57622bb92f138aa3e23d4d59e260ebbff0116
file: internal/store/process_unix.go
line: 51
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UG-O9,comment:PRRC_kwDOS0qyts7aENCV
review_hash: 87909db04f3f6cb9ec0a46874523d302d71becbbc75e376a74999a6189f8f1f6
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260727T152947Z_936cd84aa803ba5d/verification/batch-001-attempt-2.log'
source_review_id: "4788632386"
source_review_submitted_at: "2026-07-27T15:23:14Z"
---


# Issue 004: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Check which platforms this file targets and whether a non-ps fallback exists.
fd -t f 'process_.*\.go' internal/store --exec head -3 {}
rg -nP --type=go -C3 'processStartIdentity' internal
```

Repository: marcioaltoe/roundfix

Length of output: 2875

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

ast-grep outline internal/store/process.go --view expanded
printf '\n---\n'
ast-grep outline internal/store/*_test.go --view expanded
printf '\n---\n'
rg -n 'owner_identity|TerminateAndWait|processStartIdentity|lstart|TZ=UTC|LC_ALL=C|BusyBox|ps' internal README.md docs .github -g '!**/vendor/**'
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

if command -v rtk >/dev/null 2>&1; then
  runner=rtk
else
  runner=""
fi

run() {
  if [ -n "$runner" ]; then
    "$runner" "$@"
  else
    "$@"
  fi
}

printf '--- internal/store/process.go (selected slices) ---\n'
run sed -n '1,230p' internal/store/process.go | nl -ba | sed -n '1,230p'
printf '\n--- process tests ---\n'
run rg -n -C 3 'owner_identity|TerminateAndWait|processStartIdentity|lstart|ErrOwnerProcessUnsupported|DST|timezone|locale' internal/store/*_test.go
printf '\n--- process unix implementation ---\n'
run sed -n '1,140p' internal/store/process_unix.go | nl -ba | sed -n '1,140p'
```

Repository: marcioaltoe/roundfix

Length of output: 250

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

if command -v rtk >/dev/null 2>&1; then
  rtk rg -n -C 2 'owner_identity|processStartIdentity|TerminateAndWait|ErrOwnerProcessUnsupported|unix build|ps -p|lstart=' internal docs README.md
else
  rg -n -C 2 'owner_identity|processStartIdentity|TerminateAndWait|ErrOwnerProcessUnsupported|unix build|ps -p|lstart=' internal docs README.md
fi
```

Repository: marcioaltoe/roundfix

Length of output: 15620

---



</details>

**Pin the `ps` rendering environment or stop using formatted output for owner identity.** `owner_identity` is compared verbatim later, so `ps -o lstart=` can stop matching after a timezone/locale change and block a genuine owner from being stopped. Set `TZ=UTC LC_ALL=C` for this command, or switch to a raw start-time token.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/process_unix.go` around lines 38 - 51, Update
processStartIdentity to make the ps output deterministic by setting TZ=UTC and
LC_ALL=C in the command environment before execution, while preserving the
existing verbatim identity comparison and error handling.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:f0842addfe3a878fb32065f2 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `ps -o lstart=` inherited the caller timezone and produced different opaque identity tokens for the same process. The command now pins `TZ=UTC` and `LC_ALL=C`. Before the fix, the regression observed Honolulu `"Mon Jul 27 05:42:21 2026"` versus Tokyo `"Tue Jul 28 00:42:21 2026"`; after the fix, `rtk proxy env GOCACHE=/tmp/roundfix-run-936cd84aa803ba5d-gocache go test ./internal/store -run '^(TestOwnerProcessIdentityIgnoresCallerTimezone|TestOwnerProcessControllerMatchingOwnerIdentityProceeds)$' -count=1` passed.
