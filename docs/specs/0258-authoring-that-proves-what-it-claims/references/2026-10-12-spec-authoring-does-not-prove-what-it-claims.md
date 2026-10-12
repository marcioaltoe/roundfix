---
type: fix
status: promoted
created: 2026-10-12
spec: 0258-authoring-that-proves-what-it-claims
---

# Spec authoring does not prove what it claims

## Symptom

Four gaps in `write-techspec` and `write-tasks`. All of them were measured in
the fluxus, and the `skills` maintainer confirmed on 2026-10-09 that none is
covered by the synced skills.

1. **A slice's Verification does not cover its consumers.** In fluxus Spec
   0021 two green Tasks left the repository red. One moved a read endpoint,
   the other moved two trigger routes, and each broke a suite in another
   directory. Only the terminal QA gate's full run found it. (Part 1 of the
   original entry, "the exit status is the verdict", is already covered.)
2. **Nothing proves that what the Spec builds does not already exist.** Fluxus
   Spec 0028 ordered a deterministic quarantine that already existed, with
   eight tests, in the file its own TechSpec cited. It passed `write-techspec`,
   `write-tasks` and `spec check`. The agent reconciled the contradiction by
   lowering the threshold from 3 to 1 and deleting three tests, and the gate
   stayed green.
3. **ADR applicability is declared, but the relation to rejected alternatives
   is not.** Fluxus Spec 0020 raised a timeout from 60s to 120s, which is
   exactly the alternative an accepted ADR had rejected. It cost four gate
   cycles.
4. **The Coverage Map is filled in but never checked** against the PRD. A PRD
   promised per-request records, and its TechSpec aggregated them on purpose.

## Where

The `write-techspec` and `write-tasks` skills and `roundfix spec check --strict`.

## Expected

Maintainer decision, 2026-10-11 ("Skills + gap no spec check"):

- `write-tasks`: a Task that changes a shared contract (route, public
  signature, persisted format) names the consuming suites in its
  Verification, or the graph gains an integration checkpoint before the QA
  gate.
- `write-techspec`: each component declared new records the search that
  proved it absent, as a command and its result. Each applicable ADR gets one
  line on how the Spec relates to the alternatives that ADR rejected.
- The PRD-to-TechSpec Coverage Map cross-read happens before decomposition.
- `spec check --strict` reports the missing section or line structurally,
  without judging the content.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-08-08-write-tasks-nao-governa-exit-status-nem-consumidores.md`
- Secondbrain `inbox/roundfix/_triaged/2026-08-13-autoria-de-spec-nao-verifica-se-o-que-ela-manda-construir-ja-existe.md`
- Secondbrain `inbox/roundfix/_triaged/2026-08-18-a-autoria-de-spec-nao-presta-contas-de-findings-nem-verifica-o-coverage-map.md`
- Secondbrain `inbox/roundfix/_triaged/2026-10-09-skills-council-the-fool-removidos-e-quatro-entradas-reroteadas.md`
  (coverage status on 2026-10-09)
