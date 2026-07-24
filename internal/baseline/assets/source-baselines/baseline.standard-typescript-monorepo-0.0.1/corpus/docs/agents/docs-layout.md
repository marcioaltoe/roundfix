<!-- source-baseline-entry: contract.docs-layout.directory-matrix -->
# Documentation directory matrix

| Directory | Purpose | Lifecycle |
| --- | --- | --- |
| `docs/_inbox/` | Raw, untriaged notes | Promote to the correct durable location, then remove the inbox item. |
| `docs/adr/` | Accepted architecture decisions | Append decisions; never reuse an identity. |
| `docs/agents/` | Agent-facing operating guidance | Regenerate setup-owned guides; preserve repository-owned guides. |
| `docs/design/` | Design and interaction artifacts | Keep while active; archive or prune superseded exploration. |
| `docs/findings/` | Dated investigations and field evidence | Preserve history and append dated evidence or routing updates. |
| `docs/handoffs/` | Resumable session state | Replace with newer handoffs and keep only useful recent history. |
| `docs/references/` | External references with local relevance | Prune references that no longer matter. |
| `docs/specs/` | Local Specs, Task Graphs, Tasks, and QA evidence | Archive only after Tasks complete and QA passes. |
| `docs/user-guide/` | Human-facing documentation and runbooks | Update with the behavior it documents. |
<!-- /source-baseline-entry: contract.docs-layout.directory-matrix -->

<!-- source-baseline-entry: contract.findings.template -->
# Findings template

```markdown
---
status: pending # pending | partial | deferred | done
created_at: YYYY-MM-DD
updated_at: YYYY-MM-DD
---

# <Area> — <short title> (YYYY-MM-DD)

<Two to four sentences describing the session or investigation, the attempted outcome, and links to adjacent evidence.>

## 1. <Finding title — symptom, not hypothesis>

- Symptom / evidence: <observed behavior, command output, identifiers, and paths needed to reproduce>
- Root cause: <proven cause, or `unknown` with what was ruled out>
- Action / suggestion: <fix, mitigation, or route to a Spec, direct change, or upstream report>

## 2. <Next finding>

## What worked — keep

<Optional evidence about behavior worth preserving.>
```
<!-- /source-baseline-entry: contract.findings.template -->

<!-- source-baseline-entry: contract.findings.lifecycle -->
# Findings lifecycle

1. Create the dated findings file with `status: pending`, `created_at`, and `updated_at`.
2. Use `pending` while no implementation Spec exists.
3. Use `partial` when a linked Spec intentionally covers only selected observations; record why the remainder is unnecessary.
4. Use `deferred` only when the finding will not be implemented; record the reason.
5. Set `status: done` as soon as an implementation Spec is created and linked. Do not wait for Tasks, QA, archive, merge, or release.
6. Preserve the original observation. Append evidence, root-cause proof, and routing links as dated addenda.
7. Update `updated_at` for every status change or evidence addendum; never change `created_at`.
<!-- /source-baseline-entry: contract.findings.lifecycle -->
