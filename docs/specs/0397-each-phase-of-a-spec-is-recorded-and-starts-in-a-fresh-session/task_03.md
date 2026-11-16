---
task: task_03
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
status: completed
type: backend
complexity: high
---

# Task 03: A published authoring skill receives no Declared Skill Header, and this repository declares its models

## Overview

Installation learns one rule from the tree: an authoring skill whose source
sits at `skills/<skill>/SKILL.md` under the repository root is installed
without the Declared Skill Header, the installation says so, and every
declared subagent is still generated. The three checks count the declaration
as applied for such a skill and report a header found there. This repository
then declares its authoring models and commits the five generated subagents,
with its four authoring skills still equal to their mirrors. The slice is
whole on its own: in a temporary repository that publishes one of two
declared skills, only the other receives a header; and in this checkout the
product built from the tree reports no unapplied declaration. The
configuration, skills, doctor and update guides, the Roundfix skill's `setup`
reference, the glossary and ADR-0401 describe the rule.

## Requirements

1. MUST apply `_techspec.md` → API Contract 7 in the installation of the
   `skills` package. For each declared Authoring Phase whose skill has a
   regular file at `<root>/skills/<skill>/SKILL.md`, the installation writes
   the shipped file with no header and names the skill in a new list of its
   result, `Withheld`, declared right after `Headers` with the JSON name
   `withheld` and omitted when empty; a skill without that file receives its
   header as before. Every declared subagent is generated as before, an
   installed copy that held a header loses it, and a second installation
   changes no byte. `roundfix skills install --target project` prints
   `Withheld declared header: <skill> (this repository publishes skills/<skill>/SKILL.md)`
   in the place and phase order of its `Applied declared header:` lines
   (Surface Transcript 6). The result of `roundfix baseline update` carries
   `withheld` beside `headers` through the same structure; its
   `Authoring models:` text line and every existing field stay as they are.
2. MUST apply `_techspec.md` → API Contract 8 in the inspector. For a
   declared phase whose skill the repository publishes, `InspectAuthoring`
   reports no problem of that skill, neither the header nor the Claude Skill
   Link, when the installed `SKILL.md` holds no top-level `model` and no
   top-level `effort` or does not exist, and otherwise reports exactly one
   problem, fixable by the installation:
   `<skill> holds a header, and this repository publishes skills/<skill>/SKILL.md`,
   in the place the header problem of that phase takes. `roundfix skills
   check`, Doctor and the Update Check print it through their existing
   lines and change no exit code. A user-level copy of such a skill is
   reported by its version only: its `hidesDeclaredHeader` is false and its
   description never says that it hides the declared header. The inspector
   still writes, creates, links and deletes nothing.
3. MUST keep Invariant 6: without a declaration, and in a repository that
   publishes no skill, the installation, `roundfix baseline update` and the
   three checks print the bytes they print on the default branch when the
   Task runs. No existing test of the installation, the inspector, the
   user-level copies or the three reports is changed, removed or renamed.
