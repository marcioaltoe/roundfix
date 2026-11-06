---
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
prd: _prd.md
created: 2026-11-06
---

# A session starts with an Update Check, and the update mode applies it — Technical Spec

## Executive Summary

Releases reached adopters as notes in their ai-memory inboxes, and nobody read
them. This Spec moves "what changes in this repository" into the session
itself. A new read-only command, `roundfix update check`, compares the
installed Roundfix with the latest published release from one `gh api`
release listing, cached for 24 hours in Roundfix Home, prints the `### ` title
of each newer release (the `CHANGELOG.md` section titles the release workflow
publishes as GitHub Release notes), and summarizes this repository's
`roundfix baseline update` preview, run in process without writing. A Project
Config key, `update.mode`, chooses `notify` (default) or `auto`, and the report's
`Next:` line says which. The Roundfix skill gains an `update` Skill Mode that
runs the check, upgrades the binary, previews the Baseline plan with the new
binary and applies exactly that plan by its Plan Digest, asking before deleting
files or touching uncommitted work. One mandatory Baseline clause tells every
session to run the check before its first task. ADR-0371 records the decision.

The trade-off accepted: the clause is an instruction, so a session that ignores
it shows nothing. A harness hook would run deterministically, but writing hooks
into adopters' harness configuration is tooling configuration in three
formats, and Codex skips an untrusted project hook without a word. The command
guide shows the hook a repository may add, and the operational commands' daily
Version Freshness line stays as a second signal.

## Project Constraints

- Identifier strategy: not applicable — the only new identifiers are the schema
  versions `roundfix/update-check/v1` and `roundfix/update-check-cache/v1` of
  the command's own JSON document and cache file. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint or credential is
  added; the release listing runs `gh api` with the user's own GitHub CLI
  authentication, as `roundfix upgrade` does, and tests use a fake runner.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0371 (new) records the decision:
  "A session starts with an Update Check, and the Roundfix skill applies the
  update". ADR-0320 defines the Skill Modes this Spec extends with `update`,
  and its rule of up to three questions per call governs the mode's questions
  and stays unchanged; ADR-0320: "An Agent asks at most three questions per
  call or message". ADR-0354 keeps every announcement in the user's own
  instructions, and the Update Check sends no message, so it stays unchanged;
  ADR-0354: "and Roundfix's code sends no message". ADR-0338, which cites
  ADR-0320, does not apply: it governs how the `--follow-up` mode takes a
  handoff and settles messages, and the `update` mode takes no handoff and
  consumes no message. ADR-0353, which cites ADR-0338, does not apply
  either: it holds Waiting Items back from the follow-up list, which the
  `update` mode never reads. No accepted ADR is amended. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 4 answers of 2026-11-06 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files tasks 02 and 03 change, the CLI
  help test task_02 changes, and the Baseline core module, the Standard
  TypeScript profile and its formatter golden, this repository's agent
  instructions guide and its setup manifest that task_04 changes; no Makefile, workflow, module or lint
  configuration changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-06 at `9a285d77`; `_prd.md` → Validation records the runs.

1. `roundfix update` is not a command (`roundfix: unknown command "update"`,
   exit `2`). The Upgrade Command (`internal/cli/upgrade.go`,
   `performUpgrade`) reads `releases/latest` through `app.LatestRelease`
   (`internal/app/release.go`), which runs `gh api` through the package
   variable `runGitHubAPI` and maps an HTTP 404 to `app.ErrNoReleases`.
2. `maybeReportVersionFreshness` (`internal/cli/upgrade.go`) runs before the
   operational commands and `implement`: it caches `releases/latest` for 24
   hours in `<home>/.roundfix/version-check.json` with a two-second timeout,
   writes the cache on failure too, skips development builds, and prints only
   `roundfix <v> is behind latest <v>; run roundfix upgrade`.
3. Project Config and User Config decode into `configOverlay`
   (`internal/config/config.go`) with known fields only, so `update:` is
   refused at load; `Validate` and the overlay merge own every key.
4. `runBaselineUpdateCommand` (`internal/cli/baseline_update.go`) computes and
   prints the managed refresh preview as `roundfix/baseline-update-result/v1`
   and writes nothing without `--yes` or `--confirm-plan`. Its `state` is
   `current`, `plan_ready` or another state with `category`, `message` and
   `nextAction`; `fileChanges`, `history.units`, `skills.outdated`,
   `skills.missing` and `skills.retired` carry the plan.
