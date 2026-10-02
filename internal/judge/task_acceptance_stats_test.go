package judge

// Suite: acceptance statistics and auditable records. Boundary: deterministic
// arithmetic and local files; no process credentials or external requests.
import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unicode/utf8"

	"roundfix/internal/runcause"
)

var (
	acceptanceRecordFlag   = flag.String("task-acceptance-record", "", "acceptance record to audit")
	acceptanceDocumentFlag = flag.String("task-acceptance-document", "", "measurement document to audit")
	acceptanceLabelsFlag   = flag.String("task-acceptance-labels", "", "cause record supplying labels")
)

type acceptanceStatistics struct {
	Answered       int     `json:"answered"`
	Repaired       int     `json:"repaired"`
	BaseRate       float64 `json:"base_rate"`
	AUROC          float64 `json:"auroc"`
	Lower          float64 `json:"ci_lower"`
	Upper          float64 `json:"ci_upper"`
	Draws          int     `json:"bootstrap_draws"`
	Discarded      int     `json:"bootstrap_discarded"`
	DiscardedShare float64 `json:"discarded_share"`
	RandomAUROC    float64 `json:"random_auroc"`
	LengthAUROC    float64 `json:"criteria_chars_auroc"`
	Flagged        int     `json:"flagged"`
	Precision      float64 `json:"precision"`
	Recall         float64 `json:"recall"`
}
type acceptancePoint struct {
	Spec          string
	Repaired      bool
	Score, Length float64
	Flagged       bool
}

func acceptanceAUROC(points []acceptancePoint) (float64, bool) {
	positives, negatives := 0, 0
	wins := 0.0
	for _, positive := range points {
		if !positive.Repaired {
			negatives++
			continue
		}
		positives++
		for _, negative := range points {
			if negative.Repaired {
				continue
			}
			if positive.Score > negative.Score {
				wins++
			} else if positive.Score == negative.Score {
				wins += 0.5
			}
		}
	}
	if positives == 0 || negatives == 0 {
		return 0, false
	}
	return wins / (float64(positives) * float64(negatives)), true
}
func acceptanceClusterInterval(points []acceptancePoint, q acceptanceQuestion) (float64, float64, int) {
	clusters := map[string][]acceptancePoint{}
	for _, p := range points {
		clusters[p.Spec] = append(clusters[p.Spec], p)
	}
	names := make([]string, 0, len(clusters))
	for name := range clusters {
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		return 0, 0, q.Bootstrap.Resamples
	}
	rng := rand.New(rand.NewPCG(q.Bootstrap.Seed, q.Bootstrap.Seed))
	draws := make([]float64, 0, q.Bootstrap.Resamples)
	discarded := 0
	for range q.Bootstrap.Resamples {
		sample := make([]acceptancePoint, 0, len(points))
		for range names {
			sample = append(sample, clusters[names[rng.IntN(len(names))]]...)
		}
		auc, ok := acceptanceAUROC(sample)
		if !ok {
			discarded++
			continue
		}
		draws = append(draws, auc)
	}
	if len(draws) == 0 {
		return 0, 0, discarded
	}
	sort.Float64s(draws)
	rank := func(p float64) float64 { return draws[int(math.Ceil(p*float64(len(draws))))-1] }
	return rank(0.025), rank(0.975), discarded
}
func computeAcceptanceStatistics(tasks []acceptanceTask, q acceptanceQuestion) acceptanceStatistics {
	// Canonical order makes the seeded baseline independent of record/map order.
	ordered := append([]acceptanceTask(nil), tasks...)
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Spec != ordered[j].Spec {
			return ordered[i].Spec < ordered[j].Spec
		}
		return ordered[i].Task < ordered[j].Task
	})
	points := []acceptancePoint{}
	for _, task := range ordered {
		if task.Outcome != "answered" || task.Noul == nil {
			continue
		}
		points = append(points, acceptancePoint{task.Spec, task.Repaired, 1 - *task.Noul, float64(task.CriteriaChars), *task.Noul < q.FlagBelow})
	}
	stats := acceptanceStatistics{Answered: len(points), Draws: q.Bootstrap.Resamples}
	trueFlags := 0
	for _, p := range points {
		if p.Repaired {
			stats.Repaired++
		}
		if p.Flagged {
			stats.Flagged++
			if p.Repaired {
				trueFlags++
			}
		}
	}
	if stats.Answered > 0 {
		stats.BaseRate = float64(stats.Repaired) / float64(stats.Answered)
	}
	if stats.Flagged > 0 {
		stats.Precision = float64(trueFlags) / float64(stats.Flagged)
	}
	if stats.Repaired > 0 {
		stats.Recall = float64(trueFlags) / float64(stats.Repaired)
	}
	stats.AUROC, _ = acceptanceAUROC(points)
	stats.Lower, stats.Upper, stats.Discarded = acceptanceClusterInterval(points, q)
	stats.DiscardedShare = float64(stats.Discarded) / float64(stats.Draws)
	baseline := append([]acceptancePoint(nil), points...)
	rng := rand.New(rand.NewPCG(q.RandomSeed, q.RandomSeed))
	for i := range baseline {
		baseline[i].Score = rng.Float64()
	}
	stats.RandomAUROC, _ = acceptanceAUROC(baseline)
	for i := range baseline {
		baseline[i].Score = baseline[i].Length
	}
	stats.LengthAUROC, _ = acceptanceAUROC(baseline)
	return stats
}
func acceptanceVerdict(status string, s acceptanceStatistics, q acceptanceQuestion) string {
	if status == "blocked" || s.Answered < 150 || s.Repaired < 20 || s.DiscardedShare > q.Bootstrap.MaxDiscarded {
		return "inconclusive"
	}
	if s.Lower >= 0.70 && s.Flagged >= 10 && s.Precision >= 2*s.BaseRate {
		return "adopt"
	}
	return "do not adopt"
}

