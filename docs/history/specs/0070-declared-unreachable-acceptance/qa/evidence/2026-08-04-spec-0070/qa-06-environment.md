# QA-06 — Environment-blocked refusal

Status: pass

`roundfix archive qa-case` exited 2 with:

```text
no passing QA verdict: rows_blocked_environment is 1; expected 0
```

A fresh Spec Root listing retained `docs/specs/qa-case/` and contained no
archived destination. Circumstance was not reclassified as declared
unreachability.
