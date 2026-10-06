### baseline

```text
roundfix baseline [--repo <path>] [--format <text|json>]
roundfix baseline plan --profile <id> [--decision <id=value> ...] [--decision-file <path> ...] [--repo <path>] [--format <text|json>]
roundfix baseline apply --plan <file> --confirm-plan <digest> [--repo <path>] [--format <text|json>]
roundfix baseline capabilities check [--profile <id>] [--repo <path>] [--format <text|json>]
roundfix baseline profile init --id <id> [--from <built-in-id>]
roundfix baseline profile show <id> [--format <text|json>]
roundfix baseline profile validate [<id>|<path>] [--format <text|json>]
roundfix baseline skills restore --profile <id> [--skill <name> ...] [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
roundfix baseline skills reconcile --profile <id> --source <owner/repo> --revision <commit> [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
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

The managed refresh (`roundfix baseline update`) also restores a required
external skill that matches its lock but trails its Setup Snapshot. Its
read-only preview lists the skill under `Skills drifted` with a restore action,
even when managed guidance is current. Rerun with `--yes` or the preview's
`--confirm-plan` digest to restore it through the existing immutable snapshot
restore path. `--skills-source-dir` supplies an offline source commit as for
other external restorations; `--no-skills` skips the comparison and restore.

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
