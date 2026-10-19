---
spec: 0282-the-linux-lane-runs-like-ci
prd: _prd.md
created: 2026-10-19
---

# The Linux Lane runs like CI — Technical Spec

## Executive Summary

Four slices extend existing seams; one test helper and one subcommand are
new, and no package is.

1. External Test Tools. `internal/testfixture` gains `RequireTool`, which
   skips a test naming an absent tool, and a repository rule that refuses a
   test resolving a tool outside the Base Test Tools by a literal name. The
   rust-cli journey and the real acpx integration test use it.
2. The init and the configured image (ADR-0294). The lane's `docker run`
   gains `--init`, and Project Config gains `delivery.linux_lane_image`, read
   from the fetched default branch like `delivery.linux_lane`.
3. Container and change failures (ADR-0294). The lane script reports its own
   setup failures with exit `3`; the lane classifies each exit as a Lane
   Container Failure or a Lane Change Failure; `delivery.GateResult` carries
   the blocker, and a new `linux-lane-container` blocker parks in the
   `environment` Park Class.
4. The on-demand lane (ADR-0294). `roundfix deliver lane [--all]
   [<revision>]` runs the same lane against the checkout and exits with its
   outcome.

The trade-offs accepted:

- With the default image the rust-cli journey skips in the lane; it still
  runs in CI, whose runner carries `rustfmt`.
- The new `deliver lane` line changes the help, and so the Behavior Surface
  fingerprint, of every `deliver` subcommand.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; one Project Config key, one blocker name and one subcommand are
  fixed names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the lane keeps `--network none` and
  `--pull never`, and every test uses compiled stand-ins on temporary
  repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — each Build Order step follows recorded
  decisions. ADR-0294 governs Build Orders 1 to 4: "The Linux Lane runs like
  CI". ADR-0288: "The Delivery Queue gate runs the changed Go packages in a Linux
  Lane", which ADR-0294 amends. ADR-0125: "A spawned fixture is a
  compiled binary, never a written script". ADR-0259: "A test runs in
  parallel unless it names why it cannot". ADR-0271: "A test bounds by real
  time only what it asserts". ADR-0285: "Every test file compiles for
  Windows, and a repository contract proves it". ADR-0189: "An owned skill's
  version names its content, and its minimum is the bundle". ADR-0256: "A
  release waits for the skills that describe a changed surface". ADR-0280: "A
  guide Task records its own surface fingerprint, a PRD records who decided,
  and the glossary is held to its own rule". ADR-0104: "A Spec accepts on
  evidence it did not author". ADR-0240: "One QA partial policy, and rows a
  Run sandbox cannot reach". ADR-0276: "A row the authorization keeps offline
  qualifies, and a byte invariant compares with the Delivery Base".
  ADR-0265: "A Spec proves what it claims, and every active Spec is held to
  it". ADR-0183: "An attributed claim carries a receipt the check proves". Set no obligation here, because this Spec changes none of their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0184 (a TechSpec states a command surface as a transcript); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof).
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ grant of 2026-10-19 ("Mexer em Governed Paths declarados") and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0282-the-linux-lane-runs-like-ci/_authorization.md`; bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `docs/agents/specific-repository.md`,
  `internal/cli/baseline_release_gate_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-19 at `f382480b` (v0.91.0).

