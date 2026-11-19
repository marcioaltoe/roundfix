---
spec: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
prd: _prd.md
created: 2026-11-18
---

# Adopting the recommendations works under load and gives the researcher the web — Technical Spec

## Executive Summary

Two independent fixes ship together. The proof path shared by `profiles configure` and `profiles check --apply` gains the configured `defaults.profile_proof_setup_timeout`, read by a new narrow reader in `internal/config` that applies `Load`'s precedence and range check to that one key, so the commands keep working on a file another invalid key would stop from loading. The Recommended Authoring Declaration gains `tools: [WebSearch, WebFetch]` for `research`, and the authoring report of Doctor and the Update Check names those tools when a declared research entry has no `tools` key (ADR-0441). The trade-off accepted is a second reader of one key, held equal to `Load` by a test, instead of a full load that would refuse the repairs `profiles configure` exists for.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the reader reads one configuration key and the report reads one declaration entry. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; the web tools are names written into a generated file, never called by Roundfix. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0441 (new) and its amendment of ADR-0415; ADR-0298 restored; ADR-0402 kept; ADR-0181, ADR-0344, ADR-0401 and ADR-0430 listed only to record why they do not apply. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: "Mesmas do BE" of 2026-11-18 for Spec 0411, quoted in the record; it bounds the five Roundfix skill files below and grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/profiles.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/modes.md`.

## System Architecture

Four components change; nothing new is wired.

1. `internal/config` gains `ProfileProofSetupTimeout`, a reader of `defaults.profile_proof_setup_timeout` alone. It reads User Config, then the Project Config at the Git root of the working directory, decodes only `defaults.profile_proof_setup_timeout` with the `durationValue` type `Load` uses, lets Project Config win, and applies the range check of `Validate` with its message. Every other key is ignored.
2. `writeProfilesConfiguration` in `internal/cli/profiles_configure.go`, which both `profiles configure` and `profiles check --apply` call, reads the deadline after `PrepareProfilesConfig` and before the proof, and passes it in the proof configuration's `Defaults`. The proof scope stays the changed categories only (`profilesConfigureProofScope`), and no other configuration value reaches the proof. The configuration it proves is still `roundconfig.Config{Profiles: …}`, now with `Defaults: roundconfig.Defaults{ProfileProofSetupTimeout: …}`.
3. `internal/config/authoring_recommendation.go`: the research entry of the table lists `WebSearch` and `WebFetch`, `RecommendedAuthoringDeclaration` copies the list on each call, and `RecommendedAuthoringYAML` renders it.
4. `internal/cli/authoring_reports.go`: `loadedAuthoringState` also records the recommended research tools when `loaded.Config.Skills.Authoring.Claude.Research` is declared with a nil `Tools`; Doctor's `authoring-models` result and the Update Check's `authoringModels` report them when the inspection finds no problem. The report reads the loaded configuration, not the `roundskills.AuthoringDeclaration` built from it, because that copy turns a nil list into an empty one (characterization, section 6).

The generator in `skills/authoring.go` and the tool validation in `internal/config/authoring.go` do not change (ADR-0402).

Every other proof path was swept (characterization, section 1): `profiles validate`, Doctor, pre-PR review, operational Preflight and Setup already pass the configured value; the fallback candidate probe and the health checker's `Agent` probe have no caller from any command; the Baseline semantic analysis runs from `roundfix baseline`, which reads no Roundfix configuration, and its timeout preserves every Source Baseline Entry. The last three keep their behavior.

## Implementation Design

### Interfaces

```go
// internal/config/proof_setup_deadline.go

// ProfileProofSetupTimeout reads defaults.profile_proof_setup_timeout from
// User Config and Project Config with Load's precedence and range check. It
// reads no other key.
func ProfileProofSetupTimeout(homeDir string, workDir string) (time.Duration, error)
```

```go
// internal/cli/authoring_reports.go

// authoringResearchToolsGap returns the research tools of the Recommended
// Authoring Declaration when the declared research entry has no tools key;
// nil otherwise. An empty list is a declared choice.
func authoringResearchToolsGap(declared *roundconfig.AuthoringHarness) []string
```

