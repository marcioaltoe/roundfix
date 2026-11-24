---
spec: 0420-roundfix-stops-verifying-on-linux
status: active
created: 2026-11-24
surfaces: [backend, cli, infra, docs]
type: refactor
---

# Roundfix stops verifying on Linux

The maintainer runs Roundfix only on macOS, and Roundfix verifies on Linux in three places: the Delivery Queue's optional Linux Lane, which tests a candidate's changed Go packages in a Docker container before publication, the job that verifies every pull request and every push to main, and the release. On 2026-11-24 the Linux Lane parked the candidate of Spec 0418 twice with `compile: signal: killed`, on a host whose Docker virtual machine has 7.7 GiB of memory, for a platform nobody runs the product on. The maintainer decided to remove Linux verification for good. After this Spec the Delivery Queue has no Linux Lane, `roundfix deliver lane` and the lane's check in `roundfix doctor` are gone, a configuration that still names the lane's two keys keeps loading, and this repository's pull requests, pushes to main and releases are verified on macOS. Linux stays a platform Roundfix is compiled and published for, and no gate runs a test on it.

## Sources

The sources are two decisions of the maintainer of 2026-11-24, both relayed by the operator's brief and quoted in `## Decisions`: ai-memory message `01a20d73-109f-7bb2-a2b6-718854f6b25d`, which removes Linux verification, and ai-memory message `01a20f32-07af-7b11-a675-a2cbc7a64743`, which answers the author's question about eight Governed Paths. Both are copied whole in [`references/sources.md`](references/sources.md), sections 1 and 2. Sections 3 and 4 of that file hold what triggered the first decision and what the author measured, section 5 the published passages this PRD quotes, and section 6 the inventory of the lane in the code and the documents.

No open Backlog Entry or unresolved Finding shares this context. The seven open Backlog Entries concern a team's status report, the size of a session's context, the name of the guides folder, a `spec check` operand, the semantic analysis deadline, dead probe code in the health checker and a vacuous Verification; none changes the Linux Lane, the continuous integration workflows or the release. The entry about dead probe code touches the same source file as the lane's health check and nothing else of this context, so it stays apart. `docs/findings/` is empty. This project's inbox held fourteen pending messages each of the two times the author listed it; none concerns the lane or the workflows, and all stay pending for Triage.

## Prerequisites

Authored against `cefd4e82`, which holds v0.153.2. This Spec names no Prerequisite Spec: Spec 0418 and Spec 0419, which changed the same reference of the Roundfix skill, the same page of the user guide, the coverage record and the number of contract files that `TestRepositoryContractTestsDeclareTheirRelevance` pins, are delivered and archived.

