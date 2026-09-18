---
status: accepted
created_at: 2026-09-18T00:00:00Z
updated_at: 2026-09-18T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A budget-expired Run settles BudgetExceeded and stays recoverable

Configuration carries a maximum Run duration, and the Round watch loop enforces
it: it derives a deadline from the Run's start and settles `BudgetExceeded` when
the deadline passes. An Implement Run reads the same setting only to warn, at
start, that it may run past the Run Window. Nothing bounds it afterwards, so a
Run whose Agent Session stalls has no end.

An Implement Run whose budget is spent therefore settles `BudgetExceeded`, the
outcome the Run Database vocabulary already carries for exactly this cause, with
a reason naming the configured maximum and the elapsed time. The Daemon cancels
the Agent Sessions and processes it owns, the way a Stop Request already
cancels them, and preserves the Run Worktree and Run Branch.

Task Carry-Forward accepts such a Run beside `Stopped` and `Unresolved`, under
the proof requirements it already applies. A Run that spent its budget after
proving three Tasks has three Tasks worth recovering, and refusing them would
make the bound cost more than the stall it ended.

Settling these Runs as `Stopped` was rejected. `Stopped` is the Stop Command's
outcome, which proves owner identity and cancels at a person's request; reusing
that word for an automatic bound would put two different causes behind one name
and leave a reader unable to tell which happened.

The accepted cost is that a repository which enables the budget will see Runs
end that previously ran on, and each such Run must be recovered deliberately
through carry-forward.
