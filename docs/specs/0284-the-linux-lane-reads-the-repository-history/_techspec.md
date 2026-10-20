---
spec: 0284-the-linux-lane-reads-the-repository-history
prd: _prd.md
created: 2026-10-19
---

# The Linux Lane reads the repository history — Technical Spec

## Executive Summary

One slice extends the existing lane builder, `runLinuxLane` in
`internal/cli/deliver_linux_lane.go`, which the Delivery Queue gate and
`roundfix deliver lane` already share (ADR-0297). No package, command, flag
or configuration key is new.

1. The lane resolves on the host the object store Git reads in its working
   directory and mounts it read-only at `/repo-objects`.
2. The lane script checks that it can read `/repo-objects` and registers it as
   an alternate of the container repository right after `git init`.
3. A store whose path holds a colon, that is not a directory, or that borrows
   objects through `info/alternates` skips the lane with its reason before
   Docker is called.

The trade-off accepted: an owner whose checkout borrows objects from another
store, a clone made with `--shared` or `--reference`, gets a skip instead of a
lane run until the store is self-contained. Mounting borrowed stores would
need host paths as container paths, which can shadow container directories
and cannot express relative entries or Windows paths. Roundfix's own
worktrees never borrow.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; `/repo-objects` and the skip reasons are fixed strings. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the lane keeps `--network none` and
  `--pull never`, and every test uses compiled stand-ins on temporary
  repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the Build Order follows recorded
  decisions. ADR-0297 governs Build Order 1: "The Linux Lane reads the
  repository history read-only". ADR-0294: "The Linux Lane runs like CI",
  which ADR-0297 amends. ADR-0288: "The Delivery Queue gate runs the changed
  Go packages in a Linux Lane", which ADR-0297 amends. ADR-0125: "A spawned
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
  qualifies, and a byte invariant compares with the Delivery Base".
  ADR-0265: "A Spec proves what it claims, and every active Spec is held to
  it". ADR-0183: "An attributed claim carries a receipt the check proves".
  Set no obligation here, because this Spec changes none of their rules:
  ADR-0092, ADR-0120, ADR-0165, ADR-0167, ADR-0169, ADR-0184, ADR-0196,
  ADR-0197, ADR-0229, ADR-0233, ADR-0237, ADR-0247, ADR-0248, ADR-0250,
  ADR-0251, ADR-0254, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269,
  ADR-0272, ADR-0277, ADR-0278, ADR-0279, ADR-0281, ADR-0282, ADR-0286,
  ADR-0290, ADR-0295 and ADR-0296, each named with its title in
  `_prd.md` → Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ grant of 2026-10-19 ("Mexer em Governed Paths declarados") and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0284-the-linux-lane-reads-the-repository-history/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-19 at `dc7e5f72`.

- The lane (`internal/cli/deliver_linux_lane.go`).
  - `linuxLaneScript` (`:24`) unpacks the archive into `/tmp/work`, runs `git
    init`, `git add -A` and one commit, then `go test`. The container
    repository holds that one commit and nothing else.
  - The `docker run` arguments (`:150-152`) carry one `-v`, the module cache
    `<GOMODCACHE>:/go/pkg/mod:ro`, then the image and `sh -ceu <script>
    linux-lane <packages>`.
  - `runLinuxLane` (`:38`) is the one builder. The gate calls it from
    `Gate` (`internal/cli/deliver_workflow.go:1497`) with the item worktree
    as `WorkDir`; `runDeliverLane` (`internal/cli/deliver_lane.go:62`) calls
    it with the checkout. `All` selects `./...` (`:90-91`); otherwise the
    changed packages are selected (`:92-114`).
  - The skips run in this order: `no go.mod at the repository root`, `go.mod
    declares no go version`, `no changed Go package` (`:115-117`), `go is not
    on PATH` (`:119-122`), `docker is not on PATH`, `docker info failed`,
    `image <image> is not present locally`, and exit `125`.
  - `TestLinuxLaneRunsTheChangedPackagesInAContainer`
    (`internal/cli/deliver_linux_lane_test.go:132`) pins every argument and
    an independent copy of the script (`:168-169`).
