---
status: accepted
created_at: 2026-10-01T00:00:00Z
updated_at: 2026-10-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The judge suggests which sources share a Spec, and never decides it

ADR-0208 asks an author to adopt sources that share a context into one Spec,
and the maintainer asked that Jev help decide it. On 2026-10-01 one Noul
question was measured against this repository's archive: "Should both be
resolved by the same specification, delivered together?", asked of 94 pairs
of sources that one archived Spec adopted together and 94 pairs that
different Specs adopted within three days of each other. Jev 1.13, through
OpenRouter, reached AUROC 0.783 against 0.753 for a TF-IDF baseline, a
difference whose 95% interval, [-0.032, 0.093], includes zero. At a
probability of 0.3 or more it was right on 93% of the pairs it raised and
found 27% of the true pairs; the TF-IDF baseline's 27 best pairs were right
on 85%.

`roundfix spec judge` therefore asks that question, as measured, for every
pair of one source the Spec adopted and one source still open, and raises the
pair as a suggestion at a probability of 0.3 or more:

- **A suggestion, never a decision.** A suggested pair is answered by
  adopting the open source or by stating why it stays apart. It never gates,
  never moves a file and never changes an exit code (ADR-0200).
- **Low recall, stated.** Most pairs that belong together are not raised. An
  author who sees no suggestion has not been told that nothing fits.
- **The data boundary grows by one reader.** ADR-0201 kept Findings and
  Backlog Entries out of requests because no judgment needed them. The
  grouping question needs them, and the maintainer's authorization of
  2026-09-30 covers them: "trechos de PRD, TechSpec, Task files, ADRs,
  findings e Backlog Entries DESTE repositório". One reader accepts only a
  Finding or Backlog Entry the Spec adopted, an `open` Backlog Entry and an
  unresolved Finding of this repository, as regular files directly in their
  directories. An Inbox Entry, a terminal record, source code and anything
  outside the repository cannot reach a request. Everything else ADR-0201
  decides stays: the recipients, the keys, the Judge Log and the monthly
  ceiling.

## Consequences

- A run of the judge on a Spec that adopted sources asks one question per
  pair, so its cost grows with the open backlog. At the measured US$0.00004
  per pair, 100 pairs cost less than half a cent, under the same ceiling.
- A lexical baseline nearly as good overall was not adopted, because at the
  operating point Jev raised fewer wrong pairs and the maintainer asked for
  Jev. Moving the threshold or the model means measuring again.
