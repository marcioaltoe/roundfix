---
task: task_02
spec: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
status: pending
type: backend
complexity: high
---

# Task 02: The recommendation gives the Research Subagent the web, and the checks name the missing tools

## Overview

The Recommended Authoring Declaration lists `tools: [WebSearch, WebFetch]` under `research`, and Doctor and the Update Check name those tools when a declared research entry has no `tools` key, as ADR-0441 decides. The generated subagent's fixed list and the tool validation stay as ADR-0402 decided. The guides, the Roundfix skill's `setup` and `modes` references, the glossary entry **Recommended Authoring Declaration** and ADR-0415 state the rule. It answers the report of 2026-11-18 by conexus, message 01a1f18f, recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md).

## Requirements

1. MUST change `internal/config/authoring_recommendation.go` as `_techspec.md` → Data Models and API Contract 3 say: the research entry of `recommendedAuthoringDeclaration` holds `Tools: []string{"WebSearch", "WebFetch"}`; `RecommendedAuthoringDeclaration` copies that list on each call, so a caller that changes its copy changes neither a later copy nor the rendering; `RecommendedAuthoringYAML` writes `        tools: [WebSearch, WebFetch]` right after the research entry's `reasoning_effort` line, built from the list; and the two doc comments that say the declaration has no research tools say what it now holds.
2. MUST implement `_techspec.md` → API Contract 4 in `internal/cli/authoring_reports.go`, with `authoringResearchToolsGap` as `_techspec.md` → Interfaces gives it and invariant 4 bounds it:
   - `authoringReportState` and `updateCheckAuthoring` gain `RecommendedResearchTools []string`, the latter with the JSON tag `json:"recommendedResearchTools,omitempty"` after `Recommendation`;
   - `loadedAuthoringState` sets it from `authoringResearchToolsGap(loaded.Config.Skills.Authoring.Claude)`, reading the loaded configuration and never the `roundskills.AuthoringDeclaration` built from it;
   - Doctor's `authoring-models` result, when the inspection has no problem and the field is set, is `warn` with the detail `<today's applied detail>; research declares no tools, so the Research Subagent cannot read the web` and the next action `add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then run roundfix skills install --target project`;
   - the Update Check, in the same case, keeps `state` `applied`, sets `recommendedResearchTools`, and prints `Authoring models: research declares no tools, so the Research Subagent cannot read the web; add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then apply with roundfix skills install --target project.`;
   - the tool list in every text is joined with `, ` from `RecommendedAuthoringDeclaration().Research.Tools`; with an inspection problem, an empty `tools` list or any declared list, every report stays as it is today; no exit code and no top-level `next` changes.
3. MUST add `internal/cli/authoring_research_tools_test.go` with `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone`, calling `t.Parallel()` first, built on `withCLIWorkspace`, `writeAuthoringCommandFile`, `installAuthoringReportFixture`, `linkAuthoringReportFixture`, `configureAuthoringDoctorFixture`, `configureAuthoringUpdateFixture`, `runUpdateCheckTest` and `assertAuthoringRecommendationPosition`. Each subtest writes a Project Config declaring `prd: {model: opus, reasoning_effort: xhigh}` and the research entry below, installs, links, runs `doctor`, `update check --format json` and `update check`, and asserts the exact `authoring-models:` line after `skills:` and before `baseline-rules:`, the exact `authoringModels` JSON, a top-level `next` of `"none"`, and the exact text line or its absence:
   - `no_tools_key`, `research: {model: sonnet}`: the Doctor line and the text line of requirement 2, and `{"state":"applied","harness":"claude","declared":["prd","research"],"problems":[],"recommendedResearchTools":["WebSearch","WebFetch"]}`;
   - `empty_tools_list`, `{model: sonnet, tools: []}`; `web_tools`, `{model: sonnet, tools: [WebSearch, WebFetch]}`; `mcp_tools_only`, `{model: sonnet, tools: [mcp__kb__read]}`: `authoring-models: ok (claude: prd, research applied)`, `{"state":"applied","harness":"claude","declared":["prd","research"],"problems":[]}` and no `Authoring models:` line;
   - `a_problem_keeps_its_report`, `{model: sonnet}` plus a leftover generated `.claude/agents/roundfix-tasks-author.md`: today's problem report in both checks, `{"state":"unapplied","harness":"claude","declared":["prd","research"],"problems":[".claude/agents/roundfix-tasks-author.md has no declared entry"],"next":"roundfix skills install --target project"}`, and no research tools text.
4. MUST add `TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool` to `skills/authoring_test.go`, calling `t.Parallel()` first: it loads `config.RecommendedAuthoringYAML()` as Project Config with `config.ResolveConfigProposal`, generates the research subagent with `GeneratedSubagent`, and asserts its `tools:` line is exactly `tools: Read, Grep, Glob, WebSearch, WebFetch` and names none of `Bash`, `Write`, `Edit`, `NotebookEdit`, `PowerShell` and `Agent` (invariant 5). `TestTheResearchSubagentNeverHoldsAShellOrAWriteTool` stays unchanged.
5. MUST update every test that pins the recommendation on the default branch when this Task runs, among them:
   - `recommendedAuthoringYAML` in `internal/config/authoring_recommendation_test.go` gains the `tools` line; its `loads_as_project_config` subtest asserts `Research.Tools` equals `[WebSearch WebFetch]`; its `returns_a_copy` subtest also changes an element of the copy's tools and appends to them, and still finds the declaration and the rendering unchanged;
   - `authoringRecommendationBlock` in `internal/cli/authoring_recommendation_test.go` gains `          tools: [WebSearch, WebFetch]`, which `TestRunDoctorRealRepositoryCheckDoesNotMutateState`, the two `…RecommendsAnAuthoringDeclarationWhenNoneIsDeclared` tests and the two Surface Transcript tests of that file then read;
   - `authoringRecommendationMismatch` in `internal/docscontract/authoring_recommendation_test.go` also fails, naming `skills.authoring.claude.research.tools`, when the repository's research tools lack a recommended tool, and passes further tools; a new subtest `a_research_entry_without_a_recommended_tool_fails` of `TestTheRecommendedAuthoringDeclarationMatchesThisRepository` asserts that for the lists nil, empty, `[WebSearch]` and `[WebFetch, mcp__kb__read]`, and that the recommendation plus `mcp__kb__read` passes.
6. MUST make the changes of `_techspec.md` → API Contract 5 word for word in `docs/user-guide/configuration.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/update.md`, `.agents/skills/roundfix/references/setup.md` and `.agents/skills/roundfix/references/modes.md`, carrying forward every other sentence of the paragraphs it edits; in particular step 5 of the `--update` mode keeps "never commits, pushes or opens a pull request" and its `authoringModels.state` sentences.
7. MUST, through `domain-modeling`, revise the glossary entry **Recommended Authoring Declaration** in `CONTEXT.md` as `_techspec.md` → API Contract 5 says, with the ADR-0441 link after the ADR-0430 link. No other glossary entry changes.
8. MUST write the amendment of `_techspec.md` → API Contract 5 into `docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md`: the paragraph right after the paragraph that begins "**Amended by ADR-0430 (2026-11-17).**", the sentence at the end of the Considered Option that begins "**The repository's `research.tools` in the block.**", and its `updated_at` moved to the time of the change.
9. MUST regenerate, in this order and with no hand edit of a mirror, version field or record:
   - the mirrors `skills/roundfix/references/setup.md` and `skills/roundfix/references/modes.md`, with `make skills-sync`;
   - the Roundfix skill's version, with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` to the next version and records it in `skills/testdata/owned-skill-versions.json`;
   - `make skills-sync` again, then `make skills-sync-check`;
   - as the last step, after the last guide edit, the Behavior Surface Record, with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, which rewrites `docs/references/behavior-surfaces.json`.
   task_01 raised the same version fields and re-recorded the same record; each command runs on the tree this Task starts from.
