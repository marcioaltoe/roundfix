---
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
prd: _prd.md
created: 2026-11-11
---

# An older Roundfix refuses a repository a newer Roundfix wrote — Technical Spec

## Executive Summary

The Setup Manifest gains an optional `roundfixVersion` field, the Recorded
Roundfix Version. A Baseline plan the CLI builds passes the running release as
its writer version; the plan keeps the recorded value and replaces it with the
writer's only when the plan changes a file and the writer's release is higher.
`config.Load` reads the field from the repository's Setup Manifest before it
reads Project Config and returns a typed refusal when the recorded release is
newer than the running one; every command that loads configuration fails with
it, and the Baseline commands and `roundfix init --scope project` call the same
check for the repository they write. A load option lets Doctor, the Update
Check, the Run readers and `roundfix stop` continue without Project Config and
report the refusal instead. The primary trade-off is precision against churn:
the record names the release that last changed the Baseline, not the last one
that ran, so a release with no Baseline change makes no adopter commit while an
older binary is still refused once the Baseline moves (ADR-0392).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted;
  `setup-context-driven/manifest/0.0.1` gains one optional field. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the record is a local file; no
  endpoint, credential or network call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0392 (new) records the decision.
  ADR-0103 applies unchanged: a converged repository replans with no change.
  ADR-0287 applies unchanged: the settlement contract builds the plan without a
  writer version. ADR-0371 applies unchanged: the
  Update Check changes nothing and exits `0`. ADR-0173 does not apply.
  `## ADR Alternatives` relates
  each to its rejected options. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  BC answers of 2026-11-11 quoted in the authorization record; the Governed
  Paths are the Roundfix skill files task_03 changes and the two CLI test
  files whose parity tests task_02 updates. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/setup.md`, `internal/cli/baseline_human_test.go`, `internal/cli/baseline_plan_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `013af3b4`:

1. `internal/app/version.go`: `Version` is the release number for local builds
   (`0.139.0`), stamped from the tag for releases; `app.IsDevelopmentVersion`
   is true for an empty version or one containing `dev`; `app.CompareVersions`
   orders release numbers and puts a pre-release below its release.
2. `baseline.SetupManifest` (`internal/baseline/plan.go`) has no Roundfix
   version; `parseManagedSetupManifest` (`internal/baseline/classification.go`)
   decodes it with `decodeStrictJSON`, so an unknown field makes it invalid,
   and `ResolveManifestInput` reports `incompatible` with "manifest does not
   match the supported Setup Manifest schema".
3. `BuildPlan` builds the Setup Manifest with `buildSetupManifest`, marshals
   it, assembles the postimages with `assemblePostimages`, derives
   `FileChanges` with `deriveFileChanges` and computes the Plan Digest; a plan
   whose postimages equal the preimages has no file change, which `baseline
   update` reports as `current`. `validateSetupManifest` and
   `validatePlanApplyContract` check the manifest field by field and its
   postimage bytes against `marshalSetupManifestBytes`.
4. `config.Load` (`internal/config/config.go`) applies User Config, then, inside
   a Git repository, Project Config through `applyConfigFile`, whose strict
   decode refuses an unknown key; `LoadOptions` has `WorktreeLocationOptional`
   only. `config.Init` writes Project Config at the Git root without reading
   anything in the repository.
5. `loadCommandConfig`, `loadRunReaderCommandConfig` (`internal/cli/cli.go`),
   `runDoctorCommand` (`internal/cli/doctor.go`) and `runUpdateCommand`
   (`internal/cli/update_check.go`) all reach `config.Load`; `runStopCommand`
   uses `loadCommandConfig`; the Update Check passes `io.Discard` as the load's
   warning writer.
6. The Baseline commands (`runBaselineHumanCommandWithIO`,
   `runBaselineUpdateCommandWithSkillsStage`, `runBaselinePlanCommand`,
   `runBaselineApplyCommand`, `runBaselineSkillsRestoreCommandWith`,
   `runBaselineSkillsReconcileCommand`) take `--repo` and load no
   configuration; the plan builders in `baseline_update.go`, `baseline_human.go`
   (two) and `baseline_profile.go` pass no Roundfix version to `BuildPlan`.
7. `store.SchemaVersionError` refuses a newer Run Database with "a newer
   roundfix wrote it, so use that binary or run 'roundfix upgrade'".

