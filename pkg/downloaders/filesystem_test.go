package downloaders

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFilesystemDownloader(t *testing.T) {
	// Test data
	testPath := "/path/to/source"
	testTempDir := "/path/to/temp"

	// Create new filesystem downloader
	downloader := NewFilesystemDownloader(testPath, testTempDir)

	// Assertions
	assert.NotNil(t, downloader, "Downloader should not be nil")
	assert.Equal(t, testPath, downloader.Path, "Path should be set correctly")
	assert.Equal(t, testTempDir, downloader.TempDir, "TempDir should be set correctly")
	assert.Nil(t, downloader.logger, "Logger should be nil initially")
}

func TestFilesystem_SetLogger(t *testing.T) {
	downloader := NewFilesystemDownloader("test", "temp")
	logger := logrus.New()

	// Set logger
	downloader.SetLogger(logger)

	// Assertions
	assert.Equal(t, logger, downloader.logger, "Logger should be set correctly")
}

func TestFilesystem_PackageURL(t *testing.T) {
	testPath := "/path/to/source"
	downloader := NewFilesystemDownloader(testPath, "temp")

	// Get package URL
	url := downloader.PackageURL()

	// Assertions
	assert.Equal(t, testPath, url, "PackageURL should return the path")
}

func TestFilesystem_Get_Success(t *testing.T) {
	// Create temporary directories for testing
	sourceDir := t.TempDir()
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")

	// Create test files in source directory
	testFiles := []string{
		"file1.txt",
		"subdir/file2.txt",
		"subdir/nested/file3.txt",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(sourceDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err, "Should create directories")

		err = os.WriteFile(fullPath, []byte("test content for "+file), 0644)
		require.NoError(t, err, "Should create test file")
	}

	// Create downloader
	downloader := NewFilesystemDownloader(sourceDir, destDir)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress log output during tests
	downloader.SetLogger(logger)

	// Execute Get method
	result, err := downloader.Get()

	// Assertions
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, destDir, result, "Get should return the destination directory")

	// Verify files were copied
	for _, file := range testFiles {
		destPath := filepath.Join(destDir, file)
		assert.FileExists(t, destPath, "File should exist in destination: %s", file)

		// Verify content
		content, err := os.ReadFile(destPath)
		assert.NoError(t, err, "Should be able to read copied file")
		assert.Equal(t, "test content for "+file, string(content), "File content should be preserved")
	}
}

func TestFilesystem_Get_SkipsGitDirectory(t *testing.T) {
	// Create temporary directories for testing
	sourceDir := t.TempDir()
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")

	// Create test files including .git directory
	testStructure := map[string]string{
		"file1.txt":           "regular file",
		"subdir/file2.txt":    "nested file",
		".git/config":         "git config",
		".git/HEAD":           "ref: refs/heads/main",
		"subdir/.git/config":  "nested git config",
	}

	for file, content := range testStructure {
		fullPath := filepath.Join(sourceDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err, "Should create directories")

		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err, "Should create test file")
	}

	// Create downloader
	downloader := NewFilesystemDownloader(sourceDir, destDir)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress log output during tests
	downloader.SetLogger(logger)

	// Execute Get method
	result, err := downloader.Get()

	// Assertions
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, destDir, result, "Get should return the destination directory")

	// Verify regular files were copied
	assert.FileExists(t, filepath.Join(destDir, "file1.txt"), "Regular file should be copied")
	assert.FileExists(t, filepath.Join(destDir, "subdir/file2.txt"), "Nested regular file should be copied")

	// Verify .git directories were skipped
	assert.NoFileExists(t, filepath.Join(destDir, ".git/config"), ".git directory should be skipped")
	assert.NoFileExists(t, filepath.Join(destDir, ".git/HEAD"), ".git files should be skipped")
	assert.NoFileExists(t, filepath.Join(destDir, "subdir/.git/config"), "Nested .git directory should be skipped")
}

func TestFilesystem_Get_NonExistentSource(t *testing.T) {
	// Use non-existent source directory
	sourceDir := "/path/that/does/not/exist"
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")

	// Create downloader
	downloader := NewFilesystemDownloader(sourceDir, destDir)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress log output during tests
	downloader.SetLogger(logger)

	// Execute Get method
	result, err := downloader.Get()

	// Assertions
	assert.Error(t, err, "Get should return an error for non-existent source")
	assert.Equal(t, destDir, result, "Get should still return the destination directory")
}

