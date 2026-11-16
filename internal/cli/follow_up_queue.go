package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/phaserecord"
	"roundfix/internal/preflight"
	"roundfix/internal/queueforecast"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type followUpQueueRow struct {
	Row            int               `json:"row"`
	Kind           string            `json:"kind"`
	Spec           string            `json:"spec"`
	Ref            string            `json:"ref"`
	Title          string            `json:"title"`
	Status         string            `json:"status"`
	TasksCompleted int               `json:"tasks_completed"`
	TasksTotal     int               `json:"tasks_total"`
	Note           string            `json:"note"`
	Start          followUpQueueCell `json:"start"`
	End            followUpQueueCell `json:"end"`
}

type followUpQueueCell struct {
	State              string  `json:"state"`
	Actual             *string `json:"actual"`
	Likely             *string `json:"likely"`
	Limit              *string `json:"limit"`
	Condition          string  `json:"condition"`
	LikelyAfterSeconds *int64  `json:"likely_after_seconds"`
	LimitAfterSeconds  *int64  `json:"limit_after_seconds"`
	Reason             string  `json:"reason"`
}

type followUpQueueOrigin struct {
	remote string
	err    error
}

type followUpQueueInput struct {
	Context               context.Context
	Loaded                roundconfig.Loaded
	SpecRoot, ArchiveRoot string
	Progress              followUpProgress
	Active                []spec.Spec
	Skipped               []spec.SkippedSpec
	History               []queueforecast.SpecTimes
	Now                   time.Time
	Zone                  *time.Location
	Origin                *followUpQueueOrigin
	Git                   preflight.GitRunner
}

