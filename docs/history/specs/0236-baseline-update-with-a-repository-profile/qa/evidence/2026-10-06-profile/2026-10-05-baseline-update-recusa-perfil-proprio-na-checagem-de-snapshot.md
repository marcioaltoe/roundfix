---
origin: oraculum
destination: roundfix
type-hint: fix
created_at: 2026-10-05
capture: manual
---

# `roundfix baseline update` falha com perfil próprio do repositório na checagem de snapshot das skills

## O que acontece

No Oraculum, em 2026-10-05, com Roundfix 0.43.0 (`69462a0c`), o repositório usa
o perfil próprio `.roundfix/baseline/profiles/oraculum-backend.json`. O manifesto
`docs/agents/setup-context.json` registra `"profile": "oraculum-backend"`, e o
perfil passa em `roundfix baseline profile validate`. Mesmo assim, o
`roundfix baseline update` simples, inclusive só como prévia, sai com 1:

```text
roundfix: baseline update failed: load profile for snapshot comparison: Unknown built-in Baseline Profile "oraculum-backend".
```

A prévia chega a listar o plano: 3 arquivos gerenciados, 4 skills desatualizadas
e o Plan Digest. Depois disso, o comando aborta no estágio que compara as skills
instaladas com o Setup Snapshot.

## O que funciona

- `roundfix baseline update --no-skills --confirm-plan <digest>` aplica as guias,
  e uma segunda execução com `--no-skills` reporta `state: current`.
- As quatro skills embutidas no binário (`implement-task`, `qa-gate`, `roundfix`
  e `write-tasks`) não passam por `baseline skills restore --profile
  standard-typescript-monorepo`, que recusa com "not external members of this
  profile". Elas só atualizam por `roundfix skills install --target project`.
- O `roundfix doctor` mostra `skills: ok`, mas com a nota "snapshot comparison
  unavailable: Unknown built-in Baseline Profile "oraculum-backend"".

## Por que importa

O roteiro de cada release do Roundfix (notas 0.24 a 0.43 no inbox do Oraculum)
manda rodar `roundfix baseline update` até reportar `current`. Num repositório
com perfil próprio, esse passo nunca termina limpo. O operador precisa descobrir
sozinho a combinação `--no-skills` mais `skills install --target project`, e o
`release plan` continua acusando as skills até alguém fazer isso.

## O que esperaríamos

- A checagem de snapshot resolve o perfil como o resto do `update` já resolve:
  embutido pelo ID, ou o arquivo em `.roundfix/baseline/profiles/<id>.json`.
- Ou, enquanto isso não existir, o `update` pula a comparação de snapshot para
  perfil próprio com um aviso, em vez de abortar, e o próprio texto de erro
  aponta `--no-skills` e `roundfix skills install --target project`.

## Evidência

- Branch `chore/roundfix-0.43` do Oraculum, commits `cffbc2db` (guias) e
  `89b1d7c5` (skills).
- Nota de memória do Oraculum com o passo a passo:
  `roundfix-023-baseline-profile-proprio.md`, adendo de 05/10.
