---
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
prd: _prd.md
created: 2026-10-18
---

# Linux tests run before publication, and delivered work gives back its storage — Technical Spec

## Executive Summary

Three slices extend existing seams, and no package is new.

1. The Delivery Queue's repository gate gains the **Linux Lane** (ADR-0288).
   `commandDeliveryWorkflow.Gate` runs it after the repository Verification
   passes when Project Config sets `delivery.linux_lane: true`, and this
   repository sets it.
2. Run Worktree Reconciliation widens its leftover scope to paths beneath a
   declared or recorded path and names why it keeps a dirty Run (ADR-0289).
3. The Delivery Queue's merged release and the GC Command reclaim **Item
   Binary** directories, and the GC Command guide states the retention policy
   for Roundfix Home (ADR-0290).

The trade-offs accepted:

- A gate that is about three minutes longer per candidate on a machine with
  the image, against a Linux-only defect found before the Pull Request.
- A lane that tests only the changed packages and misses a Linux-only break in
  an unchanged package.
- A wider leftover scope that releases changes beneath a declared directory,
  as ADR-0212 already releases changes to a declared file.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; Project Config gains the boolean
  `delivery.linux_lane`, the lane derives a container name from the Spec slug
  and the process ID, and the dirty reason gains a suffix. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the Linux Lane starts a container that
  reaches no network (`--network none`), pulls nothing (`--pull never`) and
  receives only `HOME`, `GOCACHE`, `GOMODCACHE`, `GOFLAGS`, `GOPROXY`,
  `GOTOOLCHAIN` and `GOSUMDB`; the docker CLI keeps the Queue Owner's
  environment and receives no new credential. No request path, endpoint or
  recipient is added, and every test uses a fake `docker`. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — each Build Order step follows recorded
  decisions. ADR-0288 governs Build Order 1: "The Delivery Queue gate runs the
  changed Go packages in a Linux Lane". ADR-0289 governs Build Order 2: "A
  delivered Spec's dirty Run is superseded under its declared directories".
  ADR-0290 governs Build Order 3: "Merged items give back their Item Binaries".
  ADR-0212: "A merged Spec supersedes its Runs' Spec-directory work and
  declared leftovers", which ADR-0289 amends. ADR-0282: "this decision adds no
  reclamation for them", which ADR-0290 amends. ADR-0053: "Ambiguous, dirty,
  and unintegrated work stays intact". ADR-0115: "A superseded Run Branch is
  discarded through Roundfix". ADR-0033: "The Run Event Journal is pruned by a
  retention window". ADR-0255: "Run Retention removes terminal Runs whole".
  ADR-0266: "worktrees share one build cache". ADR-0225: "A queue item runs the
  Roundfix binary its branch builds". ADR-0281: "A test replaces an executable
  it rebuilds instead of writing into it". ADR-0285: "Every test file compiles
  for Windows". ADR-0273: "The judge's keys never enter a Run". ADR-0252: "The
  selective gate runs the repository contracts a change makes relevant".
  ADR-0253: "Every Repository Contract Test runs on main and before a release".
  ADR-0271: "A test bounds by real time only the budget it asserts". ADR-0189:
  "A version names one content". ADR-0280: "A guide Task records its own
  fingerprint". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0276:
  "A row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "Each claim carries its proof as
  a section". ADR-0161: "A merged Spec releases its Runs on the merged head", which the
  widened scope extends to leftovers beneath declared paths. ADR-0232: "Merge
  evidence releases the Runs and item branches of a merged Spec", whose release
  now also removes the Item Binary directory. ADR-0171: "A retention prune
  reports only what it reclaimed", which the Item Binary sweep follows. Set no obligation here, because this Spec changes none of their rules:
  ADR-0170 (a Task already completed is nothing to carry); ADR-0233 (version
  regeneration at merge); ADR-0240 (the QA partial policy); ADR-0250 (module
  versions and the coverage record's platforms, which no Task moves); ADR-0278
  (retirement deletes); ADR-0283 (the Run Marker, which the container never
  receives); ADR-0286 (the Correction Review); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0172 (doctor reports reclaimable Run storage from cheap reads); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0284 (carry-forward binds a Task to the planning text it references); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a pre-PR review finding parks only after validation). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH grant of 2026-10-16 ("Governed Paths declarados") and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/storage.md`, `.roundfixrc.yml`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-18 at `2ef4f080` (v0.89.0).

- Delivery gate. `commandDeliveryWorkflow.Gate`
  (`internal/cli/deliver_workflow.go:1446`) runs
  `defaults.verification` in the item worktree through the Verification
  executor, writes `<artifact dir>/delivery/<slug>/repository-gate.log`, and
  returns `GateResult{Passed: false}` for a failed command, which
  `gateCandidate` (`internal/delivery/engine.go:1032`) parks as `gate-failed`.
  It ignores the candidate head argument. Nothing runs a test for another
  operating system, and `internal/cli/deliver_item_binary_fork_linux_test.go`
  (`//go:build linux`) compiles only on Linux.
