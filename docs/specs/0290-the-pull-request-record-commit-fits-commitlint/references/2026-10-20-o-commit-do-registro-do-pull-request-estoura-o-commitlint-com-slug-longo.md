---
origin: oraculum
destination: roundfix
type-hint: fix
created_at: 2026-10-20
capture: manual
---

# O commit do registro do pull request estoura o commitlint quando o slug da Spec é longo

## Observação

Em 2026-10-20 às 16h02 BRT, a Delivery Queue do oraculum estacionou a Spec
`0088-as-lojas-autonomas-declaram-os-pontos-de-operacao` com `delivery-error`
no passo Pull Request Record Commit:

```text
delivery-error: record pull request: commit pull request record: git commit -m docs: record pull request #280 in the 0088-as-lojas-autonomas-declaram-os-pontos-de-operacao Archive Record: … commitlint … exit status 1
```

O assunto tem 107 caracteres. O hook `commit-msg` do repositório roda o
commitlint com `@commitlint/config-conventional`, cujo `header-max-length` é
100. Um `deliver retry` simples repete o mesmo assunto e falha igual.

O assunto vem fixo de `internal/cli/deliver_workflow.go`, perto da linha 1352:
`"docs: record pull request #"+number+" in the "+specSlug+" Archive Record"`.
Com o prefixo de 4 dígitos, todo slug de mais de 48 caracteres estoura o limite
quando o número do pull request tem 3 dígitos. Os outros assuntos gerados pela
fila ficam abaixo do limite com o mesmo slug: `docs: archive <slug>` tem 68
caracteres e `feat: deliver <slug>` tem 67.

## Contorno aplicado

O passo adota um commit já existente quando ele é exatamente a mudança do
registro sobre a cabeça anterior (`pullRequestRecordCommitIsExact`). No item
worktree, o Supervisor commitou a mudança já staged com o assunto curto
`docs: record pull request #280 in the 0088 Archive Record` e rodou
`roundfix deliver retry`, e a fila seguiu. A Spec
`0089-o-login-gesttione-conecta-o-claude-web-e-o-chatgpt` vai precisar do mesmo
contorno.

## Por que é acionável

- O assunto pode usar só o número da Spec (`0088`) ou caber em 100 caracteres
  cortando o slug.
- Se o commit falhar no hook, a classificação poderia apontar o commitlint e o
  contorno, em vez de um `delivery-error` genérico.