`authoringReportState` gains `RecommendedResearchTools []string`, and `updateCheckAuthoring` gains `RecommendedResearchTools []string` with the JSON name `recommendedResearchTools` and `omitempty`, after `recommendation`.

Invariants:

1. The reader returns `Load`'s `Config.Defaults.ProfileProofSetupTimeout` whenever `Load` succeeds on the same User Config and Project Config.
2. The reader fails, with `defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m`, for a value `Load` refuses by range, and with a parse error for a value `Load` refuses as a duration; it never fails for another key. A file that does not parse as YAML fails with `parse config "<path>": <parser error>`, where the parser error names the line, as in `yaml: line 3: mapping values are not allowed in this context`.
3. Every proof request of `profiles configure` and `profiles check --apply` carries the reader's value; when the reader fails, no proof runs and nothing is written.
4. The research tools report appears only when the declared research entry exists, its `Tools` is nil, and the inspection reports no problem.
5. The generated Research Subagent never holds `Bash`, `Write`, `Edit`, `NotebookEdit`, `PowerShell` or `Agent`, whatever the recommendation or a declaration lists.

### Data Models

The Recommended Authoring Declaration's research entry becomes `{model: sonnet, reasoning_effort: medium, tools: [WebSearch, WebFetch]}`. No stored data changes.

### API Contracts

1. API Contract: the proof deadline of a profile write. `profiles configure` and `profiles check --apply` call `roundconfig.ProfileProofSetupTimeout(loadOptions.HomeDir, loadOptions.WorkDir)` in `writeProfilesConfiguration`, after `PrepareProfilesConfig` and before the proof, and pass the result as `ProfileProofSetupTimeout` of the proof configuration's `Defaults`. A reader error is returned through `printProfilesConfigureError` unchanged, before any proof, and nothing is written (Surface Transcripts 1 and 2); the maintainer accepted on 2026-11-19 that a YAML syntax error in the other scope's file now refuses the command, with a message that names the file and the line.
2. API Contract: the guides and the skill on the proof deadline, word for word.
   - In `docs/user-guide/configuration.md`, section `### Profile proof setup deadline`, the sentence that begins "Every profile-readiness consumer passes it" becomes: "Every profile-readiness consumer passes it: `profiles validate`, operational Preflight, pre-PR review, Doctor, Setup, `profiles configure` and `profiles check --apply`." and is followed by: "`profiles configure` and `profiles check --apply` read the key from User Config and Project Config before they prove, with Project Config winning, even when another key of those files would stop loading, and an invalid value refuses the command before any proof."
   - In `docs/user-guide/commands/profiles.md`, the paragraph that ends "adjust `defaults.profile_proof_setup_timeout`." gains: "`profiles configure` and `profiles check --apply` prove under the same key, read from User Config and Project Config with Project Config winning."
   - In `.agents/skills/roundfix/references/profiles.md`, section `## Profile proof setup deadline`, the sentence that begins "`profiles validate`, operational Preflight, pre-PR review, Doctor, Setup and" becomes: "`profiles validate`, operational Preflight, pre-PR review, Doctor, Setup, `profiles configure` and `profiles check --apply` pass the value to each attempt's setup and close context." and is followed by the second sentence of the configuration guide above, word for word.
3. API Contract: the rendered recommendation. `RecommendedAuthoringYAML` ends with `      research:`, `        model: sonnet`, `        reasoning_effort: medium` and `        tools: [WebSearch, WebFetch]`, each on its own line. The JSON `recommendation` of the Update Check and the indented block of Doctor carry the same line.
4. API Contract: the research tools report. With the inspection finding no problem and the declared research entry having no `tools` key:
   - Doctor prints `authoring-models: warn (claude: <declared keys> applied; research declares no tools, so the Research Subagent cannot read the web; next: add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then run roundfix skills install --target project)`, where `<declared keys>` is today's list, such as `prd, research`;
   - the Update Check's text prints, in the position of the other `Authoring models:` lines, `Authoring models: research declares no tools, so the Research Subagent cannot read the web; add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then apply with roundfix skills install --target project.`;
   - its JSON `authoringModels` keeps `state` `applied`, `declared`, empty `problems`, no `next`, and adds `"recommendedResearchTools":["WebSearch","WebFetch"]`.
   The tool list in each text comes from the recommendation. With `tools: []`, a declared list of any tools, or an inspection problem, the output is today's. No exit code and no top-level `next` changes.
