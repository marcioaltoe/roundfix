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

When a **Stale User Skill Copy** exists, this line follows the Baseline line inside or outside a Git repository, including for a development build:

```text
User skills: copies older than this binary: <skill> <found> (binary <version>) at <path>[ -> <link target>], ...; refresh with <refresh> && <refresh>.
```

It does not change `next`; the `--update` Skill Mode asks before running a refresh.

When a newer Roundfix wrote this repository's Baseline, the check prints `roundfix: warning: <refusal>; Project Config not read` on stderr, reads no Project Config, and its Baseline line is `Baseline: failed:` followed by the refusal.

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
| `userSkills` | Each Stale User Skill Copy with `skill`, `path`, `resolved` (the link target, omitted for a plain folder), `found`, `embedded` and `refresh`; omitted when there is none |
| `next` | `apply` for a newer release or a `plan_ready` Baseline, otherwise `none` |

A repository that wants the harness itself to run the check may add a SessionStart hook whose command is `roundfix update check`: in Claude Code under `hooks` → `SessionStart` with the matcher `startup` in `.claude/settings.json`, whose plain standard output reaches the session's context, and in Codex the same entry in `.codex/hooks.json`, which runs only in a trusted project after the hook is reviewed with `/hooks`. A hook is repository tooling configuration, so add it only with the maintainer's authorization; the session still sees the report once, because the instruction skips a session whose context already holds it.
