//go:build unix

package store

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const processStartLayout = "Mon Jan _2 15:04:05 2006"

func readOwnedProcess(ctx context.Context, pid int) (OwnedProcess, error) {
	cmd := exec.CommandContext(
		ctx,
		"ps",
		"-p", strconv.Itoa(pid),
		"-o", "lstart=",
		"-o", "time=",
		"-o", "comm=",
	)
	environment := make([]string, 0, len(os.Environ())+1)
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(variable, "LC_ALL=") {
			environment = append(environment, variable)
		}
	}
	cmd.Env = append(environment, "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return OwnedProcess{}, fmt.Errorf("read process table entry: %w", err)
	}
	return parsePSOwnedProcess(output)
}

func parsePSOwnedProcess(output []byte) (OwnedProcess, error) {
	line := strings.TrimSpace(string(output))
	if len(line) < len(processStartLayout) {
		return OwnedProcess{}, fmt.Errorf("parse process table entry %q: missing start time", line)
	}
	started, err := time.ParseInLocation(processStartLayout, line[:len(processStartLayout)], time.Local)
	if err != nil {
		return OwnedProcess{}, fmt.Errorf("parse process start time: %w", err)
	}
	remainder := strings.TrimSpace(line[len(processStartLayout):])
	cpuField, command, ok := strings.Cut(remainder, " ")
	if !ok {
		return OwnedProcess{}, fmt.Errorf("parse process table entry %q: missing command", line)
	}
	cpuTime, err := parsePSCPUTime(cpuField)
	if err != nil {
		return OwnedProcess{}, err
	}
	return OwnedProcess{
		Started: started,
		CPUTime: cpuTime,
		Command: strings.TrimSpace(command),
	}, nil
}

func parsePSCPUTime(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("parse process CPU time: value is empty")
	}
	var duration time.Duration
	if dayField, clock, ok := strings.Cut(value, "-"); ok {
		dayCount, err := strconv.ParseUint(dayField, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse process CPU time days %q: %w", dayField, err)
		}
		duration, err = processDurationComponent(dayCount, 24*time.Hour)
		if err != nil {
			return 0, fmt.Errorf("parse process CPU time days %q: %w", dayField, err)
		}
		value = clock
	}
	fields := strings.Split(value, ":")
	if len(fields) < 2 || len(fields) > 3 {
		return 0, fmt.Errorf("parse process CPU time %q: expected [[hours:]minutes:]seconds", value)
	}
	var hoursField string
	var minutesField string
	var secondsField string
	if len(fields) == 3 {
		hoursField, minutesField, secondsField = fields[0], fields[1], fields[2]
	} else {
		minutesField, secondsField = fields[0], fields[1]
	}
	hours, err := strconv.ParseUint(hoursFieldOrZero(hoursField), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse process CPU time hours %q: %w", hoursField, err)
	}
	minutes, err := strconv.ParseUint(minutesField, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse process CPU time minutes %q: %w", minutesField, err)
	}
	seconds, err := time.ParseDuration(secondsField + "s")
	if err != nil {
		return 0, fmt.Errorf("parse process CPU time seconds %q: %w", secondsField, err)
	}
	hourDuration, err := processDurationComponent(hours, time.Hour)
	if err != nil {
		return 0, fmt.Errorf("parse process CPU time hours %q: %w", hoursField, err)
	}
	minuteDuration, err := processDurationComponent(minutes, time.Minute)
	if err != nil {
		return 0, fmt.Errorf("parse process CPU time minutes %q: %w", minutesField, err)
	}
	for _, component := range []time.Duration{hourDuration, minuteDuration, seconds} {
		if component > 0 && duration > time.Duration(math.MaxInt64)-component {
			return 0, fmt.Errorf("parse process CPU time %q: duration overflows", value)
		}
		duration += component
	}
	return duration, nil
}

func hoursFieldOrZero(value string) string {
	if value == "" {
		return "0"
	}
	return value
}

func processDurationComponent(value uint64, unit time.Duration) (time.Duration, error) {
	if value > uint64(math.MaxInt64/int64(unit)) {
		return 0, fmt.Errorf("duration component overflows")
	}
	return time.Duration(value) * unit, nil
}