5. The core module's `rule.core.cross-project-continuity` (version 1) holds
   four mandatory clauses about ai-memory; no Baseline module renders a harness
   hook, and the Setup Manifest's managed artifacts are guides and root blocks.
6. The Roundfix skill's Modes paragraph, pinned word for word by
   `TestTheRoundfixSkillReadsItsModeFromTheRequest`
   (`skills/roundfix_modes_test.go`), names `follow-up`, `status` and
   `unattended`; `references/modes.md` has a section per mode.

## System Architecture

- `internal/config`: the `Update` section and its validation (task_01).
- `internal/app`: `ListReleases`, one `gh api` release listing (task_02).
- `internal/cli/update_check.go`: the `update` group, the cache, the
  comparison, the in-process Baseline preview and the text and JSON reports
  (task_02).
- The Roundfix skill: the `setup` reference documents the command (task_02);
  the `modes` reference and `SKILL.md` add the `update` mode (task_03).
- `internal/baseline/assets/modules/core.json`: the session-start clause and
  its derived files, and this repository's managed refresh (task_04).

## Implementation Design

### Interfaces

```go
// internal/config
type Update struct{ Mode string } // "notify" or "auto"
type updateOverlay struct{ Mode *string `yaml:"mode"` }

// internal/app
type Release struct {
	Version string   // normalized, no leading v
	Titles  []string // "### " headings of the body, in order
	URL     string   // html_url
}
const ReleaseListingPageSize = 30
func ListReleases(ctx context.Context) ([]Release, error) // newest first

// internal/cli
type updateCheckDependencies struct {
	now            func() time.Time
	currentVersion func() string
	listReleases   func(context.Context) ([]app.Release, error)
}
```

`commandDependencies` gains `updateCheck updateCheckDependencies`, defaulted
like `versionFreshness`; tests replace it through the existing
command-dependency context.

### Data Models

The cache, `<home>/.roundfix/update-check.json`, mode `0644`:

```json
{
  "schemaVersion": "roundfix/update-check-cache/v1",
  "checkedAt": "2026-11-06T15:00:00Z",
  "releases": [{"version": "0.136.1", "titles": ["..."], "url": "https://github.com/marcioaltoe/roundfix/releases/tag/v0.136.1"}],
  "error": ""
}
```

### API Contracts

1. API Contract: `update.mode`. A top-level `update:` mapping with one key,
   `mode`, accepted in User Config and Project Config, Project Config
   winning; built-in default `notify`. Any value other than `notify` or `auto`
   is refused by `Validate` with `update.mode must be notify or auto, got
   "<value>"`, inside the existing load error of the file. The generated
   init templates (`DefaultConfigYAML`, `DefaultProjectConfigYAML`) keep their
   bytes. `docs/user-guide/configuration.md` adds one row to the General
   settings table, word for word:
   ``| `update.mode` | `notify` | Selects what a session does with the report of `roundfix update check`: `notify` reports a newer release and this repository's Baseline plan and applies nothing; `auto` asks the session to apply them through the Roundfix skill's `--update` mode before its first task. Accepted in User and Project Config, with Project Config winning. |``
2. API Contract: `app.ListReleases`. Runs `runGitHubAPI` with
   `repos/marcioaltoe/roundfix/releases?per_page=30`, decodes `tag_name`,
   `body`, `draft`, `prerelease` and `html_url`, skips drafts and prereleases,
   normalizes each version, takes as titles every body line that starts with
   `### `, without that marker and trimmed, and returns the releases sorted
   newest first by `CompareVersions`. A 404 or an empty result returns
   `ErrNoReleases`; any other failure returns `list releases: <error>`.
3. API Contract: the cache. A check reads `update-check.json` and uses it when
   its schema version matches and `0 <= now - checkedAt < 24h`; otherwise, or
   with `--refresh`, it calls `listReleases` under a five-second timeout and
   writes the cache with `checkedAt` = now and either the releases or the
   failure text in `error` (releases empty). Cache read or write errors are
   ignored. A development build (`app.IsDevelopmentVersion`) neither reads,
   calls nor writes. The latest release is the first cached release.
