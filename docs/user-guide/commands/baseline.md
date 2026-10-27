### baseline

Profile-change output carries valid recorded decisions as `Carried decisions:`.
Apply names each Root Backup with `baseline.preservation.backup.written`;
reconcile names the Setup Snapshot revision. Restoration uses `create-link`
and never moves an installed skill.

```text
roundfix baseline [--repo <path>] [--format <text|json>]
roundfix baseline update [--repo <path>] [--format <text|json>] [--yes | --confirm-plan <digest>] [--adopt-suggested] [--decision <id=value> ...] [--no-skills] [--no-history] [--skills-source-dir <path>]
roundfix baseline plan --profile <id> [--decision <id=value> ...] [--decision-file <path> ...] [--repo <path>] [--format <text|json>]
roundfix baseline apply --plan <file> --confirm-plan <digest> [--repo <path>] [--format <text|json>]
roundfix baseline capabilities check [--profile <id>] [--repo <path>] [--format <text|json>]
roundfix baseline profile init --id <id> [--from <built-in-id>]
roundfix baseline profile show <id> [--format <text|json>]
roundfix baseline profile validate [<id>|<path>] [--format <text|json>]
roundfix baseline skills restore --profile <id> [--skill <name> ...] [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
roundfix baseline skills reconcile --profile <id> --source <owner/repo> --revision <commit> [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
roundfix baseline starters --profile <id> [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
roundfix baseline assets sync --source-dir <path> [--check] [--format <text|json>]
```

The root command is the terminal-only human workflow for first adoption and
later updates. It detects existing state, collects instruction preservation,
one Baseline Profile, and repository decisions, then presents one consolidated
Change Plan. Rejecting the plan returns to a selected decision area and
recalculates the complete plan. Mutation requires one explicit confirmation of
the current displayed Plan Digest.

#### Capability evidence and remediation

Profile alignment groups unsatisfied Repository Capabilities by requirement
strength: blocking, advisory, and informational. A blocking divergence prevents
readiness until it is resolved. An advisory divergence never blocks readiness
or apply; any optional next action appears after that statement. Informational
divergences record evidence without requesting a decision.

Executable discovery inspects each PATH candidate and resolves a bounded
symlink chain without executing the candidate or its target. A rejected
candidate reports one of these reasons:

- `link-cycle`: the chain revisited a link or exceeded the bounded hop count.
- `broken-link`: a link could not be read or its target does not exist.
- `not-executable`: the resolved target is not a regular executable file or
  lacks an executable permission bit.

When Profile alignment is blocked, the interactive prompt offers four
outcomes:

1. Change the Baseline Profile and recalculate alignment.
2. Create a reviewed repository-owned Profile adaptation. Adaptation is
   removal-only: it can remove repository-inapplicable Profile modules or
   Repository Capabilities, but cannot add policy, waive a universal
   requirement, or weaken one.
3. Remediate in the repository and re-run. This exits without writing, prints
   remediation for every unsatisfied divergence, prints the exact capability
   re-check command, and records an outcome distinct from decline.
4. Decline without writing.

After applying the printed remediation, run the exact
`roundfix baseline capabilities check` command from the prompt before returning
to planning. The re-check evaluates the same capability evidence and produces
the same capability outcomes as a full plan. It requires no decisions, resolves
no decisions, and writes nothing: no repository file, journal entry, or
configuration.

#### Retention and completion evidence

When the Baseline identifier is unchanged but the Profile or catalog digest
changes, planning requires retention accounting for every previously managed
Normative Clause. Each clause must have an explicit disposition. An
unaccounted clause stops planning with an action-required result, and Roundfix
does not offer apply.

The final result reports five independent status axes. Each axis is `verified`
or `not run`:

| Axis | Evidence reported |
| --- | --- |
| Approved postimages | The approved repository postimages were written and verified. |
| Semantic retention | Every applicable prior Normative Clause has an accounted disposition. |
| Profile alignment | The resulting Setup Manifest and selected Baseline Profile align. |
| Repository Verification | The repository's declared Verification command ran successfully. |
| Idempotence | A repeated check found no further Baseline change to apply. |

