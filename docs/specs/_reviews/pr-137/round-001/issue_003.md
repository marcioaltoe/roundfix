---
source: coderabbit
pr: "137"
round: 1
round_created_at: "2026-08-07T03:22:12Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/promoted-backlog-entries-have-a-home
head_sha: ea93c68b70d066c1ee7f322e40ac1d547420e8be
file: internal/speccheck/backlog.go
line: 11
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XKdmw,comment:PRRC_kwDOS0qyts7eg8KO
review_hash: 33bab1bf8900144aae23c3c3fd91b28039ed3c85e048de009d3d6065cabac245
duplicate_of: ""
source_review_id: "4879615443"
source_review_submitted_at: "2026-08-07T03:12:07Z"
---

# Issue 003: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' 'YAML module declarations:'
rg -n 'gopkg\.in/yaml\.v3' go.mod go.sum || true

printf '%s\n' 'Relevant diff:'
git diff -- go.mod go.sum internal/speccheck/backlog.go
```

Repository: marcioaltoe/roundfix

Length of output: 405

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' 'Dependency entries:'
rg -n -C 1 'gopkg\.in/yaml\.v3' go.mod go.sum

printf '%s\n' 'Backlog parser symbols and frontmatter handling:'
ast-grep outline internal/speccheck/backlog.go
rg -n -C 3 'Frontmatter|front matter|yaml|Split|Delimiter|---|CRLF|ReadFile' internal/speccheck/backlog.go
```

Repository: marcioaltoe/roundfix

Length of output: 4055

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' 'Detector implementation:'
sed -n '40,180p' internal/speccheck/backlog.go

printf '%s\n' 'Related tests and Backlog Entries:'
rg -n -C 2 'backlog|Backlog|SC-BACKLOG-UNMOVED|parseBacklogFrontmatter' --glob '*_test.go' --glob '*.md' .

printf '%s\n' 'CRLF files under docs/backlog:'
if [ -d docs/backlog ]; then
  find docs/backlog -type f -print0 | xargs -0 -r file
else
  printf '%s\n' 'docs/backlog not present'
fi

printf '%s\n' 'Minimal parser-behavior check:'
python3 - <<'PY'
def parser_recognizes_frontmatter(content: bytes) -> bool:
    opening = b"---\n"
    text = content.decode()
    if not text.startswith(opening.decode()):
        return False
    rest = text[len(opening):]
    return rest.find("\n---") >= 0

for name, content in {
    "LF": b"---\nstatus: promoted\n---\n",
    "CRLF": b"---\r\nstatus: promoted\r\n---\r\n",
}.items():
    print(f"{name}: {'recognized' if parser_recognizes_frontmatter(content) else 'ignored'}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 50381

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' 'Backlog files and line-ending counts:'
python3 - <<'PY'
from pathlib import Path

root = Path("docs/backlog")
if not root.exists():
    print("docs/backlog: absent")
else:
    files = sorted(p for p in root.rglob("*") if p.is_file())
    if not files:
        print("docs/backlog: no files")
    for path in files:
        data = path.read_bytes()
        print(f"{path}: CRLF={data.count(b'\\r\\n')} LF={data.count(b'\\n')}")

def parser_state(content: bytes) -> str:
    opening = "---\n"
    text = content.decode()
    if not text.startswith(opening):
        return "ignored: opening delimiter"
    rest = text[len(opening):]
    if rest.find("\n---") < 0:
        return "ignored: closing delimiter"
    return "parsed"

print("LF:", parser_state(b"---\nstatus: promoted\nspec: example\n---\n"))
print("CRLF:", parser_state(b"---\r\nstatus: promoted\r\nspec: example\r\n---\r\n"))
PY

printf '%s\n' 'Backlog parser tests:'
sed -n '1,125p' internal/speccheck/backlog_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 437

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
from pathlib import Path

root = Path("docs/backlog")
if not root.exists():
    print("docs/backlog: absent")
else:
    files = sorted(p for p in root.rglob("*") if p.is_file())
    if not files:
        print("docs/backlog: no files")
    for path in files:
        data = path.read_bytes()
        crlf = data.count(b"\r\n")
        lf = data.count(b"\n")
        print(f"{path}: CRLF={crlf} LF={lf}")

def parser_state(content: bytes) -> str:
    opening = "---\n"
    text = content.decode()
    if not text.startswith(opening):
        return "ignored: opening delimiter"
    rest = text[len(opening):]
    if rest.find("\n---") < 0:
        return "ignored: closing delimiter"
    return "parsed"

print("LF:", parser_state(b"---\nstatus: promoted\nspec: example\n---\n"))
print("CRLF:", parser_state(b"---\r\nstatus: promoted\r\nspec: example\r\n---\r\n"))
PY

sed -n '1,125p' internal/speccheck/backlog_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 4908

---



</details>

**Handle CRLF frontmatter delimiters.** `parseBacklogFrontmatter` ignores CRLF documents, so it can miss promoted entries. Accept both LF and CRLF delimiters and add a regression test.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 golangci-lint (2.12.2)</summary>

[error] 11-11: could not import gopkg.in/yaml.v3 (internal/rounds/rounds.go:14:2: missing go.sum entry for module providing package gopkg.in/yaml.v3 (imported by roundfix/internal/baseline); to add:
	go get roundfix/internal/baseline)

(typecheck)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/backlog.go` at line 11, Update parseBacklogFrontmatter to
recognize frontmatter delimiters terminated by either LF or CRLF, preserving
existing LF parsing behavior so CRLF documents are parsed and promoted entries
are detected. Add a regression test covering a CRLF-formatted backlog document.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:200b8ced07ac1a97105edbc5 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `parseBacklogFrontmatter` checked only `---\n` and `\n---`, so a
  CRLF document returned empty frontmatter and a promoted Backlog Entry escaped
  `SC-BACKLOG-UNMOVED`.

## Resolution

- Normalized CRLF to LF before delimiter recognition and YAML parsing while
  retaining the existing LF path.
- Added a public `Check` regression case with a fully CRLF-formatted promoted
  Backlog Entry; it asserts `SC-BACKLOG-UNMOVED` and the resolved Spec
  destination.
- Reproduction evidence before the production fix: the focused test failed
  because the CRLF case returned no `SC-BACKLOG-UNMOVED` finding.
- Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/speccheck -run TestCheckBacklogUnmoved -count=1` passed after the
  fix; `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/speccheck -count=1` also passed.
- Authoritative Verification `make verify` was not run; the Daemon owns it
  after this Agent turn.
