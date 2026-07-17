---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/config/profile_config.go
line: 58
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5sl,comment:PRRC_kwDOS0qyts7Wt94-
review_hash: 4db4b4fb26d1a6092ea6660e45dd2f02b84da5a2e5aca245c440807abee2427d
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 006: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🌐 Web query:

`Does gopkg.in/yaml.v3 Decoder.Decode consume one YAML document per call, requiring another Decode call to detect trailing YAML documents?`

💡 Result:

Yes, gopkg.in/yaml.v3 Decoder.Decode consumes exactly one YAML document per call [1][2]. To process multiple YAML documents contained within a single stream, you must call Decode repeatedly in a loop until it returns io.EOF, which signals that no more documents are available in the input [1][3]. The Decoder is designed for streaming operations [1][4], and each call to Decode reads only the next available YAML-encoded value from the input stream [1][5][6]. This is distinct from the package's Unmarshal function, which only parses the first document found in a byte slice [5][2][6]. Example usage for processing multiple documents: decoder:= yaml.NewDecoder(reader) for { var doc YourType err:= decoder.Decode(&doc) if err == io.EOF { break // End of stream reached } if err!= nil { log.Fatal(err) } // Process doc }

Citations:

- 1: https://tessl.io/registry/tessl/golang-gopkg-in--yaml-v3/3.0.1/files/docs/streaming-operations.md
- 2: https://github.com/go-yaml/yaml/issues/805
- 3: https://stackoverflow.com/questions/67606439/reading-multiple-documents-from-a-yaml-file-and-marshaling-them-with-go-yaml-v3
- 4: https://code.asyncmachine.dev/pkg/gopkg.in/yaml.v3.html
- 5: https://github.com/go-yaml/yaml/blob/v3.0.1/yaml.go
- 6: https://github.com/npikall/gotpm/blob/main/vendor/gopkg.in/yaml.v3/yaml.go

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

# Inspect the relevant file and nearby code paths.
rg -n "ParseProfilesFragment|documentRootMapping|onlyProfilesWrapper|NormalizeProfilesFragment|profilesOverlay|yaml.NewDecoder|Decode\\(" internal/config/profile_config.go

echo
echo "---- file excerpt ----"
sed -n '1,260p' internal/config/profile_config.go
```

Repository: marcioaltoe/roundfix

Length of output: 8945

---



</details>

**Reject trailing YAML documents before parsing or rewriting.** `yaml.Decoder.Decode` only reads the first document here, so a fragment can hide extra content after `---`, and multi-document configs written back from `mergeProfilesConfigContent` will drop later documents. Require `io.EOF` after the first decode in both paths.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/config/profile_config.go` around lines 33 - 58, Update
ParseProfilesFragment and mergeProfilesConfigContent to perform a second
yaml.Decoder.Decode after the first document; require that it returns io.EOF,
and reject any additional YAML document before parsing or rewriting. Preserve
the existing empty-input and parse-error handling while ensuring multi-document
content is not silently accepted or dropped.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7496d19070d3c878e7808557 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
