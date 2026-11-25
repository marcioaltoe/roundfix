---
spec: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
status: active
created: 2026-11-25
surfaces: [backend, infra, docs]
type: perf
---

# The test suite meets a measured time target at every gate

Every gate of this repository waits for tests its change cannot fail, and no gate is held to a time. On 2026-11-25 the first push to main verified on the standard macOS runner failed, because `go test` killed the test binary of the CLI package at its ten-minute limit: main is red, and the release verifies on Linux under an override that names this Spec as its end. On the maintainer's Mac mini each Task of a Run waits a median of ten minutes for the repository check at its settlement, and in one week those checks held the machine's only verification slot for 15.3 hours. The maintainer refused a longer timeout and a larger budget and asked for the tests and the test sets to be reviewed (the maintainer, 2026-11-25: "É primordial uma revisão dos testes e conjuntos de testes para que tenhamos performance decente."). After this Spec each gate has a **Gate Time Target**, a number measured on the machine that runs the gate. A Task settles on the checks its change can fail, and the whole suite runs once per Spec. Continuous integration runs only for a pull request, in parallel jobs on the standard macOS runner, and a release publishes only a tree such a run verified. Each gate reports its time. The Spec ends by writing one rule from the measured result, the **Gate Time Rule**, which changes only in extreme necessity.

This PRD rewrites the first draft from the approved Idea. The maintainer has decided the shape: the targets, what a Task settlement runs, two Specs, the retirement of two time rules, the allowance and the change clause. On 2026-11-25, at about 14:25 BRT, he also authorized the delivery, and decided how the last delivery is carried, what follows a gate that passes its number and where continuous integration runs. That last answer opened one Reserved Question, what a release verifies now that nothing runs after a merge. The Praetor answered it the same day, at about 15:12 BRT, by applying that decision: the release runs no test and proves its Release Provenance. `## Decisions` records the answer as that application, and the authorization record is approved with it.

## Sources

The request is the maintainer's decision of 2026-11-25, 06:40 BRT. His later answers of the same day shaped it, and `## Decisions` quotes every one. The approved Idea of this Spec holds the research, the trade-offs and the direction he chose. [`references/sources.md`](references/sources.md) holds every measurement with its command and its machine: sections 2 to 5 are the measurements of the first draft, section 4.1.1 the sixteen failure diagnostics read for this rewrite, section 6 the published passages, section 8 the Governed Paths and section 9 what the revision of 2026-11-25 measured for a release.

No source is adopted. The Idea stage found no open Backlog Entry and no unresolved Finding that shares this context, and minted one Backlog Entry for the second Spec: the Git processes the product starts and the size of the CLI package. That entry waits on this Spec. This project's ai-memory inbox was empty when the author listed it for this rewrite, on 2026-11-25.

The Secondbrain and the web were consulted by the two earlier phases, and the Idea records what each source changed. For this rewrite the author asked both one question, what follows a time budget that is passed. The Indexus search of the Secondbrain returned no relevant page, which is a limit of this consultation; the page the Idea used, `wiki/concepts/slo-e-orcamento-de-erro.md`, already says to agree beforehand what changes when a budget is spent. The Exa search returned three practitioner posts, read as excerpts and listed in `references/sources.md`, section 6. They support the recommended option of question 3: a duration judged on single runs fails by the runner's noise, and one of them holds a release on a regression above 25 % of a rolling median. They are not primary sources, and no acceptance row rests on them. Their digest went to the project that owns the Secondbrain as ai-memory message `01a212a5-1d96-7fc2-bccb-5518f86a76c5`.

For the revision of 2026-11-25 the TechSpec author asked both what a release verifies when continuous integration runs only for a pull request. The TechSpec's `## Decisions` lists what each source changed, and `references/sources.md`, section 9, holds the passages and the measurements.

## Measured baseline

"Lab" is the planning worktree on the mini, the maintainer's Mac mini, with nothing else running and the five cores Roundfix gives a Verification there. "Field" is what the Run Database and the workflow runs recorded.

| Gate | What it runs today | Machine | Baseline | How it was measured |
| --- | --- | --- | --- | --- |
| Task settlement, and the precondition of the QA gate | `make verify-changed`, which for one changed Go file runs every package outside the Baseline set | the mini | Field, 2026-11-18 to 2026-11-25: 99 executions in 43 Runs, median 583 s, 15.3 hours in all; the 83 that passed, median 607 s, p90 777 s, p95 882 s, maximum 1,142 s. Lab: 251 s for one comment line in a leaf package, 17 s with no change | Run Database, from `started` to the next event of the same attempt; lab, the command with fresh results |
| Delivery Queue's repository gate | The same command, read from the same key of the Project Config | the mini | No record: the Run Database does not hold the gate's time, and no log of it was found. Lab: 251 s for a small change | The code and the configuration; lab |
| The whole suite, which no local gate runs today | `make verify` | the mini | Lab: 376 s | lab |
| Before a pull request opens | `make verify-docs` | the mini | Lab: 77 s, of which 70 s are the regeneration contracts of the Baseline package | lab |
| Pull request | `make verify-changed`, then `make verify-docs`, in one job | macOS runner | One run, the delivery of Spec 0420: 1,965 s (32.8 minutes). On the Linux runner until 2026-11-24, 35 runs: median 656 s, p90 860 s | step times of the workflow runs |
| Push to main | The timed suite, `make verify-docs`, the Full Contract Run and the budget check, in one job | macOS runner | One run: failed after 1,017 s, the CLI test binary killed at 603 s; the later steps were skipped. On the Linux runner, 31 runs: median 840 s, p90 996 s | the same |
| Release: verification before publication | `make verify`, then the Full Contract Run, in one job | Linux runner since 2026-11-25, under the override | On the macOS runner, none. On the Linux runner, two releases: 814 s and 813 s | the same |

The first draft said the Delivery Queue's repository gate ran `make verify`. It does not: the settlement, the QA gate's precondition and the queue's gate all run the one command this repository configures, `make verify-changed`. The 376 s therefore measure a command no local gate runs.

Where the time goes:

- **One test binary is the critical path of every gate.** The CLI package took 295 s on the Linux runner, 603 s and more on the macOS runner and 372 s in the lab beside the other packages. `go test` runs one binary per package, its ten-minute limit is per binary, and no gate ends before its slowest binary.
- **The suite is bound by the work it does, not by waiting.** One lab run held 6,131 top-level tests and 17,008 test seconds in 391 s of wall time. The 100 costliest tests hold 30 % of the work, and every one of the seven slowest packages already runs its tests in parallel.
- **The suite grew back.** ADR-0259 measured 1,409 top-level tests in the CLI package and a suite of 227 s on 2026-10-08. Seven weeks later the package holds 2,102 tests and the suite takes 391 s on the same machine. The settlement check went from 206 s to 223 s on 2026-09-30, to 350 s on 2026-10-08, to a median of 607 s now.
- **The macOS runner does the same work in more than twice the time.** The packages of one push summed about 1,090 s on the Linux runner and more than 2,230 s on the macOS runner, which has three cores.
- **The product starts most of the processes.** One run of the suite starts at least 153,850 Git processes, and about three quarters of them come from the commands under test. This is the second Spec's work.
- **The same contracts run twice in one job.** `make verify-docs` runs the regeneration contracts of the Baseline package again after another step of the job ran them: 214 s to 372 s on the macOS runner each time.
- **The settlement check caught the repository's own rules and records, and nothing else.** Since 2026-11-01 it failed sixteen times in 125 Runs. The author read the sixteen diagnostics for this rewrite. One was formatting. Two were a contract over a stale record. Thirteen were code tests that check the repository itself: two rules over the test sources, the rule that every contract declares its relevance, the coverage record, and the tests that compare this repository with its Baseline record. One of the thirteen names no test in its log. None was a test of product behavior that the change broke.

The time rules moved each time a gate blocked the work. The Delivery Queue's wait for the checks rose from 30 to 60 minutes when the gate took longer (ADR-0311: "needed 30 to 33 minutes"). The suite budget rose from 240 s to 420 s (ADR-0387: "26 of them exceeded 240 s"). The verification moved to macOS with the budget kept for later (ADR-0452: "The suite-time budget stays 420 seconds until the suite is measured on the macOS runner."). On 2026-11-25 the release's verification returned to Linux by an override. The maintainer named the pattern and inverted it: the rules that block are retired now, and the new rule is written at the end, from the result.

## Gate Time Targets

A gate, in this Spec, is a point of the delivery that runs a verification and can stop the work. The maintainer approved the targets on 2026-11-25: "Como proposto".

| Gate | Machine | Gate Time Target | Judged by |
| --- | --- | --- | --- |
| Task settlement, and the precondition of the QA gate | the mini | Field: median 120 s and p95 240 s. Lab: 45 s, and 150 s when the change selects the regeneration contracts | Field: the first 30 passing settlement checks of Runs that start after the first delivery merges, read from the Run Database. Lab: the protocol the reference document records, the median of three runs |
| The **Complete Suite Run**, once per Spec | the mini | Lab: 360 s | The lab protocol. Each execution in the field is reported |
| Before a pull request opens | the mini | Lab: 90 s | The lab protocol |
| Pull request | standard macOS runner | Median 8 minutes and p95 10 minutes | The first 20 passing runs of that event after the first delivery merges |
| Every test binary under the Go command's default limit | standard macOS runner | 300 s | The slowest binary of each run of the pull request's window |

