### migrate

```bash
roundfix migrate
```

Upgrades an existing older Run Database to this binary's schema version under
the machine-wide write lock. It needs no Git repository and reports exactly one
of four outcomes:

- An older database prints `Run Database migrated from schema version <from>
  to <to>: <path>`.
- A current database is not written and prints `Run Database is already at
  schema version <n>: <path>`.
- An absent database creates neither the database nor its directory and prints
  `No Run Database at <path>; nothing to migrate`.
- A newer database is refused without writing and names the newer binary or
  `roundfix upgrade` as the remedy.

A Run started by an older binary that is still Active meets the same migrated
database any operational command would leave.

