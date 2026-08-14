---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: internal/store/process_other.go
line: 7
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ953,comment:PRRC_kwDOS0qyts7V5taw
review_hash: 86ef909f4c9c35c8cc6a1049002ddd99d251334e001e7bf29d5f5c7fbe3325a4
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 015: _ Stability & Availability_ _ Major_ _ Heavy lift_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf 'Files:\n'
git ls-files internal/store | sed -n '1,200p'

printf '\nOutline process files:\n'
for f in $(git ls-files 'internal/store/process*'); do
  echo "== $f =="
  wc -l "$f"
  sed -n '1,220p' "$f"
done

printf '\nSearch ReclaimOrphanedRun / ProcessAlive usages:\n'
rg -n "ReclaimOrphanedRun|ProcessAlive|process_other" internal -S
```

Repository: marcioaltoe/roundfix

Length of output: 3530

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf 'Platform-specific build files:\n'
git ls-files | rg 'internal/store/process|_windows|_unix|_linux|build tags|go:build|process_.*\.go$' -n || true
```

Repository: marcioaltoe/roundfix

Length of output: 551

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf 'Platform-specific build files:\n'
git ls-files | rg 'internal/store/process|_windows|_unix|_linux|go:build|process_.*\.go$' || true
```

Repository: marcioaltoe/roundfix

Length of output: 511

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

echo '--- process files ---'
git ls-files 'internal/store/process*'
echo '--- search ---'
rg -n "ReclaimOrphanedRun|ProcessAlive" internal/store internal -S
echo '--- build tags ---'
rg -n "^//go:build|^// \+build" internal/store -S
```

Repository: marcioaltoe/roundfix

Length of output: 3833

---



</details>

**Don’t return `true` for every non-Unix PID.** On `!unix` builds, `ProcessAlive` makes `ReclaimOrphanedRun` unable to recover orphaned locks after a crash. Add a real liveness check for supported non-Unix platforms, or make those platforms explicitly unsupported instead of silently disabling recovery.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/process_other.go` around lines 5 - 7, Update ProcessAlive to
avoid unconditionally returning true on non-Unix builds: implement a real
process-liveness check for each supported non-Unix platform, or explicitly mark
unsupported platforms so ReclaimOrphanedRun does not silently skip orphan
recovery. Preserve the existing Unix behavior and ensure the result accurately
reflects whether the target PID is alive.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:529a184d189fb4e809469a07 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added Windows process liveness and limited the always-alive fallback to unsupported non-Unix platforms; `make verify` passed.