- Project Config. `Delivery` (`internal/config/config.go:125`) holds
  `DerivedPaths` and `ItemBinary`; Project Config replaces User Config for
  `delivery.item_binary`, and the Delivery Queue reads the declarations from
  the fetched default-branch commit.
- The Queue Owner's child commands. `runRoundfix`
  (`internal/cli/deliver_workflow.go:1672`) builds its environment from
  `commandDependenciesForContext(ctx).environ()` and replaces `HOME` with
  Roundfix Home through `deliveryCommandEnvironment` (`:2161`).
- Dirty Run classification. In `inspectTerminalRunMerged`
  (`internal/worktree/worktree.go:954`), a present Run Worktree with
  `git status` output (`:1038`) is `superseded` only when the merged-head or
  delivery proof is `safe` or `superseded` and `dirtyPathsInArchivedSpec`
  holds (`:1053`). Otherwise the state is `dirty` with the bare
  `reconciliationReasonDirty` (`:1062-1063`), whichever proof failed.
  `dirtyPathsInArchivedSpec` (`internal/worktree/merged_head.go:785`) reads
  the archived Task files at the Archive Record's `source_revision` and keeps a
  path in scope only when it lies under the Spec's directories or equals a
  declared `interface:`, `creates:` or `deletes:` path or a Recorded path
  (`:853-857`). A directory declaration covers no file beneath it.
  `TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty` and
  `TestMergedHeadKeepsADirtyRunWorktree` assert the bare reason.
- The 0272 Runs. `roundfix runs list --all` and `roundfix events` show seven
  Runs between 02:27Z and 08:12Z on 2026-10-17, each started after a plan
  revision. The three kept Runs started from `123f6524`, `adc564d0` and
  `6ee8d456` and left failed task_02, task_02 and task_04 work uncommitted. The
  three released Runs started after `24aeeb4a`. Every path any revision
  declared is declared at the source revision `10a9c759`, where task_04 also
  declares `deletes: docs/history/findings` and `deletes: docs/history/backlog`
  and records 173 deleted files. The kept worktrees were removed by hand, so
  their exact dirty paths are not recoverable.
- Merged release. When an item reaches `merged`, the engine calls
  `ReleaseMergedRuns` and `RemoveItemBranch` (`internal/delivery/engine.go:714-733`)
  and records a joined failure as a cleanup warning. `ReleaseMergedRuns`
  (`internal/cli/deliver_workflow.go:721`) calls `releaseMergedSpecRuns`
  (`internal/cli/deliver_release.go:14`), which applies only `safe` or
  `superseded` classifications.
- Item Binaries. `stepExecutable` (`internal/cli/deliver_workflow.go:1571`)
  builds the item binary and `copyDeliveryItemBinary` (`:1618`) writes each
  distinct build to `<artifact dir>/delivery/<slug>/item-binary/<16 hex>/<name>`.
  Nothing removes them; the `deliver` guide and its skill reference say "this
  operation does not reclaim them".
- GC Command. `runGC` (`internal/cli/gc.go:740`) prunes Journal Retention,
  scans and removes Worktree Build Caches of ended worktrees and merged items
  (`internal/cli/worktree_build_cache.go:173`), and sweeps Run Retention, which
  keeps a Run whose worktree exists (`internal/cli/run_retention.go:62`).
  `--dry-run` opens read-only readers (`gc.go:openGCStore`). It never looks
  under `delivery/`.

