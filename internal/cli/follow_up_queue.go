package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
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
	addSpec := func(slug, kind, status, note string, parked bool, wait string, item *store.DeliveryQueueItem) {
		row := followUpQueueRow{Spec: slug, Kind: kind, Status: status, Note: followUpQueueText(note)}
		if kind != "authoring" {
			var graph *spec.Graph
			var err error
			if item != nil {
				graph, err = followUpItemTaskGraph(in.Context, in.Git, in.Loaded.GitRoot, in.SpecRoot, in.ArchiveRoot, *item)
			} else {
				graph, err = spec.Load(in.SpecRoot, slug)
			}
			if err == nil {
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
		addSpec(item.SpecSlug, "delivery_queue", string(item.Stage), note, parked, wait, &item)
	}
	active := append([]spec.Spec(nil), in.Active...)
	sort.Slice(active, func(i, j int) bool { return active[i].Slug < active[j].Slug })
	for _, item := range active {
		if !seen[item.Slug] {
			addSpec(item.Slug, "planned", "not queued", "", false, "", nil)
		}
	}
	phases, _, phaseErr := phaserecord.Phases(in.Loaded.HomeDir, in.Now.AddDate(0, 0, -7))
	phaseSource := followUpSourceResult("phase-record", len(phases), phaseErr)
	if phaseErr != nil {
		phases = nil
	}
	latest := followUpAuthoringPhases(in, phases)
	held := followUpHeldSpecNumbers(in)
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
		if info, err := os.Stat(filepath.Join(in.SpecRoot, slug)); err != nil || !info.IsDir() {
			entry, recorded := latest[slug]
			if !recorded || followUpRenamedSlug(slug, entry, held, latest) {
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
		addSpec(slug, "authoring", status, note, false, "", nil)
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

// followUpItemTaskGraph reads committed Task progress without opening the item's
// worktree. An unreadable branch or graph falls back to the checkout's Spec Root.
func followUpItemTaskGraph(ctx context.Context, git preflight.GitRunner,
	gitRoot, specRoot, archiveRoot string, item store.DeliveryQueueItem) (*spec.Graph, error) {
	fallback := func() (*spec.Graph, error) { return spec.Load(specRoot, item.SpecSlug) }
	if item.Branch == "" {
		return fallback()
	}
	relativeRoot := func(root string) (string, bool) {
		if !filepath.IsAbs(root) {
			root = filepath.Join(gitRoot, root)
		}
		relative, err := filepath.Rel(gitRoot, root)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return "", false
		}
		return filepath.ToSlash(relative), true
	}
	specPath, specInside := relativeRoot(specRoot)
	archivePath, archiveInside := relativeRoot(archiveRoot)
	if !specInside || !archiveInside {
		return fallback()
	}
	if git == nil {
		git = preflight.ExecGitRunner{}
	}
	if ctx == nil {
		ctx = context.TODO()
	}
	revision, err := git.RunGit(ctx, gitRoot, "rev-parse", "--verify", "refs/heads/"+item.Branch+"^{commit}")
	if err != nil {
		return fallback()
	}
	revision = strings.TrimSpace(revision)
	readAt := func(revision string) func(string) ([]byte, error) {
		return func(file string) ([]byte, error) {
			content, err := git.RunGit(ctx, gitRoot, "show", revision+":"+filepath.ToSlash(file))
			if err != nil {
				// git show has no typed missing-path error. These reads are best
				// effort; an unavailable object is absent from this count source.
				return nil, fmt.Errorf("read Item Branch path %q: %w", file, errors.Join(os.ErrNotExist, err))
			}
			return []byte(content), nil
		}
	}
	read := readAt(revision)
	if graph, err := spec.LoadAt(specPath, item.SpecSlug, read); err == nil {
		return graph, nil
	}
	archived, err := spec.ReadArchivedSpecAt(archivePath, item.SpecSlug, read)
	if err == nil && archived.Form == spec.ArchivedRecord {
		record := archived.Record
		if graph, err := spec.LoadAt(path.Dir(record.Source), path.Base(record.Source), readAt(record.SourceRevision)); err == nil {
			return graph, nil
		}
	}
	return fallback()
}

func followUpHeldSpecNumbers(in followUpQueueInput) map[string]map[string]bool {
	held := map[string]map[string]bool{}
	hold := func(slug string) {
		if number, ok := specNumber(slug); ok {
			if held[number] == nil {
				held[number] = map[string]bool{}
			}
			held[number][slug] = true
		}
	}
	for _, item := range in.Active {
		hold(item.Slug)
	}
	for _, folder := range in.Skipped {
		hold(filepath.Base(filepath.Clean(folder.Dir)))
	}
	for _, item := range in.Progress.queue.Items {
		hold(item.SpecSlug)
	}
	if entries, err := os.ReadDir(in.ArchiveRoot); err == nil {
		for _, entry := range entries {
			if spec.IsArchiveRecordName(entry.Name()) {
				hold(strings.TrimSuffix(entry.Name(), ".md"))
			}
		}
	}
	return held
}

func followUpRenamedSlug(slug string, entry phaserecord.RecordedPhase,
	held map[string]map[string]bool, latest map[string]phaserecord.RecordedPhase) bool {
	number, ok := specNumber(slug)
	if !ok {
		return false
	}
	for other := range held[number] {
		if other != slug {
			return true
		}
	}
	for other, phase := range latest {
		if otherNumber, ok := specNumber(other); ok && otherNumber == number && other != slug {
			if phase.Time.After(entry.Time) || (phase.Time.Equal(entry.Time) && other > slug) {
				return true
			}
		}
	}
	return false
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
