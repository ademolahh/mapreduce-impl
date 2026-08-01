package sequential

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/ademolahh/map-reduce-impl/internal/shared"
)

func TestRunMapPhase(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"a.txt": "hello world",
		"b.txt": "hello go",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write fixture file %q: %v", name, err)
		}
	}

	fakeMap := func(filename, content string) []shared.KV {
		var kvs []shared.KV
		for w := range strings.FieldsSeq(content) {
			kvs = append(kvs, shared.KV{Word: w, Value: "1"})
		}
		return kvs
	}

	got, err := runMapPhase(fakeMap, dir)
	if err != nil {
		t.Fatalf("runMapPhase returned error: %v", err)
	}

	if len(got) != 4 {
		t.Errorf("expected 4 KV pairs, got %d (%+v)", len(got), got)
	}

	count := 0
	for _, kv := range got {
		if kv.Word == "hello" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected \"hello\" to appear 2 times, got %d", count)
	}
}

func TestRunReducePhase(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "result.txt")

	k := shared.KeyValue{
		{Word: "go", Value: "1"},
		{Word: "hello", Value: "1"},
		{Word: "hello", Value: "1"},
		{Word: "world", Value: "1"},
	}

	fakeReduce := func(word string, values []string) string {
		return strconv.Itoa(len(values))
	}

	if err := runReducePhase(k, fakeReduce, outPath); err != nil {
		t.Fatalf("runReducePhase returned error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	got := string(data)
	want := "go 1\nhello 2\nworld 1\n"
	if got != want {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", got, want)
	}
}