## System Architecture

- Delivery (`internal/cli`). A new `deliver_linux_lane.go` holds the lane;
  `Gate` calls it. `internal/config` gains `Delivery.LinuxLane`.
- Reconciliation (`internal/worktree`). `dirtyPathsInArchivedSpec` returns the
  first out-of-scope path and covers paths beneath a declared or recorded path;
  the dirty branch of `inspectTerminalRunMerged` builds the named reason.
- Storage (`internal/cli`). `ReleaseMergedRuns` removes the merged item's Item
  Binary directory; a new `gc_item_binary.go` scans and removes Item Binary
  directories after the Worktree Build Cache section.

## Implementation Design

### Interfaces

```go
// internal/config
type Delivery struct {
	DerivedPaths []DerivedPathDeclaration
	ItemBinary   ItemBinaryDeclaration
	LinuxLane    bool // yaml: linux_lane
}

// internal/cli
type linuxLaneResult struct {
	Ran, Passed bool
	Skip, Reason, LogPath string
	Packages []string
}
func (workflow *commandDeliveryWorkflow) runLinuxLane(ctx context.Context, workDir, specSlug, head string) (linuxLaneResult, error)
```

### Data Models

No new file format. Project Config accepts `delivery.linux_lane` (boolean,
default `false`, Project Config replaces User Config). The lane log is plain
text at `<artifact dir>/delivery/<slug>/linux-lane.log`.

### API Contracts

1. API Contract: when the Linux Lane runs.
   - `Gate` calls the lane only after the repository Verification passed and
     only when the loaded `delivery.linux_lane` is `true`; with it unset or
     `false`, `Gate` starts no `docker` or `go` process and its behavior and
     output are unchanged.
   - The lane runs once per `Gate` call, on the head the engine passes.
2. API Contract: which packages it tests.
   - The base is `git merge-base <remote>/<default branch> <head>` in the item
     worktree, with the remote and default branch the Delivery Queue already
     resolves. The changed paths are `git diff --name-only -z --no-renames
     <base> <head>`.
   - A changed `go.mod` or `go.sum` selects `./...`. Otherwise each changed
     path maps to the nearest directory, from its own directory upward to the
     root, that holds a `.go` file in `git ls-tree` at the head, written `./<dir>`
     (or `.` for the root). Paths with no such directory select nothing. The
     list is sorted and unique.
3. API Contract: when it skips. Checked in this order, each skip writes
   `roundfix: notice: Delivery Queue item <slug>: Linux Lane skipped: <reason>`
   to the Queue Owner's console log and to the lane log, and the gate passes:
   - `no go.mod at the repository root`;
   - `go.mod declares no go version`;
   - `no changed Go package`;
   - `go is not on PATH` and `docker is not on PATH`, resolved on the `PATH` of
     `commandDependenciesForContext(ctx).environ()`;
   - `docker info failed: <first non-empty output line>`;
   - `image golang:<version> is not present locally`, from a failing
     `docker image inspect golang:<version>`;
   - `docker could not start the container: <first non-empty output line>`,
     when `docker run` exits `125`.
4. API Contract: how it runs. `<version>` is the `go` directive of `go.mod`
   verbatim. The host module cache is `go env GOMODCACHE` in the item worktree.
   The lane starts, with the Queue Owner's own environment (not
   `deliveryCommandEnvironment`):

   ```text
   docker run --rm -i --pull never --network none --name roundfix-linux-lane-<slug prefix>-<pid> --user 65534:65534 -e HOME=/tmp -e GOCACHE=/tmp/go-build -e GOMODCACHE=/go/pkg/mod -e GOFLAGS=-mod=readonly -e GOPROXY=off -e GOTOOLCHAIN=local -e GOSUMDB=off -v <host module cache>:/go/pkg/mod:ro golang:<version> sh -ceu '<script>' linux-lane <package>...
   ```

   Its standard input is `git archive --format=tar <head>` from the item
   worktree. The script is fixed:

   ```text
   mkdir -p /tmp/work && cd /tmp/work && tar -xf - && git init -q && git add -A && git -c user.name=roundfix -c user.email=roundfix@localhost commit -q -m "Linux Lane" && pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if [ -z "$pkgs" ]; then echo "no package builds on linux"; exit 0; fi && exec go test -count=1 $pkgs
   ```

   The lane writes the command line and the combined output to the lane log. It
   writes nothing into the item worktree.
