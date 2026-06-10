package cli

import (
	"io"
	"os"
	"strings"
)

type terminalStyle struct {
	color bool
}

func styleFor(writer io.Writer) terminalStyle {
	return terminalStyle{color: colorEnabled(writer)}
}

func colorEnabled(writer io.Writer) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("ROUNDFIX_COLOR"))) {
	case "always", "1", "true", "yes", "on":
		return true
	case "never", "0", "false", "no", "off":
		return false
	}
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func (style terminalStyle) red(text string) string {
	return style.wrap("31", text)
}

func (style terminalStyle) cyan(text string) string {
	return style.wrap("36", text)
}

func (style terminalStyle) yellow(text string) string {
	return style.wrap("33", text)
}

func (style terminalStyle) green(text string) string {
	return style.wrap("32", text)
}

func (style terminalStyle) bold(text string) string {
	return style.wrap("1", text)
}

func (style terminalStyle) dim(text string) string {
	return style.wrap("2", text)
}

func (style terminalStyle) wrap(code string, text string) string {
	if !style.color || text == "" {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}
