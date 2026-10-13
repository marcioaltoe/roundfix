//go:build darwin

package suiteguard

import (
	"encoding/binary"
	"os/exec"
	"testing"
)

func TestDarwinResidueArgumentsDoNotCountAsEnvironment(t *testing.T) {
	t.Parallel()
	const mark = "123.0123456789abcdef"
	raw := make([]byte, 4)
	binary.NativeEndian.PutUint32(raw, 2)
	raw = append(raw, []byte("/fixture\x00\x00/fixture\x00"+testBinaryMarkEnv+"="+mark+"\x00OTHER=value\x00")...)
	environment, command, err := parseDarwinProcessArguments(raw, 123)
	if err != nil {
		t.Fatal(err)
	}
	if environmentHasMark(environment, mark) {
		t.Fatalf("argument mistaken for environment: %q", command)
	}
	raw = append(raw, []byte(testBinaryMarkEnv+"="+mark+"\x00")...)
	environment, _, err = parseDarwinProcessArguments(raw, 123)
	if err != nil || !environmentHasMark(environment, mark) {
		t.Fatalf("real environment mark not found: %q, %v", environment, err)
	}
}

func TestDarwinResidueProvesAReapedPIDAbsent(t *testing.T) {
	t.Parallel()
	command := exec.Command("/usr/bin/true")
	if err := command.Run(); err != nil {
		t.Fatal(err)
	}
	absent, err := darwinResidueProcessAbsent(command.Process.Pid)
	if err != nil || !absent {
		t.Fatalf("reaped pid %d absent = %v, error %v", command.Process.Pid, absent, err)
	}
}
