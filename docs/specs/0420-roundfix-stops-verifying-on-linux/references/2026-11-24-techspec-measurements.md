# Measurements of the TechSpec phase, Spec 0420 (2026-11-24)

Measured by the author of the TechSpec on `cefd4e82` (v0.153.2), with a binary built from that tree by `go build -buildvcs=false ./cmd/roundfix`, a temporary home directory and a temporary repository. The rehearsals of section 5 ran in a clone outside the repository and changed nothing in it. Section 8 holds what the Research Subagent and the author read on the web and in the Secondbrain.

## 1. The governed set

`speccheck.GovernedPath` was run over 100 paths in a `go test -overlay` probe whose test file lives outside the repository: every file of section 2, the planning files of this Spec, and neighbours of those files probed as a control. 18 are governed: the twelve the authorization record grants, and six that the design leaves untouched.

Governed and granted, each changed by one step of the Build Order:

- step 1: `.roundfixrc.yml`, `docs/references/coverage-record.json`;
- step 2: `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `Makefile`, `docs/agents/specific-repository.md`;
- step 3: `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runtime.md`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`.

Governed, not granted and not changed, probed as a control: `internal/cli/cli_test.go`, `.github/workflows/ci-conventions.yml`, `.github/workflows/secondbrain-sync.yml`, `go.mod`, `internal/docscontract/testdata/corpus-golden.json` and `.agents/skills/roundfix/references/release.md`. Ordinary and probed as a control: `internal/config/config_test.go`, `internal/corebudget/corebudget.go`, `CHANGELOG.md`, and the three files of `internal/verifyselect` whose tests read the workflows and do not change, `suite_budget_test.go`, `release_workflow_test.go` and `release_workflow_compare_test.go`. `docs/user-guide/release-runbook.md` is ordinary.

`internal/cli/cli_test.go` holds one test of a deprecated key, `TestRunFetchWarnsAndIgnoresDeprecatedUserConfig`. It names only `resolve.concurrent`, so registering two more keys does not move it. The corpus golden holds the finding counts of the active Specs, all zero: this Spec must reach the Run with no finding, or that file would have to change.

## 2. Files of the design, by step

Every path below is ordinary unless section 1 names it.

Step 1, the product:

- deleted: `internal/cli/deliver_linux_lane.go`, `internal/cli/deliver_lane.go`, and the eight test files of section 6;
- changed: `internal/cli/deliver.go`, `internal/cli/cli.go`, `internal/cli/deliver_workflow.go`, `internal/cli/doctor.go`, `internal/cli/health.go`, `internal/delivery/engine.go`, `internal/delivery/park_class.go`, `internal/config/config.go`, the test files of section 7, `.roundfixrc.yml`, `docs/references/coverage-record.json`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`;
- created, with the names the probe used: `internal/cli/deliver_gate_without_lane_test.go`, `internal/cli/deliver_lane_removed_test.go`, `internal/cli/doctor_without_lane_test.go`, `internal/cli/deliver_stored_lane_records_test.go`, `internal/cli/retired_lane_keys_test.go`, `internal/config/retired_lane_keys_test.go`, `internal/delivery/stored_lane_records_test.go`.

Step 2, where the verification runs:

- changed: `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `Makefile`, `README.md`, `docs/agents/specific-repository.md`, `internal/verifyselect/full_contract_run_test.go`;
- deleted: `internal/verifyselect/windows_test_compile_repocontract_test.go`;
- created: `internal/verifyselect/platform_test_compile_repocontract_test.go`, `internal/verifyselect/verifying_runner_test.go`.

Step 3, the guide and the skill:

- changed: `docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/configuration.md`, `docs/user-guide/release-runbook.md`, the four references `deliver.md`, `setup.md`, `implement.md` and `runtime.md` under `.agents/skills/roundfix/references/` and their mirrors under `skills/roundfix/references/`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `internal/verifyselect/contracts_repository_test.go`;
- created: `internal/docscontract/retired_lane_documentation_test.go`.

Step 4, the decision record and the glossary:

- created: `docs/adr/0452-roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane.md`, `internal/docscontract/linux_verification_decision_record_test.go`;
- deleted: the files of ADR-0288, ADR-0294, ADR-0297, ADR-0303 and ADR-0316;
- changed: the files of ADR-0252, ADR-0285, ADR-0295, ADR-0384, ADR-0385 and ADR-0423, `CONTEXT.md`, `.secondbrain-export`, `internal/verifyselect/contracts_repository_test.go`.

## 3. Behavior Surfaces

`skillcoverage.Map.SurfacesForPath` was run over the same paths in the same kind of probe. Twelve surfaces have a source among them:

