---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/profile_alignment.go
line: 260
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymo,comment:PRRC_kwDOS0qyts7cjgFE
review_hash: 34f8abf962cf6a0f2e394c9a0271544e7551b76cf44ef3069d3aea0755ce0761
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 015: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Inspect PostgreSQLEvidence construction and nil-slice handling.
set -euo pipefail

rg -n -C6 'type PostgreSQLEvidence struct' internal/baseline
rg -n -C4 'AcceptedContractPaths' internal/baseline
rg -n -C4 'func stringsFromAny' internal/baseline
```

Repository: marcioaltoe/roundfix

Length of output: 4328

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- stringsFromAny ---'
sed -n '1590,1635p' internal/baseline/profile_alignment.go

printf '%s\n' '--- PostgreSQL evaluation ---'
sed -n '1040,1145p' internal/baseline/profile_alignment.go

printf '%s\n' '--- RecheckCapabilities ---'
rg -n -A120 -B20 '^func RecheckCapabilities|^func \(.*\) RecheckCapabilities' internal/baseline/profile_alignment.go

printf '%s\n' '--- CLI failure construction and serialization ---'
rg -n -A100 -B20 'printBaselineCapabilitiesCheckFailure|CapabilityRecheckResult|capabilities.check' internal/cli/baseline_profile.go
```

Repository: marcioaltoe/roundfix

Length of output: 19100

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- capability dispatch and PostgreSQL initialization ---'
sed -n '560,700p' internal/baseline/profile_alignment.go
rg -n -A90 -B20 'func evaluateRepositoryCapabilities|evaluatePostgreSQLCapability|capability.stack.postgresql' internal/baseline/profile_alignment.go

printf '%s\n' '--- implementation collector ---'
sed -n '1120,1215p' internal/baseline/profile_alignment.go

printf '%s\n' '--- profile and catalog PostgreSQL declarations ---'
rg -n -C8 'capability.stack.postgresql|postgresql' internal/baseline --glob '*.go' --glob '*.json' --glob '*.yaml' --glob '*.yml'

printf '%s\n' '--- focused tests for empty and nil PostgreSQL evidence ---'
rg -n -C8 'PostgreSQL|AcceptedContractPaths|Implementation|RecheckCapabilities' internal/baseline/*test.go
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- CapabilityRecheck tests ---'
rg -n -C12 'CapabilityRecheck|capabilities.check|CapabilityRecheckResult|acceptedContractPaths|postgresql' internal --glob '*_test.go' \
  | head -n 300

printf '%s\n' '--- schema and output assertions ---'
rg -n -C8 'schemaVersion|acceptedContractPaths|implementation.*null|postgresql.*null' internal/cli internal/baseline --glob '*_test.go' --glob '*.json' \
  | head -n 300

printf '%s\n' '--- focused source invariant and JSON probe ---'
python3 - <<'PY'
import json
from pathlib import Path

source = Path("internal/baseline/profile_alignment.go").read_text()
dispatch = source[source.index("func evaluateRepositoryCapabilities"):source.index("func collectCapabilityEvidence")]
recheck = source[source.index("func RecheckCapabilities"):source.index("func resolveCapabilityRecheckProfile")]

assert "postgres := PostgreSQLEvidence{}" in dispatch
assert 'if capability.ID == "capability.stack.postgresql"' in dispatch
assert "PostgreSQL:   evaluation.postgres" in recheck

nil_shape = {
    "acceptedContractPaths": None,
    "implementation": None,
}
empty_shape = {
    "acceptedContractPaths": [],
    "implementation": [],
}
print("success zero-value shape:", json.dumps(nil_shape, separators=(",", ":")))
print("failure explicit-empty shape:", json.dumps(empty_shape, separators=(",", ":")))
assert json.dumps(nil_shape)["acceptedContractPaths"] if False else True
assert json.dumps(nil_shape) != json.dumps(empty_shape)
print("distinct JSON shapes: yes")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- CapabilityRecheck tests ---'
rg -n -C12 'CapabilityRecheck|capabilities.check|CapabilityRecheckResult|acceptedContractPaths|postgresql' internal --glob '*_test.go' \
  | head -n 300

printf '%s\n' '--- schema and output assertions ---'
rg -n -C8 'schemaVersion|acceptedContractPaths|implementation.*null|postgresql.*null' internal/cli internal/baseline --glob '*_test.go' --glob '*.json' \
  | head -n 300

printf '%s\n' '--- focused source invariant and JSON probe ---'
python3 - <<'PY'
import json
from pathlib import Path

source = Path("internal/baseline/profile_alignment.go").read_text()
dispatch = source[source.index("func evaluateRepositoryCapabilities"):source.index("func collectCapabilityEvidence")]
recheck = source[source.index("func RecheckCapabilities"):source.index("func resolveCapabilityRecheckProfile")]

assert "postgres := PostgreSQLEvidence{}" in dispatch
assert 'if capability.ID == "capability.stack.postgresql"' in dispatch
assert "PostgreSQL:   evaluation.postgres" in recheck

nil_shape = {
    "acceptedContractPaths": None,
    "implementation": None,
}
empty_shape = {
    "acceptedContractPaths": [],
    "implementation": [],
}
print("success zero-value shape:", json.dumps(nil_shape, separators=(",", ":")))
print("failure explicit-empty shape:", json.dumps(empty_shape, separators=(",", ":")))
assert json.dumps(nil_shape) != json.dumps(empty_shape)
print("distinct JSON shapes: yes")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---



</details>

**Normalize empty PostgreSQL slices in capability-check JSON.**

When a profile omits `capability.stack.postgresql`, `RecheckCapabilities` returns nil slices, which encode as `null`. The CLI failure response encodes the same fields as `[]`. Normalize both slices before returning the success result.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/profile_alignment.go` around lines 246 - 260, Normalize the
PostgreSQL evidence slices in RecheckCapabilities before returning the success
CapabilityRecheckResult, converting nil slices to empty slices so both JSON
fields encode as []. Apply this to both PostgreSQL slice fields while preserving
existing non-nil values and the CLI failure response format.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:548ac73b549aa92cc8a1f9eb -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `RecheckCapabilities` now normalizes omitted PostgreSQL contract and implementation slices to empty arrays, with JSON regression coverage. The full Baseline package test passed.
