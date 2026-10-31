package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"roundfix/internal/judge"
)

const followUpUsage = `Usage:
  roundfix follow-up [--status] [--messages <path>] [--format <text|json>]

List this repository's pending work in Follow-up Priority order, or show its status.
It is read-only: it never writes a file, commits, pushes or prompts.
Options:
  --status    Show work in progress, pending Specs, the next release and blockers
  --messages  JSON file of ai-memory messages to list with the pending work
  --format    text (default) or json

Exits 0 for a report, including one with nothing to do; 2 for invalid usage,
configuration or messages file; 1 for an unreadable Run Database.
`

type followUpSource struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Count  int    `json:"count"`
	Detail string `json:"detail"`
}

type followUpItem struct {
	Rank         int               `json:"rank"`
	Priority     int               `json:"priority"`
	PriorityName string            `json:"priority_name"`
	Kind         string            `json:"kind"`
	Ref          string            `json:"ref"`
	Type         string            `json:"type"`
	Date         string            `json:"date"`
	Title        string            `json:"title"`
	RequestedBy  *string           `json:"requested_by"`
	Blocks       *string           `json:"blocks"`
	Basis        string            `json:"basis"`
	Answers      map[string]string `json:"answers"`
	Reason       *string           `json:"reason"`
	text         string
	origin       string
}

type followUpJudge struct {
	State   string  `json:"state"`
	Reason  *string `json:"reason"`
	Model   *string `json:"model"`
	Calls   int     `json:"calls"`
	CostUSD float64 `json:"cost_usd"`
}

type followUpRepository struct {
	Root      string `json:"root"`
	Namespace string `json:"namespace"`
}

type followUpDocument struct {
	Schema     string             `json:"schema"`
	View       string             `json:"view"`
	State      string             `json:"state"`
	Repository followUpRepository `json:"repository"`
	Sources    []followUpSource   `json:"sources"`
	Judge      followUpJudge      `json:"judge"`
	Items      []followUpItem     `json:"items"`
	InProgress followUpProgress   `json:"in_progress"`
}

type followUpRequest struct {
	messages, format string
	status           bool
}

func parseFollowUpCommand(args []string) (followUpRequest, error) {
	req := followUpRequest{format: "text"}
	fs := flagSet("follow-up")
	fs.BoolVar(&req.status, "status", false, "Show status")
	fs.StringVar(&req.messages, "messages", "", "Messages file")
	fs.StringVar(&req.format, "format", "text", "Output format")
	if err := fs.Parse(args); err != nil {
		return req, err
	}
	if fs.NArg() != 0 {
		return req, fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}
	if req.format != "text" && req.format != "json" {
		return req, fmt.Errorf("unknown --format %q; use text or json", req.format)
	}
	messagesSet := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "messages" {
			messagesSet = true
		}
	})
	if req.status && messagesSet {
		return req, fmt.Errorf("--messages applies only to the follow-up list; drop it with --status")
	}
	return req, nil
}

func runFollowUpCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, followUpUsage)
		return exitOK
	}
	fail := func(err error, code int) int {
		fmt.Fprintf(stderr, "roundfix: follow-up failed: %v\nRun 'roundfix follow-up --help' for usage.\n", err)
		return code
	}
	req, err := parseFollowUpCommand(args)
	if err != nil {
		return fail(err, exitPreflight)
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return fail(err, exitPreflight)
	}
	if loaded.GitRoot == "" {
		return fail(fmt.Errorf("follow-up requires a git repository working tree"), exitPreflight)
	}
	if req.status {
		return runFollowUpStatus(ctx, req, stdout, stderr, loaded, environment)
	}
	namespace := followUpNamespace(ctx, loaded)
	doc := followUpDocument{Schema: "roundfix/follow-up/v1", View: "follow-up", State: "work_found", Repository: followUpRepository{loaded.GitRoot, namespace}, Items: []followUpItem{}, Sources: []followUpSource{}}
	items, sources, err := followUpInventory(loaded, namespace, req.messages, environment.workDir)
	if err != nil {
		return fail(err, exitPreflight)
	}
	doc.Items, doc.Sources = items, sources
	progress, sources, err := readFollowUpProgress(ctx, loaded, environment)
	if err != nil {
		return fail(err, exitRunFailed)
	}
	doc.InProgress = progress
	doc.Sources = append(doc.Sources, sources...)
	doc.Judge = followUpJudge{State: "not_asked", Reason: followUpString("no item")}
	if len(items) > 0 {
		questions, err := judge.Load()
		if err != nil {
			return fail(err, exitPreflight)
		}
		keys := make(map[string]string)
		if commandRunID(environment.environ) == "" {
			for _, entry := range environment.environ {
				name, value, _ := strings.Cut(entry, "=")
				for _, allowed := range questions.KeyVariables() {
					if name == allowed {
						keys[name] = value
					}
				}
			}
		}
		inputs := make([]judge.FollowUpItem, 0, len(items))
		for _, item := range items {
			inputs = append(inputs, judge.FollowUpItem{Ref: item.Ref, Kind: item.Kind, Type: item.Type, Origin: item.origin, Title: item.Title, Text: item.text})
		}
		result, err := judge.PrioritizeFollowUp(ctx, questions.WithMonthlyCeiling(loaded.Config.Jev.MonthlyCeilingUSD), judge.FollowUpRequest{RunID: commandRunID(environment.environ), RepoRoot: loaded.GitRoot, Project: namespace, Items: inputs, Keys: keys, HomeDir: environment.homeDir, Transport: environment.dependencies.judgeTransport, Now: environment.dependencies.judgeNow})
		answers := make(map[string]judge.FollowUpAnswer)
		answered := 0
		for _, answer := range result.Answers {
			answers[answer.Ref] = answer
			if answer.Outcome == "answered" {
				answered++
			}
		}
		doc.Judge = followUpJudge{State: "unavailable", Reason: result.Skipped, Calls: result.Calls, CostUSD: result.CostUSD}
		if doc.Judge.Reason == nil {
			doc.Judge.Reason = result.Stopped
		}
		if err != nil {
			doc.Judge.Reason = followUpString(err.Error())
		}
		if answered > 0 {
			doc.Judge.State = "answered"
			if answered < len(items) {
				doc.Judge.State = "partial"
			}
			doc.Judge.Model = followUpString(result.Model)
			for _, answer := range result.Answers {
				if answer.Outcome == "answered" && answer.Model != nil {
					doc.Judge.Model = answer.Model
					break
				}
			}
		}
		if answered == 0 && doc.Judge.Reason == nil {
			for _, answer := range result.Answers {
				if answer.Reason != nil {
					doc.Judge.Reason = answer.Reason
					break
				}
			}
			if doc.Judge.Reason == nil {
				doc.Judge.Reason = followUpString("no answer")
			}
		}
		for _, item := range items {
			if _, ok := answers[item.Ref]; !ok {
				answers[item.Ref] = judge.FollowUpAnswer{Reason: doc.Judge.Reason}
			}
		}
		doc.Items = followUpOrder(items, answers)
	}
	if len(doc.Items) == 0 && progress.count() == 0 {
		doc.State = "nothing_to_do"
	}
	if req.format == "json" {
		if err := json.NewEncoder(stdout).Encode(doc); err != nil {
			return fail(fmt.Errorf("write follow-up JSON: %w", err), exitRunFailed)
		}
	} else {
		fmt.Fprint(stdout, renderFollowUpText(doc))
	}
	return exitOK
}

