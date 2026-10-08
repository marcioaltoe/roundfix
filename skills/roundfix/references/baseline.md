## Context-Driven Baseline

Use the Baseline Command for Context-Driven Baseline adoption, update, profile
management, recovery, Repository Skill Set restoration, and canonical asset
synchronization. The CLI is the runtime authority; the
`setup-context-driven` skill contains recipes only.

At a terminal, the human workflow is:

```bash
roundfix baseline --repo . --format text
```

It detects adoption or update, collects one Baseline Profile and repository
decisions, shows one consolidated Change Plan, and writes only after explicit
confirmation of the displayed Plan Digest. Rejecting a plan returns to a
selected decision area and requires a newly calculated complete plan and
confirmation.

For a repository with a compatible Setup Manifest, use the dedicated
non-interactive managed refresh:

```bash
roundfix baseline update --repo . --format json
roundfix baseline update --repo . --yes --format json
```

Without confirmation, a changed Plan is presented and nothing is written.
The same preview lists, under `skills.outdated`, each installed
Roundfix-owned skill older than the version the binary carries. With
unchanged guidance it then reports `plan_ready` instead of `current`, and
`--yes` refreshes the Repository Skill Set.
Doctor reports `DR-SKILL-TRAILS-SNAPSHOT` with status `warn` when a required
upstream skill is present and matches its lock but differs from the embedded
Setup Snapshot. The managed refresh lists those trailing skills under
`skills.drifted` and, after confirmation, restores them to the snapshot
commit through the existing restore path. Use `roundfix baseline update` to
preview that restore; `--no-skills` skips it.

`--yes` approves the Plan Digest computed in that invocation;
`--confirm-plan <digest>` approves a previously reviewed digest, and the two
flags are mutually exclusive. `--adopt-suggested` explicitly adopts and reports
suggestions only for decisions absent from the manifest. `--no-skills` skips the
Repository Skill Set refresh. `--skills-source-dir <path>` selects an offline
Git checkout or bare object store for external skill restoration. `--repo`
defaults to the current directory, and `--format` accepts `text` or `json`.

Update exits `0` when the repository is current or an approved refresh applied
and verified; `1` for apply, verification, output, rollback, or recovery
failure; `2` for invalid input, an incompatible manifest, or an unsafe
repository; `3` when adoption, a new decision, confirmation, or retention
action is required; and `130` when canceled. JSON output uses
`roundfix/baseline-update-result/v1`. A plan with an Unrecorded Managed Region
reports its path, managed identity, reason, and every reported line the refresh
removes; text says `no lines removed` when there are none, and JSON includes the
optional `unrecordedManagedRegions` field only when at least one exists. The
same report remains in the applied result. Managed refresh never invokes
semantic classification and preserves every non-managed byte exactly.

`baseline update` lists each Retired Skill the repository still holds under
`Skills retired` (`skills.retired` in JSON) with the paths to delete; it never
deletes the copy and never changes the update's state.

When a plan includes History Relocations, it will report each tracked file whose
citations its History Relocations would break. These warnings use
`baseline.history.citation` for each citing file,
`baseline.history.citation.omitted` when more citing files exist than the report
lists, and `baseline.history.citation.unscanned` when a tracked file could not
be scanned. The Plan Digest covers these warnings. Planning never rewrites a
citation, and apply writes the same files.

When Pending History exists, the managed refresh adds a `History` section to
the plan. Text output lists the history units and their planned records or
reductions, with `History: <status>`, `History units:`, `History applied:`,
`History refused:` and `History tag:` lines when those values exist. Each reason
is one line; uncommitted changes and paths outside the coverage of an existing
History Full Tag are refused and left untouched. A refused unit does not block
the Baseline Plan or change the update's state. A repository-wide history
precondition such as an invalid tag can make the history section `blocked` with
its one-line result while the Baseline Plan keeps its existing behavior.

JSON always includes a `history` object. Its status is `skipped`, `current`,
`pending`, `applied` or `blocked`; it carries the selected units, Refused Units,
tag action and apply result when those values exist. When units are selected,
the history section is part of the same Plan Digest as the Baseline Plan, and
`--yes` or `--confirm-plan <digest>` approves both. Apply runs the Baseline Plan
first, creates the annotated `history-full` tag at `HEAD` when absent, then
converts every selected unit. The tag is never moved or pushed by the update;
push it with `git push origin history-full`.

Pass `--no-history` to leave the history section out and keep the reviewed
batch procedure under `roundfix history sanitize`.

Automation and Agents use the non-interactive plan/apply pair for first
adoption or a Profile change:

```bash
roundfix baseline plan --repo . --profile <profile-id> --decision-file <decision-file> --format json
roundfix baseline plan --repo . --profile-file .roundfix/baseline/profiles/<profile-id>.json --decision-file <decision-file> --format json
roundfix baseline apply --repo . --plan <plan-file> --confirm-plan <plan-digest> --format json
```

Planning never prompts or writes. Exit `0` emits one complete
`roundfix/baseline-plan/v1` document; exit `3` emits a
`roundfix/baseline-result/v1` next action and no partial plan. Apply accepts
only that strict portable document and its exact digest. A stale preimage,
confirmation mismatch, or unrelated Git lineage exits `3`; generate and
review a new plan instead of forcing the old one.

Requested results go to stdout; diagnostics and progress go to stderr. Baseline
reports repository formatter and Verification commands as recommendations and
never executes them. It also never installs dependencies, connects to live
infrastructure, follows unsafe links, mutates nested instruction carriers, or
lets ACP proposals authorize writes.

Re-check Repository Capability evidence without changing the repository:

```text
roundfix baseline capabilities check [--profile <id>] [--repo <path>] [--format <text|json>]
```

This reads local Repository Capability evidence through the evaluator and
divergence renderer Baseline planning uses. With no `--profile`, it resolves
the current Baseline Profile from a valid Setup Manifest; no resolvable Profile
is a named error. `--repo` defaults to the current directory and `--format`
defaults to `text`; JSON uses `roundfix/baseline-capability-recheck/v1`.
It resolves no decisions, writes no repository or journal bytes, executes no
candidate or repository command, and uses no network. Exit `0` means no
blocking divergence, `1` means output failure, `2` means invalid arguments,
repository failure, or no resolvable Baseline Profile, and `3` means a blocking
divergence.

To create a repository-owned Baseline Profile from an embedded built-in source:

```text
roundfix baseline profile init --id <id> [--from <built-in-id>]
```

This reads allowed embedded entry IDs from the selected built-in Profile
(`go-cli-tui` by default) and exclusively creates
`.roundfix/baseline/profiles/<id>.json`. The required ID is lowercase. It
does not compose profiles, copy assets, or accept executable or remote content.

Use explicit maintenance operations only when the user placed them in scope:

```bash
roundfix baseline profile show <profile-id> --format json
roundfix baseline profile validate <profile-id> --format text
roundfix baseline skills restore --repo . --profile <profile-id> --skill <skill-name> --format json
roundfix baseline skills reconcile --repo . --profile <profile-id> --source <owner/repo> --revision <40-hex-commit> --format json
roundfix baseline assets sync --source-dir <canonical-setups> --check --format json
```

Restore and reconcile accept a repository profile, and
`restore.profile-unresolved` names its searched
`.roundfix/baseline/profiles/<profile-id>.json` path.
`restore.snapshot-conflict` names an external skill whose embedded Setup
Snapshot contracts disagree.

`roundfix baseline skills reconcile` removes only lock entries absent at the
selected immutable commit. It preserves present, moved, unrelated, and
required entries, and requires the same reviewed Plan Digest confirmation as
skill restoration before applying a non-empty plan. When a Profile-required
skill is absent at the selected revision, reconciliation blocks with exit `3`
and finding `reconcile.required-removed`, prints `plannedChanges: []` and
`planDigest: null`, and writes nothing.

Skill restoration and reconciliation read `skills-lock.json` after acquiring
the source. If the lock changes during planning before its transaction
preimage is captured, the command refuses with `lock.changed-during-plan`,
exits `3`, and writes nothing.

Editing Roundfix-owned skill content no longer requires a Baseline digest or
characterization-corpus regeneration step: compatibility readiness depends on
the declared version comparison, not the skill bytes. Inside the Roundfix
source repository, keep the canonical and embedded Skill copies synchronized
with `make skills-sync`.

An expressly authorized edit to a Baseline module can still change derived
catalog pins. Regenerate those deterministic artifacts only with:

```bash
make baseline-digests
```

Every pin rewritten by that command is fallout of the authorized source edit
and needs no separate per-Spec express authorization. A hand-edited pin value
remains unauthorized. Run the command before repository Verification; do not
transcribe digest values.

A non-empty skill-restoration preview requires its exact current Plan Digest
through `--confirm-plan`. Asset refresh without `--check` requires explicit
maintainer intent. For Decision Documents, preservation, cross-clone safety,
recovery, migration, security limits, and completion evidence, follow
`docs/user-guide/context-driven-development.md#adopt-or-update-the-context-driven-baseline`.