## System Architecture

- `internal/config/recorded_version.go` (new): the reader, the release-number
  test, the typed refusal and the check that `Load`, `Init` and the Baseline
  commands share.
- `internal/config/config.go`: the check in `Load` and `Init`, one
  `LoadOptions` field and one `Loaded` field.
- `internal/cli/cli.go`, `doctor.go`, `health.go`, `update_check.go`: the
  recovery loader, the warning and the Doctor line.
- `internal/baseline/plan.go`, `classification.go`: the field, the writer rule
  and its validation.
- `internal/cli/baseline_version.go` (new) and the six Baseline command files:
  the Baseline guard and the writer version.

## Implementation Design

### Interfaces

```go
// package config
const SetupManifestPath = "docs/agents/setup-context.json"
func IsReleaseVersion(version string) bool
type NewerRepositoryError struct{ Path, Recorded, Running string }
func (err *NewerRepositoryError) Error() string
func RecordedRoundfixVersion(gitRoot string) (string, error)
func CheckRecordedRoundfixVersion(gitRoot string, running string) error
// LoadOptions gains: AllowNewerRepository bool
// Loaded gains:      RepositoryVersionErr error

// package baseline
// SetupManifest gains: RoundfixVersion string `json:"roundfixVersion,omitempty"` (declared last)
// PlanRequest gains:   WriterVersion string

// package cli
func loadRecoveryCommandConfig(environment commandEnvironment, stderr io.Writer) (roundconfig.Loaded, error)
func warnRepositoryVersion(loaded roundconfig.Loaded, stderr io.Writer) error
func refuseNewerBaselineRepository(repo string) error
// health.go gains: HealthCheckBaselineVersion = "baseline-version"
```

### Data Models

The Setup Manifest's optional `roundfixVersion` is a JSON string matching
`^[0-9]+\.[0-9]+\.[0-9]+$`, serialized after `verification` and omitted when
empty. No schema or manifest version changes; no stored record other than the
Setup Manifest changes.

### API Contracts

1. API Contract: reader. `RecordedRoundfixVersion(gitRoot)` returns `""`, nil
   for a blank root, a Setup Manifest at `<gitRoot>/docs/agents/setup-context.json`
   that is absent or unreadable, content that is not a JSON object, and an
   object without `roundfixVersion`. A present value that is not a JSON string
   matching the release pattern returns an error whose text is word for word
   `Setup Manifest "<path>" records roundfixVersion <raw JSON value>, which is not a Roundfix release version such as 0.140.0; restore the value the Roundfix that wrote it recorded`.
   A valid value returns it. It writes nothing.
2. API Contract: check. `CheckRecordedRoundfixVersion(gitRoot, running)`
   returns nil when `app.IsDevelopmentVersion(running)`, without reading;
   otherwise the reader's error; nil for no value; and
   `&NewerRepositoryError{Path, Recorded, Running: app.NormalizeVersion(running)}`
   when `app.CompareVersions(recorded, running) > 0`. Its `Error()` is word
   for word
   `Setup Manifest "<path>" records Roundfix <recorded>, newer than this binary's <running>; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'`,
   with `<path>` the absolute Setup Manifest path.
3. API Contract: load. Inside a Git repository `Load` calls
   `CheckRecordedRoundfixVersion(gitRoot, app.Version)` after User Config and
   before Project Config. A non-nil result is returned as `Load`'s error,
   unwrapped, and Project Config is not read; with `AllowNewerRepository` it is
   stored in `Loaded.RepositoryVersionErr` instead, `Load` succeeds and
   Project Config is not read. `ProjectConfigPath` is set either way. `Init`
   with scope `project` returns the same error before it writes.
4. API Contract: commands that refuse. Every command that loads configuration
   through `loadCommandConfig` or `environment.loadOptions` and `Load`, other
   than those of API Contract 5, fails with the refusal through its existing
   failure printer and exit code (for example `roundfix: spec check failed:
   <refusal>` and exit `2`). `roundfix init --scope project` prints its
   `Init failed` block with the refusal as the reason, exits `2` and leaves
   Project Config unchanged. `roundfix --version`, `help`, `roundfix upgrade`'s
   install step and `roundfix migrate` do not load repository configuration and
   are unchanged.
