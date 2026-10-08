// Suite: Claim Receipts
// Invariant: written quotes are proved exactly; held claims need paragraph-local receipts.
// Boundary IN: exported readers, real files, CheckStage, and rendered findings.
// Boundary OUT: semantic support and transcript detection.
package speccheck

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func receiptProofFixture(t *testing.T, text string) string {
	t.Helper()
	root := t.TempDir()
	writeReceiptFixture(t, root, receiptFixtureADR, "# A citation is checked against what it cites\n\nRoundfix reads the cited record when an artifact makes a claim.\n")
	writeReceiptFixture(t, root, receiptFixturePRD, text)
	return root
}

func requireReceiptProof(t *testing.T, root, content string, proven bool, reason string) {
	t.Helper()
	receipts := Receipts(receiptFixturePRD, []byte(content))
	if len(receipts) != 1 {
		t.Fatalf("receipts = %#v", receipts)
	}
	_, got, why, err := ProveReceipt(root, receipts[0])
	if err != nil || got != proven || why != reason {
		t.Fatalf("proof = %v, %q, %v", got, why, err)
	}
	result := receiptStage(t, root, StagePRD)
	findings := receiptCodeFindings(result, CodeReceiptUnproven)
	want := 1
	if proven {
		want = 0
	}
	if len(findings) != want {
		t.Fatalf("unproven findings = %#v, want %d", findings, want)
	}
	if !proven && (findings[0].Severity != SeverityError || !strings.HasSuffix(findings[0].Summary, reason)) {
		t.Fatalf("finding = %#v", findings[0])
	}
}

func TestAVerbatimReceiptIsProven(t *testing.T) {
	t.Parallel()
	text := `ADR-0116: "reads the cited record"`
	root := receiptProofFixture(t, text)
	requireReceiptProof(t, root, text, true, "")
}

func TestAReceiptThatWrapsAcrossLinesIsProven(t *testing.T) {
	t.Parallel()
	for _, text := range []string{"ADR-0116:\n  \"reads the cited record\"", "ADR-0116: \"reads the\n  cited record\""} {
		t.Run(text, func(t *testing.T) {
			root := receiptProofFixture(t, text)
			writeReceiptFixture(t, root, receiptFixtureADR, "# Citation\nreads\t the cited\nrecord\n")
			requireReceiptProof(t, root, text, true, "")
		})
	}
}

func TestAReceiptWithOneChangedWordIsUnproven(t *testing.T) {
	t.Parallel()
	text := `ADR-0116: "reads the cited records"`
	root := receiptProofFixture(t, text)
	requireReceiptProof(t, root, text, false, "does not contain that text")
}

func TestAReceiptWithAnUnresolvedSourceIsUnproven(t *testing.T) {
	t.Parallel()
	text := `ADR-0999: "reads the cited record"`
	root := receiptProofFixture(t, text)
	requireReceiptProof(t, root, text, false, "does not resolve to a file")
}

func TestAReceiptThroughASymbolicLinkDoesNotResolve(t *testing.T) {
	t.Parallel()
	for _, outside := range []bool{false, true} {
		for _, directory := range []bool{false, true} {
			name := "inside/file"
			if outside {
				name = "outside/file"
			}
			if directory {
				name = strings.ReplaceAll(name, "file", "directory")
			}
			t.Run(name, func(t *testing.T) {
				source := "linked.md"
				if directory {
					source = "linked/source.md"
				}
				text := "`" + source + "`: \"reads the cited record\""
				root := receiptProofFixture(t, text)
				targetRoot := root
				if outside {
					targetRoot = t.TempDir()
				}
				writeReceiptFixture(t, targetRoot, "real/source.md", "reads the cited record")
				target := filepath.Join(targetRoot, "real/source.md")
				link := filepath.Join(root, "linked.md")
				if directory {
					target = filepath.Dir(target)
					link = filepath.Join(root, "linked")
				}
				if err := os.Symlink(target, link); err != nil {
					t.Fatal(err)
				}
				requireReceiptProof(t, root, text, false, "does not resolve to a file")
			})
		}
	}
}

func TestAReceiptShorterThanThreeWordsIsUnproven(t *testing.T) {
	t.Parallel()
	text := `ADR-0116: "cited record"`
	root := receiptProofFixture(t, text)
	requireReceiptProof(t, root, text, false, "a receipt needs at least three words")
}

func TestAFileReceiptIsProvedAgainstTheFile(t *testing.T) {
	t.Parallel()
	text := "`internal/example.go`: \"CitationClaims parses subject attributions\""
	root := receiptProofFixture(t, text)
	writeReceiptFixture(t, root, "internal/example.go", "// CitationClaims parses subject attributions\n")
	requireReceiptProof(t, root, text, true, "")
}

