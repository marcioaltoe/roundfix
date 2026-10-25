---
task: task_01
spec: 0316-an-agent-may-start-the-setup-skill
status: completed
type: docs
complexity: medium
---

# Task 01: The setup skill no longer turns off model invocation

## Overview

The Roundfix-owned `setup-context-driven` skill carries
`disable-model-invocation: true`, so an Agent cannot start it and the skill
guide's generic person-only clause sends the person to start it by hand. This
Task removes the line from the canonical skill and its shipped mirror, raises
and records the skill's version, adds a section that keeps the person's
approval at the Plan Digest, adds the same rule to the Baseline guide, and
pins the result with a contract test.

## Requirements

1. MUST answer the maintainer's request in
   [`references/2026-10-24-an-agent-may-start-the-setup-skill.md`](references/2026-10-24-an-agent-may-start-the-setup-skill.md)
   (`_prd.md` → Sources 1) by implementing API Contracts 1-4 of
   `_techspec.md`.
2. MUST, in the canonical `.agents/skills/setup-context-driven/SKILL.md`,
   remove the frontmatter line `disable-model-invocation: true`, raise
   `metadata.version` and `version` from 0.0.5 to 0.0.6, and add the `## Who
   starts this skill` section of `_techspec.md` → Interfaces, byte for byte,
   directly before `## Boundaries`. Every other byte stays. Never touch a
   `### QA settlement` section.
3. MUST run `make skills-sync`, so `skills/setup-context-driven/SKILL.md`
   equals the canonical copy, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   so `skills/testdata/owned-skill-versions.json` records 0.0.6.
4. MUST add `skills/setup_skill_invocation_test.go` with the `Suite:`/
   `Invariant:`/`Boundary` header and `TestAnAgentMayStartTheSetupSkill` of
   `_techspec.md` → Testing Approach, calling `t.Parallel()` first in the
   test and in each subtest: `canonical` and `shipped` (no
   `disable-model-invocation` key in the frontmatter; `## Who starts this
   skill` before `## Boundaries`; "Starting it authorizes no write." and "Ask
   the maintainer to approve the exact current `planDigest`" in the
   whitespace-normalized body) and `person-only-skills-keep-the-lock`
   (`.agents/skills/implement-spec/SKILL.md` keeps
   `disable-model-invocation: true`). Its helpers must not reuse a name the
   package already declares (`skillFrontmatter` exists in `skills/skills.go`).
5. MUST add the paragraph of `_techspec.md` → Interfaces to
   `docs/user-guide/context-driven-development.md`, after the second
   paragraph of "Adopt or update the Context-Driven Baseline".
6. MUST, as the Task's last step, after every guide and skill edit,
   re-record the Behavior Surface fingerprint ("re-record after the last
   guide edit") with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST NOT change any Go production file, any file under
   `internal/baseline/`, `docs/agents/`, any other skill, `CONTEXT.md`, the
   Makefile or the CI workflows.

## Subtasks

- [ ] Remove the line, raise the version and add the section in the canonical skill.
- [ ] Sync the mirror and record the version.
- [ ] Add the contract test.
- [ ] Add the guide paragraph and re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Neither copy of the skill turns off model invocation, and both keep the Plan Digest approval.
- [ ] The skill's 0.0.6 is recorded, the mirror matches, and `implement-spec` keeps its lock.
- [ ] The guide states the rule, and the documentation contracts pass.

## Consumers

None. The skill text, its metadata and one guide paragraph change; no Go package reads the field, and the Baseline modules and rendered guides are unchanged.

## Context

- instruction: `internal/baseline/assets/modules/core.json`
- instruction: `internal/baseline/person_only_skill_test.go`
- instruction: `skills/skills.go`
- instruction: `skills/acceptance_rules_test.go`
- instruction: `.agents/skills/implement-spec/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/setup_skill_invocation_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -timeout 10m -count=1 -v ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAnAgentMayStartTheSetupSkill/canonical TestAnAgentMayStartTheSetupSkill/shipped TestAnAgentMayStartTheSetupSkill/person-only-skills-keep-the-lock TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync/setup-context-driven TestBaselineSkillContract/setup-context-driven; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -timeout 10m -count=1 -run '^(TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill|TestTheSkillGuidesStateThePersonOnlySkillClause)$' ./internal/baseline && go test -timeout 15m -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract && tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'Starting it writes nothing: the Baseline Command still writes only after the person approves the exact Plan Digest'` — expected: exit 0. Before this Task `TestAnAgentMayStartTheSetupSkill` does not exist, so the command fails at its first pass check, and the guide lacks the paragraph. Afterwards both copies pass the new test with `implement-spec` still locked, the version record, mirror and skill contracts pass, the Baseline's person-only clause tests pass unchanged, the Behavior Surface record is current, and the guide holds the rule.

## References

- `_prd.md` → Sources; Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-4; Invariants 1-3; Testing Approach; Build Order 1

## Result

- Neither canonical nor shipped `setup-context-driven` skill disables model invocation; both are version 0.0.6 and contain the Plan Digest approval section. Added `skills/setup_skill_invocation_test.go`, including the unchanged `implement-spec` person-only lock contract.
- Ran `make skills-sync` and recorded the owned-skill version entry with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`; the mirror and version record were updated.
- Added the Context-Driven Baseline guide paragraph and re-recorded `docs/references/behavior-surfaces.json` last with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` (pass).
- Focused check: `go test ./skills -run '^TestAnAgentMayStartTheSetupSkill$' -count=1` (4 passed).
