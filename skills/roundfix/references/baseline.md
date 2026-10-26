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
Doctor reports `DR-SKILL-DIFFERS-FROM-SNAPSHOT` when a required upstream skill
is present and matches its lock but differs from the embedded Setup Snapshot.
Newer or older is unknown, so the `skills` line keeps its readiness status and
gains no next action. The managed refresh lists such a skill under
`skills.differing`, keeps it, and never restores it unasked; restore one
deliberately with `roundfix baseline skills restore --profile <id> --skill <name>`.
A missing skill, or one that no longer matches its lock, is still restored under
`skills.drifted`.

The preview without `--yes` or `--confirm-plan` lists, under `skills.missing`
(`Skills missing` in text), each external skill the adopted modules require and
dispatch that is absent from `.agents/skills` or `skills-lock.json`, or no
longer matches its lock entry: a Missing Required Skill. Each entry names the
skill and the command that restores it,
`roundfix baseline skills restore --profile <id> --skill <name>`.
`skills.status` reads `missing`; with unchanged guidance the state is
`plan_ready` instead of `current` and the exit code is `3`, as for an owned
skill older than the binary's, and `--yes` restores the listed skills.
A skill the upstream catalog retired while the Baseline still dispatches it
is a Baseline defect to report, not a skill to keep restoring by hand.
`--no-skills` skips the list, and an absent or unreadable `skills-lock.json`
adds nothing to it; `roundfix doctor` reports that state.

`--yes` approves the Plan Digest computed in that invocation;
`--confirm-plan <digest>` approves a previously reviewed digest, and the two
flags are mutually exclusive. `--adopt-suggested` explicitly adopts and reports
suggestions only for decisions absent from the manifest.

`--decision <id=value>` records or changes a decision the repository's profile
selects before the plan is built; repeat it for several decisions. The value is
parsed as `baseline plan --decision` parses it, the plan and its digest cover
the changed Setup Manifest, and the result lists each one under
`declaredDecisions` with its previously recorded value. A decision the profile
does not select or fixes, `preservation.mode`, or a malformed pair exits `2` and
writes nothing. A declared decision also answers a new decision without
`--adopt-suggested`. A repository that mounts Better Auth at `/auth` runs
`roundfix baseline update --decision auth.base_path=/auth`.

`--no-skills` skips the Repository Skill Set refresh.
`--skills-source-dir <path>` selects an offline Git checkout or bare object store for external skill restoration. `--repo`
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

The history plan names every file it will remove before approval. Each pending
unit carries sorted repository-relative `paths` in JSON: the files a removal
or conversion deletes, the Archive Records `record-links` rewrites, or the live pages `live-links`
rewrites. Text
prints the same paths in the same order as `  - <path>` lines directly under
each unit's summary.

When Pending History exists, the managed refresh adds a `History` section to
the plan. Text output lists the history units and their planned records or
removals and `record-links` repairs, with `History: <status>`, `History units:`, `History applied:`,
`History refused:` and `History tag:` lines when those values exist. Each reason
is one line; uncommitted changes and removed paths held by neither the History
Full Tag nor a commit the default branch reaches are refused and left untouched. A refused unit does not block
the Baseline Plan or change the update's state. A repository-wide history
precondition such as an invalid tag can make the history section `blocked` with
its one-line result while the Baseline Plan keeps its existing behavior.

For each selected unit, the history section also lists every History Citation
that the removal breaks:

```
History citations: <n>
- cites <path>:<line> names <target>
- cites <path>:<line> names <target> (live)
- cites <path>:<line> names <target> (live; promote with --promote <target>)
- cites <path>:<line> names <target> (promoted to docs/references/<basename>)
History stale citations: <n>
- cites <path>:<line> names <target>; history-full holds it: git show history-full:<target>
```

The citation list is in report order and is part of the Plan Digest. JSON
always includes `history.citations` and `history.staleCitations` as arrays.
Each citation has `path`, `line`, `target` and `live`; `offer` is present
for a Live Document promotion offer and `promotedTo` is present when the
target was promoted. Stale History Citations are advice: they do not enter the
digest, state, category or exit code.