- The lane (`internal/cli/deliver_linux_lane.go`).
  - `runLinuxLane` (`:31`) reads the root `go.mod`'s `go` version (`:67-81`)
    and uses `golang:<version>` as the image (`:127`), after `docker info`
    (`:120`) and `docker image inspect` (`:128`).
  - The `docker run` arguments (`:139-142`) are `run --rm -i --pull never
    --network none --name <name> --user 65534:65534`, seven `-e` variables,
    the read-only module cache mount, the image and `sh -ceu <script>
    linux-lane <packages>`. There is no `--init`.
  - `linuxLaneScript` (`:23`) unpacks the archive, records it in a temporary
    repository and `exec`s `go test`, so `go test` becomes PID 1.
  - The cap and cancellation remove the container (`:173-187`); the cap's
    reason is `Linux Lane exceeded 20m0s; see <lane log>`. Exit `125` is a
    skip (`:193-195`). Every other non-zero exit is `Linux Lane failed: docker
    run exited <n>; see <lane log>` (`:196`), so `126` and `127`, which Docker
    reserves for itself, read as the change's failure.
  - `linuxLaneOptIn` (`:288-320`) reads `delivery.linux_lane` from
    `roundconfig.DeliveryConfigAtCommit` at the fetched default.
  - `TestLinuxLaneRunsTheChangedPackagesInAContainer`
    (`internal/cli/deliver_linux_lane_test.go:131`) pins the arguments and an
    independent copy of the script (`:167-168`);
    `TestLinuxLaneFailureFailsTheGate` (`:327`) pins exits `1`, `126` and
    `127` to the old failure reason.
- The gate and the queue.
  - `Gate` (`internal/cli/deliver_workflow.go:1456-1503`) returns
    `delivery.GateResult{Passed: false, Reason: lane.Reason}` on a lane
    failure.
  - `gateCandidate` (`internal/delivery/engine.go:1063-1080`) parks the bare
    `gate-failed` and drops the reason; `GateResult` (`:100-103`) has only
    `Passed` and `Reason`.
  - `ClassifyPark` (`internal/delivery/park_class.go:34`) cuts the blocker at
    its first colon and puts `gate-failed` in the `finding` class (`:68`).
    `TestEveryBlockerHasAParkClass` (`park_class_test.go:53`) refuses an
    exported blocker without a class.
- Project Config (`internal/config/config.go`). `Delivery` (`:125-129`) and
  `deliveryOverlay` (`:376-380`) carry `LinuxLane`; the parse refuses a
  non-boolean `delivery.linux_lane` (`:1627-1629`); the merge applies the
  project value over the user value (`:1765-1767`).
- Tests that need a tool. `runBaselineReleaseExternal`
  (`internal/cli/baseline_release_gate_test.go:1257-1276`) runs `rustfmt`
  for the rust-cli journey (`:1017`) and fails with `t.Fatalf` when it is
  absent. Run with a `PATH` of only `go`, `git` and `sh`, the subtest failed
  in 2.3 s with `run external rustfmt src/main.rs: exec: "rustfmt":
  executable file not found in $PATH`. `internal/agent/acpx_integration_unix_test.go:37`
  looks up `node` directly and skips with its own message. A sweep of every
  `_test.go` file for `exec.Command`, `exec.CommandContext` and
  `exec.LookPath` with a literal program found `git`, `go`, `sh`, `make`,
  `mkfifo`, `ps`, `cat`, `kill`, `id`, `cmd`, `node` and one literal inside a
  fixture's source text.
- The command surface. `deliverUsage` (`internal/cli/deliver.go:23-47`) lists
  six subcommands; the root usage (`internal/cli/cli.go:55`) reads `roundfix
  deliver <plan|start|status|resume|retry|stop> [<slug> ...]`, pinned by
  `internal/cli/deliver_retry_test.go:367`. Every `deliver <sub> --help`
  prints `deliverUsage`, so the seven `command deliver ...` surfaces of
  `docs/references/behavior-surfaces.json` share one fingerprint. No
  subcommand runs the lane on demand.
- Documentation. The `deliver` guide's "Linux Lane" section
  (`docs/user-guide/commands/deliver.md:604-670`) and the Roundfix reference
  (`.agents/skills/roundfix/references/deliver.md:594-660`) state the image,
  the skips and both failure texts; the Park Class tables (`deliver.md:411`,
  reference `:326`) list the `environment` blockers; the configuration guide's
  "Linux Lane" section (`docs/user-guide/configuration.md:781-797`) documents
  only `delivery.linux_lane`; **Linux Lane** is at `CONTEXT.md:692`.
