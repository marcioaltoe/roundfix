---
task: task_04
spec: 0192-owned-skills-that-describe-the-product-as-it-is
status: completed
type: docs
complexity: medium
---

# Task 04: The gate, lifecycle and discovery skills agree with the product

## Overview

Eight owned skills each carry a small statement that the Daemon, a command or a Baseline clause contradicts. `qa-gate` tells the QA Agent to rerun a Verification the Daemon already ran and shows a report header the eligibility policy cannot read. `archive-spec` teaches stamping by hand. Five skills show a question form the structured-question clause forbids. This Task corrects each statement where it stands and changes nothing else in those skills. The edits are small and independent; the Verification checks each one by phrase, so every wording this Task quotes is written literally.

## Requirements

1. MUST edit `.agents/skills/qa-gate/SKILL.md`, without touching its `### QA settlement` section:
   - the static-gate step says that in a Daemon-assigned gate the Daemon already ran the repository Verification and gives its result in the prompt, so the Agent must record that result and do not run the repository Verification again; and that a standalone gate runs the repository's selected Verification. The sentence that starts `Run the repository's full verification pipeline` goes;
   - the closing report template's `## Results` table uses the seeded header `| # | Status | Provenance |`, and says each row's story, actor, surface, steps and evidence go in a block below the table. The header `| # | Story / criterion / sweep | Actor and surface | Status | Evidence |` goes;
   - the closing rule says the gate permits Pull Request preparation on `pass`, or on a qualifying declared `partial` as the QA settlement table defines. The words `PR preparation only on` go;
   - both version fields move to `0.0.4`.
2. MUST edit `.agents/skills/archive-spec/SKILL.md`, without touching its `### QA settlement` section:
   - the Steps section says a normal archive also runs through `roundfix archive <slug>`, which verifies the preconditions, stamps the archive metadata and moves the folder, and contains the words `Never hand-edit` about archive front matter. The manual `**Stamp**` and `git mv` steps go;
   - the commit step names the subject `docs: archive <slug>`, the one the Delivery Queue writes. `chore(specs): archive` goes;
   - every mention of a `--release` flag and of stamping a release reference goes, because the command has no such flag;
   - the Unarchive section and every `docs/_inbox/` mention stay unchanged;
   - both version fields move to `0.0.3`.
3. MUST edit `.agents/skills/setup-context-driven/SKILL.md`:
   - keep the exact phrase `` `pending`, `partial`, `deferred`, and `done` `` and add, in the same sentence, the words `plus the terminal statuses` followed by `deprecated`, `superseded`, `closed` and `cancelled`;
   - add to the plan review list an item naming `historyMoves` and the `baseline.history.citation` warnings, reviewed when the plan relocates history;
   - both version fields move to `0.0.3`.
4. MUST remove `[--from task_NN]` from the `argument-hint` of `.agents/skills/implement-spec/SKILL.md` and move both version fields to `0.1.1`.
5. MUST edit the `council` skill:
   - the first paragraph of `.agents/skills/council/references/archetypes.md` says each archetype is defined in that file, that the facilitator builds each advisor's prompt from its definition there, and contains the words `no pre-installed agent file is required`. The words `standalone subagent at` go;
   - Step 7 of `.agents/skills/council/SKILL.md` starts with the words `Ask two questions, one at a time`. The first asks which path the user takes, with the two or three paths the synthesis named as options. The second, asked after the first is answered, asks which trigger reopens the decision, with two or three triggers drawn from the synthesis risks. Each puts the recommended option first with `(Recommended)` in its label and adds no `Other` option. Both answers are recorded verbatim under `## Decision Captured`. The compound question `Which path are you taking, and what triggers would cause you to revisit this decision?` goes;
   - both version fields move to `0.0.3`.
6. MUST edit the `write-idea` skill:
   - the ground rule in `.agents/skills/write-idea/SKILL.md` asks one question with two or three options, the recommended one first and labelled `(Recommended)`, and no longer names `D) Other — describe`;
   - `.agents/skills/write-idea/references/opportunity-scan.md` names `step 7`, the step the skill's Process gives the opportunity scan, instead of `step 6`;
   - both version fields move to `0.0.3`.
7. MUST edit `.agents/skills/business-analyst/SKILL.md`: the Mode 2 block drops the `D) Other — describe` line and labels its first option `(Recommended)`, keeping the one-line rationale; its rule says `2–3 options` per decision instead of `2–4 options`. Both version fields move to `0.0.3`.
8. MUST edit `.agents/skills/brainstorming/SKILL.md`: the "Multiple choice preferred" rule asks for two or three options with the recommended one first and labelled `(Recommended)`, and no longer names `D) Other — describe`. Both version fields move to `0.0.3`.
9. MUST say, wherever requirements 5 to 8 change a question form, that the structured question tool supplies the custom answer itself, and that without such a tool the same single question is shown in chat and a custom answer is accepted.
10. MUST regenerate the mirrors with `make skills-sync`, then run `make baseline-digests`, and name in the Result every path either command rewrote.
11. MUST keep every phrase the existing contract tests require from these skills, and MUST NOT edit any test, any Baseline module, any file under `docs/agents/`, or any skill this Task does not name.

