// Suite: saved cause measurement integrity.
// Invariant: the saved summary and single verdict are independently recomputable.
package runcause

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

var causesRecord = flag.String("causes-record", "", "Saved cause record to validate")
var causesDocument = flag.String("causes-document", "", "Measurement document to validate")

func checkCauseRecord(recordPath, documentPath string) error {
	data, err := os.ReadFile(recordPath)
	if err != nil {
		return fmt.Errorf("read cause record: %w", err)
	}
	doc, err := os.ReadFile(documentPath)
	if err != nil {
		return fmt.Errorf("read cause document: %w", err)
	}
	var report Report
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("decode cause record: %w", err)
	}
	table, err := Load()
	if err != nil {
		return err
	}
	if report.Schema != "roundfix/runs-causes/v1" {
		return fmt.Errorf("unknown cause record schema")
	}
	if report.SignaturesSHA256 != table.SHA256 {
		return fmt.Errorf("signature digest mismatch")
	}
	if report.Window.Since == nil || report.Window.Until == nil || report.Runs < 1 {
		return fmt.Errorf("record requires a closed window and at least one Run")
	}
	since, err := time.Parse(time.DateOnly, *report.Window.Since)
	if err != nil {
		return fmt.Errorf("parse record since: %w", err)
	}
	until, err := time.Parse(time.DateOnly, *report.Window.Until)
	if err != nil {
		return fmt.Errorf("parse record until: %w", err)
	}
	if !until.After(since) {
		return fmt.Errorf("record window is not increasing")
	}
	for _, item := range report.Items {
		if item.Class != "unclassified" && !slices.Contains(table.Classes, item.Class) {
			return fmt.Errorf("unknown item class %q", item.Class)
		}
	}
	summary := table.Summarize(report.Items)
	if !reflect.DeepEqual(report.Summary, summary) {
		return fmt.Errorf("summary disagrees with items: got %v, recomputed %v", report.Summary, summary)
	}
	verdict := "keep closed"
	if summary["items"] < 30 || summary["unclassified"]*3 > summary["items"] {
		verdict = "inconclusive"
	} else if summary["repository_knowledge"] > summary["implementation-defect"]+summary["environment"] {
		verdict = "reopen"
	}
	lines := []string{}
	for _, line := range strings.Split(string(doc), "\n") {
		if strings.HasPrefix(line, "Verdict: ") {
			lines = append(lines, line)
		}
	}
	if len(lines) != 1 || lines[0] != "Verdict: "+verdict {
		return fmt.Errorf("document verdict disagrees with rule: want one line %q, got %v", "Verdict: "+verdict, lines)
	}
	return nil
}

func TestCausesRecordIsConsistent(t *testing.T) {
	record, document := *causesRecord, *causesDocument
	if record == "" && document == "" {
		record, document = "testdata/causes-record.json", "testdata/causes-measurement.md"
	}
	if err := checkCauseRecord(record, document); err != nil {
		t.Fatal(err)
	}
}

func TestCausesRecordRejectsASabotagedFixture(t *testing.T) {
	t.Parallel()
	if err := checkCauseRecord("testdata/causes-record-sabotaged.json", "testdata/causes-measurement.md"); err == nil {
		t.Fatal("sabotaged summary accepted")
	}
	dir := t.TempDir()
	doc := dir + "/measurement.md"
	if err := os.WriteFile(doc, []byte("Verdict: reopen\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := checkCauseRecord("testdata/causes-record.json", doc); err == nil {
		t.Fatal("sabotaged verdict accepted")
	}
	for _, paths := range [][2]string{{dir + "/absent.json", doc}, {"testdata/causes-record.json", dir + "/absent.md"}} {
		if err := checkCauseRecord(paths[0], paths[1]); err == nil {
			t.Fatal("missing input accepted")
		}
	}
}

func TestCausesRecordDecisionBoundaries(t *testing.T) {
	t.Parallel()
	table, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		knowledge, defect, environment, unknown int
		verdict                                 string
	}{{20, 0, 3, 13, "inconclusive"}, {20, 0, 4, 12, "reopen"}, {15, 10, 5, 0, "keep closed"}, {20, 0, 0, 0, "inconclusive"}, {16, 10, 4, 0, "reopen"}} {
		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {
			start, end := "2026-10-01", "2026-10-02"
			record := Report{Schema: "roundfix/runs-causes/v1", SignaturesSHA256: table.SHA256, Window: Window{&start, &end}, Runs: 1}
			for _, group := range []struct {
				class string
				n     int
			}{{"repository-convention", c.knowledge}, {"implementation-defect", c.defect}, {"environment", c.environment}, {"unclassified", c.unknown}} {
				for range group.n {
					record.Items = append(record.Items, Item{Class: group.class})
				}
			}
			record.Summary = table.Summarize(record.Items)
			data, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			dir := t.TempDir()
			path, doc := dir+"/record.json", dir+"/doc.md"
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(doc, []byte("Verdict: "+c.verdict+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := checkCauseRecord(path, doc); err != nil {
				t.Fatal(err)
			}
		})
	}
}