Pass repeatable `--promote <path>` values to copy named files from selected
Legacy Archive Folders to `docs/references/` during apply. A value must be a
clean repository-relative file inside a selected folder, outside every Refused
Unit, not a core artifact, unique by basename among the promotions, and absent
from `docs/references/`. `--promote` with `--no-history`, an invalid value,
a refused unit or an occupied destination exits 2 with category `invalid`
and writes nothing. The update's conversions read Legacy Archive Folders with
the same summary, field, source and creation-date fallbacks, and Legacy Link
Rebase rewrites each promoted copy, using a Git Locator for an unpromoted
folder destination. After removals the update rewrites every History Citation it lists, in every Citation Source: each tracked text file outside the skill trees and testdata directories. A prose path that starts with ./ or ../ is resolved against the directory of the citing file. Outside Markdown the Git Locator is written without backticks. A Reduced
History Entry uses its original provenance path and commit when reachable;
no new reduced entry is written.

After removal units, `record-links` repairs each surviving Archive Record's
relative links to missing files and to itself. It tries the record-relative
path, then the record's original `source` folder; the first reachable commit
holding the target wins. A self-link becomes the Git Locator of the source
folder. Before that, a link whose target, resolved against the `source` folder,
is a file whose copy `docs/references/<file>` the record lists under `promoted`
and the tree holds is rewritten to a relative link to that copy, and the plan
names `docs/references/<file>` as its rewrite (ADR-0308). When neither path has
a qualifying commit, the record-relative path
and its pull request are used. Code spans, fences, URLs and Git Locators are
not relative links. A record removed by the same plan is skipped.

Text prints `- record-links: rewrites <l> link(s) in <n> Archive Record(s)`.
Removal units print `- <unit>: removes <n> file(s) (<b> bytes)`. Rewritten
citations append `; rewrites to <replacement>` and carry `rewrite` in JSON.
JSON actions are `convert`, `remove` and `rewrite`; a rewrite unit carries
`links`. `bytesAfter` is omitted, `applied.reduced` stays 0, and apply adds
`rewrittenCitations` and `rewrittenFiles`. File totals count only removals.
The Plan Digest binds removals, record repairs and replacement text.


A superseded ADR's accepted successor receives a missing lineage line exactly
as `Supersedes ADR-<number> (<date>); Git keeps its text at <locator>.`
A retired ADR whose file name is not a number followed by a slug is named by its file name without the extension. The date is the day the retired ADR was superseded: its updated_at in UTC, else the committer date of the newest commit that changed it, else the committer date of HEAD.
The locator is the removed path's Git Locator. Apply appends the line after a blank line, preserves the
successor's existing bytes and never adds the same line twice. Several lines
for one successor follow removal-path order. Each line is a History Citation:
`cites <successor>:<n> names <removed>; rewrites to <line>` in history text,
and under `citations` in JSON. Baseline text prefixes the citation with `- `;
its Plan Digest binds the exact line, and an uncommitted successor refuses
the ADR unit. Successors are cited, never listed as removal paths.

`superseded_by` accepts `ADR-NNNN` or `NNNN` to select an accepted ADR
under `docs/adr/` by number. It also accepts `NNNN-slug`, `NNNN-slug.md`,
or a path whose last element is `NNNN-slug.md`; these forms select that
exact filename under `docs/adr/`. An unknown form, a missing file or a
successor without `status: accepted` keeps the refusal below.

No lineage line is appended when the successor already names the retired ADR
in its body, or when its front matter `supersedes` names it as `ADR-NNNN`,
`NNNN`, `NNNN-slug`, `NNNN-slug.md` or a path ending in `/NNNN-slug.md`.
Here the number and slug are the retired ADR's; `supersedes` can be a string
or a list of strings. Naming another ADR does not suppress the lineage line.

