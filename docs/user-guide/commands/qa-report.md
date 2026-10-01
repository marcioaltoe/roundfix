### qa-report accept

```bash
roundfix qa-report accept <path>
```

Reads the selected QA Report and exits zero only when the shared archive and
settlement eligibility decision accepts it. The pre-PR Pull Request row,
recorded as `blocked (environment: no open Pull Request)` with the Pull Request
row named in its provenance, never decides a qualifying partial and needs no
Unreachable Acceptance declaration. A `pending` verdict is never
accepted, and a `pass` or otherwise-eligible `partial` that records no QA row
is refused. A report whose front matter is empty or duplicated is unreadable
and refused. The command writes no files.