func followUpQueueRows(in followUpQueueInput) ([]followUpQueueRow, []queueforecast.Row, []followUpSource) {
	rows := []followUpQueueRow{}
	forecasts := []queueforecast.Row{}
	starts := map[string]time.Time{}
	for _, run := range in.Progress.historyRuns {
		if run.Kind != store.KindImplement || run.SpecSlug == "" {
			continue
		}
		if old := starts[run.SpecSlug]; old.IsZero() || run.CreatedAt.Before(old) {
			starts[run.SpecSlug] = run.CreatedAt
		}
	}
	for _, item := range in.History {
		if old := starts[item.Spec]; old.IsZero() || item.Start.Before(old) {
			starts[item.Spec] = item.Start
		}
	}
	seen := map[string]bool{}
	addSpec := func(slug, kind, status, note string, parked bool, wait string) {
		row := followUpQueueRow{Spec: slug, Kind: kind, Status: status, Note: followUpQueueText(note)}
		if kind != "authoring" {
			if graph, err := spec.Load(in.SpecRoot, slug); err == nil {
				row.TasksTotal = len(graph.Tasks)
				for _, task := range graph.Tasks {
					if task.Status == spec.StatusCompleted {
						row.TasksCompleted++
					}
				}
			}
			row.Status = fmt.Sprintf("%s %d/%d Tasks", status, row.TasksCompleted, row.TasksTotal)
		}
		rows = append(rows, row)
		forecasts = append(forecasts, queueforecast.Row{Label: followUpSpecNumber(slug), Planned: kind != "authoring", Queued: kind == "delivery_queue", Started: starts[slug], Parked: parked, Wait: wait})
		seen[slug] = true
	}
	queue := append([]store.DeliveryQueueItem(nil), in.Progress.queue.Items...)
	sort.SliceStable(queue, func(i, j int) bool { return queue[i].Position < queue[j].Position })
	for _, item := range queue {
		if item.Stage == store.DeliveryStageMerged {
			seen[item.SpecSlug] = true
			continue
		}
		parked := item.Stage == store.DeliveryStageParked
		note := followUpFirstLine(item.Blocker)
		if parked {
			park := delivery.ClassifyPark(in.Progress.queue, item)
			note = string(park.Class) + ": " + park.Next
		}
		wait := ""
		if starts[item.SpecSlug].IsZero() {
			wait = followUpQueueWait(item)
		}
		addSpec(item.SpecSlug, "delivery_queue", string(item.Stage), note, parked, wait)
	}
	active := append([]spec.Spec(nil), in.Active...)
	sort.Slice(active, func(i, j int) bool { return active[i].Slug < active[j].Slug })
	for _, item := range active {
		if !seen[item.Slug] {
			addSpec(item.Slug, "planned", "not queued", "", false, "")
		}
	}
	phases, _, phaseErr := phaserecord.Phases(in.Loaded.HomeDir, in.Now.AddDate(0, 0, -7))
	phaseSource := followUpSourceResult("phase-record", len(phases), phaseErr)
	if phaseErr != nil {
		phases = nil
	}
	latest := followUpAuthoringPhases(in, phases)
	candidates := map[string]bool{}
	for _, folder := range in.Skipped {
		candidates[filepath.Base(filepath.Clean(folder.Dir))] = true
	}
	for slug := range latest {
		candidates[slug] = true
	}
	slugs := make([]string, 0, len(candidates))
	for slug := range candidates {
		if seen[slug] || !spec.IsArchiveRecordName(slug+".md") {
			continue
		}
		if _, err := os.Lstat(filepath.Join(in.SpecRoot, slug, "_prd.md")); !errors.Is(err, os.ErrNotExist) {
			continue
		}
		if _, err := os.Lstat(spec.ArchiveRecordPath(in.ArchiveRoot, slug)); !errors.Is(err, os.ErrNotExist) {
			continue
		}
		if _, recorded := latest[slug]; !recorded {
			info, err := os.Stat(filepath.Join(in.SpecRoot, slug))
			if err != nil || !info.IsDir() {
				continue
			}
		}
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)
	zone := in.Zone
	if zone == nil {
		zone = time.UTC
	}
	for _, slug := range slugs {
		status, note := "authoring", ""
		if phase, found := latest[slug]; found {
			status = "authoring: " + string(phase.Phase)
			note = string(phase.Phase) + " since " + phase.Time.In(zone).Format("01-02 15:04")
		}
		addSpec(slug, "authoring", status, note, false, "")
	}
	backlog, backlogErr := readFollowUpQueueBacklog(in.Loaded.GitRoot)
	if backlogErr == nil {
		rows = append(rows, backlog...)
	}
	for i := range rows {
		rows[i].Row = i + 1
	}
	return rows, forecasts, []followUpSource{followUpSourceResult("backlog", len(backlog), backlogErr), phaseSource}
}

func followUpSpecNumber(slug string) string {
	number, _, _ := strings.Cut(slug, "-")
	return number
}

func followUpQueueWait(item store.DeliveryQueueItem) string {
	switch {
	case strings.HasPrefix(item.Blocker, delivery.WaitRunSlot):
		return "after a Run slot frees"
	case strings.HasPrefix(item.Blocker, delivery.BlockerPrerequisiteUnmerged):
		_, detail, _ := strings.Cut(followUpFirstLine(item.Blocker), ":")
		numbers := []string{}
		for _, slug := range strings.Split(detail, ",") {
			if slug = strings.TrimSpace(slug); slug != "" {
				numbers = append(numbers, followUpSpecNumber(slug))
			}
		}
		if len(numbers) == 1 {
			return "after prerequisite " + numbers[0] + " is merged"
		}
		if len(numbers) > 1 {
			return "after prerequisites " + strings.Join(numbers, ", ") + " are merged"
		}
	}
	if item.Stage == store.DeliveryStageParked {
		return "after its park is answered"
	}
	return ""
}

func followUpQueueText(value string) string {
	value = strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ").Replace(value)
	// Repository text reaches a terminal: every control character (C0, DEL,
	// C1, so ESC and CSI sequences lose their introducer) and every line or
	// paragraph separator becomes a space before the text is cut.
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.In(r, unicode.Zl, unicode.Zp) {
			return ' '
		}
		return r
	}, value)
	runes := []rune(value)
	if len(runes) > 120 {
		return string(runes[:119]) + "…"
	}
	return value
}

