# Harness facts and test measurements for Spec 0405 — 2026-11-17

Recorded by the author of the Idea, the PRD and the TechSpec, on Opus 5.5, in
a clone of this repository on the branch `docs/plan-0405`, based on
`a8236a87`.

## What Claude Code says about fallbacks and usage limits

The author read the documentation on 2026-11-17 from
`https://code.claude.com/docs/en/<page>.md`. The installed CLI was
`2.1.326 (Claude Code)`, and `claude --help` was read without running a
session.

- **Skill frontmatter** (`skills`): `model` and `effort` apply "for the rest
  of the current turn". There is no fallback field.
- **Subagent frontmatter** (`sub-agents`, supported frontmatter fields): it
  takes `model` (`sonnet`, `opus`, `haiku`, `fable`, a full model ID or
  `inherit`) and `effort` (`low`, `medium`, `high`, `xhigh`, `max`). There is
  no fallback field.
- **Model order of a subagent** (`sub-agents`, choose a model): "The
  per-invocation `model` parameter … On an Agent tool call, the parameter
  takes a model alias, not a full model ID or `inherit`", and it comes before
  the `model` field of the definition. "If Claude passed an alias, the
  subagent keeps that alias when it resumes or receives a follow-up message."
- **Effort of a subagent** (`sub-agents`, choose an effort level): "it can
  also pass an `effort` parameter for that invocation. The parameter overrides
  the `effort` field and stays in effect when the subagent is resumed … The
  per-invocation parameter requires Claude Code v2.1.292 or later."
- **The fallback chain** (`model-config`, fallback model chains): it is
  `--fallback-model` or the setting `fallbackModel`, with at most three
  models, and it lasts one turn. "Authentication, billing, rate-limit,
  request-size, and transport errors … never trigger a switch." Since v2.1.247
  it also applies to subagents. `claude --help` reads: "Models to fall back to
  when the main model is busy or unavailable (comma-separated, in order)."
- **A subagent that ends on an API error** (`sub-agents`, API errors in
  subagents): "A subagent whose run ends on an API error, such as a usage
  limit or a repeated server error, reports that failure back to Claude." In
  the foreground, partial text comes back with a note that the subagent was
  cut off. In the background, the subagent is marked failed and its last
  output is returned.
- **Usage limit messages** (`errors`, you've hit your session limit):
  `You've hit your session limit · resets 3:45pm`, `You've hit your weekly
  limit · resets Mon 12:00am`, `You've hit your Opus limit · resets 3:45pm`,
  `You've hit your Sonnet limit · resets 3:45pm`. "The session and weekly
  limits are shared across all models … The Opus and Sonnet limits each apply
  only to requests to that model family." The same page lists spend-limit
  messages that begin `You've hit your monthly spend limit`, `individual spend
  limit` and `individual usage limit`. A Fable consent message reads `Fable
  limit reached · continuing on Fable 5.1 uses usage credits, and the prompt to
  confirm went unanswered · …`.
- **Observed in this repository on 2026-11-17:** the TechSpec Author Subagent
  of Spec 0404 ended with `You've hit your Fable limit · resets Nov 22 at 8pm`.

## Governed Paths of the design

The author measured each path with `speccheck.GovernedPath`, through a
`go test -overlay` probe in `internal/speccheck` that wrote nothing to the
repository.

- **Governed and changed by the Tasks:**
  - `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-tasks/SKILL.md` and their mirrors under `skills/`;
  - `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`;
  - `.agents/skills/roundfix/references/phase.md` and `.agents/skills/roundfix/references/setup.md`, which describes the Authoring Model Declaration.
- **Governed and not changed:** `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`, `.roundfixrc.yml`, the `Makefile`, `go.mod` and `.github/workflows/ci-verify.yml`.
- **Not governed:**
  - every Go source and test file under `internal/config`, `internal/phaserecord` and `internal/cli` the design names, `internal/cli/cli.go` included;
  - `skills/roundfix/references/phase.md`, `skills/roundfix/references/setup.md`, `skills/phase_reference_test.go` and new tests under `skills/`;
  - `docs/user-guide/commands/phase.md`, `docs/user-guide/configuration.md` and `docs/user-guide/usage.md`;
  - `CONTEXT.md`, ADR-0401, ADR-0407, the new ADRs and this Spec's files;
  - `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `skills/testdata/owned-skill-versions.json` and `skills-lock.json`.

## Existing tests a prototype of the surfaces moves

The prototype ran in a scratch worktree of `a8236a87`, removed afterwards, with
three changes:

- `fallback` accepted under an entry of the Authoring Model Declaration, with the refusal of an unknown key listing it;
- `roundfix phase model` and `roundfix phase fallback` added to the usage of `roundfix phase`;
- the same two lines added to the root help, with the summary of `phase` widened.

It then ran `go test -count=1` on `./internal/config`, `./skills` and `./internal/cli`, and `-tags docscontract` on `./internal/docscontract`.

- `./internal/config`: `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/phase_field`, whose expected refusal names the accepted keys of an entry, failed.
- `./internal/cli` (188 s): only `TestPhaseStartSurfaceTranscripts/transcript_1`, which holds the full help of `roundfix phase`, failed. `TestEveryCommandPathPrintsItsHelp` passed with the two new command paths.
- `./internal/docscontract`: four tests failed, each for a documentation or record the Tasks write:
  - `TestEveryCommandIsNamedInTheRoundfixSkill` and `TestEveryCommandIsNamedInTheUserGuide`, because no text named `phase model` or `phase fallback`;
  - `TestAnUndocumentedCommandIsReported`, because the undocumented list held the two new paths beside `reopen`;
  - `TestTheSkillCoverageMapIsCurrent`, because the two paths had no Skill Coverage Map entry and no Behavior Surface Record entry, and the fingerprints of `command phase start` and `config keys` were stale.
- `./skills`: passed.

Specs 0401 to 0404 were not on the tree, and they change the same help and the
same phase entry. Each Task therefore runs these packages again on its own
base.
