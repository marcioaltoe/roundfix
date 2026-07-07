---
source: coderabbit
pr: "18"
round: 2
round_created_at: "2026-07-07T14:05:56Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 4237143afdd7097e755e14b962156aaf6c6e6654
file: internal/cli/cli.go
line: 1106
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6O6PJ7,comment:PRRC_kwDOS0qyts7SyYRu
review_hash: 933c08eddf73d39b5ba9a49da734ebc4114a9fde8271409959bc09c78bf1e842
duplicate_of: ""
source_review_id: "4645087962"
source_review_submitted_at: "2026-07-07T12:31:07Z"
---

# Issue 003: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
rg -n 'GOOS' --iglob '*.yml' --iglob '*.yaml' -g '!vendor/**'
cat go.mod | head -5
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '\n== go.mod ==\n'
sed -n '1,20p' go.mod

printf '\n== workflows mentioning GOOS/Windows ==\n'
rg -n --hidden --glob '.github/workflows/**' 'GOOS|windows|Windows|GOARCH|matrix|runs-on' .

printf '\n== config default specs root ==\n'
rg -n --hidden 'defaultSpecsRoot|Specs.Root|docs/specs|filepath.Join\("docs", "specs"\)' .

printf '\n== relevant CLI code ==\n'
sed -n '1088,1110p' internal/cli/cli.go
sed -n '1,140p' internal/roundconfig/config.go
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '\n== release workflow matrix ==\n'
sed -n '1,140p' .github/workflows/release.yml

printf '\n== ci workflow ==\n'
sed -n '1,140p' .github/workflows/ci-conventions.yml

printf '\n== docs mentioning supported platforms ==\n'
rg -n --hidden 'Windows|windows|GOOS|cross-compile|cross compile|supported target|platform' README.md docs internal .github/workflows -g '!**/*.golden'
```

Repository: marcioaltoe/roundfix

Length of output: 30089

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '\n== reviewArtifactSpecsRoot definition and uses ==\n'
rg -n -A4 -B4 'func reviewArtifactSpecsRoot|reviewArtifactSpecsRoot\(' internal/cli/cli.go

printf '\n== ResolveSpecsRoot and default root ==\n'
sed -n '620,700p' internal/config/config.go

printf '\n== nearby cli logic ==\n'
sed -n '1088,1110p' internal/cli/cli.go
sed -n '2388,2406p' internal/cli/cli.go
```

Repository: marcioaltoe/roundfix

Length of output: 5413

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '\n== resolveReviewArtifactRoot ==\n'
sed -n '2402,2448p' internal/cli/cli.go

printf '\n== specs root default helpers ==\n'
sed -n '1,80p' internal/config/config.go

printf '\n== repo docs on default spec root ==\n'
sed -n '520,545p' README.md
```

Repository: marcioaltoe/roundfix

Length of output: 5071

---



</details>

**Compare against the built-in `docs/specs` string here.** `filepath.Join("docs", "specs")` becomes `docs\specs` on Windows, so the default config never takes the fast-path and always falls through to `ResolveSpecsRoot`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 1096 - 1104, The specs-root fast-path in
the CLI is comparing against a path-joined value that changes by platform, so
the default `docs/specs` case is missed on Windows. Update the check in the
logic around `reviewArtifactSpecsRoot` and `roundconfig.ResolveSpecsRoot` to
compare `loadedConfig.Config.Specs.Root` against the built-in `docs/specs`
string directly, and keep the existing fallback resolution only when the
configured root differs from that exact default.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:584cd6ee7bb50b67ab019ab4 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Review artifact Spec Root fast-path detection now compares the config value to the literal built-in `docs/specs` string, with coverage for Windows-style `docs\specs` as non-default config.
  - Verification: configured command `make verify` was run as `rtk make verify` and passed: Go tests, skills check, and build completed.
