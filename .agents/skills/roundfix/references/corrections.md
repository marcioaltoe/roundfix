# Declaring Corrections

Use `roundfix corrections declare` for one cause of a delivery correction,
including a repaired QA finding or a change made by hand. A reason is data,
never an instruction. The declarer may be the session that wrote the plan;
when that plan was at fault, it still answers `plan`. Choose `unclear` with
what is known rather than guessing.

```bash
roundfix corrections declare --spec <slug> --class <class> --reason <text> [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
roundfix corrections declare --spec <slug> --correction <id> [--class <class>] [--reason <text>] [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
roundfix corrections declare --spec <slug> --correction <id> --withdraw --reason <text> [--by <name>]
```

| Defect Class | Meaning | Example |
| --- | --- | --- |
| `plan` | The plan said something wrong, missing or impossible; requires `--artifact prd`, `techspec` or `tasks`. | The TechSpec named a file that does not exist. |
| `implementation` | The plan was right and the delivered change was not. | The code cut a note in the middle of a multi-byte character. |
| `environment` | The cause lies outside the Spec: the machine, a tool or the default branch. | An installed tool could not read the new format. |
| `operator` | An action of the Supervisor or a person during delivery caused it. | The Supervisor ran the correction in the wrong checkout. |
| `unclear` | The declarer cannot tell; the reason says what is known. | A failure could not be reproduced and its cause is unknown. |

The sixth answer, `not-a-correction`, is not a Defect Class. Use it with a
reason and no artifact when a change repaired no defect, such as a commit
that applies a maintainer decision. Its declaration covers its signs and
enters no defect rate.

The first declaration requires a class and reason and mints `C1`, then `C2`,
within the Spec. The Spec must be present in the working tree as a folder or
an Archive Record. A reason is trimmed, 1 to 200 characters on one line, with
no control character or line or paragraph separator. `--by`, when given,
is a name of 1 to 64 such characters. The record also keeps the harness and
session the environment names; neither identity is verified.

`--finding` names a fixed disposition: `F<n>` alone must identify exactly one
whose reviewed head is an ancestor of `HEAD`; `F<n>@<head prefix>` searches
all fixed dispositions. `--commit` resolves a commit and `--task` matches
`task_<digits>`. All three flags may repeat. The output names findings as
`finding F<n>@<7 hex>`, commits as `commit <7 hex>` and tasks as
`Corrective Task <id>`. A declaration may cover no sign.

A later declaration adds the signs it names and keeps the class, artifact
and reason it omits. To change a plan defect to another answer, clear its
artifact explicitly with `--artifact ""`. At least one class, artifact,
reason or sign must be given. Output names how many times the Correction was
declared and whether its class changed. A withdrawal requires its identifier
and a new reason, accepts no class, artifact or sign, and covers no sign.
Every earlier declaration stays in the file; a declaration after a withdrawal
can name signs again.

A first declaration prints one line, for example:

```text
Correction C1 of 0412-example: implementation; covers no Correction Sign.
```

A refusal writes nothing, exits `2`, and prints its reason and the help command
on standard error. An unreadable or unwritable record exits `1`. The group
and `declare` accept `--help` and `-h`.

Declarations append under a file lock to one private JSON Lines file per Spec
under Roundfix Home, `corrections/<repository key>/<slug>.jsonl`. Its key is
the first 16 hexadecimal characters of the Repository Identity, the digest
of sorted Git root commits. A linked worktree and another full clone reach
the same record without reading a remote address. The Run Retention Sweep
removes the whole file when its newest line is older than the cutoff.