5. API Contract: its outcome.
   - Exit 0: the gate passes, and the console log gets
     `roundfix: Delivery Queue item <slug>: Linux Lane passed <n> package(s) on golang:<version>`.
   - Exit `125`: the skip of API Contract 3.
   - Any other exit: `GateResult{Passed: false, Reason: "Linux Lane failed: docker run exited <n>; see <lane log>"}`.
   - After 20 minutes: `docker rm -f <name>` runs and the gate fails with
     `Linux Lane exceeded 20m0s; see <lane log>`. When the gate's context is
     cancelled, `docker rm -f <name>` runs and `Gate` returns the context's
     error.
6. API Contract: the dirty leftover scope. `dirtyPathsInArchivedSpec` returns
   the first dirty path, in byte order, that is out of scope, or none. A path is
   in scope when it lies under the Spec's directories, or equals or lies beneath
   (`<declared>/` prefix) a path a Task of the archived Spec declares under
   `interface:`, `creates:` or `deletes:` or records under `## Recorded paths`.
   A sibling that only shares the prefix, such as `gone-other.txt` beside a
   declared `gone`, is out of scope. Every other input and the apply-time
   re-proof are unchanged.
7. API Contract: the dirty reason. When the Spec is archived at the chosen
   merged head or at the default branch and the Run stays `dirty`:
   - if the commit proof is neither `safe` nor `superseded`, the reason is
     `Run Worktree has tracked or untracked changes; <commit proof reason>`;
   - otherwise it is `Run Worktree has tracked or untracked changes; <path> is outside the archived Spec's declared and recorded paths`.

   Both are bounded to 160 bytes by shortening the suffix's middle text, never
   the constant prefix. A Run whose Spec is not archived there keeps the bare
   constant, and an unavailable archive source keeps today's reason.
8. API Contract: the merged release. After `releaseMergedSpecRuns`,
   `ReleaseMergedRuns` removes `<artifact dir>/delivery/<slug>/item-binary/`
   when it exists. A removal failure is joined to the returned error, so it
   becomes the cleanup warning; other files under `delivery/<slug>/` stay.
9. API Contract: the GC sweep. After the Worktree Build Cache section, GC scans
   `<artifact dir>/delivery/*/item-binary` of the repository it runs in. A
   directory is reclaimable when the repository's Delivery Queue has no item
   for that slug or holds it at stage `merged`; otherwise it is kept. When at
   least one such directory exists, the dry-run prints, in this order:

   ```text
     Item Binaries reclaimable: <n>
     Item Binary bytes reclaimable: <bytes>
     Item Binaries kept: <n>
     Reclaimable Item Binaries:
       <path>
   ```

   The live run rereads the queue before each removal and prints
   `Item Binaries removed`, `Item Binary bytes reclaimed` and
   `Removed Item Binaries`, counting only what it removed. A symlinked
   `item-binary` or a path resolving outside `<artifact dir>/delivery` is kept.
   With no such directory, GC output is byte-identical to today's.
10. API Contract: everything else stays. These keep their behavior: the
    repository Verification of the gate, every park class, the item binary
    build and copy, the commit proof of reconciliation, `--discard-superseded`,
    Task Carry-Forward, Journal Retention, Run Retention, Worktree Build Caches,
    `gc compact`, `gc sanitize`, and every existing test name.

### Surface Transcripts

None. The new lines are fixed in API Contracts 3, 5, 7 and 9, and the named
tests assert them; no command gains a flag or an exit code.

### Invariants