- The Roundfix skill's version fields, the owned-skill version record, the coverage record, the Behavior Surface Record and that pinned number move with every Spec that touches them. Each Task regenerates the shared records on the tree it runs on and writes the next version and the next count, never a number copied from this plan.
- Pull request #950 is merged: this repository's Project Config holds `delivery.linux_lane: false` under six comment lines. This Spec removes the line and those comments.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The new ADR takes the number 0452, because 0450 belongs to Spec 0418 and 0451 to Spec 0419. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the product gains no endpoint, credential or network call, and it stops calling the local `docker` command. The workflows keep their permissions, their checkout and setup actions and the release's Trusted Publishing; only the runner of the steps that verify changes. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the maintainer's decision of 2026-11-24 retires the five decision records of the lane, so a new accepted ADR, numbered 0452, supersedes them and they are deleted in the same change (`docs/agents/docs-layout.md`: "a superseded ADR once an accepted successor names it with its date"). The five are ADR-0288 (ADR-0288: "The Delivery Queue's repository gate runs a Linux Lane when Project Config opts in."), ADR-0294 (ADR-0294: "The lane runs on demand."), ADR-0297 (ADR-0297: "The object store is mounted read-only."), ADR-0303 (ADR-0303: "Doctor checks the lane it will need.") and ADR-0316 (ADR-0316: "Its Park Class follows the cause."). One rule of ADR-0294 is not the lane's and survives in ADR-0452 with its words (ADR-0294: "A test that needs an External Test Tool skips naming it."). Six accepted ADRs change in part, stay accepted and gain their amendment notices in the change that adds ADR-0452 (`docs/agents/docs-layout.md`: "In the change that adds the later ADR, its header gains"): ADR-0252, whose known limit names the platform of continuous integration (ADR-0252: "A failure that only another platform produces is a known limit: CI is the Linux gate."); ADR-0285, whose contract now covers Linux as well (ADR-0285: "A repository contract compiles every test for Windows."); ADR-0384 and ADR-0385, which count the lane in the gate (ADR-0384: "The Delivery Queue's repository gate takes one unit for its repository Verification and its Linux Lane together."; ADR-0385: "the Delivery Queue's repository Verification and Linux Lane container receive it"); and ADR-0295 and ADR-0423, which name the lane or its setting as an example (ADR-0295: "Commands that already read the fetched remote after a fetch, such as the Delivery Queue and the Linux Lane, keep their reads"; ADR-0423: "so that an Item Branch cannot change it"). These apply unchanged and bound the design. ADR-0253 keeps the Full Contract Run on main and before a release (ADR-0253: "CI runs the Full Contract Run on every push to main and in the release workflow before any build or publication."). ADR-0387 keeps the suite-time budget a measured number (ADR-0387: "changing the number takes a new measurement"). ADR-0236 keeps the base a pull request is tested against (ADR-0236: "merges the default branch tip it fetched at job start"). ADR-0250 keeps the coverage record's platforms (ADR-0250: "The coverage record lists tests for every release platform."). ADR-0031 keeps the platforms Roundfix is published for (ADR-0031: "cross-compiles the Go binary for darwin/linux/windows on both architectures"). ADR-0365 and ADR-0304 forbid a change to the Run Database's schema (ADR-0365: "a schema change makes older item binaries refuse the shared Run Database"; ADR-0304: "No stage, column or schema is added."). ADR-0398 keeps the way back of an item the lane parked (ADR-0398: "A parked item always has a way back through a Delivery Retry"), and ADR-0450 keeps how that retry ends (ADR-0450: "Success requires a proven live owner unless every item is already merged or parked."). ADR-0256 makes the skill change part of this Spec (ADR-0256: "A surface whose fingerprint changed, appeared or disappeared is a Lagging Surface unless a covering skill file changed in the same range"). ADR-0056, ADR-0259, ADR-0279, ADR-0311, ADR-0349, ADR-0369, ADR-0395 and ADR-0406 are listed because each cites an ADR named above; this Spec changes none of them, and the TechSpec relates each one, marking as not applicable those whose rule this Spec does not touch. ADR-0034, ADR-0148, ADR-0271, ADR-0373, ADR-0412 and ADR-0435 are listed because each cites one of those eight; this Spec changes none of them either, and the TechSpec marks each as not applicable. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's decision of 2026-11-24, "Tirar de vez", for the Spec, its tooling paths, its delivery and a patch release, and his decision of the same day, "Os 8 arquivos (Recomendado)" (ai-memory message `01a20f32-07af-7b11-a675-a2cbc7a64743`), for the eight Governed Paths the first decision does not name, both quoted with their provenance in the authorization record. The record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. It grants no change to `go.mod`, `go.sum`, the lint or formatter configuration, any Baseline asset or any skill other than `roundfix`, and the design needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0420-roundfix-stops-verifying-on-linux/_authorization.md`; bounded files: `.roundfixrc.yml`, `Makefile`, `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runtime.md`, `skills/roundfix/SKILL.md`, `docs/references/coverage-record.json`, `docs/agents/specific-repository.md`.

## Goals

- No Roundfix command, gate or check starts a container or asks for Docker: the Delivery Queue's repository gate ends when the repository Verification passes.
- A Project Config or User Config that still names `delivery.linux_lane` or `delivery.linux_lane_image` keeps loading, and its owner is told once per key that the key is retired.
- Every pull request, every push to main and every release of this repository is verified on macOS, and no workflow runs a test on Linux.
- Linux stays a platform Roundfix is compiled and published for, its code still compiles before a pull request merges, and the README says that only macOS is tested.
- The glossary, the decision records, the user guide, the agent guides and the Roundfix skill describe a product without a Linux Lane.

## User Stories

1. As the maintainer, who runs Roundfix only on macOS, I want the Delivery Queue's gate to end at the repository Verification, so that a container on my machine never parks a correct candidate again.
2. As an operator whose configuration still names one of the lane's two keys, I want Roundfix to load it and tell me the key is retired, so that an upgrade never stops my queue.
3. As the maintainer, I want pull requests, main and releases verified on a macOS runner, so that the gate before a merge runs on the platform the product runs on.
4. As an operator with an item the lane parked before the upgrade, I want the item to keep its recorded blocker and to return through a Delivery Retry, so that the upgrade strands nothing.
5. As someone who installs Roundfix on Linux or Windows, I want the README to say that those packages are compiled and not tested, so that I know what the project proves.
6. As an operator who reads the help, `roundfix doctor`, the guide or the skill, I want no trace of a lane that no longer exists, so that I am not sent to repair Docker.

## Core Features

