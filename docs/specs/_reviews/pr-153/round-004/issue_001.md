---
source: coderabbit
pr: "153"
round: 4
round_created_at: "2026-08-10T22:48:14Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: 9752569db9792fe76459a9f007a196f92c2c5357
file: internal/speccheck/verification.go
line: 47
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YD7lf,comment:PRRC_kwDOS0qyts7fwUnk
review_hash: 33d2ee0210070bd86c4b15767f8e12007cd2fb91d5ef709362ea4ebc8757b04a
duplicate_of: ""
source_review_id: "4901527003"
source_review_submitted_at: "2026-08-10T22:47:06Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

if sh -c 'printf "" | test -z "$(printf x)"'; then
  echo "unexpected success for non-empty command substitution" >&2
  exit 1
fi

if sh -c "printf '' | test -z '\$(printf x)'"; then
  echo "unexpected success for single-quoted command substitution" >&2
  exit 1
fi

rg -n -C 2 'emptyOutputSucceedsPattern|test\s+-z|\[\s*-z' \
  internal/speccheck/verification.go
```

Repository: marcioaltoe/roundfix

Length of output: 1228

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- verification.go ---'
cat -n internal/speccheck/verification.go | sed -n '1,180p'

printf '%s\n' '--- citations.go context ---'
cat -n internal/speccheck/citations.go | sed -n '1260,1310p'

printf '%s\n' '--- related tests and references ---'
rg -n -C 3 'emptyOutputSucceedsPattern|emptyInputPasses|VacuousVerificationCommands|test -z|\[\s*-z|command substitution|empty-string' \
  internal/speccheck docs --glob '*.go' --glob '*.md' || true

printf '%s\n' '--- regex behavior probe ---'
python3 - <<'PY'
import re

pattern = re.compile(
    r'^(?:rtk\s+)?(?:cat|true|:)\s*$'
    r'|^(?:rtk\s+)?exit\s+0\s*$'
    r'|^(?:rtk\s+)?test\s+-z\s+(?:["\']?\$\(|""|\'\')'
    r'|^(?:rtk\s+)?\[\s*-z\s+(?:["\']?\$\(|""|\'\')'
)
cases = [
    'test -z "$(printf x)"',
    "test -z '$(printf x)'",
    'test -z ""',
    "test -z ''",
    '[ -z "$(printf x)" ]',
    "[ -z '$(printf x)' ]",
    '[ -z "" ]',
    "[ -z '' ]",
    'test -z "$(printf x)" trailing',
]
for case in cases:
    print(f'{case!r}: {bool(pattern.match(case))}')
PY
```

Repository: marcioaltoe/roundfix

Length of output: 44623

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- verification.go remaining logic ---'
cat -n internal/speccheck/verification.go | sed -n '190,240p'

printf '%s\n' '--- verification tests around classification ---'
cat -n internal/speccheck/verification_test.go | sed -n '120,225p'

printf '%s\n' '--- exact classification probe ---'
python3 - <<'PY'
import re

working_tree = re.compile(
    r'git\s+(?:status|diff)\b[^|&;]*(?:--porcelain|--short|--check|--quiet|--exit-code|--name-only|--name-status|--stat)'
)
succeeds = re.compile(
    r'^(?:rtk\s+)?(?:cat|true|:)\s*$'
    r'|^(?:rtk\s+)?exit\s+0\s*$'
    r'|^(?:rtk\s+)?test\s+-z\s+(?:["\']?\$\(|""|\'\')'
    r'|^(?:rtk\s+)?\[\s*-z\s+(?:["\']?\$\(|""|\'\')'
)
fails = re.compile(r'^(?:rtk\s+)?grep\b')

def pipeline_outcome(segment):
    if '|' in segment:
        segment = segment[segment.rfind('|') + 1:]
    segment = segment.strip()
    if fails.search(segment):
        return 'fails'
    if succeeds.search(segment):
        return 'passes'
    return 'unknown'

def chain_passes(command):
    parts = re.split(r'\|\||&&|;', command)
    ops = re.findall(r'\|\||&&|;', command)
    outcome = pipeline_outcome(parts[0])
    for op, part in zip(ops, parts[1:]):
        right = pipeline_outcome(part)
        if op == ';':
            outcome = right
        elif op == '&&':
            if outcome == 'passes':
                outcome = right
            elif outcome != 'fails':
                outcome = 'unknown'
        elif op == '||':
            if outcome == 'fails':
                outcome = right
            elif outcome != 'passes':
                outcome = 'unknown'
    return outcome == 'passes'

cases = [
    'git diff --name-only HEAD | test -z "$(printf x)"',
    "git diff --name-only HEAD | test -z '$(printf x)'",
    'git diff --name-only HEAD | test -z "$(git diff --name-only HEAD)"',
    'git diff --name-only HEAD | test -z ""',
    'git diff --name-only HEAD | [ -z "$(printf x)" ]',
]
for case in cases:
    print(f'{case!r}: working_tree={bool(working_tree.search(case))}, '
          f'pipeline={pipeline_outcome(case)!r}, vacuous={bool(working_tree.search(case) and chain_passes(case))}')
PY
```

Repository: marcioaltoe/roundfix

Length of output: 6671

---



</details>

**Reject arbitrary command substitutions in empty-test matching.**

The matcher classifies non-empty and single-quoted command substitutions as vacuous, although both predicates fail. Parse complete operands, or whitelist substitutions proven empty on the unchanged tree. Add regressions for both cases.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/verification.go` around lines 46 - 47, The empty-test
matcher must reject arbitrary or single-quoted command substitutions instead of
classifying them as vacuous. Update the test patterns shown for test -z and [ -z
to validate complete operands, or only accept substitutions proven empty on the
unchanged tree, while preserving direct empty-string matching. Add regressions
covering non-empty and single-quoted command substitutions.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a4d2ad2de3bdae200a087e3b -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: The empty-output success matcher in `internal/speccheck/verification.go`
  treated any `test -z`/`[ -z` operand that opened with `$(` as a guaranteed
  empty-input success, so both `test -z "$(printf x)"` (non-empty substitution)
  and `test -z '$(printf x)'` (single-quoted literal) were misreported vacuous
  even though both fail. The `$(` form was removed from
  `emptyOutputSucceedsPattern`; a whole-operand parse now applies
  `emptyTestOperandPasses`, which accepts a substitution as vacuous only when
  its inner command reads the working tree (matched by
  `workingTreeStatePattern`). Non-working-tree substitutions, single-quoted
  forms, variables, and literals fall to `emptyInputUnknown`. Added regressions
  in `verification_test.go` for non-empty substitution (not vacuous),
  single-quoted substitution (not vacuous), and working-tree substitution
  (vacuous) in both `test -z` and `[ -z ` bracket forms. Evidence: the focused
  `TestVacuousVerificationCommandIsCaughtBesideHonestSiblings` and
  `TestOneHonestCommand` pass, `go test ./internal/...` (3644) passes, `gofmt`
  is clean. Daemon owns authoritative `make verify`.
