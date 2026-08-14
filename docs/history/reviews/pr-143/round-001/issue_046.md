---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/baseline/plan_test.go
line: 2160
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiApO,comment:PRRC_kwDOS0qyts7fC8Rh
review_hash: 8424ff0a85922771c64b57ab2daa9418d588fd2c5a278eab02cf7638b02361d1
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:49Z"
---

# Issue 046: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Check that the newly excluded carriers keep golden coverage elsewhere.
rg -nP -C4 'TestFormatterComposition' --type=go
rg -nP 'docs/agents/(backend|domain|frontend|issue-tracker|monorepo|spec-routing)\.md' --type=go
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- target test section ---'
sed -n '2100,2180p' internal/baseline/plan_test.go

printf '%s\n' '--- formatter composition tests and golden references ---'
rg -n -C5 'FormatterComposition|formatter composition|golden|Golden' internal --glob '*.go' --glob '*.golden' --glob '*.yaml' --glob '*.json' || true

printf '%s\n' '--- target carrier references ---'
rg -n -C3 'docs/agents/(backend|domain|frontend|issue-tracker|monorepo|spec-routing)\.md' . --glob '*.go' --glob '*.md' --glob '*.yaml' --glob '*.json' || true

printf '%s\n' '--- relevant files ---'
git ls-files | rg '(^|/)(plan_test\.go|.*golden.*|.*formatter.*|docs/agents/(backend|domain|frontend|issue-tracker|monorepo|spec-routing)\.md)$' || true
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu
log=/tmp/coderabbit-shell-logs/shell-output-nbM61v

printf '%s\n' '--- section markers ---'
rg -n --max-count 20 '^---|TestFormatterComposition|evolvedPastFrozenCorpus' "$log" || true

printf '%s\n' '--- target test output ---'
awk '
  /^--- target test section ---/ {show=1}
  /^--- formatter composition tests and golden references ---/ {show=0}
  show {print}
' "$log"

printf '%s\n' '--- focused formatter test definitions ---'
rg -n -C8 'TestFormatterComposition|evolvedPastFrozenCorpus' internal/baseline/plan_test.go internal --glob '*_test.go' --glob '!docs/**' || true

printf '%s\n' '--- focused target path references in source and tests ---'
rg -n -C2 'docs/agents/(backend|domain|frontend|issue-tracker|monorepo|spec-routing)\.md' internal --glob '*.go' || true
```

Repository: marcioaltoe/roundfix

Length of output: 27401

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- TestFormatterComposition ---'
sed -n '247,390p' internal/baseline/plan_test.go

printf '%s\n' '--- formatter golden assertions and decision-plan tests ---'
sed -n '2440,2635p' internal/baseline/plan_test.go

printf '%s\n' '--- formatter fixture inventory for all six carriers ---'
rg -n -C2 'formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/(backend|domain|frontend|issue-tracker|monorepo|spec-routing)\.md' internal/baseline --glob '*.go' || true

printf '%s\n' '--- focused diff ---'
git diff -- internal/baseline/plan_test.go | sed -n '1,220p'
```

Repository: marcioaltoe/roundfix

Length of output: 14652

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- formatter fixture files for the six carriers ---'
git ls-files | rg '^internal/baseline/.*formatter-fixtures.*/golden/docs/agents/(backend|domain|frontend|issue-tracker|monorepo|spec-routing)\.md$' || true
git ls-files | rg 'formatter-fixtures/standard-typescript-monorepo' | head -80

printf '%s\n' '--- formatter contract declarations ---'
rg -n -C8 '"fixturePaths"|"goldenDigest"|standard-typescript-monorepo' internal/baseline --glob '*.json' --glob '*.go' | head -260

printf '%s\n' '--- complete TestFormatterComposition tail ---'
sed -n '330,470p' internal/baseline/plan_test.go

printf '%s\n' '--- fixture path references in the catalog and test data ---'
rg -n 'docs/agents/(backend|domain|frontend|issue-tracker|monorepo|spec-routing)\.md' internal/baseline --glob '*.json' --glob '*.go' --glob '*.md' | head -300
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import json
from pathlib import Path
import difflib

root = Path("internal/baseline/assets")
profile = json.loads((root / "profiles/standard-typescript-monorepo.json").read_text())
fixtures = profile["formatter"]["fixturePaths"]
targets = {
    "docs/agents/backend.md",
    "docs/agents/domain.md",
    "docs/agents/frontend.md",
    "docs/agents/issue-tracker.md",
    "docs/agents/monorepo.md",
    "docs/agents/spec-routing.md",
}
print("--- formatter contract coverage ---")
for target in sorted(targets):
    matches = [p for p in fixtures if p.endswith("/golden/" + target)]
    print(f"{target}: {len(matches)} fixture(s): {matches}")

print("--- frozen corpus versus maintained golden ---")
corpus_root = root / "source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus"
golden_root = root / "formatter-fixtures/standard-typescript-monorepo/golden"
for target in sorted(targets):
    old = (corpus_root / target).read_text().splitlines(keepends=True)
    new = (golden_root / target).read_text().splitlines(keepends=True)
    diff = list(difflib.unified_diff(
        old, new,
        fromfile="frozen/" + target,
        tofile="golden/" + target,
        n=1,
    ))
    print(f"{target}: frozen_bytes={sum(map(len, map(str.encode, old)))} "
          f"golden_bytes={sum(map(len, map(str.encode, new)))} "
          f"changed_lines={sum(1 for line in diff if line.startswith(('+', '-')) and not line.startswith(('+++', '---')))}")
    for line in diff[:16]:
        print(line.rstrip("\n"))
    if len(diff) > 16:
        print("...")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 7765

---



</details>

**Document the six frozen-corpus exceptions individually**

`TestFormatterComposition` already covers all six paths through the profile’s `formatter.fixturePaths`. Add a brief per-path rationale to `evolvedPastFrozenCorpus`; the frozen-to-golden changes are not uniformly growth, since `frontend.md` and `spec-routing.md` are smaller.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/plan_test.go` around lines 2155 - 2160, Update the
evolvedPastFrozenCorpus expectations in TestFormatterComposition to document a
brief individual rationale for each of the six formatter.fixturePaths, including
that frontend.md and spec-routing.md reflect reductions rather than growth.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:6fb668371a3d0642d063a0a7 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