Verified postimages alone do not mean the Baseline update is complete.
Completion language appears only when semantic retention is `verified` and the
idempotence check passed.

`baseline plan` is read-only, non-interactive, local, and network-free. JSON
exit `0` emits one complete `roundfix/baseline-plan/v1` document. Missing
decisions, manual classification, or unresolved alignment exit `3` with a
`roundfix/baseline-result/v1` next action and no partial plan.

`baseline apply` accepts only a strict portable JSON plan and its exact
`planDigest`. It validates clone-stable Git lineage, the embedded catalog,
profile identity, every bounded preimage, the complete managed-entry ledger,
and its derived file projection. It never substitutes a newer plan. Matching
postimages make an exact reapply an idempotent success.

`baseline profile init` creates
`.roundfix/baseline/profiles/<id>.json` from one embedded built-in profile.
`show` resolves one built-in or repository-owned profile; `validate` checks one
ID, one direct profile path, or every repository-owned profile. Repository
profiles can use embedded entry IDs only and cannot compose profiles or load
remote executable content.

`baseline skills restore` is a separate confirmation-gated Repository Skill
Set operation. Its non-empty preview exits `3` with a current Plan Digest;
`--confirm-plan` applies only that exact preview. `--source-dir` selects an
offline Git checkout or bare object store containing the declared immutable
source commit. The command reads `skills-lock.json` after the source is
acquired. If the lock changes during planning before its transaction preimage
is captured, the command refuses with `lock.changed-during-plan`, exits `3`,
and writes nothing.

Restore, lock reconciliation, and the managed refresh accept both built-in
and repository-owned Baseline Profiles. A repository profile's restorable
skills are the external skills its modules require, with contracts taken
from the embedded Setup Snapshots that agree on each skill. Its restore and
reconcile payloads carry `setup: null` because it names no single snapshot.
`restore.profile-unresolved` says the profile is neither built-in nor
resolvable and names the searched `.roundfix/baseline/profiles/<id>.json`
path. `restore.snapshot-conflict` names a skill whose embedded snapshot
contracts disagree; no contract is chosen by catalog order.

The managed refresh previews its Plan without writing. `--yes` approves the
Plan Digest computed in that invocation; `--confirm-plan <digest>` approves a
previously reviewed digest. `--adopt-suggested` explicitly adopts and reports
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

After the bundle install, an existing `skills-lock.json` entry naming an owned
skill receives the installed folder's `computedHash`; every other field is
kept. No entry or lock is created when none exists. The result reports
`skills.lockUpdated` and prints `Skills lock updated: <name>` for each name.

The history plan names every file it will remove before approval. Each pending
unit carries sorted repository-relative `paths` in JSON: the files a removal
or conversion deletes, the Archive Records `record-links` rewrites, or the live pages `live-links`
rewrites. Text
prints the same paths in the same order as `  - <path>` lines directly under
each unit's summary.

When Pending History exists, the managed refresh adds a `History` section to
the plan. Text lists each planned history unit and its result, then prints
`History: <status>`, `History units:`, `History applied:`, `History refused:`
and `History tag:` lines when those values exist. Refusal reasons are one line.
Uncommitted changes and removed paths held by neither the History Full Tag nor
a commit the default branch reaches are refused by name and left untouched; a refusal does not block the Baseline Plan
or change the update's state. A repository-wide history precondition such as an
invalid tag makes the history section `blocked` with a one-line result.

The update lists each History Citation in text and JSON, and binds the list into
the Plan Digest:

```
History citations: <n>
- cites <path>:<line> names <target>
- cites <path>:<line> names <target> (live)
- cites <path>:<line> names <target> (live; promote with --promote <target>)
- cites <path>:<line> names <target> (promoted to docs/references/<basename>)
History stale citations: <n>
- cites <path>:<line> names <target>; history-full holds it: git show history-full:<target>
```

JSON includes `history.citations` and `history.staleCitations` arrays.
Citation objects carry `path`, `line`, `target` and `live`, with
`offer` for the Live Document promotion offer and `promotedTo` after a
promotion. A rewritten citation also carries `rewrite`. Stale History Citations are reported outside the digest, state,
category and exit code.

