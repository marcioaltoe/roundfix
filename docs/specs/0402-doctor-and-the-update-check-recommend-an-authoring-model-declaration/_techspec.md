---
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
prd: _prd.md
created: 2026-11-17
---

# Doctor and the Update Check recommend an Authoring Model Declaration — Technical Spec

## Executive Summary

The fix adds one table and one report state. `internal/config` gains the Recommended Authoring Declaration, the model and reasoning effort of `idea`, `prd`, `techspec`, `tasks` and `research`, and the function that renders it as the YAML block for `.roundfixrc.yml`. `internal/cli/authoring_reports.go` gains the `undeclared` report: with no `skills.authoring.claude`, no problem from the ADR-0401 inspection, and at least one installed authoring skill the repository does not publish, Doctor prints `authoring-models: warn` followed by the block, and the Update Check prints `Authoring models:` followed by the block and sets `authoringModels.state` to `undeclared`. A repository contract test holds this repository's `.roundfixrc.yml` and the configuration guide's block equal to the table.

The primary trade-off is that the recommendation is a copy in the binary of what this repository declares, held equal by a contract test rather than read from the source. The binary cannot read this repository's configuration in an adopter, and a generated record would add a derived file and a tooling rule for ten values (ADR-0415). The second trade-off is that Doctor, which prints one line per check, now prints indented lines after one of them; the block cannot start at the first column, where its `skills:` line would read as Doctor's own `skills:` check.

## Project Constraints

- Identifier strategy: not applicable — no persistent identifier is minted. The new names are a JSON state (`undeclared`) and two JSON fields (`installed`, `recommendation`), public CLI contract under `docs/agents/cli.md`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — both checks read files of the checkout; the Update Check's release lookup is unchanged and this Spec adds no network call. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0415 (new), ADR-0401, ADR-0408 and ADR-0402 bind the design; ADR-0405 and ADR-0414 cite ADR-0401 and do not apply. ADR Alternatives relates each one. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: "Mesmas do BE (Recommended)" (2026-11-17), quoted in the record. It covers the Roundfix skill's setup and modes references and its `SKILL.md` pair. It also covers the three changed lines of this repository's `.roundfixrc.yml` declaration, which must equal the maintainer's recommended values. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `.roundfixrc.yml`, `skills/roundfix/SKILL.md`.

## System Architecture

No package, layer or directory is new. The fix extends three seams.

### Recommendation

A new file `internal/config/authoring_recommendation.go` sits beside `recommendations.go`, which holds the Recommended Profile. It holds one unexported `AuthoringHarness` value and two exported functions: one returns a fresh copy, the other renders the block with fixed bytes. The table's values are the maintainer's of 2026-11-17 (about 12:40 BRT, "Bloco recomendado: - Fable xhigh para Ideia - Fable high para PRD e TechSpec - Opus high para Tasks - Sonnet medium para pesquisa"): `fable`/`xhigh` for `idea`, `fable`/`high` for `prd` and `techspec`, `opus`/`high` for `tasks`, and `sonnet`/`medium` for `research`, with no `tools`. This repository's declaration then declared `fable`/`xhigh` for `prd` and `techspec` and `fable`/`high` for `tasks`, so task_02 changes those three entries of `.roundfixrc.yml` to the table and regenerates the three Author Subagents they produce.

### Reports

`internal/cli/authoring_reports.go` owns both checks' authoring output. A new helper lists the installed authoring skills the repository does not publish: for each phase of `roundskills.AuthoringPhases()`, `.agents/skills/<skill>/SKILL.md` is a regular file and `roundskills.PublishesAuthoringSkill` is false. `loadedAuthoringState` keeps its result. A new step runs only when the declaration is nil and the inspection returned no problem; when that list is not empty, it returns the `undeclared` report.

- `doctorAuthoringModelsResult` gains the `undeclared` branch: status `warn`, the detail and next action of API Contract 1, and the block as continuation lines.
- `CheckResult` in `internal/cli/health.go` gains `Lines []string`. `printDoctorResult` in `internal/cli/doctor.go` prints each after the result line with two leading spaces; every other check leaves it empty and prints the same bytes.
- `updateCheckAuthoring` gains `Installed` and `Recommendation`, both omitted when empty. `printUpdateCheckAuthoringModels` prints the `undeclared` line and the block, each block line with two leading spaces.

The Doctor and Update Check call sites keep their order and conditions: the line still follows `skills:` and `user-skills:`, and `Authoring models:` still follows `User skills:`.

### Contract and guides

