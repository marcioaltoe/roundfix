// Suite: auditing-binary Delivery Base evidence.
// Invariant: commit ancestry compares the build with the Delivery Base, never the candidate head.
// Boundary IN: a real Git repository, the supplied Delivery Base, and the binary identity.
// Boundary OUT: QA report rendering and warning publication, owned by internal/daemon.
package spec

import (
	"context"
	"testing"

	"roundfix/internal/app"
)

func TestAuditorEvidenceIsCurrentWhenTheBuildIsTheDeliveryBase(t *testing.T) {
	t.Parallel()
	root := initEvidenceRepo(t)
	base := commitEvidenceFile(t, root, "base.txt", "base\n")
	commitEvidenceFile(t, root, "candidate-one.txt", "one\n")
	commitEvidenceFile(t, root, "candidate-two.txt", "two\n")

	evidence := ResolveAuditorEvidence(context.Background(), root, base, app.AuditingBinary{Version: "0.17.0", Commit: base})

	if evidence.Ancestry != app.AncestryNotOlder {
		t.Fatalf("ancestry = %v, want %v", evidence.Ancestry, app.AncestryNotOlder)
	}
}

func TestAuditorEvidenceIsStaleWhenTheBuildPredatesTheDeliveryBase(t *testing.T) {
	t.Parallel()
	root := initEvidenceRepo(t)
	buildCommit := commitEvidenceFile(t, root, "build.txt", "build\n")
	deliveryBase := commitEvidenceFile(t, root, "base.txt", "base\n")

	evidence := ResolveAuditorEvidence(context.Background(), root, deliveryBase, app.AuditingBinary{Version: "0.17.0", Commit: buildCommit})

	if evidence.Ancestry != app.AncestryOlder {
		t.Fatalf("ancestry = %v, want %v", evidence.Ancestry, app.AncestryOlder)
	}
}

func TestAuditorEvidenceIsUnknownWithoutADeliveryBase(t *testing.T) {
	t.Parallel()
	root := initEvidenceRepo(t)
	buildCommit := commitEvidenceFile(t, root, "build.txt", "build\n")

	evidence := ResolveAuditorEvidence(context.Background(), root, "", app.AuditingBinary{Version: "0.17.0", Commit: buildCommit})

	if evidence.Ancestry != app.AncestryUnknown {
		t.Fatalf("ancestry = %v, want %v", evidence.Ancestry, app.AncestryUnknown)
	}
}

func TestAuditorEvidenceCarriesTheBinary(t *testing.T) {
	t.Parallel()
	root := initEvidenceRepo(t)
	deliveryBase := commitEvidenceFile(t, root, "base.txt", "base\n")
	binary := app.AuditingBinary{Version: "0.17.0", Commit: deliveryBase, Built: "2026-09-28T12:00:00Z"}

	evidence := ResolveAuditorEvidence(context.Background(), root, deliveryBase, binary)

	if evidence.Binary != binary {
		t.Fatalf("binary = %#v, want %#v", evidence.Binary, binary)
	}
	if evidence.DeliveryBase != deliveryBase {
		t.Fatalf("delivery base = %q, want %q", evidence.DeliveryBase, deliveryBase)
	}
}