The approved table also held a push to main and the verification of a release, each at the numbers of a pull request. The maintainer's answer of 2026-11-25 on continuous integration removes the first gate, and the answer on the release, the Praetor's application of that decision, removes the second. A pull request that reaches every package, as a release's does, runs what a push to main ran and is judged in the pull request's window.

What the numbers mean:

- The Complete Suite Run is the one run per Spec, in this repository, of every test of every package with fresh results, before the Spec's pull request opens. The first draft attached the 360 s to the Delivery Queue's repository gate. The target follows the run to whichever gate carries it, and its baseline is the 376 s of the whole suite in the lab.
- The time of a gate in continuous integration is the longest of its verification jobs, each from its own start to its own end. The wait for a runner is reported apart and is not counted, because the approved targets are for a run that did not wait.
- A window counts passing runs. A run that fails for a reason of the change under test is not a time measurement. A run whose test binary is killed at the limit, or whose job is stopped by its timeout, misses the target.
- The first draft also held a number per job and five minutes for a pull request without a Go change. This PRD drops both. The number per job was the change of ADR-0387 the maintainer did not take, and the five minutes were never measured on that runner.
- At the current volume the two windows take about three days.

What leaves each gate, and where it goes:

- **Task settlement.** The package tests leave it. They run once per Spec in the Complete Suite Run, and on each pull request whose change reaches them. The Task's own declared Verification keeps running the tests the Task names. Format, vet, build, the Repository Contract Tests that Contract Relevance selects and the checks of the repository against its own rules and records stay, because those are what the settlement check caught.
- **Pull request.** The packages a change cannot reach through an import leave it. What stays runs in parallel jobs instead of one after the other, and every Repository Contract Test runs on it once.
- **Push to main.** The gate itself leaves. Continuous integration runs only to release the squash merge of a pull request, and nothing runs after the merge.
- **Release.** By the answer on the release its verification leaves too: the release reads the verdict of its own pull request, which ran every package and every contract on the tree that is published.
- **Every gate.** The four code tests that compile or run the tests of other packages leave the code tests.

## The deliveries of this Spec

The maintainer decided where the rule is written: "Final é a 0421, em 2 Specs (Recomendado)". The Gate Time Rule is the last delivery of this Spec, and it needs measurements that exist only after the first delivery merges.

1. **The first delivery** is one pull request, through the Delivery Queue, with the Tasks of `## Sizing` and the QA gate. It retires the two time rules, changes what each gate runs, removes the work the tests repeat, runs continuous integration only for a pull request and in parallel jobs, ends the override of the release and makes each gate report its time.
2. **The measurement windows** open when the first delivery merges: the first 20 passing runs of a pull request and the first 30 passing settlement checks.
3. **The last delivery** writes the Gate Time Rule from the windows. The maintainer decided how the Spec carries it: "PR só de documentos (Recomendado)". It is one closing pull request of this Spec, without a Run, opened when the windows close, with files outside every Governed Path and without a release.

Between the two, no written time rule holds the gates. Three limits of the maintainer hold: no suite-time budget above 420 seconds, no longer `-timeout` on any test invocation, and no paid, larger or self-hosted runner. Nothing new fails a gate for its time in that interval. The Go command's ten-minute limit per test binary and the job timeout that exist today stay the only hard stops, and each gate shows its time against its target.

### The Gate Time Rule

The rule is one new ADR with a number of its own: the next free number when it is written, never the number of ADR-0387, of ADR-0452 or of the ADR that records the retirement. It states:

1. For each gate, the machine, the measured result with its window and dates, and the gate's number. The number is the measured result plus a declared allowance of 25 %, rounded up to the next ten seconds, and never above the approved Gate Time Target. A row with a median and a p95 gets a number for each.
2. The shape of the verification as delivered: what a Task settlement runs, where the Complete Suite Run happens, that continuous integration runs only for a pull request and never after a merge, what a pull request selects, the parallel jobs on the standard macOS runner, and what a release reads before it publishes. The rule records the maintainer's decision of 2026-11-25 on continuous integration, as that decision asks.
3. What follows a gate that passes its number, agreed beforehand. The maintainer's answer: "Julgar pela janela (Recomendado)". Core Feature 16 states it.
4. The change clause. A number rises, or the shape of the verification changes, only by the maintainer's decision on a **Measured Dossier**. A longer timeout and another kind of runner are changes of shape. A number may fall on a new measurement without a decision. A change of machine, such as a new runner image or a new mini, measures the same suite again and goes to the maintainer; it is not a relaxation.

The Measured Dossier carries four things: the large new feature; the time of the gate before and after it, on the machine the rule names; what was already cut, merged or moved to a less frequent gate; and the new number, which is the measured one. It is put to the maintainer as a Rule Challenge of the Gate Time Rule, only in the Idea or the PRD of the feature's Spec. It is never asked during a Run, in a corrective Spec or on the day a gate ran over. Without the dossier the work must fit.

If a window misses an approved target, the result is not reached and the rule waits. The miss goes to the maintainer with its measurement, and no number is written above an approved target.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The ADR that records the retirement takes ADR-0453, the next free number, and the Gate Time Rule takes the next free number when it is written. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the product gains no endpoint, credential or network call. The verification workflow keeps its permissions and its checkout and setup actions, and the release's publication job keeps its permissions and its Trusted Publishing. By the answer on the release, the release's first job gains two read-only permissions, to read a pull request and its check from GitHub with the workflow's own token. The command that reads the windows uses the GitHub CLI session the operator already has, on the mini, outside every gate. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — two accepted rules are retired by the maintainer's answers, recorded in `## Rule challenges`, and no Rule Challenge was asked at this stage. The TechSpec asked two more, which the same section records, both answered. The first is three things ADR-0452 fixes: one job that runs the steps in sequence (ADR-0452: "The CI Verification gate runs its existing steps on `macos-26`."), the same for the release's verification, and the budget sentence (ADR-0452: "The suite-time budget stays 420 seconds until the suite is measured on the macOS runner."); the rest of that ADR keeps binding, the end of the Linux Lane, the deprecated keys, the Doctor Command's checks, the External Test Tool rule and the compile contract among it. The second is ADR-0387 whole (ADR-0387: "The budget is one measured number, in the Makefile."), which also amended the budget clause of ADR-0259. One Task of this Spec writes the accepted ADR-0453, which records the retirement. It goes through the lifecycle (`docs/agents/docs-layout.md`: "A retirement deletes the file."), and a partial change is an amendment (`docs/agents/docs-layout.md`: "An accepted ADR that a later accepted ADR changes in part, without replacing it, stays `accepted`."). The maintainer's answer "Só o leve" changes what two sentences describe for this repository, and ADR-0453 gives each its notice: that a settlement is a full gate (ADR-0257: "Inside a Run the Daemon is the only full gate") and what one costs (ADR-0182: "measured at 206 to 223 seconds on this repository"). The rest is kept. Every Task still settles on the configured command (ADR-0342: "The repository Verification at settlement applies to every Task Graph."), and an Agent in a Run runs neither selected command (ADR-0257: "In a Daemon-assigned turn the Agent runs focused tests of the packages it changed and neither selected command"). The maintainer's answer of 2026-11-25 on continuous integration ends the run after a merge that one decision fixed (ADR-0253: "CI runs the Full Contract Run on every push to main and in the release workflow before any build or publication"). Every Repository Contract Test runs on every pull request instead, and ADR-0453 gives that ADR its notice; by the answer on the release, the same notice withdraws its release clause, and ADR-0453's notice in ADR-0452 withdraws that ADR's sentences on the release's verification job. A contract keeps its Contract Relevance (ADR-0252: "A malformed directive fails the gate rather than skipping a test."). The contracts rule keeps requiring `make verify-docs` before a pull request opens (`docs/agents/specific-repository.md`: "repository contracts validate at the pull request boundary"), and by the same answer continuous integration runs each of its checks once for each event. A code test that reads the live documentation trees stays a contract (ADR-0326: "A test that reads the repository's live documentation trees is a Repository Contract Test."). A test stays parallel unless it names why it cannot, and no ceiling rises (ADR-0259: "a ceiling only falls"). No wait on real time is added (ADR-0271: "A test bounds by real time only the budget it asserts."). The one build of a package stays safe for the tests that share it (ADR-0281: "every test that rebuilds an executable now writes a new file and renames it over the old path"). The Delivery Queue's check wait stays what it is (ADR-0311: "needed 30 to 33 minutes"). Continuous integration keeps owing a fresh verdict (`docs/agents/agent-instructions.md`: "CI must run the selected repository Verification from a fresh run"): every job of a pull request runs from fresh results. The TechSpec adds one obligation it keeps: the Run Marker is read and never set (ADR-0283: "Every process a Run starts for an Agent or a Verification carries the Run Marker."). ADR-0258, ADR-0273, ADR-0282, ADR-0285, ADR-0286, ADR-0295, ADR-0298, ADR-0319, ADR-0384, ADR-0385, ADR-0404, ADR-0423 and ADR-0435 cite one of the ADRs above and decide nothing this Spec changes; the TechSpec's `## ADR Alternatives` records why for each. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`, `docs/agents/specific-repository.md`, `docs/agents/agent-instructions.md`.
- Tooling authority: applicable — the plan proposes protected tooling mutation. On 2026-11-25, at about 14:25 BRT, the maintainer authorized the delivery: the nine paths with their bounds, the regenerations through their own commands, the changes of the decision records, three additions and the operations up to a patch release, adjusted by his answer on continuous integration. `## Decisions` quotes the message. The bound of `.github/workflows/release.yml` that the plan needs, the Release Provenance in the job `verify` with its runner and its read-only permissions, was granted the same day by the Praetor's application of his decision on continuous integration, which `## Decisions` quotes and the authorization record holds as that application. The record is `approved`, granted on 2026-11-25, and must land on the default branch before the Run (`docs/agents/spec-routing.md`: "A new or widened grant must land in the delivery target's ancestry before the consuming squash delivery"). Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0421-the-test-suite-meets-a-measured-time-target-at-every-gate/_authorization.md`; bounded files: `Makefile`, `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `.roundfixrc.yml`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`, `internal/cli/baseline_plan_test.go`, `internal/spec/coverage_test.go`.