- The failure. The operator's `roundfix deliver lane --all` on `main` at
  `dc7e5f72` (2026-10-19, 21:44 BRT) ended `Linux Lane change failure: go
  test exited 1`; only `TestOutputsForCommand`
  (`internal/baseline/derived_ownership_test.go:61`, object
  `81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4`) and
  `TestAuthorizationReaderResolvesPreservedHistoricalRecords`
  (`internal/spec/authorization_test.go:300`, tree of
  `6b8ea48725cbca13974eee0b400b3482202874f6`) failed, both reading an object
  by name. No test failed for a missing ref.
- CI. `.github/workflows/ci-verify.yml` checks out with `fetch-depth: 0`.
- Object stores. `git rev-parse --path-format=absolute --git-path objects`,
  run in a linked worktree reached through the symbolic link `/tmp`, printed
  the main repository's `/private/tmp/.../main/.git/objects`: the common
  directory's store, with symbolic links resolved. Roundfix creates
  worktrees with `git worktree add` only (`internal/worktree`), and nothing in
  `internal` writes `objects/info/alternates`. Git ignored an `info/alternates`
  holding only a `#` comment and a blank line.
- Probes with Docker 29.6.1 and `golang:1.26`, `--init --pull never --network
  none --user 65534:65534`, on a scratch module whose test runs `git show
  <first commit>:old.txt` after a later commit removed the file:
  - today's script failed the test: `fatal: path 'old.txt' does not exist`,
    `go test` exit `1`;
  - API Contract 2's script with `-v <objects>:/repo-objects:ro` passed;
  - the same script without the mount exited `3` with
    `roundfix-linux-lane: container: cannot read the repository history at
    /repo-objects`;
  - a `git clone --shared` of the module, its store mounted alone, failed the
    test and every Git call printed `error: unable to normalize alternate
    object path: <host path>`; adding its borrowed store at the same host path
    made it pass.

## System Architecture

`internal/cli` only. `deliver_linux_lane.go` gains the object store
resolution, the mount, the three skips and the new script. The gate and
`deliver lane` keep calling `runLinuxLane`, so both get the mount with no
change to `deliver_workflow.go` or `deliver_lane.go`.

## Implementation Design

### Interfaces

```go
// internal/cli
const linuxLaneHistory = "/repo-objects"

// linuxLaneObjectStore names, on the host, the object store Git reads in
// workDir: git rev-parse --path-format=absolute --git-path objects, cleaned.
func linuxLaneObjectStore(ctx context.Context, runner preflight.GitRunner, workDir string) (string, error)

// linuxLaneHistoryMount returns the -v value for objects, or the reason the
// lane skips when it cannot mount the store whole.
func linuxLaneHistoryMount(objects string) (mount, skip string, err error)
```

### Data Models

None. No file, record or configuration key changes shape.

### API Contracts

