package baseline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProbe2(t *testing.T) {
	repo, plan := newHistoryMoveTransactionRepository(t, map[string]string{
		"_archived/specs/0001-alpha/_prd.md": "alpha prd\n",
	})
	move := plan.HistoryMoves[0]
	interrupted := beginTestTransaction(t, repo, plan)
	if err := interrupted.revalidatePreimages(context.Background()); err != nil {
		t.Fatalf("revalidate: %v", err)
	}
	if refusal, err := interrupted.relocateHistoryMove(context.Background(), 0); err != nil || refusal != nil {
		t.Fatalf("relocate = (%+v, %v)", refusal, err)
	}
	if err := os.Remove(filepath.Join(repo, filepath.FromSlash(move.To))); err != nil {
		t.Fatalf("remove dest: %v", err)
	}
	sidecar := filepath.Join(interrupted.stateDir, historyMoveContentName(move.Ordinal))
	os.WriteFile(sidecar, []byte("corrupted sidecar bytes\n"), 0o600)
	abandonTestTransaction(t, interrupted)
	recovered, err := BeginTransaction(context.Background(), repo, plan)
	fmt.Printf("PROBE2 recovered!=nil=%v err=%v containsMsg=%v\n", recovered != nil, err, err != nil && strings.Contains(err.Error(), "sidecar for ordinal 0 does not match the journal"))
}
