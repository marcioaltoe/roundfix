### update check

```bash
roundfix update check [--refresh] [--format <text|json>]
```

The read-only Update Check compares the installed Roundfix with the latest
published release, lists each newer release's `### ` titles newest first,
and summarizes this repository's Baseline preview with the installed binary.
Every report exits `0`. Usage errors and invalid configuration exit `2`.

The release listing uses one GitHub CLI request with a five-second timeout.
It is shared across repositories through `~/.roundfix/update-check.json` for
24 hours, including a failed lookup. `--refresh` reads it now. Cache read or
write failures do not block the report. A development build skips release
lookup and cache access. The cache is the command's only write; the repository
and its Baseline state stay unchanged.

`update.mode` in Project Config overrides User Config: `notify` is the default,
and `auto` asks the session to apply the update through the Roundfix skill's
`--update` mode. The command reports the mode and applies nothing itself.

The first output line is one of:

```text
Roundfix <installed> → <latest> is available (update.mode <mode>).
Roundfix <installed> is current (latest release <latest>; update.mode <mode>).
Roundfix <installed>: latest release not reached (<reason>); the session continues.
Roundfix <installed> is a development build; no release comparison (update.mode <mode>).
```

A failed lookup's reason has whitespace collapsed and is limited to 300
characters; no published release reads `no releases published`. When a newer
release exists, the first line is followed by:

```text
Changes since <installed>:
- <version>: <title>; <title>
```

Each newer release has one line, newest first. A release without a `### `
title shows its URL instead. When all 30 releases on the listing's page are
newer than the installed version, the block ends with:

```text
- earlier releases: https://github.com/marcioaltoe/roundfix/releases
```

The Baseline line follows the release report. Inside a Git repository, the
check runs `roundfix baseline update --repo <git root> --format json` in
process as a preview, without any approval or decision flags. Its line is one
of:

```text
Baseline: current with the installed Roundfix.
Baseline: plan ready with the installed Roundfix: <f> file change(s), <h> history unit(s), <k> skill change(s).
Baseline: <state>: <message>
Baseline: <state>.
Baseline: not checked (outside a Git repository).
```

The counts come from `fileChanges`, `history.units`, and the sum of
`skills.outdated`, `skills.missing`, and `skills.retired`. Other states show
the preview's message with whitespace collapsed, or just the state without
a message. The final outside-repository line skips the preview.

When a **Stale User Skill Copy** exists, or a declared phase's skill has a
copy in Claude Code's user-level directory that hides its Declared Skill
Header, this line follows the Baseline line, including for a development
build. Older copies are named inside or outside a Git repository:

```text
User skills: copies that an Agent harness may load in place of the project copy: <skill> <found> (binary <version>) at <path>[ -> <link target>]; <reason>, ...; remove with <remove> && <remove>, or refresh with <refresh> && <refresh>.
```

Reasons and directory order follow [skills check](skills.md). Removals come
before refreshes; the `, or refresh with ...` clause appears only for older
copies. The line does not change `next`. Neither command runs here; the
`--update` Skill Mode asks which to run. When Project Config was not loaded,
only older copies are named.

When the installed authoring files do not match the Authoring Model
Declaration, this line follows `User skills:` when present, otherwise the
Baseline line. Problems follow the [skills check problem order](skills.md)
and are joined by `; `:

```text
Authoring models: the installed files do not match the declaration: <problem>; apply with roundfix skills install --target project.
```

For an authoring skill the repository publishes, the declaration counts as
applied when the installed skill holds no model and no effort.

The `; apply with ...` clause appears only when installation can repair at least
one problem. Links and files Roundfix did not generate need the user's
correction. An applied declaration adds no text line, except the research tools
line below. The comparison writes nothing, changes neither the report's `next`
nor its exit code, and runs only when Project Config was read.

With no `skills.authoring.claude` and no inspection problem, an installed
authoring skill the repository does not publish gets this line in the same
position, followed by the Recommended Authoring Declaration:

```text
Authoring models: the authoring skills are installed without a declaration: write-idea, write-prd, write-techspec, write-tasks; add the Recommended Authoring Declaration below to .roundfixrc.yml, then apply with roundfix skills install --target project.
  skills:
    authoring:
      claude:
        idea:
          model: fable
          reasoning_effort: xhigh
          fallback:
            model: opus
            reasoning_effort: xhigh
        prd:
          model: fable
          reasoning_effort: high
          fallback:
            model: opus
            reasoning_effort: xhigh
        techspec:
          model: fable
          reasoning_effort: high
          fallback:
            model: opus
            reasoning_effort: xhigh
        tasks:
          model: opus
          reasoning_effort: high
          fallback:
            model: sonnet
            reasoning_effort: high
        research:
          model: sonnet
          reasoning_effort: medium
          tools: [WebSearch, WebFetch]
```

The line names only installed authoring skills not published under
`skills/<name>/SKILL.md`, in phase order; a published skill does not count.
The JSON report sets `authoringModels.state` to `undeclared`, with empty
`declared` and `problems` arrays, the `installed` names and the unindented
YAML string in `recommendation`, ending with a newline. Its `next` member is
omitted. Neither this report nor the Doctor warning changes an exit code or
the Update Check's top-level `next`. The recommendation writes nothing.

Review the block and add it to `.roundfixrc.yml`, then run the named install
command yourself. If the file already has `skills:`, merge into that mapping
rather than adding a duplicate key. The block gives the Research Subagent
`WebSearch` and `WebFetch` and no MCP tool. A leftover generated file keeps the
existing problem report without this block. With no declaration, problem or
installed unpublished authoring skill, no text line or JSON member is added.
Outside a Git repository or when Project Config was not read, neither is added.

When `skills.authoring.claude.research` is declared without a `tools` key and
the inspection finds no problem, the check prints this line in the same
position:

```text
Authoring models: research declares no tools, so the Research Subagent cannot read the web; add tools: [WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml, then apply with roundfix skills install --target project.
```

The JSON report adds `recommendedResearchTools`, the tools to add, and keeps
`authoringModels.state` `applied`. An empty list, `tools: []`, and any declared
list add neither. Neither changes an exit code or the top-level `next`, and
nothing is written.

When Project Config was read and no newer Roundfix wrote the repository's
Baseline, the check compares each configured Agent Selection Profile with the
Recommended Profile from the running binary's snapshot. A category that
differs without a Profile Deviation is named as:

```text
Profiles: <category> (<source>) differs from the Recommended Profile of <snapshot>: configured <profile>; recommended <profile>; adopt with roundfix profiles check --apply --scope <source>.
```

A category whose Profile Deviation was declared against an older snapshot is
named as:

```text
Profiles: <category> (<source>) is pinned to an older snapshot (<from>: <reason>); recommended <profile> as of <snapshot>; re-pin with profiles.<category>.deviation.from: <snapshot> in <Project|User> Config, or adopt with roundfix profiles check --apply --scope <source>.
```

The comparison uses the running binary's snapshot, collapses whitespace in
`<reason>` and limits it to 300 characters. A current category or one pinned
to the shipped snapshot adds no `Profiles:` line. The comparison writes
nothing and changes neither `next` nor the exit code. `--apply` adopts every
differing category in its scope, including a category pinned to an older
snapshot. When a newer Roundfix wrote the repository's Baseline, no profile
comparison runs.

When a newer Roundfix wrote this repository's Baseline, the check prints `roundfix: warning: <refusal>; Project Config not read` on stderr, reads no Project Config, and its Baseline line is `Baseline: failed:` followed by the refusal.

After the profile recommendations and before `Next:`, the check reports each
harness whose Flow Command List is current, behind or not readable:

```text
Flow permissions: claude current
Flow permissions: claude behind, 2 to add, 0 to remove — roundfix permissions --harness claude
Flow permissions: claude not readable — roundfix permissions --harness claude
```

The optional JSON part `flowPermissions` holds one object per such harness, in
Claude Code, Codex, OpenCode order, with `harness`, `state` (`current`, `behind`
or `unreadable`), `add` and `remove` counts, and `next`. The preview command in
`next` is empty when current. Add counts include missing Guard Entries. The
part is omitted when none qualifies, including when no list is applied or no
harness is present. It changes neither the top-level `next` nor the exit code.
The update check never applies the Flow Command List.

Only when a release is newer or the Baseline has a plan, the last line names
the configured mode:

```text
Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.
Next: update.mode auto: apply it now with the Roundfix skill's --update mode.
```

With `--format json`, standard output is one `roundfix/update-check/v1`
document:

