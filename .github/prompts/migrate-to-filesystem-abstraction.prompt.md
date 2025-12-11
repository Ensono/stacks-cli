## Plan: Implement FileSystem Interface Pattern for Testability Using BillyFS

Refactor filesystem operations in stacks-cli to use [BillyFS](https://pkg.go.dev/github.com/go-git/go-billy/v5), a battle-tested filesystem abstraction already available as a transitive dependency. This enables mock injection in tests via [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem) interface and [`memfs`](https://pkg.go.dev/github.com/go-git/go-billy/v5/memfs) for in-memory testing.

---

### Progress Checklist

> **Instructions for LLM**: As you complete each step, update this checklist by changing `[ ]` to `[x]`. Add notes about any issues encountered or deviations from the plan in the "Notes" column. Run tests after each step and record the result.

#### Phase 1: Foundation

| Step | Task                          | Status | Test Result | Notes |
| ---- | ----------------------------- | ------ | ----------- | ----- |
| 1    | Add direct BillyFS dependency | [ ]    |             |       |

#### Phase 2: Core Utilities

| Step | Task                              | Status | Test Result | Notes |
| ---- | --------------------------------- | ------ | ----------- | ----- |
| 2    | Refactor `internal/util/fs.go`    | [ ]    |             |       |
| 3    | Update `internal/util/fs_test.go` | [ ]    |             |       |

#### Phase 3: Configuration Layer

| Step | Task                               | Status | Test Result | Notes |
| ---- | ---------------------------------- | ------ | ----------- | ----- |
| 4    | Refactor `pkg/config/config.go`    | [ ]    |             |       |
| 5    | Update `pkg/config/config_test.go` | [ ]    |             |       |

#### Phase 4: Scaffolding Engine

| Step | Task                                | Status | Test Result | Notes |
| ---- | ----------------------------------- | ------ | ----------- | ----- |
| 6    | Refactor `pkg/scaffold/scaffold.go` | [ ]    |             |       |
| 7    | Refactor `pkg/scaffold/project.go`  | [ ]    |             |       |

#### Phase 5: Downloaders

| Step | Task                        | Status | Test Result | Notes |
| ---- | --------------------------- | ------ | ----------- | ----- |
| 8    | Refactor `pkg/downloaders/` | [ ]    |             |       |

#### Phase 6: Remaining Packages

| Step | Task                                      | Status | Test Result | Notes |
| ---- | ----------------------------------------- | ------ | ----------- | ----- |
| 9    | Refactor `pkg/export/export.go`           | [ ]    |             |       |
| 10   | Refactor `pkg/interactive/interactive.go` | [ ]    |             |       |
| 11   | Refactor `pkg/setup/setup.go`             | [ ]    |             |       |

#### Phase 7: Final Validation

| Step | Task                                     | Status | Test Result | Notes |
| ---- | ---------------------------------------- | ------ | ----------- | ----- |
| 12   | Run full test suite (`eirctl build`)     | [ ]    |             |       |
| 13   | Run integration tests (`eirctl inttest`) | [ ]    |             |       |

---

### Why BillyFS?

- **Already a dependency** - transitive via go-git, no new dependencies required
- **Battle-tested** - widely used in the Go ecosystem (go-git, etc.)
- **Built-in test support** - provides [`memfs`](https://pkg.go.dev/github.com/go-git/go-billy/v5/memfs) package for in-memory filesystem testing
- **Clean interface** - [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem) covers all common operations we need

### Patterns used in eirctl

#### Defined Interface

```go
type osFSOpsIface interface {
	Rename(oldpath string, newpath string) error
	WriteFile(name string, data []byte, perm os.FileMode) error
	Create(name string) (io.Writer, error)
}
```

#### Concrete Implementation

```go
type osFsOps struct { }
```

#### Use in Tests

```go
type mockOsFsOps struct {
	rename func(oldpath string, newpath string) error
	write  func(name string, data []byte, perm os.FileMode) error
}
```

### Incremental Refactoring Strategy

Refactor incrementally, running tests after each discrete area/file to maintain a tight feedback loop. This approach:

- Catches regressions early before they compound
- Keeps PRs reviewable and focused
- Allows course-correction if BillyFS gaps are discovered

### Steps

#### Phase 1: Foundation

1. **Add direct BillyFS dependency** - run `go get github.com/go-git/go-billy/v5` to make it a direct dependency with [`osfs`](https://pkg.go.dev/github.com/go-git/go-billy/v5/osfs) and [`memfs`](https://pkg.go.dev/github.com/go-git/go-billy/v5/memfs) packages available.
   - Run: `eirctl test:unit` ✓

#### Phase 2: Core Utilities (highest impact)

2. **Refactor `internal/util/fs.go`** - convert utility functions to accept [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem) as a parameter. This file contains ~18 filesystem operations and is used throughout the codebase.

   - Run: `eirctl test:unit` ✓

3. **Update `internal/util/fs_test.go`** - migrate tests to use [`memfs.New()`](https://pkg.go.dev/github.com/go-git/go-billy/v5/memfs#New) instead of `t.TempDir()`.
   - Run: `eirctl test:unit` ✓

#### Phase 3: Configuration Layer

4. **Refactor `pkg/config/config.go`** - add [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem) field to `Config` struct, update `Save()`, `WriteVariableFile()`, `WriteCommand()` methods.

   - Run: `eirctl test:unit` ✓

5. **Update `pkg/config/config_test.go`** - migrate to [`memfs.New()`](https://pkg.go.dev/github.com/go-git/go-billy/v5/memfs#New).
   - Run: `eirctl test:unit` ✓

#### Phase 4: Scaffolding Engine

6. **Refactor `pkg/scaffold/scaffold.go`** - inject [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem) into `Scaffold` struct.

   - Run: `eirctl test:unit` ✓

7. **Refactor `pkg/scaffold/project.go`** - update pipeline file operations to use billy.
   - Run: `eirctl test:unit` ✓

#### Phase 5: Downloaders

8. **Refactor `pkg/downloaders/`** - update Git, NuGet, and Filesystem downloaders to use [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem).
   - Run: `eirctl test:unit` ✓

#### Phase 6: Remaining Packages

9. **Refactor `pkg/export/export.go`** - update static file export operations.

   - Run: `eirctl test:unit` ✓

10. **Refactor `pkg/interactive/interactive.go`** - update config saving.

    - Run: `eirctl test:unit` ✓

11. **Refactor `pkg/setup/setup.go`** - update config file creation.
    - Run: `eirctl test:unit` ✓

#### Phase 7: Final Validation

12. **Run full test suite** - `eirctl build` to run all unit tests, linting, and compilation.
13. **Run integration tests** - `eirctl inttest` to validate end-to-end scaffolding workflows.

### Handle BillyFS Limitations

BillyFS doesn't support all `os` operations directly (e.g., `os.Chmod`, `os.RemoveAll`). For these, either:

- Use `billy.Capable` interface checks for optional capabilities
- Create a thin wrapper interface that extends `billy.Filesystem` with additional methods
- Use `util.RemoveAll()` helper that works with billy's `Remove()` recursively

### BillyFS Usage Patterns

```go
// Production code
import (
    "github.com/go-git/go-billy/v5"
    "github.com/go-git/go-billy/v5/osfs"
)

type Config struct {
    Fs billy.Filesystem
    // ... other fields
}

// Initialize with real filesystem
cfg := &Config{
    Fs: osfs.New("/"),
}

// Use throughout code
data, err := util.ReadFile(cfg.Fs, "/path/to/file")
```

```go
// Test code
import (
    "github.com/go-git/go-billy/v5/memfs"
)

func TestSomething(t *testing.T) {
    fs := memfs.New()
    // Pre-populate test files
    util.WriteFile(fs, "/test/file.txt", []byte("content"), 0644)

    cfg := &Config{Fs: fs}
    // Run test with in-memory filesystem
}
```

### BillyFS Capability Gap Audit

Audit of 77+ filesystem operations in the codebase revealed the following gaps:

#### Operations Requiring Helper Functions (Easy)

| Operation          | Count            | Solution                                                                   |
| ------------------ | ---------------- | -------------------------------------------------------------------------- |
| `os.ReadFile`      | 14 prod, 20 test | Implement `ReadFile(fs, path)` using `Open()` + `io.ReadAll()`             |
| `os.WriteFile`     | 14 prod, 20 test | Implement `WriteFile(fs, path, data, perm)` using `OpenFile()` + `Write()` |
| `os.RemoveAll`     | 2 prod, 9 test   | Implement `RemoveAll(fs, path)` using `ReadDir()` + `Remove()` recursively |
| `filepath.WalkDir` | 1 (file.go)      | Implement `WalkDir(fs, root, fn)` using `ReadDir()` recursively            |

#### Operations Requiring Interface Extension (Wrapper)

| Operation   | Count | Location          | Solution                                            |
| ----------- | ----- | ----------------- | --------------------------------------------------- |
| `os.Chmod`  | 2     | file.go, fetch.go | Create `ChmodCapable` interface; use type assertion |
| `os.Lchown` | 1     | file.go           | Create `ChownCapable` interface; skip on memfs      |

#### Operations NOT Requiring Migration (Process/Environment)

These are not filesystem operations - they read system state, not files:

| Operation                 | Count | Reason to Keep                                 |
| ------------------------- | ----- | ---------------------------------------------- |
| `os.Getwd` / `os.Chdir`   | 5     | Process working directory, not file I/O        |
| `os.Getenv` / `os.Setenv` | 8     | Environment variables                          |
| `os.UserHomeDir`          | 1     | System path resolution                         |
| `os.TempDir`              | 2     | System path resolution                         |
| `zip.OpenReader`          | 1     | Reads from OS path (temp file during download) |

### Recommended Pattern: Helper Functions in `internal/util/`

Create helper functions that work with [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem) in `internal/util/fs.go`:

```go
package util

import (
    "io"
    "io/fs"
    "os"
    "path/filepath"

    "github.com/go-git/go-billy/v5"
)

// ReadFile reads the named file from the filesystem and returns its contents.
func ReadFile(bfs billy.Filesystem, name string) ([]byte, error) {
    f, err := bfs.Open(name)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return io.ReadAll(f)
}

// WriteFile writes data to the named file, creating it if necessary.
func WriteFile(bfs billy.Filesystem, name string, data []byte, perm os.FileMode) error {
    f, err := bfs.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
    if err != nil {
        return err
    }
    _, err = f.Write(data)
    if err1 := f.Close(); err1 != nil && err == nil {
        err = err1
    }
    return err
}

// RemoveAll removes path and any children it contains.
func RemoveAll(bfs billy.Filesystem, path string) error {
    info, err := bfs.Lstat(path)
    if os.IsNotExist(err) {
        return nil
    }
    if err != nil {
        return err
    }
    if !info.IsDir() {
        return bfs.Remove(path)
    }
    entries, err := bfs.ReadDir(path)
    if err != nil {
        return err
    }
    for _, entry := range entries {
        err = RemoveAll(bfs, bfs.Join(path, entry.Name()))
        if err != nil {
            return err
        }
    }
    return bfs.Remove(path)
}

// WalkDir walks the file tree rooted at root, calling fn for each file or directory.
func WalkDir(bfs billy.Filesystem, root string, fn fs.WalkDirFunc) error {
    info, err := bfs.Lstat(root)
    if err != nil {
        return fn(root, nil, err)
    }
    return walkDir(bfs, root, &billyDirEntry{info}, fn)
}

func walkDir(bfs billy.Filesystem, path string, d fs.DirEntry, fn fs.WalkDirFunc) error {
    if err := fn(path, d, nil); err != nil {
        if err == fs.SkipDir && d.IsDir() {
            return nil
        }
        return err
    }
    if !d.IsDir() {
        return nil
    }
    entries, err := bfs.ReadDir(path)
    if err != nil {
        return fn(path, d, err)
    }
    for _, entry := range entries {
        childPath := bfs.Join(path, entry.Name())
        err = walkDir(bfs, childPath, &billyDirEntry{entry}, fn)
        if err != nil && err != fs.SkipDir {
            return err
        }
    }
    return nil
}

// billyDirEntry adapts os.FileInfo to fs.DirEntry
type billyDirEntry struct {
    info os.FileInfo
}

func (e *billyDirEntry) Name() string               { return e.info.Name() }
func (e *billyDirEntry) IsDir() bool                { return e.info.IsDir() }
func (e *billyDirEntry) Type() fs.FileMode          { return e.info.Mode().Type() }
func (e *billyDirEntry) Info() (fs.FileInfo, error) { return e.info, nil }
```

### Recommended Pattern: Optional Capabilities via Type Assertion

For operations like `Chmod` that aren't universally supported:

```go
// ChmodCapable is an optional interface for filesystems that support chmod.
type ChmodCapable interface {
    Chmod(name string, mode os.FileMode) error
}

// Chmod changes the mode of the named file if the filesystem supports it.
// Returns nil without error if the filesystem doesn't support chmod.
func Chmod(bfs billy.Filesystem, name string, mode os.FileMode) error {
    if cc, ok := bfs.(ChmodCapable); ok {
        return cc.Chmod(name, mode)
    }
    // Silently succeed on filesystems that don't support chmod (e.g., memfs)
    return nil
}
```

This pattern:

- **Gracefully degrades** - tests using `memfs` won't fail on chmod calls
- **Preserves functionality** - production code using `osfs` still sets permissions
- **Follows Go idioms** - uses interface type assertion pattern

### Summary: What Gets Migrated vs What Stays

| Category                                           | Action                      |
| -------------------------------------------------- | --------------------------- |
| File read/write (`ReadFile`, `WriteFile`)          | Migrate to helper functions |
| Directory ops (`MkdirAll`, `ReadDir`, `Stat`)      | Use billy methods directly  |
| Recursive ops (`RemoveAll`, `WalkDir`)             | Migrate to helper functions |
| Permissions (`Chmod`, `Lchown`)                    | Optional capability pattern |
| Process context (`Getwd`, `Getenv`, `UserHomeDir`) | Keep using `os` package     |
| Temp files (`TempDir`, zip extraction)             | Keep using `os` package     |
