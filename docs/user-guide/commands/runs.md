### runs and runs list

```bash
roundfix runs                     # Run Browser at an interactive terminal
roundfix runs list [--state active|terminal|all] [--limit N] [--all]
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