5. API Contract: the guides, the skill, the glossary and ADR-0415 on the research tools, word for word.
   - `docs/user-guide/configuration.md`, section `## Authoring models`: the recommended block gains its `tools` line, and the paragraph "The repository can add this option under `research`:" with the YAML block after it is replaced by: "The block gives the Research Subagent the built-in web tools `WebSearch` and `WebFetch` and no MCP tool, because MCP tool names differ from one machine to another. A repository may add MCP tools that read to `research.tools`, or write `tools: []` to keep the Research Subagent off the network. A research entry without a `tools` key makes Doctor and the Update Check name the recommended tools; neither writes the file."
   - `docs/user-guide/commands/doctor.md`: the example block gains its `tools` line; the sentence "The block supplies no `research.tools`; network tools remain the repository's choice." becomes "The block gives the Research Subagent `WebSearch` and `WebFetch`; a repository may add MCP tools that read, or write `tools: []` to keep it off the network."; and right before the paragraph that begins "With no such installed skill" comes the paragraph "When `skills.authoring.claude.research` is declared without a `tools` key and the inspection finds no problem, the line names the recommended research tools:", a `text` block holding the Doctor line of API Contract 4 with `<declared keys>` written as `prd, research`, and the paragraph "The keys before `applied` are the declared ones. An empty list, `tools: []`, is the repository's choice and adds no warning, and a declared list is never compared with the recommendation. Nothing is written and the warning never changes Doctor's exit code."
   - `docs/user-guide/commands/update.md`: "An applied declaration adds no text line." becomes "An applied declaration adds no text line, except the research tools line below."; the example block gains its `tools` line; "The block omits `research.tools`, so it grants no network tools." becomes "The block gives the Research Subagent `WebSearch` and `WebFetch` and no MCP tool."; after the paragraph that begins "Review the block" come the paragraph "When `skills.authoring.claude.research` is declared without a `tools` key and the inspection finds no problem, the check prints this line in the same position:", a `text` block holding the Update Check line of API Contract 4, and the paragraph "The JSON report adds `recommendedResearchTools`, the tools to add, and keeps `authoringModels.state` `applied`. An empty list, `tools: []`, and any declared list add neither. Neither changes an exit code or the top-level `next`, and nothing is written."; in the JSON table, the `authoringModels.recommendation` row's "with no `tools` key" becomes "with `tools: [WebSearch, WebFetch]` under `research`", and a row follows it: "| `authoringModels.recommendedResearchTools` | The research tools of the Recommended Authoring Declaration to add; present only when `applied` and the declared research entry has no `tools` key |".
   - `.agents/skills/roundfix/references/setup.md`: after the sentence that ends "before `roundfix skills install --target project`; neither fails or writes.", insert: "The research entry also accepts `tools`, a list of `WebSearch`, `WebFetch` and MCP tool names beginning with `mcp__` that the Research Subagent receives after `Read, Grep, Glob`. The recommendation lists `WebSearch` and `WebFetch`, `tools: []` keeps the subagent off the network, and when the research entry has no `tools` key Doctor (`authoring-models: warn`) and the Update Check (`authoringModels.recommendedResearchTools`) name the tools to add; neither fails or writes."
   - `.agents/skills/roundfix/references/modes.md`, step 5 of the `--update` mode: after "When it is `undeclared`, show its `recommendation` and leave `.roundfixrc.yml` to the user; never write it." insert "When it holds `recommendedResearchTools`, show the line the check prints and leave `.roundfixrc.yml` to the user as well."
   - `CONTEXT.md`, entry **Recommended Authoring Declaration**: "including the fallbacks of the four Authoring Phases," becomes "including the fallbacks of the four Authoring Phases and the research tools `WebSearch` and `WebFetch`,", and the link `[ADR-0441](docs/adr/0441-the-recommended-authoring-declaration-gives-the-research-subagent-the-built-in-web-tools.md)` follows the ADR-0430 link inside the parentheses.
   - ADR-0415: right after the paragraph that begins "**Amended by ADR-0430 (2026-11-17).**", a paragraph that starts with the bold notice `**Amended by [ADR-0441](0441-the-recommended-authoring-declaration-gives-the-research-subagent-the-built-in-web-tools.md) (2026-11-18).**` and goes on, word for word: "The block lists `tools: [WebSearch, WebFetch]` under `research`, still with no MCP tool, and Doctor and the Update Check name those tools when a declared research entry has no `tools` key. The table, the configuration guide's block and the repository contract hold the tools too; the checks still never write and never fail." The Considered Option that begins "**The repository's `research.tools` in the block.**" gains, at its end: "ADR-0441 takes its built-in part: the block lists `WebSearch` and `WebFetch`, and MCP tools stay out." Its `updated_at` moves to the time of the change.

