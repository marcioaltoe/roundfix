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


#### phase model

```bash
roundfix phase model <phase> --spec <slug>
```

Prints the model and reasoning effort a phase's Author Subagent starts on.
Takes `idea`, `prd`, `techspec` or `tasks`, requires the same Spec slug
pattern as `phase start`, resolves the Git root and loads Project Config.
It writes nothing.

If the phase declares a model, it reads fallback entries oldest first. An
entry applies when it belongs to this Spec and repository, its limit name
occurs in the declared model ignoring case, this phase declares a fallback,
and that fallback's model does not contain the limit name. Repository
Identity must match; when either identity is empty, the Git root must match.
The oldest applicable entry supplies `<phase0>` and `<limit>`; the current
phase's declaration supplies the fallback model and effort.

One of these lines is printed on standard output:

```text
Start roundfix-<phase>-author with model <m> and effort <e>, the fallback of <declared>: <phase0> of <slug> fell back when the <limit> limit ran out.
Start roundfix-<phase>-author with no model and no effort: the declared model <declared>.
No Author Subagent model for <phase>: the Authoring Model Declaration declares no model for it.
```

Without a fallback effort, `and effort <e>` reads `and no effort`. When the
line names the declared model or no Author Subagent model, start the subagent
passing no model and no effort. Another Spec starts from its own record.

| Exit | Meaning |
| --- | --- |
| `0` | A fallback, declared model or absent declaration printed |
| `1` | The record directory could not be read |
| `2` | Usage, Git working directory or configuration error |

A usage error prints these lines on standard error:

```text
roundfix: phase model failed: <reason>
Run 'roundfix phase --help' for usage.
```

A record read failure prints only `roundfix: phase model failed: <error>`
on standard error. Malformed record lines are skipped.

#### phase fallback

```bash
roundfix phase fallback <phase> --spec <slug> --error <message>
```

Decides and records an Authoring Fallback from the error an Author Subagent
ended on. Takes the same phases and Spec slug as `phase model` and loads
Project Config. `--error` must be non-empty and at most 4,096 bytes.

The first matching error line names a limit: `You've hit your <name> limit`
(with either apostrophe), or a line beginning `<name> limit reached`, in any
letter case. A one-word name other than `session` or `weekly` is a model
limit; those two names are shared limits, and multi-word names are other
limits. The kept line has control characters replaced by spaces, is trimmed,
and is cut to 200 characters on a rune boundary. The error is untrusted data.

Checks run in this order. The first refusal prints
`No fallback for <phase> of <slug>: <reason>.` on standard output and writes
nothing:

| Condition | Reason |
| --- | --- |
| No matching limit line | `the error is not a usage limit` |
| Shared limit | `a <name> limit is shared by every model` |
| Other limit | `a <name> limit is not a model's usage limit` |
| No declared model | `<phase> declares no model` |
| Limit name absent from declared model, ignoring case | `the <Name> limit is not that of the declared model <model>` |
| No declared fallback | `<phase> declares no fallback` |
| Limit name also occurs in fallback model | `the fallback <model> is a <Name> model too` |

Otherwise it appends exactly one fallback entry before printing:

```text
Start roundfix-<phase>-author with model <m> and effort <e>, the fallback of <declared>: "<Line>"
```

Without a fallback effort, `and effort <e>` reads `and no effort`. Tell the
user which phase switched, from which model to which, and the limit line.
Without waiting, start the same Author Subagent with the printed model and
effort, the same brief, its returned output and the files it wrote; say in
the brief that it runs on the fallback. If no fallback applies or the restart
cannot start, report the error and author in the session only if the user
says so. The entry remains even if the restart fails.

| Exit | Meaning |
| --- | --- |
| `0` | One fallback entry appended and its start line printed |
| `1` | The record could not be written |
| `2` | Usage, Git working directory or configuration error; nothing written |
| `3` | No fallback applies; nothing written |

A usage error prints these lines on standard error:

```text
roundfix: phase fallback failed: <reason>
Run 'roundfix phase --help' for usage.
```

A write failure prints only `roundfix: phase fallback failed: <error>` on
standard error. Both commands support `--help`. If an older Roundfix does not
know them, start passing no model and no effort and report an error as before.
`roundfix phase start` retains its arguments, output, record and exit codes.
