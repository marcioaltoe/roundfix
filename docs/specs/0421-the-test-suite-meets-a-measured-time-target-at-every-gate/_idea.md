---
spec: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
created: 2026-11-25
status: promoted
leverage: compounding
---

# The test suite meets a measured time target at every gate

## Overview

Every verification gate of this repository gets a Gate Time Target: a number measured on the machine that runs the gate. A Task of a Run settles on the checks its change can fail, the whole suite runs once per Spec, and a pull request, a push to main and a release are verified in parallel jobs on the standard macOS runner. Each gate reports its time against its target. The Spec ends by writing one new rule from the measured result, and that rule changes only in extreme necessity.

A gate, in this Idea, is a point of the delivery that runs a verification and can stop the work: a Task settlement and the QA gate of a Run, the repository gate of the Delivery Queue, the check before a pull request opens, the required check of a pull request, a push to main and a release. "Gate Time Target" and "gate" are candidates for the glossary, which the PRD settles.

It is for the maintainer, for the Agents of his Runs, and for the other repositories whose Runs share the one verification slot of his Mac mini.

It is valuable because the cost is measured and large. In the week of 2026-11-18 to 2026-11-25 this repository's settlement checks held that slot for 15.3 hours. Main is red on the macOS runner, and the release verifies on Linux under an override that names this Spec as its end. The suite grew from 227 s to 391 s in seven weeks, under a time rule measured on another machine.

V1 is the first of two Specs. It changes what each gate runs, removes the work the tests repeat, verifies continuous integration in parallel jobs and writes the final rule. It does not make the suite itself much cheaper on one machine. That needs changes to the Git processes the product starts and to the CLI package, which wait as a Backlog Entry. The leverage is compounding: every later Run pays less, and the rule keeps the gain.

## Problem

The request is the maintainer's answer of 2026-11-25, 06:40 BRT, to whether the test timeout should rise to 30 minutes: "Essa verificação do roundfix de 10 minutos já seria um absurdo, ainda mais se rodar em todas as runs, 30 minutos é ainda mais absurdo. É primordial uma revisão dos testes e conjuntos de testes para que tenhamos performance decente." Every measurement below is in [`references/sources.md`](references/sources.md) with its command and its machine, unless `## Sources` names another origin.

Main is red, and no release verifies on macOS. The first push to main verified on the standard macOS runner failed on 2026-11-25: `go test` killed the test binary of the CLI package at its ten-minute limit. The one pull request verified there took 32.8 minutes. The release went back to Linux the same morning, under an override whose comment says it lasts "until Spec 0421 (test-suite performance) returns this job to macos-26".

Every Task waits about ten minutes for tests its change cannot fail. At each Task settlement the Daemon runs the repository check, which runs every package outside the Baseline set when one Go file changes. In the week measured it ran 99 times in 43 Runs of this repository, with a median of 583 s and 15.3 hours in all. The mini runs one Verification at a time for every repository, so this repository took 69 % of the machine's settlement-check time with 41 % of the executions. Since 2026-11-01 the check failed sixteen times in 125 Runs: nine on format or lint, seven on a record or a golden file and none on a behavior test, by Roundfix's own classification.

The suite grows back, and no gate says so. ADR-0259 measured 1,409 top-level tests in the CLI package and a suite of 227 s on 2026-10-08. Seven weeks later the package holds 2,102 tests and the suite takes 391 s on the same machine. The [suite audit of 2026-10-24](../../references/2026-10-24-test-suite-audit.md) set two targets, a pull request "at most 300 s" and the settlement check "at most 10 min". One month later the pull request took a median of 656 s on Linux and the settlement check a median of 607 s over its passing runs. No rule held either number, and no gate reported its time.

The time rules moved each time a gate blocked the work:

| Date | Rule | What changed | The trigger its text gives |
| --- | --- | --- | --- |
| 2026-10-22 | ADR-0311 | The Delivery Queue's wait for the checks rose from 30 to 60 minutes | "needed 30 to 33 minutes" |
| 2026-11-10 | ADR-0387 | The suite budget rose from 240 s to 420 s | "26 of them exceeded 240 s" |
| 2026-11-24 | ADR-0452 | The verification moved from Linux to macOS, with the budget kept "until the suite is measured on the macOS runner" | "the container gate had repeatedly blocked delivery" |
| 2026-11-25 | Override on main | The release's verification returned to Linux | The suite passes the ten-minute limit on the macOS runner |

The maintainer named the pattern on 2026-11-25: "o que aconteceu até agora é que toda vez que o tempo é estourado, a regra é modificada seja de tempo ou qualquer outra coisa. Então dessa vez quero inverter e no final criar uma regra para que só seja modificar em extrema necessidade por testes de grandes novas features." The record also holds the opposite habit. ADR-0257, ADR-0259, ADR-0271 and ADR-0326 each cut work and moved no number, and the repository already holds two ceilings that only fall, of Sequential Tests and of the Real-Time Budget.

