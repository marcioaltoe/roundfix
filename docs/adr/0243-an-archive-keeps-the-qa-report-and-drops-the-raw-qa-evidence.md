---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An archive keeps the QA Report and drops the raw QA evidence

On 2026-10-06 the History Root held 50.9 MB in 5,448 files, and 21.0 MB in
2,713 of those files were raw QA evidence under a Spec's `qa/evidence/`:
transcripts, JSON dumps, harness scripts, fixtures, screenshots and SQLite
files. A Jev judgment of a stratified sample of the same tree classed raw
evidence as transient, read again by nobody, and every QA Report, PRD,
TechSpec, Task Graph and Task file as a repository record. Of the files
the History Root gained in the ten days before, 1,440 of 2,217 were raw
evidence. ADR-0215 decided that "A removal is its own later change, and the
maintainer's explicit approval of that removal is recorded in it". On
2026-10-06 the maintainer approved that change ("Sim, como Spec").

The Archive Command now cuts the raw evidence as part of a normal archive:

- **What goes.** Everything under the Spec's `qa/evidence/` directory, of
  any type. Nothing else in the Spec is cut: QA Reports and every other file
  under `qa/`, `references/`, `measurement/`, `design/`, the PRD, TechSpec,
  Task Graph, Task files, `_authorization.md` and review files all stay.
  Binaries and databases outside `qa/evidence/` stay too, because the only
  ones in the History Root are adopted references and design images.
- **What stays in its place.** An Evidence Manifest, `qa/evidence-manifest.md`.
  It records the date of the cut, the revision whose tree held the files,
  the repository-relative Spec directory they lived in at that revision,
  and, for each dropped file, its path, size in bytes and SHA-256 digest.
  The bytes stay in Git history, and the manifest says where to find them
  and how to recognize them.
- **Links.** A relative Markdown link in the Spec that reaches into
  `qa/evidence/` is rewritten to reach the Evidence Manifest, keeping its
  text. A path written as prose or code is not a link and is left alone;
  the manifest lists it.
- **When.** The cut happens in the same archive operation, after
  eligibility and before the move. An active Spec keeps its evidence, so QA,
  review and settlement read it unchanged until the archive.
- **Which archives.** A normal archive on a QA verdict cuts. An archive made
  with the QA Archive Override does not, because ADR-0154 keeps that
  record's QA files as observed. A superseded Spec without a Task Graph
  moves unchanged, as before.
- **The Delivery Queue.** An archive commit is still an exact Spec move when
  its only other changes are the stamp, the link rewrites of ADR-0230, the
  evidence links that now reach the manifest, the dropped evidence files and
  the manifest that lists exactly those files with matching sizes and
  digests.

Already-archived Specs are cut only on request, one Spec at a time, with
`roundfix archive <slug> --drop-evidence`. It is a dry run by default and
writes only with `--apply`. It refuses an active Spec, keeps a QA Archive
Override or superseded Spec unchanged, reports nothing to do for a Spec
already cut, and refuses while a file other than Markdown, outside the
History Root and the Spec roots, names the evidence directory. It never runs
implicitly, never commits, and nothing in the repository runs it on the
existing History Root.

No archived evidence may be read by a test, fixture or build step. The
Baseline already prohibits reading a file under the Spec root. The
coverage record had listed three Go packages inside archived evidence, and
a repository-copy helper failed on a tracked file that was missing from the
working tree. Both are fixed before the cut ships.

We rejected three alternatives. Keeping evidence and excluding it from the
Secondbrain mirror alone (done separately on 2026-10-06) leaves the
repository growing by evidence with every archive. Cutting after a number of
days needs a clock and a sweep that nothing owns, while the archive is
already the moment a Spec stops being live work. Rewriting Git history to
reclaim old evidence bytes was declined by the maintainer on 2026-10-06
("Não agora").

## Consequences

- Removing a file from the tree does not shrink Git history. The cut stops
  the working tree and every later clone's checkout from carrying evidence,
  and the Secondbrain mirror from receiving it. The pack keeps every byte
  already committed.
- A reader who needs a dropped file runs `git show` at the recorded revision
  and checks it against the recorded digest. Evidence that was never
  committed cannot be recovered, and the delivery flow commits QA evidence
  before it archives.
- The QA Report becomes the durable record of QA. A conclusion that lives
  only in an evidence file is lost to a reader of the archive, so the gate
  writes each row's observed result in the report. Material meant to be
  reread belongs in `references/`, `measurement/` or `docs/references/`.
- ADR-0215 is applied, not superseded: this Spec is the later change, and it
  records the maintainer's approval. ADR-0230 is extended: the Delivery
  Queue also accepts the evidence cut as part of an exact move. ADR-0154
  stands: an override archive keeps its QA files. ADR-0223's refusal for a
  file that names an active Spec's directory gains a counterpart for the
  archived evidence directory.
- The Baseline's rule that completed or archived legacy Specs stay
  byte-identical gains one exception, the runtime's evidence cut.
