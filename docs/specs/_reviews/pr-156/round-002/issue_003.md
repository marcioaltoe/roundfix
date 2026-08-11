---
source: coderabbit
pr: "156"
round: 2
round_created_at: "2026-08-11T19:18:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-journal-cheap-to-write-and-keep
head_sha: 5d670eebf70c3a291890416717848f2df5b2ce0d
file: internal/store/journal_writer.go
line: 101
severity: review
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YWRBD,comment:PRRC_kwDOS0qyts7gLDvW
review_hash: a5102445b5924cb1b9e3c422a074266f0824f2f49a86df19e27bdd77ef0de13f
duplicate_of: ""
source_review_id: "4909904494"
source_review_submitted_at: "2026-08-11T19:17:47Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Critical_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🔴 Critical_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

# Every reference to commitDone, with context.
rg -nP --type=go -C4 '\bcommitDone\b'

# Any receive or close on it.
rg -nP --type=go '<-\s*w\.commitDone|close\(\s*w\.commitDone\s*\)'

# Callers that rely on flush reporting durability.
rg -nP --type=go -C3 '\bFlushJournal\s*\(|\bCloseJournal\s*\('
```

Repository: marcioaltoe/roundfix

Length of output: 383

---

🏁 Script executed:

```shell
set -euo pipefail

printf '%s\n' '--- tracked file ---'
git ls-files -- internal/store/journal_writer.go internal/cli/runui.go

printf '%s\n' '--- journal writer symbols and commitDone references ---'
rg -n -C 5 'commitDone|func \(w \*journalWriter\) (publish|flush|drainLocked|commitDrained|close)|FlushJournal|CloseJournal' \
  internal/store/journal_writer.go internal/cli/runui.go

printf '%s\n' '--- relevant source ranges ---'
sed -n '80,235p' internal/store/journal_writer.go
sed -n '235,330p' internal/store/journal_writer.go
sed -n '350,430p' internal/store/journal_writer.go

printf '%s\n' '--- repository-wide references ---'
rg -n -C 3 'commitDone|FlushJournal|CloseJournal' --glob '*.go' .
```

Repository: marcioaltoe/roundfix

Length of output: 32597

---

🏁 Script executed:

```shell
set -euo pipefail

python3 - <<'PY'
from pathlib import Path

source = Path("internal/store/journal_writer.go").read_text()
refs = [line for line in source.splitlines() if "commitDone" in line]
print("commitDone references:", len(refs))
for line in refs:
    print(line)

assert len(refs) == 4, "expected only the declaration, initialization, and two comment lines"
assert "<-w.commitDone" not in source
assert "<- w.commitDone" not in source
assert "close(w.commitDone)" not in source

# Model the exact interleaving described by the review:
# A has already drained a batch and is committing outside the mutex.
pending = []
in_flight = True
durable = False

def drain():
    global pending
    if not pending:
        return None
    if in_flight:
        return None
    batch, pending = pending, []
    return batch

def flush():
    batch = drain()
    if batch is None:
        return True
    raise AssertionError("the modeled path should return before a commit")

# Immediate publish: the first flush observes A, then the event is appended.
assert flush() is True
pending.append("immediate-event")
assert flush() is True
assert pending == ["immediate-event"]
assert durable is False

# Count-triggered flush has the same result when pending is non-empty.
pending = ["count-triggered-event"]
assert flush() is True
assert pending == ["count-triggered-event"]

print("modeled result: flush returns nil while pending events remain and inFlight is true")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 470

---

🏁 Script executed:

```shell
set -euo pipefail

python3 - <<'PY'
from pathlib import Path

source = Path("internal/store/journal_writer.go").read_text()
refs = [line for line in source.splitlines() if "commitDone" in line]
print("commitDone references:", len(refs))
for line in refs:
    print(line)

assert len(refs) == 3, "unexpected commitDone reference count"
assert "close(w.commitDone)" not in source
assert "<-w.commitDone" not in source
assert "<- w.commitDone" not in source

# Model the relevant journal state transition.
# A has drained its batch and is committing outside the mutex.
pending = []
in_flight = True

def drain():
    global pending
    if len(pending) == 0:
        return None
    if in_flight:
        return None
    batch, pending = pending, []
    return batch

def flush():
    batch = drain()
    if batch is None:
        return None  # current implementation reports success

# Immediate-durability publish.
flush()
pending.append("immediate-event")
flush()
assert pending == ["immediate-event"]
print("immediate path: flush returns nil with pending event and inFlight=true")

# Count-triggered flush.
pending = ["count-triggered-event"]
flush()
assert pending == ["count-triggered-event"]
print("count path: flush returns nil with pending event and inFlight=true")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 450

---



</details>

**Wait for the in-flight commit in `flush`.**

`commitDone` is only declared and initialized. When `inFlight` is true, `drainLocked` returns `(nil, nil)`, so `flush` returns `nil` while pending events remain uncommitted. This breaks immediate durability and count-triggered flushes. Wait on `commitDone`, then re-run the drain loop. Close and replace the channel in `commitDrained`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/journal_writer.go` around lines 86 - 101, Update flush to wait
on commitDone whenever drainLocked reports an in-flight commit instead of
returning while events remain pending, then retry the drain loop after the
channel signals. In commitDrained, close the current commitDone channel and
replace it with a new channel for the next commit, preserving synchronization
for subsequent waiters.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:26c254d38beebe0955c77e25 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:

## Resolution

- Status: `resolved`
- Date: 2026-08-11

### Finding
`commitDone` was declared and initialized but never awaited or closed. When
`inFlight` was true, `drainLocked` returned `(nil, nil)`, so `flush` returned
`nil` while pending events remained uncommitted, breaking immediate durability
and count-triggered flushes during a concurrent in-flight commit.

### Fix
- `flush` now inspects whether `drainLocked` reported an in-flight commit. When
  it did, `flush` captures `commitDone` (under the mutex) and waits on it via
  the new `waitCommit` helper, then re-runs the drain loop so events appended or
  re-queued behind the in-flight commit are not skipped. Waiting respects
  `ctx.Done()` so a cancelled context still returns an error.
- `drainLocked` now reports the in-flight condition through a second return
  value instead of conflating it with "nothing pending".
- `commitDrained` closes the current `commitDone` channel and replaces it with a
  fresh one so waiters observe the commit's completion and subsequent commits
  have a new channel (exactly one close per drain).

### Evidence
- `internal/store/journal_writer.go`: `flush`, `waitCommit`, `drainLocked`, and
  `commitDrained` updated as described.
- `go test ./internal/store/ -count=1` and the full `make verify` pass.

## Verification
- Focused: `go test ./internal/store/ -count=1` (all journal/batch tests pass).
- Authoritative: daemon runs `make verify`.