- The container, probed with Docker 29.6.1 and `golang:1.26` (go 1.26.9,
  linux/arm64), each run with `--network none --user 65534:65534`:
  - An orphan that exited stayed `Z` under PID 1 without `--init`, and
    `kill(pid, 0)` returned no error; with `--init`, PID 1 was `docker-init`
    and `kill` returned `no such process`.
  - Present: `sh`, `tar`, `git`, `go`, `gofmt`, `make`, `ps`, `kill`,
    `mkfifo`, `cat`, `id`. Absent: `rustfmt`, `cargo`, `node`, `bun`, `npm`,
    `docker`.
  - Exits: `127` for a missing command and for `sh -ceu 'rustfmt --version'`,
    `126` for a non-executable command, `137` for a killed process, `125` for
    an absent image under `--pull never`; `go test` exited `1` for a failing
    test, a build failure and a panic.
  - The script of API Contract 2 passed on a one-package module, and exited
    `3` with `roundfix-linux-lane: container: missing tool tar` when `PATH`
    held only the Go toolchain, and with `roundfix-linux-lane: container:
    cannot unpack the candidate` on a non-tar input.

## System Architecture

- `internal/testfixture`: a new `external_tool.go` holds `RequireTool`; a new
  `external_tool_test.go` holds its tests and the External Test Tool rule.
- `internal/config`: `config.go` gains `Delivery.LinuxLaneImage`, its overlay
  field, its validation and its merge.
- `internal/cli`: `deliver_linux_lane.go` gains `--init`, the configured
  image, the new script and the outcome classification; `deliver_workflow.go`
  hands the blocker to the engine; a new `deliver_lane.go` implements
  `deliver lane`, dispatched from `deliver.go` and named in `cli.go`'s root
  usage.
- `internal/delivery`: `engine.go` gains `GateResult.Blocker` and
  `BlockerLinuxLaneContainer`; `park_class.go` classifies it.

## Implementation Design

### Interfaces

```go
// internal/testfixture
func RequireTool(t testing.TB, name string) string

// internal/config
type Delivery struct {
	DerivedPaths   []DerivedPathDeclaration
	ItemBinary     ItemBinaryDeclaration
	LinuxLane      bool
	LinuxLaneImage string // empty: golang:<go.mod go version>
}

// internal/delivery
const BlockerLinuxLaneContainer = "linux-lane-container"
type GateResult struct {
	Passed  bool
	Reason  string
	Blocker string // empty parks the bare gate-failed
}

// internal/cli
type linuxLaneRequest struct {
	WorkDir, Head, Image string // task_02
	Label, LogPath, Name string // task_04: "" Label is the on-demand lane
	All                  bool   // task_04: ./... instead of changed packages
}
```

### Data Models

No new file or record format. The Project Config key
`delivery.linux_lane_image` is a YAML string. A parked item's blocker is a
string as before; `linux-lane-container: <detail>` is a new value of it.

### API Contracts

1. API Contract: `delivery.linux_lane_image`.
   - A YAML string, default empty. A non-string node, or a non-empty value
     that does not match the reference pattern below, is refused with
     `parse config "<label>": delivery.linux_lane_image must be a container
     image reference, such as golang:1.26`.
   - The pattern is the distribution reference grammar without IPv6 hosts:
     an optional `domain[:port]/` (domain components
     `[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?` joined by `.`), one or more
     path components `[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*` joined by `/`, an
     optional `:` tag `[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}` and an optional `@`
     digest `[A-Za-z][A-Za-z0-9]*([-_+.][A-Za-z][A-Za-z0-9]*)*:[0-9a-fA-F]{32,}`.
     A value can therefore never begin with `-` or hold whitespace.
   - Project Config replaces User Config; an absent project key inherits the
     user value; an explicit empty string selects the default.
   - The gate reads it from the same `DeliveryConfigAtCommit` call, at the
     fetched default, that reads `delivery.linux_lane`.
