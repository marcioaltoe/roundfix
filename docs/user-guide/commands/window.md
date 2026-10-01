### window

```bash
roundfix window set <HH:MM|YYYY-MM-DDTHH:MM> [--force]
roundfix window show
roundfix window clear
```

The `window` command manages one durable Run Window for the current repository
in the Run Database. The Run Window bounds when an `implement` Run may start;
`budget.max_run_duration` bounds how long a Run may run after it starts, and an
Implement Run's allowance renews at each Task settlement. The window does not
apply to `fetch`, `resolve`, or `watch`.

`set` accepts a local `HH:MM` and stores its next occurrence: tomorrow when
that time has already passed today. It also accepts an absolute local
`YYYY-MM-DDTHH:MM` instant. A past absolute instant is refused with exit `2`.
If a window already exists, `set` reports the standing cutoff and leaves it
unchanged; pass `--force` to replace it. `show` prints the repository, cutoff,
current time, and remaining duration, or reports that no window is set. Both
states exit `0`. `clear` reports whether it removed the window.

When `implement` runs after the cutoff, Preflight Validation refuses before it
creates a Run. It exits `2`, leaves stdout empty, and names the cutoff, current
time, and recovery commands:

```text
Preflight failed

Reason:
  the Run Window for /path/to/repo closed at 2026-08-27 07:00 BRT; the time is 2026-08-27 07:14 BRT

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

Next action:
  move the window with `roundfix window set <HH:MM> --force`, or remove it with `roundfix window clear`
```

If the window is open but the remaining time is shorter than
`budget.max_run_duration`, `implement` still creates the Run and writes one
stderr crossing report:

```text
Run Window: closes 2026-08-27 07:00 BRT, in 12m; max_run_duration is 45m, so this Run may run past it.
```

The crossing report does not refuse or shorten the Run. A Run created before
the cutoff continues to its own terminal outcome.

