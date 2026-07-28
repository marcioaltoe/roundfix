---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/runevent/event.go
line: 154
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpy,comment:PRRC_kwDOS0qyts7aUVEF
review_hash: 990a61872e1ef746c6c0417d6fb1da5bce4cf9c9857001967262d9fb5f58dc60
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 011: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
rg -n '^go ' go.mod
rg -n 'omitzero' --type=go -C2
```

Repository: marcioaltoe/roundfix

Length of output: 168

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf 'go.mod:\n'
sed -n '1,20p' go.mod

printf '\ninternal/runevent/event.go:\n'
nl -ba internal/runevent/event.go | sed -n '110,155p'

printf '\nwatch publisher usage:\n'
rg -n 'ReviewStatusPayload|publishAcceptedEvidence' internal/watch/watch.go -n -C 4
```

Repository: marcioaltoe/roundfix

Length of output: 811

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf 'internal/runevent/event.go (lines 110-155):\n'
sed -n '110,155p' internal/runevent/event.go

printf '\ninternal/watch/watch.go (publishAcceptedEvidence usage):\n'
rg -n 'publishAcceptedEvidence|ReviewStatusPayload' internal/watch/watch.go -C 6

printf '\nSearch for omitzero usage in repo:\n'
rg -n 'omitzero' .
```

Repository: marcioaltoe/roundfix

Length of output: 3770

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
import json
from pathlib import Path

go_mod = Path('go.mod').read_text()
go_version = None
for line in go_mod.splitlines():
    if line.startswith('go '):
        go_version = line.split()[1]
        break
print(json.dumps({"go_version": go_version}))
PY
```

Repository: marcioaltoe/roundfix

Length of output: 181

---



</details>

**Use `omitzero` or pointers for the optional timestamps** `StartedAt` and `Deadline` are `time.Time` values, so `omitempty` still serializes the zero epoch. This module targets Go 1.26, so `json:",omitzero"` is available here; otherwise switch these fields to `*time.Time`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/runevent/event.go` around lines 120 - 139, Update the optional
StartedAt and Deadline fields in ReviewStatusPayload to use the json omitzero
tag so zero time.Time values are omitted from serialization. Keep their
time.Time types and existing field names unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:52c05d319d9e07e3c23e2c50 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Changed optional `time.Time` event fields to `json:",omitzero"` and added a serialization regression proving zero timestamps are omitted. `go test ./internal/runevent` passed.