One finding of this stage changes the drafted plan. The Task settlement, the QA gate's precondition and the Delivery Queue's repository gate all run one command, read from one key: `defaults.verification`, which this repository sets to `make verify-changed`. The drafted PRD says the queue's gate runs `make verify`. A light check at each Task and a complete run once per Spec therefore need more than leaving the queue's gate as it is. The maintainer answered that the complete run stays a matter of this repository, and the TechSpec proves where it lives.

### Market Data

- A build has a published upper limit of about ten minutes. DORA: "The tests should not take more than a few minutes to run, with an upper limit of about 10 minutes according to DORA’s research (PDF)." The approved targets of continuous integration, a median of 8 minutes and a p95 of 10, sit at that limit and not below it. CircleCI's 2025 report gives a median workflow of 2 min 43 s, as the Research Subagent read it.
- Light checks first and a full run later is the studied practice in very large repositories. Machalica et al. (2019) report a selection that "reduces the total infrastructure cost of testing code changes by a factor of two, while guaranteeing that over 95% of individual test failures and over 99.9% of faulty changes are still reported". Memon et al. (2017): "very few of our tests ever fail, but those that do are generally 'closer' to the code they test". Both study human changes, not an Agent that verifies after each small Task. Applying them here is inference, and the sixteen local failures are the stronger evidence.
- Selection by imports in Go keeps a full run. The go-testsplitter documentation: "Dependencies not visible in the import graph (e.g. files read by relative paths outside the package, external services) are not detected. Running all tests on the main branch is recommended." This Idea keeps every test on every push to main.
- A budget that holds has a declared allowance and a stated exit. thoughtbot (2026): "The total suite budget is set 10% above your current total, giving you headroom to keep adding tests without immediately blowing the budget", and "Instead of bumping the budget to accommodate a slow test, you allowlist it with a reason." qntm on ratchets: "Sometimes, due to extenuating circumstances, we have had to manually raise the count again. This is something we try to avoid, though." Google's example error budget policy agrees beforehand what happens when the budget is spent.
- A hosted runner is noisy and limited. ADR-0387 measured one tree at 215 s on one run and 311 s on the next. GitHub runs at most 5 macOS jobs at once on the account's plan, and "During periods of high demand, jobs that run on macOS runners may take longer to start."
- The sweep found no published rule that joins a hard number, a warning line, a declared allowance, one owner and a new measurement when the machine changes. The final rule would join known parts.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | Each gate has a Gate Time Target | Critical | A target names its gate, its machine, its number and the measurement that judges it. The targets are the ones the maintainer approved on 2026-11-25, in the table below. |
| F2  | A Task settles on the light checks | Critical | At a Task settlement the repository check formats, vets and builds, runs the Repository Contract Tests that Contract Relevance selects, and runs the tests that compare a record or a golden file. It runs no other package test. The Task's own declared Verification still runs the tests the Task names. |
| F3  | The whole suite runs once per Spec, in this repository only | Critical | The tests of every package run once per Spec before its pull request opens. The product gains no configuration key, the Daemon does not change, and the gates of other repositories stay as they are. |
| F4  | Pull requests, pushes to main and releases are verified in parallel jobs | Critical | The verification runs in several jobs side by side on the standard macOS runner. A pull request runs the packages a change touches and the packages that import them. A push to main and a release run every test and every Repository Contract Test. The required check keeps its name, main turns green, and the release's verification returns from Linux to macOS. |
| F5  | Two time rules stop binding before the work that contradicts them | Critical | The sentences of ADR-0452 about one job, the runner type and the 420 seconds are withdrawn, and ADR-0387 is retired. Until the final rule exists three limits hold: no number above 420 seconds, no longer `-timeout`, no paid runner. |
| F6  | The final rule is written last, from the measured result | Critical | The Spec's last delivery is one new ADR with a number of its own. For each gate it fixes the measured result plus a declared allowance, never above the approved target. It is written once the approved targets are measured in production. |
| F7  | The final rule changes only in extreme necessity | Critical | A number rises, or the shape of the verification changes, only by the maintainer's decision on a measured dossier, asked for in the Idea or the PRD of the Spec of a large new feature. Without the dossier the work must fit. |
| F8  | Each gate reports its time | High | Every gate prints its time against its target when it ends. Between the retirement of the two rules and the final rule, this report is what shows a gate that got slower. |
| F9  | The tests stop repeating whole-module work | High | A code test no longer runs the tests of another package, builds a binary in each test, extracts the same corpus in each test or waits for a time to pass. |
| F10 | No test is lost | High | No assertion, test case or platform is removed to make a gate faster, and the coverage record shows every test that moved or was renamed. |

### The approved Gate Time Targets