## Goals

- Each gate has a Gate Time Target measured on the machine that runs it, and meets it in its window.
- Every change that reaches main, and every tree a release publishes, passed its verification on the standard macOS runner before the merge, the release's with every test and every contract, without a longer timeout, a larger budget or a paid runner. Nothing runs after a merge.
- A Task of a Run settles after the checks its change can fail, and the whole suite is paid once per Spec instead of once per Task.
- A gate that gets slower is seen in the change that causes it, because each gate reports its time.
- The Spec ends with one rule written from the measured result, which holds each gate to its number and changes only on a Measured Dossier.

## User Stories

1. As the maintainer, I want each Task of a Run to settle without waiting ten minutes for tests its change cannot fail, so that this repository stops holding the machine's only verification slot for fifteen hours a week.
2. As the maintainer, I want a pull request verified within ten minutes on the standard runner, so that the Delivery Queue's wait for the required check is short.
3. As the maintainer, I want continuous integration to run only to release the squash merge of a pull request, and a release to publish only a tree such a run verified on the macOS runner, so that nothing runs after a merge and no release repeats the tests its own pull request ran.
4. As the maintainer, I want the time rule written last, from the measured result, and changed only on a dossier I decide, so that a gate that runs over stops being answered by moving the rule.
5. As an Agent in a Run, or the Supervisor that reads its result, I want each gate to print its time and its target, so that I know whether a slow gate is the machine or my change.
6. As whoever adds a test, I want a rule that refuses a test that starts a nested test run, builds a binary of its own or extracts a corpus of its own, so that the suite does not return to where it is today.

## Core Features

Feature 6 follows the answer on the release, the Praetor's application of the maintainer's decision on continuous integration. Features 5, 16 and 18 follow the maintainer's answers of 2026-11-25.

1. **Each gate has a Gate Time Target, recorded with its measurement.** A reference document of the repository holds the table of `## Gate Time Targets`: the gate, the machine, the number, what judges it and the dated measurement it came from. It also holds the lab protocol, the command that reads each window, and the state of each window. Every number a gate reads lives there.
2. **A Task settles on the light checks.** The maintainer's answer: "Só o leve". The Incremental Verification, which the Daemon runs at each Task settlement and as the precondition of the QA gate, formats, vets and builds as today and runs the Repository Contract Tests that Contract Relevance selects as today. It also runs the checks of the repository against its own rules and records that live in code packages today: the rules over the test sources, the rule that every contract declares its relevance, the comparison of the coverage record, and the tests that compare this repository with its Baseline record. It runs no other package test. The Task's declared Verification still runs the tests the Task names. The command prints what it ran and one line that says where the package tests run.
3. **The Complete Suite Run happens once per Spec, in this repository only.** The maintainer's answer: "Só no Roundfix (Recomendado)". Every test of every package runs once per Spec, with fresh results, before the Spec's pull request opens. The product gains no configuration key, the Daemon does not change, and the gates of other repositories stay as they are. The TechSpec proves where the run lives: at the Delivery Queue's repository gate, if this repository's own command can tell the candidate of a Spec from a Task settlement, or at the Spec's QA gate, where a failure returns to the Run's Agent instead of parking the item. If neither works without a new key, a decision request returns to the maintainer before any Task. A Run delivered outside the Delivery Queue gets the run from wherever the TechSpec puts it: if that is the queue's gate, such a Run has no Complete Suite Run before its pull request, and the repository's own guide says so.
4. **A pull request runs what its change can reach.** It is verified with the tests of every package that holds a changed file and of every package that imports one, with every Repository Contract Test and with the checks of feature 2. A change to the module's dependencies or to the build selects every package, as today. So does the pull request of a release, because its version files lie outside every package.
5. **Continuous integration runs only for a pull request, in parallel jobs on the standard macOS runner.** The maintainer's answer: "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main." The verification workflow stops running for a push to main. The steps one job runs in sequence today run in several jobs side by side, each from a fresh checkout with fresh test results. The TechSpec decides how many from the measurement, within the five macOS jobs the account runs at once. The tests of the CLI package run as two or more shards, and a contract proves that the shards run every test of the package exactly once. The required check keeps its name and passes only when every job passed, so no setting of the repository on GitHub changes and the Delivery Queue waits for the same check. Each contract runs once for each event. Outside continuous integration nothing changes: `make verify-docs` stays mandatory before a pull request opens.
6. **A release publishes only a tree that a pull request verified.** The override of 2026-11-25 ends in the first delivery, as its own comment says, and no job of the release workflow runs a test. Before anything is built, the release proves its **Release Provenance**: the tag points at the commit that set the version, that commit is the squash merge of a pull request into main, the pull request's branch held main's tip when it merged, and its required check passed. Any other state refuses the release. The publication job keeps its runner, its permissions and its steps, and still needs that proof. Before any tag, the manual dispatch of the release workflow, which publishes nothing, runs the same proof.
7. **Two time rules are retired in the change that first depends on it.** One short accepted ADR, ADR-0453, records the retirement in the change that brings the parallel jobs. It supersedes ADR-0387, whose file is deleted, and it amends ADR-0452 by withdrawing the three things `## Project Constraints` names. It gives ADR-0259, ADR-0257 and ADR-0182 the amendment notices their changed sentences need. It also records that continuous integration runs only for a pull request, and gives ADR-0253 and ADR-0252 the notices that decision needs. The suite-time budget stops judging a gate, and its timing becomes part of the report of feature 8. The ADR states the three limits that hold until the Gate Time Rule exists, names the approved Gate Time Targets and says that the rule comes as a new ADR with its own number.
8. **Each gate reports its time.** Every gate prints a **Gate Time Report** when it ends: one line with its name, the seconds it took, its Gate Time Target and its state. On the mini the line informs and never fails a gate, because the machine does other work. Each verification job keeps the time of every test it ran, and its summary lists the twenty slowest. One command, run on the mini outside every gate, reads the windows from the Run Database and from the workflow runs and prints each gate's median and p95 against its target.
9. **A code test does not compile or run the tests of another package.** The four tests `## Categories of work` names that start a nested `go test` from inside a code test stop doing so in the default suite. Each becomes a Repository Contract Test with a Contract Relevance that names its inputs, or proves the same thing in its own process. The three that list the packages of the module with the Go tool share one listing per run of their package. A repository rule refuses a new code test that starts `go test` or `go vet`.
10. **A package builds a binary once and extracts a corpus once.** A test that needs the product binary or a stub binary takes the one its package built for that run, and a repository rule refuses a build inside a test outside that one place. A test that reads the pinned history reads one extraction per package and never writes into it.
11. **A test that only reads an expensive result shares it.** The asset-sync and transaction tests of the Baseline package, which each prepare the same source and run the same whole sync before the part they assert, take a prepared copy. A test that asserts the sync itself keeps running it.
12. **A test waits on its event.** The four tests that wait about ten seconds each for an allowance to pass prove the same behavior without the wait: on the event they read, or with an allowance the test gives the code under test.
13. **No test is lost.** No assertion, test case or platform is removed to make a gate faster, and no test is skipped on the macOS runner. The coverage record is re-recorded by its own command, and the only names it may lose or gain are the ones the Task Graph lists, each a test that moved to a contract, was merged into a named test with its assertions, or was renamed.
14. **The first delivery is measured on the macOS runner before it merges.** The delivering pull request is itself verified by the new jobs, and the Delivery Queue waits for its required check as today. The time of each job of that run is recorded in the reference document before the release.
15. **The Gate Time Rule is written last, from the measured result.** When the windows have closed and every approved target is met, the last delivery writes the rule `### The Gate Time Rule` describes, as an accepted ADR with a number of its own, and the three limits of the interval end with it. The reference document then holds each gate's measured result and its number.
16. **A gate whose window passes its number holds the next release.** The maintainer's answer: "Julgar pela janela (Recomendado)". Under the Gate Time Rule a Gate Time Report has three states: within the measured result, in the allowance, and over the number. A single run over its number never fails for its time, because a hosted runner is noisy. After the rule a gate's window is its latest 20 passing runs, or its latest 30 passing settlement checks. The windows are read with the command of feature 8 before each release. A gate whose window is over its number holds the release until a change brings it back, by cutting, merging or moving work, or until the maintainer decides on a Measured Dossier.
17. **The Gate Time Rule changes only in extreme necessity.** The maintainer's answer: "Sua decisão, com dossiê (Recomendado)". The change clause of `### The Gate Time Rule` is part of the rule, and the repository's own guide repeats it where an author of a Spec reads it.
18. **The last delivery changes documents only.** The maintainer's answer: "PR só de documentos (Recomendado)". The first delivery puts every number a gate reads in the reference document and writes the guide's sentences so that they stay true after the rule. A contract of the first delivery proves that every number of the reference document is at or below its approved target and, once a window is recorded, equals the measured result plus the allowance. The last delivery therefore changes no Governed Path, no test and no workflow: it adds the ADR, ends the retirement ADR's interval through the lifecycle, and updates the reference document and the glossary's link.
19. **The glossary and the guide say what each gate runs.** The glossary gains the terms of `## Glossary`. The repository's own guide and the comments of the Makefile say what each gate runs after this Spec, and the rules of features 9 and 10 are written there.

