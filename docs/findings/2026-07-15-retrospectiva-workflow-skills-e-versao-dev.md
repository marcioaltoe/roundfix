# Roundfix — retrospectiva do ciclo 0027/0028 e melhorias de workflow, skills e versão dev

Consolidação dos achados da sessão 2026-07-14/15 (specs 0027 e 0028 implementadas, QA retroativo
da 0023, seis specs arquivadas) mais um achado do operador sobre identificação de builds dev.
Itens já detalhados em outros findings desta pasta são referenciados, não duplicados.

## 1. Build dev sem identificação — `--version` não diz qual binário está em uso

Achado do operador (2026-07-15):

```text
❯ make install
rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix
rtk go install -buildvcs=false ./cmd/roundfix
❯ roundfix --version
roundfix 0.0.0-dev
```

Todo build dev imprime o mesmo `0.0.0-dev` — impossível saber se o `make install` surtiu efeito
ou qual versão de desenvolvimento está instalada na máquina. Agrava-se porque `-buildvcs=false`
(necessário neste ambiente) remove o carimbo VCS que o Go faria sozinho.

**Melhoria desejada**: `roundfix --version` de um build dev deve exibir data e hora **local** do
build e a identidade do código, para provar que o build funcionou e qual dev version está em uso.

Sugestão de forma: stampar via `-ldflags` no `Makefile` (`build` e `install`) a versão, o short
SHA do commit e o timestamp do build, exibindo algo como:

```text
roundfix 0.0.0-dev (a1b2c3d, built 2026-07-15 14:32:05 -0300)
```

Mantendo a primeira linha estável para parsing e a identidade extra na mesma linha ou em linha
adicional — decidir no techspec respeitando o contrato de saída (`agentic-cli-design`:
introspecção honesta; a string de versão é API pública).

## 2. Cópia global obsoleta da skill roundfix shadowing a do repositório

`~/.claude/skills/roundfix` existe como cópia global e foi a versão servida ao Supervisor nesta
sessão — ainda ensinando o contrato pré-0027: review Runs em Run Worktree, nota "treating Run as
Clean", "Roundfix never comments on threads". Qualquer agente que a leia segue um contrato que
não existe mais (o repositório, via `.claude/skills → .agents/skills`, está correto e sincronizado).

**Ação**: remover a cópia global (ou trocá-la por symlink para a do repositório) e considerar,
no `skill-governance`, um aviso de que cópias globais da skill roundfix não devem existir.

## 3. Guidance não previne gate não-hermético — falta enforcement

O autor da 0027 escreveu `go build ./...` (sem `-buildvcs=false`) e `rg` em task files mesmo com
a regra documentada em `docs/agents/autonomous-work.md` — custo: 2 tasks failed com trabalho
100% pronto + 2 settles (ver `2026-07-14-implement-0027-launch-findings.md` §4). A 0028
(task_08) corrigiu a *guidance* do write-tasks; falta o *enforcement*:

1. **write-tasks**: passo final mecânico de lint nos task files gerados (padrões proibidos:
   `go build ./...` sem buildvcs, `rg `, pipelines com `wc`).
2. **roundfix (candidato)**: o preflight do `implement` avisar sobre padrões de Verification
   sabidamente não-herméticos antes de qualquer agente gastar um ciclo neles.

## 4. Discrição do Agente sobre emendar comandos de Verification é não-especificada

No mesmo run, agentes das tasks 01–08 **reescreveram silenciosamente** os gates quebrados para
equivalentes seguros no ambiente (rg→grep, buildvcs) e passaram; o agente das tasks 09/10
**recusou** emendar e settlou failed com o trabalho completo. Ambos defensáveis no contrato
atual — o desfecho dependeu do temperamento do agente. O contrato de Batch deve dizer
explicitamente: substituição por comando comprovadamente equivalente é permitida **e registrada
no task file** (o Daemon reexecuta o que estiver no arquivo de qualquer forma), ou é proibida.
Recomendação: permitida-e-registrada.

## 5. Débito de QA aparece só na hora do archive

A 0023 shipou antes do `--qa` existir e a ausência de QA Report só apareceu quando o
`roundfix archive` recusou — o QA retroativo funcionou (relatório
`_archived/0023-.../qa/qa-report-2026-07-15.md`, verdict pass) mas custou uma sessão.
Regra futura já correta (todo implement com `--qa`); adicionar uma linha no release runbook:
"nenhuma spec mergeia sem QA Report".

## 6. Candidatos de produto já registrados (referências)

- **Divergência probe × prompt** na disponibilidade de modelo + lista anunciada é dinâmica
  (mudou entre 14 e 15/07 e "desquebrou" o pin `gpt-5.6-sol`); `doctor` poderia imprimir os
  modelos anunciados. Ver `2026-07-14-implement-0027-launch-findings.md` §1 e o QA report da 0023.
- **`--detach` morre silencioso pré-handshake** com stderr vazio. Ver o mesmo findings, §2.
- **Settle resolve superfície kept obsoleta** em vez do checkout onde a task está `failed`.
  Ver o mesmo findings, §3.
- **Settle varre o worktree inteiro** — 0028 shipou transparência (paths + aviso); se
  insuficiente na prática, a escalada é settle com pathspec.

## 7. O que funcionou — manter

- Verificação do Daemon prevaleceu sobre o relato do agente em todos os casos.
- Commits por task tornaram a recuperação cirúrgica: um run morto no meio perdeu zero trabalho
  settlado (ff-integration + settle recuperaram tudo).
- O loop findings → spec → implement está apertado: três features da 0028 corrigiram modos de
  falha reproduzidos ao vivo na mesma sessão (lock órfão, sweep do settle, diagnóstico de
  adapter).
- Supervisão por `events --follow` filtrado + lançamento nohup+disown foi confiável enquanto o
  `--detach` não é corrigido.
