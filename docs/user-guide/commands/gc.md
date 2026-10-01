### gc

```bash
roundfix gc --dry-run
roundfix gc
```

Non-interactive storage reclamation. Resolves `store.journal_retention`,
prunes eligible terminal Runs' Run Event Journal rows and
`<artifact-dir>/runs/<run-id>` directories, removes orphaned `runs/<id>`
directories, and reports Runs, journal rows, and artifact bytes reclaimed.
Live and dry-run reports count only Runs that still hold Run Event Journal rows
or an artifact directory. `--dry-run` lists that reclaimable set without
deleting. Retention never deletes Active Runs, `runs` rows, active-run locks,
or Review artifacts under the Spec Root.