## User Experience

The lines below show the intended shape. The TechSpec fixes every byte in its Surface Transcripts.

The Incremental Verification at a Task settlement, for a change to one file of one package:

```text
verify-select: packages: none at a settlement; the Complete Suite Run runs them once per Spec
verify-select: contracts: 92 selected (90 always, 2 by change), 31 not selected, 0 boundary
gate time: task settlement 31s (target: median 120s, p95 240s)
```

The Complete Suite Run:

```text
gate time: complete suite run 352s (target: 360s)
```

A verification job on the macOS runner:

```text
gate time: pull request, job cli-2 412s (target: median 480s, p95 600s)
```

The command that reads the windows, while one is still open:

```text
gate times: task settlement, 30 of 30 passing: median 104s of 120s, p95 212s of 240s, within
gate times: pull request, 14 of 20 passing: window open
```

## Categories of work

Test seconds are from the lab run of the whole suite; `references/sources.md`, section 5, has every figure.

| Category | What was found | Costliest cases | Where it goes |
| --- | --- | --- | --- |
| Slow: a Go tool over the module, or a build, inside a code test | Seven tests run the Go tool over the module from inside a code test: four start a nested `go test`, three list its packages. The product binary is built inside a test in six places, and a stub binary once per call by a fixture fifteen test files use | `TestCoverageCollectionMatchesTheToolchainOnThisPlatform` 123; `TestJournalConsumerCorpusReplaysEveryConsumer` 103; `TestAgentSelectionProfilesMacro` 55 in the lab and more than 4 minutes on the macOS runner | Task 1 |
| Slow: each test prepares the same input | Thirteen tests extract the same pinned corpus. Eighteen asset-sync tests each build a source and run the whole sync two to four times | The five corpus tests of one package 310; the asset-sync tests 619; `TestTransactionFailureMatrix` 388 | Task 2 |
| Slow: the commands under test start many Git processes | At least 153,850 Git processes per run. The time is spread over 16,056 leaf tests, so no list of tests removes it | `rev-parse` 36,787 calls; one asset sync about 390; one reconcile test up to 1,170 | The second Spec |
| Duplicated | The regeneration contracts run twice in one job. The whole suite is verified at every Task settlement, again at the queue's gate, again on the pull request, again on main and again for the release of the same commit | 214 s to 372 s per repeated contract run on the macOS runner; 99 settlement checks a week at a median of 583 s | Tasks 3 and 4 |
| Waiting for a time to pass | Four tests wait about ten seconds each. The Real-Time Budget rule already removed the rest | About 45 in all | Task 2 |
| Not hermetic | None found. The whole suite passes with an empty home directory, without the Codex path variable and without the user's tools | One more subtest skips there, for a formatter that is an External Test Tool | Nothing to do |
| Serial | None worth a Task. The seven slowest packages are Parallel Test Packages, and eight tests at once in a package measured the same as sixteen | 395 s against 391 s for the whole suite | Nothing to do on the mini; the macOS runner is measured in the delivering pull request |
| Out of a gate | The package tests at a settlement; the packages a pull request's change cannot reach; the four nested test runs in every code-test run | 251 s in the lab for one changed line, against 17 s for the checks that stay | Tasks 1 and 3 |

## Sizing

The first delivery is four implementation Tasks and the QA gate, in one context: the gates of this repository and the tests they run.

| Task | What it changes | Measured basis |
| --- | --- | --- |
| 1. Code tests stop running other packages' tests and building per test | The seven tests named beside, the tests that build the product binary or a stub inside a test, and the two repository rules that keep both out | Nested `go test`: `TestCoverageCollectionMatchesTheToolchainOnThisPlatform` 123 test seconds; `TestJournalConsumerCorpusReplaysEveryConsumer` 103, which is 93 % of the wall time of its package; `TestPartitionFollowsTheMakefileRecipes` 23; `TestDerivedQAVerificationPassesTheChecker` 11. A listing of the module: `TestCoverageEquivalence` 13; `TestCoverageRecordCountsNoPackageUnderDocs` 11; `TestPartitionCoversEveryTestExactlyOnce` 9. A build: `TestAgentSelectionProfilesMacro` 55 in the lab |
| 2. Tests share what they only read, and wait on events | One extraction of the pinned history per package; prepared sources for the asset-sync and transaction tests of the Baseline package; the four waiting tests | Thirteen tests extract the corpus, five of them 310 test seconds; eighteen asset-sync tests 619; `TestTransactionFailureMatrix` 388; the waiting tests about 45 |
| 3. A Task settles on the light checks, the Complete Suite Run gets its home, and each gate reports its time | What the Incremental Verification selects; where the Complete Suite Run happens; what a pull request selects; the Gate Time Report; the command that reads the windows; the reference document | Lab: 17 s for the checks that stay, against 251 s today; the sixteen failures the settlement check caught; the CLI package imports 40 of the 49 other packages |
| 4. Continuous integration verifies in parallel jobs, and two time rules are retired | The two workflows, with continuous integration on a pull request only and the Release Provenance of a release, the shards of the CLI package and their contract, the per-test times, ADR-0453 with the retirement and the notices, the end of the budget's judgment, the glossary and the repository's own guide | macOS runner: the CLI package 630 s as one binary; the contracts 372 s once and 214 s again; the other packages 406 s |
| QA gate | The lab protocol of every target, the contracts, the control cases of the settlement, the coverage record | |

The last delivery follows the windows. It is one closing pull request of documents by the Supervisor, without a Run, as the maintainer decided.

The second Spec follows this one and is not minted yet. It waits as a Backlog Entry for the Git processes the product starts for each command (ADR-0259: "they belong to a later Spec that batches Git in production code") and for the size of the CLI package. Under the Gate Time Rule it may only lower numbers.

## Non-Goals / Out of Scope

- A longer `-timeout` on any `go test` invocation, and a suite-time budget above 420 seconds: the maintainer refused both.
- A larger, paid or self-hosted runner, a paid plan and any setting of the repository on GitHub. The required check keeps its name.
- A written time rule between the retirement and the Gate Time Rule, and any new failure of a gate for its time in that interval.
- A new key of the Project Config or of the User Config, any change to the Daemon and any other change of the product: the maintainer's answers "Só o leve" and "Só no Roundfix (Recomendado)".
- A product capability for two verifications, and gate time recorded by the product.
- Batching the Git processes of the product, splitting the CLI package and copying repository fixtures from a template: the second Spec.
- Removing an assertion, a test case or a platform to gain time, and skipping a test on the macOS runner.
- Relying on cached test results in continuous integration: every job keeps owing a fresh verdict.
- A change to the contracts rule outside continuous integration: `make verify-docs` stays mandatory before a pull request opens.
- Any run of continuous integration after a merge: the maintainer's answer of 2026-11-25.
- A setting of the repository on GitHub that requires a branch to be up to date before it merges. The release checks that for its own commit.
- The rest of ADR-0452: only what `## Project Constraints` and Core Feature 7 name is withdrawn, with its sentences on the release's verification job by the answer on the release.
- The race detector in any gate.
- Any Baseline module, template or profile, any skill, `go.mod` and `go.sum`.
- The gates of every other repository that uses Roundfix.
- What happens to main and to the releases while this Spec is not delivered: decided apart, by the override of 2026-11-25.

## Success Metrics

