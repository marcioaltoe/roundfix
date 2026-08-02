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
line: 136
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ym0,comment:PRRC_kwDOS0qyts7cjgFU
review_hash: 1dc5229bc85be4f0b01633d434a963cf3bd395ee5d9e160651c3041b1b0031d5
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 022: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Compare captured retention-source digests with the digests recorded by shipped Source Baselines.
set -euo pipefail

# Inspect the Source Baseline identity schema for declared profile/catalog digests.
rg -nP -C 4 'type SourceBaselineIdentity struct' --type=go internal/baseline

# Show shipped baseline identity documents.
fd -t f 'baseline.json' | rg 'source-baselines' | while IFS= read -r f; do
  echo "== $f"
  cat "$f"
done

# Inspect the characterization fixture that names changed digests.
fd -t f 'same-baseline-changed-profile-and-catalog-digests.golden.json' --exec sed -n '1,80p' {}
```

Repository: marcioaltoe/roundfix

Length of output: 4146

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

echo '== source_contracts.go =='
sed -n '1,190p' internal/baseline/source_contracts.go

echo '== plan.go retention flow =='
sed -n '1040,1145p' internal/baseline/plan.go

echo '== retention symbols =='
rg -n -C 3 'captureCurrentRetentionSources|resolvePlanRetention|retentionSources|ProfileDigest|CatalogDigest' internal/baseline --glob '*.go'

echo '== characterization references =='
rg -n -C 5 'same-baseline-changed-profile-and-catalog-digests|existingProfileDigest|existingCatalogDigest' .
```

Repository: marcioaltoe/roundfix

Length of output: 32816

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

echo '== characterization test =='
sed -n '210,270p' internal/baseline/plan_characterization_test.go

echo '== retention tests =='
sed -n '2180,2410p' internal/baseline/plan_test.go

echo '== source-baseline digest construction =='
rg -n -C 5 'sourceBaselineIndexRecord|ManifestDigest|CorpusDigest|SourceBaselineIdentity|source-baselines/index' internal/baseline --glob '*.go'

echo '== manifest creation and catalog loading =='
sed -n '130,190p' internal/baseline/catalog.go
sed -n '2020,2060p' internal/baseline/plan.go
```

Repository: marcioaltoe/roundfix

Length of output: 32282

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

echo '== relevant diff =='
git diff -- internal/baseline/source_contracts.go internal/baseline/plan.go internal/baseline/plan_characterization_test.go

echo '== source-baseline assets =='
fd -t f . internal/baseline/assets/source-baselines --max-depth 4 | sort | while IFS= read -r f; do
  echo "== $f"
  case "$f" in
    *.json) sed -n '1,80p' "$f" ;;
  esac
done

echo '== deterministic key-space check =='
python3 - <<'PY'
from pathlib import Path
import json
import re

source = Path("internal/baseline/source_contracts.go").read_text()
plan = Path("internal/baseline/plan.go").read_text()
fixture = json.loads(Path(
    "internal/baseline/testdata/plan-characterization/"
    "same-baseline-changed-profile-and-catalog-digests.golden.json"
).read_text())

capture = re.search(
    r"tuple := BaselineSourceTuple\{(?P<body>.*?)\n\t\t\}",
    source,
    re.S,
)
lookup = re.search(
    r"tuple := BaselineSourceTuple\{(?P<body>.*?)\n\t\t\t\}",
    plan,
    re.S,
)

assert capture and lookup
assert "ProfileDigest: profile.Digest" in capture.group("body")
assert "CatalogDigest: c.Digest()" in capture.group("body")
assert "ProfileDigest: existing.ProfileDigest" in lookup.group("body")
assert "CatalogDigest: existing.CatalogDigest" in lookup.group("body")

existing = fixture["fixture"]
captured = {
    "Baseline": existing["existingBaseline"],
    "ProfileDigest": "<current profile.Digest>",
    "CatalogDigest": "<current catalog.Digest()>",
}
looked_up = {
    "Baseline": existing["existingBaseline"],
    "ProfileDigest": existing["existingProfileDigest"],
    "CatalogDigest": existing["existingCatalogDigest"],
}
assert captured != looked_up
print("capture and lookup use different digest sources for the changed-digest fixture")
print("captured:", captured)
print("lookup:  ", looked_up)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 11584

---



</details>

**Preserve retention evidence across digest changes.**

`captureCurrentRetentionSources` indexes sources with the current `profile.Digest` and `c.Digest()`, while `resolvePlanRetention` looks up the prior manifest digests. A changed digest therefore skips clause accounting. The source-baseline schema does not store historical profile or catalog digests, so use a stable baseline identity or persist the historical tuple. The characterization fixture changes only manifest digests and exits before this lookup because managed artifacts remain unchanged; it does not cover this path. Add a regression case with artifact drift and changed digests that asserts clause accounting.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/source_contracts.go` around lines 127 - 136, Update
captureCurrentRetentionSources and resolvePlanRetention so retention evidence
remains matchable when profile or catalog digests change, using a stable
baseline identity or persisting the historical tuple required by prior-manifest
lookups. Preserve duplicate detection and clause accounting for unchanged and
changed digests. Add a regression case with artifact drift and changed manifest
digests that verifies retention clauses are accounted for.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fbea1141380cd5a1c1cac1b2 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Retention sources are now indexed and resolved by immutable Source Baseline identity instead of current versus prior profile/catalog digests, preserving duplicate detection across digest changes. The regression recaptures evidence after both digests change, introduces managed-artifact drift, and proves clause accounting. `go test ./internal/baseline ./internal/cli` passed with a writable GOCACHE.
