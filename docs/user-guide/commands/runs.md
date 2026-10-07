### runs, runs list and runs show

```bash
roundfix runs                     # Run Browser at an interactive terminal
roundfix runs list [--state active|terminal|all] [--limit N] [--all]
roundfix runs show <run-id> [--json]
```

Bare `runs` opens the machine-wide read-only Run Browser: every repository's
Runs newest first, Active only by default, `a` toggles active/all, `Enter`
attaches read-only, `q`/`Esc`/`Ctrl-C` quits with exit `0` and no side
effects. In a non-interactive context it exits `2` naming `runs list`.

`runs list` prints one Run per line, newest first:
`<run-id>  <state>  <kind>  <target>  <agent>  <started-utc>  <duration>
<branch>`. Targets are `pr:<number>` or `spec:<slug>`. Default scope is the
current repository's 20 newest Active Runs; `--state` widens the state filter,
`--limit N` changes the bound (`0` unbounded), `--all` lists every repository
and adds a repository column. When Runs are hidden, exactly one trailing
stderr note names the hidden count and the widening flag. Empty results print
`No Runs found.` and exit `0`.

When the selected repository scope contains retained terminal Run Worktrees or
Run Branches, Runs List keeps the stdout row shape unchanged and writes one
diagnostic to stderr:

```text
(2 terminal Run Worktrees retained; run 'roundfix reconcile' to inspect)
```

The count can appear even when the default Active view has no rows. It is
discovery guidance, not a `safe` or unsafe classification; use the Reconcile
Command to classify each retained Run. When present, this retained-worktree
diagnostic is the one trailing stderr note.


#### Show a Run's token usage

`roundfix runs show <run-id>` reads the Run Database without creating a Run,
starting an Agent or changing recorded state. It prints a header, one
tab-separated row per Work Item scope, then a Run total:

```text
Run <run-id>: implement <spec-slug> Clean
task_01	codex/gpt-6.1-sol/high	5639755 tokens (request-sum) from 1 of 1 prompt(s)	cost not reported
qa	codex/gpt-6.1-sol/high	none reported by 1 prompt(s)	cost not reported
Tokens: 5639755 from 1 of 2 prompt(s); cost not reported
```

A scope row contains its ID, Agent Selections (`runtime/model/effort`, joined
by `, `), token phrase with the reported bases, and cost phrase. Unreported
prompts stay absent. A Run without usage rows prints `Tokens: no prompts
recorded; cost not reported`. See [Token usage](../usage.md#token-usage) for counting rules.

`--json` emits one document with schema `roundfix/runs-show/v1`:

| Field | Value |
| --- | --- |
| `schema` | `roundfix/runs-show/v1` |
| `run_id`, `kind`, `spec_slug`, `state` | Recorded Run identity and outcome |
| `scopes` | Array of scope usage objects |
| `total` | Run usage totals |

Each scope has `scope_kind`, `scope_id`, `selections` (objects containing
`runtime`, `model`, `reasoning_effort`) and the totals below. `total` has the
same totals without those three scope fields:

| Totals field | Value |
| --- | --- |
| `prompts`, `reported_prompts` | Recorded and reported prompt counts |
| `tokens` | Sum of reported counts; null when none reported |
| `bases` | Distinct reported bases |
| `input_tokens`, `output_tokens`, `cached_read_tokens`, `cached_write_tokens`, `thought_tokens` | Sum when every reported prompt carried that split; otherwise null |
| `costs` | Array of `currency` and `amount` objects; empty when absent |
| `sessions`, `cost_sessions` | Recorded Agent Sessions and those reporting cost |

Exit `0` means the Run exists. Exit `2` means usage is invalid (including a
missing or extra Run ID), the Run is unknown, or the database cannot be read.
An unknown Run prints only this diagnostic on stderr:

```text
roundfix: runs show failed: Run "run_missing" does not exist
Run 'roundfix runs show --help' for usage.
```

#### Why Verification failed

```bash
roundfix runs causes [--since <YYYY-MM-DD>] [--until <YYYY-MM-DD>] [--format <text|json>]
```

This read-only command reports failed Verification attempts and corrective
Tasks for this repository's terminal Spec Runs, oldest first. The window
uses Run creation dates: `--since` includes its UTC midnight and `--until`
excludes its UTC midnight. Omit either flag for an open end. The command
writes nothing and makes no network request. SQLite may create the WAL and
shared-memory sidecars when opening a WAL database read-only; the database,
lock and artifact logs stay unchanged.

Each item names its class, the first matching signature, and its check. The
five classes are `scope-or-authorization`, `shared-section-contract`,
`repository-convention`, `implementation-defect` and `environment`. The first
three count as repository knowledge. An unmatched item is `unclassified`,
which is counted separately and is not a class. The summary counts every
class and gives the embedded signature table's digest.

A corrective Task is a graph node numbered after the QA Task, excluding the
QA Task itself. The command reads the active Task Graph first, then its
archive. A missing graph produces `corrective: null` in JSON and lists the
Spec under `specs_not_found`. Corrective Tasks appear once, attached to their
first Run in the window. Their title and Overview determine the trigger:
`pre-pr-review`, `qa-gate`, `verification`, or `unknown`, in that order.

An Archive Record supplies the graph and Task files from Git at its
`source_revision`. If that commit is absent, the Spec stays archived in the
report with its disposition; it is not listed under `specs_not_found`.
JSON adds `archived_specs` entries with `spec`, `disposition` and
`source_available`. Text reports the unavailable revision. Corrective
classification remains unknown without the graph.

`--format json` emits schema `roundfix/runs-causes/v1`, including the window,
signature digest, items, per-Task passed and failed verdict counts, feedback
rounds, Run counts, QA and corrective flags, summary, and missing Specs. An
unclassified item's signature and evidence are null; a corrective item's
attempt is null. Diagnostic text, absolute paths and keys are excluded.

Exit `0` means the report ran, including an empty window or absent Run
Database. Exit `2` means invalid flags, dates, format, a non-increasing window,
or no Git repository. Exit `1` means the Run Database could not be read.
