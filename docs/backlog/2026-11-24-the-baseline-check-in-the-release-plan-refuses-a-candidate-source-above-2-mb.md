---
type: fix
status: open
created: 2026-11-24
spec: null
reason: null
---

# The Baseline check in the release plan refuses a candidate source above 2 MB

## Symptom

On 2026-11-23, about 18:00 BRT, the release plan of Vortex v0.69.0 (gesttione-solutions/vortex) failed its Baseline verification with `HTTP candidate source bytes exceed 2097152`. The release went on. Expected by the sender: the verification reads the candidate source without exceeding the limit, or the limit and the message say what to do.

## Open question

Evaluate whether the 2 MB per-HTTP-source limit is too low for the current Baseline, and decide whether the fix belongs to Roundfix or to the repository. Not measured here: the size of the Baseline source today and where the limit is defined.

## Follow-up owed

Answer the Vortex project (`default/gesttione-solutions-vortex`) with that decision. The sender is not running `baseline update` until the answer comes; no urgency.

## Source

Cross-project message `01a20b17-73b6-7481-bbc2-be217b275724` from the Vortex project, 2026-11-24. Nature: reported. Recorded by triage on 2026-11-24.