2. API Contract: the lane's container.
   - The image is `delivery.linux_lane_image` when set. Otherwise it is
     `golang:<go.mod go version>`, and only then does the `go.mod declares no
     go version` skip apply.
   - The arguments become `run --rm -i --init --pull never --network none
     --name <name> --user 65534:65534`, then the seven `-e` variables, the
     mount, the image, `sh -ceu <script> linux-lane` and the packages, all
     otherwise unchanged.
   - The script becomes exactly:
     `fail() { printf 'roundfix-linux-lane: container: %s\n' "$1" >&2; exit 3; }; for tool in tar git go; do command -v "$tool" >/dev/null 2>&1 || fail "missing tool $tool"; done; mkdir -p /tmp/work && cd /tmp/work || fail "cannot prepare /tmp/work"; tar -xf - || fail "cannot unpack the candidate"; { git init -q && git add -A && git -c user.name=roundfix -c user.email=roundfix@localhost commit -q -m "Linux Lane"; } || fail "cannot record the candidate"; pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if [ -z "$pkgs" ]; then echo "no package builds on linux"; exit 0; fi && exec go test -count=1 $pkgs`
   - Every skip, the archive, the cache mount, the environment and the
     container removal keep their behavior.
3. API Contract: the outcome of a run that started.
   - A Lane Container Failure is, in this order: the cap (`the container
     exceeded 20m0s`); exit `3` (the last output line's text after
     `roundfix-linux-lane: container: `, or `the lane script exited 3 before
     go test` when no such line exists); exit `126` (`docker run exited 126:
     the image could not invoke the lane command`); exit `127` (`docker run
     exited 127: the image lacks the lane command`); an exit of `128` or more
     (`docker run exited <n>: the container process ended by signal <n-128>`).
   - Any other non-zero exit is a Lane Change Failure, with the cause `go test
     exited <n>`.
   - The gate's reason is `Linux Lane container failure: <cause>; see <lane
     log>` or `Linux Lane change failure: go test exited <n>; see <lane log>`.
     Its blocker is `linux-lane-container: <cause>; see <lane log>` or
     `gate-failed: Linux Lane change failure: go test exited <n>; see <lane
     log>`.
   - The lane log's last line is `Linux Lane container failure: <cause>`,
     `Linux Lane change failure: go test exited <n>` or `Linux Lane passed <n>
     package(s) on <image>`.
   - Cancellation still returns the context's error after removing the
     container.
4. API Contract: the queue.
   - `gateCandidate` parks `result.Blocker` when it is not empty, and the
     bare `gate-failed` otherwise, so a repository Verification failure parks
     as before.
   - `ClassifyPark` puts `linux-lane-container` in `environment` with the
     next action `repair the Linux Lane container cause in the blocker
     (Docker, the image, or delivery.linux_lane_image on the default branch),
     then run roundfix deliver retry <slug>`.
   - `gate-failed` keeps its class and its generic next action. Retry keeps
     its rules: an archived item whose head did not move returns to `gating`.
5. API Contract: External Test Tools.
   - `RequireTool` returns `exec.LookPath(name)`'s path, or calls
     `t.Skipf("requires %s, which is not on PATH", name)`.
   - The Base Test Tools are `go`, `gofmt`, `git`, `make` (this repository's
     `verification.tools`) and `sh`, `cat`, `kill`, `ps`, `id`, `mkfifo` and
     `cmd`.
   - The External Test Tool rule parses every `_test.go` file under
     `internal/`, `cmd/` and `skills/`, resolving the `os/exec` import name.
     It reports `<file>:<line>: <tool> is not a Base Test Tool; resolve it
     with testfixture.RequireTool` for each `exec.Command`,
     `exec.CommandContext` or `exec.LookPath` call whose program is a string
     literal holding no path separator and naming no Base Test Tool.
   - The rust-cli step calls `RequireTool(t, "rustfmt")` before formatting;
     the steps before it still run. The real acpx integration test calls
     `RequireTool(t, "node")` after its opt-in check.
6. API Contract: `roundfix deliver lane [--all] [<revision>]`.
   - It loads the checkout's configuration like every `deliver` subcommand,
     resolves `<revision>` (default `HEAD`) with `git rev-parse --verify
     <revision>^{commit}`, and never reads `delivery.linux_lane`, the queue
     or a Spec.
   - It runs the lane with the checkout's `delivery.linux_lane_image`, the
     log `<artifact dir>/delivery/linux-lane.log` and the container name
     `roundfix-linux-lane-<pid>`. `--all` selects `./...` and reads no merge
     base; otherwise packages are selected as the gate selects them.
   - Its lines go to standard output: `roundfix: notice: Linux Lane skipped:
     <reason>`, `roundfix: Linux Lane passed <n> package(s) on <image>`,
     `roundfix: Linux Lane change failure: go test exited <n>; see <lane log>`
     and `roundfix: Linux Lane container failure: <cause>; see <lane log>`.
   - Exit `0` passed, `1` Lane Change Failure, `3` skipped or Lane Container
     Failure, `130` interrupted (after the container's removal), and `2` with
     the `deliver lane` Preflight failed layout for an extra argument, an
     unknown flag, a revision that is not a commit, an invalid configuration
     or an unknown default branch without `--all`.
   - `deliverUsage` gains `  roundfix deliver lane [--all] [<revision>]`,
     the command line `  lane    Run the Linux Lane on demand and exit with
     its outcome` and the flag line `  --all  Test every package in the Linux
     Lane`. The root usage line becomes `  roundfix deliver
     <plan|start|status|resume|retry|stop|lane> [<slug> ...]`.
7. API Contract: everything else stays. The gate's skip and pass lines, the
   package selection, the opt-in source, the cap's length, every other
   blocker and Park Class, the repository Verification and every existing
   test name keep their behavior.

### Surface Transcripts

Each runs in a temporary repository whose `go.mod` declares `go 1.26.1` and
whose `.roundfixrc.yml` sets `delivery.linux_lane: false`, with compiled
Docker and Go stand-ins first on `PATH`; the Docker stand-in exits `0`, exits
`1`, or prints `roundfix-linux-lane: container: missing tool git` and exits
`3`. `<lane log>` is `<artifact dir>/delivery/linux-lane.log`.

1. Surface Transcript: the on-demand lane passes.

   ```transcript
   $ roundfix deliver lane --all
   stdout:
   roundfix: Linux Lane passed 1 package(s) on golang:1.26.1
   stderr:
   exit: 0
   ```

2. Surface Transcript: the on-demand lane meets a Lane Change Failure.

   ```transcript
   $ roundfix deliver lane --all
   stdout:
   roundfix: Linux Lane change failure: go test exited 1; see <lane log>
   stderr:
   exit: 1
   ```

3. Surface Transcript: the on-demand lane meets a Lane Container Failure.

   ```transcript
   $ roundfix deliver lane --all
   stdout:
   roundfix: Linux Lane container failure: missing tool git; see <lane log>
   stderr:
   exit: 3
   ```

4. Surface Transcript: the on-demand lane skips without Docker on `PATH`.

   ```transcript
   $ roundfix deliver lane --all
   stdout:
   roundfix: notice: Linux Lane skipped: docker is not on PATH
   stderr:
   exit: 3
   ```

### Invariants

```text
1. Every lane container runs with --init.
2. The gate's image comes only from the fetched default branch and User Config, never from the item branch.
3. A test never fails because an External Test Tool is absent.
4. A run's non-zero exit is exactly one of a Lane Container Failure or a Lane Change Failure, and the blocker and the lane log name which.
5. A repository Verification failure still parks the bare gate-failed.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Feature 1, Success Metric 1, Invariant 1 → `--init` (API Contract 2)
- Goal 2, User Story 2, Core Feature 1, Success Metric 1, Invariant 2 → `delivery.linux_lane_image` and the container (API Contracts 1-2)
- Goal 3, User Story 3, Core Feature 2, Success Metric 2, Invariant 3 → `RequireTool` and the External Test Tool rule (API Contract 5)
- Goal 4, User Story 4, Core Feature 3, Success Metric 3, Invariants 4-5 → the outcome classification and the queue (API Contracts 3-4)
- Goal 5, User Story 5, Core Feature 4, Success Metric 4 → `roundfix deliver lane` and Surface Transcripts 1-4 (API Contract 6)
- Core Feature 5 → the configuration and `deliver` guides, the Roundfix `deliver` reference, `docs/agents/specific-repository.md`, `CONTEXT.md`
- Success Metrics 1-4 and API Contract 7 → the named tests and the QA scope audit

