# QA-09 — Pass compatibility

Status: pass

An ordinary `verdict: pass` fixture with no declaration section archived with
exit 0 and the established stdout:

```text
archived qa-case -> docs/specs/_archived/qa-case
```

The archived PRD gained `status`, `archived`, and `source_slug`; it did not
gain `unproven`. A second invocation exited 2 because the active source no
longer existed, confirming persisted move semantics.
