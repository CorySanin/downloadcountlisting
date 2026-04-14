package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store := New(dbPath)
	defer os.Remove(dbPath)

	path := "/test/path"
	filename := "file.txt"

	store.IncrementDownload(Download{
		DownloadIndex: DownloadIndex{
			Path:     path,
			Filename: filename,
		},
		AccessDomain: "localhost",
		UserAgent:    "curl",
	})

	ch := make(chan map[string]Totals, 1)
	store.GetTotalsByPath(path, ch)
	totals := <-ch

	if totals == nil {
		t.Fatal("Expected totals, got nil")
	}

	if totals[filename].All != 1 {
		t.Errorf("Expected 1 total download, got %d", totals[filename].All)
	}

	if totals[filename].Recent != 1 {
		t.Errorf("Expected 1 recent download, got %d", totals[filename].Recent)
	}

	err := store.RemoveDownloads([]DownloadIndex{{Path: path, Filename: filename}})
	if err != nil {
		t.Fatalf("Failed to remove downloads: %v", err)
	}

	store.GetTotalsByPath(path, ch)
	totals = <-ch
	if _, ok := totals[filename]; ok {
		t.Errorf("Expected file to be removed from totals")
	}

	for range 5 {
		store.IncrementDownload(Download{
			DownloadIndex: DownloadIndex{Path: path, Filename: filename},
		})
	}
	store.GetTotalsByPath(path, ch)
	totals = <-ch
	if totals[filename].All != 5 {
		t.Errorf("Expected 5 total downloads, got %d", totals[filename].All)
	}

	if err := store.Optimize(); err != nil {
		t.Errorf("Optimize failed: %v", err)
	}
}

func TestRecentDownloads(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test-recent.db")
	store := New(dbPath)
	defer os.Remove(dbPath)

	path := "/test/recent"
	filename := "old_file.txt"

	_, err := store.db.Exec("INSERT INTO downloads (Path, Filename, Timestamp) VALUES (?, ?, ?)",
		path, filename, time.Now().Add(-100*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan map[string]Totals, 1)
	store.GetTotalsByPath(path, ch)
	totals := <-ch

	if totals[filename].All != 1 {
		t.Errorf("Expected 1 total download, got %d", totals[filename].All)
	}
	if totals[filename].Recent != 0 {
		t.Errorf("Expected 0 recent downloads, got %d", totals[filename].Recent)
	}
}
