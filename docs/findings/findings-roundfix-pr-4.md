# Findings — Roundfix Watch do PR #4

Data da análise: 2026-07-14  
Run: `run_20260714T010341Z_c74f28986b116a5c`  
PR: `#4` (`gesttione-solutions/oraculum`)  
Branch: `ma/0002-postgres-0003-auth-axis-vortex`  
Comando: `roundfix watch --source coderabbit --pr 4 --agent codex --until-clean --max-rounds 3 --no-input`

## Resumo

O Roundfix executou o `fetch` e iniciou o `resolve`, mas não chegou ao commit, push ou novo review. O daemon bloqueou o assentamento porque o `make verify` falhou duas vezes no mesmo erro de TypeScript. Sem um batch verificado, o contrato do Roundfix impede commit e push.

O agente registrou os 17 issues como `Decision: resolved`, mas o código do worktree ainda continha `toSorted`, exatamente a API que fazia o typecheck falhar. O daemon prevaleceu sobre o relato do agente e finalizou os issues como `failed`.

## Fluxo esperado e fluxo observado

| Etapa | Esperado | Observado |
| --- | --- | --- |
| Fetch | Buscar o review atual | Round 001 buscou 17 issues às `2026-07-14T01:03:44Z` |
| Resolve | Corrigir, testar e marcar os issues | O agente editou o worktree e gravou `Decision: resolved` nos arquivos de issue |
| Verify | Passar `make verify` antes de assentar | Falhou nas tentativas 1 e 2 com `TS2550` em `toSorted` |
| Commit | Criar commit do batch verificado | Não ocorreu; não existe commit novo no branch |
| Push | Enviar o commit ao upstream | Não ocorreu; o journal registra `Final Push blocked` |
| Novo review | Aguardar o CodeRabbit analisar o novo HEAD | Não ocorreu porque não houve push |

`--max-rounds 3` define o limite máximo. Não obriga três rounds quando o batch atual termina falho; o Run terminou após o primeiro round com `Unresolved`.

## Evidências

### Verificação do daemon

Os dois logs de verificação apontam para o mesmo erro:

```text
src/infra/repositories/postgres/registro-organizacoes-postgres.repository.ts(102,34):
error TS2550: Property 'toSorted' does not exist on type 'Organizacao[]'.

src/infra/repositories/postgres/repositorio-autenticacao-mcp-postgres.repository.ts(69,42):
error TS2550: Property 'toSorted' does not exist on type 'Organizacao[]'.

make: *** [verify] Error 1
```

Arquivos de evidência:

- `/Users/marcio/.roundfix/artifacts/d1255e752e9d9f49/runs/run_20260714T010341Z_c74f28986b116a5c/verification/batch-001-attempt-1.log`
- `/Users/marcio/.roundfix/artifacts/d1255e752e9d9f49/runs/run_20260714T010341Z_c74f28986b116a5c/verification/batch-001-attempt-2.log`

O journal registra:

- tentativa 2 de `make verify`: `failed`, exit status 2;
- `Batch 001 failed; 17 Unresolved Review Issue(s) remain`;
- `Final Push blocked: 17 Unresolved Review Issue(s) remain`;
- estado final: `Unresolved`.

### Divergência entre o relato do agente e o worktree

Os arquivos `issue_015.md` e `issue_016.md` dizem que `toSorted` foi substituído por `sort` e que o typecheck passou. Porém, a inspeção do worktree mantido mostrou que o código ainda contém:

```ts
return [...organizacoesAtivas].toSorted(...)
```

nos dois arquivos apontados pelo typecheck. Portanto, a evidência textual nos Review Issues não correspondia ao estado verificável do código quando o daemon executou o gate.

### Estado Git

- Checkout principal: permanece no commit `8ce948f` e alinhado com `origin/ma/0002-postgres-0003-auth-axis-vortex`.
- Worktree do Run: mantido com alterações não commitadas em 19 arquivos.
- Não há commit de correção nem push.
- Os artefatos de review ficaram não rastreados em `docs/specs/_reviews/pr-4/round-001/`.

## Problemas encontrados

### 1. O agente declarou uma correção que não estava no código

