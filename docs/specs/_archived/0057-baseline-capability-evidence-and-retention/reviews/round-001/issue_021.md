---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/source_contracts.go
line: 84
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymx,comment:PRRC_kwDOS0qyts7cjgFQ
review_hash: d43fa966b99596d959015dc2ebd0a4ca1690f2b4d16bb9317e7b9f5b1399fd70
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 021: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Inspect golangci-lint configuration for revive rules and godot settings.
set -euo pipefail

fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' --exec cat -n {}
```

Repository: marcioaltoe/roundfix

Length of output: 8663

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- source_contracts.go ---'
cat -n internal/baseline/source_contracts.go | sed -n '1,120p'

printf '%s\n' '--- related disposition declarations and comments ---'
rg -n -C 3 'Clause(Retained|Moved|Replaced|RepositoryDocument|RepositoryExtension|ReasonedRejection|Unaccounted)|type ClauseDisposition' internal .

printf '%s\n' '--- linter/version metadata ---'
rg -n -C 2 'golangci|revive|godot|go [0-9]' .golangci.yml .golangci.yaml .golangci.toml go.mod go.work 2>/dev/null || true
```

Repository: marcioaltoe/roundfix

Length of output: 22562

---

🌐 Web query:

`revive Go linter exported rule grouped exported const comment behavior`

💡 Result:

In the revive linter, the `exported` rule checks that all exported symbols, including constants, have documentation comments. When constants are declared in a group, revive inspects the associated `ast.CommentGroup` for that declaration block [1]. The linter's behavior for grouped constants is as follows: - The rule looks for a documentation comment attached to the constant declaration [1]. - It specifically processes the `ast.CommentGroup` to find the first non-empty, non-directive, and non-deprecated line to use as the valid comment [1]. - If you have a single comment block preceding a grouped constant declaration (e.g., a `const (...)` block), revive generally expects that block to provide the necessary documentation. If you encounter issues where grouped constants are flagged despite having comments, note that the `exported` rule is highly configurable [2][3]. You can selectively disable checks for specific declaration types using configuration flags in your `revive.toml` or configuration file [2][3]: - `disable-checks-on-constants`: Disables all checks on constant declarations [2][3]. If you only want to suppress specific warnings while keeping the rule active for other symbols, you can use linter directives in your code to ignore the rule for specific lines or blocks [4][5]: //revive:disable:exported const ( ConstA = 1 ConstB = 2) //revive:enable:exported This allows you to bypass the requirement for individual comments on grouped constants if your project's style prefers a single group-level comment.

Citations:

- 1: https://github.com/mgechev/revive/blob/f31bbab8/rule/exported.go
- 2: https://github.com/mgechev/revive/blob/HEAD/RULES_DESCRIPTIONS.md
- 3: https://github.com/mgechev/revive/blob/f31bbab8/RULES_DESCRIPTIONS.md
- 4: https://pkg.go.dev/github.com/mgechev/revive/rule
- 5: https://github.com/mgechev/revive/

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
from pathlib import Path
import re

source = Path("internal/baseline/source_contracts.go").read_text()
config = Path(".golangci.yml").read_text()

const_block = re.search(r"(?m)^const \(\n(?P<body>.*?)^\)", source, re.S)
assert const_block, "Clause disposition const block not found"

names = re.findall(
    r"(?m)^\s*(Clause(?:Retained|Moved|Replaced|RepositoryDocument|RepositoryExtension|ReasonedRejection|Unaccounted))\b",
    const_block.group("body"),
)
assert len(names) == 7, names

before = source[: const_block.start()].splitlines()
has_group_comment = bool(before and before[-1].lstrip().startswith("//"))
revive_enabled = bool(re.search(r"(?m)^\s*-\s*revive\s*(?:#.*)?$", config))
godot_enabled = bool(re.search(r"(?m)^\s*-\s*godot\s*(?:#.*)?$", config))

print({"exportedConstants": names,
       "groupHasDocComment": has_group_comment,
       "reviveEnabled": revive_enabled,
       "godotEnabled": godot_enabled})
