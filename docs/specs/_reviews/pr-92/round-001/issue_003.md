---
source: coderabbit
pr: "92"
round: 1
round_created_at: "2026-08-03T19:34:00Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spawn-economy
head_sha: 7765d1f6d62e59ebf68ca2e4e2e273733da58425
file: internal/baseline/skills_restore_git.go
line: 121
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WGXLH,comment:PRRC_kwDOS0qyts7c9hTU
review_hash: 7c9fb4ab0abda37c6df8107371e090874d2fa98cf637bf5a8a3f15ca357d8f5e
duplicate_of: ""
source_review_id: "4847882119"
source_review_submitted_at: "2026-08-03T19:33:11Z"
---

# Issue 003: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Bound the declared blob size and poison the reader after a framing error.**

Two hardening points on `Read`:

1. Line 110 allocates exactly the size declared by the header. The header comes from an externally fetched GitHub commit, so a single oversized blob allocates that much memory before `restoreMaxBytes` is ever consulted at line 314. Reject the size first.
2. A malformed header, a short content read, or a bad terminator leaves the stream desynchronized while `closed` stays false. A future caller that continues reading would receive misaligned bytes and compute a wrong digest. Mark the reader unusable on the first protocol failure.

<details>
<summary>🛡️ Proposed hardening</summary>

```diff
 type batchObjectReader struct {
 	stdin    io.WriteCloser
 	stdout   *bufio.Reader
 	cmd      *exec.Cmd
 	stderr   *bytes.Buffer
+	maxBytes int
+	broken   bool
 	closed   bool
 	closeErr error
 }
```

```diff
 	size, err := strconv.ParseInt(fields[2], 10, strconv.IntSize)
 	if err != nil || size < 0 {
+		reader.broken = true
 		return nil, fmt.Errorf("invalid batch object size %q", fields[2])
 	}
+	if reader.maxBytes > 0 && size > int64(reader.maxBytes) {
+		reader.broken = true
+		return nil, fmt.Errorf("batch object %s size %d exceeds the read limit %d", object, size, reader.maxBytes)
+	}
 	content := make([]byte, int(size))
 	if _, err := io.ReadFull(reader.stdout, content); err != nil {
+		reader.broken = true
 		return nil, fmt.Errorf("read batch object content: %w", err)
 	}
 	terminator, err := reader.stdout.ReadByte()
 	if err != nil {
+		reader.broken = true
 		return nil, fmt.Errorf("read batch object terminator: %w", err)
 	}
 	if terminator != '\n' {
+		reader.broken = true
 		return nil, fmt.Errorf("invalid batch object terminator %#x", terminator)
 	}
```

Then guard entry alongside the existing closed check:

```diff
-	if reader.closed {
+	if reader.closed {
 		return nil, errors.New("batch object reader is closed")
 	}
+	if reader.broken {
+		return nil, errors.New("batch object reader is desynchronized")
+	}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/skills_restore_git.go` around lines 106 - 121, Harden the
batch-object read path around the visible Read logic: reject declared sizes
above the configured restoreMaxBytes limit before allocating content, and mark
the reader closed on malformed headers, oversized sizes, short content reads,
terminator read failures, or invalid terminators so subsequent reads fail
immediately. Extend the existing closed-state guard at Read entry to return the
established unusable-reader error.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ea85af8dc851fd4483014306 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `Read` allocated the header-declared blob size before the 8 MiB restoration limit was checked, and framing failures left the reader available for another request on a desynchronized stream.

## Result

- Added a per-reader byte limit that rejects blobs larger than `restoreMaxBytes` before allocation.
- Marked the reader desynchronized after request-write, header, type, size, content, or terminator failures; later reads now fail before writing another request. A valid Git `missing` reply remains retryable because it consumes a complete protocol frame.
- Added regressions for header failure, malformed headers, oversized declarations, short content, invalid terminators, and the subsequent-read guard.
- Focused check: `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache-run_20260803T193424Z-baseline go test ./internal/baseline -run '^TestBatchObjectReader(ReportsProcessDeathMidStream|RejectsOversizedObjectsAndPoisonsProtocolFailures)$' -count=1` — passed.
- Package check: `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache-run_20260803T193424Z-baseline go test ./internal/baseline -count=1` — passed.
- The Daemon owns the configured `make verify` run.
