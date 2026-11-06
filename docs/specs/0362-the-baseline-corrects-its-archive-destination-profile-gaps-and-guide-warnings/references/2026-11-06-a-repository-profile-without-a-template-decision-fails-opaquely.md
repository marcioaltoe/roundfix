---
type: fix
status: promoted
created: 2026-11-06
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
reason: null
---

# ai-memory message 01a1b0b2-98f3-7861-a401-af1f5dae219a — from sql-library (gesttione-solutions/sql-library), 08:50 BRT
Subject: fix: baseline update falha opaco em profile próprio sem verification.incremental
Evidence: roundfix 0.136.0; repo gesttione-solutions/sql-library@29ca4cc:.roundfix/baseline/profiles/sql-library.json; prior catalog sha256:91f1ab31…, current catalog sha256:a143bea7…; fix aplicado em gesttione-solutions/sql-library@61366e5
Observação: num repositório adotado com Baseline Profile próprio (`.roundfix/baseline/profiles/sql-library.json`, derivado de go-cli-tui, cuja lista `decisions` não tinha `verification.incremental`), `roundfix baseline update` (roundfix 0.136.0, catálogo sha256:91f1ab31… → sha256:a143bea7…) falha com:
  roundfix: baseline update failed: render managed entry "guide.agent-instructions": missing render value for token "verification.incremental"
  Category: invalid
  Next action: repair the repository or manifest input and rerun roundfix baseline update
Repetir com `--adopt-suggested` dá o mesmo erro, e `roundfix baseline profile validate sql-library` diz que o profile é válido.
Por que é um defeito: a documentação embutida da skill setup-context-driven promete que "A single-gate Setup Manifest exits 3, names the missing incremental decision, and writes nothing", com a correção de rodar de novo com `--adopt-suggested`. Com profile próprio isso não acontece. O template `template.guide.agent-instructions` passou a exigir um token que o profile não declara, e a falha sai como erro de render, sem dizer o que o mantenedor deve mudar.
Como contornei: acrescentei `verification.incremental` à lista `decisions` do profile próprio e rodei `roundfix baseline update --decision verification.incremental='make verify' --decision 'branch.prefix=<type>/' --no-history --confirm-plan <digest>`. O plano saiu limpo (7 arquivos) e aplicou. O commit é gesttione-solutions/sql-library@61366e5.
Resultado esperado: um de dois. (a) `profile validate` e `baseline update` detectam que o template selecionado exige uma decisão que o profile próprio não lista e dizem isso: o nome da decisão, a linha a acrescentar em `decisions` e o comando `--decision`. (b) Ou a decisão obrigatória nova entra sozinha no profile próprio. De quebra: `baseline.inventory.nested-carrier-conflict` acusa `docs/agents/sql.md` só porque a prosa cita o nome do marcador `setup-context-driven` entre crases, sem marcador de verdade; parece falso positivo.
Não precisa responder.
