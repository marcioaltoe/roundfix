---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: internal/cli/doctor_test.go
line: 935
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kaZ,comment:PRRC_kwDOS0qyts7eJyla
review_hash: 72d94874481a23d9a5de2a599c7b8c88d56f85b60fe6d4670d9be1cd37c68e99
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:10Z"
---

# Issue 007: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Assert the unversioned skill name on both CLI surfaces.**

The unversioned case only checks the status text. A regression can omit the owned skill name and still pass this test. Assert that both `doctor` and `skills check` outputs contain `skill` when `state` is `skills.ReadinessUnversioned`.

<details>
<summary>Proposed test addition</summary>

```diff
+			if test.state == skills.ReadinessUnversioned {
+				for surface, output := range map[string]string{
+					"Doctor":       doctorStdout.String(),
+					"skills check": checkOutput,
+				} {
+					if !strings.Contains(output, skill) {
+						t.Errorf("%s output %q does not list unversioned skill %q", surface, output, skill)
+					}
+				}
+			}
```
</details>

As per coding guidelines, test observable behavior and public API contracts.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/doctor_test.go` around lines 930 - 935, Update the unversioned
test case in the doctor and skills-check test tables to assert that both CLI
outputs include the owned skill name “skill” when state is
skills.ReadinessUnversioned, while preserving the existing status and exit-code
assertions.
```

</details>

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3ccb1946a34dc7dd606fc808 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The shared readiness test asserted the unversioned status and exit code but did not require either public surface to identify the affected Skill. The unversioned case now checks both Doctor stdout and combined Skills Check output for `roundfix`. Focused evidence: `rtk go test ./internal/cli -count=1 -run '^TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness$'` with the repository-local Go cache exited 0 and reported four passing subtests. Authoritative `make verify` remains Daemon-owned.
