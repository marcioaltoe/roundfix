---
status: accepted
knowledge: true
created_at: 2026-11-24T00:00:00Z
updated_at: 2026-11-24T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Round 2 of a pre-PR review opens a new session when the configured model changed

## Context

On 2026-11-23 Fiscus reported that round 2 continued a reviewer session after
its review profile changed from `gpt-6-luna` to `gpt-6.1-sol`. Codex put a
model-change notice before `No findings.`. Roundfix correctly blocked that
answer as unclassifiable: a no-findings verdict cannot account for other
content. The lineage recorded the session name and selection index but no
model, so it could not prove that the configured model still matched.

## Decision

This decision amends [ADR-0197](0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md)
in part. Round 2, outside a Correction Review, prepares the open round-1
session on its recorded selection index only when that index is inside the
review profile, the record holds a model, and the selection still names that
model, compared as exact strings after trimming. Only the model is compared.

Otherwise Roundfix prepares nothing on the recorded name: it ends that session
before preparing any new one, then tries the normal selection loop from the
preferred selection with fresh random names. This includes a changed model,
an older record without a model, and an index outside the profile. Failure to
end the session returns its existing error before any prompt; preparation of
a new session retains the existing retry and fallback behavior. A failed
preparation of the recorded session retains its existing fresh-session path.

The new session receives the same round-2 delta prompt, with round-1 findings
and dispositions. It remains round 2 and records `continued: false`; the
two-round ceiling and Correction Review rules stay unchanged.

Every round that sends a prompt records the receiving session's runtime model
in the optional `lineage.model` field, after `lineage.selection`. This is the
preferred model or the fallback actually used, the same model for a continued
session, the new model for a fresh round-2 session, and the Correction Review's
own model. The field appears in the printed record, checkout record file and
Reviewed Head Record. A record written without a prompt gains no model,
including a `ceiling-closed` record. Older records remain readable without a
migration, and older readers can ignore the new field.

The classifier remains unchanged under
[ADR-0174](0174-a-pre-pr-review-reads-the-final-message-and-keeps-a-record-per-checkout.md).
No output line, flag, exit code, help text or Project Config key is added.

## Considered options

- Remove or ignore the runtime's notice before classification: rejected because
  a fresh session covers the reported case while preserving strict verdicts.
- Set the recorded model back on the continued session: rejected because that
  runs the review on a model the operator no longer configures.
- Compare reasoning effort or the whole Agent Selection: rejected because the
  reported notice concerns the model alone.
- Restart the Reviewer Lineage at round 1: rejected because it rereads the whole
  candidate and resets the two-round ceiling.
- Continue a session whose record holds no model: rejected because sameness
  cannot be proven.

## Consequences

An open lineage written by an earlier Roundfix loses continuation once. The
delta prompt carries the findings and dispositions either way. Different names
for the same model count as a change; a changed reasoning effort alone does
not prevent continuation. Other notices beside a no-findings verdict still
block and remain visible.