| Surface | Source the design changes | Covering file the design changes |
| --- | --- | --- |
| `command deliver lane` (disappears) | `internal/cli/deliver.go`, the two deleted lane files, `internal/cli/deliver_workflow.go`, the guide's `deliver` page | `.agents/skills/roundfix/references/deliver.md` |
| `command deliver plan`, `resume`, `retry`, `start`, `status`, `stop` | the same, without the two deleted files | `.agents/skills/roundfix/references/deliver.md` |
| `command doctor` | `internal/cli/doctor.go`, the guide's `doctor` page | `.agents/skills/roundfix/references/setup.md` |
| `config keys` | `internal/config/config.go`, the guide's configuration page | `deliver.md`, `setup.md` and `implement.md` are three of its eleven covering files |
| `guide docs/user-guide/commands/deliver.md` | itself | `deliver.md` |
| `guide docs/user-guide/commands/doctor.md` | itself | `setup.md` |
| `guide docs/user-guide/configuration.md` | itself | `deliver.md`, `setup.md`, `implement.md` |

`roundfix doctor --help` does not name the lane, so the fingerprint of `command doctor` does not move; the fingerprints of the six `deliver` commands, of `config keys` and of the three guide pages do, and `command deliver lane` leaves the record.

The release runbook adds two surfaces whose only covering file is `.agents/skills/roundfix/references/release.md`, which the grant does not hold: `guide docs/user-guide/release-runbook.md` and `command release plan`. The runbook is the single source of both. That reference names the release workflow only where it says that a repository may publish from a workflow; it describes neither the workflow's steps nor their order.

## 4. What the unchanged binary does

The removed command's refusal already exists for any other word:

```text
$ roundfix deliver bogus
roundfix: unknown deliver command "bogus"
Run 'roundfix deliver --help' for usage.
exit 2
```

`roundfix deliver lane --help` prints the `deliver` usage and exits `0`, and so does `roundfix deliver bogus --help`: a help flag is read before the subcommand. The usage holds three lines about the lane: `  roundfix deliver lane [--all] [<revision>]`, `  lane    Run the Linux Lane on demand and exit with its outcome` and `  --all  Test every package in the Linux Lane`. The root help holds `  roundfix deliver <plan|start|status|resume|retry|stop|lane> [<slug> ...]`.

Configuration, in a temporary repository:

- `delivery.linux_lane: "yes"` is refused: `delivery.linux_lane must be boolean: cannot unmarshal non-boolean value`;
- `delivery.linux_lane_image: "not a ref!!"` is refused: `delivery.linux_lane_image must be a container image reference, such as golang:1.26`;
- a valid value in Project Config or in User Config loads with no line on standard error;
- an unknown key under `delivery` is refused with `field bogus_key not found in type config.deliveryOverlay`;
- `delivery: {}` loads, which is what a `delivery` mapping holding only the two keys becomes once they are removed before the decode;
- a key already registered as deprecated, `jev.router_min_credit_usd`, prints its one line and `roundfix runs list` then prints `No Runs found.` and exits `0`.

`roundfix doctor` with `delivery.linux_lane: true` prints `linux-lane: warn (DR-LANE-IMAGE-UNKNOWN: could not inspect go.mod at HEAD; next: set delivery.linux_lane_image)` between the `verification-cores` line and the `codex` line.

Stored records: `roundfix deliver status` prints `Warning: <slug> <warning>` for an item's stored warning and `Park: <slug> <class>: <next>` for a parked item. A blocker that begins `linux-lane-container:` is classed `environment` today, with the next step `repair the Linux Lane container cause in the blocker (Docker, the image, or delivery.linux_lane_image on the default branch), then run roundfix deliver retry <slug>`. Without its case in the classifier it would be `unclassified`, with `resolve the blocker, then run roundfix deliver retry <slug>`. A blocker `gate-failed: Linux Lane change failure: …` is classed `finding` by the `gate-failed` case and names nothing of the lane in its next step.

## 5. Rehearsals in a clone

- **The five decision records deleted and committed.** `roundfix spec check 0420-roundfix-stops-verifying-on-linux --stage prd --strict` still proves every Claim Receipt of the PRD that quotes one of them: the check reads a deleted record from Git. Of the `internal/docscontract` tests, only `TestGlossaryLinksNameAcceptedADRs` fails, on the nine links of the six glossary entries that cite those records. So the records, the glossary entries and ADR-0452 change in one step.
- **The lane removed from the `deliver` usage, the dispatch and the root help.** Only `TestTheSkillCoverageMapIsCurrent` fails: six stale fingerprints and the surface `command deliver lane`, which is in the map and in the record and is no longer computed. After the surface's entry leaves `docs/references/skill-coverage.json` and the record is written by its command, every `internal/docscontract` test passes while the guide and the skill still describe the lane. So the product step and the documentation step can be two steps, and the product step carries the record and the map.
- **The coverage record.** `TestCoverageEquivalence` reports a test the record lists and the tree lacks as a regression, and only logs a test the tree gained. Only a step that removes a test needs the record written again.