### Surface Transcripts

1. Surface Transcript: `profiles configure` refuses a configuration whose deadline `Load` would refuse, before any proof. Run in a Git repository whose `.roundfixrc.yml` holds `defaults:` with `profile_proof_setup_timeout: 0s`, with a valid fragment of one `backend` profile and a fresh `HOME`:

```transcript
$ roundfix profiles configure --scope user --file <fragment> --yes
stdout:
stderr:
roundfix: profiles failed: defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m
Run 'roundfix profiles --help' for usage.
exit: 2
```

The lines were copied from the prototype's run (characterization, section 4). On the unchanged tree the same command proves the fragment.

2. Surface Transcript: `profiles configure` refuses a configuration file that does not parse as YAML and names the file and the line of the error, before any proof. Run in a Git repository whose `.roundfixrc.yml` holds the three lines `defaults:`, `  profile_proof_setup_timeout: 90s` and `broken: a: b`, with the fragment of Surface Transcript 1 and a fresh `HOME`:

```transcript
$ roundfix profiles configure --scope user --file <fragment> --yes
stdout:
stderr:
roundfix: profiles failed: parse config "<repository>/.roundfixrc.yml": yaml: line 3: mapping values are not allowed in this context
Run 'roundfix profiles --help' for usage.
exit: 2
```

The line number is the one the YAML parser reports, the same `Load` reports for that file. The lines were copied from a run of the prototype; on the unchanged tree the same command proves the fragment and writes User Config.

## Coverage Map

- Goal 1 → `ProfileProofSetupTimeout`, `writeProfilesConfiguration` (API Contract 1).
- Goal 2 → the recommendation table and `RecommendedAuthoringYAML` (API Contract 3).
- Goal 3 → `authoringResearchToolsGap`, Doctor's `authoring-models` result and the Update Check's `authoringModels` (API Contract 4).
- Goal 4 → API Contracts 2 and 5.
- Core Feature 1 → API Contract 1, Surface Transcripts 1 and 2. Core Feature 2 → API Contract 3 and invariant 5. Core Feature 3 → API Contract 4. Core Feature 4 → API Contracts 2 and 5, ADR-0441.
- Success Metric 1 → `TestProfileAdoptionPassesTheConfiguredSetupDeadline`, `TestTheProofSetupDeadlineReaderAgreesWithLoad`.
- Success Metric 2 → `TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript`, `TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig`, `TestProfilesConfigureProofScope` unchanged.
- Success Metric 3 → `TestRecommendedAuthoringDeclarationRendersValidProjectConfig`, `TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool`, `TestTheRecommendedAuthoringDeclarationMatchesThisRepository/configuration_guide`.
- Success Metric 4 → `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone`, and the QA run of `bin/roundfix update check`.
- Success Metric 5 → `TestTheRecommendedAuthoringDeclarationMatchesThisRepository`.
- Success Metric 6 → the QA outside-evidence row, from the references.

Deviations from the PRD: None.

## Absence Evidence

