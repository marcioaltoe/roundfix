package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/queueforecast"
	"roundfix/internal/releaseplan"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type followUpPendingSpec struct {
	Spec           string  `json:"spec"`
	QueuePosition  *int    `json:"queue_position"`
	Stage          *string `json:"stage"`
	TasksCompleted int     `json:"tasks_completed"`
	TasksTotal     int     `json:"tasks_total"`
}

type followUpReleaseChange struct {
	Commit  string `json:"commit"`
	Subject string `json:"subject"`
}

type followUpRelease struct {
	State           string                  `json:"state"`
	Base            string                  `json:"base"`
	ProposedVersion string                  `json:"proposed_version"`
	Changes         []followUpReleaseChange `json:"changes"`
	Reason          *string                 `json:"reason"`
}

type followUpBlocker struct {
	Spec      string `json:"spec"`
	Blocker   string `json:"blocker"`
	ParkClass string `json:"park_class"`
	Next      string `json:"next"`
}

type followUpMaintainer struct {
	Kind   string `json:"kind"`
	Spec   string `json:"spec,omitempty"`
	Detail string `json:"detail"`
	Answer string `json:"answer,omitempty"`
}

type followUpForecastMeasure struct {
	Count         int    `json:"count"`
	LikelySeconds *int64 `json:"likely_seconds"`
	LimitSeconds  *int64 `json:"limit_seconds"`
}

type followUpForecastBasis struct {
	DeliveredSpecs   int                     `json:"delivered_specs"`
	Minimum          int                     `json:"minimum"`
	FirstEnd         *string                 `json:"first_end"`
	LastEnd          *string                 `json:"last_end"`
	LikelyPercentile int                     `json:"likely_percentile"`
	LimitPercentile  int                     `json:"limit_percentile"`
	StartToMerge     followUpForecastMeasure `json:"start_to_merge"`
	StartGap         followUpForecastMeasure `json:"start_gap"`
}

type followUpStatusDocument struct {
	Schema             string                      `json:"schema"`
	View               string                      `json:"view"`
	State              string                      `json:"state"`
	Repository         followUpRepository          `json:"repository"`
	Sources            []followUpSource            `json:"sources"`
	InProgress         followUpProgress            `json:"in_progress"`
	PendingSpecs       []followUpPendingSpec       `json:"pending_specs"`
	Queue              []followUpQueueRow          `json:"queue"`
	ForecastBasis      followUpForecastBasis       `json:"forecast_basis"`
	Release            followUpRelease             `json:"release"`
	Blockers           []followUpBlocker           `json:"blockers"`
	Maintainer         []followUpMaintainer        `json:"maintainer"`
	AuthoringQuestions followUpAuthoringQuestions  `json:"authoring_questions"`
	DelegatedDecisions []followUpDelegatedDecision `json:"delegated_decisions"`
	now                time.Time
	zone               *time.Location
}

func runFollowUpStatus(ctx context.Context, req followUpRequest, stdout, stderr io.Writer, loaded roundconfig.Loaded, environment commandEnvironment) int {
	return runFollowUpStatusWithHistoryGit(ctx, req, stdout, stderr, loaded, environment, preflight.ExecGitRunner{})
}

