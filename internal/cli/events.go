package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

type eventsRequest struct {
	runID  string
	follow bool
	filter runevent.StreamCategoryFilter
}

func runEventsCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("events"))
		return exitOK
	}
	req, err := parseEventsCommand(args)
	if err != nil {
		printEventsFailure(err, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printEventsFailure(err, stderr)
		return exitPreflight
	}
	reader, err := store.OpenReader(ctx, loaded.HomeDir)
	if err != nil {
		printEventsFailure(err, stderr)
		return exitPreflight
	}
	defer func() {
		_ = reader.Close()
	}()

	run, found, err := reader.Run(ctx, req.runID)
	if err != nil {
		printEventsFailure(err, stderr)
		return exitPreflight
	}
	if !found {
		printEventsFailure(validationError{message: fmt.Sprintf("Run %q does not exist", req.runID)}, stderr)
		return exitPreflight
	}

	encoder := json.NewEncoder(stdout)
	cursor, err := replayEventStream(ctx, reader, run.ID, 0, req.filter, encoder)
	if err != nil {
		printEventsFailure(err, stderr)
		return exitRunFailed
	}
	if !req.follow || store.IsTerminalState(run.State) {
		return exitOK
	}

	follower := journalFollower{
		source: reader,
		sleep:  attachSleep,
		accept: func(entry store.JournalEvent) error {
			return encodeStreamEntry(entry, req.filter, encoder)
		},
	}
	if _, _, err := follower.follow(ctx, run.ID, cursor); err != nil {
		if isStopRequest(ctx, err) {
			return exitSIGINT
		}
		printEventsFailure(err, stderr)
		return exitRunFailed
	}
	return exitOK
}

func parseEventsCommand(args []string) (eventsRequest, error) {
	req := eventsRequest{filter: runevent.AllStreamCategories()}
	positionals := []string{}
	filterSet := false
	filterRaw := ""
	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--follow" || arg == "-follow":
			req.follow = true
		case arg == "--filter" || arg == "-filter":
			if index+1 >= len(args) {
				return req, validationError{message: "missing value for --filter"}
			}
			index++
			filterSet = true
			filterRaw = args[index]
		case strings.HasPrefix(arg, "--filter="):
			filterSet = true
			filterRaw = strings.TrimPrefix(arg, "--filter=")
		case strings.HasPrefix(arg, "-filter="):
			filterSet = true
			filterRaw = strings.TrimPrefix(arg, "-filter=")
		case strings.HasPrefix(arg, "-"):
			return req, validationError{message: fmt.Sprintf("unknown flag %q", arg)}
		default:
			positionals = append(positionals, arg)
		}
	}
	if filterSet {
		filter, err := runevent.ParseStreamCategoryFilter(filterRaw)
		if err != nil {
			return req, validationError{message: err.Error()}
		}
		req.filter = filter
	}
	if len(positionals) == 0 {
		return req, validationError{message: "missing run id"}
	}
	if len(positionals) > 1 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", positionals[1])}
	}
	req.runID = strings.TrimSpace(positionals[0])
	if req.runID == "" {
		return req, validationError{message: "missing run id"}
	}
	return req, nil
}

func replayEventStream(ctx context.Context, source journalEventSource, runID string, cursor int64, filter runevent.StreamCategoryFilter, encoder *json.Encoder) (int64, error) {
	return replayJournalEvents(ctx, source, runID, cursor, func(entry store.JournalEvent) error {
		return encodeStreamEntry(entry, filter, encoder)
	})
}

func encodeStreamEntry(entry store.JournalEvent, filter runevent.StreamCategoryFilter, encoder *json.Encoder) error {
	record, ok, err := runevent.ProjectStreamEvent(entry.Cursor, entry.Event, filter)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if err := encoder.Encode(record); err != nil {
		return fmt.Errorf("write Run Event Stream record: %w", err)
	}
	return nil
}

func printEventsFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "roundfix events failed: %v\n", err)
}
