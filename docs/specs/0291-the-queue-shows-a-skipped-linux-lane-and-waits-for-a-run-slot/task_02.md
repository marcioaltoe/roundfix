---
task: task_02
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
status: completed
type: backend
complexity: medium
---

# Task 02: Doctor checks the Linux Lane, and the references name the Lane Environment Skip

## Overview

`roundfix doctor` said nothing while the image the enabled Linux Lane needed
was missing (`_prd.md` → Sources 1). This Task adds a `linux-lane:` line,
printed only while the checkout's effective `delivery.linux_lane` is `true`,
that reports Docker and the lane's image without pulling or starting
anything, and brings the Roundfix references and the glossary up to date for
both lane surfaces. It is verifiable on its own through the lane's compiled
Docker stand-in and the doctor's fake dependencies.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-the-queue-skips-the-linux-lane-without-a-trace.md`](references/2026-10-20-the-queue-skips-the-linux-lane-without-a-trace.md)
   (`_prd.md` → Sources 1), by implementing API Contract 4 of
   `_techspec.md` in `internal/cli/doctor.go`, `internal/cli/health.go` and
   `internal/cli/deliver_linux_lane.go`, with the `linuxLane` doctor
   dependency and `HealthCheckLinuxLane`.
2. MUST resolve the image through one function in
   `internal/cli/deliver_linux_lane.go` that both `runLinuxLane` and the
   doctor check call: the configured image when set, else `golang:<go
   directive>` of `go.mod` at the given revision. `runLinuxLane` keeps its
   skip reasons and their order; no second copy of the `go` directive parsing
   remains.
3. MUST state the hazard in the code's behavior, as API Contract 4 does. The
   data sources are the checkout's effective configuration, `go.mod` at
   `HEAD` and the Docker CLI's answers; the destination is one standard
   output line. The check resolves `docker` on the command environment's PATH,
   runs only `docker info` and `docker image inspect <image>`, each bounded
   by ten seconds, never pulls, runs or removes anything, and prints no
   environment value. Its findings are `warn`, so they never change doctor's
   exit code.
4. MUST print the line after `storage:` and before `codex:` only when the
   dependency is set and `loaded.Config.Delivery.LinuxLane` is `true`. Every
   existing doctor test keeps its name and expected output, because their
   configurations leave the lane off; this Task declares no break.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows, ending every child it starts in `t.Cleanup`, and
   declaring no real-time budget under ten seconds, in
   `internal/cli/doctor_linux_lane_test.go`:
   - `TestDoctorLinuxLaneCheckReportsDockerAndTheImage`, with the subtests
     `_techspec.md` → Testing Approach lists, reusing `linuxLaneTools` from
     `internal/cli/deliver_linux_lane_test.go`, each asserting the whole
     rendered line, for example `linux-lane: warn (DR-LANE-IMAGE-MISSING:
     image golang:1.26.1 is not present locally; next: docker pull
     golang:1.26.1)`, and the stand-in's recorded arguments;
   - `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`, through
     `runDoctorCommand` with fake doctor dependencies.
6. MUST document the change:
   - in `docs/user-guide/commands/doctor.md`, the `linux-lane:` line in the
     check list, when it is printed, the image it inspects, and the four codes
     `DR-LANE-DOCKER-MISSING`, `DR-LANE-DOCKER-UNREACHABLE`,
     `DR-LANE-IMAGE-MISSING` and `DR-LANE-IMAGE-UNKNOWN` with status `warn`
     and next action in its findings table;
   - in `.agents/skills/roundfix/references/setup.md`, the same four codes in
     its findings table and the `linux-lane` check where doctor's checks are
     listed;
   - in `.agents/skills/roundfix/references/runtime.md`, the `linux-lane`
     check in the list of checks doctor runs;
   - in `.agents/skills/roundfix/references/deliver.md`, the Lane Environment
     Skip as task_01 documented it in the user guide: the expected skips, the
     warning `linux-lane-skipped: <reason>` on the item, the Pull Request
     paragraph, and `roundfix doctor`'s `linux-lane:` line.

   MUST search the guides, the agent guides and the skill copies for wording
   that lists doctor's checks without the Linux Lane, or says a lane skip
   leaves only the notice, change each one in the files this Task declares,
   and name any other file found in the Task Result. No text inside any `###
   QA settlement` section may change.
