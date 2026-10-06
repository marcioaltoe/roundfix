package lighttier

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/openrouterkey"
)

func TestLightSpendLog(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	t.Run("append sum and modes", func(t *testing.T) {
		home := t.TempDir()
		line := SpendLine{Repository: "repo", RunID: "run", Spec: "demo", Task: "task_01", Session: "session", Model: "deepseek/flash", CostUSD: 0.25, CostSource: "opencode"}
		for i := 0; i < 2; i++ {
			if err := AppendSpend(home, now, line); err != nil {
				t.Fatal(err)
			}
		}
		if err := AppendSpend(home, now, SpendLine{CostSource: "unreported"}); err != nil {
			t.Fatal(err)
		}
		got, err := ReadMonth(home, now)
		if err != nil || got != 0.5 {
			t.Fatalf("sum=%v err=%v", got, err)
		}
		path := filepath.Join(home, ".roundfix/openrouter/implement/2026-10.jsonl")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		rows := bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n"))
		if len(rows) != 3 {
			t.Fatalf("lines=%d", len(rows))
		}
		var row SpendLine
		if err := json.Unmarshal(rows[0], &row); err != nil {
			t.Fatal(err)
		}
		line.Schema, line.Time = "roundfix/light-spend/v1", now
		if !reflect.DeepEqual(row, line) {
			t.Fatalf("row=%+v want=%+v", row, line)
		}
		for _, tc := range []struct {
			path string
			mode os.FileMode
		}{{path, 0600}, {filepath.Dir(path), 0700}} {
			info, err := os.Stat(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != tc.mode {
				t.Fatalf("%s mode=%o want=%o", tc.path, info.Mode().Perm(), tc.mode)
			}
		}
		if err := os.Chmod(path, 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := AppendSpend(home, now, line); err != nil {
			t.Fatal(err)
		}
		for _, tc := range []struct {
			path string
			mode os.FileMode
		}{{path, 0600}, {filepath.Dir(path), 0700}} {
			info, err := os.Stat(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != tc.mode {
				t.Fatalf("existing %s mode=%o want=%o", tc.path, info.Mode().Perm(), tc.mode)
			}
		}
	})
	t.Run("missing", func(t *testing.T) {
		home := t.TempDir()
		got, err := ReadMonth(home, now)
		if err != nil || got != 0 {
			t.Fatalf("sum=%v err=%v", got, err)
		}
		if _, err := os.Stat(filepath.Join(home, ".roundfix")); !os.IsNotExist(err) {
			t.Fatalf("read created storage: %v", err)
		}
	})
	for _, data := range []string{"not json\n", "{\"cost_usd\":0.25}\n{\n", "{\"cost_usd\":-0.1}\n", "null\n", "{}\n", "{\"cost_usd\":\"bad\"}\n", "{\"time\":\"invalid\",\"cost_usd\":0}\n"} {
		t.Run("unreadable "+data, func(t *testing.T) {
			home := t.TempDir()
			if err := AppendSpend(home, now, SpendLine{}); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(home, ".roundfix/openrouter/implement/2026-10.jsonl"), []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			if got, err := ReadMonth(home, now); err == nil || got != 0 {
				t.Fatalf("sum=%v err=%v", got, err)
			}
		})
	}
	t.Run("UTC month boundary", func(t *testing.T) {
		home := t.TempDir()
		local := time.Date(2026, 9, 30, 23, 30, 0, 0, time.FixedZone("west", -3*60*60))
		if err := AppendSpend(home, local, SpendLine{CostUSD: 0.75}); err != nil {
			t.Fatal(err)
		}
		if got, err := ReadMonth(home, now); err != nil || got != 0.75 {
			t.Fatalf("October sum=%v err=%v", got, err)
		}
		september := time.Date(2026, 9, 30, 23, 59, 59, 0, time.UTC)
		if err := AppendSpend(home, september, SpendLine{CostUSD: 0.5}); err != nil {
			t.Fatal(err)
		}
		if got, err := ReadMonth(home, september); err != nil || got != 0.5 {
			t.Fatalf("September sum=%v err=%v", got, err)
		}
		data, err := os.ReadFile(filepath.Join(home, ".roundfix/openrouter/implement/2026-10.jsonl"))
		if err != nil {
			t.Fatal(err)
		}
		var row SpendLine
		if err := json.Unmarshal(bytes.TrimSpace(data), &row); err != nil {
			t.Fatal(err)
		}
		if !row.Time.Equal(local) || row.Time.Location() != time.UTC {
			t.Fatalf("time=%v", row.Time)
		}
	})
}

func TestOpenModelImplementationSpendRecordNamesItsKeyVariable(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	home := t.TempDir()
	line := SpendLine{Model: "deepseek/deepseek-v4.1-flash", CostUSD: 0.25, CostSource: "opencode", KeyVariable: openrouterkey.Implement}
	if err := AppendSpend(home, now, line); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, ".roundfix/openrouter/implement/2026-10.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw := string(data)
	if !strings.Contains(raw, `"key_variable":"`+openrouterkey.Implement+`"`) {
		t.Fatalf("record does not name its key variable: %s", raw)
	}
	for _, sentinel := range []string{"implement-sentinel", "shared-sentinel"} {
		if strings.Contains(raw, sentinel) {
			t.Fatalf("record carries key value sentinel %q: %s", sentinel, raw)
		}
	}
	legacy := "{\"schema\":\"roundfix/light-spend/v1\",\"time\":\"" + now.Format(time.RFC3339Nano) + "\",\"cost_usd\":0.5,\"cost_source\":\"opencode\"}\n"
	if err := os.WriteFile(path, append(data, []byte(legacy)...), 0600); err != nil {
		t.Fatal(err)
	}
	total, err := ReadMonth(home, now)
	if err != nil || math.Abs(total-0.75) > 1e-9 {
		t.Fatalf("legacy record did not count: sum=%v err=%v", total, err)
	}
}
