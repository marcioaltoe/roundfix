---
source: coderabbit
pr: "42"
round: 1
round_created_at: "2026-07-29T02:33:36Z"
status: invalid
terminal_reason: "AdapterLineageError already implements a nil-safe Unwrap method that returns Err, so the reported missing error-chain exposure does not exist."
head_repository: marcioaltoe/roundfix
head_branch: ma/claude-adapter-standardization
head_sha: 7155ba4d2ef353257a1bacf697027202d4750492
file: internal/agent/acpx_runner.go
line: 399
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UmtXZ,comment:PRRC_kwDOS0qyts7ayCzv
review_hash: 2a6ac49c24f941ea4ba181614f8efbb8e453b3669733a6b138b66cddd6b824c6
duplicate_of: ""
source_review_id: "4803488138"
source_review_submitted_at: "2026-07-29T02:32:46Z"
---

# Issue 001: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**`AdapterLineageError` carries `Err` but never exposes it.**

`newAdapterLineageError` stores the underlying `exec` failure in `Err`, yet `Error()` omits it and the type has no `Unwrap()`, so `errors.Is`/`errors.As` cannot traverse to the root cause (e.g. a context or exec error) and the information is silently dropped.

As per coding guidelines, "Custom error types that wrap another error must implement `Unwrap()` so `errors.Is` and `errors.As` can traverse the error chain."






<details>
<summary>🛠️ Proposed fix</summary>

```diff
 func (err *AdapterLineageError) Classification() string {
 	return AdapterLineageUnknown
 }
+
+func (err *AdapterLineageError) Unwrap() error {
+	if err == nil {
+		return nil
+	}
+	return err.Err
+}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
type AdapterLineageError struct {
	Runtime         string
	Command         string
	Package         string
	Version         string
	RequiredPackage string
	RequiredVersion string
	Install         string
	Legacy          bool
	Err             error
}

func (err *AdapterLineageError) Error() string {
	if err == nil {
		return ""
	}
	runtimeID := adapterErrorRuntime(err.Runtime)
	requiredPackage := adapterErrorRequiredPackage(err.RequiredPackage)
	requiredVersion := adapterErrorRequiredVersion(err.RequiredVersion)
	message := fmt.Sprintf("effective %s adapter command %q did not prove required package lineage %s", runtimeDisplayName(runtimeID), strings.TrimSpace(err.Command), requiredPackage)
	if err.Legacy || isLegacyAdapterPackage(err.Package) {
		message = fmt.Sprintf("effective %s adapter command %q reported legacy package %s", runtimeDisplayName(runtimeID), strings.TrimSpace(err.Command), err.Package)
	} else if packageName := strings.TrimSpace(err.Package); packageName != "" && packageName != requiredPackage {
		message = fmt.Sprintf("effective %s adapter command %q reported unknown package %s", runtimeDisplayName(runtimeID), strings.TrimSpace(err.Command), packageName)
	}
	if strings.TrimSpace(err.Version) != "" {
		message += " version " + strings.TrimSpace(err.Version)
	}
	return fmt.Sprintf("%s; required %s %s or newer; update with: %s", message, requiredPackage, requiredVersion, err.InstallCommand())
}

func (err *AdapterLineageError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner.go` around lines 370 - 399, Update
AdapterLineageError to expose its underlying Err by implementing an Unwrap()
error method that returns err.Err, preserving nil safety as appropriate. Keep
the existing Error() message and lineage-specific formatting unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c1bf6bb259f55d6406b5ccca -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes:
  - Current code at `internal/agent/acpx_runner.go` defines `(*AdapterLineageError).Unwrap()` and returns `err.Err`, with a nil receiver returning `nil`.
  - `rtk env GOCACHE=/private/tmp/roundfix-review-001-gocache.QR9F0C go test ./internal/agent` passed.
