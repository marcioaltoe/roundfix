---
type: feat
status: open
created: 2026-11-24
spec: null
reason: null
---

# A rule or ADR must not block the project, and is revised in the Spec phases

## Intent

Add to the documentation Baseline, which applies to every repository, a rule that a rule or ADR serves to keep the implementation within what is specified and must not stop the project's progress. Meaning to preserve:

1. A rule or ADR exists to keep the implementation inside the specification and cannot prevent the project from evolving.
2. When a rule gets in the way there are three exits: deprecate (delete) it, rewrite it, or keep it unchanged and use a one-off exception.
3. In the Idea, PRD, TechSpec and Tasks phases, a rule or ADR that is badly worded or outdated is revised there, instead of becoming an exception request during implementation.

The Jarvis project asked for the normal Roundfix flow (a Spec or a patch release of the Baseline, as Roundfix decides). Changes outside the Baseline or needing other authorization come back as a `PEDIDO DE DECISÃO`. The motivating case was the Vortex project, Spec 0199, Task 02.

## Maintainer decision

The maintainer, on 2026-11-20, 13:45 BRT (message `01a1f9d3-c947-7062-8976-bc2f7050fda6`; resent as part B1 of `01a1fa5d-8462-7192-834d-aa199d2a1785`):

> "Sobre as regras: se estiver atrapalhando o trabalho, podemos revisar e deprecar (excluir), reescrever ou manter inalterada e usar a exceção, mas uma regra não deve impedir o progresso/evolução do projeto; ela serve para manter a implementação dentro do que é especificado. Se estivermos na fase de idea, PRD, techspec e tasks e uma ADR/regra estiver mal redigida ou deprecada, deve ser revista. (Isso vale para qualquer projeto.)"

Where to record it: "Pedir ao Roundfix".

## Follow-up owed

Tell the Jarvis project (`default/jarvis`) the Spec number or the Baseline version in which the rule lands.

## Source

Cross-project messages from the Jarvis project, 2026-11-20: `01a1f9d3-c947-7062-8976-bc2f7050fda6` and `01a1fa5d-8462-7192-834d-aa199d2a1785` (B1). Nature: measured as to the decision (the maintainer's words); the design is not decided. Recorded by triage on 2026-11-24.