5. API Contract: commands that warn. `loadRecoveryCommandConfig` is
   `loadCommandConfig` with `AllowNewerRepository` set and no warning of its
   own. `warnRepositoryVersion` writes, only when `RepositoryVersionErr` is
   set, one stderr line, word for word
   `roundfix: warning: <refusal>; Project Config not read`.
   - `loadRunReaderCommandConfig` sets `AllowNewerRepository` and calls
     `warnRepositoryVersion` with its stderr before its worktree-location
     warning (`runs`, `runs show`, `attach`, `events`, the Run Browser).
   - `runStopCommand` loads with `loadRecoveryCommandConfig(environment,
     stderr)` and calls `warnRepositoryVersion(loaded, stderr)`.
   - `runUpdateCommand` loads with `loadRecoveryCommandConfig(environment,
     io.Discard)`, as today without warnings, and then calls
     `warnRepositoryVersion(loaded, stderr)`; its report and exit code follow
     ADR-0371.
   - `runDoctorCommand` sets `AllowNewerRepository` on its load options and,
     when `RepositoryVersionErr` is set, puts first among its results
     `CheckResult{Name: HealthCheckBaselineVersion, Status: failed, Detail:
     <refusal> + "; Project Config not read", NextAction: "roundfix
     upgrade"}`, printed as
     `baseline-version: failed (<refusal>; Project Config not read; next: roundfix upgrade)`;
     every other check runs, and Doctor exits `1`. Without the refusal no line
     is added.
6. API Contract: writer. `BuildPlan` sets the built manifest's
   `RoundfixVersion` to the value the repository's current Setup Manifest
   records, when it is a release number (absent, unreadable or malformed:
   empty). After `FileChanges` is derived, when it is not empty and
   `app.NormalizeVersion(request.WriterVersion)` is a release number above the
   carried value (or the carried value is empty), the plan sets that release,
   re-marshals the manifest, re-assembles the postimages and the ledger with
   `assemblePostimages`, derives `FileChanges` again and then computes the
   Plan Digest. A development, pre-release, malformed or empty writer version
   and a writer at or below the carried value change nothing.
7. API Contract: validation. `parseManagedSetupManifest` and
   `validateSetupManifest` accept the field when it is empty or a release
   number and reject any other value (`validateSetupManifest` with
   `Setup Manifest roundfixVersion must be a release version such as 0.140.0`).
8. API Contract: CLI writer. The plan builders of `roundfix baseline update`,
   `roundfix baseline plan` and both of `roundfix baseline` pass
   `WriterVersion: app.Version`. `RecalculatePlan` passes the request it is
   given unchanged. `releaseplan_checks.go` passes none.
9. API Contract: Baseline guard. `refuseNewerBaselineRepository(repo)`
   resolves `repo` to an absolute path, finds its Git root with
   `config.FindGitRoot`, returns nil without one, and otherwise returns
   `config.CheckRecordedRoundfixVersion(root, app.Version)`. Each Baseline
   command calls it right after its arguments parse, before it loads the
   catalog, reads the manifest, a plan file or a skill source, and on an error:
   - `baseline update` writes its failure result with state `failed`,
     category `version`, the refusal as message and next action
     `run roundfix upgrade, then rerun roundfix baseline update`, exit `2`;
   - `roundfix baseline` writes `baselineHumanActionResult("preflight",
     <refusal>, "run roundfix upgrade, then rerun roundfix baseline")` through
     `writeBaselineHumanAction`, exit `2`;
   - `baseline plan` and `baseline apply` print their usual failure with the
     refusal, exit `2`;
   - `baseline skills restore` and `baseline skills reconcile` write their
     failure payload with finding code `restore.roundfix-version-newer` or
     `reconcile.roundfix-version-newer`, the refusal as message and action
     `Run roundfix upgrade, then rerun roundfix baseline skills restore.` or
     `... reconcile.`, exit `2`.
   No refused command writes a file. The Update Check's in-process preview
   reaches the `baseline update` refusal, so its Baseline line reads
   `Baseline: failed: <refusal>` and its JSON `baseline.state` is `failed`.
