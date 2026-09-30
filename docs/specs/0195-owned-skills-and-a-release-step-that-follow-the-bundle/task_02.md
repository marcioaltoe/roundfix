---
task: task_02
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
status: pending
type: test
complexity: medium
---

# Task 02: An owned skill's content cannot change under one version

## Overview

Readiness compares versions, so a version must name one content. Today nothing stops an owned skill from changing while its version stays the same. This Task adds a record of every version each owned skill has shipped, with the digest of that version's folder, and a test that refuses a changed skill under a recorded version and a version that is not recorded. It also states the rule in the repository's own guide, where a Spec author reads it before declaring a skill edit.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST add `skills/owned_skill_versions_test.go` with the four tests the TechSpec's Testing Approach 2 names, implementing "The record rule" exactly: its four cases, and recording only through the flag `-record-skill-versions`.
2. MUST compute each digest with the package's existing skill folder hash over the embedded folder, so the canonical copy and the mirror give the same value.
3. MUST NOT reuse the package's `-update` flag for recording, and MUST NOT give any of the four tests a name that starts with `TestAuthorialSkillSync`. `make baseline-digests` must never record a version.
4. MUST create `skills/testdata/owned-skill-versions.json` in the form the TechSpec's Data Models section gives, by running the recording command once. It holds, for each of the 14 owned skills, the version its embedded `SKILL.md` declares on this Task's base and that folder's digest.
5. MUST run the three negative tests on an in-memory record, so they change no file.
6. MUST add to `docs/agents/specific-repository.md` the rule the TechSpec's "Fixed texts" gives, and change nothing else in that file.
7. MUST NOT edit the Makefile, any skill, or any existing test, and MUST NOT embed the record in the binary.

## Subtasks

- [ ] Add the record test with its recording flag.
- [ ] Add the three negative tests on an in-memory record.
- [ ] Record the 14 current versions.
- [ ] State the rule in the repository guide.

## Acceptance Criteria

- [ ] The record lists all 14 owned skills with the version and digest they ship.
- [ ] A skill whose digest differs from the recorded digest for its version is refused, and recording does not replace that digest.
- [ ] A version that is not recorded is refused, and the failure names the recording command.
- [ ] `make baseline-digests` leaves the record byte-identical.

## Context

- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `skills/skills.go`
- instruction: `skills/baseline_skill_contract_test.go`
- interface: `docs/agents/specific-repository.md`
- creates: `skills/owned_skill_versions_test.go`
- creates: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestAChangedOwnedSkillUnderARecordedVersionIsRefused|TestAnUnrecordedOwnedSkillVersionIsRefused|TestRecordingNeverReplacesARecordedVersion)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAChangedOwnedSkillUnderARecordedVersionIsRefused TestAnUnrecordedOwnedSkillVersionIsRefused TestRecordingNeverReplacesARecordedVersion; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for name in roundfix write-idea write-prd write-techspec write-tasks setup-context-driven implement-task implement-spec brainstorming council business-analyst archive-spec qa-gate evidence-gate; do grep -q -- "\"$name\": \[" skills/testdata/owned-skill-versions.json || { printf 'not recorded: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/specific-repository.md|content changes only together with its version" "docs/agents/specific-repository.md|-record-skill-versions"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four named tests exists and the record file is absent, so the command fails.

## References

- [_techspec.md](_techspec.md) — Data Models; The record rule; Fixed texts; Testing Approach 2
- `_prd.md` → Goal 2; Core Feature 2; Success Metric 3
- `_techspec.md` → API Contract 3
- [references/2026-09-30-an-owned-skill-older-than-the-bundle-passes-readiness.md](references/2026-09-30-an-owned-skill-older-than-the-bundle-passes-readiness.md)
- ADR-0189

## Result
