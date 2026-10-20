---
task: task_01
spec: 0284-the-linux-lane-reads-the-repository-history
status: pending
type: backend
complexity: medium
---

# Task 01: The Linux Lane mounts the repository object store read-only

## Overview

On 2026-10-19 the on-demand Linux Lane run on `main` failed two tests that
read historical commits by object name, because the lane records the
candidate in a fresh repository with no history (`_prd.md` → Sources 1). This
Task mounts the object store Git reads in the lane's working directory
read-only at `/repo-objects`, registers it as the container repository's
alternate, and skips a store the lane cannot mount whole. It is verifiable on
its own through the lane's compiled Docker stand-in, which records the
arguments.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-19,
   [`references/2026-10-19-the-linux-lane-cannot-read-the-repository-history.md`](references/2026-10-19-the-linux-lane-cannot-read-the-repository-history.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1 and 2 of
   `_techspec.md` in `internal/cli/deliver_linux_lane.go` (ADR-0297), with
   `linuxLaneObjectStore` and `linuxLaneHistoryMount` as `_techspec.md` →
   Interfaces sketches them.
2. MUST keep one builder: the mount, the skips and the script live in
   `runLinuxLane`, so the Delivery Queue gate and `roundfix deliver lane`, in
   changed-package and `--all` modes, get them alike. Neither
   `deliver_workflow.go` nor `deliver_lane.go` builds a mount of its own.
3. MUST state the hazard in the code's behavior, as API Contract 1 does:
   - the data source is the host object store of the repository the lane runs
     in, resolved with `git rev-parse --path-format=absolute --git-path
     objects`, so a linked worktree reads its common directory's store;
   - its destination is the lane container only, through `-v
     <objects>:/repo-objects:ro`; the store is never mounted writable, the
     working tree and refs are never mounted, and the alternates line is
     written only inside the container's `/tmp/work/.git`;
   - a store path holding a colon, a path that is not a directory and a store
     whose `info/alternates` holds an entry each skip with API Contract 1's
     exact reason before any Docker or Go command starts, so Docker never
     creates a host directory for a missing source.
4. MUST change the script to API Contract 2's exact text. This is a declared
   break of `TestLinuxLaneRunsTheChangedPackagesInAContainer`, whose
   expected arguments gain `-v`, `<objects>:/repo-objects:ro` right after the
   module cache mount and whose independent script copy becomes API
   Contract 2's; it keeps its name, and `<objects>` is
   `filepath.EvalSymlinks` of the fixture repository's `.git/objects`.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows and ending every child it starts in `t.Cleanup`:
   - in `internal/cli/deliver_linux_lane_test.go`,
     `TestLinuxLaneMountsTheRepositoryHistoryReadOnly`, with subtests for the
     gate run from a linked item worktree made with `git worktree add` (the
     store is the main repository's), `deliver lane --all` and `deliver lane`
     on changed packages. Each asserts that `docker run` carries exactly two
     `-v` values, the second `<objects>:/repo-objects:ro`, that no `-v`
     source is the worktree, and that the script holds the readable check
     and `echo /repo-objects > .git/objects/info/alternates` after `git init
     -q`;
   - in the same file, `TestLinuxLaneSkipsAnObjectStoreItCannotMount`, with
     subtests: an `info/alternates` entry naming another repository's store,
     through the gate, writes the notice and starts no Docker or Go command;
     an `info/alternates` holding only a `#` comment and a blank line still
     mounts; and, through `linuxLaneHistoryMount`, a path holding a colon
     (built under a temporary directory, never created) and a path that is a
     regular file each return their exact reason;
   - in `internal/cli/deliver_lane_test.go`,
     `TestDeliverLaneObjectStoreSkipTranscript`, asserting Surface Transcript
     1 of `_techspec.md` whole: the one standard output line, the empty
     standard error and exit `3`, with `<objects>` the only controlled
     variation.
6. MUST document the change in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md`:
   - "Linux Lane": the read-only mount of the object store, written as `-v
     <objects>:/repo-objects:ro`, beside the module cache, the alternate registered after
     `git init`, the new setup cause `cannot read the repository history at
     /repo-objects`, and the three skip reasons in their place in the ordered
     list, after `no changed Go package`;
   - "Run the Linux Lane on demand": the command uses the same history
     mount, and a checkout that borrows objects through `info/alternates`
     skips until `git repack -a` and the removal of `objects/info/alternates`
     make its store self-contained.

   MUST search the guides, the agent guides and the skill copies for wording
   that says the lane mounts only the module cache, and change each. No text
   inside any `### QA settlement` section may change.
7. MUST revise **Linux Lane** in `CONTEXT.md` through domain-modeling: it
   reads the repository history through the object store mounted read-only,
   linking `docs/adr/0297-the-linux-lane-reads-the-repository-history-read-only.md`.
8. MUST run these commands, because this Task edits a user guide page and the
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

- [ ] Resolve the object store and decide its mount or skip.
- [ ] Add the mount argument and the new script; update the pinned arguments.
- [ ] Add the mount, skip and transcript tests.
- [ ] Document the mount and the skips, revise the glossary, sync and record.

## Acceptance Criteria

- [ ] Every lane `docker run`, from the gate or `deliver lane`, in either mode, mounts exactly the module cache and the object store, both read-only, and the container repository borrows from `/repo-objects`.
- [ ] A linked worktree mounts its common directory's store.
- [ ] A store that borrows through `info/alternates`, holds a colon or is not a directory skips with its reason and starts no Docker command, and `roundfix deliver lane` reproduces Surface Transcript 1.
- [ ] The guide, the reference and the glossary describe the mount and the skips.

## Consumers

- the Linux Lane's `docker run` arguments, script and skip reasons — consumers: `./internal/cli`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0297-the-linux-lane-reads-the-repository-history-read-only.md`
- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `internal/cli/deliver_lane.go`
- instruction: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_linux_lane_test.go`
- interface: `internal/cli/deliver_lane_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLinuxLaneMountsTheRepositoryHistoryReadOnly|TestLinuxLaneSkipsAnObjectStoreItCannotMount|TestDeliverLaneObjectStoreSkipTranscript|TestLinuxLaneRunsTheChangedPackagesInAContainer|TestLinuxLaneSkipsWithANamedReason|TestLinuxLaneRunsTheConfiguredImage|TestDeliverLaneSelectsPackagesAndImage)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLinuxLaneMountsTheRepositoryHistoryReadOnly TestLinuxLaneSkipsAnObjectStoreItCannotMount TestDeliverLaneObjectStoreSkipTranscript TestLinuxLaneRunsTheChangedPackagesInAContainer TestDeliverLaneSelectsPackagesAndImage; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- 'echo /repo-objects > .git/objects/info/alternates' internal/cli/deliver_linux_lane.go || { printf 'missing alternates line\n' >&2; exit 1; }; for phrase in ':/repo-objects:ro' 'cannot read the repository history at /repo-objects' 'borrows objects through info/alternates, which the Linux Lane does not mount' 'cannot be mounted: its path holds a colon' 'cannot be mounted: it is not a directory'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "$phrase" || { printf 'missing %s in the deliver guide\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "$phrase" || { printf 'missing %s in the deliver reference\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Linux Lane**:' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '0297-the-linux-lane-reads-the-repository-history-read-only.md' || { printf 'Linux Lane glossary entry not revised\n' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/cli || exit 1; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'Deliver|Gate|LinuxLane' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards the mount, skip and transcript tests pass with the existing lane and `deliver lane` tests, the script registers the alternate, the guide, the reference and the glossary name the mount and every skip, the CLI builds for Windows, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 2; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Invariants 1-5; Interfaces; Current behavior; Testing Approach; Build Order 1
- ADR-0297; ADR-0294; ADR-0288; ADR-0125; ADR-0259; ADR-0271; ADR-0285; ADR-0189; ADR-0256; ADR-0280
