package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
)

const (
	detachHandshakeFDEnv   = "ROUNDFIX_DETACH_HANDSHAKE_FD"
	detachConsoleTempEnv   = "ROUNDFIX_DETACH_CONSOLE_TEMP"
	detachHandshakeFD      = 3
	detachHandshakeTimeout = 10 * time.Second
	detachWaitTimeout      = 2 * time.Second
)

type detachChild struct {
	handshake *os.File
	tempPath  string
}

type detachHandshake struct {
	runID      string
	consoleLog string
	err        error
}

func applyDetachSemantics(req commandRequest) commandRequest {
	if req.detach || req.detachChild != nil {
		req.noInput = true
	}
	return req
}

func (req commandRequest) reportDetachedRunCreated(runID string) error {
	if req.detachChild == nil {
		return nil
	}
	return req.detachChild.reportRunCreated(runID, req.artifactDir)
}

func newDetachChildFromEnv() (*detachChild, error) {
	fdValue := strings.TrimSpace(os.Getenv(detachHandshakeFDEnv))
	tempPath := strings.TrimSpace(os.Getenv(detachConsoleTempEnv))
	if fdValue == "" && tempPath == "" {
		return nil, nil
	}
	if fdValue == "" || tempPath == "" {
		return nil, errors.New("detached child is missing handshake environment")
	}
	fd, err := strconv.Atoi(fdValue)
	if err != nil || fd < 0 {
		return nil, fmt.Errorf("detached child handshake fd %q is invalid", fdValue)
	}
	file := os.NewFile(uintptr(fd), "roundfix-detach-handshake")
	if file == nil {
		return nil, fmt.Errorf("open detached child handshake fd %d", fd)
	}
	_ = os.Unsetenv(detachHandshakeFDEnv)
	_ = os.Unsetenv(detachConsoleTempEnv)
	return &detachChild{handshake: file, tempPath: tempPath}, nil
}

func (child *detachChild) Close() error {
	if child == nil || child.handshake == nil {
		return nil
	}
	err := child.handshake.Close()
	child.handshake = nil
	return err
}

func (child *detachChild) reportRunCreated(runID string, artifactDir string) error {
	if child == nil || child.handshake == nil {
		return nil
	}
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return errors.New("detached Run id is empty")
	}
	consoleLog := detachedConsoleLogPath(artifactDir, runID)
	if err := os.MkdirAll(filepath.Dir(consoleLog), 0o755); err != nil {
		return fmt.Errorf("create Detached Run console log directory: %w", err)
	}
	if err := os.Rename(child.tempPath, consoleLog); err != nil {
		return fmt.Errorf("move Detached Run console log into place: %w", err)
	}
	if _, err := fmt.Fprintf(child.handshake, "%s\t%s\n", runID, consoleLog); err != nil {
		return fmt.Errorf("write Detached Run handshake: %w", err)
	}
	return child.Close()
}

func runDetachedCommand(args []string, req commandRequest, loaded roundconfig.Loaded, stdout, stderr io.Writer) int {
	artifactDir, err := roundconfig.ValidateArtifactDirectory(req.artifactDir, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	tempFile, err := createDetachedConsoleTemp(artifactDir)
	if err != nil {
		printPreflightFailure(req.name, err, stderr)
		return exitPreflight
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = tempFile.Close()
	}()

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		printPreflightFailure(req.name, fmt.Errorf("create Detached Run handshake pipe: %w", err), stderr)
		_ = os.Remove(tempPath)
		return exitPreflight
	}
	defer func() {
		_ = readPipe.Close()
	}()

	executable, err := os.Executable()
	if err != nil {
		printPreflightFailure(req.name, fmt.Errorf("resolve current executable: %w", err), stderr)
		_ = writePipe.Close()
		_ = os.Remove(tempPath)
		return exitPreflight
	}
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		printPreflightFailure(req.name, fmt.Errorf("open %s for Detached Run stdin: %w", os.DevNull, err), stderr)
		_ = writePipe.Close()
		_ = os.Remove(tempPath)
		return exitPreflight
	}
	defer func() {
		_ = devNull.Close()
	}()

	cmd := exec.Command(executable, removeDetachArg(args)...)
	cmd.Env = detachedChildEnv(os.Environ(), detachHandshakeFD, tempPath)
	cmd.Stdin = devNull
	cmd.Stdout = tempFile
	cmd.Stderr = tempFile
	cmd.ExtraFiles = []*os.File{writePipe}
	configureDetachedSysProcAttr(cmd)

	if err := cmd.Start(); err != nil {
		printPreflightFailure(req.name, fmt.Errorf("start Detached Run child: %w", err), stderr)
		_ = writePipe.Close()
		_ = os.Remove(tempPath)
		return exitPreflight
	}
	_ = writePipe.Close()
	_ = tempFile.Close()

	handshakeCh := make(chan detachHandshake, 1)
	go func() {
		handshakeCh <- readDetachedHandshake(readPipe)
	}()

	timer := time.NewTimer(detachHandshakeTimeout)
	defer timer.Stop()
	select {
	case handshake := <-handshakeCh:
		if handshake.err != nil {
			waitErr, waited := waitDetachedProcess(cmd, detachWaitTimeout)
			if !waited {
				killDetachedProcess(cmd)
				waitErr = errors.New("Detached Run child closed the handshake without exiting")
			}
			relayDetachedConsole(stderr, tempPath)
			_ = os.Remove(tempPath)
			return exitCodeFromWait(waitErr)
		}
		if err := cmd.Process.Release(); err != nil {
			printPreflightFailure(req.name, fmt.Errorf("release Detached Run child: %w", err), stderr)
			return exitPreflight
		}
		printDetachedReport(stdout, handshake.runID, handshake.consoleLog)
		return exitOK
	case <-timer.C:
		killDetachedProcess(cmd)
		waitErr, _ := waitDetachedProcess(cmd, detachWaitTimeout)
		relayDetachedConsole(stderr, tempPath)
		_ = os.Remove(tempPath)
		if waitErr != nil {
			return exitCodeFromWait(waitErr)
		}
		return exitRunFailed
	}
}