The repeatable `--promote <path>` flag copies a selected Legacy Archive
Folder's named file to `docs/references/` during apply. Each value must be a
clean repository-relative file inside a selected folder, outside Refused Units,
non-core, unique by basename among promotions, and absent from
`docs/references/`. The update refuses `--promote` with `--no-history`,
invalid values, refused units or occupied destinations with category `invalid`
and writes nothing.

After removals the update rewrites every History Citation it lists, in every Citation Source: each tracked text file outside the skill trees and testdata directories. A prose path that starts with ./ or ../ is resolved against the directory of the citing file. Outside Markdown the Git Locator is written without backticks. A removed link becomes its label followed by a Git Locator; prose paths become locators too. A Reduced History Entry uses its provenance path and commit when that commit remains reachable from the default branch or a tag. The rewrite reads only this repository and writes only the planned citing files.

History units are Legacy Archive Folder conversions, then `findings`,
`backlog`, `adr`, `reviews`, `handoffs` and `superseded` removals, followed by `record-links` and `live-links` repairs. They include
Reduced History Entries, legacy family trees, terminal active Findings and
Backlog Entries, retired ADRs and finished Review Artifacts the Baseline used
to relocate. A superseded ADR naming an accepted successor leaves, and the plan
appends a missing lineage line to that successor. An ADR without an accepted
successor is refused by its own path while the other ADRs leave. The
`superseded` unit holds only uncited Archive Records without
regeneration grants. Its citation check covers every surviving tracked Markdown
file outside skill trees and `testdata`, including bounded slugs in front matter.
Legacy Spec folder relocations and the live review root rename still apply.

After removal units, `record-links` repairs each surviving Archive Record's
relative links to missing files and to itself. It tries the record-relative
path, then the record's original `source` folder; the first reachable commit
holding the target wins. Before that, a link whose target, resolved against
the `source` folder, is a file whose copy `docs/references/<file>` the record
lists under `promoted` and the tree holds is rewritten to a relative link to
that copy, and the plan names `docs/references/<file>` as its rewrite
(ADR-0308). A self-link becomes the Git Locator of the source
folder. When neither path has a qualifying commit, the record-relative path
and its pull request are used. Code spans, fences, URLs and Git Locators are
not relative links. A record removed by the same plan is skipped.

The repair line is `- record-links: rewrites <l> link(s) in <n> Archive Record(s)`.
A JSON `rewrite` unit carries `links`; `History units` file and byte totals
count only files that leave the tree.


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

Removal lines read `- <unit>: removes <n> file(s) (<b> bytes)`. A rewritten
citation appends `; rewrites to <replacement>`. `History applied:` drops the
old reduction count and appends `; rewrote <c> citation(s) in <m> file(s)` after
the History Full Tag clause. JSON unit actions are `convert`, `remove` or `rewrite` and
omit `bytesAfter`; `applied.reduced` stays 0 and the applied object adds
`rewrittenCitations` and `rewrittenFiles`. The Plan Digest binds the removals
and each citation's `rewrite`, with no `bytesAfter`. A plan writes nothing and
a second update reports the applied units current.

Adopter note: if an earlier `roundfix baseline update` or Sanitize Batch
already broke citations, run `roundfix baseline update` with the History Full
Tag present, read the `History stale citations` lines, and restore a target
you keep with `git show history-full:<path>` into `docs/references/` before
fixing the link. Specs already written to link an owner's copy need nothing,
because the owner's archive performs an Archive Relink.

JSON always includes a `history` object whose status is `skipped`, `current`,
`pending`, `applied` or `blocked`. It reports selected units, refusals, tag
action and apply counts when those values exist. Selected history is bound into
the same Plan Digest as the Baseline Plan, so `--yes` or `--confirm-plan
<digest>` approves both. Apply runs the Baseline Plan first, creates the
annotated `history-full` tag at `HEAD` when absent, and converts every selected
unit. The update never moves or pushes the tag; push it with
`git push origin history-full`.

