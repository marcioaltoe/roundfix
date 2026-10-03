### migrate

```bash
roundfix migrate [--check]
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

With `--check`, Roundfix reads the database without migrating or writing it.
The migration check has these outcomes:

- A current database prints `Run Database is at schema version <n>, the version
  this binary supports: <path>` on stdout and exits `0`.
- An older database prints `roundfix: migrate check: <error>` on stderr and
  exits `2`; the error names `roundfix migrate` as the remedy.
- A newer database prints `roundfix: migrate check: <error>` on stderr and
  exits `2`; the error names `roundfix upgrade` as the remedy.
- An absent database prints `No Run Database at <path>; nothing to migrate` on
  stdout and exits `0`, creating nothing.
- Any other read failure prints `roundfix: migrate check failed: <error>` on
  stderr and exits `1`.

For a current database, the transcript is:

```text
$ roundfix migrate --check
stdout:
Run Database is at schema version <n>, the version this binary supports: <path>
stderr:
exit: 0
```

For an older database, the transcript is:

```text
$ roundfix migrate --check
stdout:
stderr:
roundfix: migrate check: Run Database "<path>" has schema version 21, older than the schema version <n> this binary supports; run 'roundfix migrate' to upgrade it
exit: 2
```

For no Run Database, the transcript is:

```text
$ roundfix migrate --check
stdout:
No Run Database at <path>; nothing to migrate
stderr:
exit: 0
```
