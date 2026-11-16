## Phase Record

```bash
roundfix phase start <phase> [--spec <slug>]
```

Run `roundfix phase start` when a session starts a Supervisor Phase: `follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery` or `release`. `--spec` takes the Spec folder slug and is required for every phase but `follow-up` and `release`. The command appends one phase entry to the Phase Record and prints one line. It loads no Project Config and no User Config. It exits `2` and writes nothing for an unknown phase, a missing or malformed `--spec`, an extra argument or a directory outside a Git repository, and `1` when the record cannot be written. When the installed Roundfix does not know the command, say so and continue.

The Phase Record is the day files in the directory `phase-record` under Roundfix Home, one JSON object per line with the schema `roundfix/phase-record/v1`. A phase entry names the phase, the Spec, the repository, the harness and its session, read from the environment; without a session the entry is written with none. No command lists the record: read the files.

When the same session last recorded a different phase or Spec in the last seven days, the command prints the Boundary Notice on standard error and marks the entry. The Boundary Notice refuses nothing: standard output and the exit code stay as they are. At a Phase Boundary, write the Session Handoff and continue in a new session.

The Daemon appends a prompt entry when a prompt of an Agent Session ends, and `roundfix review` after each prompt of a pre-PR review: the Run, the Spec, the scope (`task`, `qa`, `review`, `session` or `pre-pr-review`), the repository, the working directory, the runtime, the session the adapter reported and the times the prompt began and ended. A prompt that never returns leaves no entry. When the record cannot be written, one warning is printed and the Run or the review ends as it would have.