4. MUST add these tests, each top-level test calling `t.Parallel()` first,
   in temporary repositories and homes through the packages' existing seams,
   without editing `internal/cli/cli_test.go`:
   - `skills`: `TestInstallationWithholdsTheHeaderOfAPublishedSkill`, with
     the subtests `one_published` (fixture D of `_techspec.md` → Surface
     Transcripts: the installed `write-prd` equals the shipped file, the
     installed `write-techspec` holds the header, both subagents exist, and
     the result names each skill in its list), `second_run_changes_no_byte`,
     `every_authoring_skill_published`, `none_published` and
     `a_previous_header_is_removed`; and
     `TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader`, with the
     subtests `no_header_is_no_problem` (in a repository without
     `.claude/skills`, where the unpublished skill still gets its link
     problem), `a_header_is_the_problem` (the whole list of problems with
     their fixable marks) and `the_tree_is_unchanged`.
   - `skills`: `TestTheSetupReferenceNamesThePublishedSkillRule`, which
     reads both copies through `roundfixLayoutSources` and requires each
     phrase requirement 6 lists, exactly once, whitespace-normalized, in one
     named subtest each.
   - `internal/cli`: `TestThePublishedSkillSurfaceTranscripts`, with the
     subtests `transcript_6` and `transcript_7`, each asserting every line
     of standard output and standard error and the exit code of that
     Surface Transcript. Fixture E is fixture D after Surface Transcript 6
     with the line `model: opus` added to the installed `write-prd` and with
     `.claude/skills` made a symbolic link to `../.agents/skills`: without
     that link the inspector also names the Claude Skill Link of
     `write-techspec`, a line the transcript does not show (measured on the
     planning branch). The home holds no user-level skill copy.
   - `internal/cli`: `TestTheChecksCountAPublishedSkillAsApplied`, with the
     subtests `skills_check` (an empty standard error), `doctor` (the
     `authoring-models` line is `ok`) and `update_check` (the state
     `applied` and no `Authoring models:` line), in a repository that
     publishes every declared skill and holds the Claude Skill Link;
     `TestBaselineUpdateReportsAWithheldHeader`, which asserts the whole
     `skills.authoring` member of the JSON result for fixture D, `withheld`
     right after `headers`; and
     `TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly`, with
     the subtests `a_current_copy_is_not_named` and
     `an_older_copy_hides_no_header`.
   - Every test file compiles for Windows.
5. MUST document the rule in the guides and then re-record the Behavior
   Surface Record, after the last guide edit and before requirement 9, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`:
   - `docs/user-guide/configuration.md`, in the section on the authoring
     models: a paragraph saying that a repository that holds
     `skills/<skill>/SKILL.md` for a declared authoring skill receives no
     header for it, that installation still generates every declared
     subagent, the sentence "The rule is read from the tree; no
     configuration key turns it on.", and the sentence "A Roundfix older
     than the release that introduced the rule writes the header into such a
     skill, so upgrade it before installing there.";
   - `docs/user-guide/commands/skills.md`: the line
     `Withheld declared header: <skill> (this repository publishes skills/<skill>/SKILL.md)`
     in the `text` block of the installation's lines, with when it prints;
     the problem
     `<skill> holds a header, and this repository publishes skills/<skill>/SKILL.md`
     among the problems `skills check` names; that a user-level copy of a
     published skill hides no header; and that the JSON result of
     `roundfix baseline update` carries `withheld` in `skills.authoring`;
   - `docs/user-guide/commands/doctor.md` and
     `docs/user-guide/commands/update.md`: the sentence "For an authoring
     skill the repository publishes, the declaration counts as applied when
     the installed skill holds no model and no effort." where each guide
     describes its authoring line, and, in the update guide's description of
     `userSkills[].hidesDeclaredHeader`, that it is never true for a skill
     the repository publishes.
6. MUST add to `.agents/skills/roundfix/references/setup.md`, as a paragraph
   of its own right after the paragraph that holds "write the Declared Skill
   Header into the four installed authoring skills", word for word:

   > A repository that publishes an authoring skill under `skills/<name>/` receives no Declared Skill Header for it: installation writes the shipped skill, prints `Withheld declared header: <name>` with the reason, and still generates every declared subagent. The checks count the declaration as applied when that installed skill holds no `model` and no `effort`, and report a header found there. A Roundfix older than the release that introduced this rule writes the header into a published skill, so upgrade it before installing there.

   pinned as `a_published_skill_receives_no_header` ("A repository that
   publishes an authoring skill under `skills/<name>/` receives no Declared
   Skill Header for it"), `the_subagents_are_generated` ("still generates
   every declared subagent"), `the_checks_agree` (the whole second sentence)
   and `an_older_binary_writes_it` (the whole third sentence). Every other
   sentence of the reference stays, and every phrase
   `TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase` pins there
   still occurs once. MUST then run `make skills-sync` and the version
   record
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version of `roundfix` from the tree the Task runs on;
   never write a number or a mirror by hand. MUST NOT add or change text
   inside any `### QA settlement` section.
