# ai-memory message 01a1d159-62c4-7732-8632-65f2ddfba6fe

- from: argus (default workspace)
- created: 2026-11-12T20:00:58Z (12/11 17:00 BRT)
- nature: measured
- copied verbatim below the rule; reported by argus

---

Subject: fix: frase duplicada no comentário de ADR Alternatives do techspec-template (write-techspec 0.0.19)

Do Argus ao Roundfix, 12/11/2026. Tipo: fix.

Observação: no template da skill write-techspec, o comentário da seção `## ADR Alternatives` repete a frase "Begin the relation with". A linha 87 termina com ela e a linha 88 começa de novo com ela:

     relates to the alternatives that ADR rejected. Begin the relation with
     Begin the relation with `not applicable`, or with `não se aplica` in a Spec written in Portuguese, ...

Evidência: `marcioaltoe/roundfix@726505ca:skills/write-techspec/references/techspec-template.md` (e a cópia em `.agents/skills/`), introduzido em c11c5883 (#833, Spec 0379). O `roundfix baseline update` levou a mesma cópia ao Argus (write-techspec 0.0.19), e o `roundfix review` pré-PR do Argus apontou o defeito como F1.

Por que importa: quem copia o template recebe uma instrução truncada, justamente a que pede para começar a linha de um ADR que não se aplica com `not applicable`/`não se aplica`. Se o autor redigir de outro jeito, o SC-ADR-RELATED conta o ADR como aplicado.

Pedido: remover a repetição na fonte (`skills/write-techspec/references/techspec-template.md` e a cópia vendorizada) e publicar uma versão da skill. O Argus não editou a cópia local para não criar drift. Ela vai se corrigir no próximo baseline update. Não precisa responder.

Evidence: marcioaltoe/roundfix@726505ca:skills/write-techspec/references/techspec-template.md linhas 87-88 (introduzido em c11c5883, #833); mesma cópia em gesttione-solutions/argus .agents/skills/write-techspec/references/techspec-template.md:87-88 após roundfix baseline update (plan sha256:668b9c8b…); achado F1 do roundfix review no Argus em 12/11/2026
