---
source: coderabbit
pr: "156"
round: 2
round_created_at: "2026-08-11T19:18:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-journal-cheap-to-write-and-keep
head_sha: 5d670eebf70c3a291890416717848f2df5b2ce0d
file: internal/tui/cockpit.go
line: 333
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YWRBJ,comment:PRRC_kwDOS0qyts7gLDvf
review_hash: d44668b3b9ab48e8dfc65ac39a27a261c6b716162b753b4c149f6a6be01f90d0
duplicate_of: ""
source_review_id: "4909904494"
source_review_submitted_at: "2026-08-11T19:17:48Z"
---

# Issue 004: _ Performance & Scalability_ _ Major_ _ Heavy lift_

## Review Comment

_🚀 Performance & Scalability_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

ast-grep outline internal/tui/cockpit.go --match cockpitModel --view expanded
rg -n -C4 --type=go '\brefresh(TaskJournalEvents|BatchClocks)\s*\(' internal/tui
rg -n -C4 --type=go 'journalHeaderPageSize|RunEventHeadersAfter' internal/tui
```

Repository: marcioaltoe/roundfix

Length of output: 17726

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- cockpit update and refresh flow ---'
sed -n '175,345p' internal/tui/cockpit.go
printf '%s\n' '--- batch clock refresh ---'
sed -n '560,645p' internal/tui/cockpit.go
printf '%s\n' '--- header query implementations ---'
rg -n -C8 --type=go 'func .*RunEventHeadersAfter|RunEventHeadersAfter\(' .
printf '%s\n' '--- relevant tests ---'
sed -n '1,230p' internal/tui/cockpit_forward_cursor_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 36396

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- model construction and initial refresh ---'
sed -n '120,180p' internal/tui/cockpit.go
printf '%s\n' '--- complete header query implementation ---'
sed -n '1240,1305p' internal/store/journal.go
printf '%s\n' '--- fake paging behavior ---'
sed -n '1,75p' internal/tui/viewport_test.go
printf '%s\n' '--- constructor call sites ---'
rg -n -C5 --type=go 'newCockpitModel\(|tea\.NewProgram|NewProgram\(' internal/tui
```

Repository: marcioaltoe/roundfix

Length of output: 18291

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
from pathlib import Path

source = Path("internal/tui/cockpit.go").read_text()

required = [
    "case cockpitTickMsg:",
    "model.poll()",
    "model.refreshWorkItems()",
    "for {\n\t\theaders, err := model.cfg.Source.RunEventHeadersAfter(model.ctx, runID, model.taskJournalCursor, journalHeaderPageSize)",
    "for {\n\t\theaders, err := model.cfg.Source.RunEventHeadersAfter(model.ctx, runID, model.batchTimeCursor, journalHeaderPageSize)",
]
for marker in required:
    if marker not in source:
        raise SystemExit(f"missing expected synchronous-flow marker: {marker!r}")

def calls_for_one_refresh(total, page_size=500):
    cursor = 0
    calls = 0
    while True:
        page = list(range(cursor + 1, min(total, cursor + page_size) + 1))
        calls += 1
        if page:
            cursor = page[-1]
        if len(page) < page_size:
            return calls, cursor

for total, expected_calls in [(0, 1), (500, 2), (501, 2), (1001, 3)]:
    calls, cursor = calls_for_one_refresh(total)
    assert calls == expected_calls, (total, calls, expected_calls)
    assert cursor == total, (total, cursor)

print("synchronous tick path present")
print("one-refresh page simulation: 1001 headers -> 3 synchronous header queries")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 262

---



</details>

**Bound journal-header work per refresh and startup.**

`Update` handles `cockpitTickMsg` synchronously. Both refresh functions read every 500-row page before returning, and `newCockpitModel` performs the initial refresh synchronously. A large backlog can delay rendering and input handling.

Process at most one page per invocation. Schedule continuation messages for remaining pages. Move initial backlog processing onto the Bubble Tea update path. Apply this to both `internal/tui/cockpit.go#L303-L333` and `internal/tui/cockpit.go#L615-L638`. Add a test with more than 500 headers that proves one refresh does not read the next page.

<details>
<summary>📍 Affects 1 file</summary>

- `internal/tui/cockpit.go#L303-L333` (this comment)
- `internal/tui/cockpit.go#L615-L638`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/tui/cockpit.go` around lines 303 - 333, The journal-header refresh
loops at internal/tui/cockpit.go:303-333 and internal/tui/cockpit.go:615-638
must process at most one page per invocation; schedule a continuation message
when a full page indicates more headers remain, preserving cursor/retry
behavior. Update newCockpitModel so the initial backlog refresh is dispatched
through Bubble Tea’s update path rather than performed synchronously, and add a
test with more than 500 headers proving one refresh does not read the next page.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/tui/cockpit.go</file>
<line_range>303-333</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/tui/cockpit.go</file>
<line_range>615-638</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:911dbc47d2ac337df984114b -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:

## Resolution

- Status: `resolved`
- Date: 2026-08-11

### Finding
Both journal-header refresh loops (`refreshTaskJournalEvents` and
`refreshBatchClocks`) read every 500-row page synchronously before returning,
and `newCockpitModel` performed the initial refresh synchronously. A large
backlog could delay rendering and input handling on a single `Update` frame.

### Fix
- Both refresh functions now fold at most **one** page per invocation and return
  whether a full page was read (more may remain).
- `refreshWorkItems`/`refreshTasks`/`refreshIssues`/`poll` propagate that
  "more" signal up to `Update`.
- Added a `cockpitJournalContinueMsg` handled in `Update`: when a refresh returns
  a full page, it schedules a continuation message to fold the next page on the
  next frame instead of draining synchronously in one call.
- `newCockpitModel` still performs the bounded one-page opening fold synchronously
  (the test seam), but records whether more pages remain in `journalMore` and
  `Init()` dispatches the remaining backlog through the Bubble Tea update path.
- `continueJournalPage` folds one more page of whichever cursor is active (Task
  journal for spec Runs, Batch clocks otherwise).

### Test
Added `TestCockpitJournalHeaderRefreshBoundedToOnePage` in
`internal/tui/cockpit_forward_cursor_test.go`, which seeds a journal spanning
multiple pages (>500 headers) and asserts the opening fold and each individual
refresh read exactly one page and report more until the short final page, rather
than draining the whole journal in one refresh.

### Evidence
- `internal/tui/cockpit.go` and `internal/tui/cockpit_forward_cursor_test.go`
  updated as described.
- `go test ./internal/tui/ -count=1` passes (including the new test).
- Full `make verify` passes.

## Verification
- Focused: `go test ./internal/tui/ -count=1` (all cockpit/forward-cursor tests pass).
- Authoritative: daemon runs `make verify`.