1. `ProfileProofSetupTimeout` in `internal/config` — command: `rg -n 'func ProfileProofSetupTimeout|profile_proof_setup_timeout' internal/config --glob '!*_test.go'`; result: only the `defaultsOverlay` tag, the `Validate` message and the overlay application in `config.go`; no function reads the key outside `Load`.
2. The deadline in `writeProfilesConfiguration` — command: `rg -n 'Defaults|SetupTimeout|ProfileProofSetupTimeout' internal/cli/profiles_configure.go internal/cli/profiles_check_apply.go`; result: no match.
3. `authoringResearchToolsGap` and the `recommendedResearchTools` member — command: `rg -n 'ResearchTools|recommendedResearchTools|Research\.Tools' internal/cli --glob '!*_test.go'`; result: only `authoring_models.go:36`, which copies the declared tools into the installation's declaration.
4. Research tools in the recommendation — command: `rg -n 'Tools' internal/config/authoring_recommendation.go`; result: no match; the comment at line 17 says "It declares no research tools".
5. The new tests — command: `rg -n 'TestProfileAdoptionPassesTheConfiguredSetupDeadline|TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript|TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig|TestTheProofSetupDeadlineReaderAgreesWithLoad|TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone|TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool' internal skills`; result: no match.

## ADR Alternatives

- ADR-0441: the design takes its decision; its rejected options, web tools in the fixed list, MCP names in the block, comparing a declared list, rewriting the generated text, a snapshot date and writing the file, stay rejected.
- ADR-0415: amended, not superseded; it rejected the repository's `research.tools` in the block because they were MCP tools of one machine, and ADR-0441 takes only the built-in part of that option. Its other rejected options, build-time reading, a generated record, the Model Catalog, a hard-coded block, writing the block, failing the check and a `skills check` warning, stay rejected.
- ADR-0298: the design restores its decision that `profiles configure` passes the configured deadline, and extends the list to `profiles check --apply`, which reuses the same write; its rejected options, proving only the Preferred Selection, raising the built-in deadline, unbounded retries and randomized backoff, stay rejected, and the fallback probe keeps the built-in value as it decided.
- ADR-0402: the design keeps its fixed list and its accepted names; its rejected option "web tools in the fixed list" stays rejected, and a listed web tool reaches the subagent only through the repository's declaration, as it decided.
- ADR-0181: not applicable — it governs the Recommended Profile comparison, which this Spec leaves unchanged; ADR-0441 only borrows its principle that a declared choice is not argued with.
- ADR-0344: not applicable — the close deadline, the Unclosed Proof Session record and the fallback probe's built-in deadline are unchanged; the proofs of the two commands now close under the configured deadline as it already prescribed for every consumer of the key.
- ADR-0401: not applicable — the inspection of the installed files and the generated subagents are unchanged; the research tools report runs only when that inspection finds no problem.
- ADR-0430: not applicable — the fallbacks of the four Authoring Phases are unchanged, and `research` still declares none.

## Integration Points

None external. The reader reads the two configuration files the commands already resolve. The web tools are names in a generated subagent file that the harness resolves; Roundfix never calls them.

## Testing Approach

- `internal/config/proof_setup_deadline_test.go`: `TestTheProofSetupDeadlineReaderAgreesWithLoad` with the subtests `built_in`, `user`, `project`, `project_wins`, `another_key_load_refuses`, `zero`, `above_ten_minutes` and `not_a_duration`, each comparing the reader with `Load` on the same temporary home and repository.
- `internal/cli/profile_proof_setup_deadline_test.go`: `TestProfileAdoptionPassesTheConfiguredSetupDeadline` with the subtests `configure/built_in`, `configure/project_scope_reads_user`, `configure/user_scope_reads_project`, `configure/project_wins`, `check_apply/project` and `check_apply/user`, on `profileReadinessExactRunner`, asserting the `SetupTimeout` of every recorded request; `TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript` asserting Surface Transcript 1 and `TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig` asserting Surface Transcript 2, each with no proof run and no file written. All six subtests fail on the unchanged tree (characterization, section 4). `TestProfilesConfigureProofScope` and `TestProfilesConfigureChangeSummary` stay unchanged and keep passing.
- `internal/cli/authoring_research_tools_test.go`: `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone` with the subtests `no_tools_key`, `empty_tools_list`, `web_tools`, `mcp_tools_only` and `a_problem_keeps_its_report`, built on `writeAuthoringCommandFile`, `installAuthoringReportFixture`, `linkAuthoringReportFixture`, `configureAuthoringDoctorFixture`, `configureAuthoringUpdateFixture`, `runUpdateCheckTest` and `assertAuthoringRecommendationPosition`.
- `skills/authoring_test.go`: `TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool`; `TestTheResearchSubagentNeverHoldsAShellOrAWriteTool` stays unchanged.
- Updated pins: `authoringRecommendationBlock` in `internal/cli/authoring_recommendation_test.go`; `recommendedAuthoringYAML` and the `loads_as_project_config` and `returns_a_copy` subtests in `internal/config/authoring_recommendation_test.go`; `authoringRecommendationMismatch` and a new subtest `a_research_entry_without_a_recommended_tool_fails` in `internal/docscontract/authoring_recommendation_test.go`.
- The sabotages of the QA gate, each run by the author on the prototype: reading Project Config before User Config, passing a zero deadline, removing the range check, dropping the file name from a YAML error, reporting an empty list as missing, never reporting, reporting beside an inspection problem, and a recommendation without `WebFetch`.
- No test reads the clock, the network, the real home or a real Agent.

