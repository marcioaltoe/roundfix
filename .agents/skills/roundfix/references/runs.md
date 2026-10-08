## Detached Runs

Use `--detach` on `resolve`, `watch`, or `implement` when the caller must not
own the Run lifetime, such as scripts and CI jobs. The foreground command
starts a Detached Run, prints exactly this five-line stdout report, and exits
`0`:

```text
Run ID: <run-id>
Console Log: <path>
Attach: roundfix attach <run-id>
Supervisor monitor: roundfix events <run-id> --follow --filter outcome
Stop: roundfix stop <run-id>
```

The console log path is under the Artifact Directory at
`<artifact_dir>/runs/<run-id>/console.log`; it receives the detached child's
stderr, Agent output, and terminal outcome messages. `Attach` is the human
Live Run View; `Supervisor monitor` is the stable terminal outcome
subscription; `Stop` is the Stop Command surface. Detached Runs behave as
normal non-TTY Runs after startup: Run Events, Worktrees, integration,
outcomes, and locks keep their normal contracts. The detached child that wins
terminal completion sends the configured outcome notification. A competing
completion observes the stored outcome and does not send another notification.
Supervisors and scripts use the printed
`roundfix events <run-id> --follow --filter outcome` command, humans use the
Live Run View with `roundfix attach <run-id>`, and neither treats the console
log as a state API.

Detach implies non-interactive mode. `--interactive` is rejected before Run
creation, and `--no-input` is implied. Startup uses a two-phase handshake: the
child writes a liveness marker immediately on entering child mode, before
configuration load and Preflight Validation, then writes the Run id after the
Run exists. The foreground command waits 10 seconds only for liveness; Run
creation has its own 5-minute ceiling. A slow but live Preflight Validation no
longer fails detach startup only because it takes more than 10 seconds.

Detached startup failures write no stdout and always print an explicit stderr
diagnostic before the console relay or empty-output note:

```text
roundfix: Detached Run child produced no liveness signal within 10s; killed (exit: <exit or signal>)
roundfix: Detached Run child did not create a Run within 5m0s; killed (exit: <exit or signal>)
roundfix: Detached Run child exited before the handshake (<exit or signal>); console output follows
roundfix: Detached Run child exited before the handshake (<exit or signal>) and produced no output
```

The winning terminal transition for an operational Run through `resolve`,
`watch`, or `implement` sends exactly one outcome notification. Identical
completion replay and conflicting completion do not republish it. `fetch`,
`settle`, `archive`, and commands that create no Run do not notify.
Notification failures write one stderr warning shaped as
`roundfix: outcome notification failed: <reason>` and one Daemon-source Run
Event; they never change the Run report, terminal outcome, or exit code.
Every notification attempt appends exactly one durable receipt Run Event with
its route, completion time, and status `sent`, `skipped`, or `failed`. A sent
receipt means the local route accepted the request, not that a person saw it.

## Run discovery and Attach

Use `roundfix runs list` to discover Runs recorded in the Run Database. By
default it lists this repository's 20 newest Active Runs, newest first. Each
line uses stable plain-text columns:

```text
<run-id>  <state>  <kind>  <target>  <agent>  <started-utc>  <duration>  <branch>
```

Run ids are full and untruncated; start times are absolute UTC (RFC 3339);
durations render like `42m` / `1h12m`, with `running <elapsed>` for Active
Runs; missing fields render `-`. Targets are `pr:<number>` for review Runs
and `spec:<slug>` for Spec Runs. Agents widen the view with `--state
<active|terminal|all>` (default `active`) and `--limit N` (default `20`, `0`
unbounded, applied after the state filter). `--all` lists every repository
and adds the repository path as a final column. The flags compose. When the
state filter or the bound hides Runs, exactly one trailing stderr note names
the hidden count and the widening flag:

```text
(23 terminal Run(s) hidden; use --state all)
(15 older Run(s) hidden; use --limit 0)
```

With no matches, stdout is exactly:

```text
No Runs found.
```

`runs list` exits `0` for matching and empty results. Invalid flags,
unexpected arguments, repository-resolution failures, and Run Database
open/list failures exit `2` with diagnostics on stderr. Outside a Git
repository, `runs list` without `--all` exits `2` and names `--all` as the
alternative.

`runs list`, the Run Browser, and Attach refuse a Run Database whose schema
version differs from the binary's without writing it. Upgrade an older Run
Database with `roundfix migrate`, then retry discovery or Attach. A newer Run
Database needs the newer Roundfix binary that wrote it; use that binary or run
`roundfix upgrade`.

At an interactive terminal, bare `roundfix runs` and `roundfix attach`
without a Run ID open the Run Browser: machine-wide, every repository's Runs
newest first, Active Runs only by default, with a header naming the
`ACTIVE`/`ALL` state filter and rows showing short run id, state, kind,
target, Agent, relative start, duration, branch, and repository. No git
repository is required. `↑↓` moves, `Enter` attaches the selected Run
through the read-only Live Run View — leaving it returns to a refreshed
browser — `a` toggles active/all, and `q`/`Esc`/`Ctrl-C` quits with exit `0`
and no side effects. The empty Active view really does mean nothing is
running anywhere: `No active Runs — press a to include terminal Runs.` In a
non-interactive context, bare `runs` exits `2` and names
`roundfix runs list`; `attach` without a Run ID, including `--no-input`,
exits `2` and names `roundfix runs list` as the discovery command. The Run
Browser is the human surface — agents use the bounded `runs list`, which
stays repository-scoped with `--all` for every repository.

Use `roundfix attach <run-id>` to replay a Run's Run Event Journal and follow
new Run Events read-only. Attach never creates Runs, fetches, starts Agents,
commits, pushes, stops, or resolves Review Source threads. An unknown Run ID
exits `2` with an error stating that picker numbers are not stable Run ids —
pass a run id or run `roundfix attach` to pick interactively.

