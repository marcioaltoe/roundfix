---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/selection_assignment.go
line: 457
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWis,comment:PRRC_kwDOS0qyts7fswOQ
review_hash: 8aa8a3bbfe0d93c03537ab57bfc717a2c632a6d9f2014539adbddb6882a2870f
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:14Z"
---

# Issue 008: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**An empty advertised model list silently disables the membership refusal.**

The guard `len(request.Catalogue.Models) > 0` treats two different states as identical:

- No catalogue was read, so there is no evidence. Skipping the check is correct.
- A catalogue was read and the runtime advertised zero models. Skipping the check removes the net that ADR-0119 keeps for a runtime that declines to refuse.

`readRuntimeCatalogueWithEvidence` returns `RuntimeCatalogue{Models: models}` where `models` is empty when `capabilities.Models` is empty (line 204). That value is indistinguishable from the zero value, so the refusal disappears without any signal. Record whether the catalogue was observed, then refuse when an observed catalogue advertises nothing.



<details>
<summary>🛡️ Proposed fix to distinguish an observed empty catalogue</summary>

```diff
 type RuntimeCatalogue struct {
 	Models       []string
 	Efforts      []string
+	Observed     bool
 	Contaminated bool
 }
```

```diff
 	return RuntimeCatalogue{Models: models, Efforts: efforts}
```
becomes
```diff
-	return RuntimeCatalogue{Models: models, Efforts: efforts}
+	return RuntimeCatalogue{Models: models, Efforts: efforts, Observed: true}
```

```diff
 	requestedModel := strings.TrimSpace(request.Runtime.Model)
-	if len(request.Catalogue.Models) > 0 && !request.Catalogue.AdvertisesModel(requestedModel) {
+	if request.Catalogue.Observed && !request.Catalogue.AdvertisesModel(requestedModel) {
 		return SelectionProof{}, &ModelNotAdvertisedError{
```

Update `recordAdvertisement` to gate on `Observed` instead of `len(catalogue.Models) == 0`.
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/selection_assignment.go` around lines 450 - 457, Distinguish
an unread catalogue from an observed catalogue advertising zero models by adding
and propagating an observation signal through readRuntimeCatalogueWithEvidence
and RuntimeCatalogue. Update recordAdvertisement to gate membership refusal on
Observed rather than len(catalogue.Models) > 0, while preserving the existing
advertised-model validation for non-empty catalogues.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e63d9c6212416d7648f62c1f -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Distinguished an unread catalogue from an observed catalogue advertising zero models. Added `Observed bool` to `RuntimeCatalogue`; `runtimeCatalogueFromCapabilities` now returns `Observed: true`; `applySessionSelection` refuses when `request.Catalogue.Observed && !AdvertisesModel` instead of `len(Models) > 0`; `recordAdvertisement` gates on `!catalogue.Observed` instead of `len == 0`. Test fixtures that represent observed catalogues set `Observed: true`. Focused: `go test ./internal/agent -run 'TestProof|TestRuntime|TestSelectionCatalogueCharacterization'` ok.
