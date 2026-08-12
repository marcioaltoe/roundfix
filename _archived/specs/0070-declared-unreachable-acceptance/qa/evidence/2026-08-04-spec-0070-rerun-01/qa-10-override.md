# QA-10 — QA override compatibility

Status: pass

The public Archive Command remained the negative control: a failed QA report
with `qa_override: true` still exited 2 with the established fail-verdict
diagnostic.

The explicit `archive-spec` operational path was then followed in a disposable
Git repository. It verified the completed Task, observed the failed QA report,
retained the explicit override, stamped `status: archived` and
`archived: 2026-08-04`, and moved the Spec to `_archived` with `git mv`.
A fresh read confirmed:

```yaml
status: archived
archived: 2026-08-04
qa_override: true
```

No commit or push was performed; scratch staged state contained only the four
fixture files under the archived Spec.