1. Success Metric: on the lab protocol of the reference document, the Incremental Verification at a settlement, for one comment line added to one file of the CLI package, ends within the lab number of its Gate Time Target, and its output names no package test beyond the checks of Core Feature 2. On the unchanged tree the same command takes 251 s and runs every package outside the Baseline set.
2. Success Metric, a control that holds on both trees: on a scratch copy, each kind of failure the settlement check caught still fails the Incremental Verification at a settlement, with the diagnostic that names it. The kinds are a file that needs formatting, a stale record of command surfaces, a Real-Time Budget above its file's ceiling, a contract file that is not declared, a test added without re-recording the coverage record, and a `Makefile` change without the regeneration of the Baseline record.
3. Success Metric: the selection a pull request uses, for one comment line added to one file of a package that only the CLI package imports, names that package and the CLI package and no other; for a change to the module's dependency file, and for the version files of a release, it names every package. On the unchanged tree the first change selects every package outside the Baseline set.
4. Success Metric: on the lab protocol, the Complete Suite Run ends within its Gate Time Target, and a scratch delivery of a Spec with two Tasks shows the package tests running once, at the gate the TechSpec names, and at neither settlement. On the unchanged tree the whole suite takes 376 s in the lab and every settlement runs the package tests.
5. Success Metric: the repository rule fails, naming the file, on a code test that starts `go test` or `go vet` and on a test that builds a binary outside its package's one place, and passes on the tree. On the unchanged tree the rule does not exist and four tests start a nested `go test` in the default suite.
6. Success Metric: the contract over the shards of the CLI package fails when a test of the package is in no shard or in two, and passes on the tree. On the unchanged tree there are no shards.
7. Success Metric: the tests that read the workflows pass only when the verification workflow triggers on a pull request and on nothing else; when every job that runs a test names the standard macOS runner; when the required check passes only if every verification job passed; when a pull request runs every Repository Contract Test once; when no job of the release workflow runs a test and its publication job needs the Release Provenance; when no test invocation carries a longer `-timeout` than today; and when no suite-time budget above 420 seconds is in the tree. On the unchanged default branch the verification workflow also triggers on a push to main, and the release's verification names the Linux runner.
8. Success Metric: the decision records agree with the retirement. ADR-0387 is absent, ADR-0453 is accepted and names what it supersedes with its date, ADR-0452 carries its amendment notice and keeps every other decision, and the contract that pins that record passes. On the unchanged tree ADR-0387 is accepted.
9. Success Metric: the coverage record lists every test name the unchanged tree lists except the ones the Task Graph names with their destination, and the documentation, glossary and contract checks pass.
10. Success Metric: every gate prints its Gate Time Report in the form of `## User Experience`, proved by one named test per gate that asserts the whole line, and the command that reads the windows prints the lines of `## User Experience` from a recorded fixture of runs and settlements.
11. Success Metric: the contract over the reference document passes with the approved targets alone and with a recorded window whose number is the measured result plus the allowance, and fails when a number is above its approved target or above the measured result plus the allowance.
12. Success Metric, outside evidence: the Gate Time Target of the pull request is at or below the upper limit DORA publishes for a build, quoted in `## Acceptance evidence`, and the runner label every verifying job names is one of the labels GitHub's documentation lists for the standard macOS runner on arm64, with neither the `-large` nor the `-xlarge` suffix of a billed runner.
13. Success Metric: the Release Provenance step passes for the squash merge of a pull request whose branch held main's tip and whose required check passed. It fails, naming the reason, for a pull request that merged behind main, for a check that failed or is missing, for a commit that no pull request merged, for a tag on a commit that did not set the version, and when GitHub cannot be read. On the unchanged tree the release runs the suite and reads nothing of the pull request; of the last fifteen release commits, fourteen meet the four conditions and the one of `v0.148.1` merged behind main.

## Acceptance evidence

The outside-evidence row of Success Metric 12 rests on published literature, read on 2026-11-25 and copied in `references/sources.md`, section 6. DORA, "Capabilities: Continuous integration": "The tests should not take more than a few minutes to run, with an upper limit of about 10 minutes according to DORA’s research (PDF). If your build takes longer than this, you should improve the efficiency of your tests, add more compute resources so you can run them in parallel, or split out longer-running tests into a separate build using the deployment pipeline pattern." GitHub Docs, "GitHub-hosted runners reference": "Use of the standard GitHub-hosted runners is free and unlimited on public repositories." The same page lists, for the standard macOS runner on arm64, the workflow labels `macos-latest`, `macos-15`, `macos-26` and `xcode-27`, with 3 processors and 7 GB of memory. The QA gate checks the row by reading the reference document and the two workflow files against the passages quoted here, without the network.

The ten-minute limit this Spec keeps is the Go command's own (`go help testflag`: "If a test binary runs longer than duration d, panic. If d is 0, the timeout is disabled. The default is 10 minutes (10m).").

## Unreachable Acceptance

- criterion: the pull request that delivers this Spec passes the required check on the macOS runner, with every job reporting its time
  reason: the jobs run on GitHub's runner and exist only once the pull request is open, after the QA gate
  satisfied-by: the Delivery Queue waits for the required check of the delivering pull request before it merges; the operator reads the job times of that run and records them in the reference document before the release
- criterion: the release workflow proves the Release Provenance of a real release commit from GitHub and reaches its publication job
  reason: the release runs only on GitHub, for a manual dispatch or a pushed tag, after the release's pull request merges
  satisfied-by: the operator runs the release workflow's manual dispatch, which publishes nothing, after the release's pull request merges and before pushing the tag
- criterion: the first 20 passing runs of a pull request after the merge meet the Gate Time Target of a pull request, and none is stopped for its time
  criterion: in the Run Database, the first 30 passing settlement checks of Runs that start after the merge meet the field numbers of the Gate Time Target of a Task settlement
  reason: the runs and the settlements exist only after the merge
  satisfied-by: the operator reads the windows with the command the reference document records and writes each result into that document; a miss goes to the maintainer with its measurement
- criterion: the Gate Time Rule is an accepted ADR with a number of its own, written from the two windows, and the reference document holds each gate's measured result and its number
  reason: the rule fixes measurements that exist only after the windows close, about three days after the merge
  satisfied-by: the last delivery of this Spec, the closing pull request the maintainer decided; the contract over the reference document proves its numbers

## Risks

1. **A break reaches the end of the Spec.** A Task can break a behavior test its declared Verification does not name, and the Run learns it at the Complete Suite Run. The record gives both sides. In the sixteen failures the author read, none was a behavior test, and earlier the same check found little (ADR-0257: "failed 6 times in 125"). Without any check at a settlement, two Runs of an adopter ended with a broken tree (ADR-0342: "Both were caught only by the branch's own gate before the Pull Request."), which is where this plan finds such a break too. The reference document counts each Complete Suite Run that fails on a test the old check would have run; two in thirty Specs returns the composition of the settlement to the maintainer.
2. **The light checks may not fit 45 s.** The checks that stay take 17 s in the lab today without the rule and record checks of Core Feature 2. The lab run and the diagnostics put those at about 25 s when all of them are selected, and that sum was never measured as one command. The lab number is approved, so a miss is found by the QA gate and goes to the maintainer.
3. **The macOS runner was measured twice.** Every number of that runner comes from one pull request and one failed push. Core Feature 14 makes the delivering pull request the first real measurement, before anything merges.
4. **macOS jobs queue.** GitHub runs at most five macOS jobs at once on the account's plan, across its repositories. Two pull requests that overlap wait for each other. The wait is reported and not counted.
5. **The lab is not the field.** The same command takes 251 s in the lab and a median of 607 s in the field, where other Runs use the machine. The field numbers are confirmed only by the window.
6. **For some days the repository has no written time rule.** The three limits and the reported times are all that hold. The suite grew 72 % in seven weeks under a rule, so the interval is kept short: the windows take about three days.
7. **The last delivery depends on someone opening it.** It happens days after the merge, outside the Run that delivered the rest. The maintainer chose the closing pull request by the Supervisor over a Spec of its own, and the Spec's `## Unreachable Acceptance` keeps it owed.
8. **A strict change clause can stop a feature.** The dossier is the exit the clause names for itself, so it must be possible to assemble one: the command of Core Feature 8 gives the time before and after, on the machine the rule names.
9. **Moving a test to a contract changes when it runs.** A relevance that misses an input lets a break pass a Task settlement. The pull request finds it before the merge, because it runs every contract; until 2026-11-25 the net was after the merge (ADR-0253: "A miss of the selection rule now surfaces on the push to main, after the merge, rather than before it").
10. **This Spec does not make the whole suite fast on one machine.** Its repairs remove about 6 % of the test seconds; the rest of the gain is composition. The Complete Suite Run stays near six minutes until the second Spec.
11. **Nothing verifies main after a merge.** The repository does not require a branch to be up to date before it merges, so two pull requests that each pass can meet on main in a tree neither verified (GitHub Docs, "About protected branches", on such a check: "Status checks may fail after you merge your branch if there are incompatible changes with the base branch."). The next pull request that reaches the broken package finds it. At the latest the release's pull request does, because it runs every package and every contract and must hold main's tip to be released.
12. **A release depends on GitHub's record of a pull request.** If GitHub cannot be read the release is refused and run again. A release commit that merged behind main is refused for good, and the release is cut again as the next patch version: one of the last fifteen releases would have been.

## Idea-stage judgment

- item: the request in the session
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: yes (0.88)
- unverified_assumption: no (0.13)
- market_question: no (0.03)
- verdict: indicated
- stage: ran
- shape: split — two Specs: this one, for what each gate runs, the work the tests repeat, parallel jobs in continuous integration and the final rule; and a second, held as a Backlog Entry, for the Git processes of the product and the split of the CLI package.