10. API Contract: guides. task_01 adds the Doctor and upgrade text and
    task_02 the Baseline and Update Check text, word for word where quoted:
    - `docs/user-guide/commands/baseline.md`, right before the heading
      `#### Capability evidence and remediation`, a heading `#### Recorded
      Roundfix Version` and a paragraph: what the field records, when a plan
      writes it, that a plan without a file change keeps it, that development
      and pre-release builds write none, and that the six Baseline commands
      refuse a newer value before they read anything else, with the
      three-line `text` block task_02 prescribes, with `<recorded>`,
      `<running>` and `<path>` placeholders;
    - `docs/user-guide/commands/doctor.md`, right before the heading
      `## Worktree Build Caches`, the sentence "When a newer Roundfix wrote
      this repository's Baseline, Doctor prints this line first, runs every
      other check without Project Config, and exits `1`:" and a `text` block
      holding the Doctor line of API Contract 5 with placeholders;
    - `docs/user-guide/commands/update.md`, right before the paragraph that
      begins "Only when a release is newer or the Baseline has a plan", the
      sentence "When a newer Roundfix wrote this repository's Baseline, the
      check prints `roundfix: warning: <refusal>; Project Config not read` on
      stderr, reads no Project Config, and its Baseline line is `Baseline:
      failed:` followed by the refusal.";
    - `docs/user-guide/commands/upgrade.md`, at its end, a paragraph that
      begins "A binary older than the Recorded Roundfix Version of a
      repository" and says that every command that loads the repository's
      configuration refuses there and names `roundfix upgrade`, that `roundfix
      upgrade` reads no configuration before it installs, and that `roundfix
      doctor`, `roundfix update check`, the Run readers and `roundfix stop`
      warn and continue without Project Config.
11. API Contract: Roundfix references. The canonical
    `.agents/skills/roundfix/references/baseline.md` gains, right before the
    heading `## Vendored skills`, a heading `### Recorded Roundfix Version` and
    one paragraph whose first sentence is word for word "The Setup Manifest
    records as `roundfixVersion` the release of the Roundfix that last applied
    a Baseline change to the repository." and which says when it moves, that
    an older binary refuses the configuration-loading commands, `roundfix init
    --scope project` and the Baseline commands with the refusal of API
    Contract 2, and that the remedy is `roundfix upgrade`, never editing the
    field. The canonical `.agents/skills/roundfix/references/setup.md` gains,
    right after the paragraph under `## Update check` that begins
    "`roundfix update check [--refresh]", one paragraph whose first sentence
    is word for word "When a newer Roundfix wrote the repository's Baseline,
    the check warns on stderr and reports `Baseline: failed:` with the
    refusal; run `roundfix upgrade` before anything else in that repository."
    and which names the Doctor line of API Contract 5.
12. API Contract: glossary. `CONTEXT.md` gains, right after the
    **Setup Manifest** entry, word for word, with `ADR-0392` in the parentheses
    written as a Markdown link to
    `docs/adr/0392-an-older-roundfix-refuses-a-repository-whose-baseline-a-newer-roundfix-wrote.md`:

    ```text
    **Recorded Roundfix Version**:
    The release of the Roundfix that last applied a Baseline change to a repository, recorded as `roundfixVersion` in its Setup Manifest. A plan that changes a file records the running release when it is higher, and a plan that changes none keeps the value. An older Roundfix refuses, naming both releases and `roundfix upgrade`, every command that loads the repository's configuration and every Baseline command, while the Doctor Command, the Update Check, the Run readers and `roundfix stop` report it and continue without Project Config (ADR-0392).
    _Avoid_: Minimum version, required version, writer stamp
    ```

    and the **Setup Manifest** entry's definition gains, at the end of its
    first paragraph, the sentence "It also carries the Recorded Roundfix
    Version (ADR-0392)." with the same link.

### Surface Transcripts

The fixture is a temporary directory marked as a Git repository whose
`docs/agents/setup-context.json` is `{"roundfixVersion":"999999.0.0"}` and
whose `.roundfixrc.yml` holds `future_section:` with `enabled: true`; the home
is an empty temporary directory. `<refusal>` is the refusal of API Contract 2
for that repository and the running release.

1. Surface Transcript: a Run reader warns and continues.

   ```transcript
   $ roundfix runs list
   stdout:
   No Runs found.
   stderr:
   roundfix: warning: <refusal>; Project Config not read
   exit: 0
   ```