1. API Contract: the history mount.
   - After the `no changed Go package` skip and before the `go is not on PATH`
     skip, `runLinuxLane` resolves the store with `linuxLaneObjectStore` on
     the request's `WorkDir`. A Git failure returns `read Linux Lane object
     store: <error>`, an error like the lane's other Git failures.
   - `linuxLaneHistoryMount` checks, in this order: a colon in the path after
     its volume name skips `the object store <objects> cannot be mounted: its
     path holds a colon`; a path that is not a directory skips `the object
     store <objects> cannot be mounted: it is not a directory`; an
     `<objects>/info/alternates` holding a line that is neither blank nor
     starts with `#` skips `the object store <objects> borrows objects
     through info/alternates, which the Linux Lane does not mount`. An absent
     `info/alternates` is no entry; any other read error returns `read Linux
     Lane object store alternates: <error>`. Otherwise the mount is
     `<objects>:/repo-objects:ro`.
   - Each skip passes through the lane's existing `skip`, so the gate writes
     `roundfix: notice: Delivery Queue item <slug>: Linux Lane skipped:
     <reason>` and `deliver lane` writes `roundfix: notice: Linux Lane
     skipped: <reason>` and exits `3`. No Docker or Go command starts.
   - The `docker run` arguments gain `-v`, `<objects>:/repo-objects:ro`
     right after the module cache mount; nothing else in them changes.
   - Data source: the host's object store of the repository the lane runs
     in, every branch's objects included. Destination: the lane container
     only, read-only, without network, removed when it ends. The store is
     never written: the mount is `ro`, and the alternates line is written
     into the container's own `/tmp/work/.git`. The working tree and refs are
     never mounted. A missing source never reaches `-v`, so Docker never
     creates a host directory for it.
2. API Contract: the lane script becomes exactly:
   `fail() { printf 'roundfix-linux-lane: container: %s\n' "$1" >&2; exit 3; }; for tool in tar git go; do command -v "$tool" >/dev/null 2>&1 || fail "missing tool $tool"; done; [ -d /repo-objects ] && [ -r /repo-objects ] && [ -x /repo-objects ] || fail "cannot read the repository history at /repo-objects"; mkdir -p /tmp/work && cd /tmp/work || fail "cannot prepare /tmp/work"; tar -xf - || fail "cannot unpack the candidate"; { git init -q && echo /repo-objects > .git/objects/info/alternates && git add -A && git -c user.name=roundfix -c user.email=roundfix@localhost commit -q -m "Linux Lane"; } || fail "cannot record the candidate"; pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if [ -z "$pkgs" ]; then echo "no package builds on linux"; exit 0; fi && exec go test -count=1 $pkgs`
   - `cannot read the repository history at /repo-objects` is a new setup
     cause; it exits `3` and is a Lane Container Failure through the existing
     classification.
3. API Contract: everything else stays. The image, `--init`, the module cache
   mount, the environment, the package selection, the other skips and their
   order, the cap, the outcome classification, the blockers, the Park
   Classes, the gate's and `deliver lane`'s lines and exit codes, and every
   existing test name keep their behavior.

### Surface Transcripts

The transcript runs in a temporary repository whose `go.mod` declares `go
1.26.1`, whose `.roundfixrc.yml` sets `delivery.linux_lane: false`, and whose
`.git/objects/info/alternates` names another repository's object store, with
compiled Docker and Go stand-ins first on `PATH`. `<objects>` is that
repository's `.git/objects` with symbolic links resolved.

1. Surface Transcript: the on-demand lane skips a store that borrows objects.

   ```transcript
   $ roundfix deliver lane --all
   stdout:
   roundfix: notice: Linux Lane skipped: the object store <objects> borrows objects through info/alternates, which the Linux Lane does not mount
   stderr:
   exit: 3
   ```

### Invariants

```text
1. Every lane docker run mounts exactly the module cache and the object store, both read-only.
2. The object store is the one Git reads in the lane's working directory, the common directory's store for a linked worktree.
3. The container repository lists /repo-objects as its alternate before the candidate is recorded.
4. A store the lane cannot mount whole skips the lane with its reason and starts no Docker command.
5. The gate and deliver lane, in changed-package and --all modes, build the same mount.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1, Invariants 1-3 → the history mount and the script (API Contracts 1-2)
- Goal 2, Core Feature 1, Invariant 1 → the `ro` mount of the object store only (API Contract 1)
- Goal 3, Success Metric 1, Invariant 5 → `runLinuxLane`, shared by `Gate` and `runDeliverLane` (API Contract 1)
- Goal 4, Core Feature 2, Success Metric 2, Invariant 4 → `linuxLaneHistoryMount` and Surface Transcript 1 (API Contract 1)
- Core Feature 3 → the `deliver` guide, the Roundfix `deliver` reference, **Linux Lane** in `CONTEXT.md`
- Success Metric 3 → the QA gate's real lane row, or the PRD's Unreachable Acceptance
- Core Features 1-2 and API Contract 3 → the existing lane, gate and `deliver lane` tests and the QA scope audit

## Absence Evidence

1. The history mount — command: `grep -rn 'repo-objects' internal docs/user-guide .agents/skills`; result: no match at `dc7e5f72`.
2. `linuxLaneObjectStore` — command: `grep -rn 'git-path' internal/cli --include='*.go' --exclude='*_test.go'`; result: no match at `dc7e5f72`; no production code in `internal/cli`, the lane included, asks Git where its object store is.
3. `linuxLaneHistoryMount` and its skips — command: `grep -rn 'alternates' internal/cli docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md`; result: no match at `dc7e5f72`.

## ADR Alternatives