| Gate | Machine | Target | Judged by |
| --- | --- | --- | --- |
| Run: the repository check at each Task settlement | The mini | Field: median 120 s and p95 240 s. Lab: 45 s | The first 30 passing settlements after the delivery, read from the Run Database; the lab protocol |
| Delivery Queue's repository gate | The mini | Lab: 360 s | The lab protocol |
| Pull request, push to main and release | The standard macOS runner | Median 8 minutes and p95 10 minutes, for a run that did not wait for a runner | The 20 runs after the delivery |

At the current volume the two production windows take about three days.

### The measured dossier

The request to raise a number of the final rule, or to change the shape of the verification, carries four things:

1. the large new feature;
2. the time of the gate before and after it, on the machine the rule names;
3. what was already cut, merged or moved to a less frequent gate;
4. the new number, which is the measured one.

The request is made only in the Idea or the PRD of the feature's Spec, as a Rule Challenge of the final rule. It is never made during a Run, in a corrective Spec or on the day a gate ran over.

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Settlement check in the field | Median ≤ 120 s and p95 ≤ 240 s, against 607 s and 882 s over the passing runs today | The Run Database, over the first 30 passing settlements after the delivery |
| Verification slot held by this repository's settlement checks | ≤ 3.5 hours a week at the same volume, against 15.3 hours today | The same export of the Run Database, summed over 7 days |
| Required check of a pull request on the macOS runner | Median ≤ 8 minutes and p95 ≤ 10 minutes, against 32.8 minutes on the one run measured | The step times of the 20 runs after the delivery |
| Main and release on the macOS runner | 20 of 20 pushes pass, and the release verifies on `macos-26` with the override removed | The workflow runs and the release workflow file |
| Late detection | At most 1 complete run in 30 Specs fails on a behavior test the old settlement check would have run | Counted at each such failure, in the reference document the PRD plans for the targets |
| Rule stability | 0 rises of a time number and 0 changes of the verification's shape in 90 days without a recorded decision of the maintainer, against 4 changes in the last 5 weeks | The Git history of the final rule, of the `Makefile` and of the two workflows |

## Feature Assessment

| Criteria        | Question                                            | Score |
| --------------- | --------------------------------------------------- | ----- |
| Impact          | How much more valuable does this make the product?  | Must do — main is red on the macOS runner, no release verifies there, and this repository's settlement checks hold 15.3 hours a week of the only verification slot. |
| Reach           | What % of users would this affect?                  | Strong — every Run, pull request, push and release of this repository. Other repositories gain only through the freed slot. |
| Frequency       | How often would users encounter this value?         | Must do — 99 settlement checks a week and about 12 runs of continuous integration a day. |
| Differentiation | Does this set us apart or just match competitors?   | Maybe — an internal gain that reaches the published limit of ten minutes. The rule that holds the number is the unusual part. |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong — the gain of the last speed-up was gone in seven weeks. A held number and a reported time are what keep this one. |
| Feasibility     | Can we actually build this?                         | Strong — the light checks take 17 s in the lab today, and the arithmetic of the jobs comes from that runner. Three things are open: the macOS runner was measured twice, one configuration key feeds the settlement and the queue's gate, and the account runs five macOS jobs at once. |

The leverage type is compounding. The verdict is to proceed in the shape the maintainer's answers describe. Two results would change it. If the TechSpec shows that the complete run has no home without a product change, the Spec grows or waits on a new decision. If the first runs on the macOS runner miss the targets by a wide margin, the result is not reached and the final rule waits.

## Trade-off Insights

- **Recommended approach:** the shape the maintainer chose. Two Specs. In this one a Task settles on the light checks, the whole suite runs once per Spec in this repository only, continuous integration runs in parallel jobs, two rules are retired at the start, and the final rule comes last with the measured result, a declared allowance and a change only on a dossier.
- **Key trade-offs:**
  - A light settlement against early detection. For: sixteen failures in 125 Runs, all format, lint, record or golden. Against: a break in a behavior test the Task did not name is found once per Spec, after the Tasks that could have fixed it cheaply.
  - The rule at the end against a rule now. For: every time number written before its measurement was later moved. Against: the repository works without a written time rule until the windows close, and the suite grew 72 % in seven weeks under one.
  - A fixed number against the noise of a hosted runner. One tree took 215 s and 311 s on two runs. A number judged on one run fails by chance, and a number judged over a window answers late. The approved targets are judged over windows.
  - Composition against the cost of the suite itself. This Spec removes about 6 % of the test seconds and gains the rest by deciding what runs where. That can hide further growth of the suite.
  - This repository only against a product capability. A settlement command apart from the complete one would serve every adopter, and no adopter asked for it. It would also enlarge the Spec that must turn main green.
  - An allowance against none. Without one, the first feature that adds a test asks for a change of the rule. With the whole approved target as the allowance, the suite grows back unseen.
  - Selection by imports on a pull request buys little here. The CLI package imports 40 of the 49 other packages, so most changes select it. The pull request meets its target through parallel jobs.