## Subtasks

- [ ] Correct the static gate, the report header and the closing rule in `qa-gate`.
- [ ] Make `archive-spec` archive through the command only.
- [ ] Correct `setup-context-driven` and `implement-spec`.
- [ ] Align the question forms in `council`, `write-idea`, `business-analyst` and `brainstorming`.
- [ ] Correct the two reference files.
- [ ] Raise the eight versions and regenerate the mirrors and digests.

## Acceptance Criteria

- [ ] `qa-gate` tells a Daemon-assigned gate to record the Daemon's Verification result, shows the `| # | Status | Provenance |` header in its closing template, and permits Pull Request preparation on `pass` or a qualifying declared `partial`.
- [ ] `archive-spec` archives only through `roundfix archive <slug>`, forbids hand-editing archive front matter, and names the `docs: archive <slug>` subject.
- [ ] `setup-context-driven` names all eight Finding statuses and the history relocation review items.
- [ ] `implement-spec` names no `--from` argument.
- [ ] `council` captures its decision with two single questions, and its archetype catalog depends on no pre-installed agent file.
- [ ] None of `council`, `write-idea`, `business-analyst` and `brainstorming` shows a `D) Other` option, and each labels its recommended option `(Recommended)`.
- [ ] The `### QA settlement` section is byte-identical in `qa-gate`, `archive-spec` and the Roundfix skill.
- [ ] Every contract test that reads these skills passes unedited, and each mirror is byte-identical to its canonical skill.

## Context

- instruction: `docs/history/specs/0188-a-grant-that-authorizes-delivery-without-governed-paths/qa/qa-report-2026-09-30.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `.agents/skills/implement-spec/SKILL.md`
- interface: `skills/implement-spec/SKILL.md`
- interface: `.agents/skills/council/SKILL.md`
- interface: `skills/council/SKILL.md`
- interface: `.agents/skills/council/references/archetypes.md`
- interface: `skills/council/references/archetypes.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-idea/references/opportunity-scan.md`
- interface: `skills/write-idea/references/opportunity-scan.md`
- interface: `.agents/skills/business-analyst/SKILL.md`
- interface: `skills/business-analyst/SKILL.md`
- interface: `.agents/skills/brainstorming/SKILL.md`
- interface: `skills/brainstorming/SKILL.md`

## Verification

- `for pair in ".agents/skills/qa-gate/SKILL.md|do not run the repository Verification again" ".agents/skills/qa-gate/SKILL.md|or on a qualifying declared" ".agents/skills/archive-spec/SKILL.md|docs: archive <slug>" ".agents/skills/archive-spec/SKILL.md|Never hand-edit" ".agents/skills/setup-context-driven/SKILL.md|plus the terminal statuses" ".agents/skills/setup-context-driven/SKILL.md|historyMoves" ".agents/skills/setup-context-driven/SKILL.md|baseline.history.citation" ".agents/skills/council/SKILL.md|Ask two questions, one at a time" ".agents/skills/council/SKILL.md|(Recommended)" ".agents/skills/council/references/archetypes.md|no pre-installed agent file is required" ".agents/skills/write-idea/SKILL.md|(Recommended)" ".agents/skills/write-idea/references/opportunity-scan.md|step 7" ".agents/skills/business-analyst/SKILL.md|(Recommended)" ".agents/skills/business-analyst/SKILL.md|2–3 options" ".agents/skills/brainstorming/SKILL.md|(Recommended)"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; for pair in ".agents/skills/qa-gate/SKILL.md|Run the repository's full verification pipeline" ".agents/skills/qa-gate/SKILL.md|Story / criterion / sweep" ".agents/skills/qa-gate/SKILL.md|PR preparation only on" ".agents/skills/archive-spec/SKILL.md|chore(specs): archive" ".agents/skills/archive-spec/SKILL.md|**Stamp**" ".agents/skills/archive-spec/SKILL.md|--release" ".agents/skills/implement-spec/SKILL.md|--from task_NN" ".agents/skills/council/SKILL.md|Which path are you taking, and what triggers" ".agents/skills/council/references/archetypes.md|standalone subagent at" ".agents/skills/write-idea/SKILL.md|D) Other" ".agents/skills/write-idea/references/opportunity-scan.md|step 6" ".agents/skills/business-analyst/SKILL.md|D) Other" ".agents/skills/business-analyst/SKILL.md|2–4 options" ".agents/skills/brainstorming/SKILL.md|D) Other"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done; grep -q "^version: 0.0.4$" .agents/skills/qa-gate/SKILL.md && grep -q "^  version: 0.0.4$" .agents/skills/qa-gate/SKILL.md || { printf 'version fields are not 0.0.4 in %s\n' .agents/skills/qa-gate/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.3$" .agents/skills/archive-spec/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/archive-spec/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/archive-spec/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.3$" .agents/skills/setup-context-driven/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/setup-context-driven/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/setup-context-driven/SKILL.md >&2; exit 1; }; grep -q "^version: 0.1.1$" .agents/skills/implement-spec/SKILL.md && grep -q "^  version: 0.1.1$" .agents/skills/implement-spec/SKILL.md || { printf 'version fields are not 0.1.1 in %s\n' .agents/skills/implement-spec/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.3$" .agents/skills/council/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/council/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/council/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.3$" .agents/skills/write-idea/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/write-idea/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/write-idea/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.3$" .agents/skills/business-analyst/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/business-analyst/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/business-analyst/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.3$" .agents/skills/brainstorming/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/brainstorming/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/brainstorming/SKILL.md >&2; exit 1; }; diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null || { printf 'mirror differs: %s\n' skills/qa-gate >&2; exit 1; }; diff -r .agents/skills/archive-spec skills/archive-spec >/dev/null || { printf 'mirror differs: %s\n' skills/archive-spec >&2; exit 1; }; diff -r .agents/skills/setup-context-driven skills/setup-context-driven >/dev/null || { printf 'mirror differs: %s\n' skills/setup-context-driven >&2; exit 1; }; diff -r .agents/skills/implement-spec skills/implement-spec >/dev/null || { printf 'mirror differs: %s\n' skills/implement-spec >&2; exit 1; }; diff -r .agents/skills/council skills/council >/dev/null || { printf 'mirror differs: %s\n' skills/council >&2; exit 1; }; diff -r .agents/skills/write-idea skills/write-idea >/dev/null || { printf 'mirror differs: %s\n' skills/write-idea >&2; exit 1; }; diff -r .agents/skills/business-analyst skills/business-analyst >/dev/null || { printf 'mirror differs: %s\n' skills/business-analyst >&2; exit 1; }; diff -r .agents/skills/brainstorming skills/brainstorming >/dev/null || { printf 'mirror differs: %s\n' skills/brainstorming >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestSettlementGuidanceIsOneTable|TestProjectConstraintQAGate|TestArchiveSpecContract|TestSpecReferenceLifecycleSkillContracts|TestThinSetupSkill|TestBaselineSkillContract|TestToolingAuthorizationJourney|TestAuthorialSkillSync)$" ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSettlementGuidanceIsOneTable TestProjectConstraintQAGate TestArchiveSpecContract TestSpecReferenceLifecycleSkillContracts TestThinSetupSkill TestBaselineSkillContract TestToolingAuthorizationJourney TestAuthorialSkillSync; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -tags docscontract -v -run "^(TestBaselineDocumentationContract)$" ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineDocumentationContract; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task the first required phrase is absent from `qa-gate`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The gate, lifecycle and discovery skills; The structured question form
- `_prd.md` → Goal 1; Core Feature 4; Core Feature 5; Core Feature 7; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Testing Approach 3; Testing Approach 4; Build Order 4