func TestAReceiptInsideAFencedBlockIsNotRead(t *testing.T) {
	t.Parallel()
	text := "```text\nADR-0116: \"words that do not exist\"\nADR-0116 requires a receipt.\n```\n"
	root := receiptProofFixture(t, text)
	if got := Receipts(receiptFixturePRD, []byte(text)); len(got) != 0 {
		t.Fatalf("receipts = %#v", got)
	}
	if got := ReceiptedClaims(receiptFixturePRD, []byte(text)); len(got) != 0 {
		t.Fatalf("claims = %#v", got)
	}
	result := receiptStage(t, root, StagePRD)
	if got := receiptCodeFindings(result, CodeReceiptUnproven); len(got) != 0 {
		t.Fatalf("findings = %#v", got)
	}
}

func heldReceiptFixture(t *testing.T, text string) string {
	t.Helper()
	root := receiptProofFixture(t, text)
	// A guide in a plain directory has unreadable history, which holds the Spec.
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	return root
}

func TestAHeldAttributionWithoutAReceiptIsAGap(t *testing.T) {
	t.Parallel()
	root := heldReceiptFixture(t, receiptFixtureClaim)
	result := receiptStage(t, root, StagePRD)
	findings := receiptCodeFindings(result, CodeReceiptMissing)
	if len(findings) != 1 || findings[0].Severity != SeverityGap {
		t.Fatalf("gaps = %#v", findings)
	}
}

func TestAReceiptForAnotherRecordDoesNotCoverTheClaim(t *testing.T) {
	t.Parallel()
	root := heldReceiptFixture(t, receiptFixtureClaim+` ADR-0117: "another valid passage"`)
	writeReceiptFixture(t, root, "docs/adr/0117-another.md", "# Another\nanother valid passage\n")
	result := receiptStage(t, root, StagePRD)
	if got := receiptCodeFindings(result, CodeReceiptMissing); len(got) != 1 {
		t.Fatalf("gaps = %#v", got)
	}
	if got := receiptCodeFindings(result, CodeReceiptUnproven); len(got) != 0 {
		t.Fatalf("unproven = %#v", got)
	}
}

func TestAHeldAttributionWithItsReceiptHasNoGap(t *testing.T) {
	t.Parallel()
	root := heldReceiptFixture(t, receiptFixtureClaim+` ADR-0116: "reads the cited record"`)
	result := receiptStage(t, root, StagePRD)
	for _, code := range []string{CodeReceiptMissing, CodeReceiptUnproven} {
		if got := receiptCodeFindings(result, code); len(got) != 0 {
			t.Fatalf("%s = %#v", code, got)
		}
	}
}

func TestAReceiptInAnotherParagraphDoesNotCoverTheClaim(t *testing.T) {
	t.Parallel()
	for _, separator := range []string{"\n\n", "\n- "} {
		t.Run(separator, func(t *testing.T) {
			root := heldReceiptFixture(t, receiptFixtureClaim+separator+`ADR-0116: "reads the cited record"`)
			if got := receiptCodeFindings(receiptStage(t, root, StagePRD), CodeReceiptMissing); len(got) != 1 {
				t.Fatalf("gaps = %#v", got)
			}
		})
	}
}

func TestASpecThatIsNotHeldStillProvesItsReceipts(t *testing.T) {
	t.Parallel()
	text := receiptFixtureClaim + ` ADR-0116: "reads the cited records"`
	root := receiptProofFixture(t, text)
	requireReceiptProof(t, root, text, false, "does not contain that text")
	if got := receiptCodeFindings(receiptStage(t, root, StagePRD), CodeReceiptMissing); len(got) != 0 {
		t.Fatalf("gaps = %#v", got)
	}
}

