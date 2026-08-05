---
status: open
created_at: 2026-08-05
updated_at: 2026-08-05
---

# O loop devolve controle por motivos mecânicos, não por decisões (2026-08-05)

Uma sessão longa no Vortex levou três Specs da fila (`0013`, `0014`, `0015`) a duas mergeadas e uma
a um passo do QA, com 32 achados de review resolvidos e três gates de QA executados. O Roundfix fez
o trabalho pesado bem. Este finding é sobre **onde ele parou e devolveu o controle**.

Foram **12 paradas**. **Três** eram decisões que só o mantenedor podia tomar. Das nove restantes,
**cinco vieram do Roundfix ou da fronteira que ele não modela** — e todas devolveram o controle sem
carregar julgamento nenhum.

Complementa
[an accepted gap has no terminal state](2026-08-04-an-accepted-gap-has-no-terminal-state-so-the-loop-cannot-close.md),
que já cobre o status terminal ausente e a modelagem dos portões de merge. Aqui há evidência nova e
comportamentos diferentes.

## 1. O Preflight reporta um defeito por vez, e só no start do Run

### Evidência

A Spec `0014` tinha os Task files escritos em 2026-08-01. Três dias depois, ao iniciar o Run:

```text
Preflight failed
  Task "task_01" ... Task Context entry: kind "pattern": expected "instruction" or "interface" label
```

Corrigido. Novo start:

```text
Preflight failed
  Task "task_01" ... kind "interface": path "packages/backend/drizzle/migrations/": path must be clean
```

Dois starts perdidos, em sequência, por defeitos que **um parser encontra de uma vez**. Os dois
estavam no mesmo arquivo, e eram quatro ocorrências do primeiro tipo espalhadas por três Task
files.

### Comportamento pedido

1. **Enumerar todas as violações num passe.** Preflight de schema é validação estática; parar na
   primeira transforma N defeitos em N tentativas.
2. **Expor a validação como comando próprio** — `roundfix validate --spec <slug>` — para que a
   autoria rode antes de reportar o breakdown, em vez de o defeito dormir até o Run.

Custo hoje: o intervalo entre autorar e executar é a fila inteira. Um Task file com defeito de
schema fica latente por dias.

## 2. Rejeição de hook de commit não é diagnosticada como categoria própria

### Evidência

Aconteceu **quatro vezes**, em três PRs. A saída é um dump cru do `lint-staged` terminando em:

```text
✖ Failed to run tasks for staged files!
⋯ Reverting to original state because of errors…
roundfix: watch failed after Run start: git commit -m fix: resolve Roundfix batch 002 failed: exit status 1
Roundfix completed the Watch Run with a terminal failure
```

O pior caso: o primeiro Run da `0015` completou **cinco Tasks** e morreu na sexta porque um arquivo
passou 500 linhas **por três**. As cinco ficaram numa Worktree retida.

### Por que o diagnóstico importa

Essa saída não distingue duas situações com remediações **opostas**:

- *a correção está errada* → refazer o trabalho;
- *a correção está certa e um hook local rejeitou o commit* → ajustar o que o hook cobra, preservar
  o trabalho.

Um Supervisor humano lê o dump e descobre. Um loop autônomo não tem como.

### Comportamento pedido

Classificar rejeição de hook como **outcome próprio**, nomeando a regra que rejeitou e o caminho de
retomada, com o trabalho preservado explicitamente. O trabalho **já fica** staged — o que falta é
o Roundfix dizer isso em vez de terminar com "terminal failure".

## 3. A ordem Verification → commit → hook põe a falha no lugar mais caro

O contrato diz *"Each Task's Verification commands gate one commit"* e que a Verification roda após
o turno do Agente, com feedback bounded de volta à mesma sessão.

Quando a Verification declarada do repositório é mais **permissiva** que o hook de commit, toda
violação que só o hook pega atravessa a Verification e mata o Run — depois do turno, quando a falha
é terminal em vez de retentativa.

Isso é configuração do repositório, não bug do Roundfix. Mas o Roundfix **pode detectar o
desalinhamento**: se um hook de commit reprova o que a Verification declarada aprova, o loop tem
uma falha estrutural garantida. Um aviso no `doctor` — *"o hook de pre-commit aplica regras que a
Verification declarada não aplica"* — teria evitado quatro paradas.