7. MUST change in `CONTEXT.md`, through `domain-modeling`, only the
   definition line of **Declared Skill Header**, to these words, where each
   `(ADR-NNNN)` is the glossary's Markdown link to that ADR's file under
   `docs/adr/` as the entry writes its links today:

   > The `model` and `effort` lines that installation writes into the frontmatter of an installed authoring skill from the Authoring Model Declaration, so that a phase the user starts by typing the skill runs on them. The shipped skills never hold one, a user-level installation receives none, and a user-level copy of the skill hides it in Claude Code. An authoring skill whose source the repository itself publishes under `skills/<name>/` receives none either: installation withholds it and says so, the checks report a header found there, and no user-level copy of that skill hides one (ADR-0401, ADR-0405, ADR-0408).

8. MUST add to
   `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`,
   as a paragraph of its own right after the Decision paragraph that begins
   "**The Author Subagent.**", this notice, word for word, set the ADR's
   `updated_at` to the day the Task runs in the form `YYYY-MM-DDT00:00:00Z`,
   and change nothing else in it. In the notice, `ADR-0408` is a Markdown
   link whose text is `ADR-0408` and whose destination is the file
   `0408-a-repository-that-publishes-an-authoring-skill-receives-no-header-for-it-and-a-fix-is-authored-in-the-session.md`
   beside it, as the notice ADR-0386 carries writes its link:

   > **Amended by ADR-0408 (2026-11-16).** Installation writes no Declared Skill Header into an authoring skill whose source the repository itself publishes under `skills/<name>/`, and still generates every declared subagent; this repository declares the section on those terms. An authoring skill the agent activated on the fix or refactor route runs its phase in the session and starts the Author Subagent only when the user asks. A feature's phases are still handed over, a typed skill still runs on its header where one is written, and Roundfix still ships no model name.

9. MUST then adopt the declaration in this repository, as the Task's last
   step and only after requirements 1 to 8 hold and their tests pass:
   - add to `.roundfixrc.yml`, as a new top-level section at the end of the
     file, `skills.authoring.claude` with exactly the entries of
     `_techspec.md` → API Contract 12: `idea`, `prd` and `techspec` with
     `model: fable` and `reasoning_effort: xhigh`; `tasks` with
     `model: fable` and `reasoning_effort: high`; and `research` with
     `model: sonnet`, `reasoning_effort: medium` and
     `tools: [WebSearch, WebFetch, mcp__exa__web_search_exa, mcp__exa__web_fetch_exa]`.
     No other key of the file changes;
   - run `make skills-sync`, and then run the installation exactly as
     `go run -buildvcs=false ./cmd/roundfix skills install --target project`,
     the product built from this tree. MUST NOT run `skills install` or
     `baseline update` with a `roundfix` found on `PATH`, in this Task or in
     its checks: a binary older than this Task writes the header into the
     four canonical authoring skills;
   - the installation leaves five new files, the generated subagents under
     `.claude/agents/`, and no other change: every file under
     `.agents/skills/` and `skills/` and `skills-lock.json` are
     byte-identical to what they were before it. When it changed one of
     them, the bundle was stale or the rule is not in the tree: undo that
     change, find the cause, and run the two commands again. Never edit a
     generated subagent by hand.
10. Of the Governed Paths the Task may change only `.roundfixrc.yml`,
    `.agents/skills/roundfix/SKILL.md`,
    `.agents/skills/roundfix/references/setup.md` and
    `skills/roundfix/SKILL.md`, the two `SKILL.md` files only through the
    two commands of requirement 6. MUST NOT change the `SKILL.md` of
    `write-idea`, `write-prd`, `write-techspec` or `write-tasks` in either
    copy, add a `model` or `effort` key to any shipped skill, change the
    baseline guide, `skills-lock.json`, any Baseline file, any other ADR,
    `docs/references/coverage-record.json`, the Makefile, the CI workflows,
    `go.mod`, `go.sum` or the lint configuration.