The judge indicated the stage and the drafted PRD recommended dismissing it. The maintainer chose to run it on 2026-11-25, at 07:58 BRT by the record of his decision: "Rodar a Ideia". The judgment is the one the PRD phase obtained; this stage did not run the judge again.

Besides the split, the stage moved the new rule from a Task of the first delivery to the last delivery of the Spec, retired ADR-0387 where the draft proposed to change it, and made the home of the complete run a proof the TechSpec owes. The research agrees with the judge that the shape was open. It disagrees mildly on the second criterion: the draft rested on one assumption that proved false, about what the queue's gate runs, and on one still unchecked, the classification of the sixteen failures.

## Rule challenges

- rule: ADR-0452, `docs/adr/0452-roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane.md`
- imposes: one job runs the verification steps one after the other on `macos-26`, for a pull request, a push to main and a release, and the suite-time budget stays 420 seconds until the suite is measured on that runner
- prevents: a push to main that passes, because the packages of one push sum more than 2,230 s on three cores and the largest test binary passes the ten-minute limit; a pull request verified within ten minutes; a release verified on macOS
- asked: 2026-11-25, prd
- answer: retire — "Retirar a regra e ao final, quando chegarmos ao resultado, produzimos uma nova regra. Pois o que aconteceu até agora é que toda vez que o tempo é estourado, a regra é modificada seja de tempo ou qualquer outra coisa. Então dessa vez quero inverter e no final criar uma regra para que só seja modificar em extrema necessidade por testes de grandes novas features.", 2026-11-25
- effect: the PRD states the retirement as a requirement. A Task of this Spec withdraws the sentences about the grouping in one job, the runner type and the 420 seconds, in the change that first depends on it and through the ADR lifecycle; the rest of ADR-0452 keeps binding (ADR-0452: "macOS runs the tests; Linux and Windows builds and tests are compiled but not executed."). The jobs are designed from the measurement, and the final rule, the last delivery of this Spec, states what replaces the withdrawn sentences.

- rule: ADR-0387, `docs/adr/0387-the-suite-time-budget-is-measured-on-ci-and-judged-after-the-contract-runs.md`
- imposes: one number, 420 seconds, judges the wall time of the whole suite in the job that ran it, and changing the number takes a new measurement
- prevents: judging the suite once it runs in parallel jobs, because no job runs it whole; a budget measured on the macOS runner, where the packages of one push sum more than 2,230 s against the 209 s to 343 s of the Linux runs the 420 seconds came from; a green main in this Spec while the rule stands as written
- asked: 2026-11-25, idea
- answer: retire — "Retirar e criar uma nova no fim, não use o mesmo número de referencia.", 2026-11-25
- effect: the PRD states the retirement as a requirement, and a Task of this Spec retires ADR-0387 in the change that first depends on it, through the ADR lifecycle. Asked which number must not be reused, the maintainer answered "O número da ADR": the final rule is a new ADR with a number of its own, neither a rewrite of ADR-0387 nor of ADR-0452. Until it exists the three limits of 2026-11-25 hold, and each gate shows its time against its target.

- rule: "repository contracts validate at the pull request boundary", `docs/agents/specific-repository.md`, with the run after a merge that ADR-0253 decides
- imposes: the checks `make verify-docs` lists run in every event of continuous integration, whether or not another step of that event ran them, and every contract runs again on each push to main
- prevents: five jobs with time to spare, because the repeat of the regeneration contracts took 214 s on the macOS runner; a sixth job, which is one more than the account runs at once
- asked: 2026-11-25, techspec
- answer: change — "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main.", 2026-11-25
- effect: the PRD states the change as a requirement, in Core Features 4 and 5. Continuous integration runs only for a pull request, where each contract runs once, and the run on a push to main ends. Outside continuous integration the rule stands: `make verify-docs` stays mandatory before a pull request opens. A Task of this Spec rewrites the rule's sentence on where the Full Contract Run runs, and ADR-0453 gives ADR-0253 and ADR-0252 their notices.

- rule: ADR-0253, `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`, its release clause, with the sentences of ADR-0452 on the release's verification job
- imposes: the release workflow runs the whole suite and the Full Contract Run in a job of its own, on the macOS runner, before any build or publication
- prevents: a release that runs nothing after the merge, as the maintainer's answer asks of continuous integration; a release that does not repeat, on the same tree, what its own pull request ran; a release that refuses a commit no pull request verified
- asked: 2026-11-25, techspec
- answer: change — the Praetor's application of the maintainer's decision 01a213bf-b859, answer to 01a213d8-c4d8: "Vale a opção 1, com o alcance exatamente como no seu pedido", applying "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main.", 2026-11-25
- effect: the PRD states the change as a requirement, in Core Feature 6: the release runs no test and proves its Release Provenance before anything is built. A Task of this Spec changes the job `verify` of the release workflow within the bound the authorization record holds, and ADR-0453 gives ADR-0253 the notice that withdraws its release clause and ADR-0452 the notice that withdraws its sentences on the release's verification job. Reusing the verdict of a pull request for the release of the same commit is no longer out of scope. `## Decisions` records whose answer this is.

## Skills

- unchanged: command release plan — the release runbook gains the Release Provenance, the update of the release's branch before its merge, the manual dispatch before the tag and the reading of the windows before a release; the release reference describes the release plan and the tag-triggered workflow without naming the workflow's steps, so it needs no text change
- unchanged: guide docs/user-guide/release-runbook.md — the same runbook edits; its covering release reference names neither the workflow's steps nor their order, so it needs no text change
- unchanged: command resolve — the repair of the test that waits for the detach handshake may give a test its own allowance through internal/cli/detach.go; no flag, output, exit code or detached behavior changes, so the review-runs reference needs no text change
- unchanged: command watch — the same possible change of internal/cli/detach.go; no flag, output, exit code or detached behavior changes, so the review-runs reference needs no text change

## Glossary

- adds: **Gate Time Target**
- adds: **Gate Time Report**
- adds: **Gate Time Rule**
- adds: **Complete Suite Run**
- adds: **Measured Dossier**
- adds: **Release Provenance**
- changes: **Full Contract Run**

## Decisions

Every decision below is the maintainer's, of 2026-11-25, except the last dated line. Each of his reached this project as a `DECISÃO DE MARCIO` block that the Praetor, the session that coordinates the software projects, sent through ai-memory from `default/praetor`. The last dated line is not such a block: it is the Praetor's application of one of his decisions, recorded as the Praetor asks. The text in quotes is verbatim, and the times are the ones the messages state.

- 2026-11-25: The full verification of main and of a release fails on macOS because the CLI package passes the default ten-minute limit; raise the timeout to 30 minutes, return main and the release to Linux, or take a larger macOS runner? — "Essa verificação do roundfix de 10 minutos já seria um absurdo, ainda mais se rodar em todas as runs, 30 minutos é ainda mais absurdo. É primordial uma revisão dos testes e conjuntos de testes para que tenhamos performance decente." Message `01a2121b-9124-7571-8f8f-de0f40876364`, 06:40 BRT.
- 2026-11-25: Does the Idea stage run? — "Rodar a Ideia". Message `01a21256`, 07:58 BRT.
- 2026-11-25: Are the Gate Time Targets the right ones? — "Como proposto". The same message. The option not taken put the second Spec inside this one.
- 2026-11-25: What does a Task settlement run? — "Só o leve". The same message. Its scope as relayed: "nenhuma chave nova de configuração e nenhuma mudança no Daemon".
- 2026-11-25: Does the rule of ADR-0452, one job on the macOS runner, change? — "Retirar a regra e ao final, quando chegarmos ao resultado, produzimos uma nova regra." The same message; `## Rule challenges` holds the whole answer.
- 2026-11-25: Where is the "final" at which the new rule is written? — "Final é a 0421, em 2 Specs (Recomendado)". Message `01a2127e-c61d-7a93-b702-ccbc6ebd2d45`, 09:15 BRT.
- 2026-11-25: Does ADR-0387 get the same treatment as ADR-0452? — "Retirar e criar uma nova no fim, não use o mesmo número de referencia." Asked which number, he answered "O número da ADR". The same message, 09:17 BRT.
- 2026-11-25: Which number does the final rule fix at each gate? — "Medido + folga declarada (Recomendado)". The same message, 09:15 BRT.
- 2026-11-25: How is the "extreme necessity" judged? — "Sua decisão, com dossiê (Recomendado)". The same message, 09:15 BRT.
- 2026-11-25: Light per Task and complete once per Spec: only in the Roundfix repository, or a capability of the product? — "Só no Roundfix (Recomendado)". The same message, 09:17 BRT.
- 2026-11-25: Is the delivery of this Spec authorized: nine Governed Paths with their bounds, the regenerations by their own commands, the retirement of ADR-0387 and of parts of ADR-0452 by a short new ADR, continuous integration in parallel jobs on macOS, up to a patch release? — "Autorizar (Recomendado)". Message `01a213bf-b859-7f60-91a1-201ca006c98a`, about 14:25 BRT. Its scope adds the three additions the TechSpec found and adjusts the whole by the fourth answer below.
- 2026-11-25: How is the final rule delivered? — "PR só de documentos (Recomendado)". The same message.
- 2026-11-25: What happens when a gate passes its number? — "Julgar pela janela (Recomendado)". The same message.
- 2026-11-25: Does the rule "repository contracts validate at the pull request boundary" change in continuous integration? — "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main." The same message. `## Rule challenges` holds its effect.
- 2026-11-25: What does a release verify, now that continuous integration runs only for a pull request: its Release Provenance, a minimal publication set, or the whole verification as already granted? — "Vale a opção 1, com o alcance exatamente como no seu pedido". Aplicação da DECISÃO de Marcio 01a213bf-b859 pelo Praetor, resposta 01a213d8-c4d8: ai-memory message `01a213da-6946-7920-bae3-fedc57a48923` from `default/praetor`, about 15:12 BRT. It is not a `DECISÃO DE MARCIO` block; the Praetor applied the maintainer's answer on continuous integration without a new question, and a later disagreement of the maintainer becomes an amendment. The message is quoted below.