// The runner belongs to this read, not the command dependencies. The child
// owns only its Git result; Run processing stays on the calling goroutine.
func runFollowUpStatusWithHistoryGit(ctx context.Context, req followUpRequest, stdout, stderr io.Writer, loaded roundconfig.Loaded, environment commandEnvironment, runner preflight.GitRunner) int {
	now := environment.dependencies.followUpNow()
	zone := environment.dependencies.runsListLocation()
	if zone == nil {
		zone = time.Local
	}
	root, rootErr := roundconfig.ResolveSpecsRootOrAbsent(loaded, loaded.GitRoot)
	input := deliveryHistoryInput{Checkout: loaded.GitRoot, SpecRoot: root.Path, ArchiveRoot: spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), RetentionDays: loaded.Config.Store.RunRetentionDays, Now: now, Git: runner}
	type gitResult struct {
		read deliveryHistoryGitRead
		err  error
	}
	historyCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	result := make(chan gitResult, 1)
	go func(gitInput deliveryHistoryInput) {
		if rootErr != nil {
			result <- gitResult{err: rootErr}
			return
		}
		branch, err := runner.RunGit(historyCtx, gitInput.Checkout, "rev-parse", "--abbrev-ref", "HEAD")
		if err != nil {
			result <- gitResult{err: fmt.Errorf("read history branch: %w", err)}
			return
		}
		gitInput.Branch = strings.TrimSpace(branch)
		read, err := readDeliveryHistoryGit(historyCtx, gitInput)
		result <- gitResult{read, err}
	}(input)
	progress, sources, err := readFollowUpProgress(ctx, loaded, environment)
	if err != nil {
		cancel()
	}
	historyGit := <-result
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: follow-up failed: %v\nRun 'roundfix follow-up --help' for usage.\n", err)
		return exitRunFailed
	}
	remote, remoteErr := (preflight.ExecGitRunner{}).RunGit(ctx, loaded.GitRoot, "remote", "get-url", "origin")
	namespaceRoot, namespaceErr := roundconfig.RepositoryRoot(loaded.GitRoot)
	if namespaceErr != nil {
		namespaceRoot = loaded.GitRoot
	}
	origin := followUpQueueOrigin{remote, remoteErr}
	doc := followUpStatusDocument{
		Schema: "roundfix/follow-up/v1", View: "status", State: "work_found",
		Repository: followUpRepository{loaded.GitRoot, followUpNamespaceFromRemote(namespaceRoot, remote, remoteErr)},
		InProgress: progress, Sources: sources, now: now, zone: zone,
		Blockers: []followUpBlocker{}, Maintainer: []followUpMaintainer{},
	}
	pending, source := followUpPendingSpecs(loaded, progress)
	doc.AuthoringQuestions = readFollowUpAuthoringQuestions(ctx, loaded, runner)
	doc.PendingSpecs = pending
	doc.Sources = append(doc.Sources, source)
	for _, item := range progress.DeliveryQueue {
		if item.Stage == "parked" {
			doc.Blockers = append(doc.Blockers, followUpBlocker{item.Spec, item.Blocker, item.ParkClass, item.Next})
		}
	}
	if question, found := delivery.PendingQuestionFor(progress.queue); found {
		doc.Maintainer = append(doc.Maintainer, followUpMaintainer{Kind: "pending_question", Spec: question.SpecSlug, Detail: question.Blocker, Answer: question.Answer})
	}
	plan, planErr := releaseplan.Build(ctx, releaseplan.Request{}, newReleasePlanGitSource(environment.workDir, environment.dependencies.releasePlanGitRunner))
	doc.Release = followUpRelease{Changes: []followUpReleaseChange{}}
	releaseSource := followUpSourceResult("release-plan", 1, planErr)
	switch {
	case errors.Is(planErr, releaseplan.ErrInvalidReleaseRange):
		doc.Release.State = "no_unreleased_changes"
		releaseSource = followUpSourceResult("release-plan", 0, nil)
	case planErr != nil:
		doc.Release.State = "unavailable"
		doc.Release.Reason = followUpString(planErr.Error())
	default:
		doc.Release.State = string(plan.State)
		doc.Release.Base = plan.Base.Tag
		doc.Release.ProposedVersion = plan.ProposedVersion
		for _, change := range plan.Changes {
			doc.Release.Changes = append(doc.Release.Changes, followUpReleaseChange{change.CommitSHA, change.Subject})
		}
		switch plan.State {
		case releaseplan.StateApprovalRequired:
			doc.Maintainer = append(doc.Maintainer, followUpMaintainer{Kind: "release_approval", Detail: plan.Approval.Question})
		case releaseplan.StateManualClassificationRequired:
			doc.Release.Reason = followUpString(fmt.Sprintf(
				"%d commit(s) since %s need a manual release classification",
				len(plan.Classification.BlockingCommits), plan.Base.Tag))
			doc.Maintainer = append(doc.Maintainer, followUpMaintainer{Kind: "release_classification",
				Detail: fmt.Sprintf("classify them with roundfix release plan --from %s --to %s --impact <none|patch|minor|major> --reason <text>",
					plan.Base.Tag, plan.Target.Name)})
		}
	}
	doc.Sources = append(doc.Sources, releaseSource)
	input.Runs = progress.historyRuns
	history, historySource := deliveryHistoryFromTimes(input, historyGit.read, historyGit.err)
	active, skipped, _ := spec.ListActiveDetailed(root.Path)
	var delegatedErr error
	doc.DelegatedDecisions, delegatedErr = readFollowUpDelegatedDecisions(root.Path, active, skipped)
	if delegatedErr != nil {
		for i, source := range doc.Sources {
			if source.Name != "specs" {
				continue
			}
			detail := strings.ReplaceAll(delegatedErr.Error(), "\n", "; ")
			if source.Detail != "" {
				detail = source.Detail + "; " + detail
			}
			doc.Sources[i] = followUpSource{Name: "specs", State: "unavailable", Detail: detail}
		}
	}
	rows, forecastRows, queueSources := followUpQueueRows(followUpQueueInput{Context: ctx, Loaded: loaded, SpecRoot: root.Path, ArchiveRoot: input.ArchiveRoot, Progress: progress, Active: active, Skipped: skipped, History: history, Now: now, Zone: zone, Origin: &origin})
	basis := queueforecast.NewBasis(history)
	forecasts := queueforecast.Forecast(basis, forecastRows, now)
	for i, forecast := range forecasts {
		rows[i].Start = followUpForecastCell(forecast.Start)
		rows[i].End = followUpForecastCell(forecast.End)
	}
	doc.Queue = rows
	doc.ForecastBasis = followUpBasis(basis)
	doc.Sources = append(doc.Sources, queueSources...)
	doc.Sources = append(doc.Sources, historySource)
	if progress.count() == 0 && len(doc.PendingSpecs) == 0 && len(doc.Blockers) == 0 && len(doc.Maintainer) == 0 && (doc.Release.State == "no_unreleased_changes" || doc.Release.State == string(releaseplan.StateNoRelease)) {
		doc.State = "nothing_to_do"
	}
	if req.format == "json" {
		err = json.NewEncoder(stdout).Encode(doc)
	} else {
		_, err = fmt.Fprint(stdout, renderFollowUpStatusText(doc))
	}
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: follow-up failed: write status report: %v\nRun 'roundfix follow-up --help' for usage.\n", err)
		return exitRunFailed
	}
	return exitOK
}

