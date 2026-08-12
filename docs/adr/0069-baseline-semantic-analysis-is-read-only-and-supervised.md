---
status: accepted
created_at: 2026-07-24T21:27:41Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# ADR-0069: Baseline semantic analysis is read-only and supervised

Ambiguous root-instruction classification and free-form plan suggestions use
Codex `gpt-5.6-sol` with `xhigh` reasoning after Exact Agent Selection Proof.
If it is unavailable or returns an invalid or incomplete result, Roundfix
discards the output and restarts analysis from the same immutable snapshot with
Codex `gpt-5.5` and `xhigh`; if both attempts fail, the maintainer classifies
the content manually. ACP output is always a proposal that requires
deterministic validation, consolidated human review, and Change Plan
confirmation, so this read-only retry does not continue Agent work over
possibly modified state and does not weaken ADR-0050. Each attempt receives the
same bounded, digest-identified JSON snapshot in a private empty directory,
with terminal and tools denied and no checkout access. Raw ACP output remains
ephemeral; only a complete, schema-valid proposal can enter the deterministic
review.