10. MUST NOT change `skills/authoring.go`, `internal/config/authoring.go`, `internal/cli/authoring_models.go`, `.roundfixrc.yml`, `internal/config/recommendations.go`, any other skill or skill reference, any Governed Path outside `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, or the `### QA settlement` section of any skill.

## Subtasks

- [ ] List the research tools in the recommendation and update the tests that pin it.
- [ ] Report a research entry without tools in Doctor and the Update Check, with its test.
- [ ] Add the generated-subagent test and the repository contract on the research tools.
- [ ] Edit the three guides, the canonical `setup` and `modes` references, the glossary entry and ADR-0415.
- [ ] Synchronize the mirrors, record the skill version, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The rendered recommendation ends its research entry with `tools: [WebSearch, WebFetch]`, and the subagent generated from it holds those two tools and no shell or write tool.
- [ ] A research entry without a `tools` key draws the exact Doctor warning, Update Check line and JSON member, and an empty list, a declared list or an inspection problem draws none of them.
- [ ] This repository's declaration lists every recommended research tool, and the configuration guide's block equals the rendering.
- [ ] The guides, the skill references, the glossary entry and ADR-0415 say what the code does; the mirrors equal the canonical references, the Roundfix skill's version is recorded, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0441-the-recommended-authoring-declaration-gives-the-research-subagent-the-built-in-web-tools.md`
- instruction: `docs/adr/0402-the-research-subagent-reads-and-never-runs-a-command-or-writes-a-file.md`
- instruction: `skills/authoring.go`
- interface: `internal/config/authoring_recommendation.go`
- interface: `internal/config/authoring_recommendation_test.go`
- interface: `internal/cli/authoring_reports.go`
- creates: `internal/cli/authoring_research_tools_test.go`
- interface: `internal/cli/authoring_recommendation_test.go`
- interface: `internal/docscontract/authoring_recommendation_test.go`
- interface: `skills/authoring_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Recommended Authoring Declaration and its rendering, and the `authoringModels` report of the Update Check with Doctor's `authoring-models` line — consumers: `./internal/config`, `./internal/cli`, `./skills`, `./internal/docscontract`

## Verification

- `grep -q '^func TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone(' internal/cli/authoring_research_tools_test.go && grep -q '^func TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool(' skills/authoring_test.go && grep -q 'recommendedResearchTools' internal/cli/authoring_reports.go && go test -count=1 -timeout 15m ./internal/config && go test -count=1 -timeout 30m ./skills && go test -count=1 -timeout 30m -run '^(TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone|Test.*Authoring.*|TestRunDoctorRealRepositoryCheckDoesNotMutateState|TestSurfaceTranscript.*|TestDoctor.*|TestUpdateCheck.*)$' ./internal/cli` — expected: exit 0. Before this Task the two tests and the JSON member do not exist, so the command fails at the first `grep`; afterwards the whole `internal/config` and `skills` packages and every authoring, Doctor and Update Check test of `internal/cli` pass.
- `for file in docs/user-guide/commands/doctor.md docs/user-guide/commands/update.md; do text="$(tr -s '[:space:]' ' ' < "$file")"; for phrase in 'research declares no tools, so the Research Subagent cannot read the web' 'add tools: \[WebSearch, WebFetch\] under skills.authoring.claude.research in .roundfixrc.yml' 'The block gives the Research Subagent .WebSearch. and .WebFetch.' 'tools: \[WebSearch, WebFetch\]'; do printf '%s\n' "$text" | grep -qE -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qE -- 'authoringModels.recommendedResearchTools. [|] The research tools of the Recommended Authoring Declaration to add' || { echo 'update guide lacks the JSON row' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qE -- 'or write .tools: \[\]. to keep the Research Subagent off the network' || { echo 'configuration guide lacks the paragraph' >&2; exit 1; }; for file in .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$file" | grep -qE -- 'and the Update Check \(.authoringModels.recommendedResearchTools.\) name the tools to add' || { printf 'missing sentence in %s\n' "$file" >&2; exit 1; }; done; for file in .agents/skills/roundfix/references/modes.md skills/roundfix/references/modes.md; do tr -s '[:space:]' ' ' < "$file" | grep -qE -- 'When it holds .recommendedResearchTools., show the line the check prints' || { printf 'missing sentence in %s\n' "$file" >&2; exit 1; }; done; grep -F -- '**Recommended Authoring Declaration**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Recommended Authoring Declaration\*\*(:| —( |$))' || { echo 'glossary lacks Recommended Authoring Declaration' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qE -- 'the research tools .WebSearch. and .WebFetch.' || { echo 'glossary entry lacks the research tools' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md | grep -qF -- 'Amended by [ADR-0441](0441-the-recommended-authoring-declaration-gives-the-research-subagent-the-built-in-web-tools.md) (2026-11-18).' || { echo 'ADR-0415 lacks the amendment' >&2; exit 1; }; make skills-sync-check && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document holds the sentences, so the command fails; afterwards the guides, the references and their mirrors, the glossary entry and ADR-0415 agree with the code, and the documentation contract suite passes with the configuration guide's block equal to the rendering and the Behavior Surface Record current.

## References

- `_prd.md` → Goals (second to fourth bullets); Core Features 2 to 4; Success Metric 3; Success Metric 4; Success Metric 5; Non-Goals; Glossary
- `_techspec.md` → System Architecture 3 and 4; Interfaces; invariants 4 and 5; Data Models; API Contract 3; API Contract 4; API Contract 5; Testing Approach; Build Order 2; Glossary
- [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md) sections 5 to 8
- ADR-0441, ADR-0415, ADR-0402