4. API Contract: the text report, on standard output, every report exiting
   `0`. Line 1 is one of:
   - `Roundfix <installed> → <latest> is available (update.mode <mode>).`
   - `Roundfix <installed> is current (latest release <latest>; update.mode <mode>).`
   - `Roundfix <installed>: latest release not reached (<reason>); the session continues.`
   - `Roundfix <installed> is a development build; no release comparison (update.mode <mode>).`

   `<reason>` is the lookup error with whitespace collapsed, at most 300
   characters; `no releases published` for `ErrNoReleases`. Only when a
   release is available, `Changes since <installed>:` follows, then one line
   per release newer than the installed one, newest first:
   `- <version>: <title>; <title>` with its titles, or `- <version>: <url>`
   without one; and, when the listing held `ReleaseListingPageSize`
   releases all newer than the installed one,
   `- earlier releases: https://github.com/marcioaltoe/roundfix/releases`.
   Then the Baseline line (API Contract 6). Last, only when a release is
   available or the Baseline line is `plan ready`, the `Next:` line:
   - notify: `Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.`
   - auto: `Next: update.mode auto: apply it now with the Roundfix skill's --update mode.`
5. API Contract: `--format json` prints one `roundfix/update-check/v1`
   document: `schemaVersion`; `mode`; `installed`; `latest` (empty when
   unknown); `state` (`update_available`, `current`, `unknown`,
   `development`); `lookup` with `source` (`network`, `cache` or `skipped`),
   `checkedAt` (RFC 3339, omitted when skipped) and `error`; `releases` (the
   newer releases as `version`, `titles`, `url`; an empty array otherwise);
   `releasesTruncated`; `baseline` with `state`, `fileChanges`,
   `historyUnits`, `skillChanges` and `message`; and `next` (`apply` when a
   release is available or the Baseline state is `plan_ready`, else `none`).
6. API Contract: the Baseline line. Inside a Git repository the check runs
   `runBaselineUpdateCommand` in process with `--repo <git root> --format
   json`, its standard error discarded, and decodes the result:
   - `current`: `Baseline: current with the installed Roundfix.`
   - `plan_ready`: `Baseline: plan ready with the installed Roundfix: <f> file change(s), <h> history unit(s), <k> skill change(s).`,
     where `<f>` counts `fileChanges`, `<h>` counts `history.units` and `<k>`
     counts `skills.outdated`, `skills.missing` and `skills.retired`;
   - any other state: `Baseline: <state>: <message>` with the result's
     message, whitespace collapsed, or `Baseline: <state>.` without one.
   Outside a Git repository: `Baseline: not checked (outside a Git repository).`
   with JSON `baseline.state` `not_checked`. The preview never receives
   `--yes`, `--confirm-plan`, `--adopt-suggested` or `--decision`.
7. API Contract: the command surface. `roundfix update --help`, `roundfix
   update check --help` and `-h` print the usage of Surface Transcript 1.
   `roundfix update` alone, an unknown subcommand, an unexpected argument or a
   `--format` other than `text` or `json` print the preflight failure of
   `printPreflightFailure` and exit `2`, the bare group with the reason
   `expected subcommand "check"` and the pointer `Run 'roundfix update
   --help' for usage.` An invalid configuration prints the same failure for
   `update check` and exits `2`. The root usage gains the line
   `  roundfix update check [--refresh] [--format <text|json>]` after
   `  roundfix upgrade [--check]`, and its command list gains
   `  update     Report a newer release and this repository's Baseline plan`
   after the `upgrade` row.
