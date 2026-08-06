---
status: done
created_at: 2026-08-04
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# 2026-08-04 — A Spec archives with `pass` while a user story was never exercised

status: pending

## What was observed

fluxus Spec 0012 reached `verdict: pass` and archived on that verdict. Its
report:

```
verdict: pass
rows_blocked_environment: 3
rows_blocked_finding: 0
Rows: 113 pass, 3 environment-blocked, 0 fail, 0 finding-blocked
```

One of those three environment-blocked rows is **US-06, a PRD user story**.
The gate could not boot the service — the Run Worktree carries no database,
auth, public URL, mail, or storage configuration — so the story was never
walked from a user seat. The report is honest about it and says so plainly.
The archive precondition is still satisfied, because it reads `verdict: pass`
and `rows_blocked_finding: 0`.

This is not a one-off. The same wall produced `verdict: partial` on Spec 0013
a week earlier, and the same three rows will block on Specs 0014 and 0015,
whose acceptance also observes running HTTP surfaces. The operator has already
had to make a standing judgement call to keep delivery moving.

## Root cause

`environment` and `finding` blocks are correctly separated, and only the
latter gates the verdict. That is the right call for a single Spec — an
unreachable third-party API should not fail work that is complete.

What is missing is that the *cost accumulates silently*. Each Spec is
individually defensible; the fleet-level result is a set of archived Specs
whose most user-visible rows were never executed, and nothing aggregates that
into a visible number. A `pass` reads as "verified" long after anyone
remembers which rows were not.

The proximate cause is mundane and fixable: there is no non-production runtime
profile a Run Worktree can boot with. `worktree.copy` is deliberately empty in
this repository, and correctly so — the checkout's `.env` holds the production
`DATABASE_URL`, so copying it would point migrations and repository tests at
production. The gate is left choosing between no environment and a dangerous
one.

## What would settle it

The narrow fix is a first-class, non-secret runtime profile for gates: a
declared way to supply a Spec's Verification and QA with an ephemeral database
and placeholder credentials that are obviously not production. `worktree.copy`
plus `worktree.bootstrap` almost gets there, but copying a real `.env` is the
wrong shape — the profile should be authored for testing, not borrowed from
the developer's machine.

The broader fix is to stop letting the cost disappear:

- Carry the environment-blocked row count into the archive stamp, so an
  archived Spec records what its gate could not reach rather than only that it
  passed.
- Report the fleet-level total somewhere an operator sees it. Three rows per
  Spec across a dozen Specs is a coverage hole worth a decision; one row at a
  time never triggers one.

Neither changes the verdict rule, which is sound. They keep a `pass` from
quietly meaning less than it did last month — the same property this Spec's
own subject matter exists to defend, since 0012 is about a reconciliation that
reported clean while never having run.

## Spec pointer

None yet.
