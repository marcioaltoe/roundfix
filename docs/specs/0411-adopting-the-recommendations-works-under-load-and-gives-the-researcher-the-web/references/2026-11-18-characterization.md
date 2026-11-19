# Characterization for Spec 0411

Measured by the author on 2026-11-18, about 23:30 to 00:40 BRT, on the planning clone of the branch `docs/plan-0411` based on `a6a4a15b` (v0.148.0 plus the plans of 0407 and 0408), and on a prototype in a scratch worktree of that commit, removed afterwards. Every test ran with `GOMAXPROCS=2`, a fake Agent runner and no network.

## 1. Every path that proves or probes an Agent Selection

The sweep read every non-test construction of `agent.ProbeRequest`, every call of `proveProfileSelections` and `proveProfileSelectionsWithOptions`, every `roundconfig.Config{` literal under `internal/cli`, and every use of `acpxPreflightSetupTimeout` (`rg -n 'ProbeRequest\{|proveProfileSelections(WithOptions)?\(|roundconfig\.Config\{|acpxPreflightSetupTimeout' internal --glob '!*_test.go'`).

| # | Path | Deadline it uses | Reached from a command |
| - | ---- | ---------------- | ---------------------- |
| 1 | `profiles validate` (`internal/cli/profiles_validate.go:144`, proof at `:220-224`) | `config.Defaults.ProfileProofSetupTimeout` of the loaded configuration | yes |
| 2 | Doctor (`internal/cli/doctor.go:54`) | the loaded configuration | yes |
| 3 | pre-PR review (`internal/cli/review.go:1637`) | the loaded configuration | yes |
| 4 | operational Preflight (`internal/cli/profile_preflight.go:62`) | the loaded configuration | yes |
| 5 | Setup (`internal/cli/setup.go:313`) | the proposed configuration that `ResolveConfigProposal` resolves (`setup.go:270`) | yes |
| 6 | `profiles configure` (`internal/cli/profiles_configure.go:90-96`) | `roundconfig.Config{Profiles: proofProfiles}`: `Defaults` is zero, so `SetupTimeout` is 0 and the runner falls back to the built-in 30 s (`internal/agent/selection_assignment.go:118-121`) | yes: **defect** |
| 7 | `profiles check --apply` (`internal/cli/profiles_check_apply.go:69`) | the same `writeProfilesConfiguration` as row 6 | yes: **defect** |
| 8 | the fallback candidate probe (`internal/agent/fallback.go:58,64,96`) | the constant `acpxPreflightSetupTimeout` | no: `ProbeFallback` is called only by `probeRuntimeSelection` (`internal/cli/selection.go:94-113`), and no code calls `probeRuntimeSelection`; its last caller left in `d7fdd8cf` (#32) |
| 9 | Baseline semantic analysis (`internal/baselineacp/analyzer.go:330`) | no `SetupTimeout`, so the built-in 30 s | yes, from `roundfix baseline`, which reads no Roundfix configuration; an unavailable analysis preserves every Source Baseline Entry (ADR-0344, "Baseline semantic analysis keeps its stop") |
| 10 | the health checker's `Agent` probe (`internal/cli/health.go:186`) | the request it is given | no: `runtimeHealthChecker.Agent` has no non-test caller |

Rows 6 and 7 are the defect the adopters met. Row 8 is unreachable, so it changes nothing for a user; this Spec leaves it alone. Row 9 runs on a repository whose configuration the `baseline` command never loads; making it read the key would add a configuration load and a new way for `baseline` to fail, so this Spec leaves it alone and the report names it.

## 2. The configured deadline never reaches `profiles configure` or `profiles check --apply`

A test with the fake exact runner `profileReadinessExactRunner` records each `ProbeRequest`. On the unchanged tree every request of `profiles configure` and `profiles check --apply` carries `SetupTimeout` 0, whatever User Config or Project Config declares; `profiles validate` on the same configuration carries the configured value (`TestProfileProofPassesTheConfiguredSetupDeadline`). This matches the field reports: Fiscus measured each Codex proof at about 45 s under a load of 13 to 23 and saw the adoption fail while `profiles validate` passed, and Fluxus raised the key to `3m` and still saw `profiles check --apply` time out (`references/2026-11-18-maintainer-answers.md`).

## 3. Reading the whole configuration would break `profiles configure`

A first prototype resolved the configuration the command leaves behind with `ResolveConfigProposal`, as Setup does. Two existing tests then failed:

- `TestProfilesConfigureProofScope`: the Project Config holds an untouched category whose runtime `stale-adapter` is invalid, and `profiles configure` must still rewrite another category; the full resolution refused with `profiles.frontend.preferred.runtime "stale-adapter" is invalid`.
- `TestProfilesConfigureChangeSummary/declared_removals_include_present_and_absent_categories`: removing every category leaves no profile, and the full resolution refused with `profiles must define at least one Agent Selection Profile`.

`profiles configure` is the command that repairs a configuration `Load` refuses, so it must not depend on the rest of the file loading. The second prototype reads only `defaults.profile_proof_setup_timeout`, from User Config and then Project Config, with `Load`'s parser for the duration and its range check. With it both tests pass unchanged.

## 4. The second prototype

- `internal/config/proof_setup_deadline.go`: `ProfileProofSetupTimeout(homeDir, workDir string) (time.Duration, error)`.
- `internal/cli/profiles_configure.go`: `writeProfilesConfiguration` calls it and passes the value in the proof configuration's `Defaults`; an error refuses before any proof.
- New tests and results:
  - `TestProfileAdoptionPassesTheConfiguredSetupDeadline`: 6 subtests (`configure/built_in`, `configure/project_scope_reads_user`, `configure/user_scope_reads_project`, `configure/project_wins`, `check_apply/project`, `check_apply/user`); all 6 fail on the unchanged tree and pass with the prototype.
  - `TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript`: with `defaults.profile_proof_setup_timeout: 0s` in Project Config and a user-scope write, the prototype prints nothing on standard output, prints `roundfix: profiles failed: defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m` and `Run 'roundfix profiles --help' for usage.` on standard error, exits 2, proves nothing and writes nothing. The unchanged tree proves and writes.
  - `TestTheProofSetupDeadlineReaderAgreesWithLoad`: 8 subtests; the reader returns `Load`'s value wherever `Load` succeeds, keeps the configured value where only another key stops `Load`, and fails where `Load` fails on the key.
- Sabotages, each run alone and restored:
  - reading Project Config before User Config fails `TestTheProofSetupDeadlineReaderAgreesWithLoad/project_wins` and `TestProfileAdoptionPassesTheConfiguredSetupDeadline/configure/project_wins`;
  - passing a zero deadline (`0 * setupTimeout`) fails all 6 subtests of `TestProfileAdoptionPassesTheConfiguredSetupDeadline`;
  - removing the range check fails `TestTheProofSetupDeadlineReaderAgreesWithLoad/zero` and `.../above_ten_minutes`.
- `go test -run 'Profile|Recommend|Setup' ./internal/cli ./internal/config` passed with the prototype.

## 5. The Research Subagent the recommendation produces cannot read the web

- `internal/config/authoring_recommendation.go:13` declares `research` with a model and an effort and no `tools`; its comment at `:17` says "It declares no research tools", and `RecommendedAuthoringYAML` renders the block "without a tools key" (`:30`).
- `skills/authoring.go:163-165` writes `tools: Read, Grep, Glob` followed by the declared `research.tools`, while the generated description and body say the subagent reads "the codebase, the knowledge base and the web" and returns URLs (`:157`, `:168-172`).
- The key exists and is validated (`internal/config/authoring.go:135-152`: `WebSearch`, `WebFetch` or `mcp__…`, no duplicate), and the configuration guide documents it (`docs/user-guide/configuration.md:146-151`, `:182-187`). Doctor and the Update Check never mention it: `docs/user-guide/commands/doctor.md:227-228` says "The block supplies no `research.tools`" and `docs/user-guide/commands/update.md:148-149` says "The block omits `research.tools`, so it grants no network tools".
- An adopter that pasted the block sees `authoring-models: ok` and `authoringModels.state` `applied` (`TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/declared`), so nothing tells it the subagent has no web tool. Conexus found it only through a pre-PR review (`references/2026-11-18-maintainer-answers.md`, "Segunda fonte").
- Published evidence: the Claude Code tools reference, <https://code.claude.com/docs/en/tools>, read on 2026-11-18, says of a subagent definition: "`tools` only: the subagent gets only the listed tools." The subagent guide, <https://code.claude.com/docs/en/sub-agents>, says the field lists the "Tools the subagent can use" and "Inherits every tool available to subagents if omitted", and names `WebFetch` and `WebSearch` among the built-in tools a background subagent keeps. A generated file whose list is `Read, Grep, Glob` therefore cannot search or fetch the web.

## 6. An empty tool list stays distinct from an absent key

A throwaway test package, removed afterwards, resolved two Project Configs with `ResolveConfigProposal`: `research: {model: sonnet}` loads with `Research.Tools == nil`, and `research: {model: sonnet, tools: []}` loads with a non-nil empty list. The loaded configuration therefore tells "declared no tools" from "never chose", and the report can stay silent for an empty list. `authoringDeclaration` in `internal/cli/authoring_models.go:36` copies the list with `append([]string{}, …)`, which loses the difference, so the report reads `loaded.Config.Skills.Authoring.Claude` directly.

## 7. The research prototype

- `RecommendedAuthoringDeclaration` lists `Tools: []string{"WebSearch", "WebFetch"}` for `research`, copied on each call; `RecommendedAuthoringYAML` adds `        tools: [WebSearch, WebFetch]` after the research effort.
- `internal/cli/authoring_reports.go` computes the missing recommendation when the declared research entry has no `tools` key and the inspection finds no problem: Doctor prints `authoring-models: warn (claude: prd, research applied; research declares no tools, so the Research Subagent cannot read the web; next: add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then run roundfix skills install --target project)`; the Update Check prints `Authoring models: research declares no tools, so the Research Subagent cannot read the web; add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then apply with roundfix skills install --target project.` and adds `"recommendedResearchTools":["WebSearch","WebFetch"]` to `authoringModels`, whose `state` stays `applied`.
- Tests that pin the old block and had to change: `authoringRecommendationBlock` in `internal/cli/authoring_recommendation_test.go` (consumed by `TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared`, `TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared`, `TestSurfaceTranscript1UpdateCheckRecommendsTheDeclaration`, `TestSurfaceTranscript9UpdateCheckRecommendsTheFallbacks` and `TestRunDoctorRealRepositoryCheckDoesNotMutateState` in `internal/cli/doctor_test.go`, which needs no edit of its own), `recommendedAuthoringYAML` and the `loads_as_project_config` and `returns_a_copy` subtests in `internal/config/authoring_recommendation_test.go`, and the guide block that `TestTheRecommendedAuthoringDeclarationMatchesThisRepository/configuration_guide` compares. Each failed before its update and passed after it. `TestTheSkillCoverageMapIsCurrent` failed after the guide edit until the Behavior Surface Record is re-recorded, as expected.
- New tests and results: `TestTheChecksNameTheRecommendedResearchToolsWhenResearchDeclaresNone` (subtests `no_tools_key`, `empty_tools_list`, `web_tools`, `mcp_tools_only`, `a_problem_keeps_its_report`) and `TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool` pass with the prototype; the repository contract `authoringRecommendationMismatch` extended with the research tools passes on this repository's `.roundfixrc.yml`, which lists `WebSearch, WebFetch` and two MCP tools.
- Sabotages: reporting an empty list as missing fails `.../empty_tools_list`; never reporting fails `.../no_tools_key`; adding the tools to the Update Check report beside an inspection problem fails `.../a_problem_keeps_its_report`; a recommendation that lists only `WebSearch` fails `TestTheRecommendedResearchSubagentReadsTheWebAndHoldsNoShellOrWriteTool`.
- `go test -run 'Authoring|Research|Recommend|Doctor|UpdateCheck|SurfaceTranscript' ./internal/cli ./internal/config ./skills` passed with the prototype, and `-tags docscontract ./internal/docscontract` passed except `TestTheSkillCoverageMapIsCurrent`, which the Tasks re-record.

## 8. Snapshots

`ModelRecommendationSnapshotVersion` (`internal/config/recommendations.go:6`, `2026-11-17`) dates the Recommended Profiles only; `CheckRecommendations` and the `recommendations:` Doctor line read profiles and never the authoring declaration. The Recommended Authoring Declaration has no snapshot, and changing its research tools moves no profile, so the snapshot stays `2026-11-17`.

## 9. The whole prototype

With both halves, every document change of the TechSpec's API Contracts 2 and 5, the ADR-0415 amendment, the synchronized mirrors, the raised Roundfix skill version and the re-recorded Behavior Surface Record applied to the prototype, the four Task Verification commands of task_01 and task_02 exited 0, and `go test -count=1 ./internal/cli ./internal/config ./skills` passed whole (`internal/cli` in 974 s with `GOMAXPROCS=2`). Through the prototype binary, in temporary repositories with a fresh `HOME` and a `gh` that exits 1, the Update Check rows of task_03 requirement 5 and Surface Transcript 1 printed exactly what the Tasks prescribe.

## 10. A configuration file that does not parse

Measured on 2026-11-19 on a prototype rebuilt from the saved diff. With a Project Config of three lines, `defaults:`, `  profile_proof_setup_timeout: 90s` and `broken: a: b`, a user-scope `profiles configure` printed nothing on standard output, printed `roundfix: profiles failed: parse config "<repository>/.roundfixrc.yml": yaml: line 3: mapping values are not allowed in this context` and `Run 'roundfix profiles --help' for usage.` on standard error, exited 2 and proved and wrote nothing. The YAML parser's line is exact for this input; for an unclosed flow sequence or a tab it names an earlier line (`yaml: line 1: did not find expected ',' or ']'` for a list opened on line 2), which is the line `Load` reports too. Returning the parser's error without the file name failed `TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig`. The prototype was removed afterwards.