8. API Contract: the documents.
   - `docs/user-guide/commands/update.md` (new), headed `### update check`,
     states the command line, the four first lines of API Contract 4, the
     Baseline lines of API Contract 6, the `Next:` lines, the JSON fields of
     API Contract 5, the 24-hour cache in `~/.roundfix/update-check.json` and
     `--refresh`, and this paragraph word for word: "A repository that wants
     the harness itself to run the check may add a SessionStart hook whose
     command is `roundfix update check`: in Claude Code under `hooks` →
     `SessionStart` with the matcher `startup` in `.claude/settings.json`,
     whose plain standard output reaches the session's context, and in Codex
     the same entry in `.codex/hooks.json`, which runs only in a trusted
     project after the hook is reviewed with `/hooks`. A hook is repository
     tooling configuration, so add it only with the maintainer's
     authorization; the session still sees the report once, because the
     instruction skips a session whose context already holds it."
     `docs/user-guide/commands.md` gains the row
     ``| `update` | [command guide](commands/update.md) |`` after the
     `upgrade` row, and `docs/user-guide/usage.md` gains
     ``| `update check` | Report a newer release and this repository's Baseline plan (read-only) |``
     after the `upgrade` row.
   - The Roundfix `setup` reference gains, before `## Pending History after
     upgrade`, a section `## Update check` with this paragraph word for word:
     "`roundfix update check [--refresh] [--format <text|json>]` is the
     read-only Update Check a session runs before its first task, as the
     Baseline's agent instructions say. It compares the installed Roundfix
     with the latest published release, lists the `### ` titles of each newer
     release's notes, newest first, and summarizes this repository's
     `roundfix baseline update` preview with the installed binary. The release
     listing is one `gh api` call cached for 24 hours in
     `~/.roundfix/update-check.json` and shared by every repository on the
     machine; `--refresh` reads it now, and a failed lookup is one line ending
     `the session continues.` and is cached for the same 24 hours. Every
     report exits `0`; a usage error or an invalid configuration exits `2`.
     `update.mode` (Project Config over User Config) selects `notify`, the
     default, or `auto`; the `Next:` line names the `--update` Skill Mode when
     a release is newer or the Baseline has a plan. JSON uses
     `roundfix/update-check/v1` with `state` (`update_available`, `current`,
     `unknown` or `development`), `next` (`apply` or `none`), `releases`,
     `lookup` and `baseline`."
     The `SKILL.md` index row of `setup` lists `` `update` `` after
     `` `upgrade` ``.
   - The Roundfix `SKILL.md` Modes paragraph becomes, word for word: "A
     request names a Skill Mode with the first word after the skill's name:
     `follow-up`, `status`, `unattended` or `update`, with or without a
     leading `--`, in every harness: `/roundfix follow-up` or `/roundfix
     --follow-up` in Claude Code, `$roundfix status` or `$roundfix --status`
     in Codex, and `@roundfix unattended`, `@roundfix --unattended` or
     `/roundfix unattended` in OpenCode. A Skill Mode is not the CLI command
     `roundfix follow-up`, which the follow-up mode runs as one of its steps.
     Read the mode from the request text; a request whose first word after
     the skill's name is anything else, such as `/roundfix watch 123`, uses
     the skill as the rest of this file describes." It gains, after the
     `--unattended` bullet, "- `--update`: run the Update Check, show what a
     newer Roundfix and this repository's Baseline plan change, and apply
     them, asking before a step that deletes files." The frontmatter
     `description` names "its `--follow-up`, `--status`, `--unattended` and
     `--update` Skill Modes", and the index row of `modes` names the
     `--update` Skill Mode beside the other three.
   - The `modes` reference: "How a mode is chosen" names `follow-up`,
     `status`, `unattended` or `update`; its sentence "Each mode runs
     Roundfix's read-only follow-up command; when the installed Roundfix lacks
     it, the mode says so, names `roundfix upgrade` and stops." becomes "The
     `--follow-up`, `--status` and `--unattended` modes run Roundfix's
     read-only follow-up command; when the installed Roundfix lacks it, the
     mode says so, names `roundfix upgrade` and stops. The `--update` mode
     runs `roundfix update check` instead and, when the installed Roundfix
     lacks it, runs `roundfix upgrade` first." A new last section
     `## --update` holds, word for word:

     ```text
     Run this mode in the repository the session works in, never in a Daemon-assigned Run.

     1. Check: run `roundfix update check --refresh --format json`. When the installed Roundfix does not know the command, run `roundfix upgrade` and then the check again. Report the installed and latest versions, each newer release with its titles, and the Baseline line. When `state` is `unknown` or `development` and `baseline.state` is not `plan_ready`, report that nothing can be applied now and stop.
     2. Upgrade: when `state` is `update_available`, run `roundfix upgrade`. When it fails, report its manual fallback and stop without applying a Baseline plan.
     3. Preview: run `roundfix baseline update --repo <root> --format json` with the binary now installed and show its state, each file change, each history unit with its paths, each skill it refreshes, restores or deletes, and its `planDigest`. When the state is `current`, report the versions and stop. When it needs a new decision or an adoption, ask the user as the stop-and-ask clause of `docs/agents/agent-instructions.md` says, and never pass `--adopt-suggested` unasked.
     4. Apply: apply exactly the previewed plan with `roundfix baseline update --repo <root> --confirm-plan <planDigest> --format json`. Ask first, in one question with the recommended option first, when the plan deletes files, that is when `history.units` or `skills.retired` is not empty: apply the whole plan, apply it without history (preview again with `--no-history` and confirm that digest), or stop. Ask first, too, when `git status --porcelain` shows uncommitted changes. When the mode runs at session start because `update.mode` is `auto`, apply no Baseline plan over uncommitted changes; report it for a later `--update`.
     5. Report the versions before and after, the applied Plan Digest and the changed files. This mode never commits, pushes or opens a pull request; committing follows the repository's own rules.
     ```