## Build Order

1. The proof deadline of a profile write: the reader, the wiring in `writeProfilesConfiguration`, their tests, the guide and skill sentences of API Contract 2, the Roundfix skill's version and the Behavior Surface Record (API Contracts 1 and 2, Surface Transcripts 1 and 2).
2. The research tools of the recommendation and their report, with the guides, the skill, the glossary entry, the ADR-0415 amendment, the Roundfix skill's version and the Behavior Surface Record (API Contracts 3 to 5; depends on: 1, because both change `docs/user-guide/configuration.md`, the Roundfix skill's `SKILL.md` in both copies, `skills/testdata/owned-skill-versions.json` and `docs/references/behavior-surfaces.json`).
3. The QA gate (depends on: 1, 2).

Files shared with Specs 0407 and 0408, authored beside this one and delivered before it: `CONTEXT.md` (step 2 edits only the **Recommended Authoring Declaration** entry), `docs/references/behavior-surfaces.json` (re-recorded by steps 1 and 2), `skills/testdata/owned-skill-versions.json` and the Roundfix skill's `SKILL.md` (0408 raises the Roundfix skill version; steps 1 and 2 each raise it again on the tree they run on). No Verification of this Spec depends on another Spec.

## Risks & Considerations

- A second reader of one key could drift from `Load`; `TestTheProofSetupDeadlineReaderAgreesWithLoad` holds them equal on every case `Load` accepts.
- A repository whose `.roundfixrc.yml` has a YAML syntax error now sees `profiles configure` refuse when that file is the other scope's; before, only the target file had to parse. Every other command already refuses such a file, and the maintainer accepted it on 2026-11-19 ("Aceitar, com mensagem clara (Recommended)") with a message that names the file and the line (Surface Transcript 2).
- A repository that deliberately omitted the `tools` key sees a Doctor warning until it writes `tools: []`.
- `WebSearch` is unavailable on providers without the server-side search; the generated file then lists a tool the harness does not grant, and the Claude Code tools reference says such a tool "is never granted, even when listed in `tools`", while the other listed tools stay.
- The Fiscus and Fluxus adoption under real load is declared unreachable for the gate (`_prd.md` → Unreachable Acceptance).

## Glossary

- changes: **Recommended Authoring Declaration**

## Decisions

- One Spec for two defects that share no component, by the maintainer's answer "Juntar à Spec 0411 (Recommended)" (`_prd.md` → Sources, Decisions).
- The proof deadline is read by a narrow reader, not by `Load` or `ResolveConfigProposal`, because both refuse a file `profiles configure` must be able to repair (characterization, section 3).
- `profiles check --apply` uses the same reader through `writeProfilesConfiguration`, although it already loaded the configuration, so one path serves both commands.
- The fallback candidate probe, unreachable since `d7fdd8cf`, and the Baseline semantic analysis, which reads no configuration, keep the built-in deadline; each is a Backlog Entry of 2026-11-19 (`_prd.md` → Non-Goals).
- The recommendation lists the built-in web tools; the generator's fixed list does not change. See ADR-0441.
- Only a research entry without a `tools` key is reported. See ADR-0441.
- The authoring declaration has no snapshot, and `ModelRecommendationSnapshotVersion` stays `2026-11-17`, because no Recommended Profile changes.
