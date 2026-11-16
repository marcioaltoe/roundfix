### phase

```bash
roundfix phase start <phase> [--spec <slug>]
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

Standard output is one line:

```text
Recorded <phase> of <spec> for claude session <session>
Recorded <phase> of <spec> without a harness session: the environment names none
```

Only one result line is printed. ` of <spec>` is omitted when no Spec was
named. The harness and session come from `CLAUDE_CODE_SESSION_ID` in the
command environment; when it is absent, both fields are empty strings.

When the same non-empty session last recorded a different phase or Spec in
the current UTC day or the six preceding dates, the command marks the entry
and prints the Boundary Notice on standard error:

```text
roundfix: notice: this session already recorded <phase> of <spec>; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.
```

The notice omits ` of <spec>` when the earlier entry names no Spec. Three
exceptions print no notice: an earlier `follow-up`; `delivery` followed by
`release` of the same Spec or with no Spec; and the same phase of the same
Spec again. Another session, or an entry without a session, never matches.
The command refuses nothing at a boundary: standard output and the exit code
stay as they are.

| Exit | Meaning |
| --- | --- |
| `0` | One phase entry recorded, including at a boundary |
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