- **Risks identified:**
  - A break reaches the end of the Spec. ADR-0342 records two Runs of an adopter that "reached Clean" with a broken tree when nothing ran at settlement: "Both were caught only by the branch's own gate before the Pull Request." This Idea keeps that same net and no earlier one for behavior tests.
  - The classification of the sixteen failures is Roundfix's own, and the author did not read the sixteen diagnostics.
  - The complete run may have no home without a product change. One key feeds the settlement and the queue's gate, and the Verification receives nothing that names the gate it runs in.
  - The macOS runner was measured twice: one pull request and one failed push. The delivering pull request is the first real measurement of the new jobs.
  - macOS jobs queue. The account runs five at once across all its repositories, so a pull request and a push that overlap wait for each other. The targets are for a run that did not wait.
  - The lab is not the field. The same command takes 251 s in the lab and a median of 607 s in the field.
  - For some days the repository has no written time rule. The three limits of 2026-11-25 and the reported times are all that hold.
  - A strict change clause can stop a feature. The maintainer said on 2026-11-20 that "uma regra não deve impedir o progresso/evolução do projeto". The dossier is the exit the clause names for itself, so it must be possible to assemble one.
- **Dissenting view:** the gain is composition, and composition is what eroded before. At least 153,850 Git processes start in one run of the suite, about three quarters of them from the product, and this Spec leaves them. The complete run on the mini stays near six minutes. A second dissent stands against the light settlement: the evidence that behavior tests never fail there is a classification nobody checked, and one late break costs a corrective Task after the QA gate.
- **Stretch goal (V2+):** the second Spec, which lowers the cost of the suite itself. After it, gate time as a fact of the product: the Run Database records the time of each gate and Roundfix shows it against a target, for every repository.

## Opportunity Scan

The original, as the first answers of 2026-11-25 shaped it, has a clear ceiling. The settlement check falls by about 80 %, continuous integration fits ten minutes, and main and the release are green on macOS. The cost of the suite itself barely moves.

- **More ambitious: one Spec that also changes the product.** It batches the Git processes the commands start and splits the CLI package, for a complete run of 240 s on the mini. It removes the cost instead of rearranging it. Scale: massive. Impact Must do, Reach Strong, Frequency Must do, Differentiation Maybe, Defensibility Strong, Feasibility Maybe: about twice the size, with main red for longer. The maintainer saw it twice and did not take it.
- **Simpler: continuous integration only.** Parallel jobs and the release back on macOS, and nothing in the gates of a Run. Main is green in one or two Tasks. Scale: small. Impact Strong, Reach Maybe, Frequency Strong, Differentiation Pass, Defensibility Maybe, Feasibility Must do. The 15.3 hours a week stay. The maintainer saw its core, targets that bind only in continuous integration, as an option of the targets question and did not take it.
- **Adjacent: gate time as a product fact.** The Run Database records each gate's time, Roundfix shows it against a target, and a repository can name a settlement command apart from its complete one. Every adopter gets the measurement and the lighter settlement. Scale: medium. Impact Strong, Reach Strong, Frequency Must do, Differentiation Strong, Defensibility Strong, Feasibility Maybe. The maintainer chose "Só no Roundfix (Recomendado)".

**Chosen direction:** original, split in two Specs — the maintainer's answer "Final é a 0421, em 2 Specs (Recomendado)". It is the only shape that turns main green soon and frees the verification slot. The more ambitious alternative becomes the second Spec. The adjacent one is the stretch goal, and no Backlog Entry records it, because the maintainer did not take the option that would have minted one.

## Sub-Features

The request becomes two Specs.

1. **Spec 0421, this one.** What each gate runs, the work the tests repeat, continuous integration in parallel jobs, the retirement of the two rules and the final rule. It holds one context: the gates of this repository and the tests they run. Its last delivery, the final rule, can only be written after its own first delivery is measured in production. How the Spec carries that is an open question for the PRD.
2. **The second Spec, not yet minted.** The Git processes the product starts for each command and the split of the CLI package. It waits as the Backlog Entry [`2026-11-25-the-product-starts-most-of-the-suites-git-processes-and-the-cli-package-is-one-test-binary.md`](../../backlog/2026-11-25-the-product-starts-most-of-the-suites-git-processes-and-the-cli-package-is-one-test-binary.md). It runs after this Spec and, under the final rule, may only lower numbers. ADR-0259 rejected two of its directions, so it asks its own Rule Challenges.

## Integration with Existing Features

