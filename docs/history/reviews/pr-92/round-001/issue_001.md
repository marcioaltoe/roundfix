---
source: coderabbit
pr: "92"
round: 1
round_created_at: "2026-08-03T19:34:00Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spawn-economy
head_sha: 7765d1f6d62e59ebf68ca2e4e2e273733da58425
file: internal/agent/acpx_runner.go
line: 669
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WGXKv,comment:PRRC_kwDOS0qyts7c9hS0
review_hash: 6f62d890de2dd34d62b698a478c01a9b2317252d155cc0a0dcfe941daa2ddb22
duplicate_of: ""
source_review_id: "4847882119"
source_review_submitted_at: "2026-08-03T19:33:10Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- file outline ---'
ast-grep outline internal/agent/acpx_runner.go --match 'func acpxConfigPath($_)' --view expanded || true

printf '%s\n' '--- target implementation ---'
sed -n '620,690p' internal/agent/acpx_runner.go

printf '%s\n' '--- references and related tests ---'
rg -n -C 4 'acpxConfigPath|environmentValue|USERPROFILE|HOME|config\.json' internal/agent --glob '*.go'

printf '%s\n' '--- module and working tree ---'
sed -n '1,80p' go.mod
git status --short
git diff --stat
```

Repository: marcioaltoe/roundfix

Length of output: 13244

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- target implementation ---'
sed -n '620,690p' internal/agent/acpx_runner.go

printf '%s\n' '--- references and tests ---'
rg -n -C 4 'acpxConfigPath|environmentValue|USERPROFILE|HOME|config\.json' internal/agent --glob '*.go'

printf '%s\n' '--- module and tree ---'
sed -n '1,80p' go.mod
git status --short
git diff --stat
```

Repository: marcioaltoe/roundfix

Length of output: 13163

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

python3 - <<'PY'
from pathlib import Path
p = Path("internal/agent/acpx_runner.go")
s = p.read_text()
start = s.index("func acpxConfigPath(")
end = s.index("\n}", start) + 2
print(s[start:end])
print("uses os.UserHomeDir:", "os.UserHomeDir()" in s[start:end])
print("uses USERPROFILE:", "USERPROFILE" in s[start:end])
PY

printf '%s\n' '--- call sites ---'
rg -n -C 8 'acpxConfigPath\(' . --glob '*.go'

printf '%s\n' '--- environment construction around callers ---'
rg -n -C 12 'environmentValue\(|environment :=|append\(.*HOME|HOME=' internal/agent --glob '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 28755

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- adapter resolution call path ---'
rg -n -C 18 'func resolveAdapterCommandWithEnv|configuredAdapterInvocation\(' internal/agent/acpx_runner.go internal/agent --glob '*.go'

printf '%s\n' '--- environment helpers and platform handling ---'
sed -n '130,180p' internal/agent/codex_spawn.go
rg -n -C 8 'goos|runtime\.GOOS|USERPROFILE|HOMEDRIVE|HOMEPATH' internal/agent --glob '*.go'

printf '%s\n' '--- standard-library contract ---'
go doc os.UserHomeDir
go env GOOS GOARCH
```

Repository: marcioaltoe/roundfix

Length of output: 29296

---



</details>

**Do not use the process home for an explicit environment**

`resolveAdapterCommandWithEnv` passes its supplied environment to `acpxConfigPath`, but `os.UserHomeDir()` reads the process environment. When `HOME` is absent, this can select the host config instead of the supplied environment's config. On Windows, it also ignores a supplied `USERPROFILE`. Resolve both variables from `environment`, or return an error when neither exists. Add regression tests for these cases.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner.go` around lines 659 - 669, Update acpxConfigPath
to resolve the home directory exclusively from the supplied environment,
honoring HOME and Windows USERPROFILE without falling back to os.UserHomeDir;
return an error when neither variable is present. Add regression tests covering
an absent HOME with a supplied environment and USERPROFILE-based resolution.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5cfedfc41dd590609cecb49a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `acpxConfigPath` read `HOME` from the supplied environment but fell back to `os.UserHomeDir`, which could select the host process config and ignored an explicit `USERPROFILE`.

## Result

- Changed `acpxConfigPath` to resolve `HOME`, then `USERPROFILE`, exclusively from the supplied environment and return an error when neither exists.
- Added a regression test for `HOME` precedence, `USERPROFILE` resolution, and a supplied environment with no home variable.
- Focused check: `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache-run_20260803T193424Z-agent go test ./internal/agent -run '^TestACPXConfigPathUsesOnlyTheExplicitEnvironment$' -count=1` — passed.
- Package check: `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache-run_20260803T193424Z-agent go test ./internal/agent -count=1` — passed.
- The Daemon owns the configured `make verify` run.