7. MUST add **Lane Environment Skip** to `CONTEXT.md` through the
   `domain-modeling` skill: a skip of an enabled Linux Lane for any reason but
   having nothing to test, or of a lane whose fetched configuration cannot be
   read, which passes the gate and is recorded on the item, shown by `roundfix
   deliver status`, the queue's Pull Request body and `roundfix doctor`,
   linking ADR-0303.
8. MUST run these commands, because this Task edits a user guide page and the
   Roundfix references. Before them, raise both version fields of the
   Roundfix skill's canonical `SKILL.md` (`metadata.version` and the top-level
   `version`) by one patch over the value the tree holds when this Task
   starts, because the record command refuses changed content under an
   unchanged version. The Agent runs these commands itself before handing
   back; no record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.

## Subtasks

- [ ] Extract the one image resolver and call it from the lane.
- [ ] Add the `linux-lane` doctor check and its dependency.
- [ ] Add the check and doctor command tests.
- [ ] Document the check and the Lane Environment Skip in the guide, the references and the glossary.
- [ ] Sync the skill, record its version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] With the lane on and the image absent, doctor prints `linux-lane: warn (DR-LANE-IMAGE-MISSING: …; next: docker pull <image>)` and exits as before.
- [ ] Docker missing or unreachable gives its own code, and a healthy setup prints `linux-lane: ok (docker reachable; image <image> present)`.
- [ ] With the lane off, doctor's output is unchanged.
- [ ] The guide, the `setup`, `runtime` and `deliver` references and the glossary describe the check and the Lane Environment Skip.

## Consumers

- the doctor `linux-lane` check — consumers: `./internal/cli`
- the Linux Lane image resolver — consumers: `./internal/cli`
- the Roundfix references and the glossary — consumers: `./skills`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0303-a-linux-lane-skipped-for-its-environment-is-shown-where-the-delivery-is-read.md`
- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `internal/cli/deliver_linux_lane_test.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `internal/cli/readiness.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/deliver_linux_lane.go`
- creates: `internal/cli/doctor_linux_lane_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDoctorLinuxLaneCheckReportsDockerAndTheImage|TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn|TestRunDoctorMissingRepositoryRoot|TestDoctorReportsPrePRReviewPolicy|TestLinuxLaneSkipsWithANamedReason|TestLinuxLaneRunsTheConfiguredImage|TestLinuxLaneFailureFailsTheGate)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorLinuxLaneCheckReportsDockerAndTheImage TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn TestRunDoctorMissingRepositoryRoot TestLinuxLaneSkipsWithANamedReason TestLinuxLaneRunsTheConfiguredImage TestLinuxLaneFailureFailsTheGate; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/doctor.md .agents/skills/roundfix/references/setup.md; do grep -qF -- 'DR-LANE-DOCKER-MISSING' "$f" && grep -qF -- 'DR-LANE-DOCKER-UNREACHABLE' "$f" && grep -qF -- 'DR-LANE-IMAGE-MISSING' "$f" && grep -qF -- 'DR-LANE-IMAGE-UNKNOWN' "$f" || { printf 'missing a DR-LANE code in %s\n' "$f" >&2; exit 1; }; done; grep -qF -- 'linux-lane' .agents/skills/roundfix/references/runtime.md || { echo 'missing linux-lane in the runtime reference' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- 'linux-lane-skipped: <reason>' || { echo 'missing warning in the deliver reference' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Lane Environment Skip**:' || { echo 'missing glossary term' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/cli || exit 1; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the doctor tests do not exist, so the command fails at its first pass check. Afterwards the check reports every case, the line appears only with the lane on, existing doctor and lane tests pass unchanged, the guide and the references name the four codes, the `runtime` reference names the check, the `deliver` reference names the warning, the glossary defines the term, the CLI builds for Windows, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 2; Core Feature 5; Success Metric 3; Glossary
- `_techspec.md` → API Contract 4; API Contract 9; Invariant 3; Interfaces; Current behavior; Testing Approach; Build Order 2; Glossary
- ADR-0303; ADR-0294; ADR-0288; ADR-0259; ADR-0271; ADR-0285; ADR-0189; ADR-0256; ADR-0280


## Result

Implemented this Task's slice for Daemon verification; status and settlement
remain Daemon-owned. No declared Verification command, repository Verification,
commit, push, or pull request operation was run.

### Implementation and acceptance evidence

- Image absent with the lane enabled: the new `linuxLane` dependency returns
  `DR-LANE-IMAGE-MISSING` through `readinessResult`, with status `warn` and
  `docker pull <image>` as its next action. The new check test asserts the
  entire rendered line and exactly `info` then `image inspect <image>`;
  the command test proves this warning preserves exit `0` and sits between
  `storage:` and `codex:`.
- Missing Docker, unreachable daemon, and healthy setup: the check test
  asserts each complete line and the recorded argument sequence, including
  no image inspection after a failed `info`. Further cases cover a configured
  image without a Go directive, no Go directive without an image, no module,
  and no repository. Dirty `go.mod` never controls the selected image.
  Each Docker read cancels its child after ten seconds. The check never
  starts, pulls, or removes a container/image or prints the environment.
- Lane disabled: the command test compares the whole output against the
  enabled case with just the lane line removed, asserts no check invocation
  while disabled, and covers an unset dependency. Existing doctor tests keep
  their names and assertions. Both lane and doctor call `linuxLaneImage`;
  the former's skip texts and ordering are preserved.
- Documentation and glossary: the doctor guide and canonical setup reference
  list all four warning codes, next actions, image selection, conditional line
  and read-only bounds. Runtime lists the conditional check; deliver names
  expected skips, the durable warning, the exact pull request paragraph,
  and doctor's diagnostic surface. `CONTEXT.md` defines Lane Environment Skip
  and links ADR-0303, following domain-modeling's glossary format.

### Focused checks and required regeneration

All commands below exited `0`, except the initial glossary check noted below.

- `go test ./internal/cli -run '^(TestDoctorLinuxLaneCheckReportsDockerAndTheImage|TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn)$' -count=1`:
  the two new tests passed using compiled stand-ins and temporary repositories.
- `go test ./internal/cli -run '^(TestDoctor|TestRunDoctor|TestLinuxLane)' -count=1`:
  focused doctor and Linux Lane consumers passed.
- `GOOS=windows go test -c -buildvcs=false -o /tmp/roundfix-task02-cli-windows.test.exe ./internal/cli`:
  the CLI and its tests compile for Windows; no Windows execution claimed.
- Raised both canonical Roundfix skill version fields from `0.1.118` to
  `0.1.119`, then ran `make skills-sync` and
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
  The mirror and version record are generated, not hand-edited.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`:
  generated the doctor guide's Behavior Surface Record fingerprint.