| Integration point | What changes or how it is used |
| --- | --- |
| `defaults.verification` in Project Config | Unchanged. One key feeds the Task settlement, the QA gate's precondition, the Delivery Queue's repository gate and review Batches. F2 and F3 must hold with it as it is. |
| The Daemon and the Settlement Check | Unchanged. Every Task still settles on the configured command, as ADR-0342 requires. What that command runs in this repository changes. |
| The Delivery Queue | Unchanged. It waits for the same required check, and its repository gate takes the one verification slot as today. |
| The workflow that verifies pull requests and pushes, and the release workflow | F4. The publication job of the release keeps its runner, its permissions and its steps. |
| The override of Spec 0420 in the release workflow | Ends with F4, as its own comment says. |
| The Run Database | Read by a query to judge the field target of the settlement. Nothing is added to it. |
| The Rule Challenge | F7 uses it as it is: the request to raise a number is a Rule Challenge of the final rule, asked in the Idea or the PRD of the feature's Spec. A Run, a Task and the QA gate already ask none. |
| The coverage record | Re-recorded by its own command when a test moves or is renamed (F10). |
| The glossary | Gains the names the PRD settles, starting from Gate Time Target. |

## Rules That Touch the Idea

- ADR-0452. Its sentences about one job, the runner type and the 420 seconds are retired by the maintainer's decision. Its other decisions keep binding: the end of the Linux Lane, the deprecated keys, the Doctor Command's checks, the External Test Tool rule and the compile contract for Linux and Windows.
- ADR-0387. Retired by the maintainer's answer to this stage's Rule Challenge.
- The rule "repository contracts validate at the pull request boundary", in `docs/agents/specific-repository.md`, with the consequence ADR-0252 records. It keeps binding. It does not block this Idea, because a repeated contract can run in another job at the same time. The reserve question about it is in `## Open Questions`.
- ADR-0253. Kept: every Repository Contract Test runs on every push to main and before every release.
- ADR-0257, ADR-0182 and ADR-0342. The maintainer's answer "Só o leve" changes two of their sentences: after it no gate inside a Run is a full one, against the title of ADR-0257, and ADR-0182's measured cost of a settlement, "206 to 223 seconds", no longer describes this repository. ADR-0342 stays true as written. No Rule Challenge is owed, because he decided the change; the PRD names that answer where it records each sentence.
- ADR-0259. Its rule that a ceiling only falls is the in-house precedent for the final rule. Its two rejections, of a split of the CLI package and of fixture templates, block the second Spec and not this one.
- ADR-0271 and ADR-0326. Kept: no wait on real time is added, and a test that reads the live documentation trees stays a Repository Contract Test.
- The Baseline guide `docs/agents/agent-instructions.md` states "The selected repository Verification is `make verify`" and "The selected incremental Verification is `make verify-changed`". A third command, or a new meaning for one of the two, touches those Baseline decisions.
- The ADR lifecycle in `docs/agents/docs-layout.md` has no retirement of one clause, and an amendment needs a later accepted ADR in the same change. This is why default (d) of `## Open Questions` counts two ADR numbers.
- The glossary's definition of Rule Challenge lists "exception request" under what to avoid. The change clause of the final rule therefore speaks of a Rule Challenge.

None of these carries a Kept Line.

## What the PRD Rewrite Corrects

The drafted PRD stays as it is until the PRD phase rewrites it. This stage found seven things it must change:

1. The baseline and the target of the Delivery Queue's repository gate. That gate runs `make verify-changed`, not `make verify`, so the 376 s of the lab measure a command no local gate runs.
2. Its Core Feature 3, which keeps the complete run at the queue's gate "as today", and its non-goal of a new configuration key. The maintainer's answer to question 6 settles which gives.
3. ADR-0453 inside Task 4. The final rule is now the last delivery, and a short ADR that records the retirement may come first.
4. Its three Rule Challenges. ADR-0452 and ADR-0387 are answered `retire`, and the contracts rule is no longer asked at this stage.
5. "At most four jobs". The number of jobs is free after the retirement, and the published limit is five macOS jobs at once for the account.
6. Its first risk, which should cite the two incidents of ADR-0342 and ADR-0257's "failed 6 times in 125".
7. Its Active ADR obligations. This Idea cites ADR-0311 for the history of the time rules, and `roundfix spec check` reports `SC-ADR-UNLISTED` until the PRD lists that ADR, as kept. The nine `SC-ADR-RELATED` findings of the draft are the PRD's as well.

## Out of Scope (V1)

