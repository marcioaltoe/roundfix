---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: internal/spec/archive.go
line: 119
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdY,comment:PRRC_kwDOS0qyts7gSdx7
review_hash: f467b69ddb100506fa118cab892621e9f62318bf59a65ac83d55c354c36b8759
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 009: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- tracked candidates ---'
git ls-files | rg '(^|/)(archive|spec|.*test).*\.go$|internal/spec'

printf '%s\n' '--- symbols and call sites ---'
rg -n --glob '*.go' 'ArchiveSpecRoot|ArchiveDir|ArchiveKindSpec|archivedDirName|SpecsRoot|specsRoot' .

printf '%s\n' '--- archive.go outline ---'
ast-grep outline internal/spec/archive.go --view expanded

printf '%s\n' '--- focused archive.go ---'
cat -n internal/spec/archive.go | sed -n '1,180p'
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu

printf '%s\n' '--- archive implementation and tests ---'
cat -n internal/spec/archive.go | sed -n '1,145p'
cat -n internal/spec/archive_test.go | sed -n '280,390p'
cat -n internal/spec/archive_layout_characterization_test.go | sed -n '1,220p'

printf '%s\n' '--- external archive CLI test ---'
cat -n internal/cli/archive_test.go | sed -n '200,245p'

printf '%s\n' '--- resolved SpecsRoot contract ---'
cat -n internal/config/config.go | sed -n '150,180p'
cat -n internal/config/config.go | sed -n '930,995p'

printf '%s\n' '--- audit archive call path ---'
cat -n internal/specaudit/audit.go | sed -n '300,350p'
cat -n internal/specaudit/audit.go | sed -n '1125,1170p'

printf '%s\n' '--- path behavior probe ---'
python3 - <<'PY'
import posixpath

def archive_spec_root(specs_root):
    clean = posixpath.normpath(specs_root)
    docs_root = posixpath.dirname(clean)
    if posixpath.basename(clean) != "specs" or posixpath.basename(docs_root) != "docs":
        return posixpath.join(clean, "_archived")
    return posixpath.join(posixpath.dirname(docs_root), "_archived/specs")

for path in [
    "/repo/docs/specs",
    "/tmp/other/docs/specs",
    "/tmp/other/specs",
    "docs/specs",
]:
    print(f"{path} -> {archive_spec_root(path)}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 27753

---



</details>

**Preserve external Specs Root archive locations.**

When an external `specsRoot` ends in `docs/specs`, this suffix check selects the default archive path. `/tmp/other/docs/specs` therefore resolves to `/tmp/other/_archived/specs` instead of `/tmp/other/docs/specs/_archived`. Pass an explicit default-root decision or the repository root, and add a regression test for this path shape.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/archive.go` around lines 107 - 119, Update ArchiveSpecRoot so
the default archive path is selected only when the configured specsRoot belongs
to the actual repository default root, not merely when its path ends in
docs/specs; pass an explicit default-root decision or repository-root context
from callers. Preserve external roots such as /tmp/other/docs/specs as
/tmp/other/docs/specs/_archived, and add a regression test covering that path
shape.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c97a59c743fb0e28f547b3cd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Fixed the external Spec Root archive-location bug in `internal/spec/archive.go`. `ArchiveSpecRoot(specsRoot string)` used a path-suffix heuristic that misclassified an external root ending in `docs/specs` (e.g. `/tmp/other/docs/specs`) as the repository default, resolving it to `/tmp/other/_archived/specs` instead of `/tmp/other/docs/specs/_archived`. The function now takes an explicit `external bool` decision and returns the beside-root archive `<specsRoot>/_archived` whenever the root is external; the default `_archived/specs` layout is selected only for a non-external built-in root. Threaded the external decision from callers that resolve the configured root:
  - `internal/spec/archive.go` — `ArchiveRequest` carries `External`; `cli/archive.go` passes `resolvedSpecsRoot.External`.
  - `internal/cli/spec_check.go` — `validateSpecAuditSlug` receives `resolvedSpecsRoot.External`.
  - `internal/specaudit/audit.go` — the repo-relative audit helper resolves within its own git tree and keeps `false` (correct there, since it operates on a repo-relative path).
  - Updated call sites in `internal/spec/archive_test.go` and `internal/cli/archive_test.go`.
  Added regression coverage: `TestArchiveSpecRootExternalKeepsArchiveBesideActiveRoot` (including the `/tmp/other/docs/specs` shape → `/tmp/other/docs/specs/_archived`) and `TestArchiveSpecRootDefaultLayout`. Focused evidence: `go test ./internal/spec/... ./internal/cli/... ./internal/specaudit/...` passed; full `make verify` passes.