## Live Run View

The Live Run View uses the same cockpit for review and spec Runs, whether the
Run is owned by `resolve`, `watch`, or `implement`, or replayed read-only
through Attach. The cockpit reads the Run Event Journal; Attach replays that
Journal and then follows new Run Events without mutating or stopping the Run.

- The `WORK QUEUE` pane lists Work Items on the left: Review Issues for review
  Runs and Tasks for spec Runs.
- The `SESSION.TIMELINE` pane is the wider right pane. It groups Run Events by
  Batch and event kind, including Agent plan/tool/think/status events and
  Daemon milestones such as verification, commit, QA, push, and outcome.
- Batch groups collapse automatically by state: a Batch whose state is
  `completed`, `failed`, or `stopped` folds to one `▶` summary row, and every
  other Batch renders expanded under `▼`. Collapse is state-driven; no key
  toggles it.
- Every structured event renders as exactly one bounded summary row behind an
  aligned timestamp gutter. Raw payloads (tool JSON, diffs, markdown bodies)
  never render inline; full content stays in the Detail Modal.
- The timeline pane header carries a `Live · detail hidden` /
  `Live · detail open` indicator that follows the Detail Modal state.
- Empty panes explain themselves per Run kind, naming what would populate
  them — a Fetch Run, for example, reports that it writes Review artifacts to
  disk and starts no Agent.
- State is color-coded in capable terminals: cyan section labels and active
  borders, green done, amber running/waiting/pending, red locked, failed, or
  blocking, and muted gray timestamps and paths. Under `ROUNDFIX_COLOR=never`
  or `NO_COLOR`, the layout and text markers are unchanged, so every state
  distinction survives without color.
- The Phase Row stays above both panes. Review Runs show
  `FETCH > TRIAGE > AGENT > VERIFY > PUSH`; spec Runs show
  `AGENT > VERIFY > COMMIT`, plus `QA` only when the Run opted into QA. Status
  markers are text: `[done]`, `[run]`, `[wait]`, and `[locked]`.
- `Enter` opens the selected Work Item's Detail Modal; `D` toggles it; `Esc`
  closes it. Review detail shows the Review Issue artifact. Spec detail shows
  the Task file body read-only.
- Normal footer keys are `Tab focus`, `↑↓ move/scroll`, `PgUp/PgDn page`,
  `Enter issue` or `Enter Task`, `D show detail`, `End follow`, and the mode
  key. The modal footer keys are `Esc close`, `j/k scroll`, `PgUp/PgDn page`,
  and the mode key.
- Owning active Runs use `Ctrl-C stop`. Attach uses `q detach` in the footer
  and detaches with `q` or `Ctrl-C`; detaching never stops the Run. Owning
  terminal Runs use `q close`.
- Below the two-pane width, the cockpit collapses to `SESSION.TIMELINE` with a
  one-line Work Queue summary and a footer hint to widen the terminal.


### Token usage

`roundfix runs show <run-id> [--json]` reads durable token usage without writing
the Run Database. Text output lists each Work Item scope's Agent selections,
tokens and counting basis, prompt coverage and adapter-reported cost, followed
by a `Tokens:` total. Unreported prompts stay visible and never become zero;
no usage rows says `no prompts recorded`. Cost is grouped by currency, with
Agent Session coverage, and is never priced from tokens. Usage survives Run
Event Journal retention. An unknown Run or invalid arguments exit `2`.

`--json` emits schema `roundfix/runs-show/v1` with `run_id`, `kind`, `spec_slug`,
`state`, `scopes` and `total`. Each scope carries `scope_kind`, `scope_id`,
`selections`, `prompts`, `reported_prompts`, nullable `tokens`, `bases`, nullable
`input_tokens`, `output_tokens`, `cached_read_tokens`, `cached_write_tokens`
and `thought_tokens`, `costs` (currency/amount), `sessions` and `cost_sessions`.
A split is null unless every reported prompt supplied it. The total has the
same usage fields without scope identity and selections.

### Why Verification failed

```bash
roundfix runs causes [--since <YYYY-MM-DD>] [--until <YYYY-MM-DD>] [--format <text|json>]
```

Read-only: explains failed Verification attempts and corrective Tasks in this
repository's terminal Spec Runs, without writes or network requests. Dates
bound Run creation times at UTC midnight, since inclusive and until exclusive;
omitting an end leaves it open. SQLite may create WAL sidecars while opening
the database read-only; database, lock and artifact logs stay unchanged.

The first embedded signature that matches determines one of five classes:
`scope-or-authorization`, `shared-section-contract`, `repository-convention`,
`implementation-defect`, or `environment`. The first three count as repository
knowledge. An `unclassified` item is not a class and carries no signature.

Corrective Tasks are graph nodes numbered after the QA Task. The active graph
is read first, then the archive; a missing graph makes corrective null and
names the Spec in `specs_not_found`. Triggers are `pre-pr-review`, `qa-gate`,
`verification`, then `unknown`. JSON includes per-Task verdicts, feedback and
Run counts, and the table digest. Exit `0` includes empty results and an absent
database, `2` means invalid usage or no Git repository, and `1` means an
unreadable database.

An unknown Run ID exits `2`. When its `run_<YYYYMMDD>T<HHMMSS>Z_<hex>`
creation time predates the configured Run Retention cutoff, the diagnostic
adds `; Run Retention may have removed it, because it removes terminal Runs
that completed more than <N> days ago`, using `store.run_retention_days` from
User Config (30 by default). Other IDs, including `run_missing`, keep their
existing refusal. The hint describes a possible removal, not proof of one.
