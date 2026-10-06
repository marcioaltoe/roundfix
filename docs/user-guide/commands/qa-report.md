### qa-report accept

```bash
roundfix qa-report accept <path>
```

Reads the selected QA Report and exits zero only when the shared archive and
settlement eligibility decision accepts it. The pre-PR Pull Request row,
recorded as `blocked (environment: no open Pull Request)` with the Pull Request
row named in its provenance, never decides a qualifying partial and needs no
Unreachable Acceptance declaration. An outside-evidence row blocked only because
the Run sandbox denied network access, recorded as `blocked (environment: network denied: <host>)` with the outside-evidence row named in its provenance,
also never decides a qualifying partial and needs no Unreachable Acceptance
declaration. Any other blocked outside-evidence row still blocks Pull Request
preparation. A `pending` verdict is never
accepted, and a `pass` or otherwise-eligible `partial` that records no QA row
is refused. A report whose front matter is empty or duplicated is unreadable
and refused. The command writes no files.