| Field | Meaning |
| --- | --- |
| `schemaVersion` | `roundfix/update-check/v1` |
| `mode` | Effective `update.mode`: `notify` or `auto` |
| `installed`, `latest` | Normalized installed and latest versions; `latest` is empty when unknown |
| `state` | `update_available`, `current`, `unknown`, or `development` |
| `lookup.source` | `network`, `cache`, or `skipped` |
| `lookup.checkedAt` | RFC 3339 lookup time, omitted when skipped |
| `lookup.error` | Lookup failure text, or an empty string |
| `releases` | Newer releases with `version`, `titles`, and `url`; an empty array otherwise |
| `releasesTruncated` | All 30 listed releases are newer than the installed version |
| `baseline.state` | Preview state, or `not_checked` outside a Git repository |
| `baseline.fileChanges`, `baseline.historyUnits`, `baseline.skillChanges` | Plan counts |
| `baseline.message` | Preview message with whitespace collapsed |
| `userSkills` | Each reported user-level copy with `skill`, `path`, `resolved` (the link target, omitted for a plain folder), `found`, `embedded`, `refresh`, `remove`, `stale` and optional `hidesDeclaredHeader`; omitted when there is none |
| `userSkills[].remove` | Command that removes the copy, before the refresh advice in text reports |
| `userSkills[].stale` | Whether the copy's comparable version is older than the binary; always present, including `false` |
| `userSkills[].hidesDeclaredHeader` | `true` when Claude Code's user-level copy hides a declared phase's header, whatever its version; never true for a skill the repository publishes; omitted otherwise |
| `authoringModels` | Present after `userSkills` when a key is declared, a problem exists, or an installed unpublished authoring skill has no declaration; omitted outside Git or when Project Config was not read |
| `authoringModels.state` | `applied`, `unapplied`, or `undeclared` |
| `authoringModels.harness` | `claude` |
| `authoringModels.declared` | Declared keys in phase order, with `research` last; an empty array for leftover generated files or `undeclared` |
| `authoringModels.problems` | Problem texts in report order; always an array, empty when `applied` or `undeclared` |
| `authoringModels.next` | `roundfix skills install --target project` when a problem is fixable by installation; omitted otherwise |
| `authoringModels.installed` | Installed authoring skills the repository does not publish, in phase order; present only when `undeclared` |
| `authoringModels.recommendation` | Recommended Authoring Declaration as YAML starting at `skills:` and ending with a newline, with `tools: [WebSearch, WebFetch]` under `research`; present only when `undeclared` |
| `authoringModels.recommendedResearchTools` | The research tools of the Recommended Authoring Declaration to add; present only when `applied` and the declared research entry has no `tools` key |
| `recommendations` | Optional profile comparison, omitted when Project Config was not read, the comparison failed, or no category differs |
| `recommendations.snapshot` | Snapshot from the running binary used for the comparison |
| `recommendations.categories[].category` | Agent Work Category whose configured profile differs |
| `recommendations.categories[].state` | `differs` or `pinned_to_older_snapshot` |
| `recommendations.categories[].source` | `user` or `project`, the configuration source of the category |
| `recommendations.categories[].configured` | Configured Profile with `preferred` and `fallbacks` |
| `recommendations.categories[].recommended` | Recommended Profile with `preferred` and `fallbacks` |
| `recommendations.categories[].deviation` | Profile Deviation, present only when `state` is `pinned_to_older_snapshot` |
| `recommendations.categories[].repin` | Configuration key and snapshot that re-pin an older deviation, present only when `state` is `pinned_to_older_snapshot` |
| `recommendations.categories[].adopt` | `roundfix profiles check --apply --scope <user|project>` for the category's source |
| `flowPermissions` | Optional harness list with `harness`, `state`, `add`, `remove` and `next`; omitted when empty |
| `next` | `apply` for a newer release or a `plan_ready` Baseline, otherwise `none` |

A repository that wants the harness itself to run the check may add a SessionStart hook whose command is `roundfix update check`: in Claude Code under `hooks` → `SessionStart` with the matcher `startup` in `.claude/settings.json`, whose plain standard output reaches the session's context, and in Codex the same entry in `.codex/hooks.json`, which runs only in a trusted project after the hook is reviewed with `/hooks`. A hook is repository tooling configuration, so add it only with the maintainer's authorization; the session still sees the report once, because the instruction skips a session whose context already holds it.