The first decision, verbatim. It answers the request `01a2109e-f31b-74e3-a162-d0f700280f5c`:

```text
DECISÃO DE MARCIO
pedido: 01a2109e-f31b-74e3-a162-d0f700280f5c
pergunta: A verificação completa do main e do release falha em macOS porque o pacote internal/cli passa do tempo-limite padrão de 10 min. Opções: -timeout 30m nos alvos test e test-budget; voltar main e release ao Linux; runner macOS maior.
resposta (literal): "Essa verificação do roundfix de 10 minutos já seria um absurdo, ainda mais se rodar em todas as runs, 30 minutos é ainda mais absurdo. É primordial uma revisão dos testes e conjuntos de testes para que tenhamos performance decente."
alcance: NÃO aplicar -timeout 30m nem subir o SUITE_BUDGET_SECONDS. Abrir uma Spec de performance da suíte de testes do Roundfix (o que roda em cada Run, em PR, em main e em release; testes lentos, duplicados, que esperam timeout, que não são herméticos; paralelismo; o que pode sair do gate de cada Run). Autoria da Spec sem Run; mande ao Praetor o plano com uma meta de tempo medida (proposta sua, por gate: Run, PR, main, release) e, separado, o que fazer com o main vermelho e a v0.153.3 enquanto a Spec não sai, como PEDIDO DE DECISÃO com opções. Nada de mudança em Makefile ou workflow sem nova decisão.
quando: 25/11/2026 06:40 BRT
```

The author did not read the message `01a21256`. The brief of the Idea stage relays its four answers in the words quoted above and gives the id in this short form.

The five blocks of message `01a2127e-c61d-7a93-b702-ccbc6ebd2d45`, verbatim. They answer the request `01a2126b-9f73-7791-983f-0b6a4cf191b1`:

```text
DECISÃO DE MARCIO
pedido: 01a2126b-9f73-7791-983f-0b6a4cf191b1 (pergunta 2)
pergunta: Onde é o "final" em que a regra nova é escrita? Opções: final é a 0421, em 2 Specs; final é a segunda Spec; uma Spec só.
resposta: "Final é a 0421, em 2 Specs (Recomendado)"
alcance: a regra nova é a última entrega da 0421, com as metas medidas em produção; a parte 2 (processos Git, divisão do pacote da CLI) vira backlog e, sob a regra, só pode baixar números.
quando: 25/11/2026 09:15 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 3)
pergunta: A ADR-0387 (orçamento único de 420 s) recebe o mesmo tratamento da ADR-0452? Opções: retirar e reescrever no fim; mudar agora; manter.
resposta (literal): "Retirar e criar uma nova no fim, não use o mesmo número de referencia." Esclarecido em seguida ("Qual número não deve ser reaproveitado?"): "O número da ADR".
alcance: a ADR-0387 é retirada agora; a regra final é uma ADR nova, com número próprio, e não uma reescrita da 0387 nem da 0452. No intervalo valem os limites de 25/11 (nada acima de 420 s, sem -timeout maior, sem runner pago), com cada gate mostrando o seu tempo contra a meta.
quando: 25/11/2026 09:17 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 4)
pergunta: Que número a regra final fixa em cada gate? Opções: medido + folga declarada; as metas como estão; medido sem folga.
resposta: "Medido + folga declarada (Recomendado)"
alcance: o resultado medido mais uma folga declarada no PRD, nunca acima da meta aprovada.
quando: 25/11/2026 09:15 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 5)
pergunta: Como se julga a "extrema necessidade"? Opções: decisão de Marcio com dossiê medido, só na autoria da feature; exceção com prazo; caso a caso.
resposta: "Sua decisão, com dossiê (Recomendado)" (a opção era: decisão de Marcio, com dossiê medido e só na autoria da feature)
alcance: o pedido só nasce na Ideia ou no PRD da Spec da feature grande, nunca durante um Run, numa corretiva ou no dia do estouro, com o tempo do gate antes e depois, o que já foi cortado, juntado ou movido e o número novo medido; sem o dossiê, o trabalho precisa caber.
quando: 25/11/2026 09:15 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 6)
pergunta: "Leve por Task, completo uma vez por Spec": só no repositório do Roundfix ou capacidade do produto? Opções: só no Roundfix; produto já na 0421; produto em Spec própria.
resposta: "Só no Roundfix (Recomendado)"
alcance: sem mudar o produto; a TechSpec prova onde a suíte inteira roda uma vez por Spec; se não houver como sem chave nova, volta PEDIDO antes de qualquer Task.
quando: 25/11/2026 09:17 BRT
```

The four blocks of message `01a213bf-b859-7f60-91a1-201ca006c98a`, verbatim. They answer the requests `01a212a7-9fe6-7fa3-9b84-d251948359e0` and `01a212ca-e97f`:

```text
DECISÃO DE MARCIO
pedido: 01a212a7-9fe6-7fa3-9b84-d251948359e0 (bloco 1, autorização) e acréscimos (i), (ii) e (iii) de 01a212ca-e97f
pergunta: Autoriza a entrega da 0421: 9 arquivos protegidos com limite, regenerações pelos comandos próprios, retirada da ADR-0387 e de partes da ADR-0452 por ADR nova curta, CI em jobs paralelos no macOS, até o release patch; nada acima de 420 s, sem -timeout maior, sem runner pago?
resposta: "Autorizar (Recomendado)"
alcance: o bloco (a) a (d) do 01a212a7 como descrito, mais os acréscimos (i) registro de cobertura (saem os 2 testes do orçamento, entram os do código novo), (ii) a regeneração condicional do registro de superfícies e (iii) a nota de emenda na ADR-0252, AJUSTADOS pela decisão do bloco 4 abaixo. Fora: o que o 01a212a7 já punha fora.
quando: 25/11/2026 ~14:25 BRT

DECISÃO DE MARCIO
pedido: 01a212a7-9fe6 (bloco 2)
pergunta: Como a regra final é entregue? Opções: PR só de documentos; Spec curta 0422.
resposta: "PR só de documentos (Recomendado)"
alcance: PR de fechamento da própria 0421, sem Run, aberto quando as janelas fecharem, só com arquivos não protegidos, sem release.
quando: 25/11/2026 ~14:25 BRT

DECISÃO DE MARCIO
pedido: 01a212a7-9fe6 (bloco 3)
pergunta: O que acontece quando um gate passa do número? Opções: julgar pela janela e segurar o release; cobrar cada execução; só avisar.
resposta: "Julgar pela janela (Recomendado)"
alcance: nenhuma execução isolada falha por tempo; antes do release lê-se a janela, e um gate acima do número segura o release até o tempo voltar ou Marcio decidir com dossiê.
quando: 25/11/2026 ~14:25 BRT

DECISÃO DE MARCIO
pedido: 01a212ca-e97f (pergunta 4, contratos duas vezes no CI)
pergunta: A regra "repository contracts validate at the pull request boundary" muda no CI? Opções: mudar só no CI (uma vez por evento); manter; retirar.
resposta (literal): "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main."
alcance (leitura do Praetor; se discordar ou se faltar autorização, mande PEDIDO): o CI do repositório roda só no PR, como gate do squash merge; o push na main deixa de disparar a verificação. No PR, cada contrato roda uma vez. Fora do CI nada muda: make verify-docs segue obrigatório antes do PR. O "Full Contract Run da main" e as medições de push na main deixam de existir: a janela da regra final passa a ser PR e fechamento de Task (e release, se houver verificação nele). O release: proponha na TechSpec o que ele verifica sob essa regra (por exemplo, nada além do que o PR já provou no commit marcado, ou um conjunto mínimo de publicação); se isso mudar o alcance do release.yml além do já autorizado, mande PEDIDO. A ADR nova da regra final registra esta decisão.
quando: 25/11/2026 ~14:25 BRT
```

The author read these blocks in the copy the Supervisor session saved from the message, not in the inbox.

The answer on the release, verbatim. It answers the request `01a213d8-c4d8`, and it is the Praetor's application of the fourth block above, not a block of the maintainer:

