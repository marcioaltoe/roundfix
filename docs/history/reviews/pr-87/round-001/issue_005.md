---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T15:34:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: a12c1a665c5970773e04c4a145c6b9b0c5a0e686
file: internal/daemon/task_engine.go
line: 1637
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WBeN2,comment:PRRC_kwDOS0qyts7c2V0d
review_hash: 2cd60de407cb31bee0ef956917e7adefe4e18183d95d6db388abe35bcb5bdeae
duplicate_of: ""
source_review_id: "4845660382"
source_review_submitted_at: "2026-08-03T15:14:34Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Check how ensureCommitPath handles an empty path argument.
set -euo pipefail

ast-grep run --pattern 'func ensureCommitPath($$$) $_ { $$$ }' --lang go internal/daemon
ast-grep run --pattern 'func artifactCommitPath($$$) $_ { $$$ }' --lang go internal/daemon
```

Repository: marcioaltoe/roundfix

Length of output: 1175

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- task engine: report settlement and commit path ---'
sed -n '1450,1665p' internal/daemon/task_engine.go

printf '%s\n' '--- task engine: path helpers ---'
sed -n '1488,1510p;1780,1810p' internal/daemon/task_engine.go

printf '%s\n' '--- relevant tests and status assertions ---'
rg -n -C 8 'missing report|wantCommit|commitQAReport|ensureCommitPath|StatusFailed' internal/daemon --glob '*_test.go'
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- commitQAReport implementation and callers ---'
rg -n -C 20 'func \(engine \*Engine\) commitQAReport|commitQAReport\(' internal/daemon/task_engine.go

printf '%s\n' '--- focused QA tests ---'
rg -n -C 12 -i 'missing report|report path|qa.?gate|qa report|wantCommit' internal/daemon/task_engine_test.go internal/daemon/engine_test.go

printf '%s\n' '--- commit abstraction calls in the focused implementation range ---'
rg -n -C 10 'ensureCommitPath|artifactCommitPath|Commit\(' internal/daemon/task_engine.go
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- complete QA commit path ---'
sed -n '1735,1805p' internal/daemon/task_engine.go

printf '%s\n' '--- stageable-path filtering ---'
rg -n -C 18 'func FilterStageablePaths|FilterStageablePaths\(' internal/daemon

printf '%s\n' '--- commit request recording and missing-report assertions ---'
sed -n '4700,4800p' internal/daemon/task_engine_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 22616

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- exact stage path resolution ---'
sed -n '1394,1455p' internal/daemon/task_engine.go

printf '%s\n' '--- empty-path and QA-commit tests ---'
rg -n -C 10 'FilterStageablePaths|stagePathInWorktree|empty path|missing report|QA Report commit' internal/daemon/*_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 18887

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- commit path handling ---'
rg -n -C 12 'type .*Commit|func .*Commit|git.*add|Paths:' internal --glob '*.go' | head -n 240

printf '%s\n' '--- standalone empty-path verifier ---'
python3 - <<'PY'
import posixpath

def path_stays_inside(relative):
    return relative != ".." and not relative.startswith("../") and not relative.startswith("/")

clean = posixpath.normpath("")
print(f"clean({''!r}) = {clean!r}")
print(f"pathStaysInside(clean) = {path_stays_inside(clean)}")
print("FilterStageablePaths classification: kept as '.'; it is not dropped")
print("A commit path of '.' can include unrelated worktree changes.")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 19597

---



</details>

**Commit the settled gate Task without staging an empty report path.**

The missing-report branch leaves the settled Task uncommitted. Removing it directly is unsafe because `ensureCommitPath("", ...)` becomes `"."`, which `FilterStageablePaths` keeps and can stage unrelated worktree changes. Ignore empty paths, then commit the gate Task.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/task_engine.go` around lines 1633 - 1637, Update the
missing-report handling in the task commit flow around commitQAReport: ignore
the empty report path when building or staging commit paths so
ensureCommitPath("", ...) cannot produce "." or stage unrelated changes, while
still committing the settled gate Task. Remove the early return for reportPath
== "" and preserve normal report-path processing when a path is present.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:36bdb85258a9ba7936ed78f0 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The missing-report path now reaches the QA commit boundary, omits an empty report path before filtering, and commits only the settled gate Task file. Verification attempt 1 exposed a stale CLI matrix expectation that still coupled report presence to commit presence; the integration test now expects the settled gate Task commit in every verdict case while requiring a report path only when a report exists. The focused daemon and CLI missing-report regressions passed, and `go test ./internal/cli -count=1` passed with 932 tests.
