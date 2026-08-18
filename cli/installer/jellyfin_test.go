package installer

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestTailReturnsLastLines(t *testing.T) {
	result := tail("one\ntwo\nthree\nfour", 2)
	if result != "three\nfour" {
		t.Fatalf("expected last two lines, got %q", result)
	}
}

func TestTailReturnsEverythingWhenShorterThanLimit(t *testing.T) {
	result := tail("one\ntwo", 40)
	if result != "one\ntwo" {
		t.Fatalf("expected all lines, got %q", result)
	}
}

func TestTailIgnoresTrailingNewlines(t *testing.T) {
	result := tail("one\ntwo\n\n", 1)
	if result != "two" {
		t.Fatalf("expected last non empty line, got %q", result)
	}
}

func TestTailEmpty(t *testing.T) {
	result := tail("", 10)
	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

func TestNewestFilePicksMostRecentlyModified(t *testing.T) {
	dir := t.TempDir()
	older := path.Join(dir, "log_20260101.log")
	newer := path.Join(dir, "log_20260102.log")
	writeFile(t, older, "old")
	writeFile(t, newer, "new")

	past := time.Now().Add(-time.Hour)
	if err := os.Chtimes(older, past, past); err != nil {
		t.Fatal(err)
	}

	result, err := newestFile(dir)
	if err != nil {
		t.Fatal(err)
	}
	if result != newer {
		t.Fatalf("expected %s, got %s", newer, result)
	}
}

func TestNewestFileSkipsDirectories(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(path.Join(dir, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	file := path.Join(dir, "log.log")
	writeFile(t, file, "content")

	result, err := newestFile(dir)
	if err != nil {
		t.Fatal(err)
	}
	if result != file {
		t.Fatalf("expected %s, got %s", file, result)
	}
}

func TestNewestFileEmptyDir(t *testing.T) {
	_, err := newestFile(t.TempDir())
	if err == nil {
		t.Fatal("expected an error for an empty directory")
	}
}

func TestServerLogIncludesFatalError(t *testing.T) {
	dataDir := t.TempDir()
	logDir := path.Join(dataDir, "data", "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, path.Join(logDir, "log_20260818.log"),
		"[18:41:10] [INF] Main: starting\n"+
			"[18:41:10] [FTL] Main: Unhandled Exception\n"+
			"System.InvalidOperationException: The path has insufficient free space. Available: 855MiB, Required: 2GiB.\n")

	result := jellyfinForLog(dataDir).serverLog()

	if !strings.Contains(result, fatalMarker) {
		t.Fatalf("expected the fatal marker in the log, got %q", result)
	}
	if !strings.Contains(result, "insufficient free space") {
		t.Fatalf("expected the underlying reason in the log, got %q", result)
	}
}

func TestServerLogWhenDirectoryMissing(t *testing.T) {
	result := jellyfinForLog(t.TempDir()).serverLog()

	if !strings.Contains(result, "no server log in") {
		t.Fatalf("expected a readable message when the log dir is missing, got %q", result)
	}
}

func jellyfinForLog(dataDir string) *Jellyfin {
	return NewJellyfin("", dataDir, NewExecutor(zap.NewNop()), zap.NewNop())
}

func writeFile(t *testing.T, name, content string) {
	t.Helper()
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