## Subtasks

- [ ] Withhold the header of a published skill in the installation and report it.
- [ ] Make the inspector and the user-level copy reader follow the rule.
- [ ] Add the tests in the `skills` and `internal/cli` packages.
- [ ] Write the four guides and re-record the Behavior Surface Record.
- [ ] Add the `setup` paragraph, sync the mirror and record the skill version.
- [ ] Change the glossary entry and add the notice to ADR-0401.
- [ ] Declare this repository's authoring models and install with the product built from the tree.

## Acceptance Criteria

- [ ] In a temporary repository that holds the published source of one declared authoring skill and not of another, installation writes the header into the second only, says which header it withheld, generates every declared subagent, and a second run changes no byte.
- [ ] The three checks report the declaration as applied when a published skill holds no model and no effort, report a header placed by hand in it as a problem the installation fixes, and never say that a user-level copy of it hides a header; the two Surface Transcripts hold line for line.
- [ ] In a repository that publishes no skill, and without a declaration, the installation and the three checks print what they printed before.
- [ ] In this repository the product built from the tree prints no warning about the Authoring Model Declaration, the five generated files exist with their markers, no shipped `SKILL.md` holds a model or an effort, and each mirror equals its canonical copy.
- [ ] The four guides and both copies of the `setup` reference describe the rule, `CONTEXT.md` carries the changed **Declared Skill Header**, and ADR-0401 carries one notice that names ADR-0408.

## Context

- instruction: `docs/adr/0408-a-repository-that-publishes-an-authoring-skill-receives-no-header-for-it-and-a-fix-is-authored-in-the-session.md`
- instruction: `docs/adr/0402-the-research-subagent-reads-and-never-runs-a-command-or-writes-a-file.md`
- instruction: `docs/adr/0405-a-user-level-skill-copy-is-reported-with-the-command-that-removes-it.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `skills/authoring_inspect_test.go`
- instruction: `internal/cli/authoring_reports_test.go`
- interface: `skills/authoring.go`
- interface: `skills/authoring_inspect.go`
- interface: `skills/user_copies.go`
- creates: `skills/published_skill_test.go`
- interface: `internal/cli/authoring_models.go`
- interface: `internal/cli/authoring_reports.go`
- interface: `internal/cli/user_skills.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/published_skill_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`
- interface: `.roundfixrc.yml`
- creates: `.claude/agents/roundfix-idea-author.md`
- creates: `.claude/agents/roundfix-prd-author.md`
- creates: `.claude/agents/roundfix-techspec-author.md`
- creates: `.claude/agents/roundfix-tasks-author.md`
- creates: `.claude/agents/roundfix-researcher.md`

## Consumers

- The installation's result and the inspector's problems — consumers: `./skills`, `./internal/cli`
- The Roundfix skill's `setup` reference and the owned-skill version record — consumers: `./skills`
- The configuration, skills, doctor and update guides, the Behavior Surface Record, the glossary and ADR-0401 — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestInstallationWithholdsTheHeaderOfAPublishedSkill|TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader|TestTheSetupReferenceNamesThePublishedSkillRule|TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestInstallationWithholdsTheHeaderOfAPublishedSkill/one_published TestInstallationWithholdsTheHeaderOfAPublishedSkill/second_run_changes_no_byte TestInstallationWithholdsTheHeaderOfAPublishedSkill/none_published TestInstallationWithholdsTheHeaderOfAPublishedSkill/a_previous_header_is_removed TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader/no_header_is_no_problem TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader/a_header_is_the_problem TestTheSetupReferenceNamesThePublishedSkillRule TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the three new tests do not exist, so their pass lines are missing; afterwards a published skill receives no header and no problem, a header found in it is the one problem, both copies of the `setup` reference carry the rule beside every phrase they already held, and the skills suite passes.

- `out="$(go test -count=1 -timeout 40m -v -run '^(TestThePublishedSkillSurfaceTranscripts|TestTheChecksCountAPublishedSkillAsApplied|TestBaselineUpdateReportsAWithheldHeader|TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThePublishedSkillSurfaceTranscripts/transcript_6 TestThePublishedSkillSurfaceTranscripts/transcript_7 TestTheChecksCountAPublishedSkillAsApplied/skills_check TestTheChecksCountAPublishedSkillAsApplied/doctor TestTheChecksCountAPublishedSkillAsApplied/update_check TestBaselineUpdateReportsAWithheldHeader TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly/a_current_copy_is_not_named TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly/an_older_copy_hides_no_header; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m -run 'Skill|Doctor|UpdateCheck|Authoring' ./internal/cli` — expected: exit 0. Before this Task the four tests do not exist, so their pass lines are missing; afterwards Surface Transcripts 6 and 7 hold line for line, the three checks count a published skill as applied, the Baseline update names the withheld header, and the CLI's skills, Doctor, Update Check and authoring tests pass unchanged.

