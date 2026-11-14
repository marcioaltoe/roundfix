---
spec: 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
prd: _prd.md
created: 2026-11-14
---

# The Verification probe prepares its tree as the Daemon does — Technical Spec

## Executive Summary

`probeSpecVerifications` in `internal/cli/spec_check.go` keeps its one
disposable checkout of `HEAD` per invocation and, right before the first
command it probes, provisions that checkout with the repository's
`worktree.copy` list and Worktree Bootstrap through a new exported
`worktree.Provision`, the two steps `ProvisionItem` already runs for a
Delivery Queue item worktree and `Create` and `createTaskWithOptions` run for
a Run and a Task Worktree. A failed bootstrap ends the check with exit `1`
and a dedicated standard-error message, and no command runs. A Spec report
whose commands ran after a bootstrap names it in one text line and one JSON
key; without a bootstrap the report is unchanged. The trade-off is one
bootstrap (a dependency install, seconds with a warm cache, minutes cold) per
check that probes a command in a repository that configures one, in exchange
for verdicts that match the Daemon's (ADR-0406).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the check runs the bootstrap command the repository already
  configures for its Run Worktrees. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0406 (new), ADR-0148, which it
  amends, and ADR-0034, which it extends, apply; ADR-0384, ADR-0385,
  ADR-0056 and ADR-0315 do not (see ADR Alternatives). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BF batch 2 answers of 2026-11-14 quoted in the authorization record;
  the Governed Paths are the Roundfix skill files task_01 changes. No
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0396-the-verification-probe-prepares-its-tree-as-the-daemon-does/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `4f8a5b7d`:

1. `runSpecCheck` loads the configuration (`loadCommandConfig`) and, with
   `--run-verification`, calls `probeSpecVerifications(ctx, loaded.GitRoot,
   resolvedSpecsRoot.Path, slugs)`. A non-nil error prints through
   `printSpecCheckFailure` (`roundfix: spec check failed: <err>` and the
   usage line) and returns `exitRunFailed` (`1`) before any report is
   printed.
2. `probeSpecVerifications` creates one `speccheck.DisposableCheckout`
   (`git worktree add --detach HEAD` in a temporary directory) for every
   Spec of the invocation. It copies no file and runs no bootstrap. For each
   Task the Daemon would still run (`daemon.TaskAwaitsWork`) that is not the
   QA gate (`daemon.TaskIsQAGate`), it calls `daemon.ProbeCommands`, which
   runs each command through `daemon.ExecVerifier` (`sh -c` in the
   checkout, environment `judge.WithoutKeyVariables(os.Environ())`).
3. The Daemon's path: `worktree.Create` (Run Worktree) and
   `createTaskWithOptions` (Task Worktree) run `copyProvisionedFiles` and
   then `runBootstrap` after `git worktree add`; `ProvisionItem` runs the
   same two steps for an item worktree. `Engine.verifyTaskPreWork` then
   probes in `plan.WorkDir` through `probeCommandsForRun` with the same
   `ExecVerifier`, adding only the Run Marker and the Verification Core
   Budget. `runBootstrap` inherits the process environment, serializes
   bootstraps inside one process, bounds the command with its timeout and
   returns a `*worktree.BootstrapError` that keeps the last 4096 bytes of
   output in `Tail`.
4. In a temporary repository whose Project Config sets
   `worktree.bootstrap: "touch .deps-installed"`, the installed 0.143.4
   reported `- task_01: honest — "test -f .deps-installed" (exited non-zero
   before work)` and exited `0`; the Daemon runs that command after the
   bootstrap, where it exits zero, and refuses the Task as vacuous.
5. The text report of a probed Spec prints `Verification tree: HEAD`, the
   uncommitted-source lines, then one line per command. Its JSON
   `verification` object holds `executed`, `tree`, `commands`, `uncommitted`
   and, when not empty, `notProbed`.

## System Architecture

The change stays in two existing files. `internal/worktree/worktree.go`
gains `Provision`, and `ProvisionItem` calls it after validating its ref, so
the item worktree and the check share one function; `Create` and
`createTaskWithOptions` keep their own calls to the same two steps, because
they return a different ref on a copy failure and a bootstrap failure. In
`internal/cli/spec_check.go`, `runSpecCheck` passes a
`verificationTreePreparation` built from `loaded.Config.Worktree.Copy` and
`worktreeBootstrapSpec(loaded.Config)`, the helper `roundfix implement`
already uses, to `probeSpecVerifications`, which prepares the checkout
lazily. No package, command, flag or configuration key is added.

## Implementation Design

