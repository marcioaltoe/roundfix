---
source: coderabbit
pr: "113"
round: 1
round_created_at: "2026-08-05T02:12:07Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0068-implementation
head_sha: c9af2617f988bd63e1bd8f22c6178758a8e5fd40
file: internal/cli/spec_check.go
line: 83
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WhZqG,comment:PRRC_kwDOS0qyts7dk96u
review_hash: dbcbf36eea5b2435e83ea069badc47957033f2d5d5d9b648355244c5c16695e1
duplicate_of: ""
source_review_id: "4860420451"
source_review_submitted_at: "2026-08-05T02:11:26Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Compare JSON envelope fields across CLI renderers.
set -uo pipefail

rg -nP --type=go '`json:"(schema|schemaVersion|type|ok)[",]' -g '!**/*_test.go'
echo '--- documents with a schema constant ---'
rg -nP --type=go 'SchemaVersion\s*=\s*"' -g '!**/*_test.go'
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
#!/bin/bash
set -u
printf '%s\n' '--- candidate file ---'
sed -n '1,130p' internal/cli/spec_check.go
printf '%s\n' '--- later location ---'
sed -n '360,415p' internal/cli/spec_check.go
printf '%s\n' '--- sibling symbols ---'
rg -n --type=go 'RenderJSON|reconcile|specAuditDocument|spec audit|SchemaVersion|schemaVersion|json:"(schema|type|ok)' .
printf '%s\n' '--- CLI entry points ---'
rg -n --type=go 'spec.?audit|spec.?check|audit' internal/cli
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -u
printf '%s\n' '--- spec audit implementation ---'
sed -n '120,245p' internal/cli/spec_check.go
sed -n '280,405p' internal/cli/spec_check.go
printf '%s\n' '--- Spec Check JSON renderer ---'
sed -n '1,135p' internal/speccheck/report.go
printf '%s\n' '--- Reconcile envelope and renderer ---'
sed -n '1,75p' internal/cli/reconcile.go
sed -n '225,255p' internal/cli/reconcile.go
sed -n '360,380p' internal/cli/reconcile.go
printf '%s\n' '--- Spec audit tests ---'
sed -n '220,370p' internal/cli/spec_check_test.go
printf '%s\n' '--- Relevant guideline text in repository ---'
rg -n -i -C 2 'schemaVersion.*type.*ok|every CLI response|CLI responses|schema-only' docs README.md internal skills 2>/dev/null | head -160
```

Repository: marcioaltoe/roundfix

Length of output: 27862

---

🏁 Script executed:

```shell
#!/bin/bash
set -u
python3 - <<'PY'
import re
from pathlib import Path

targets = {
    "internal/cli/spec_check.go": "specAuditDocument",
    "internal/speccheck/report.go": "jsonDocument",
    "internal/cli/reconcile.go": "reconcileReport",
}
required = {"schemaVersion", "type", "ok"}

for filename, typename in targets.items():
    text = Path(filename).read_text()
    match = re.search(r"type\s+" + re.escape(typename) + r"\s+struct\s*\{(.*?)\n\}", text, re.S)
    if not match:
        print(f"{filename}: {typename}: NOT FOUND")
        continue
    fields = set(re.findall(r'`json:"([^",]+)', match.group(1)))
    print(f"{filename}: {typename}: fields={sorted(fields)}")
    print(f"  missing required fields={sorted(required - fields)}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 638

---



</details>

**Add the required JSON envelope fields to `specAuditDocument`.**

`spec audit --format json` must emit `schemaVersion`, `type`, and `ok`. Set `ok` from `!specAuditNeedsAttention(result)` so consumers do not need the exit code to determine the audit result. Existing sibling envelopes do not satisfy this requirement.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/spec_check.go` around lines 78 - 83, Extend specAuditDocument
with the required JSON envelope fields schemaVersion, type, and ok, and populate
them in the spec audit JSON output path. Set ok to the inverse of
specAuditNeedsAttention(result), while preserving the existing schema, slug,
survivors, and undelivered fields.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:b4c521d74abe239541f72811 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The JSON document preserved its command schema but omitted the repository's required machine-readable envelope fields.

## Resolution

- Added `schemaVersion`, `type`, and `ok` while preserving `schema`, `slug`, `survivors`, and `undelivered`; `ok` is the inverse of `specAuditNeedsAttention`.
- Added clean and attention-requiring JSON assertions.
- Focused evidence: the regression first failed to compile because the envelope constant and fields were absent; `rtk env GOCACHE=/private/tmp/roundfix-review0068-packages-cache go test ./internal/specaudit ./internal/worktree ./internal/cli ./internal/baseline -count=1` later exited 0.
- Daemon Verification: `make verify` was not run; the Daemon owns that command.
