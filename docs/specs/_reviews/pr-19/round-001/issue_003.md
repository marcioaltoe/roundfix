---
source: coderabbit
pr: "19"
round: 1
round_created_at: "2026-07-07T22:44:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/run-browser
head_sha: f2726493dff5e63e604139d27d147973ff650cf5
file: internal/tui/runbrowser.go
line: 306
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6PCMCm,comment:PRRC_kwDOS0qyts7S9VwH
review_hash: c49fa19f5f74c8c512d98ab1f198b5737e72ffe818edcf471ca10c3327389caf
duplicate_of: ""
source_review_id: "4648487653"
source_review_submitted_at: "2026-07-07T19:59:41Z"
---

# Issue 003: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Wrap the `prog.Run()` error and guard the type assertion.**

Two coding-guideline gaps in this block:
- The `prog.Run()` error is returned bare instead of wrapped with `%w`/an operation name.
- `final.(runBrowserProgram)` is an unchecked type assertion; a mismatch would panic outside `main`.

As per coding guidelines, `**/*.go` code should "wrap errors with `%w`" and "avoid `panic`/`log.Fatal` outside unrecoverable startup in `main`," and `internal/**/*.go` errors "must name the failed operation and the next useful action when one is known."




<details>
<summary>🛠️ Proposed fix</summary>

```diff
+	"fmt"
 	"io"
 	"strings"
 	"time"
@@
 	final, err := prog.Run()
 	if err != nil {
-		return BrowserOutcome{}, err
+		return BrowserOutcome{}, fmt.Errorf("run browser session: %w", err)
 	}
-	outcome := final.(runBrowserProgram).browser.Outcome()
+	program, ok := final.(runBrowserProgram)
+	if !ok {
+		return BrowserOutcome{}, fmt.Errorf("run browser session: unexpected model type %T", final)
+	}
+	outcome := program.browser.Outcome()
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	final, err := prog.Run()
	if err != nil {
		return BrowserOutcome{}, fmt.Errorf("run browser session: %w", err)
	}
	program, ok := final.(runBrowserProgram)
	if !ok {
		return BrowserOutcome{}, fmt.Errorf("run browser session: unexpected model type %T", final)
	}
	outcome := program.browser.Outcome()
	if outcome.RunID == "" && !outcome.Cancelled {
		// The program ended without a key outcome: context cancellation
		// quit it, which is a cancel, never a selection.
		outcome.Cancelled = true
	}
	return outcome, nil
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/tui/runbrowser.go` around lines 291 - 301, Wrap the prog.Run()
failure with an operation-specific error using %w and include the failed action
plus the next useful step, and replace the unchecked final.(runBrowserProgram)
assertion with a safe type check that returns an error instead of panicking if
the value is not the expected browser program. Update the runBrowser flow in the
function handling prog.Run() and BrowserOutcome extraction so all failures are
returned cleanly from internal/tui/runbrowser.go.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c28ba37f7cf7634cdd62af2b -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `RunBrowserSession` returned the `prog.Run()` error without operation context and used an unchecked final-model assertion.

## Resolution

- Status: `resolved`
- Changes: Wrapped `prog.Run()` failures with `run browser session: %w`, added safe final-model extraction, and covered unexpected model extraction with a regression test.
- Verification: `rtk make verify` passed in this session for the configured `make verify` gate.
