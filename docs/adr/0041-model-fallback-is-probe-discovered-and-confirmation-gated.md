# ADR-0041: Model fallback is probe-discovered and confirmation-gated

Status: Accepted

A configured Agent Model can fail selection Preflight Validation while the
runtime has other working models, and switching models changes token
consumption in ways the user must own. On a selection failure Roundfix
therefore probes the runtime's Model Catalog newest-first with disposable
Agent Sessions and offers the newest proven Agent Model at its highest
proven reasoning effort — but starts work only after explicit human
confirmation: an interactive prompt, or, in non-interactive contexts, a
Preflight Validation failure naming the exact explicit-flags re-run. An
`--allow-fallback` flag, a config key, or any autonomous fallback was
rejected because a standing pre-authorization silently converts a model
outage into an unplanned spend; static fallback lists were rejected per
ADR-0039. Refines ADR-0037's no-fallback stance: fallback exists, but only
through a human decision.
