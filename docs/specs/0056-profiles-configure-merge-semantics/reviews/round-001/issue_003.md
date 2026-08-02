---
source: coderabbit
pr: "67"
round: 1
round_created_at: "2026-08-02T11:30:39Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/profiles-configure-merge-semantics
head_sha: ffcc15ebed0a055d329cb3215ae0878b90931948
file: internal/cli/profiles_configure.go
line: 391
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vwh2H,comment:PRRC_kwDOS0qyts7ceBLo
review_hash: 068f1d0f2e2c061b9ee5d2402b6d117f3914c2c6867a001c2d7a464c1e69c5e8
duplicate_of: ""
source_review_id: "4838273774"
source_review_submitted_at: "2026-08-02T11:29:41Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Emit the JSON refusal report even when the stderr write fails.**

If the diagnostic write to stderr fails, the function returns immediately. A caller that passed `--json` then receives no stdout document at all, while the exit code is still `1`. The machine-readable contract breaks for the exact case the tests assert on. Write the JSON report first, then the human diagnostic.

<details>
<summary>🛠️ Proposed reordering</summary>

```diff
 func printProfilesConfigureRefusal(req profilesConfigureRequest, result roundconfig.ProfileConfigResult, stdout, stderr io.Writer) int {
 	result.Changed = false
-	if _, err := fmt.Fprintf(stderr, "Profile configuration unchanged: confirmation declined for %s\n", result.Path); err != nil {
-		return exitRunFailed
-	}
 	if req.json {
 		response := profilesConfigureResponseForResult(result, "")
 		response.Refused = true
 		if err := json.NewEncoder(stdout).Encode(response); err != nil {
 			return printProfilesConfigureOutputError(err, stderr)
 		}
 	}
+	fmt.Fprintf(stderr, "Profile configuration unchanged: confirmation declined for %s\n", result.Path)
 	return exitRunFailed
 }
```
</details>

Note: if you adopt the diff above, keep the `errcheck` contract by assigning the ignored return explicitly or by retaining an error check that does not skip the JSON output. As per coding guidelines: "Keep stdout limited to machine-readable result data; send logs, progress messages, and diagnostics to stderr" and "Always check returned errors; never discard them with `_`".

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/profiles_configure.go` around lines 378 - 391, In
printProfilesConfigureRefusal, emit the JSON response through
profilesConfigureResponseForResult and stdout before attempting the stderr
diagnostic, so --json always produces machine-readable output even when stderr
writing fails. Preserve the existing output-error handling and exitRunFailed
result, and explicitly handle the fmt.Fprintf return error without discarding
it.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:0401a91b430722dbb2c31dfb -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `printProfilesConfigureRefusal` returned on a stderr failure before encoding the requested JSON document. The function now encodes the refusal response first and then reports the human diagnostic, while preserving exit code `1` and explicit write-error handling.
- Focused evidence: `TestPrintProfilesConfigureRefusalWritesJSONBeforeDiagnostic` reproduced the defect before the fix with `decode profiles configure JSON "": EOF`; the focused CLI command then passed 17 tests after the reorder.