- `make baseline-digests`: repository-required sanctioned regeneration;
  reported `ok: true`, `changed: false`.
- `go test -tags docscontract ./internal/docscontract -run '^(TestUserGuideLinksResolve|TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce)$' -count=1`:
  first reported the new definition's undefined capitalized phrase
  `Pull Request`. Using the glossary's existing lowercase `pull request`
  wording resolved it; the focused checks then passed.
- Direct byte comparison: all four edited canonical skill files match their
  generated mirrors; every `### QA settlement` section is unchanged.
- `git diff --check`: no whitespace errors.

### Scope and search notes

The initial worktree had only this Task file modified (Daemon-owned status).
Only the declared Task slice and its generated records were changed; Task
Graph and other Task files were left untouched. The full declared Verification
and repository gates remain pending with the Daemon. No real Docker, network,
or operator Roundfix home was used by the checks.

Searched `docs/user-guide`, `docs/agents`, `.agents/skills/roundfix` and
`skills/roundfix` for doctor check descriptions and skip/notice wording.
Updated the exhaustive lists in the declared setup/runtime references and
synchronized their mirrors. Other hits outside this Task's files:
`docs/user-guide/commands/setup.md` describes only Setup's five shared readiness
checks (still accurate); `docs/user-guide/usage.md` has a narrower readiness
summary and an existing blanket offline claim; and
`docs/user-guide/commands/baseline.md` has an existing blanket offline claim.
Those unrelated summaries remain follow-up notes. The expected-skip-only
notice text in the prior Task's deliver guide agrees with this reference.

Execution environment: the sandbox could not start commands because the
assigned writable root contains a symlink. Commands ran with authorized
escalation in the resolved assigned Run checkout; this is an environment
constraint, not a product failure.
