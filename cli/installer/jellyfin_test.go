package installer

import (
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap"
)

type RunnerStub struct {
	output  string
	err     error
	command string
}

func (r *RunnerStub) Run(app string, args ...string) (string, error) {
	r.command = strings.Join(append([]string{app}, args...), " ")
	return r.output, r.err
}

func TestServerLogIncludesFatalError(t *testing.T) {
	runner := &RunnerStub{output: "[FTL] Main: Unhandled Exception\n" +
		"System.InvalidOperationException: The path has insufficient free space. " +
		"Available: 855MiB, Required: 2GiB.\n"}

	result := jellyfinWith(runner).serverLog()

	if !strings.Contains(result, fatalMarker) {
		t.Fatalf("expected the fatal marker in the log, got %q", result)
	}
	if !strings.Contains(result, "insufficient free space") {
		t.Fatalf("expected the underlying reason in the log, got %q", result)
	}
}

func TestServerLogReadsTheServerUnit(t *testing.T) {
	runner := &RunnerStub{output: "some log"}

	jellyfinWith(runner).serverLog()

	if !strings.Contains(runner.command, serverUnit) {
		t.Fatalf("expected the server unit to be read, got %q", runner.command)
	}
	if !strings.Contains(runner.command, "journalctl") {
		t.Fatalf("expected the journal to be read, got %q", runner.command)
	}
}

func TestServerLogWhenJournalCannotBeRead(t *testing.T) {
	runner := &RunnerStub{output: "permission denied", err: fmt.Errorf("exit status 1")}

	result := jellyfinWith(runner).serverLog()

	if !strings.Contains(result, "cannot read the journal") {
		t.Fatalf("expected a readable message when the journal fails, got %q", result)
	}
	if !strings.Contains(result, "permission denied") {
		t.Fatalf("expected the command output to be kept, got %q", result)
	}
}

func jellyfinWith(runner CommandRunner) *Jellyfin {
	return NewJellyfin("", "/var/snap/jellyfin/current", runner, zap.NewNop())
}