- `tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'The rule is read from the tree; no configuration key turns it on.' || { echo 'the configuration guide lacks the published-skill rule' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'Withheld declared header: <skill> (this repository publishes skills/<skill>/SKILL.md)' || { echo 'the skills guide lacks the withheld line' >&2; exit 1; }; for f in docs/user-guide/commands/doctor.md docs/user-guide/commands/update.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'For an authoring skill the repository publishes, the declaration counts as applied when the installed skill holds no model and no effort.' || { echo "missing the published-skill sentence: $f" >&2; exit 1; }; done; grep -F -- '**Declared Skill Header**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Declared Skill Header\*\*(:| —( |$))' || { echo 'missing term: Declared Skill Header' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'receives none either: installation withholds it and says so, the checks report a header found there, and no user-level copy of that skill hides one' || { echo 'Declared Skill Header unchanged' >&2; exit 1; }; f=docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md; test "$(grep -cF -- '**Amended by [ADR-0408](0408-a-repository-that-publishes-an-authoring-skill-receives-no-header-for-it-and-a-fix-is-authored-in-the-session.md) (2026-11-16).**' "$f")" -eq 1 || { echo "missing amendment: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'a typed skill still runs on its header where one is written, and Roundfix still ships no model name.' || { echo "amendment text differs: $f" >&2; exit 1; }; if grep -q '^updated_at: 2026-11-14T11:00:00Z' "$f"; then echo "updated_at not raised: $f" >&2; exit 1; fi; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guides, the glossary entry and ADR-0401 do not carry the rule; afterwards they do, ADR-0401 holds one notice with a later `updated_at`, and the documentation contracts, the Behavior Surface Record included, pass.

- `grep -q '^skills:' .roundfixrc.yml || { echo 'the Project Config declares no authoring models' >&2; exit 1; }; for key in idea prd techspec tasks research; do f=".claude/agents/roundfix-$key-author.md"; test "$key" = research && f=.claude/agents/roundfix-researcher.md; test -f "$f" || { echo "missing: $f" >&2; exit 1; }; grep -qF -- "<!-- roundfix:generated skills.authoring.claude.$key -->" "$f" || { echo "no marker: $f" >&2; exit 1; }; done; for skill in roundfix write-idea write-prd write-techspec write-tasks; do diff -r ".agents/skills/$skill" "skills/$skill" >/dev/null || { echo "mirror differs: $skill" >&2; exit 1; }; done; d="$(mktemp -d)" || exit 1; trap 'rm -rf "$d"' EXIT; go build -buildvcs=false -o "$d/roundfix" ./cmd/roundfix || exit 1; mkdir "$d/home" || exit 1; err="$(HOME="$d/home" "$d/roundfix" skills check 2>&1 >/dev/null)" || { printf '%s\n' "$err"; exit 1; }; if printf '%s\n' "$err" | grep -q 'Authoring Model Declaration'; then printf '%s\n' "$err"; exit 1; fi; out="$(go test -count=1 -timeout 20m -v -run '^TestNoShippedSkillDeclaresAModelOrEffort$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestNoShippedSkillDeclaresAModelOrEffort ' || { echo 'missing pass: TestNoShippedSkillDeclaresAModelOrEffort' >&2; exit 1; }` — expected: exit 0. Before this Task the Project Config declares nothing and no generated subagent exists; afterwards the five generated files carry their markers, the five changed or installed skills equal their mirrors, the product built from the tree reports no unapplied declaration with an empty temporary home, and no shipped skill holds a model or an effort.