```text
Praetor → Roundfix, 25/11/2026 ~15:12 BRT, resposta ao 01a213d8-c4d8. Decidido pelo Praetor aplicando a decisão de Marcio de 14:25 BRT (01a213bf-b859), sem nova pergunta: "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main." A opção 3 roda verificação depois do merge, contra a regra; a opção 2 publica uma tag em qualquer commit sem prova de PR. A opção 1 é a única que cumpre a regra sem abrir esse buraco.

Vale a opção 1, com o alcance exatamente como no seu pedido: no release.yml, só o job verify: tirar make verify, make verify-contracts e o setup do Go; pôr o passo de procedência pela API do GitHub com o token do workflow (a tag aponta para o commit que mudou a versão; é o squash merge de um PR na main; o branch do PR continha a ponta da main e tem a mesma árvore; o "Verification gate" desse PR passou; qualquer outra situação recusa, inclusive GitHub fora do ar); runner ubuntu-latest sem o comentário do override; permissões só de leitura (contents, checks, pull-requests). O job de publicação fica igual. Notas de emenda na ADR-0253 e na ADR-0452 pela ADR-0453; tirar dos não-objetivos do PRD o reaproveitamento do veredito. Runbook: atualizar o branch do PR de release antes do merge e disparo manual antes da tag. Nenhuma configuração do GitHub.

Registre no _authorization.md como: aplicação da DECISÃO de Marcio 01a213bf-b859 pelo Praetor, resposta 01a213d8-c4d8. Se Marcio discordar depois, vira emenda. O risco que você registrou (dois PRs verdes que se encontram numa árvore nunca verificada) fica escrito na ADR nova da regra final.
```

The author read this message in the copy the Supervisor session saved from it, not in the inbox.

Autonomous route: no person present; decisions taken from the maintainer's decisions quoted above, the approved Idea and the defaults it records, the measurements of `references/sources.md`, the accepted ADRs and the glossary.

Settled from those sources, without a question:

- The override of 2026-11-25 ends in the first delivery. Its own comment says so.
- What is withdrawn from ADR-0452 is the grouping in one job, the runner type and the 420 seconds. The question the maintainer answered offered "fica sem regra sobre agrupamento e sobre o tipo de runner".
- The retirement is recorded by a short ADR in the change that first depends on it, and the rule is another ADR: two numbers, as the block of 09:17 BRT says.
- A number may fall on a new measurement without a decision, and a change of machine goes to the maintainer. The Idea showed this default with its questions and no answer changed it.
- A missed target means the rule waits and the miss goes to the maintainer, as the same answer requires: "nunca acima da meta aprovada".
- The 360 s judge the Complete Suite Run at whichever gate carries it.
- The time of a gate in continuous integration is the longest of its jobs, without the wait for a runner. The approved target reads "for a run that did not wait for a runner" in the Idea.
- Which checks stay at a settlement follows from "Só o leve" and from the sixteen diagnostics the author read.
- The targets and the windows of a push to main leave with the gate. The answer on continuous integration says so in its scope: "as medições de push na main deixam de existir".
- ADR-0453 gives ADR-0253 its notice for the run on a push to main. The same scope ends that run, and the decision record that fixed it changes through the lifecycle.
- The Gate Time Rule records the decision on continuous integration, and ADR-0453 records it first, for the interval, because the first delivery is the change that depends on it.

Delegated Decisions. The first four are the first draft's; the maintainer's later answers kept three of them.

- delegated: 2026-11-25, prd
  - question: Does the work fit one Spec?
  - options: one Spec for the gates and the repeated work of the tests, then a second Spec for the Git processes of the product and the size of the CLI package; one Spec that holds both; two Specs authored now
  - adopted: one Spec now and a second after it — the second needs the per-test times this one starts recording on the macOS runner, and the two change different things
  - under: Project Config
  - answer: 2026-11-25, kept — "Final é a 0421, em 2 Specs (Recomendado)"

- delegated: 2026-11-25, prd
  - question: Which measurement judges a target on the maintainer's machine?
  - options: a lab protocol and the field record of the Run Database, each with its number; the lab protocol alone; the field record alone
  - adopted: both, each with its number — the lab number is the one the QA gate can check before delivery, and the field number is the one the maintainer waits for
  - under: Project Config
  - answer: 2026-11-25, kept — "Como proposto"

- delegated: 2026-11-25, prd
  - question: Over which runs is a target of continuous integration judged?
  - options: the twenty runs after the delivery, by their median and their p95; the delivering run alone; every run of a month
  - adopted: the twenty runs after the delivery — one run cannot show the spread between runners that ADR-0387 measured, and a month delays the answer past the next Specs
  - under: Project Config
  - answer: 2026-11-25, kept — "Como proposto"

- delegated: 2026-11-25, prd
  - question: Where do the measurements live?
  - options: in a reference file of the Spec now and in a reference document of the repository that a Task writes; in the PRD alone
  - adopted: in both files — the PRD stays readable, and the gates need a document that outlives the Spec to name their targets
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-25, prd
  - question: How large is the declared allowance of the Gate Time Rule?
  - options: 25 % over the measured result, never above the approved target; 10 %, as the published test budget the Idea cites; 22 %, as ADR-0387 used
  - adopted: 25 % — it is the example the maintainer saw with the question he answered, and a narrower one meets its first request for change at the first feature
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-25, prd
  - question: Are the twenty runs of a window counted for each event or for both together?
  - options: twenty passing runs of a pull request and twenty passing pushes to main; twenty runs of both events together
  - adopted: twenty of each — the approved targets are stated for each gate, a pull request runs fewer tests than a push, and the rule fixes one number for each
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-25, prd
  - question: Which secondary rows of the first draft's target table stand?
  - options: keep 150 s for a settlement that selects the regeneration contracts, 90 s before a pull request opens and 300 s for a test binary, and drop the number per job and the five minutes for a pull request without a Go change; keep all five; keep none
  - adopted: keep three and drop two — the number per job was the change of ADR-0387 the maintainer did not take, and the five minutes were never measured on the macOS runner
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-25, prd
  - question: Does anything fail a gate for its time before the Gate Time Rule exists?
  - options: nothing new, only the ten-minute limit per test binary and the job timeout that exist today; the old budget judging each job at 420 seconds; a failure at each approved target
  - adopted: nothing new — a budget for each job is the change of ADR-0387 the maintainer did not take, and a failure at a target would write a time rule before its measurement
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-25, techspec
  - question: Which packages does a change to a documentation tree select on a pull request?
  - options: none for `docs/history/`, `docs/adr/`, `docs/specs/` and `docs/user-guide/`, and every package for any other path outside a package; every package for every path outside a package, as today; none for any path under `docs/`
  - adopted: none for the four trees no code test may read — the repository's own rule keeps code tests out of them, the pull request of a release still runs every test, and the other trees have no such rule
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-25, techspec
  - question: What does a review Batch run?
  - options: the packages its change reaches, because it runs outside the Run's branch; the light checks, like a Task settlement
  - adopted: the packages its change reaches — the maintainer's answer "Só o leve" is about a Task settlement, and reading the Run Marker alone would let a queue owner started from a marked shell skip the Complete Suite Run
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-25, techspec
  - question: Where does the Full Contract Run happen once continuous integration runs only for a pull request?
  - options: on every pull request, each contract once; only on a pull request whose change selects every package, such as a release's; nowhere, with only the contracts a change selects
  - adopted: on every pull request — a relevance that misses an input is then found before the merge, in the change that causes it, and the repository's rule already ran the regeneration contracts on every pull request, so the two contract jobs keep the estimates they had
  - under: Project Config
  - answer: none yet

The maintainer's answer on continuous integration leaves one event. The entry above that counts twenty runs for each event therefore decides nothing now, and the reason of the entry on documentation trees was rewritten on 2026-11-25, because it rested on the run on a push to main.

## Open Questions

None. The four questions of the earlier drafts are answered: the authorization, how the last delivery is carried, what follows a gate over its number and the contracts rule. The Reserved Question of the TechSpec revision, what a release verifies now that continuous integration runs only for a pull request, was answered on 2026-11-25 by the Praetor's application of the maintainer's decision on continuous integration: the release runs no test and proves its Release Provenance, as Core Feature 6 states. `## Decisions` quotes every answer, and the measurements behind the release question stay in `references/sources.md`, section 9.

Settled by the TechSpec, in its `## System Architecture` and `## Decisions`:

- The Complete Suite Run lives at the Delivery Queue's repository gate. This repository's command tells a settlement by the Run Marker with the Run's branch and the candidate of a Spec by its Archive Record, without a new key.
- Five jobs and two shards, gathered by one job that keeps the required check's name. The workflow keeps one trigger, a pull request into main.
- Release Provenance is one step of the release workflow, tested through the harness that already runs that workflow's steps. It is not a command of the product, and the release commit is the commit that last changed the version file.
- The suite-time budget's targets and variable leave the Makefile.
- The Gate Time Rule amends ADR-0453.
- The checks of Core Feature 2 become Repository Contract Tests declared `//verify:always`, except the comparison of the coverage record, which stays a code test the Makefile names.
- The number of tests that run at once in a package stays 16; the delivering pull request measures it on the macOS runner.
- ADR-0259 rejected splitting the CLI package and copying fixture templates. This Spec does neither. The second Spec may need both, and its author asks then.