func followUpPendingSpecs(loaded roundconfig.Loaded, progress followUpProgress) ([]followUpPendingSpec, followUpSource) {
	rows := []followUpPendingSpec{}
	root, active, err := readFollowUpSpecs(loaded)
	queued := make(map[string]store.DeliveryQueueItem)
	for _, item := range progress.queue.Items {
		queued[item.SpecSlug] = item
	}
	queue := append([]followUpQueueItem{}, progress.DeliveryQueue...)
	sort.Slice(queue, func(i, j int) bool { return queue[i].Position < queue[j].Position })
	for _, item := range queue {
		position, stage := item.Position, item.Stage
		rows = append(rows, followUpPendingSpec{Spec: item.Spec, QueuePosition: &position, Stage: &stage})
	}
	for _, item := range active {
		if _, found := queued[item.Slug]; !found {
			rows = append(rows, followUpPendingSpec{Spec: item.Slug})
		}
	}
	if err == nil {
		for index := range rows {
			var graph *spec.Graph
			var loadErr error
			if item, found := queued[rows[index].Spec]; found {
				archiveRoot := spec.ArchiveSpecRoot(root, root == filepath.Join(loaded.GitRoot, "docs", "specs"))
				graph, loadErr = followUpItemTaskGraph(nil, nil, loaded.GitRoot, root, archiveRoot, item)
			} else {
				graph, loadErr = spec.Load(root, rows[index].Spec)
			}
			if loadErr != nil {
				err = errors.Join(err, loadErr)
				continue
			}
			rows[index].TasksTotal = len(graph.Tasks)
			for _, task := range graph.Tasks {
				if task.Status == spec.StatusCompleted {
					rows[index].TasksCompleted++
				}
			}
		}
	}
	return rows, followUpSourceResult("specs", len(rows), err)
}