## 6. The 36 tests the coverage record may lose

Each is a test of the lane, of its command, of its doctor check or of its two keys. The record on `cefd4e82` holds 36 names that contain `Lane`, and they are these.

`internal/cli/deliver_linux_lane_test.go` (13): `TestLinuxLaneRunsTheChangedPackagesInAContainer`, `TestLinuxLaneSkipsWithANamedReason`, `TestLinuxLaneFailureFailsTheGate`, `TestLinuxLaneCancellationRemovesTheContainer`, `TestGateWithoutTheLaneStartsNoDocker`, `TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead`, `TestLinuxLaneOptInPropagatesCancellationAndAnUnreadableConfig`, `TestLinuxLaneRunsTheConfiguredImage`, `TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure`, `TestLinuxLaneOutcomeClassification`, `TestLinuxLaneMountsTheRepositoryHistoryReadOnly`, `TestLinuxLaneSkipsAnObjectStoreItCannotMount`, `TestLinuxLaneReportsThePackagesTheContainerSelected`.

`internal/cli/deliver_lane_test.go` (11): `TestDeliverLanePassedTranscript`, `TestDeliverLaneChangeFailureTranscript`, `TestDeliverLaneContainerFailureTranscript`, `TestDeliverLaneSkipTranscript`, `TestDeliverLaneSelectsPackagesAndImage`, `TestDeliverLaneCancellationRemovesContainer`, `TestDeliverLaneObjectStoreSkipTranscript`, `TestDeliverLanePackageCountTranscript`, `TestDeliverLaneNoPackageCountTranscript`, `TestDeliverLaneGoListFailureTranscript`, `TestDeliverLaneGoListSignalTranscript`.

`internal/cli/deliver_linux_lane_go_list_test.go` (2): `TestLinuxLaneScriptNamesAFailedGoList`, `TestLinuxLaneGoListFailureIsNamed`.

`internal/cli/deliver_lane_skip_visible_test.go` (2): `TestGateReturnsTheLaneEnvironmentSkip`, `TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript`.

`internal/cli/doctor_linux_lane_test.go` (2): `TestDoctorLinuxLaneCheckReportsDockerAndTheImage`, `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`.

`internal/config/linux_lane_test.go` (1): `TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt`.

`internal/config/linux_lane_image_test.go` (2): `TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt`, `TestDeliveryLinuxLaneImageRefusesAnInvalidReference`.

`internal/delivery/lane_skip_warning_test.go` (2): `TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody`, `TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne`.

`internal/cli/deliver_core_budget_test.go` (1 of its 2): `TestDeliverLanePassesItsVerificationCoreBudget`. `TestTheRepositoryGateCommandReceivesTheCoreBudget`, in the same file, stays.

These 36 tests are also the characterization of what the lane does on the unchanged tree: the gate's call, the command's transcripts, the doctor line and the two keys each have a named test there.

## 7. Tests of another behavior that name the lane

Each keeps its name and loses only the lane from a fixture, a stub or a pinned line.

- `internal/cli/deliver_retry_test.go`: `TestTopLevelUsageNamesDeliverRetry` pins the root help line.
- `internal/cli/deliver_owner_log_test.go`: `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate` writes `linux_lane: false` and asserts the gate's lane-skip field.
- `internal/cli/doctor_disk_space_test.go`: `TestDoctorWarnsWhenALocationIsLowOnDisk` has a row that places the warning before the lane line.
- `internal/cli/doctor_verification_cores_test.go`: `TestDoctorReportsTheVerificationCoreBudget` pins `linux-lane: ok`.
- `internal/cli/machine_verification_status_test.go`: `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns` and `TestDoctorReportsTheMachineVerificationCapacity`.
- `internal/cli/shared_cache_trim_test.go`: the helper `sharedCacheImplementWorkspace` and `TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification`.
- `internal/cli/external_skill_lock_test.go` and `internal/cli/external_skill_presence_test.go`: each clears the doctor's lane dependency.
- `internal/cli/releaseplan_release_preparation_test.go`: one comment.
- `internal/delivery/gate_blocker_test.go`: `TestGateParksWithTheBlockerItNames` has four rows of the lane.
- `internal/delivery/low_disk_warning_test.go`: `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate`.
- `internal/delivery/park_class_test.go`: `TestParkClassesAndNextActions` has the lane's row and its assertion.

No file under a `testdata` directory and no golden names the lane.