1. **The repository gate has no Linux Lane.** When the repository Verification of a Delivery Queue candidate passes, the gate passes. It reads no lane setting from the fetched default branch, calls no `docker` command, writes no lane log and prints no lane notice. It records no `linux-lane-skipped` warning on the item, adds no paragraph about Linux to the pull request body, and can no longer park an item as `linux-lane-container` or as a `gate-failed` with a lane cause. The gate takes its unit of the machine's verification capacity and its core budget for the repository Verification alone.
2. **The on-demand command is gone.** `roundfix deliver lane` is refused as any unknown `deliver` command is: exit `2`, the unknown-command line and the pointer to the help on standard error. `roundfix --help` and `roundfix deliver --help` no longer name it. The maintainer's answer settles that the command leaves no notice of its own (the maintainer, 2026-11-24: "Tirar de vez").
3. **The two keys are retired, never unknown.** `delivery.linux_lane` and `delivery.linux_lane_image` join the keys Roundfix registers as deprecated (`docs/user-guide/configuration.md`: "Removed keys that Roundfix registers as deprecated never break an existing config"). In Project Config and in User Config each is ignored whatever its value, and a command that loads the configuration reports it with one line on standard error in the existing form, naming the key and saying that the Linux Lane was retired and the key can be removed. A configuration the Delivery Queue reads at a commit of the default branch is never refused for naming one. A value that today's validation refuses, such as a text where a boolean is expected, is ignored the same way. Every other unknown key keeps failing strict validation.
4. **Doctor and health no longer know the lane.** `roundfix doctor` prints no `linux-lane:` line and none of the four `DR-LANE-` codes, whatever the configuration says, and the health check named `linux-lane` no longer exists. Every other line keeps its text and its order.
5. **What the lane recorded stays readable.** An item parked with a lane blocker before the upgrade keeps that text, `roundfix deliver status` prints it, and a Delivery Retry returns the item as it returns any parked item, to a queue whose gate no longer runs a lane. A `linux-lane-skipped` warning recorded on an item, a merged one included, is printed as recorded. The Run Database's schema does not change, no stored record is rewritten, and the lane logs already on disk are left to the storage rules that hold them today.
6. **Pull requests and main are verified on macOS.** The job `Verification gate` of the workflow `CI | Verify` runs on a standard GitHub-hosted macOS runner on arm64, never on a larger runner. Its steps, their order, their conditions and their commands stay: the full-history checkout, the merge of the current base with its `tested-base` notice, the incremental Verification and the documentation contracts on a pull request, and the repository Verification, the documentation contracts, the Full Contract Run and the suite-time budget on a push to main. The workflow's name, the job's name and its timeout stay, so the two required checks of the default branch and the workflows the Delivery Queue waits for need no change. The workflows that validate a pull request title and mirror files to the Secondbrain run no test and keep their runner.
7. **A release is verified on macOS.** The release runs the repository Verification and the Full Contract Run on a macOS runner before any build or publication, for a pushed tag and for a manual dispatch alike. No step that runs a test, a vet or a contract runs on Linux. What a release builds and publishes does not change: the same five platform packages, the same assets, the same order, the same Trusted Publishing.
8. **Linux and Windows are compiled, and only macOS is tested.** The Repository Contract Test that compiles every test of the repository for Windows also compiles them for Linux, so the code of every release platform that the verifying host does not build is compiled before a pull request merges. The production code and the tests that build only for Linux stay in the tree and in the coverage record, compiled and never run. The README states, where it names the platform packages, that the test suite runs on macOS only and that the Linux and Windows packages are compiled and not tested.
9. **The suite-time budget is the macOS runner's.** The budget keeps judging every push to main, now over the time the suite takes on the macOS runner. Its number and the comment that explains it state a measurement taken on that runner, because the number is a measured one (ADR-0387: "changing the number takes a new measurement"). `## Open Questions`, item 1, holds how that measurement is obtained before the first judged run.
10. **This repository's configuration loses the lane.** The Project Config of this repository no longer holds `delivery.linux_lane` nor the comments that explain it.
11. **The rules that change are changed with the work.** One Task of this Spec writes the accepted ADR-0452, and no planning change does. It states that the Delivery Queue has no Linux Lane and that the two keys are retired, that this repository verifies on macOS and compiles its other release platforms, and it restates the rule about an External Test Tool. In the same change the five decision records of the lane are deleted, the six ADRs that change in part gain ADR-0452 in their amendment lists with a dated notice after the text that changed, and every citation of a deleted record in a live document becomes a citation of ADR-0452 or the Git Locator the Spec Consistency Check prints. The glossary loses its four entries about the lane, and the two entries `## Glossary` names cite ADR-0452. The Standard Export List loses the lines of the deleted records, and ADR-0452 carries its Knowledge Mark and joins the list when the mark is true. The sentence of the repository's own guide and the comment of the Makefile that call continuous integration the Linux gate say what is true after this Spec.
12. **The guide and the skill say what the product does.** The `deliver`, `doctor` and configuration pages of the user guide and both copies of the Roundfix skill's `deliver`, `setup`, `implement` and `runtime` references lose every description of the lane, its command, its check and its codes, and the configuration page and the `setup` reference list the two keys among the deprecated ones with their warning. The Park Class table of the guide and of the skill no longer lists `linux-lane-container`. The skill's version rises through the record command, and the Behavior Surface Record and the Skill Coverage Map are re-recorded without the surface of the removed command, as the CLI guide requires of a change to a command (`docs/agents/cli.md`: "a change to command behavior updates that skill or reference in the same pull request").
13. **Each removal is characterized and declared.** Before any behavior changes, named tests state on the unchanged tree what the gate, the command, the configuration and doctor do with the lane, so that each Success Metric below fails before the work and passes after it. Every existing test, golden or fixture the change removes or moves is named with the contract it moves to. The coverage record is re-recorded by its own command, and the only test names it may lose are the ones the Task Graph lists, each a test of the lane, of its command, of its doctor check or of its two keys, as the maintainer's second decision bounds it. A test of another behavior that the change rewrites keeps its name, and none is removed to make the change pass.

