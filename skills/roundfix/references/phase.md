## Phase Record

```bash
roundfix phase start <phase> [--spec <slug>] [--model <identifier>] [--effort <level>] [--author-subagent]
```

Run `roundfix phase start` when a session starts a Supervisor Phase: `follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery` or `release`. `--spec` takes the Spec folder slug and is required for every phase but `follow-up` and `release`. The command appends one phase entry to the Phase Record and prints one line. It loads no Project Config and no User Config. It exits `2` and writes nothing for an unknown phase, a missing or malformed `--spec`, an extra argument or a directory outside a Git repository, and `1` when the record cannot be written. When the installed Roundfix does not know the command, say so and continue.

The Phase Record is the day files in the directory `phase-record` under Roundfix Home, one JSON object per line with the schema `roundfix/phase-record/v1`. A phase entry names the phase, the Spec, the repository, the harness and its session, read from the environment; without a session the entry is written with none. No command lists the record: read the files.

For a call without invalid or inapplicable values, no warning is needed.
The command never reads the record back: whatever the session recorded before, it writes `boundary` as `false` and prints nothing on standard error. An entry an earlier Roundfix marked keeps its mark and stays readable. A phase entry ends no session: the session goes on with the next phase and ends only when the user asks or the work stops.

### Authoring values

On `idea`, `prd`, `techspec` and `tasks`, add `--model <identifier>` with
only the exact model identifier the harness states for this session, never
an alias or a guess; omit it when the session's instructions state none.
Add `--author-subagent` when the phase's Author Subagent runs the phase.
`--effort <level>` supplies the reasoning effort only when `CLAUDE_EFFORT`
is unset. If the installed Roundfix refuses the flags, run the command
without them and say so in the report.

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

With a session, the authoring result line is:

```text
Recorded <phase> of <slug> for claude session <session>: model <value> (<source>), effort <value> (<source>), run by <the Author Subagent|the session>
```

An unknown value prints as `unknown` without a parenthesized source.
Malformed values or flags on another phase produce these one-line warnings
on standard error; the entry is still written and the exit code is unchanged:

```text
roundfix: warning: phase start: --model is not a model identifier of at most 200 characters; the model is recorded as unknown
roundfix: warning: phase start: <--effort|CLAUDE_EFFORT> is not an effort level; the effort is recorded as unknown
roundfix: warning: phase start: --model, --effort and --author-subagent apply to idea, prd, techspec and tasks; ignored for <phase>
```

The Daemon appends a prompt entry when a prompt of an Agent Session ends, and `roundfix review` after each prompt of a pre-PR review: the Run, the Spec, the scope (`task`, `qa`, `review`, `session` or `pre-pr-review`), the repository, the working directory, the runtime, the session the adapter reported and the times the prompt began and ended. A prompt that never returns leaves no entry. When the record cannot be written, one warning is printed and the Run or the review ends as it would have.
