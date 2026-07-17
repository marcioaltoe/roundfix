---
source: coderabbit
pr: "32"
round: 3
round_created_at: "2026-07-17T14:20:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: 204bbd00fbc648be0df0b8bf2f883b9e2dc490c8
file: internal/cli/implement.go
line: 462
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ry2aM,comment:PRRC_kwDOS0qyts7Wzc-d
review_hash: 9010fdc0d74f72c034247eeee7fecc29e0832e01e0374e6d81faca86bf4ea231
duplicate_of: ""
source_review_id: "4723452116"
source_review_submitted_at: "2026-07-17T14:16:02Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
if command -v rtk >/dev/null 2>&1; then
  rtk rg -n '\btar\.TypeRegA\b|^go [0-9.]+' internal/cli/implement.go go.mod
else
  rg -n '\btar\.TypeRegA\b|^go [0-9.]+' internal/cli/implement.go go.mod
fi
```

Repository: marcioaltoe/roundfix

Length of output: 239

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail
GOROOT="$(go env GOROOT)"
grep -nE 'TypeRegA|TypeReg ' "$GOROOT/src/archive/tar/common.go" || true
sed -n '1,220p' "$GOROOT/src/archive/tar/common.go" | sed -n '/TypeRegA/,+12p'
```

Repository: marcioaltoe/roundfix

Length of output: 668

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail
if command -v rtk >/dev/null 2>&1; then
  rtk sed -n '430,475p' internal/cli/implement.go
  rtk rg -n 'TypeRegA|TypeReg\b|TypeDir\b|TypeLink\b' internal/cli/implement.go
else
  sed -n '430,475p' internal/cli/implement.go
  rg -n 'TypeRegA|TypeReg\b|TypeDir\b|TypeLink\b' internal/cli/implement.go
fi
```

Repository: marcioaltoe/roundfix

Length of output: 1640

---



</details>

**Remove `tar.TypeRegA`** `tar.TypeRegA` is deprecated in Go 1.26 and will trip SA1019 here; `tar.TypeReg` is sufficient.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 golangci-lint (2.12.2)</summary>

[error] 462-462: SA1019: tar.TypeRegA has been deprecated since Go 1.11 and an alternative has been available since Go 1.1: Use TypeReg instead. 

(staticcheck)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/implement.go` at line 462, Update the tar header type switch
containing `case tar.TypeReg, tar.TypeRegA` to handle only `tar.TypeReg`; remove
the deprecated `tar.TypeRegA` reference while preserving the existing
regular-file handling.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:85c9ce49de5209408ee2e340 -->

_Sources: Coding guidelines, Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `internal/cli/implement.go` still referenced deprecated `tar.TypeRegA` in the archive unpacker. Removed the deprecated case and kept `tar.TypeReg` regular-file handling.

## Resolution

- Changed `unpackTarArchive` to handle regular files with `tar.TypeReg` only.
- Evidence:
  - `rtk rg -n '\btar\.TypeRegA\b' internal/cli/implement.go` — no matches.
  - `rtk go test ./internal/cli ./internal/store` — passed.
  - `rtk make verify` — passed.
