package showmetadatafetcher

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// makeTree creates the given relative paths under root. Entries ending in "/"
// are created as directories, everything else as empty files (parent
// directories are created as needed).
func makeTree(t *testing.T, root string, paths ...string) {
	t.Helper()
	for _, p := range paths {
		full := filepath.Join(root, p)
		if p[len(p)-1] == '/' {
			if err := os.MkdirAll(full, 0o755); err != nil {
				t.Fatalf("mkdir %s: %v", full, err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
		}
		if err := os.WriteFile(full, nil, 0o644); err != nil {
			t.Fatalf("write %s: %v", full, err)
		}
	}
}

func TestGetQueriesForDirectory(t *testing.T) {
	root := t.TempDir()
	makeTree(t, root,
		"Death-note/Show.txt",
		"The-Quintessential-Quintuplets/season1/Show.txt",
		"The-Quintessential-Quintuplets/season2/Show.txt",
	)

	// os.ReadDir returns entries sorted by name, so the order is deterministic.
	want := []queryPathPair{
		{
			path:  filepath.Join(root, "Death-note"),
			query: "Death-note",
		},
		{
			path:  filepath.Join(root, "The-Quintessential-Quintuplets", "season1"),
			query: "The-Quintessential-Quintupletsseason1",
		},
		{
			path:  filepath.Join(root, "The-Quintessential-Quintuplets", "season2"),
			query: "The-Quintessential-Quintupletsseason2",
		},
	}

	got, err := GetQueriesForDirectory(root)
	if err != nil {
		t.Fatalf("GetQueriesForDirectory(%q) returned error: %v", root, err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("GetQueriesForDirectory(%q)\n Expected: %+v\n Got: %+v", root, want, got)
	}
}

func TestGetQueriesForDirectory_ShowWithSeasonsHasNoShowLevelPair(t *testing.T) {
	root := t.TempDir()
	makeTree(t, root,
		"TQQ/season1/ep1.mkv",
		"TQQ/season2/ep1.mkv",
	)

	got, err := GetQueriesForDirectory(root)
	if err != nil {
		t.Fatalf("GetQueriesForDirectory(%q) returned error: %v", root, err)
	}

	for _, pair := range got {
		if pair.path == filepath.Join(root, "TQQ") {
			t.Errorf("got a show-level pair %+v for a show that is split into seasons", pair)
		}
	}
	if len(got) != 2 {
		t.Errorf("expected 2 season pairs, got %d: %+v", len(got), got)
	}
}

func TestGetQueriesForDirectory_FlatShowIgnoresFileNames(t *testing.T) {
	root := t.TempDir()
	makeTree(t, root,
		"TQQ/ep1.mkv",
		"TQQ/ep2.mkv",
		"TQQ/ep3.mkv",
	)

	want := []queryPathPair{
		{path: filepath.Join(root, "TQQ"), query: "TQQ"},
	}

	got, err := GetQueriesForDirectory(root)
	if err != nil {
		t.Fatalf("GetQueriesForDirectory(%q) returned error: %v", root, err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("GetQueriesForDirectory(%q)\n Expected: %+v\n Got: %+v", root, want, got)
	}
}

func TestGetQueriesForDirectory_EmptyRoot(t *testing.T) {
	root := t.TempDir()

	got, err := GetQueriesForDirectory(root)
	if err != nil {
		t.Fatalf("GetQueriesForDirectory(%q) returned error: %v", root, err)
	}
	if len(got) != 0 {
		t.Errorf("expected no pairs for an empty directory, got %+v", got)
	}
}

func TestGetQueriesForDirectory_MissingRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "Yotsuba-fan-IQ")

	got, err := GetQueriesForDirectory(root)
	if err == nil {
		t.Fatalf("expected an error for a missing directory, got nil (pairs: %+v)", got)
	}
	if len(got) != 0 {
		t.Errorf("expected no pairs on error, got %+v", got)
	}
}