9. API Contract: the session-start clause. `rule.core.cross-project-continuity`
   rises to version 2 and gains, after
   `clause.core.consume-messages-after-recording`, the mandatory clause
   `clause.core.check-for-updates-at-session-start` with this guidance, word
   for word: "At the start of each session, before its first task, run
   `roundfix update check` once and relay its report in a few lines, unless
   the session's context already holds that report. The check reads the
   network at most once a day and never holds up the session: an unreachable
   release lookup or a failed check is one line, and work continues. With
   `update.mode` set to `notify`, the default, relay the report and apply
   nothing; the user applies an update with the Roundfix skill's `--update`
   mode. With `update.mode` set to `auto`, run that mode before the first task
   when the report's `Next:` line asks for it; the mode still asks before a
   step that deletes files and applies no Baseline plan over uncommitted
   changes. In a Daemon-assigned turn, skip the check." It replaces no clause.
10. API Contract: the glossary. `CONTEXT.md` adds, after **Upgrade Command**:

    ```text
    **Update Check**:
    The read-only report of `roundfix update check` that a session runs before its first task: the installed Roundfix against the latest published release, the titles of each newer release from its published notes, and the repository's Baseline plan with the installed binary. It reads the release listing at most once a day per machine, reports a lookup that fails in one line, and changes nothing (ADR-0371).
    _Avoid_: Release note, version freshness line, upgrade check
    ```

    and **Skill Mode** becomes:

    ```text
    **Skill Mode**:
    A way of running the Roundfix skill that the user's request selects with the first word after the skill's name, `follow-up`, `status`, `unattended` or `update`, with or without a leading `--`; the skill reads it from the request text in every harness, it is not the CLI command `roundfix follow-up`, and a request whose first word names no mode uses the skill as before (ADR-0320). The `update` mode runs the Update Check, upgrades the binary and applies the Baseline plan it showed by its Plan Digest, asking before a step that deletes files (ADR-0371).
    _Avoid_: Command flag, CLI mode, Follow Mode
    ```

    In both entries each parenthesized `ADR-NNNN` is written as the
    glossary's Markdown link to that ADR file under `docs/adr/`, as the
    neighbouring entries write it. Both texts, with those links, passed
    `TestGlossaryPhrasesAreDefinedTerms` on a scratch copy.

### Surface Transcripts

Every transcript runs `bin/roundfix` with a fresh temporary `HOME` and
`TMPDIR`, outside a Git repository, with a fake `gh` first on `PATH` that
prints a fixed release listing (transcripts 2 and 3) or exits `1` with
`error connecting to api.github.com` (transcript 4). The listing holds
`99.0.0` (body `### The first fixture release`), `98.0.0` (body
`### The second fixture release` and `### Its second title`), a draft
`100.0.0` and a prerelease `99.1.0-rc.1`. `<installed>` is the built
binary's version, lower than `98.0.0`.

1. Surface Transcript: help.

   ```transcript
   $ roundfix update check --help
   stdout:
   Usage:
     roundfix update check [--refresh] [--format <text|json>]

   Reports, read-only, whether a newer Roundfix release exists, the title of each
   release since the installed version, and this repository's Baseline plan with
   the installed binary. The release listing comes from the GitHub CLI and is
   cached for 24 hours in Roundfix Home; a lookup that fails is reported in one
   line. Every report exits 0. update.mode in Project Config or User Config
   selects notify (the default) or auto, which the Next: line names.

   Commands:
     check  Report the latest release, what it changes, and this repository's Baseline plan

   Options:
     --refresh  Read the release listing now instead of the cached one
     --format   Output format: text (default) or json
   stderr:
   exit: 0
   ```

