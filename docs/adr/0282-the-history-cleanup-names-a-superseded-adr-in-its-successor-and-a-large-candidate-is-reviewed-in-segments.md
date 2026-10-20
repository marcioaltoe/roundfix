---
status: accepted
created_at: 2026-10-17T00:00:00Z
updated_at: 2026-10-20T09:00:00Z
deprecated_at: null
superseded_by: null
---

# The history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments

On 2026-10-17 `roundfix baseline update` refused the whole `adr` unit in three
of four adopters because one superseded ADR had a successor that did not name
it, and the guide the same update wrote then turned `SC-HISTORY-RETAINED` red
over every retired ADR left behind and over links from Live Documents to files
an earlier layout removed. The same day the Delivery Queue parked Spec 0272 as
`review-blocked`: its candidate held 158 KB of deleted history files and 270 KB
of regenerated goldens, and the review refused the diff before any reviewer
read it. Four smaller gaps of the same cycle cost Runs or reviews: a delivery
convention that asked a QA Task for settlement sections, a plan that edited an
owned skill without its version record, a QA gate whose archive plan sent files
to the judge, and an item binary rewritten in place while a copy still ran.

## Decision

**The cleanup writes the lineage a successor lacks.** When a superseded ADR's
`superseded_by` names an accepted ADR that does not name it, the history plan
removes the superseded file and, in the same plan, appends to the successor one
line, `Supersedes ADR-NNNN (YYYY-MM-DD); Git keeps its text at <Git Locator>.`,
dated with the committer date of the planned `HEAD` in UTC. The line is listed
as a History Citation, so the Plan Digest binds it. A superseded ADR with no
accepted successor stays, refused by its own path with the fix, and the unit
removes every other retired ADR. ADR-0278's rule that a successor names what it
supersedes with its date holds after the plan applies.

**Live Documents get their Git Locators.** A `live-links` unit, after
`record-links`, rewrites each relative link from a Live Document or a
`docs/user-guide/` page to a missing file into the Git Locator that
`SC-HISTORY-RETAINED` already prints for it. Both `history sanitize` and
`baseline update` plan it. The gate itself does not change.

**The review omits what no reviewer should read, and splits what is still too
large.** The pre-PR review diff omits a file deleted under the History Root,
other than an Archive Record, as `retired-history`, and a file matched by an
exact or glob `paths` entry of a `delivery.derived_paths` declaration of the
base commit's Project Config as `generated`. A directory entry is not used,
because a directory such as `docs/agents/` also holds authored files, and the
candidate's own Project Config is not used, because a candidate must not choose
what its review skips. A `codex` round-1 diff still above the review bound is
reviewed as up to four **Review Segments**: consecutive first-parent commit
ranges, each under the bound, sent in order to one review session, whose
findings form one Pre-PR Review Record. A single commit above the bound, a
fifth segment, and a round-2 delta above the bound still block.

**A QA Task settles on its QA Report.** Delivery Convention C2 says the Daemon
writes `## Result` and `## Carry-forward provenance` for non-QA Tasks only and
`## Recorded paths` only when a Task changed an undeclared path, so a QA Task
without them is not a defect.

**Spec check holds an owned skill edit to its version.** A pending non-QA Task
that declares a file of an owned skill without its mirror, both `SKILL.md`
files, the version record and the `-record-skill-versions` step is the gap
`SC-SKILL-VERSION-UNTASKED`, with no horizon, from the Tasks stage.

**The archive plan has an offline form, and an item binary runs from its own
copy.** `roundfix archive <slug> --plan --offline` lists the cut without
Archive Advice and sends nothing to the judge; the `qa-gate` skill uses it.
After the `delivery.item_binary` build, the Delivery Queue writes a copy of the
built file to a new file under the item's artifact directory, renames it into
place, and probes and runs that copy, so a later build that writes its output
in place never rewrites a file a running item binary maps. This extends
ADR-0281's rule from tests to the Delivery Queue.

## Rejected

- **Refusing the whole `adr` unit until a person edits the successor.** It
  keeps every other retired ADR and turns the adopter's check red.
- **Holding back the guide that enables `SC-HISTORY-RETAINED`.** The check is
  right; the plan was incomplete.
- **Raising the review bound.** Spec 0200's 1,295,055-byte candidate failed
  three times with an agent protocol error, which is why the bound exists.
- **Splitting the diff by file or by bytes.** A commit range keeps each
  segment a change its author made, so a finding's anchor names a line of a
  real intermediate tree.
- **Reading derived declarations from the candidate.** A candidate could
  declare its own code generated and skip its review.
- **Skipping the judge only when its key is absent.** ADR-0273 keeps the keys
  out of a Run, yet a QA attempt on 2026-10-16 still reached the judge; an
  explicit flag does not depend on the environment.
- **Removing the built file before the next build.** It leaves a window in
  which the declared path is missing, and it still runs the file a build
  rewrites.

## Consequences

- An adopter whose ADR log has a successor without its lineage line gets that
  line in the change that removes the superseded ADR.
- A segmented review costs up to four prompts in one session; its record lists
  the segments and each finding's segment.
- Every pending Task that edits an owned skill now fails `spec check --strict`
  until it declares the version record and its step.
- Item binary copies accumulate under the artifact directory, one per distinct
  build of an item; this decision adds no reclamation for them.

**Amended by [ADR-0290](0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md) (2026-10-18).** Item Binary copies no longer accumulate past a merge: the Delivery Queue removes a merged item's `item-binary` directory when it releases the item's Runs, and the GC Command reclaims the directory of an item that is merged or no longer queued.

**Amended by [ADR-0299](0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md) (2026-10-20).** The lineage line names a retired ADR whose file name is not `NNNN-<slug>.md` by that file name without `.md`, never by its leading digits, and the numeric forms of `superseded_by` and `supersedes` select only numbered ADRs. The line is dated with the day the ADR was superseded: its `updated_at` in UTC, else the committer date of the newest commit that changed it, and the planned `HEAD`'s committer date only when neither exists. The `live-links` unit leaves a link to a file a pending Task declares in `creates:`, and it resolves a missing link by its file name when exactly one removed history file has that name.