A superseded ADR with no accepted successor stays with this reason:
`superseded ADR has no accepted successor; name the accepted ADR that replaces it in superseded_by, or set its status to deprecated`.
History text prints `refused <path>: <reason>`; baseline text prefixes it with
`- `. JSON lists `{"unit":"<path>","reason":"<reason>"}` under `refused`.
The refusal names that ADR's repository-relative path, leaves the other
retired ADRs eligible and does not change the command's exit code.

After `record-links`, `live-links` repairs missing relative links in Live
Documents and `docs/user-guide/` pages, including the configured Spec Root.
It writes the Git Locator printed by `SC-HISTORY-RETAINED`, changes only link
spans and skips Archive Records, removed pages, self-links, code spans,
fences, URLs and existing Git Locators. A second plan after apply has no
`live-links` unit.
The repair leaves a link to a file that a Task of an active Spec declares in creates: and has not completed. It resolves a missing link by its file name when exactly one removed history file has that name.

The repair line is `- live-links: rewrites <l> link(s) in <n> file(s)`.
JSON lists unit `live-links`, action `rewrite`, sorted page `paths`, `files`
and `links`; its citations and replacement text join the Plan Digest.
Text lists each page as `  - <path>` immediately after the repair line.

JSON always includes a `history` object. Its status is `skipped`, `current`,
`pending`, `applied` or `blocked`; it carries the selected units, Refused Units,
tag action and apply result when those values exist. When units are selected,
the history section is part of the same Plan Digest as the Baseline Plan, and
`--yes` or `--confirm-plan <digest>` approves both. Apply runs the Baseline Plan
first, creates the annotated `history-full` tag at `HEAD` when absent, then
converts folders, removes retired files and repairs record links. The tag is never moved or pushed by the update;
push it with `git push origin history-full`.

Pass `--no-history` to leave the history section out and keep the reviewed
batch procedure under `roundfix history sanitize`.

While Pending History remains, an update with --no-history still writes the retirement rule into the guide, so SC-HISTORY-RETAINED reports every retained history file on each active Spec, and the path to green is the history conversion: rerun roundfix baseline update without --no-history, or apply roundfix history sanitize batches until nothing is pending.

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

`roundfix baseline capabilities --help`, `roundfix baseline skills --help` and
`roundfix baseline assets --help`, or `-h`, print the `baseline` usage,
which lists their commands, and exit `0`.

Run bare, or with a command they do not have, they exit `2`: the usage error names the group once on standard error and points at `roundfix baseline <group> --help`, and an unknown `baseline` command points at `roundfix baseline --help`.

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

## Vendored skills

A vendored skill is installed from an upstream source and never edited in the repository, because the next sync replaces its folder; the restore and the managed refresh only install, compare and restore it. When one diverges from what the repository needs, correct it in this order:

1. A rule that governs over the skill's default: a Baseline clause, or a Repository-Specific Normative Rule. This answers a convention difference such as the package manager, the test runner or a library choice.
2. A repository-authored skill that supplies or overrides only what is missing, installed beside the vendored one.
3. A repository-authored skill that replaces the vendored one, only when the vendored skill is useless or almost useless to the repository.

Send an objective upstream defect to the upstream as a pull request; it never justifies the third choice.

### Starter Files

A Starter File is repository tooling a Baseline module offers. The `python`
profile offers `Makefile`, `.roundfixrc.yml` and `.github/workflows/verify.yml`.

```bash
roundfix baseline starters --profile python
roundfix baseline starters --profile python --confirm-plan <digest>
```

Preview writes nothing, lists absent paths as `create` and existing paths as
`keep`, and prints the Plan Digest. A non-empty preview exits 3; an empty
preview is `current` and exits 0. Confirmation with `--confirm-plan` writes
only absent paths, never overwriting, comparing or removing existing files.
A stale digest exits 3 without writing; preview again. A write failure removes
every file the call created and exits 1; created directories may remain.

Confirming the exact preview is the maintainer's express authorization for
exactly those paths. An Agent runs the confirmation only with the
maintainer's express authorization. Once written, the repository owns the
file: it is outside the Change Plan and Setup Manifest and no refresh changes
it. `--format json` emits `roundfix/baseline-starters/v1`; invalid input,
an unknown profile or an unsupported format exits 2.