A new repository contract test `internal/docscontract/authoring_recommendation_test.go` reads this repository's `.roundfixrc.yml` and `docs/user-guide/configuration.md`. This repository's declaration changes to the table in the same Task, and `roundfix skills install --target project`, run from the tree, regenerates `.claude/agents/roundfix-prd-author.md`, `roundfix-techspec-author.md` and `roundfix-tasks-author.md`; the four authoring skills receive no header, because this repository publishes them (ADR-0408). The configuration guide's example block becomes the rendered recommendation. The Doctor and Update Check guides, the Roundfix skill's setup reference and `--update` mode, the glossary and ADR-0401 describe the report.

## Implementation Design

### Interfaces

```go
// internal/config/authoring_recommendation.go
// RecommendedAuthoringDeclaration returns a fresh copy of the
// Recommended Authoring Declaration for the claude harness.
func RecommendedAuthoringDeclaration() AuthoringHarness

// RecommendedAuthoringYAML renders it as the block for .roundfixrc.yml,
// starting at "skills:" and ending with a newline, with no tools key.
func RecommendedAuthoringYAML() string

// internal/cli/health.go
type CheckResult struct {
	// existing fields unchanged
	Lines []string // printed after the result line, each with two leading spaces
}

// internal/cli/authoring_reports.go
type updateCheckAuthoring struct {
	// existing fields unchanged, in the same order
	Installed      []string `json:"installed,omitempty"`
	Recommendation string   `json:"recommendation,omitempty"`
}
```

Numbered invariants:

1. The `undeclared` report exists only when all hold: Project Config was read and the repository root is known; `authoringDeclaration` returns nil; `InspectAuthoring` returns no problem; at least one authoring skill is installed under `.agents/skills/` and not published under `skills/`.
2. `installed` lists those skills in phase order: `write-idea`, `write-prd`, `write-techspec`, `write-tasks`.
3. Every input that does not meet invariant 1 produces the bytes it produced before this Spec, in Doctor and in both Update Check formats.
4. Neither check writes a file, changes its exit code, or changes the Update Check's `next`.
5. `RecommendedAuthoringYAML()` loaded as Project Config validates, and its `skills.authoring.claude` equals `RecommendedAuthoringDeclaration()`; it holds no `tools` key.
6. This repository's `.roundfixrc.yml` declares, for each of the five keys, the model and reasoning effort of `RecommendedAuthoringDeclaration()`, and the first fenced YAML block under `## Authoring models` in `docs/user-guide/configuration.md` equals `RecommendedAuthoringYAML()`.
7. A `CheckResult` with no `Lines` prints the same bytes as before this Spec.

### Data Models

No stored data changes. The Update Check JSON's `authoringModels` gains the state `undeclared` and the members `installed` and `recommendation`, present only in that state, after `next`. With `undeclared`, `declared` and `problems` are `[]` and `next` is omitted, because installation alone does not apply a declaration that does not exist. `updateCheckSchema` stays `roundfix/update-check/v1`, because the change only adds members and a state.

### API Contracts

1. API Contract: `roundfix doctor` — in the `undeclared` case, the authoring line is `authoring-models: warn (the authoring skills are installed without an Authoring Model Declaration: <skills>; next: add the Recommended Authoring Declaration below to .roundfixrc.yml, then run roundfix skills install --target project)`, where `<skills>` joins `installed` with `, `. Each line of `RecommendedAuthoringYAML()` follows it with two leading spaces, before the next check's line. The exit code is unchanged.
2. API Contract: `roundfix update check` text — in the `undeclared` case, the line is `Authoring models: the authoring skills are installed without a declaration: <skills>; add the Recommended Authoring Declaration below to .roundfixrc.yml, then apply with roundfix skills install --target project.`, followed by each line of the block with two leading spaces, in the position of today's `Authoring models:` line. The `Next:` line and the exit code are unchanged.
3. API Contract: `roundfix update check --format json` — `authoringModels` is `{"state":"undeclared","harness":"claude","declared":[],"problems":[],"installed":[<skills>],"recommendation":"<block>"}`, where `<block>` is `RecommendedAuthoringYAML()` as a JSON string. Every other state keeps its members.
4. API Contract: the Roundfix skill's `--update` mode — when `authoringModels.state` is `undeclared`, the report shows the recommendation and leaves `.roundfixrc.yml` to the user; the mode never writes it.

### Surface Transcripts

