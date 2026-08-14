---
source: coderabbit
pr: "156"
round: 2
round_created_at: "2026-08-11T19:18:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-journal-cheap-to-write-and-keep
head_sha: 5d670eebf70c3a291890416717848f2df5b2ce0d
file: internal/cli/runui.go
line: 141
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YWRA7,comment:PRRC_kwDOS0qyts7gLDvL
review_hash: 2e40caaecdc37637f5b3481bb2d184d862096f099e17cfba716b4d392e697bf0
duplicate_of: ""
source_review_id: "4909904494"
source_review_submitted_at: "2026-08-11T19:17:47Z"
---

# Issue 002: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

ast-grep outline internal/cli/runui.go --match runUI --view expanded
rg -n -C3 --type=go '\b(?:ui|runUI)\.Close\s*\(\)' internal
rg -n -C3 --type=go 'func \(ui \*runUI\) Close\s*\(' internal/cli/runui.go
```

Repository: marcioaltoe/roundfix

Length of output: 9746

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- runUI lifecycle and Close implementation ---'
cat -n internal/cli/runui.go | sed -n '1,155p'

printf '%s\n' '--- startRunUI call sites and surrounding contexts ---'
rg -n -C5 --type=go 'startRunUI\(' internal/cli

printf '%s\n' '--- context cancellation and teardown context usage ---'
rg -n -C4 --type=go 'context\.(WithCancel|WithTimeout|WithoutCancel|Background)\(' internal/cli/runui.go internal/cli/implement.go internal/cli/cli.go

printf '%s\n' '--- writer and FlushJournal definitions/callers ---'
rg -n -C5 --type=go 'func .*FlushJournal|FlushJournal\(' internal
```

Repository: marcioaltoe/roundfix

Length of output: 40914

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- FlushJournal implementation and context usage ---'
cat -n internal/store/journal.go | sed -n '1280,1435p'
rg -n -C4 --type=go 'func .*flush|ctx\.Value|Value\(' internal/store

printf '%s\n' '--- tracing and request metadata definitions ---'
rg -n -C4 --type=go 'trace|span|request.?id|correlation|WithValue|slog.*Context|Context' internal/store internal/cli | head -n 500

printf '%s\n' '--- runUI tests and Close-related assertions ---'
rg -n -C6 --type=go 'runUI|Close\(\)|FlushJournal|JournalShutdownTimeout' internal/cli/*_test.go internal/store/*_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- journal writer flush and commit path ---'
cat -n internal/store/journal_writer.go | sed -n '150,325p'

printf '%s\n' '--- context-bearing call-site declarations ---'
for file in internal/cli/implement.go internal/cli/cli.go internal/cli/cli_test.go; do
  printf '%s\n' "--- $file ---"
  rg -n -B12 -A3 --type=go 'ui\.Close\(\)' "$file"
done

printf '%s\n' '--- all context value and tracing references, excluding tests ---'
rg -n --type=go --glob '!**/*_test.go' 'ctx\.Value|context\.WithValue|traceparent|trace_id|span_id|request_id|correlation_id|requestId|traceID|spanID' internal || true
```

Repository: marcioaltoe/roundfix

Length of output: 19233

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

cat >/tmp/verify_context_without_cancel.go <<'EOF'
package main

import (
	"context"
	"fmt"
	"time"
)

type key struct{}

func main() {
	parent := context.WithValue(context.Background(), key{}, "request-metadata")
	parent, cancel := context.WithCancel(parent)
	detached := context.WithoutCancel(parent)
	bounded, boundedCancel := context.WithTimeout(detached, time.Minute)
	defer cancel()
	defer boundedCancel()

	cancel()
	_, hasDeadline := detached.Deadline()
	fmt.Printf("value_preserved=%t detached_done_nil=%t bounded_deadline=%t\n",
		detached.Value(key{}) == "request-metadata",
		detached.Done() == nil,
		hasDeadline == false && bounded.Err() == nil,
	)
}
EOF

go run /tmp/verify_context_without_cancel.go
```

Repository: marcioaltoe/roundfix

Length of output: 224

---



</details>

**Preserve the caller context for the teardown flush.**

Change `runUI.Close` to accept `ctx` and derive the timeout context from `context.WithoutCancel(ctx)`. This keeps teardown independent of cancellation while preserving caller context values and the existing shutdown deadline.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/runui.go` around lines 131 - 141, Update runUI.Close to accept a
ctx parameter and derive its journal flush timeout context from
context.WithoutCancel(ctx) instead of context.Background(). Propagate the caller
context through each Close invocation while preserving the existing
JournalShutdownTimeout deadline and FlushJournal behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:d2bda5d50e493a77369d7fe5 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:

## Resolution

- Status: `resolved`
- Date: 2026-08-11

### Finding
`runUI.Close` derived its journal flush timeout from
`context.WithoutCancel(context.Background())`, which preserved the shutdown
deadline but would have been independent of caller context values. The finding
asked to propagate the caller context while stripping cancellation.

### Fix
`runUI.Close` now accepts `ctx context.Context` and derives the flush timeout
from `context.WithTimeout(context.WithoutCancel(ctx), store.JournalShutdownTimeout)`.
This keeps the existing `JournalShutdownTimeout` deadline and `FlushJournal`
behavior while propagating the caller's context (and any values it carries)
into the teardown flush. All call sites were updated to pass `ctx`
(`internal/cli/implement.go` and `internal/cli/cli.go`), including the deferred
calls, and the non-TTY test at `internal/cli/cli_test.go` was updated to
`ui.Close(ctx)`.

### Evidence
- `internal/cli/runui.go`: `func (ui *runUI) Close(ctx context.Context)` uses
  `context.WithoutCancel(ctx)`.
- `make verify` passes; `go test ./internal/cli/ -count=1` (1037 passed).

## Verification
- Focused: `go build ./internal/cli/...` and `go test ./internal/cli/ -count=1`.
- Authoritative: daemon runs `make verify`.

