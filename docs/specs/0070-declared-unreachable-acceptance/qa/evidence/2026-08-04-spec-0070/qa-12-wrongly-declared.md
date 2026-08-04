# QA-12 — Wrongly declared reachable row

Status: pass

The disposable PRD declared `Archive Command help is reachable` unreachable.
The gate exercised the journey anyway:

```text
roundfix archive --help
```

It exited 0 and described the declared-only partial acceptance. The fixture's
QA state therefore recorded a wrongly-declared-row finding and `verdict: fail`
rather than a declared block. `roundfix archive qa-case` then exited 2 with
the fail-verdict diagnostic, and a fresh read retained the active Spec.