func TestReceiptFindingsRenderSurfaceTranscriptsOneAndTwo(t *testing.T) {
	t.Parallel()
	t.Run("one", func(t *testing.T) {
		root := heldReceiptFixture(t, `ADR-0116: "reads the cited records"`)
		report := RenderText(receiptStage(t, root, StagePRD), VerificationCoverage{})
		want := `[error] SC-RECEIPT-UNPROVEN: docs/specs/0200-example/_prd.md quotes ADR-0116 as "reads the cited records", but docs/adr/0116-a-citation-is-checked-against-what-it-cites.md does not contain that text
  at docs/specs/0200-example/_prd.md:1
  at docs/adr/0116-a-citation-is-checked-against-what-it-cites.md:1
  fix: Copy the passage verbatim from docs/adr/0116-a-citation-is-checked-against-what-it-cites.md into the receipt, or remove the receipt.`
		if !strings.Contains(report, want) {
			t.Fatalf("report lacks transcript finding:\n%s", report)
		}
	})
	t.Run("two", func(t *testing.T) {
		root := heldReceiptFixture(t, receiptFixtureClaim)
		report := RenderText(receiptStage(t, root, StagePRD), VerificationCoverage{})
		want := `[gap] SC-RECEIPT-MISSING: docs/specs/0200-example/_prd.md attributes "ADR-0116 requires the check to read the cited record." to ADR-0116 without a receipt
  at docs/specs/0200-example/_prd.md:1
  at docs/adr/0116-a-citation-is-checked-against-what-it-cites.md:1
  fix: Add ADR-0116: "<verbatim passage>" to the same paragraph in docs/specs/0200-example/_prd.md.`
		if !strings.Contains(report, want) {
			t.Fatalf("report lacks transcript finding:\n%s", report)
		}
	})
}

func TestReceiptPairingPreservesClaimsAndSourceLines(t *testing.T) {
	t.Parallel()
	text := "\n- " + receiptFixtureClaim + "\n  ADR-0116: \"reads the cited record\"\n  ADR-0116: \"when an artifact makes\"\n"
	claims := CitationClaims("artifact", []byte(text))
	paired := ReceiptedClaims("artifact", []byte(text))
	if len(paired) != 1 || !reflect.DeepEqual(paired[0].Claim, claims[0]) {
		t.Fatalf("paired = %#v, claims = %#v", paired, claims)
	}
	want := []Receipt{{Artifact: "artifact", Line: 3, Source: "ADR-0116", Quote: "reads the cited record"}, {Artifact: "artifact", Line: 4, Source: "ADR-0116", Quote: "when an artifact makes"}}
	if !reflect.DeepEqual(paired[0].Receipts, want) {
		t.Fatalf("receipts = %#v, want %#v", paired[0].Receipts, want)
	}
}

func TestReceiptProofIsCaseSensitiveAndDoesNotStripMarkup(t *testing.T) {
	t.Parallel()
	for _, quote := range []string{"Reads the cited record", "reads cited record", "reads the **cited** record"} {
		t.Run(quote, func(t *testing.T) {
			text := `ADR-0116: "` + quote + `"`
			root := receiptProofFixture(t, text)
			requireReceiptProof(t, root, text, false, "does not contain that text")
		})
	}
}

func TestReceiptPathsMustBeCleanRelativeRegularFiles(t *testing.T) {
	t.Parallel()
	for _, source := range []string{"../outside.md", "/etc/passwd", "internal/../source.md", "./source.md", "directory", "absent.md"} {
		t.Run(source, func(t *testing.T) {
			root := t.TempDir()
			writeReceiptFixture(t, root, "source.md", "three exact words")
			if err := os.Mkdir(filepath.Join(root, "directory"), 0755); err != nil {
				t.Fatal(err)
			}
			_, proven, reason, err := ProveReceipt(root, Receipt{Source: source, Quote: "three exact words"})
			if err != nil || proven || reason != "does not resolve to a file" {
				t.Fatalf("proof = %v %q %v", proven, reason, err)
			}
		})
	}
}

func TestReceiptsResolveArchivedAndInactiveRecords(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeReceiptFixture(t, root, "docs/history/adr/0116-archived.md", "---\nstatus: superseded\n---\nthree exact words\n")
	path, proven, reason, err := ProveReceipt(root, Receipt{Source: "ADR-0116", Quote: "three exact words"})
	if err != nil || !proven || reason != "" || path != "docs/history/adr/0116-archived.md" {
		t.Fatalf("proof = %q %v %q %v", path, proven, reason, err)
	}
}

func TestMissingReceiptsIgnoreUnresolvedAndInactiveClaims(t *testing.T) {
	t.Parallel()
	root := heldReceiptFixture(t, "ADR-0999 requires another rule. ADR-0116 requires the cited record.")
	writeReceiptFixture(t, root, receiptFixtureADR, "---\nstatus: proposed\n---\n# A citation\nreads the cited record\n")
	if got := receiptCodeFindings(receiptStage(t, root, StagePRD), CodeReceiptMissing); len(got) != 0 {
		t.Fatalf("gaps = %#v", got)
	}
}

