---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T15:34:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: a12c1a665c5970773e04c4a145c6b9b0c5a0e686
file: internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
line: 33
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WBeNK,comment:PRRC_kwDOS0qyts7c2Vzh
review_hash: 57aeed177ef58973f2c840bb334831b5cc8ebb14467660e293069da0696071e0
duplicate_of: ""
source_review_id: "4845660382"
source_review_submitted_at: "2026-08-03T15:14:34Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
go test ./internal/baseline -run '^TestAssetsSyncCompatibilityMatchesMaintainedPythonContract$'
```

Repository: marcioaltoe/roundfix

Length of output: 197

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- fixture diff ---'
git diff -- internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json | sed -n '1,240p'

printf '%s\n' '--- relevant fixture structure ---'
rg -n -C 3 '"digest"|"afterIdentity"|"identity"|"go-cli"|"rust-cli"|"typescript-bun"' \
  internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json | sed -n '1,260p'

printf '%s\n' '--- test and implementation references ---'
rg -n -C 4 'AssetsSyncCompatibilityMatchesMaintainedPythonContract|afterIdentity|postState|asset-sync' \
  internal/baseline --glob '*.go' --glob '*.py' | sed -n '1,320p'
```

Repository: marcioaltoe/roundfix

Length of output: 17794

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import hashlib
import json
from pathlib import Path

p = Path("internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json")
data = json.loads(p.read_text())

for key in ("plannedByteSequence", "postState"):
    print(key)
    for item in data.get(key, []):
        print(item)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 1529

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- compatibility test and regeneration logic ---'
sed -n '279,390p' internal/baseline/assets_sync_test.go
sed -n '218,390p' internal/baseline/compatibility_corpus_test.go

printf '%s\n' '--- fixture output sections ---'
rg -n '"(normalizedOutput|plannedByteSequence|postState|repositoryPreimage|managedEntryLedger)"' \
  internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
sed -n '1,120p' internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
sed -n '2025,2110p' internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json

printf '%s\n' '--- hashing and snapshot generation symbols ---'
ast-grep outline internal/baseline/assets_sync.go
rg -n -C 5 'canonicalSHA256|setup-snapshot|normalizedOutput|afterIdentity|postState|sha256' \
  internal/baseline --glob '*.go' | sed -n '1,360p'
```

Repository: marcioaltoe/roundfix

Length of output: 42065

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import hashlib
import json
from pathlib import Path

p = Path("internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json")
fixture = json.loads(p.read_text())

def hashes(value):
    variants = {
        "compact": json.dumps(value, separators=(",", ":"), ensure_ascii=False).encode(),
        "pretty": (json.dumps(value, indent=2, ensure_ascii=False) + "\n").encode(),
    }
    return {name: "sha256:" + hashlib.sha256(raw).hexdigest()
            for name, raw in variants.items()}

for setup in fixture["manifest"]["setups"]:
    print(setup["id"], setup["digest"])
    print(hashes(setup))
PY
```

Repository: marcioaltoe/roundfix

Length of output: 900

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- generated setup files ---'
fd -i -t f 'go-cli\.json|rust-cli\.json|typescript-bun\.json' .
while IFS= read -r file; do
  [ -z "$file" ] && continue
  printf '%s ' "$file"
  wc -c < "$file"
  sha256sum "$file"
done < <(fd -i -t f 'go-cli\.json|rust-cli\.json|typescript-bun\.json' .)

printf '%s\n' '--- snapshot struct ---'
sed -n '80,125p' internal/baseline/assets_sync.go

printf '%s\n' '--- exact identity consumers ---'
rg -n -C 8 'plannedByteSequence|postState|afterIdentity|repositoryPreimage' \
  internal/baseline --glob '*.go' | sed -n '1,360p'
```

Repository: marcioaltoe/roundfix

Length of output: 595

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import hashlib
import json
from pathlib import Path

fixture = json.loads(
    Path("internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json")
    .read_text()
)

for setup in fixture["manifest"]["setups"]:
    ordered = {
        "schemaVersion": setup.get("schemaVersion"),
        "id": setup.get("id"),
        "version": setup.get("version"),
        "source": setup.get("source"),
        "digest": setup.get("digest"),
        "skills": setup.get("skills"),
        "activationBundles": setup.get("activationBundles"),
    }
    ordered = {k: v for k, v in ordered.items() if v is not None}
    raw = (json.dumps(ordered, indent=2, ensure_ascii=False) + "\n").encode()
    print(setup["id"], len(raw), "sha256:" + hashlib.sha256(raw).hexdigest())
PY
```

Repository: marcioaltoe/roundfix

Length of output: 423

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

