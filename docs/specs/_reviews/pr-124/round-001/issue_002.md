---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/cli/cli_test.go
line: 6805
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wuazm,comment:PRRC_kwDOS0qyts7d35s1
review_hash: 4db0122da6b05acbe5ad75f85597414f526bcceaff9dd174bd6287d732d53e65
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Assert that an outcome event was found.**

The loop decodes every `KindDaemonOutcome` event into `outcome`. If no such event exists, `outcome` stays zero-valued and the assertion at Line 6800 fails with the message "outcome payload = ...State:\"\"...". That message does not distinguish a missing event from a wrong state.

The neighboring tests at Line 6705 and Line 6949 track a `foundOutcome` flag for this reason. Apply the same pattern here.



<details>
<summary>💚 Proposed fix</summary>

```diff
 	_, events := journaledRunEvents(t, homeDir, stderr.String())
 	var outcome runevent.OutcomePayload
+	foundOutcome := false
 	for _, journaled := range events {
 		if journaled.Event.Kind != runevent.KindDaemonOutcome {
 			continue
 		}
 		if err := json.Unmarshal(journaled.Event.Payload, &outcome); err != nil {
 			t.Fatalf("decode unrecognised green signal outcome: %v", err)
 		}
+		foundOutcome = true
 	}
-	if outcome.State != store.StateTimedOut ||
+	if !foundOutcome ||
+		outcome.State != store.StateTimedOut ||
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	var outcome runevent.OutcomePayload
	foundOutcome := false
	for _, journaled := range events {
		if journaled.Event.Kind != runevent.KindDaemonOutcome {
			continue
		}
		if err := json.Unmarshal(journaled.Event.Payload, &outcome); err != nil {
			t.Fatalf("decode unrecognised green signal outcome: %v", err)
		}
		foundOutcome = true
	}
	if !foundOutcome ||
		outcome.State != store.StateTimedOut ||
		outcome.EvidenceState != string(reviewsource.EvidencePending) ||
		!strings.Contains(outcome.Reason, "signal was not recognised") ||
		!strings.Contains(outcome.Reason, detail) {
		t.Fatalf("unrecognised green signal outcome payload = %#v", outcome)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 6790 - 6805, Update the test around
the event loop decoding KindDaemonOutcome to track a foundOutcome flag, set it
when an outcome event is encountered, and assert that the flag is true before
validating outcome fields. Follow the neighboring tests’ pattern so missing
events produce a distinct failure from an incorrect payload.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ddd218f5cd252d94378f2b5a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The CLI regression now records `foundOutcome` only after decoding a
  `KindDaemonOutcome` event and requires it before asserting payload fields.
  A missing event is therefore distinguishable from a zero-valued payload.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test ./internal/cli -count=1 -run
  '^TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush$'` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
