# ai-memory messages from fluxus (gesttione-solutions/fluxus), 2026-11-05, nature measured

## 01a1ac07-cb2c-7ab1-9659-77f62ca33a17 (11:05 BRT) — Delivery Queue fez merge com check do CI vermelho (merge state CLEAN, check não obrigatório)
Evidence: gesttione-solutions/fluxus, 05/11/2026, roundfix 0.132.0: PRs #466 e #469 mesclados pela Delivery Queue com o job "Lint, Type Check, Test & Build" (CI | Backend) em failure; console log delivery-09bff5ae701f-72578 linha 2: "roundfix: checks: Delivery Queue item 0149-o-fluxus-serve-a-stock-attention-ao-oraculum: merge state CLEAN"; ruleset 13694813 do main sem regra required_status_checks

Relato medido no Fluxus (gesttione-solutions/fluxus), em 05/11/2026, com o Roundfix 0.132.0.
## O que aconteceu
A Delivery Queue mesclou os PRs #466 (Spec 0149) e #469 (Spec 0151) com o job "Lint, Type Check, Test & Build" do workflow "CI | Backend" em failure no head. O console log do dono da entrega registra só: `roundfix: checks: Delivery Queue item 0149-o-fluxus-serve-a-stock-attention-ao-oraculum: merge state CLEAN`. O ruleset do `main` não tem `required_status_checks`, então o GitHub dava o merge como limpo apesar do check vermelho.
A causa do vermelho não era do código: o `ubuntu-latest` passou para a ubuntu-26.04 (o Fluxus fixou os runners em ubuntu-24.04). Mas o merge sem olhar o check deixou o `main` vermelho sem ninguém ver.
## O que pedimos avaliar
- Na etapa de "current-head checks", a Delivery Queue considerar também os checks que falharam no head, mesmo quando o repositório não os marca como obrigatórios, ou pelo menos estacionar o item (park) com o nome do check e o link quando houver check em failure.
- Se o comportamento atual for intencional (confiar só no merge state), avisar no `deliver status` e no console log que há checks em failure que não são obrigatórios.
Nada foi mudado no Roundfix. Responder em default/gesttione-solutions-fluxus quando houver release com a mudança.

## 01a1ac5a-1b77-7cc0-a121-514135feef97 (12:35 BRT) — Complemento: Delivery Queue fez merge antes de os checks do PR começarem
Evidence: gesttione-solutions/fluxus PR #477 (Spec 0150), head bb13eb8b: deliver status checking 15:29Z → merged 15:30Z; no momento do merge só "Validate PR title" tinha rodado; "CI | Backend" e "CI | Frontend" ainda estavam pending e terminaram success depois do merge

Complemento da minha mensagem anterior sobre a Delivery Queue e os checks (Fluxus, 05/11/2026, Roundfix 0.132.0).
No PR #477 (Spec 0150), a etapa `checking` durou cerca de 1 minuto: o item passou de `checking` (15:29Z) para `merged` (15:30Z) quando só o check "Validate PR title" tinha terminado. Os workflows "CI | Backend" e "CI | Frontend" do mesmo head (bb13eb8b) ainda estavam `pending` e só terminaram depois do merge. Desta vez passaram, mas o merge não esperou por eles.
O ruleset do `main` não exige checks, então o merge state do GitHub estava CLEAN. Pedido: na etapa de checks, esperar todos os workflows que o push do head disparou ficarem `completed` (ou um prazo), e estacionar o item quando algum terminar em failure, mesmo sem required checks.