## User Experience

The lines below show the intended shape. The TechSpec fixes every byte in its Surface Transcripts.

A configuration that still names the keys, on any command that loads it:

```text
config: delivery.linux_lane is deprecated and ignored; the Linux Lane was retired, so remove it
config: delivery.linux_lane_image is deprecated and ignored; the Linux Lane was retired, so remove it
```

The removed command, with exit `2`:

```text
roundfix: unknown deliver command "lane"
Run 'roundfix deliver --help' for usage.
```

`roundfix doctor` on a checkout whose configuration still says `linux_lane: true`: the two warnings above on standard error, and no `linux-lane:` line among the checks.

`roundfix deliver status` for an item the lane parked before the upgrade keeps the item row as it was recorded, with the blocker that begins `linux-lane-container:`, and the next step the status names for it leads to `roundfix deliver retry <slug>`.

## Non-Goals / Out of Scope

- Publishing fewer platforms: the Linux and Windows packages, the platform list and the release assets do not change.
- Running a test on Linux anywhere: no container, no emulator, no self-hosted runner and no second job.
- A replacement for the lane on any platform, and any new gate of the Delivery Queue.
- A larger or paid runner, a paid plan, and any setting of the repository on GitHub: the required checks keep their names.
- Removing the production code that exists only for Linux, such as the process identity read from procfs and the desktop notification, or the tests that build only for Linux.
- A change to what the repository Verification, the incremental Verification, the documentation contracts or the Full Contract Run execute, or to how tests are selected.
- A change to the workflows that validate a pull request title and mirror files to the Secondbrain.
- Rewriting a blocker or a warning already stored in a Run Database, and deleting lane logs.
- A date after which the two retired keys are refused.
- Any Baseline module, template or profile: none names the lane, and the starter workflow the Baseline gives a new Python repository keeps the Linux runner it names for that repository.
- What Spec 0418 and Spec 0419 delivered: the Delivery Retry, the status lines and the review rounds keep their behavior and their text.
- Every pending message of this project's inbox.
- Messages to the Praetor or to another repository: the operator sends them.

## Success Metrics

1. Success Metric: through the production gate of the Delivery Queue over a temporary repository whose default branch says `linux_lane: true`, with a fake runner and a fake `docker` on PATH that records each call, the gate passes when the repository Verification passes, the fake `docker` was never called, and the item carries no lane warning. On the unchanged tree the same gate calls the fake `docker`.
2. Success Metric: the built binary refuses `roundfix deliver lane` with exit `2` and the two lines of `## User Experience` on standard error, and neither `roundfix --help` nor `roundfix deliver --help` lists `lane` among the `deliver` commands. On the unchanged tree the command runs the lane and both helps list it.
3. Success Metric: a Project Config and a User Config that name the two keys, with a valid value and with a value of the wrong type, load without error, standard error holds exactly one warning per key in the form of `## User Experience`, and a key Roundfix does not know still fails strict validation with today's message. On the unchanged tree a valid value prints no warning and a value of the wrong type is refused.
4. Success Metric: with `linux_lane: true` in the checkout's configuration and the recording fake `docker` on PATH, `roundfix doctor` prints no `linux-lane:` line, the fake `docker` was never called, and every other check line is byte-identical to the one the unchanged tree prints. On the unchanged tree doctor prints the `linux-lane:` line.
5. Success Metric, a control that holds on both trees: on a temporary Run Database that holds one item parked with a `linux-lane-container` blocker and one merged item with a `linux-lane-skipped` warning, `roundfix deliver status` prints both recorded texts and exits `0`, and a Delivery Retry of the parked item succeeds as it does for an item that any other gate failure parked: it exits `0` and the item leaves the `parked` stage.
6. Success Metric: the contracts that read the workflows pass only when the job `Verification gate` and every job that runs the repository Verification, the incremental Verification, the documentation contracts, the Full Contract Run, `go test` or `go vet` name a standard macOS runner, when no job on another runner runs one of those commands, when the release runs the repository Verification and the Full Contract Run before its first build or publication step, and when the names `CI | Verify` and `Verification gate` and the order of that job's steps are the ones the unchanged tree has. On the unchanged tree these contracts fail on the runner.
7. Success Metric: the contract that compiles every test for Windows also fails when a test file compiles for macOS and for Windows and not for Linux, and names the file. On the unchanged tree it passes over such a file.
8. Success Metric: ADR-0452 is accepted, the five decision records of the lane are absent, and the six ADRs that change in part name ADR-0452 in their amendment lists with a dated notice. The glossary defines none of the four lane terms and cites ADR-0452 for the two entries `## Glossary` names. The user guide, both copies of the Roundfix skill, the agent guides, the README, the Makefile, the workflows and this repository's Project Config contain none of `Linux Lane`, `deliver lane`, `Linux gate`, `linux-lane` and `DR-LANE-`, except in the two entries that list the retired keys with their warning. The coverage record lists no test of the lane and lost no other name. The documentation, skills and glossary contracts pass.
9. Success Metric, outside evidence: the runner label the two workflows name is one of the labels GitHub's documentation lists for the standard macOS runner on arm64, quoted in `## Acceptance evidence`, and carries neither the `-large` nor the `-xlarge` suffix of a billed runner.

