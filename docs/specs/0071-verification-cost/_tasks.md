---
schema: spec-tasks/v1
spec: 0071-verification-cost
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04]
    - id: task_08
      file: task_08.md
      needs: [task_03]
    - id: task_09
      file: task_09.md
      needs: [task_08]
    - id: task_07
      file: task_07.md
      needs: [task_03, task_04, task_06, task_08, task_09]
---

# Tasks — Verification cost

| id      | title                                            | type    | complexity | needs                     |
| ------- | ------------------------------------------------ | ------- | ---------- | ------------------------- |
| task_01 | Record which tests the suite executes             | test    | medium     | —                         |
| task_02 | Let the CLI package take its environment          | backend | high       | task_01                   |
| task_03 | Run the CLI tests in parallel                     | test    | high       | task_02                   |
| task_04 | Free the Baseline package from process state      | backend | high       | task_01                   |
| task_05 | Stop charging every Task for the whole suite      | docs    | medium     | task_01                   |
| task_06 | Assert a suite-time budget                        | infra   | medium     | task_03, task_04          |
| task_08 | Finish declaring parallelism in the CLI package   | test    | high       | task_03                   |
| task_09 | Split the gate into a fast local and full CI tier | infra   | high       | task_08                   |
| task_07 | Publish the before-and-after                      | docs    | low        | task_03, task_04, task_06, task_08, task_09 |

Waves: 1 → task_01 · 2 → task_02, task_04, task_05 · 3 → task_03 · 4 → task_06 ·
5 → task_08 · 6 → task_09 · 7 → task_07

Wave 2 fans out because its three Tasks touch disjoint files: the CLI package,
the Baseline package, and the Task files plus the authoring skill.

**Measurement corrected which Tasks move the headline.** The Spec assumed
sequential execution was the cost and parallelism the lever. After tasks 02–05,
`internal/baseline` improved 40% while `internal/cli` got 27% *worse* and the
suite went from 136.9s to 168.5s. The eight heaviest end-to-end journeys
complete in 31.7s together — they were never the bottleneck. Only 207 of the
CLI package's 488 tests declare parallelism, and of the 281 still sequential
just nine retain a real blocker.

So tasks 08 and 09 are what move the number now. task_08 finishes the
declaration the prefactor made possible.

**Measurement corrected task_09 too.** It was written to strip the heavy
end-to-end journeys out of the local gate. After task_08, a second consecutive
`make verify` costs 5.1s: Go's test result cache re-runs only the packages
whose compiled output changed, and the prefactors in tasks 02 and 04 removed
the `t.Setenv` and `t.Chdir` calls that made that cache untrustworthy. A gate
that already costs seconds has nothing to gain from losing tests, so task_09
kept the local test set whole and moved the tier boundary onto the cache
instead: local trusts it, CI refuses it with `-count=1`. The Pull Request
workflow it adds closes a real hole — until now nothing ran the suite before a
merge, only the release workflow on a tag.

task_01 changes no behavior and lands first. It records which test functions
the suite executes, so "coverage is unchanged" becomes an assertion instead of
a claim. The timing baseline is already committed under `baseline/` and is not
re-derived — comparing against a re-measured "before" would make the comparison
circular.
