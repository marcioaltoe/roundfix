---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: invalid
terminal_reason: "The repository does not enable an exhaustive-switch gate, and the omitted unmanaged case intentionally falls through to retain its warning."
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/classification.go
line: 1071
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymb,comment:PRRC_kwDOS0qyts7cjgE3
review_hash: b893535bf51ad2402927ed507590e4f98862fd366304f2de8efad85fa4fa0f48
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 007: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Check whether the exhaustive linter is enabled.
set -euo pipefail
fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' --exec rg -n -C 3 'exhaustive' {}
```

Repository: marcioaltoe/roundfix

Length of output: 567

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- enum declarations and uses ---'
rg -n -C 6 'type carrierClassificationKind|carrierCurrentManaged|carrierStaleManaged|carrierRepositoryExtension|carrierUnmanagedNested' internal/baseline/classification.go

printf '%s\n' '--- target switch ---'
sed -n '1035,1085p' internal/baseline/classification.go

printf '%s\n' '--- exhaustive configuration ---'
fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' --exec rg -n -C 3 'exhaustive' {}
```

Repository: marcioaltoe/roundfix

Length of output: 6232

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
from pathlib import Path
import re

source = Path("internal/baseline/classification.go").read_text()
config_files = list(Path(".").glob(".golangci.yml")) + list(Path(".").glob(".golangci.yaml")) + list(Path(".").glob(".golangci.toml")) + list(Path(".").glob(".golangci.json"))

enum_block = re.search(
    r"type carrierClassificationKind string\s+const \(\s*(.*?)\n\)",
    source,
    re.S,
)
assert enum_block, "carrierClassificationKind declaration not found"
enum_values = re.findall(r"^\s*(carrier\w+)\s+carrierClassificationKind\b", enum_block.group(1), re.M)

switch_block = re.search(
    r"switch byPath\[warning\.Path\] \{\s*(.*?)\n\s*\}",
    source,
    re.S,
)
assert switch_block, "target switch not found"
switch_values = re.findall(r"\b(carrier\w+)\b", switch_block.group(1))
switch_values = sorted(set(switch_values))
has_default = bool(re.search(r"^\s*default\s*:", switch_block.group(1), re.M))

config_text = "\n".join(path.read_text() for path in config_files)
print("enum_values:", enum_values)
print("switch_values:", switch_values)
print("missing_enum_cases:", sorted(set(enum_values) - set(switch_values)))
print("has_default:", has_default)
print("exhaustive_enabled:", bool(re.search(r"^\s*-\s*exhaustive\b", config_text, re.M)))

assert set(enum_values) - set(switch_values) == {"carrierUnmanagedNested"}
assert not has_default
assert re.search(r"^\s*-\s*exhaustive\b", config_text, re.M)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 551

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

files="$(fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' .)"
printf '%s\n' "$files"

while IFS= read -r file; do
  [ -n "$file" ] || continue
  printf '%s\n' "--- $file ---"
  rg -n -C 2 'exhaustive' "$file"
done <<< "$files"
```

Repository: marcioaltoe/roundfix

Length of output: 573

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- linter configuration and invocations ---'
git ls-files | rg '(^|/)(\.golangci\.(yml|yaml|toml|json)|golangci|workflows/.*\.(yml|yaml))$' || true
rg -n -C 3 'golangci-lint|\.golangci|exhaustive' .github .agents Makefile* Taskfile* Justfile* docs 2>/dev/null || true

printf '%s\n' '--- tracked status of the discovered config ---'
git ls-files --stage .agents/skills/golang-lint/assets/.golangci.yml
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---



</details>

**Add the missing `carrierUnmanagedNested` case**

Keep its warning and add a short comment that states this intent. This satisfies the exhaustive-switch policy.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/classification.go` around lines 1066 - 1071, Add
carrierUnmanagedNested to the switch in the baseline warning filter so its
warning is also retained, and add a brief comment documenting that this case is
intentionally kept to satisfy the exhaustive-switch policy.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:b84314bf34cf6eba5a713e94 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The cited exhaustive configuration exists only inside the `golang-lint` Skill asset. The switch suppresses warnings only for positively classified managed or repository-extension carriers; `carrierUnmanagedNested` correctly falls through and keeps its warning.
