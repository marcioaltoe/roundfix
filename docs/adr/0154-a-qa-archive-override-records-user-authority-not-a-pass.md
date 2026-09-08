---
status: accepted
created_at: 2026-09-08T19:34:57Z
updated_at: 2026-09-08T19:34:57Z
deprecated_at: null
superseded_by: null
---

# A QA archive override records user authority, not a pass

The maintainer requires Spec archival with a QA override when explicitly
requested or authorized. The existing archive skill already describes an
archive-anyway exception, while the canonical rule and Archive Command require
QA eligibility. The canonical policy now permits the exception and the runtime
must expose a supported, auditable path for it.

An explicit request or applicable prior authorization identifies the covered
Spec and unmet QA prerequisite. Preserve its approval source/date and actual
QA outcome or absence, stamp `qa_override: true`, and retain supplied rationale.
Do not require a new confirmation for an already applicable grant. A generic
archive request alone does not infer an override, and enabling the capability
is not blanket authorization for any existing Spec.

The exception can waive the terminal QA Task's archive completion/evidence
requirement without rewriting its status, Result or report. Non-QA Tasks must
still be complete and adopted references self-contained. Structural errors,
unsafe destinations and execution ownership remain separate checks. A normal
declined-gate or qualifying declared-partial case retains its own eligibility
rules and does not acquire an invented override.

The consequence is a truthful archived record whose QA remains unproved or
failed as recorded. It is not a successful Run or an authorization to publish,
merge, release or skip required checks. Any later delivery action must satisfy
its own authority and gates. CLI flags, authorization schema and archived
metadata details beyond the existing marker remain implementation work;
the current command does not yet support the new path.