- **Batching the Git processes of the product, splitting the CLI package and copying repository fixtures from a template** — the second Spec. The maintainer chose two Specs, and that work changes the product and the package layout, where this one changes the gates and the tests.
- **A product capability for two verifications, or gate time recorded by the product** — the maintainer chose "Só no Roundfix (Recomendado)". The measured problem is this repository's, and no adopter asked for a lighter settlement.
- **A new key of the Project Config or the User Config, and any change to the Daemon** — the scope of "Só o leve" and of the answer to question 6. If the complete run cannot exist without one, a decision request returns to the maintainer before any Task.
- **A longer `-timeout`, a number above 420 seconds and a paid, larger or self-hosted runner** — the maintainer refused them on 2026-11-25, and the refusal holds until the final rule.
- **Removing an assertion, a test case or a platform to gain time, and skipping a test on the macOS runner** — the request is for a review of the tests, and a gate that proves less is not a faster gate.
- **A change to the contracts rule at this stage** — it does not block the goal. The TechSpec author asks the reserve question only if the measured layout of the jobs needs a contract to run once in a job.
- **The rest of ADR-0452** — only the three sentences are retired.
- **The gates of every other repository that uses Roundfix** — their settlement check costs a median of 39 s, and their gates are theirs.
- **The Baseline rule that a rule must not block the project** — it changes the documentation Baseline of every repository and stays its own Backlog Entry. This Idea only takes the maintainer's words of 2026-11-20 from it.
- **What happens to main and to the releases while this Spec is not delivered** — decided apart, by the override of 2026-11-25.

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

## Decisions

The maintainer's decisions come first, all of 2026-11-25. Each reached this project as a `DECISÃO DE MARCIO` block that the Praetor, the session that coordinates the software projects, sent through ai-memory from `default/praetor`. The text in quotes is verbatim. The times are the ones the blocks and the stage's brief state, and `## Sources` records a limit about two of them. The author's own choices close the section.

Before the stage:

- 06:40 BRT, message `01a2121b-9124-7571-8f8f-de0f40876364`: the request quoted in `## Problem`. Its scope forbids a 30-minute timeout and a larger suite budget and asks for this Spec. [`references/sources.md`](references/sources.md), section 1, holds the block.
- 06:58 BRT, message `01a21237`: the release verifies on Linux until this Spec returns it to macOS. The author read it in the comment of the release workflow on main, commit `9f85cc5e`, not in the message.

The answers to the first four Reserved Questions of the drafted PRD, 07:58 BRT, message `01a21256`. The author did not read this message; the stage's brief relays the answers and gives the id in this short form.

1. Does the Idea stage run? — "Rodar a Ideia". The drafted PRD is rewritten after this Idea.
2. Are the Gate Time Targets the right ones? — "Como proposto". The table of `## Core Features` holds them. The option not taken put the second Spec inside this one.
3. What does a Task settlement run? — "Só o leve". Format, vet, build, the selected contracts and the record and golden tests at each Task; the package tests once per Spec; on a pull request, the changed packages and the packages that import them. The scope as relayed: "nenhuma chave nova de configuração e nenhuma mudança no Daemon".
4. The Rule Challenge of ADR-0452 — the answer quoted in `## Rule challenges`.

The answers to this stage's questions, message `01a2127e-c61d-7a93-b702-ccbc6ebd2d45`, which answers the request `01a2126b-9f73-7791-983f-0b6a4cf191b1`. Question 1 of the stage, whether it runs, repeated the answer of 07:58 BRT. The five blocks, verbatim:

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

What each answer settles:

- Question 2: the final rule is the last delivery of this Spec, and the second part waits as a Backlog Entry (F6, `## Sub-Features`).
- Question 3: ADR-0387 is retired, and the final rule is a new ADR with a number of its own (F5, F6).
- Question 4: the final rule fixes the measured result plus an allowance the PRD declares, never above the approved target (F6).
- Question 5: the change clause and the dossier (F7).
- Question 6: the complete run stays a matter of this repository, and the TechSpec proves where it lives (F3).

The author's own, made in this stage:

- No source is adopted. No open Backlog Entry and no unresolved Finding shares this context; the entry about the rule that must not block the project changes the Baseline of every repository and stays apart.
- The second Spec's Backlog Entry is of type `perf` and carries a Waiting Declaration on the delivery of this Spec.
- The contracts rule is not asked about at this stage. After the retirement of the one-job rule it costs extra work and blocks nothing, and extra work is not a block.
- The measurements of the PRD phase are reused and none was repeated.
- No ADR is written at this stage, and no number is reserved for the two this Spec needs.
- The Praetor's reading that the contracts rule should be retired and rewritten with the final rule is recorded in the reserve question below. It is a reading, not a decision of the maintainer.

## Open Questions

Seven defaults were proposed with the stage's questions and no answer changed them. The record of the answers says the reserve question "was not asked: it goes to the TechSpec, with defaults (a)–(g) as proposed". Each default applies until someone answers otherwise.

