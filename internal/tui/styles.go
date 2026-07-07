package tui

import (
	"strings"

	"roundfix/internal/runevent"

	"charm.land/lipgloss/v2"
)

// Tokens is the cockpit style contract (ui-redesign-plan §7): cyan for
// section labels and active borders, green for done, amber for running,
// waiting, and pending, red only for locked, failed, or blocking states,
// muted gray for timestamps, paths, and background text. This file is the
// single style source: exact palette indices live only here.
type Tokens struct {
	SectionLabel lipgloss.Style // cyan
	ActiveBorder lipgloss.Style // cyan panel border
	Done         lipgloss.Style // green
	Running      lipgloss.Style // amber
	Waiting      lipgloss.Style // amber
	Pending      lipgloss.Style // amber
	Blocked      lipgloss.Style // red
	Failed       lipgloss.Style // red
	Locked       lipgloss.Style // red
	Muted        lipgloss.Style // gray: timestamps, paths
	Selection    lipgloss.Style // accent card border
}

// ResolveTokens maps the effective color mode to styled or identity tokens.
// Callers resolve the mode through the existing ROUNDFIX_COLOR / NO_COLOR
// contract and pass the result here. Identity tokens keep border structure
// (borders are layout, not color) but render text unchanged, so every state
// distinction rides on the existing text markers.
func ResolveTokens(colorEnabled bool) Tokens {
	panelBorder := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	cardBorder := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	if !colorEnabled {
		return Tokens{ActiveBorder: panelBorder, Selection: cardBorder}
	}
	var (
		cyan  = lipgloss.Color("39")
		green = lipgloss.Color("78")
		amber = lipgloss.Color("214")
		red   = lipgloss.Color("203")
		gray  = lipgloss.Color("244")
	)
	amberText := lipgloss.NewStyle().Foreground(amber)
	redText := lipgloss.NewStyle().Foreground(red)
	return Tokens{
		SectionLabel: lipgloss.NewStyle().Foreground(cyan),
		ActiveBorder: panelBorder.BorderForeground(cyan),
		Done:         lipgloss.NewStyle().Foreground(green),
		Running:      amberText,
		Waiting:      amberText,
		Pending:      amberText,
		Blocked:      redText,
		Failed:       redText,
		Locked:       redText,
		Muted:        lipgloss.NewStyle().Foreground(gray),
		Selection:    cardBorder.BorderForeground(cyan),
	}
}

// EventSummary returns one bounded line for a Run Event: the first line of
// its summary, truncated to width. It never renders the raw payload.
func EventSummary(event runevent.RunEvent, width int) string {
	summary := event.Summary
	if index := strings.IndexByte(summary, '\n'); index >= 0 {
		summary = summary[:index]
	}
	summary = strings.TrimRight(summary, "\r")
	return truncateDisplay(summary, width)
}

// cockpitTokens is the styled token set the current renderers alias below.
// Per-surface token adoption and the color-mode wiring land with the
// fidelity tasks that restyle each pane.
var cockpitTokens = ResolveTokens(true)

// Legacy style names fold into the tokens where a token exists. The
// remaining styles are cockpit chrome pending per-surface restyling; they
// stay here so no other file defines styles.
var (
	styleAccent       = cockpitTokens.SectionLabel
	styleMuted        = cockpitTokens.Muted
	styleError        = cockpitTokens.Failed
	styleActiveBorder = cockpitTokens.ActiveBorder
	styleBright       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	styleTool         = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	styleBar          = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("238"))
	styleBarFill      = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("27")).Bold(true)
	styleBarRest      = lipgloss.NewStyle().Foreground(lipgloss.Color("17")).Background(lipgloss.Color("153"))
	styleFooter       = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Background(lipgloss.Color("234"))
	styleBorder       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("238")).Padding(0, 1)
)