Impacto: Blocks-Completion.  
O agente atualizou a narrativa dos issues e afirmou que corrigiu `toSorted`, mas o worktree continuou com a implementação incompatível. O gate independente detectou a divergência.

Sugestão: antes de atualizar o status de um issue, o batch deve reler os arquivos alterados, confirmar o diff efetivo e executar o mesmo comando de verificação que o daemon executará.

### 2. O agente usou verificações focadas como evidência final, mas não garantiu `make verify`

Impacto: Blocks-Completion.  
Os testes focados, typecheck local e lint foram relatados como aprovados, mas o comando obrigatório do repositório era `make verify`. O daemon executou o gate correto e encontrou a falha.

Sugestão: tratar `make verify` como pré-condição do fechamento do batch. Verificações focadas podem complementar a evidência, mas não podem substituir o comando configurado pelo daemon.

### 3. O status final dos artefatos ficou ambíguo

Impacto: Friction.  
Os arquivos têm `Decision: resolved` no bloco de triagem, mas `status: failed` no frontmatter. O primeiro descreve a decisão do agente; o segundo é o resultado terminal do daemon. A distinção existe, mas não está clara para quem inspeciona o artefato.

Sugestão: registrar explicitamente no arquivo o motivo terminal, incluindo o caminho do log de verificação e o comando que falhou, quando o daemon converter um issue triado como `resolved` em `failed`.

### 4. O Watch não avançou para o segundo round

Impacto: Blocks-Completion.  
O limite de três rounds não é uma garantia de três tentativas. Como o primeiro batch não passou pelo gate, não houve commit para o Watch publicar nem novo HEAD para o CodeRabbit revisar. Encerrar em `Unresolved` nesse ponto é coerente com o contrato atual, mas não atende à expectativa de recuperação automática após uma falha de verificação.

Sugestão: documentar essa semântica no comando e, se a intenção for retry automático, separar retries de verificação de rounds de Review Source. Um retry de verificação não deve consumir um round nem atualizar status como se um novo review tivesse ocorrido.

## Solução imediata sugerida

1. No worktree mantido, substituir as duas chamadas a `toSorted` por cópias ordenadas com `sort`, compatíveis com o target TypeScript do projeto.
2. Executar `make verify` completo no worktree e confirmar que o código efetivo, não apenas os textos dos issues, contém a correção.
3. Reexecutar o fluxo de resolução pelo Roundfix para que o daemon assente o batch, crie o commit e faça o push.
4. Somente depois do push, aguardar o novo status do CodeRabbit e iniciar o round seguinte.

Não marcar os issues manualmente como `resolved` nem fazer push manual antes de uma verificação completa; isso contornaria o gate que evitou publicar um commit que ainda falhava no typecheck.

## Estado final

Verdicto do Run: `Unresolved after 1 Round(s): 0 resolved, 0 invalid, 17 failed, 0 unresolved.`  
Commit/push: não realizados.  
Novo review do CodeRabbit: não iniciado.

## Findings adicionais durante a recuperação

### 5. O autofix do lint reintroduzia a falha depois da correção

Impacto: Blocks-Completion.  
Uma troca manual de `toSorted` por `sort` passava no typecheck antes do lint, mas `bun run lint:fix` aplicava `unicorn(no-array-sort)` e restaurava `toSorted`. O `Makefile` executa `lint:fix` antes de `typecheck`, portanto o verify sempre voltava a falhar.

Evidência: `rtk bun run --cwd packages/backend lint` reportou `unicorn(no-array-sort): Use Array#toSorted()`, enquanto `tsconfig.base.json` define `target: ES2022` e o typecheck rejeita `toSorted`.

Solução aplicada: manter `sort` sobre uma cópia nova e adicionar exceções locais `oxlint-disable-next-line unicorn/no-array-sort`, documentando a compatibilidade com ES2022. Depois disso, `make verify` passou.

### 6. A correção precisa ser validada depois do autofix, não antes

Impacto: Blocks-Completion.  
Uma verificação executada antes de `lint:fix` não representa o estado que o daemon verificará. O pipeline precisa testar a ordem real do repositório: formatter, lint autofix, OnionCry, typecheck e testes.

