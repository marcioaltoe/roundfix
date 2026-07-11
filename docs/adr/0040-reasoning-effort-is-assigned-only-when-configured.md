# ADR-0040: Reasoning effort is assigned only when configured

Status: Accepted

Some Agent Models manage reasoning themselves and reject every
`reasoning_effort` value through the ACP adapter (the codex gpt-5.6 family),
so requiring a concrete Default Reasoning Effort makes those models
unusable. Roundfix therefore treats an empty Default Reasoning Effort as a
valid selection meaning the Agent Model manages reasoning: selection skips
the reasoning set call on both the disposable preflight session and the live
Agent Session, and the effective selection records the model-managed state.
A non-empty configured or flag-passed value keeps the ADR-0037/ADR-0039
contract — it is assigned explicitly and any runtime rejection fails
Preflight Validation without fallback. Tolerating rejections of explicit
values was rejected because identical Roundfix inputs must not silently run
at a different reasoning level. Refines ADR-0037 and ADR-0039.
