# Static gate evidence

Build: `e91bf4088b7547ab1f1c4a15c78d1427e769f032`.

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
skip is on an implementation or acceptance path exercised by this gate.

The same unpiped gate ran twice in the scratch owned-Skill journey: once after
the first sanctioned regeneration and once after the confirmed no-op second
regeneration. Both runs exited 0 with the same 3,382-test result and the same
remaining gates green.