func checkAcceptanceRecord(repo, home string, recordBytes, document, labels []byte) error {
	aq, questionBytes, err := loadAcceptanceQuestion()
	if err != nil {
		return err
	}
	q, err := Load()
	if err != nil {
		return err
	}
	var record acceptanceRecord
	if err := json.Unmarshal(recordBytes, &record); err != nil {
		return err
	}
	if record.Schema != "roundfix/task-acceptance-remeasurement/v1" || (record.Status != "blocked" && record.Status != "measured") || record.PinnedModel != q.PinnedModel {
		return errors.New("record schema, status or pin differs")
	}
	if record.LabelsHash != acceptanceHash(labels) {
		return errors.New("labels_sha256 differs")
	}
	if record.QuestionHash != acceptanceHash(questionBytes) {
		return errors.New("question_sha256 differs")
	}
	var cause runcause.Report
	if err := json.Unmarshal(labels, &cause); err != nil {
		return err
	}
	if cause.Schema != "roundfix/runs-causes/v1" {
		return errors.New("invalid labels schema")
	}
	labelMap := map[string]runcause.TaskCounts{}
	for _, label := range cause.Tasks {
		if label.QA || label.VerdictsFailed+label.VerdictsPassed == 0 {
			continue
		}
		id := label.Spec + "/" + label.Task
		if _, ok := labelMap[id]; ok {
			return fmt.Errorf("duplicate label %s", id)
		}
		labelMap[id] = label
	}
	roots, err := acceptanceRoots(repo, home)
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, task := range record.Tasks {
		id := task.Spec + "/" + task.Task
		label, ok := labelMap[id]
		if !ok || seen[id] {
			return fmt.Errorf("unknown or duplicate Task %s", id)
		}
		seen[id] = true
		if task.Failed != label.VerdictsFailed || task.Repaired != (label.VerdictsFailed >= 1) {
			return fmt.Errorf("label differs for %s", id)
		}
		if task.Outcome == "answered" {
			if !q.pinned(task.Model) {
				return fmt.Errorf("unpinned model for %s", id)
			}
			if task.Noul == nil || math.IsNaN(*task.Noul) || *task.Noul < 0 || *task.Noul > 1 {
				return fmt.Errorf("invalid noul for %s", id)
			}
		} else if task.Outcome != "skipped" || task.Noul != nil {
			return fmt.Errorf("invalid outcome for %s", id)
		}
		state, data, err := buildAcceptanceState(roots, task.Spec, task.Task, aq)
		if err != nil {
			return fmt.Errorf("rebuild %s: %w", id, err)
		}
		if task.StateHash != acceptanceHash(data) || task.CriteriaChars != utf8.RuneCountInString(state.Criteria) || !q.Language.isEnglish(state.Title+"\n"+state.Criteria) {
			return fmt.Errorf("state differs for %s", id)
		}
	}
	for _, excluded := range record.Excluded {
		id := excluded.Spec + "/" + excluded.Task
		if _, ok := labelMap[id]; !ok || seen[id] || excluded.Reason == "" {
			return fmt.Errorf("invalid exclusion %s", id)
		}
		seen[id] = true
		state, _, err := buildAcceptanceState(roots, excluded.Spec, excluded.Task, aq)
		if err == nil && q.Language.isEnglish(state.Title+"\n"+state.Criteria) {
			return fmt.Errorf("unjustified exclusion %s", id)
		}
	}
	if record.Status == "measured" && len(seen) != len(labelMap) {
		return errors.New("record omits eligible labels")
	}
	want := computeAcceptanceStatistics(record.Tasks, aq)
	// Marshal then compare each field: no statistic is silently omitted.
	wantBytes, _ := json.Marshal(want)
	var wire struct {
		Statistics map[string]float64 `json:"statistics"`
	}
	if err := json.Unmarshal(recordBytes, &wire); err != nil {
		return err
	}
	gotMap := wire.Statistics
	var wantMap map[string]float64
	if err := json.Unmarshal(wantBytes, &wantMap); err != nil {
		return err
	}
	if len(gotMap) != len(wantMap) {
		return errors.New("statistics fields differ")
	}
	for name, value := range wantMap {
		if _, ok := gotMap[name]; !ok {
			return fmt.Errorf("missing statistic %s", name)
		}
		if fmt.Sprintf("%.4f", gotMap[name]) != fmt.Sprintf("%.4f", value) {
			return fmt.Errorf("statistic %s differs: got %.4f, want %.4f", name, gotMap[name], value)
		}
	}
	if record.Verdict != acceptanceVerdict(record.Status, want, aq) {
		return errors.New("record verdict differs")
	}
	verdicts := []string{}
	for _, line := range strings.Split(string(document), "\n") {
		if strings.HasPrefix(line, "Verdict: ") {
			verdicts = append(verdicts, strings.TrimPrefix(line, "Verdict: "))
		}
	}
	if len(verdicts) != 1 || verdicts[0] != record.Verdict {
		return errors.New("document must carry exactly one matching Verdict: line")
	}
	return nil
}
func acceptanceRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
func acceptanceFixtureRepo(t *testing.T) string {
	t.Helper()
	repo, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	source := "testdata/task-acceptance/repo"
	err = filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		acceptanceWrite(t, filepath.Join(repo, relative), string(data))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Isolate Project Config discovery from the checkout's own configuration.
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0700); err != nil {
		t.Fatal(err)
	}
	return repo
}
func TestTaskAcceptanceRecordIsConsistent(t *testing.T) {
	record, document, labels := *acceptanceRecordFlag, *acceptanceDocumentFlag, *acceptanceLabelsFlag
	repo := acceptanceFixtureRepo(t)
	home := t.TempDir()
	if record != "" || document != "" || labels != "" {
		if record == "" || document == "" || labels == "" {
			t.Fatal("all three -task-acceptance-* flags are required")
		}
		repo = *measureRepo
		if repo == "" {
			var err error
			repo, err = filepath.Abs("../..")
			if err != nil {
				t.Fatal(err)
			}
		}
		// The explicit record audit resolves the same user/project Spec Root.
		home = os.Getenv("HOME")
		if home == "" {
			t.Fatal("HOME is required for explicit record audit")
		}
	} else {
		record = "testdata/task-acceptance/record.json"
		document = "testdata/task-acceptance/measurement.md"
		labels = "testdata/task-acceptance/labels.json"
	}
	if err := checkAcceptanceRecord(repo, home, acceptanceRead(t, record), acceptanceRead(t, document), acceptanceRead(t, labels)); err != nil {
		t.Fatal(err)
	}
}
func TestTaskAcceptanceRecordRejectsASabotagedFixture(t *testing.T) {
	err := checkAcceptanceRecord(acceptanceFixtureRepo(t), t.TempDir(), acceptanceRead(t, "testdata/task-acceptance/record-sabotaged.json"), acceptanceRead(t, "testdata/task-acceptance/measurement.md"), acceptanceRead(t, "testdata/task-acceptance/labels.json"))
	if err == nil || !strings.Contains(err.Error(), "statistic auroc differs") {
		t.Fatalf("sabotaged AUROC accepted or failed elsewhere: %v", err)
	}
}
func TestAUROCCountsTiesAsHalf(t *testing.T) {
	for _, tc := range []struct {
		name    string
		points  []acceptancePoint
		want    float64
		defined bool
	}{
		{"all ties", []acceptancePoint{{Repaired: true, Score: 0.4}, {Score: 0.4}}, 0.5, true},
		{"mixed ties", []acceptancePoint{{Repaired: true, Score: 0.6}, {Repaired: true, Score: 0.4}, {Score: 0.4}, {Score: 0.2}}, 0.875, true},
		{"reversed", []acceptancePoint{{Repaired: true, Score: 0.1}, {Score: 0.9}}, 0, true},
		{"one class", []acceptancePoint{{Repaired: true, Score: 0.1}}, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, defined := acceptanceAUROC(tc.points)
			if got != tc.want || defined != tc.defined {
				t.Fatalf("AUROC=(%v,%v), want (%v,%v)", got, defined, tc.want, tc.defined)
			}
		})
	}
}
func TestClusterBootstrapIsDeterministic(t *testing.T) {
	q, _, err := loadAcceptanceQuestion()
	if err != nil {
		t.Fatal(err)
	}
	// Two homogeneous Specs: half the cluster draws lack one class. Task-level
	// resampling would almost never discard a draw, even though it is repeatable.
	points := []acceptancePoint{}
	for range 20 {
		points = append(points, acceptancePoint{Spec: "positive", Repaired: true, Score: 0.8}, acceptancePoint{Spec: "negative", Score: 0.2})
	}
	lower, upper, discarded := acceptanceClusterInterval(points, q)
	l2, u2, d2 := acceptanceClusterInterval(points, q)
	if lower != l2 || upper != u2 || discarded != d2 {
		t.Fatal("same seed produced different interval")
	}
	if lower != 1 || upper != 1 || discarded < 800 || discarded > 1200 {
		t.Fatalf("not a whole-Spec bootstrap: interval=[%v,%v], discarded=%d", lower, upper, discarded)
	}
	// An all-tied sample has an analytically known interval, including unequal
	// cluster sizes; percentile boundaries remain ties at one half.
	ties := []acceptancePoint{{Spec: "a", Repaired: true, Score: 0.5}, {Spec: "a", Score: 0.5}, {Spec: "b", Repaired: true, Score: 0.5}, {Spec: "b", Score: 0.5}, {Spec: "b", Score: 0.5}}
	lower, upper, discarded = acceptanceClusterInterval(ties, q)
	if lower != 0.5 || upper != 0.5 || discarded != 0 {
		t.Fatalf("tied interval=(%v,%v,%d)", lower, upper, discarded)
	}
}
func TestTaskAcceptanceVerdictRules(t *testing.T) {
	q, _, err := loadAcceptanceQuestion()
	if err != nil {
		t.Fatal(err)
	}
	base := acceptanceStatistics{Answered: 150, Repaired: 20, BaseRate: 20.0 / 150, Lower: 0.70, Flagged: 10, Precision: 0.8}
	for _, tc := range []struct {
		name, status, want string
		change             func(*acceptanceStatistics)
	}{
		{"adopt boundary", "measured", "adopt", func(*acceptanceStatistics) {}},
		{"blocked", "blocked", "inconclusive", func(*acceptanceStatistics) {}},
		{"few Tasks", "measured", "inconclusive", func(s *acceptanceStatistics) { s.Answered = 149 }},
		{"few repaired", "measured", "inconclusive", func(s *acceptanceStatistics) { s.Repaired = 19 }},
		{"too many discarded", "measured", "inconclusive", func(s *acceptanceStatistics) { s.DiscardedShare = q.Bootstrap.MaxDiscarded + 0.001 }},
		{"discard boundary", "measured", "adopt", func(s *acceptanceStatistics) { s.DiscardedShare = q.Bootstrap.MaxDiscarded }},
		{"weak interval", "measured", "do not adopt", func(s *acceptanceStatistics) { s.Lower = 0.6999 }},
		{"few flagged", "measured", "do not adopt", func(s *acceptanceStatistics) { s.Flagged = 9 }},
		{"weak precision", "measured", "do not adopt", func(s *acceptanceStatistics) { s.Precision = 2*s.BaseRate - 0.001 }},
		{"precision boundary", "measured", "adopt", func(s *acceptanceStatistics) { s.Precision = 2 * s.BaseRate }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := base
			tc.change(&s)
			if got := acceptanceVerdict(tc.status, s, q); got != tc.want {
				t.Fatalf("verdict=%s, want %s", got, tc.want)
			}
		})
	}
}
func TestTaskAcceptanceStatisticsBaselinesAndThreshold(t *testing.T) {
	q, _, err := loadAcceptanceQuestion()
	if err != nil {
		t.Fatal(err)
	}
	f := func(v float64) *float64 { return &v }
	tasks := []acceptanceTask{
		{Spec: "a", Task: "task_01", Repaired: true, Noul: f(0.2), CriteriaChars: 100, Outcome: "answered"},
		{Spec: "a", Task: "task_02", Noul: f(0.5), CriteriaChars: 50, Outcome: "answered"},
		{Spec: "b", Task: "task_01", Repaired: true, Noul: f(0.4), CriteriaChars: 90, Outcome: "answered"},
		{Spec: "b", Task: "task_02", Noul: f(0.9), CriteriaChars: 40, Outcome: "answered"},
		{Spec: "b", Task: "task_03", Repaired: true, Noul: f(0.1), Outcome: "skipped"},
	}
	got := computeAcceptanceStatistics(tasks, q)
	if got.Answered != 4 || got.Repaired != 2 || got.AUROC != 1 || got.LengthAUROC != 1 || got.Flagged != 2 || got.Precision != 1 || got.Recall != 1 || got.BaseRate != 0.5 {
		t.Fatalf("statistics=%+v", got)
	}
	// Independent baseline replay checks seed and canonical Task ordering.
	rng := rand.New(rand.NewPCG(q.RandomSeed, q.RandomSeed))
	scores := []float64{rng.Float64(), rng.Float64(), rng.Float64(), rng.Float64()}
	wins := 0.0
	for _, i := range []int{0, 2} {
		for _, j := range []int{1, 3} {
			if scores[i] > scores[j] {
				wins++
			}
		}
	}
	if got.RandomAUROC != wins/4 {
		t.Fatalf("random baseline=%v, want %v", got.RandomAUROC, wins/4)
	}
	tasks[0], tasks[3] = tasks[3], tasks[0]
	if !reflect.DeepEqual(got, computeAcceptanceStatistics(tasks, q)) {
		t.Fatal("record order changed statistics")
	}
}