## Absence Evidence

1. `RequireTool` — command: `grep -rn 'func RequireTool' internal`; result: no match at `f382480b`.
2. The External Test Tool rule — command: `grep -rln 'Base Test Tool' internal docs/agents`; result: no match at `f382480b`.
3. `delivery.linux_lane_image` — command: `grep -rn 'linux_lane_image' internal docs`; result: no match at `f382480b`.
4. `--init` in the lane — command: `grep -n '"--init"' internal/cli/deliver_linux_lane.go`; result: no match at `f382480b`.
5. The `linux-lane-container` blocker — command: `grep -rn 'linux-lane-container\|BlockerLinuxLaneContainer' internal/delivery`; result: no match at `f382480b`.
6. `GateResult.Blocker` — command: `grep -n -A4 'type GateResult struct' internal/delivery/engine.go`; result: the struct holds only `Passed` and `Reason` at `f382480b`.
7. `roundfix deliver lane` — command: `grep -rn 'deliver lane' internal/cli docs/user-guide .agents/skills`; result: no match at `f382480b`.

## ADR Alternatives

- ADR-0294: this Spec's decision; it rejects a repository image build, a reaping lane script, every non-zero exit as the change's, turning the absent-image and `125` skips into failures, an item-branch image and a `deliver start` flag.
- ADR-0288: amended by ADR-0294 with a dated note; it takes none of ADR-0288's rejected alternatives: the lane still never pulls, never mounts the checkout, runs once per candidate rather than at each Task, and stays opt-in. Its decision that any non-zero exit other than `125` fails the gate is narrowed: exits `126`, `127`, `3`, `128` and more, and the cap now fail it as a Lane Container Failure.
- ADR-0125: no rejected alternative is taken; the Docker and Go stand-ins stay compiled binaries.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first or names the process-wide state it changes.
- ADR-0271: no rejected alternative is taken; no new test bounds a wait by real time, and the cap is classified by a unit test, not by waiting for it.
- ADR-0285: no rejected alternative is taken; every new test file compiles for Windows.
- ADR-0189: no rejected alternative is taken; each Roundfix skill edit raises its recorded version.
- ADR-0256: no rejected alternative is taken; each changed surface's covering skill changes with it.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the repository's CI workflow and published sources.
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

