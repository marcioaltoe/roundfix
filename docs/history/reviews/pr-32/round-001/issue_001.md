---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/agent/acpx_runner.go
line: 1079
severity: review
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5sG,comment:PRRC_kwDOS0qyts7Wt94S
review_hash: 9d8693ae5c4576bd718dfedab8753d217941d2d6e5e7a20883f22fab5966182c
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-32/round-002/issue_001.md
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:30Z"
---


# Issue 001: _ Stability & Availability_ _ Critical_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🔴 Critical_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

if command -v rtk >/dev/null 2>&1; then
  alias_cmd='rtk'
else
  alias_cmd=''
fi

run() {
  if [[ -n "$alias_cmd" ]]; then
    $alias_cmd "$@"
  else
    "$@"
  fi
}

run git ls-files internal/agent/acpx_runner.go internal/agent/agent.go internal/daemon/agent_session_owner.go internal/daemon/task_engine.go
echo '--- outline acpx_runner.go ---'
run ast-grep outline internal/agent/acpx_runner.go --view expanded
echo '--- outline agent.go ---'
run ast-grep outline internal/agent/agent.go --view expanded
echo '--- relevant ranges acpx_runner.go ---'
run sed -n '1040,1105p' internal/agent/acpx_runner.go
echo '--- relevant ranges agent.go ---'
run sed -n '320,390p' internal/agent/agent.go
echo '--- search stateMu / acpx initialization ---'
run rg -n 'stateMu|acpxRunner\(|ACPXRunner|DefaultRunner' internal/agent internal/daemon
```

Repository: marcioaltoe/roundfix

Length of output: 29580

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

if command -v rtk >/dev/null 2>&1; then
  run_cmd() { rtk "$@"; }
else
  run_cmd() { "$@"; }
fi

run_cmd git ls-files internal/agent/acpx_runner.go internal/agent/agent.go internal/daemon/agent_session_owner.go internal/daemon/task_engine.go
echo '--- outline acpx_runner.go ---'
run_cmd ast-grep outline internal/agent/acpx_runner.go --view expanded
echo '--- outline agent.go ---'
run_cmd ast-grep outline internal/agent/agent.go --view expanded
echo '--- relevant ranges acpx_runner.go ---'
sed -n '1040,1105p' internal/agent/acpx_runner.go
echo '--- relevant ranges agent.go ---'
sed -n '320,390p' internal/agent/agent.go
echo '--- search stateMu / acpx initialization ---'
run_cmd rg -n 'stateMu|acpxRunner\(|ACPXRunner|DefaultRunner' internal/agent internal/daemon
```

Repository: marcioaltoe/roundfix

Length of output: 29580

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail
pwd
git ls-files internal/agent/acpx_runner.go internal/agent/agent.go
```

Repository: marcioaltoe/roundfix

Length of output: 231

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail
echo hi
```

Repository: marcioaltoe/roundfix

Length of output: 161

---



</details>

**Make `DefaultRunner.acpxRunner()` initialize shared state once.** The nil checks and `Now` assignment are unsynchronized, so concurrent first use can publish different `*ACPXRunner` or `*sync.Mutex` values and race on the shared runner.

- `internal/agent/acpx_runner.go#L1074-L1079`: replace the lazy `stateMu` pointer with a zero-value `sync.Mutex`, or initialize it in the constructor.
- `internal/agent/agent.go#L360-L366`: avoid mutating `runner.acpx` and `runner.acpx.Now` on every call; construct the runner once and reuse it.

<details>
<summary>📍 Affects 2 files</summary>

- `internal/agent/acpx_runner.go#L1074-L1079` (this comment)
- `internal/agent/agent.go#L360-L366`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner.go` around lines 1074 - 1079, Make shared ACPX
runner initialization safe for concurrent first use: in
internal/agent/acpx_runner.go:1074-1079, replace the lazily initialized stateMu
pointer with a zero-value sync.Mutex or initialize it in the ACPXRunner
constructor, and update lockState to use it directly. In
internal/agent/agent.go:360-366, change DefaultRunner.acpxRunner() to construct
and publish runner.acpx only once, reusing the existing runner without mutating
runner.acpx or runner.acpx.Now on each call.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/agent/acpx_runner.go</file>
<line_range>1074-1079</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/agent/agent.go</file>
<line_range>360-366</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:daf53b221d182734a7ae90d4 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
