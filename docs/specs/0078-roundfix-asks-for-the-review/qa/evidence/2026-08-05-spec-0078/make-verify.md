# Authoritative repository Verification

Command: `rtk make verify`

Build: `bdf6ff8d4d680188a97986ee1340ab9ff052a499` plus the in-progress QA report.

Initial result: exit `0`.

- Go test sweep: 3,452 tests passed across 26 packages.
- Isolated Spec corpus budget: 1 test passed.
- Skill repository checks: 4 tests passed.
- `roundfix skills check`: passed for every required owned Skill.
- Build: `bin/roundfix` produced successfully.
- `bin/roundfix spec check`: Spec 0078 reported no findings. Its two
  informational skips were the absent optional Vocabulary Contract section and
  absent `references/_index.md`; neither is an error or a failed check.

The command ran unpiped. Its process exit status, not a downstream formatter or
pager, supplies the verdict.

## Final-worktree rerun and classification

After the report closed, the exact unpiped command ran again. All 3,452 Go
tests passed, but `TestCheckCorpusBudget` measured the full Spec corpus at
`1.489615292s` against its one-second wall-clock budget, so `make verify`
exited 2 before its later targets.

This is environment-caused timing variance rather than a deterministic code or
report failure:

- the initial full command passed on the same build;
- no product or test code changed between attempts;
- the minimal unchanged-worktree reproduction `rtk go test -count=1
  -parallel=1 ./internal/speccheck -run '^TestCheckCorpusBudget$' -v` then
  exited 0;
- the final-worktree failure named elapsed wall time only, while the full
  product test sweep remained green.

Unblocking action: rerun the exact unpiped full gate on a host/session without
the transient load. QA did not retry the whole gate until green. Equivalent
current evidence is the initial complete pass plus the focused current
budget-test pass.