func readFollowUpQueueBacklog(root string) ([]followUpQueueRow, error) {
	dir := filepath.Join(root, "docs", "backlog")
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read backlog directory: %w", err)
	}
	type datedRow struct {
		row     followUpQueueRow
		created string
	}
	dated := []datedRow{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		fields, err := readFollowUpMarkdown(path)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(fields.Status) != "open" {
			continue
		}
		note := fields.title
		if waiting, since, _ := followUpWaiting(fields); waiting != "" {
			note += "; waits on " + waiting + " since " + since
		}
		dated = append(dated, datedRow{row: followUpQueueRow{Kind: "backlog", Ref: deliverPlanRelativePath(root, path), Title: followUpQueueText(fields.title), Status: "not yet a Spec", Note: followUpQueueText(note), Start: followUpQueueCell{State: "none"}, End: followUpQueueCell{State: "none"}}, created: fields.Created})
	}
	sort.Slice(dated, func(i, j int) bool {
		if dated[i].created != dated[j].created {
			return dated[i].created < dated[j].created
		}
		return dated[i].row.Ref < dated[j].row.Ref
	})
	rows := make([]followUpQueueRow, 0, len(dated))
	for _, item := range dated {
		rows = append(rows, item.row)
	}
	return rows, nil
}

func followUpCanonicalDirectory(path string) (string, bool) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", false
	}
	info, err := os.Stat(canonical)
	return filepath.Clean(canonical), err == nil && info.IsDir()
}

func followUpAuthoringPhases(in followUpQueueInput, phases []phaserecord.RecordedPhase) map[string]phaserecord.RecordedPhase {
	latest := map[string]phaserecord.RecordedPhase{}
	roots, rootsErr := roundconfig.RepositoryRoots(in.Loaded.GitRoot)
	if rootsErr != nil {
		// The checkout still counts when linked-worktree metadata cannot be read.
		roots = nil
	}
	roots = append(roots, in.Loaded.GitRoot)
	accepted := map[string]bool{}
	for _, root := range roots {
		if path, ok := followUpCanonicalDirectory(root); ok {
			accepted[path] = true
		}
	}
	runner := in.Git
	if runner == nil {
		runner = preflight.ExecGitRunner{}
	}
	ctx := in.Context
	if ctx == nil {
		ctx = context.TODO()
	}
	var raw string
	var remoteErr error
	if in.Origin != nil {
		raw, remoteErr = in.Origin.remote, in.Origin.err
	} else {
		raw, remoteErr = runner.RunGit(ctx, in.Loaded.GitRoot, "remote", "get-url", "origin")
	}
	host, repository, remoteOK := parseReadinessRemote(raw)
	remoteOK = remoteOK && remoteErr == nil
	ordered := append([]phaserecord.RecordedPhase(nil), phases...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Time.After(ordered[j].Time) })
	checked := 0
	for _, entry := range ordered {
		if entry.Time.After(in.Now) || entry.Time.Before(in.Now.AddDate(0, 0, -7)) {
			continue
		}
		switch entry.Phase {
		case "idea", "prd", "techspec", "tasks":
		default:
			continue
		}
		path, exists := followUpCanonicalDirectory(entry.Repository)
		if !exists {
			continue
		}
		matches, known := accepted[path]
		if !known {
			matches = false
			if remoteOK && checked < 8 {
				checked++
				remote, err := runner.RunGit(ctx, path, "remote", "get-url", "origin")
				otherHost, otherRepo, ok := parseReadinessRemote(remote)
				matches = err == nil && ok && host == otherHost && repository == otherRepo
			}
			accepted[path] = matches
		}
		if matches {
			if _, seen := latest[entry.Spec]; !seen {
				latest[entry.Spec] = entry
			}
		}
	}
	return latest
}