func TestTaskAcceptanceRecordChecksProvenanceAndEveryFigure(t *testing.T) {
	original := acceptanceRead(t, "testdata/task-acceptance/record.json")
	labels := acceptanceRead(t, "testdata/task-acceptance/labels.json")
	document := acceptanceRead(t, "testdata/task-acceptance/measurement.md")
	repo := acceptanceFixtureRepo(t)
	for _, tc := range []struct {
		name, want string
		change     func(*acceptanceRecord)
	}{
		{"labels digest", "labels_sha256", func(r *acceptanceRecord) { r.LabelsHash = "altered" }},
		{"question digest", "question_sha256", func(r *acceptanceRecord) { r.QuestionHash = "altered" }},
		{"answer model", "unpinned model", func(r *acceptanceRecord) { r.Tasks[0].Model = "jev-1.14.0" }},
		{"state hash", "state differs", func(r *acceptanceRecord) { r.Tasks[0].StateHash = "altered" }},
		{"criteria length", "state differs", func(r *acceptanceRecord) { r.Tasks[0].CriteriaChars++ }},
		{"repair label", "label differs", func(r *acceptanceRecord) { r.Tasks[0].Repaired = false }},
		{"missing Task", "omits eligible labels", func(r *acceptanceRecord) { r.Tasks = r.Tasks[1:] }},
		{"verdict", "record verdict differs", func(r *acceptanceRecord) { r.Verdict = "adopt" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var r acceptanceRecord
			if err := json.Unmarshal(original, &r); err != nil {
				t.Fatal(err)
			}
			tc.change(&r)
			changed, err := json.Marshal(r)
			if err != nil {
				t.Fatal(err)
			}
			err = checkAcceptanceRecord(repo, t.TempDir(), changed, document, labels)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("mutation accepted or failed elsewhere: %v", err)
			}
		})
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(original, &wire); err != nil {
		t.Fatal(err)
	}
	var stats map[string]float64
	if err := json.Unmarshal(wire["statistics"], &stats); err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(stats))
	for name := range stats {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		t.Run("figure "+name, func(t *testing.T) {
			var copyStats map[string]float64
			if err := json.Unmarshal(wire["statistics"], &copyStats); err != nil {
				t.Fatal(err)
			}
			copyStats[name]++
			data, err := json.Marshal(copyStats)
			if err != nil {
				t.Fatal(err)
			}
			changedWire := map[string]json.RawMessage{}
			for key, value := range wire {
				changedWire[key] = value
			}
			changedWire["statistics"] = data
			changed, err := json.Marshal(changedWire)
			if err != nil {
				t.Fatal(err)
			}
			err = checkAcceptanceRecord(repo, t.TempDir(), changed, document, labels)
			if err == nil || !strings.Contains(err.Error(), "statistic "+name+" differs") {
				t.Fatalf("figure accepted or failed elsewhere: %v", err)
			}
		})
	}
	t.Run("missing zero statistic", func(t *testing.T) {
		delete(stats, "bootstrap_discarded")
		data, err := json.Marshal(stats)
		if err != nil {
			t.Fatal(err)
		}
		wire["statistics"] = data
		changed, err := json.Marshal(wire)
		if err != nil {
			t.Fatal(err)
		}
		if err := checkAcceptanceRecord(repo, t.TempDir(), changed, document, labels); err == nil {
			t.Fatal("missing zero statistic accepted")
		}
	})
	t.Run("duplicate verdict line", func(t *testing.T) {
		duplicate := append(append([]byte(nil), document...), []byte("\nVerdict: inconclusive\n")...)
		err := checkAcceptanceRecord(repo, t.TempDir(), original, duplicate, labels)
		if err == nil || !strings.Contains(err.Error(), "exactly one") {
			t.Fatalf("duplicate verdict accepted: %v", err)
		}
	})
	t.Run("current Task differs", func(t *testing.T) {
		changedRepo, home, _ := acceptanceTestRepo(t)
		// Copy all fixture Tasks, then mutate only criteria of the first one.
		err := filepath.WalkDir(filepath.Join(repo, "docs/specs"), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			relative, err := filepath.Rel(repo, path)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			acceptanceWrite(t, filepath.Join(changedRepo, relative), string(data))
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(changedRepo, "docs/specs/0300-example/task_01.md")
		text := string(acceptanceRead(t, path))
		acceptanceWrite(t, path, strings.Replace(text, "exits with code 0", "exits with code 2", 1))
		err = checkAcceptanceRecord(changedRepo, home, original, document, labels)
		if err == nil || !strings.Contains(err.Error(), "state differs") {
			t.Fatalf("stale state accepted: %v", err)
		}
	})
}
