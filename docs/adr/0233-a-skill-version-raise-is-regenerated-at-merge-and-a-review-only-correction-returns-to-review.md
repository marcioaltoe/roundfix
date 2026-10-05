---
status: accepted
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A skill version raise is regenerated at merge, and a review-only correction returns to review

On 2026-10-04 four queued Specs raised the Roundfix Skill version from the
tree each Task started on. Two Runs computed `0.1.26` while neither item had
merged, and a third item recorded `0.1.25` over the default branch's
`0.1.27`. The operator resolved each Pull Request conflict by hand, taking the
default branch's version lines and raising to the next patch. A Codex Task
agent had also written a version digest by hand, and the record command,
which never replaces a digest, then failed every later attempt. The same day
the pre-PR review of an archived candidate raised two findings on the archived
Spec's own records; the operator fixed one in a docs-only commit, dismissed
the other with evidence, and the retry refused with
`corrective-spec-required`, so the operator ran review round 2 and opened the
Pull Request by hand.

Three rules settle it:

- The owned-skill record command chooses the version. When a skill's content
  is not recorded under its declared version and that version is not higher
  than every recorded one, the command writes one patch above the highest
  recorded version into both version fields of the canonical `SKILL.md` and
  its mirror, then records that content. It still never replaces a recorded
  digest, and the check without the flag still refuses unrecorded content.
- A derived-path declaration may name line-scoped paths with a line pattern.
  At a Pull Request conflict, a conflicted line-scoped path whose every
  conflict hunk has the same number of lines on both sides and differs only on matching lines (identical lines between them are kept) takes the
  default branch's side of each hunk; any other hunk keeps it a source
  conflict. The
  regeneration may then change a line-scoped path only on matching lines. This
  repository declares the two version fields of every `SKILL.md` as line-scoped
  for the record command, so two items that raise the same skill merge with
  the next free version instead of parking.
- A Delivery Retry of a `corrective-spec-required` item whose head moved
  returns it to `reviewing` when Git proves the head descends from the parked
  candidate, every standing finding of the review recorded at that candidate
  has one disposition (fixed by a commit the head contains, or dismissed with
  evidence), and every path changed from the candidate to the head lies under
  an archived Spec the blocker names. The next review is round 2 of the same
  Reviewer Lineage. Anything else still requires a corrective Spec.

Refreshing each item from the default branch before its Run was rejected: the
two Runs that recorded `0.1.26` both started from a default branch already
at `0.1.25`, so the collision arises between unmerged items and only a merge-time regeneration
removes it. Declaring every `SKILL.md` a derived path was rejected because a
conflict confined to declared paths takes the default branch's whole file and
would drop the item's skill edit.

This decision refines
[ADR-0189](0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md),
whose record still only adds a higher version;
[ADR-0192](0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md),
whose whole-file rule is unchanged for declared paths;
[ADR-0165](0165-a-blocking-review-after-archive-parks-publication-for-a-corrective-spec.md),
whose corrective Spec is still required for a finding about the delivered
change; and
[ADR-0197](0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md),
whose round 2 reviews the correction. It supersedes none of them.

## Consequences

A version record entry an item wrote for its own unmerged version is replaced
whenever a conflict regenerates the record from the default branch's bytes;
without a conflict, a skipped version may stay recorded. A line pattern that
also matches a body line resolves that line from the default branch; the
repository gate after the merge is the net, as for any derived path. The
review command still prints its corrective-Spec line for findings on an
archived Spec, because only the retry knows what the correction changed.
