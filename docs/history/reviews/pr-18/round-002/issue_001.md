---
source: coderabbit
pr: "18"
round: 2
round_created_at: "2026-07-07T14:05:56Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 4237143afdd7097e755e14b962156aaf6c6e6654
file: internal/cli/attach.go
line: 262
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6O6PJw,comment:PRRC_kwDOS0qyts7SyYRd
review_hash: eb8ffbb97f231d7f5a03ee69d2074efe4cbfb8a445b8ba73b568f562a006347b
duplicate_of: ""
source_review_id: "4645087962"
source_review_submitted_at: "2026-07-07T12:31:07Z"
---

# Issue 001: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Blocking stdin read ignores context cancellation.**

`ctx` is only checked once up front; the subsequent `bufio.NewReader(stdin).ReadString('\n')` blocks regardless of `ctx.Done()`. If the command's context is canceled (signal/timeout) while the picker is waiting for input, this hangs. As per coding guidelines, "keep goroutines owned and cancellable" and "prefer context-first signatures for blocking operations."





<details>
<summary>🔒️ Proposed fix: make the read cancellable</summary>

```diff
-	fmt.Fprint(stderr, renderAttachRunPicker(runs))
-	line, err := bufio.NewReader(stdin).ReadString('\n')
-	if err != nil && err != io.EOF {
-		return "", fmt.Errorf("read attach Run picker: %w", err)
-	}
+	fmt.Fprint(stderr, renderAttachRunPicker(runs))
+	type readResult struct {
+		line string
+		err  error
+	}
+	resultCh := make(chan readResult, 1)
+	go func() {
+		line, err := bufio.NewReader(stdin).ReadString('\n')
+		resultCh <- readResult{line, err}
+	}()
+	var line string
+	select {
+	case <-ctx.Done():
+		return "", ctx.Err()
+	case res := <-resultCh:
+		if res.err != nil && res.err != io.EOF {
+			return "", fmt.Errorf("read attach Run picker: %w", res.err)
+		}
+		line = res.line
+	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func collectAttachRunSelection(ctx context.Context, runs []store.Run, stdin io.Reader, stderr io.Writer) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if stdin == nil {
		stdin = strings.NewReader("")
	}
	if stderr == nil {
		stderr = io.Discard
	}
	fmt.Fprint(stderr, renderAttachRunPicker(runs))
	type readResult struct {
		line string
		err  error
	}
	resultCh := make(chan readResult, 1)
	go func() {
		line, err := bufio.NewReader(stdin).ReadString('\n')
		resultCh <- readResult{line, err}
	}()
	var line string
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resultCh:
		if res.err != nil && res.err != io.EOF {
			return "", fmt.Errorf("read attach Run picker: %w", res.err)
		}
		line = res.line
	}
	choice := strings.TrimSpace(line)
	if choice == "" || isAttachPickerCancelChoice(choice) {
		return "", errAttachPickerCanceled
	}
	if index, err := strconv.Atoi(choice); err == nil {
		if index >= 1 && index <= len(runs) {
			return runs[index-1].ID, nil
		}
		return "", validationError{message: fmt.Sprintf("Run picker choice %d is out of range", index)}
	}
	for _, run := range runs {
		if choice == run.ID {
			return run.ID, nil
		}
	}
	return "", validationError{message: fmt.Sprintf("Run picker choice %q is not a listed Run", choice)}
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/attach.go` around lines 231 - 262, Blocking stdin read in
collectAttachRunSelection ignores context cancellation, so the picker can hang
after ctx is canceled. Update collectAttachRunSelection to make the
ReadString('\n') path cancellable by watching ctx.Done() while waiting for
input, and return ctx.Err() promptly if the context is canceled. Keep the
existing selection parsing and validation logic intact, but ensure the blocking
read is interrupted rather than waiting indefinitely.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:fba3e6de03d134aeffd4cfdd -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - `collectAttachRunSelection` now waits on `ctx.Done()` while stdin is blocked and returns the context cancellation error promptly; regression coverage exercises cancellation during a blocked read.
  - Verification: configured command `make verify` was run as `rtk make verify` and passed: Go tests, skills check, and build completed.