## Integration Points

Docker, only through the existing `docker` executable on the Queue Owner's
or operator's `PATH`; no test, Verification command or QA row needs it.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, except where it sets
`PATH` and says so, and compiles for Windows. The lane tests reuse
`linuxLaneTools` and `newLinuxLaneFixture`.

- `internal/testfixture/external_tool_test.go`:
  `TestRequireToolReturnsAToolOnPath`, `TestRequireToolSkipsNamingAnAbsentTool`,
  `TestEveryExternalTestToolGoesThroughRequireTool` (the repository rule) and
  `TestTheExternalTestToolRuleReportsALiteralLookup` (a seeded tree).
- `internal/config/linux_lane_image_test.go`:
  `TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt`,
  `TestDeliveryLinuxLaneImageRefusesAnInvalidReference`.
- `internal/cli/deliver_linux_lane_test.go`:
  - declared changes: `TestLinuxLaneRunsTheChangedPackagesInAContainer` pins
    `--init` and the new script; `TestLinuxLaneFailureFailsTheGate` keeps
    exits `1` and `2` as Lane Change Failures and moves `126` and `127` out;
    `TestLinuxLaneOptInPropagatesCancellationAndAnUnreadableConfig` follows
    the opt-in's new return.
  - new: `TestLinuxLaneRunsTheConfiguredImage`,
    `TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure`,
    `TestLinuxLaneOutcomeClassification` (the cap and the signal exits as a
    unit).
