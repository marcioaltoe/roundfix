---
type: fix
status: promoted
created: 2026-11-06
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
reason: null
---

# ai-memory message 01a1b112-223d-71e0-8812-2975136c4111 (from argus, 2026-11-06 10:35 BRT, nature measured)
Subject: doctor: residue fica partial para sempre quando o PID do dono foi reutilizado
Evidence: roundfix 0.136.0 `roundfix doctor` em /Users/marcio/dev/argus, 2026-11-06 ~13:20 UTC; `ps -o pid,lstart,command -p 16788,57492,842`

O `roundfix doctor` 0.136.0, rodado em qualquer repositório (eu rodei no argus), sempre termina com `residue: partial`. Ele cita três Runs de outros repositórios, todas terminadas:
- `run_20261106T110411Z_8808faa6b603c904` (0002-suffixed-artifacts-require-their-role-folder, Clean), owner process 16788
- `run_20261104T212407Z_00b4e007d03f3127` (0061-apuracao-do-icms-com-varios-e100, Clean), owner process 57492
- `run_20261027T132639Z_4985b8e1c8d33f32` (0027-conferencia-de-estoque, Unresolved), owner process 842
O texto de cada uma é: `could not read process table for Run <id>: owner process <pid> failed step "prove owner process identity": owner process identity cannot be proven: live process start identity does not match the recorded owner identity`.
Por que o estado está errado: o `ps` mostra que hoje esses PIDs pertencem a processos do macOS sem relação com o Roundfix (842 TextInputMenuAgent desde 29/10; 16788 proactived e 57492 mediaanalysisd de 06/11). O dono gravado de cada Run já morreu, e o sistema reutilizou o PID. Quando a identidade de início do processo vivo não bate com a gravada, isso já prova que o dono original não existe mais. Mesmo assim, o doctor trata o caso como "não consegui ler" e devolve `partial`. Esse estado não tem saída: nenhum `reconcile` nem `gc` muda o resultado, e o aviso fica para sempre em todos os repositórios.
Correção sugerida: quando o start identity do PID vivo não bate, classifique o dono como morto (`owner gone: pid reused`) e não registre falha de leitura. Assim o residue fica `ok`, a não ser que haja outro resíduo de verdade. Guarde o `partial` para quando a tabela de processos realmente não pode ser lida.
