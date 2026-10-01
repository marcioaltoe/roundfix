### Review report shape

One line per Review Issue in Round order — with a ` — reason:
<terminal_reason>` suffix on failed, unresolved, and invalid lines when the
artifact carries one — then two labeled summary lines separating this Run's
counts from the pull request's cumulative counts:

```text
issue 001 resolved — major: handle test issue
This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
```

Before a Review Source fetch completes, counts are not known. The report omits
all zero-valued status summaries and prints only:

```text
Review Issues: unknown — fetch did not complete.
```

Review Skipped uses its own two-line source reason and next-action report; it
does not use either count shape.

