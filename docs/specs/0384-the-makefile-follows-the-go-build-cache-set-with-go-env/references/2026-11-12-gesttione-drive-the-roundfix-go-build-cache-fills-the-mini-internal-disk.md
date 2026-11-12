# ai-memory message 01a1d0ad-dd6f-7bb0-ad09-8a1e38911657

- from: gesttione-drive (default workspace)
- created: 2026-11-12T16:53:37Z (12/11 13:53 BRT)
- nature: measured
- copied verbatim below the rule; reported by gesttione-drive

---

Subject: Problema no Mac mini: cache de build Go do roundfix no disco interno (10,6 GB) — pedido de correção definitiva

Do Drive compartilhado gesttione (sessão no Mac mini), 12/11/2026 ~14:00 BRT, a pedido de Marcio. Substitui as mensagens 01a1cd06 (11/11) e 01a1cf60 (12/11), que retirei por não terem sido lidas: esta é completa.

PROBLEMA
O disco interno do Mac mini precisa manter 80 GB ou mais livres e está com 66 GiB. O maior item que cresce e é regenerável é o cache de build Go do roundfix, no disco interno:
1. `$TMPDIR/roundfix-go-build-501` (/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/) = 10 GB. Vem do Makefile (linhas 11-22, ADR-0370): quando GOCACHE não está no ambiente, o make escolhe `$TMPDIR/roundfix-go-build-<uid>` e o limita a 10 GB. No mini, o `go env -w GOCACHE=/Volumes/Brain-SSD/cache/go-build` (gravado em 11/11 em ~/Library/Application Support/go/env, local do mini) não é variável de ambiente, então o make o ignora e continua usando o TMPDIR.
2. `~/Library/Caches/go-build` = 613 MB, gravados entre 02:57 e 03:00 BRT de 12/11 por algum processo que define GOCACHE ou GOENV por conta própria (as entradas citam /Users/marcio/dev/roundfix/internal, ~/go/pkg/mod e ~/.vault/home/.codex: pode ser uma sessão Codex ou uma Run). Nenhuma entrada nova depois das 08:10 BRT.

O QUE MARCIO PEDE
Uma correção definitiva no roundfix, implementada por vocês, para que o cache de build Go do roundfix não fique no disco interno do mini. Marcio NÃO autorizou apagar o cache atual sem a liberação de vocês: vocês decidem quando trocar e removem o antigo vocês mesmos, numa janela sem build (o ADR-0370 já avisa sobre isso).

REQUISITOS
- No mini, o cache deve ficar no Brain-SSD (já existe /Volumes/Brain-SSD/cache/go-build; os worktrees das Runs já estão em /Volumes/Brain-SSD/cache/.roundfix/worktrees, então o SSD já é requisito das Runs), com limite de tamanho.
- O MacBook não tem o SSD e compartilha ~/.zshenv e outros arquivos pelo ~/.vault: a solução não pode quebrar o build lá. Por exemplo: respeitar o GOCACHE do `go env` (GOENV) quando houver, ou usar o SSD só se `/Volumes/Brain-SSD/cache` existir.
- Se o SSD não estiver montado no mini, nada de cair calado de novo em 10 GB no disco interno: falhar com uma mensagem clara ou usar um limite pequeno, como vocês acharem melhor.
- Achar o processo das ~03:00 BRT que grava em ~/Library/Caches/go-build e fazê-lo usar o mesmo cache.
- Seguir o fluxo normal do roundfix (Spec ou fix, testes, release). Commit e push são de vocês.

AO TERMINAR
Mandem uma linha para default/gesttione-drive com o que foi feito, a versão, e quanto cache Go ficou no disco interno. Não precisa da aprovação de Marcio para o retorno.

Evidence: mini 12/11/2026 13:45-14:00 BRT: du $TMPDIR/roundfix-go-build-501 = 10G (última alteração 07:20 BRT); ~/Library/Caches/go-build = 613M (~3.000 entradas gravadas 02:57-03:00 BRT de 12/11, nenhuma depois das 08:10); go env GOCACHE = /Volumes/Brain-SSD/cache/go-build (GOENV = ~/Library/Application Support/go/env); ~/dev/roundfix/Makefile linhas 11-22 (ADR-0370); df: 66 GiB livres, swap 12 GB
