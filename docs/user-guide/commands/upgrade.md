### upgrade

```bash
roundfix upgrade [--check]
```

Resolves the latest release through the GitHub CLI. Successful stdout outcomes
are `upgraded 1.0.0 → 1.1.0`, `already current 1.0.0`, `no releases published`,
and, with `--check`, `upgrade available 1.0.0 → 1.1.0`. Failures leave the
current binary untouched and print a manual fallback on stderr. Operational
commands run a best-effort daily freshness check that prints one stderr line
when the binary is behind.

After every successful release outcome, including `--check`, a recommendation
notice is written to standard error. It never changes standard output or the
exit code. After an install, the installed executable runs `profiles check`
in the process working directory under a ten-second timeout; otherwise the
running executable compares in process. A check that cannot run prints one
line, `roundfix: recommendations not checked: <reason>`. Help, usage errors
and failed upgrades print no notice. An upgrade performed by an older
executable prints none; the notice starts with a subsequent upgrade.