for file in \
  internal/baseline/assets/setups/go-cli.json \
  internal/baseline/assets/setups/rust-cli.json \
  internal/baseline/assets/setups/typescript-bun.json
do
  printf '%s\n' "--- $file ---"
  wc -c "$file"
  sha256sum "$file"
  sed -n '1,18p' "$file"
done
```

Repository: marcioaltoe/roundfix

Length of output: 308

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import hashlib
from pathlib import Path

for name in ("go-cli.json", "rust-cli.json", "typescript-bun.json"):
    p = Path("internal/baseline/assets/setups") / name
    data = p.read_bytes()
    print(f"{p}: size={len(data)} sha256:sha256:{hashlib.sha256(data).hexdigest()}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 573

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import hashlib
from pathlib import Path

for name in ("go-cli.json", "rust-cli.json", "typescript-bun.json"):
    path = Path("internal/baseline/assets/setups") / name
    data = path.read_bytes()
    print(f"{path}: size={len(data)} sha256:{hashlib.sha256(data).hexdigest()}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 552

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import json
from pathlib import Path

fixture = json.loads(
    Path("internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json").read_text()
)

for setup in fixture["manifest"]["setups"]:
    name = setup["id"] + ".json"
    generated = json.loads(
        Path("internal/baseline/assets/setups", name).read_text()
    )
    print(f"--- {name} ---")
    print("fixture keys:", list(setup))
    print("asset keys:  ", list(generated))
    print("equal:", setup == generated)
    for key in sorted(set(setup) | set(generated)):
        if setup.get(key) != generated.get(key):
            print("different:", key)
            if key in setup:
                print("  fixture:", repr(setup[key])[:500])
            if key in generated:
                print("  asset:  ", repr(generated[key])[:500])
PY
```

Repository: marcioaltoe/roundfix

Length of output: 5453

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

sed -n '85,125p' internal/baseline/assets_sync.go
```

Repository: marcioaltoe/roundfix

Length of output: 1408

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
import hashlib
import json
from collections import OrderedDict
from pathlib import Path

fixture = json.loads(
    Path("internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json").read_text()
)

def omit_empty(value):
    if isinstance(value, dict):
        return OrderedDict(
            (k, omit_empty(v))
            for k, v in value.items()
            if v not in ("", None, [], {})
        )
    if isinstance(value, list):
        return [omit_empty(v) for v in value]
    return value

for setup in fixture["manifest"]["setups"]:
    ordered = OrderedDict()
    for key in ("schemaVersion", "id", "version", "source", "digest", "skills", "activationBundles"):
        if key in setup:
            ordered[key] = setup[key]
    ordered = omit_empty(ordered)
    raw = (json.dumps(ordered, indent=2, ensure_ascii=False) + "\n").encode()
    print(
        setup["id"],
        f"size={len(raw)}",
        f"sha256:{hashlib.sha256(raw).hexdigest()}",
    )
PY
```

Repository: marcioaltoe/roundfix

Length of output: 438

---



</details>

**Refresh the fixture output identities and sizes.**

The updated manifest serializes to:

- `go-cli.json`: 14099 bytes, `sha256:994a789533435eb45b79999403cbcc889f5f9f937d8a5792a711420013821eff`
- `rust-cli.json`: 12388 bytes, `sha256:a89cc5015fc926c81b5774255b67944bb6f361d33bb7fbe08d819ad023971517`
- `typescript-bun.json`: 40301 bytes, `sha256:37592a49c0e4d7296078986bb82e141eb6cfbcb6e04b343a11a797d923b9962b`

Copy each identity to its matching `plannedByteSequence[*].afterIdentity` and `postState[*].identity`, update the three `postState[*].size` values, and rerun `go test ./internal/baseline -run '^TestAssetsSyncCompatibilityMatchesMaintainedPythonContract$'`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` around
lines 25 - 33, Refresh the asset-sync fixture identities and sizes for go-cli,
rust-cli, and typescript-bun: update each matching
plannedByteSequence[*].afterIdentity and postState[*].identity with the
specified SHA-256, and set postState[*].size to the specified byte count. Then
run the targeted TestAssetsSyncCompatibilityMatchesMaintainedPythonContract test
to verify the fixture.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5bcb204155bd788c78903f46 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Refreshed every planned/post-state output identity and size, then ran the sanctioned Baseline digest regeneration to include the branch's current Roundfix-skill digest and its deterministic catalog/characterization fallout. The final TypeScript identity is `sha256:33ab52a83944dab5d274f9aedfdbe34d26890b452d7647bf4a2fe0848eb390cb` because the regenerator also brought the stale Roundfix skill entry current. `make baseline-digests` then reported no changes, and `go test ./internal/baseline -count=1` passed with 489 tests.
