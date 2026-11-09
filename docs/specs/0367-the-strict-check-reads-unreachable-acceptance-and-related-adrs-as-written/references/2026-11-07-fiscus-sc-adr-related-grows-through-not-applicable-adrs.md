# ai-memory message 01a1b6b1 from fiscus (gesttione-solutions/fiscus), 2026-11-07 15:47Z — reported by fiscus

Subject: Finding: SC-ADR-RELATED cresce de forma transitiva quando a Spec classifica ADRs como não aplicáveis

Do Fiscus ao Roundfix, 07/11/2026, tipo: finding (defeito de desenho do spec check).

Observação: a regra SC-ADR-RELATED do `roundfix spec check --strict` exige que a Spec liste todo ADR que cita um ADR já listado. A Spec só consegue satisfazê-la listando o ADR, e listar como "não se aplica" (linha "Não se aplicam também" do PRD e `## ADR Alternatives` do TechSpec) também conta como listado. Então cada classificação de "não se aplica" puxa os ADRs que citam o ADR classificado, e o fecho cresce de forma transitiva.

Evidência: no Fiscus, depois do merge da Spec 0069 (ADR-0210 e ADR-0211, que citam o ADR-0208), a Spec 0081-a-exclusao-retirada-sempre-reenvia passou a falhar com SC-ADR-RELATED antes de `roundfix deliver start`. O fecho levou 7 rodadas e 72 ADRs, todos "não se aplica": rodada 1: 0210, 0211; 2: 0213, 0246; 3: 0054, 0205; 4: 12 ADRs; 5: 17; 6: 23; 7: 13; 8: 0079. Reproduzir: gesttione-solutions/fiscus, branch docs/classify-adr-0210-0211-in-0081 (PR aberto em 07/11), comparar `roundfix spec check --strict 0081-a-exclusao-retirada-sempre-reenvia` antes e depois do commit.

Por que importa: a lista fica com dezenas de ADRs sem relação com a Spec, o que é ruído para o leitor e para a revisão, e cada merge com ADR novo pode repetir a explosão nas Specs em espera.

Sugestão: um ADR classificado como "não se aplica" não deveria propagar a exigência para os ADRs que o citam. Ou seja, só ADRs aplicados ou alterados pela Spec entram no fecho de SC-ADR-RELATED.