```text
1. The Linux Lane never pulls an image, never reaches a network, and never writes into the item worktree.
2. With delivery.linux_lane unset, the gate starts no docker or go process.
3. A dirty Run is released only when every dirty path lies under the Spec's directories or at or beneath a path the archived Spec declared or recorded, and its commit proof holds.
4. A dirty Run kept while its Spec is archived names the failing proof or path.
5. An Item Binary is reclaimed only for an item that is merged or not in the repository's Delivery Queue.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Stories 1-2, Core Features 1-2, Success Metric 1 → `runLinuxLane`, `Gate`, `Delivery.LinuxLane`, `.roundfixrc.yml` (API Contracts 1-5)
- Goal 2, User Story 3, Core Feature 3, Success Metric 2 → `dirtyPathsInArchivedSpec` (API Contract 6)
- Goal 3, User Story 4, Core Feature 4, Success Metric 2 → the dirty branch of `inspectTerminalRunMerged` (API Contract 7)
- Goal 4, User Story 5, Core Features 5-6, Success Metric 3 → `ReleaseMergedRuns`, `gc_item_binary.go`, the `gc` guide's retention policy (API Contracts 8-9)
- Core Feature 7 → the `deliver`, `reconcile` and `gc` guides, the configuration guide, the Roundfix references, `CONTEXT.md`
- Success Metrics 1-3, API Contract 10 and Invariants 1-5 → the named tests and the scope audit

## Absence Evidence

1. Linux Lane — command: `grep -rn 'linux_lane\|LinuxLane\|docker' internal/cli internal/config --include='*.go'`; result: no match at `2ef4f080`.
2. Directory-scoped leftovers and the named dirty reason — command: `grep -n 'outside the archived Spec' internal/worktree/*.go`; result: no match at `2ef4f080`.
3. Item Binary reclamation — command: `grep -rn 'item-binary' internal/cli --include='*.go'`; result: at `2ef4f080` only the build log path, the copy destination and the temporary name in `deliver_workflow.go` and three tests that build copies; none removes the directory.

## ADR Alternatives

- ADR-0288: this Spec's decision; it rejects the lane at Task settlement, a repository contract test that calls `docker`, pulling a missing image, mounting the checkout, and turning the lane on for every Go repository.
- ADR-0289: this Spec's decision; it rejects releasing every dirty Run of a merged Spec, reading the start revision's declarations, comparing bytes with the delivered head, and requiring a declared directory to exist.
- ADR-0290: this Spec's decision; it rejects reclaiming Run Worktrees by age, removing older copies while an item runs, cross-repository reclamation, and removing a merged item's whole artifact directory.
- ADR-0212: amended by ADR-0289 with a dated note; the scope gains paths beneath a declared or recorded path, and no rejected alternative of ADR-0212 is taken.
- ADR-0282: amended by ADR-0290 with a dated note; the copy-and-rename rule stays, and "Removing the built file before the next build" stays rejected because reclamation waits for the merge.
- ADR-0053: no rejected alternative is taken; dirty work outside the delivered scope stays intact.
- ADR-0115: no rejected alternative is taken; reconcile stays read-only without `--apply`.
- ADR-0033: no rejected alternative is taken; Journal Retention is unchanged.
- ADR-0255: no rejected alternative is taken; Run Retention is unchanged.
- ADR-0266: no rejected alternative is taken; the Item Binary sweep follows the Worktree Build Cache sweep's ownership rule.
- ADR-0225: no rejected alternative is taken; the item binary still drives only its own steps.
- ADR-0281: no rejected alternative is taken; no test rewrites an executable in place, and fake `docker` binaries come from `testfixture.FixtureBinary`.
- ADR-0285: no rejected alternative is taken; every new test compiles for Windows.
- ADR-0273: no rejected alternative is taken; no key enters the container.
- ADR-0252: no rejected alternative is taken; no Repository Contract Test is added.
- ADR-0253: no rejected alternative is taken; the Full Contract Run is unchanged.
- ADR-0271: no rejected alternative is taken; the 20-minute cap is never waited on by a test, which cancels instead.
- ADR-0189: no rejected alternative is taken; each skill edit raises its recorded version.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0256: no rejected alternative is taken; each changed surface's skill reference changes with it.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on repository history and published documentation.
- ADR-0276: no rejected alternative is taken; network-free rows qualify and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0161: no rejected alternative is taken; release still waits for the merged head's proof.
- ADR-0232: no rejected alternative is taken; merge evidence still drives the release, which now also removes the Item Binary directory.
- ADR-0171: no rejected alternative is taken; the Item Binary sweep counts only what it removed.
- ADR-0170: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0240: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; no test changes its platform set, so the coverage record is not re-recorded.
- ADR-0278: no rejected alternative is taken; this Spec retires nothing.
- ADR-0283: no rejected alternative is taken; the container receives no Run Marker and no key.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0172: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0284: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

- Docker, through its CLI only: `docker info`, `docker image inspect`,
  `docker run` and `docker rm -f`, with the Queue Owner's environment. No API
  socket is opened by Roundfix itself.
- The host Go toolchain, for `go env GOMODCACHE` only.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows, uses
temporary repositories and homes, and never starts a real container.

- `internal/config/linux_lane_test.go`: `TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt`.
- `internal/cli/deliver_linux_lane_test.go`, with a fake `docker` and a fake
  `go` built by `testfixture.FixtureBinary` and placed on the `PATH` of the
  injected environment:
  - `TestLinuxLaneRunsTheChangedPackagesInAContainer`;
  - `TestLinuxLaneSkipsWithANamedReason`;
  - `TestLinuxLaneFailureFailsTheGate`;
  - `TestLinuxLaneCancellationRemovesTheContainer`;
  - `TestGateWithoutTheLaneStartsNoDocker`.
- `internal/worktree/merged_spec_leftover_scope_test.go`:
  - `TestMergedSpecRunWithLeftoversBeneathADeclaredDirectoryIsSuperseded`;
  - `TestMergedSpecRunNamesTheLeftoverOutsideItsScope`;
  - `TestMergedSpecRunNamesTheUnrepresentedCommit`.
- `internal/worktree/merged_spec_leftovers_test.go`: the declared break.
  `TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty` moves from the bare
  reason to the named reason of API Contract 7. `TestMergedHeadKeepsADirtyRunWorktree`
  keeps the bare reason, because its Spec is not archived.
- `internal/cli/deliver_item_binary_reclaim_test.go`: `TestMergedReleaseReclaimsTheItemBinaries`.
- `internal/cli/gc_item_binary_test.go`:
  - `TestGCDryRunListsReclaimableItemBinaries`;
  - `TestGCRemovesOnlyItemBinariesOfMergedOrUnqueuedItems`.

## Build Order

1. The Linux Lane: the config key, the lane, the gate call and its tests, the
   `deliver` and configuration guides, the Roundfix `deliver` reference, the
   glossary, the opt-in in `.roundfixrc.yml`, the skill version and the
   Behavior Surface Record (no dependency).
2. The leftover scope and the named dirty reason with their tests, the
   `reconcile` guide, the Roundfix `reconcile` reference, the glossary, the
   skill version and the record (depends on: 1, which raises the same Roundfix
   version, glossary and record).
3. Item Binary reclamation at the merged release and in the GC Command, the
   retention policy in the `gc` guide, the `deliver` guide paragraph, the
   Roundfix `storage` and `deliver` references, the glossary, the skill version
   and the record (depends on: 2, which raises the same version, glossary and
   record).
4. QA gate (depends on: 3).

No other Spec is a prerequisite.

## Risks & Considerations

- A host module cache missing a module makes the lane fail. The lane runs only
  after the repository gate built the candidate on the same machine, which
  fills that cache.
- Tests that behave differently as an unprivileged user on Linux surface in the
  lane, as they do in CI, which also runs unprivileged.
- The image `golang:1.26` was absent on the maintainer's machine on
  2026-10-18, so the lane skips there until the operator pulls it once.
- The three Tasks share the Roundfix skill version, the glossary and the
  Behavior Surface Record, so they run in series.

## Glossary

- adds: **Linux Lane**
- adds: **Item Binary**
- changes: **Run Worktree Reconciliation**
- changes: **GC Command**

## Decisions

- The Linux Lane is a step of the Delivery Queue gate, opt-in, once per
  candidate. See ADR-0288.
- Exit `125` of `docker run` skips; every other non-zero exit fails. See
  ADR-0288.
- A declared or recorded path covers the dirty paths beneath it, and a kept
  dirty Run names why. See ADR-0289.
- Item Binaries go at the merged release and through the GC Command, and the
  GC Command guide states the retention policy. See ADR-0290.
