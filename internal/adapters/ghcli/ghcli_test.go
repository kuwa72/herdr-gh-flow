package ghcli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kuwa72/lead-cli/internal/ports"
)

// mkDummyGh installs an executable dummy `gh` that logs "$@" and serves
// canned JSON. It returns the temp bin dir (to prepend to PATH).
func mkDummyGh(t *testing.T, exitCode int, stderrMsg string) (binDir, logPath string) {
	t.Helper()
	binDir = t.TempDir()
	logPath = filepath.Join(binDir, "gh-args.log")
	script := "#!/bin/sh\n" +
		"printf '<%s>\\n' \"$@\" >> \"$GH_LOG\"\n" +
		"if [ -n \"" + stderrMsg + "\" ]; then echo \"" + stderrMsg + "\" >&2; fi\n" +
		"if [ \"" + itoa(exitCode) + "\" != \"0\" ]; then exit " + itoa(exitCode) + "; fi\n" +
		"if [ \"$1 $2\" = \"issue list\" ]; then\n" +
		"  printf '[{\"number\":50,\"title\":\"Tracking issue\"},{\"number\":36,\"title\":\"ports adapter\"}]'\n" +
		"elif [ \"$1 $2\" = \"issue view\" ]; then\n" +
		"  printf '{\"number\":36,\"title\":\"ports adapter\",\"body\":\"hello body\",\"state\":\"OPEN\"}'\n" +
		"else\n" +
		"  echo \"unexpected: $@\" >&2; exit 3\n" +
		"fi\n"
	if err := os.WriteFile(filepath.Join(binDir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GH_LOG", logPath)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return binDir, logPath
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func readLog(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading dummy log: %v", err)
	}
	return string(b)
}

func TestListOpen_CallsGhWithExpectedArgs(t *testing.T) {
	_, logPath := mkDummyGh(t, 0, "")

	got, err := New().ListOpen(context.Background())
	if err != nil {
		t.Fatalf("ListOpen: %v", err)
	}
	if len(got) != 2 || got[0].Number != 50 || got[0].Title != "Tracking issue" ||
		got[1].Number != 36 || got[1].Title != "ports adapter" {
		t.Fatalf("ListOpen = %+v, want issues 50 and 36 with titles", got)
	}

	log := readLog(t, logPath)
	for _, want := range []string{
		"<issue>", "<list>", "<--state>", "<open>",
		"<--limit>", "<50>", "<--json>", "<number,title>",
	} {
		if !strings.Contains(log, want) {
			t.Errorf("gh args log missing %q, got:\n%s", want, log)
		}
	}
}

func TestView_CallsGhViewWithJSONFields(t *testing.T) {
	_, logPath := mkDummyGh(t, 0, "")

	got, err := New().View(context.Background(), 36)
	if err != nil {
		t.Fatalf("View: %v", err)
	}
	if got.Number != 36 || got.Title != "ports adapter" || got.Body != "hello body" || got.State != "OPEN" {
		t.Fatalf("View = %+v, want number/title/body/state", got)
	}

	log := readLog(t, logPath)
	for _, want := range []string{
		"<issue>", "<view>", "<36>", "<--json>", "<number,title,body,state>",
	} {
		if !strings.Contains(log, want) {
			t.Errorf("gh args log missing %q, got:\n%s", want, log)
		}
	}
}

func TestListOpen_MissingGhReturnsTypedError(t *testing.T) {
	t.Setenv("PATH", t.TempDir()) // no gh on PATH

	_, err := New().ListOpen(context.Background())
	if err == nil {
		t.Fatal("ListOpen with no gh = nil error, want typed not-found error")
	}
	if !ports.IsBinaryNotFound(err) {
		t.Fatalf("ListOpen error = %v (%T), want BinaryNotFoundError for graceful fallback", err, err)
	}
}

func TestView_PropagatesGhFailure(t *testing.T) {
	mkDummyGh(t, 1, "issue not found")

	_, err := New().View(context.Background(), 999)
	if err == nil {
		t.Fatal("View with failing gh = nil error, want propagation")
	}
	if ports.IsBinaryNotFound(err) {
		t.Fatalf("View error = %v, must not be classified as missing binary", err)
	}
	if !strings.Contains(err.Error(), "issue not found") {
		t.Errorf("View error = %q, want gh stderr surfaced", err)
	}
}
