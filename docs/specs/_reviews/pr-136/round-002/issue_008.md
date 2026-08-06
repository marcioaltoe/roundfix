---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md
line: 99
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5Yh,comment:PRRC_kwDOS0qyts7eY0jP
review_hash: 29248527e95094ac89faf22d186833510b623f88fc711727ae0917f62e094561
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 008: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

rg -n -C 6 \
  '\b(JournalSink|Flush|Close|FinalizeRun|appendRunEvent|AppendRunEvents)\b' \
  --glob '*.go' .
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- target spec ---'
sed -n '1,240p' docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md

printf '%s\n' '--- exact production references ---'
rg -n \
  '\b(JournalSink|FinalizeRun|appendRunEvent|AppendRunEvents)\b' \
  --glob '*.go' \
  --glob '!**/*_test.go' . || true

printf '%s\n' '--- lifecycle and event persistence references ---'
rg -n -C 4 \
  '\b(CompleteRun|StateClean|StateStopped|Append.*Event|RunEvent|event.*journal|journal.*event|cancel|teardown|shutdown)\b' \
  --glob '*.go' \
  --glob '!**/*_test.go' internal cmd . 2>/dev/null | head -n 1200
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- current journal append and sink implementation ---'
sed -n '540,620p' internal/store/journal.go
sed -n '1040,1100p' internal/store/journal.go
sed -n '1210,1265p' internal/store/journal.go

printf '%s\n' '--- JournalSink construction and lifetime ---'
rg -n -C 10 \
  'JournalSink\{|newReviewRequester|NewFanout|runUI|ui\.Close|CompleteRun\(' \
  internal/cli internal/watch internal/agent internal/store \
  --glob '*.go' --glob '!**/*_test.go' | head -n 1600

printf '%s\n' '--- relevant lifecycle functions ---'
sed -n '2360,2450p' internal/cli/cli.go
sed -n '2640,2740p' internal/cli/cli.go
sed -n '3230,3300p' internal/cli/cli.go
sed -n '3460,3550p' internal/cli/cli.go
sed -n '4660,4710p' internal/cli/cli.go
sed -n '4790,4845p' internal/cli/cli.go
sed -n '1,90p' internal/cli/runui.go
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- terminal completion and Run settlement ---'
sed -n '4700,4855p' internal/cli/cli.go
sed -n '650,825p' internal/store/store.go

printf '%s\n' '--- fanout close and sink interface ---'
rg -n -C 12 \
  'type Sink|type Fanout|func \(.*Fanout.*Close|func \(.*Fanout.*Publish|func \(.*Fanout.*Flush' \
  internal/runevent --glob '*.go'
sed -n '100,145p' internal/cli/runui.go

printf '%s\n' '--- all production JournalSink values and terminal paths ---'
rg -n \
  'JournalSink\{|CompleteRun\(|ui\.Close\(\)|defer ui\.Close' \
  internal --glob '*.go' --glob '!**/*_test.go'
```

Repository: marcioaltoe/roundfix

Length of output: 22299

---



</details>

**Define a shared flush lifecycle before accepting batched writes.**

`runUI.Close` only closes the fanout; it does not flush the critical sink. Terminal completion and Agent teardown create new `JournalSink` values, so an immediate terminal publish cannot flush a pending batch held by an earlier sink unless buffering is shared per `Store`. Add an explicit flush and close contract, invoke it before terminal settlement and teardown, and propagate errors.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md` around lines
84 - 99, Define a Store-scoped JournalSink lifecycle so buffered events are
shared across terminal completion and Agent teardown rather than creating
isolated sinks. Add explicit Flush and Close operations, invoke them before
terminal settlement and teardown, and propagate any flush or close errors
through the existing Run failure path.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a67dc92963485fc7413e557b -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Spec 0081 now makes the JournalWriter Store-scoped, so every sink
  handle shares one batch. It defines explicit Flush and Close operations,
  flushes before Agent teardown, closes before terminal settlement, commits
  terminal state and `daemon.outcome` atomically, and routes Flush or Close
  errors through the existing Run-failure path. Post-terminal receipts bypass
  the closed buffer with immediate transactions.
- Focused evidence: independent source inspection confirmed current callers
  construct separate `JournalSink{Store: runStore}` values and `runUI.Close`
  only closes the fanout, establishing the reported lifecycle gap; `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go run -buildvcs=false
  ./cmd/roundfix spec check` passed with no findings for Spec 0081.
- Daemon Verification: `make verify` not run; Daemon-owned.
