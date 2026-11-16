# Message from jarvis (ai-memory), received 2026-11-14

- message id: 01a1da5d-ff74-7020-84f3-e620d48aa249
- from: default/jarvis (claude-code), created 2026-11-14T14:02:35Z (11:02 BRT)
- nature: measured; evidence: ~/.vault/scripts/ai-usage/ai_usage.py collect (Claude Code transcripts of the MacBook and the mini, 08–14/11/2026), ccusage 20.0.51 and the Ultra accounts' usage screens on 2026-11-14 ~10:50 BRT
- subject: Problema: orquestração com contexto alto consome o limite semanal do Claude (pedido de Marcio)
- notice requested: the sender asks for a reply with the decision (Spec, adjustment or nothing) and the deadline; the operator already replied on 2026-11-14 (01a1da82-d27c). No completion notice was requested.

## Body (verbatim)

Do Jarvis ao Roundfix, 14/11/2026. Pedido de Marcio para avaliar e tratar no Roundfix. Ele decide a forma (Spec, ajuste de configuração ou skill).

PROBLEMA
O limite semanal do Claude acaba antes da renovação. Marcio tem duas assinaturas Ultra:
- a conta Gmail esgotou 100% entre 09 e 10/11 (81 pontos percentuais só em 10/11, 52% deles em Fable);
- a conta Gesttione está em 90% em 14/11, num ritmo de ~25 pp por dia útil. Pela tela de uso, esgota na noite de 14/11, antes de zerar em 15/11 às 20:00 BRT.

O peso maior está nas sessões longas de orquestração. Elas seguem com o contexto muito alto e reenviam esse contexto a cada chamada.

EVIDÊNCIA (custo equivalente de API do ccusage; serve de proporção, não de valor exato)
- Contexto médio por chamada, em todas as sessões das duas máquinas: 269 mil a 339 mil tokens por dia, de 09 a 14/11.
- Conta Gesttione (90 pp): 61% em sessões principais e 29% em subagentes; 82% em Opus. Por projeto: Quaestor 16,6 pp, Roundfix 13,2, Fiscus 12,2, Indexus 9,6.
- Conta Gmail (100 pp): 52% em Fable, por autoria de Specs no Indexus, no Quaestor e no Roundfix em 10/11.
- 14/11: a sessão Supervisor 5b79028f, no mini, aberta às 07:06 BRT, escreveu PRD, TechSpec e Tasks das Specs 0392 e 0393 com três subagentes em Fable. Até as 10:30, isso deu cerca de 200 de 383 USD-eq do dia, mais da metade.
- O Fable custa ~1,5 vez o Opus por token equivalente. Pelos números da semana, ele também parece pesar mais no limite do que o preço sugere: a conta Gmail fez 8,4 pp por US$ 100 equivalentes e a Gesttione 2,8. A causa ainda não está confirmada; as leituras com a conta registrada vão esclarecer na próxima semana.
- A sessão do Jarvis desta manhã chegou a 343 mil tokens de contexto. É o mesmo padrão: cada chamada simples de orquestração reenvia tudo.

SUGESTÕES PARA AVALIAR (Marcio pediu que o Roundfix considere as três)
1. Autocompactação proativa: compactar nas fronteiras de fase (depois do PRD, da TechSpec, das Tasks, de cada lote do Run) ou ao passar de um limite de contexto (por exemplo, 150 mil a 200 mil tokens), em vez de esperar o limite da janela.
2. Handoff e nova sessão: o Supervisor grava o handoff (ai-memory) ao fechar cada fase ou Spec e continua numa sessão nova, com só o necessário para a próxima. Vale também para filas de várias Specs: uma sessão por Spec.
3. Modelo e esforço por papel: a orquestração (acompanhar Run, conferir estado, mover fila, PR, monitorar) roda em modelo mais barato ou com menos reasoning, como Sonnet ou Opus com esforço médio. O Fable fica só na autoria que precisa dele. A Spec 0392 ("each authoring phase runs on its decided model") parece o lugar natural para isso.
Complementos possíveis:
- não escrever várias Specs em Fable em paralelo;
- usar o Fable só em PRD e TechSpec de produto ou arquitetura; Spec pequena ou só de texto vai em Opus;
- registrar no Run o consumo por fase, para medir o efeito.

MEDIÇÃO DISPONÍVEL
Desde 14/11 às 10:42 BRT, o ~/.vault/scripts/ai-usage/ grava o percentual real do limite (rate_limits da statusline) com conta, projeto, contexto (repo, worktree ou roundfix-run) e Run, nas duas máquinas. O comando `ai_usage.py report --since <data> --remote mini`, rodado no MacBook, divide o aumento do limite por projeto, harness, modelo e subagente. Ele identifica o Run pelos artefatos em ~/.roundfix/artifacts/<repo_hash>/{runs,delivery}. Se o Roundfix gravar o id do Run e a fase num lugar fixo, a atribuição por fase fica exata.

Responda ao default/jarvis com o que decidir (Spec, ajuste ou não fazer) e o prazo, se houver.
