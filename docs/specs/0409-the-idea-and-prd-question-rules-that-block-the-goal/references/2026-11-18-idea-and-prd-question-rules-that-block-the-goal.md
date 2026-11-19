---
type: feat
status: promoted
created: 2026-11-18
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
reason: null
---

# The Idea and PRD stages question the rules that block the goal

## Request

The maintainer, 2026-11-18, verbatim: "E eu quero também que esteja no baseline, que no momento da criação do IDEA, principalmente na criação do IDEA e PRD, a gente possa questionar decisões gravadas na memória do projeto, em ADRs e outras, se isso impactar no futuro ou no alcance do objetivo daquilo que está sendo construído. As ADRs, regras do specific repository e outras referências que possam existir dentro do projeto, elas são para guiar a implementação, são guard-rails de guiação da implementação. No momento que a gente está trabalhando no IDEA e no PRD, a gente está trabalhando com criatividade. Se houver uma regra que está impedindo a criatividade, solução de problema ou identificação de problema, causada até por uma regra mal escrita ou mal definida ou deprecada, é importante que a gente tenha a possibilidade, que você traga essa possibilidade de questionamento usando as ferramentas de Ask a question ou compatíveis no Codex ou no Open Code, para que o usuário decida manter ou não com base na análise que vai ser feita. Essa regra hoje impõe isso. E ela impede que eu faça isso, isso e isso para atingir o nosso objetivo. Devo mudar, devo deprecar ou excluir? Alguma coisa assim. Isso tem que existir também no baseline. Tenho visto com frequência esse problema na construção de soluções, porque tem uma regra, então, nem é proposto a mudança da regra, deprecar, regra, simplesmente não faz, porque existe uma regra. Em algum momento, talvez depois de uma, duas, três specs, é falado, não, não foi feito, porque existe essa regra, sim, lá do início do projeto."

## Decisions

- 2026-11-18: when → "Spec após a release"; authorization → "Mesmas do BE (Recommended)" (to be recorded in the Spec's `_authorization.md` when it is written).

- 2026-11-18 (`/roundfix --follow-up`): the single release v0.148.0 is out, so the wait ended; Idea-stage judgment (Jev) indicated (shape_open 0.86, unverified_assumption 0.64, market_question 0.12); the maintainer chose "Ideia agora (Recommended)". Delivery order: third, after the corrective qa-gate Spec and the Delivery Queue fix.

## Where

The Baseline guides and the authoring skills (write-idea, write-prd, and brainstorming), across Claude Code, Codex and OpenCode question tools.
