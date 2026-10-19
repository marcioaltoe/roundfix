---
task: task_02
spec: 0282-the-linux-lane-runs-like-ci
status: completed
type: backend
complexity: medium
---

# Task 02: The Linux Lane runs under an init on the configured image

## Overview

On 2026-10-18 `TestCLITestHelpersEndWithTheirTestBinary` timed out in the
Linux Lane because no init reaped an orphaned helper, and the lane could only
run `golang:<go.mod go version>` (`_prd.md` → Sources 1). This Task adds
`--init` to the lane's `docker run` and the Project Config key
`delivery.linux_lane_image`, read from the fetched default branch. It is
verifiable on its own through config tests and the lane's compiled Docker
stand-in, which records the arguments.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-18,
   `references/2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md`
   (`_prd.md` → Sources 1), for the init and the image, by implementing API
   Contracts 1 and 2 of `_techspec.md` except the new script, which task_03
   owns (ADR-0294).
2. MUST add `Delivery.LinuxLaneImage` and its `linux_lane_image` overlay field
   to `internal/config/config.go`, with API Contract 1's validation, refusal
   text and precedence. The value comes from Project Config and User Config
   only; whoever reads it next is `docker image inspect` and `docker run`, so
   a value that could be read as a flag or split into words MUST be refused
   at parse time.
3. MUST make the gate read the image from the same fetched default-branch
   configuration that supplies `delivery.linux_lane`, never from the item
   branch or the owner's checkout. The opt-in read returns the image with the
   switch; `TestLinuxLaneOptInPropagatesCancellationAndAnUnreadableConfig`
   follows the new return and keeps its two cases.
4. MUST make `runLinuxLane` take a `linuxLaneRequest` carrying `WorkDir`,
   `Head` and `Image` (`_techspec.md` → Interfaces). With an image, it skips
   reading the `go` version and uses the image verbatim in `docker image
   inspect`, `docker run` and the passed line; without one, it keeps
   `golang:<go.mod go version>` and the `go.mod declares no go version` skip.
   An absent configured image MUST skip with `image <image> is not present
   locally`; the lane MUST NOT pull.
5. MUST add `--init` immediately after `-i` in the `docker run` arguments.
   This is a declared break of `TestLinuxLaneRunsTheChangedPackagesInAContainer`,
   whose expected arguments gain `--init` and keep its name; its script copy
   stays as it is until task_03.
6. MUST add these tests, each calling `t.Parallel()` first and compiling for
   Windows:
   - in the new `internal/config/linux_lane_image_test.go`,
     `TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt` (default,
     project, user, project over user, explicit empty, and
     `DeliveryConfigAtCommit`) and
     `TestDeliveryLinuxLaneImageRefusesAnInvalidReference` (`-v`, `golang
     1.26`, `Golang:1.26`, ` golang:1.26`, `golang:1.26 --privileged`, a
     list, a mapping and an integer are refused with the exact text, while
     `golang:1.26`, `registry.example.com/team/go-ci:1.26`,
     `localhost:5000/go-ci` and a digest reference are accepted);
   - in `internal/cli/deliver_linux_lane_test.go`,
     `TestLinuxLaneRunsTheConfiguredImage`: with
     `delivery.linux_lane_image: registry.example.com/team/go-ci:1.26` at the
     fetched default and another image on the item branch, the stand-in sees
     `image inspect` and `run` with the default's image and `--init`, and the
     console line names it; with the image absent, the lane skips with its
     name.
7. MUST document the change:
   - `docs/user-guide/configuration.md`, "Linux Lane" section: the key, its
     default, its grammar, its refusal and the fetched-default rule;
   - `docs/user-guide/commands/deliver.md` and
     `.agents/skills/roundfix/references/deliver.md`, "Linux Lane": `--init`,
     the configured image in place of the fixed `golang:<version>`, the skip
     `image <image> is not present locally` and the passed line `... on
     <image>`.

   MUST search the guides, the agent guides and the skill copies for
   `golang:<version>` and `golang:<go.mod go version>` and keep each only
   where it names the default. No text inside any `### QA settlement` section
   may change.
8. MUST revise **Linux Lane** in `CONTEXT.md` through domain-modeling: it runs
   under an init on the image `delivery.linux_lane_image` names.
9. MUST run these commands, because this Task edits user guide pages and the
   Roundfix reference. No record or mirror is hand-edited.
   Before them, raise both version fields of each edited owned skill's
   canonical `SKILL.md` (`metadata.version` and the top-level `version`)
   by one patch over the value the tree holds when this Task starts,
   because the record command refuses changed content under an unchanged
   version. The Agent runs these commands itself before handing back: the
   Daemon only runs the Verification and never regenerates a mirror or a
   record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the configuration key, its validation and its tests.
- [ ] Read the image from the fetched default and pass it in the lane request.
- [ ] Add `--init` and update the pinned arguments.
- [ ] Add the configured-image lane test.
- [ ] Document the key and the init, revise the glossary, sync and record.

## Acceptance Criteria

- [ ] Every lane `docker run` carries `--init`.
- [ ] The lane runs the image the fetched default names, never the item branch's, and an invalid value is refused at parse time.
- [ ] The guides, the reference and the glossary describe the init and the key.

## Consumers

- `roundconfig.Delivery.LinuxLaneImage` and the `delivery.linux_lane_image` key — consumers: `./internal/config`, `./internal/cli`

## Context

- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- instruction: `internal/config/delivery.go`
- instruction: `internal/config/linux_lane_test.go`
- interface: `internal/config/config.go`
- creates: `internal/config/linux_lane_image_test.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_linux_lane_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt|TestDeliveryLinuxLaneImageRefusesAnInvalidReference|TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt TestDeliveryLinuxLaneImageRefusesAnInvalidReference TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestLinuxLaneRunsTheConfiguredImage|TestLinuxLaneRunsTheChangedPackagesInAContainer|TestLinuxLaneSkipsWithANamedReason|TestLinuxLaneOptInPropagatesCancellationAndAnUnreadableConfig|TestGateWithoutTheLaneStartsNoDocker|TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLinuxLaneRunsTheConfiguredImage TestLinuxLaneRunsTheChangedPackagesInAContainer TestLinuxLaneSkipsWithANamedReason; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- '"--init"' internal/cli/deliver_linux_lane.go || { printf 'missing --init\n' >&2; exit 1; }; for f in docs/user-guide/configuration.md docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'delivery.linux_lane_image' || { printf 'missing delivery.linux_lane_image in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- '--init' || { printf 'missing --init in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Linux Lane**:' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'delivery.linux_lane_image' || { printf 'Linux Lane glossary entry not revised\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the config and lane tests pass, the lane carries `--init`, the guides, the reference and the glossary name the key, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; User Stories 1-2; Core Feature 1; Core Feature 5; Success Metric 1; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 7; Invariants 1-2; Interfaces; Current behavior; Testing Approach; Build Order 2
- ADR-0294; ADR-0288; ADR-0125; ADR-0259; ADR-0285; ADR-0189; ADR-0256; ADR-0280

## Result

Implemented this Task's init and image slice; Task status and declared
Verification remain Daemon-owned.

- Every lane container now receives `--init` immediately after `-i`.
  `TestLinuxLaneRunsTheChangedPackagesInAContainer` retains its name and
  independent script copy, and pins the new arguments. The production script
  is byte-identical to HEAD; task_03 owns its replacement.
- Added `Delivery.LinuxLaneImage` and its pointer overlay, validating YAML
  strings against API Contract 1's reference grammar before decoding. Project
  Config overrides User Config, an absent key inherits, and explicit empty
  chooses the existing default. The new config tests cover all six loading
  cases including `DeliveryConfigAtCommit`, exact refusals for the required
  invalid inputs plus null and boolean nodes, and the named accepted images.
- The opt-in returns the image from its existing fetched-default
  `DeliveryConfigAtCommit` read; the gate passes it through `linuxLaneRequest`.
  The configured-image test proves both inspect and run use that image despite
  different item-branch and owner-config images, works without a go version,
  prints the configured image on pass, and skips a missing image before run.
  The fallback, named skips, cancellation and unreadable-config cases remain
  covered by the focused checks.
- Updated the configuration guide, deliver guide, canonical Roundfix reference
  and Linux Lane glossary entry. The guide/agent-guide/skill-copy search leaves
  `golang:<version>` and `golang:<go.mod go version>` only as defaults. Raised
  both canonical Roundfix version fields from `0.1.103` to `0.1.104`, synced
  the mirror, and generated the version and Behavior Surface records. A direct
  comparison confirms the mirrors match and QA settlement sections did not
  change.

Focused implementation evidence (all final checks exited 0):

- `go test ./internal/config -run '^TestDeliveryLinuxLane' -count=1`
- `go test ./internal/cli -run '^TestLinuxLane' -count=1`
- `go test ./internal/cli -run '^TestGate(WithoutTheLane|SkipsTheLane)' -count=1`
- `GOOS=windows GOARCH=amd64 go test -c -o /tmp/roundfix-task02-config.test.exe ./internal/config`
- `GOOS=windows GOARCH=amd64 go test -c -o /tmp/roundfix-task02-cli.test.exe ./internal/cli`
- Required generation: `make skills-sync`;
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
  `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
- `make baseline-digests` exited 0 and left no derived Baseline changes.
- `git diff --check` and direct script/settlement/mirror comparisons passed.

The initial new config test build failed because `LinuxLaneImage` did not
exist. After implementation, one test expectation used the wrong existing
proposal label; it was corrected to `Project Config proposal`, and the focused
config checks passed.

Environment: the sandbox rejected process startup because the assigned
writable root contains a symlink. Commands ran with approved escalation in
the same resolved SSD worktree. Tests used compiled stand-ins, not Docker;
no image was pulled and no network service was reached.

Declared Verification and full repository Verification were not run. No
other Task or Task Graph manifest was edited; no commit, push or PR was made.

### Verification feedback repair, attempt 1

Inspected the Daemon diagnostic artifact and the guide link check. The
configuration guide's inline domain regex contained bracket/parenthesis
syntax that the link scanner treated as a local Markdown link. Replaced that
one domain-component expression with equivalent prose: ASCII letters or
digits, internal hyphens, dot-separated components and a numeric port. The
image grammar and production implementation are unchanged.

Fresh focused evidence:

- `go test -tags docscontract ./internal/docscontract -run '^TestUserGuideLinksResolve$' -count=1` exited 0.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` exited 0 and regenerated the Behavior Surface Record for the changed guide.
- `git diff --check` exited 0.

Only the configuration guide, its generated record and this Result were
changed during the repair. Task status and declared Verification remain
Daemon-owned; the full declared sequence was not rerun.

## Carry-forward provenance

- Source Run: `run_20261019T224006Z_3d81ced188078879`
- Source commit: `4a9490da81ecd664cceb961fd9ec5749f414f77b`
