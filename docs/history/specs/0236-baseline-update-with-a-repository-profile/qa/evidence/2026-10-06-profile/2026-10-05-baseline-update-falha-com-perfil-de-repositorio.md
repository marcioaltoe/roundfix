---
origin: pantheon
destination: roundfix
type-hint: finding
created_at: 2026-10-05
capture: manual
---

# `baseline update` falha em todo repositório com perfil próprio: a comparação de snapshot só aceita perfil built-in

## Observação

No Pantheon, adotado em 2026-10-02 com o perfil de repositório
`.roundfix/baseline/profiles/pantheon-devops.json` (Roundfix 0.26), o
`roundfix baseline update` deixou de funcionar no 0.43.0
(`69462a0c`, build 2026-10-05):

```
$ roundfix baseline update --repo . --format json
exit 1
roundfix: baseline update failed: load profile for snapshot comparison: Unknown built-in Baseline Profile "pantheon-devops".
category: execution
nextAction: repair snapshot comparison and rerun roundfix baseline update
```

Antes de falhar, o update já tinha calculado um plano válido: duas mudanças
(`docs/agents/setup-context.json` e `docs/agents/skill-dispatch.md`), digest
`sha256:1b6b08fa03805907e4aa470d6aba7905890d3a660dd86d0e2624ae850cb48745`, e
`skills.outdated` com `implement-task` 0.0.2→0.0.3, `qa-gate` 0.0.6→0.0.8,
`roundfix` 0.1.10→0.1.33 e `write-tasks` 0.0.7→0.0.8. No mesmo minuto,
`roundfix baseline profile validate pantheon-devops` responde
`valid (repository)`.

## Onde

`internal/baseline/skills_trailing.go`, `TrailingSetupSkills`, entregue em
`dd0e5de9` (Spec 0215). Ela carrega o perfil com `loadRestoreProfile(catalog,
profileID)`, que só resolve perfis built-in do catálogo embutido. O chamador é
`internal/cli/baseline_update.go:282`; o `doctor` usa a mesma comparação
(`internal/cli/doctor.go:135`) e só anexa o erro ao detalhe.

## Por que é acionável

Todo repositório adotado com perfil de repositório (o caminho que a
documentação de "Profile alignment and adaptation" recomenda quando nenhum
built-in serve) perde o update gerenciado, e com ele a atualização das skills
do Roundfix. A comparação de snapshot existe para skills externas com
`TreeDigest`; um perfil de repositório pode resolver as mesmas skills pelo
perfil built-in de origem, ou a comparação pode ser pulada com aviso quando o
perfil não é built-in, em vez de abortar o update inteiro.

Evidência capturada localmente no Pantheon; o repositório não foi alterado.