assert not has_group_comment
assert revive_enabled and godot_enabled
PY
```

Repository: marcioaltoe/roundfix

Length of output: 682

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

config_file="$(fd -H -t f '\.golangci\.(yml|yaml|toml|json)$' | head -n 1)"
test -n "$config_file"

python3 - "$config_file" <<'PY'
from pathlib import Path
import re
import sys

config_path = Path(sys.argv[1])
source = Path("internal/baseline/source_contracts.go").read_text()
config = config_path.read_text()

const_block = re.search(r"(?m)^const \(\n(?P<body>.*?)^\)", source, re.S)
assert const_block, "Clause disposition const block not found"

names = re.findall(
    r"(?m)^\s*(Clause(?:Retained|Moved|Replaced|RepositoryDocument|RepositoryExtension|ReasonedRejection|Unaccounted))\b",
    const_block.group("body"),
)
assert len(names) == 7, names

before = source[:const_block.start()].splitlines()
has_group_comment = bool(before and before[-1].lstrip().startswith("//"))
revive_enabled = bool(re.search(r"(?m)^\s*-\s*revive\s*(?:#.*)?$", config))
godot_enabled = bool(re.search(r"(?m)^\s*-\s*godot\s*(?:#.*)?$", config))

print({
    "config": str(config_path),
    "exportedConstants": names,
    "groupHasDocComment": has_group_comment,
    "reviveEnabled": revive_enabled,
    "godotEnabled": godot_enabled,
})
assert not has_group_comment
assert revive_enabled and godot_enabled
PY
```

Repository: marcioaltoe/roundfix

Length of output: 251

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

config_file="$(fd -H -t f '\.golangci\.(yml|yaml|toml|json)$' | head -n 1)"
test -n "$config_file"

python3 - "$config_file" <<'PY'
from pathlib import Path
import re
import sys

config_path = Path(sys.argv[1])
source = Path("internal/baseline/source_contracts.go").read_text()
config = config_path.read_text()

type_pos = source.index("type ClauseDisposition string")
const_start = source.index("const (", type_pos)
const_end = source.index("\n)", const_start)
const_body = source[const_start:const_end]

names = re.findall(
    r"(?m)^\s*(Clause(?:Retained|Moved|Replaced|RepositoryDocument|RepositoryExtension|ReasonedRejection|Unaccounted))\b",
    const_body,
)
assert len(names) == 7, names

preceding_lines = source[:const_start].splitlines()
has_group_comment = bool(
    preceding_lines and preceding_lines[-1].lstrip().startswith("//")
)
revive_enabled = bool(re.search(r"(?m)^\s*-\s*revive\s*(?:#.*)?$", config))
godot_enabled = bool(re.search(r"(?m)^\s*-\s*godot\s*(?:#.*)?$", config))

print({
    "config": str(config_path),
    "exportedConstants": names,
    "groupHasDocComment": has_group_comment,
    "reviveEnabled": revive_enabled,
    "godotEnabled": godot_enabled,
})
assert not has_group_comment
assert revive_enabled and godot_enabled
PY
```

Repository: marcioaltoe/roundfix

Length of output: 473

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- tracked linter configuration ---'
git ls-files | rg '(^|/)\.golangci\.(yml|yaml|toml|json)$|golangci'

printf '%s\n' '--- references to the discovered configuration and lint commands ---'
rg -n -C 3 'golangci|\.agents/skills/golang-lint/assets/\.golangci|revive|godot' \
  .github Makefile Taskfile.yml justfile README.md docs .agents 2>/dev/null || true
```

Repository: marcioaltoe/roundfix

Length of output: 42739

---



</details>

**Document the exported clause dispositions.**

Add one doc comment immediately above the `const` block. End the comment with a period.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/source_contracts.go` around lines 76 - 84, Add a Go doc
comment immediately above the clause disposition const block describing the
exported ClauseDisposition values, and ensure the comment ends with a period.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:ae89cc4b4a92d87fd1ad1523 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added a group-level Go doc comment for the exported ClauseDisposition values. `gofmt` completed and `go test ./internal/baseline ./internal/cli` passed with a writable GOCACHE.
