---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: invalid
terminal_reason: "The broad historical allowlist belongs to an immutable archived Task, while ADR-0081 owns computed regeneration fallout."
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md
line: 76
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v51,comment:PRRC_kwDOS0qyts7eEK7D
review_hash: 15f3b3f898a03978b8c19a897f78984b5ce28f61afac5d8f25faa0b9d76f4cba
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 010: _ Security & Privacy_ _ Major_ _ Quick win_

## Review Comment

_🔒 Security & Privacy_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -eu
for f in \
  docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md \
  docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md
do
  printf '\n--- %s ---\n' "$f"
  sed -n '1,130p' "$f"
done
printf '\n--- references to DERIVED_DIGEST_PATHS and task 0065 ---\n'
grep -RIn --exclude-dir=.git 'DERIVED_DIGEST_PATHS\|0065-loop-order-and-verification-honesty' docs/specs docs/agents internal 2>/dev/null | head -200
```

Repository: marcioaltoe/roundfix

Length of output: 26514

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu
for f in \
  docs/specs/_archived/0065-loop-order-and-verification-honesty/_prd.md \
  docs/specs/_archived/0065-loop-order-and-verification-honesty/_techspec.md \
  docs/adr/0081* \
  Makefile \
  internal/baseline/derived_ownership_test.go
do
  for match in "$f"; do
    [ -e "$match" ] || continue
    printf '\n--- %s ---\n' "$match"
    case "$match" in
      *derived_ownership_test.go) grep -n -A35 -B15 'DERIVED_DIGEST_PATHS' "$match" | head -220 ;;
      Makefile) grep -n -A12 -B8 'DERIVED_DIGEST_PATHS' "$match" ;;
      *) sed -n '1,220p' "$match" ;;
    esac
  done
done
printf '\n--- ownership declarations ---\n'
find internal/baseline -type f -maxdepth 3 -print0 | xargs -0 grep -nH -E 'derived|ownership|assets/(modules|setups|source-baselines|formatter-fixtures|profiles)|testdata' | head -240
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu
printf '%s\n' '--- exact derived paths named by Spec 0065 ---'
grep -RIn --exclude-dir=.git -E \
  'internal/baseline/assets/setups/typescript-bun\.json|internal/baseline/testdata/(catalog\.digest|catalog\.normalized\.json|parity-corpus/v1/fixtures/asset-sync\.json|parity-corpus/v1/manifest\.json)' \
  docs/specs/_archived/0065-loop-order-and-verification-honesty \
  docs/adr/0081-sanctioned-digest-regeneration-is-fallout-of-the-authorized-edit.md

printf '%s\n' '--- task 01 and task 05 authorization/verification sections ---'
sed -n '25,85p' docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md
sed -n '35,85p' docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md

printf '%s\n' '--- QA evidence with changed-path conclusions for Spec 0065 ---'
grep -RIn --include='*.md' -E \
  'changed paths|diff-tree|only.*authorized|DERIVED_DIGEST_PATHS|task_01|task_05' \
  docs/specs/_archived/0065-loop-order-and-verification-honesty/qa \
  | head -240

printf '%s\n' '--- derived ownership records at the five scan roots ---'
find internal/baseline/assets/setups internal/baseline/testdata \
     internal/baseline/assets/source-baselines internal/baseline/assets/formatter-fixtures \
     internal/baseline/assets/profiles \
     -name '_ownership.y*ml' -print | sort
```

Repository: marcioaltoe/roundfix

Length of output: 15330

---



</details>

<!-- cr-reachability -->

**Authorization Bypass (CWE-862):** Missing Authorization

**Reachability:** External · **Exploitability:** Moderate

**Restrict changed-path checks to exact derived files.**

Both task commands allow unrelated files under entire Baseline directories. Allow only the five ADR-0081 fallout files named by the TechSpec.

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md#L73-L76` (this comment)
- `docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md#L74-L77`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md`
around lines 73 - 76, Restrict the changed-path allowlists in task_01.md lines
73-76 and task_05.md lines 74-77 by replacing the broad internal/baseline
directory patterns with only the five exact ADR-0081 fallout file paths named in
the TechSpec; preserve the existing documentation paths and expected exit
behavior in both commands.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/_archived/0065-loop-order-and-verification-honesty/task_01.md</file>
<line_range>73-76</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0065-loop-order-and-verification-honesty/task_05.md</file>
<line_range>74-77</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3323253f8e2aa956b8d5cdab -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The historical command admitted broad Baseline roots, but the Task is
  completed and archived. ADR-0081 makes sanctioned regeneration output a
  computed consequence of the authorized source edit rather than a widened
  per-Spec authorization boundary, and archive policy forbids retrofitting the
  old declaration.
- Daemon Verification: `make verify` not run; Daemon-owned.