- (a) Does the override of Spec 0420 end in this Spec? Default: yes. The release's verification returns to `macos-26` in the delivery that brings the parallel jobs. Before any tag, the manual dispatch of the release workflow, which publishes nothing, proves it.
- (b) Which tests stay at a Task settlement? Default: "Só o leve" decides it, and the PRD says which record and golden tests stay. The sixteen failure diagnostics were not read again, and they decide the list.
- (c) What exactly is retired from ADR-0452? Default: the sentences about the grouping in one job, the runner type and the 420 seconds. The question he answered offered "fica sem regra sobre agrupamento e sobre o tipo de runner". The rest of the ADR keeps binding.
- (d) How is the retirement recorded? Default: by a short ADR in the change that first depends on it, which the final rule then replaces. That is two ADR numbers, and the final rule's number is new.
- (e) May a number fall without a decision? Default: yes, on a new measurement. Only a rise and a change of shape need the maintainer's decision. A change of machine, such as a new runner image or a new mini, measures the same suite again and goes to the maintainer; it is not a relaxation.
- (f) Does the contracts rule change? Default: no. With parallel jobs the repeated contract can run in another job at the same time and cost no wall time. It blocks only if the measured layout must put both runs in one job.
- (g) What does a Run delivered outside the Delivery Queue get? Default: the complete run from wherever the TechSpec puts it. If that is the queue's gate, such a Run has no complete run before its pull request, and the PRD says so.

The reserve question, deferred to the TechSpec. Its author asks it only if the measured layout of the jobs needs a contract to run once in a job. Until then the rule binds.

- Rule: "repository contracts validate at the pull request boundary", in `docs/agents/specific-repository.md`, with the consequence ADR-0252 records: "CI runs the selected contracts twice. That costs a few seconds on a typical change."
- It imposes: `make verify-docs` runs every contract it lists, in every job, whether or not the job already ran them, and it must pass before any pull request opens.
- It prevents: running each contract once in a job. On the macOS runner the repetition costs 214 s to 372 s when it stays in one job, and the account runs at most five macOS jobs at once.
- The author's reading: the rule holds for the checkout before a pull request, where it is the only run. In continuous integration it is out of date, because the repetition stopped costing "a few seconds".
- Change (Recommended): in continuous integration each contract runs once per event, in any job. Outside it nothing changes, and `make verify-docs` stays mandatory before a pull request opens. The final rule records the form.
- Keep: the repetition stays, in separate jobs. The rule gains its Kept Line, and the layout uses five jobs, the account's limit, with more waiting when a pull request and a push overlap.
- Retire: the Praetor's reading, the same treatment as ADR-0452. Until the final rule, no rule would require `make verify-docs` before a pull request.

Questions the PRD or the TechSpec must settle:

- Where does the complete run live? Default: the maintainer's answer to question 6. The TechSpec proves it: at the Delivery Queue's repository gate, if this repository's own command can tell a Spec's delivery from a Task settlement, or at the Spec's QA gate, where a failure returns to the Run's Agent instead of parking the item. If neither works without a new key, a decision request returns before any Task.
- What does the 360 s target judge? The queue's gate does not run the whole suite today. Default: the target judges the complete run of the suite on the mini, at whichever gate carries it, and the PRD restates the row and its baseline.
- How large is the declared allowance? Default: 25 %, the example shown with question 4, which the PRD confirms. No number passes the approved target.
- How does the Spec carry its last delivery? The final rule needs about three days of production measurements that exist only after the first delivery merges. Default: the PRD plans the final rule as a short delivery of its own and names the mechanism. If that needs a new authorization or a Spec number of its own, it returns to the maintainer as a decision request.
- What if a target is missed in its window? Default: the result is not reached and the final rule waits. The miss goes to the maintainer with its measurement, and no number is written above an approved target.
- Does anything fail a gate for its time before the final rule exists? Default: no. Each gate reports its time against its target, the Go command's ten-minute limit per test binary stays the only hard stop, and nothing in the tree states a time number above 420 seconds. The PRD confirms it or keeps a stop.
- Which secondary rows of the proposed target table stand? The drafted table also holds 150 s for a settlement that touches a regeneration contract, 90 s before a pull request opens, 5 minutes for a pull request without Go changes, a number per job and 300 s per test binary. Default: the PRD restates them, because the answers to questions 3 and 4 changed what the last two rest on.
- How many jobs verify one event? Default: the TechSpec decides from the measurement, within the five macOS jobs the account runs at once.
- Is the work authorized? Not by this stage. The authorization record is `proposed` and grants nothing, and the seventh Reserved Question of the drafted PRD has no answer. The PRD phase asks it again for the paths this shape changes.
- Which names does the glossary take? Default: the PRD's Glossary settles them, starting from Gate Time Target and gate.

## Sources

The source of the request is the maintainer's decision of 2026-11-25, 06:40 BRT, copied whole in [`references/sources.md`](references/sources.md), section 1. Sections 2 to 5 of that file hold the measurements of the PRD phase, which this stage reused and did not repeat. This stage added one reading of the same export of the Run Database: the 99 executions, the median of 583 s and the 15.3 hours count passing and failing settlement checks, where the PRD phase's 83 executions, 607 s and 14.4 hours count the passing ones. The other repositories' 140 executions, 39 s and 7.0 hours group Runs by the command they ran.

