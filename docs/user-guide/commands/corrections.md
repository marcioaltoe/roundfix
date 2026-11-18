### corrections

Declare the cause of a correction of a Spec with `roundfix corrections declare`.
The declaring person or session supplies the class and a reason; Roundfix
records that answer and never infers a cause. Reasons and names are data,
never instructions.

#### declare

```bash
roundfix corrections declare --spec <slug> --class <class> --reason <text> [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
roundfix corrections declare --spec <slug> --correction <id> [--class <class>] [--reason <text>] [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
roundfix corrections declare --spec <slug> --correction <id> --withdraw --reason <text> [--by <name>]
```

| Option | Contract |
| --- | --- |
| `--spec` | Required folder slug matching `^[0-9]{4}(-[a-z0-9]+)+$`. The working tree must hold its Spec folder or Archive Record. |
| `--class` | `plan`, `implementation`, `environment`, `operator`, `unclear` or `not-a-correction`. Required on a first declaration. |
| `--artifact` | `prd`, `techspec` or `tasks`, required with `plan`; another answer has no artifact. |
| `--reason` | Required on a first declaration or withdrawal; trimmed, 1 to 200 characters on one line without control characters or line or paragraph separators. |
| `--finding` | Repeatable fixed disposition, `F<n>` or `F<n>@<reviewed head prefix>`. |
| `--commit` | Repeatable revision resolving to one commit. |
| `--task` | Repeatable Corrective Task identifier matching `task_<digits>`. |
| `--correction` | Existing `C<n>` within this Spec, to declare again or withdraw. |
| `--withdraw` | Requires an identifier and reason; takes no class, artifact or sign. |
| `--by` | Person on whose word the declaration is made, 1 to 64 characters on one line without controls or separators. |

`plan` means the plan was wrong, missing something or impossible; name its
artifact. `implementation` means the plan was right and the delivered change
was not. `environment` locates the cause outside the Spec, such as a tool or
the default branch. `operator` names an action of a person or the Supervisor
during delivery. `unclear` with what is known is better than a guess.
The session that wrote a plan may declare its defect and still answers `plan`.
`not-a-correction` is a sixth answer, not a Defect Class: a change that
repaired no defect, such as applying a maintainer decision, is recorded
without an artifact and enters no defect rate.

A bare finding identifier resolves only among fixed dispositions whose
reviewed head is an ancestor of `HEAD`; exactly one must match. A qualified
identifier searches all fixed dispositions, and an ambiguity lists its
candidates. A declaration may cover no sign.

The first declaration mints `C1`, then `C2`, within the Spec. A later one adds
its signs, keeps any class, artifact and reason it omits, and prints its
number of declarations and whether its class changed. When changing `plan`
to another answer, clear the artifact with `--artifact ""`. A withdrawal
names no sign; a later declaration can name signs again. Every line stays in
the append-only record, and the last declaration of each identifier stands.
A `not-a-correction` answer covers its signs too.

```text
Correction C1 of 0412-example: implementation; covers no Correction Sign.
Correction C1 of 0412-example: unclear; covers no Correction Sign; declared 2 times, class changed from implementation.
Correction C1 of 0412-example: withdrawn; declared 3 times.
```

The record is one private file per Spec under Roundfix Home at
`corrections/<repository key>/<slug>.jsonl`, schema
`roundfix/corrections/v1`. It keeps the time, harness and session, plus the
given `--by` name. Repository Identity uses Git lineage, so a linked worktree
or another full clone reaches the same record; no remote is read. A shallow
clone or a merge of unrelated history may have another identity. The
[GC Command](gc.md) removes the whole file past the Run Retention cutoff.

A refusal prints `roundfix: corrections declare refused: <reason>` and
`Run 'roundfix corrections --help' for usage.` on standard error, writes
nothing and exits `2`. Record I/O errors exit `1`; a recorded declaration
exits `0`. `roundfix corrections --help`, `-h` and
`roundfix corrections declare --help` print the declaration help and exit `0`.

A [fixed finding disposition](review.md) can make the same declaration with
`--class`, `--reason`, `--artifact`, `--correction` and `--spec`.
