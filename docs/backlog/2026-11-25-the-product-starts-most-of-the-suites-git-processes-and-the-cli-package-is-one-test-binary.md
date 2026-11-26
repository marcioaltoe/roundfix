---
type: perf
status: open
created: 2026-11-25
spec: null
reason: null
waiting_on: the delivery of Spec 0421, the per-test times its jobs record on the macOS runner and the final rule it writes from the measured result
since: 2026-11-25
---

# The product starts most of the suite's Git processes, and the CLI package's tests are one test binary

## Slow

The whole test suite of this repository, wherever it runs whole: once per Spec on the maintainer's Mac mini, and on the standard macOS runner for each pull request whose change reaches every package, a release's among them. Two costs hold it up, and Spec 0421 measures both without changing either.

- The commands under test start most of the Git processes of a run. The cost is spread over 16,056 leaf tests, so no list of slow tests removes it.
- The CLI package is one test binary. `go test` runs one binary per package and its ten-minute limit is per binary, so no gate ends before this one does.

Spec 0421 reaches its targets by changing what each gate runs and by removing the work the tests repeat. That repeated work is about 6 % of the suite's test seconds, so the cost of the suite itself stays: the complete run on the mini stays near six minutes.

## Measured

On the mini on 2026-11-25, at `33ece0cf`, by the author of Spec 0421. The numbers are copied here because a Spec folder leaves the tree when the Spec is archived; Spec 0421 keeps each command in its `references/sources.md`, sections 4.3 and 5 to 5.3.

Git processes:

- One run of the suite starts at least 153,850 Git processes, counted by a script named `git` placed first on `PATH`. ADR-0259 counted at least 83,660 on 2026-10-08.
- 35,875 of them carry the identity arguments the test fixture passes and 117,975 do not: about three quarters are started by the commands under test.
- By subcommand: `rev-parse` 36,787, `ls-tree` 15,384, `cat-file` 12,974, `add` 9,908, `commit` 9,846, `show` 7,964, `symbolic-ref` 5,761, `init` 4,682.
- One asset sync of the Baseline package starts about 390, and one test of the reconcile command 400 to 1,170.

The CLI package:

- 2,102 top-level tests, against 1,409 on 2026-10-08. It holds 6,091 of the suite's 17,008 test seconds.
- Its binary ran for 390 s in a whole-suite run of 391 s, so it is the critical path. Alone, without its Baseline tests, it took 159 s.
- On the Linux runner it took 295 s. On the macOS runner `go test` killed it at 603 s on a push, and on a pull request it took 562 s plus 68 s for its Baseline tests.
- Its costliest families, in test seconds: review 1,052 over 308 tests, Baseline 1,037 over 145, reconcile 692 over 68, implement 614 over 122, release plan 468 over 30, archive 379 over 123.

The whole suite: `make verify` took 376 s in the lab on the mini.

Not measured: how much of the suite's time the Git processes hold. The count is a count. One Git process that only prints its version costs 3.8 ms there, which puts the starts alone at about 580 process seconds, and a process that reads or writes a repository costs more. Published results for copying a prepared repository in place of `git init` range from 4.5 % to 21 % of a suite's time, as the Research Subagent of Spec 0421 read them. The Spec that adopts this entry measures the share before it promises a number.

## Target

- The complete run of the suite on the mini, under the lab protocol of Spec 0421, takes at most 240 s, against 376 s today. This is the number the maintainer was shown on 2026-11-25 for folding this work into Spec 0421, an option he did not take.
- The CLI package's tests stop being one binary whose time is the suite's time. The number comes from the per-test times Spec 0421 starts recording on the macOS runner.
- Every gate's number in the final rule of Spec 0421 falls or stays. This work raises none.

All three are proposals. The Spec that adopts this entry fixes them after it measures.

## Shape

Not binding. Three directions were named while Spec 0421 was planned:

- Batch the Git calls a command makes. ADR-0259 left it open: "The at least 83,660 Git processes per run and their files are the next cost, and they belong to a later Spec that batches Git in production code."
- Split the CLI package, so `go test` can run its parts side by side and a gate can run only the part a change reaches.
- Copy a prepared repository where a test builds one commit by commit.

## Rules it meets

ADR-0259 rejected the second and the third direction: "Splitting `internal/cli` into several packages would let `go test` overlap them, but it moves a thousand tests and leaves every sequential chain in place", and "Fixture templates copied instead of recreated were measured and left out: the suite creates about 2,000 repositories with `git init`, a small part of the at least 83,660 Git processes it starts." Both bind. The Spec that adopts this entry puts its own Rule Challenge for each one it needs, with the facts ADR-0259 did not have: the seven slowest packages are Parallel Test Packages with at most 40 Sequential Tests, and a run now starts 4,682 `git init` processes.

The final rule of Spec 0421 also binds this work once it exists: a number of that rule rises only by the maintainer's decision, in extreme necessity, for the tests of a large new feature.

## Decisions

The maintainer, on 2026-11-25. The times are the ones the decision records state.

- 07:58 BRT, to whether the Gate Time Targets were the right ones: "Como proposto". The option he did not take put this work inside Spec 0421 (ai-memory message `01a21256` from `default/praetor`, as the brief of Spec 0421's Idea stage relays it; the brief gives the id in this short form).
- 09:15 BRT, to the question of where the "final" is at which the new rule is written: "Final é a 0421, em 2 Specs (Recomendado)". The scope, as the Praetor recorded it: "a regra nova é a última entrega da 0421, com as metas medidas em produção; a parte 2 (processos Git, divisão do pacote da CLI) vira backlog e, sob a regra, só pode baixar números." (ai-memory message `01a2127e-c61d-7a93-b702-ccbc6ebd2d45` from `default/praetor`, answering request `01a2126b-9f73-7791-983f-0b6a4cf191b1`).

## Source

Minted on 2026-11-25 by the Idea stage of Spec 0421, `0421-the-test-suite-meets-a-measured-time-target-at-every-gate`, from the second decision above. No open Backlog Entry or unresolved Finding fitted this intent: the fourteen entries open that day concern other components, and `docs/findings/` was empty. Nature: measured as to the numbers under `## Measured`; the targets and the shape are not decided.