Sugestão: o agente deve executar `make verify` como último teste do batch e reler o diff depois do comando. O daemon pode reforçar isso registrando no journal os arquivos alterados pelo autofix e comparando-os com a evidência do agente.

## Sugestões de melhorias do Roundfix

- Fazer o agente executar o comando de verificação configurado pelo projeto antes de marcar os issues como `resolved`.
- Após a verificação, reabrir o conteúdo dos arquivos alterados e validar que a evidência registrada corresponde ao diff efetivo.
- Quando a verificação falhar, acrescentar automaticamente ao Review Issue o comando, exit status, caminho do diagnóstico e resumo do erro; hoje `Decision: resolved` e `status: failed` ficam lado a lado sem explicar a transição.
- Diferenciar no modelo de Run `retry de verificação` de `novo round de Review Source`. O retry não deve ser contado como novo review e não deve permitir commit/push.
- Detectar no preflight conflitos entre o target TypeScript e regras de lint autofix, como `ES2022` combinado com `unicorn/no-array-sort`.
- Exibir no Live Run View uma linha explícita para cada transição: `fetch`, `resolve`, `verify`, `commit`, `push`, `review wait` e `round complete`.
- Quando o batch falhar antes do commit, oferecer um caminho de recuperação documentado que reutilize o worktree mantido sem exigir marcação manual dos issues ou push fora do daemon.

### 8. Atualização antecipada do status e comentários no Review Source

O Roundfix deve atualizar o status da Review Issue no GitHub o mais cedo possível após a decisão. Para issues independentes, isso pode ocorrer ao final do fix; para um batch em que a atualização intermediária teria custo significativo, a atualização pode ocorrer ao final do batch, desde que a decisão e a evidência estejam registradas.

Qualquer status diferente de `resolved` deve gerar um comentário na issue do GitHub. O comentário deve explicar o motivo e o próximo passo:

- `invalid` ou issue ignorada: motivo verificável pelo qual o finding não se aplica;
- `failed`: etapa que falhou, comando, exit status, caminho do diagnóstico e ação necessária;
- `duplicated`: referência ao issue original, quando essa transição for autorizada pelo daemon;
- `unresolved`: motivo pelo qual a issue permanece aberta e quando será reavaliada.

O comentário deve ser idempotente, evitar duplicações em retries e aparecer no journal do Run com a referência da Review Issue. A resolução no GitHub não pode ser feita antes de o código e o gate obrigatório confirmarem o resultado; atualizar cedo significa tornar a decisão e o progresso visíveis, não marcar sucesso prematuramente.

### 7. Guardrail de estado da branch antes de `fetch` e `watch`

Proposta: o modo de review do Roundfix não deve criar worktrees. Antes de iniciar o `fetch` — e, consequentemente, antes de iniciar o `watch` — o preflight determinístico deve reconciliar todas as worktrees existentes relacionadas à branch do PR. As alterações dessas worktrees precisam estar mergeadas na branch responsável pelo PR; se o merge não puder ser feito com segurança, o comando deve bloquear e informar a worktree, branch, commits divergentes e a ação necessária.

O mesmo preflight deve bloquear quando existir qualquer Run ativo vinculado à branch, ao PR ou a uma worktree relacionada. O bloqueio precisa ocorrer antes de qualquer fetch, criação de sessão de agente, comentário de review ou alteração no código.

O único bypass deve ser explícito, por exemplo `--force`. Quando usado, o Roundfix deve registrar automaticamente um comentário na pull request contendo:

```text
Roundfix Watch preflight bypassed with --force.
Existing worktrees or active Runs related to <branch> were not reconciled before review.
```

O comentário deve incluir a branch, Run IDs, worktrees afetadas e o motivo do bypass, para que o review continue auditável.

Essa regra precisa existir em três camadas:

1. guardrails determinísticos do preflight de `watch`;
2. guardrails determinísticos do comando `fetch`, que é uma etapa do Watch;
3. skill e resources de `watch`, `fetch`, `resolve` e review de PR, para que agentes conheçam o contrato e não contornem o bloqueio.

O contrato deve ser implementado como bloqueio por padrão, não como recomendação textual. O preflight só pode prosseguir sem `--force` depois de verificar que não há worktrees divergentes nem Runs ativos relacionados à branch.