Tests that read the workflows, all in `internal/verifyselect` and none of them naming a runner: `TestMainAndReleaseRunEveryContract` reads the release's job `release` for `make verify` followed by `Verify every contract` and a later `Publication preflight`, so the new layout moves its assertions; `TestTheCIPushJudgesTheSuiteTimeBudgetLast`, `TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts` and the tests of `release_workflow_test.go` and `release_workflow_compare_test.go` read steps this design leaves where they are. `TestTheSuiteTimeBudgetIsRecordedThenJudged` reads the Makefile's default of 420 seconds, which stays. `TestRepositoryContractTestsDeclareTheirRelevance` pins 14 `docscontract` files in `internal/docscontract` and requires every name of the Makefile's `REPO_CONTRACT_TESTS` to be a discovered test.

## 8. Sources read on 2026-11-24

GitHub's runner image for the label `macos-26`, https://raw.githubusercontent.com/actions/runner-images/main/images/macos/macos-26-arm64-Readme.md, image 20261117.0516.1 on macOS 26.7.2. It lists four of the eight External Test Tools: "Bash 3.2.57(1)-release", "jq 1.8.2", "Node.js 24.22.3" and "Rustfmt 1.11.0-stable". The other four, `cmp`, `grep`, `gzip` and `sh`, are part of macOS and the image does not list them; that they are present is the author's knowledge of the system, not a quotation. The image also lists Go 1.26.11 among its cached tools, Git 2.56.1 and Xcode Command Line Tools 26.6, and lists no Docker.

The `actions/runner-images` README, https://github.com/actions/runner-images: the standard arm64 image of macOS 26 answers to "`macos-latest`, `macos-26` or `macos-26-xlarge`" and the one of macOS 15 to "`macos-15`, or `macos-15-xlarge`". "The `-latest` migration process is gradual and happens over 1-2 months". "To avoid unwanted migration, users can specify a specific OS version in the yaml file". "We support (at maximum) 2 GA images and 1 beta image at a time." "All macOS 14 images retired from GitHub Actions on November 2, 2026." The migration of `macos-latest` from macOS 15 to macOS 26 ran from 2026-06-15 and was reported complete on 2026-07-22 (issue 14167 of the same repository). So `macos-26` is the newer of the two supported versioned labels and the one `macos-latest` names today.

`actions/setup-go`, `src/package-managers.ts` and `src/cache-restore.ts`: the action caches the two directories that `go env GOMODCACHE` and `go env GOCACHE` print, and its cache key holds the runner's operating system. The step `Align GOCACHE with the restored cache` reads `go env GOCACHE`, so it holds no Linux path; the first run on macOS starts with an empty cache.

GitHub Docs, "Workflow syntax", https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax: jobs a job needs must "complete successfully before this job will run", and "If a job fails or is skipped, all jobs that need it are skipped unless the jobs use a conditional expression that causes the job to continue."

npm Docs, "Trusted publishers", https://docs.npmjs.com/trusted-publishers: "Trusted publishing currently supports only cloud-hosted runners." The page names no operating system, so the publishing job stays where it is proven. The trusted publisher of each package is bound to the repository and the workflow file name, which do not change.

Published practice for a release in two jobs, found with Exa: several public release workflows run the tests in a job that holds no publishing permission and give `id-token: write` only to a publishing job that needs it. One of them states the rule in its header: "`id-token: write` is granted on the publish job only — never to `verify`" (https://github.com/daloyjs/daloy/blob/main/.github/workflows/release.yml). This is the layout API Contract 7 of the TechSpec takes.

Secondbrain, read only: `wiki/index.md`, then two searches with the Indexus tool `pesquisar_acervos`, which answered with wiki results.

- `wiki/concepts/testes-paralelos-em-go.md` records that this suite is bound by the processes it starts, not by processors: more than 83,000 `git` processes in one run, measured on a machine of ten processors. A runner of three processors and a different cost of starting a process can therefore take a time that the Linux number does not predict. It is the reason the budget's number is not estimated (Delegated Decision 1).
- `wiki/sources/digest-roundfix-toolchain-formatter-e-subprocessos-go-no-macos-2026-09-09.md` records failures of test subprocesses that only macOS showed, with no proven cause. It supports verifying where the product runs and is the reason Risks names a first run on a new macOS host.
- `wiki/concepts/tags-protegidas-e-release-por-workflow-no-github.md` records another project of the maintainer that stopped running the same verification several times for one delivery and its release. In the two-job layout, re-running a failed publication no longer repeats the verification.

The Research Subagent ran three sweeps: the code of the lane and the tests that name it, the 35 decision records with the glossary and the documentation layout, and the web pages above, which it read with the web search and fetch tools. The author ran every command of this file, the Secondbrain searches and two Exa searches, which found the runner image's README, the migration issue and the release workflows quoted here.
