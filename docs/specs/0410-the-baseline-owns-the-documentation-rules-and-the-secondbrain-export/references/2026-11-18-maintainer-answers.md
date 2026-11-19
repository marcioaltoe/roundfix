Spec das regras de documentação e da exportação ao Secondbrain no Baseline (backlogs docs/backlog/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md e docs/backlog/2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md, PR #914).
Pedido original: está na seção `## Request` da primeira Backlog Entry e na seção `## Tangled` da segunda, as duas adotadas por esta Spec.
Julgamento da Ideia (18/11, roundfix spec judge --stage idea, typesafe/jev-1.13-20260917), item mensagem 01a1eeca: undecided — shape_open unsure (0.68), unverified_assumption no (0.17), market_question no (0.06).

Respostas do Marcio — /roundfix --follow-up, 18/11/2026 ~19:55 BRT (AskUserQuestion), as que tratam desta Spec, verbatim:
2. Pedido do Indexus (01a1eeca + 01a1eec1: regras de documentação/pastas/exportação ao Secondbrain no Baseline via update/doctor + Lições Aprendidas): "Spec com etapa de Ideia (Recommended)".
3. Backlog 2026-11-16 pastas de docs dos adotantes: "Juntar à Spec do Indexus (Recommended)".
4. Autorização da Spec do Indexus: "Mesmas do BE (Recommended)" (push, PR, merge, release após a Spec; até 5 correções; Makefile/CI/go.mod/lint só se a correção exigir).
6. Ordem: "Corretiva → regras → Indexus (Recommended)"; o fix da fila entra logo após a corretiva.

Retomada 18/11/2026 ~22:35 BRT (AskUserQuestion), rótulos verbatim:
- Condução da Ideia 0410: "Ideia 0410 agora (Recommended)" — a pesquisa da Ideia 0410 (Fable) começa já, o PRD 0409 espera ela terminar (nunca duas sessões Fable em paralelo); perguntas e aprovação da Ideia com o Marcio presente.

Ideia 0410 — respostas do Marcio (AskUserQuestion, 18/11/2026 ~23:20 BRT), rótulos verbatim:
1. Pastas (pedido de 16/11): "Só a função de cada pasta (Recommended)" — o Baseline diz o que cada pasta guarda (inclusive docs/design/ e onde vão arquivos como DATABASE.md); o doctor aponta o que está fora; nada é renomeado nem movido; renomear/juntar volta ao backlog com o custo medido.
2. Variação da lista de exportação: "Padrão + variação declarada (Recommended)" — o Baseline escreve o bloco padrão e o baseline update o mantém; o repositório pode tirar itens ou acrescentar arquivos de conhecimento num bloco próprio com o motivo em uma linha; o doctor aponta pasta de processo exportada e bloco padrão alterado.
3. ADRs na exportação: "Só as marcadas" (NÃO a recomendada) — cada ADR ganha um campo que diz se interessa à base de conhecimento.
3a. Quem marca: "Agente marca, você vê no PR (Recommended)" — ao escrever uma ADR o autor a marca (conhecimento/regra de negócio ou ferramenta); nas existentes, o baseline update gera a marcação num PR por repositório, que o Marcio pode corrigir; ADR sem marca não é exportada.
4. ADR que deixou de valer: "Apagar substituída, emendada no topo (Recommended)" — a regra de apagar fica; doctor aponta a superseded que sobrou e o baseline update a remove; a ADR emendada ganha no cabeçalho a lista das ADRs que a emendam.
5. Momento da Lição Aprendida: "Finding fechado e archive c/ tropeço (Recommended)" — Finding que fecha e Spec arquivada que teve correção, QA reprovado ou retentativa; sem lição reutilizável, nada é gravado.
6. Quem decide a Lição: "Agente grava, você vê no PR (Recommended)" — o Agente aplica o teste do Baseline (tentativa, problema, causa provada, solução, ainda não dito em ADR/referência/glossário) e grava no mesmo PR; quem revisa pode apagar.

Aprovação da Ideia 0410 (AskUserQuestion, 18–19/11/2026 ~23:55 BRT), rótulos verbatim:
- Ideia: "Aprovar e seguir ao PRD (Recommended)" — uma Spec, 4–5 Tasks + QA (a 5ª: o Roundfix aplica as regras a si), lista padrão README, CONTEXT.md, docs/user-guide/, docs/references/, docs/design/, docs/lessons/ e só ADRs marcadas; exceção à regra "não reescrever ADR só para metadados".
- Quem propõe a marca das ADRs existentes: "O Agente da sessão do update (Recommended)".
- Lições do Roundfix: "Sim, mais docs/lessons/ (Recommended)" — Roundfix exporta README.md, docs/user-guide/ e docs/lessons/.
- docs/design/: "Sim, na lista padrão (Recommended)".

Perguntas em aberto do PRD e do TechSpec, com o Marcio ausente (registrado no brief do operador de 19/11/2026 ~02:25 BRT), rótulo verbatim:
- "Seguir o padrão do autor (Recommended)" — para uma pergunta que não é fronteira do mantenedor (dado enviado para fora, escopo novo, gasto), o autor escolhe a opção que recomenda e a registra no PRD como decisão do operador delegada pelo mantenedor, com a pergunta e as opções; uma pergunta de fronteira para e vai ao relatório.
