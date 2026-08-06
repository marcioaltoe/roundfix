---
source: coderabbit
pr: "125"
round: 1
round_created_at: "2026-08-05T19:38:23Z"
status: invalid
terminal_reason: "The archived matrix explicitly labels the watch row as fresh table-test coverage and separately records supervised equivalence; it does not claim a hosted public watch observation."
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 1384d928c1d0af4d0bca06e506eb8b2953f9f341
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/preflight-table.md
line: 14
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WxMJN,comment:PRRC_kwDOS0qyts7d790b
review_hash: b44320e200f6fda7a42f6f9565e4238e1d32200c68452fff6efddad5b143c616
duplicate_of: ""
source_review_id: "4868102567"
source_review_submitted_at: "2026-08-05T19:37:35Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -eu

target='docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/preflight-table.md'

printf '%s\n' '--- target file ---'
cat -n "$target"

printf '%s\n' '--- related evidence files ---'
git ls-files 'docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/*'

printf '%s\n' '--- relevant references in the archived Spec 0078 evidence ---'
rg -n -i -C 3 'watch|pushTriggersReview|request_review|true \| false|false \| true' \
  'docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078' \
  || true

printf '%s\n' '--- public watch invocations in the repository ---'
rg -n -i -C 2 'roundfix([[:space:]]+|[-_])?watch|(^|[[:space:]])watch([[:space:]]|$)' \
  docs scripts .github 2>/dev/null \
  || true
```

Repository: marcioaltoe/roundfix

Length of output: 50379

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

target='docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/preflight-table.md'

cat -n "$target"
git ls-files 'docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/*'
rg -n -i -C 3 'watch|pushTriggersReview|request_review|true \| false|false \| true' \
  'docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078' || true
```

Repository: marcioaltoe/roundfix

Length of output: 16423

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

target='docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/preflight-table.md'
evidence='docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078'

printf '%s\n' '--- repository status ---'
git status --short

printf '%s\n' '--- target diff summary and diff ---'
git diff --stat -- "$target"
git diff -- "$target"

printf '%s\n' '--- focused related evidence ---'
cat -n "$evidence/repository-config.md"
cat -n "$evidence/equivalent-pr-evidence.md"

printf '%s\n' '--- archive-preservation guidance references ---'
rg -n -i -C 2 'archived|proven-completed|byte-identical|legacy Spec' \
  docs CONTEXT.md 2>/dev/null \
  | head -n 160 || true

printf '%s\n' '--- row-specific public watch commands in Spec 0078 evidence ---'
rg -n -i -C 2 '(^|`|[[:space:]])(rtk[[:space:]]+)?(bin/)?roundfix[[:space:]]+watch|public .*watch|built .*watch' \
  "$evidence" || true
```

Repository: marcioaltoe/roundfix

Length of output: 23968

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

evidence='docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078'

git status --short
git diff --stat -- "$evidence/preflight-table.md"
git diff -- "$evidence/preflight-table.md"

cat -n "$evidence/repository-config.md"
cat -n "$evidence/equivalent-pr-evidence.md"

rg -n -i -C 2 'archived|proven-completed|byte-identical|legacy Spec' \
  docs CONTEXT.md 2>/dev/null | head -n 160 || true

rg -n -i -C 2 'roundfix[[:space:]]+watch|public .*watch|built .*watch' \
  "$evidence" || true
```

Repository: marcioaltoe/roundfix

Length of output: 23758

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
from pathlib import Path
import re

root = Path("docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078")
table = (root / "preflight-table.md").read_text()
config = (root / "repository-config.md").read_text()
equivalent = (root / "equivalent-pr-evidence.md").read_text()

rows = []
for line in table.splitlines():
    if line.startswith("| ") and re.match(r"\|\s*(true|false)\s*\|\s*(true|false)\s*\|", line):
        cells = [c.strip() for c in line.strip("|").split("|")]
        rows.append(cells)

target = next(row for row in rows if row[0:2] == ["true", "false"])
public_watch_mentions = re.findall(r"(?i)(?:built public command|public .*watch|roundfix watch)", config + equivalent)

print("target_row:", target[:2])
print("watch_cell_has_command:", bool(re.search(r"(?i)(roundfix|watch --)", target[3])))
print("repository_config_public_watch_pair:", "false/true" if "Pair `false/true`" in config else "not stated")
print("related_public_watch_mentions:", len(public_watch_mentions))
print("equivalent_evidence_explicitly_supervised:", "supervised equivalents" in equivalent)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 347

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
from pathlib import Path
import re

root = Path("docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078")
table = (root / "preflight-table.md").read_text()
config = (root / "repository-config.md").read_text()
equivalent = (root / "equivalent-pr-evidence.md").read_text()

rows = []
for line in table.splitlines():
    if line.startswith("| ") and re.match(r"\|\s*(true|false)\s*\|\s*(true|false)\s*\|", line):
        rows.append([c.strip() for c in line.strip("|").split("|")])

target = next(row for row in rows if row[:2] == ["true", "false"])
print("target_row:", target[:2])
print("watch_cell_has_command:", bool(re.search(r"(?i)(roundfix|watch --)", target[3])))
print("repository_config_public_watch_pair:",
      "false/true" if "Pair `false/true`" in config else "not stated")
print("equivalent_evidence_explicitly_supervised:",
      "supervised equivalents" in equivalent)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 314

---



</details>

**Record public `watch` evidence for the `true | false` pair before archiving.**

The `Public watch observation` cell contains only predicate and table-test coverage. The available built `watch` evidence covers `false | true`. Add the result to active QA evidence, or label this column as test-only; do not edit the archived artifact.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/preflight-table.md`
around lines 9 - 14, Do not modify the archived table; record concrete built
`watch` evidence for the `pushTriggersReview=true` and `request_review=false`
pair in the active QA evidence, or explicitly relabel that column as test-only.
Update the corresponding active evidence or labeling around the `Public watch
observation` column while preserving the existing `false | true` evidence.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:19f666c8fbba300685fcc4ee -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The cited cell says "Covered by the same Preflight predicate and the fresh table test," not that a hosted `watch` invocation occurred. The adjacent evidence records which CLI journeys were observed and which were equivalent supervised coverage, so changing this archived artifact would make its provenance less accurate.
- Evidence: Full-file inspection of `preflight-table.md` shows the exact qualification plus the fresh 13-test command and result; no production or evidence change is required.
