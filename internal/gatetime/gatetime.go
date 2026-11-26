// Package gatetime reports repository verification times without failing individual runs.
package gatetime

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

type Limits struct {
	Seconds       int `json:"seconds,omitempty"`
	MedianSeconds int `json:"medianSeconds,omitempty"`
	P95Seconds    int `json:"p95Seconds,omitempty"`
}
type Window struct {
	Kind    string     `json:"kind"`
	Event   string     `json:"event"`
	Size    int        `json:"size"`
	Since   *time.Time `json:"since"`
	Passing int        `json:"passing"`
	Closed  *time.Time `json:"closed"`
}
type Lab struct {
	Seconds             int `json:"seconds"`
	RegenerationSeconds int `json:"regenerationSeconds,omitempty"`
}
type Gate struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Machine  string  `json:"machine"`
	Approved Limits  `json:"approved"`
	Lab      *Lab    `json:"lab"`
	Window   Window  `json:"window"`
	Measured *Limits `json:"measured"`
	Number   *Limits `json:"number"`
}
type Document struct {
	SchemaVersion    int     `json:"schemaVersion"`
	AllowancePercent int     `json:"allowancePercent"`
	Rule             *string `json:"rule"`
	Gates            []Gate  `json:"gates"`
}

func Load(path string) (Document, error) {
	var doc Document
	data, err := os.ReadFile(path)
	if err != nil {
		return doc, err
	}
	err = json.Unmarshal(data, &doc)
	return doc, err
}
func Number(measured, approved int) (int, bool) {
	if measured < 0 || approved <= 0 || measured > approved {
		return 0, false
	}
	return min(approved, int(math.Ceil(float64(measured)*1.25/10))*10), true
}
func values(l Limits) []int { return []int{l.Seconds, l.MedianSeconds, l.P95Seconds} }
func Validate(doc Document, approved map[string]Limits) []error {
	var errs []error
	if doc.SchemaVersion != 1 || doc.AllowancePercent != 25 {
		errs = append(errs, fmt.Errorf("invalid schema or allowance"))
	}
	seen := map[string]bool{}
	for _, gate := range doc.Gates {
		want, ok := approved[gate.ID]
		if !ok || want != gate.Approved || seen[gate.ID] {
			errs = append(errs, fmt.Errorf("changed approved target %s", gate.ID))
		}
		seen[gate.ID] = true
		if (gate.Number == nil) != (gate.Measured == nil) {
			errs = append(errs, fmt.Errorf("unpaired number and measurement %s", gate.ID))
		}
		if gate.Number != nil && gate.Measured != nil {
			for i, n := range values(*gate.Number) {
				a := values(gate.Approved)[i]
				m := values(*gate.Measured)[i]
				if a == 0 {
					if n != 0 || m != 0 {
						errs = append(errs, fmt.Errorf("unexpected number %s", gate.ID))
					}
					continue
				}
				expected, reached := Number(m, a)
				if !reached || n != expected {
					errs = append(errs, fmt.Errorf("wrong number %s", gate.ID))
				}
			}
		}
	}
	for id := range approved {
		if !seen[id] {
			errs = append(errs, fmt.Errorf("missing gate %s", id))
		}
	}
	return errs
}
func limitText(l Limits) string {
	if l.Seconds > 0 {
		return fmt.Sprintf("%ds", l.Seconds)
	}
	return fmt.Sprintf("median %ds, p95 %ds", l.MedianSeconds, l.P95Seconds)
}
func Report(gate Gate, seconds int, job string) string {
	name := gate.Name
	if job != "" {
		name += ", job " + job
	}
	detail := "target: " + limitText(gate.Approved)
	if gate.Number != nil && gate.Measured != nil {
		m, n := gate.Measured.P95Seconds, gate.Number.P95Seconds
		if n == 0 {
			m, n = gate.Measured.Seconds, gate.Number.Seconds
		}
		state := "over the number"
		if seconds <= m {
			state = "within the measured result"
		} else if seconds <= n {
			state = "in the allowance"
		}
		detail = "number: " + limitText(*gate.Number) + "; " + state
	}
	return fmt.Sprintf("gate time: %s %ds (%s)", name, seconds, detail)
}

type testEvent struct {
	Action, Package, Test, Output string
	Elapsed                       float64
}

