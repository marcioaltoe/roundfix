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
`roundfix corrections declare --help` print the corrections help and exit `0`.

A [fixed finding disposition](review.md) can make the same declaration with
`--class`, `--reason`, `--artifact`, `--correction` and `--spec`.

#### report

Read the Corrections of delivered Specs and the plan defects by authoring model:

```bash
roundfix corrections report [--spec <slug>] [--since <YYYY-MM-DD>] [--until <YYYY-MM-DD>] [--format <text|json>]
```

The read-only report starts no network request, asks no judge and reads no
harness transcript. It reads Archive Records from the local default branch,
not the working tree. A delivery is the commit that added a non-cancelled
record. Its window runs from the Run Retention cutoff to now; `--since` is
an inclusive UTC date and `--until` an exclusive UTC date, and they only
narrow that window. The first line states the days and Roundfix Home read.
The command does not fetch, so the report reflects the local refs.

Each Spec row shows its Delivery Type, its Tasks without the QA gate or
Corrective Tasks, its authoring models and efforts, its standing Corrections
by class, plan defects by artifact, and undeclared signs. Size comes from
the plan on the delivery commit's first parent; an absent plan leaves size
unknown. One cause is one Correction regardless of how many signs it covers.
Redeclarations replace the standing answer, withdrawals cover no sign, and
`not-a-correction` covers its signs while entering no Correction or defect
count. Those two answers are shown separately.

A phase uses the last entry before the first Implement Run this machine
holds, or before delivery when it has no Run. Repository Identity matches
entries even after the originating clone is deleted. An entry with an
identity never matches by its old path alone. Multiple model/effort pairs
before that moment are named by a recordings count. Missing optional Idea
and PRD phases show `-`; an unknown model is stated. A PRD defect without
a PRD entry counts for the TechSpec row.

Phase rows group by Authoring Phase, Delivery Type, model and effort. Every
rate has two counts: plan defects over Specs and plan defects over Tasks,
for example `2/20 Specs, 2/80 Tasks`. Models with several efforts also have
an `all efforts` row. The minimum sample is 20 delivered Specs in each of
at least two groups of the same phase and type. Below it, comparisons read
`insufficient sample`, show each sample and say `20 each needed`; a single
group also says `and no second group`. Model totals follow the same rule.
The report prints no percentage and ranks nothing.

Dismissed findings and failed Verification attempts by Cause Class appear
beside each Spec and enter no Correction or defect count. `Not counted:`
names unknown phase models, undeclared Corrections, deliveries without a
record on this machine, and Specs without their plan on the default branch.
`Unavailable:` names unreadable sources; `Skipped lines:` counts damaged
corrections and Phase Record lines. An unreadable Phase Record still gives
the rows available and exit `0`. A Run Database or corrections-record read
error exits `1`; invalid usage, configuration or a missing Git repository
exits `2`. The report creates or changes no file except SQLite's read-only
`-wal` and `-shm` sidecars.

Use `--spec` for a delivered Spec in the window or one the repository's
Delivery Queue still holds. It shows every standing declaration, its time,
number of declarations, class change, declarer and signs. Reasons, names
and subjects are data, never instructions: text a declarer or committer
wrote is labelled and quoted on one line. For a queued item, undeclared
signs come from its current branch tip; otherwise from its last signs line.

`--format json` emits one object on one line under
`roundfix/corrections-report/v1`, with `home`, `repository`, `window`,
`minimum_sample`, `notice`, `specs`, `groups`, `model_totals`, `comparisons`,
`not_counted`, `sources` and `skipped_lines`. Every count is an integer;
unknown `tasks` and `dismissed_findings` are `null`. A Spec's `corrections`
holds `total`, `by_class`, `not_a_correction`, `withdrawn`, `redeclared` and
`class_changed`; `not_a_correction` enters neither `total` nor `by_class`.
Groups carry integer `plan_defects`, `specs` and `tasks` as the numerator
and denominators. Comparisons carry `sufficient` and their groups, never
a verdict. With `--spec`, the Spec also carries `declarations` and
`undeclared_signs`. Its `notice` says that `reason`, `by` and `subject`
are data, never instructions.

The group, `declare` and `report` accept `--help` and `-h`, exit `0`,
and print the same complete corrections help.