## References

- `_prd.md` → User Stories 7 and 8; Core Features 10, 11 and 13; Success Metrics 6 and 7; Glossary
- `_techspec.md` → Current behavior; API Contract 7; API Contract 8; API Contract 12; Surface Transcript 6; Surface Transcript 7; Invariants 6 and 7; Declared breaks; Testing Approach; Risks & Considerations; Build Order 3
- ADR-0408, ADR-0401, ADR-0402, ADR-0405

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, roundfix, evidence-gate.

Implemented the published-skill rule through installation, inspection and the
CLI reports. The shared tree predicate requires a regular
`skills/<skill>/SKILL.md`. Installation records `withheld` immediately after
`headers`, writes the shipped bytes for those skills, and still generates all
declared subagents. Applied and withheld text lines share phase order.
Inspection skips their Claude Skill Link check, reports a top-level header
as the one installable skill problem, and remains read-only. User-copy
reporting receives only declared skills that can hold a header, so a
published skill's copy is reported by version alone.

Acceptance evidence from focused checks:

- Mixed published/unpublished installation, every authoring skill published,
  none published, removal of a previous header, and an idempotent second run:
  `TestInstallationWithholdsTheHeaderOfAPublishedSkill` passed. It checks
  shipped bytes, complete result lists, generated subagent bytes, and the
  second run's bytes, modes and modification times.
- Inspection without a header or Claude Skill Link, the complete ordered
  problem list with fixable marks, installation repair, an absent installed
  file, effort alone, nested keys, and an unchanged tree:
  `TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader` passed.
  `TestThePublishedSkillSurfaceTranscripts` passed both `transcript_6` and
  `transcript_7`, comparing all stdout, stderr and exit status. Fixture E
  includes the required whole-directory Claude Skill Link.
- Skills Check, Doctor and Update Check accept the clean published skills;
  Doctor and Update Check also report a planted header without changing
  their exit status. `TestTheChecksCountAPublishedSkillAsApplied`,
  `TestBaselineUpdateReportsAWithheldHeader`, and
  `TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly` passed.
  The Baseline test compares the whole `skills.authoring` JSON member,
  including key order, and retains the existing text summary.
- Existing installation, inspector and report tests passed unchanged:
  `TestInstallAppliesTheDeclarationAndASecondRunChangesNothing`,
  `TestInspectAuthoringNamesEachProblem`,
  `TestSkillsInstallAppliesTheAuthoringModelDeclaration`,
  `TestBaselineUpdateAppliesTheAuthoringModelDeclaration`,
  `TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted`,
  `TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit`,
  `TestDoctorReportsTheAuthoringModels`,
  `TestUpdateCheckReportsTheAuthoringModels`, and
  `TestUserSkillReportsNameAHiddenHeaderWithoutRefresh`. No existing test
  was edited, removed or renamed.
- The four guides carry the required rule and exact normalized sentences.
  The setup paragraph was added verbatim and its mirror generated;
  `TestTheSetupReferenceNamesThePublishedSkillRule` and
  `TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase` passed.
  Only the Declared Skill Header definition line changed in the glossary;
  ADR-0401 gained only the required amendment paragraph and `updated_at`.
  Focused glossary, accepted-ADR link and user-guide link checks passed.