- `internal/delivery/gate_blocker_test.go`: `TestGateParksWithTheBlockerItNames`;
  `internal/delivery/park_class_test.go` gains the `linux-lane-container`
  row.
- `internal/cli/deliver_lane_test.go`: `TestDeliverLanePassedTranscript`,
  `TestDeliverLaneChangeFailureTranscript`,
  `TestDeliverLaneContainerFailureTranscript`,
  `TestDeliverLaneSkipTranscript`, `TestDeliverLaneSelectsPackagesAndImage`;
  `internal/cli/deliver_retry_test.go:367` moves to the new root line.

## Build Order

1. External Test Tools: `RequireTool`, the rule, the two callers, the
   repository rule in `docs/agents/specific-repository.md` and the glossary
   (no dependency).
2. The init and the configured image: the key, `--init`, the image, their
   tests, the configuration and `deliver` guides, the Roundfix reference, the
   glossary, the skill version and the Behavior Surface Record (depends on:
   1, which changes `CONTEXT.md`).
3. Container and change failures: the script, the classification, the
   blocker, the Park Class, their tests, the guide, the reference, the
   glossary, the skill version and the record (depends on: 2, which changes
   the same lane file, guide, reference and records).
4. The on-demand lane: the subcommand, the request refactor, its tests, the
   help, the guide, the reference, the Skill Coverage Map entry, the
   glossary, the skill version and the record (depends on: 3, which changes
   the same lane file, guide, reference and records).
5. QA gate (depends on: 4).

No other Spec is a prerequisite. Specs 0281 and 0283 are authored beside this
one; whichever merges first may move the Roundfix skill version and the
Behavior Surface Record, so this plan is rebased on the default branch and
checked again before its Run.

## Risks & Considerations

- A configured image without `tar`, `git` or `go` fails every run as a Lane
  Container Failure naming the tool, which is the intended signal.
- An older configured toolchain that cannot build the module fails as a Lane
  Change Failure, because `go test` reports it; the log shows the toolchain.
- A test that resolves an External Test Tool through a helper's variable, not
  a literal, escapes the rule; the QA gate runs the whole suite under a
  `PATH` without the External Test Tools to catch it.
- Steps 2 to 4 raise the Roundfix skill version and re-record shared files,
  so they run in series.

## Glossary

- adds: **External Test Tool**
- adds: **Lane Container Failure**
- adds: **Lane Change Failure**
- changes: **Linux Lane**

## Decisions

- `--init` rather than a reaping script. See ADR-0294.
- `delivery.linux_lane_image`, from the fetched default, held to the
  reference grammar. See ADR-0294.
- Setup failures exit `3`; `3`, `126`, `127`, `128` and more and the cap are
  the container's; a `linux-lane-container` blocker in `environment`. See
  ADR-0294.
- `GateResult.Blocker` keeps blocker names in `internal/delivery` while the
  lane chooses which one applies, and leaves the repository Verification's
  park unchanged.
- `roundfix deliver lane`, exiting `0`, `1` or `3`. See ADR-0294.
