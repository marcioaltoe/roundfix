package cli

import (
	"context"
	"io"
	"os"
	"sync"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

// runUI is the per-command Run Event consumer wiring for the owning
// resolve/watch process. In a TTY it opens the interactive cockpit on its
// own read-only journal connection (ADR 0009: the cockpit never consumes
// the live sink) and silences raw progress text; otherwise the journal
// stays critical alongside the plain-text writer and progress goes to
// stderr, byte-compatible with the previous behavior.
type runUI struct {
	sink     runevent.Sink
	progress io.Writer

	fanout        *runevent.Fanout
	reader        *store.Store
	cockpitCancel context.CancelFunc
	cockpitDone   chan error

	waitOnce  sync.Once
	closeOnce sync.Once
}

func startRunUI(ctx context.Context, view roundtui.LiveRunView, runID string, homeDir string, runStore *store.Store, stderr io.Writer) (*runUI, error) {
	journal := store.JournalSink{Store: runStore}
	if !liveTUIEnabled(stderr) {
		fanout := runevent.NewFanout([]runevent.Sink{journal, agent.WriterSink{Writer: stderr}}, nil)
		return &runUI{sink: fanout, progress: stderr, fanout: fanout}, nil
	}

	reader, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		return nil, err
	}
	fanout := runevent.NewFanout([]runevent.Sink{journal}, nil)
	// The cockpit outlives command-context cancellation on purpose: after a
	// Stop Request the user keeps the view open to inspect what happened,
	// closing it with q. Close cancels it as the backstop.
	cockpitCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	ui := &runUI{
		sink:          fanout,
		progress:      io.Discard,
		fanout:        fanout,
		reader:        reader,
		cockpitCancel: cancel,
		cockpitDone:   make(chan error, 1),
	}
	go func() {
		ui.cockpitDone <- roundtui.RunCockpit(cockpitCtx, stderr, roundtui.CockpitConfig{
			Mode:   roundtui.CockpitOwning,
			View:   view,
			RunID:  runID,
			Source: reader,
			OnStop: interruptSelf,
		})
	}()
	return ui, nil
}

// Wait keeps the cockpit on screen until the user closes it (q after the
// Run reaches a terminal state). Non-TTY mode returns immediately.
func (ui *runUI) Wait() {
	if ui.cockpitDone == nil {
		return
	}
	ui.waitOnce.Do(func() {
		<-ui.cockpitDone
	})
}

// Close drains the fanout, quits the cockpit, and releases the reader. It
// is idempotent so commands can both defer it and call it before printing
// their closing summary.
func (ui *runUI) Close() {
	ui.closeOnce.Do(func() {
		if ui.fanout != nil {
			ui.fanout.Close()
		}
		if ui.cockpitCancel != nil {
			ui.cockpitCancel()
			ui.waitOnce.Do(func() {
				<-ui.cockpitDone
			})
		}
		if ui.reader != nil {
			_ = ui.reader.Close()
		}
	})
}

// interruptSelf turns the owning cockpit's Ctrl-C key into the same SIGINT
// the terminal would deliver outside raw mode, so Stop Request semantics
// and exit codes stay exactly as today.
func interruptSelf() {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		return
	}
	_ = process.Signal(os.Interrupt)
}