func renderTests(r io.Reader, w io.Writer, top int) (int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 4<<20)
	var tests []testEvent
	outputs := map[string][]string{}
	slowest := 0
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return 0, err
		}
		key := event.Package + "/" + event.Test
		if event.Action == "output" {
			outputs[key] = append(outputs[key], event.Output)
		}
		if event.Action != "pass" && event.Action != "fail" && event.Action != "skip" {
			continue
		}
		if event.Test == "" {
			fmt.Fprintf(w, "%s %s %.2fs\n", event.Action, event.Package, event.Elapsed)
			slowest = max(slowest, int(math.Ceil(event.Elapsed)))
		} else {
			tests = append(tests, event)
			if event.Action == "fail" {
				for _, line := range outputs[key] {
					fmt.Fprint(w, line)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	sort.SliceStable(tests, func(i, j int) bool { return tests[i].Elapsed > tests[j].Elapsed })
	fmt.Fprintln(w, "slowest tests:")
	for _, e := range tests[:min(top, len(tests))] {
		fmt.Fprintf(w, "%.2fs %s/%s\n", e.Elapsed, e.Package, e.Test)
	}
	return slowest, nil
}
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "gate-time: expected report, tests or windows")
		return 2
	}
	flags := flag.NewFlagSet("gate-time "+args[0], flag.ContinueOnError)
	flags.SetOutput(stderr)
	targets := flags.String("targets", "docs/references/gate-time-targets.json", "target document")
	gateID := flags.String("gate", "", "gate id")
	job := flags.String("job", "", "job name")
	seconds := flags.Int("seconds", -1, "elapsed seconds")
	since := flags.Int64("since", 0, "start Unix seconds")
	fixture := flags.String("fixture", "", "recorded windows")
	latest := flags.Bool("latest", false, "read latest windows")
	events := flags.String("events", "", "Go JSON events")
	top := flags.Int("top", 20, "slowest test count")
	summary := flags.String("summary", "", "summary output file")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		return 2
	}
	fail := func(err error) int { fmt.Fprintln(stderr, "gate-time: "+err.Error()); return 2 }
	if args[0] == "tests" {
		if *top < 0 {
			return fail(fmt.Errorf("top must be nonnegative"))
		}
		f, err := os.Open(*events)
		if err != nil {
			return fail(err)
		}
		defer f.Close()
		var text strings.Builder
		slowest, err := renderTests(f, &text, *top)
		if err != nil {
			return fail(err)
		}
		fmt.Fprint(stdout, text.String())
		if *summary != "" {
			f, err := os.OpenFile(*summary, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
			if err != nil {
				return fail(err)
			}
			_, writeErr := io.WriteString(f, text.String())
			closeErr := f.Close()
			if writeErr != nil {
				return fail(writeErr)
			}
			if closeErr != nil {
				return fail(closeErr)
			}
		}
		doc, err := Load(*targets)
		if err != nil {
			return fail(err)
		}
		for _, g := range doc.Gates {
			if g.ID == "test-binary" {
				fmt.Fprintln(stdout, Report(g, slowest, *job))
				return 0
			}
		}
		return fail(fmt.Errorf("unknown gate test-binary"))
	}
	doc, err := Load(*targets)
	if err != nil {
		return fail(err)
	}
	switch args[0] {
	case "report":
		if (*seconds >= 0) == (*since != 0) {
			return fail(fmt.Errorf("use seconds or since"))
		}
		if *since != 0 {
			*seconds = int(time.Now().Unix() - *since)
		}
		for _, g := range doc.Gates {
			if g.ID == *gateID {
				fmt.Fprintln(stdout, Report(g, *seconds, *job))
				return 0
			}
		}
		return fail(fmt.Errorf("unknown gate %q", *gateID))
	case "windows":
		var source Source
		if *fixture != "" {
			data, err := os.ReadFile(*fixture)
			if err != nil {
				return fail(err)
			}
			if err = json.Unmarshal(data, &source); err != nil {
				return fail(err)
			}
		} else {
			source, err = readSources(doc, *latest)
			if err != nil {
				return fail(err)
			}
		}
		lines, within := Windows(doc, source)
		for _, line := range lines {
			fmt.Fprintln(stdout, line)
		}
		if !within {
			return 1
		}
		return 0
	default:
		return fail(fmt.Errorf("unknown command %q", args[0]))
	}
}
