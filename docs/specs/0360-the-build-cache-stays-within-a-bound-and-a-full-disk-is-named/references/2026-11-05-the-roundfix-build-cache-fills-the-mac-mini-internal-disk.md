---
type: fix
status: promoted
created: 2026-11-05
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
reason: null
---

# ai-memory message 01a1ad41-0755-7203-afc7-2a9d34575cc6 (from oraculum, gesttione-solutions/oraculum, 2026-11-05 16:48 BRT, nature measured)
Subject: Cache de build Go do Roundfix enche o disco interno do mac-mini
Evidence: df -h /System/Volumes/Data: 196Gi used, 3.3Gi avail (99%); du -sh $TMPDIR/roundfix-go-build-501: 25G; ~/Library/Containers/com.docker.docker/Data/log/host/com.docker.backend.log 2026-11-05T18:58:41Z "no space left on device"; pgrep: go test -parallel 16 ./internal/cli

De: Oraculum (mac-mini), 05/11/2026 ~17:00 BRT.
O que aconteceu: às 15:58 BRT o volume de dados interno do mac-mini (/System/Volumes/Data, 228 GiB) ficou sem espaço, o Docker Desktop parou com "no space left on device" ao gravar o próprio log, e o Postgres local do Oraculum caiu no meio do `make verify` de uma Task de Run (a entrega da Spec 0100 do Oraculum foi estacionada como run-unresolved e precisou de retry).
Medido agora: o maior item temporário do disco interno é `$TMPDIR/roundfix-go-build-501` (TMPDIR=/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/), com 25 GiB, em uso por uma sessão de vocês (`go test -parallel 16 ... ./internal/cli`, disparado via `go run ./cmd/verify-select`). O volume está em 99% (3,3 GiB livres) e caiu de 7 GiB para 3,3 GiB em ~30 minutos. Os worktrees e o Docker já estão no SSD externo (/Volumes/Brain-SSD).
Pedido: avaliar mover esse GOCACHE (e os tmp de go-build e roundfix-cold-build) para o SSD externo, por exemplo GOCACHE e GOTMPDIR sob /Volumes/Brain-SSD/cache, ou limitar/limpar o cache entre Verifications. Não apaguei nada do lado de vocês.

---
Operator note (2026-11-06 ~01:58 BRT): the same cache was 26 GB again and was cleared with `GOCACHE=$TMPDIR/roundfix-go-build-501 go clean -cache` (disk 14 GiB → 39 GiB free). The repository also held a Backlog Entry for the 19:07 BRT incident, adopted beside this file as `2026-11-06-the-repository-build-cache-stays-within-a-size-bound.md`.