func renderFollowUpStatusText(doc followUpStatusDocument) string {
	var out strings.Builder
	fmt.Fprintf(&out, "Status for %s in %s\n", doc.Repository.Namespace, doc.Repository.Root)
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
	fmt.Fprintf(&out, "Pending Specs: %d\n", len(doc.PendingSpecs))
	for _, item := range doc.PendingSpecs {
		position, stage := "-", "not queued"
		if item.QueuePosition != nil {
			position = strconv.Itoa(*item.QueuePosition)
		}
		if item.Stage != nil {
			stage = *item.Stage
		}
		fmt.Fprintf(&out, "  %s %s %s %d/%d Tasks\n", position, item.Spec, stage, item.TasksCompleted, item.TasksTotal)
	}
	out.WriteString(renderFollowUpSpecQueue(doc))
	switch doc.Release.State {
	case "no_unreleased_changes":
		out.WriteString("Next release: no unreleased changes\n")
	case string(releaseplan.StateNoRelease):
		out.WriteString("Next release: none (no_release)\n")
	case "unavailable":
		fmt.Fprintf(&out, "Next release: unavailable: %s\n", *doc.Release.Reason)
	case string(releaseplan.StateManualClassificationRequired):
		fmt.Fprintf(&out, "Next release: not proposed (manual_classification_required): %s\n", *doc.Release.Reason)
	default:
		fmt.Fprintf(&out, "Next release: %s from %s (%s)\n", doc.Release.ProposedVersion, doc.Release.Base, doc.Release.State)
	}
	for _, change := range doc.Release.Changes {
		fmt.Fprintf(&out, "  %s %s\n", shortReleasePlanSHA(change.Commit), change.Subject)
	}
	fmt.Fprintf(&out, "Blockers: %d\n", len(doc.Blockers))
	for _, blocker := range doc.Blockers {
		fmt.Fprintf(&out, "  %s %s: %s\n", blocker.Spec, blocker.ParkClass, blocker.Next)
	}
	fmt.Fprintf(&out, "Waiting on the maintainer: %d\n", len(doc.Maintainer))
	for _, item := range doc.Maintainer {
		if item.Kind == "pending_question" {
			fmt.Fprintf(&out, "  question %s parked %s: %s\n", item.Spec, item.Detail, item.Answer)
		} else if item.Kind == "release_classification" {
			fmt.Fprintf(&out, "  release classification: %s\n", item.Detail)
		} else {
			fmt.Fprintf(&out, "  release %s: %s\n", doc.Release.ProposedVersion, item.Detail)
		}
	}
	source := "default"
	if doc.AuthoringQuestions.Source == "project-config" {
		source = "Project Config"
	}
	fmt.Fprintf(&out, "Open authoring questions: %s (%s)", doc.AuthoringQuestions.Preference, source)
	if doc.AuthoringQuestions.Checkout != "" {
		fmt.Fprintf(&out, "; this checkout's Project Config gives %s, which counts only once the default branch holds it", doc.AuthoringQuestions.Checkout)
	}
	if doc.AuthoringQuestions.Detail != "" {
		fmt.Fprintf(&out, "; %s", doc.AuthoringQuestions.Detail)
	}
	out.WriteByte('\n')
	out.WriteString(renderFollowUpDelegatedDecisions(doc.DelegatedDecisions))
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

func followUpTime(at time.Time) *string {
	if at.IsZero() {
		return nil
	}
	value := at.UTC().Format(time.RFC3339)
	return &value
}

func followUpForecastCell(cell queueforecast.Cell) followUpQueueCell {
	out := followUpQueueCell{State: string(cell.State)}
	switch cell.State {
	case queueforecast.Actual:
		out.Actual = followUpTime(cell.Actual)
	case queueforecast.Forecasted:
		out.Likely, out.Limit = followUpTime(cell.Likely), followUpTime(cell.Limit)
	case queueforecast.Condition:
		out.Condition = cell.Condition
		if cell.HasAfter {
			likely, limit := int64(cell.LikelyAfter/time.Second), int64(cell.LimitAfter/time.Second)
			out.LikelyAfterSeconds, out.LimitAfterSeconds = &likely, &limit
		}
	case queueforecast.None:
		out.Reason = cell.Reason
	}
	return out
}

func followUpBasis(basis queueforecast.Basis) followUpForecastBasis {
	measure := func(values []time.Duration) followUpForecastMeasure {
		m := followUpForecastMeasure{Count: len(values)}
		if likely, limit, ok := queueforecast.Pair(values); ok {
			l, u := int64(likely/time.Second), int64(limit/time.Second)
			m.LikelySeconds, m.LimitSeconds = &l, &u
		}
		return m
	}
	return followUpForecastBasis{DeliveredSpecs: basis.Delivered, Minimum: 11, FirstEnd: followUpTime(basis.FirstEnd), LastEnd: followUpTime(basis.LastEnd), LikelyPercentile: 50, LimitPercentile: 85, StartToMerge: measure(basis.StartToMerge), StartGap: measure(basis.StartGap)}
}

func followUpQueueDuration(seconds int64, limit bool) string {
	minutes := seconds / 60
	if limit && seconds%60 != 0 {
		minutes++
	}
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%dh%02dm", minutes/60, minutes%60)
}

func followUpQueueInstant(value *string, zone *time.Location, limit bool) string {
	if value == nil {
		return "-"
	}
	at, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return "-"
	}
	if limit && !at.Equal(at.Truncate(time.Minute)) {
		at = at.Truncate(time.Minute).Add(time.Minute)
	}
	return at.In(zone).Format("01-02 15:04")
}

