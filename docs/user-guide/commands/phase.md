### phase

```bash
roundfix phase start <phase> [--spec <slug>] [--model <identifier>] [--effort <level>] [--author-subagent]
```

Records the Supervisor Phase this session starts. The command appends one
phase entry to the [Phase Record](../usage.md#phase-record), resolves the Git
root of the working directory, and loads no Project Config and no User
Config. A malformed `.roundfixrc.yml` does not stop it.

Phases: `follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery`, `release`.

Options:

- `--spec <slug>` — Spec folder slug matching `^[0-9]{4}(-[a-z0-9]+)+$`,
  for example `0412-example`. Required for `idea`, `prd`, `techspec`, `tasks`
  and `delivery`; optional for `follow-up` and `release`.
- `--model <identifier>` — exact model identifier the harness states for the
  session, at most 200 characters matching
  `^[A-Za-z0-9][A-Za-z0-9._:/@\[\]-]{0,199}$`. Omit it when none is stated;
  never infer it from configuration or another environment variable.
- `--effort <level>` — reasoning effort when `CLAUDE_EFFORT` is unset,
  matching `^[a-z0-9][a-z0-9_-]{0,31}$`.
- `--author-subagent` — the phase's Author Subagent runs it.

The last three flags apply only to `idea`, `prd`, `techspec` and `tasks`.

Standard output is one line:

```text
Recorded <phase> of <spec> for claude session <session>: model <m>, effort <e>, run by <the Author Subagent|the session>
Recorded <phase> of <spec> without a harness session: the environment names none
```

For an Authoring Phase, a learned `<m>` or `<e>` prints as
`<value> (<source>)`; an unknown value prints as `unknown`. Other phases
keep `Recorded <phase> of <spec> for claude session <session>`.
Only one result line is printed. ` of <spec>` is omitted when no Spec was
named. The harness and session come from `CLAUDE_CODE_SESSION_ID` in the
command environment; when it is absent, both fields are empty strings.

The command reads no earlier entry. Whatever the session recorded before, it
writes the entry's `boundary` as `false`. A phase entry ends no session.

Only an Authoring Phase appends these six keys after `boundary`, in order:

| Key | Value |
| --- | --- |
| `repository_id` | Repository Identity digest; empty when local Git cannot compute it |
| `model` | `--model`, or empty when unknown |
| `model_source` | `reported` or `unknown` |
| `effort` | `CLAUDE_EFFORT`, else `--effort`, or empty when unknown |
| `effort_source` | `environment`, `reported` or `unknown` |
| `runner` | `author-subagent`, `session` or `unknown` |

The model accepts `^[A-Za-z0-9][A-Za-z0-9._:/@\[\]-]{0,199}$`, the
Authoring Model Declaration's model-name pattern, at most 200 characters.
The effort accepts `^[a-z0-9][a-z0-9_-]{0,31}$`; an invalid environment
value stays unknown even if a valid flag was given. No configuration or
model environment variable supplies the model. Without a session, the
model and effort are empty, their sources and the runner are `unknown`, and
standard output keeps the line without a harness session. An old entry
without these keys reads with empty values and the three sources `unknown`.
An entry of `follow-up`, `delivery` or `release` keeps its nine keys.

The Repository Identity is computed from the root commits of local Git
history, with or without a session, never a remote address or its
credentials. Full clones and linked working trees share it; a shallow clone
has a different identity because its history is truncated, so an entry
written there does not match a full clone in the Corrections Report.

Malformed values or flags on another phase produce these one-line warnings
on standard error; the entry is still written and the exit code is unchanged:

```text
roundfix: warning: phase start: --model is not a model identifier of at most 200 characters; the model is recorded as unknown
roundfix: warning: phase start: <--effort|CLAUDE_EFFORT> is not an effort level; the effort is recorded as unknown
roundfix: warning: phase start: --model, --effort and --author-subagent apply to idea, prd, techspec and tasks; ignored for <phase>
```

| Exit | Meaning |
| --- | --- |
| `0` | One phase entry recorded |
| `1` | The record could not be written |
| `2` | Unknown phase, missing or malformed `--spec`, extra argument, missing or unknown subcommand, or working directory outside a Git repository; nothing is written |

A refusal prints no standard output and writes no record directory. Standard
error names the reason and the help command:

```text
roundfix: phase start failed: <reason>
Run 'roundfix phase --help' for usage.
```

A write failure prints `roundfix: phase start failed: <error>` on standard
error. `roundfix phase --help`, `roundfix phase -h`, and either flag after
`phase start` print usage and exit `0`.