## Result

Implemented the gate, archive, Baseline, implementation, council, idea,
business-analysis, and brainstorming wording corrections in the canonical
`.agents/skills/` tree. The protected `### QA settlement` sections in
`qa-gate`, `archive-spec`, and `roundfix` were left unchanged. The question
forms now use two or three recommended options without an `Other` option and
document that the structured question tool supplies custom answers, with the
same single-question chat fallback accepting a custom answer.

Focused evidence for the acceptance criteria:

- `qa-gate`: contains the Daemon Verification recording rule and standalone
  selected-Verification rule, the seeded `| # | Status | Provenance |` table
  header plus below-table row details, and Pull Request preparation on `pass`
  or qualifying declared `partial`.
- `archive-spec`: normal archives run through `roundfix archive <slug>`,
  archive front matter is never hand-edited, and the Delivery Queue subject is
  `docs: archive <slug>`; the `Unarchive` section and `docs/_inbox/` mentions
  remain unchanged.
- `setup-context-driven` lists all eight Finding statuses and reviews
  `historyMoves` with `baseline.history.citation` warnings when history moves;
  `implement-spec` has no `--from` argument.
- `council` captures two sequential questions with synthesis paths and risk
  triggers, and its archetypes are defined in the catalog without a
  pre-installed agent file. `write-idea`, `business-analyst`, and
  `brainstorming` use the required recommended-option form.
- All eight version pairs were raised as required. The extracted protected QA
  settlement sections compare byte-identically, and every canonical mirror
  compares byte-identically.

Focused checks run:

- `rtk make skills-sync`: passed; rewrote these mirror paths:
  `skills/qa-gate/SKILL.md`, `skills/archive-spec/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`, `skills/implement-spec/SKILL.md`,
  `skills/council/SKILL.md`, `skills/council/references/archetypes.md`,
  `skills/write-idea/SKILL.md`, `skills/write-idea/references/opportunity-scan.md`,
  `skills/business-analyst/SKILL.md`, and `skills/brainstorming/SKILL.md`.
- `rtk make baseline-digests`: passed; no tracked digest paths were rewritten.
- Focused phrase, stale-phrase, version, mirror, protected-settlement, and
  `git diff --check` inspections: passed.

The Task's declared Verification remains for the Daemon; it was not run in
this child-agent turn.

## Carry-forward provenance

- Source Run: `run_20260930T144240Z_33c18e8a522f7217`
- Source commit: `ae33fa0fb2cd35c21b3c8cf80a9b41fa6f21eb79`