2. Surface Transcript: a newer release, `notify`.

   ```transcript
   $ roundfix update check
   stdout:
   Roundfix <installed> → 99.0.0 is available (update.mode notify).
   Changes since <installed>:
   - 99.0.0: The first fixture release
   - 98.0.0: The second fixture release; Its second title
   Baseline: not checked (outside a Git repository).
   Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.
   stderr:
   exit: 0
   ```

3. Surface Transcript: the same with `update: {mode: auto}` in User Config.

   ```transcript
   $ roundfix update check
   stdout:
   Roundfix <installed> → 99.0.0 is available (update.mode auto).
   Changes since <installed>:
   - 99.0.0: The first fixture release
   - 98.0.0: The second fixture release; Its second title
   Baseline: not checked (outside a Git repository).
   Next: update.mode auto: apply it now with the Roundfix skill's --update mode.
   stderr:
   exit: 0
   ```

4. Surface Transcript: the lookup fails.

   ```transcript
   $ roundfix update check
   stdout:
   Roundfix <installed>: latest release not reached (<reason>); the session continues.
   Baseline: not checked (outside a Git repository).
   stderr:
   exit: 0
   ```

5. Surface Transcript: the bare group.

   ```transcript
   $ roundfix update
   stdout:
   stderr:
   Preflight failed

   Reason:
     expected subcommand "check"

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix update --help' for usage.
   exit: 2
   ```

### Invariants

```text
1. roundfix update check writes no repository file and no Baseline state; its only write is the cache file in Roundfix Home.
2. Two checks less than 24 hours apart make at most one release lookup unless the second passes --refresh; a failed lookup counts as a lookup.
3. Every report of the check exits 0, whatever the lookup or the Baseline preview returned.
4. A release marked draft or prerelease is never the latest release and never listed.
5. The init templates and every existing configuration key keep their bytes and meaning.
6. The Baseline preview the check runs never carries --yes, --confirm-plan, --adopt-suggested or --decision.
```

### Declared breaks

Measured with a prototype of task_04 on a scratch copy at `9a285d77`, and by
reading the tests that pin the changed surfaces for tasks 01 to 03.

- Behavior, task_01: `update:` in Project Config or User Config is accepted
  instead of refused; the `config keys` Behavior Surface changes.
- Behavior, task_02: the root usage gains one usage line and one command row;
  `roundfix update` exists; two Behavior Surfaces (`command update check`,
  `guide docs/user-guide/commands/update.md`) and their Skill Coverage Map
  entries are added.
- Behavior, task_04: every adopter's next managed refresh adds the clause to
  `docs/agents/agent-instructions.md`; this repository's guide and setup
  manifest change through the same refresh.
- Tests changed, task_02: `TestRunCommandHelp` (the help table in
  `internal/cli/cli_test.go`) gains an `update check` row; every other test on
  the default branch when the Task runs that pins the root usage text gains
  the new line.
- Tests changed, task_03: `TestTheRoundfixSkillReadsItsModeFromTheRequest`
  (`skills/roundfix_modes_test.go`) pins the new Modes paragraph, the
  `--update` bullet and the new "How a mode is chosen" phrase.
- Tests changed, task_04: `characterizedBaselineForce`
  (`internal/baseline/clause_characterization_test.go`) gains
  `clause.core.check-for-updates-at-session-start` as `mandatory`; without it
  `TestBaselineClauseForceIsCharacterized` and three sibling tests fail, as
  the prototype showed.
- No test is removed or renamed, so the coverage record is not re-recorded.

## Coverage Map

- Goal 1, Core Feature 2, Success Metrics 1 and 3 → API Contracts 2 to 7,
  `TestUpdateCheckTranscripts`, `TestUpdateCheckSummarizesTheBaselinePreview`,
  `TestListReleasesReadsTitlesAndSkipsDraftsAndPrereleases`, Invariants 1, 3, 4
  and 6.