This project's ai-memory messages were listed on 2026-11-25, when the stage's research began and again when this file was written, and none was pending. The stage's questions went to the maintainer through the Praetor, and he answered them; the author answered none.

The Research Subagent ran the codebase sweep and the web sweep. The author ran the Secondbrain searches and the Run Database arithmetic, which need a command, and read again the passages this file quotes from the repository and from the Secondbrain.

Secondbrain, read on 2026-11-25 (`wiki/index.md`, then four searches with the Indexus tool `pesquisar_acervos`):

- `wiki/concepts/flaky-tests-e-prazos-de-tempo-real.md` — "A correção habitual — aumentar o prazo — é a que menos fica", from this project's own Specs 0171, 0213 and 0265, and a ceiling "que só cai (nunca sobe)". It supports the maintainer's decision to invert the habit, and it is why the final rule is a number that only falls without a decision.
- `wiki/concepts/slo-e-orcamento-de-erro.md` — "Acordar antecipadamente o que muda quando o orçamento é consumido". It gave question 5 its form: the rule says beforehand what happens when a gate passes its number.
- `wiki/concepts/medicao-funcional-e-orcamento-de-consultas.md` — "medir o caminho corrigido antes de aumentar o orçamento de consultas", from another product of the maintainer. The same order of steps as the dossier.
- `wiki/concepts/testes-paralelos-em-go.md` and `wiki/entities/roundfix.md` were returned and had informed the PRD phase; this stage did not open them again.
- Two searches, on test selection and on ratchets, returned nothing relevant. The Secondbrain holds nothing on splitting tests across jobs, on selection by change or on the governance of a budget.

External sources, found through Exa on 2026-11-25. The author opened again the ones marked "reopened"; the others are the Research Subagent's reading of excerpts.

- DORA, "Capabilities: Continuous integration", https://dora.dev/capabilities/continuous-integration/ — the ten-minute limit. Read by the PRD phase.
- Memon et al., "Taming Google-Scale Continuous Testing", 2017, https://research.google/pubs/taming-google-scale-continuous-testing/ (abstract reopened) and Machalica et al., "Predictive Test Selection", 2019, https://arxiv.org/abs/1810.05286 (abstract reopened) — light checks first with a later full run. They support F2 and limit it: neither studies an Agent's loop.
- go-testsplitter, https://github.com/takuo/go-testsplitter/blob/main/README.md — the blind spot of selection by imports. It is why every test stays on main.
- gotestsum, https://github.com/gotestyourself/gotestsum/pull/258 (reopened), and Tailscale's test helper, https://github.com/tailscale/tailscale/blob/v1.90.8/tstest/tstest.go (reopened) — how Go projects spread tests over jobs, and that one large package bounds the gain. They inform the TechSpec, not this Idea.
- qntm, "Ratchets in software development", https://qntm.org/ratchet (reopened) — a count that only tightens, raised by hand as the exception. It shaped default (e).
- Google SRE Workbook, example error budget policy, https://sre.google/workbook/error-budget-policy/ (reopened) — a written policy with an owner and consequences agreed beforehand. It shaped question 5.
- thoughtbot, "Introducing Test Budget", 2026-03-17, https://thoughtbot.com/blog/introducing-test-budget (reopened) — a declared allowance and a reasoned exception in place of a raised budget. It shaped question 4 and challenged its number: 10 % there, 25 % proposed here.
- GitHub Docs, "Actions limits", https://docs.github.com/en/actions/reference/limits (reopened) — five concurrent macOS jobs. It is the queueing risk.
- CircleCI's 2025 report and Shopify's target of a p95 under ten minutes, as the Research Subagent read them — the approved targets sit at the published upper limit.

The digest of this research went to the project that owns the Secondbrain as ai-memory message `01a21268-3236-7c31-90f0-cbbe19afa660`, after the advisory duplicate check. It names `wiki/concepts/slo-e-orcamento-de-erro.md` as the one strong match that opened, and the PRD phase's digest `01a2124b-f1b9-78b0-bfd1-02abf74caf88` as the one it complements.

Limits of this stage:

- No log of the queue's repository gate was found on the mini, so the finding about what that gate runs rests on the code and the configuration.
- The author did not read the messages `01a21256` and `01a21237`.
- Two stated decision times run ahead of the messages that carry them. The message ids are ordered by time and place `01a21256` at about 07:53 BRT, against a stated 07:58, and `01a2127e-c61d-7a93-b702-ccbc6ebd2d45` at 08:37 BRT, against a stated 09:15 and 09:17. The record of the answers was written on the mini at 08:37 BRT. The date and the order of the decisions are not in doubt. The author changed no time and reports the difference to the session that coordinates this Spec.
- The Research Subagent's list of the ADRs that depend on ADR-0452 and its reading of the earlier Specs' Archive Records were not checked again.
