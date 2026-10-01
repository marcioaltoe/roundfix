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
