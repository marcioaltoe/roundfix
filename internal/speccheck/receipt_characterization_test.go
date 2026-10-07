package speccheck

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

func TestReceiptCharacterizationKeepsEveryParsedClaim(t *testing.T) {
	pinned := gittest.PinnedHistory(t, "../..", "docs/history/specs/0181-gates-that-refuse-only-what-someone-can-act-on", "docs/history/specs/0182-delivery-that-reviews-and-retries-from-where-the-item-stands")
	var got []Claim
	for _, item := range []struct{ name, slug, artifact string }{
		{"0181-prd.md", "0181-gates-that-refuse-only-what-someone-can-act-on", "_prd.md"},
		{"0181-techspec.md", "0181-gates-that-refuse-only-what-someone-can-act-on", "_techspec.md"},
		{"0182-prd.md", "0182-delivery-that-reviews-and-retries-from-where-the-item-stands", "_prd.md"},
		{"0182-techspec.md", "0182-delivery-that-reviews-and-retries-from-where-the-item-stands", "_techspec.md"},
	} {
		content, err := os.ReadFile(filepath.Join("testdata/receipt-characterization", item.name))
		if err != nil {
			t.Fatal(err)
		}
		artifact := "docs/history/specs/" + item.slug + "/" + item.artifact
		original, err := os.ReadFile(filepath.Join(pinned, artifact))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(content, original) {
			t.Fatalf("fixture %s differs from archived artifact", item.name)
		}
		got = append(got, CitationClaims(artifact, content)...)
	}
	path := "testdata/receipt-characterization/claims-golden.json"
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var want []Claim
	if err := json.Unmarshal(b, &want); err != nil {
		t.Fatal(err)
	}
	for i := range got {
		got[i].sentence = ""
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("claims = %#v, want %#v", got, want)
	}
}
