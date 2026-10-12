package cli

// Suite: Delivery Queue Archive Relink staging and exact retirement proof.
// Invariant: every external change equals the rewrite of its parent bytes.
// Boundary IN: real Git repositories, Archive Command and queue archive stage.
// Boundary OUT: review, repository gates and publication.

import (
	"bytes"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

const deliveryRelinkPath = "docs/specs/0002-secondary/_prd.md"
const deliveryRelinkTask = "docs/specs/0002-secondary/task_01.md"
const deliveryRelinkPromotion = "docs/references/knowledge.md"

func deliveryRelinkFixture(t *testing.T, ambiguous bool) (*commandDeliveryWorkflow, string, string, string, string) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	w := retirementWorkflow(repo, home)
	source, destination, err := w.archivePaths(repo, implementTestSlug)
	if err != nil {
		t.Fatal(err)
	}
	mustMkdir(t, filepath.Join(repo, source, "references"))
	mustWrite(t, filepath.Join(repo, source, "references/knowledge.md"), "Reusable knowledge\n")
	if ambiguous {
		mustMkdir(t, filepath.Join(repo, source, "other"))
		mustWrite(t, filepath.Join(repo, source, "other/knowledge.md"), "Reusable knowledge\n")
	}
	mustMkdir(t, filepath.Dir(filepath.Join(repo, deliveryRelinkPath)))
	mustWrite(t, filepath.Join(repo, deliveryRelinkPath), "# Secondary\n\n[Source](../"+implementTestSlug+"/references/knowledge.md#lesson)\n[Owner](../"+implementTestSlug+"/_prd.md)\n")
	mustWrite(t, filepath.Join(repo, deliveryRelinkTask), "Read `"+source+"/references/knowledge.md` and `"+source+"/`.")
	commitArchiveFixture(t)
	parent := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	return w, repo, parent, source, destination
}

func archiveDeliveryRelinkFixture(t *testing.T, repo, parent string) {
	t.Helper()
	_, err := spec.Archive(spec.ArchiveRequest{SpecsRoot: filepath.Join(repo, "docs/specs"), BuiltInRoot: true, Slug: implementTestSlug, SourceRevision: parent, RepositoryRoot: repo})
	if err != nil {
		t.Fatal(err)
	}
}

func assertDeliveryRelinkProof(t *testing.T, w *commandDeliveryWorkflow, repo, parent, source, destination string, want bool) {
	t.Helper()
	ctx := commandContextForTest(t, t.Context())
	exact, err := w.archiveDiffIsExact(ctx, repo, source, destination)
	if err != nil || exact != want {
		t.Fatalf("working retirement=%v err=%v, want %v", exact, err, want)
	}
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "docs: archive relink fixture")
	head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	exact, err = w.archiveCommitIsExact(ctx, repo, parent, head, source, destination)
	if err != nil || exact != want {
		t.Fatalf("committed retirement=%v err=%v, want %v", exact, err, want)
	}
	reconciled, err := w.reconcileArchiveCommit(ctx, repo, implementTestSlug, parent, head)
	if err != nil || reconciled.ExactSpecMove != want {
		t.Fatalf("reconciled retirement=%+v err=%v, want exact %v", reconciled, err, want)
	}
}

func TestDeliveryArchiveStageCommitsRelinkedFiles(t *testing.T) {
	t.Parallel()
	w, repo, parent, source, destination := deliveryRelinkFixture(t, false)
	gittest.PersistIdentity(t, repo)
	setCommandEnvForTest(t, cliTestHelperEnv, "1")
	var log bytes.Buffer
	w.log = &log
	result, err := w.Archive(commandContextForTest(t, t.Context()), repo, implementTestSlug, parent)
	if err != nil || !result.ExactSpecMove || result.Parent != parent || result.Head == parent {
		t.Fatalf("archive stage=%+v err=%v", result, err)
	}
	r := archiveRecordAt(t, filepath.Join(repo, destination+".md"))
	if !reflect.DeepEqual(r.Promoted, []string{deliveryRelinkPromotion}) || !reflect.DeepEqual(r.Relinked, []string{deliveryRelinkPath, deliveryRelinkTask}) {
		t.Fatalf("promoted=%v relinked=%v", r.Promoted, r.Relinked)
	}
	paths := append([]string{destination + ".md"}, r.Promoted...)
	paths = append(paths, r.Relinked...)
	for _, path := range paths {
		if got := gittest.Run(t, repo, "show", result.Head+":"+path); got != mustRead(t, filepath.Join(repo, path)) {
			t.Fatalf("archive did not commit working bytes of %s", path)
		}
	}
	assertPathMissing(t, filepath.Join(repo, source))
	if status := strings.TrimSpace(gittest.Run(t, repo, "status", "--porcelain")); status != "" {
		t.Fatalf("dirty archive: %s", status)
	}
	exact, err := w.archiveCommitIsExact(commandContextForTest(t, t.Context()), repo, parent, result.Head, source, destination)
	if err != nil || !exact {
		t.Fatalf("committed stage proof=%v err=%v", exact, err)
	}
}