func createDetachedConsoleTemp(artifactDir string) (*os.File, error) {
	dir := filepath.Join(artifactDir, "runs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create Detached Run console directory: %w", err)
	}
	file, err := os.CreateTemp(dir, ".roundfix-detach-*.log")
	if err != nil {
		return nil, fmt.Errorf("create Detached Run console log: %w", err)
	}
	return file, nil
}

func detachedConsoleLogPath(artifactDir string, runID string) string {
	return filepath.Join(artifactDir, "runs", runID, "console.log")
}

func readDetachedHandshake(reader io.Reader) detachHandshake {
	line, err := bufio.NewReader(reader).ReadString('\n')
	if err != nil {
		return detachHandshake{err: err}
	}
	runID, consoleLog, ok := strings.Cut(strings.TrimSuffix(line, "\n"), "\t")
	if !ok || strings.TrimSpace(runID) == "" || strings.TrimSpace(consoleLog) == "" {
		return detachHandshake{err: fmt.Errorf("invalid Detached Run handshake %q", line)}
	}
	return detachHandshake{runID: strings.TrimSpace(runID), consoleLog: strings.TrimSpace(consoleLog)}
}

func printDetachedReport(stdout io.Writer, runID string, consoleLog string) {
	fmt.Fprintf(stdout, "Run detached: %s\n", runID)
	fmt.Fprintf(stdout, "Console log: %s\n", consoleLog)
	fmt.Fprintf(stdout, "Follow: roundfix attach %s\n", runID)
	fmt.Fprintf(stdout, "Stop: roundfix stop %s\n", runID)
}

func relayDetachedConsole(stderr io.Writer, path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(stderr, "read Detached Run console log %q: %v\n", path, err)
		}
		return
	}
	_, _ = stderr.Write(content)
}

func waitDetachedProcess(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err := <-waitCh:
		return err, true
	case <-timer.C:
		return nil, false
	}
}

func killDetachedProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	if killDetachedProcessGroup(cmd) {
		return
	}
	_ = cmd.Process.Kill()
}

func exitCodeFromWait(err error) int {
	if err == nil {
		return exitOK
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()
		if code >= 0 {
			return code
		}
	}
	return exitRunFailed
}

func removeDetachArg(args []string) []string {
	result := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--detach" || arg == "-detach" ||
			strings.HasPrefix(arg, "--detach=") || strings.HasPrefix(arg, "-detach=") {
			continue
		}
		result = append(result, arg)
	}
	return result
}

func detachedChildEnv(env []string, fd int, tempPath string) []string {
	clean := make([]string, 0, len(env)+2)
	for _, entry := range env {
		key, _, _ := strings.Cut(entry, "=")
		if key == detachHandshakeFDEnv || key == detachConsoleTempEnv {
			continue
		}
		clean = append(clean, entry)
	}
	clean = append(clean, fmt.Sprintf("%s=%d", detachHandshakeFDEnv, fd))
	clean = append(clean, detachConsoleTempEnv+"="+tempPath)
	return clean
}
