---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T15:34:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: a12c1a665c5970773e04c4a145c6b9b0c5a0e686
file: internal/spec/spec.go
line: 155
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WBeOR,comment:PRRC_kwDOS0qyts7c2V1E
review_hash: 5ee0881768a0058f98d8ceaf4b3da69b20521b0a8e6206c08ccf58dee410a294
duplicate_of: ""
source_review_id: "4845660382"
source_review_submitted_at: "2026-08-03T15:14:34Z"
---

# Issue 007: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Confirm yaml.v3 node kinds and tags for the QA declaration shapes.
set -euo pipefail

# Verify the declared yaml.v3 version in the module manifest.
fd -t f '^go.mod$' -d 1 --exec rg -n 'yaml' {}

# Locate other UnmarshalYAML implementations for comparison of node-tag handling.
ast-grep run --pattern 'func ($_ $_) UnmarshalYAML(node *yaml.Node) error { $$$ }' --lang go .
```

Repository: marcioaltoe/roundfix

Length of output: 12750

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- yaml.v3 module source availability ---'
gopath="$(go env GOPATH 2>/dev/null || true)"
if [ -n "$gopath" ]; then
  find "$gopath/pkg/mod" -maxdepth 4 -type d -path '*gopkg.in/yaml.v3@v3.0.1' -print 2>/dev/null || true
fi

printf '%s\n' '--- repository dependency metadata ---'
rg -n 'gopkg.in/yaml.v3 v3\.0\.1|yaml\.v3' go.mod go.sum
```

Repository: marcioaltoe/roundfix

Length of output: 496

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

cat >/tmp/yaml_node_probe.go <<'EOF'
package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type manifestOptionalString struct {
	Present bool
	Value   string
}

func (value *manifestOptionalString) UnmarshalYAML(node *yaml.Node) error {
	fmt.Printf("node kind=%d tag=%q value=%q\n", node.Kind, node.Tag, node.Value)
	value.Present = true
	if node.Kind == yaml.ScalarNode && node.Tag == "!!null" {
		return nil
	}
	if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
		return fmt.Errorf("must be a string")
	}
	value.Value = node.Value
	return nil
}

type manifest struct {
	QA       manifestOptionalString `yaml:"qa"`
	QAReason manifestOptionalString `yaml:"qa_reason"`
}

func main() {
	for _, input := range []string{
		"qa: declined\n",
		"qa_reason:\n",
		"qa: 123\n",
		"qa: true\n",
	} {
		var value manifest
		err := yaml.Unmarshal([]byte(input), &value)
		fmt.Printf("input=%q err=%v qa=(present:%t value:%q) reason=(present:%t value:%q)\n",
			input, err, value.QA.Present, value.QA.Value, value.QAReason.Present, value.QAReason.Value)
	}
}
EOF

go run /tmp/yaml_node_probe.go
```

Repository: marcioaltoe/roundfix

Length of output: 661

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- manifest types and decoder ---'
cat -n internal/spec/spec.go | sed -n '100,210p'

printf '%s\n' '--- manifest-related tests and empty-reason cases ---'
rg -n -C 5 'qa_reason|reason cannot be empty|qa:|declined|manifestOptionalString|invalid frontmatter' --glob '*_test.go' internal .
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- manifest decode and QA validation ---'
rg -n -C 12 'manifestFrontmatter|QAReason|qaDeclaration|reason required|cannot be empty|qa_reason' internal/spec/spec.go internal/spec/spec_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 19272

---



</details>

**Track explicit `qa_reason:` presence outside `UnmarshalYAML`.** `yaml.v3` skips `manifestOptionalString.UnmarshalYAML` for a null mapping value, so `qa_reason:` leaves `Present` false and bypasses the required non-empty-reason check. Decode the field through a presence-aware node or pointer before validation.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/spec.go` around lines 145 - 155, Update the manifest decoding
and validation flow around manifestOptionalString.UnmarshalYAML so explicit
qa_reason: null or empty mapping values are detected even when yaml.v3 skips
UnmarshalYAML. Decode qa_reason through a presence-aware yaml.Node or pointer,
set the optional field’s Present state based on key presence, and ensure the
existing non-empty-reason validation runs for explicitly present values.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:09ddb3832c3b4ba6f0e80dcd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `manifestFrontmatter.UnmarshalYAML` now records `qa` and `qa_reason` key presence from the mapping node even when yaml.v3 skips the field-level decoder for null values. A regression proves standalone `qa_reason:` is rejected, and `go test ./internal/spec -count=1` passed with 126 tests.