func TestFilesystem_Get_EmptySourceDirectory(t *testing.T) {
	// Create empty source directory
	sourceDir := t.TempDir()
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")

	// Create downloader
	downloader := NewFilesystemDownloader(sourceDir, destDir)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress log output during tests
	downloader.SetLogger(logger)

	// Execute Get method
	result, err := downloader.Get()

	// Assertions
	assert.NoError(t, err, "Get should not return an error for empty source directory")
	assert.Equal(t, destDir, result, "Get should return the destination directory")
	assert.DirExists(t, destDir, "Destination directory should be created")

	// Verify destination is empty (except for potential hidden files)
	entries, err := os.ReadDir(destDir)
	assert.NoError(t, err, "Should be able to read destination directory")
	assert.Empty(t, entries, "Destination directory should be empty")
}

func TestFilesystem_Get_PreservesFilePermissions(t *testing.T) {
	// Skip on Windows as it doesn't have the same permission model
	if os.Getenv("GOOS") == "windows" || os.PathSeparator == '\\' {
		t.Skip("Skipping permission test on Windows")
	}

	// Create temporary directories for testing
	sourceDir := t.TempDir()
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")

	// Create test file with specific permissions
	testFile := filepath.Join(sourceDir, "executable.sh")
	err := os.WriteFile(testFile, []byte("#!/bin/bash\necho hello"), 0755)
	require.NoError(t, err, "Should create executable file")

	// Create downloader
	downloader := NewFilesystemDownloader(sourceDir, destDir)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress log output during tests
	downloader.SetLogger(logger)

	// Execute Get method
	result, err := downloader.Get()

	// Assertions
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, destDir, result, "Get should return the destination directory")

	// Verify file permissions are preserved
	destFile := filepath.Join(destDir, "executable.sh")
	assert.FileExists(t, destFile, "File should exist in destination")

	sourceInfo, err := os.Stat(testFile)
	require.NoError(t, err, "Should get source file info")

	destInfo, err := os.Stat(destFile)
	require.NoError(t, err, "Should get destination file info")

	assert.Equal(t, sourceInfo.Mode(), destInfo.Mode(), "File permissions should be preserved")
}

func TestFilesystem_Get_WithRelativePaths(t *testing.T) {
	// Save current directory and restore at end
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Should get current directory")
	defer os.Chdir(originalDir)

	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Change to temp directory and use relative paths
	err = os.Chdir(tempDir)
	require.NoError(t, err, "Should change to temp directory")

	// Create relative source directory
	relativeSource := "relative_source"
	err = os.MkdirAll(relativeSource, 0755)
	require.NoError(t, err, "Should create relative source directory")

	// Create test file in relative source
	testFile := filepath.Join(relativeSource, "test.txt")
	err = os.WriteFile(testFile, []byte("relative test content"), 0644)
	require.NoError(t, err, "Should create test file in relative source")

	// Create downloader with relative paths
	destDir := "relative_dest"
	downloader := NewFilesystemDownloader(relativeSource, destDir)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress log output during tests
	downloader.SetLogger(logger)

	// Execute Get method
	result, err := downloader.Get()

	// Assertions
	assert.NoError(t, err, "Get should not return an error with relative paths")
	assert.Equal(t, destDir, result, "Get should return the relative destination directory")

	// Verify file was copied
	destFile := filepath.Join(destDir, "test.txt")
	assert.FileExists(t, destFile, "File should exist in relative destination")

	// Verify content
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err, "Should be able to read copied file")
	assert.Equal(t, "relative test content", string(content), "File content should be preserved")
}

// Test that the downloader implements the Downloader interface
func TestFilesystem_ImplementsDownloaderInterface(t *testing.T) {
	downloader := NewFilesystemDownloader("test", "temp")
	
	// This will fail compilation if the interface is not implemented correctly
	var _ interface {
		Get() (string, error)
		PackageURL() string
		SetLogger(*logrus.Logger)
	} = downloader
	
	// If we get here, the interface is implemented correctly
	assert.True(t, true, "Filesystem downloader implements the required interface")
}

// Benchmark test for performance measurement
func BenchmarkFilesystem_Get(b *testing.B) {
	// Create temporary directories
	sourceDir := b.TempDir()
	tempBaseDir := b.TempDir()

	// Create test files
	for i := 0; i < 100; i++ {
		testFile := filepath.Join(sourceDir, fmt.Sprintf("file%d.txt", i))
		os.WriteFile(testFile, []byte("benchmark test content"), 0644)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress log output during benchmarks

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		destDir := filepath.Join(tempBaseDir, fmt.Sprintf("bench%d", i))
		downloader := NewFilesystemDownloader(sourceDir, destDir)
		downloader.SetLogger(logger)
		
		_, err := downloader.Get()
		if err != nil {
			b.Fatal(err)
		}
	}
}
