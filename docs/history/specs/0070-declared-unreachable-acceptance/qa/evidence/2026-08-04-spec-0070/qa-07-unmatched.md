# QA-07 — Unmatched declaration refusal

Status: pass

The QA report claimed two declared-blocked rows while the PRD declared one.
`roundfix archive qa-case` exited 2 with:

```text
no passing QA verdict: rows_blocked_declared is 2, but Spec declares 1 unreachable acceptance; shortfall is 1
```

A fresh Spec Root listing retained the active Spec and contained no archived
destination.