func TestDeliveryAcceptsARelinkedRetirement(t *testing.T) {
	t.Parallel()
	w, repo, parent, source, destination := deliveryRelinkFixture(t, false)
	archiveDeliveryRelinkFixture(t, repo, parent)
	wantPRD := "# Secondary\n\n[Source](../../references/knowledge.md#lesson)\n[Owner](../../history/specs/" + implementTestSlug + ".md)\n"
	if got := mustRead(t, filepath.Join(repo, deliveryRelinkPath)); got != wantPRD {
		t.Fatalf("relinked PRD=%q, want %q", got, wantPRD)
	}
	wantTask := "Read `" + deliveryRelinkPromotion + "` and `" + destination + ".md`."
	if got := mustRead(t, filepath.Join(repo, deliveryRelinkTask)); got != wantTask {
		t.Fatalf("relinked Task=%q, want %q", got, wantTask)
	}
	assertDeliveryRelinkProof(t, w, repo, parent, source, destination, true)
}

func TestDeliveryRefusesATamperedRelinkedFile(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"extra byte", "extra newline", "unlisted change", "unlisted relink", "missing promotion pairing", "history path", "new relinked path", "duplicate relink"} {
		t.Run(kind, func(t *testing.T) {
			w, repo, parent, source, destination := deliveryRelinkFixture(t, false)
			historyPath := spec.ArchiveDir(spec.ArchiveKindADR) + "/retired.md"
			if kind == "history path" {
				mustMkdir(t, filepath.Dir(filepath.Join(repo, historyPath)))
				mustWrite(t, filepath.Join(repo, historyPath), "Read `"+source+"/references/knowledge.md`.\n")
				commitArchiveFixture(t)
				parent = strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			}
			archiveDeliveryRelinkFixture(t, repo, parent)
			switch kind {
			case "extra byte", "extra newline":
				path := filepath.Join(repo, deliveryRelinkPath)
				extra := "x"
				if kind == "extra newline" {
					extra = "\n"
				}
				mustWrite(t, path, mustRead(t, path)+extra)
			case "unlisted change":
				mustWrite(t, filepath.Join(repo, "unrelated.md"), "Extra\n")
			case "unlisted relink", "missing promotion pairing", "history path", "new relinked path", "duplicate relink":
				recordPath := filepath.Join(repo, destination+".md")
				r := archiveRecordAt(t, recordPath)
				switch kind {
				case "unlisted relink":
					r.Relinked = r.Relinked[:1]
				case "missing promotion pairing":
					other := "docs/references/other.md"
					mustWrite(t, filepath.Join(repo, other), mustRead(t, filepath.Join(repo, deliveryRelinkPromotion)))
					r.Promoted = append(r.Promoted, other)
				case "history path", "new relinked path":
					path := historyPath
					if kind == "new relinked path" {
						path = "docs/references/new.md"
					}
					mustWrite(t, filepath.Join(repo, path), "Read `"+deliveryRelinkPromotion+"`.\n")
					r.Relinked = append(r.Relinked, path)
				case "duplicate relink":
					r.Relinked = append(r.Relinked, deliveryRelinkPath)
				}
				content, err := spec.RenderArchiveRecord(r)
				if err != nil {
					t.Fatal(err)
				}
				mustWrite(t, recordPath, string(content))
			}
			assertDeliveryRelinkProof(t, w, repo, parent, source, destination, false)
		})
	}
}

func TestDeliveryRefusesAnAmbiguousPromotionPairing(t *testing.T) {
	t.Parallel()
	w, repo, parent, source, destination := deliveryRelinkFixture(t, true)
	archiveDeliveryRelinkFixture(t, repo, parent)
	assertDeliveryRelinkProof(t, w, repo, parent, source, destination, false)
}