- Goal 2, Success Metric 2 → API Contract 3,
  `TestUpdateCheckReadsTheNetworkAtMostOncePerDay`, Invariant 2.
- Goal 3, Core Feature 3, Success Metric 5 → API Contract 8 (`modes`),
  `TestTheUpdateModeAppliesOnlyThePlanItShowed`.
- Goal 4, Core Feature 1, Success Metric 4 → API Contract 1,
  `TestUpdateModeIsNotifyOrAuto`, Invariant 5.
- Goal 5, Core Feature 4, Success Metric 6 → API Contract 9,
  `TestTheSessionStartClauseRunsTheUpdateCheck`, the managed refresh.
- Core Feature 5, Success Metric 7 → API Contracts 8 and 10, the
  documentation, skill and glossary contracts.

Deviations from the PRD: None.

## Absence Evidence

1. The `update` command, `ListReleases` and the cache file — command:
   `grep -rn '"update"\|ListReleases\|update-check' internal/cli internal/app --include='*.go'`;
   result: no match; `bin/roundfix update check` prints `roundfix: unknown
   command "update"`.
2. The `update.mode` key — command:
   `grep -rn 'yaml:"update"\|update\.mode' internal/config docs/user-guide`;
   result: no match; a Project Config with `update:` is refused.
3. The `update` Skill Mode — command:
   `grep -rn -- '--update' .agents/skills/roundfix skills/roundfix`;
   result: no match.
4. A session-start clause or harness hook in the Baseline — command:
   `grep -rn 'SessionStart\|update check\|session start' internal/baseline/assets/modules`;
   result: no match.
5. The new test files — command:
   `ls internal/config/update_mode_test.go internal/cli/update_check_test.go internal/baseline/update_check_clause_test.go docs/user-guide/commands/update.md`;
   result: no such files.

## ADR Alternatives

- ADR-0371: the decision this Spec records; it rejects a Baseline-shipped
  harness hook, the key as a Setup Manifest decision, the full release notes,
  applying from the command in `auto`, an uncached check and keeping the
  inbox notes, and the design takes none of them.
- ADR-0320: it rejected one question per call, four per call and each tool's
  own limit; the `update` mode asks at most one question per call in its apply
  step and takes none of the rejected options.
- ADR-0338: it rejected consuming messages only in one step, leaving an
  unhandled message pending and listing a handoff without accepting it; the
  `update` mode takes no handoff and hands in no message, so the design takes
  none of them and the record does not apply.
- ADR-0353: it rejected listing Waiting Items by default; the `update` mode
  reads no follow-up list, so the design takes none of its alternatives and
  the record does not apply.
- ADR-0354: it rejected writing an assistant's format into the skills, a
  Baseline clause or configuration key for an announcement destination,
  announcing from a Run and a report line without a definition; the Update
  Check is a read-only report, not an announcement, sends nothing and runs in
  no Run, so the design takes none of them.

## Integration Points

The release listing calls the GitHub CLI (`gh api`) with the user's own
authentication; the check survives its absence, a missing login or no
network, in one line. No test or Verification reaches it: tests replace
`runGitHubAPI` or the command dependency, and the QA puts a fake `gh` first on
`PATH`.

## Testing Approach

Every Go test uses temporary homes and repositories, the fake release lister
or a replaced `runGitHubAPI`, and a fixed clock; none reaches the network,
the real `~/.roundfix` or `gh`. Each package's rule on `t.Parallel()` applies.

- `internal/config`, `TestUpdateModeIsNotifyOrAuto`: rows `absent`,
  `project auto`, `user auto`, `project over user`, `invalid project` and
  `invalid user` (the refusal text), and `templates unchanged`.
- `internal/app`, `TestListReleasesReadsTitlesAndSkipsDraftsAndPrereleases`
  (order, titles, a body without a title, draft and prerelease skipped, the
  endpoint string) and `TestListReleasesMapsGitHub404ToNoReleases`.