func followUpString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func followUpOrder(items []followUpItem, answers map[string]judge.FollowUpAnswer) []followUpItem {
	ordered := append([]followUpItem{}, items...)
	for index := range ordered {
		item := &ordered[index]
		answer := answers[item.Ref]
		item.Priority = 4
		item.Basis = "type"
		item.Answers = nil
		item.Blocks = nil
		item.Reason = answer.Reason
		if answer.Outcome == "answered" {
			item.Basis = "judge"
			item.Answers = answer.Verdicts
			item.Reason = nil
			yes := func(key string) bool { return answer.Verdicts[key] == "yes" }
			switch {
			case yes("blocks_other_project") && (yes("is_fix") || yes("is_small")):
				item.Priority = 1
			case yes("is_fix"):
				item.Priority = 2
			case yes("affects_production"):
				item.Priority = 3
			}
			if yes("blocks_other_project") {
				item.Blocks = item.RequestedBy
			}
		} else {
			item.Priority = followUpTypeClass(item.Type, item.Kind)
			if item.Reason == nil {
				item.Reason = followUpString("no answer")
			}
		}
		item.PriorityName = []string{"", "unblocking", "health", "feature", "other"}[item.Priority]
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		a, b := ordered[i], ordered[j]
		if a.Priority != b.Priority {
			return a.Priority < b.Priority
		}
		smallA, smallB := a.Answers["is_small"] == "yes", b.Answers["is_small"] == "yes"
		if smallA != smallB {
			return smallA
		}
		dateA, dateB := followUpDate(a.Date), followUpDate(b.Date)
		if dateA.IsZero() != dateB.IsZero() {
			return !dateA.IsZero()
		}
		if !dateA.Equal(dateB) {
			return dateA.Before(dateB)
		}
		return a.Ref < b.Ref
	})
	for index := range ordered {
		ordered[index].Rank = index + 1
	}
	return ordered
}

func followUpDate(value string) time.Time {
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date
		}
	}
	return time.Time{}
}

func followUpDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func renderFollowUpText(doc followUpDocument) string {
	var out strings.Builder
	fmt.Fprintf(&out, "Follow-up for %s in %s\n", doc.Repository.Namespace, doc.Repository.Root)
	if doc.State == "nothing_to_do" {
		out.WriteString("Nothing to do: no pending item and nothing in progress.\n")
	} else {
		fmt.Fprintf(&out, "Items: %d\n", len(doc.Items))
		for _, item := range doc.Items {
			fmt.Fprintf(&out, "  %d. %s %s %s %s: %s", item.Rank, item.PriorityName, item.Type, followUpDash(item.Date), item.Ref, item.Title)
			if item.RequestedBy != nil {
				fmt.Fprintf(&out, " — from %s", *item.RequestedBy)
			}
			if item.Blocks != nil {
				fmt.Fprintf(&out, " — blocks %s", *item.Blocks)
			}
			out.WriteByte('\n')
		}
		fmt.Fprintf(&out, "In progress: %d\n", doc.InProgress.count())
		for _, run := range doc.InProgress.Runs {
			fmt.Fprintf(&out, "  run %s %s %s\n", run.ID, run.State, followUpDash(run.Spec))
		}
		for _, item := range doc.InProgress.DeliveryQueue {
			fmt.Fprintf(&out, "  queue %d %s %s %s\n", item.Position, item.Spec, item.Stage, followUpDash(item.Blocker))
		}
		for _, pr := range doc.InProgress.PullRequests {
			fmt.Fprintf(&out, "  pull request #%d %s: %s\n", pr.Number, pr.Branch, pr.Title)
		}
		git := doc.InProgress.Git
		if git.hasChanges() {
			fmt.Fprintf(&out, "  git %s %s: %d changed path(s), %d unpushed commit(s)\n", git.Branch, git.Head, len(git.ChangedPaths), git.UnpushedCommits)
		}
		answered := 0
		for _, item := range doc.Items {
			if item.Basis == "judge" {
				answered++
			}
		}
		switch doc.Judge.State {
		case "answered", "partial":
			fmt.Fprintf(&out, "Judge: answered %d of %d item(s) with %s", answered, len(doc.Items), *doc.Judge.Model)
			if answered < len(doc.Items) {
				fmt.Fprintf(&out, "; %d ordered by type", len(doc.Items)-answered)
			}
			out.WriteByte('\n')
		case "not_asked":
			out.WriteString("Judge: not asked: no item\n")
		default:
			fmt.Fprintf(&out, "Judge: unavailable: %s; ordered by type\n", *doc.Judge.Reason)
		}
	}
	out.WriteString("Sources:\n")
	for _, source := range doc.Sources {
		fmt.Fprintf(&out, "  %s: %s", source.Name, source.State)
		if source.State == "checked" {
			fmt.Fprintf(&out, " %d", source.Count)
		}
		if source.Detail != "" {
			fmt.Fprintf(&out, ": %s", source.Detail)
		}
		out.WriteByte('\n')
	}
	return out.String()
}
