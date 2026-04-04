package storage

import "testing"

func TestFilenameToVersionSuccess(t *testing.T) {
	filenames := []string{"001-start.sql", "123-drop-all.sql", "1000-so_many_migrations.sql", "2-.sql"}
	expected := []int{1, 123, 1000, 2}
	for i, v := range filenames {
		ver, err := filenameToVersion(v)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		} else if expected[i] != ver {
			t.Errorf("Expected %d, got %d", expected[i], ver)
		}
	}
}

func TestFilenameToVersionFailure(t *testing.T) {
	filenames := []string{"initialize.sql", "index.html", "one-two-three.sql"}
	for _, v := range filenames {
		_, err := filenameToVersion(v)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	}
}
