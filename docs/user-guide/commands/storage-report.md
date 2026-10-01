### storage report

```bash
roundfix storage report
```

Measures bytes and row counts by repository, state, table, and Artifact Root.
It reads the machine-wide Run Database and recorded Artifact Roots from
Roundfix Home. The report is read-only, accepts no flags, and needs no Git
repository. It never migrates the database, locks for writes, or changes any
byte.

