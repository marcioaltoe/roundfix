---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/selection_assignment.go
line: 226
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWim,comment:PRRC_kwDOS0qyts7fswOJ
review_hash: 80635464c7ec39f45f272bd49268bf5f8393a76093230b2a025bfdadfce804d6
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:14Z"
---

# Issue 007: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Find callers and readers of the symbols added by this PR.
set -euo pipefail

for symbol in readRuntimeCatalogue contaminatedClaudeCapabilityPayload; do
  echo "== ${symbol} =="
  rg -nP --type=go -C3 "\b${symbol}\s*\("
done

echo '== Contaminated readers =='
rg -nP --type=go -C3 '\bContaminated\b'
```

Repository: marcioaltoe/roundfix

Length of output: 185

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

echo '== symbol references =='
for symbol in readRuntimeCatalogue readRuntimeCatalogueWithEvidence contaminatedClaudeCapabilityPayload Contaminated; do
  echo "== ${symbol} =="
  rg -n --type=go -C2 "\b${symbol}\b" . || true
done

echo '== selection assignment structure =='
sed -n '1,280p' internal/agent/selection_assignment.go

echo '== characterization fixture and nearby tests =='
sed -n '100,180p' internal/agent/selection_catalogue_characterization_test.go

echo '== relevant linter configuration =='
rg -n -C2 'unused|golangci|staticcheck' .golangci.yml .golangci.yaml golangci.yml golangci.yaml Makefile justfile 2>/dev/null || true
```

Repository: marcioaltoe/roundfix

Length of output: 17481

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

echo '== catalogue flow =='
rg -n -C4 'recordAdvertisement|Catalogue:|\.Catalogue|Contaminated|ProveExactSelection|applySessionSelection' internal/agent --glob '*.go' || true

echo '== assignment implementation =='
sed -n '260,520p' internal/agent/selection_assignment.go

echo '== all characterization test names and helper calls =='
rg -n '^(func (Test|Benchmark)|func [a-zA-Z_].*\*testing\.T)|contaminatedClaudeCapabilityPayload|honestClaudeCapabilityPayload' internal/agent/selection_catalogue_characterization_test.go || true

echo '== tracked lint/config files =='
git ls-files | rg '(^|/)(\.golangci\.(yml|yaml)|golangci\.(yml|yaml)|Makefile|Taskfile|justfile|go\.mod)$' || true
```

Repository: marcioaltoe/roundfix

Length of output: 48433

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

python3 - <<'PY'
from pathlib import Path
import re

targets = {
    "readRuntimeCatalogue": "internal/agent/selection_assignment.go",
    "contaminatedClaudeCapabilityPayload": "internal/agent/selection_catalogue_characterization_test.go",
}
for name, expected_file in targets.items():
    declaration = re.compile(rf'\bfunc\s+(?:\([^)]*\)\s+)?{re.escape(name)}\s*\(')
    call = re.compile(rf'\b{re.escape(name)}\s*\(')
    occurrences = []
    for path in Path(".").rglob("*.go"):
        text = path.read_text(errors="replace")
        for match in call.finditer(text):
            line = text.count("\n", 0, match.start()) + 1
            kind = "declaration" if declaration.match(text, match.start() - 5) else "call-or-reference"
            occurrences.append((str(path), line, kind))
    print(f"{name}:")
    for item in occurrences:
        print("  ", item)
    print("  call/reference count excluding declaration:",
          sum(kind == "call-or-reference" for _, _, kind in occurrences))

print("== unused-linter references ==")
for path in Path(".").rglob("*"):
    if path.is_file():
        try:
            text = path.read_text(errors="replace")
        except OSError:
            continue
        if re.search(r'(?i)(golangci|staticcheck|unused)', text):
            print(path)
PY

echo '== all tracked configuration and CI references =='
git ls-files -a | rg '(^|/)(\.golangci\.(yml|yaml)|golangci\.(yml|yaml)|Makefile|Taskfile|justfile|\.github/.*)$' || true
rg -n -i -C2 'golangci|staticcheck|unused' --glob '!*.sum' --glob '!vendor/**' . || true

echo '== module version =='
sed -n '1,20p' go.mod
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

python3 - <<'PY'
from pathlib import Path
import re

targets = [
    "readRuntimeCatalogue",
    "contaminatedClaudeCapabilityPayload",
]
for name in targets:
    declaration = re.compile(
        rf'^\s*func\s+(?:\([^)]*\)\s+)?{re.escape(name)}\s*\(',
        re.M,
    )
    call = re.compile(rf'\b{re.escape(name)}\s*\(')
    refs = []
    for path in Path(".").rglob("*.go"):
        text = path.read_text(errors="replace")
        for match in call.finditer(text):
            line_start = text.rfind("\n", 0, match.start()) + 1
            line = text.count("\n", 0, match.start()) + 1
            if declaration.match(text, line_start):
                kind = "declaration"
            else:
                kind = "call-or-reference"
            refs.append((str(path), line, kind))
    print(name)
    for ref in refs:
        print(" ", ref)
    print(" non-declaration references:",
          sum(ref[2] == "call-or-reference" for ref in refs))

print("RuntimeCatalogue.Contaminated references:")
field = re.compile(r'\bContaminated\b')
for path in Path("internal/agent").rglob("*.go"):
    text = path.read_text(errors="replace")
    for match in field.finditer(text):
        line = text.count("\n", 0, match.start()) + 1
        print(f"  {path}:{line}")
PY

echo '== active CI lint references =='
rg -n -C3 'golangci-lint|staticcheck|unused|go vet' \
  .github/workflows Makefile 2>/dev/null || true

echo '== root lint configuration candidates =='
find . -maxdepth 2 -type f \( -name '.golangci.yml' -o -name '.golangci.yaml' -o -name '.golangci.toml' -o -name '.golangci.json' \) -print
```

Repository: marcioaltoe/roundfix

Length of output: 1354

---



</details>

**Remove the unreferenced helper functions.** Remove `readRuntimeCatalogue` and `contaminatedClaudeCapabilityPayload`. Keep `RuntimeCatalogue.Contaminated`; `recordAdvertisement` uses it and `SelectionProof` returns it.

<details>
<summary>📍 Affects 2 files</summary>

- `internal/agent/selection_assignment.go#L155-L226` (this comment)
- `internal/agent/selection_catalogue_characterization_test.go#L152-L161`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/selection_assignment.go` around lines 155 - 226, Remove the
unreferenced readRuntimeCatalogue helper from
internal/agent/selection_assignment.go and remove
contaminatedClaudeCapabilityPayload from
internal/agent/selection_catalogue_characterization_test.go; retain
RuntimeCatalogue.Contaminated, recordAdvertisement, and SelectionProof behavior
unchanged.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/agent/selection_assignment.go</file>
<line_range>155-226</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/agent/selection_catalogue_characterization_test.go</file>
<line_range>152-161</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:2eb40a3aac30d1889d977fd8 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Removed the unreferenced `readRuntimeCatalogue` helper from internal/agent/selection_assignment.go and `contaminatedClaudeCapabilityPayload` from internal/agent/selection_catalogue_characterization_test.go. Grep confirms neither has any call or reference beyond its declaration; `RuntimeCatalogue.Contaminated`, `recordAdvertisement`, and `SelectionProof` behavior are unchanged. Focused: `go test ./internal/agent -run 'TestRuntimeCatalogue|TestSelectionCatalogueCharacterization'` ok.