2. Surface Transcript: the Update Check of a development build names the
   refused Baseline preview.

   ```transcript
   $ roundfix update check
   stdout:
   Roundfix 0.0.0-dev is a development build; no release comparison (update.mode notify).
   Baseline: failed: <refusal>
   stderr:
   roundfix: warning: <refusal>; Project Config not read
   exit: 0
   ```

   The development build is the Update Check's installed version only; the
   refusal compares the binary's `app.Version`.

### Invariants

1. A repository without `roundfixVersion` behaves exactly as at `013af3b4` on
   every command; a binary at or above the recorded release behaves exactly
   as at `013af3b4`.
2. No refused command reads or writes Project Config, the Setup Manifest or
   any other repository file after the check.
3. No plan built without a writer version records the field, so every plan
   golden, parity fixture and this repository's Setup Manifest keep their
   bytes; a plan that changes no file records nothing new.
4. Nothing lowers a recorded value; a development or pre-release build never
   records one and is never refused.
5. The Update Check still exits `0` and changes nothing; `--version`, `help`
   and `roundfix upgrade`'s install step read no repository configuration.

### Declared breaks

- A repository whose Setup Manifest records a newer release: every
  configuration-loading command, `init --scope project` and the six Baseline
  commands refuse; Doctor exits `1` with a new first line; the Update Check,
  the Run readers and `stop` print a stderr warning and read no Project Config.
- A Setup Manifest with `roundfixVersion` is invalid for every binary before
  this release (strict decode); their Baseline commands report it as
  `incompatible`.
- A Baseline plan the CLI applies with a file change adds `roundfixVersion` to
  the Setup Manifest it writes. Three tests compare a CLI plan with one
  `BuildPlan` builds directly and differ by that field:
  `TestBaselinePlanProfileFile` (`internal/cli/baseline_plan_test.go`),
  `TestHumanAutomationPlanParity` (also run by
  `TestBaselineMacroJourneysPublicCLI`) and
  `TestBaselineHumanProfileAdaptation` (`internal/cli/baseline_human_test.go`);
  task_02 adds `WriterVersion: app.Version` to their direct requests, which
  keeps their parity contract. The prototype ran `./internal/baseline` and the
  rest of `./internal/cli` unchanged; a Task that finds another such test on
  the default branch when it runs changes it the same way and says so.

## Coverage Map

- Goal 1, User Stories 1 and 2, Core Features 2 and 4, Success Metric 1,
  Success Metric 2 → API Contracts 1 to 4 and 9, Invariant 2.
- Success Metric 6 (a malformed value refused, a development build not) → API
  Contracts 1, 2 and 7, Invariant 4.
- Goal 2, User Story 3, Core Feature 3, Success Metric 3 → API Contract 5,
  Surface Transcripts 1 and 2, Invariant 5.
- Goal 3, Core Feature 1, Success Metrics 4 and 5 → API Contracts 6 to 8,
  Invariants 3 and 4.
- Core Feature 5, Success Metric 7 → API Contracts 10 to 12.

Deviations from the PRD: None.

## Absence Evidence

1. Recorded Roundfix Version reader and refusal — command: `rg -n 'roundfixVersion|RecordedRoundfixVersion|NewerRepositoryError|AllowNewerRepository|baseline-version' internal docs/user-guide .agents/skills/roundfix`; result: no match at `013af3b4`.
2. Writer version in Baseline plans — command: `rg -n 'WriterVersion|app\.Version' internal/baseline internal/cli/baseline_*.go -g '!*_test.go'`; result: no match in `internal/baseline` and only output uses of `app.Name` in the Baseline command files at `013af3b4`; no plan carries a Roundfix release.

## ADR Alternatives

- ADR-0392: the design is the decision; it takes none of the six options the
  record rejects (a version on every plan, a hand-raised minimum, a separate
  marker file, a Project Config key, refusing read-only commands, explaining
  unknown Project Config keys).
- ADR-0103: it rejected an append-only adoption record that never converges;
  the field rides only a plan that changes a file, so a second run still
  reports the repository current.
- ADR-0287: it rejected regenerating the Setup Manifest only at release; the
  settlement contract still builds the plan without a writer version, so a
  release never makes the Setup Manifest stale.
- ADR-0173: not applicable — it makes a Baseline Plan report the citations
  its History Relocations break; this Spec changes no relocation, citation or
  history unit, and the field rides the plan beside them.
- ADR-0371: it rejected an Update Check that applies anything; the check only
  reports the refusal and keeps `next` and its exit code.

