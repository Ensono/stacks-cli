package util

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
)

func TestReadWriteFile(t *testing.T) {
	bfs := memfs.New()
	path := filepath.Join("nested", "file.txt")
	content := []byte("hello")

	if err := WriteFile(bfs, path, content, 0o644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	data, err := ReadFile(bfs, path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	if string(data) != string(content) {
		t.Fatalf("unexpected content: %s", data)
	}
}

func TestRemoveAll(t *testing.T) {
	bfs := memfs.New()
	if err := bfs.MkdirAll("a/b/c", os.ModePerm); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	if _, err := bfs.Create(filepath.Join("a", "b", "c", "file.txt")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := RemoveAll(bfs, "a"); err != nil {
		t.Fatalf("RemoveAll error: %v", err)
	}

	if _, err := bfs.Stat("a"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected directory to be removed, got %v", err)
	}

	if err := RemoveAll(bfs, "doesnotexist"); err != nil {
		t.Fatalf("RemoveAll missing path: %v", err)
	}
}

func TestWalkDir(t *testing.T) {
	bfs := memfs.New()
	if err := bfs.MkdirAll("root/sub", os.ModePerm); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	if _, err := bfs.Create("root/sub/file.txt"); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	var visited []string
	err := WalkDir(bfs, "root", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		visited = append(visited, path)
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir error: %v", err)
	}

	expected := map[string]bool{
		"root":                                   true,
		filepath.Join("root", "sub"):             true,
		filepath.Join("root", "sub", "file.txt"): true,
	}

	for _, entry := range visited {
		delete(expected, entry)
	}

	if len(expected) != 0 {
		t.Fatalf("WalkDir did not visit: %v", expected)
	}
}

func TestChmodNoop(t *testing.T) {
	bfs := memfs.New()
	if err := Chmod(bfs, "doesnotexist", 0o644); err != nil {
		t.Fatalf("Chmod unexpected error: %v", err)
	}
}
