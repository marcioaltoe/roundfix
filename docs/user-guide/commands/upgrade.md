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

