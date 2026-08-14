---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/reviewsource/coderabbit/coderabbit.go
line: 389
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wua0Y,comment:PRRC_kwDOS0qyts7d35tz
review_hash: 1248bde2bbc9a7cbf0a97b5c6a2aaeadf7e29fcf11b041f7632a7ff9215abe9a
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 008: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🌐 Web query:

`gh api --paginate multiple JSON arrays invalid JSON --slurp`

💡 Result:

When using gh api with --paginate, the command makes multiple sequential HTTP requests, returning each page as a separate, independent JSON object or array [1][2]. Because these pages are returned as separate top-level JSON structures, the raw output is technically invalid JSON if viewed as a single document [3][4]. The --slurp flag is designed to address this by wrapping all pages of JSON arrays or objects into a single outer JSON array [1][2]. Key considerations and limitations include: 1. Incompatibility with --jq or --template: The --slurp option cannot be used in combination with --jq or --template [5][2]. If you need to filter or process the data while paginating, you should omit --slurp and pipe the output to the standalone jq tool [6][7]. 2. Alternative for filtering: To filter paginated results and produce a single array, use a command similar to the following, which isolates the target array from each page before merging them with jq's own --slurp (-s) flag [6][7]: gh api graphql -f query='QUERY' --paginate --jq '.data.path.to.nodes[]' | jq -s 3. Behavior with REST vs. GraphQL: While --slurp helps combine outputs, it simply wraps the pages into an array [1][2]. For REST APIs that return arrays, this results in an array of arrays (e.g., [[...], [...]]) [5]. For more complex GraphQL structures, the nesting of objects and arrays remains, and you may still need to use jq to flatten or transform the output to your desired schema [6].

Citations:

- 1: https://cli.github.com/manual/gh_api
- 2: https://github.com/cli/cli/blob/trunk/pkg/cmd/api/api.go
- 3: https://github.com/cli/cli/pull/8620
- 4: https://github.com/cli/cli/pull/7190
- 5: https://github.com/cli/cli/issues/10459
- 6: https://github.com/cli/cli/issues/1268
- 7: https://github.com/cli/cli/issues/1268?timeline_page=1

---

🏁 Script executed:

```shell
#!/bin/bash
# Description: Find existing gh api --paginate usages and how their output is decoded.
set -euo pipefail

rg -n -C 10 -- '--paginate' --glob '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 258

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- repository gh pagination usages ---'
rg -n -C 12 --glob '*.go' -- '--paginate' .

printf '%s\n' '--- Go module and tool-version declarations ---'
rg -n -C 3 'gh |github.com/cli|go [0-9]|gh version|GitHub CLI' go.mod go.sum .github docs internal 2>/dev/null || true

printf '%s\n' '--- IssueComments and related decoders ---'
rg -n -C 20 'func \(client GHClient\) (IssueComments|ReviewComments)|json\.Unmarshal\(output' internal/reviewsource/coderabbit/coderabbit.go
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- Go source usages only ---'
rg -n --glob '*.go' -- '--paginate' . || true

printf '%s\n' '--- module and tool declarations ---'
for f in go.mod go.sum Makefile .github/workflows/*; do
  if [ -f "$f" ]; then
    rg -n -C 2 'gh|github.com/cli|go [0-9]' "$f" || true
  fi
done

printf '%s\n' '--- target implementation ---'
sed -n '350,425p' internal/reviewsource/coderabbit/coderabbit.go

printf '%s\n' '--- related call sites ---'
rg -n -C 8 'IssueComments\(|Evidence\(|CodeRabbit access|access error' internal/reviewsource internal
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- release-plan pagination implementation ---'
sed -n '60,115p' internal/cli/releaseplan_git_source.go
sed -n '485,515p' internal/cli/releaseplan_command_test.go

printf '%s\n' '--- IssueComments tests and fixtures ---'
sed -n '1215,1270p' internal/reviewsource/coderabbit/coderabbit_test.go
rg -n -C 8 'IssueComments|issues/.*/comments|paginate.*slurp|slurp.*paginate' internal/reviewsource/coderabbit internal/cli

printf '%s\n' '--- installed GitHub CLI information ---'
if command -v gh >/dev/null 2>&1; then
  gh version
  gh api --help | rg -n -C 3 'paginate|slurp'
else
  printf '%s\n' 'gh is not installed'
fi

printf '%s\n' '--- JSON-shape reproduction ---'
python3 - <<'PY'
import json

page1 = '[{"id": 1}]'
page2 = '[{"id": 2}]'

for label, payload in {
    "concatenated pages": page1 + page2,
    "slurped pages": "[" + page1 + "," + page2 + "]",
}.items():
    try:
        value = json.loads(payload)
        print(label, "parses as", type(value).__name__, "shape", value)
        if not all(isinstance(item, dict) for item in value):
            print("  incompatible with []struct{}: elements are not objects")
    except json.JSONDecodeError as exc:
        print(label, "fails:", exc)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 22386

---



</details>

**Handle paginated issue-comment JSON correctly**

`gh api --paginate` emits each page as a separate JSON array. `json.Unmarshal` therefore fails when the endpoint returns multiple pages, and `Evidence` reports an access error instead of classifying a refusal. Do not add `--slurp` alone because it produces an array of arrays. Flatten the slurped pages or decode each page with `json.Decoder`, and add a multi-page regression test.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 375 - 389,
Update GHClient.IssueComments to parse all JSON arrays emitted by gh api
--paginate, decoding each page with json.Decoder or flattening slurped pages
rather than unmarshalling the combined output as one array. Preserve the
existing IssueComment mapping and parse-error behavior, and add a regression
test covering multiple pages so refusal classification does not become an access
error.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:654f87c6764ae833d1474401 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: GitHub CLI’s current manual confirms that `--paginate` fetches every
  page and that `--slurp` would wrap pages rather than flatten page arrays.
  `IssueComments` now decodes each top-level page array with `json.Decoder`,
  preserves page order and mapping, rejects empty or malformed output, and
  retains the contextual parse error.
- Focused evidence: the multi-page regression failed against the former single
  `json.Unmarshal`; the multi-page and malformed-later-page regressions then
  passed with `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/reviewsource/coderabbit -count=1`.
- Daemon Verification: `make verify` not run; Daemon-owned.
