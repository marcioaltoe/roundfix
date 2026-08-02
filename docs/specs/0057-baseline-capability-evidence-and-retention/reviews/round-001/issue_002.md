---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/catalog.go
line: 174
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0YmS,comment:PRRC_kwDOS0qyts7cjgEs
review_hash: 457f54cd98f1f156af7aae70667b23b4b7ffe6a16de0c17a494760220590cc62
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:29Z"
---

# Issue 002: _ Performance & Scalability_ _ Minor_ _ Quick win_

## Review Comment

_🚀 Performance & Scalability_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Check whether embedded catalog loading is cached, and size the source-baseline corpus.
set -euo pipefail

# Locate the embedded catalog entry point and any caching (sync.Once / package-level var).
rg -nP -C 6 'func LoadEmbeddedCatalog\s*\(' --type=go
rg -nP -C 3 '\b(sync\.Once|OnceValue|OnceValues)\b' --type=go internal/baseline

# Count indexed source baselines and corpus size to estimate per-load work.
fd -t f 'index.json' | rg 'source-baselines' | while IFS= read -r f; do
  echo "== $f"
  jq '.baselines | length' "$f" 2>/dev/null || cat "$f" | head -40
done
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
set -u

printf '%s\n' '== catalog entry points =='
rg -n -P -C 8 'func\s+LoadEmbeddedCatalog\s*\(' --glob '*.go' . || true

printf '%s\n' '== retention capture and SourceBaseline call sites =='
rg -n -P -C 8 'captureCurrentRetentionSources|SourceBaseline' --glob '*.go' . || true

printf '%s\n' '== catalog caching primitives =='
rg -n -P -C 4 '\b(sync\.Once|sync\.OnceValue|sync\.OnceValues|OnceValue|OnceValues)\b' --glob '*.go' internal . || true

printf '%s\n' '== source-baseline indexes =='
fd -t f 'index\.json$' . | while IFS= read -r f; do
  case "$f" in
    *source-baselines*)
      printf '== %s ==\n' "$f"
      python3 - "$f" <<'PY'
import json, pathlib, sys
p = pathlib.Path(sys.argv[1])
try:
    data = json.loads(p.read_text())
except Exception as e:
    print(f"parse error: {e}")
    raise SystemExit
for key in ("baselines", "entries", "sources"):
    value = data.get(key)
    if isinstance(value, list):
        print(f"{key}_length={len(value)}")
        break
else:
    print("top_level_keys=" + ",".join(map(str, data.keys())))
PY
      ;;
  esac
done
true
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
set -u

printf '%s\n' '== internal/baseline/catalog.go outline =='
ast-grep outline internal/baseline/catalog.go || true

printf '%s\n' '== LoadCatalog and retention implementation =='
rg -n -P '^(func|type)\s+(LoadCatalog|LoadEmbeddedCatalog|captureCurrentRetentionSources|SourceBaseline|Catalog|sourceBaseline)' internal/baseline/catalog.go
sed -n '90,230p' internal/baseline/catalog.go
sed -n '230,430p' internal/baseline/catalog.go

printf '%s\n' '== exact cache declarations in baseline package =='
rg -n -P '(^|[^[:alnum:]_])(sync\.Once|sync\.OnceValue|sync\.OnceValues|var\s+\w+\s*=\s*&?Catalog|var\s+\w+\s+\*Catalog)' internal/baseline --glob '*.go' || true

printf '%s\n' '== embedded source-baseline assets and index shape =='
fd -t f . internal/baseline | rg 'source-baselines|index\.json$' | head -80
find . -type f -path '*source-baselines*index.json' -print | while IFS= read -r f; do
  printf '== %s ==\n' "$f"
  wc -c "$f"
  python3 - "$f" <<'PY'
import json, pathlib, sys
p = pathlib.Path(sys.argv[1])
data = json.loads(p.read_text())
print("top_level_keys:", list(data))
for k, v in data.items():
    if isinstance(v, list):
        print(f"{k}_length:", len(v))
PY
done
```

Repository: marcioaltoe/roundfix

Length of output: 15712

---

🏁 Script executed:

```shell
set -u

printf '%s\n' '== retention capture and SourceBaseline definitions =='
rg -n -P '^func\s+\([^)]*\)\s+(captureCurrentRetentionSources|SourceBaseline)\b|^func\s+(captureCurrentRetentionSources|SourceBaseline)\b' internal/baseline --glob '*.go' || true

printf '%s\n' '== relevant source-baseline loader symbols =='
rg -n -P 'captureCurrentRetentionSources|sourceBaseline|SourceBaseline|source-baselines/index\.json|ReadoptionSourceBaseline' internal/baseline --glob '*.go' \
  | grep -v '_test.go' | head -160

printf '%s\n' '== source-baseline implementation file candidates =='
fd -t f -e go internal/baseline | while IFS= read -r f; do
  if rg -q 'captureCurrentRetentionSources|SourceBaseline' "$f"; then
    printf '%s\n' "$f"
  fi
done

printf '%s\n' '== asset counts and sizes =='
python3 - <<'PY'
from pathlib import Path
root = Path("internal/baseline/assets/source-baselines")
for p in sorted(root.rglob("*")):
    if p.is_file():
        print(f"{p.relative_to(root)}\t{p.stat().st_size}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 18308

---

🏁 Script executed:

```shell
set -u

printf '%s\n' '== source_contracts.go retention implementation =='
sed -n '90,325p' internal/baseline/source_contracts.go

printf '%s\n' '== catalog load validation sequence =='
sed -n '1,125p' internal/baseline/catalog_load.go
sed -n '1380,1515p' internal/baseline/catalog_validate.go

printf '%s\n' '== corpus totals =='
python3 - <<'PY'
from pathlib import Path
root = Path("internal/baseline/assets/source-baselines")
for name in ["corpus", "manifest.json", "baseline.json", "accounting.json"]:
    p = root / "baseline.standard-typescript-monorepo-0.0.1" / name
    if p.is_dir():
        files = [x for x in p.rglob("*") if x.is_file()]
        print(name, "files=", len(files), "bytes=", sum(x.stat().st_size for x in files))
    elif p.is_file():
        print(name, "bytes=", p.stat().st_size)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 17502

---



</details>

**Avoid eager retention-source loading during every catalog load.**

`LoadEmbeddedCatalog` has no cache. `validateSourceBaselines` already calls `SourceBaseline` for each indexed record, and `captureCurrentRetentionSources` repeats that work. This reparses metadata and recomputes the corpus digest on every load. Reuse the validated sources or defer capture until retention data is needed.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/catalog.go` around lines 172 - 174, Update
LoadEmbeddedCatalog to avoid calling captureCurrentRetentionSources after
validateSourceBaselines, which already resolves each indexed record through
SourceBaseline. Reuse the validated sources for retention data or defer
captureCurrentRetentionSources until retention information is actually needed,
while preserving the existing validation and error behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d83731d0c90bb425b4cc1b9d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Catalog validation now retains each verified Source Baseline and retention capture reuses those values instead of reparsing and rehashing them. The full Baseline package test passed.
