---
status: accepted
created_at: 2026-10-08T00:00:00Z
updated_at: 2026-10-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Baseline update sanitizes pending history in the change it plans

ADR-0248 made `roundfix history sanitize` an operator procedure: tag the full
history, then apply batches through Pull Requests. Adopters reach it only by
reading release notes. On 2026-10-08 Fluxus still had 57 pending units (23.9 MB)
after upgrading to 0.59. The same day the maintainer chose "No baseline update
(Recommended)": after `roundfix upgrade`, the repository's
`roundfix baseline update` detects the **Pending History** and includes the
sanitize in its own plan. `roundfix upgrade` only says so.

**One plan, one digest.** The Managed Refresh plans the pending units the same
way the History Sanitize Command does, through the same Lenient Legacy Reading,
records builder and reductions. The update's Plan Digest binds the history
section together with the Baseline Plan whenever the section selects a unit.
When it selects none, the digest is the Baseline Plan Digest byte for byte.
`--confirm-plan <digest>` and `--yes` approve both parts together. A second
update over the result reports the history `current`.

**All units at once.** The update converts every unit it can convert in the
change the operator reviews, because a Managed Refresh is already one reviewed
change and it has no batch flag. An operator who wants ADR-0248's reviewed
batches passes `--no-history` and uses `roundfix history sanitize --apply
--batch <n>`, which is unchanged.

**Refused Units never block.** A unit the update cannot convert is a Refused
Unit. It is listed with its reason on one line and left untouched, and it does
not change the update's state or exit code, as a Retired Skill does not. The
Baseline part applies. A unit whose paths hold uncommitted changes is refused
with that reason, because its removed bytes would not be in Git. So is a unit
whose paths an existing History Full Tag does not hold. The rest of the
working tree keeps the update's existing contract: preimages bind it, and no
clean tree is required.

**The tag is created, never moved, and never pushed.** When `history-full` is
absent and a unit is selected, apply creates the annotated History Full Tag at
`HEAD` before converting anything. The plan names that tag and its commit, and
both the plan and the result say it must be pushed with
`git push origin history-full`. An existing tag is never moved or replaced. A
lightweight tag or one that is not an ancestor of `HEAD` blocks only the
history section, with a reason. Baseline update still never commits or pushes.
A tag is a local ref, and the operator pushes it.

**Order of writes.** The Baseline Plan applies first, because its refusals are
the strictest. A refusal there writes nothing, not even the tag. Then the tag is
created and the units are converted. A history write failure exits 1, names the
unit and keeps the Baseline result. It tells the operator to restore the
History Root with `git restore` and `git clean`.

The upgrade notice counts the Pending History in process from the filesystem
of the working directory's repository and names it by kind. Outside a
repository it says to run `roundfix baseline update` in each adopted
repository. It never changes the upgrade's output or exit code.

## Consequences

Two alternatives were rejected. Running the sanitize from `roundfix upgrade`
would write into whatever directory the operator happened to upgrade from,
outside any reviewed plan. Keeping the sanitize a separate manual procedure
leaves adopters with 99% of a history they were told is irrelevant, or with a
four-step release-note recipe. This narrows ADR-0248's "one Pull Request per
batch" for adopters who sanitize through the update. Their one batch is the
update's change. ADR-0100's preservation proof still covers every instruction
carrier the Baseline Plan touches. The history section changes only files
under the History Root, as History Relocations already do inside the same
plan. The History Sanitize Command keeps its contract, and its
Refused Unit reasons now print on one line.

Under Lenient Legacy Reading, a legacy `unproven` list of maps is accepted.
Each map becomes one stable text line in the Archive Record: `row <row>: <claim>`
followed by the other keys in lexical order. Active Specs keep the string-only
list.
