---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/cli/cli.go
line: 2114
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9X6,comment:PRRC_kwDOS0qyts7dnSbZ
review_hash: de3aa3964dbcf68df75194f539e24a78cbd57db7d3e84586992949d175f76fd2
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 007: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Find every hardcoded Run Branch prefix and the canonical constant.
rg -nP -C 3 'roundfix/run-' --type=go
rg -nP -C 3 'runBranchPrefix\s*=' internal/worktree
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
#!/bin/bash
set -e
printf '%s\n' '--- candidate files ---'
fd -i 'worktree|cli.go' . --type f | head -80
printf '%s\n' '--- branch naming symbols and literals ---'
rg -n -C 4 'BranchName|runBranchPrefix|roundfix/run-|RunIDFromBranchName|filterPendingRunWorkByTarget' --glob '*.go' .
printf '%s\n' '--- relevant file outlines ---'
for f in $(fd -i 'worktree.go|cli.go' . --type f | head -20); do
  printf '\n### %s\n' "$f"
  wc -l "$f"
  ast-grep outline "$f" 2>/dev/null | rg 'BranchName|filterPendingRunWorkByTarget|func|type' | head -120 || true
done
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -e
printf '%s\n' '--- cli imports and target function ---'
sed -n '1,80p' internal/cli/cli.go
sed -n '2070,2248p' internal/cli/cli.go
printf '%s\n' '--- canonical branch helpers ---'
sed -n '1,45p' internal/store/store.go
sed -n '1880,1910p' internal/worktree/worktree.go
printf '%s\n' '--- package declarations and imports ---'
sed -n '1,35p' internal/worktree/worktree.go
rg -n 'internal/(store|worktree)' internal/cli/cli.go
printf '%s\n' '--- production hardcoded prefix occurrences only ---'
rg -n 'roundfix/run-' --glob '*.go' --glob '!**/*_test.go' .
```

Repository: marcioaltoe/roundfix

Length of output: 12138

---

🏁 Script executed:

```shell
#!/bin/bash
set -e
python3 - <<'PY'
from pathlib import Path
import re

store = Path("internal/store/store.go").read_text()
worktree = Path("internal/worktree/worktree.go").read_text()
cli = Path("internal/cli/cli.go").read_text()

canonical = re.search(r'RunBranchPrefix\s*=\s*"([^"]+)"', store)
alias = re.search(r'runBranchPrefix\s*=\s*store\.RunBranchPrefix', worktree)
builder = re.search(r'func BranchName\(runID string\) string \{\s*return runBranchPrefix \+ strings\.TrimSpace\(runID\)', worktree)
consumer = re.search(r'runID := strings\.TrimPrefix\(work\.Branch, "([^"]+)"\)', cli)

assert canonical, "store.RunBranchPrefix was not found"
assert alias, "worktree.runBranchPrefix does not alias store.RunBranchPrefix"
assert builder, "worktree.BranchName does not use the canonical alias"
assert consumer, "cli branch-prefix consumer was not found"

canonical_prefix = canonical.group(1)
consumer_prefix = consumer.group(1)
print(f"canonical store.RunBranchPrefix: {canonical_prefix!r}")
print(f"worktree alias present: {bool(alias)}")
print("worktree.BranchName uses the alias: yes")
print(f"cli filter literal: {consumer_prefix!r}")
print(f"prefixes currently equal: {canonical_prefix == consumer_prefix}")

for changed in ("roundfix/run-v2-", "roundfix/job-"):
    generated = changed + "run_01"
    extracted = generated[len(consumer_prefix):] if generated.startswith(consumer_prefix) else generated
    is_run_branch = extracted != generated and extracted.strip() != ""
    print(f"if canonical prefix changed to {changed!r}: generated={generated!r}, cli recognizes={is_run_branch}")
PY
printf '%s\n' '--- Go language version ---'
sed -n '1,20p' go.mod
```

Repository: marcioaltoe/roundfix

Length of output: 1156

---



</details>

**Use the canonical Run Branch parser.**

The CLI duplicates `store.RunBranchPrefix`. If the prefix changes, this code can derive an incorrect Run ID or classify the branch as unrelated. Add an inverse helper next to `runworktree.BranchName` and use it here.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 2110 - 2114, Replace the inline
"roundfix/run-" trimming in the worktree filtering flow with a canonical inverse
helper added alongside runworktree.BranchName. Use that helper to parse
work.Branch into a Run ID, preserving the existing handling for unrelated or
blank branch names, and rely on store.RunBranchPrefix through the shared helper
rather than duplicating the prefix.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:6297fa86be6db23e9d3307d3 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Confirmed CLI code duplicated the Run Branch prefix contract with a string literal.
  - Added `worktree.RunIDFromBranchName`, switched Branch Integrity filtering to it, and added canonical, whitespace, foreign-branch, and empty-ID tests.
  - Focused evidence: `rtk go test ./internal/worktree -run 'TestRunIDFromBranchName|TestApplyRunBranchCandidate'` passed (9 tests); complete affected package suites passed (1,247 tests).
  - The Daemon owns authoritative `make verify` after this Agent turn.
