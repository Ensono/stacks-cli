package util

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
)

// ChmodCapable describes filesystems that expose a Chmod method.
type ChmodCapable interface {
	Chmod(name string, mode os.FileMode) error
}

// ReadFile reads the named file from the provided billy Filesystem.
func ReadFile(bfs billy.Filesystem, name string) ([]byte, error) {
	file, err := bfs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// WriteFile writes data to the named file, creating parent directories if required.
func WriteFile(bfs billy.Filesystem, name string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(name)
	if dir != "." && dir != "" {
		if err := bfs.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := bfs.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	_, writeErr := file.Write(data)
	closeErr := file.Close()

	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

// RemoveAll removes path and any children using the billy Filesystem.
func RemoveAll(bfs billy.Filesystem, name string) error {
	info, err := bfs.Lstat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return bfs.Remove(name)
	}

	entries, err := bfs.ReadDir(name)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		childPath := bfs.Join(name, entry.Name())
		if err := RemoveAll(bfs, childPath); err != nil {
			return err
		}
	}

	return bfs.Remove(name)
}

// WalkDir traverses the filesystem rooted at root, calling fn for each file or directory.
func WalkDir(bfs billy.Filesystem, root string, fn fs.WalkDirFunc) error {
	info, err := bfs.Lstat(root)
	if err != nil {
		return fn(root, nil, err)
	}

	return walkDir(bfs, root, &billyDirEntry{info: info}, fn)
}

func walkDir(bfs billy.Filesystem, path string, entry fs.DirEntry, fn fs.WalkDirFunc) error {
	if err := fn(path, entry, nil); err != nil {
		if err == fs.SkipDir && entry.IsDir() {
			return nil
		}
		return err
	}

	if !entry.IsDir() {
		return nil
	}

	entries, err := bfs.ReadDir(path)
	if err != nil {
		return fn(path, entry, err)
	}

	for _, info := range entries {
		childPath := bfs.Join(path, info.Name())
		if err := walkDir(bfs, childPath, &billyDirEntry{info: info}, fn); err != nil {
			if err == fs.SkipDir {
				continue
			}
			return err
		}
	}

	return nil
}

type billyDirEntry struct {
	info os.FileInfo
}

func (e *billyDirEntry) Name() string               { return e.info.Name() }
func (e *billyDirEntry) IsDir() bool                { return e.info.IsDir() }
func (e *billyDirEntry) Type() fs.FileMode          { return e.info.Mode().Type() }
func (e *billyDirEntry) Info() (fs.FileInfo, error) { return e.info, nil }

// Chmod changes the permissions of the named file if supported by the filesystem.
func Chmod(bfs billy.Filesystem, name string, mode os.FileMode) error {
	if fsCapable, ok := bfs.(ChmodCapable); ok {
		if err := fsCapable.Chmod(name, mode); err != nil {
			if strings.Contains(err.Error(), "does not implement billy.Chmod") {
				return nil
			}
			return err
		}
	}
	return nil
}
