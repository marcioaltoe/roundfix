---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An archive leaves an Archive Record, and the Spec folder stays in Git

`docs/history` held about 51 MB in 5,400 files on 2026-10-06 and grew about
22 MB in ten days. A Jev judgment of 123 sampled artifacts rated about 8% of
that text as knowledge worth reusing, 60% as records useful only inside the
repository and 31% as throwaway evidence. On 2026-10-06 the maintainer
decided: "Não quero nada no histórico que não seja relevante para o
secondbrain". The archive therefore no longer moves a Spec folder into the
History Root. It writes one **Archive Record**, `<archive-root>/<slug>.md`, of
about 2 KB, and removes the Spec folder in the same change. That change is the
archive commit, or the working-tree change the operator commits. The folder's
bytes stay in Git at the commit the record names as `source_revision`.

The record holds what readers need without the folder:

- identity: the slug, title, creation and archive dates and the disposition
  (`pass`, `partial`, `qa-override` or `superseded`);
- the QA Task, QA Report and verdict, and the `unproven` actions;
- every QA Archive Override field, the reason included;
- the supersession target;
- the ADRs the PRD's Decisions cite and the adopted sources by name;
- the sanctioned regeneration commands of its authorization record;
- the files promoted upstream;
- a one-paragraph outcome.

Only the outcome paragraph is shortened to keep the record near 2 KB. The
archive refuses when the Spec folder differs from `HEAD`, so the recorded
revision holds exactly what was removed.

Knowledge moves upstream before the cut, and a person or Agent decides what
moves. `roundfix archive <slug> --plan` lists what the archive would remove.
When a judge key is set and the monthly ceiling allows, Jev advises on each
file outside the core artifacts and outside `qa/evidence/`. The advice comes
through the Spec judge's stage key, ceiling and log, and Jev answers
`reusable_knowledge`, `repository_record` or `transient_evidence`. The advice
never refuses, never moves a file and never decides an archive. An operator
promotes a file with `--promote <path>`, which copies it to `docs/references/`
in the archive change. A cross-project lesson goes to the Secondbrain inbox
under the Secondbrain guide. Inside the Delivery Queue, the QA Task copies a
file upstream before the archive.

Readers work from the record or from Git, never from a folder that may be
gone:

- The Delivery Queue's exact-archive proof accepts a commit that removes
  exactly the Spec folder and adds a record that agrees with the removed
  `_prd.md` and names the parent as its `source_revision`. Promoted copies
  must equal removed blobs.
- Merge evidence holds when the default-branch head carries the record. Task
  completion follows from the disposition.
- Reads that need a Task file, a TechSpec or a QA Report use the Run's or the
  candidate's own history at the record's `source_revision`, which a squash
  merge keeps on that branch and in the Pull Request's head ref.
- Repository tests that characterize the archived corpus read it from Git at
  the pinned commit `40a7893d872c8a6705f6d7745e6efe430ae9deeb`, the v0.49.0
  release on `main`, which holds every Spec folder archived through Spec 0237.

Folders archived before this decision stay readable until a separate,
batched migration replaces them. That migration needs the maintainer's
approval, given on 2026-10-06 as "Aplicar em lotes via PR", after a tag marks
the last commit with the full history.

Two alternatives were rejected. Cutting only `qa/evidence/` behind a
manifest, as the parked draft Spec 0238 proposed, would keep about 60% of the
bytes the maintainer called irrelevant. It would also build a manifest and
link rewrite the full cut then removes. Letting Jev decide what to keep would
make a probabilistic judge gate an archive. On the sampled `references/`
files Jev's mean confidence was 0.57, and the judge fails open when no key is
set.

## Consequences

This narrows ADR-0215's "Nothing archived is removed" for Specs archived
from now on: the maintainer's approval of the removal is the decision quoted
above. ADR-0230's link pass and ADR-0121's identity ledger still serve
folders archived earlier and the Baseline's history layout. A record has no
relative links to rewrite. A Spec's delivery Pull Request and squash commit
do not exist when the archive commit is written. Readers find them as the
default-branch commit that added the record, and the migration fills them in
where it knows them. The `### QA settlement` tables of the archive-spec,
qa-gate and Roundfix skills change identically. A passing, partial or
override archive leaves an Archive Record, with the override's reason when
overridden, instead of keeping the Spec, its QA report and evidence.