## Acceptance evidence

The outside-evidence row of Success Metric 9 rests on published literature, read on 2026-11-24 and copied in `references/sources.md`, section 5. GitHub Docs, "GitHub-hosted runners reference": "Use of the standard GitHub-hosted runners is free and unlimited on public repositories." The same page lists, for the standard macOS runner on arm64, the workflow labels `macos-latest`, `macos-15`, `macos-26` and `xcode-27`, with 3 processors and 7 GB of memory. GitHub Docs, "GitHub Actions billing": "Larger runners are billed for public repositories." The repository is public: the GitHub CLI answered `"visibility":"PUBLIC"` on 2026-11-24. The QA gate checks the row by reading the two workflow files against the labels quoted here, without the network.

Published practice names the state Linux is left in. Rust, "Platform Support", separates a platform that is tested from one that is only built: for Tier 1, "automated testing ensures that each tier 1 target builds and passes tests after each change", and for Tier 2, "Automated tests are not always run so it's not guaranteed to produce a working build". After this Spec macOS is in the first state, and Linux and Windows are in the second.

## Unreachable Acceptance

- criterion: the pull request that delivers this Spec passes the required check Verification gate on the macOS runner
  reason: the check runs on GitHub's runner and exists only once the pull request is open, after the QA gate
  satisfied-by: the Delivery Queue waits for the required checks of the delivering pull request before it merges, and an item whose check fails parks without merging
- criterion: the first push to main after the merge passes every step of the job on the macOS runner, the suite-time budget included
  reason: the push exists only after the merge, and its time can only be measured on GitHub's runner
  satisfied-by: the operator reads that run before cutting the release, and a budget overrun is corrected with the measured time before the release
- criterion: the release workflow passes its verification on the macOS runner and reaches its publication steps
  reason: the workflow runs only on GitHub, for a manual dispatch or a pushed tag
  satisfied-by: the operator runs the release workflow's manual dispatch for the next version, which publishes nothing, before pushing the tag

## Risks

1. **The standard macOS runner is smaller than the Linux one.** It has 3 processors and 7 GB of memory, against 4 and 16 GB. The suite may take longer than the 8 to 16 minutes a run takes today, and the failure that started this Spec was a compile killed for memory in a 7.7 GiB machine. The delivering pull request is itself verified on the new runner, so a suite that cannot pass there is found before anything merges. A larger runner is billed and is not an option of this Spec: if the standard runner cannot run the suite, the item parks and the question returns to the maintainer.
2. **The budget's number was measured on another machine.** The 420 seconds rest on Linux runs. A wrong number fails the last step of every push to main while every other step still runs.
3. **macOS jobs queue.** GitHub runs at most 5 macOS jobs at once on the account's plan and says that "jobs that run on macOS runners may take longer to start". The Delivery Queue's wait for checks is 60 minutes by default, and a queue of several items publishes one pull request at a time.
4. **A runner label moves.** `macos-latest` named macOS 15 until mid-2026 and names macOS 26 since. The TechSpec chooses between the moving label and a versioned one and says why.
5. **A tool the tests need may be absent from the macOS image.** The tests resolve eight External Test Tools and skip when one is absent, so a test that ran on the Linux runner could stop running without a failure. The TechSpec names each of the eight on the image the job uses, from the image's published list, and what the workflow does about one that is absent.
6. **A defect that exists only on Linux is no longer found by any gate.** The lane was created for one, a `text file busy` of 2026-10-17 that only Linux produced. The decision accepts it: nobody runs the product on Linux, the README says the package is not tested, and the compile contract still refuses code that does not build there.
7. **The release is the only way to publish.** A mistake in its new shape appears only when a release runs. The manual dispatch, which publishes nothing, is run before the next tag.