1. Surface Transcript: in a Git repository where `roundfix skills install --target project` installed the skills and `.roundfixrc.yml` declares no `skills`, the Update Check recommends the declaration.

   ```transcript
   $ roundfix update check
   stdout:
   ...
   Authoring models: the authoring skills are installed without a declaration: write-idea, write-prd, write-techspec, write-tasks; add the Recommended Authoring Declaration below to .roundfixrc.yml, then apply with roundfix skills install --target project.
     skills:
       authoring:
         claude:
           idea:
             model: fable
             reasoning_effort: xhigh
           prd:
             model: fable
             reasoning_effort: high
           techspec:
             model: fable
             reasoning_effort: high
           tasks:
             model: opus
             reasoning_effort: high
           research:
             model: sonnet
             reasoning_effort: medium
   ...
   stderr:
   exit: 0
   ```

The first `...` stands for the version and Baseline lines, which depend on the build and the lookup. The last stands for the `Next:` line, present when a Baseline plan is ready.

## Coverage Map

- Goal 1 (the warning) → Reports, API Contracts 1–3, Task 1.
- Goal 2 (the block and the next step) → Recommendation, Reports, API Contracts 1–3, Surface Transcript 1, Task 1.
- Goal 3 (no failure, no exit change, no write) → invariant 4, Task 1.
- Goal 4 (the maintainer's values, declared here too) → Recommendation, Contract and guides, Tasks 1 and 2.
- Goal 5 (one source, no silent drift) → Recommendation, Contract and guides, invariants 5 and 6, Tasks 1 and 2.
- Goal 6 (today's cases keep their bytes) → invariants 3 and 7, Task 1.
- Core Feature 1 → Task 1. Core Feature 2 → Task 1. Core Feature 3 → Task 2. Core Feature 4 → Tasks 1 and 2.
- Success Metrics 1, 2 and 3 → Task 1. Success Metric 4 → Task 2. Success Metric 5 → QA gate.

Deviations from the PRD: None.

## Absence Evidence

1. Recommended Authoring Declaration — command: `rtk proxy grep -rniE "RecommendedAuthoring|recommended authoring" internal skills CONTEXT.md docs/user-guide`; result: no match. `internal/config/recommendations.go` recommends Agent Selection Profiles per Work Category only, and `docs/user-guide/configuration.md` shows an example marked "Roundfix ships no model name".
2. A report for installed authoring skills without a declaration — command: `rtk proxy grep -n "len(state.Declared) > 0 || len(state.Problems) > 0" internal/cli/authoring_reports.go`; result: one match in `loadedAuthoringState`, the only gate of both reports; with no declared key and no problem, Doctor prints no `authoring-models:` line and the Update Check sets no `authoringModels`. Measured in Argus with Roundfix 0.147.1: `authoringModels` absent.
3. Continuation lines in Doctor output — command: `rtk proxy grep -n "Lines" internal/cli/health.go`; result: no match. `CheckResult` holds `Name`, `Status`, `Detail`, `NextAction` and `Err`, and `printDoctorResult` prints one line.
4. Glossary term **Recommended Authoring Declaration** — command: `rtk proxy grep -n "Recommended Authoring Declaration" CONTEXT.md`; result: no match.

## ADR Alternatives

- ADR-0415: the design is that ADR's decision. It rejects reading this repository's configuration at build time, a generated record, the Model Catalog, a hard-coded block in the command code, writing `.roundfixrc.yml`, failing the check, recommending `research.tools`, and the same warning in `roundfix skills check`.
- ADR-0401: amended by ADR-0415. The design takes, in a narrowed form, the alternative ADR-0401 rejected as "Model names shipped as defaults": the names ship as a recommendation the user pastes, never as a default, as the Recommended Profile already does. ADR-0401's other rejections stand: no check fails for the declaration, and the section stays in Project Config only. Task 2 adds the amendment line to ADR-0401.
- ADR-0408: kept. A skill the repository publishes counts as not installed for the `undeclared` report, so this repository, which publishes all four, is exempt; the rejected configuration key that withholds the header stays rejected, and no key silences the warning either.
- ADR-0402: kept. The block declares `research` without `tools`, so the generated Research Subagent receives `Read`, `Grep` and `Glob` only; the rejected "web tools in the fixed list" stays rejected, in the recommendation too.
- ADR-0405: not applicable — it reports user-level copies of a declared phase's skill; the `undeclared` report exists only when no phase is declared, so no copy hides a header, and the `user-skills:` line and `userSkills` keep their rules.
- ADR-0414: not applicable — it decides when a Supervisor session ends; this Spec changes no session end and no phase record.

## Integration Points

- The file system of the checkout, through the existing `os.Stat` checks of `skills/authoring.go`. No new command and no network call.
- Adopters receive the report with the binary and the skill text with `roundfix skills install`.

## Testing Approach

- **Task 1, recommendation.** In `internal/config`, a new `authoring_recommendation_test.go` asserts the exact rendered bytes, loads them as Project Config in a temporary directory with the existing loader and compares the result with `RecommendedAuthoringDeclaration()`, and proves the returned value is a copy.
- **Task 1, reports.** In `internal/cli`, a new `authoring_recommendation_test.go` uses the existing fixtures (`withCLIWorkspace`, `installAuthoringReportFixture`, `configureAuthoringDoctorFixture`, `configureAuthoringUpdateFixture`, `runUpdateCheckTest`, `snapshotDoctorPath`):
  - `undeclared`: the whole Doctor line and every block line in position, exit unchanged; the whole JSON member and the text line and block; the report's `next` unchanged; the repository unchanged;
  - `partly_published`: `skills/write-prd/SKILL.md` and `skills/write-tasks/SKILL.md` in the tree name only the other two;
  - `published`, `no_authoring_skill`, `leftover_generated`, `declared`: today's bytes;
  - the Surface Transcript test, one named test asserting the transcript's lines through `runCLI`.
- **Declared expectation changes.** The probe of 2026-11-17 on a scratch worktree of `bcd716df` ran the whole `./internal/cli` package with the report applied and found two tests that change: `TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted` (its `doctor` and `update_check` cases install the authoring skills without a declaration) and `TestRunDoctorRealRepositoryCheckDoesNotMutateState` (its ready fixture installs them too). The first keeps its purpose with a repository that has no installed authoring skill; the second adds the warn line and the block to its whole expected output. No other test of the package changed.
- **Task 2, contract.** `internal/docscontract/authoring_recommendation_test.go`, under the `docscontract` tag and `//verify:always`, holds invariant 6 with one subtest each for the repository's configuration and the guide. Sabotage: changing one effort in a copy of `.roundfixrc.yml` or one line of the guide's block fails the named subtest.
- **Task 2, skill.** `skills/authoring_delegation_test.go` pins the phrase "Roundfix ships no model name" in the setup reference; its expectation changes to the new sentence, a declared change.
- The Surface Transcript runs in the QA gate through the built binary on a temporary repository.

## Build Order

1. Recommendation and reports: `internal/config/authoring_recommendation.go`, the `undeclared` report in `internal/cli/authoring_reports.go`, `Lines` in `health.go` and `doctor.go`, their tests and the two changed tests, and the Doctor and Update Check guides with the Behavior Surface Record.
2. Contract and documentation: this repository's declaration changed to the table with its three regenerated Author Subagents, the repository contract test, the configuration guide's block, the Roundfix skill's setup reference and `--update` mode with the mirror, the version record and the changed skills test, the glossary through `domain-modeling`, ADR-0401's amendment line and the Behavior Surface Record (depends on: 1 — the contract test calls the rendering and the guides describe the report step 1 builds).
3. QA gate: the PRD's metrics, the transcript and the outside-evidence fixture (depends on: 1, 2).

## Risks & Considerations

- **A paste into a file that already has `skills:`.** The block starts at `skills:`. A repository whose `.roundfixrc.yml` holds `skills: {authoring: {}}` reaches the report too and must merge by hand; the guides say so. A duplicate key is refused when the configuration loads, so the mistake cannot pass silently.
- **The recommendation names `fable` and `opus`.** The alias depends on the Claude Code installed; a model name the harness refuses fails the subagent call, as ADR-0401 measured. The block is a recommendation the user reads before pasting.
- **Spec 0401 changes `modes.md`.** Both Specs change the Roundfix skill; the second to deliver records the next version on the merge base and rebases its sentence.
- **Doctor parsers.** Tools that read Doctor's output line by line now see indented lines in one case; every check line still starts at the first column with its name.

## Glossary

- adds: **Recommended Authoring Declaration**
- changes: **Authoring Model Declaration**

## Decisions

- The table holds the maintainer's values of 2026-11-17, and this repository's declaration changes to them.
- The recommendation lives in one table in `internal/config`, held equal to this repository's declaration and the configuration guide by a repository contract test. See ADR-0415.
- The block omits `research.tools`. See ADR-0415 and ADR-0402.
- A problem the inspection finds keeps today's report; `undeclared` applies only without one, so a leftover generated file is fixed by installation first.
- Doctor prints the block as indented continuation lines, and the Update Check JSON carries it as text.
- `roundfix skills check` keeps its output. See ADR-0415.
- The Release Ceiling is `minor`: Doctor gains a warn case and the Update Check JSON gains a state and two members.
