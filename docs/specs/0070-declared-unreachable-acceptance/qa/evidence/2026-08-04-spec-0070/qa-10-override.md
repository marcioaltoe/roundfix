# QA-10 — QA override compatibility

Status: pass

Source tracing confirmed the existing override belongs to the `archive-spec`
operational skill, not the pass-oriented Archive Command. The CLI negative
control exited 2 on a failed report even with an existing `qa_override` marker,
matching the unchanged command boundary.

Task 06 supplies scoped authority to verify the existing override path. In a
disposable Git repository, the maintainer path verified completed Tasks,
observed the failed QA report, retained the explicit `qa_override: true`
record, stamped `status: archived` and `archived: 2026-08-04`, and moved the
Spec with `git mv` to `docs/specs/_archived/qa-case/`. A fresh read confirmed:

```yaml
status: archived
archived: 2026-08-04
qa_override: true
```

No commit or push was performed; the staged scratch state contained only the
five moved Spec files. Self-containment was not bypassed: the fixture had no
index and therefore followed the skill's legacy self-containment case.
