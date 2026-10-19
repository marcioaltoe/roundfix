---
spec: 0282-the-linux-lane-runs-like-ci
status: active
created: 2026-10-19
surfaces: [backend, cli, docs]
---

# The Linux Lane runs like CI

On 2026-10-18 the Linux Lane ran for the first time, on Spec 0280's
candidate, and parked the item `gate-failed` although every package the
candidate changed passed in the container. Two environment gaps failed it:

- `TestGuidanceCompositionJourney/rust-cli` failed with `exec: "rustfmt":
  executable file not found in $PATH`, because the `golang:1.26` image has no
  Rust toolchain and the CI runner does.
- `TestCLITestHelpersEndWithTheirTestBinary` timed out after 594 s, because
  the container runs without an init and an orphaned helper was never reaped.

The maintainer turned the lane off the same day, in commit `a49bd654`. The
operator could not tell from the queue that the container, not the change,
had failed: the gate's reason never reached the blocker, which read only
`gate-failed`.

After this Spec:

- The lane runs under `--init`, on the image `delivery.linux_lane_image`
  names, with `golang:<go.mod go version>` as the default.
- A test that needs an **External Test Tool** skips naming it when the tool
  is absent, in the container exactly as on a machine without it.
- A **Lane Container Failure** parks the item `linux-lane-container` in the
  `environment` Park Class; a **Lane Change Failure** parks it `gate-failed`
  with its reason in the `finding` Park Class. The lane log names the
  outcome too.
- `roundfix deliver lane` runs the lane on demand and exits with its
  outcome, so the operator can prove a run passes before turning the lane
  back on in a change of its own.

## Sources

One Backlog Entry, adopted into this Spec's `references/` (see
`references/_index.md`):

1. `references/2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md`,
   recorded with the commit that turned the lane off (`a49bd654`). It names
   the two failures above and asks for `--init`, a configurable image or
   tests that skip with a reason, a container failure told apart from a
   change failure, and then the lane turned back on.

The other open Backlog Entries share no context: the Git Locator check and
the Run Database WAL report belong to Spec 0283. No unresolved Finding
exists under `docs/findings/`.

## Prerequisites