Pass `--no-history` to leave the history section out and use the History
Sanitize Command's reviewed-batch procedure separately.

While Pending History remains, an update with --no-history still writes the retirement rule into the guide, so SC-HISTORY-RETAINED reports every retained history file on each active Spec, and the path to green is the history conversion: rerun roundfix baseline update without --no-history, or apply roundfix history sanitize batches until nothing is pending.

#### Starter Files

A Starter File is a repository tooling file offered by a Baseline module.
The `python` profile offers `Makefile`, `.roundfixrc.yml` (Project Config)
and `.github/workflows/verify.yml` (CI workflow).

```bash
roundfix baseline starters --profile python
roundfix baseline starters --profile python --confirm-plan <digest>
```

The first command previews each path as `create` when absent or `keep` when
anything already exists there, and prints a Plan Digest. It writes nothing.
A non-empty preview exits 3; an empty preview reports `current` and exits 0.
Confirm the exact printed digest with `--confirm-plan` to create only the
absent files. An existing path is never overwritten, compared or removed.

Confirming the exact preview is the maintainer's express authorization for
exactly those paths. An Agent runs the confirmation only with that express
authorization. A stale digest is refused with exit 3 and nothing written;
preview again. If a write fails, the command removes every file it created
before returning exit 1; directories it created may remain.

Once written, the repository owns each Starter File. It is outside the
Change Plan and Setup Manifest, and no refresh changes it. JSON output uses
`roundfix/baseline-starters/v1`; invalid arguments, an unknown profile or an
unsupported format exit 2.

### Retired skills

A Retired Skill is a skill the Baseline no longer requires although the
upstream catalog may still list it: `council` and `the-fool` since 0.52.0.
No module requires or dispatches it, and the asset sync drops it from every
Setup Snapshot. `roundfix baseline update` never deletes an installed copy.
When the repository still holds one, the preview and the applied update list
it under `Skills retired` (`skills.retired` in JSON) with the paths to
delete; the state and exit code do not change, and `--no-skills` skips the
list. To remove the copies:

```bash
git rm -r -q --ignore-unmatch .agents/skills/council .agents/skills/the-fool .claude/skills/council .claude/skills/the-fool
rm -rf .agents/skills/council .agents/skills/the-fool .claude/skills/council .claude/skills/the-fool
```

Then delete the `the-fool` entry from `skills-lock.json` by hand, because the
skills CLI's `remove` can leave a project lock entry behind, and commit. The
next `roundfix baseline update` lists no retired skill.

`baseline skills reconcile` removes only lock entries proven absent from one
source repository at the exact 40-hex commit passed with `--revision`, and only
when the selected Profile does not require them. A non-empty preview exits `3`
with its Plan Digest; `--confirm-plan` applies that exact plan and retains every
installed skill tree. When a Profile-required skill is absent at the selected
revision, the command blocks with exit `3` and finding
`reconcile.required-removed`, prints `plannedChanges: []` and
`planDigest: null`, and writes nothing. Mutable revisions are refused before
source acquisition. Reconciliation reads `skills-lock.json` after the source
is acquired. If the lock changes during planning before its transaction
preimage is captured, the command refuses with `lock.changed-during-plan`,
exits `3`, and writes nothing. Doctor remains offline and read-only; it never
reconciles or edits the lock.

`baseline assets sync` is a maintainer operation over an explicit canonical
setups directory. `--check` is read-only. Refresh validates the generated
catalog before updating only Go-owned canonical setup snapshots.

Baseline requested output goes to stdout; diagnostics and progress go to
stderr. Exit categories are `0` success or current no-op, `1` execution,
verification, output, recovery, or incomplete-rollback failure, `2` invalid
input/schema or unsafe repository, `3` another owner action or renewed
approval is required, and `130` cancellation. Baseline reports repository
formatter and Verification commands as recommendations but never runs them.

For the adoption, automation, Decision Document, cross-clone, migration,
recovery, and security procedures, read
[CONTEXT-driven development](../context-driven-development.md#adopt-or-update-the-context-driven-baseline).
