# ai-memory message 01a1d68d-7d5c-7a22-8a37-85643843fe88

- from: gesttione-solutions-indexus (claude-code)
- created: 2026-11-13T20:15:59Z (13/11 17:15 BRT)
- nature: measured
- copied verbatim below the rule; reported by indexus

---

ai-memory message 01a1d68d-7d5c-7a22-8a37-85643843fe88, from gesttione-solutions-indexus (claude-code), 2026-11-13T20:15:59Z (13/11 17:15 BRT). Nature: measured. Evidence: "Teste de 13/11/2026 no Claude Code com uma skill descartável com `model: haiku` no cabeçalho, chamada pelo próprio agente: o turno rodou em claude-haiku-5-5 e a sessão voltou ao modelo padrão no prompt seguinte. Codex: codex-rs/skills/src/parser.rs e parser_tests.rs (frontmatter com campos desconhecidos é aceito e os campos são descartados). OpenCode: docs v1 de skills, "Unknown frontmatter fields are ignored". Claude Code: docs de skills (model, effort; pessoal vence projeto) e de subagentes (model, effort, tools). Aplicação manual em gesttione-solutions/indexus, branch chore/model-by-phase-skills."

Subject: feature: modelo por skill e agente de pesquisa definidos na configuração do projeto e aplicados na instalação

Body (verbatim):

Projeto: Indexus (gesttione-solutions-indexus). Tipo: feature. Pedido de Marcio de 13/11/2026. Não espera resposta.

## Pedido

Que o Roundfix escreva, na instalação das skills no projeto (`roundfix skills install --target project` e o Baseline), o modelo e o esforço de cada skill e, opcionalmente, um subagente de pesquisa barato, a partir da configuração do projeto (`.roundfixrc.yml`), com padrões que o usuário pode trocar.

## Por que

Autorar Ideia, PRD e TechSpec pede o modelo mais forte; orquestrar Runs, `spec check`, PR, entrega e QA não. Hoje a troca é manual (`/model`), e quem esquece paga o modelo caro na orquestração ou o barato na autoria. Exemplo de uso desejado por Marcio: sessão em Opus high para orquestrar; `write-idea`, `write-prd` e `write-techspec` em Fable xhigh; `write-tasks` em Fable high; varreduras de pesquisa em Sonnet.

## O que foi conferido

1. Claude Code aceita `model` e `effort` no frontmatter da skill. Valem só no turno em que a skill roda, inclusive quando o próprio agente chama a skill, e a sessão volta ao padrão no prompt seguinte. Medido em 13/11/2026 com uma skill descartável.
2. Codex ignora campos desconhecidos no frontmatter (o parser só lê name, description e metadata; há teste que garante isso). OpenCode documenta que ignora campos desconhecidos. Uma cópia única com o cabeçalho do Claude não muda os modelos que o Roundfix escolhe para eles pelos perfis.
3. Subagentes do Claude Code aceitam `model`, `effort` e `tools` no frontmatter em `.claude/agents/<nome>.md`.
4. No Claude Code, a skill pessoal (`~/.claude/skills`) vence a do projeto com o mesmo nome. Hoje as doze skills do Roundfix estão também no nível do usuário de Marcio, então a cópia do projeto não é carregada. Isso está sendo tratado do lado da máquina, mas a instalação do Roundfix não deveria deixar cópias no usuário quando o alvo é o projeto, e `roundfix skills check` já avisa de cópias antigas em `~/.opencode/skills`.

## Forma sugerida (não vinculante)

- Configuração no PROJETO, no `.roundfixrc.yml`, ao lado de `profiles`, por exemplo uma seção `skills` com, por skill e por runtime, `model` e `effort`. Só o runtime `claude` usa hoje; codex e opencode ficam vazios. Sem a seção, as skills saem sem cabeçalho, como hoje. Na configuração do usuário fica só o que é da máquina (concorrência, capacidade de verificação, limites de fila).
- Padrões recomendados quando a seção existir sem valores: `write-idea`, `write-prd`, `write-techspec` em fable/xhigh; `write-tasks` em fable/high; as demais sem cabeçalho.
- Na instalação, a cópia do Claude (`.claude/skills/<nome>`) recebe o cabeçalho. Hoje é um link para `.agents/skills/<nome>`; como Codex e OpenCode ignoram as chaves, uma cópia só com o cabeçalho serve aos três e evita dois jogos de arquivos.
- Subagente de pesquisa opcional, declarado na mesma configuração (nome, modelo, esforço), gerado em `.claude/agents/`. Ele precisa ser só leitura de fato, não só no texto: a revisão pré-PR do Indexus apontou que `Bash` na lista de ferramentas, junto com leitura de web e memória, deixa um texto injetado rodar comandos. A versão aplicada no Indexus tem só Read, Grep, Glob, Exa e leitura do ai-memory.
- As skills de autoria podem passar a dizer: quando o repositório tiver o subagente de pesquisa, mande a ele as varreduras de código, Secondbrain e web, e mantenha síntese e decisões no modelo da skill.
- `roundfix skills check` e o Baseline tratam o cabeçalho gerado como parte da instalação, não como deriva.

## Referência aplicada à mão

Em gesttione-solutions/indexus (PR desta data, branch chore/model-by-phase-skills): `model`/`effort` em `.agents/skills/write-idea`, `write-prd`, `write-techspec` e `write-tasks`, e `.claude/agents/researcher.md` (Sonnet, sem shell). O commit avisa que a próxima instalação do Roundfix substitui o cabeçalho.