func TestReceiptStagesReadOnlyAvailableAuthoringArtifacts(t *testing.T) {
	t.Parallel()
	root := heldReceiptFixture(t, "No attribution.")
	writeReceiptFixture(t, root, "docs/specs/0200-example/_techspec.md", receiptFixtureClaim+` ADR-0116: "reads the cited records"`)
	if got := receiptCodeFindings(receiptStage(t, root, StagePRD), CodeReceiptUnproven); len(got) != 0 {
		t.Fatalf("prd findings = %#v", got)
	}
	for _, stage := range []Stage{StageTechSpec, StageAll, StageTasks} {
		t.Run(string(stage), func(t *testing.T) {
			if got := receiptCodeFindings(receiptStage(t, root, stage), CodeReceiptUnproven); len(got) != 1 {
				t.Fatalf("findings = %#v", got)
			}
		})
	}
}

func TestNoPRDSkipsBothReceiptDetectors(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs/specs/0200-example"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, stage := range []Stage{StagePRD, StageTechSpec, StageAll} {
		result := receiptStage(t, root, stage)
		for _, code := range []string{CodeReceiptMissing, CodeReceiptUnproven} {
			found := false
			for _, skip := range result.Skipped {
				if skip.Code == code {
					found = true
				}
			}
			if !found {
				t.Fatalf("stage %q lacks skip %s: %#v", stage, code, result.Skipped)
			}
		}
	}
}

func TestReceiptSyntaxRejectsBareFieldsAndMoreThanOneLineBreak(t *testing.T) {
	t.Parallel()
	for _, text := range []string{
		"`field`: \"three exact words\"",
		"ADR-0116: \"three\nexact\nwords\"",
		"ADR-0116:\n \"three exact\nwords\"",
		"prefixADR-0116: \"three exact words\"",
	} {
		if got := Receipts("artifact", []byte(text)); len(got) != 0 {
			t.Fatalf("%q produced receipts %#v", text, got)
		}
	}
}

func TestReceiptSourcesRefuseAmbiguousRecordsAndADRDirectoryLinks(t *testing.T) {
	t.Parallel()
	t.Run("ambiguous record", func(t *testing.T) {
		root := receiptProofFixture(t, `ADR-0116: "reads the cited record"`)
		writeReceiptFixture(t, root, "docs/adr/0116-duplicate.md", "reads the cited record")
		requireReceiptProof(t, root, `ADR-0116: "reads the cited record"`, false, "does not resolve to a file")
	})
	for _, location := range []string{"docs/adr", "docs/history/adr"} {
		t.Run(location, func(t *testing.T) {
			root := t.TempDir()
			source := t.TempDir()
			writeReceiptFixture(t, source, "0116-record.md", "reads the cited record")
			if err := os.MkdirAll(filepath.Dir(filepath.Join(root, location)), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(source, filepath.Join(root, location)); err != nil {
				t.Fatal(err)
			}
			_, proven, reason, err := ProveReceipt(root, Receipt{Source: "ADR-0116", Quote: "reads the cited record"})
			if err != nil || proven || reason != "does not resolve to a file" {
				t.Fatalf("proof = %v %q %v", proven, reason, err)
			}
		})
	}
}

func TestShortReceiptIsRefusedBeforeReadingItsSource(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	_, proven, reason, err := ProveReceipt(root, Receipt{Source: "missing.md", Quote: "two words"})
	if err != nil || proven || reason != "a receipt needs at least three words" {
		t.Fatalf("proof = %v %q %v", proven, reason, err)
	}
}

func TestAFileInAReceiptDirectoryPositionDoesNotResolve(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeReceiptFixture(t, root, "file.md", "reads the cited record")
	_, proven, reason, err := ProveReceipt(root, Receipt{Source: "file.md/source.md", Quote: "reads the cited record"})
	if err != nil || proven || reason != "does not resolve to a file" {
		t.Fatalf("proof = %v %q %v", proven, reason, err)
	}
}

func TestABacktickedDecisionIdentifierIsNotAPathReceipt(t *testing.T) {
	t.Parallel()
	if got := Receipts("artifact", []byte("`ADR-0116`: \"reads the cited record\"")); len(got) != 0 {
		t.Fatalf("receipts = %#v, want no path receipt", got)
	}
}

func TestAReceiptSourceTokenCannotContainUnicodeWhitespace(t *testing.T) {
	t.Parallel()
	if got := Receipts("artifact", []byte("`docs/a\u00a0b.md`: \"reads the cited record\"")); len(got) != 0 {
		t.Fatalf("receipts = %#v", got)
	}
}

func TestAReceiptAllowsUnicodeWhitespaceAfterItsColon(t *testing.T) {
	t.Parallel()
	got := Receipts("artifact", []byte("ADR-0116:\u00a0\"reads the cited record\""))
	if len(got) != 1 || got[0].Quote != "reads the cited record" {
		t.Fatalf("receipts = %#v", got)
	}
}