### Interfaces

```go
// internal/worktree
func Provision(ctx context.Context, userRoot, worktreePath string, opts ItemProvisionOptions) error
func ProvisionItem(ctx context.Context, ref ItemRef, opts ItemProvisionOptions) error // validates, then Provision

// internal/cli
type verificationTreePreparation struct {
	CopyList  []string
	Bootstrap runworktree.BootstrapSpec
}

func probeSpecVerifications(ctx context.Context, repoRoot, specsRoot string, slugs []string,
	preparation verificationTreePreparation) ([]specCheckVerificationReport, error)
func printSpecCheckBootstrapFailure(err error, bootstrapErr *runworktree.BootstrapError, stderr io.Writer)
```

### Data Models

`specCheckVerificationReport` gains `Bootstrap string`
`json:"bootstrap,omitempty"`, placed right after `Tree`. It holds the
trimmed `worktree.bootstrap` command when the Spec's report probed at least
one command and a bootstrap is configured, and is empty otherwise. No other
schema, record or stored state changes; `speccheck.SchemaVersion` stays.

### API Contracts

1. API Contract: `worktree.Provision` runs `copyProvisionedFiles(userRoot,
   worktreePath, opts.CopyList)` and, when it succeeds,
   `runBootstrap(ctx, worktreePath, opts.Bootstrap, opts.BootstrapOutput)`,
   returning the first error unchanged. `ProvisionItem` returns
   `validateItemRef`'s error or `Provision(ctx, ref.UserRoot, ref.Path,
   opts)`; its behavior does not change.
2. API Contract: `probeSpecVerifications` keeps one `DisposableCheckout` and
   adds a once-only `prepare` step. For a Task it probes, when
   `len(task.Verification) > 0`, it calls `prepare` before
   `daemon.ProbeCommands`. The first call runs `runworktree.Provision(ctx,
   repoRoot, checkoutDir, runworktree.ItemProvisionOptions{CopyList:
   preparation.CopyList, Bootstrap: preparation.Bootstrap, BootstrapOutput:
   io.Discard})`; later calls do nothing. A failure returns `prepare
   Verification tree at HEAD: %w` and no command of any Spec runs after it.
   A Task that is completed, the QA gate, a Task without commands and a Spec
   without `_tasks.md` never trigger `prepare`.
3. API Contract: when a Spec's report probes a command and
   `strings.TrimSpace(preparation.Bootstrap.Command)` is not empty,
   `report.Bootstrap` is that command. `renderSpecCheckText` prints
   `Verification tree bootstrap: %q` on the line after
   `Verification tree: HEAD` only when `Bootstrap` is not empty; JSON omits
   the key when it is empty.
4. API Contract: in `runSpecCheck`, an error from `probeSpecVerifications`
   that holds a `*runworktree.BootstrapError` (`errors.As`) goes to
   `printSpecCheckBootstrapFailure`, which writes, in order: `roundfix: spec
   check failed: <err>`; `No Verification command was probed: a Run stops at
   the same Worktree Bootstrap before any Agent work.`; and, when the error's
   `Tail` without trailing newlines is not empty, `Worktree Bootstrap output
   tail:` followed by that tail and one newline. The command returns
   `exitRunFailed` with nothing on standard output, for `--format text` and
   `--format json` alike. Every other probe error keeps
   `printSpecCheckFailure`.
5. API Contract: Text A, a new paragraph placed right after the paragraph
   that ends "If the parser cannot start, the existing command probe still
   runs." in `docs/user-guide/commands/spec.md` and in both copies of the
   Roundfix skill's `references/spec.md`, word for word:

   ```markdown
   Before the first command it probes, the check prepares the disposable checkout as the Daemon prepares a Task Worktree: it copies the `worktree.copy` files from the repository and runs the Worktree Bootstrap once, with the configured command and `worktree.bootstrap_timeout`, so a command that needs installed dependencies or a copied file gets the verdict the Daemon gives it. A check that probes no command runs no bootstrap, and one bootstrap serves every Spec of the invocation. A Spec report whose commands ran after the bootstrap prints `Verification tree bootstrap: "<command>"` after `Verification tree: HEAD`, and JSON adds `verification.bootstrap`, a key present only then. When the bootstrap fails, no command is probed: standard output stays empty, standard error names the failure, says that no Verification command was probed and shows the bootstrap's last output, and the check exits `1`. A repository without `worktree.copy` and `worktree.bootstrap` is checked as before.
   ```

6. API Contract: Text B, appended to the definition paragraph of
   **Worktree Bootstrap** in `CONTEXT.md`, after its ADR-0034 link and on
   the same line, word for word:

   ```markdown
   `roundfix spec check --run-verification` runs it once in its disposable checkout, after the same copy and before the first command it probes, so its verdicts match the Daemon's; a failure there ends the check before any command is probed ([ADR-0406](docs/adr/0406-the-authoring-verification-probe-prepares-its-checkout-as-the-daemon-prepares-a-task-worktree.md)).
   ```

7. API Contract: Text C, appended as the last paragraph of ADR-0148, and
   Text D, appended as the last paragraph of ADR-0034; each ADR's
   `updated_at` becomes the time the Task writes it. Word for word:

   ```markdown
   **Amended by [ADR-0406](0406-the-authoring-verification-probe-prepares-its-checkout-as-the-daemon-prepares-a-task-worktree.md) (2026-11-14).** The authoring caller no longer supplies a bare working directory: before its first probed command it prepares its disposable checkout with the `worktree.copy` files and the Worktree Bootstrap a Task Worktree receives, through the worktree package's provisioning function, so both callers probe a prepared tree. A failed bootstrap ends the check before any command is probed. The shared prober, its classification and the Daemon's Run bookkeeping still hold.
   ```

   ```markdown
   **Amended by [ADR-0406](0406-the-authoring-verification-probe-prepares-its-checkout-as-the-daemon-prepares-a-task-worktree.md) (2026-11-14).** The Worktree Bootstrap also runs once in the disposable checkout of `roundfix spec check --run-verification`, after `worktree.copy` and before the first command the check probes, and a failure there ends the check before any command is probed. Run and Task Worktrees, their bootstrap-failed outcome and the command's ownership of what it does still hold.
   ```

### Surface Transcripts

Both transcripts run in a repository built from
`internal/speccheck/testdata/repo` with the Spec `clean` and one pending
Task, as `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes` and
`TestSpecCheckStopsWhenTheWorktreeBootstrapFails` build it. The `...` line
stands for the fixture's `Skipped:` block, which depends on the fixture
repository and not on the probe.

1. Surface Transcript: a prepared tree. Project Config copies `local.env` and
   sets `worktree.bootstrap: "test -f local.env && touch .deps-installed"`;
   the Task's commands are `test -f .deps-installed`, `test -f local.env` and
   `test -f output-the-task-writes`.

   ```transcript
   $ roundfix spec check clean --run-verification
   stdout:
   Spec clean
   No findings. Authored Verification commands executed: 3.
   ...
   Verification tree: HEAD
   Verification tree bootstrap: "test -f local.env && touch .deps-installed"
   - task_01: vacuous — "test -f .deps-installed" (exited zero before work)
   - task_01: vacuous — "test -f local.env" (exited zero before work)
   - task_01: honest — "test -f output-the-task-writes" (exited non-zero before work)
   stderr:
   exit: 1
   ```

2. Surface Transcript: a failed bootstrap. Project Config sets
   `worktree.bootstrap: "echo bootstrap-broke >&2; exit 7"`; the Task's one
   command writes a marker file and exits `1`. The same lines appear with
   `--format json`.

   ```transcript
   $ roundfix spec check clean --run-verification
   stdout:
   stderr:
   roundfix: spec check failed: prepare Verification tree at HEAD: worktree bootstrap failed: echo bootstrap-broke >&2; exit 7: exit status 7; bootstrap work may have been applied
   No Verification command was probed: a Run stops at the same Worktree Bootstrap before any Agent work.
   Worktree Bootstrap output tail:
   bootstrap-broke
   exit: 1
   ```

### Invariants

1. A repository that configures neither `worktree.copy` nor
   `worktree.bootstrap` gets, from `spec check --run-verification`, the same
   standard output, standard error and exit code as at `4f8a5b7d`.
2. A command never runs in the check's checkout before its preparation
   succeeded, and never after it failed.
3. The preparation runs at most once per invocation and only when a command
   will be probed.
4. The check and the item worktree reach copy and bootstrap through the same
   `worktree.Provision`; the commands still run through
   `daemon.ProbeCommands` and `daemon.ExecVerifier`.

### Declared breaks

- In a repository with a Worktree Bootstrap or `worktree.copy`, a command
  that passes on the prepared tree turns from `honest` to `vacuous`, and the
  check exits `1` where it exited `0`: the defect this Spec fixes.
- Such a check runs the bootstrap command once, which can take minutes and
  can fail where the check used to succeed.
- A missing `worktree.copy` source prints `Run Worktree copy skipped missing
  <path>` on standard error, as the Daemon does.
- A Spec report whose commands ran after a bootstrap has one more text line
  and one more JSON key.

## Coverage Map

- Goal 1, Core Feature 1, Core Feature 2, Success Metric 1 → API Contracts
  1-3, Invariant 4, Surface Transcript 1,
  `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes`.
- Goal 2, Core Feature 3, Success Metric 2 → API Contracts 2 and 4,
  Invariant 2, Surface Transcript 2,
  `TestSpecCheckStopsWhenTheWorktreeBootstrapFails`.
- Goal 3, Success Metric 3 → API Contract 2, Invariant 3,
  `TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand`.
- Goal 4, Success Metric 4 → API Contract 3, Invariant 1,
  `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted` and the
  existing `SpecCheck` tests of `./internal/cli`.
- Goal 5, Core Feature 4, Success Metric 5 → API Contracts 5-7.

Deviations from the PRD: None.

## Absence Evidence

1. `worktree.Provision`, `verificationTreePreparation` and the preparation of the check's checkout — command: `grep -rn 'func Provision(\|verificationTreePreparation\|copyProvisionedFiles\|runBootstrap' internal/worktree/worktree.go internal/cli/spec_check.go internal/speccheck/checkout.go`; result: at `4f8a5b7d` only the `func copyProvisionedFiles`, `func runBootstrap` and their calls in `Create`, `ProvisionItem` and `createTaskWithOptions` in `internal/worktree/worktree.go`; nothing in `spec_check.go` or `checkout.go`.
2. The bootstrap line and key — command: `grep -rn 'Verification tree bootstrap\|json:"bootstrap' internal`; result: no match at `4f8a5b7d`.

## ADR Alternatives

- ADR-0406: the design is the decision; it takes none of the five options the
  record rejects.
- ADR-0148: it rejected a second prober for the authoring check; the design
  keeps the one prober and changes only the directory the authoring caller
  hands it, recorded by the amendment notice task_01 writes.
- ADR-0034: it records no rejected alternative; the bootstrap keeps its
  command, timeout and failure semantics and gains one more place where it
  runs, recorded by the amendment notice task_01 writes.
- ADR-0384: not applicable — it gates Daemon Verifications with the Machine
  Verification Capacity and already lists `spec check --run-verification`
  and the Worktree Bootstrap as ungated; this Spec gates nothing.
- ADR-0385: not applicable — it gives Daemon Verifications a core budget and
  already lists `spec check --run-verification` and the Worktree Bootstrap
  as receiving none; this Spec adds no budget.
- ADR-0056: not applicable — it separates Task Capacity from the Daemon's
  Verification Capacity and leaves shared bootstrap resources to ADR-0034;
  the check runs outside any Run and takes no capacity unit.
- ADR-0315: not applicable — it refuses at authoring a test selector that
  another argument of its invocation contains; the probe still classifies
  such a command as before, and this Spec changes no selector rule.

## Integration Points

- The repository's own `worktree.bootstrap` command, run with `sh -c` in the
  disposable checkout, with the process environment, under the configured
  timeout and the in-process bootstrap permit. Roundfix owns running it, not
  what it does (ADR-0034).
- Every caller of `roundfix spec check --run-verification`: authors, the
  `write-tasks` skill, the QA gate's rows and agents inside Runs.

## Testing Approach

New tests in the new file `internal/cli/spec_check_bootstrap_test.go`, with a
suite header like the other `spec_check_*_test.go` files, each calling
`t.Parallel()` first and building its repository with
`newSpecCheckVerificationWorkspace` or `newSpecCheckStatusWorkspace`, then
committing a `.roundfixrc.yml`:

- `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes`: `.gitignore`
  holds `.deps-installed` and `local.env`; an untracked `local.env` in the
  repository; Project Config `worktree.copy: [local.env]` and the bootstrap
  of Surface Transcript 1. Subtest `surface_transcript_1` asserts exit `1`,
  empty standard error, the first two lines and every line from
  `Verification tree: HEAD` on; subtest `json_names_the_bootstrap` decodes
  the JSON and asserts `verification.tree` `HEAD` and
  `verification.bootstrap` equal to the command.
- `TestSpecCheckStopsWhenTheWorktreeBootstrapFails`: the bootstrap of
  Surface Transcript 2; the command `touch <marker under t.TempDir()>; exit
  1`. Subtests `surface_transcript_2_text` and `surface_transcript_2_json`
  assert exit `1`, empty standard output, standard error byte for byte, and
  that the marker does not exist.
- `TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand`:
  the bootstrap appends `ran` to a counter file under `t.TempDir()`.
  Subtest `nothing_to_probe`: a completed Task and a pending QA gate; exit
  `0`, empty standard error, no bootstrap line, zero counter lines. Subtest
  `two_Specs_and_two_Tasks`: two pending Tasks, and a second Spec `second`
  copied from `clean` with `spec: clean` replaced by `spec: second`; `spec
  check clean second --run-verification` prints two bootstrap lines and four
  `honest` lines, and the counter holds one line.
- `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted`: no
  `.roundfixrc.yml`; commands `test -f output-the-task-writes` and `true`.
  Subtest `text` asserts exit `1`, empty standard error, the first two lines
  `Spec clean` and `No findings. Authored Verification commands executed:
  2.` and, from `Verification tree: HEAD` on, exactly that line and the two
  command lines; subtest `json` asserts that standard output ends with
  `"verification":{"executed":true,"tree":"HEAD","commands":[{"task":"task_01","command":"test -f output-the-task-writes","verdict":"honest"},{"task":"task_01","command":"true","verdict":"vacuous"}],"uncommitted":[]}}`.
  Both outputs were captured at `4f8a5b7d` and pass there.

The existing `SpecCheck` tests of `./internal/cli` and the provisioning tests
of `./internal/worktree` pass unchanged; no existing test is edited.

Sabotages, each on a temporary copy of the audited head made with `git
archive` and committed in a fresh Git repository, each observed by the named
test (all run by the author on 2026-11-14 against the prototype):

1. Never call `prepare` (replace its call with a function that returns nil) →
   `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes`
   (`surface_transcript_1`: the two commands are `honest`).
2. `Provision` runs the bootstrap without the copy →
   `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes` (the bootstrap
   fails on the missing `local.env`).
3. Ignore a `prepare` error and probe →
   `TestSpecCheckStopsWhenTheWorktreeBootstrapFails` (the marker exists and
   standard output is a report).
4. Prepare right after the checkout is created →
   `TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand`
   (`nothing_to_probe`: one counter line).
5. Reset the once-only flag for each Spec →
   `TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand`
   (`two_Specs_and_two_Tasks`: two counter lines).
6. Print the bootstrap line for every executed report →
   `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted` (`text`).
7. Drop `omitempty` from the JSON tag →
   `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted` (`json`).
8. Send the bootstrap failure through `printSpecCheckFailure` →
   `TestSpecCheckStopsWhenTheWorktreeBootstrapFails` (both subtests).

## Build Order

1. `worktree.Provision`, the lazy preparation, the report field, the failure
   message, the four tests, Texts A to D, the skill mirror, the Roundfix skill
   version and the Behavior Surface Record (task_01).
2. The QA gate (task_02) (depends on: 1).

No other Spec is required. Specs 0392 and 0393 deliver first; the shared
files are named in `_prd.md` → Prerequisites. task_01 anchors every edit on
a function name or a sentence, never on a line number, raises the skill
version with the record command on the tree it runs on, and re-records the
Behavior Surface Record as its last edit. The code and the text stay in one
Task: by the grouping criterion the text alone would not earn its own
verification cycle, and the code already requires the Behavior Surface
Record and the skill that covers the `spec check` surface.

## Risks & Considerations

- A repository whose bootstrap needs the network (a cold `bun install` or
  `uv sync`) makes the check fail where the network is denied, for example in
  an Agent's sandbox; the message names the bootstrap, and the remedy is the
  same one the Daemon needs.
- A bootstrap that writes outside its checkout (a shared database) runs once
  more per check; ADR-0034 already leaves that to the command, and
  `worktree.concurrency: 1` projects already accept one bootstrap per Run.
- The bootstrap's own output is discarded on success and only its last 4096
  bytes are shown on failure; a longer diagnosis needs the command run by
  hand in a checkout.
- The disposable checkout lives under the system temporary directory, as
  before; a bootstrap that resolves paths relative to the repository's
  parent sees a different parent than the Run Worktree's configured
  location.

## Glossary

- changes: **Worktree Bootstrap**

## Decisions

- Provision the check's checkout through `worktree.Provision`, shared with
  `ProvisionItem`. See ADR-0406.
- Prepare once, right before the first probed command. See ADR-0406.
- A failed bootstrap probes nothing and exits `1` with its own message. See
  ADR-0406.
- Discard the bootstrap output on success; show the kept tail on failure.
- Report the bootstrap in one text line and one JSON key, only when it ran
  for that Spec's commands. See ADR-0406.
- ADR-0148 and ADR-0034 receive amendment notices. See ADR-0406.