## 4. O grafo não declara colisão de arquivo, e a concorrência conflita

### Evidência

Na `0015`, `task_06` e `task_07` não têm dependência entre si e rodaram em paralelo sob
`worktree.concurrency`. As duas editaram o mesmo arquivo de teste:

```text
Task task_06 completed.
Task task_07 completed.
Task task_07 failed: integration conflict:
  .../use-cases/__tests__/sync-sales-products.use-case.test.ts
QA Task task_08 withheld; unmet dependencies: task_07
```

Nada se perdeu — a `task_07` estava commitada na Worktree dela e foi recuperada por `cherry-pick`.
O conflito era pequeno: cada lado **acrescentou um teste diferente no mesmo ponto**. Mas parou o
Run e exigiu resolução manual de conflito.

### Comportamento pedido

O manifesto declara `needs` por dependência **lógica**. Duas Tasks sem dependência lógica que
tocam o mesmo arquivo têm dependência **mecânica** que ninguém declara.

Duas saídas, em ordem de preferência:

1. **Derivar do `## Context` dos Task files** — eles já listam os arquivos que cada Task toca — e
   serializar automaticamente Tasks com interseção, sem exigir declaração nova.
2. **Tentar merge de três vias antes de declarar conflito de integração.** Aqui a resolução
   correta era trivial (manter os dois testes) e o `git` resolveria sozinho se as adições não
   fossem adjacentes.

## 5. O `reconcile` é o comportamento que mais funcionou — e vale proteger

Contraponto deliberado, porque a maior parte deste documento é crítica.

Duas vezes nesta sessão houve Worktree retida com trabalho não integrado. Nas duas o `reconcile`
**se recusou a liberar**, com evidência textual:

```text
classification: unintegrated
evidence: Run Branch is not integrated into the target branch
action: preserve
```

E só depois de eu integrar por `merge --ff-only`, reclassificou:

```text
classification: safe
evidence: Run Branch is integrated and Run Worktree is clean
action: would release with --apply
```

**Isso salvou cinco Tasks verificadas.** O handoff da sessão anterior recomendava
`git worktree remove --force`, que as teria descartado. A distinção entre "integrado" e "não
integrado", com recusa explícita e motivo legível, é exatamente o comportamento que um loop
autônomo precisa: recusar por padrão e explicar por quê.

Vale como precedente para os itens 1 a 4: **a recusa informativa é melhor que o dump cru.**

## 6. Nota sobre `--until-clean` com achado permanentemente `failed`

Já registrado no finding de ontem, mas esta sessão deu o número: o PR #116 acumulou **três** issues
`failed` — todas a mesma lacuda aceita de tooling MySQL — e `--until-clean` não podia convergir. O
Supervisor teve que abandonar o comando documentado e dirigir rodadas isoladas à mão, com
`--max-rounds 1`.

O #119, por contraste, fechou em três rodadas (19 → 8 → 6 threads) com **zero** `failed`, e
`--until-clean` teria funcionado.

A diferença entre os dois PRs não é qualidade de código: é a existência de um achado válido que a
Spec não estava autorizada a corrigir.

## Resumo do pedido

| # | Comportamento | Custo medido |
| --- | --- | --- |
| 1 | Preflight completo num passe + comando de validação na autoria | 2 starts perdidos |
| 2 | Rejeição de hook como outcome próprio, com trabalho preservado | 4 paradas, 5 Tasks presas |
| 3 | `doctor` avisa quando o hook é mais estrito que a Verification declarada | causa raiz de (2) |
| 4 | Serializar Tasks com interseção de arquivo, ou tentar merge de três vias | 1 parada |
| 5 | — | (manter como está) |
| 6 | Status terminal de lacuna aceita — ver finding de 2026-08-04 | `--until-clean` inutilizável |

Os itens 1 a 4 são todos da mesma família: **quando o loop encontra um obstáculo mecânico, ele
termina em vez de explicar o suficiente para alguém — ou ele mesmo — seguir**. O `reconcile` já
mostra como fazer diferente.