func followUpQueueCellText(cell followUpQueueCell, zone *time.Location) string {
	switch cell.State {
	case "actual":
		return followUpQueueInstant(cell.Actual, zone, false)
	case "forecast":
		return "~" + followUpQueueInstant(cell.Likely, zone, false) + ", by " + followUpQueueInstant(cell.Limit, zone, true)
	case "condition":
		if cell.LikelyAfterSeconds != nil && cell.LimitAfterSeconds != nil {
			return "+" + followUpQueueDuration(*cell.LikelyAfterSeconds, false) + ", by +" + followUpQueueDuration(*cell.LimitAfterSeconds, true) + " after its start"
		}
		return followUpDash(cell.Condition)
	default:
		return followUpDash(cell.Reason)
	}
}

func renderFollowUpSpecQueue(doc followUpStatusDocument) string {
	var out strings.Builder
	fmt.Fprintf(&out, "Spec queue: %d\n", len(doc.Queue))
	if len(doc.Queue) == 0 {
		return out.String()
	}
	zone := doc.zone
	if zone == nil {
		zone = time.UTC
	}
	name := doc.now.In(zone).Format("MST")
	table := [][]string{{"#", "Spec", "Status", "Start (" + name + ")", "End (" + name + ")", "Note"}}
	hasSpec := false
	for _, row := range doc.Queue {
		number := "-"
		if row.Kind != "backlog" {
			number = followUpSpecNumber(row.Spec)
			hasSpec = true
		}
		table = append(table, []string{strconv.Itoa(row.Row), number, row.Status, followUpQueueCellText(row.Start, zone), followUpQueueCellText(row.End, zone), followUpDash(row.Note)})
	}
	widths := make([]int, 5)
	for _, row := range table {
		for col := range widths {
			widths[col] = max(widths[col], utf8.RuneCountInString(row[col]))
		}
	}
	for _, row := range table {
		out.WriteString("  ")
		for col, value := range row {
			out.WriteString(value)
			if col < len(widths) {
				out.WriteString(strings.Repeat(" ", widths[col]-utf8.RuneCountInString(value)+2))
			}
		}
		out.WriteByte('\n')
	}
	if !hasSpec {
		return out.String()
	}
	basis := doc.ForecastBasis
	measure := basis.StartToMerge
	if measure.LikelySeconds == nil || measure.LimitSeconds == nil {
		fmt.Fprintf(&out, "Forecast basis: not enough history, %d of 11 delivered Specs.\n", basis.DeliveredSpecs)
		return out.String()
	}
	date := func(value *string) string {
		at, _ := time.Parse(time.RFC3339, *value)
		return at.In(zone).Format("01-02")
	}
	fmt.Fprintf(&out, "Forecast basis: %d Specs delivered from %s to %s on this machine. Likely is the median and the limit the 85th percentile: %s and %s from start to merge", basis.DeliveredSpecs, date(basis.FirstEnd), date(basis.LastEnd), followUpQueueDuration(*measure.LikelySeconds, false), followUpQueueDuration(*measure.LimitSeconds, true))
	gap := basis.StartGap
	if gap.LikelySeconds != nil && gap.LimitSeconds != nil {
		fmt.Fprintf(&out, ", %s and %s from one start to the next", followUpQueueDuration(*gap.LikelySeconds, false), followUpQueueDuration(*gap.LimitSeconds, true))
	}
	out.WriteString(".\n")
	return out.String()
}