- `internal/cli`, `TestUpdateCheckTranscripts`: rows `available notify`,
  `available auto`, `current`, `unknown`, `development`, `truncated`,
  `untitled release` and `json`, each asserting the whole standard output,
  empty standard error and exit `0`; `TestUpdateCheckReadsTheNetworkAtMostOncePerDay`:
  rows `second check within a day`, `a day later`, `refresh` and
  `failed lookup`, counting lister calls; `TestUpdateCheckSummarizesTheBaselinePreview`:
  rows `current` (`newBaselineUpdateRepository`), `plan ready`
  (`staleBaselineUpdateRepository`), `no manifest`, `outside repository` and
  `writes nothing` (the stale repository's tracked files and status
  unchanged after the check); `TestUpdateCheckUsage`: rows `bare group`,
  `unknown subcommand`, `unexpected argument`, `invalid format`, `invalid
  configuration` and `help`.
- `skills`, `TestTheUpdateModeAppliesOnlyThePlanItShowed` and the changed
  `TestTheRoundfixSkillReadsItsModeFromTheRequest`.
- `internal/baseline`, `TestTheSessionStartClauseRunsTheUpdateCheck`: the
  rendered agent instructions of a resolved profile contain the clause; the
  force characterization entry.

Sabotages the QA runs, each on its own temporary copy of the audited head:

1. `readUpdateCheckCache` ignoring freshness (always stale):
   `TestUpdateCheckReadsTheNetworkAtMostOncePerDay` fails on `second check
   within a day`.
2. A failed lookup not written to the cache: the same test fails on `failed
   lookup`.
3. `--refresh` ignored: the same test fails on `refresh`.
4. Drafts kept in `ListReleases`:
   `TestListReleasesReadsTitlesAndSkipsDraftsAndPrereleases` fails.
5. The `Next:` line ignoring the mode: `TestUpdateCheckTranscripts` fails on
   `available auto`.
6. `--yes` added to the in-process preview:
   `TestUpdateCheckSummarizesTheBaselinePreview` fails on `writes nothing`.
7. `Validate` accepting any `update.mode`: `TestUpdateModeIsNotifyOrAuto`
   fails on `invalid project`.
8. The clause removed from `core.json` (with derived files regenerated):
   `TestTheSessionStartClauseRunsTheUpdateCheck` fails.
9. The `--confirm-plan` step removed from the `modes` reference:
   `TestTheUpdateModeAppliesOnlyThePlanItShowed` fails.

## Build Order

1. task_01: `update.mode` in `internal/config` and the configuration guide.
2. task_02: `ListReleases`, the `update check` command, its guide and the
   `setup` reference, the coverage map, the glossary entry; depends on step 1
   (the command reads the key).
3. task_03: the `update` Skill Mode in the `modes` reference and `SKILL.md`,
   the changed glossary entry; depends on step 2 (the mode runs the command,
   and both change the Roundfix `SKILL.md` pair and its version record).
4. task_04: the core clause, its derived files and this repository's managed
   refresh; depends on step 3 (the clause names the `--update` mode).
5. task_05: the QA gate; depends on step 4.

No other Spec is required. Specs 0360, 0362 and 0363 deliver first; the
shared files are named in `_prd.md` → Prerequisites, and each Task regenerates
derived files from the tree it runs on.

## Risks & Considerations

- An instruction can be skipped: a session that ignores the clause shows
  nothing until the next one; the operational commands' Version Freshness line
  still names the newer release, and the guide shows an optional hook.
- The Baseline preview of a repository with much Pending History takes longer
  than the 0.4 s measured here; the check still runs it once per session.
- `auto` replaces the machine's binary for every repository, and a running
  Delivery Queue keeps its own executable until it restarts.
- The release titles are only as good as the `### ` headings of
  `CHANGELOG.md`; a release without one shows its link.
- Adopters see nothing until their binary knows `update check`: the first
  update after this release is still a manual `roundfix upgrade`.

## Glossary

- adds: **Update Check**
- changes: **Skill Mode**

## Decisions

Autonomous route: no person present; decisions taken from `_prd.md` → Decisions, the operator's cycle AZ batch 4 brief for Spec 0361 and the characterization above.

- One command group, `update`, with one subcommand, `check`, modeled on the
  `storage` group: a read-only name beside the `upgrade` installer, which
  stays unchanged.
- The in-process preview reuses `runBaselineUpdateCommand` and its JSON
  result, so the Baseline line can never differ from what `roundfix baseline
  update` prints.
- The cache is separate from the Version Freshness cache, so neither
  command's contract moves.
- The clause joins `rule.core.cross-project-continuity` rather than a new
  rule, keeping one coverage area for what a session does when it starts.
