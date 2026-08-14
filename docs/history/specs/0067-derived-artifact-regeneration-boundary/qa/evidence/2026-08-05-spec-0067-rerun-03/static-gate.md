# Static gate evidence

Build: `9e668cfe4649c323694c55f8124ed73840260910`.

The repository's authoritative gate ran unpiped from the Run Worktree root:

```text
rtk make verify
exit 0
Go test: 3382 passed in 26 packages
isolated corpus budget: 1 passed
Skill tests: 4 passed
Roundfix skill check: passed
go build -buildvcs=false: passed
Spec 0067 consistency check: No findings
```

The two Spec 0067 skips concern the absent Vocabulary Contract section and
`references/_index.md`. The PRD does not require either artifact, and neither
skip is on the implementation or acceptance path exercised by this gate.