## Integration Points

- An adopter's Setup Manifest gains the field at its next applied Baseline
  change, through `roundfix baseline update` or the `--update` Skill Mode.
- Binaries before this release cannot read the field and report the Setup
  Manifest as incompatible in their Baseline commands.

## Testing Approach

New tests (each calls `t.Parallel()` first):

- `internal/config/recorded_version_test.go`:
  `TestARepositoryRecordedByANewerRoundfixRefusesLoad` (a recorded
  `999999.0.0` and a Project Config with an unknown top-level key: `Load`
  returns a `*NewerRepositoryError` with the message of API Contract 2 word
  for word, not the unknown-key error),
  `TestAnEqualOrOlderRecordedRoundfixVersionLoadsProjectConfig` (the running
  release and `0.0.1`: Project Config's `update.mode: auto` is read),
  `TestARepositoryWithoutARecordedRoundfixVersionLoads` (absent manifest, no
  field, not JSON), `TestAMalformedRecordedRoundfixVersionRefusesLoad`
  (`"0.140"`, `"v0.140.0"`, `"0.140.0-rc.1"`, `140`),
  `TestAllowNewerRepositoryLoadsWithoutProjectConfig` (subtests `valid
  project config`, a Project Config that only says `update.mode: auto`, and
  `unknown project key`, the same plus an unknown top-level key: no error,
  `RepositoryVersionErr` set, `update.mode` stays `notify`) and
  `TestADevelopmentBuildIgnoresTheRecordedRoundfixVersion` (a newer and a
  malformed value, `0.0.0-dev` passes and `0.139.0` does not). The newer
  release is a test constant above any shipped version; the running one reads
  `app.Version`.
- `internal/cli/recorded_version_test.go`:
  `TestACommandRefusesARepositoryRecordedByANewerRoundfix` (subtests `spec
  check x`, `follow-up`, `init --scope project --force` in the Surface
  Transcripts fixture: exit `2`, the refusal printed, `future_section` never
  mentioned, Project Config bytes unchanged) and
  `TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix` (subtest
  `runs list` asserts Surface Transcript 1 whole; subtest `update check`, with
  a development installed version, asserts exit `0`, the stderr warning line
  whole and the release line; subtest `doctor`, with
  `withDoctorFakeLoadedAndReadiness`, an empty profile proof and a
  `loadConfig` that fails the test unless `AllowNewerRepository` is set and
  then calls `config.Load`, asserts the first line of API Contract 5 word for
  word, a `node: ok` line after it and exit `1`).
- `internal/baseline/recorded_version_test.go`:
  `TestAPlanRecordsItsWriterVersionOnlyWithAFileChange` (on the plan
  characterization repository: an adoption plan without a writer version has
  no `roundfixVersion` in its manifest postimage; with `0.140.0` it records
  it and applies `verified`; replans by `0.141.0`, empty and `0.0.0-dev` on
  the applied repository have no file change and keep `0.140.0`; after a
  managed guide is made stale, writers `0.139.0`, `0.141.0-rc.1` and empty
  keep `0.140.0` and `v0.141.0` records `0.141.0`, each in the postimage
  bytes) and `TestASetupManifestRoundfixVersionMustBeAReleaseVersion` (the
  adoption manifest recording `0.140.0` parses; `0.140`, `v0.140.0` and
  `0.140.0-rc.1` do not).
- `internal/cli/baseline_version_test.go`:
  `TestBaselineUpdateRecordsTheRunningRoundfixVersion` (on
  `staleBaselineUpdateRepository`, whose CLI adoption already recorded a
  value: the test first rewrites its Setup Manifest with `RoundfixVersion`
  empty through `WriteBaselineSetupManifest`, commits, and asserts the field is
  absent; then `baseline update --yes` writes
  `"roundfixVersion": "<app.NormalizeVersion(app.Version)>"` and a second run
  has zero file changes),
  `TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix` (subtests
  `adopt`, `update`, `plan`, `apply`, `skills restore`, `skills reconcile`,
  each with `--repo` and `--format json` in the Surface Transcripts fixture:
  exit `2`, the JSON-escaped refusal in stdout, the repository tree unchanged)
  and `TestTheUpdateCheckNamesARepositoryRecordedByANewerRoundfix` (Surface
  Transcript 2 whole; JSON `baseline.state` `failed`, `baseline.message` the
  refusal, `next` `none`).

Sabotages, each on a temporary copy of the audited head made with `git
archive` and committed in a fresh Git repository, each observed by the named test on the named case (all run by the
author on 2026-11-11 against the prototype):

1. Remove the check from `Load` → `TestARepositoryRecordedByANewerRoundfixRefusesLoad`.
2. Refuse an equal release (`>= 0`) → `TestAnEqualOrOlderRecordedRoundfixVersionLoadsProjectConfig`
   (the running-release case).
3. Read Project Config before the check → `TestARepositoryRecordedByANewerRoundfixRefusesLoad`
   (the unknown-key error comes first).
4. Read Project Config under `AllowNewerRepository`, ignoring its error →
   `TestAllowNewerRepositoryLoadsWithoutProjectConfig/valid_project_config`.
5. Accept any string value → `TestAMalformedRecordedRoundfixVersionRefusesLoad`.
6. Drop the development-build skip → `TestADevelopmentBuildIgnoresTheRecordedRoundfixVersion`.
7. Remove the check from `Init` → `TestACommandRefusesARepositoryRecordedByANewerRoundfix/init_--scope_project_--force`.
8. Drop `AllowNewerRepository` from `loadRunReaderCommandConfig` → `TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix/runs_list`.
9. Drop `AllowNewerRepository` from Doctor's load options → `TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix/doctor`.
10. Record the writer's release without a file change → `TestAPlanRecordsItsWriterVersionOnlyWithAFileChange`
    (the `0.141.0` replan of the applied repository).
11. Record a writer below the carried value → the same test (writer `0.139.0`).
12. Stop carrying the recorded value → the same test (the empty-writer replan
    proposes a change).
13. Drop the release-number test from `parseManagedSetupManifest` → `TestASetupManifestRoundfixVersionMustBeAReleaseVersion`.
14. Drop `WriterVersion` from the `baseline update` plan builder →
    `TestBaselineUpdateRecordsTheRunningRoundfixVersion` (the field stays
    absent after the update).
15. Remove the guard from `baseline update` → `TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix/update`.
16. Remove the guard from `baseline apply` → `TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix/apply`
    (it then fails reading the plan file, without the refusal).

## Build Order

1. The reader, the refusal in configuration loading and `Init`, the warning
   commands, their tests, the glossary entries and the Doctor and upgrade
   guide text (task_01).
2. The writer rule, the field's validation, the CLI writer version, the
   Baseline guard, their tests and the Baseline and Update Check guide text
   (task_02) (depends on: 1).
3. The Roundfix `baseline` and `setup` references and the skill version
   (task_03) (depends on: 2).
4. The QA gate (task_04) (depends on: 3).

No other Spec is required. Specs 0377, 0376, 0378, 0379, 0380 and 0381 deliver
first and edit other sections of the shared files `_prd.md` → Prerequisites
names; each Task works on the tree it runs on, cites functions rather than
lines, and raises the skill version with the record command.

## Risks & Considerations

- A repository written by a newer release blocks a teammate on an older one
  until `roundfix upgrade`; that is the intended refusal, and the message
  names the fix.
- A build from source between releases carries the last release's number, so
  a repository it writes accepts that release; the gap is one release at most.
- A hand edit of the field to a lower number lets an older binary act again;
  the reference says the remedy is `roundfix upgrade`, never editing the field.
- A Daemon started before the repository moved keeps running with its loaded
  configuration; its next command refuses.
- Reading the Setup Manifest costs one small file read per command.

## Glossary

- adds: **Recorded Roundfix Version**
- changes: **Setup Manifest**

## Decisions

- The record lives in the Setup Manifest and moves only with a plan that
  changes a file. See ADR-0392, ADR-0103 and ADR-0287.
- One check in `config.Load` before Project Config; the Baseline commands and
  `Init` call it for the repository they write. See ADR-0392.
- Doctor, the Update Check, the Run readers and `stop` continue without
  Project Config; Doctor fails its first line. See ADR-0392 and ADR-0371.
- The message mirrors `store.SchemaVersionError`'s newer-schema refusal.
- The writer version is a `PlanRequest` input the CLI sets from `app.Version`,
  so library tests and the settlement contract pass none and no golden moves.
