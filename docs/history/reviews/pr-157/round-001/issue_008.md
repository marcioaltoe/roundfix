---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: internal/spec/archive_layout_characterization_test.go
line: 93
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdV,comment:PRRC_kwDOS0qyts7gSdx4
review_hash: c29f2f887b077917bc348fdbdc37d1be3a7484a4825a12fd04957324369e8cce
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 008: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Replace detached test data with contract assertions.**

These tests compare local values with duplicated local literals. They do not inspect archive consumers or the current Secondbrain guidance. They will pass if another package composes an archive path or conditional guidance returns.

Remove these records, or assert the affected behavior and guidance files directly.






As per coding guidelines: `**/*_test.go`: NEVER test implementation details; test observable behavior and public API contracts.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/archive_layout_characterization_test.go` around lines 69 - 93,
Replace the detached checks in
TestArchiveLayoutCharacterizationEnumeratesEveryPathComposer and
TestArchiveLayoutCharacterizationCapturesConditionalSecondbrainClause with
assertions against the actual archive-consuming behavior and current Secondbrain
guidance files. Remove archivePathComposersAfterTask04,
conditionalSecondbrainClauseBeforeSpec0085, and duplicated expected literals,
and verify the externally observable paths and escape-hatch guidance through
their public contract.
```

</details>

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:0b175774fb4ab027bbca0b78 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: In `internal/spec/archive_layout_characterization_test.go`, removed the two detached, tautological checks:
  - `TestArchiveLayoutCharacterizationEnumeratesEveryPathComposer` and the `archivePathComposersAfterTask04`/`archivePathComposer` data — they only compared a local slice against a re-declared copy of itself (the four `_archived/<kind>` literals also live in `ArchiveDir`, so the check added nothing and would pass even if another package composed archive paths).
  - `TestArchiveLayoutCharacterizationCapturesConditionalSecondbrainClause` and the `conditionalSecondbrainClauseBeforeSpec0085` constant — they asserted a substring of a local literal, and the guidance this PR made unconditional no longer carries that conditional clause.
  The two remaining characterization tests (`RecordsEveryRetiredFamily` and `PinsCorpusGoldenBeforeRelocation`) assert real, externally observable state via `os.Stat` and the corpus golden file, and stay. Removed `encoding/json`/`strings`/`reflect` usage is still present for those two tests, so imports remain valid. Focused evidence: `go test ./internal/spec/ -run 'TestArchiveLayoutCharacterization|TestArchiveSpecRoot'` passed; `make verify` passes.
