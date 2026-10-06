---
type: feat
status: promoted
created: 2026-10-06
spec: 0242-an-archive-that-leaves-an-archive-record
---

# History keeps only what the Secondbrain needs

## Problem

`docs/history` holds about 51 MB of file content in 5,400+ files and grew
about 22 MB in ten days (measured 2026-10-06). Almost none of it is knowledge
anyone will reuse:

| Kind under `docs/history/specs` | Size | Files |
| --- | ---: | ---: |
| `qa/evidence/` (raw transcripts, JSON, logs, scripts) | 21.05 MB | 2,713 |
| `task_*.md` | 8.64 MB | 1,267 |
| QA reports | 6.22 MB | 350 |
| `references/` (mostly adopted Backlog Entries and Findings) | 5.48 MB | 216 |
| `_techspec.md` / `_prd.md` / `_tasks.md` / `_authorization.md` | 8.55 MB | 745 |
| `design/`, `reviews/`, `baseline/`, `measurement/`, other | 3.9 MB | ~140 |

Outside `specs/`: `reviews/` 2.5 MB, `findings/` 0.8 MB, `backlog/`,
`handoffs/` and superseded `adr/` about 0.35 MB.

On 2026-10-06 the Jev judge (`typesafe/jev-1.13-20260917`, 123 sampled
artifacts, report `~/.roundfix-operator/history-value-2026-10-06.md`) rated
about 8% of the history text as reusable knowledge for the Secondbrain, 60% as
records useful only inside the repository, and 31% as throwaway evidence. The
reusable 8% was mostly a few `references/` and `measurement/` files (model-cost
measurements, the stale-merge CI lesson, an external skill analysis, a pilot
report), and none of it existed upstream until copied on 2026-10-06.

The maintainer's decision (2026-10-06): "Não quero nada no histórico que não
seja relevante para o secondbrain", and the clean-up "deve atuar também no que
já existe no diretório".

`docs/agents/docs-layout.md` already says archived Specs are downstream and
deletable once their durable knowledge moved upstream. The mirror to the
Secondbrain already excludes everything but `references/` and `measurement/`
(`.secondbrain-export`, PR #419, mirror now 7.0 MB). The repository itself
still keeps everything.

## Expected

1. **Target state.** After archive, a Spec leaves in `docs/history` only:
   - a small **Archive Record** per Spec, for example
     `docs/history/specs/<slug>.md` of at most about 2 KB. It holds the slug,
     title, dates, delivery Pull Request and commit, the ADRs the Spec
     recorded, the Backlog Entries and Findings it adopted (by name), and a
     one-paragraph outcome. It exists for provenance and for the machine
     readers listed below.
   - nothing else. Knowledge worth keeping moves upstream before the cut: to
     `docs/references/`, an ADR, the glossary, an agent guide, or the
     Secondbrain inbox.

   The same rule applies to the other history kinds. Reviews, handoffs,
   declined or promoted backlog and findings records become Archive Records or
   are deleted. Superseded ADRs stay, because they are decision rationale.
2. **At archive.** `roundfix archive` classifies every file it is about to
   drop. The Jev judge is advisory, with the stage key and the judge ceiling.
   The archive moves the reusable ones upstream, or opens them as a proposal
   for the operator, then writes the Archive Record and deletes the Spec
   folder in the same archive commit. Override archives follow the same rule;
   the override reason goes into the Archive Record.
3. **Existing history.** A one-time sanitize command does the same for
   everything already under `docs/history`. Examples:
   `roundfix history sanitize` (dry run) and `--apply`.
   - It works in batches, one archive record per Spec, and reports bytes and
     files removed.
   - It never runs implicitly.
   - It goes through a Pull Request per batch, so the diff stays reviewable.
   - Before the first `--apply`, tag the last commit that holds the full
     history (for example `history-full-2026-10`), so everything stays
     recoverable from Git.
4. **Readers fixed first.** These read archived Specs today:
   - the 0227 merge evidence, which needs the archived `_prd.md` on main;
   - the Delivery Queue's archived-retry checks (0219, 0224, 0228, 0232);
   - `speccheck` citations and related-ADR checks;
   - `specaudit`;
   - the docscontract corpus;
   - `docs/references/coverage-record.json` with `TestCoverage`;
   - `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`, which reads the
     0001 transcript;
   - Baseline assets and history-layout code;
   - skills and guides that link into history;
   - ADRs and `CHANGELOG.md` links (about 80 tracked files name
     `docs/history` today; see `docs/references/archived-evidence-measurement.md`).

   Each one must work from the Archive Record or from Git, and
   `make verify` and `make verify-docs` must pass with no Spec folder under
   `docs/history/specs`.
5. **Mirror.** With history reduced to Archive Records, `.secondbrain-export`
   can go back to mirroring `docs/` whole.

## Notes

- Expected result: `docs/history` drops from about 51 MB to well under 1 MB of
  Archive Records, plus whatever moves upstream.
- The Git pack does not shrink. Rewriting history in this repository or in
  the Secondbrain stays a separate maintainer decision ("Não agora",
  2026-10-06).
- Spec 0238 (drop raw `qa/evidence` at archive, with an evidence manifest) was
  authored on 2026-10-06 and, by the maintainer's decision the same day
  ("Incorporar ao saneamento"), is not delivered on its own: its reader
  analysis (coverage record, repo-copy test helper, owned-skill test) carries
  into this entry. Its draft is kept on branch
  `docs/archive-keeps-the-report-not-the-raw-evidence`. It was a strict subset
  of this entry. If this entry is adopted, 0238 should be
  folded into it rather than delivered on its own, so the manifest and
  link-rewrite work is not built and then removed.
- Size the work: this likely needs two Specs. One covers the Archive Record,
  the readers and the archive-time cut. The other covers the sanitize command
  and migrating existing history. Each has at most 4 Tasks plus QA.
- Resolve after the current cycle (0235, 0236, 0237).