## Validation

The decisions were checked on 2026-11-24 against the code of the lane, the configuration loader, doctor, the four workflows, the Makefile, the glossary, the ADRs `## Project Constraints` names, the user guide and the Roundfix skill on `808f9001`. After Spec 0418 and Spec 0419 merged, the author measured again on `cefd4e82` the governed set, the lane's files, its tests in the coverage record, the four workflows and this repository's Project Config.

- The code bounds the removal. The lane lives in two files of the command layer and is called from one place of the gate. Its blocker and its warning live in text the Run Database already holds for every item, so no schema changes. The configuration loader already has a registry of deprecated keys that are removed before the strict decode and reported with one line each, which five keys use today.
- The workflows correct the brief. Two workflows verify on Linux, not one: the release runs the repository Verification and the Full Contract Run before it publishes. Core Feature 7 exists for that reason.
- One rule of a retired record is not about the lane: the rule that a test skips, naming the tool, when an External Test Tool is absent. Two glossary entries cite it, and the test fixture enforces it, so ADR-0452 restates it.
- Spec 0418 changed how a Delivery Retry ends and what status prints for a queue. Core Feature 5 and Success Metric 5 were read again against its decision record, and ADR-0450 joined the obligations of `## Project Constraints`.
- Only this repository sets `linux_lane` among the fourteen repositories of the maintainer's machine that carry a Roundfix Project Config, and no Baseline asset names it. The effect on another repository is the warning of Core Feature 3 if its configuration ever names a key, and the absence of a command it never ran.
- Secondbrain, read only: `wiki/index.md`, then the Indexus tool `pesquisar_acervos`, which answered with wiki results. `wiki/sources/digest-vortex-portas-do-docker-libpq-e-saida-do-psql-2026-11-04.md` records a test that failed because the Linux runner image lacks a tool the maintainer's machines have; it is the reason Risk 5 asks for the eight tools on the macOS image. `wiki/sources/digest-roundfix-checkout-completo-para-testes-de-evidencia-historica-2026-09-09.md` and `wiki/sources/digest-roundfix-ci-rerun-merge-antigo-2026-10-06.md` record why the job checks out the full history and merges the current base; Core Feature 6 keeps both steps. `wiki/sources/digest-roundfix-macos-sigkill-binario-assinado-2026-10-16.md` records a defect that only macOS produces, found on the maintainer's machine and never by the Linux job; it supports verifying where the product runs.
- Exa and the published pages: GitHub's runner reference, billing page and limits page, the runner images' README, Rust's platform support page, Go's porting policy and Go's source installation page. The runner reference and the billing page settled that a standard macOS runner costs nothing for a public repository, so the choice of runner is not a question of a paid service. The same reference changed Risk 1: the runner has less memory than the machine where the lane was killed. Rust's tiers gave the two states Core Feature 8 and the README name.

The Research Subagent ran the codebase sweep and the web sweep through Exa. The author ran the commands, the governed-path probe and the Secondbrain search, which the Research Subagent has no tool for, and read again the two GitHub pages `## Acceptance evidence` quotes.

## Rule challenges

None. The five decision records of the lane and its four glossary entries forbid the goal, and the maintainer's decision already retires them (the maintainer, 2026-11-24: "Tirar de vez"), so no Rule Challenge is asked; Core Feature 11 records the change.

## Glossary

- changes: **Base Test Tools**
- changes: **External Test Tool**

## Skills

- unchanged: guide docs/user-guide/release-runbook.md — the page changes only the order in which it lists the release workflow's verification and its tag validation, and the covering reference describes neither the workflow's steps nor their order
- unchanged: command release plan — the command, its help, its output and its exit codes do not change; the page that is its source changes only that order

## Decisions

- 2026-11-24: Roundfix tests on Linux in a Docker container before the pull request and the maintainer does not use Roundfix on Linux; what becomes of it? — "Tirar de vez" (ai-memory message `01a20d73-109f-7bb2-a2b6-718854f6b25d`, from the Praetor, 09:10 BRT).
- 2026-11-24: Spec 0420 needs to change eight protected files the first decision does not name, the texts of the Roundfix skill about the lane, the coverage record of the tests and one sentence of the repository's guide; may it? — "Os 8 arquivos (Recomendado)" (ai-memory message `01a20f32-07af-7b11-a675-a2cbc7a64743`, from the Praetor, 17:18 BRT).

The first decision, verbatim:

```text
DECISÃO DE MARCIO
pedido: PEDIDO de 24/11 09:05 BRT (Linux Lane da 0418)
pergunta: Roundfix, portão de testes em Linux que roda no Docker do mini antes do PR: Marcio não usa o Roundfix em Linux. Como fica, a começar pela 0418 que está parada nele? (Tirar do local, só no CI; tirar de vez; subir a memória por ora)
resposta: "Tirar de vez" — com o contexto que Marcio escreveu antes: "Eu estou quase pedindo para remover esses testes do linux, não uso o roundfix no linux, o quaestor-legacy só roda no mini e o quaestor vai rodar em container. Isso está gerando tanto problema. Talvez só deixar algum tipo de teste no CI antes do merge na main."
alcance: (1) Spec própria no repositório roundfix para remover a verificação Linux, tanto a Linux Lane local quanto o job Linux do CI, ficando só macOS; inclui os caminhos de tooling necessários (Makefile, CI, configuração da Lane), implementação, PR, merge e release patch. (2) A 0418 pode seguir sem o portão Linux: se a Lane não puder ser desligada sem a Spec nova, priorize essa Spec antes de retomar a 0418, dentro das filas novas já autorizadas (01a20d32). Se a remoção afetar outros repositórios que usam o Roundfix (Lane em quaestor, vortex etc.), descreva o efeito no PRD e avise o Praetor.
quando: 24/11/2026 09:10 BRT
```

The second decision, verbatim. It answers the Reserved Question the author returned in the first round of this phase, which the Supervisor session sent as the decision request `01a20d89-cf32`:

```text
DECISÃO DE MARCIO
pedido: 01a20d89-cf32
pergunta: Roundfix, Spec 0420 (tirar de vez os testes Linux): ela precisa mudar 8 arquivos protegidos que a decisão não citou — os textos da skill do Roundfix que descrevem o portão Linux, o registro de cobertura dos testes e uma frase do guia que chama o CI de "portão Linux". Pode? (Os 8 arquivos; menos a frase do guia; não, só desligar)
resposta: "Os 8 arquivos (Recomendado)"
alcance: o alcance do pedido 01a20d89-cf32, até o release patch: .agents/skills/roundfix/SKILL.md, .agents/skills/roundfix/references/{deliver,setup,implement,runtime}.md, skills/roundfix/SKILL.md, docs/references/coverage-record.json (só pelo comando próprio; só saem testes da Lane, do comando dela, da checagem do doctor e das duas chaves), docs/agents/specific-repository.md (só a frase "CI is the Linux gate"), com as regenerações listadas no pedido; mais .roundfixrc.yml, Makefile, .github/workflows/ci-verify.yml e release.yml (CI para macOS). Sem go.mod, lint, Baseline nem outra skill. Chaves antigas aceitas com aviso.
quando: 24/11/2026 17:18 BRT
```

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0420 of 2026-11-24 and its addendum of the same day, the maintainer's two decisions above, the code, the workflows and the decision records on `cefd4e82`, and GitHub's published documentation.

Settled by a source, with no question left:

- The route and the Delivery Type. The Spec removes a capability and adds none, so among the four Delivery Types it is a `refactor`. It changes product behavior, a command and two keys, so it takes `write-prd`, `write-techspec` and `write-tasks`. A refactor is never judged for the Idea stage, so this PRD carries no Idea-stage judgment. Source: the operator's brief and the Spec routing guide.
- The removal of the lane from the gate, of its command and of its check in doctor. Source: the decision, "Tirar de vez".
- No Rule Challenge. Source: the same decision, as `## Rule challenges` says.
- What a configuration that still names the keys meets: the product's existing contract for a removed key, not a new one. Source: the configuration page of the user guide and the loader's registry of deprecated keys. The scope of the second decision says the same in its last sentence.
- The runner costs nothing. Source: GitHub's documentation, quoted in `## Acceptance evidence`, and the repository's visibility.
- The Release Ceiling, `patch`. Source: the decision. A `refactor` delivery that changes released code is classified by hand in the release plan; the operator classifies it as a patch, with the decision as the reason. The operator's addendum names v0.153.3.
- The Governed Paths the Spec changes: the twelve the authorization record lists. Source: the two decisions.
- No Prerequisite Spec. Source: the operator's addendum and the default branch, where Spec 0418 and Spec 0419 are archived.
- What becomes of an item the lane parked or warned about before the upgrade: its recorded text is printed as it is, a Delivery Retry returns it, and no stored record is rewritten. Source: the two decision records that bound it (ADR-0398: "A parked item always has a way back through a Delivery Retry"; ADR-0365: "a schema change makes older item binaries refuse the shared Run Database").
- The size: one context, and up to four implementation Tasks plus the QA gate. Source: the operator's brief.

The Delegated Decisions below are the questions no source settles. The author returned each with its options and read none as reserved. The Supervisor session adopted each as recommended on 2026-11-24, under the project's Authoring Question Preference, `follow`.