- ADR-0297: this Spec's decision; it rejects mounting borrowed stores, running without history when the store cannot be mounted, mounting the checkout or its `.git`, fetching or bundling history, and copying refs.
- ADR-0294: amended by ADR-0297 with a dated note; it takes none of ADR-0294's rejected alternatives: no repository image is built, no reaping script is added, the exit classification and the absent-image and `125` skips stay, and the image still comes from the fetched default.
- ADR-0288: amended by ADR-0297 with a dated note; its rejected "Mounting the checkout" stays rejected: only the object store is mounted, read-only, resolved through the common directory so a worktree's `.git` file is never followed inside the container. The lane still never pulls and stays opt-in.
- ADR-0125: no rejected alternative is taken; the Docker and Go stand-ins stay compiled binaries.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows and build their colon path without creating it.
- ADR-0189: no rejected alternative is taken; the Roundfix skill edit raises its recorded version.
- ADR-0256: no rejected alternative is taken; the covering skill changes with the `deliver` surface.
- ADR-0280: no rejected alternative is taken; the guide Task records its own fingerprint.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the CI workflow, the operator's run and published sources.
- ADR-0240: no rejected alternative is taken; the real lane run is declared under Unreachable Acceptance.
- ADR-0276: no rejected alternative is taken; network-free rows qualify and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0184: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0278: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0295: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0296: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

Docker, only through the existing `docker` executable on the Queue Owner's or
operator's `PATH`, and Git on the host for `rev-parse`. No test, Verification
command or QA row needs Docker, except the QA row that runs the real lane only
when it can.

## Testing Approach

Every new top-level test calls `t.Parallel()` first and compiles for Windows.
The lane tests reuse `linuxLaneTools`, `newLinuxLaneFixture` and
`invokeDeliverLane`; the expected store is `filepath.EvalSymlinks` of the
main repository's `.git/objects`.

- `internal/cli/deliver_linux_lane_test.go`:
  - declared change: `TestLinuxLaneRunsTheChangedPackagesInAContainer` pins
    the new `-v <objects>:/repo-objects:ro` pair after the module cache and
    the new script in its independent copy, keeping its name;
  - new: `TestLinuxLaneMountsTheRepositoryHistoryReadOnly`, with subtests for
    the gate run from a linked item worktree (the store is the main
    repository's), `deliver lane --all` and `deliver lane` on changed
    packages; each asserts exactly two `-v` values, the second
    `<objects>:/repo-objects:ro`, no `-v` source equal to the worktree, and
    the script's readable check and alternates line;
  - new: `TestLinuxLaneSkipsAnObjectStoreItCannotMount`, with subtests: an
    `info/alternates` entry through the gate (the notice, no Docker or Go
    start); an `info/alternates` holding only a comment and a blank line
    still mounts; and through `linuxLaneHistoryMount`, a path holding a colon
    (built, never created) and a path that is a regular file, each with its
    exact reason.
- `internal/cli/deliver_lane_test.go`: new
  `TestDeliverLaneObjectStoreSkipTranscript`, asserting Surface Transcript 1
  whole, with `<objects>` the only controlled variation.

## Build Order

1. The history mount: `linuxLaneObjectStore`, `linuxLaneHistoryMount`, the
   three skips, the mount argument, the script, their tests, the `deliver`
   guide, the Roundfix `deliver` reference, **Linux Lane** in `CONTEXT.md`,
   the skill version and the Behavior Surface Record (no dependency).
2. QA gate (depends on: 1).

No other Spec is a prerequisite. Spec 0283 may merge first and move the
Roundfix skill version and the Behavior Surface Record, so this plan is
rebased on the default branch and checked again before its Run.

## Risks & Considerations

- A store whose directories the container user `65534` cannot read, such as
  one created under a `077` umask on a Linux host, fails the readable check as
  a Lane Container Failure naming the history, or a later object read as a
  Lane Change Failure; the lane log shows which.
- `GIT_ALTERNATE_OBJECT_DIRECTORIES` in the owner's environment is not
  checked; a store borrowing that way would fail the history tests as a Lane
  Change Failure. Nothing in Roundfix sets it.
- A test that needs a ref, not an object, would still fail in the lane; none
  does today, and adding refs then is a change of its own.

## Glossary

- changes: **Linux Lane**

## Decisions

- Mount the object store read-only at `/repo-objects` and register it as the
  container repository's alternate. See ADR-0297.
- Skip a store that borrows through `info/alternates`, holds a colon or is not
  a directory. See ADR-0297.
- Resolve the store with `git rev-parse --path-format=absolute --git-path
  objects`, which follows a linked worktree to its common directory.
- No refs. See ADR-0297.
- Check placement before the executable lookups, so a store skip never calls
  Docker.