- After those checks, appended only `skills.authoring.claude` to Project
  Config with the exact declared values. Ran `make skills-sync`, then
  `go run -buildvcs=false ./cmd/roundfix skills install --target project`
  with an isolated temporary home and existing Go caches, network disabled.
  Installation printed four withheld lines and generated the five marked
  files under `.claude/agents/`. Before/after SHA-256 snapshots of every file
  under `.agents/skills/`, `skills/` and `skills-lock.json` were identical;
  the only added status entries were the five generated subagents.
  A separate isolated-home
  `go run -buildvcs=false ./cmd/roundfix skills check` exited 0 with empty
  stderr. The five skill trees match their mirrors.
  `TestNoShippedSkillDeclaresAModelOrEffort` passed before adoption; the
  installation byte audit proves it preserved those same shipped bytes.

Focused commands and outcomes:

- Initial red check:
  `go test ./skills -run '^TestInstallationWithholdsTheHeaderOfAPublishedSkill$/one_published$' -count=1`
  failed to compile because `AuthoringInstall.Withheld` did not yet exist.
- `go test ./skills ./internal/cli -run '^(TestInstallationWithholdsTheHeaderOfAPublishedSkill|TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader|TestThePublishedSkillSurfaceTranscripts|TestTheChecksCountAPublishedSkillAsApplied|TestBaselineUpdateReportsAWithheldHeader|TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly)$' -count=1`
  exited 0 in both changed packages after implementation.
- `go test ./skills -run '^(TestInstallationWithholdsTheHeaderOfAPublishedSkill|TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader|TestTheSetupReferenceNamesThePublishedSkillRule|TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase|TestInstallAppliesTheDeclarationAndASecondRunChangesNothing|TestInspectAuthoringNamesEachProblem|TestNoShippedSkillDeclaresAModelOrEffort)$' -count=1`
  exited 0 after synchronization and recording.
- `go test ./internal/cli -run '^(TestThePublishedSkillSurfaceTranscripts|TestTheChecksCountAPublishedSkillAsApplied|TestBaselineUpdateReportsAWithheldHeader|TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly|TestSkillsInstallAppliesTheAuthoringModelDeclaration|TestBaselineUpdateAppliesTheAuthoringModelDeclaration|TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted|TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit|TestDoctorReportsTheAuthoringModels|TestUpdateCheckReportsTheAuthoringModels|TestUserSkillReportsNameAHiddenHeaderWithoutRefresh)$' -count=1`
  exited 0.
- Required generation: `make skills-sync` and
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0. The recorder raised both Roundfix skill version fields from
  `0.1.242` to `0.1.243` and wrote the owned-version record; no version or
  mirror was hand-edited.
- After the last guide edit and before adoption,
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 and regenerated the Behavior Surface Record.
- `go test -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce|TestUserGuideLinksResolve|TestTheSkillCoverageMapIsCurrent)$' -count=1`
  exited 0.
- `GOOS=windows go test -c -o <temporary-directory>/skills.test.exe ./skills`
  and `GOOS=windows go test -c -o <temporary-directory>/cli.test.exe ./internal/cli`
  both exited 0; the new test files compile for Windows.
- `git diff --check` exited 0. The path audit found only Task-declared
  changes. Four authoring `SKILL.md` pairs remain byte-identical to HEAD,
  and the pre-existing Project Config bytes remain unchanged before the
  appended section.

The pre-existing task-file change was its Daemon-owned status; it was
preserved. Commands ran through the physical checkout after the configured
symlinked writable root prevented sandbox startup. No Task status, graph,
other Task, Baseline, lock, or prohibited path was edited. No declared
Verification command, repository Verification, commit, push, PR, or external
message was run. Daemon Verification and settlement remain outstanding.