- delegated: 2026-11-24, prd
  - question: After the Linux job leaves, what verifies a pull request and a push to main?
  - options: the same Verification gate job with the same steps on a standard GitHub-hosted macOS runner; a light job on a Linux runner that only builds and vets; no test in continuous integration, with the local macOS verification of the Delivery Queue as the only gate
  - adopted: the same Verification gate job with the same steps on a standard GitHub-hosted macOS runner — the decision leaves only macOS, the maintainer's context asks for a test before the merge into main, a standard runner is free for a public repository, and keeping the job's name keeps the required checks
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-24, prd
  - question: The release also runs the repository Verification and the Full Contract Run on Linux before it publishes. What becomes of that?
  - options: the release verifies on a macOS runner before any build or publication and what it builds and publishes does not change; the whole release runs on a macOS runner; the release stops verifying and relies on the run of the push to main
  - adopted: the release verifies on a macOS runner before any build or publication and what it builds and publishes does not change — it keeps the rule that the Full Contract Run precedes every release, removes the last test that runs on Linux, and leaves the publication steps and Trusted Publishing as they are proven today
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-24, prd
  - question: Roundfix is still published for Linux. What does the repository prove about Linux after this Spec?
  - options: the contract that compiles every test for Windows also compiles them for Linux, nothing is tested there, and the README says so; nothing is compiled for Linux before the release's own cross-compile and the tests that build only for Linux are removed; the Linux packages stop being published
  - adopted: the contract that compiles every test for Windows also compiles them for Linux, nothing is tested there, and the README says so — it runs on macOS in seconds and executes nothing on Linux, it gives Linux the state the repository already chose for Windows, and without it a build error in code that only Linux compiles would first appear when a release is being published
  - under: Project Config
  - answer: none yet

The four entries below are the questions of the TechSpec phase that no source settles. Its author read none as reserved and adopted each as recommended on 2026-11-24, as the operator's brief for that phase directs, under the same preference.

- delegated: 2026-11-24, techspec
  - question: How is the suite-time budget measured on the macOS runner before the first push to main is judged?
  - options: the number stays at 420 seconds, its comment says that it was measured on the Linux runner and awaits a measurement on the macOS runner, and the first push to main is that measurement; a pull request opened before the Run only to time the suite on the macOS runner, which needs a push and a pull request nobody authorized; a larger number chosen without a measurement
  - adopted: the number stays at 420 seconds, its comment says that it was measured on the Linux runner and awaits a measurement on the macOS runner, and the first push to main is that measurement — a Task reaches no runner, the delivering pull request runs the incremental Verification and not the suite the budget times, and the budget's decision record lets the number change only with a measurement
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-24, techspec
  - question: Which label names the runner, the moving `macos-latest` or a versioned one?
  - options: the versioned label `macos-26`; the versioned label `macos-15`; the moving label `macos-latest`
  - adopted: the versioned label `macos-26` — a versioned label makes a new image arrive by a reviewed change, and of the two versioned labels GitHub supports it is the newer one, the image `macos-latest` names today and the one that will be retired last
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-24, techspec
  - question: How does the release keep its verification before the publication: a job of its own that the publishing job waits for, or one job on the macOS runner?
  - options: a job on the macOS runner that the publishing job needs, with every publishing step unchanged on its runner; three jobs, so that the tag is still validated before the verification; the whole release in one job on the macOS runner
  - adopted: a job on the macOS runner that the publishing job needs, with every publishing step unchanged on its runner — it is the change the authorization record bounds for that workflow, the verification runs without the publishing permissions, and a failed publication is re-run without repeating it; its cost is that a wrong tag is refused after the verification, and the release runbook says so
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-24, techspec
  - question: Which Park Class and next step does status show for a stored `linux-lane-container` blocker once no code produces it?
  - options: the `environment` class it has today, with a next step that names the Delivery Retry and says the lane was retired; no case for it, so that it shows as `unclassified` with the generic next step; the class and the next step of today, word for word
  - adopted: the `environment` class it has today, with a next step that names the Delivery Retry and says the lane was retired — the generic step tells the operator to resolve a blocker that can no longer be resolved, today's step tells the operator to repair Docker, and both status and the follow-up's Pending Question print this text
  - under: Project Config
  - answer: none yet

## Open Questions

None. The TechSpec settles the five questions this section held.

1. The suite-time budget on the macOS runner: the first Delegated Decision of the TechSpec phase, in `## Decisions`.
2. The runner's label: the second.
3. The release's job layout: the third.
4. The Park Class of a stored lane blocker: the fourth.
5. The Skills Declaration: settled by measurement. Every Behavior Surface the removal moves has a covering reference that changes, and `## Skills` holds the two entries of the release runbook, whose order of steps the third decision changes.
