---
source: coderabbit
pr: "125"
round: 1
round_created_at: "2026-08-05T19:38:23Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 1384d928c1d0af4d0bca06e506eb8b2953f9f341
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/make-verify.md
line: 42
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WxMJH,comment:PRRC_kwDOS0qyts7d790U
review_hash: fac4026faa31df54a146b459fb68d68aae3e10344893bc79f1261f2730b6795a
duplicate_of: ""
source_review_id: "4868102567"
source_review_submitted_at: "2026-08-05T19:37:35Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -eu
printf '%s\n' '--- target files ---'
git ls-files -- \
  'docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/make-verify.md' \
  'docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/non-goals.md'
printf '%s\n' '--- make-verify.md ---'
cat -n docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/make-verify.md
printf '%s\n' '--- non-goals.md ---'
cat -n docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/non-goals.md
printf '%s\n' '--- related verification references ---'
rg -n -C 2 'rtk make verify|make verify|authoritative|exited 2|3,452|1\.489615292|TestCheckCorpusBudget' \
  docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078
```

Repository: marcioaltoe/roundfix

Length of output: 9667

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu
python3 - <<'PY'
from pathlib import Path
root = Path("docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078")
make_verify = (root / "make-verify.md").read_text()
non_goals = (root / "non-goals.md").read_text()

checks = {
    "make_verify_records_exit_2": "make verify exited 2" in make_verify,
    "make_verify_records_focused_exit_0": "then\n  exited 0" in make_verify or "then exited 0" in make_verify,
    "non_goals_calls_result_authoritative": "rtk make verify` result" in non_goals and "authoritative" in non_goals,
    "non_goals_identifies_run": any(token in non_goals for token in (
        "2026-08-05", "run", "commit", "exit 0", "exit 2", "passed", "failed"
    )),
}
for name, value in checks.items():
    print(f"{name}={value}")

print("--- conclusion inputs ---")
print("full_gate_failed_later=" + str(checks["make_verify_records_exit_2"]))
print("focused_budget_test_passed=" + str(checks["make_verify_records_focused_exit_0"]))
print("second_file_uses_ambiguous_authoritative_claim=" + str(
    checks["non_goals_calls_result_authoritative"] and not checks["non_goals_identifies_run"]
))
PY
```

Repository: marcioaltoe/roundfix

Length of output: 443

---

🏁 Script executed:

```shell
#!/bin/bash
set -eu
python3 - <<'PY'
from pathlib import Path
import re

root = Path("docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078")
make_text = (root / "make-verify.md").read_text()
non_text = (root / "non-goals.md").read_text()

patterns = {
    "initial_exit_0": r"Initial result:\s*exit\s*`0`",
    "later_full_gate_exit_2": r"exact unpiped command ran again.*?make verify`\s*exited\s*2",
    "later_focused_exit_0": r"minimal unchanged-worktree reproduction.*?exited\s*0",
    "non_goals_mentions_authoritative_gate": r"authoritative\s+`rtk make verify`\s+result",
    "non_goals_has_run_identifier": r"(?:run|attempt|commit|build|initial|final|exit)\s+[`A-Za-z0-9_.:-]+",
}
for name, pattern in patterns.items():
    text = make_text if name != "non_goals_mentions_authoritative_gate" and name != "non_goals_has_run_identifier" else non_text
    print(f"{name}={bool(re.search(pattern, text, re.S | re.I))}")

print("--- non-goals exact statement ---")
for i, line in enumerate(non_text.splitlines(), 1):
    if 14 <= i <= 18:
        print(f"{i}: {line}")

print("--- candidate run identifiers in evidence directory ---")
for path in sorted(root.iterdir()):
    if path.is_file():
        text = path.read_text(errors="replace")
        hits = re.findall(r"(?:Build|commit|run|attempt|Initial result|Final-worktree|exit)\b.{0,100}", text, re.I)
        if hits:
            print(path.name)
            for hit in hits[:5]:
                print("  " + " ".join(hit.split()))
PY
printf '%s\n' '--- related audit status ---'
sed -n '45,60p' docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/constraint-audit.md
printf '%s\n' '--- evidence files ---'
find docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078 -maxdepth 1 -type f -printf '%f\n' | sort
```

Repository: marcioaltoe/roundfix

Length of output: 2911

---



</details>

**Keep the verification status consistent across QA evidence.**

`make-verify.md` records an initial exit `0` and a later full-gate exit `2`. In `non-goals.md`, identify which run is authoritative and link its result. Do not claim a current passing gate until the exact command passes again.

<details>
<summary>📍 Affects 2 files</summary>

- `docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/make-verify.md#L21-L42` (this comment)
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/non-goals.md#L15-L18`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/make-verify.md`
around lines 21 - 42, Keep verification status consistent across both evidence
files: in
docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/make-verify.md
lines 21-42, clearly designate the authoritative run and link its result; in
docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/non-goals.md
lines 15-18, identify that same authoritative result. Do not describe the gate
as currently passing unless the exact unpiped make verify command has passed
again.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/make-verify.md</file>
<line_range>21-42</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/evidence/2026-08-05-spec-0078/non-goals.md</file>
<line_range>15-18</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:18cdca86ffc2d1f6536eaeb2 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Separated the latest exact full-gate attempt (exit 2 on the isolated timing budget) from the earlier passing run, and linked the Non-Goals evidence to that same chronology. The evidence no longer implies that the latest exact gate passed.
- Evidence: `make-verify.md` names the latest attempt and its exit code explicitly; `non-goals.md` cites it and disclaims a current passing-gate claim; `rtk git diff --check` exited 0.