None. This Spec is authored against `f382480b` (v0.91.0) and delivers after
Spec 0281 and before Spec 0283. Neither changes the Linux Lane, the `deliver`
guide's lane section or the configuration guide's lane section; whichever
merges first may raise the Roundfix skill version and re-record the Behavior
Surface Record, so this plan is rebased and checked again before its Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. One Project Config key, one blocker name and one subcommand are
  added; each is a fixed name. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. The lane keeps `--network none` and
  `--pull never`, and every test uses compiled Docker and Go stand-ins on
  temporary repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the lane, its tests and its
  documentation follow recorded decisions. ADR-0294 (this Spec): "The Linux
  Lane runs like CI". ADR-0288: "The Delivery Queue gate runs the changed Go packages in a Linux
  Lane", which ADR-0294 amends. ADR-0125: "A spawned
  fixture is a compiled binary, never a written script". ADR-0259: "A test
  runs in parallel unless it names why it cannot". ADR-0271: "A test bounds
  by real time only what it asserts". ADR-0285: "Every test file compiles for
  Windows, and a repository contract proves it". ADR-0189: "An owned skill's
  version names its content, and its minimum is the bundle". ADR-0256: "A
  release waits for the skills that describe a changed surface". ADR-0280: "A
  guide Task records its own surface fingerprint, a PRD records who decided,
  and the glossary is held to its own rule". ADR-0104: "A Spec accepts on
  evidence it did not author". ADR-0240: "One QA partial policy, and rows a
  Run sandbox cannot reach". ADR-0276: "A row the authorization keeps offline
  qualifies, and a byte invariant compares with the Delivery Base". ADR-0265:
  "A Spec proves what it claims, and every active Spec is held to it".
  ADR-0183: "An attributed claim carries a receipt the check proves". Set no obligation here, because this Spec changes none of their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0184 (a TechSpec states a command surface as a transcript); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof). Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ approval of 2026-10-19 through AskUserQuestion, option "Mexer em
  Governed Paths declarados" ("'Concedo' para os caminhos que cada Spec
  declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0282-the-linux-lane-runs-like-ci/_authorization.md`; bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `docs/agents/specific-repository.md`,
  `internal/cli/baseline_release_gate_test.go`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- A lane run reaps every orphan, so a test that waits for a helper to end
  passes in the container as it does on the host.
- The operator chooses the lane's image in Project Config, and an item branch
  cannot choose it.
- A test that needs a tool the image lacks skips with the tool's name, so a
  missing tool never fails the lane.
- The queue's blocker and the lane log say whether the container or the
  change failed.
- The operator can run the lane on demand and read its outcome from the exit
  status.

## User Stories

1. As the maintainer, I want the lane to run under an init, so that a helper
   orphaned inside the container does not time out a passing test.
2. As an operator, I want to name the lane's image in Project Config, so that
   I can give the lane the tools my tests need without a repository image
   build.
3. As a test author, I want a test that needs an absent tool to skip naming
   it, so that the lane and a bare machine report the same thing.
4. As the operator of a parked item, I want the blocker to say whether the
   container or the change failed, so that I fix the right thing.
5. As the maintainer, I want to run the lane on demand and get its outcome as
   an exit status, so that I turn the lane back on only after a run passes.

## Core Features

1. The init and the configured image.
   - `docker run` gains `--init`.
   - `delivery.linux_lane_image` is a string, default empty. Empty keeps
     `golang:<go.mod go version>`. A non-empty value must be a container
     image reference and is refused otherwise.
   - The gate reads it from the fetched default-branch Project Config, with
     User Config beneath it, like `delivery.linux_lane`.
   - The lane never pulls: an absent configured image is skipped as the
     default one is.
2. External Test Tools.
   - `testfixture.RequireTool` resolves a tool on `PATH` and skips the test
     with `requires <tool>, which is not on PATH` when it is absent.
   - The rust-cli journey resolves `rustfmt` through it, and the real acpx
     integration test resolves `node` through it.
   - A repository rule refuses a test that resolves a tool outside the Base
     Test Tools by a literal name in any other way.
3. Container and change failures.
   - The lane script checks `tar`, `git` and `go`, unpacks and records the
     candidate, and exits `3` with a `roundfix-linux-lane: container: <cause>`
     line when any of that fails.
   - Exit `3`, `126` or `127`, an exit of `128` or more, and the 20-minute cap
     are a Lane Container Failure. Any other non-zero exit is a Lane Change
     Failure.
   - A Lane Container Failure parks `linux-lane-container: <cause>; see <lane
     log>`, Park Class `environment`. A Lane Change Failure parks
     `gate-failed: Linux Lane change failure: go test exited <n>; see <lane
     log>`, Park Class `finding`. A repository Verification failure keeps the
     bare `gate-failed`.
   - The lane log ends with the outcome's line.
4. The on-demand lane.
   - `roundfix deliver lane [--all] [<revision>]` runs the lane against a
     revision of the current checkout, `HEAD` by default, with no queue, no
     Spec and whatever `delivery.linux_lane` says.
   - It reads the image from the checkout's Project Config over User Config,
     selects the changed packages as the gate does, or `./...` with `--all`.
   - It exits `0` when it passes, `1` on a Lane Change Failure, `3` when it
     skips or meets a Lane Container Failure, `2` on a usage or preflight
     error, and `130` when interrupted after removing its container.
5. Documentation.
   - The configuration guide and the `deliver` guide, the Roundfix `deliver`
     reference, the repository rules in `docs/agents/specific-repository.md`
     and `CONTEXT.md` describe each change.
   - The glossary adds **External Test Tool**, **Lane Container Failure** and
     **Lane Change Failure**, and revises **Linux Lane**.

## Non-Goals / Out of Scope

- Turning `delivery.linux_lane` back on in this repository. The operator does
  it in a separate change after this Spec merges and an on-demand run passes.
- Building, publishing or pulling an image, or adding a pull step to the lane.
- Changing the lane's skips, the module cache mount, `--network none`, the
  variables the container receives, the package selection or the cap's
  length.
- Making a missing Base Test Tool skip: a test that cannot find `go`, `git`,
  `gofmt` or `make` still fails.
- Changing the CI workflows, or making CI fail when an External Test Tool is
  absent.
- Carrying the repository Verification's failure reason into the blocker.

## Success Metrics

1. Success Metric: with compiled Docker and Go stand-ins, the lane's
   `docker run` arguments hold `--init` before the image, and with
   `delivery.linux_lane_image: registry.example.com/team/go-ci:1.26` at the
   fetched default the run inspects and uses that image while a different
   value on the item branch is ignored. `-v`, `golang 1.26`, `Golang:1.26`
   and a list are each refused with `delivery.linux_lane_image must be a
   container image reference`.
2. Success Metric: with a `PATH` that holds only `go`, `git` and `sh`,
   `TestGuidanceCompositionJourney/rust-cli` skips with `requires rustfmt,
   which is not on PATH`, and the External Test Tool rule finds no test that
   resolves a tool outside the Base Test Tools by a literal name.
3. Success Metric: a stand-in container that exits `3` with a marker, `126`,
   `127` or `137` parks the item `linux-lane-container: <cause>; see <lane
   log>` with Park Class `environment`; one that exits `1` or `2` parks it
   `gate-failed: Linux Lane change failure: go test exited <n>; see <lane
   log>` with Park Class `finding`; and the lane log ends with that outcome
   line.
4. Success Metric: `roundfix deliver lane` reproduces the four Surface
   Transcripts of the TechSpec with exit `0`, `1`, `3` and `3`, in a
   repository whose Project Config sets `delivery.linux_lane: false`, and
   `--all` hands `./...` to the container.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-19.

- This repository's CI workflow, which no Task of this Spec authors.
  `.github/workflows/ci-verify.yml` runs on `runs-on: ubuntu-latest` and
  installs no Rust toolchain, yet the rust-cli journey runs there; the runner
  image carries `rustfmt`. The QA gate reads the file at the audited head.
- Docker Docs, "docker container run", section "Specify an init process"
  (https://docs.docker.com/reference/cli/docker/container/run/): "Specifying
  an init process ensures the usual responsibilities of an init system, such
  as reaping zombie processes, are performed inside the created container."
- Docker Docs, "Running containers", section "Exit status"
  (https://docs.docker.com/engine/containers/run/): "Exit code `126`
  indicates that the specified contained command can't be invoked." and
  "Exit code `127` indicates that the contained command can't be found."
- krallin/tini README (https://github.com/krallin/tini): "All Tini does is
  spawn a single child (Tini is meant to be run in a container), and wait for
  it to exit all the while reaping zombies and performing signal forwarding."
- actions/runner-images, "Ubuntu 24.04" README
  (https://github.com/actions/runner-images/blob/main/images/ubuntu/Ubuntu2404-Readme.md):
  under "Rust Tools", "Packages", the line "Rustfmt 1.10.0".

## Unreachable Acceptance

- criterion: a real Linux Lane run of the audited head, in a locally present
  `golang:1.26` image, passes with `TestCLITestHelpersEndWithTheirTestBinary`
  passing and `TestGuidanceCompositionJourney/rust-cli` skipping.
  reason: no Verification may need Docker, and a Run sandbox may not reach
  the Docker daemon; the QA gate runs it only when it can.
  satisfied-by: the operator runs `roundfix deliver lane` on `main` after
  this Spec merges and before the change that sets `delivery.linux_lane:
  true`.

## Validation

The decisions were checked on 2026-10-19 against these sources.

- Secondbrain.
  - Read `wiki/index.md`, then ran three `qmd query` runs from the Secondbrain
    root: "linux lane docker container init zombie orphan reaping test",
    "test skips when external tool rustfmt missing CI container image" and
    "docker image configurable CI parity go test container".
  - They returned the mirrors of this repository's `CONTEXT.md`, ADR-0288,
    the `deliver` guide and the adopted Backlog Entry, which agree with this
    repository and hold no decision on an init or an image.
  - They also returned the `docker` skill mirror,
    `projects/skills/mirror/skills/08-release/docker/SKILL.md`, whose
    checklist asks for "Signals handled correctly (exec form CMD, or
    `dumb-init`/`tini`)". It supports an init as PID 1 rather than a
    hand-written reaper.
  - The Secondbrain was read only. No research Inbox Entry was captured,
    because the lead's briefing for this Spec forbids writing there.
- Exa.
  - Supporting: the Docker `--init` and exit-status passages, the tini README
    and the runner image README, quoted above. Docker's
    "Run multiple processes in a container" page
    (https://docs.docker.com/engine/containers/multi-service_container/) adds
    that `--init` "handles reaping of all processes when the container
    exits".
  - Grammar: the distribution reference library
    (https://github.com/distribution/reference/blob/main/regexp.go) defines a
    reference as a name of lower-case path components with an optional
    domain, an optional `:tag` of `[\w][\w.-]{0,127}` and an optional
    `@digest`; the configured image is held to it.
  - Challenging: tini forwards a signal only to its immediate child unless
    `-g` is passed, so a container that ignores the cap's cancellation would
    keep running. The lane already removes the container with `docker rm -f`
    on the cap and on cancellation (ADR-0288), so `--init` changes nothing
    there.
- Measurement on 2026-10-19 with Docker 29.6.1 and `golang:1.26`
  (`golang@sha256:37cfcca63222ec499bec5e8ad7aad52dbf7529db36d5812b8eae2f0621aa92d6`,
  go 1.26.9, linux/arm64); the TechSpec records each probe.
  - Without `--init`, an exited orphan stayed a zombie under PID 1 and
    `kill(pid, 0)` succeeded on it; with `--init`, PID 1 was `docker-init`
    and the orphan was gone.
  - The image carries `sh`, `tar`, `git`, `go`, `gofmt`, `make`, `ps`,
    `kill`, `mkfifo`, `cat` and `id`, and lacks `rustfmt`, `cargo`, `node`,
    `bun`, `npm` and `docker`.
  - `docker run` exited `127` for a missing command, `126` for a
    non-executable one, `137` when the process was killed and `125` for an
    absent image under `--pull never`; `go test` exited `1` for a failing
    test, a build failure and a panic.

## Glossary

- adds: **External Test Tool**
- adds: **Lane Container Failure**
- adds: **Lane Change Failure**
- changes: **Linux Lane**

## Decisions

- 2026-10-19 (about 16:20 BRT): the cycle's decomposition — "Três Specs
  (Recomendado)".
- 2026-10-19 (about 16:20 BRT): this Spec's design — "--init + imagem
  configurável + skip com motivo (Recomendado)", option text "O lane roda com
  --init; a imagem passa a ser configurável em .roundfixrc; testes que
  precisam de ferramenta externa pulam com motivo quando ela falta; falha de
  container é distinguida de falha da mudança. Sem build de imagem própria."
- 2026-10-19 (about 16:20 BRT): the authorizations — "Religar
  delivery.linux_lane: true", "docker pull/build da imagem", "Mexer em
  Governed Paths declarados" (option text "'Concedo' para os caminhos que cada
  Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
  squash + release + push".
- 2026-10-19 (about 16:20 BRT): costs — "Regras do ciclo AH/AI
  (Recomendado)"; failures — "Regras AE–AH e parar na proposta
  (Recomendado)"; authoring — "Opus + strict check, enfileirar sem esperar
  você (Recomendado)".
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present. Decisions were taken from the lead's
briefing of 2026-10-19 for Spec 0282 and from the code and probes
characterized in the TechSpec.

Taken from those sources, with no further interview:

- The configuration key is `delivery.linux_lane_image`, read from the fetched
  default like the opt-in, and its value is held to the image reference
  grammar. See ADR-0294.
- The lane's own setup exits `3` on failure, and `3`, `126`, `127`, `128` or
  more and the cap are the container's; a new `linux-lane-container` blocker
  carries them in the `environment` Park Class. See ADR-0294.
- The skips ADR-0288 defined stay skips. See ADR-0294.
- The on-demand command is `roundfix deliver lane`, because the lane belongs
  to the Delivery Queue's gate and the operator needs it without a queue.
  See ADR-0294.
- Reaching Docker is not required of any Verification; the real run is the
  operator's, declared under Unreachable Acceptance.

## Open Questions

None.
